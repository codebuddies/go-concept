package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-concept/db"
	"go-concept/handlers"
)

func main() {
	if err := db.RunMigrations(db.URL, "migrations"); err != nil {
		log.Printf("failed to run migrations: %s\n", err)
		os.Exit(1)
	} else {
		log.Println("migrations completed successfully")
	}

	db, err := db.NewDB(db.URL)
	if err != nil {
		log.Printf("unable to open db: %s\n", err)
		os.Exit(1)
	}

	server := handlers.NewServer(&handlers.Handler{DB: db})

	go func() {
		log.Println("server listening at 3000")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	graceful(server, 5*time.Second)
}

func graceful(hs *http.Server, timeout time.Duration) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := hs.Shutdown(ctx); err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		log.Println("Server stopped")
	}
}
