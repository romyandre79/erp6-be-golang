package admin

import (
	"erp6-be-golang/core/configs"
	"erp6-be-golang/core/helpers"
	"erp6-be-golang/models"
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
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		// Try to get from query param (for WebSocket)
		tokenQuery := c.Query("token")
		if tokenQuery != "" {
			authHeader = "Bearer " + tokenQuery
		} else {
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

		// Fix: If token is invalid/expired, try to get UserID and set isonline=0
		// We parse WITHOUT validation to get the claims
		token, _ := jwt.ParseWithClaims(rawToken, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(configs.ConfigApps.JwtSecret), nil
		})
		
		if token != nil {
			if claims, ok := token.Claims.(*CustomClaims); ok {
				if GlobalDB != nil && claims.UserID > 0 {
					log.Printf("AuthMiddleware: Auto Offline for UserID %d due to invalid token\n", claims.UserID)
					GlobalDB.Model(&models.Useraccess{}).
						Where("useraccessid = ?", claims.UserID).
						Update("isonline", 0)
				}
			}
		}

		return helpers.FailResponse(c, fiber.StatusUnauthorized, "INVALID_TOKEN", err.Error())
	}

	// ✅ Simpan data user di context biar bisa diakses di handler lain
	log.Printf("AuthMiddleware: Setting userid %d\n", claims.UserID)
	c.Locals("userid", int(claims.UserID))
	c.Locals("username", claims.Username)

	return c.Next()
}
