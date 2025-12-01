package middleware

import (
	"github.com/example/team-stack/backend/pkg/response"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func Recovery(log *zap.SugaredLogger) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		defer func() {
			if r := recover(); r != nil {
				log.Errorw("panic", "err", r)
				_ = response.Fail(c, fiber.StatusInternalServerError, "ERR_PANIC", "internal error")
			}
		}()
		return c.Next()
	}
}
