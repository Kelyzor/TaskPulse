package main

import (
	"log"
	"taskpulse/internal/handlers"
	"taskpulse/internal/logger"
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
	_ = godotenv.Load()

	log := logger.InitLogger()
	defer log.Sync()

	db := InitDB()
	h := handlers.Handler{DB: db, Logger: log}

	router := gin.Default()
	router.Use(middleware.CORSMiddleware())

	router.GET("/health", h.HealthCheck)

	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", middleware.RateLimitMiddleware(3), h.RegisterUser)
			auth.POST("/login", middleware.RateLimitMiddleware(5), h.LoginUser)
		}
		users := api.Group("/users")
		{
			users.GET("/", h.GetUsers)
			users.GET("/me", middleware.AuthMiddleware(), h.GetMe)
		}
		tasks := api.Group("/tasks")
		{
			tasks.GET("", middleware.AuthMiddleware(), h.GetTasks)
			tasks.POST("", middleware.AuthMiddleware(), h.CreateTask)
			tasks.GET("/:id", middleware.AuthMiddleware(), h.GetTask)
			tasks.PUT("/:id", middleware.AuthMiddleware(), h.UpdateTask)
			tasks.DELETE("/:id", middleware.AuthMiddleware(), h.DeleteTask)
		}
		api.DELETE("/deleteUser/:id", h.DeleteUser)
	}

	router.Run(":8080")
}
