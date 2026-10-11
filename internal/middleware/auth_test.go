package middleware

import(
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequireRole(t *testing.T){
	// dummy handler that returns 200 OK if the middleware lets the request pass
	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		w.WriteHeader(http.StatusOK)
	})

	protected := RequireRole("admin","manager")(dummyHandler)

	tests := []struct{
		name string
		roleHeader string
		expectedStatus int
	}{
		{
			name: "missing role header should return 401",
			roleHeader: "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:"unauthorized role should return 403",
			roleHeader: "employee",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:"authorized role should return 200",
			roleHeader: "manager",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests{
		t.Run(tt.name, func(t *testing.T){
			// Simulate an incoming HTTP request
			req := httptest.NewRequest("POST","/test",nil)
			if tt.roleHeader != "" {
				req.Header.Set("X-User-Role",tt.roleHeader)
			}

			// Capture whatever HTTP response the middleware sends back
			rec := httptest.NewRecorder()
			// run the request through the middleware
			protected.ServeHTTP(rec,req)

			// inspect the recorded status code 
			if rec.Code != tt.expectedStatus{
				t.Errorf("expected status %d, got %d",tt.expectedStatus,rec.Code)
			}
		})
	}
}
