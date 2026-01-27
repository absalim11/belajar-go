package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

// Category represents the category model
type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// In-memory storage
var (
	categories = []Category{}
	nextID     = 1
	mu         sync.Mutex
)

func main() {
	mux := http.NewServeMux()

	// Routes using Go 1.22+ enhanced routing
	mux.HandleFunc("GET /categories", listCategories)
	mux.HandleFunc("POST /categories", createCategory)
	mux.HandleFunc("GET /categories/{id}", getCategory)
	mux.HandleFunc("PUT /categories/{id}", updateCategory)
	mux.HandleFunc("DELETE /categories/{id}", deleteCategory)

	port := ":8080"
	fmt.Printf("Server starting on http://localhost%s\n", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal(err)
	}
}

// GET /categories
func listCategories(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

// POST /categories
func createCategory(w http.ResponseWriter, r *http.Request) {
	var cat Category
	if err := json.NewDecoder(r.Body).Decode(&cat); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	mu.Lock()
	cat.ID = nextID
	nextID++
	categories = append(categories, cat)
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cat)
}

// GET /categories/{id}
func getCategory(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for _, cat := range categories {
		if cat.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(cat)
			return
		}
	}

	http.Error(w, "Category not found", http.StatusNotFound)
}

// PUT /categories/{id}
func updateCategory(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updatedCat Category
	if err := json.NewDecoder(r.Body).Decode(&updatedCat); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for i, cat := range categories {
		if cat.ID == id {
			updatedCat.ID = id
			categories[i] = updatedCat
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedCat)
			return
		}
	}

	http.Error(w, "Category not found", http.StatusNotFound)
}

// DELETE /categories/{id}
func deleteCategory(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for i, cat := range categories {
		if cat.ID == id {
			categories = append(categories[:i], categories[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Category not found", http.StatusNotFound)
}
