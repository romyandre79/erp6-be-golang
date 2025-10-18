package common

import (
	"erp6-be-golang/core/plugin"
	"erp6-be-golang/models"
	"erp6-be-golang/plugins/admin"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// RegisterRoutes digunakan untuk mendaftarkan route plugin common.
// Isi sesuai kebutuhan modul.
func RegisterRoutes(app *fiber.App, db *gorm.DB) {
	// TODO: implement routes
	common := app.Group("/common")
	common.Use(admin.AuthMiddleware)
	{
		plugin.RegisterModelRoutes(common, db, models.Addresstype{}, "addresstype")
		plugin.RegisterModelRoutes(common, db, models.Identitytype{}, "identitytype")
		plugin.RegisterModelRoutes(common, db, models.Romawi{}, "romawi")
		plugin.RegisterModelRoutes(common, db, models.Unitofmeasure{}, "unitofmeasure")
		plugin.RegisterModelRoutes(common, db, models.Materialtype{}, "materialtype")
		plugin.RegisterModelRoutes(common, db, models.Materialgroup{}, "materialgroup")
		plugin.RegisterModelRoutes(common, db, models.Materialstatus{}, "materialstatus")
		plugin.RegisterModelRoutes(common, db, models.Ownership{}, "ownership")
		plugin.RegisterModelRoutes(common, db, models.Addressbook{}, "addressbook")
		plugin.RegisterModelRoutes(common, db, models.Plant{}, "plant")
		plugin.RegisterModelRoutes(common, db, models.Sloc{}, "sloc")
		plugin.RegisterModelRoutes(common, db, models.Storagebin{}, "storagebin")
		plugin.RegisterModelRoutes(common, db, models.Product{}, "product")
		plugin.RegisterModelRoutes(common, db, models.Productplant{}, "productplant")
		plugin.RegisterModelRoutes(common, db, models.Address{}, "address")
		plugin.RegisterModelRoutes(common, db, models.Addresscontact{}, "addresscontact")
		plugin.RegisterModelRoutes(common, db, models.Groupcustomer{}, "groupcustomer")
		plugin.RegisterModelRoutes(common, db, models.Productsales{}, "productsales")

	}

}
