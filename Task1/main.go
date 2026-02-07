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
	port   string `mapstructure:"PORT"`
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
		port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	// log.Println("PORT:", config.port)
	// log.Println("DB_CONN:", config.DBConn)

	// koneksi ke database
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// cek server /health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "200 - success",
			"message": "server is running normally",
		})
	})

	// Dependency injection
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	reportRepo := repositories.NewReportRepository(db)
	reportService := services.NewReportService(reportRepo)
	reportHandler := handlers.NewReportHandler(reportService)

	/* Products */
	http.HandleFunc("/api/product", productHandler.HandleProducts)
	http.HandleFunc("/api/product/{id}", productHandler.HandleProductByID)
	// "/api/product/{id}" -> Bisa pake idStr := r.PathValue("id")

	/* Category */
	http.HandleFunc("/api/category", categoryHandler.HandleCategory)
	http.HandleFunc("/api/category/", categoryHandler.HandleCategoryByID)
	// ""/api/category/" -> Bisa jg pake idStr := strings.TrimPrefix(r.URL.Path, "/api/category/")

	/* Transaction */
	http.HandleFunc("/api/checkout", transactionHandler.HandleCheckout)

	/* Report */
	http.HandleFunc("/api/report", reportHandler.GetReport)
	http.HandleFunc("/api/report/hari-ini", reportHandler.GetReport)
	http.HandleFunc("/api/report/today", reportHandler.GetReport)

	// PORT := ":8080"
	if config.port != "80" {
		fmt.Println("server running di http://localhost:" + config.port)
	} else {
		fmt.Println("server running di https://localhost:" + config.port)
	}

	err = http.ListenAndServe(":"+config.port, nil)
	if err != nil {
		fmt.Println("Error running server")
	}
}
