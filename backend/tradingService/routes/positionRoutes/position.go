package positionRoutes

import (
	"net/http"

	handler "github.com/codingbot24-s/handlers/positionHandler"
	"github.com/codingbot24-s/middlewares"
	"github.com/gorilla/mux"
)


func SetUpPositionRoutes (r *mux.Router) *mux.Router {
	r.Handle("/position", middlewares.AuthMiddleware(http.HandlerFunc(handler.CreatePositionHandler))).Methods("POST")
	r.Handle("/position/{userid}", middlewares.AuthMiddleware(http.HandlerFunc(handler.GetUserPositionsHandler))).Methods("GET")
	return r
} 


