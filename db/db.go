package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/nathanbizkit/article-management/env"
)

// New returns a database pool connection
func New(e *env.ENV) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		e.DBUser, e.DBPass, e.DBHost, e.DBPort, e.DBName)

	var d *sql.DB
	var err error
	for i := 0; i < 10; i++ {
		d, err = sql.Open("postgres", psqlInfo)
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		return nil, err
	}

	err = d.Ping()
	if err != nil {
		return nil, err
	}

	return d, nil
}

// RunInTx wraps database operations with a transaction
func RunInTx(db *sql.DB, fn func(tx *sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	err = fn(tx)
	if err == nil {
		return tx.Commit()
	}

	rollbackErr := tx.Rollback()
	if rollbackErr != nil {
		return errors.Join(err, rollbackErr)
	}

	return err
}
