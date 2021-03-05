package main

import (
	"log"
	"os"

	"b-pay/config/database"
	"b-pay/config/middleware"
	"b-pay/config/migration"
	savingController "b-pay/controllers/savingcontroller"
	transactionController "b-pay/controllers/transactioncontroller"
	userController "b-pay/controllers/usercontroller"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Loads .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf(err.Error())
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	database.InitDB()
	migration.AutoMigrate(database.DB)
	// Initialize Gin with default settings.
	r := gin.Default()

	v1 := r.Group("/v1")
	{
		// Can be accessed without token
		public := v1.Group("/public")
		{
			// User Registration
			public.POST("/register", userController.RegisterUserHandler)
			// User Login
			public.POST("/login", userController.LoginHandler)
		}

		// Can be accessed with token.
		protected := v1.Group("/protected")
		protected.Use(middleware.AuthJWT())
		{
			saving := protected.Group("/s")
			{
				// Create a Saving account
				saving.POST("/create", savingController.CreateSavingHandler)
				// Get all Saving account owned by the User who accessed it.
				saving.GET("/", savingController.IndexSavingHandler)
				// Log into a Saving account.
				saving.POST("/login/:id", savingController.LoginSavingHandler)
				// Show a Saving data info.
				saving.GET("/:id", savingController.ShowSavingHandler)
				// Update a Saving data info. (Only Name and PIN)
				saving.PATCH("/update/:id", savingController.UpdateSavingHandler)
				// Delete a Saving account.
				saving.DELETE("/delete/:id", savingController.DeleteSavingHandler)
			}

			transaction := protected.Group("/t")
			{
				// Add a Transaction to a saving account.
				transaction.POST("/add", transactionController.CreateTransactionHandler)
			}

		}
	}

	r.Run(":" + port)
}
