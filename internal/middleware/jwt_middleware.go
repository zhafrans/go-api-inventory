package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"inventory-api/internal/config"
	"inventory-api/internal/utils"
)

func JWTMiddleware(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization header",
			})
		}
		
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		
		claims, err := utils.ValidateJWT(tokenString, cfg.JWTSecret)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}
		
		c.Locals("userID", claims.UserID)
		c.Locals("userName", claims.UserName)
		c.Locals("userEmail", claims.Email)
		c.Locals("userRole", claims.Role)
		
		return c.Next()
	}
}