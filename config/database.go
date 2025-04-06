package config

import (
	"database/sql"
	"errors"
	"fmt"
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
		return nil, errors.New("DATABASE_URL is not set")
	}

	// Membuka koneksi ke database menggunakan connection string yang didapat dari Railway
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Cek koneksi ke database
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Println("✅ Database Connected!")
	DB = db // Assign the connection to the global variable
	return db, nil
}
