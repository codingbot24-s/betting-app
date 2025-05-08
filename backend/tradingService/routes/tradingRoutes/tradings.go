package tradingRoutes

import (
	"net/http"

	"github.com/codingbot24-s/handlers/tradingHandlers"
	"github.com/codingbot24-s/middlewares"
	"github.com/gorilla/mux"
)


func SetupTradingRoutes(router *mux.Router) *mux.Router {
	router.Handle("/buy", middlewares.AuthMiddleware(http.HandlerFunc(tradingHandlers.Buy))).Methods("POST")
	
	return router
}





