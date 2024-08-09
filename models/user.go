package models

import (
	"database/sql"
	"fmt"
	"os"
	"sync"
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

var (
	cacheMutex sync.Mutex
	userCache  = make(map[string]*User) // userCache caches User instances.
)

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

	removeFromCache(u.Email)

	return nil
}

func GetUserByEmail(email string) (*User, error) {
	if cachedUser, isCached := getFromCache(email); isCached {
		return cachedUser, nil
	}

	db, err := DbConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var user User
	query := `SELECT id, username, password, email, registration_date FROM users WHERE email = $1`
	row := db.QueryRow(query, email)
	if err = row.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.RegistrationDate); err != nil {
		return nil, err
	}

	addToCache(email, &user)

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

	addToCache(u.Email, u)

	return nil
}

func ListAllUsers() ([]User, error) {
	db, err := DbConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `SELECT id, username, password, email, registration_date FROM users`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.RegistrationDate); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func addToCache(email string, user *User) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	userCache[email] = user
}

func getFromCache(email string) (*User, bool) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	user, exists := userCache[email]
	return user, exists
}

func removeFromCache(email string) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	delete(userCache, email)
}