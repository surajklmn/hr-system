package employee

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/surajklmn/hr-system/internal/database"
)

func setupTestDB(t *testing.T) *Repository {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "test.db")
	db, err := database.Open(dbPath)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	t.Cleanup(func() {
		db.Close()
	})
	return NewRepository(db)
}

func TestHandler_Create(t *testing.T) {
	repo := setupTestDB(t)
	svc := NewService(repo)
	handler := NewHandler(svc)

	t.Run("successful employee creation returns 201", func(t *testing.T) {
		reqBody := `{"name":"Alice Smith","email":"alice@example.com","department":"Engineering","role":"employee"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/employees", bytes.NewBufferString(reqBody))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("expected status 201, got %d (body: %s)", rec.Code, rec.Body.String())
		}

		var created Employee
		if err := json.NewDecoder(rec.Body).Decode(&created); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if created.ID <= 0 || created.Email != "alice@example.com" {
			t.Fatalf("unexpected employee record: %+v", created)
		}
	})

	t.Run("validation failure returns 400", func(t *testing.T) {
		reqBody := `{"name":"","email":"alice2@example.com","department":"Engineering"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/employees", bytes.NewBufferString(reqBody))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", rec.Code)
		}
		if !strings.Contains(rec.Body.String(), "validation failed") || !strings.Contains(rec.Body.String(), "name is required") {
			t.Fatalf("expected validation error message, got: %s", rec.Body.String())
		}
	})

	t.Run("duplicate email returns 409", func(t *testing.T) {
		reqBody := `{"name":"Alice Duplicate","email":"alice@example.com","department":"Engineering"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/employees", bytes.NewBufferString(reqBody))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		if rec.Code != http.StatusConflict {
			t.Fatalf("expected status 409, got %d (body: %s)", rec.Code, rec.Body.String())
		}
		if !strings.Contains(rec.Body.String(), "email already exists") {
			t.Fatalf("expected duplicate email error message, got: %s", rec.Body.String())
		}
	})

	t.Run("server database failure returns 500", func(t *testing.T) {
		// Close underlying db to simulate database failure
		repo.db.Close()

		reqBody := `{"name":"Bob Jones","email":"bob@example.com","department":"HR"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/employees", bytes.NewBufferString(reqBody))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("expected status 500, got %d (body: %s)", rec.Code, rec.Body.String())
		}
		if strings.TrimSpace(rec.Body.String()) != "internal server error" {
			t.Fatalf("expected 'internal server error', got: %q", rec.Body.String())
		}
	})
}
