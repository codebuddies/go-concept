package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"go-concept/db"
)

var (
	testServer   *httptest.Server
	testHandlers *Handler
)

func TestMain(m *testing.M) {
	// setup & migrate test db
	testDB, err := db.NewDB("./test_handlers.db")
	if err != nil {
		panic(err)
	}

	if err := db.RunMigrations("./test_handlers.db", "../migrations"); err != nil {
		panic(err)
	}

	// setup test server
	testHandlers = &Handler{DB: testDB}
	testServer = httptest.NewServer(NewServer(testHandlers).Handler)
	defer testServer.Close()

	// run tests
	status := m.Run()

	// clean up db file. Comment out to inspect test db.
	if err := os.Remove("./test_handlers.db"); err != nil {
		panic(err)
	}

	os.Exit(status)
}

func TestHealthHandler(t *testing.T) {
	// success
	res, err := http.Get(testServer.URL + "/api")
	if err != nil {
		require.NoError(t, err)
	}

	requireResponse(t, res, http.StatusOK, "{\"meta\":{},\"response\":{\"message\":\"I'm healthy\"}}")

	// failure: need to create a new server since it only fails if the database
	// does not exist.
	badTestServer := httptest.NewServer(NewServer(&Handler{}).Handler)
	defer badTestServer.Close()

	res, err = http.Get(badTestServer.URL + "/api")
	if err != nil {
		require.NoError(t, err)
	}
	requireResponse(t, res, http.StatusInternalServerError, "{\"meta\":{},\"response\":{\"message\":\"runtime error: invalid memory address or nil pointer dereference\"}}")
}
