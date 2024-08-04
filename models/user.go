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
    // Declare a mutex to make the cache safe for concurrent use
    cacheMutex sync.Mutex
    // Create a simple in-memory cache
    userCache = make(map[string]*User)
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

    // Invalidate cache after creating a new user
    removeFromCache(u.Email)

    return nil
}

func GetUserByEmail(email string) (*User, error) {
    // Attempt to retrieve user from cache
    cachedUser, isCached := getFromCache(email)
    if isCached {
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
    err = row.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.RegistrationDate)
    if err != nil {
        return nil, err
    }

    // Add user to cache
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

    // Update cache after modifying the user
    addToCache(u.Email, u)

    return nil
}

// addToCache adds a user to the in-memory cache.
func addToCache(email string, user *User) {
    cacheMutex.Lock()
    defer cacheMutex.Unlock()
    userCache[email] = user
}

// getFromCache tries to get a user from the in-memory cache.
func getFromCache(email string) (*User, bool) {
    cacheMutex.Lock()
    defer cacheMutex.Unlock()
    user, exists := userCache[email]
    return user, exists
}

// removeFromCache removes a user from the cache.
func removeFromCache(email string) {
    cacheMutex.Lock()
    defer cacheMutex.Unlock()
    delete(userCache, email)
}