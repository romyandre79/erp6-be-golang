package admin

import (
	"erp6-be-golang/core/configs"
	"erp6-be-golang/core/helpers"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware validates JWT tokens for protected routes
// This middleware can be used by any plugin
func AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return helpers.FailResponse(c, fiber.StatusUnauthorized, "INVALID_HEADER", "")
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return helpers.FailResponse(c, fiber.StatusUnauthorized, "INVALID_HEADER", "")
	}

	rawToken := strings.TrimPrefix(authHeader, "Bearer ")
	jwtSecret := []byte(configs.ConfigApps.JwtSecret)

	// ✅ Parse pakai CustomClaims
	token, err := jwt.ParseWithClaims(rawToken, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, helpers.FailResponse(c, fiber.StatusUnauthorized, "INVALID_TOKEN", "")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return helpers.FailResponse(c, fiber.StatusUnauthorized, "INVALID_TOKEN", err.Error())
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return helpers.FailResponse(c, fiber.StatusUnauthorized, "INVALID_TOKEN", "")
	}

	// ✅ Simpan data user di context biar bisa diakses di handler lain
	c.Locals("userid", claims.UserID)
	c.Locals("username", claims.Username)

	return c.Next()
}
