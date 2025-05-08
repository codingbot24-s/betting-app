package middlewares



import (
	"context"
	"net/http"

	"github.com/codingbot24-s/helpers"
)

type contextKey string

const UserIDKey contextKey = "userID"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// first get the token from the header
		authToken := r.Header.Get("Authorization")
		if authToken == "" {
			http.Error(w, "No token provided", http.StatusUnauthorized)
			return
		}
		// validate the token
		jwtToken, err := helpers.ValidateToken(authToken)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		// extract userid from token
		userID, err := helpers.GetUserIDFromToken(jwtToken)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		
		
		// pass it to the next handler
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}