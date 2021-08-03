// Copyright 2021 Job Stoit. All rights reserved.

package ql

import (
	"context"
	"database/sql"

	"github.com/jobstoit/gqb/config"
)

// Querier is an interface for the *sql.DB and *sql.Tx to execute incomming queries
type Querier interface {
	Exec(q string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, q string, args ...interface{}) (sql.Result, error)
	Query(q string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, q string, args ...interface{}) (*sql.Rows, error)
	QueryRow(q string, args ...interface{}) (*sql.Row, error)
	QueryRowContext(ctx context.Context, q string, args ...interface{}) (*sql.Row, error)
}

// Open updates the database to the current configuration and returns the db
func Open(dvr string, cs string, cfg []byte) (*sql.DB, error) {
	db, err := sql.Open(dvr, cs)
	if err != nil {
		return nil, err
	}

	q := `CREATE TABLE IF NOT EXISTS config-version (
		id INT PRIMARY KEY NOT NULL,
		config TEXT NOT NULL
	);`

	if _, err := db.Exec(q); err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	q = `SELECT id, config FROM config-version ORDER BY id DESC`

	var version int
	var prevCfg []byte
	if err := db.QueryRow(q).Scan(&version, &prevCfg); err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	currentModel := config.Read(cfg)
	if version == 0 {
		initMigrate(db, currentModel)
	}

	// TODO
	//	if version > 0 && !bytes.Equal(prevCfg, []byte{}) && bytes.Equal(prevCfg, cfg) {
	//		previousModel := config.Read(prevCfg)
	//		for _, tp := range previousModel.Types {
	//			//
	//		}
	//	}

	return db, nil
}

func initMigrate(db *sql.DB, mdl config.Model) {
	// TODO
}
