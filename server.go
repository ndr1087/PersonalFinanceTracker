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
    router.Run(":" + os.Getenv("PORT"))
}

func registerUser(c *gin.Context) {}
func loginUser(c *gin.Context) {}

func getTransactions(c *gin.Context) {}
func addTransaction(c *gin.Context) {}
func deleteTransaction(c *gin.Context) {}

func getBudgets(c *gin.Context) {}
func addBudget(c *gin.Context) {}
func updateBudget(c *gin.Context) {}

func loadEnv() error {
    return nil
}