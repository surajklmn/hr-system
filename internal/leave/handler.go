package leave

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

	// Store the leaveReuqest fileds later
	var req CreateLeaveRequest
	// Decode the incomding json stream into go struct
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		// 400 bad request
		http.Error(w, err.Error(), http.StatusBadRequest)
		// return if the decoding fails
		return
	}

	leave, err := h.service.Create(r.Context(), req)

	if err != nil {
		if errors.Is(err, ErrEmployeeNotFound) {
			// 404
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		// 400
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	//Return 201 created message
	w.WriteHeader(http.StatusCreated)

	// Translate the GO struct back into JSON
	json.NewEncoder(w).Encode(leave)
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
		http.Error(w, "invalid leave id", http.StatusBadRequest)
		return
	}
	leave, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			// 404 error if the leave can't be fetched using GetByID
			http.Error(w, "leave request not found", http.StatusNotFound)
			// Always return after an error
			return
		}
		// 500 error
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// 200 status  code
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(leave)

}

func (h *Handler) Approve(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil {
		http.Error(w, "invalid leave id", http.StatusBadRequest)
		return
	}

	leave, err := h.service.Approve(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			http.Error(w, "leave request not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, ErrInvalidStatusTransition) {
			// 409 conflict
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// 200 status  code
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(leave)
}

func (h *Handler) Reject(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil {
		http.Error(w, "invalid leave id", http.StatusBadRequest)
		return
	}

	leave, err := h.service.Reject(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			http.Error(w, "leave request not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, ErrInvalidStatusTransition) {
			// 409 conflict
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// 200 status  code
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(leave)
}
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/leaves", h.Create)
	mux.HandleFunc("GET /api/v1/leaves/{id}", h.GetByID)

	// Use mux.Handle for middlerware, mux.HandleFunc for raw functions
	mux.Handle("POST /api/v1/leaves/{id}/approve", middleware.RequireRole("admin", "manager")(http.HandlerFunc(h.Approve)))
	mux.Handle("POST /api/v1/leaves/{id}/reject", middleware.RequireRole("admin", "manager")(http.HandlerFunc(h.Reject)))
}
