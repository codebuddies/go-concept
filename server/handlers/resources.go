package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"go-concept/db"
)

type (
	resourcesListRequest struct {
		Page    int    `json:"page"`
		PerPage int    `json:"per_page"`
		Query   string `json:"q"`
	}

	resourcesCreateRequest struct {
		Title       string          `json:"title"`
		Description *string         `json:"description"`
		URL         string          `json:"url"`
		Referrer    *string         `json:"referrer"`
		Credit      *string         `json:"credit"`
		Type        db.ResourceType `json:"type"`
	}
)

func (h *Handler) ResourcesList(w http.ResponseWriter, r *http.Request) {
	var req resourcesListRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, fmt.Errorf("could not parse request: %s", err))
		return
	}

	resources, err := h.DB.GetResources(req.Page, req.PerPage)
	if err != nil {
		writeHTTPError(w, http.StatusInternalServerError, err)
		return
	}

	writeHTTPResponse(w, http.StatusOK, resources)
}

func (h *Handler) ResourcesCreate(w http.ResponseWriter, r *http.Request) {
	var req resourcesCreateRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, fmt.Errorf("could not parse request: %s", err))
		return
	}

	resource := db.Resource{
		Title:       req.Title,
		Description: req.Description,
		URL:         req.URL,
		Referrer:    req.Referrer,
		Credit:      req.Credit,
		Type:        req.Type,
	}

	if err := h.DB.InsertResources(&resource); err != nil {
		writeHTTPError(w, http.StatusInternalServerError, err)
		return
	}

	writeHTTPResponse(w, http.StatusOK, resource)
}

func (h *Handler) ResourcesGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		writeHTTPError(w, http.StatusBadRequest, fmt.Errorf("could not parse id: %s", vars["id"]))
		return
	}
	resource, err := h.DB.GetResource(id)
	if err != nil {
		if err == sql.ErrNoRows {
			writeHTTPError(w, http.StatusNotFound, err)
			return
		}

		writeHTTPError(w, http.StatusInternalServerError, err)
		return
	}

	writeHTTPResponse(w, http.StatusOK, resource)
}
