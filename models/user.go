package models

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

type User struct {
	ID               int
	Username         string
	Password         string
	Email            string
	RegistrationDate time.Time
}

func DbConnection() (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME")))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (u *User) CreateUser() error {
	db, err := DbConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	query := `INSERT INTO users (username, password, email, registration_date) VALUES ($1, $2, $3, $4)`
	_, err = db.Exec(query, u.Username, u.Password, u.Email, u.RegistrationDate)
	if err != nil {
		return err
	}
	return nil
}

func GetUserByEmail(email string) (*User, error) {
	db, err := DbConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var user User
	query := `SELECT id, username, password, email, registration_date FROM users WHERE email = $1`
	row := db.QueryRow(query, email)
	err = row.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.RegistrationDate)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *User) UpdateUser() error {
	db, err := DbConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	query := `UPDATE users SET username = $2, password = $3, email = $4, registration_date = $5 WHERE id = $1`
	_, err = db.Exec(query, u.ID, u.Username, u.Password, u.Email, u.RegistrationDate)
	if err != nil {
		return err
	}
	return nil
}