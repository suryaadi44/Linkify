package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Links struct {
	Username    string `bson:"_id"`
	Description string `bson:"description,omitempty"`
	Links       []Link `bson:"links,omitempty"`
}

type Link struct {
	ID        primitive.ObjectID `bson:"id"`
	Title     string             `bson:"title"`
	URL       string             `bson:"url"`
	Thumbnail string             `bson:"thumbnail,omitempty"`
}
