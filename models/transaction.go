package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Transaction struct {
	ID          int
	UserID      int
	Amount      float64
	Date        time.Time
	Category    string
	Description string
}

type Database struct {
	conn *sql.DB
}

func NewDatabase() *Database {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME")))
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	return &Database{conn: db}
}

func (db *Database) CreateTransaction(t *Transaction) error {
	query := `INSERT INTO transactions (user_id, amount, date, category, description) VALUES ($1, $2, $3, $4, $5)`
	_, err := db.conn.Exec(query, t.UserID, t.Amount, t.Date, t.Category, t.Description)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) GetTransactions(userID int) ([]Transaction, error) {
	transactions := []Transaction{}
	query := `SELECT id, user_id, amount, date, category, description FROM transactions WHERE user_id = $1`
	rows, err := db.conn.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t Transaction
		err := rows.Scan(&t.ID, &t.UserID, &t.Amount, &t.Date, &t.Category, &t.Description)
		if err != nil {
		return nil, err
	}
		transactions = append(transactions, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (db *Database) UpdateTransaction(t *Transaction) error {
	query := `UPDATE transactions SET amount = $2, date = $3, category = $4, description = $5 WHERE id = $1`
	_, err := db.conn.Exec(query, t.ID, t.Amount, t.Date, t.Category, t.Description)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) DeleteTransaction(id int) error {
	query := `DELETE FROM transactions WHERE id = $1`
	_, err := db.conn.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	db := NewProductDatabase()
	defer db.conn.Close()

	tx := &Transaction{
		UserID:      1,
		Amount:      199.99,
		Date:        time.Now(),
		Category:    "Electronics",
		Description: "New headphones",
	}
	err := db.CreateTransaction(tx)
	if err != nil {
		log.Fatalf("Could not create transaction: %v", err)
	}

	transactions, err := db.GetTransactions(tx.UserID)
	if err != nil {
		log.Fatalf("Could not get transactions: %v", err)
	}
	fmt.Println("Transactions:", transactions)
}