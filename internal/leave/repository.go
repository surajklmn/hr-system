package leave

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var ErrNotFound = errors.New("leave request not found")

type Repository struct {
	db *sql.DB
}

// This constructor helps to acceess the private db inside Repository struct and open a connetion.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// It takes an empty leaveRequest struct and populates the data into db
func (repo *Repository) Create(ctx context.Context, leave *LeaveRequest) (*LeaveRequest, error) {
	query := `INSERT INTO leave_requests(employee_id,leave_type,start_date, end_date, reason, status, created_at, updated_at) VALUES(?,?,?,?,?,?,?,?);`

	// ExecContext, it executes the query without returning any row
	result, err := repo.db.ExecContext(ctx, query, leave.EmployeeID, leave.LeaveType, leave.StartDate, leave.EndDate, leave.Reason, leave.Status, leave.CreatedAt, leave.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Autoincrement the id
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	leave.ID = id
	return leave, nil
}

func (repo *Repository) GetByID(ctx context.Context, id int64) (*LeaveRequest, error) {
	query := `SELECT id, employee_id, leave_type, start_date, end_date, reason, status, created_at, updated_at FROM leave_requests WHERE ID = ?`
	// executes a query that returns at most one row
	row := repo.db.QueryRowContext(ctx, query, id)

	var leave LeaveRequest
	err := row.Scan(&leave.ID, &leave.EmployeeID, &leave.LeaveType, &leave.StartDate, &leave.EndDate, &leave.Reason, &leave.Status, &leave.CreatedAt, &leave.UpdatedAt)
	if err != nil {

		// We check if the error is same as the employee not found, if it's the same then we return our self defined error
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &leave, nil
}

func (repo *Repository) UpdateStatus(ctx context.Context, id int64, fromStatus, toStatus Status, updatedAt time.Time) error {
	query := `UPDATE leave_requests SET status = ?, updated_at = ? WHERE id = ? AND status = ?`

	result, err := repo.db.ExecContext(ctx, query, toStatus, updatedAt, id, fromStatus)
	if err != nil {
		return err
	}

	// Returns the rows affected by an update
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// If id doesn't exists or the status wasn't fromStatus
	if rows == 0 {
		existing, err := repo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if existing.Status != fromStatus {
			return ErrInvalidStatusTransition
		}
	}

	return nil
}
