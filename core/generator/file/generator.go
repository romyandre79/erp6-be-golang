package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var initTemplate = `package {{name}}

import (
	"erp6-be-golang/core/plugin"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type {{structName}}Plugin struct{}

func (h {{structName}}Plugin) Name() string { return "{{name}}" }

func (h {{structName}}Plugin) Init(db *gorm.DB) error { return nil }

func (h {{structName}}Plugin) Routes(app *fiber.App, db *gorm.DB) {
	RegisterRoutes(app, db)
}

func init() {
	plugin.Register({{structName}}Plugin{})

	{{initModels}}
}
`

// Template isi file routes.go (kosong)
var routesTemplate = `package {{name}}

import (
	"erp6-be-golang/core/plugin"
	"erp6-be-golang/models"
	"erp6-be-golang/plugins/admin"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// RegisterRoutes digunakan untuk mendaftarkan route plugin {{name}}.
// Isi sesuai kebutuhan modul.
func RegisterRoutes(app *fiber.App, db *gorm.DB) {
	// TODO: implement routes
	{{initroutes}}
}
`

var modelControllerTemplate = `package {{name}}

import "erp6-be-golang/core/events"

func Init{{model}} () {
	events.Register("BeforeGet:{{model}}", func(data interface{}) error {
	    // TODO: implement before get
		return nil
	})
	events.Register("AfterGet:{{model}}", func(data interface{}) error {
	    // TODO: implement after get
		return nil
	})
	events.Register("BeforeSave:{{model}}", func(data interface{}) error {
	    // TODO: implement before save
		return nil
	})
	events.Register("AfterSave:{{model}}", func(data interface{}) error {
	    // TODO: implement after save
		return nil
	})
	events.Register("BeforeUpdate:{{model}}", func(data interface{}) error {
	    // TODO: implement before update
		return nil
	})
	events.Register("AfterUpdate:{{model}}", func(data interface{}) error {
	    // TODO: implement after update
		return nil
	})
	events.Register("BeforeDelete:{{model}}", func(data interface{}) error {
	    // TODO: implement before delete
		return nil
	})
	events.Register("AfterDelete:{{model}}", func(data interface{}) error {
	    // TODO: implement after delete
		return nil
	})
}
`

// Membuat file dengan isi tertentu, skip jika sudah ada
func createFile(path, content string) error {
	if _, err := os.Stat(path); err == nil {
		fmt.Printf("⚠️  File already exists: %s\n", path)
		return err
	}
	return os.WriteFile(path, []byte(content), 0644)
}

func GeneratePlugin(pluginName string, modelList []string) error {
	caser := cases.Title(language.English)
	structName := caser.String(pluginName)
	basePath := filepath.Join("plugins", pluginName)
	// Buat folder plugin
	if err := os.MkdirAll(basePath, os.ModePerm); os.IsExist(err) {
		fmt.Printf("❌ Failed to create folder: %v\n", err)
		return err
	}

	initModels := ""
	initRoutes := pluginName + ` := app.Group("/` + pluginName + `")`
	initRoutes += "\n" + pluginName + ".Use(admin.AuthMiddleware)\n{\n"
	for _, val := range modelList {
		initModels += "Init" + val + "()\n"
		initRoutes += `plugin.RegisterModelRoutes(` + pluginName + `, db, models.` + strings.Title(strings.ToLower(val)) + `{}, "` + val + `")`
		initRoutes += "\n"
	}
	initRoutes += "\n}\n"
	initContent := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(initTemplate, "{{name}}", pluginName), "{{structName}}", structName), "{{initModels}}", initModels)
	routesContent := strings.ReplaceAll(strings.ReplaceAll(routesTemplate, "{{name}}", pluginName), "{{initroutes}}", initRoutes)

	files := map[string]string{
		"init.go":       initContent,
		"routes.go":     routesContent,
		"controller.go": "package " + pluginName + "\n\n// TODO: implement controllers",
		"helpers.go":    "package " + pluginName + "\n\n// TODO: implement helpers",
		"middleware.go": "package " + pluginName + "\n\n// TODO: implement middleware",
	}

	for name, content := range files {
		path := filepath.Join(basePath, name)
		if err := createFile(path, content); err != nil {
			fmt.Printf("❌ Failed to create file %s: %v\n", name, err)
			return err
		} else {
			fmt.Printf("✅ Created %s\n", path)
		}
	}

	for _, val := range modelList {
		modelControllerContent := strings.ReplaceAll(strings.ReplaceAll(modelControllerTemplate, "{{model}}", val), "{{name}}", pluginName)
		path := filepath.Join(basePath, val+"controller.go")
		if err := createFile(path, modelControllerContent); err != nil {
			fmt.Printf("❌ Failed to create file %s: %v\n", val, err)
			return err
		} else {
			fmt.Printf("✅ Created %s\n", path)
		}
	}

	return nil
}
