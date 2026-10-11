package employee

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/surajklmn/hr-system/internal/middleware"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateEmployeeRequest
	// Decode the json into GO struct, in memory
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		// 400 bad request
		http.Error(w, err.Error(), http.StatusBadRequest)
		// return if the decoding fails
		return
	}
	createdEmp, err := h.service.Create(r.Context(), req)
	if err != nil {
		if errors.Is(err, ErrDuplicateEmail) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	//Return 201 created message
	w.WriteHeader(http.StatusCreated)

	// Translate the GO struct back into JSON
	json.NewEncoder(w).Encode(createdEmp)

}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	// First we need to get the id from the URL path
	// This return a id as a string
	idStr := r.PathValue("id")

	// Convert string to int
	// ParseInt takes base and bit size, basically converting string to base 10 64 bit integer here
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		// check for gibberish employee id, /api/v1/employees/abc
		http.Error(w, "invalid employee id", http.StatusBadRequest)
		return
	}

	emp, err := h.service.GetByID(r.Context(), id)
	if errors.Is(err, ErrNotFound) {
		// 404 error if the emp can't be fetched using GetByID
		http.Error(w, "employee not found", http.StatusNotFound)
		// Always return after an error
		return

	}
	if err != nil {
		// 500 error
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// 200 status  code
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(emp)

}

// to make our life easier, so we don't have to write this in main.go
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	// Only admin can create an entry for employees
	mux.Handle("POST /api/v1/employees", middleware.RequireRole("admin")(http.HandlerFunc(h.Create)))
	mux.HandleFunc("GET /api/v1/employees/{id}", h.GetByID)
}
