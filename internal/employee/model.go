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
	RoleAdmin Role = "admin"
	RoleManager Role = "manager"
	RoleEmployee Role = "employee"
)

type Status string

const (
	StatusActive Status = "ACTIVE"
	StatusInactive Status = "INACTIVE"
)

type Employee struct{
	ID int64 `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Department string `json:"department"`
	Role Role `json:"role"`
	Status Status `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateEmployeeRequest struct{
	Name string `json:"name"`
	Email string `json:"email"`
	Department string `json:"department"`
	Role Role `json:"role"`
}
// Validates the required filed while creating and employee request

func (r *CreateEmployeeRequest) Validate() error{
	r.Name = strings.TrimSpace(r.Name)
	if r.Name == ""{
		return errors.New("name is required")
	}
	r.Email = strings.TrimSpace(strings.ToLower(r.Email))
	if r.Email == ""{
		return errors.New("email is required")
	}

	if !strings.Contains(r.Email,"@"){
		_, err := mail.ParseAddress(r.Email)
		if err != nil{
			return errors.New("invalid email address")
		}
	}

	r.Department = strings.TrimSpace(r.Department)
	if r.Department == ""{
		return errors.New("department is required")
	}

// switch because too many ifs
	switch r.Role {
	case "",RoleAdmin, RoleManager, RoleEmployee:
		// Do nothing
	default:
		return fmt.Errorf("invalid role: must be %s, %s, or %s",RoleAdmin,RoleManager, RoleEmployee)
	}
	return nil
}

