package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// DB is a global variable to hold the database connection
var DB *sql.DB

// ConnectDB is the function to connect to the database
func ConnectDB() (*sql.DB, error) {
	// Ambil connection string dari environment variable yang disediakan Railway
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	// Membuka koneksi ke database menggunakan connection string yang didapat dari Railway
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Cek koneksi ke database
	if err = db.Ping(); err != nil {
		return nil, err
	}

	fmt.Println("✅ Database Connected!")
	DB = db // Assign the connection to the global variable
	return db, nil
}
