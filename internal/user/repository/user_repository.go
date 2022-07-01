package repository

import (
	"context"

	entity "github.com/suryaadi44/linkify/internal/user/entitiy"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	db *mongo.Database
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u *UserRepository) CreateUser(ctx context.Context, user entity.User) error {
	collection := u.db.Collection("users")

	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}
