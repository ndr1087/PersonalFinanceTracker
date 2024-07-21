package main

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    _ "github.com/lib/pq"
)

var databaseConnection *sql.DB

type Budget struct {
    ID     int
    Name   string
    Amount float64
}

func initializeDatabase() {
    var err error
    databaseConnectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

    databaseConnection, err = sql.Open("postgres", databaseConnectionString)
    if err != nil {
        panic(err)
    }

    if err = databaseConnection.Ping(); err != nil {
        panic(err)
    }
}

func insertBudget(budget *Budget) error {
    insertQuery := `INSERT INTO budgets (name, amount) VALUES ($1, $2) RETURNING id`
    err := databaseConnection.QueryRow(insertQuery, budget.Name, budget.Amount).Scan(&budget.ID)
    if err != nil {
        return err
    }
    return nil
}

func fetchBudgetByID(id int) (*Budget, error) {
    budget := &Budget{}
    selectQuery := `SELECT id, name, amount FROM budgets WHERE id = $1`
    err := databaseConnection.QueryRow(selectQuery, id).Scan(&budget.ID, &budget.Name, &budget.Amount)
    if err != nil {
        return nil, err
    }
    return budget, nil
}

func modifyBudget(budget *Budget) error {
    updateQuery := `UPDATE budgets SET name = $1, amount = $2 WHERE id = $3`
    _, err := databaseConnection.Exec(updateQuery, budget.Name, budget.Amount, budget.ID)
    if err != nil {
        return err
    }
    return nil
}

func removeBudget(id int) error {
    deleteQuery := `DELETE FROM budgets WHERE id = $1`
    _, err := databaseConnection.Exec(deleteQuery, id)
    if err != nil {
        return err
    }
    return nil
}

func retrieveAllBudgets() ([]Budget, error) {
    var budgets []Budget
    fetchAllQuery := `SELECT id, name, amount FROM budgets`
    rows, err := databaseConnection.Query(fetchAllQuery)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var budget Budget
        if err := rows.Scan(&budget.ID, &budget.Name, &budget.Amount); err != nil {
            return nil, err
        }
        budgets = append(budgets, budget)
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return budgets, nil
}

func main() {
    initializeDatabase()

    budgets, err := retrieveAllBudgets()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Budgets:", budgets)
}