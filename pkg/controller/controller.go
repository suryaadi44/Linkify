package controller

import (
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	userControllerPkg "github.com/suryaadi44/linkify/internal/user/controller"
	userRepositoryPkg "github.com/suryaadi44/linkify/internal/user/repository"
	userServicePkg "github.com/suryaadi44/linkify/internal/user/service"

	linkControllerPkg "github.com/suryaadi44/linkify/internal/link/controller"
	linkRepositoryPkg "github.com/suryaadi44/linkify/internal/link/repository"
	linkServicePkg "github.com/suryaadi44/linkify/internal/link/service"
)

func InitializeController(app *fiber.App, db *mongo.Database) {
	app.Use(cors.New())
	app.Use(favicon.New())
	app.Use(logger.New(logger.Config{
		Format:     "${time} [API] ${ip}:${port} ${status} - ${latency} ${method} ${path}\n",
		TimeFormat: "2006/01/02 15:04:05",
	}))
	app.Use(recover.New())

	api := app.Group("/api")

	userRepository := userRepositoryPkg.NewUserRepository(db)
	linkRepository := linkRepositoryPkg.NewLinkRepository(db)

	userService := userServicePkg.NewUserService(*userRepository, *linkRepository)
	linkService := linkServicePkg.NewLinkService(*linkRepository, *userService)

	userController := userControllerPkg.NewUserController(api, *userService)
	linkController := linkControllerPkg.NewLinkController(api, *linkService)

	userController.InitializeController()
	linkController.InitializeController()
}
