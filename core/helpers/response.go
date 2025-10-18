package helpers

import (
	"erp6-be-golang/core/i18n"

	"github.com/gofiber/fiber/v2"
)

type SuccessData struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorData struct {
	Code    int    `json:"code"`
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

// SuccessResponse kirim response sukses
func SuccessResponse(c *fiber.Ctx, message string, data interface{}) error {
	lang := c.FormValue("lang")
	if lang == "" {
		lang = "id"
	}
	return c.Status(fiber.StatusOK).JSON(SuccessData{
		Code:    fiber.StatusOK,
		Message: i18n.Translate(lang, message, nil),
		Data:    data,
	})
}

// FailResponse kirim response error
func FailResponse(c *fiber.Ctx, code int, message string, details string) error {
	lang := c.FormValue("lang")
	if lang == "" {
		lang = "id"
	}
	return c.Status(code).JSON(ErrorData{
		Code:    code,
		Error:   i18n.Translate(lang, message, nil),
		Details: i18n.Translate(lang, details, nil),
	})
}
