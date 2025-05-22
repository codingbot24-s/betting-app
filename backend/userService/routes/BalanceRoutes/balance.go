package BalanceRoutes

import (
	"net/http"

	handlers "github.com/codingbot24-s/handlers/balanceHandlers"
	"github.com/codingbot24-s/middlewares"
	"github.com/gorilla/mux"
)



func SetupBalanceRoutes(router *mux.Router) *mux.Router {

	// Private route for getting the balance with user id in the context
	router.Handle("/balance",middlewares.AuthMiddleware(http.HandlerFunc(handlers.GetBalancePrivate))).Methods("GET")

	// private route for adding and withdrawing balance
	router.Handle("/balance/add", middlewares.AuthMiddleware(http.HandlerFunc(handlers.AddBalance))).Methods("POST")
	router.Handle("/balance/withdraw", middlewares.AuthMiddleware(http.HandlerFunc(handlers.WithdrawBalance))).Methods("POST")
	
	// private route for getting the transaction history
	router.Handle("/balance/history", middlewares.AuthMiddleware(http.HandlerFunc(handlers.GetTransactionHistory))).Methods("GET")


	// public route for getting the balance
	router.Handle("/balance/public", http.HandlerFunc(handlers.GetBalancePublic)).Methods("GET")

	return router
}
