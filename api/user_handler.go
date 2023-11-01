package api

import (
	"net/http"

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

func (uh *UserHandler) HandleCreateUser(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return fiber.ErrBadRequest
	}
	if errs := params.Validate(); len(errs) > 0 {
		return c.Status(http.StatusBadRequest).JSON(errs)
	}
	foundUser, err := uh.userStore.GetUserByEmail(c.Context(), params.Email)
	if err != nil {
		return fiber.ErrBadRequest
	}
	if foundUser != nil {
		return fiber.NewError(http.StatusConflict, "email taken")
	}
	user, err := types.NewUserFromParams(params)
	if err != nil {
		return fiber.ErrBadRequest
	}
	dbUser, err := uh.userStore.CreateUser(c.Context(), user)
	if err != nil {
		return err
	}
	return c.JSON(dbUser)
}

func (uh *UserHandler) HandleGetUserById(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := uh.userStore.GetUserById(c.Context(), id)
	if err != nil {
		return fiber.NewError(http.StatusNotFound, err.Error())
	}
	return c.JSON(user)
}

func (uh *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := uh.userStore.GetUsers(c.Context())
	if err != nil {
		return fiber.ErrNotFound
	}
	return c.JSON(users)
}

func (uh *UserHandler) HandleUpdateUser(c *fiber.Ctx) error {
	var (
		params types.UpdateUserParams
		userID = c.Params("id")
	)
	if err := c.BodyParser(&params); err != nil {
		return fiber.ErrBadRequest
	}
	filter := db.Map{"_id": userID}
	if err := uh.userStore.UpdateUser(c.Context(), filter, params); err != nil {
		return err
	}
	return c.JSON(genericResponse{
		Type:    "message",
		Message: "user successfully updated",
	})
}

func (uh *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	if err := uh.userStore.DeleteUser(c.Context(), userID); err != nil {
		return err
	}
	return c.JSON(genericResponse{
		Type:    "message",
		Message: "user successfully deleted",
	})
}
