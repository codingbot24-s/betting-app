package routes

import (
	"net/http"

	"github.com/codingbot24-s/handlers"
	"github.com/codingbot24-s/middleware"
	"github.com/gorilla/mux"
)


func SetupMarketRoutes(r *mux.Router) *mux.Router {
	// create market 
	r.Handle("/market",middleware.AdminCheckerMiddleware(http.HandlerFunc(handlers.CreatedMarketHandler))).Methods("POST")
	// list active markets
	r.HandleFunc("/markets",handlers.ListActiveMarketsHandler).Methods("GET")
	// list closed markets
	r.Handle("/markets/closed",middleware.AdminCheckerMiddleware(http.HandlerFunc(handlers.ListClosedMarketsHandler))).Methods("GET")
	// resolve market
	r.Handle("/market/resolved/{id}",middleware.AdminCheckerMiddleware(http.HandlerFunc(handlers.ResolvedMarketsHandler))).Methods("POST")
	// create a route for sending the market status to tradingservice with taking a market id  
	r.HandleFunc("/market/{id}",handlers.SendMarketStatus()).Methods("GET")

	// close the active market 
	r.Handle("/market/{id}/close",middleware.AdminCheckerMiddleware(http.HandlerFunc(handlers.CloseMarketHandler))).Methods("POST")
	return r


}
