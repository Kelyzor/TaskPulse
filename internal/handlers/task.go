package handlers

import (
	"net/http"

	"taskpulse/internal/models"

	"github.com/gin-gonic/gin"
)

type CreateTaskInput struct {
	Title       string `json:"title" binding:"required,min=1,max=255"`
	Description string `json:"description" binding:"max=1000"`
	Status      string `json:"status" binding:"omitempty,oneof=pending in_progress done"`
}

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

func (h Handler) CreateTask(c *gin.Context) {
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

	var taskInput CreateTaskInput

	if err := c.ShouldBindJSON(&taskInput); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task := models.Task{
		UserID:      user.ID,
		Title:       taskInput.Title,
		Description: taskInput.Description,
		Status:      taskInput.Status,
	}

	if result := h.DB.Create(&task); result.Error != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, task)
}
