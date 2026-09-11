package handlers

import (
	"net/http"

	"taskpulse/internal/models"

	"github.com/gin-gonic/gin"
)

func (h Handler) GetTasks(c *gin.Context) {
	email, exists := c.Get("email")

	if !exists {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var user models.User

	if result := h.DB.Where("email = ?", email).First(&user); result.Error != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	var tasks []models.Task

	if result := h.DB.Where("user_id = ?", user.ID).Find(&tasks); result.Error != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, tasks)
}
