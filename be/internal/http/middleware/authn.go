package middleware

import (
	"github.com/gofiber/fiber/v3"
	"sinde.ru/internal/auth"
	"sinde.ru/internal/http/responses"
	"sinde.ru/internal/models"
	"strings"
)

const (
	AuthUserLocalKey   = "auth.user"
	AuthClaimsLocalKey = "auth.claims"
)

func RequireAuth(service *auth.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		token := strings.TrimSpace(c.Cookies(service.AccessCookieName()))
		if token == "" {
			return responses.Error(c, fiber.StatusUnauthorized, "требуется аутентификация")
		}
		device := auth.BuildDeviceContext(c.IP(), c.Get("User-Agent"), c.Get("Accept-Language"))
		user, claims, err := service.AuthenticateAccessToken(c.Context(), token, device)
		if err != nil {
			return responses.Error(c, fiber.StatusUnauthorized, "сессия недействительна или истекла")
		}
		c.Locals(AuthUserLocalKey, user)
		c.Locals(AuthClaimsLocalKey, claims)
		return c.Next()
	}
}
func RequireCSRFCookie(service *auth.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		headerToken := strings.TrimSpace(c.Get("X-CSRF-Token"))
		cookieToken := strings.TrimSpace(c.Cookies(service.CSRFCookieName()))
		if headerToken == "" || cookieToken == "" || headerToken != cookieToken {
			return responses.Error(c, fiber.StatusForbidden, "проверка CSRF не пройдена")
		}
		return c.Next()
	}
}
func RequireRole(role string) fiber.Handler {
	return func(c fiber.Ctx) error {
		user := CurrentUser(c)
		if user == nil {
			return responses.Error(c, fiber.StatusUnauthorized, "требуется аутентификация")
		}
		for _, item := range user.Roles {
			if item == role {
				return c.Next()
			}
		}
		return responses.Error(c, fiber.StatusForbidden, "доступ запрещен")
	}
}
func CurrentUser(c fiber.Ctx) *models.User {
	value := c.Locals(AuthUserLocalKey)
	user, _ := value.(*models.User)
	return user
}
func CurrentClaims(c fiber.Ctx) *auth.Claims {
	value := c.Locals(AuthClaimsLocalKey)
	claims, _ := value.(*auth.Claims)
	return claims
}
