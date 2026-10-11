package leave

import "testing"

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
		t.Errorf("expected valide request to pass, but got error: %v", err)
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
		t.Errorf("expected an error when end date come before the start date: %v", err)
	}
}
