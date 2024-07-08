package main

import (
	"database/sql"
	"errors"
	"os"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/lib/pq"
)

var (
	db *sql.DB
)

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

	err = db.QueryRow("INSERT INTO users(email, password) VALUES($1, $2) RETURNING id", email, string(hashedPassword)).Scan(&user.ID)
	if err != nil {
		return nil, err
	}

	user.Email = email
	return &user, nil
}

func getUserByID(userID int) (*User, error) {
	var user User

	row := db.QueryRow("SELECT id, email FROM users WHERE id = $1", userID)
	err := row.Scan(&user.ID, &user.Email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func updateUser(userID int, email, password string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("UPDATE users SET email = $1, password = $2 WHERE id = $3", email, hashedPassword, userID)
	if err != nil {
		return nil, err
	}

	return getUserByID(userID)
}

func authenticateUser(email, password string) (bool, *User, error) {
	var user User
	var hashedPassword string

	row := db.QueryRow("SELECT id, email, password FROM users WHERE email = $1", email)
	err := row.Scan(&user.ID, &user.Email, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil, errors.New("user not found")
		}
		return false, nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false, &user, errors.New("authentication failed")
	}

	return true, &user, nil
}

func main() {
	initDB()
}