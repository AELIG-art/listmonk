// Package quota accounts for SMTP messages across workers and restarts.
package quota

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

const Key = "app.monthly_email_limit"

var ErrLimited = errors.New("Monthly email quota exceeded")
var ErrAccounting = errors.New("monthly email quota accounting unavailable")

type Store struct {
	db       *sql.DB
	override *int64
}

type Snapshot struct {
	Limit  int64
	Used   int64
	Locked bool
	Reset  string
}

func New(db *sql.DB) (*Store, error) {
	s := &Store{db: db}
	if raw, ok := os.LookupEnv("MAX_QUOTA"); ok {
		n, err := strconv.ParseInt(raw, 10, 64)
		if err != nil || n < 0 {
			return nil, errors.New("MAX_QUOTA must be a nonnegative integer")
		}
		s.override = &n
	}
	return s, nil
}

func (s *Store) Override() *int64 { return s.override }

const monthSQL = "date_trunc('month', CURRENT_TIMESTAMP AT TIME ZONE 'UTC')::date"
const limitSQL = "COALESCE($1::bigint, (SELECT value::text::bigint FROM settings WHERE key='app.monthly_email_limit'))"

func (s *Store) Snapshot() (Snapshot, error) {
	v := Snapshot{Locked: s.override != nil}
	err := s.db.QueryRow("SELECT "+limitSQL+", COALESCE((SELECT used FROM email_monthly_usage WHERE month="+monthSQL+"),0), to_char("+monthSQL+" + interval '1 month','YYYY-MM-DD')", s.override).Scan(&v.Limit, &v.Used, &v.Reset)
	return v, err
}

// Send reserves one slot before SMTP, including during unlimited months.
// PostgreSQL serializes competing updates to the same monthly counter.
func (s *Store) Send(send func() error) error {
	var month time.Time
	err := s.db.QueryRow("WITH cap AS (SELECT "+limitSQL+" AS n) "+`
 INSERT INTO email_monthly_usage (month,used) SELECT `+monthSQL+`,1 FROM cap WHERE n <> 0
 ON CONFLICT (month) DO UPDATE SET used=email_monthly_usage.used+1
 WHERE (SELECT n FROM cap)<0 OR email_monthly_usage.used < (SELECT n FROM cap)
 RETURNING month`, s.override).Scan(&month)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrLimited
	}
	if err != nil {
		return fmt.Errorf("%w: %w", ErrAccounting, err)
	}
	if err := send(); err != nil {
		if _, rollback := s.db.Exec("UPDATE email_monthly_usage SET used=used-1 WHERE month=$1", month); rollback != nil {
			return fmt.Errorf("%w: SMTP: %v, refund: %w", ErrAccounting, err, rollback)
		}
		return err
	}
	return nil
}

// BlocksCampaign counts distinct eligible subscribers after the resume checkpoint.
func (s *Store) BlocksCampaign(id int, v Snapshot) (bool, error) {
	if v.Limit < 0 {
		return false, nil
	}
	var n int64
	err := s.db.QueryRow(`SELECT COUNT(DISTINCT sl.subscriber_id)
 FROM campaigns c JOIN campaign_lists cl ON cl.campaign_id=c.id
 JOIN lists l ON l.id=cl.list_id JOIN subscriber_lists sl ON sl.list_id=l.id
 JOIN subscribers sub ON sub.id=sl.subscriber_id
 WHERE c.id=$1 AND sl.subscriber_id>c.last_subscriber_id AND sub.status!='blocklisted'
 AND CASE WHEN c.type='optin' THEN sl.status='unconfirmed' AND l.optin='double'
 WHEN l.optin='double' THEN sl.status='confirmed' ELSE sl.status!='unsubscribed' END`, id).Scan(&n)
	return n > v.Limit-v.Used, err
}
