package middleware

import (
	"strings"

	"exam-tasks-backend/jwt"
	"exam-tasks-backend/response"

	"github.com/gin-gonic/gin"
)

const (
	CtxUserID = "user_id"
	CtxLogin  = "login"
	CtxRole   = "role"
)

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			response.Error(c, 401, "Authorization header is required")
			c.Abort()
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
			response.Error(c, 401, "Authorization header must be in the form 'Bearer <token>'")
			c.Abort()
			return
		}

		claims, err := jwt.ValidateToken(secret, strings.TrimSpace(parts[1]))
		if err != nil {
			response.Error(c, 401, "Invalid or expired token")
			c.Abort()
			return
		}

		c.Set(CtxUserID, claims.UserID)
		c.Set(CtxLogin, claims.Login)
		c.Set(CtxRole, claims.Role)
		c.Next()
	}
}
