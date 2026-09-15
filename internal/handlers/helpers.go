package handlers

import (
	"errors"
	"net/http"
	"taskpulse/internal/models"

	"github.com/gin-gonic/gin"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrUserNotFound = errors.New("user not found")
)

func (h Handler) GetFromUserContext(c *gin.Context) (*models.User, error) {
	email, exists := c.Get("email")

	if !exists {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return nil, ErrUnauthorized
	}

	var user models.User

	if result := h.DB.Where("email = ?", email).First(&user); result.Error != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return nil, ErrUserNotFound
	}

	return &user, nil
}
