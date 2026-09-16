package handlers

import (
	"net/http"
	"os"
	"taskpulse/internal/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func (h Handler) RegisterUser(c *gin.Context) {
	var input models.RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var existingUser models.User
	if result := h.DB.Where("email = ?", input.Email).First(&existingUser); result.Error == nil {
		ErrorResponse(c, http.StatusConflict, "Email already exists")
		return
	}

	hashedPassword, err := HashPassword(input.Password)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	user := models.User{
		Email:        input.Email,
		CreatedAt:    time.Now(),
		PasswordHash: hashedPassword,
	}

	if result := h.DB.Create(&user); result.Error != nil {
		ErrorResponse(c, http.StatusInternalServerError, result.Error.Error())
		return
	}

	SuccessResponse(c, http.StatusCreated, user)
}

func (h Handler) LoginUser(c *gin.Context) {
	var input models.RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var user models.User
	if result := h.DB.Where("email = ?", input.Email).First(&user); result.Error != nil {
		ErrorResponse(c, http.StatusUnauthorized, "invalid email or password")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
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
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, gin.H{"token": token})
}

func (h Handler) GetMe(c *gin.Context) {
	email, exists := c.Get("email")

	if !exists {
		ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	var user models.User
	if result := h.DB.Where("email = ?", email).First(&user); result.Error != nil {
		ErrorResponse(c, http.StatusNotFound, "user not found")
		return
	}

	SuccessResponse(c, http.StatusOK, user)
}

func (h Handler) GetUsers(c *gin.Context) {
	var users []models.User
	if result := h.DB.Find(&users); result.Error != nil {
		ErrorResponse(c, http.StatusInternalServerError, result.Error.Error())
		return
	}
	SuccessResponse(c, http.StatusOK, users)
}

func (h Handler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	result := h.DB.Delete(&models.User{}, id)
	if result.Error != nil {
		ErrorResponse(c, http.StatusInternalServerError, result.Error.Error())
		return
	}

	if result.RowsAffected == 0 {
		ErrorResponse(c, http.StatusNotFound, "user not found")
		return
	}

	SuccessResponse(c, http.StatusOK, gin.H{"message": "user deleted"})
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
