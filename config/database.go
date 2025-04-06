package config

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// DB is a global variable to hold the database connection
var DB *sql.DB

// ConnectDB is the function to connect to the database
func ConnectDB() (*sql.DB, error) {
	// Get connection string from the environment variable
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, errors.New("DATABASE_URL is not set")
	}

	// Parse the DATABASE_URL to convert it into a DSN
	parsedURL, err := url.Parse(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DATABASE_URL: %w", err)
	}

	// Extract username and password
	user := parsedURL.User.Username()
	password, _ := parsedURL.User.Password()

	// Extract host and port
	host := parsedURL.Host
	if !strings.Contains(host, ":") {
		host += ":3306" // Default MySQL port
	}

	// Extract database name
	dbName := strings.TrimPrefix(parsedURL.Path, "/")

	// Build the DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", user, password, host, dbName)

	// Open connection to the database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Check the connection
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Println("✅ Database Connected!")
	DB = db // Assign the connection to the global variable
	return db, nil
}
