package main

import (
	"database/sql"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Budget struct {
	ID        int       `json:"id"`
	Amount    float64   `json:"amount"`
	Category  string    `json:"category"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	UserID    int       `json:"user_id"`
}

var db *sql.DB

func Initialize() {
	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}
}

func CreateBudget(budget Budget) (int, error) {
	var budgetID int
	query := `INSERT INTO budgets (amount, category, start_date, end_date, user_id) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := db.QueryRow(query, budget.Amount, budget.Category, budget.StartDate, budget.EndDate, budget.UserID).Scan(&budgetID)
	if err != nil {
		return 0, err
	}
	return budgetID, nil
}

func GetBudget(id int) (*Budget, error) {
	budget := &Budget{}
	query := `SELECT id, amount, category, start_date, end_date, user_id FROM budgets WHERE id = $1`
	row := db.QueryRow(query, id)
	err := row.Scan(&budget.ID, &budget.Amount, &budget.Category, &budget.StartDate, &budget.EndDate, &budget.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return budget, nil
}

func UpdateBudget(id int, budget Budget) error {
	query := `UPDATE budgets SET amount = $1, category = $2, start_date = $3, end_date = $4, user_id = $5 WHERE id = $6`
	_, err := db.Exec(query, budget.Amount, budget.Category, budget.StartDate, budget.EndDate, budget.UserID, id)
	return err
}

func DeleteBudget(id int) error {
	query := `DELETE FROM budgets WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}

func main() {
	Initialize()
}