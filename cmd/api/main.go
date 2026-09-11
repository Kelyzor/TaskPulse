package main

import (
	"log"
	"taskpulse/internal/handlers"
	"taskpulse/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	dsn := "host=db user=postgres password=secretpassword dbname=taskpulse_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect DB:", err)
	}

	return db
}

func main() {
	// Загружаем переменные из .env
	_ = godotenv.Load()

	db := InitDB()
	h := handlers.Handler{DB: db}

	router := gin.Default()

	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", h.RegisterUser)
			auth.POST("/login", h.LoginUser)
		}
		users := api.Group("/users")
		{
			users.GET("/", h.GetUsers)
			users.GET("/me", middleware.AuthMiddleware(), h.GetMe)
		}
		api.DELETE("/deleteUser/:id", h.DeleteUser)
	}

	router.Run(":8080")
}
