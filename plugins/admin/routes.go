package admin

import (
	"erp6-be-golang/core/configs"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func RegisterRoutes(app *fiber.App) {
	r := app.Group("/auth")

	r.Post("/login", loginHandler)
	r.Post("/logout", logoutHandler)

	protected := r.Group("/api")
	protected.Use(AuthMiddleware)
	protected.Get("/me", meHandler)
}

func loginHandler(c *fiber.Ctx) error {
	var body loginRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid payload"})
	}

	// sementara hardcode
	if body.Username != "admin" || body.Password != "1234" {
		return c.Status(401).JSON(fiber.Map{"error": "invalid credentials"})
	}

	claims := jwt.MapClaims{
		"username": body.Username,
		"role":     "admin",
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(configs.ConfigApps.JwtSecret))
	if err != nil {
		return c.SendStatus(500)
	}

	return c.JSON(fiber.Map{"token": t})
}

func logoutHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "logged out"})
}

func meHandler(c *fiber.Ctx) error {
	claims, err := GetUserFromContext(c)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"user": claims["username"],
		"role": claims["role"],
	})
}
