package service

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"strings"

	"exam-tasks-backend/dto"
	"exam-tasks-backend/models"
	"exam-tasks-backend/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrInvalidID = errors.New("invalid task id")

type TaskService struct {
	repo *repository.TaskRepository
}

func NewTaskService(repo *repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) GetPublicTasks(ctx context.Context, subject, taskType string) ([]models.TaskResponse, error) {
	tasks, err := s.repo.FindPublic(ctx, subject, taskType)
	if err != nil {
		return nil, err
	}
	return models.ToResponseList(tasks), nil
}

func (s *TaskService) GetAdminTasks(
	ctx context.Context,
	search, subject, taskType string,
	page, limit int,
) (tasks []models.TaskResponse, total int, err error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 5
	}

	all, err := s.repo.FindForAdmin(ctx, subject, taskType)
	if err != nil {
		return nil, 0, err
	}

	query := strings.TrimSpace(search)
	query = strings.ToLower(query)
	query = strings.TrimPrefix(query, "#")

	if query != "" {
		filtered := make([]models.Task, 0, len(all))
		for _, t := range all {
			if strings.Contains(strconv.Itoa(t.Number), query) {
				filtered = append(filtered, t)
			}
		}
		all = filtered
	}

	sort.SliceStable(all, func(i, j int) bool { return all[i].Number < all[j].Number })

	total = len(all)
	start := (page - 1) * limit
	if start > total {
		start = total
	}
	end := start + limit
	if end > total {
		end = total
	}
	pageSlice := all[start:end]
	return models.ToResponseList(pageSlice), total, nil
}

func (s *TaskService) CreateTask(ctx context.Context, req dto.CreateTaskRequest) (models.TaskResponse, error) {
	maxNum, err := s.repo.FindMaxNumber(ctx)
	if err != nil {
		return models.TaskResponse{}, err
	}

	task := &models.Task{
		Subject:           req.Subject,
		TaskType:          req.TaskType,
		Number:            maxNum + 1,
		Condition:         req.Condition,
		ConditionImageURL: req.ConditionImageURL,
		Answer:            req.Answer,
		Solution:          req.Solution,
		SolutionImageURL:  req.SolutionImageURL,
	}

	if err := s.repo.Create(ctx, task); err != nil {
		return models.TaskResponse{}, err
	}
	return task.ToResponse(), nil
}

func (s *TaskService) UpdateTask(ctx context.Context, id string, req dto.UpdateTaskRequest) (models.TaskResponse, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.TaskResponse{}, ErrInvalidID
	}

	existing, err := s.repo.FindByID(ctx, oid)
	if err != nil {
		return models.TaskResponse{}, err
	}

	existing.Subject = req.Subject
	existing.TaskType = req.TaskType
	existing.Condition = req.Condition
	existing.ConditionImageURL = req.ConditionImageURL
	existing.Answer = req.Answer
	existing.Solution = req.Solution
	existing.SolutionImageURL = req.SolutionImageURL

	if err := s.repo.Update(ctx, oid, existing); err != nil {
		return models.TaskResponse{}, err
	}
	return existing.ToResponse(), nil
}

func (s *TaskService) DeleteTask(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidID
	}
	return s.repo.Delete(ctx, oid)
}
