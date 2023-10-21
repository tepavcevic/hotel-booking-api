package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tepavcevic/hotel-reservation/db"
	"github.com/tepavcevic/hotel-reservation/types"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (uh *UserHandler) HandleGetUserById(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := uh.userStore.GetUserById(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(user)
}

func (uh *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	u := types.User{
		ID:        "23456789",
		FirstName: "James",
		LastName:  "Smith",
	}
	return c.JSON(u)
}
