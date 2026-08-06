package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware(origins []string) gin.HandlerFunc {
	allowAll := false
	for _, o := range origins {
		if o == "*" {
			allowAll = true
			break
		}
	}

	cfg := cors.Config{
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders: []string{"Content-Length"},
		MaxAge:        12 * time.Hour,
	}

	if allowAll {
		cfg.AllowAllOrigins = true
	} else {
		cfg.AllowOrigins = origins
	}

	return cors.New(cfg)
}
