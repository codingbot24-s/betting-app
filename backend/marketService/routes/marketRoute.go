package routes

import (
	"net/http"

	"github.com/codingbot24-s/handlers"
	"github.com/codingbot24-s/middleware"
	"github.com/gorilla/mux"
)


func SetupMarketRoutes(r *mux.Router) *mux.Router {
	r.Handle("/market",middleware.AdminCheckerMiddleware(http.HandlerFunc(handlers.CreatedMarketHandler))).Methods("POST")
	return r
}