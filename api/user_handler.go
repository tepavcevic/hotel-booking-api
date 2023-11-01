package api

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/tepavcevic/hotel-reservation/db"
	"github.com/tepavcevic/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fiber.NewError(404, "user not found")
		}
		return err
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
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fiber.NewError(400, "invalid user id")
	}
	if err := c.BodyParser(&params); err != nil {
		return fiber.ErrBadRequest
	}
	filter := bson.M{"_id": oid}
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
