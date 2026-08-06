package handlers

import (
	"errors"
	"net/http"

	"exam-tasks-backend/dto"
	"exam-tasks-backend/response"
	"exam-tasks-backend/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	token, admin, err := h.svc.Login(c.Request.Context(), req.Login, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			response.Error(c, http.StatusUnauthorized, "Неверный логин или пароль")
			return
		}
		response.Error(c, http.StatusInternalServerError, "Login failed: "+err.Error())
		return
	}

	payload := dto.LoginResponse{
		Token: token,
		Admin: dto.AdminPublic{
			ID:    admin.Seq,
			Login: admin.Login,
			Role:  admin.Role,
		},
	}
	response.Success(c, "Login successful", payload)
}
