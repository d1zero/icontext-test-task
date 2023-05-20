package v1

import (
	"github.com/gofiber/fiber/v2"
	"icontext-test-task/internal/entity"
	"icontext-test-task/internal/interfaces"
)

type SignController struct {
	userService interfaces.UserService
}

func (c *SignController) sign() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var p entity.Sign

		if err := ctx.BodyParser(&p); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}

		result, err := c.userService.SignBody(p)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return ctx.JSON(result)
	}
}

func (c *SignController) RegisterRoutes(group fiber.Router) {
	group.Post("hmacsha512", c.sign())
}

func NewSignController(
	userService interfaces.UserService,
) *SignController {
	return &SignController{
		userService: userService,
	}
}
