package models

import (
	"database/sql"
	"errors"

	"github.com/farhanjaa/AVWarehouse/entities"
)

// UserModel adalah model untuk operasi pengguna
type UserModel struct {
	db *sql.DB
}

// NewUserModel membuat instance UserModel dengan koneksi database yang diberikan
func NewUserModel(db *sql.DB) (*UserModel, error) {
	if db == nil {
		return nil, errors.New("database connection is nil")
	}
	return &UserModel{db: db}, nil
}

// Where mencari user berdasarkan field tertentu
func (u UserModel) Where(user *entities.User, fieldName, fieldValue string) error {
	// Validasi field yang diizinkan
	allowedFields := map[string]bool{"id": true, "nama_lengkap": true, "email": true, "username": true}
	if !allowedFields[fieldName] {
		return errors.New("invalid field name")
	}

	// Query dengan prepared statement untuk mencegah SQL Injection
	query := "SELECT id, nama_lengkap, email, username, password FROM users WHERE " + fieldName + " = ? LIMIT 1"
	stmt, err := u.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Eksekusi query
	row := stmt.QueryRow(fieldValue)

	// Scan hasil query
	err = row.Scan(&user.Id, &user.NamaLengkap, &user.Email, &user.Username, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("user not found")
		}
		return err
	}

	return nil
}

// Create menyimpan user baru ke database
func (u UserModel) Create(user entities.User) (int64, error) {
	result, err := u.db.Exec("INSERT INTO users (nama_lengkap, email, username, password) VALUES(?,?,?,?)",
		user.NamaLengkap, user.Email, user.Username, user.Password)

	if err != nil {
		return 0, err
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastInsertId, nil
}
