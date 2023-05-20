package v1

import (
	"github.com/gofiber/fiber/v2"
	"icontext-test-task/internal/entity"
	"icontext-test-task/internal/interfaces"
)

type RedisController struct {
	userService interfaces.UserService
}

func (c *RedisController) incr() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var p entity.Value

		if err := ctx.BodyParser(&p); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}

		result, err := c.userService.IncrementValue(ctx.Context(), p)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return ctx.JSON(result)
	}
}

func (c *RedisController) RegisterRoutes(group fiber.Router) {
	group.Post("incr", c.incr())
}

func NewRedisController(
	userService interfaces.UserService,
) *RedisController {
	return &RedisController{
		userService: userService,
	}
}
