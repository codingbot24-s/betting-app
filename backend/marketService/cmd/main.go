package main

import (
	"fmt"
	"net/http"

	"github.com/codingbot24-s/db"
	"github.com/codingbot24-s/routes"
	"github.com/gorilla/mux"
)

func main() {
	db.ConnectDB()
	fmt.Println("Connected to database")
	router := mux.NewRouter()
	r := routes.SetupMarketRoutes(router)
	fmt.Println("Router initialized Now starting server")

	err := http.ListenAndServe(":8083", r)

	if err != nil {
		panic("failed to start server")
	}
}
