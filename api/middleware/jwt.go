package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tepavcevic/hotel-reservation/db"
)

func JWTAuthentication(userStore db.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := c.GetReqHeaders()["Authorization"]
		if !ok {
			return fmt.Errorf("unauthorized")
		}
		claims, err := validateToken(token[0])
		if err != nil {
			return err
		}
		expires := claims["expires"].(float64)
		if time.Now().Unix() > int64(expires) {
			return fmt.Errorf("token expires, please sign in")
		}
		userID := claims["userID"].(string)
		user, err := userStore.GetUserById(c.Context(), userID)
		if err != nil {
			return fmt.Errorf("unauthorized")
		}
		c.Context().SetUserValue("user", user)
		return c.Next()
	}
}

func validateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Printf("unexpected signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("unauthorized")
		}
		secret := "dsds" // os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("failed to parse token:", err)
		return nil, fmt.Errorf("unauthorized")
	}
	if !token.Valid {
		return nil, fmt.Errorf("unauthorized")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}
	return claims, nil
}
