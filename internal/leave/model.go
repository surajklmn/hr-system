package leave

import(
	"time"
	"errors"
)
const dateLayout = "2006-01-02"

type LeaveRequest struct{
	ID int64 `json:"id"`
	EmployeeID int64 `json:"employee_id"`
	LeaveType  string `json:"leave_type"`
	StartDate  string `json:"start_date"`
	EndDate string `json:"end_date"`
	Reason string `json:"reason"`
	Status string `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}


type CreateLeaveRequest struct{
	EmployeeID int64`json:"employee_id"`
	LeaveType  string `json:"leave_type"`
	StartDate  string `json:"start_date"`
	EndDate string `json:"end_date"`
	Reason string `json:"reason"`
}

func(r *CreateLeaveRequest) Validate() error{
	if r.EmployeeID <= 0{
		return errors.New("EmployeeID can't be less than 0")
	}
	
	switch r.LeaveType{
	case "ANNUAL", "SICK", "UNPAID":
		// Do nothing
	default:
		return errors.New("invalid leave type: must be annual, sick or unpaid")
	}
	start, err := time.Parse(dateLayout,r.StartDate)
	// I wasn't sure which convention to use for year-date
	if err!=nil{
		return errors.New("invalid start_date format, expected YYYY-MM-DD")
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
