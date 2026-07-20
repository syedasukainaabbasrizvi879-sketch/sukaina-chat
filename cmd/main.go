package main

import (
	"log"
	"net/http"
	"os"

	"sukaina-chat/internal/database"
	"sukaina-chat/internal/handlers"
	"sukaina-chat/internal/websocket"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Frontend origin allow kiya
		w.Header().Set("Access-Control-Allow-Origin", "https://syedasukainaabbasrizvi879-sketch.github.io")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
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

	// Setup router
	mux := http.NewServeMux()

	mux.HandleFunc("/health", handlers.Health)
	mux.HandleFunc("/api/v1/auth/register", handlers.Register)
	mux.HandleFunc("/api/v1/auth/login", handlers.Login)
	mux.HandleFunc("/api/v1/messages", handlers.GetMessages)
	mux.HandleFunc("/ws", hub.HandleWebSocket)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Sukaina Chat Backend is Running!"))
	})

	// Get port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Sukaina Chat running on port %s", port)
	
	// Server start with CORS
	if err := http.ListenAndServe(":"+port, corsMiddleware(mux)); err != nil {
		log.Fatal("Server failed:", err)
	}
}
