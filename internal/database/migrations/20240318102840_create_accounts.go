package migrations

import (
	"database/sql"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upCreateAccounts, downCreateAccounts)
}

func upCreateAccounts(tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.Exec(`CREATE TABLE IF NOT EXISTS accounts (
		id VARCHAR(14) NOT NULL,
		name VARCHAR(80) NOT NULL,
		document_number VARCHAR(255),
		created_at INT(11) NOT NULL,
		updated_at INT(11) NOT NULL,
		PRIMARY KEY (id)
	);`)

	return err
}

func downCreateAccounts(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE IF EXISTS accounts`)
	return err
}
