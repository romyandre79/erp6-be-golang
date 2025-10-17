package admin

import (
	"erp6-be-golang/core/configs"
	"erp6-be-golang/core/helpers"
	"erp6-be-golang/models"
	"erp6-be-golang/response"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type CustomClaims struct {
	UserID   uint   `json:"userid"`
	Username string `json:"username"`
	jwt.RegisteredClaims
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

	var user models.Useraccess
	if err := db.Where("username = ? AND recordstatus = 1", body.Username).
		First(&user).Error; err != nil {
		return helpers.FailResponse(c, fiber.StatusUnauthorized, "user not found or inactive", err.Error())
	}

	// cek password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		return helpers.FailResponse(c, fiber.StatusUnauthorized, "invalid credentials", "password mismatch")
	}

	// update last login
	user.Lastlogin = time.Now()
	user.Isonline = 1
	_ = db.Save(&user)

	// TTL dari env
	tokenTTL := configs.ConfigApps.JwtTtlHour
	ttl, _ := strconv.Atoi(tokenTTL)

	// ✅ Gunakan RegisteredClaims (bukan MapClaims)
	claims := CustomClaims{
		UserID:   uint(user.Useraccessid),
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(ttl))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	fmt.Printf("Claims: %+v\n", token.Claims)

	t, err := token.SignedString([]byte(configs.ConfigApps.JwtSecret))
	if err != nil {
		return helpers.FailResponse(c, fiber.StatusInternalServerError, "failed to generate token", err.Error())
	}

	return helpers.SuccessResponse(c, "login success", fiber.Map{
		"token": t,
		"user": fiber.Map{
			"userid":   user.Useraccessid,
			"username": user.Username,
			"realname": user.Realname,
			"email":    user.Email,
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
	userID := c.Locals("userid")

	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "No user found in context",
		})
	}

	var user models.Useraccess
	if err := db.First(&user, userID).Error; err != nil {
		return helpers.FailResponse(c, fiber.StatusNotFound, "user not found", err.Error())
	}

	user.Lastlogin = time.Now()
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
	userID := c.Locals("userid")

	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "No user found in context",
		})
	}

	var user models.Useraccess
	if err := db.Debug().
		Preload("UserGroups.GroupAccess.GroupMenus", "isread = ?", 1).
		Preload("UserGroups.GroupAccess.GroupMenus.MenuAccess").
		First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	resp := response.UserResponse{
		UserAccessID: uint(user.Useraccessid),
		Username:     user.Username,
		RealName:     user.Realname,
		Email:        user.Email,
	}

	for _, g := range user.Usergroups {
		resp.Groups = append(resp.Groups, response.GroupResponse{
			GroupAccessID: g.Groupaccess.Groupaccessid,
			GroupName:     g.Groupaccess.Groupname,
			Description:   g.Groupaccess.Description,
		})
		for _, h := range g.Groupaccess.Groupmenus {
			resp.Menus = append(resp.Menus, response.GroupMenuResponse{
				MenuAccessID: h.Menuaccess.Menuaccessid,
				MenuName:     h.Menuaccess.Menuname,
				ParentId:     &h.Menuaccess.Parentid,
				IsRead:       h.Isread,
				IsWrite:      h.Iswrite,
				IsPost:       h.Ispost,
				IsReject:     h.Isreject,
				IsUpload:     h.Isupload,
				IsDownload:   h.Isdownload,
				IsPurge:      h.Ispurge,
			})
		}
	}

	return helpers.SuccessResponse(c, "DATA_RETRIEVE", resp)
}
