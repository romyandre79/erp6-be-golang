// admin/auth.go
package admin

import (
	"erp6-be-golang/core/configs"
	"erp6-be-golang/core/helpers"
	adminModels "erp6-be-golang/plugins/admin/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func RegisterRoutes(app *fiber.App, db *gorm.DB) {
	r := app.Group("/auth")

	// passing db via closure
	r.Post("/login", func(c *fiber.Ctx) error { return loginHandler(c, db) })
	r.Post("/logout", func(c *fiber.Ctx) error { return logoutHandler(c, db) })

	protected := r.Group("/admin")
	protected.Use(AuthMiddleware)
	protected.Get("/me", func(c *fiber.Ctx) error { return meHandler(c, db) })
}

// loginHandler godoc
// @Summary Login
// @Description Login user dengan username & password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body loginRequest true "Login Request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
func loginHandler(c *fiber.Ctx, db *gorm.DB) error {
	var body loginRequest
	if err := c.BodyParser(&body); err != nil {
		return helpers.FailResponse(c, fiber.StatusBadRequest, "invalid payload", err.Error())
	}

	var user adminModels.UserAccess
	if err := db.Where("username = ? AND recordstatus = 1", body.Username).
		First(&user).Error; err != nil {
		return helpers.FailResponse(c, fiber.StatusUnauthorized, "user not found or inactive", err.Error())
	}

	// cek password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		return helpers.FailResponse(c, fiber.StatusUnauthorized, "invalid credentials", "password mismatch")
	}

	// update last login
	user.LastLogin = time.Now()
	user.IsOnline = true
	_ = db.Save(&user)

	// generate JWT
	claims := jwt.MapClaims{
		"userid":   user.UserAccessID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(configs.ConfigApps.JwtSecret))
	if err != nil {
		return helpers.FailResponse(c, fiber.StatusInternalServerError, "failed to generate token", err.Error())
	}

	if err := db.Preload("Groups.GroupAccess.GroupMenus.MenuAccess").
		First(&user, claims["userid"]).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "user not found"})
	}

	return helpers.SuccessResponse(c, "login success", fiber.Map{
		"token": t,
		"user": fiber.Map{
			"userid":   user.UserAccessID,
			"username": user.Username,
			"realname": user.RealName,
			"email":    user.Email,
			"groups":   user.Groups,
		},
	})
}

// logoutHandler godoc
// @Summary Logout
// @Description Logout user (dummy endpoint, nanti bisa update isonline)
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]string
// @Router /auth/logout [post]
func logoutHandler(c *fiber.Ctx, db *gorm.DB) error {
	claims, err := GetUserFromContext(c)
	if err != nil {
		return helpers.FailResponse(c, fiber.StatusUnauthorized, "unauthorized", err.Error())
	}

	var user adminModels.UserAccess
	if err := db.First(&user, claims["userid"]).Error; err != nil {
		return helpers.FailResponse(c, fiber.StatusNotFound, "user not found", err.Error())
	}

	user.LastLogin = time.Now()
	if err := db.Save(&user).Error; err != nil {
		return helpers.FailResponse(c, fiber.StatusInternalServerError, "failed to update logout time", err.Error())
	}

	// opsional: tambahkan ke blacklist (redis / db)
	// blacklistToken(claims["jti"].(string))

	return helpers.SuccessResponse(c, "logout success", nil)
}

// meHandler godoc
// @Summary Get current user
// @Description Mendapatkan user dari JWT token
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Router /auth/me [get]
func meHandler(c *fiber.Ctx, db *gorm.DB) error {
	claims, err := GetUserFromContext(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	var user adminModels.UserAccess
	if err := db.Preload("Groups.GroupAccess.GroupMenus.MenuAccess").
		First(&user, claims["userid"]).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "user not found"})
	}

	return c.JSON(fiber.Map{
		"userid":   user.UserAccessID,
		"username": user.Username,
		"realname": user.RealName,
		"email":    user.Email,
		"groups":   user.Groups,
	})
}
