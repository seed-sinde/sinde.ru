package auth
import (
	"time"
	"github.com/gofiber/fiber/v3"
)
func (s *Service) setAuthCookies(c fiber.Ctx, bundle *TokenBundle) {
	c.Cookie(&fiber.Cookie{
		Name:     s.cfg.AccessCookieName,
		Value:    bundle.AccessToken,
		HTTPOnly: true,
		Secure:   s.cfg.CookieSecure,
		SameSite: s.cfg.CookieSameSite,
		Path:     "/",
		Domain:   s.cfg.CookieDomain,
		Expires:  bundle.AccessExpiresAt,
	})
	c.Cookie(&fiber.Cookie{
		Name:     s.cfg.RefreshCookieName,
		Value:    bundle.RefreshToken,
		HTTPOnly: true,
		Secure:   s.cfg.CookieSecure,
		SameSite: s.cfg.CookieSameSite,
		Path:     "/",
		Domain:   s.cfg.CookieDomain,
		Expires:  bundle.RefreshExpiresAt,
	})
	c.Cookie(&fiber.Cookie{
		Name:     s.cfg.CSRFCookieName,
		Value:    bundle.CSRFToken,
		HTTPOnly: false,
		Secure:   s.cfg.CookieSecure,
		SameSite: s.cfg.CookieSameSite,
		Path:     "/",
		Domain:   s.cfg.CookieDomain,
		Expires:  bundle.RefreshExpiresAt,
	})
	c.Cookie(&fiber.Cookie{
		Name:     s.cfg.SessionHintCookieName,
		Value:    "1",
		HTTPOnly: false,
		Secure:   s.cfg.CookieSecure,
		SameSite: s.cfg.CookieSameSite,
		Path:     "/",
		Domain:   s.cfg.CookieDomain,
		Expires:  bundle.RefreshExpiresAt,
	})
}
func (s *Service) clearAuthCookies(c fiber.Ctx) {
	expired := time.Unix(0, 0)
	for _, name := range []string{
		s.cfg.AccessCookieName,
		s.cfg.RefreshCookieName,
		s.cfg.CSRFCookieName,
		s.cfg.SessionHintCookieName,
	} {
		c.Cookie(&fiber.Cookie{
			Name:     name,
			Value:    "",
			HTTPOnly: name != s.cfg.CSRFCookieName,
			Secure:   s.cfg.CookieSecure,
			SameSite: s.cfg.CookieSameSite,
			Path:     "/",
			Domain:   s.cfg.CookieDomain,
			Expires:  expired,
		})
	}
}
