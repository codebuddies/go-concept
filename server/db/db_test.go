package db

import (
	"os"
	"testing"

	_ "github.com/db-journey/sqlite3-driver"
	_ "github.com/mattn/go-sqlite3"
)

var testDB *DB

func TestMain(m *testing.M) {
	// setup & migrate test db
	db, err := NewDB("./test_db.db")
	if err != nil {
		panic(err)
	}
	testDB = db

	if err := RunMigrations("./test_db.db", "../migrations"); err != nil {
		panic(err)
	}

	// run tests
	status := m.Run()

	// clean up db file. Comment out to inspect test db.
	if err := os.Remove("./test_db.db"); err != nil {
		panic(err)
	}

	os.Exit(status)
}

func TestPing(t *testing.T) {
	if err := testDB.Ping(); err != nil {
		t.Error(err)
	}
}
