package router

import (
	"os"
	"path/filepath"

	"exam-tasks-backend/config"
	"exam-tasks-backend/handlers"
	"exam-tasks-backend/middleware"

	"github.com/gin-gonic/gin"
)

func Build(
	cfg *config.Config,
	authHandler *handlers.AuthHandler,
	taskHandler *handlers.TaskHandler,
) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORSMiddleware(cfg.CORSOrigins))

	webRoot := os.Getenv("WEB_ROOT")
	if webRoot == "" {
		webRoot = "web"
	}
	r.Static("/css", filepath.Join(webRoot, "css"))
	r.Static("/js", filepath.Join(webRoot, "js"))
	r.StaticFile("/", filepath.Join(webRoot, "html", "main.html"))
	r.StaticFile("/login", filepath.Join(webRoot, "html", "login.html"))
	r.StaticFile("/admin", filepath.Join(webRoot, "html", "admin.html"))

	api := r.Group("/api/v1")
	{

		api.GET("/tasks", taskHandler.GetPublicTasks)

		api.POST("/auth/login", authHandler.Login)

		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		admin.Use(middleware.AdminMiddleware())
		{
			admin.GET("/tasks", taskHandler.GetAdminTasks)
			admin.POST("/tasks", taskHandler.CreateTask)
			admin.PUT("/tasks/:id", taskHandler.UpdateTask)
			admin.DELETE("/tasks/:id", taskHandler.DeleteTask)
		}
	}

	return r
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
