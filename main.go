package main

import (
	"log"
	"net/http"

	"github.com/farhanjaa/AVWarehouse/config"
	authcontroller "github.com/farhanjaa/AVWarehouse/controllers/authcontroller"
	"github.com/farhanjaa/AVWarehouse/controllers/categorycontroller"
	"github.com/farhanjaa/AVWarehouse/controllers/productcontroller"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	config.ConnectDB()

	// 0.Login Page
	http.HandleFunc("/", authcontroller.Index)
	http.HandleFunc("/login", authcontroller.Login)
	http.HandleFunc("/logout", authcontroller.Logout)
	http.HandleFunc("/register", authcontroller.Register)

	// 2. Categories
	http.HandleFunc("/categories", categorycontroller.Index)
	http.HandleFunc("/categories/add", categorycontroller.Add)
	http.HandleFunc("/categories/edit", categorycontroller.Edit)
	http.HandleFunc("/categories/delete", categorycontroller.Delete)

	// Products
	http.HandleFunc("/products", productcontroller.Index)
	http.HandleFunc("/products/add", productcontroller.Add)
	http.HandleFunc("/products/detail", productcontroller.Detail)
	http.HandleFunc("/products/edit", productcontroller.Edit)
	http.HandleFunc("/products/delete", productcontroller.Delete)

	log.Println("Server running on port 8080")
	http.ListenAndServe(":8080", nil)
}
