package repository

import (
	"context"
	"fmt"

	"github.com/suryaadi44/linkify/internal/link/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	link.ID = primitive.NewObjectID()
	_, err := collection.UpdateOne(ctx, bson.M{"_id": username}, bson.M{"$push": bson.M{"links": link}})
	if err != nil {
		return err
	}

	return nil
}

func (l *LinkRepository) EditLinkById(ctx context.Context, username string, link entity.Link) error {
	collection := l.db.Collection("links")

	// edit link in document links field with specified id
	debug, err := collection.UpdateOne(ctx, bson.M{"_id": username, "links.id": link.ID}, bson.M{"$set": bson.M{"links.$": link}})
	if err != nil {
		return err
	}
	fmt.Println(debug)

	return nil
}

func (l *LinkRepository) IsLinkIdExist(ctx context.Context, username string, id primitive.ObjectID) bool {
	collection := l.db.Collection("links")

	// check if _id exist in document links field (array)
	count, err := collection.CountDocuments(ctx, bson.M{"_id": username, "links.id": id})
	if err != nil {
		return false
	}

	return count > 0
}

func (l *LinkRepository) DeleteLinkById(ctx context.Context, username string, id primitive.ObjectID) error {
	collection := l.db.Collection("links")

	//delete link with specified id from document links field (array)
	_, err := collection.UpdateOne(ctx, bson.M{"_id": username}, bson.M{"$pull": bson.M{"links": bson.M{"id": id}}})
	if err != nil {
		return err
	}

	return nil
}
