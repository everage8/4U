package service

import (
	"context"
	"errors"
	"fmt"

	"exam-tasks-backend/jwt"
	"exam-tasks-backend/models"
	"exam-tasks-backend/repository"
)

var ErrInvalidCredentials = errors.New("invalid login or password")

type AuthService struct {
	admins        *repository.AdminRepository
	jwtSecret     string
	jwtExpiryHours int
}

func NewAuthService(admins *repository.AdminRepository, jwtSecret string, jwtExpiryHours int) *AuthService {
	return &AuthService{
		admins:         admins,
		jwtSecret:      jwtSecret,
		jwtExpiryHours: jwtExpiryHours,
	}
}

func (s *AuthService) Login(ctx context.Context, login, password string) (string, *models.Admin, error) {
	admin, err := s.admins.FindByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, repository.ErrAdminNotFound) {
			return "", nil, ErrInvalidCredentials
		}
		return "", nil, fmt.Errorf("auth login: find admin: %w", err)
	}

	if err := jwt.ComparePassword(password, admin.Password); err != nil {
		return "", nil, ErrInvalidCredentials
	}

	token, err := jwt.GenerateToken(s.jwtSecret, s.jwtExpiryHours, admin.Seq, admin.Login, admin.Role)
	if err != nil {
		return "", nil, fmt.Errorf("auth login: sign token: %w", err)
	}
	return token, admin, nil
}

func (s *AuthService) SeedDefaultAdmin(ctx context.Context, login, password string) error {
	count, err := s.admins.Count(ctx)
	if err != nil {
		return fmt.Errorf("seed admin: count: %w", err)
	}
	if count > 0 {
		return nil
	}

	hashed, err := jwt.HashPassword(password)
	if err != nil {
		return fmt.Errorf("seed admin: hash password: %w", err)
	}

	admin := &models.Admin{
		Seq:      1,
		Login:    login,
		Password: hashed,
		Role:     "admin",
	}
	if err := s.admins.Create(ctx, admin); err != nil {
		return fmt.Errorf("seed admin: insert: %w", err)
	}
	return nil
}
