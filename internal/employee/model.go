package employee

import (
	"errors"
	"strings"
	"time"
)

type Employee struct{
	ID int64 `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Department string `json:"department"`
	Role string `json:"role"`
	Status string `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateEmployeeRequest struct{
	Name string `json:"name"`
	Email string `json:"email"`
	Department string `json:"department"`
	Role string `json:"role"`
}
// Validates the required filed while creating and employee request

func (r *CreateEmployeeRequest) Validate() error{
	if strings.TrimSpace(r.Name) == ""{
		return errors.New("name is required")
	}

	if strings.TrimSpace(r.Email) == ""{
		return errors.New("email is required")
	}

	// We can do this another way by using mail.ParseAddress(repo.Email)
	// Will do it in a refactor
	if !strings.Contains(r.Email,"@"){
		return errors.New("invalid email address")
	}

	if strings.TrimSpace(r.Department) == ""{
		return errors.New("department is required")
	}

// switch because too many ifs
	switch r.Role {
	case "","admin", "manager", "employee":
		// Do nothing
	default:
		return errors.New("invalide role: must be admin, manager, or employee")
	}
	return nil
}

