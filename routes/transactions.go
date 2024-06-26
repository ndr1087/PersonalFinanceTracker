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

var dbConnection *sql.DB

func initializeDatabase() *sql.DB {
    databaseConnectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))
    database, err := sql.Open("postgres", databaseConnectionString)
    if err != nil {
        panic(err)
    }
    if err = database.Ping(); err != nil {
        panic(err)
    }
    return database
}

func insertTransaction(transactionEntry Transaction) error {
    insertionQuery := `INSERT INTO transactions (amount, type, category, created_at) VALUES ($1, $2, $3, $4)`
    _, err := dbConnection.Exec(insertionQuery, transactionEntry.Amount, transactionEntry.Type, transactionEntry.Category, transactionEntry.CreatedAt)
    return err
}

func fetchAllTransactions() ([]Transaction, error) {
    selectQuery := "SELECT id, amount, type, category, created_at FROM transactions"
    rows, err := dbConnection.Query(selectQuery)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    allTransactions := make([]Transaction, 0)
    for rows.Next() {
        var singleTransaction Transaction
        if err := rows.Scan(&singleTransaction.ID, &singleTransaction.Amount, &singleTransaction.Type, &singleTransaction.Category, &singleTransaction.CreatedAt); err != nil {
            return nil, err
        }
        allTransactions = append(allTransactions, singleTransaction)
    }
    return allTransactions, nil
}

func modifyTransaction(updatedTransaction Transaction) error {
    updateQuery := `UPDATE transactions SET amount = $1, type = $2, category = $3, created_at = $4 WHERE id = $5`
    _, err := dbConnection.Exec(update, updatedTransaction.Amount, updatedTransaction.Type, updatedTransaction.Category, updatedTransaction.CreatedAt, updatedTransaction.ID)
    return err
}

func removeTransaction(transactionID int) error {
    deletionQuery := "DELETE FROM transactions WHERE id = $1"
    _, err := dbConnection.Exec(deletionQuery, transactionID)
    return err
}

func main() {
    dbConnection = initializeDatabase() // Use the initialized db connection throughout the application
    // Application logic goes here
}