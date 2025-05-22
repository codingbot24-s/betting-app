package middleware

import (
	"net/http"

	"github.com/codingbot24-s/helpers"
)



func AdminCheckerMiddleware(next http.Handler) http.Handler {
	// Parse the user id then check

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the authorization header
		token := r.Header.Get("Authorization")

		if token == "" {
			http.Error(w, "Unauthorized token not found", http.StatusUnauthorized)
			return
		}
		// verify the token

		role, err := helpers.VerifyToken(token)
		if err != nil || role != "admin" {
			http.Error(w, "Unauthorized token not include admin role", http.StatusUnauthorized)
			return
		}

		// if the token is valid, call the next handler
		next.ServeHTTP(w, r)
	})

}
