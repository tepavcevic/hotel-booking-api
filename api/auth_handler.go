package api

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tepavcevic/hotel-reservation/db"
	"github.com/tepavcevic/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	userStore db.UserStore
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{userStore: userStore}
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User  *types.User `json:"user"`
	Token string      `json:"token"`
}

type genericResponse struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func invalidCredentials(c *fiber.Ctx) error {
	return c.Status(http.StatusBadRequest).JSON(genericResponse{
		Type:    "error",
		Message: "invalid credentials",
	})
}

func (ah *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var authParams AuthParams
	if err := c.BodyParser(&authParams); err != nil {
		return fiber.ErrBadRequest
	}
	user, err := ah.userStore.GetUserByEmail(c.Context(), authParams.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return invalidCredentials(c)
		}
		return err
	}
	if !types.IsValidPassword(user.PasswordHash, authParams.Password) {
		return invalidCredentials(c)
	}
	authResponse := AuthResponse{
		User:  user,
		Token: CreateTokenFromUser(user),
	}
	return c.JSON(authResponse)
}

func CreateTokenFromUser(user *types.User) string {
	expires := time.Now().Add(4 * time.Hour).Unix()
	claims := jwt.MapClaims{
		"userID":  user.ID,
		"email":   user.Email,
		"expires": expires,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	fmt.Println(secret)
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("create token:", err)
	}
	return tokenStr
}
