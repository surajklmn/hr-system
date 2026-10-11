package leave

import (
	"errors"
	"fmt"
	"time"
	"slices"
)

const dateLayout = "2006-01-02"

type LeaveType string 

const (
	LeaveTypeAnnual LeaveType = "ANNUAL"
	LeaveTypeSick LeaveType = "SICK"
	LeaveTypeUnpaid LeaveType = "UNPAID"
)

var validLeaveTypes = map[LeaveType]bool{
	LeaveTypeAnnual: true,
	LeaveTypeSick: true,
	LeaveTypeUnpaid: true,
}

func (lt LeaveType) IsValid() bool {
	return validLeaveTypes[lt]
}

type Status string

const (
	StatusPending Status = "PENDING"
	StatusApproved Status = "APPROVED"
	StatusRejected Status = "REJECTED"
)

var allowedTransitions = map[Status][]Status{
	StatusPending: {StatusApproved, StatusRejected},
	StatusApproved:{},
	StatusRejected:{},
}

func (current Status) CanTransitionTo(next Status) bool{
	return slices.Contains(allowedTransitions[current], next)
	}

type LeaveRequest struct{
	ID int64 `json:"id"`
	EmployeeID int64 `json:"employee_id"`
	LeaveType  LeaveType `json:"leave_type"`
	StartDate  string `json:"start_date"`
	EndDate string `json:"end_date"`
	Reason string `json:"reason"`
	Status Status `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}


type CreateLeaveRequest struct{
	EmployeeID int64`json:"employee_id"`
	LeaveType  LeaveType `json:"leave_type"`
	StartDate  string `json:"start_date"`
	EndDate string `json:"end_date"`
	Reason string `json:"reason"`
}

func(r *CreateLeaveRequest) Validate() error{
	if r.EmployeeID <= 0{
		return errors.New("EmployeeID must be a positive number")
	}

	
	switch r.LeaveType{
	case LeaveTypeAnnual, LeaveTypeSick, LeaveTypeUnpaid:
		// Do nothing
	default:
		return fmt.Errorf("invalid leave type: must be %s, %s or %s",LeaveTypeAnnual,LeaveTypeSick, LeaveTypeUnpaid)
	}
	start, err := time.Parse(dateLayout,r.StartDate)
	// I wasn't sure which convention to use for year-date
	if err!=nil{
		return fmt.Errorf("invalid start date format, expected YYYY-MM-DD")
	}
	end, err := time.Parse(dateLayout,r.EndDate)
	if err!=nil{
		return errors.New("invalid end_date format, expected YYYY-MM-DD")
	}
	
	// Before(start) reports if the start time come before the end time
	if end.Before(start){
		return errors.New("end_date cannot be before start_date")
	}
	
	return nil
}
