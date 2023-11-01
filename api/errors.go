package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	if apiErr, ok := err.(*fiber.Error); ok {
		return c.Status(apiErr.Code).JSON(apiErr)
	}
	defaultErr := fiber.NewError(http.StatusInternalServerError, err.Error())
	return c.Status(defaultErr.Code).JSON(defaultErr)
}
