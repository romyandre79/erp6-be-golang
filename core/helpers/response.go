package helpers

import "github.com/gofiber/fiber/v2"

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
	return c.Status(fiber.StatusOK).JSON(SuccessData{
		Code:    fiber.StatusOK,
		Message: message,
		Data:    data,
	})
}

// FailResponse kirim response error
func FailResponse(c *fiber.Ctx, code int, message string, details string) error {
	return c.Status(code).JSON(ErrorData{
		Code:    code,
		Error:   message,
		Details: details,
	})
}
