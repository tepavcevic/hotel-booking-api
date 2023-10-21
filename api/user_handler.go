package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tepavcevic/hotel-reservation/types"
)

func HandleGetUsers(c *fiber.Ctx) error {
	u := types.User{
		ID:        "23456789",
		FirstName: "James",
		LastName:  "Smith",
	}
	return c.JSON(u)
}

func HandleGetUserById(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"User": "James"})
}
