package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"exam-tasks-backend/config"
	"exam-tasks-backend/database"
	"exam-tasks-backend/handlers"
	"exam-tasks-backend/repository"
	"exam-tasks-backend/router"
	"exam-tasks-backend/service"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	cfg := config.Load()
	log.Printf("config: loaded (port=%s, mongo=%s/%s, admin=%s)",
		cfg.Port, cfg.MongoURI, cfg.MongoDB, cfg.AdminLogin)

	db, err := database.Connect(cfg.MongoURI, cfg.MongoDB)
	if err != nil {
		log.Fatalf("mongo: %v", err)
	}
	log.Printf("mongo: connected to %s/%s", cfg.MongoURI, cfg.MongoDB)

	if err := ensureIndexes(context.Background(), db); err != nil {
		log.Fatalf("mongo: ensure indexes: %v", err)
	}
	log.Printf("mongo: indexes ensured")

	adminRepo := repository.NewAdminRepository(db)
	taskRepo := repository.NewTaskRepository(db)

	authSvc := service.NewAuthService(adminRepo, cfg.JWTSecret, cfg.JWTExpiryHours)
	taskSvc := service.NewTaskService(taskRepo)

	authHandler := handlers.NewAuthHandler(authSvc)
	taskHandler := handlers.NewTaskHandler(taskSvc)

	seedCtx, seedCancel := context.WithTimeout(context.Background(), 5*time.Second)
	if err := authSvc.SeedDefaultAdmin(seedCtx, cfg.AdminLogin, cfg.AdminPassword); err != nil {
		seedCancel()
		log.Fatalf("seed admin: %v", err)
	}
	seedCancel()
	log.Printf("seed: default admin ready (login=%s)", cfg.AdminLogin)

	r := router.Build(cfg, authHandler, taskHandler)
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Printf("http: listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http: ListenAndServe: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Printf("shutdown: signal received, draining…")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("http: shutdown error: %v", err)
	}
	if err := database.Disconnect(shutdownCtx, db.Client()); err != nil {
		log.Printf("mongo: disconnect error: %v", err)
	}
	log.Printf("shutdown: done")
}

func ensureIndexes(ctx context.Context, db *mongo.Database) error {
	taskIdxCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err := db.Collection("tasks").Indexes().CreateOne(taskIdxCtx, mongo.IndexModel{
		Keys: bson.D{{Key: "subject", Value: 1}, {Key: "task_type", Value: 1}},
	})
	if err != nil {
		return err
	}

	adminIdxCtx, cancel2 := context.WithTimeout(ctx, 5*time.Second)
	defer cancel2()
	_, err = db.Collection("admins").Indexes().CreateOne(adminIdxCtx, mongo.IndexModel{
		Keys:    bson.D{{Key: "login", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	return err
}
