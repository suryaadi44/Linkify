package controller

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/suryaadi44/linkify/internal/user/dto"
	"github.com/suryaadi44/linkify/internal/user/service"

	global "github.com/suryaadi44/linkify/pkg/dto"
)

type UserController struct {
	app         *fiber.App
	UserService service.UserService
}

func NewUserController(app *fiber.App, userService service.UserService) *UserController {
	return &UserController{
		app:         app,
		UserService: userService,
	}
}

func (u *UserController) InitializeController() {
	u.app.Post("/user/register", u.RegisterUser)
}

func (u *UserController) RegisterUser(c *fiber.Ctx) error {
	var user dto.RegisterForm
	if err := c.BodyParser(&user); err != nil {
		return c.Status(500).JSON(global.NewBaseResponse(500, err.Error()))
	}

	err := u.UserService.CreateUser(c.Context(), user)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key error") {
			return c.Status(400).JSON(global.NewBaseResponse(400, "User already exists"))
		}

		return c.Status(500).JSON(global.NewBaseResponse(500, err.Error()))
	}
	return c.Status(201).JSON(global.NewBaseResponse(201, "User created successfully"))
}
