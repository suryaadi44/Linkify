package controller

import (
	"os"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/suryaadi44/linkify/internal/link/dto"
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
	//without auth
	l.Router.Get("/link/:username", l.GetLink)

	//with auth
	l.Router.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(global.NewBaseResponse(fiber.StatusUnauthorized, "Unauthorized"))
		},
	}))
	l.Router.Patch("/link", l.AddLink)
}

func (l *LinkController) GetLink(c *fiber.Ctx) error {
	username := c.Params("username")
	links, err := l.LinkService.GetLink(c.Context(), username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(global.NewBaseResponse(fiber.StatusInternalServerError, err.Error()))
	}

	return c.JSON(links)
}

func (l *LinkController) AddLink(c *fiber.Ctx) error {
	var link dto.Link
	if err := c.BodyParser(&link); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(global.NewBaseResponse(fiber.StatusInternalServerError, err.Error()))
	}

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	uid := claims["username"].(string)

	err := l.LinkService.AddLink(c.Context(), uid, link)
	if err != nil {
		if err == service.ErrLinkExists {
			return c.Status(fiber.StatusBadRequest).JSON(global.NewBaseResponse(fiber.StatusBadRequest, err.Error()))
		}

		return c.Status(fiber.StatusInternalServerError).JSON(global.NewBaseResponse(fiber.StatusInternalServerError, err.Error()))
	}

	return c.Status(fiber.StatusCreated).JSON(global.NewBaseResponse(fiber.StatusCreated, "Link added successfully"))
}
