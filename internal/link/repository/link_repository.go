package repository

import (
	"context"

	"github.com/suryaadi44/linkify/internal/link/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type LinkRepository struct {
	db *mongo.Database
}

func NewLinkRepository(db *mongo.Database) *LinkRepository {
	return &LinkRepository{
		db: db,
	}
}

func (l *LinkRepository) GetLinkByUsername(ctx context.Context, username string) (*entity.Links, error) {
	collection := l.db.Collection("links")

	var links entity.Links
	err := collection.FindOne(ctx, bson.M{"_id": username}).Decode(&links)
	if err != nil {
		return nil, err
	}

	return &links, err
}

func (l *LinkRepository) CreateDefaultLink(ctx context.Context, username string) error {
	collection := l.db.Collection("links")

	_, err := collection.InsertOne(ctx, entity.Links{
		Username:    username,
		Description: "",
	})

	if err != nil {
		return err
	}

	return nil
}

func (l *LinkRepository) AddLink(ctx context.Context, username string, link entity.Link) error {
	collection := l.db.Collection("links")

	// add link to document links field (array)
	_, err := collection.UpdateOne(ctx, bson.M{"_id": username}, bson.M{"$push": bson.M{"links": link}})
	if err != nil {
		return err
	}

	return nil
}

func (l *LinkRepository) IsLinkExists(ctx context.Context, username string, linkTitle string) bool {
	collection := l.db.Collection("links")

	// check if link exists in document links field (array)
	var links entity.Links
	err := collection.FindOne(ctx, bson.M{"_id": username, "links.title": linkTitle}).Decode(&links)
	return err == nil
}
