package auth
import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)
func LoadConfig() (Config, error) {
	get := func(key string) string {
		return strings.TrimSpace(os.Getenv(key))
	}
	cfg := Config{
		PublicBaseURL:           firstNonEmpty(get("AUTH_PUBLIC_BASE_URL"), get("NUXT_PUBLIC_BASE_URL")),
		JWTIssuer:               firstNonEmpty(get("AUTH_JWT_ISSUER"), "example.com"),
		JWTAudience:             firstNonEmpty(get("AUTH_JWT_AUDIENCE"), "app-web"),
		AccessTTL:               parseDurationMinutes(get("AUTH_ACCESS_TTL_MINUTES"), 15),
		RefreshTTL:              parseDurationHours(get("AUTH_REFRESH_TTL_HOURS"), 24*14),
		VerifyTTL:               parseDurationHours(get("AUTH_VERIFY_TTL_HOURS"), 24),
		ResetTTL:                parseDurationMinutes(get("AUTH_RESET_TTL_MINUTES"), 30),
		MFATicketTTL:            parseDurationMinutes(get("AUTH_MFA_TICKET_TTL_MINUTES"), 5),
		TwoFactorSetupTTL:       parseDurationMinutes(get("AUTH_2FA_SETUP_TTL_MINUTES"), 10),
		CookieDomain:            get("AUTH_COOKIE_DOMAIN"),
		CookieSecure:            parseBoolDefault(get("AUTH_COOKIE_SECURE"), true),
		CookieSameSite:          firstNonEmpty(get("AUTH_COOKIE_SAME_SITE"), "Strict"),
		TOTPIssuer:              firstNonEmpty(get("AUTH_TOTP_ISSUER"), "app-name"),
		PasswordMinLength:       parseIntDefault(get("AUTH_PASSWORD_MIN_LENGTH"), 12),
		AllowedOrigins:          splitCSV(firstNonEmpty(get("AUTH_ALLOWED_ORIGINS"), get("NUXT_PUBLIC_BASE_URL"))),
		RedisAddr:               firstNonEmpty(get("REDIS_ADDR"), "127.0.0.1:6379"),
		RedisPassword:           get("REDIS_PASSWORD"),
		RedisDB:                 parseIntDefault(get("REDIS_DB"), 0),
		RegisterRateLimit:       parseIntDefault(get("AUTH_REGISTER_LIMIT"), 10),
		LoginRateLimitIP:        parseIntDefault(get("AUTH_LOGIN_LIMIT_IP"), 20),
		LoginRateLimitMail:      parseIntDefault(get("AUTH_LOGIN_LIMIT_EMAIL"), 10),
		TwoFactorSetupRateLimit: parseIntDefault(get("AUTH_2FA_SETUP_LIMIT"), 5),
		MFALoginRateLimitIP:     parseIntDefault(get("AUTH_MFA_LOGIN_LIMIT_IP"), 20),
		MFALoginRateLimitTicket: parseIntDefault(get("AUTH_MFA_LOGIN_LIMIT_TICKET"), 6),
		VerifyRateLimit:         parseIntDefault(get("AUTH_VERIFY_REQUEST_LIMIT"), 5),
		ResetRateLimit:          parseIntDefault(get("AUTH_RESET_LIMIT"), 10),
		RateWindow:              parseDurationMinutes(get("AUTH_RATE_WINDOW_MINUTES"), 15),
		MailFrom:                firstNonEmpty(get("AUTH_MAIL_FROM"), "no-reply@example.com"),
		MailerDriver:            firstNonEmpty(get("AUTH_MAILER_DRIVER"), "log"),
		SMTPHost:                get("AUTH_SMTP_HOST"),
		SMTPPort:                parseIntDefault(get("AUTH_SMTP_PORT"), 587),
		SMTPUsername:            get("AUTH_SMTP_USERNAME"),
		SMTPPassword:            get("AUTH_SMTP_PASSWORD"),
		SMTPTLSMode:             firstNonEmpty(get("AUTH_SMTP_TLS_MODE"), "starttls"),
		SMTPHELO:                get("AUTH_SMTP_HELO"),
	}
	cfg.AccessCookieName = cookieNameForMode(get("AUTH_ACCESS_COOKIE_NAME"), cfg.CookieSecure, "__Host-access_token", "access_token")
	cfg.RefreshCookieName = cookieNameForMode(get("AUTH_REFRESH_COOKIE_NAME"), cfg.CookieSecure, "__Host-refresh_token", "refresh_token")
	cfg.CSRFCookieName = cookieNameForMode(get("AUTH_CSRF_COOKIE_NAME"), cfg.CookieSecure, "csrf_token", "csrf_token")
	cfg.SessionHintCookieName = cookieNameForMode(
		get("AUTH_SESSION_HINT_COOKIE_NAME"),
		cfg.CookieSecure,
		"__Host-auth_session_hint",
		"auth_session_hint",
	)
	jwtSecret := firstNonEmpty(get("AUTH_JWT_SECRET"), get("AUTH_JWT_SECRET_B64"))
	if jwtSecret == "" {
		return Config{}, fmt.Errorf("AUTH_JWT_SECRET is required")
	}
	if err := validateSecretMaterial("AUTH_JWT_SECRET", jwtSecret); err != nil {
		return Config{}, err
	}
	cfg.JWTSecret = deriveKey(jwtSecret)
	encryptionSecret := firstNonEmpty(get("AUTH_ENCRYPTION_SECRET"), get("AUTH_ENCRYPTION_SECRET_B64"))
	if encryptionSecret == "" {
		return Config{}, fmt.Errorf("AUTH_ENCRYPTION_SECRET is required")
	}
	if err := validateSecretMaterial("AUTH_ENCRYPTION_SECRET", encryptionSecret); err != nil {
		return Config{}, err
	}
	cfg.EncryptionKey = deriveKey(encryptionSecret)
	if cfg.PublicBaseURL == "" {
		return Config{}, fmt.Errorf("AUTH_PUBLIC_BASE_URL or NUXT_PUBLIC_BASE_URL is required")
	}
	if len(cfg.JWTSecret) < 32 {
		return Config{}, fmt.Errorf("AUTH_JWT_SECRET must derive to at least 32 bytes")
	}
	if len(cfg.EncryptionKey) != 32 {
		return Config{}, fmt.Errorf("AUTH_ENCRYPTION_SECRET must derive to 32 bytes")
	}
	return cfg, nil
}
func validateSecretMaterial(name string, secret string) error {
	raw, err := base64.StdEncoding.DecodeString(secret)
	if err == nil {
		if len(raw) < 32 {
			return fmt.Errorf("%s must decode to at least 32 bytes", name)
		}
		return nil
	}
	if len([]byte(secret)) < 32 {
		return fmt.Errorf("%s must be at least 32 bytes of high-entropy material", name)
	}
	return nil
}
func deriveKey(secret string) []byte {
	raw, err := base64.StdEncoding.DecodeString(secret)
	if err == nil && len(raw) >= 32 {
		sum := sha256.Sum256(raw)
		return sum[:]
	}
	sum := sha256.Sum256([]byte(secret))
	return sum[:]
}
func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			return trimmed
		}
	}
	return ""
}
func cookieNameForMode(explicit string, secure bool, secureDefault string, insecureDefault string) string {
	if trimmed := strings.TrimSpace(explicit); trimmed != "" {
		return trimmed
	}
	if secure {
		return secureDefault
	}
	return insecureDefault
}
func splitCSV(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		out = append(out, trimmed)
	}
	return out
}
func parseDurationMinutes(raw string, fallback int) time.Duration {
	v := parseIntDefault(raw, fallback)
	return time.Duration(v) * time.Minute
}
func parseDurationHours(raw string, fallback int) time.Duration {
	v := parseIntDefault(raw, fallback)
	return time.Duration(v) * time.Hour
}
func parseBoolDefault(raw string, fallback bool) bool {
	if strings.TrimSpace(raw) == "" {
		return fallback
	}
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return fallback
	}
}
func parseIntDefault(raw string, fallback int) int {
	if strings.TrimSpace(raw) == "" {
		return fallback
	}
	v, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil {
		return fallback
	}
	return v
}
