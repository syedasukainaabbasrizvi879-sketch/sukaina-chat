package main

import (
    "log"
    "net/http"
    "os"

    "sukaina-chat/internal/database"
    "sukaina-chat/internal/handlers"
    "sukaina-chat/internal/websocket"
)

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        next(w, r)
    }
}

func main() {
    // Get database URL
    databaseURL := os.Getenv("DATABASE_URL")
    if databaseURL == "" {
        databaseURL = "postgres://sukaina:password@localhost:5432/sukaina_chat?sslmode=disable"
    }

    // Connect to database
    if err := database.Connect(databaseURL); err != nil {
        log.Fatal("Database connection failed:", err)
    }

    // Initialize schema
    if err := database.InitSchema(); err != nil {
        log.Fatal("Schema initialization failed:", err)
    }

    // Create WebSocket hub
    hub := websocket.NewHub()
    go hub.Run()

    // Setup routes
    http.HandleFunc("/health", corsMiddleware(handlers.Health))
    http.HandleFunc("/api/v1/auth/register", corsMiddleware(handlers.Register))
    http.HandleFunc("/api/v1/auth/login", corsMiddleware(handlers.Login))
    http.HandleFunc("/api/v1/messages", corsMiddleware(handlers.GetMessages))
    http.HandleFunc("/ws", hub.HandleWebSocket)

    // Get port
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("🚀 Sukaina Chat running on port %s", port)
    if err := http.ListenAndServe(":"+port, nil); err != nil {
        log.Fatal("Server failed:", err)
    }
}
