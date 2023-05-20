package v1

import (
	"github.com/gofiber/fiber/v2"
)

func HandleError() fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		return ctx.Status(fiber.StatusInternalServerError).JSON(newErrResp(err))
	}
}
