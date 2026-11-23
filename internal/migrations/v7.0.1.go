package migrations

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/knadh/koanf/v2"
	"github.com/knadh/stuffbin"
)

func V7_0_1(db *sqlx.DB, fs stuffbin.FileSystem, ko *koanf.Koanf, lo *log.Logger) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS email_monthly_usage (month DATE PRIMARY KEY, used BIGINT NOT NULL CHECK (used >= 0));
INSERT INTO settings (key,value) VALUES ('app.monthly_email_limit','-1') ON CONFLICT (key) DO NOTHING;`)
	return err
}
