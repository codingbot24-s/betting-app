package main

import (
    
    "log"
    "net/http"

    "github.com/codingbot24-s/db"
    "github.com/codingbot24-s/helpers"
    "github.com/codingbot24-s/routes/BalanceRoutes"
    "github.com/codingbot24-s/routes/authRoutes"
    "github.com/gorilla/mux"
)

func main() {
    log.Println("Initializing User Service...")

    // Initialize database
    db.InitDB()

    // Setup router
    router := mux.NewRouter()
    router = authRoutes.SetupAuthRoutes(router)
    router = BalanceRoutes.SetupBalanceRoutes(router)
    log.Println("Routes configured successfully")

    // Start Kafka consumer in a separate goroutine
    go func() {
        log.Println("Starting Kafka consumer...")
        helpers.ReadFromKafka(db.DB)
    }()

    // Start HTTP server
    addr := ":8080"
    server := &http.Server{
        Addr:    addr,
        Handler: router,
    }

    log.Printf("Starting HTTP server on %s", addr)
    if err := server.ListenAndServe(); err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }
}