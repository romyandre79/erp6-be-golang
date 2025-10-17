package plugin

import (
	"erp6-be-golang/core/events"
	"erp6-be-golang/core/helpers"
	"fmt"
	"reflect"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// RegisterModelRoutes membuat endpoint CRUD otomatis berdasarkan model GORM.
func RegisterModelRoutes(router fiber.Router, db *gorm.DB, model interface{}, name string) {
	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	modelName := strings.ToLower(modelType.Name())

	basePath := fmt.Sprintf("/%s", modelName)
	fmt.Printf("[CRUD] Registering routes for %s\n", modelName)

	router.Get(basePath, listHandler(db, model))
	router.Get(basePath+"/:id", getHandler(db, model, name))
	router.Post(basePath, createHandler(db, model, name))
	router.Put(basePath+"/:id", updateHandler(db, model, name))
	router.Delete(basePath+"/:id", deleteHandler(db, model, name))
}

func listHandler(db *gorm.DB, model interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		modelType := reflect.TypeOf(model)
		sliceType := reflect.SliceOf(reflect.TypeOf(model))
		slicePtr := reflect.New(sliceType).Interface()

		dbWithPreload := autoPreload(db, modelType)

		if err := dbWithPreload.Find(slicePtr).Error; err != nil {
			return helpers.FailResponse(c, 500, "INVALID_ERROR", err.Error())
		}
		return helpers.SuccessResponse(c, "DATA_RETRIEVED", slicePtr)
	}
}

func getHandler(db *gorm.DB, model interface{}, name string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		obj := reflect.New(reflect.TypeOf(model)).Interface()

		if err := events.Trigger("BeforeGet:"+name, obj); err != nil {
			return helpers.FailResponse(c, 500, "INVALID_RETRIEVE", err.Error())
		}

		dbWithPreload := autoPreload(db, reflect.TypeOf(model))

		if err := dbWithPreload.First(obj, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return helpers.FailResponse(c, 404, "INVALID_DATA", "Record not found")
			}
			return helpers.FailResponse(c, 500, "INVALID_DATA", err.Error())
		}

		if err := events.Trigger("AfterGet:"+name, obj); err != nil {
			return helpers.FailResponse(c, 500, "INVALID_RETRIEVE", err.Error())
		}

		return helpers.SuccessResponse(c, "DATA_RETRIEVED", obj)
	}
}

func createHandler(db *gorm.DB, model interface{}, name string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		obj := reflect.New(reflect.TypeOf(model)).Interface()

		if err := events.Trigger("BeforeCreate:"+name, obj); err != nil {
			return helpers.FailResponse(c, 400, "INVALID_CREATE", err.Error())
		}

		if err := c.BodyParser(obj); err != nil {
			return helpers.FailResponse(c, 500, "INVALID_CREATE", err.Error())
		}

		if err := db.Create(obj).Error; err != nil {
			return helpers.FailResponse(c, 500, "INVALID_CREATE", err.Error())
		}

		if err := events.Trigger("AfterCreate:"+name, obj); err != nil {
			return helpers.FailResponse(c, 400, "INVALID_CREATE", err.Error())
		}
		return helpers.SuccessResponse(c, "DATA_SAVED", obj)
	}
}

func updateHandler(db *gorm.DB, model interface{}, name string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		obj := reflect.New(reflect.TypeOf(model)).Interface()

		if err := events.Trigger("BeforeUpdate:"+name, obj); err != nil {
			return helpers.FailResponse(c, 400, "INVALID_UPDATE", err.Error())
		}

		if err := db.First(obj, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(404).JSON(fiber.Map{"error": "Record not found"})
			}
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		if err := c.BodyParser(obj); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON payload", "detail": err.Error()})
		}

		if err := db.Save(obj).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		if err := events.Trigger("AfterUpdate:"+name, obj); err != nil {
			return helpers.FailResponse(c, 400, "INVALID_UPDATE", err.Error())
		}

		return helpers.SuccessResponse(c, "DATA_UPDATED", obj)
	}
}

func deleteHandler(db *gorm.DB, model interface{}, name string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		obj := reflect.New(reflect.TypeOf(model)).Interface()
		if err := events.Trigger("BeforeDelete:"+name, obj); err != nil {
			return helpers.FailResponse(c, 400, "INVALID_DELETE", err.Error())
		}

		if err := db.Delete(obj, id).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		if err := events.Trigger("AfterDelete:"+name, obj); err != nil {
			return helpers.FailResponse(c, 400, "INVALID_DELETE", err.Error())
		}

		return helpers.SuccessResponse(c, "DATA_DELETED", obj)
	}
}

func autoPreload(db *gorm.DB, modelType reflect.Type) *gorm.DB {
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		tag := field.Tag.Get("gorm")

		// Preload untuk relasi (1-N, N-1, N-N)
		if strings.Contains(tag, "foreignKey") ||
			strings.Contains(tag, "many2many") ||
			strings.Contains(tag, "references") {
			db = db.Preload(field.Name)
		}
	}
	return db
}
