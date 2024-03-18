package migrations

import (
	"database/sql"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upCreateTransactions, downCreateTransactions)
}

func upCreateTransactions(tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.Exec(`CREATE TABLE IF NOT EXISTS transactions (
		id VARCHAR(14) NOT NULL,
		amount DECIMAL(5,2),
		operation_type tinyint(2),
    	event_date int,
		created_at INT(11) NOT NULL,
		updated_at INT(11) NOT NULL,
		PRIMARY KEY (id)
	);`)

	return err
}

func downCreateTransactions(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE IF EXISTS transactions`)
	return err
}
