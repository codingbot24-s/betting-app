package main

import (
	
	"fmt"
	"log"
	"net/http"

	"github.com/codingbot24-s/db"
	"github.com/codingbot24-s/helpers"
	"github.com/codingbot24-s/routes/tradingRoutes"
	"github.com/segmentio/kafka-go"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	db.ConnectToDB()

	r = tradingRoutes.SetupTradingRoutes(r)

	writer := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "order.created",
		Balancer: &kafka.LeastBytes{},
	}

	helpers.StartOutboxProcessor(db.DB, writer)

	fmt.Println("Starting Trading Service on port 8081")
	if err := http.ListenAndServe(":8081", r); err != nil {
		log.Fatal(err)
	}

}
