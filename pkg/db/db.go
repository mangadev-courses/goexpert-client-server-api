package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func DbClient() *sql.DB {
	db, err := sql.Open("sqlite3", "cotacao.db")
	if err != nil {
		panic(err)
	}

	err = CreateTable(db)
	if err != nil {
		panic(err)
	}

	return db
}

func CreateTable(db *sql.DB) error {
	_, err := db.Exec(`
			CREATE TABLE IF NOT EXISTS cotacao (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				bid TEXT NOT NULL,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			)
		`)
	return err
}

func InsertBid(ctx context.Context, db *sql.DB, bid string) error {
	stmt, err := db.Prepare("INSERT INTO cotacao (bid) VALUES (?)")
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, bid)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}
