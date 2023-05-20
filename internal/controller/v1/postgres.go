package v1

import (
	"github.com/gofiber/fiber/v2"
	"icontext-test-task/internal/entity"
	"icontext-test-task/internal/interfaces"
)

type PostgresController struct {
	userService interfaces.UserService
}

func (c *PostgresController) users() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var p entity.User

		if err := ctx.BodyParser(&p); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}

		result, err := c.userService.CreateUser(ctx.Context(), p)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return ctx.JSON(result)
	}
}

func (c *PostgresController) RegisterRoutes(group fiber.Router) {
	group.Post("users", c.users())
}

func NewPostgresController(
	userService interfaces.UserService,
) *PostgresController {
	return &PostgresController{
		userService: userService,
	}
}
