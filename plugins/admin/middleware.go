package admin

import (
	"erp6-be-golang/core/configs"
	"erp6-be-golang/core/helpers"
	"strings"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// ValidateToken validates a JWT token and returns claims
func ValidateToken(rawToken string) (*CustomClaims, error) {
	jwtSecret := []byte(configs.ConfigApps.JwtSecret)

	token, err := jwt.ParseWithClaims(rawToken, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}

// AuthMiddleware validates JWT tokens for protected routes
// This middleware can be used by any plugin
func AuthMiddleware(c *fiber.Ctx) error {
	log.Println("AuthMiddleware: Starting")
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		// Try to get from query param (for WebSocket)
		tokenQuery := c.Query("token")
		if tokenQuery != "" {
			log.Println("AuthMiddleware: Found token in query")
			authHeader = "Bearer " + tokenQuery
		} else {
			log.Println("AuthMiddleware: No token found")
			return helpers.FailResponse(c, fiber.StatusUnauthorized, "INVALID_HEADER", "")
		}
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return helpers.FailResponse(c, fiber.StatusUnauthorized, "INVALID_HEADER", "")
	}

	rawToken := strings.TrimPrefix(authHeader, "Bearer ")
	
	claims, err := ValidateToken(rawToken)
	if err != nil {
		log.Printf("AuthMiddleware: Invalid token: %v\n", err)
		return helpers.FailResponse(c, fiber.StatusUnauthorized, "INVALID_TOKEN", err.Error())
	}

	// ✅ Simpan data user di context biar bisa diakses di handler lain
	log.Printf("AuthMiddleware: Setting userid %d\n", claims.UserID)
	c.Locals("userid", int(claims.UserID))
	c.Locals("username", claims.Username)

	return c.Next()
}
