package main

import (
	"belajar-go/internal/category"
	"belajar-go/internal/product"
	"belajar-go/internal/transaction"
	"belajar-go/pkg/database"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// Initialize database connection
	db, err := database.NewPostgresConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Database connected successfully")

	// Initialize Category dependencies
	categoryRepo := category.NewRepository(db)
	categoryService := category.NewService(categoryRepo)
	categoryHandler := category.NewHandler(categoryService)

	// Initialize Product dependencies
	productRepo := product.NewRepository(db)
	productService := product.NewService(productRepo)
	productHandler := product.NewHandler(productService)

	// Initialize Transaction dependencies
	transactionRepo := transaction.NewRepository(db)
	transactionService := transaction.NewService(transactionRepo)
	transactionHandler := transaction.NewHandler(transactionService)

	// Setup router
	mux := http.NewServeMux()

	// Category Routes
	mux.HandleFunc("GET /categories", categoryHandler.GetAll)
	mux.HandleFunc("POST /categories", categoryHandler.Create)
	mux.HandleFunc("GET /categories/{id}", categoryHandler.GetByID)
	mux.HandleFunc("PUT /categories/{id}", categoryHandler.Update)
	mux.HandleFunc("DELETE /categories/{id}", categoryHandler.Delete)

	// Product Routes
	mux.HandleFunc("GET /products", productHandler.GetAll)
	mux.HandleFunc("POST /products", productHandler.Create)
	mux.HandleFunc("GET /products/{id}", productHandler.GetByID)
	mux.HandleFunc("PUT /products/{id}", productHandler.Update)
	mux.HandleFunc("DELETE /products/{id}", productHandler.Delete)

	// Transaction Routes
	mux.HandleFunc("POST /api/checkout", transactionHandler.Checkout)

	// Report Routes
	mux.HandleFunc("GET /api/report/hari-ini", transactionHandler.GetDailySalesReport)

	// Health Check Route
	mux.HandleFunc("GET /health", healthCheck)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server starting on http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}

// Health Check Handler
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "OK",
		"message": "Server is running smoothly",
	})
}
