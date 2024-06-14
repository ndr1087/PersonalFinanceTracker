package main

import (
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	if err := loadEnv(); err != nil {
		panic("Failed to load environment variables")
	}

	router := gin.Default()

	// User route registrations
	registerUserRoutes(router)

	// Transaction route registrations
	registerTransactionRoutes(router)

	// Budget route registrations
	registerBudgetRoutes(router)

	// Starting the server
	startServer(router)
}

// registerUserRoutes adds routes related to user operations.
func registerUserRoutes(router *gin.Engine) {
	router.POST("/users/register", registerUser)
	router.POST("/users/login", loginUser)
}

// registerTransactionRoutes adds routes related to transaction operations.
func registerTransactionRoutes(router *gin.Engine) {
	router.GET("/transactions", getTransactions)
	router.POST("/transactions", addTransaction)
	router.DELETE("/transactions/:id", deleteTransaction)
}

// registerBudgetRoutes adds routes related to budget operations.
func registerBudgetRoutes(router *gin.Engine) {
	router.GET("/budgets", getBudgets)
	router.POST("/budgets", addBudget)
	router.PUT("/budgets/:id", updateBudget)
}

// startServer starts the Gin server on a specified PORT.
func startServer(router *gin.Engine) {
	router.Run(":" + os.Getenv("PORT"))
}

// Handlers for user routes
func registerUser(c *gin.Context)       {}
func loginUser(c *gin.Context)          {}

// Handlers for transaction routes
func getTransactions(c *gin.Context)    {}
func addTransaction(c *gin.Context)     {}
func deleteTransaction(c *gin.Context)  {}

// Handlers for budget routes
func getBudgets(c *gin.Context)         {}
func addBudget(c *gin.Context)          {}
func updateBudget(c *gin.Context)       {}

// loadEnv simulates the loading of environment variables. In real applications, this could 
// be replaced with actual logic to load configurations, like from a file or environment.
func loadEnv() error {
	return nil
}