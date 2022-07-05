package repository

import (
	"context"
	"fmt"

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

func (l *LinkRepository) EditLinkByIndex(ctx context.Context, username string, index int, link entity.Link) error {
	collection := l.db.Collection("links")

	// edit link in document links field at specified index
	_, err := collection.UpdateOne(ctx, bson.M{"_id": username}, bson.M{"$set": bson.M{fmt.Sprintf("links.%d", index): link}})
	if err != nil {
		return err
	}

	return nil
}

func (l *LinkRepository) IsLinkIndexExist(ctx context.Context, username string, index int) bool {
	collection := l.db.Collection("links")

	// Agregate links with projection to check if index exist
	matchStage := bson.D{
		{"$match", bson.D{
			{"_id", username},
		}},
	}

	projectionStage := bson.D{
		{"$project", bson.D{
			{"links", bson.D{
				{"$arrayElemAt", bson.A{"$links", index}},
			}},
		}},
	}

	indexCursor, err := collection.Aggregate(ctx, mongo.Pipeline{matchStage, projectionStage})
	if err != nil {
		return false
	}

	// decode cursor to check if index exist
	var result []bson.M
	err = indexCursor.All(ctx, &result)
	if err != nil {
		return false
	}

	if _, ok := result[0]["links"]; !ok {
		return false
	}

	return true
}
