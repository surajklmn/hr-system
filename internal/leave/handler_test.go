package leave

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/surajklmn/hr-system/internal/database"
	"github.com/surajklmn/hr-system/internal/employee"
)

func setupTestDB(t *testing.T) (*Repository, *employee.Service, *sql.DB) {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "test.db")
	db, err := database.Open(dbPath)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	t.Cleanup(func() {
		db.Close()
	})
	empRepo := employee.NewRepository(db)
	empSvc := employee.NewService(empRepo)
	leaveRepo := NewRepository(db)
	return leaveRepo, empSvc, db
}

func TestHandler_Create(t *testing.T) {
	leaveRepo, empSvc, db := setupTestDB(t)
	svc := NewService(leaveRepo, empSvc)
	handler := NewHandler(svc)

	// Create an active employee
	activeEmp, err := empSvc.Create(context.Background(), employee.CreateEmployeeRequest{
		Name:       "Active User",
		Email:      "active@example.com",
		Department: "Engineering",
		Role:       "employee",
	})
	if err != nil {
		t.Fatalf("failed to create employee: %v", err)
	}

	// Create an inactive employee directly in db to test Exists behavior
	now := time.Now().UTC()
	res, err := db.Exec(`INSERT INTO employees (name, email, department, role, status, created_at, updated_at) VALUES (?, ?, ?, ?, 'INACTIVE', ?, ?)`,
		"Inactive User", "inactive@example.com", "Sales", "employee", now, now)
	if err != nil {
		t.Fatalf("failed to insert inactive employee: %v", err)
	}
	inactiveEmpID, err := res.LastInsertId()
	if err != nil {
		t.Fatalf("failed to get last insert id: %v", err)
	}

	t.Run("successful leave creation returns 201", func(t *testing.T) {
		reqBody := `{"employee_id":1,"leave_type":"ANNUAL","start_date":"2026-11-01","end_date":"2026-11-05","reason":"Vacation"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/leaves", bytes.NewBufferString(reqBody))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("expected status 201, got %d (body: %s)", rec.Code, rec.Body.String())
		}

		var created LeaveRequest
		if err := json.NewDecoder(rec.Body).Decode(&created); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if created.ID <= 0 || created.EmployeeID != activeEmp.ID {
			t.Fatalf("unexpected leave record: %+v", created)
		}
	})

	t.Run("filing leave for inactive employee succeeds because employee exists", func(t *testing.T) {
		reqBody := `{"employee_id":` + json.Number(string([]byte{byte('0' + inactiveEmpID)})).String() + `,"leave_type":"ANNUAL","start_date":"2026-11-10","end_date":"2026-11-12","reason":"Personal"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/leaves", bytes.NewBufferString(reqBody))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("expected status 201, got %d (body: %s)", rec.Code, rec.Body.String())
		}
	})

	t.Run("validation failure returns 400", func(t *testing.T) {
		reqBody := `{"employee_id":1,"leave_type":"ANNUAL","start_date":"2026-11-15","end_date":"2026-11-05"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/leaves", bytes.NewBufferString(reqBody))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d (body: %s)", rec.Code, rec.Body.String())
		}
		if !strings.Contains(rec.Body.String(), "validation failed") || !strings.Contains(rec.Body.String(), "end_date cannot be before start_date") {
			t.Fatalf("expected validation error message, got: %s", rec.Body.String())
		}
	})

	t.Run("employee not found returns 404", func(t *testing.T) {
		reqBody := `{"employee_id":99999,"leave_type":"ANNUAL","start_date":"2026-11-01","end_date":"2026-11-05"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/leaves", bytes.NewBufferString(reqBody))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("expected status 404, got %d (body: %s)", rec.Code, rec.Body.String())
		}
		if !strings.Contains(rec.Body.String(), "employee not found") {
			t.Fatalf("expected 'employee not found', got: %s", rec.Body.String())
		}
	})

	t.Run("server database failure returns 500", func(t *testing.T) {
		// Close underlying db to simulate database failure
		db.Close()

		reqBody := `{"employee_id":1,"leave_type":"ANNUAL","start_date":"2026-11-01","end_date":"2026-11-05"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/leaves", bytes.NewBufferString(reqBody))
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
