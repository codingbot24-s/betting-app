package userService

// This is a simple user service that starts an HTTP server on port 8080.

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func StartServer() {

	router := mux.NewRouter()

	fmt.Println("Starting User Service on port 8080")

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
}
