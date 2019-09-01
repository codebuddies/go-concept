package handlers

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func NewServer(handler *Handler) *http.Server {
	r := mux.NewRouter()
	r.Use(RecoveryMiddleware, LoggingMiddleware)

	// health
	r.HandleFunc("/api", handler.HealthHandler).Methods("GET")

	// /resources
	r.HandleFunc("/api/resources", handler.ResourcesList).Methods("GET")
	r.HandleFunc("/api/resources", handler.ResourcesCreate).Methods("POST")
	r.HandleFunc("/api/resources/{id}", handler.ResourcesGet).Methods("GET")
	r.HandleFunc("/api/resources/{id}", handler.HealthHandler).Methods("PUT")
	r.HandleFunc("/api/resources/{id}", handler.HealthHandler).Methods("DELETE")

	server := http.Server{
		Handler:      r,
		Addr:         ":3000",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return &server
}
