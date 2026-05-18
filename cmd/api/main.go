package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/zonafirmann/next-go-pos-backend/internal/config"
	"github.com/zonafirmann/next-go-pos-backend/internal/handlers"

	// Menggunakan alias 'posMiddleware' agar tidak bentrok dengan middleware bawaan chi
	posMiddleware "github.com/zonafirmann/next-go-pos-backend/internal/middleware"
)

func main() {
	// 1. Initialize Database Connection
	db := config.ConnectDB()
	defer db.Close(context.Background())

	// 2. Initialize Chi Router
	r := chi.NewRouter()

	// 3. Mount Global Middleware (Security & Logging)
	r.Use(middleware.Logger)    // Logs every HTTP request to the terminal
	r.Use(middleware.Recoverer) // Prevents the server from crashing due to unexpected panics

	// Configure CORS to allow cross-origin requests from the Next.js frontend
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300, // Cache CORS preflight request for 5 minutes
	}))

	// Initialize Route Handlers
	authHandler := handlers.AuthHandler{DB: db}
	productHandler := handlers.ProductHandler{DB: db}

	// 4. Public Routes (No authentication required)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "online", "message": "Next-Go Enterprise POS Backend Engine is running", "version": "1.0.0"}`))
	})

	// 5. API Routes (Prefixed with /api)
	r.Route("/api", func(r chi.Router) {
		// Public authentication routes
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)

		// Protected routes (Requires valid JWT Token)
		r.Group(func(r chi.Router) {
			// Mount our custom JWT middleware (Satpam) to this specific group
			r.Use(posMiddleware.RequireAuth)

			r.Get("/products", productHandler.GetProducts)
			r.Post("/checkout", productHandler.Checkout)
		})
	})

	// 6. Port Allocation and Server Boot
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 [SERVER] API Engine is preparing for liftoff on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("❌ [SERVER] Failed to start: %v", err)
	}
}
