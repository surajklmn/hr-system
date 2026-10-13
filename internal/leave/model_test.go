package leave

import (
	"errors"
	"testing"
)

// a valid request test, it should pass
func TestValidate_Valid(t *testing.T) {
	req := CreateLeaveRequest{
		EmployeeID: 1,
		LeaveType:  "ANNUAL",
		StartDate:  "2026-11-01",
		EndDate:    "2026-11-05",
	}

	err := req.Validate()
	if err != nil {
		t.Errorf("expected valid request to pass, but got error: %v", err)
	}
}

func TestValidate_EndDateBeforeStartDate(t *testing.T) {
	req := CreateLeaveRequest{
		EmployeeID: 1,
		LeaveType:  "ANNUAL",
		StartDate:  "2026-11-10",
		EndDate:    "2026-11-05",
	}

	err := req.Validate()
	if err == nil {
		t.Errorf("expected an error when end date comes before start date")
	}
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected error to wrap ErrValidation, got: %v", err)
	}
}

func TestValidate_ValidationErrors(t *testing.T) {
	tests := []struct {
		name string
		req  CreateLeaveRequest
	}{
		{
			name: "invalid employee ID",
			req: CreateLeaveRequest{
				EmployeeID: 0,
				LeaveType:  LeaveTypeAnnual,
				StartDate:  "2026-11-01",
				EndDate:    "2026-11-05",
			},
		},
		{
			name: "invalid leave type",
			req: CreateLeaveRequest{
				EmployeeID: 1,
				LeaveType:  "MATERNITY",
				StartDate:  "2026-11-01",
				EndDate:    "2026-11-05",
			},
		},
		{
			name: "invalid start date format",
			req: CreateLeaveRequest{
				EmployeeID: 1,
				LeaveType:  LeaveTypeAnnual,
				StartDate:  "01-11-2026",
				EndDate:    "2026-11-05",
			},
		},
		{
			name: "invalid end date format",
			req: CreateLeaveRequest{
				EmployeeID: 1,
				LeaveType:  LeaveTypeAnnual,
				StartDate:  "2026-11-01",
				EndDate:    "2026/11/05",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !errors.Is(err, ErrValidation) {
				t.Errorf("expected error to wrap ErrValidation, got: %v", err)
			}
		})
	}
}
