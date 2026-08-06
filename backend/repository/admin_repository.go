package repository

import (
	"context"
	"errors"
	"fmt"

	"exam-tasks-backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var ErrAdminNotFound = errors.New("admin not found")

type AdminRepository struct {
	coll *mongo.Collection
}

func NewAdminRepository(db *mongo.Database) *AdminRepository {
	return &AdminRepository{coll: db.Collection("admins")}
}

func (r *AdminRepository) FindByLogin(ctx context.Context, login string) (*models.Admin, error) {
	var a models.Admin
	if err := r.coll.FindOne(ctx, bson.M{"login": login}).Decode(&a); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrAdminNotFound
		}
		return nil, fmt.Errorf("find admin %q: %w", login, err)
	}
	return &a, nil
}

func (r *AdminRepository) Count(ctx context.Context) (int64, error) {
	n, err := r.coll.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, fmt.Errorf("count admins: %w", err)
	}
	return n, nil
}

func (r *AdminRepository) Create(ctx context.Context, admin *models.Admin) error {
	if _, err := r.coll.InsertOne(ctx, admin); err != nil {
		return fmt.Errorf("insert admin: %w", err)
	}
	return nil
}
