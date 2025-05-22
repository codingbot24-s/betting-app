package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/codingbot24-s/db"
	"github.com/codingbot24-s/helpers"
	"github.com/codingbot24-s/routes"
	"github.com/gorilla/mux"
	"github.com/segmentio/kafka-go"
)

func main() {
	db.ConnectDB()
	fmt.Println("Connected to database")

	router := mux.NewRouter()
	r := routes.SetupMarketRoutes(router)
	fmt.Println("Router initialized Now starting server")

	db.StartMarketStatusUpdater()
	writer := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "market-resolved",
		Balancer: &kafka.RoundRobin{},
	}
	helpers.StartProcessor(context.Background(), db.DB, writer)
	fmt.Println("starting server on ")
	err := http.ListenAndServe(":8083", r)

	if err != nil {
		panic("failed to start server")
	}
}
