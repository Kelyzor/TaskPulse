package handlers

import (
	"net/http"
	"strconv"
	"time"

	"taskpulse/internal/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
		h.Logger.Error("failed to get user", zap.Error(err))
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
		h.Logger.Error("failed to fetch tasks list", zap.Error(result.Error))
		ErrorResponse(c, http.StatusInternalServerError, result.Error.Error())
		return
	}

	h.Logger.Info("tasks list fetched successfully", zap.Int("count", len(tasks)))
	SuccessResponse(c, http.StatusOK, tasks)
}

func (h Handler) CreateTask(c *gin.Context) {
	user, err := h.GetFromUserContext(c)

	if err != nil {
		h.Logger.Error("failed to get user", zap.Error(err))
		return
	}

	var taskInput CreateTaskInput

	if err := c.ShouldBindJSON(&taskInput); err != nil {
		h.Logger.Error("invalid task input", zap.Error(err))
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	task := models.Task{
		UserID:      user.ID,
		Title:       taskInput.Title,
		Description: taskInput.Description,
		Status:      taskInput.Status,
		Priority:    taskInput.Priority,
		DueDate:     taskInput.DueDate,
	}

	if result := h.DB.Create(&task); result.Error != nil {
		h.Logger.Error("failed to create task", zap.Error(result.Error))
		ErrorResponse(c, http.StatusInternalServerError, result.Error.Error())
		return
	}

	h.Logger.Info("task successfully created")
	SuccessResponse(c, http.StatusCreated, task)
}

func (h Handler) GetTask(c *gin.Context) {
	user, err := h.GetFromUserContext(c)

	if err != nil {
		h.Logger.Error("failed to get user", zap.Error(err))
		return
	}

	taskID := c.Param("id")

	var task models.Task

	if result := h.DB.Where("id = ?", taskID).First(&task); result.Error != nil {
		h.Logger.Error("task not found", zap.Error(result.Error))
		ErrorResponse(c, http.StatusNotFound, "task not found")
		return
	}

	if user.ID != task.UserID {
		h.Logger.Warn("user access denied", zap.String("email", user.Email))
		ErrorResponse(c, http.StatusForbidden, "Access denied")
		return
	}

	h.Logger.Info("task fetched successfully", zap.String("email", user.Email))
	SuccessResponse(c, http.StatusOK, task)
}

func (h Handler) UpdateTask(c *gin.Context) {
	user, err := h.GetFromUserContext(c)

	if err != nil {
		h.Logger.Error("failed to get user", zap.Error(err))
		return
	}

	taskID := c.Param("id")

	var task models.Task

	if result := h.DB.Where("id = ?", taskID).First(&task); result.Error != nil {
		h.Logger.Error("task not found", zap.Error(result.Error))
		ErrorResponse(c, http.StatusNotFound, "task not found")
		return
	}

	if task.UserID != user.ID {
		h.Logger.Warn("user access denied", zap.String("email", user.Email))
		ErrorResponse(c, http.StatusForbidden, "Access denied")
		return
	}

	var updateInput UpdateTaskInput

	if err := c.ShouldBindJSON(&updateInput); err != nil {
		h.Logger.Error("invalid task input", zap.Error(err))
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
		h.Logger.Error("failed to update task", zap.Error(result.Error))
		ErrorResponse(c, http.StatusInternalServerError, result.Error.Error())
		return
	}

	h.Logger.Info("task successfully updated", zap.String("email", user.Email))
	SuccessResponse(c, http.StatusOK, task)
}

func (h Handler) DeleteTask(c *gin.Context) {
	user, err := h.GetFromUserContext(c)

	if err != nil {
		h.Logger.Error("failed to get user", zap.Error(err))
		return
	}

	taskID := c.Param("id")
	var task models.Task

	if result := h.DB.Where("id = ?", taskID).First(&task); result.Error != nil {
		h.Logger.Error("task not found", zap.Error(result.Error))
		ErrorResponse(c, http.StatusNotFound, "task not found")
		return
	}

	if task.UserID != user.ID {
		h.Logger.Warn("user access denied", zap.String("email", user.Email))
		ErrorResponse(c, http.StatusForbidden, "Access denied")
		return
	}

	if result := h.DB.Delete(&task); result.Error != nil {
		h.Logger.Error("failed to delete task", zap.Error(result.Error))
		ErrorResponse(c, http.StatusInternalServerError, result.Error.Error())
		return
	}

	h.Logger.Info("task successfully deleted", zap.String("email", user.Email))
	SuccessResponse(c, http.StatusOK, gin.H{"message": "task deleted"})
}
