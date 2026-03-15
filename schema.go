package main

import (
	"database/sql"
	"fmt"
)

// migrateDB creates the tables required by SQLRepo if they don't already exist.
// Column names match the db struct tags on Post and DocPage (S4).
// Called from main() after db.Ping() and before modules are wired.
func migrateDB(db *sql.DB) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS posts (
			id           TEXT PRIMARY KEY,
			slug         TEXT NOT NULL UNIQUE,
			status       TEXT NOT NULL DEFAULT 'draft',
			published_at DATETIME,
			scheduled_at DATETIME,
			created_at   DATETIME NOT NULL,
			updated_at   DATETIME NOT NULL,
			title        TEXT NOT NULL,
			body         TEXT NOT NULL,
			tags         TEXT NOT NULL DEFAULT '[]'
		)`,
		`CREATE TABLE IF NOT EXISTS doc_pages (
			id           TEXT PRIMARY KEY,
			slug         TEXT NOT NULL UNIQUE,
			status       TEXT NOT NULL DEFAULT 'draft',
			published_at DATETIME,
			scheduled_at DATETIME,
			created_at   DATETIME NOT NULL,
			updated_at   DATETIME NOT NULL,
			title        TEXT NOT NULL,
			body         TEXT NOT NULL,
			section      TEXT NOT NULL DEFAULT '',
			sort_order   INTEGER NOT NULL DEFAULT 0
		)`,
	}
	for _, s := range stmts {
		if _, err := db.Exec(s); err != nil {
			return fmt.Errorf("migrateDB: %w", err)
		}
	}
	return nil
}
