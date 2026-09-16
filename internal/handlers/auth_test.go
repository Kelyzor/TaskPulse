package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"taskpulse/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to test db: %v", err)
	}

	db.AutoMigrate(&models.User{})

	return db
}

func TestRegisterUser_Success(t *testing.T) {
	db := setupTestDB(t)
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	router := gin.Default()
	handler := Handler{
		DB:     db,
		Logger: logger,
	}
	router.POST("/register", handler.RegisterUser)

	input := models.RegisterInput{
		Email:    "test@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestRegisterUser_DuplicateEmail(t *testing.T) {
	db := setupTestDB(t)
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	db.Create(&models.User{
		Email:        "test@example.com",
		PasswordHash: "hash",
	})

	router := gin.Default()
	handler := Handler{
		DB:     db,
		Logger: logger,
	}
	router.POST("/register", handler.RegisterUser)

	input := models.RegisterInput{
		Email:    "test@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestLoginUser_Success(t *testing.T) {
	db := setupTestDB(t)
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	hashedPassword, _ := HashPassword("password123")
	db.Create(&models.User{
		Email:        "test@example.com",
		PasswordHash: hashedPassword,
	})

	router := gin.Default()
	handler := Handler{
		DB:     db,
		Logger: logger,
	}
	router.POST("/login", handler.LoginUser)

	input := models.RegisterInput{
		Email:    "test@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLoginUser_InvalidPassword(t *testing.T) {
	db := setupTestDB(t)
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	hashedPassword, _ := HashPassword("password123")
	db.Create(&models.User{
		Email:        "test@example.com",
		PasswordHash: hashedPassword,
	})

	router := gin.Default()
	handler := Handler{
		DB:     db,
		Logger: logger,
	}
	router.POST("/login", handler.LoginUser)

	input := models.RegisterInput{
		Email:    "test@example.com",
		Password: "invalidPassword",
	}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLoginUser_InvalidEmail(t *testing.T) {
	db := setupTestDB(t)
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	hashedPassword, _ := HashPassword("password123")
	db.Create(&models.User{
		Email:        "test@example.com",
		PasswordHash: hashedPassword,
	})

	router := gin.Default()
	handler := Handler{
		DB:     db,
		Logger: logger,
	}
	router.POST("/login", handler.LoginUser)

	input := models.RegisterInput{
		Email:    "wrondTest@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
