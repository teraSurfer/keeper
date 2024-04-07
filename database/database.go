package database

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var ddl string

type Database interface {
	Health() map[string]string
	Close()
}

type service struct {
	db *sql.DB
}

func New(dbName string) Database {
	ctx := context.Background()
	db, err := sql.Open("sqlite3", dbName)

	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		log.Fatal(err)
	}

	// create tables
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		// failed to create tables.
		log.Fatal(err)
	}

	return &service{db: db}
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.db.PingContext(ctx)
	if err != nil {
		log.Fatalf(fmt.Sprintf("db down: %v", err))
	}

	return map[string]string{
		"message": "It's healthy",
	}
}

func (s *service) Close() {
	err := s.db.Close()
	if err != nil {
		log.Fatal(err)
	}
}
