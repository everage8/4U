package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"exam-tasks-backend/dto"
	"exam-tasks-backend/models"
	"exam-tasks-backend/repository"
	"exam-tasks-backend/response"
	"exam-tasks-backend/service"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	svc *service.TaskService
}

func NewTaskHandler(svc *service.TaskService) *TaskHandler {
	return &TaskHandler{svc: svc}
}

func (h *TaskHandler) GetPublicTasks(c *gin.Context) {
	subject := c.Query("subject")
	taskType := c.Query("type")

	tasks, err := h.svc.GetPublicTasks(c.Request.Context(), subject, taskType)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to load tasks: "+err.Error())
		return
	}

	if tasks == nil {
		tasks = []models.TaskResponse{}
	}
	response.Success(c, "OK", tasks)
}

func (h *TaskHandler) GetAdminTasks(c *gin.Context) {
	search := c.Query("search")
	subject := c.Query("subject")
	taskType := c.Query("type")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))

	tasks, total, err := h.svc.GetAdminTasks(c.Request.Context(), search, subject, taskType, page, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to load tasks: "+err.Error())
		return
	}

	payload := gin.H{
		"tasks": tasks,
		"total": total,
	}
	response.Success(c, "OK", payload)
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req dto.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	created, err := h.svc.CreateTask(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create task: "+err.Error())
		return
	}
	response.Created(c, "Task created", created)
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	updated, err := h.svc.UpdateTask(c.Request.Context(), id, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidID):
			response.Error(c, http.StatusBadRequest, "Invalid task id")
		case errors.Is(err, repository.ErrTaskNotFound):
			response.Error(c, http.StatusNotFound, "Task not found")
		default:
			response.Error(c, http.StatusInternalServerError, "Failed to update task: "+err.Error())
		}
		return
	}
	response.Success(c, "Task updated", updated)
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteTask(c.Request.Context(), id); err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidID):
			response.Error(c, http.StatusBadRequest, "Invalid task id")
		case errors.Is(err, repository.ErrTaskNotFound):
			response.Error(c, http.StatusNotFound, "Task not found")
		default:
			response.Error(c, http.StatusInternalServerError, "Failed to delete task: "+err.Error())
		}
		return
	}
	response.Success(c, "Task deleted", gin.H{"success": true})
}
