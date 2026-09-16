package handlers

import (
	"net/http"
	"os"
	"taskpulse/internal/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Handler struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

func (h Handler) RegisterUser(c *gin.Context) {
	var input models.RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		h.Logger.Error("invalid register input", zap.Error(err))
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var existingUser models.User
	if result := h.DB.Where("email = ?", input.Email).First(&existingUser); result.Error == nil {
		h.Logger.Warn("registration attempt with existing email", zap.String("email", input.Email))
		ErrorResponse(c, http.StatusConflict, "Email already exists")
		return
	}

	hashedPassword, err := HashPassword(input.Password)
	if err != nil {
		h.Logger.Error("failed to hash password", zap.Error(err))
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	user := models.User{
		Email:        input.Email,
		CreatedAt:    time.Now(),
		PasswordHash: hashedPassword,
	}

	if result := h.DB.Create(&user); result.Error != nil {
		h.Logger.Error("failed to create user", zap.Error(result.Error), zap.String("email", input.Email))
		ErrorResponse(c, http.StatusInternalServerError, result.Error.Error())
		return
	}

	h.Logger.Info("user registered successfully", zap.String("email", input.Email), zap.Uint("user_id", user.ID))
	SuccessResponse(c, http.StatusCreated, user)
}

func (h Handler) LoginUser(c *gin.Context) {
	var input models.RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		h.Logger.Error("invalid login input", zap.Error(err))
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var user models.User
	if result := h.DB.Where("email = ?", input.Email).First(&user); result.Error != nil {
		h.Logger.Warn("login with invalid email or password", zap.String("email", input.Email))
		ErrorResponse(c, http.StatusUnauthorized, "invalid email or password")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		h.Logger.Warn("login with invalid email or password", zap.String("email", input.Email))
		ErrorResponse(c, http.StatusUnauthorized, "invalid email or password")
		return
	}

	claims := models.Claims{
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		h.Logger.Error("failed to generate JWT token", zap.Error(err))
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.Logger.Info("user login successfully", zap.String("email", input.Email), zap.Uint("user_id", user.ID))
	SuccessResponse(c, http.StatusOK, gin.H{"token": token})
}

func (h Handler) GetMe(c *gin.Context) {
	email, exists := c.Get("email")

	if !exists {
		h.Logger.Error("invalid email")
		ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	var user models.User
	if result := h.DB.Where("email = ?", email).First(&user); result.Error != nil {
		h.Logger.Error("user with email not found", zap.Error(result.Error))
		ErrorResponse(c, http.StatusNotFound, "user not found")
		return
	}

	h.Logger.Info("user profile fetched", zap.String("email", email.(string)))
	SuccessResponse(c, http.StatusOK, user)
}

func (h Handler) GetUsers(c *gin.Context) {
	var users []models.User
	if result := h.DB.Find(&users); result.Error != nil {
		h.Logger.Error("failed to find users", zap.Error(result.Error))
		ErrorResponse(c, http.StatusInternalServerError, result.Error.Error())
		return
	}
	h.Logger.Info("users list fetched", zap.Int("count", len(users)))
	SuccessResponse(c, http.StatusOK, users)
}

func (h Handler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	result := h.DB.Delete(&models.User{}, id)
	if result.Error != nil {
		h.Logger.Error("failed to delete user", zap.Error(result.Error))
		ErrorResponse(c, http.StatusInternalServerError, result.Error.Error())
		return
	}

	if result.RowsAffected == 0 {
		h.Logger.Error("user not found", zap.Error(result.Error))
		ErrorResponse(c, http.StatusNotFound, "user not found")
		return
	}

	h.Logger.Info("user deleted successfully", zap.String("user_id", id))
	SuccessResponse(c, http.StatusOK, gin.H{"message": "user deleted"})
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
