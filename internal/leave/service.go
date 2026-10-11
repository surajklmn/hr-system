package leave

import (
	"context"
	"errors"
	"time"
)

var ErrInvalidStatusTransition = errors.New("leave request is not in PENDING status")
var ErrEmployeeNotFound = errors.New("employee not found")

// An interface so we can access Exists method in employee.Service from leave.Service
type EmployeeChecker interface {
	Exists(ctx context.Context, id int64) (bool, error)
}

type Service struct {
	repo       *Repository
	empChecker EmployeeChecker
}

func NewService(repo *Repository, empChecker EmployeeChecker) *Service {
	return &Service{
		repo:       repo,
		empChecker: empChecker,
	}
}

func (s *Service) Create(ctx context.Context, req CreateLeaveRequest) (*LeaveRequest, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()

	exists, err := s.empChecker.Exists(ctx, req.EmployeeID)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, ErrEmployeeNotFound
	}
	leave := &LeaveRequest{
		EmployeeID: req.EmployeeID,
		LeaveType:  req.LeaveType,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
		Reason:     req.Reason,
		Status:     StatusPending,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	return s.repo.Create(ctx, leave)
}

func (s *Service) GetByID(ctx context.Context, id int64) (*LeaveRequest, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Approve(ctx context.Context, id int64) (*LeaveRequest, error) {

	// Fetch leave from DB
	leave, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// check if the leave status is pending, this is just a rule, a leave's status must be pending to be approved
	if !leave.Status.CanTransitionTo(StatusApproved) {
		return nil, ErrInvalidStatusTransition
	}
	// Update status in DB
	now := time.Now().UTC()
	err = s.repo.UpdateStatus(ctx, id, StatusPending, StatusApproved, now)
	if err != nil {
		return nil, err
	}

	// Update the stuct which is in memory
	leave.Status = StatusApproved
	leave.UpdatedAt = now
	return leave, nil
}

func (s *Service) Reject(ctx context.Context, id int64) (*LeaveRequest, error) {
	leave, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	// If the current leave status is other than pending than we abort
	if !leave.Status.CanTransitionTo(StatusRejected) {
		return nil, ErrInvalidStatusTransition
	}

	// Thought about using the local timezone, but it can be a bad thing to code local time zone in server side.
	now := time.Now().UTC()
	err = s.repo.UpdateStatus(ctx, id, StatusPending, StatusRejected, now)
	if err != nil {
		return nil, err
	}

	leave.Status = StatusRejected
	leave.UpdatedAt = now
	return leave, nil
}
