package migrations

import (
	"database/sql"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upAlterTransactionsAddBalance, downAlterTransactionsAddBalance)
}

func upAlterTransactionsAddBalance(tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	// This code is executed when the migration is applied.
	_, err := tx.Exec(`ALTER TABLE transactions ADD COLUMN balance DECIMAL(5,2) DEFAULT 0`)
	return err
}

func downAlterTransactionsAddBalance(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	_, err := tx.Exec(`ALTER TABLE transactions DROP COLUMN balance DECIMAL(5,2) DEFAULT 0`)
	return err
}
