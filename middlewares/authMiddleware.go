package middlewares

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("your_secret_key") // Replace with a secure secret key

// Middleware to verify JWT token
func AuthMiddleware(c *fiber.Ctx) error {
    // Get the Authorization header
    authHeader := c.Get("Authorization")

    if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "message": "Missing or malformed token",
        })
    }

    tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

    // Parse and validate the JWT token

    token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
        return secretKey, nil
    })

    if err != nil || !token.Valid {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "message": "Invalid or expired token",
        })
    }

	claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "message": "Invalid token claims",
        })
    }

	userID := claims["user_id"]
    if userID != nil {
        c.Locals("user", map[string]interface{}{"user_id": userID})
    }
    return c.Next()
}
