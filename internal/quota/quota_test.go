package quota

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

// QUOTA_TEST_DATABASE_URL points at a PostgreSQL database where the test
// can create its own schema. No application services or SMTP server are used.
func TestMonthlyAccounting(t *testing.T) {
	raw := os.Getenv("QUOTA_TEST_DATABASE_URL")
	if raw == "" {
		t.Skip("set QUOTA_TEST_DATABASE_URL to run the PostgreSQL integration test")
	}
	u, err := url.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	admin, err := sql.Open("postgres", u.String())
	if err != nil {
		t.Fatal(err)
	}
	defer admin.Close()
	schema := fmt.Sprintf("quota_test_%d", time.Now().UnixNano())
	if _, err := admin.Exec("CREATE SCHEMA " + schema); err != nil {
		t.Fatal(err)
	}
	defer admin.Exec("DROP SCHEMA " + schema + " CASCADE")
	q := u.Query()
	q.Set("search_path", schema)
	u.RawQuery = q.Encode()
	db, err := sql.Open("postgres", u.String())
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec(`CREATE TABLE settings(key TEXT PRIMARY KEY,value JSONB);
 INSERT INTO settings VALUES ('app.monthly_email_limit','2');
 CREATE TABLE email_monthly_usage(month DATE PRIMARY KEY,used BIGINT NOT NULL CHECK (used>=0));
 INSERT INTO email_monthly_usage VALUES (date_trunc('month',CURRENT_TIMESTAMP AT TIME ZONE 'UTC') - interval '1 month',99);
 CREATE TABLE campaigns(id int,type text,last_subscriber_id int);
 CREATE TABLE campaign_lists(campaign_id int,list_id int);
 CREATE TABLE lists(id int,optin text);
 CREATE TABLE subscribers(id int,status text);
 CREATE TABLE subscriber_lists(subscriber_id int,list_id int,status text);
 INSERT INTO campaigns VALUES(1,'regular',0);
 INSERT INTO lists VALUES(1,'single'),(2,'double');
 INSERT INTO campaign_lists VALUES(1,1),(1,2);
 INSERT INTO subscribers VALUES(1,'enabled'),(2,'enabled'),(3,'blocklisted'),(4,'enabled');
 INSERT INTO subscriber_lists VALUES(1,1,'confirmed'),(1,2,'confirmed'),(2,1,'unconfirmed'),(3,1,'confirmed'),(4,2,'unconfirmed');`)
	if err != nil {
		t.Fatal(err)
	}
	store := Store{db: db}
	v, err := store.Snapshot()
	if err != nil || v.Used != 0 || v.Limit != 2 {
		t.Fatalf("new month: %+v, %v", v, err)
	}
	checkCampaign := func(want bool) {
		t.Helper()
		blocked, err := store.BlocksCampaign(1, v)
		if err != nil || blocked != want {
			t.Fatalf("blocked=%v want %v, err=%v", blocked, want, err)
		}
	}
	checkCampaign(false)
	v.Used = 1
	checkCampaign(true)
	if _, err = db.Exec("UPDATE campaigns SET last_subscriber_id=1"); err != nil {
		t.Fatal(err)
	}
	checkCampaign(false)
	smtpErr := errors.New("SMTP failure")
	if err := store.Send(func() error { return smtpErr }); !errors.Is(err, smtpErr) {
		t.Fatal(err)
	}
	var submitted atomic.Int32
	var wg sync.WaitGroup
	for i := 0; i < 12; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := store.Send(func() error { submitted.Add(1); return nil })
			if err != nil && !errors.Is(err, ErrLimited) {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
	v, err = store.Snapshot()
	if err != nil || v.Used != 2 || submitted.Load() != 2 {
		t.Fatalf("concurrent sends: %+v, accepted %d, err=%v", v, submitted.Load(), err)
	}
	t.Setenv("MAX_QUOTA", "0")
	locked, err := New(db)
	if err != nil {
		t.Fatal(err)
	}
	v, err = locked.Snapshot()
	if err != nil || !v.Locked || v.Limit != 0 {
		t.Fatalf("override: %+v, %v", v, err)
	}
	if err := locked.Send(func() error { t.Error("sent with zero quota"); return nil }); !errors.Is(err, ErrLimited) {
		t.Fatal(err)
	}
}
