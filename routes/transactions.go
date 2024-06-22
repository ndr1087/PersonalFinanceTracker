package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Transaction struct {
	ID        int
	Amount    float64
	Type      string
	Category  string
	CreatedAt time.Time
}

var db *sql.DB

func init() {
	var err error
	connStr := fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}
}

func AddTransaction(t Transaction) error {
	query := `INSERT INTO transactions (amount, type, category, created_at) VALUES ($1, $2, $3, $4)`
	_, err := db.Exec(query, t.Amount, t.Type, t.Category, t.CreatedAt)
	return err
}

func GetAllTransactions() ([]Transaction, error) {
	rows, err := db.Query("SELECT id, amount, type, category, created_at FROM transactions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var t Transaction
		if err := rows.Scan(&t.ID, &t.Amount, &t.Type, &t.Category, &t.CreatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}
	return transactions, nil
}

func UpdateTransaction(t Transaction) error {
	query := `UPDATE transactions SET amount = $1, type = $2, category = $3, created_at = $4 WHERE id = $5`
	_, err := db.Exec(query, t.Amount, t.Type, t.Category, t.CreatedA	t, t.ID)
	return err
}

func DeleteTransaction(id int) error {
	query := "DELETE FROM transactions WHERE id = $1"
	_, err := db.Exec(query, id)
	return err
}

func main() {
}