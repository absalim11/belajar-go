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

// Produk represents the product model
type Produk struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Harga int    `json:"harga"`
	Stok  int    `json:"stok"`
}

// In-memory storage
var (
	categories = []Category{}
	nextID     = 1
	mu         sync.Mutex

	produk = []Produk{
		{ID: 1, Nama: "Indomie Godog", Harga: 3500, Stok: 10},
		{ID: 2, Nama: "Vit 1000ml", Harga: 3000, Stok: 40},
		{ID: 3, Nama: "kecap", Harga: 12000, Stok: 20},
	}
	nextProdID = 4
)

func main() {
	mux := http.NewServeMux()

	// Category Routes
	mux.HandleFunc("GET /categories", listCategories)
	mux.HandleFunc("POST /categories", createCategory)
	mux.HandleFunc("GET /categories/{id}", getCategory)
	mux.HandleFunc("PUT /categories/{id}", updateCategory)
	mux.HandleFunc("DELETE /categories/{id}", deleteCategory)

	// Product Routes
	mux.HandleFunc("GET /products", listProducts)
	mux.HandleFunc("POST /products", createProduct)
	mux.HandleFunc("GET /products/{id}", getProduct)
	mux.HandleFunc("PUT /products/{id}", updateProduct)
	mux.HandleFunc("DELETE /products/{id}", deleteProduct)

	port := ":8080"
	fmt.Printf("Server starting on http://localhost%s\n", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal(err)
	}
}

// --- Category Handlers ---

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

// --- Product Handlers ---

// GET /products
func listProducts(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(produk)
}

// POST /products
func createProduct(w http.ResponseWriter, r *http.Request) {
	var p Produk
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	mu.Lock()
	p.ID = nextProdID
	nextProdID++
	produk = append(produk, p)
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

// GET /products/{id}
func getProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for _, p := range produk {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}

	http.Error(w, "Product not found", http.StatusNotFound)
}

// PUT /products/{id}
func updateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updatedP Produk
	if err := json.NewDecoder(r.Body).Decode(&updatedP); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for i, p := range produk {
		if p.ID == id {
			updatedP.ID = id
			produk[i] = updatedP
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedP)
			return
		}
	}

	http.Error(w, "Product not found", http.StatusNotFound)
}

// DELETE /products/{id}
func deleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for i, p := range produk {
		if p.ID == id {
			produk = append(produk[:i], produk[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Product not found", http.StatusNotFound)
}
