package admin

import (
	"erp6-be-golang/core/configs"
	"erp6-be-golang/core/generator"
	"erp6-be-golang/core/helpers"
	"erp6-be-golang/core/scheduler"
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
	err = generator.GeneratePlugin(db, pluginName, modelList)
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

	err = generator.GenerateStructWithMeta(db, tableName)
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
	err = generator.GenerateStructWithMultiMeta(db, tableList)
	if err != nil {
		return helpers.FailResponse(c, fiber.StatusUnauthorized, "INVALID_GENERAL_TABLE", err.Error())
	}
	return helpers.SuccessResponse(c, "DATA SAVED", nil)
}

func ExecuteFlowHandler(c *fiber.Ctx, db *gorm.DB) error {
	flowName := c.FormValue("flowname")
	search := c.FormValue("search")
	menuName := c.FormValue("menu")
	debug := c.FormValue("debug") // Enable step-by-step results
	log.Info(c.FormValue("flow"))

	if flowName == "" || search == "" || menuName == "" {
		return helpers.FailResponse(c, 401, "INVALID_FLOW_REQUEST", "INVALID_FLOW_VALUE_REQUEST")
	}

	bSearch, err := strconv.ParseBool(search)
	if err != nil {
		return helpers.FailResponse(c, 401, "INVALID_CONVERSION", err.Error())
	}

	if debug == "true" {
		c.Locals("enable_workflow_events", true)
	}

	err = generator.ExecuteFlow(c, db, flowName, bSearch, nil)
	if err != nil {
		return helpers.FailResponse(c, 401, "INVALID_FLOW", err.Error())
	}

	// Auto-reload scheduler if this was a workflow modification
	if strings.Contains(flowName, "modif") && strings.Contains(flowName, "workflow") {
		// Use debounced reload to prevent multiple reloads from frontend batch requests
		scheduler.ReloadSchedulerDebounced(db)
	}

	// If debug mode, return step results
	if debug == "true" {
		wfEngine := c.Locals("wfEngine")
		if wfEngine != nil {
			return helpers.SuccessResponse(c, "FLOW_EXECUTED", fiber.Map{
				"stepResults": wfEngine,
			})
		}
	}

	return nil
}

func LoadThemeHandler(c *fiber.Ctx, db *gorm.DB) error {
	err := generator.ExecuteFlow(c, db, "searchcombotheme", true, nil)
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
	isDesign := c.Query("design") == "true"

	if userID == nil {
		return helpers.FailResponse(c, fiber.StatusNotFound, "INVALID_USER", "NO_USER_FOUND")
	}

	// --- STEP 4: Ambil hanya data menu yang dibolehkan ---
	var menus models.Menuaccess
	// Pre-fetch menu to get ID for lock check
	if err := db.Where("menuname = ?", menuName).First(&menus).Error; err != nil {
		return helpers.FailResponse(c, fiber.StatusInternalServerError, "MENU_QUERY_FAILED", err.Error())
	}

	// CHECK SCHEMA LOCK
	var lock models.Recordlock
	err := db.Where("tablename = ? AND recordid = ?", "sys_menu", menus.Menuaccessid).First(&lock).Error
	if err == nil {
		// Found lock!

		// Check Expiry (5 minutes)
		if time.Since(lock.Lockedat) > 5*time.Minute {
			// Expired! Delete it and proceed
			db.Delete(&lock)
		} else {
			currentUserID, _ := userID.(int)

			// 1. If locked by ANOTHER user -> BLOCK
			if lock.Lockedby != currentUserID {
				var lockingUser models.Useraccess
				db.Where("useraccessid = ?", lock.Lockedby).First(&lockingUser)
				return helpers.FailResponse(c, fiber.StatusLocked, "SCHEMA_LOCKED",
					fmt.Sprintf("This menu is currently being designed by %s since %s",
						lockingUser.Username, lock.Lockedat.Format("15:04")))
			}

			// 2. If locked by ME but accessing RUNTIME (not design) -> BLOCK
			if lock.Lockedby == currentUserID && !isDesign {
				return helpers.FailResponse(c, fiber.StatusLocked, "SCHEMA_LOCKED",
					"You are currently editing this menu in Form Designer.")
			}
		}
	}

	// Continue loading menu modules
	if err := db.Model(&menus).Association("Modules").Find(&menus.Modules); err != nil {
		return helpers.FailResponse(c, fiber.StatusInternalServerError, "MODULE_QUERY_FAILED", err.Error())
	}

	// CHECK PERMISSIONS
	type Permissions struct {
		IsWrite    int `json:"iswrite" gorm:"column:iswrite"`
		IsRead     int `json:"isread" gorm:"column:isread"`
		IsPurge    int `json:"ispurge" gorm:"column:ispurge"`
		IsUpload   int `json:"isupload" gorm:"column:isupload"`
		IsDownload int `json:"isdownload" gorm:"column:isdownload"`
	}
	var perms Permissions

	// Ensure userID is int
	var uid int
	if v, ok := userID.(int); ok {
		uid = v
	} else if v, ok := userID.(uint); ok {
		uid = int(v)
	}

	err = db.Table("usergroup ug").
		Select(`
			COALESCE(MAX(gm.iswrite), 0) as iswrite, 
			COALESCE(MAX(gm.isread), 0) as isread, 
			COALESCE(MAX(gm.ispurge), 0) as ispurge, 
			COALESCE(MAX(gm.isupload), 0) as isupload, 
			COALESCE(MAX(gm.isdownload), 0) as isdownload
		`).
		Joins("JOIN groupaccess ga ON ga.groupaccessid = ug.groupaccessid").
		Joins("JOIN groupmenu gm ON gm.groupaccessid = ga.groupaccessid").
		Where("ug.useraccessid = ? AND gm.menuaccessid = ? AND ga.recordstatus = 1", uid, menus.Menuaccessid).
		Scan(&perms).Error

	if err != nil {
		log.Errorf("Permission Query Failed for User %d Menu %d: %v", uid, menus.Menuaccessid, err)
	} else {
		//log.Infof("Permissions for User %d Menu %d (%s): %+v", uid, menus.Menuaccessid, menus.Menuname, perms)
	}

	// Construct Response
	response := struct {
		models.Menuaccess
		Permissions
	}{
		Menuaccess:  menus,
		Permissions: perms,
	}

	return helpers.SuccessResponse(c, "DATA RETRIEVED", response)
}

func LockRecordHandler(c *fiber.Ctx, db *gorm.DB) error {
	type LockRequest struct {
		TableName string `json:"tablename"`
		RecordID  int    `json:"recordid"`
		LockType  string `json:"locktype"`
	}
	var req LockRequest
	if err := c.BodyParser(&req); err != nil {
		return helpers.FailResponse(c, fiber.StatusBadRequest, "INVALID_PAYLOAD", err.Error())
	}

	userID := c.Locals("userid").(int)
	sessionID := c.Get("X-Session-ID", "")

	// Check existing lock
	var existingLock models.Recordlock
	err := db.Where("tablename = ? AND recordid = ?", req.TableName, req.RecordID).First(&existingLock).Error

	if err == nil {
		// Check Expiry (5 minutes)
		if time.Since(existingLock.Lockedat) > 5*time.Minute {
			// Expired! Delete it
			db.Delete(&existingLock)
		} else {
			if existingLock.Lockedby == userID {
				// Refresh
				existingLock.Lockedat = time.Now()
				db.Save(&existingLock)
				return helpers.SuccessResponse(c, "LOCK_REFRESHED", nil)
			}
			// Locked by other
			var lockingUser models.Useraccess
			db.Where("useraccessid = ?", existingLock.Lockedby).First(&lockingUser)
			return helpers.FailResponse(c, fiber.StatusConflict, "RECORD_LOCKED",
				fmt.Sprintf("Locked by %s", lockingUser.Username))
		}
	}

	// Create Lock
	lock := models.Recordlock{
		Tablename: req.TableName,
		Recordid:  req.RecordID,
		Lockedby:  userID,
		Locktype:  req.LockType,
		Sessionid: sessionID,
		Lockedat:  time.Now(),
	}
	if req.LockType == "" {
		lock.Locktype = "edit"
	}

	if err := db.Create(&lock).Error; err != nil {
		return helpers.FailResponse(c, fiber.StatusInternalServerError, "LOCK_FAILED", err.Error())
	}

	return helpers.SuccessResponse(c, "LOCK_ACQUIRED", nil)
}

func UnlockRecordHandler(c *fiber.Ctx, db *gorm.DB) error {
	type UnlockRequest struct {
		TableName string `json:"tablename"`
		RecordID  int    `json:"recordid"`
	}
	var req UnlockRequest
	if err := c.BodyParser(&req); err != nil {
		return helpers.FailResponse(c, fiber.StatusBadRequest, "INVALID_PAYLOAD", err.Error())
	}

	userID := c.Locals("userid").(int)

	result := db.Where("tablename = ? AND recordid = ? AND lockedby = ?",
		req.TableName, req.RecordID, userID).Delete(&models.Recordlock{})

	if result.Error != nil {
		return helpers.FailResponse(c, fiber.StatusInternalServerError, "UNLOCK_FAILED", result.Error.Error())
	}

	return helpers.SuccessResponse(c, "LOCK_RELEASED", nil)
}

func ExecuteTableOperationHandler(c *fiber.Ctx, db *gorm.DB) error {
	menuName := c.FormValue("menu")
	operation := c.FormValue("operation")
	tableJSON := c.FormValue("table")

	// Validate required parameters
	if menuName == "" || operation == "" || tableJSON == "" {
		return helpers.FailResponse(c, fiber.StatusBadRequest, "INVALID_REQUEST", "Missing required parameters: menu, operation, or table")
	}

	// Check permission
	IsPermission, err := CheckUserPermission(c, db, menuName, PermWrite)
	if err != nil || !IsPermission {
		return helpers.FailResponse(c, fiber.StatusUnauthorized, "INVALID_AUTHORIZE", "User does not have permission to execute table operations")
	}

	// Parse table JSON
	tableDef, err := generator.ParseTableJSON(tableJSON)
	if err != nil {
		return helpers.FailResponse(c, fiber.StatusBadRequest, "INVALID_TABLE_JSON", err.Error())
	}

	var sqlStatements []string
	var result *generator.ExecutionResult
	driver := generator.GetDatabaseDriver(db)

	// Generate SQL based on operation type
	switch strings.ToLower(operation) {
	case "create":
		sql, err := generator.GenerateCreateTableSQL(db, tableDef)
		if err != nil {
			return helpers.FailResponse(c, fiber.StatusInternalServerError, "SQL_GENERATION_FAILED", err.Error())
		}
		sqlStatements = append(sqlStatements, sql)

	case "alter":
		sqls, err := generator.GenerateAlterTableSQL(db, tableDef.Table.Name, tableDef.Table.Columns)
		if err != nil {
			return helpers.FailResponse(c, fiber.StatusInternalServerError, "SQL_GENERATION_FAILED", err.Error())
		}
		if len(sqls) == 0 {
			return helpers.SuccessResponse(c, "NO_CHANGES_DETECTED", fiber.Map{
				"message": "No schema changes detected",
			})
		}
		sqlStatements = sqls

	case "drop":
		sql := generator.GenerateDropTableSQL(tableDef.Table.Name, driver)
		sqlStatements = append(sqlStatements, sql)

	default:
		return helpers.FailResponse(c, fiber.StatusBadRequest, "INVALID_OPERATION", fmt.Sprintf("Invalid operation: %s. Must be 'create', 'alter', or 'drop'", operation))
	}

	// Execute SQL statements
	var executionResults []map[string]interface{}
	for _, sql := range sqlStatements {
		result, err = generator.ExecuteDDLStatement(db, sql)
		if err != nil {
			log.Error(fmt.Sprintf("Failed to execute SQL for table %s: %v", tableDef.Table.Name, err))
			return helpers.FailResponse(c, fiber.StatusInternalServerError, "SQL_EXECUTION_FAILED", fmt.Sprintf("Error: %s, SQL: %s", err.Error(), sql))
		}
		executionResults = append(executionResults, map[string]interface{}{
			"sql":            result.GeneratedSQL,
			"success":        result.Success,
			"execution_time": result.ExecutionTime,
		})
	}

	// Log successful execution
	log.Info(fmt.Sprintf("User executed %s operation on table %s", operation, tableDef.Table.Name))

	return helpers.SuccessResponse(c, "OPERATION_SUCCESSFUL", fiber.Map{
		"operation": operation,
		"table":     tableDef.Table.Name,
		"results":   executionResults,
	})
}
