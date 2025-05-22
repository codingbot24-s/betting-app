package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/codingbot24-s/db"
	"github.com/codingbot24-s/helpers"
	"github.com/codingbot24-s/routes/positionRoutes"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	db.ConnectToDB()

	r = positionRoutes.SetUpPositionRoutes(r)

	go func() {
		helpers.ReadFromKafka(db.DB)		
	} ()	

	fmt.Println("Starting Trading Service on port 8081")
	if err := http.ListenAndServe(":8081", r); err != nil {
		log.Fatal(err)
	}

}
