package sqlite

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"

	_ "modernc.org/sqlite"
)

type logger interface{}

type DB struct {
	*sql.DB
	logger logger
	dsn    string
	ctx    context.Context
	cancel context.CancelFunc
}

var InMemoryDSN = ":memory:"

//go:embed migrations/*.sql
var migrationsFS embed.FS

func NewDB(dsn string, logger logger) (*DB, error) {
	if dsn == "" {
		dsn = InMemoryDSN
	}

	sqlDB, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	db := &DB{dsn: dsn, DB: sqlDB, logger: logger, cancel: cancel, ctx: ctx}

	if dsn != InMemoryDSN {
		if err := os.MkdirAll(filepath.Dir(dsn), 0700); err != nil {
			db.Close()
			return nil, err
		}
	}

	_, err = db.ExecContext(db.ctx, `PRAGMA foreign_keys = ON`)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("cannot enable foreign keys: %w", err)
	}

	if err = db.migrate(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func (db *DB) Close() error {
	db.cancel()

	if db.DB != nil {
		return db.DB.Close()
	}
	return nil
}

func (db *DB) migrate() error {
	if _, err := db.ExecContext(db.ctx, `CREATE TABLE IF NOT EXISTS migrations (filename TEXT PRIMARY KEY)`); err != nil {
		return fmt.Errorf("cannot create migrations table: %w", err)
	}

	names, err := fs.Glob(migrationsFS, "migrations/*.sql")
	if err != nil {
		return err
	}
	sort.Strings(names)

	for _, name := range names {
		if err := db.migrateFile(name); err != nil {
			return fmt.Errorf("migration error: name=%q err=%w", name, err)
		}
	}

	return err
}

func (db *DB) migrateFile(filename string) error {
	tx, err := db.BeginTx(db.ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var c int
	if err := tx.QueryRowContext(db.ctx, `
		SELECT COUNT(*)
		FROM migrations
		WHERE filename = ?
		`, filename).Scan(&c); err != nil {
		return err
	} else if c > 0 {
		return nil
	}

	if buf, err := fs.ReadFile(migrationsFS, filename); err != nil {
		return err
	} else if _, err := tx.ExecContext(db.ctx, string(buf)); err != nil {
		return err
	}

	if _, err := tx.ExecContext(db.ctx, `
		INSERT INTO migrations (filename)
		VALUES (?)
		`, filename); err != nil {
		return err
	}

	return tx.Commit()
}
