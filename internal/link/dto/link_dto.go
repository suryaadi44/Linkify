package dto

import "github.com/suryaadi44/linkify/internal/link/entity"

type LinksResponse struct {
	Username    string `json:"username"`
	Picture     string `json:"picture"`
	Description string `json:"description"`
	Links       []Link `json:"links,omitempty"`
}

type Link struct {
	ID        int    `json:"id,omitempty"`
	Title     string `json:"title"`
	URL       string `json:"url"`
	Thumbnail string `json:"thumbnail,omitempty"`
}

func NewLinksResponse(links entity.Links, picture string) *LinksResponse {
	var link []Link
	for idx, l := range links.Links {
		link = append(link, Link{
			ID:        idx,
			Title:     l.Title,
			URL:       l.URL,
			Thumbnail: l.Thumbnail,
		})
	}
	return &LinksResponse{
		Username:    links.Username,
		Picture:     picture,
		Description: links.Description,
		Links:       link,
	}
}

func NewLinkEntity(link Link) *entity.Link {
	return &entity.Link{
		Title:     link.Title,
		URL:       link.URL,
		Thumbnail: link.Thumbnail,
	}
}
