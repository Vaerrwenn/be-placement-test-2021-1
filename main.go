package main

import (
	"log"
	"os"

	"b-pay/config/database"
	"b-pay/config/migration"
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
		public := v1.Group("/public")
		{
			public.POST("/register", userController.RegisterUserHandler)
			public.POST("/login", userController.LoginHandler)
		}
	}

	r.Run(":" + port)
}
