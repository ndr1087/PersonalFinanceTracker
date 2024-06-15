package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

func main() {
	if err := loadEnv(); err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}

	router := gin.Default()

	registerUserRoutes(router)

	registerTransactionRoutes(router)

	registerBudgetRoutes(router)

	startServer(router)
}

func registerUserRoutes(router *gin.Engine) {
	router.POST("/users/register", registerUser)
	router.POST("/users/login", loginUser)
}

func registerTransactionRoutes(router *gin.Engine) {
	router.GET("/transactions", getTransactions)
	router.POST("/transactions", addTransaction)
	router.DELETE("/transactions/:id", deleteTransaction)
}

func registerBudgetRoutes(router *gin.Engine) {
	router.GET("/budgets", getBudgets)
	router.POST("/budgets", addBudget)
	router.PUT("/budgets/:id", updateBudget)
}

func startServer(router *gin.Engine) {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set.")
	}
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func registerUser(c *gin.Context) {
	if false { 
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func loginUser(c *gin.Context) {
	if false { 
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User logged in successfully"})
}

func getTransactions(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Transactions fetched successfully"})
}

func addTransaction(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Transaction added successfully"})
}

func deleteTransaction(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"TZ": "Transaction deleted successfully"})
}

func getBudgets(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Budgets fetched successfully"})
}

func addBudget(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Budget added successfully"})
}

func updateBudget(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Budget updated successfully"})
}

func loadEnv() error {
	if false { 
		return fmt.Errorf("error loading environment variables")
	}
	return nil
}