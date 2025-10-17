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
}
`

// Template isi file routes.go (kosong)
var routesTemplate = `package {{name}}

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// RegisterRoutes digunakan untuk mendaftarkan route plugin {{name}}.
// Isi sesuai kebutuhan modul.
func RegisterRoutes(app *fiber.App, db *gorm.DB) {
	// TODO: implement routes
}
`

// Membuat file dengan isi tertentu, skip jika sudah ada
func createFile(path, content string) error {
	if _, err := os.Stat(path); err == nil {
		fmt.Printf("⚠️  File already exists: %s\n", path)
		return nil
	}
	return os.WriteFile(path, []byte(content), 0644)
}

func GeneratePlugin(pluginName string) error {
	caser := cases.Title(language.English)
	structName := caser.String(pluginName)
	basePath := filepath.Join("plugins", pluginName)
	if _, err := os.Stat(basePath); err == nil {
		fmt.Printf("⚠️  Plugin '%s' already exists.\n", pluginName)
		return err
	}

	// Buat folder plugin
	if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
		fmt.Printf("❌ Failed to create folder: %v\n", err)
		return err
	}

	initContent := strings.ReplaceAll(strings.ReplaceAll(initTemplate, "{{name}}", pluginName), "{{structName}}", structName)
	routesContent := strings.ReplaceAll(routesTemplate, "{{name}}", pluginName)

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

	return nil
}
