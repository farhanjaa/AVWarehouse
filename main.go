package main

import (
	"log"
	"net/http"
	"os"

	"github.com/farhanjaa/AVWarehouse/config"
	authcontroller "github.com/farhanjaa/AVWarehouse/controllers/authcontroller"
	"github.com/farhanjaa/AVWarehouse/controllers/categorycontroller"
	"github.com/farhanjaa/AVWarehouse/controllers/productcontroller"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Menghubungkan ke database tanpa mengirimkan connection string secara eksplisit
	_, err := config.ConnectDB() // Sekarang tidak perlu mengirim parameter
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	// Mengambil port dari environment variable yang disediakan oleh Railway
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default ke 8080 jika tidak ditemukan
	}

	// Menyiapkan route HTTP
	http.HandleFunc("/", authcontroller.Index)
	http.HandleFunc("/login", authcontroller.Login)
	http.HandleFunc("/logout", authcontroller.Logout)
	http.HandleFunc("/register", authcontroller.Register)

	// Routes untuk kategori
	http.HandleFunc("/categories", categorycontroller.Index)
	http.HandleFunc("/categories/add", categorycontroller.Add)
	http.HandleFunc("/categories/edit", categorycontroller.Edit)
	http.HandleFunc("/categories/delete", categorycontroller.Delete)

	// Routes untuk produk
	http.HandleFunc("/products", productcontroller.Index)
	http.HandleFunc("/products/add", productcontroller.Add)
	http.HandleFunc("/products/detail", productcontroller.Detail)
	http.HandleFunc("/products/edit", productcontroller.Edit)
	http.HandleFunc("/products/delete", productcontroller.Delete)

	// Menjalankan server di port yang telah ditentukan
	log.Println("Server running on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
