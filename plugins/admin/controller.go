package admin

import (
	"erp6-be-golang/core/configs"
	gendb "erp6-be-golang/core/generator/db"
	genfile "erp6-be-golang/core/generator/file"
	"erp6-be-golang/core/helpers"
	"erp6-be-golang/models"
	"erp6-be-golang/response"
	"strconv"
	"strings"
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
func LoginHandler(c *fiber.Ctx, db *gorm.DB) error {
	var body loginRequest
	if err := c.BodyParser(&body); err != nil {
		return helpers.FailResponse(c, fiber.StatusBadRequest, "INVALID_PAYLOAD", err.Error())
	}

	var user models.Useraccess
	if err := db.Where("username = ? AND recordstatus = 1", body.Username).
		First(&user).Error; err != nil {
		return helpers.FailResponse(c, fiber.StatusUnauthorized, "INVALID_USER", err.Error())
	}

	// cek password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		return helpers.FailResponse(c, fiber.StatusUnauthorized, "INVALID_PASSWORD", err.Error())
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

	t, err := token.SignedString([]byte(configs.ConfigApps.JwtSecret))
	if err != nil {
		return helpers.FailResponse(c, fiber.StatusInternalServerError, "INVALID_TOKEN", err.Error())
	}

	return helpers.SuccessResponse(c, "SUCCESS_LOGIN", fiber.Map{
		"token": t,
		"user": fiber.Map{
			"userid":     user.Useraccessid,
			"username":   user.Username,
			"realname":   user.Realname,
			"email":      user.Email,
			"languageid": user.Languageid,
			"language":   user.Language,
			"themeid":    user.Themeid,
			"theme":      user.Theme,
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
func LogoutHandler(c *fiber.Ctx, db *gorm.DB) error {
	userID := c.Locals("userid")

	if userID == nil {
		return helpers.FailResponse(c, fiber.StatusNotFound, "INVALID_USER", "NO_USER_FOUND")
	}

	var user models.Useraccess
	if err := db.First(&user, userID).Error; err != nil {
		return helpers.FailResponse(c, fiber.StatusNotFound, "INVALID_USER", err.Error())
	}

	user.Lastlogin = time.Now()
	if err := db.Save(&user).Error; err != nil {
		return helpers.FailResponse(c, fiber.StatusInternalServerError, "INVALID_LOGOUT", err.Error())
	}

	// opsional: tambahkan ke blacklist (redis / db)
	// blacklistToken(claims["jti"].(string))

	return helpers.SuccessResponse(c, "SUCCESS_LOGOUT", nil)
}

// meHandler godoc
// @Summary Get current user
// @Description Mendapatkan user dari JWT token
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Router /auth/me [get]
func MeHandler(c *fiber.Ctx, db *gorm.DB) error {
	userID := c.Locals("userid")

	if userID == nil {
		return helpers.FailResponse(c, fiber.StatusNotFound, "INVALID_USER", "NO_USER_FOUND")
	}

	var user models.Useraccess
	if err := db.Debug().
		Preload("UserGroups.GroupAccess.GroupMenus", "isread = ?", 1).
		Preload("UserGroups.GroupAccess.GroupMenus.MenuAccess").
		First(&user, userID).Error; err != nil {
		return helpers.FailResponse(c, fiber.StatusNotFound, "INVALID_USER", "NO_USER_FOUND")
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

	return helpers.SuccessResponse(c, "SUCCESS_DATA_RETRIEVE", resp)
}

func CreateModulesHandler(c *fiber.Ctx, db *gorm.DB) error {
	menuName := c.FormValue("menu")
	pluginName := c.FormValue("plugin")
	modelName := c.FormValue("model")
	IsPermission, err := CheckUserPermission(c, db, menuName, PermWrite)
	if err != nil || !IsPermission {
		return helpers.FailResponse(c, fiber.StatusUnauthorized, "INVALID_AUTHORIZE", err.Error())
	}

	modelList := strings.Split(modelName, ",")
	err = genfile.GeneratePlugin(pluginName, modelList)
	if err != nil || !IsPermission {
		return helpers.FailResponse(c, fiber.StatusNotFound, "INVALID_PLUGIN_GENERATE", err.Error())
	}
	return helpers.SuccessResponse(c, "DATA_SAVED", nil)
}

func GenerateTableHandler(c *fiber.Ctx, db *gorm.DB) error {
	menuName := c.FormValue("menu")
	tableName := c.FormValue("table")
	IsPermission, err := CheckUserPermission(c, db, menuName, PermWrite)
	if err != nil || !IsPermission {
		return helpers.FailResponse(c, fiber.StatusUnauthorized, "INVALID_AUTHORIZE", err.Error())
	}

	if tableName == "" {
		return helpers.FailResponse(c, fiber.StatusUnauthorized, "INVALID_TABLE", "")
	}

	err = gendb.GenerateStructWithMeta(db, tableName)
	if err != nil {
		return helpers.FailResponse(c, fiber.StatusUnauthorized, "INVALID_GENERAL_TABLE", err.Error())
	}

	return helpers.SuccessResponse(c, "DATA_SAVED", nil)
}

func GenerateMultiTableHandler(c *fiber.Ctx, db *gorm.DB) error {
	menuName := c.FormValue("menu")
	tableName := c.FormValue("table")
	IsPermission, err := CheckUserPermission(c, db, menuName, PermWrite)
	if err != nil || !IsPermission {
		return helpers.FailResponse(c, fiber.StatusUnauthorized, "INVALID_AUTHORIZE", err.Error())
	}

	if tableName == "" {
		return helpers.FailResponse(c, fiber.StatusUnauthorized, "INVALID_TABLE", "")
	}

	tableList := strings.Split(tableName, ",")
	err = gendb.GenerateStructWithMultiMeta(db, tableList)
	if err != nil {
		return helpers.FailResponse(c, fiber.StatusUnauthorized, "INVALID_GENERAL_TABLE", err.Error())
	}
	return helpers.SuccessResponse(c, "DATA_SAVED", nil)
}

func ExecuteFlowHandler(c *fiber.Ctx, db *gorm.DB) {
	// TODO Implementation like Capella php version v510
}
