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
	Recordstatus int8   `json:"recordstatus"`
}

type TableDefinition struct {
	TableName   string                 `json:"tableName"`
	Columns     []ColumnDefinition     `json:"columns"`
	Indexes     []IndexDefinition      `json:"indexes"`
	ForeignKeys []ForeignKeyDefinition `json:"foreignKeys"`
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
			Recordstatus: wfDef.Recordstatus,
			Updatedate:   time.Now(),
		}
		if err := tx.Create(&workflow).Error; err != nil {
			tx.Rollback()
			return helpers.FailResponse(c, fiber.StatusInternalServerError, "CREATE_WORKFLOW_FAILED", err.Error())
		}

		// Track workflow in moduleworkflow
		moduleWorkflow := models.ModuleWorkflow{
			ModuleID:   moduleID,
			WorkflowID: workflow.Workflowid,
			CreatedAt:  time.Now(),
		}
		if err := tx.Create(&moduleWorkflow).Error; err != nil {
			tx.Rollback()
			return helpers.FailResponse(c, fiber.StatusInternalServerError, "TRACK_WORKFLOW_FAILED", err.Error())
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

	// 3. Get and delete workflows
	var moduleWorkflows []models.ModuleWorkflow
	tx.Where("moduleid = ?", moduleID).Find(&moduleWorkflows)

	var workflowIDs []int
	for _, mw := range moduleWorkflows {
		workflowIDs = append(workflowIDs, mw.WorkflowID)
	}

	if len(workflowIDs) > 0 {
		// Delete workflows
		if err := tx.Where("workflowid IN ?", workflowIDs).Delete(&models.Workflow{}).Error; err != nil {
			tx.Rollback()
			return helpers.FailResponse(c, fiber.StatusInternalServerError, "DELETE_WORKFLOWS_FAILED", err.Error())
		}
	}

	// 4. Delete moduleworkflow records
	if err := tx.Where("moduleid = ?", moduleID).Delete(&models.ModuleWorkflow{}).Error; err != nil {
		tx.Rollback()
		return helpers.FailResponse(c, fiber.StatusInternalServerError, "DELETE_MODULE_WORKFLOWS_FAILED", err.Error())
	}

	// 5. Delete widgets (cascade via moduleid foreign key, but explicit delete for clarity)
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
		"workflows_count": len(workflowIDs),
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
		Joins("JOIN moduleworkflow ON workflow.workflowid = moduleworkflow.workflowid").
		Where("moduleworkflow.moduleid = ?", moduleID).
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
		Joins("JOIN moduleworkflow ON workflow.workflowid = moduleworkflow.workflowid").
		Where("moduleworkflow.moduleid = ?", moduleID).
		Find(&workflows)

	var workflowDefs []WorkflowDefinition
	for _, wf := range workflows {
		workflowDefs = append(workflowDefs, WorkflowDefinition{
			Wfname:       wf.Wfname,
			Wfdesc:       wf.Wfdesc,
			Wfminstat:    wf.Wfminstat,
			Wfmaxstat:    wf.Wfmaxstat,
			Flow:         wf.Flow,
			Recordstatus: wf.Recordstatus,
		})
	}

	// Get tables - Note: We can't export full table schemas, only table names
	// Users will need to manually define table schemas if they want to recreate tables
	var moduleTables []models.ModuleTables
	db.Where("moduleid = ?", moduleID).Find(&moduleTables)

	var tableDefs []TableDefinition
	for _, mt := range moduleTables {
		// Create a basic table definition with just the name
		// Users should manually add column definitions if needed
		tableDefs = append(tableDefs, TableDefinition{
			TableName: mt.NameTable,
			Columns:   []ColumnDefinition{}, // Empty - user must define manually
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
