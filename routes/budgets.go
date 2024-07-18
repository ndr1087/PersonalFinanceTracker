package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

type Budget struct {
	ID      int
	Name    string
	Amount  float64
}

func initDB() {
	var err error
	dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	db, err = sql.Open("postgres", dbInfo)
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}
}

func createBudget(budget *Budget) error {
	query := `INSERT INTO budgets (name, amount) VALUES ($1, $2) RETURNING id`
	err := db.QueryRow(query, budget.Name, budget.Amount).Scan(&budget.ID)
	if err != nil {
		return err
	}
	return nil
}

func getBudgetByID(id int) (*Budget, error) {
	budget := &Budget{}
	query := `SELECT id, name, amount FROM budgets WHERE id = $1`
	err := db.QueryRow(query, id).Scan(&budget.ID, &budget.Name, &budget.Amount)
	if err != nil {
		return nil, err
	}
	return budget, nil
}

func updateBudget(budget *Budget) error {
	query := `UPDATE budgets SET name = $1, amount = $2 WHERE id = $3`
	_, err := db.Exec(query, budget.Name, budget.Amount, budget.ID)
	if err != nil {
		return err
	}
	return nil
}

func deleteBudget(id int) error {
	query := `DELETE FROM budgets WHERE id = $1`
	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	initDB()
}