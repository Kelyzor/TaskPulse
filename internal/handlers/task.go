package handlers

import (
	"net/http"
	"strconv"
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
	user, err := h.GetFromUserContext(c)

	if err != nil {
		return
	}

	status := c.Query("status")

	query := h.DB.Where("user_id = ?", user.ID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	pageInt, _ := strconv.Atoi(page)
	limitInt, _ := strconv.Atoi(limit)

	offset := (pageInt - 1) * limitInt

	var tasks []models.Task
	if result := query.
		Offset(offset).
		Limit(limitInt).
		Find(&tasks); result.Error != nil {
		ErrorResponse(c, http.StatusInternalServerError, result.Error.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, tasks)
}

func (h Handler) CreateTask(c *gin.Context) {
	user, err := h.GetFromUserContext(c)

	if err != nil {
		return
	}

	var taskInput CreateTaskInput

	if err := c.ShouldBindJSON(&taskInput); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	task := models.Task{
		UserID:      user.ID,
		Title:       taskInput.Title,
		Description: taskInput.Description,
		Status:      taskInput.Status,
	}

	if result := h.DB.Create(&task); result.Error != nil {
		ErrorResponse(c, http.StatusInternalServerError, result.Error.Error())
		return
	}

	SuccessResponse(c, http.StatusCreated, task)
}

func (h Handler) GetTask(c *gin.Context) {
	user, err := h.GetFromUserContext(c)

	if err != nil {
		return
	}

	taskID := c.Param("id")

	var task models.Task

	if result := h.DB.Where("id = ?", taskID).First(&task); result.Error != nil {
		ErrorResponse(c, http.StatusNotFound, "task not found")
		return
	}

	if user.ID != task.UserID {
		ErrorResponse(c, http.StatusForbidden, "Access denied")
		return
	}

	SuccessResponse(c, http.StatusOK, task)
}

func (h Handler) UpdateTask(c *gin.Context) {
	user, err := h.GetFromUserContext(c)

	if err != nil {
		return
	}

	taskID := c.Param("id")

	var task models.Task

	if result := h.DB.Where("id = ?", taskID).First(&task); result.Error != nil {
		ErrorResponse(c, http.StatusNotFound, "task not found")
		return
	}

	if task.UserID != user.ID {
		ErrorResponse(c, http.StatusForbidden, "Access denied")
		return
	}

	var updateInput UpdateTaskInput

	if err := c.ShouldBindJSON(&updateInput); err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if result := h.DB.Model(&task).Updates(models.Task{
		Title:       updateInput.Title,
		Description: updateInput.Description,
		Status:      updateInput.Status,
		Priority:    updateInput.Priority,
		DueDate:     updateInput.DueDate,
	}).Scan(&task); result.Error != nil {
		ErrorResponse(c, http.StatusInternalServerError, result.Error.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, task)
}

func (h Handler) DeleteTask(c *gin.Context) {
	user, err := h.GetFromUserContext(c)

	if err != nil {
		return
	}

	taskID := c.Param("id")
	var task models.Task

	if result := h.DB.Where("id = ?", taskID).First(&task); result.Error != nil {
		ErrorResponse(c, http.StatusNotFound, "task not found")
		return
	}

	if task.UserID != user.ID {
		ErrorResponse(c, http.StatusForbidden, "Access denied")
		return
	}

	if result := h.DB.Delete(&task); result.Error != nil {
		ErrorResponse(c, http.StatusInternalServerError, result.Error.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, gin.H{"message": "task deleted"})
}
