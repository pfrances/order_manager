package sqlite_test

import (
	"order_manager/internal/sqlite"
	"testing"
)

func MustOpenDB(t *testing.T) *sqlite.DB {
	t.Helper()

	db, err := sqlite.NewDB(sqlite.InMemoryDSN, nil)
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	return db
}

func MustCloseDB(t *testing.T, db *sqlite.DB) {
	t.Helper()

	if err := db.Close(); err != nil {
		t.Fatalf("failed to close db: %v", err)
	}
}
