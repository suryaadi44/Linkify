package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/suryaadi44/linkify/internal/user/dto"
	"github.com/suryaadi44/linkify/internal/user/service"

	global "github.com/suryaadi44/linkify/pkg/dto"
)

type UserController struct {
	Router      fiber.Router
	UserService service.UserService
}

func NewUserController(Router fiber.Router, userService service.UserService) *UserController {
	return &UserController{
		Router:      Router,
		UserService: userService,
	}
}

func (u *UserController) InitializeController() {
	u.Router.Post("/user/register", u.RegisterUser)
}

func (u *UserController) RegisterUser(c *fiber.Ctx) error {
	var user dto.RegisterForm
	if err := c.BodyParser(&user); err != nil {
		return c.Status(500).JSON(global.NewBaseResponse(500, err.Error()))
	}

	if exists := u.UserService.IsEmailExists(c.Context(), user.Email); exists {
		return c.Status(400).JSON(global.NewBaseResponse(400, "Email already registered"))
	}

	if exists := u.UserService.IsUsernameExists(c.Context(), user.Username); exists {
		return c.Status(400).JSON(global.NewBaseResponse(400, "Username already registered"))
	}

	err := u.UserService.CreateUser(c.Context(), user)
	if err != nil {
		return c.Status(500).JSON(global.NewBaseResponse(500, err.Error()))
	}
	return c.Status(201).JSON(global.NewBaseResponse(201, "User created successfully"))
}
