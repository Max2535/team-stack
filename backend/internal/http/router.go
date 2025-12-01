package http

import (
	"github.com/example/team-stack/backend/internal/app/ports"
	"github.com/example/team-stack/backend/internal/config"
	"github.com/example/team-stack/backend/internal/http/middleware"
	"github.com/example/team-stack/backend/pkg/response"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func RegisterRoutes(app *fiber.App, cfg *config.Config, log *zap.SugaredLogger, userSvc ports.UserService, jwtm ports.JWTManager, cache ports.Cache) {
	app.Use(middleware.Recovery(log))
	app.Use(middleware.RequestLogger(log))

	api := app.Group("/api")

	api.Get("/health", func(c *fiber.Ctx) error {
		return response.OK(c, map[string]string{"status": "ok"})
	})

	v1 := api.Group("/v1")

	v1.Post("/auth/login", func(c *fiber.Ctx) error {
		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.BodyParser(&req); err != nil {
			return response.Fail(c, fiber.StatusBadRequest, "ERR_BAD_REQUEST", "invalid body")
		}
		u, token, err := userSvc.Login(c.Context(), req.Email, req.Password)
		if err != nil {
			return response.Fail(c, fiber.StatusUnauthorized, "ERR_LOGIN", "invalid credentials")
		}
		return response.OK(c, map[string]any{
			"user":  u,
			"token": token,
		})
	})

	protected := v1.Group("", middleware.Authenticate(jwtm))

	protected.Get("/users", middleware.RequireRoles("admin"), func(c *fiber.Ctx) error {
		users, err := userSvc.List(c.Context())
		if err != nil {
			return response.Fail(c, fiber.StatusInternalServerError, "ERR_USERS", "cannot list users")
		}
		return response.OK(c, users)
	})
}
