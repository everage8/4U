package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func JSON(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, Response{
		Status:  statusText(status),
		Message: message,
		Data:    data,
	})
}

func Success(c *gin.Context, message string, data interface{}) {
	JSON(c, http.StatusOK, message, data)
}

func Created(c *gin.Context, message string, data interface{}) {
	JSON(c, http.StatusCreated, message, data)
}

func Error(c *gin.Context, status int, message string) {
	JSON(c, status, message, nil)
}

func statusText(status int) string {
	switch {
	case status == http.StatusCreated:
		return "created"
	case status >= 200 && status < 300:
		return "success"
	case status >= 400 && status < 500:
		return "error"
	default:
		return "error"
	}
}
