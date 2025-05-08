package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/codingbot24-s/db"
	"github.com/gorilla/mux"
	"github.com/codingbot24-s/routes"
)

func main() {
	r := mux.NewRouter()
	db.ConnectToDB()

	routes.SetupTeamRoutes(r)

	fmt.Println("Starting Team Service on port 8082")
	if err := http.ListenAndServe(":8082", r); err != nil {
		log.Fatal(err)
	}
}
