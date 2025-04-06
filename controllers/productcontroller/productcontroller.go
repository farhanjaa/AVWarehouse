package productcontroller

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/farhanjaa/AVWarehouse/config"
	"github.com/farhanjaa/AVWarehouse/entities"
	"github.com/farhanjaa/AVWarehouse/models/categorymodel"
	"github.com/farhanjaa/AVWarehouse/models/productmodel"
)

// Middleware untuk memeriksa sesi login
func checkSession(w http.ResponseWriter, r *http.Request) bool {
	session, err := config.Store.Get(r, config.SESSION_ID)
	if err != nil {
		log.Println("❌ Error mendapatkan sesi:", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return false
	}

	// Cek apakah pengguna sudah login
	if loggedIn, ok := session.Values["loggedIn"].(bool); !ok || !loggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return false
	}

	return true
}

func Index(w http.ResponseWriter, r *http.Request) {
	if !checkSession(w, r) {
		return
	}

	products := productmodel.GetAll()
	data := map[string]any{
		"products": products,
	}

	temp, err := template.ParseFiles("product/index.html")
	if err != nil {
		log.Println("❌ Error parsing template:", err)
		http.Error(w, "Error saat memuat halaman", http.StatusInternalServerError)
		return
	}

	temp.Execute(w, data)
}

func Detail(w http.ResponseWriter, r *http.Request) {
	if !checkSession(w, r) {
		return
	}

	idString := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	product := productmodel.Detail(id)
	data := map[string]any{
		"product": product,
	}

	temp, err := template.ParseFiles("product/detail.html")
	if err != nil {
		log.Println("❌ Error parsing template:", err)
		http.Error(w, "Error saat memuat halaman", http.StatusInternalServerError)
		return
	}

	temp.Execute(w, data)
}

func Add(w http.ResponseWriter, r *http.Request) {
	if !checkSession(w, r) {
		return
	}

	if r.Method == "GET" {
		temp, err := template.ParseFiles("product/create.html")
		if err != nil {
			log.Println("❌ Error parsing template:", err)
			http.Error(w, "Error saat memuat halaman", http.StatusInternalServerError)
			return
		}

		categories := categorymodel.GetAll()
		data := map[string]any{
			"categories": categories,
		}

		temp.Execute(w, data)
	}

	if r.Method == "POST" {
		var product entities.Product

		categoryId, err := strconv.Atoi(r.FormValue("category_id"))
		if err != nil {
			http.Error(w, "ID kategori tidak valid", http.StatusBadRequest)
			return
		}

		stock, err := strconv.Atoi(r.FormValue("stock"))
		if err != nil {
			http.Error(w, "Stock tidak valid", http.StatusBadRequest)
			return
		}

		product.Name = r.FormValue("name")
		product.Category.Id = uint(categoryId)
		product.Stock = int64(stock)
		product.Description = r.FormValue("description")
		product.CreatedAt = time.Now()
		product.UpdateAt = time.Now()

		if ok := productmodel.Create(product); !ok {
			http.Redirect(w, r, r.Header.Get("Referer"), http.StatusTemporaryRedirect)
			return
		}

		http.Redirect(w, r, "/products", http.StatusSeeOther)
	}
}

func Edit(w http.ResponseWriter, r *http.Request) {
	if !checkSession(w, r) {
		return
	}

	if r.Method == "GET" {
		temp, err := template.ParseFiles("product/edit.html")
		if err != nil {
			log.Println("❌ Error parsing template:", err)
			http.Error(w, "Error saat memuat halaman", http.StatusInternalServerError)
			return
		}

		idString := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			http.Error(w, "ID tidak valid", http.StatusBadRequest)
			return
		}

		product := productmodel.Detail(id)
		categories := categorymodel.GetAll()
		data := map[string]any{
			"categories": categories,
			"product":    product,
		}
		temp.Execute(w, data)
	}

	if r.Method == "POST" {
		var product entities.Product

		idString := r.FormValue("id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			http.Error(w, "ID tidak valid", http.StatusBadRequest)
			return
		}

		categoryId, err := strconv.Atoi(r.FormValue("category_id"))
		if err != nil {
			http.Error(w, "ID kategori tidak valid", http.StatusBadRequest)
			return
		}

		stock, err := strconv.Atoi(r.FormValue("stock"))
		if err != nil {
			http.Error(w, "Stock tidak valid", http.StatusBadRequest)
			return
		}

		product.Name = r.FormValue("name")
		product.Category.Id = uint(categoryId)
		product.Stock = int64(stock)
		product.Description = r.FormValue("description")
		product.CreatedAt = time.Now()
		product.UpdateAt = time.Now()

		if ok := productmodel.Update(id, product); !ok {
			http.Redirect(w, r, r.Header.Get("Referer"), http.StatusTemporaryRedirect)
			return
		}

		http.Redirect(w, r, "/products", http.StatusSeeOther)
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	if !checkSession(w, r) {
		return
	}

	idString := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	if err := productmodel.Delete(id); err != nil {
		log.Println("❌ Error menghapus produk:", err)
		http.Error(w, "Gagal menghapus produk", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/products", http.StatusSeeOther)
}
