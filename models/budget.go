package main

import (
    "database/sql"
    "fmt"
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

var (
    db    *sql.DB
    cache map[int]*Budget // Cache for Budgets
)

func Initialize() {
    initializeDatabase()
    initializeCache()
}

func initializeDatabase() {
    var err error
    db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
    if err != nil {
        panic(fmt.Sprintf("Failed to connect to database: %v", err))
    }

    if err = db.Ping(); err != nil {
        panic(fmt.Sprintf("Failed to ping database: %v", err))
    }
}

func initializeCache() {
    cache = make(map[int]*Budget) // Initialize the cache
}

func CreateBudget(budget Budget) (int, error) {
    budgetID, err := insertBudgetIntoDB(budget)
    if err != nil {
        return 0, err
    }

    budget.ID = budgetID
    cache[budgetID] = &budget // Cache the new budget

    return budgetID, nil
}

func insertBudgetIntoDB(budget Budget) (int, error) {
    var budgetID int
    query := `INSERT INTO budgets (amount, category, start_date, end_date, user_id) VALUES ($1, $2, $3, $4, $5) RETURNING id`
    err := db.QueryRow(query, budget.Amount, budget.Category, budget.StartDate, budget.EndDate, budget.UserID).Scan(&budgetID)
    return budgetID, err
}

func GetBudget(id int) (*Budget, error) {
    if budget, found := cache[id]; found {
        return budget, nil // Return from cache if found
    }

    return fetchBudgetFromDB(id)
}

func fetchBudgetFromDB(id int) (*Budget, error) {
    query := `SELECT id, amount, category, start_date, end_date, user_id FROM budgets WHERE id = $1`
    budget := &Budget{}
    err := db.QueryRow(query, id).Scan(&budget.ID, &budget.Amount, &budget.Category, &budget.StartDate, &budget.EndDate, &budget.UserID)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil // No result
        }
        return nil, err
    }

    cache[id] = budget // Cache fetched budget
    return budget, nil
}

func UpdateBudget(id int, budget Budget) error {
    if err := updateBudgetInDB(id, budget); err != nil {
        return err
    }
    cache[id] = &budget // Update the cache with the new budget data
    return nil
}

func updateBudgetInDB(id int, budget Budget) error {
    query := `UPDATE budgets SET amount = $1, category = $2, start_date = $3, end_date = $4, user_id = $5 WHERE id = $6`
    _, err := db.Exec(query, budget.Amount, budget.Category, budget.StartDate, budget.EndDate, budget.UserID, id)
    return err
}

func DeleteBudget(id int) error {
    if err := deleteBudgetFromDB(id); err != nil {
        return err
    }
    delete(cache, id) // Remove from cache
    return nil
}

func deleteBudgetFromDB(id int) error {
    _, err := db.Exec(`DELETE FROM budgets WHERE id = $1`, id)
    return err
}

func main() {
    Initialize()
    // Further logic can be added here
}