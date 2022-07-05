package dto

import (
	"github.com/suryaadi44/linkify/internal/link/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LinksResponse struct {
	Username    string `json:"username"`
	Picture     string `json:"picture"`
	Description string `json:"description"`
	Link        []Link `json:"link,omitempty"`
}

type Link struct {
	ID        primitive.ObjectID `json:"id"`
	Title     string             `json:"title"`
	URL       string             `json:"url"`
	Thumbnail string             `json:"thumbnail,omitempty"`
}

func NewLinksResponse(links entity.Links, picture string) *LinksResponse {
	var link []Link
	for _, l := range links.Link {
		link = append(link, Link{
			ID:        l.ID,
			Title:     l.Title,
			URL:       l.URL,
			Thumbnail: l.Thumbnail,
		})
	}
	return &LinksResponse{
		Username:    links.Username,
		Picture:     picture,
		Description: links.Description,
		Link:        link,
	}
}
