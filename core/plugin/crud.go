package plugin

import (
	"erp6-be-golang/core/events"
	"erp6-be-golang/core/helpers"
	"erp6-be-golang/core/i18n"
	"errors"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/utils"
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
		lang := c.FormValue("lang")
		if lang == "" {
			lang = "id"
		}

		search := c.Query("search", "")
		page, _ := strconv.Atoi(c.Query("page", "0"))
		limit, _ := strconv.Atoi(c.Query("limit", "10"))
		offset, _ := strconv.Atoi(c.Query("offset", "0"))

		modelType := reflect.TypeOf(model)
		sliceType := reflect.SliceOf(modelType)
		slicePtr := reflect.New(sliceType).Interface()

		dbWithPreload := autoPreload(db, modelType)
		query := dbWithPreload.Model(model)

		// 🔍 General search di semua field string
		if search != "" {
			var conditions []string
			var values []interface{}

			// Loop semua field dari struct
			for i := 0; i < modelType.NumField(); i++ {
				field := modelType.Field(i)

				// Abaikan field non-string
				if field.Type.Kind() != reflect.String {
					continue
				}

				// Ambil nama kolom dari tag `gorm:"column:xxx"`
				columnName := field.Tag.Get("gorm")
				if columnName == "" {
					columnName = field.Name
				} else {
					// Ekstrak dari format `column:nama_kolom;otherTag`
					if strings.Contains(columnName, "column:") {
						parts := strings.Split(columnName, ";")
						for _, p := range parts {
							if strings.HasPrefix(p, "column:") {
								columnName = strings.TrimPrefix(p, "column:")
								break
							}
						}
					} else {
						columnName = field.Name
					}
				}

				conditions = append(conditions, fmt.Sprintf("%s LIKE ?", columnName))
				values = append(values, "%"+search+"%")
			}

			if len(conditions) > 0 {
				query = query.Where(strings.Join(conditions, " OR "), values...)
			}
		}

		// Hitung total
		var total int64
		query.Count(&total)

		// Pagination
		if page > 0 {
			start := (page - 1) * limit
			query = query.Offset(start).Limit(limit)
		} else if offset > 0 {
			query = query.Offset(offset).Limit(limit)
		}

		// Eksekusi query
		if err := query.Find(slicePtr).Error; err != nil {
			return helpers.FailResponse(c, 400, "INVALID_ERROR", err.Error())
		}

		// Meta info
		var meta fiber.Map
		if page > 0 {
			meta = fiber.Map{
				"page":       page,
				"limit":      limit,
				"total":      total,
				"totalPages": int(math.Ceil(float64(total) / float64(limit))),
			}
		} else if offset > 0 {
			meta = fiber.Map{
				"offset":  offset,
				"limit":   limit,
				"total":   total,
				"hasMore": (int64(offset) + int64(limit)) < total,
			}
		} else {
			meta = fiber.Map{
				"page":       1,
				"limit":      10,
				"total":      total,
				"totalPages": int(math.Ceil(float64(total) / float64(limit))),
			}
		}

		return helpers.SuccessResponse(c, "DATA RETRIEVED", fiber.Map{
			"data": slicePtr,
			"meta": meta,
		})
	}
}

func getHandler(db *gorm.DB, model interface{}, name string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		lang := c.FormValue("lang")
		if lang == "" {
			lang = "id"
		}
		id := c.Params("id")
		obj := reflect.New(reflect.TypeOf(model)).Interface()

		if err := events.Trigger("BeforeGet:"+name, obj); err != nil {
			return helpers.FailResponse(c, 400, "INVALID_RETRIEVE", err.Error())
		}

		dbWithPreload := autoPreload(db, reflect.TypeOf(model))

		if err := dbWithPreload.First(obj, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return helpers.FailResponse(c, 400, "INVALID DATA", "DATA NOT_FOUND")
			}
			return helpers.FailResponse(c, 400, "INVALID DATA", err.Error())
		}

		if err := events.Trigger("AfterGet:"+name, obj); err != nil {
			return helpers.FailResponse(c, 400, "INVALID_RETRIEVE", err.Error())
		}

		return helpers.SuccessResponse(c, "DATA RETRIEVED", obj)
	}
}

func createHandler(db *gorm.DB, model interface{}, name string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		lang := c.FormValue("lang")
		if lang == "" {
			lang = "id"
		}
		obj := reflect.New(reflect.TypeOf(model)).Interface()

		if err := events.Trigger("BeforeCreate:"+name, obj); err != nil {
			return helpers.FailResponse(c, 400, "INVALID_CREATE", err.Error())
		}

		if err := c.BodyParser(obj); err != nil {
			return helpers.FailResponse(c, 400, "INVALID_CREATE", err.Error())
		}

		if err := db.Create(obj).Error; err != nil {
			return helpers.FailResponse(c, 400, "INVALID_CREATE", friendlySQLError(err, lang))
		}

		if err := events.Trigger("AfterCreate:"+name, obj); err != nil {
			return helpers.FailResponse(c, 400, "INVALID_CREATE", err.Error())
		}
		return helpers.SuccessResponse(c, "DATA SAVED", obj)
	}
}

func updateHandler(db *gorm.DB, model interface{}, name string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		lang := c.FormValue("lang")
		if lang == "" {
			lang = "id"
		}
		id := c.Params("id")
		obj := reflect.New(reflect.TypeOf(model)).Interface()

		if err := events.Trigger("BeforeUpdate:"+name, obj); err != nil {
			return helpers.FailResponse(c, 400, "INVALID_UPDATE", err.Error())
		}

		if err := db.First(obj, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return helpers.FailResponse(c, 400, "INVALID_UPDATE", "Record Not Found")
			}
			return helpers.FailResponse(c, 400, "INVALID_UPDATE", err.Error())
		}

		if err := c.BodyParser(obj); err != nil {
			return helpers.FailResponse(c, 400, "INVALID_UPDATE", err.Error())
		}

		if err := db.Save(obj).Error; err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		if err := events.Trigger("AfterUpdate:"+name, obj); err != nil {
			return helpers.FailResponse(c, 400, "INVALID_UPDATE", err.Error())
		}

		return helpers.SuccessResponse(c, "DATA UPDATED", obj)
	}
}

func deleteHandler(db *gorm.DB, model interface{}, name string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		lang := c.FormValue("lang")
		if lang == "" {
			lang = "id"
		}
		id := c.Params("id")
		obj := reflect.New(reflect.TypeOf(model)).Interface()
		if err := events.Trigger("BeforeDelete:"+name, obj); err != nil {
			return helpers.FailResponse(c, 400, "INVALID_DELETE", err.Error())
		}

		if err := db.Delete(obj, id).Error; err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		if err := events.Trigger("AfterDelete:"+name, obj); err != nil {
			return helpers.FailResponse(c, 400, "INVALID_DELETE", err.Error())
		}

		return helpers.SuccessResponse(c, "DATA DELETED", obj)
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

func friendlySQLError(err error, lang string) string {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return i18n.Translate(lang, "DATA NOT_FOUND", nil)
	}

	msg := err.Error()

	// Case 1: Foreign key constraint
	if strings.Contains(msg, "foreign key constraint fails") {
		// Cari nama kolom foreign key-nya kalau ada
		re := regexp.MustCompile("FOREIGN KEY \\(`(.*?)`\\)")
		matches := re.FindStringSubmatch(msg)
		if len(matches) > 1 {
			return i18n.Translate(lang, "INVALID_FIELD_REQUIRED", map[string]interface{}{"field": matches[1]})
		}
		return i18n.Translate(lang, "INVALID_REFERENCE", nil)
	}

	// Case 2: Duplicate entry
	if strings.Contains(msg, "Duplicate entry") {
		re := regexp.MustCompile("Duplicate entry '(.*?)'")
		matches := re.FindStringSubmatch(msg)
		if len(matches) > 1 {
			return i18n.Translate(lang, "INVALID DATA EXISTS", map[string]interface{}{"field": matches[1]})
		}
		return i18n.Translate(lang, "INVALID DATA DUPLICATE", nil)
	}

	// Case 3: Column cannot be null
	if strings.Contains(msg, "cannot be null") {
		re := regexp.MustCompile("Column '(.*?)'")
		matches := re.FindStringSubmatch(msg)
		if len(matches) > 1 {
			return i18n.Translate(lang, "INVALID_FIELD_REQUIRED", map[string]interface{}{"field": matches[1]})
		}
		return i18n.Translate(lang, "INVALID_REFERENCE", nil)
	}

	// Default fallback
	return utils.FileWithLineNum() + " | " + msg
}
