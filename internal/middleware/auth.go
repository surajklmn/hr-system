package middleware


import(
	"net/http"
	"strings"
)
// RequireRole returns a middleware that verifies the incoming request has an X-User-Role header matching one of the allowed roles.
func RequireRole(allowedRoles ...string) func(http.Handler) http.Handler{
	return func(next http.Handler)http.Handler{
		return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
// Reader the role from header
			role:= strings.TrimSpace(r.Header.Get("X-User-Role"))
// If the header is missing then it returns 401 Unauthorized
			if role == ""{
				http.Error(w,"unauthorized: missing X-User-Role header",http.StatusUnauthorized)
				return
			}

			authorized := false

			for _, allowed := range allowedRoles{
				// This compares role and allowed, it's case insensitive check
				if strings.EqualFold(role, allowed){
					authorized = true
					break
				}
			}
			if !authorized {
				// If role isnot allowed, 403 forbidden
				http.Error(w,"forbidden: not enough permissions",http.StatusForbidden)
				return
			}
			next.ServeHTTP(w,r)

		})
	}
}
