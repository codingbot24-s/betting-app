package authRoutes

import (
	"net/http"

	handlers "github.com/codingbot24-s/handlers/authhandler"
	"github.com/codingbot24-s/middlewares"
	"github.com/gorilla/mux"
)

func SetupAuthRoutes(router *mux.Router) *mux.Router {

	// API routes
	api := router.PathPrefix("/api").Subrouter()

	// Auth routes
	api.HandleFunc("/auth/register", handlers.CreateUser).Methods("POST")
	api.HandleFunc("/auth/login", handlers.Login).Methods("POST")
	api.Handle("/auth/me", middlewares.AuthMiddleware(http.HandlerFunc(handlers.GetSingleUser))).Methods("GET")
	return router
}
