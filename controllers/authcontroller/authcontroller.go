package authcontroller

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/farhanjaa/AVWarehouse/config"
	"github.com/farhanjaa/AVWarehouse/entities"
	"github.com/farhanjaa/AVWarehouse/libraries"
	"github.com/farhanjaa/AVWarehouse/models"

	"golang.org/x/crypto/bcrypt"
)

type UserInput struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
}

var userModel = models.NewUserModel()
var validation = libraries.NewValidation()

// **Helper function untuk mengecek keberadaan file template**
func templateExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return !os.IsNotExist(err)
}

// **Helper function untuk render template secara aman**
func renderTemplate(w http.ResponseWriter, filepath string, data interface{}) {
	if !templateExists(filepath) {
		log.Printf("Template file not found: %s\n", filepath)
		http.Error(w, "Template file not found", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles(filepath)
	if err != nil {
		log.Printf("Error parsing template: %s, Error: %v\n", filepath, err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error executing template: %s, Error: %v\n", filepath, err)
		http.Error(w, "Error displaying page", http.StatusInternalServerError)
	}
}

// **Index (Homepage)**
func Index(w http.ResponseWriter, r *http.Request) {
	session, err := config.Store.Get(r, config.SESSION_ID)
	if err != nil {
		log.Println("Error getting session:", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Pastikan sesi valid dan pengguna login
	if loggedIn, ok := session.Values["loggedIn"].(bool); !ok || !loggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data := map[string]interface{}{
		"nama_lengkap": session.Values["nama_lengkap"],
	}

	renderTemplate(w, "views/home/index.html", data)
}

// **Login**
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		renderTemplate(w, "views/login/login.html", nil)
		return
	}

	// **POST - proses login**
	if err := r.ParseForm(); err != nil {
		log.Println("Error parsing form:", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	userInput := &UserInput{
		Username: r.Form.Get("username"),
		Password: r.Form.Get("password"),
	}

	// Validasi input
	if errMessages := validation.Struct(userInput); errMessages != nil {
		renderTemplate(w, "views/login/login.html", map[string]interface{}{"validation": errMessages})
		return
	}

	// Cek apakah user ada di database
	var user entities.User
	err := userModel.Where(&user, "username", userInput.Username)
	if err != nil {
		renderTemplate(w, "views/login/login.html", map[string]interface{}{"error": "Username atau password salah!"})
		return
	}

	// **Verifikasi password**
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password)); err != nil {
		renderTemplate(w, "views/login/login.html", map[string]interface{}{"error": "Username atau password salah!"})
		return
	}

	// **Set session dengan nilai yang aman**
	session, err := config.Store.Get(r, config.SESSION_ID)
	if err != nil {
		log.Println("Error creating session:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	session.Values["loggedIn"] = true
	session.Values["email"] = user.Email
	session.Values["username"] = user.Username
	session.Values["nama_lengkap"] = user.NamaLengkap

	if err := session.Save(r, w); err != nil {
		log.Println("Error saving session:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// **Logout**
func Logout(w http.ResponseWriter, r *http.Request) {
	session, err := config.Store.Get(r, config.SESSION_ID)
	if err != nil {
		log.Println("Error getting session:", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// **Destroy session dengan benar**
	session.Values = map[interface{}]interface{}{}
	session.Options.MaxAge = -1

	if err := session.Save(r, w); err != nil {
		log.Println("Error destroying session:", err)
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// **Register**
func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		renderTemplate(w, "views/login/register.html", nil)
		return
	}

	// **POST - proses registrasi**
	if err := r.ParseForm(); err != nil {
		log.Println("Error parsing form:", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user := entities.User{
		NamaLengkap: r.Form.Get("nama_lengkap"),
		Email:       r.Form.Get("email"),
		Username:    r.Form.Get("username"),
		Password:    r.Form.Get("password"),
		Cpassword:   r.Form.Get("cpassword"),
	}

	// **Validasi input user**
	if errMessages := validation.Struct(user); errMessages != nil {
		renderTemplate(w, "views/login/register.html", map[string]interface{}{"validation": errMessages, "user": user})
		return
	}

	// **Cek apakah username/email sudah ada**
	var existingUser entities.User
	if err := userModel.Where(&existingUser, "username", user.Username); err == nil && existingUser.Username != "" {
		renderTemplate(w, "views/login/register.html", map[string]interface{}{"error": "Username sudah digunakan!", "user": user})
		return
	}

	if err := userModel.Where(&existingUser, "email", user.Email); err == nil && existingUser.Email != "" {
		renderTemplate(w, "views/login/register.html", map[string]interface{}{"error": "Email sudah digunakan!", "user": user})
		return
	}

	// **Hash password dengan aman**
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashPassword)

	// **Simpan user ke database & dapatkan lastInsertId**
	lastInsertId, err := userModel.Create(user)
	if err != nil {
		log.Println("Error inserting user:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	log.Printf("User registered successfully with ID: %d", lastInsertId)

	renderTemplate(w, "views/login/register.html", map[string]interface{}{
		"pesan":  "Registrasi berhasil",
		"lastId": lastInsertId, // Menampilkan ID user yang baru terdaftar
	})
}
