package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/suryaadi44/linkify/internal/link/service"

	global "github.com/suryaadi44/linkify/pkg/dto"
)

type LinkController struct {
	Router      fiber.Router
	LinkService service.LinkService
}

func NewLinkController(Router fiber.Router, LinkService service.LinkService) *LinkController {
	return &LinkController{
		Router:      Router,
		LinkService: LinkService,
	}
}

func (l *LinkController) InitializeController() {
	l.Router.Get("/link/:username", l.GetLink)
}

func (l *LinkController) GetLink(c *fiber.Ctx) error {
	username := c.Params("username")
	links, err := l.LinkService.GetLink(c.Context(), username)
	if err != nil {
		return c.Status(500).JSON(global.NewBaseResponse(500, err.Error()))
	}

	return c.JSON(links)
}
