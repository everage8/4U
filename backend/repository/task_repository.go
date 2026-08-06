package repository

import (
	"context"
	"errors"
	"fmt"

	"exam-tasks-backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrTaskNotFound = errors.New("task not found")

type TaskRepository struct {
	coll *mongo.Collection
}

func NewTaskRepository(db *mongo.Database) *TaskRepository {
	return &TaskRepository{coll: db.Collection("tasks")}
}

func (r *TaskRepository) FindPublic(ctx context.Context, subject, taskType string) ([]models.Task, error) {
	filter := buildSubjectTypeFilter(subject, taskType)
	return r.findManySorted(ctx, filter)
}

func (r *TaskRepository) FindForAdmin(ctx context.Context, subject, taskType string) ([]models.Task, error) {
	filter := buildSubjectTypeFilter(subject, taskType)
	return r.findManySorted(ctx, filter)
}

func (r *TaskRepository) FindMaxNumber(ctx context.Context) (int, error) {
	opts := options.FindOne().SetSort(bson.D{{Key: "number", Value: -1}}).SetProjection(bson.M{"number": 1})
	var doc struct {
		Number int `bson:"number"`
	}
	if err := r.coll.FindOne(ctx, bson.M{}, opts).Decode(&doc); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return 0, nil
		}
		return 0, fmt.Errorf("find max number: %w", err)
	}
	return doc.Number, nil
}

func (r *TaskRepository) Create(ctx context.Context, task *models.Task) error {
	res, err := r.coll.InsertOne(ctx, task)
	if err != nil {
		return fmt.Errorf("insert task: %w", err)
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		task.ID = oid
	}
	return nil
}

func (r *TaskRepository) Update(ctx context.Context, id primitive.ObjectID, task *models.Task) error {
	update := bson.M{
		"$set": bson.M{
			"subject":             task.Subject,
			"task_type":           task.TaskType,
			"condition":           task.Condition,
			"condition_image_url": task.ConditionImageURL,
			"answer":              task.Answer,
			"solution":            task.Solution,
			"solution_image_url":  task.SolutionImageURL,
		},
	}
	res, err := r.coll.UpdateByID(ctx, id, update)
	if err != nil {
		return fmt.Errorf("update task %s: %w", id.Hex(), err)
	}
	if res.MatchedCount == 0 {
		return ErrTaskNotFound
	}
	return nil
}

func (r *TaskRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	res, err := r.coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("delete task %s: %w", id.Hex(), err)
	}
	if res.DeletedCount == 0 {
		return ErrTaskNotFound
	}
	return nil
}

func (r *TaskRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Task, error) {
	var t models.Task
	if err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&t); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrTaskNotFound
		}
		return nil, fmt.Errorf("find task %s: %w", id.Hex(), err)
	}
	return &t, nil
}

func (r *TaskRepository) findManySorted(ctx context.Context, filter bson.M) ([]models.Task, error) {
	opts := options.Find().SetSort(bson.D{{Key: "number", Value: 1}})
	cur, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("find tasks: %w", err)
	}
	defer cur.Close(ctx)

	var out []models.Task
	if err := cur.All(ctx, &out); err != nil {
		return nil, fmt.Errorf("decode tasks: %w", err)
	}
	return out, nil
}

func buildSubjectTypeFilter(subject, taskType string) bson.M {
	filter := bson.M{}
	if subject != "" {
		filter["subject"] = subject
	}
	if taskType != "" {
		filter["task_type"] = taskType
	}
	return filter
}
