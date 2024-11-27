package db

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestInsertBid(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open in-memory database: %v", err)
	}
	defer db.Close()

	err = CreateTable(db)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	ctx := context.Background()
	bid := "5.123"
	err = InsertBid(ctx, db, bid)

	if err != nil {
		t.Errorf("InsertBid failed: %v", err)
	}

	var retrievedBid string
	err = db.QueryRow(`SELECT bid FROM cotacao LIMIT 1`).Scan(&retrievedBid)
	if err != nil {
		t.Fatalf("Failed to query bid: %v", err)
	}

	if retrievedBid != bid {
		t.Errorf("Expected bid %s, got %s", bid, retrievedBid)
	}
}
