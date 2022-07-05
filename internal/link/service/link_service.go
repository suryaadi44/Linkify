package service

import (
	"context"
	"errors"
	"log"

	"github.com/suryaadi44/linkify/internal/link/dto"
	linkRepository "github.com/suryaadi44/linkify/internal/link/repository"
	userService "github.com/suryaadi44/linkify/internal/user/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LinkService struct {
	linkRepository linkRepository.LinkRepository
	userService    userService.UserService
}

var (
	ErrLinkExists   = errors.New("Link already exists")
	ErrLinkNotFound = errors.New("Link not found")
)

func NewLinkService(linkRepository linkRepository.LinkRepository, userService userService.UserService) *LinkService {
	return &LinkService{
		linkRepository: linkRepository,
		userService:    userService,
	}
}

func (l LinkService) GetLink(ctx context.Context, username string) (*dto.LinksResponse, error) {
	links, err := l.linkRepository.GetLinkByUsername(ctx, username)
	if err != nil {
		log.Println("[Link] Error fetching links: ", err)
		return nil, err
	}

	picture, err := l.userService.GetUserPictureByUsername(ctx, username)
	if err != nil {
		log.Println("[Link] Error fetching picture: ", err)
		return nil, err
	}

	return dto.NewLinksResponse(*links, picture), nil
}

func (l LinkService) AddLink(ctx context.Context, username string, link dto.Link) error {
	// Add link to document links field (array)
	linkEntity := dto.NewLinkEntity(link)
	err := l.linkRepository.AddLink(ctx, username, *linkEntity)
	if err != nil {
		log.Println("[Link] Error adding link: ", err)
		return err
	}

	return nil
}

func (l LinkService) EditLinkById(ctx context.Context, username string, link dto.Link) error {
	//convert id string to object id
	id, err := primitive.ObjectIDFromHex(link.ID)
	if err != nil {
		log.Println("[Link] Error converting id string to object id: ", err)
		return err
	}

	//Check if link exist
	if !l.linkRepository.IsLinkIdExist(ctx, username, id) {
		return ErrLinkNotFound
	}

	// Edit link in document links field (array)
	linkEntity := dto.NewLinkEntity(link)
	linkEntity.ID = id
	err = l.linkRepository.EditLinkById(ctx, username, *linkEntity)
	if err != nil {
		log.Println("[Link] Error editing link: ", err)
		return err
	}

	return nil
}

func (l LinkService) DeleteLinkById(ctx context.Context, username string, idString string) error {
	//convert id string to object id
	id, err := primitive.ObjectIDFromHex(idString)
	if err != nil {
		log.Println("[Link] Error converting id string to object id: ", err)
		return err
	}

	//Check if link exist
	if !l.linkRepository.IsLinkIdExist(ctx, username, id) {
		return ErrLinkNotFound
	}

	// Delete link from document links field (array)
	err = l.linkRepository.DeleteLinkById(ctx, username, id)
	if err != nil {
		log.Println("[Link] Error deleting link: ", err)
		return err
	}

	return nil
}
