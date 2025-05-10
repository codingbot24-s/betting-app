package middleware

import (
	"fmt"
	"net/http"

	"github.com/codingbot24-s/helpers"
)

// if any service wt to call this middleware, its compulsory to send the token cointanined in the header

func AdminCheckerMiddleware(next http.Handler) http.Handler {
	// Parse the user id then check 

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the authorization header
		token := r.Header.Get("Authorization")
		fmt.Println("request recived with this ",token)
		if token == "" {
			http.Error(w, "Unauthorized token not found", http.StatusUnauthorized)
		}
		// verify the token
		fmt.Println("verifying token started")
		role, err := helpers.VerifyToken(token)
		if err != nil || role != "admin" {
			http.Error(w, "Unauthorized token not include admin role", http.StatusUnauthorized)
			return
		}

		// if the token is valid, call the next handler
		next.ServeHTTP(w, r)
	})
	

}
