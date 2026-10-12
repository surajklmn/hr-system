package employee

import (
	"context"
	"errors"
	"time"
)

// Service manages business logic and validation for employee operations
type Service struct {
	//lowercase so private, we will have to use Constructors to access this field
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Create validates employee input, applies defaults, and persists the record.
// s is for service, we could have used svc too, idk
func (s *Service) Create(ctx context.Context, req CreateEmployeeRequest) (emp *Employee, err error) {
	err = req.Validate()
	if err != nil {
		return nil, err
	}
	if req.Role == "" {
		req.Role = RoleEmployee
	}
	now := time.Now().UTC()

	// create an instance of struct with values that we want to set
	// values like created_at and updated_at gets it's 0 value
	emp = &Employee{
		Name:       req.Name,
		Email:      req.Email,
		Department: req.Department,
		Role:       req.Role,
		Status:     StatusActive,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	//s.repo.Create return (*Employee, error)
	createdEmp, err := s.repo.Create(ctx, emp)
	if err != nil {
		return nil, err
	}
	return createdEmp, nil
}

// Get employee by id, if the handlers ask for #4 employee then it asks the repo s.repo.GetByID(ctx,id)
func (s *Service) GetByID(ctx context.Context, id int64) (*Employee, error) {
	return s.repo.GetByID(ctx, id)

}

// Method to check if the employee's entry already exists or not
func (s *Service) Exists(ctx context.Context, id int64) (bool, error) {
	_, err := s.GetByID(ctx, id)
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, ErrNotFound):
		return false, nil
	default:
		return false, err
	}

}

func (s *Service) IsActive(ctx context.Context, id int64)(bool, error){
	emp, err := s.GetByID(ctx, id)
	switch {
	case err == nil:
		return emp.Status == StatusActive, nil
	case errors.Is(err, ErrNotFound):
		return false, ErrNotFound
	
	default:
		return false, err
	}
}
