package models

import (
	"1stproject/config"
	"1stproject/entities"
	"database/sql"
	"errors"
)

type UserModel struct {
	db *sql.DB
}

func NewUserModel() *UserModel {
	conn, err := config.ConnectDB()

	if err != nil {
		panic(err)
	}

	return &UserModel{
		db: conn,
	}
}

func (u UserModel) Where(user *entities.User, fieldName, fieldValue string) error {
	// 🔒 Hanya izinkan field tertentu untuk menghindari SQL Injection
	allowedFields := map[string]bool{"id": true, "nama_lengkap": true, "email": true, "username": true}
	if !allowedFields[fieldName] {
		return errors.New("invalid field name")
	}

	// 🔒 Gunakan prepared statement untuk mencegah SQL Injection
	query := "SELECT id, nama_lengkap, email, username, password FROM users WHERE " + fieldName + " = ? LIMIT 1"
	stmt, err := u.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// 🔍 Eksekusi query
	row := stmt.QueryRow(fieldValue)

	// 🔄 Scan hasil query
	err = row.Scan(&user.Id, &user.NamaLengkap, &user.Email, &user.Username, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("user not found")
		}
		return err
	}

	return nil
}

func (u UserModel) Create(user entities.User) (int64, error) {
	result, err := u.db.Exec("insert into users (nama_lengkap, email, username, password) values(?,?,?,?)",
		user.NamaLengkap, user.Email, user.Username, user.Password)

	if err != nil {
		return 0, err
	}

	lastInsertId, _ := result.LastInsertId()

	return lastInsertId, nil
}
