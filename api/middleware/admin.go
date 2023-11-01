package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tepavcevic/hotel-reservation/types"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return fiber.ErrUnauthorized
	}
	if !user.IsAdmin {
		return fiber.ErrUnauthorized
	}
	return c.Next()
}
