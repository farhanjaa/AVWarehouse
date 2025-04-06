package categorycontroller

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/farhanjaa/AVWarehouse/config"
	"github.com/farhanjaa/AVWarehouse/entities"
	"github.com/farhanjaa/AVWarehouse/models/categorymodel"
)

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

// Fungsi untuk validasi keberadaan file template
func validateTemplateFile(filename string) error {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return err
	}
	return nil
}

func Index(w http.ResponseWriter, r *http.Request) {
	if !checkSession(w, r) {
		return
	}

	categories := categorymodel.GetAll()
	data := map[string]any{
		"categories": categories,
	}

	templateFile := "views/category/index.html"
	if err := validateTemplateFile(templateFile); err != nil {
		http.Error(w, "Template file not found", http.StatusInternalServerError)
		log.Println("❌ Missing template file:", templateFile)
		return
	}

	temp, err := template.ParseFiles(templateFile)
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		log.Println("❌ Template parsing error:", err)
		return
	}

	temp.Execute(w, data)
}

func Add(w http.ResponseWriter, r *http.Request) {
	if !checkSession(w, r) {
		return
	}

	templateFile := "views/category/create.html"
	if err := validateTemplateFile(templateFile); err != nil {
		http.Error(w, "Template file not found", http.StatusInternalServerError)
		log.Println("❌ Missing template file:", templateFile)
		return
	}

	if r.Method == http.MethodGet {
		temp, err := template.ParseFiles(templateFile)
		if err != nil {
			http.Error(w, "Error parsing template", http.StatusInternalServerError)
			log.Println("❌ Template parsing error:", err)
			return
		}
		temp.Execute(w, nil)
		return
	}

	if r.Method == http.MethodPost {
		category := entities.Category{
			Name:      r.FormValue("name"),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if ok := categorymodel.Create(category); !ok {
			log.Println("❌ Failed to create category")
			http.Error(w, "Failed to create category", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/categories", http.StatusSeeOther)
	}
}

func Edit(w http.ResponseWriter, r *http.Request) {
	if !checkSession(w, r) {
		return
	}

	templateFile := "views/category/edit.html"
	if err := validateTemplateFile(templateFile); err != nil {
		http.Error(w, "Template file not found", http.StatusInternalServerError)
		log.Println("❌ Missing template file:", templateFile)
		return
	}

	if r.Method == http.MethodGet {
		idString := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idString)
		if err != nil || id <= 0 {
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			log.Println("❌ Invalid category ID:", idString)
			return
		}

		category := categorymodel.Detail(id)
		if category.Id == 0 {
			http.Error(w, "Category not found", http.StatusNotFound)
			log.Println("❌ Category not found:", id)
			return
		}

		data := map[string]any{
			"category": category,
		}

		temp, err := template.ParseFiles(templateFile)
		if err != nil {
			http.Error(w, "Error parsing template", http.StatusInternalServerError)
			log.Println("❌ Template parsing error:", err)
			return
		}

		temp.Execute(w, data)
		return
	}

	if r.Method == http.MethodPost {
		idString := r.FormValue("id")
		id, err := strconv.Atoi(idString)
		if err != nil || id <= 0 {
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			log.Println("❌ Invalid category ID:", idString)
			return
		}

		category := entities.Category{
			Name:      r.FormValue("name"),
			UpdatedAt: time.Now(),
		}

		if ok := categorymodel.Update(id, category); !ok {
			log.Println("❌ Failed to update category")
			http.Error(w, "Failed to update category", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/categories", http.StatusSeeOther)
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	if !checkSession(w, r) {
		return
	}

	idString := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idString)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		log.Println("❌ Invalid category ID:", idString)
		return
	}

	if err := categorymodel.Delete(id); err != nil {
		log.Println("❌ Failed to delete category:", err)
		http.Error(w, "Failed to delete category", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/categories", http.StatusFound)
}
