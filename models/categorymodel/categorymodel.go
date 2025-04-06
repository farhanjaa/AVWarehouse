package categorymodel

import (
	"database/sql"
	"log"
	"os"

	"github.com/farhanjaa/AVWarehouse/config"
	"github.com/farhanjaa/AVWarehouse/entities"

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

func GetAll() []entities.Category {
	db := GetDBConnection()
	rows, err := db.Query(`SELECT * FROM categories`)
	if err != nil {
		log.Println("❌ Error mengambil data kategori:", err)
		return nil
	}
	defer rows.Close()

	var categories []entities.Category

	for rows.Next() {
		var category entities.Category
		if err := rows.Scan(&category.Id, &category.Name, &category.CreatedAt, &category.UpdatedAt); err != nil {
			log.Println("❌ Error membaca hasil query:", err)
			continue
		}
		categories = append(categories, category)
	}

	return categories
}

func Create(category entities.Category) bool {
	db := GetDBConnection()
	result, err := db.Exec(`
		INSERT INTO categories (name, created_at, updated_at)
		VALUES (?, ?, ?)`,

		category.Name, category.CreatedAt, category.UpdatedAt,
	)

	if err != nil {
		log.Println("❌ Error saat menambahkan kategori:", err)
		return false
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		log.Println("❌ Error mendapatkan ID terakhir:", err)
		return false
	}

	return lastInsertID > 0
}

func Detail(id int) entities.Category {
	db := GetDBConnection()
	row := db.QueryRow(`SELECT id, name FROM categories WHERE id = ?`, id)

	var category entities.Category
	if err := row.Scan(&category.Id, &category.Name); err != nil {
		log.Println("❌ Error mengambil detail kategori:", err)
		return entities.Category{} // Kembalikan struct kosong jika terjadi error
	}

	return category
}

func Update(id int, category entities.Category) bool {
	db := GetDBConnection()
	query, err := db.Exec(`
		UPDATE categories SET name = ?, updated_at = ? 
		WHERE id = ?`, category.Name, category.UpdatedAt, id)

	if err != nil {
		log.Println("❌ Error saat memperbarui kategori:", err)
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
	_, err := db.Exec(`DELETE FROM categories WHERE id = ?`, id)
	if err != nil {
		log.Println("❌ Error saat menghapus kategori:", err)
	}
	return err
}
