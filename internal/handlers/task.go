package handlers

import (
	"net/http"
	"time"

	"taskpulse/internal/models"

	"github.com/gin-gonic/gin"
)

type CreateTaskInput struct {
	Title       string     `json:"title" binding:"required,min=1,max=255"`
	Description string     `json:"description" binding:"max=1000"`
	Status      string     `json:"status" binding:"omitempty,oneof=pending in_progress done"`
	Priority    string     `json:"priority" binding:"omitempty,oneof=low medium high"`
	DueDate     *time.Time `json:"due_date"`
}

type UpdateTaskInput struct {
	Title       string     `json:"title" binding:"required,min=1,max=255"`
	Description string     `json:"description" binding:"max=1000"`
	Status      string     `json:"status" binding:"omitempty,oneof=pending in_progress done"`
	Priority    string     `json:"priority" binding:"omitempty,oneof=low medium high"`
	DueDate     *time.Time `json:"due_date"`
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

func (h Handler) GetTask(c *gin.Context) {
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

	taskID := c.Param("id")

	var task models.Task

	if result := h.DB.Where("id = ?", taskID).First(&task); result.Error != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	if user.ID != task.UserID {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.IndentedJSON(http.StatusOK, task)
}

func (h Handler) UpdateTask(c *gin.Context) {
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

	taskID := c.Param("id")

	var task models.Task

	if result := h.DB.Where("id = ?", taskID).First(&task); result.Error != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	if task.UserID != user.ID {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var updateInput UpdateTaskInput

	if err := c.ShouldBindJSON(&updateInput); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if result := h.DB.Model(&task).Updates(models.Task{
		Title:       updateInput.Title,
		Description: updateInput.Description,
		Status:      updateInput.Status,
		Priority:    updateInput.Priority,
		DueDate:     updateInput.DueDate,
	}).Scan(&task); result.Error != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, task)
}

func (h Handler) DeleteTask(c *gin.Context) {
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

	taskID := c.Param("id")
	var task models.Task

	if result := h.DB.Where("id = ?", taskID).First(&task); result.Error != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	if task.UserID != user.ID {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if result := h.DB.Delete(&task); result.Error != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "task deleted"})
}
