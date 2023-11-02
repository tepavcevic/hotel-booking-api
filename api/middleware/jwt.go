package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tepavcevic/hotel-reservation/db"
)

func JWTAuthentication(userStore db.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := c.GetReqHeaders()["Authorization"]
		if !ok {
			return fiber.ErrUnauthorized
		}
		claims, err := validateToken(token[0])
		if err != nil {
			return err
		}
		expires := claims["expires"].(float64)
		if time.Now().Unix() > int64(expires) {
			return fiber.NewError(http.StatusUnauthorized, "token expired, please sign in")
		}
		userID := claims["userID"].(string)
		user, err := userStore.GetUserById(c.Context(), userID)
		if err != nil {
			return fiber.ErrUnauthorized
		}
		c.Context().SetUserValue("user", user)
		return c.Next()
	}
}

func validateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Printf("unexpected signing method: %v", token.Header["alg"])
			return nil, fiber.ErrUnauthorized
		}
		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("failed to parse token:", err)
		return nil, fiber.ErrUnauthorized
	}
	if !token.Valid {
		return nil, fiber.ErrUnauthorized
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fiber.ErrUnauthorized
	}
	return claims, nil
}
