package api

import (
	"errors"
	"fmt"
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

func (ah *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var authParams AuthParams
	if err := c.BodyParser(&authParams); err != nil {
		return err
	}
	user, err := ah.userStore.GetUserByEmail(c.Context(), authParams.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("invalid credentials")
		}
		return err
	}
	if !types.IsValidPassword(user.PasswordHash, authParams.Password) {
		return fmt.Errorf("invalid credentials")
	}
	token, err := createTokenFromUser(user)
	if err != nil {
		return fmt.Errorf("authentication: %w", err)
	}
	authResponse := AuthResponse{
		User:  user,
		Token: token,
	}
	return c.JSON(authResponse)
}

func createTokenFromUser(user *types.User) (string, error) {
	expires := time.Now().Add(4 * time.Hour).Unix()
	claims := jwt.MapClaims{
		"userID":  user.ID,
		"email":   user.Email,
		"expires": expires,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := "dsds" // os.Getenv("JWT_SECRET")
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", errors.New("failed to sign token")
	}
	return tokenStr, nil
}
