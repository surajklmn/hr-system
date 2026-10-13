package employee

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"
)

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleManager  Role = "manager"
	RoleEmployee Role = "employee"
)

type Status string

const (
	StatusActive   Status = "ACTIVE"
	StatusInactive Status = "INACTIVE"
)

type Employee struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Department string    `json:"department"`
	Role       Role      `json:"role"`
	Status     Status    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreateEmployeeRequest struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Department string `json:"department"`
	Role       Role   `json:"role"`
}

var ErrValidation = errors.New("validation failed")

// Validates the required filed while creating and employee request

func (r *CreateEmployeeRequest) Validate() error {
	r.Name = strings.TrimSpace(r.Name)
	if r.Name == "" {
		return fmt.Errorf("%w: name is required", ErrValidation)
	}
	r.Email = strings.TrimSpace(strings.ToLower(r.Email))
	if r.Email == "" {
		return fmt.Errorf("%w: email is required", ErrValidation)
	}

	addr, err := mail.ParseAddress(r.Email)
	if err != nil || addr.Address != r.Email {
		return fmt.Errorf("%w: invalid email address", ErrValidation)
	}

	r.Department = strings.TrimSpace(r.Department)
	if r.Department == "" {
		return fmt.Errorf("%w: department is required", ErrValidation)
	}

	// switch because too many ifs
	switch r.Role {
	case "", RoleAdmin, RoleManager, RoleEmployee:
		// Do nothing
	default:
		return fmt.Errorf("%w: invalid role: must be %s, %s, or %s", ErrValidation, RoleAdmin, RoleManager, RoleEmployee)
	}
	return nil
}
