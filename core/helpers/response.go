package helpers

import (
	"io"
	"net/http"
	"time"

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

func GetRemoteData(url string, timeoutSeconds int) ([]byte, error) {
	client := http.Client{
		Timeout: time.Duration(timeoutSeconds) * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Baca seluruh body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
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
