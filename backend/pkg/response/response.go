package response

import "github.com/gofiber/fiber/v2"

type ErrorDetail struct {
    Message string `json:"message"`
    Code    string `json:"code"`
}

type ApiResponse[T any] struct {
    Success bool         `json:"success"`
    Data    *T           `json:"data,omitempty"`
    Error   *ErrorDetail `json:"error,omitempty"`
}

func OK[T any](c *fiber.Ctx, data T) error {
    res := ApiResponse[T]{Success: true, Data: &data}
    return c.Status(fiber.StatusOK).JSON(res)
}

func Fail(c *fiber.Ctx, status int, code, msg string) error {
    res := ApiResponse[struct{}]{
        Success: false,
        Error: &ErrorDetail{
            Message: msg,
            Code:    code,
        },
    }
    return c.Status(status).JSON(res)
}
