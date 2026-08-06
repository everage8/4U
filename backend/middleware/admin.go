package middleware

import (
	"exam-tasks-backend/response"

	"github.com/gin-gonic/gin"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, ok := c.Get(CtxRole)
		if !ok {
			response.Error(c, 401, "Authentication required")
			c.Abort()
			return
		}
		roleStr, _ := role.(string)
		if roleStr != "admin" {
			response.Error(c, 403, "Admin access required")
			c.Abort()
			return
		}
		c.Next()
	}
}
