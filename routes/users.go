package main

import (
	"database/sql"
	"errors"
	"os"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/lib/pq"
)

var db *sql.DB

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

func initDB() {
	var err error
	dsn := os.Getenv("DATABASE_URL")
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * 60) // seconds

	err = db.Ping()
	if err != nil {
		panic(err)
	}
}

func registerUser(email, password string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	var user User
	err = db.QueryRow("INSERT INTO users(email, password) VALUES($1, $2) RETURNing id", email, string(hashedPassword)).Scan(&user.ID)
	if err != nil {
		return nil, err
	}

	user.Email = email
	return &source, nil
}

func main() {
	initDB()
	// Application initialization logic
}