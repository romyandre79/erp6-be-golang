package admin

import (
	"archive/zip"
	"encoding/json"
	"erp6-be-golang/core/helpers"
	"erp6-be-golang/models"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ModuleManifest represents the module.json structure
type ModuleManifest struct {
	Modulename    string                 `json:"modulename"`
	Description   string                 `json:"description"`
	Moduleversion string                 `json:"moduleversion"`
	Createdby     string                 `json:"createdby"`
	Themeid       int                    `json:"themeid"`
	Dependencies  []int                  `json:"dependencies"`
	Menuaccess    []MenuAccessDefinition `json:"menuaccess"`
	Widgets       []WidgetDefinition     `json:"widgets"`
	Workflows     []WorkflowDefinition   `json:"workflows"`
	Tables        []TableDefinition      `json:"tables"`
	GroupMenus    []GroupMenuDefinition  `json:"groupmenus"`
}


type MenuAccessDefinition struct {
	Menuname       string `json:"menuname"`
	Menucode       string `json:"menucode"`
	Description    string `json:"description"`
	Parentid       *int   `json:"parentid"`
	Menuurl        string `json:"menuurl"`
	Sortorder      int    `json:"sortorder"`
	Menuicon       string `json:"menuicon"`
	Menutype       string `json:"menutype"`
	Menuform       string `json:"menuform"`
	Menuformdetail string `json:"menuformdetail"`
	Recordstatus   int    `json:"recordstatus"`
}

type WidgetDefinition struct {
	Widgetname    string `json:"widgetname"`
	Widgettitle   string `json:"widgettitle"`
	Widgetversion string `json:"widgetversion"`
	Widgetby      string `json:"widgetby"`
	Description   string `json:"description"`
	Widgetform    string `json:"widgetform"`
	Recordstatus  int    `json:"recordstatus"`
}

type WorkflowDefinition struct {
	Wfname       string `json:"wfname"`
	Wfdesc       string `json:"wfdesc"`
	Wfminstat    int8   `json:"wfminstat"`
	Wfmaxstat    int8   `json:"wfmaxstat"`
	Flow         string `json:"flow"`
	Moduleid     int    `json:"moduleid"`
	Recordstatus int8   `json:"recordstatus"`
}

type TableDefinition struct {
	TableName   string                 `json:"tableName"`
	Columns     []ColumnDefinition     `json:"columns"`
	Indexes     []IndexDefinition      `json:"indexes"`
	ForeignKeys []ForeignKeyDefinition `json:"foreignKeys"`
	Data        []map[string]interface{} `json:"data,omitempty"`
}

type ColumnDefinition struct {
	Name          string      `json:"name"`
	Type          string      `json:"type"`
	Length        int         `json:"length"`
	PrimaryKey    bool        `json:"primaryKey"`
	AutoIncrement bool        `json:"autoIncrement"`
	Nullable      bool        `json:"nullable"`
	Unique        bool        `json:"unique"`
	Default       interface{} `json:"default"`
}

type IndexDefinition struct {
	Name    string   `json:"name"`
	Columns []string `json:"columns"`
}

type ForeignKeyDefinition struct {
	Name       string `json:"name"`
	Column     string `json:"column"`
	References struct {
		Table  string `json:"table"`
		Column string `json:"column"`
	} `json:"references"`
	OnDelete string `json:"onDelete"`
	OnUpdate string `json:"onUpdate"`
}

type GroupMenuDefinition struct {
	Menuname   string `json:"menuname"`
	Isread     int    `json:"isread"`
	Iswrite    int    `json:"iswrite"`
	Ispost     int    `json:"ispost"`
	Isreject   int    `json:"isreject"`
	Isupload   int    `json:"isupload"`
	Isdownload int    `json:"isdownload"`
	Ispurge    int    `json:"ispurge"`
}


// UploadModulePackageHandler handles module package upload (ZIP file)
func UploadModulePackageHandler(c *fiber.Ctx, db *gorm.DB) error {
	menuName := c.FormValue("menu")

	// Check permission
	IsPermission, err := CheckUserPermission(c, db, menuName, PermUpload)
	if err != nil || !IsPermission {
		return helpers.FailResponse(c, fiber.StatusUnauthorized, "INVALID_AUTHORIZE", err.Error())
	}

	// Get the uploaded file
	file, err := c.FormFile("module")
	if err != nil {
		return helpers.FailResponse(c, fiber.StatusBadRequest, "MISSING_FILE", "Module ZIP file is required")
	}

	// Save to temp directory
	tempDir := "./tmp/modules"
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return helpers.FailResponse(c, fiber.StatusInternalServerError, "CREATE_TEMP_DIR_FAILED", err.Error())
	}

	tempPath := filepath.Join(tempDir, file.Filename)
	if err := c.SaveFile(file, tempPath); err != nil {
		return helpers.FailResponse(c, fiber.StatusInternalServerError, "SAVE_FILE_FAILED", err.Error())
	}
	defer os.Remove(tempPath) // Clean up ZIP file

	// Open and read ZIP file
	r, err := zip.OpenReader(tempPath)
	if err != nil {
		return helpers.FailResponse(c, fiber.StatusBadRequest, "INVALID_ZIP", "Could not open ZIP file")
	}
	defer r.Close()

	// Find and parse module.json
	var manifest ModuleManifest
	foundManifest := false

	for _, f := range r.File {
		if f.Name == "module.json" {
			rc, err := f.Open()
			if err != nil {
				return helpers.FailResponse(c, fiber.StatusInternalServerError, "READ_MANIFEST_FAILED", err.Error())
			}
			defer rc.Close()

			content, err := io.ReadAll(rc)
			if err != nil {
				return helpers.FailResponse(c, fiber.StatusInternalServerError, "READ_MANIFEST_FAILED", err.Error())
			}

			if err := json.Unmarshal(content, &manifest); err != nil {
				return helpers.FailResponse(c, fiber.StatusBadRequest, "INVALID_MANIFEST", "Invalid module.json format: "+err.Error())
			}

			foundManifest = true
			break
		}
	}

	if !foundManifest {
		return helpers.FailResponse(c, fiber.StatusBadRequest, "MISSING_MANIFEST", "module.json not found in ZIP")
	}

	// Call InstallModuleHandler with parsed manifest
	return installModuleFromManifest(c, db, manifest)
}

// installModuleFromManifest installs a module from parsed manifest
func installModuleFromManifest(c *fiber.Ctx, db *gorm.DB, manifest ModuleManifest) error {
	// Start transaction
	tx := db.Begin()
	if tx.Error != nil {
		return helpers.FailResponse(c, fiber.StatusInternalServerError, "TRANSACTION_START_FAILED", tx.Error.Error())
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Validate dependencies exist
	if len(manifest.Dependencies) > 0 {
		var count int64
		tx.Model(&models.Modules{}).Where("moduleid IN ?", manifest.Dependencies).Count(&count)
		if int(count) != len(manifest.Dependencies) {
			tx.Rollback()
			return helpers.FailResponse(c, fiber.StatusBadRequest, "INVALID_DEPENDENCIES", "One or more dependencies do not exist")
		}
	}

	// 2. Check if module already exists
	var existingModule models.Modules
	err := tx.Where("modulename = ?", manifest.Modulename).First(&existingModule).Error
	if err == nil {
		tx.Rollback()
		return helpers.FailResponse(c, fiber.StatusConflict, "MODULE_EXISTS", "Module with this name already exists")
	}

	// 3. Insert module record
	newModule := models.Modules{
		Modulename:    manifest.Modulename,
		Description:   manifest.Description,
		Createdby:     manifest.Createdby,
		Moduleversion: manifest.Moduleversion,
		Installdate:   time.Now(),
		Themeid:       manifest.Themeid,
		Recordstatus:  1,
	}

	if err := tx.Create(&newModule).Error; err != nil {
		tx.Rollback()
		return helpers.FailResponse(c, fiber.StatusInternalServerError, "CREATE_MODULE_FAILED", err.Error())
	}

	moduleID := newModule.Moduleid

	// 4. Insert dependency records
	for _, depID := range manifest.Dependencies {
		moduleRel := models.Modulerelation{
			Moduleid:   moduleID,
			Relationid: depID,
		}
		if err := tx.Create(&moduleRel).Error; err != nil {
			tx.Rollback()
			return helpers.FailResponse(c, fiber.StatusInternalServerError, "CREATE_DEPENDENCY_FAILED", err.Error())
		}
	}

	// 5. Insert menu access entries
	// 5. Insert menu access entries
	menuMap := make(map[string]int)
	for _, menuDef := range manifest.Menuaccess {
		menuAccess := models.Menuaccess{
			Menuname:       menuDef.Menuname,
			Menucode:       menuDef.Menucode,
			Description:    menuDef.Description,
			Moduleid:       moduleID,
			Parentid:       menuDef.Parentid,
			Menuurl:        menuDef.Menuurl,
			Sortorder:      menuDef.Sortorder,
			Menuicon:       menuDef.Menuicon,
			Menutype:       menuDef.Menutype,
			Menuform:       menuDef.Menuform,
			Menuformdetail: menuDef.Menuformdetail,
			Recordstatus:   menuDef.Recordstatus,
			Updatedate:     time.Now(),
		}
		if err := tx.Create(&menuAccess).Error; err != nil {
			tx.Rollback()
			return helpers.FailResponse(c, fiber.StatusInternalServerError, "CREATE_MENU_FAILED", err.Error())
		}
		menuMap[menuAccess.Menuname] = menuAccess.Menuaccessid
	}

	// 5b. Insert GroupMenu for groupaccessid = 2
	for _, gmDef := range manifest.GroupMenus {
		if menuID, ok := menuMap[gmDef.Menuname]; ok {
			groupMenu := models.Groupmenu{
				Groupaccessid: 2, // Hardcoded as per requirement
				Menuaccessid:  menuID,
				Isread:        gmDef.Isread,
				Iswrite:       gmDef.Iswrite,
				Ispost:        gmDef.Ispost,
				Isreject:      gmDef.Isreject,
				Isupload:      gmDef.Isupload,
				Isdownload:    gmDef.Isdownload,
				Ispurge:       gmDef.Ispurge,
				Updatedate:    time.Now(),
			}

			// Check existence first to be safe, though usually new module means new menus
			// But skipping check for simplicity in bulk insert scenario, usually safer to just create
			if err := tx.Create(&groupMenu).Error; err != nil {
				// Log error but continue? Or fail? Fail is safer to ensure consistency
				tx.Rollback()
				return helpers.FailResponse(c, fiber.StatusInternalServerError, "CREATE_GROUPMENU_FAILED", err.Error())
			}
		}
	}

	// 6. Insert widgets
	for _, widgetDef := range manifest.Widgets {
		widget := models.Widget{
			Widgetname:    widgetDef.Widgetname,
			Widgettitle:   widgetDef.Widgettitle,
			Widgetversion: widgetDef.Widgetversion,
			Widgetby:      widgetDef.Widgetby,
			Description:   widgetDef.Description,
			Widgetform:    widgetDef.Widgetform,
			Moduleid:      moduleID,
			Recordstatus:  widgetDef.Recordstatus,
			Createdate:    time.Now(),
			Updatedate:    time.Now(),
		}
		if err := tx.Create(&widget).Error; err != nil {
			tx.Rollback()
			return helpers.FailResponse(c, fiber.StatusInternalServerError, "CREATE_WIDGET_FAILED", err.Error())
		}
	}

	// 7. Insert workflows and track them
	for _, wfDef := range manifest.Workflows {
		workflow := models.Workflow{
			Wfname:       wfDef.Wfname,
			Wfdesc:       wfDef.Wfdesc,
			Wfminstat:    wfDef.Wfminstat,
			Wfmaxstat:    wfDef.Wfmaxstat,
			Flow:         wfDef.Flow,
			Moduleid:     wfDef.Moduleid,
			Recordstatus: wfDef.Recordstatus,
			Updatedate:   time.Now(),
		}
		if err := tx.Create(&workflow).Error; err != nil {
			tx.Rollback()
			return helpers.FailResponse(c, fiber.StatusInternalServerError, "CREATE_WORKFLOW_FAILED", err.Error())
		}
	}

	// 8. Create database tables and track them
	for _, tableDef := range manifest.Tables {
		// Build CREATE TABLE SQL
		createSQL := buildCreateTableSQL(tableDef)

		// Execute CREATE TABLE
		if err := tx.Exec(createSQL).Error; err != nil {
			tx.Rollback()
			return helpers.FailResponse(c, fiber.StatusInternalServerError, "CREATE_TABLE_FAILED", fmt.Sprintf("Table: %s, Error: %s", tableDef.TableName, err.Error()))
		}

		// Track table in moduletables
		moduleTable := models.ModuleTables{
			ModuleID:  moduleID,
			NameTable: tableDef.TableName,
			CreatedAt: time.Now(),
		}
		if err := tx.Create(&moduleTable).Error; err != nil {
			tx.Rollback()
			return helpers.FailResponse(c, fiber.StatusInternalServerError, "TRACK_TABLE_FAILED", err.Error())
		}
	}

	// 8b. Insert table data
	for _, tableDef := range manifest.Tables {
		if len(tableDef.Data) > 0 {
			if err := tx.Table(tableDef.TableName).CreateInBatches(tableDef.Data, 100).Error; err != nil {
				tx.Rollback()
				return helpers.FailResponse(c, fiber.StatusInternalServerError, "INSERT_DATA_FAILED", fmt.Sprintf("Table: %s, Error: %s", tableDef.TableName, err.Error()))
			}
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return helpers.FailResponse(c, fiber.StatusInternalServerError, "TRANSACTION_COMMIT_FAILED", err.Error())
	}

	return helpers.SuccessResponse(c, "MODULE_INSTALLED", fiber.Map{
		"moduleid":   moduleID,
		"modulename": manifest.Modulename,
		"version":    manifest.Moduleversion,
	})
}

// buildCreateTableSQL builds CREATE TABLE SQL from table definition
func buildCreateTableSQL(tableDef TableDefinition) string {
	var sql strings.Builder
	sql.WriteString(fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s` (\n", tableDef.TableName))

	// Add columns
	for i, col := range tableDef.Columns {
		sql.WriteString("  `" + col.Name + "` ")

		// Column type
		if col.Length > 0 {
			sql.WriteString(fmt.Sprintf("%s(%d)", col.Type, col.Length))
		} else {
			sql.WriteString(col.Type)
		}

		// Nullable
		if !col.Nullable {
			sql.WriteString(" NOT NULL")
		}

		// Auto increment
		if col.AutoIncrement {
			sql.WriteString(" AUTO_INCREMENT")
		}

		// Default value
		if col.Default != nil {
			defaultStr := fmt.Sprintf("%v", col.Default)
			if defaultStr == "CURRENT_TIMESTAMP" || strings.Contains(defaultStr, "CURRENT_TIMESTAMP") {
				sql.WriteString(fmt.Sprintf(" DEFAULT %s", defaultStr))
			} else {
				sql.WriteString(fmt.Sprintf(" DEFAULT '%s'", defaultStr))
			}
		}

		// Unique
		if col.Unique {
			sql.WriteString(" UNIQUE")
		}

		if i < len(tableDef.Columns)-1 || len(tableDef.Indexes) > 0 || len(tableDef.ForeignKeys) > 0 || hasPrimaryKey(tableDef.Columns) {
			sql.WriteString(",\n")
		} else {
			sql.WriteString("\n")
		}
	}

	// Add primary key
	pkCols := getPrimaryKeyColumns(tableDef.Columns)
	if len(pkCols) > 0 {
		sql.WriteString(fmt.Sprintf("  PRIMARY KEY (`%s`)", strings.Join(pkCols, "`, `")))
		if len(tableDef.Indexes) > 0 || len(tableDef.ForeignKeys) > 0 {
			sql.WriteString(",\n")
		} else {
			sql.WriteString("\n")
		}
	}

	// Add indexes
	for i, idx := range tableDef.Indexes {
		sql.WriteString(fmt.Sprintf("  KEY `%s` (`%s`)", idx.Name, strings.Join(idx.Columns, "`, `")))
		if i < len(tableDef.Indexes)-1 || len(tableDef.ForeignKeys) > 0 {
			sql.WriteString(",\n")
		} else {
			sql.WriteString("\n")
		}
	}

	// Add foreign keys
	for i, fk := range tableDef.ForeignKeys {
		sql.WriteString(fmt.Sprintf("  CONSTRAINT `%s` FOREIGN KEY (`%s`) REFERENCES `%s` (`%s`)",
			fk.Name, fk.Column, fk.References.Table, fk.References.Column))

		if fk.OnDelete != "" {
			sql.WriteString(fmt.Sprintf(" ON DELETE %s", fk.OnDelete))
		}
		if fk.OnUpdate != "" {
			sql.WriteString(fmt.Sprintf(" ON UPDATE %s", fk.OnUpdate))
		}

		if i < len(tableDef.ForeignKeys)-1 {
			sql.WriteString(",\n")
		} else {
			sql.WriteString("\n")
		}
	}

	sql.WriteString(") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;")
	return sql.String()
}

func hasPrimaryKey(columns []ColumnDefinition) bool {
	for _, col := range columns {
		if col.PrimaryKey {
			return true
		}
	}
	return false
}

func getPrimaryKeyColumns(columns []ColumnDefinition) []string {
	var pkCols []string
	for _, col := range columns {
		if col.PrimaryKey {
			pkCols = append(pkCols, col.Name)
		}
	}
	return pkCols
}

// UninstallModuleHandler handles module removal with cascade cleanup
func UninstallModuleHandler(c *fiber.Ctx, db *gorm.DB) error {
	menuName := c.FormValue("menu")
	moduleIDStr := c.Params("moduleid")
	dropTables := c.QueryBool("drop_tables", false)

	// Check permission
	IsPermission, err := CheckUserPermission(c, db, menuName, PermPurge)
	if err != nil || !IsPermission {
		return helpers.FailResponse(c, fiber.StatusUnauthorized, "INVALID_AUTHORIZE", err.Error())
	}

	// Parse module ID
	var moduleID int
	if _, err := fmt.Sscanf(moduleIDStr, "%d", &moduleID); err != nil {
		return helpers.FailResponse(c, fiber.StatusBadRequest, "INVALID_MODULE_ID", "Module ID must be a number")
	}

	// Check if module exists
	var module models.Modules
	if err := db.First(&module, moduleID).Error; err != nil {
		return helpers.FailResponse(c, fiber.StatusNotFound, "MODULE_NOT_FOUND", "Module does not exist")
	}

	// Check if other modules depend on this module
	var dependentModules []models.Modules
	db.Table("modules").
		Joins("JOIN modulerelation ON modules.moduleid = modulerelation.moduleid").
		Where("modulerelation.relationid = ?", moduleID).
		Find(&dependentModules)

	if len(dependentModules) > 0 {
		var depNames []string
		for _, dep := range dependentModules {
			depNames = append(depNames, dep.Modulename)
		}
		depNamesJSON, _ := json.Marshal(depNames)
		return helpers.FailResponse(c, fiber.StatusConflict, "MODULE_HAS_DEPENDENCIES", fmt.Sprintf("Cannot uninstall module. Other modules depend on it: %s", string(depNamesJSON)))
	}

	// Start transaction
	tx := db.Begin()
	if tx.Error != nil {
		return helpers.FailResponse(c, fiber.StatusInternalServerError, "TRANSACTION_START_FAILED", tx.Error.Error())
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Get and optionally drop tables
	var moduleTables []models.ModuleTables
	tx.Where("moduleid = ?", moduleID).Find(&moduleTables)

	if dropTables {
		for _, mt := range moduleTables {
			dropSQL := fmt.Sprintf("DROP TABLE IF EXISTS `%s`", mt.NameTable)
			if err := tx.Exec(dropSQL).Error; err != nil {
				tx.Rollback()
				return helpers.FailResponse(c, fiber.StatusInternalServerError, "DROP_TABLE_FAILED", fmt.Sprintf("Table: %s, Error: %s", mt.NameTable, err.Error()))
			}
		}
	}

	// 2. Delete moduletables records
	if err := tx.Where("moduleid = ?", moduleID).Delete(&models.ModuleTables{}).Error; err != nil {
		tx.Rollback()
		return helpers.FailResponse(c, fiber.StatusInternalServerError, "DELETE_MODULE_TABLES_FAILED", err.Error())
	}

	// 3. Delete workflows associated with module
	var workflows []models.Workflow
	tx.Where("moduleid = ?", moduleID).Find(&workflows)

	if err := tx.Where("moduleid = ?", moduleID).Delete(&models.Workflow{}).Error; err != nil {
		tx.Rollback()
		return helpers.FailResponse(c, fiber.StatusInternalServerError, "DELETE_WORKFLOWS_FAILED", err.Error())
	}

	// 4. Delete widgets (cascade via moduleid foreign key, but explicit delete for clarity)
	if err := tx.Where("moduleid = ?", moduleID).Delete(&models.Widget{}).Error; err != nil {
		tx.Rollback()
		return helpers.FailResponse(c, fiber.StatusInternalServerError, "DELETE_WIDGETS_FAILED", err.Error())
	}

	// 6. Delete dependency records
	if err := tx.Where("moduleid = ? OR relationid = ?", moduleID, moduleID).Delete(&models.Modulerelation{}).Error; err != nil {
		tx.Rollback()
		return helpers.FailResponse(c, fiber.StatusInternalServerError, "DELETE_DEPENDENCIES_FAILED", err.Error())
	}

	// 7. Delete module record (this will cascade delete menuaccess entries)
	if err := tx.Delete(&module).Error; err != nil {
		tx.Rollback()
		return helpers.FailResponse(c, fiber.StatusInternalServerError, "DELETE_MODULE_FAILED", err.Error())
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return helpers.FailResponse(c, fiber.StatusInternalServerError, "TRANSACTION_COMMIT_FAILED", err.Error())
	}

	return helpers.SuccessResponse(c, "MODULE_UNINSTALLED", fiber.Map{
		"moduleid":        moduleID,
		"modulename":      module.Modulename,
		"tables_dropped":  dropTables,
		"tables_count":    len(moduleTables),
		"workflows_count": len(workflows),
	})
}

// GetModuleDetailsHandler retrieves complete module information
func GetModuleDetailsHandler(c *fiber.Ctx, db *gorm.DB) error {
	menuName := c.Query("menu")
	moduleIDStr := c.Params("moduleid")

	// Check permission
	IsPermission, err := CheckUserPermission(c, db, menuName, PermRead)
	if err != nil || !IsPermission {
		return helpers.FailResponse(c, fiber.StatusUnauthorized, "INVALID_AUTHORIZE", err.Error())
	}

	// Parse module ID
	var moduleID int
	if _, err := fmt.Sscanf(moduleIDStr, "%d", &moduleID); err != nil {
		return helpers.FailResponse(c, fiber.StatusBadRequest, "INVALID_MODULE_ID", "Module ID must be a number")
	}

	// Get module
	var module models.Modules
	if err := db.First(&module, moduleID).Error; err != nil {
		return helpers.FailResponse(c, fiber.StatusNotFound, "MODULE_NOT_FOUND", "Module does not exist")
	}

	// Get menu access entries
	var menuAccess []models.Menuaccess
	db.Where("moduleid = ?", moduleID).Find(&menuAccess)

	// Get widgets
	var widgets []models.Widget
	db.Where("moduleid = ?", moduleID).Find(&widgets)

	// Get workflows
	var workflows []models.Workflow
	db.Table("workflow").
		Where("moduleid = ?", moduleID).
		Find(&workflows)

	// Get tables
	var tables []models.ModuleTables
	db.Where("moduleid = ?", moduleID).Find(&tables)

	return helpers.SuccessResponse(c, "MODULE_DETAILS", fiber.Map{
		"module":     module,
		"menuaccess": menuAccess,
		"widgets":    widgets,
		"workflows":  workflows,
		"tables":     tables,
	})
}

// GetModuleDependenciesHandler checks module dependencies
func GetModuleDependenciesHandler(c *fiber.Ctx, db *gorm.DB) error {
	menuName := c.Query("menu")
	moduleIDStr := c.Params("moduleid")

	// Check permission
	IsPermission, err := CheckUserPermission(c, db, menuName, PermRead)
	if err != nil || !IsPermission {
		return helpers.FailResponse(c, fiber.StatusUnauthorized, "INVALID_AUTHORIZE", err.Error())
	}

	// Parse module ID
	var moduleID int
	if _, err := fmt.Sscanf(moduleIDStr, "%d", &moduleID); err != nil {
		return helpers.FailResponse(c, fiber.StatusBadRequest, "INVALID_MODULE_ID", "Module ID must be a number")
	}

	// Get modules this module depends on
	var dependsOn []models.Modules
	db.Table("modules").
		Joins("JOIN modulerelation ON modules.moduleid = modulerelation.relationid").
		Where("modulerelation.moduleid = ?", moduleID).
		Find(&dependsOn)

	// Get modules that depend on this module
	var dependents []models.Modules
	db.Table("modules").
		Joins("JOIN modulerelation ON modules.moduleid = modulerelation.moduleid").
		Where("modulerelation.relationid = ?", moduleID).
		Find(&dependents)

	return helpers.SuccessResponse(c, "MODULE_DEPENDENCIES", fiber.Map{
		"depends_on": dependsOn,
		"dependents": dependents,
	})
}

// ExportModuleHandler exports an existing module as a ZIP package
func ExportModuleHandler(c *fiber.Ctx, db *gorm.DB) error {
	menuName := c.Query("menu")
	moduleIDStr := c.Params("moduleid")

	// Check permission
	IsPermission, err := CheckUserPermission(c, db, menuName, PermRead)
	if err != nil || !IsPermission {
		return helpers.FailResponse(c, fiber.StatusUnauthorized, "INVALID_AUTHORIZE", err.Error())
	}

	// Parse module ID
	var moduleID int
	if _, err := fmt.Sscanf(moduleIDStr, "%d", &moduleID); err != nil {
		return helpers.FailResponse(c, fiber.StatusBadRequest, "INVALID_MODULE_ID", "Module ID must be a number")
	}

	// Get module
	var module models.Modules
	if err := db.First(&module, moduleID).Error; err != nil {
		return helpers.FailResponse(c, fiber.StatusNotFound, "MODULE_NOT_FOUND", "Module does not exist")
	}

	// Get dependencies
	var dependencies []int
	var moduleRelations []models.Modulerelation
	db.Where("moduleid = ?", moduleID).Find(&moduleRelations)
	for _, rel := range moduleRelations {
		dependencies = append(dependencies, rel.Relationid)
	}

	// Get menu access entries
	var menuAccessList []models.Menuaccess
	db.Where("moduleid = ?", moduleID).Find(&menuAccessList)

	var menuAccessDefs []MenuAccessDefinition
	for _, ma := range menuAccessList {
		menuAccessDefs = append(menuAccessDefs, MenuAccessDefinition{
			Menuname:       ma.Menuname,
			Menucode:       ma.Menucode,
			Description:    ma.Description,
			Parentid:       ma.Parentid,
			Menuurl:        ma.Menuurl,
			Sortorder:      ma.Sortorder,
			Menuicon:       ma.Menuicon,
			Menutype:       ma.Menutype,
			Menuform:       ma.Menuform,
			Menuformdetail: ma.Menuformdetail,
			Recordstatus:   ma.Recordstatus,
		})
	}

	// Get widgets
	var widgetList []models.Widget
	db.Where("moduleid = ?", moduleID).Find(&widgetList)

	var widgetDefs []WidgetDefinition
	for _, w := range widgetList {
		widgetDefs = append(widgetDefs, WidgetDefinition{
			Widgetname:    w.Widgetname,
			Widgettitle:   w.Widgettitle,
			Widgetversion: w.Widgetversion,
			Widgetby:      w.Widgetby,
			Description:   w.Description,
			Widgetform:    w.Widgetform,
			Recordstatus:  w.Recordstatus,
		})
	}

	// Get workflows
	var workflows []models.Workflow
	db.Table("workflow").
		Where("moduleid = ?", moduleID).
		Find(&workflows)

	var workflowDefs []WorkflowDefinition
	for _, wf := range workflows {
		workflowDefs = append(workflowDefs, WorkflowDefinition{
			Wfname:       wf.Wfname,
			Wfdesc:       wf.Wfdesc,
			Wfminstat:    wf.Wfminstat,
			Wfmaxstat:    wf.Wfmaxstat,
			Flow:         wf.Flow,
			Moduleid:     wf.Moduleid,
			Recordstatus: wf.Recordstatus,
		})
	}

	// Get tables - Note: We can't export full table schemas, only table names
	// Users will need to manually define table schemas if they want to recreate tables
	var moduleTables []models.ModuleTables
	db.Where("moduleid = ?", moduleID).Find(&moduleTables)

	var tableDefs []TableDefinition
	for _, mt := range moduleTables {
		// Extract full table schema and data
		tableDef, err := extractTableSchema(db, mt.NameTable)
		if err != nil {
			return helpers.FailResponse(c, fiber.StatusInternalServerError, "EXTRACT_SCHEMA_FAILED", fmt.Sprintf("Table: %s, Error: %s", mt.NameTable, err.Error()))
		}
		tableDefs = append(tableDefs, *tableDef)
	}

	// Get GroupMenu data for groupaccessid = 2
	var groupMenus []struct {
		models.Groupmenu
		Menuname string
	}
	db.Table("groupmenu").
		Select("groupmenu.*, menuaccess.menuname").
		Joins("JOIN menuaccess ON groupmenu.menuaccessid = menuaccess.menuaccessid").
		Where("groupmenu.groupaccessid = ? AND menuaccess.moduleid = ?", 2, moduleID).
		Scan(&groupMenus)

	var groupMenuDefs []GroupMenuDefinition
	for _, gm := range groupMenus {
		groupMenuDefs = append(groupMenuDefs, GroupMenuDefinition{
			Menuname:   gm.Menuname,
			Isread:     gm.Isread,
			Iswrite:    gm.Iswrite,
			Ispost:     gm.Ispost,
			Isreject:   gm.Isreject,
			Isupload:   gm.Isupload,
			Isdownload: gm.Isdownload,
			Ispurge:    gm.Ispurge,
		})
	}

	// Create module manifest
	manifest := ModuleManifest{
		Modulename:    module.Modulename,
		Description:   module.Description,
		Moduleversion: module.Moduleversion,
		Createdby:     module.Createdby,
		Themeid:       module.Themeid,
		Dependencies:  dependencies,
		Menuaccess:    menuAccessDefs,
		Widgets:       widgetDefs,
		Workflows:     workflowDefs,
		Tables:        tableDefs,
		GroupMenus:    groupMenuDefs,
	}

	// Convert to JSON
	manifestJSON, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return helpers.FailResponse(c, fiber.StatusInternalServerError, "JSON_MARSHAL_FAILED", err.Error())
	}

	// Create temp directory
	tempDir := "./tmp/module-exports"
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return helpers.FailResponse(c, fiber.StatusInternalServerError, "CREATE_TEMP_DIR_FAILED", err.Error())
	}

	// Write module.json to temp file
	manifestPath := filepath.Join(tempDir, "module.json")
	if err := os.WriteFile(manifestPath, manifestJSON, 0644); err != nil {
		return helpers.FailResponse(c, fiber.StatusInternalServerError, "WRITE_FILE_FAILED", err.Error())
	}
	defer os.Remove(manifestPath)

	// Create ZIP file
	zipFileName := fmt.Sprintf("%s-v%s.zip", strings.ReplaceAll(strings.ToLower(module.Modulename), " ", "-"), module.Moduleversion)
	zipPath := filepath.Join(tempDir, zipFileName)

	zipFile, err := os.Create(zipPath)
	if err != nil {
		return helpers.FailResponse(c, fiber.StatusInternalServerError, "CREATE_ZIP_FAILED", err.Error())
	}
	defer zipFile.Close()
	defer os.Remove(zipPath)

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Add module.json to ZIP
	fileToZip, err := os.Open(manifestPath)
	if err != nil {
		return helpers.FailResponse(c, fiber.StatusInternalServerError, "OPEN_FILE_FAILED", err.Error())
	}
	defer fileToZip.Close()

	info, err := fileToZip.Stat()
	if err != nil {
		return helpers.FailResponse(c, fiber.StatusInternalServerError, "STAT_FILE_FAILED", err.Error())
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return helpers.FailResponse(c, fiber.StatusInternalServerError, "CREATE_HEADER_FAILED", err.Error())
	}

	header.Name = "module.json"
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return helpers.FailResponse(c, fiber.StatusInternalServerError, "CREATE_ZIP_ENTRY_FAILED", err.Error())
	}

	if _, err := io.Copy(writer, fileToZip); err != nil {
		return helpers.FailResponse(c, fiber.StatusInternalServerError, "COPY_TO_ZIP_FAILED", err.Error())
	}

	zipWriter.Close()
	zipFile.Close()

	// Read ZIP file for download
	zipData, err := os.ReadFile(zipPath)
	if err != nil {
		return helpers.FailResponse(c, fiber.StatusInternalServerError, "READ_ZIP_FAILED", err.Error())
	}

	// Send ZIP file as download
	c.Set("Content-Type", "application/zip")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", zipFileName))
	return c.Send(zipData)
}

// extractTableSchema reverses engineer table schema and data
func extractTableSchema(db *gorm.DB, tableName string) (*TableDefinition, error) {
	var columns []ColumnDefinition
	var foreignKeys []ForeignKeyDefinition
	var indexes []IndexDefinition

	// 1. Get Columns
	// Note: This is a simplified extraction relative to the full reverse engineer
	// We want raw SQL types suitable for CREATE TABLE
	driver := db.Dialector.Name()
	var query string

	type ColumnInfo struct {
		ColumnName    string
		DataType      string
		IsNullable    string
		ColumnKey     string
		ColumnDefault *string
		Extra         string
		MaxLength     *int
	}
	var colInfos []ColumnInfo

	if driver == "mysql" {
		query = `
			SELECT 
				COLUMN_NAME as column_name, 
				DATA_TYPE as data_type, 
				IS_NULLABLE as is_nullable, 
				COLUMN_KEY as column_key, 
				COLUMN_DEFAULT as column_default, 
				EXTRA as extra,
				CHARACTER_MAXIMUM_LENGTH as max_length
			FROM INFORMATION_SCHEMA.COLUMNS 
			WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? 
			ORDER BY ORDINAL_POSITION`
		
		if err := db.Raw(query, tableName).Scan(&colInfos).Error; err != nil {
			return nil, err
		}
	} else {
		// Fallback for other drivers (simplified) - reusing GORM migrator if possible would be better but complex
		// For now implementing MySQL support as primary requested context
		return nil, fmt.Errorf("schema export currently supports MySQL only")
	}

	for _, info := range colInfos {
		col := ColumnDefinition{
			Name:          info.ColumnName,
			Type:          info.DataType,
			Nullable:      info.IsNullable == "YES",
			PrimaryKey:    strings.Contains(info.ColumnKey, "PRI"),
			AutoIncrement: strings.Contains(info.Extra, "auto_increment"),
			Unique:        strings.Contains(info.ColumnKey, "UNI"),
		}

		if info.MaxLength != nil {
			col.Length = *info.MaxLength
		}

		if info.ColumnDefault != nil {
			col.Default = *info.ColumnDefault
		}

		columns = append(columns, col)
	}

	// 2. Get Data
	var data []map[string]interface{}
	if err := db.Table(tableName).Find(&data).Error; err != nil {
		return nil, err
	}

	return &TableDefinition{
		TableName:   tableName,
		Columns:     columns,
		ForeignKeys: foreignKeys, // TODO: Implement FK extraction if needed
		Indexes:     indexes,     // TODO: Implement Index extraction if needed
		Data:        data,
	}, nil
}
