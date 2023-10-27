package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/tepavcevic/hotel-reservation/types"
)

func getAuthUser(ctx *fiber.Ctx) (*types.User, error) {
	user, ok := ctx.Context().UserValue("user").(*types.User)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}
	return user, nil
}
