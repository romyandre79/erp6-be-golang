package plugin

import (
	"erp6-be-golang/models"
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var registry = map[string]Plugin{}

func Register(p Plugin) {
	if _, exists := registry[p.Name()]; exists {
		log.Printf("Plugin %s has registered", p.Name())
		return
	}
	registry[p.Name()] = p
}

func LoadActivePlugins(db *gorm.DB, app *fiber.App) {

	active := []string{}
	var modules []models.Modules
	db.Where("recordstatus = ?", 1).Find(&modules)

	for _, m := range modules {
		active = append(active, m.Modulename)
	}

	for _, name := range active {
		if p, ok := registry[name]; ok {
			log.Printf("Load plugin: %s", name)
			_ = p.Init(db)
			p.Routes(app, db)
		} else {
			log.Printf("Error plugin: %s", name)
		}
	}
}
