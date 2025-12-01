package middleware

import (
	"strings"

	"github.com/example/team-stack/backend/internal/app/ports"
	"github.com/example/team-stack/backend/pkg/response"
	"github.com/gofiber/fiber/v2"
)

const authClaimsKey = "authClaims"

func Authenticate(jwtm ports.JWTManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := strings.TrimSpace(c.Get(fiber.HeaderAuthorization))
		if header == "" {
			return response.Fail(c, fiber.StatusUnauthorized, "ERR_UNAUTHORIZED", "missing authorization header")
		}

		const bearerPrefix = "bearer "
		lower := strings.ToLower(header)
		if !strings.HasPrefix(lower, bearerPrefix) {
			return response.Fail(c, fiber.StatusUnauthorized, "ERR_UNAUTHORIZED", "invalid authorization header")
		}

		token := strings.TrimSpace(header[len(bearerPrefix):])
		if token == "" {
			return response.Fail(c, fiber.StatusUnauthorized, "ERR_UNAUTHORIZED", "invalid authorization header")
		}

		claims, err := jwtm.Verify(token)
		if err != nil {
			return response.Fail(c, fiber.StatusUnauthorized, "ERR_UNAUTHORIZED", "invalid or expired token")
		}

		c.Locals(authClaimsKey, claims)
		return c.Next()
	}
}

func RequireRoles(roles ...string) fiber.Handler {
	allowed := map[string]struct{}{}
	for _, r := range roles {
		allowed[r] = struct{}{}
	}

	return func(c *fiber.Ctx) error {
		claims, ok := c.Locals(authClaimsKey).(*ports.AuthClaims)
		if !ok || claims == nil {
			return response.Fail(c, fiber.StatusUnauthorized, "ERR_UNAUTHORIZED", "missing authentication context")
		}

		if len(allowed) == 0 {
			return c.Next()
		}

		if _, ok := allowed[claims.Role]; !ok {
			return response.Fail(c, fiber.StatusForbidden, "ERR_FORBIDDEN", "insufficient permissions")
		}

		return c.Next()
	}
}
