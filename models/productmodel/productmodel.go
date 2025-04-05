package productmodel

import (
	"1stproject/config"
	"1stproject/entities"
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// GetDBConnection ensures the DB connection is available
func GetDBConnection() *sql.DB {
	if config.DB == nil {
		// Load .env file
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		dbUser := os.Getenv("DB_USER")
		dbPass := os.Getenv("DB_PASS")
		dbHost := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")
		dbName := os.Getenv("DB_NAME")

		dsn := dbUser + ":" + dbPass + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?parseTime=true"
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			log.Fatal("❌ Failed to connect to database:", err)
		}
		config.DB = db
	}
	return config.DB
}

func GetAll() []entities.Product {
	db := GetDBConnection()
	rows, err := db.Query(`
		SELECT
			products.id,
			products.name,
			categories.name as category_name,
			products.stock,
			products.description,
			products.created_at,
			products.updated_at
		FROM products
		JOIN categories ON products.category_id = categories.id
	`)

	if err != nil {
		log.Println("❌ Error mengambil data produk:", err)
		return nil
	}
	defer rows.Close()

	var products []entities.Product

	for rows.Next() {
		var product entities.Product
		err := rows.Scan(
			&product.Id,
			&product.Name,
			&product.Category.Name,
			&product.Stock,
			&product.Description,
			&product.CreatedAt,
			&product.UpdateAt,
		)

		if err != nil {
			log.Println("❌ Error membaca hasil query:", err)
			continue
		}

		products = append(products, product)
	}

	return products
}

func Create(product entities.Product) bool {
	db := GetDBConnection()
	result, err := db.Exec(`
		INSERT INTO products(
			name, category_id, stock, description, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?)`,

		product.Name,
		product.Category.Id,
		product.Stock,
		product.Description,
		product.CreatedAt,
		product.UpdateAt,
	)

	if err != nil {
		log.Println("❌ Error saat menambahkan produk:", err)
		return false
	}

	LastInsertId, err := result.LastInsertId()
	if err != nil {
		log.Println("❌ Error mendapatkan ID terakhir:", err)
		return false
	}

	return LastInsertId > 0
}

func Detail(id int) entities.Product {
	db := GetDBConnection()
	row := db.QueryRow(`
		SELECT
			products.id,
			products.name,
			categories.name as category_name,
			products.stock,
			products.description,
			products.created_at,
			products.updated_at
		FROM products
		JOIN categories ON products.category_id = categories.id
		WHERE products.id = ?`, id)

	var product entities.Product
	err := row.Scan(
		&product.Id,
		&product.Name,
		&product.Category.Name,
		&product.Stock,
		&product.Description,
		&product.CreatedAt,
		&product.UpdateAt,
	)

	if err != nil {
		log.Println("❌ Produk tidak ditemukan atau terjadi error:", err)
		return entities.Product{} // Kembalikan struct kosong jika error
	}

	return product
}

func Update(id int, product entities.Product) bool {
	db := GetDBConnection()
	query, err := db.Exec(`
		UPDATE products SET 
			name = ?, 
			category_id = ?,
			stock = ?,
			description = ?,
			updated_at = ?
		WHERE id = ?`,
		product.Name,
		product.Category.Id,
		product.Stock,
		product.Description,
		product.UpdateAt,
		id,
	)

	if err != nil {
		log.Println("❌ Error saat memperbarui produk:", err)
		return false
	}

	result, err := query.RowsAffected()
	if err != nil {
		log.Println("❌ Error mendapatkan jumlah baris yang terpengaruh:", err)
		return false
	}

	return result > 0
}

func Delete(id int) error {
	db := GetDBConnection()
	_, err := db.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		log.Println("❌ Error saat menghapus produk:", err)
	}
	return err
}
