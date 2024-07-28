package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type FinancialTransaction struct {
	ID          int
	UserID      int
	Amount      float64
	TransactionDate time.Time
	Category    string
	Description string
}

type TransactionDatabase struct {
	connection *sql.DB
}

func NewTransactionDatabase() *TransactionDatabase {
	databaseConnection, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME")))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	if err := databaseConnection.Ping(); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	return &TransactionDatabase{connection: databaseConnection}
}

func (tdb *TransactionDatabase) AddTransaction(ft *FinancialTransaction) error {
	insertQuery := `INSERT INTO transactions (user_id, amount, date, category, description) VALUES ($1, $2, $3, $4, $5)`
	_, err := tdb.connection.Exec(insertQuery, ft.UserID, ft.Amount, ft.TransactionDate, ft.Category, ft.Description)
	if err != nil {
		log.Printf("Insert transaction failed: %v", err)
		return err
	}
	log.Println("Transaction inserted successfully")
	return nil
}

func (tdb *TransactionDatabase) FetchTransactions(userID int) ([]FinancialTransaction, error) {
	transactionsList := []FinancialTransaction{}
	retrieveQuery := `SELECT id, user_id, amount, date, category, description FROM transactions WHERE user_id = $1`
	rows, err := tdb.connection.Query(retrieveQuery, userID)
	if err != nil {
		log.Printf("Retrieving transactions failed: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ft FinancialTransaction
		err := rows.Scan(&ft.ID, &ft.UserID, &ft.Amount, &ft.TransactionDate, &ft.Category, &ft.Description)
		if err != nil {
			log.Printf("Transaction scan failed: %v", err)
			return nil, err
		}
		transactionsList = append(transactionsList, ft)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Iteration through records failed: %v", err)
		return nil, err
	}

	log.Printf("Transactions retrieved: %d", len(transactionsList))
	return transactionsList, nil
}

func (tdb *TransactionDatabase) ModifyTransaction(ft *FinancialTransaction) error {
	updateQuery := `UPDATE transactions SET amount = $2, date = $3, category = $4, description = $5 WHERE id = $1`
	_, err := tdb.connection.Exec(updateQuery, ft.ID, ft.Amount, ft.TransactionDate, ft.Category, ft.Description)
	if err != nil {
		log.Printf("Updating transaction failed: %v", err)
		return err
	}
	log.Println("Transaction updated successfully")
	return nil
}

func (tdb *TransactionDatabase) RemoveTransaction(transactionID int) error {
	deleteQuery := `DELETE FROM transactions WHERE id = $1`
	_, err := tdb.connection.Exec(deleteQuery, transactionID)
	if err != nil {
		log.Printf("Deleting transaction failed: %v", err)
		return err
	}
	log.Println("Transaction deleted successfully")
	return nil
}

func main() {
	transDB := NewTransactionDatabase()
	defer transDB.connection.Close()

	newTransaction := &FinancialTransaction{
		UserID:         1,
		Amount:         199.99,
		TransactionDate: time.Now(),
		Category:       "Electronics",
		Description:    "New headphones",
	}
	err := transDB.AddTransaction(newTransaction)
	if err != nil {
		log.Fatalf("Transaction creation failed: %v", err)
	}

	transactions, err := transDB.FetchTransactions(newTransaction.UserID)
	if err != nil {
		log.Fatalf("Failed to fetch transactions: %v", err)
	}
	fmt.Println("Transactions:", transactions)
}