package db

import (
	"context"
	"log"
	"os"
	"path"
	"time"

	"github.com/db-journey/migrate"
	_ "github.com/db-journey/sqlite3-driver"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type (
	DB struct {
		DB *sqlx.DB
	}
)

var URL = "./go-concept.db"

func NewDB(dbPath string) (*DB, error) {
	db, err := sqlx.Connect("sqlite3", dbPath+"?_foreign_keys=true")
	if err != nil {
		return nil, err
	}
	return &DB{DB: db}, nil
}

func RunMigrations(dbPath, migrationsPath string) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	var migrator *migrate.Handle

	backoff := []time.Duration{5, 10, 15, 20, 25}
	for i, timeout := range backoff {
		migrator, err = migrate.Open("sqlite3://"+dbPath, path.Join(dir, migrationsPath))
		if err != nil {
			log.Printf("migrator failed to open (try %d out of %d)", i+1, len(backoff))
			if i+1 == len(backoff) {
				return err
			}
			time.Sleep(timeout * time.Second)
		} else {
			break
		}
	}

	defer migrator.Close()

	return migrator.Up(context.Background())
}

// Ping is a tiny method to make sure the db is alive
func (db *DB) Ping() error {
	_, err := db.DB.Exec("SELECT 1")
	return err
}
