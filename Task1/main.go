// ### 🔗 Endpoint yang Wajib Ada

// * **GET** `/categories` → Ambil semua kategori
// * **POST** `/categories` → Tambah kategori
// * **PUT** `/categories/{id}` → Update kategori
// * **GET** `/categories/{id}` → Ambil detail satu kategori
// * **DELETE** `/categories/{id}` → Hapus kategori

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	port string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock  int    `json:"stock"`
}

var categories = []Category{
	{ID: 1, Name: "Food", Description: "Kategori makanan"},
	{ID: 2, Name: "Beverage", Description: "Kategori minuman"},
	{ID: 3, Name: "Personal care", Description: "Kategori perawatan diri"},
}

var products = []Product{
	{ID: 1, Name: "Indomie", Price: 5000, Stock: 87},
	{ID: 2, Name: "Sabun", Price: 7000, Stock: 36},
	{ID: 3, Name: "Yakult", Price: 3000, Stock: 100},
}

// GET /api/categories
func getCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

// GET /api/products
func getProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// POST /api/categories
func createCategory(w http.ResponseWriter, r *http.Request) {
	var newCategory Category
	// fmt.Println(r.Body)
	err := json.NewDecoder(r.Body).Decode(&newCategory)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// masukkan data ke dalam variable category baru
	newCategory.ID = len(categories) + 1
	categories = append(categories, newCategory)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201
	json.NewEncoder(w).Encode(newCategory)
}

// POST /api/products
func createProduct(w http.ResponseWriter, r *http.Request) {
	var newProduct Product
	// fmt.Println(r.Body)
	err := json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// masukkan data ke dalam variable category baru
	newProduct.ID = len(products) + 1
	products = append(products, newProduct)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201
	json.NewEncoder(w).Encode(newProduct)
}

// GET /api/categories/{id}
func getCategoryByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid category ID", http.StatusBadRequest)
		return
	}

	for _, p := range categories {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}

	http.Error(w, "category not found", http.StatusNotFound)
}

// GET /api/products/{id}
func getProductByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid product ID", http.StatusBadRequest)
		return
	}

	for _, p := range products {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}

	http.Error(w, "product not found", http.StatusNotFound)
}

// PUT /api/categories/{id}
func updateCategoryByID(w http.ResponseWriter, r *http.Request) {
	// get id dari request
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")

	// convert jadi int pake strconv.Atoi(id)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid category ID", http.StatusBadRequest)
		return
	}

	// get data dari request
	var updateCategory Category
	err = json.NewDecoder(r.Body).Decode(&updateCategory)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// loop category, cari id, ganti sesuai data dari request
	for i := range categories {
		if categories[i].ID == id {
			updateCategory.ID = id
			categories[i] = updateCategory

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updateCategory)
			return
		}
	}

	http.Error(w, "category belum ada", http.StatusNotFound)
}

// PUT /api/products/{id}
func updateProductByID(w http.ResponseWriter, r *http.Request) {
	// get id dari request
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")

	// convert jadi int pake strconv.Atoi(id)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid product ID", http.StatusBadRequest)
		return
	}

	// get data dari request
	var updateProduct Product
	err = json.NewDecoder(r.Body).Decode(&updateProduct)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// loop product, cari id, ganti sesuai data dari request
	for i := range products {
		if products[i].ID == id {
			updateProduct.ID = id
			products[i] = updateProduct

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updateProduct)
			return
		}
	}

	http.Error(w, "product belum ada", http.StatusNotFound)
}

// DELETE /api/categories/{id}
func deleteCategoryById(w http.ResponseWriter, r *http.Request) {
	// get id
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")

	// convert jadi int pake strconv.Atoi(id)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid category ID", http.StatusBadRequest)
		return
	}

	// loop category, cari ID dan index yang mau dihapus
	for i, p := range categories {
		if p.ID == id {
			// buat slice baru dengan data sebelum dan sesudah index
			categories = append(categories[:i], categories[i+1:]...)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "sukses delete",
			})
			return
		}
	}

	http.Error(w, "category belum ada", http.StatusNotFound)
}

// DELETE /api/products/{id}
func deleteProductById(w http.ResponseWriter, r *http.Request) {
	// get id
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")

	// convert jadi int pake strconv.Atoi(id)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid product ID", http.StatusBadRequest)
		return
	}

	// loop product, cari ID dan index yang mau dihapus
	for i, p := range products {
		if p.ID == id {
			// buat slice baru dengan data sebelum dan sesudah index
			products = append(products[:i], products[i+1:]...)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "sukses delete",
			})
			return
		}
	}

	http.Error(w, "product belum ada", http.StatusNotFound)
}

/*
	Main Function
*/
func main() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	_, error := os.Stat(".env")
	if error == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		port : viper.GetString("PORT"),
		DBConn : viper.GetString("DB_CONN"),
	}
	// cek server /health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "200 - success",
			"message": "server is running normally",
		})
	})

	// GET /api/categories
	// POST /api/categories
	http.HandleFunc("/api/categories", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getCategories(w, r)
		} else if r.Method == "POST" {
			createCategory(w, r)
		}
	})

	// GET /api/categories/{id}
	// PUT /api/categories/{id}
	// DELETE /api/categories/{id}
	http.HandleFunc("/api/categories/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getCategoryByID(w, r)
		} else if r.Method == "PUT" {
			updateCategoryByID(w, r)
		} else if r.Method == "DELETE" {
			deleteCategoryById(w, r)
		}
	})

	// GET /api/products
	// POST /api/products
	http.HandleFunc("/api/products", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getProducts(w, r)
		} else if r.Method == "POST" {
			createProduct(w, r)
		}
	})

	// GET /api/products/{id}
	// PUT /api/products/{id}
	// DELETE /api/products/{id}
	http.HandleFunc("/api/products/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getProductByID(w, r)
		} else if r.Method == "PUT" {
			updateProductByID(w, r)
		} else if r.Method == "DELETE" {
			deleteProductById(w, r)
		}
	})

	// PORT := ":8080"
	fmt.Println("server running di https://localhost:"+config.port)

	err := http.ListenAndServe(":"+config.port, nil)
	if err != nil {
		fmt.Println("Error running server")
	}
}
