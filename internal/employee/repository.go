package employee

import (
	// This helps us to keep track of the users, if the users isn't there to accept a request then r.Context() drops the connection, this is used for timeout, cancellation and carrying data across layers
	"context"
	"database/sql"
	"errors"
	"strings"
)

var ErrNotFound = errors.New("employee not found")
var ErrDuplicateEmail = errors.New("email already exists")

type Repository struct {
	// unexported db(this is private to this package)
	db *sql.DB 
}

// This constructor helps to acceess the private db inside Repository struct and open a connetion.
func NewRepository(db *sql.DB) *Repository{
	return &Repository{db:db}
}

func (repo *Repository) Create(ctx context.Context, emp *Employee)(*Employee, error){

	query := `INSERT INTO employees (name, email, department, role, status,created_at, updated_at)VALUES(?,?,?,?,?,?,?)`
	// Run queries that modify data without returning rows , INSERT, UPDATE, DELETE
	result, err :=repo.db.ExecContext(ctx, query, emp.Name,emp.Email, emp.Department, emp.Role,emp.Status, emp.CreatedAt, emp.UpdatedAt)
	if err != nil{
		if strings.Contains(err.Error(),"UNIQUE constraint failed: employees.email"){
		return nil, ErrDuplicateEmail
	}
	return nil, err
}
	// Autoincrement the integer ID
	id, err:= result.LastInsertId()
	if err != nil{
		return nil, err
	}

	emp.ID = id
	return emp, nil
}

// We shouldn't use uint64 in autoincrementing id it might cause weird errors
func (repo *Repository) GetByID(ctx context.Context, id int64)(*Employee, error){
	query := `SELECT id, name, email, department, role, status, created_at, updated_at FROM employees WHERE id = ?`

	// Use QueryRowContext for queries that fetch from the database
	// and return row
	row := repo.db.QueryRowContext(ctx, query, id)

	var emp Employee

// Copy columns from the matched rows into the values pointed here
	err := row.Scan(&emp.ID, &emp.Name, &emp.Email, &emp.Department, &emp.Role, &emp.Status, &emp.CreatedAt, &emp.UpdatedAt)
	if err != nil{
		// We check if the error is same as the employee not found, if it's the same then we return our self defined error
		if errors.Is(err, sql.ErrNoRows){
			return nil, ErrNotFound
		}
		return nil,err
	}

	return &emp, nil

}
