package handlers

import (
	"log"
	"net/http"

	"go-concept/db"
)

type (
	Handler struct {
		DB *db.DB
	}
)

// HealthHandler returns a 200 if the service is healthy and a 500 if it is not
func (h *Handler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	if err := h.DB.Ping(); err != nil {
		log.Printf("health check failed: %s\n", err)
		writeHTTPResponse(w, http.StatusInternalServerError, map[string]string{"message": "I'm unhealthy"})
	} else {
		writeHTTPResponse(w, http.StatusOK, map[string]string{"message": "I'm healthy"})
	}
}
