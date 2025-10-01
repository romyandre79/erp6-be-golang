package admin

import (
	"erp6-be-golang/core/configs"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware validates JWT tokens for protected routes
// This middleware can be used by any plugin
func AuthMiddleware(c *fiber.Ctx) error {
	jwtSecret := []byte(configs.ConfigApps.JwtSecret)
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(401).JSON(fiber.Map{
			"error": "Authorization header required",
			"code":  "MISSING_TOKEN",
		})
	}

	const prefix = "Bearer "
	if len(authHeader) <= len(prefix) || authHeader[:len(prefix)] != prefix {
		return c.Status(401).JSON(fiber.Map{
			"error": "Invalid authorization header format",
			"code":  "INVALID_HEADER",
		})
	}
	rawToken := authHeader[len(prefix):]

	token, err := jwt.Parse(rawToken, func(t *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.NewError(401, "Invalid signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "Invalid or expired token",
			"code":  "INVALID_TOKEN",
		})
	}

	if !token.Valid {
		return c.Status(401).JSON(fiber.Map{
			"error": "Token is not valid",
			"code":  "INVALID_TOKEN",
		})
	}

	// Store token in context for use in handlers
	c.Locals("user", token)
	return c.Next()
}

// OptionalAuthMiddleware is an optional JWT middleware that doesn't fail if no token is provided
// Useful for routes that can work with or without authentication
func OptionalAuthMiddleware(c *fiber.Ctx) error {
	jwtSecret := []byte(configs.ConfigApps.JwtSecret)
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		// No token provided, continue without authentication
		return c.Next()
	}

	const prefix = "Bearer "
	if len(authHeader) <= len(prefix) || authHeader[:len(prefix)] != prefix {
		// Invalid format, continue without authentication
		return c.Next()
	}
	rawToken := authHeader[len(prefix):]

	token, err := jwt.Parse(rawToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.NewError(401, "Invalid signing method")
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		// Invalid token, continue without authentication
		return c.Next()
	}

	// Store token in context if valid
	c.Locals("user", token)
	return c.Next()
}

// GetUserFromContext extracts user information from the JWT token in context
func GetUserFromContext(c *fiber.Ctx) (map[string]interface{}, error) {
	userToken := c.Locals("user")
	if userToken == nil {
		return nil, fiber.NewError(401, "No user found in context")
	}

	token, ok := userToken.(*jwt.Token)
	if !ok {
		return nil, fiber.NewError(401, "Invalid user token in context")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fiber.NewError(401, "Invalid token claims")
	}

	return claims, nil
}

// RequireRole creates a middleware that checks if the user has a specific role
func RequireRole(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, err := GetUserFromContext(c)
		if err != nil {
			return err
		}

		userRole, ok := claims["role"].(string)
		if !ok || userRole != requiredRole {
			return c.Status(403).JSON(fiber.Map{
				"error": "Insufficient permissions",
				"code":  "INSUFFICIENT_ROLE",
			})
		}

		return c.Next()
	}
}

// RequireAnyRole creates a middleware that checks if the user has any of the specified roles
func RequireAnyRole(requiredRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, err := GetUserFromContext(c)
		if err != nil {
			return err
		}

		userRole, ok := claims["role"].(string)
		if !ok {
			return c.Status(403).JSON(fiber.Map{
				"error": "No role found in token",
				"code":  "NO_ROLE",
			})
		}

		for _, role := range requiredRoles {
			if userRole == role {
				return c.Next()
			}
		}

		return c.Status(403).JSON(fiber.Map{
			"error": "Insufficient permissions",
			"code":  "INSUFFICIENT_ROLE",
		})
	}
}
