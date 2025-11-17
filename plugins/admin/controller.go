package admin

import (
	"erp6-be-golang/core/configs"
	gendb "erp6-be-golang/core/generator/db"
	genfile "erp6-be-golang/core/generator/file"
	"erp6-be-golang/core/helpers"
	"erp6-be-golang/models"
	"erp6-be-golang/response"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
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
func LoginHandler(c *fiber.Ctx, db *gorm.DB) error {
	var body loginRequest
	if err := c.BodyParser(&body); err != nil {
		return helpers.FailResponse(c, fiber.StatusBadRequest, "INVALID_PAYLOAD", err.Error())
	}

	// --- STEP 1: Ambil data user tanpa preload berat ---
	var user models.Useraccess
	if err := db.
		Joins("join theme on theme.themeid = useraccess.themeid").
		Joins("join language on language.languageid = useraccess.languageid").
		Where("useraccess.username = ? AND useraccess.recordstatus = 1", body.Username).
		First(&user).Error; err != nil {
		return helpers.FailResponse(c, fiber.StatusUnauthorized, "INVALID_USER", err.Error())
	}

	// --- STEP 2: Validasi password ---
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		return helpers.FailResponse(c, fiber.StatusUnauthorized, "INVALID_PASSWORD", err.Error())
	}

	// --- STEP 3: Update last login ---
	user.Lastlogin = time.Now()
	user.Isonline = 1
	_ = db.Save(&user)

	// --- STEP 5: Buat token JWT ---
	tokenTTL := configs.ConfigApps.JwtTtlHour
	ttl, _ := strconv.Atoi(tokenTTL)

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

	baseUrl := os.Getenv("LOCAL_BASE_URL") + "useraccess/"
	// --- STEP 6: Response ---
	return helpers.SuccessResponse(c, "SUCCESS_LOGIN", fiber.Map{
		"token": t,
		"user": fiber.Map{
			"userid":     user.Useraccessid,
			"username":   user.Username,
			"realname":   user.Realname,
			"email":      user.Email,
			"photo":      baseUrl + user.Userphoto,
			"languageid": user.Languageid,
			"language":   user.Language,
			"themeid":    user.Themeid,
			"theme":      user.Theme,
		},
	})
}

func MeHander(c *fiber.Ctx, db *gorm.DB) error {
	userID := c.Locals("userid")

	if userID == nil {
		return helpers.FailResponse(c, fiber.StatusNotFound, "INVALID_USER", "NO_USER_FOUND")
	}

	// --- STEP 4: Ambil hanya data menu yang dibolehkan ---
	var menus []models.Menuaccess
	err := db.
		Table("usergroup").
		Select(`
			DISTINCT m.menuaccessid, m.menuname, m.menucode, m.menuform, m.menutype, m.description, m.parentid, m.menuicon, m.menuurl, m.sortorder, n.moduleid, n.modulename
		`).
		Joins("JOIN groupaccess g ON g.groupaccessid = usergroup.groupaccessid").
		Joins("JOIN groupmenu gm ON gm.groupaccessid = g.groupaccessid").
		Joins("JOIN menuaccess m ON m.menuaccessid = gm.menuaccessid").
		Joins("JOIN modules n ON n.moduleid = m.moduleid").
		Where("usergroup.useraccessid = ? AND m.recordstatus = 1 AND gm.isread = 1", userID).
		Order("m.parentid asc, m.sortorder ASC").
		Scan(&menus).Error
	if err != nil {
		return helpers.FailResponse(c, fiber.StatusInternalServerError, "MENU_QUERY_FAILED", err.Error())
	}

	return helpers.SuccessResponse(c, "SUCCESS_LOGIN", fiber.Map{
		"menus": menus,
	})
}

// logoutHandler godoc
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

func CreateModulesHandler(c *fiber.Ctx, db *gorm.DB) error {
	menuName := c.FormValue("menu")
	pluginName := c.FormValue("plugin")
	modelName := c.FormValue("model")
	IsPermission, err := CheckUserPermission(c, db, menuName, PermWrite)
	if err != nil || !IsPermission {
		return helpers.FailResponse(c, fiber.StatusUnauthorized, "INVALID_AUTHORIZE", err.Error())
	}

	modelList := strings.Split(modelName, ",")
	err = genfile.GeneratePlugin(db, pluginName, modelList)
	if err != nil || !IsPermission {
		return helpers.FailResponse(c, fiber.StatusNotFound, "INVALID_PLUGIN_GENERATE", err.Error())
	}
	return helpers.SuccessResponse(c, "DATA SAVED", nil)
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

	return helpers.SuccessResponse(c, "DATA SAVED", nil)
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
	return helpers.SuccessResponse(c, "DATA SAVED", nil)
}

func ExecuteFlowHandler(c *fiber.Ctx, db *gorm.DB) error {
	flowName := c.FormValue("flow")
	search := c.FormValue("search")
	menuName := c.FormValue("menu")
	log.Info(c.FormValue("flow"))

	if flowName == "" || search == "" || menuName == "" {
		return helpers.FailResponse(c, 401, "INVALID_FLOW_REQUEST", "INVALID_FLOW_VALUE_REQUEST")
	}

	if strings.Contains(flowName, "search") {
		IsPermission, err := CheckUserPermission(c, db, menuName, PermRead)
		if err != nil || !IsPermission {
			return helpers.FailResponse(c, fiber.StatusUnauthorized, "INVALID_AUTHORIZE", err.Error())
		}
	}

	if strings.Contains(flowName, "modif") {
		IsPermission, err := CheckUserPermission(c, db, menuName, PermWrite)
		if err != nil || !IsPermission {
			return helpers.FailResponse(c, fiber.StatusUnauthorized, "INVALID_AUTHORIZE", err.Error())
		}
	}

	if strings.Contains(flowName, "purge") {
		IsPermission, err := CheckUserPermission(c, db, menuName, PermPurge)
		if err != nil || !IsPermission {
			return helpers.FailResponse(c, fiber.StatusUnauthorized, "INVALID_AUTHORIZE", err.Error())
		}
	}

	if strings.Contains(flowName, "approve") {
		IsPermission, err := CheckUserPermission(c, db, menuName, PermPost)
		if err != nil || !IsPermission {
			return helpers.FailResponse(c, fiber.StatusUnauthorized, "INVALID_AUTHORIZE", err.Error())
		}
	}

	if strings.Contains(flowName, "reject") {
		IsPermission, err := CheckUserPermission(c, db, menuName, PermReject)
		if err != nil || !IsPermission {
			return helpers.FailResponse(c, fiber.StatusUnauthorized, "INVALID_AUTHORIZE", err.Error())
		}
	}

	if strings.Contains(flowName, "down") {
		IsPermission, err := CheckUserPermission(c, db, menuName, PermDownload)
		if err != nil || !IsPermission {
			return helpers.FailResponse(c, fiber.StatusUnauthorized, "INVALID_AUTHORIZE", err.Error())
		}
	}

	bSearch, err := strconv.ParseBool(search)
	if err != nil {
		return helpers.FailResponse(c, 401, "INVALID_CONVERSION", err.Error())
	}

	err = gendb.ExecuteFlow(c, db, flowName, bSearch)
	if err != nil {
		return helpers.FailResponse(c, 401, "INVALID_FLOW", err.Error())
	}

	return nil
}

func DownTemplateHandler(c *fiber.Ctx, db *gorm.DB) error {
	menuName := c.FormValue("menu")

	if menuName == "" {
		return helpers.FailResponse(c, 401, "INVALID_FLOW_REQUEST", "INVALID_FLOW_VALUE_REQUEST")
	}

	IsPermission, err := CheckUserPermission(c, db, menuName, PermUpload)
	if err != nil || !IsPermission {
		return helpers.FailResponse(c, fiber.StatusUnauthorized, "INVALID_AUTHORIZE", err.Error())
	}

	filePath := fmt.Sprintf("public/template/%s.xlsx", menuName)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return helpers.FailResponse(c, fiber.StatusNotFound, "TEMPLATE_NOT_FOUND", fmt.Sprintf("Template untuk '%s' tidak ditemukan", menuName))
	}

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s_template.xlsx", menuName))

	return c.SendFile(filePath, true)
}

func DashboardListHandler(c *fiber.Ctx, db *gorm.DB) error {
	userID := c.Locals("userid")
	moduleName := c.Query("module")

	if userID == nil {
		return helpers.FailResponse(c, fiber.StatusNotFound, "INVALID_USER", "NO_USER_FOUND")
	}

	var menus []response.Widget
	err := db.
		Table("widget a").
		Select(`
			DISTINCT a.widgetid, a.widgetname, a.widgettitle, a.widgetversion, a.widgetform, a.widgetby, a.description, dashgroup, position, a.moduleid, modulename
		`).
		Joins("JOIN userdash b ON b.widgetid = a.widgetid").
		Joins("JOIN usergroup c ON c.groupaccessid = b.groupaccessid").
		Joins("JOIN modules d ON d.moduleid = a.moduleid").
		Where("c.useraccessid = ? AND d.modulename = ? AND a.recordstatus = 1", userID, moduleName).
		Order("b.dashgroup, b.position").
		Scan(&menus).Error
	if err != nil {
		return helpers.FailResponse(c, fiber.StatusInternalServerError, "MENU_QUERY_FAILED", err.Error())
	}

	return helpers.SuccessResponse(c, "DATA RETRIEVED", menus)
}

func DashboardSingleHandler(c *fiber.Ctx, db *gorm.DB) error {
	userID := c.Locals("userid")
	widgetname := c.Query("widgetname")

	if userID == nil {
		return helpers.FailResponse(c, fiber.StatusNotFound, "INVALID_USER", "NO_USER_FOUND")
	}

	var menus response.Widget
	err := db.
		Table("widget a").
		Select(`
			DISTINCT a.widgetid, a.widgetname, a.widgettitle, a.widgetversion, a.widgetform, a.widgetby, a.description, dashgroup, position, a.moduleid, modulename
		`).
		Joins("JOIN userdash b ON b.widgetid = a.widgetid").
		Joins("JOIN usergroup c ON c.groupaccessid = b.groupaccessid").
		Joins("JOIN modules d ON d.moduleid = a.moduleid").
		Where("c.useraccessid = ? AND a.widgetname = ? AND a.recordstatus = 1", userID, widgetname).
		Order("b.dashgroup, b.position").
		Scan(&menus).Error
	if err != nil {
		return helpers.FailResponse(c, fiber.StatusInternalServerError, "MENU_QUERY_FAILED", err.Error())
	}

	return helpers.SuccessResponse(c, "DATA RETRIEVED", menus)
}

func MenuSingleNameHandler(c *fiber.Ctx, db *gorm.DB) error {
	userID := c.Locals("userid")
	menuName := c.Query("menuname")

	if userID == nil {
		return helpers.FailResponse(c, fiber.StatusNotFound, "INVALID_USER", "NO_USER_FOUND")
	}

	// --- STEP 4: Ambil hanya data menu yang dibolehkan ---
	var menus models.Menuaccess
	err := db.
		Table("menuaccess").
		Preload("Modules").
		Where("menuaccess.menuname = ?", menuName).
		Find(&menus).Error
	if err != nil {
		return helpers.FailResponse(c, fiber.StatusInternalServerError, "MENU_QUERY_FAILED", err.Error())
	}

	return helpers.SuccessResponse(c, "DATA RETRIEVED", menus)
}
