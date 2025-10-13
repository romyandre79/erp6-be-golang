package helpers

import "github.com/gofiber/fiber/v2"

type SuccessData struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorData struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

// SuccessResponse kirim response sukses
func SuccessResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(SuccessData{
		Message: message,
		Data:    data,
	})
}

// FailResponse kirim response error
func FailResponse(c *fiber.Ctx, code int, message string, details string) error {
	return c.Status(code).JSON(ErrorData{
		Error:   message,
		Details: details,
	})
}
