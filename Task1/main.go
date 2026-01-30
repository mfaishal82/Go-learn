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
	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/repositories"
	"kasir-api/services"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	port string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

/*
	Main Function
*/
func main() {
	// konfigurasi viper
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	_, err := os.Stat(".env")
	if err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		port : viper.GetString("PORT"),
		DBConn : viper.GetString("DB_CONN"),
	}

	// log.Println("PORT:", config.port)
	// log.Println("DB_CONN:", config.DBConn)

	// koneksi ke database
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Dependency injection
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// cek server /health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "200 - success",
			"message": "server is running normally",
		})
	})

	// Routes products
	http.HandleFunc("/api/product", productHandler.HandleProducts)
	http.HandleFunc("/api/product/", productHandler.HandleProductByID)
	http.HandleFunc("/api/category", categoryHandler.HandleCategory)
	http.HandleFunc("/api/category/", categoryHandler.HandleCategoryByID)

	// PORT := ":8080"
	fmt.Println("server running di https://localhost:"+config.port)

	err = http.ListenAndServe(":"+config.port, nil)
	if err != nil {
		fmt.Println("Error running server")
	}
}
