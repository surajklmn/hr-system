package employee

import (
	"errors"
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name      string
		req       CreateEmployeeRequest
		expectErr bool
	}{{
		name: "valid employee",
		req: CreateEmployeeRequest{
			Name:       "Ram Bahadur",
			Email:      "ram@honda.com",
			Department: "Sales",
			Role:       "manager",
		},
		expectErr: false,
	},
		{
			name: "missing name",
			req: CreateEmployeeRequest{
				Name:       "",
				Email:      "ram@honda.com",
				Department: "Sales",
				Role:       "manager",
			},
			expectErr: true, // We expect an error if name is missing
		},
		{
			name: "missing email",
			req: CreateEmployeeRequest{
				Name:       "Ram Bahadur",
				Email:      "",
				Department: "Sales",
				Role:       "manager",
			},
			expectErr: true,
		},
		{
			name: "missing department",
			req: CreateEmployeeRequest{
				Name:       "Ram Bahadur",
				Email:      "ram@honda.com",
				Department: "",
				Role:       "manager",
			},
			expectErr: true,
		},
		{
			name: "invalid email format missing @ and domain",
			req: CreateEmployeeRequest{
				Name:       "Ram Bahadur",
				Email:      "bad-email",
				Department: "Sales",
				Role:       "manager",
			},
			expectErr: true,
		},
		{
			name: "invalid email format missing domain",
			req: CreateEmployeeRequest{
				Name:       "Ram Bahadur",
				Email:      "foo@",
				Department: "Sales",
				Role:       "manager",
			},
			expectErr: true,
		},
		{
			name: "invalid email format @example.com",
			req: CreateEmployeeRequest{
				Name:       "Ram Bahadur",
				Email:      "@example.com",
				Department: "Sales",
				Role:       "manager",
			},
			expectErr: true,
		},
		{
			name: "valid email address",
			req: CreateEmployeeRequest{
				Name:       "Ram Bahadur",
				Email:      "foo@bar.com",
				Department: "Sales",
				Role:       "manager",
			},
			expectErr: false,
		},
		{

			name: "unauthorized role",
			req: CreateEmployeeRequest{
				Name:       "Ram Bahadur",
				Email:      "ram@honda.com",
				Department: "Sales",
				Role:       "superhero",
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()

			if tt.expectErr && err == nil {
				t.Errorf("expected an error, but got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("expected no error, but got: %v", err)
			}
			if tt.expectErr && !errors.Is(err, ErrValidation) {
				t.Errorf("expected error to wrap ErrValidation, but got: %v", err)
			}
		})
	}
}
