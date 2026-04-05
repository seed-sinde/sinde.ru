package auth
import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"image/png"
	"net"
	"net/mail"
	"net/url"
	"reflect"
	"strings"
	"time"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/redis/go-redis/v9"
	"sinde.ru/db"
	dbservices "sinde.ru/db/services"
	"sinde.ru/internal/media"
	"sinde.ru/internal/models"
	"sinde.ru/internal/store"
)
type Service struct {
	cfg        Config
	repo       *Repository
	redis      *redis.Client
	mailer     Mailer
	reputation IPReputation
	now        func() time.Time
}
type Dependencies struct {
	Redis      *redis.Client
	Mailer     Mailer
	Reputation IPReputation
	Now        func() time.Time
}
func NewService(repo *Repository, cfg Config, deps Dependencies) (*Service, error) {
	rdb := deps.Redis
	if rdb == nil {
		rdb = redis.NewClient(&redis.Options{
			Addr:     cfg.RedisAddr,
			Password: cfg.RedisPassword,
			DB:       cfg.RedisDB,
		})
	}
	mailer := deps.Mailer
	if mailer == nil {
		mailer = NewMailer(cfg)
	}
	reputation := deps.Reputation
	if reputation == nil {
		reputation = NoopReputation{}
	}
	nowFn := deps.Now
	if nowFn == nil {
		nowFn = time.Now
	}
	return &Service{
		cfg:        cfg,
		repo:       repo,
		redis:      rdb,
		mailer:     mailer,
		reputation: reputation,
		now:        nowFn,
	}, nil
}
func (s *Service) Close() error {
	if s.redis == nil {
		return nil
	}
	return s.redis.Close()
}
func (s *Service) AccessCookieName() string {
	return s.cfg.AccessCookieName
}
func (s *Service) RefreshCookieName() string {
	return s.cfg.RefreshCookieName
}
func (s *Service) CSRFCookieName() string {
	return s.cfg.CSRFCookieName
}
func (s *Service) SessionHintCookieName() string {
	return s.cfg.SessionHintCookieName
}
func (s *Service) SetSessionCookies(c fiber.Ctx, bundle *TokenBundle) {
	s.setAuthCookies(c, bundle)
}
func (s *Service) ClearSessionCookies(c fiber.Ctx) {
	s.clearAuthCookies(c)
}
func (s *Service) IsAllowedOrigin(raw string) bool {
	value := strings.TrimSpace(raw)
	if value == "" {
		return true
	}
	for _, item := range s.cfg.AllowedOrigins {
		if strings.EqualFold(strings.TrimSpace(item), value) {
			return true
		}
	}
	return false
}
func (s *Service) Register(ctx context.Context, input RegisterInput, device DeviceContext) (*RegisterResult, error) {
	if strings.TrimSpace(input.Honeypot) != "" {
		return nil, ErrRegistrationBlocked
	}
	if err := s.rateLimit(ctx, "register:ip", device.IP, s.cfg.RegisterRateLimit); err != nil {
		return nil, err
	}
	email, err := normalizeEmail(input.Email)
	if err != nil {
		return nil, err
	}
	if err := ValidatePasswordPolicy(input.Password, s.cfg.PasswordMinLength, email); err != nil {
		return nil, err
	}
	if existingUser, err := s.repo.GetUserByEmail(ctx, email); err == nil {
		if existingUser.EmailVerifiedAt != nil {
			return nil, ErrEmailAlreadyExists
		}
		if err := s.issueActionTokenEmail(ctx, existingUser, "verify_email", s.cfg.VerifyTTL, "/auth/verify-email"); err != nil {
			return nil, err
		}
		return &RegisterResult{
			User:            ToAuthUser(existingUser),
			VerificationTTL: s.cfg.VerifyTTL.String(),
		}, nil
	}
	passwordHash, err := HashPassword(input.Password)
	if err != nil {
		return nil, err
	}
	user := &models.User{
		UserID:             uuid.New(),
		Email:              email,
		EmailNormalized:    email,
		PasswordHash:       passwordHash,
		Status:             "pending_verification",
		DisplayName:        strings.TrimSpace(input.DisplayName),
		Locale:             fallbackString(input.Locale, "ru-RU"),
		Timezone:           fallbackString(input.Timezone, "UTC"),
		Roles:              []string{"user"},
		IsTwoFactorEnabled: false,
		Profile:            []byte("{}"),
		Settings:           []byte("{}"),
		SessionVersion:     1,
	}
	if err := s.repo.CreateUser(ctx, user); err != nil {
		if isUniqueViolation(err) {
			return nil, ErrEmailAlreadyExists
		}
		return nil, err
	}
	if err := s.issueActionTokenEmail(ctx, user, "verify_email", s.cfg.VerifyTTL, "/auth/verify-email"); err != nil {
		return nil, err
	}
	_ = s.logSecurityEvent(ctx, &user.UserID, nil, "auth", "auth.registered", "info", device, map[string]any{
		"email": email,
	})
	return &RegisterResult{
		User:            ToAuthUser(user),
		VerificationTTL: s.cfg.VerifyTTL.String(),
	}, nil
}
func (s *Service) RequestEmailVerification(ctx context.Context, email string) error {
	normalized, err := normalizeEmail(email)
	if err != nil {
		return nil
	}
	if err := s.rateLimit(ctx, "verify:mail", normalized, s.cfg.VerifyRateLimit); err != nil {
		return err
	}
	user, err := s.repo.GetUserByEmail(ctx, normalized)
	if err != nil || user.EmailVerifiedAt != nil {
		return nil
	}
	return s.issueActionTokenEmail(ctx, user, "verify_email", s.cfg.VerifyTTL, "/auth/verify-email")
}
func (s *Service) VerifyEmail(ctx context.Context, token string) error {
	action, err := s.repo.ConsumeActionToken(ctx, "verify_email", hashToken(strings.TrimSpace(token)), s.now())
	if err != nil {
		return err
	}
	return s.repo.MarkUserEmailVerified(ctx, action.UserID, s.now())
}
func (s *Service) Login(ctx context.Context, input LoginInput, device DeviceContext) (*LoginResult, *TokenBundle, error) {
	email, err := normalizeEmail(input.Email)
	if err != nil {
		return nil, nil, ErrInvalidCredentials
	}
	if err := s.rateLimit(ctx, "login:ip", device.IP, s.cfg.LoginRateLimitIP); err != nil {
		return nil, nil, err
	}
	if err := s.rateLimit(ctx, "login:mail", email, s.cfg.LoginRateLimitMail); err != nil {
		return nil, nil, err
	}
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil || !VerifyPassword(input.Password, user.PasswordHash) {
		_ = s.logLoginAttempt(ctx, nil, email, "failed", "invalid_credentials", device, 0, "")
		return nil, nil, ErrInvalidCredentials
	}
	if user.BlockedAt != nil || user.Status == "blocked" {
		_ = s.logLoginAttempt(ctx, &user.UserID, email, "failed", "blocked", device, 0, "")
		return nil, nil, ErrUserBlocked
	}
	if user.EmailVerifiedAt == nil {
		_ = s.logLoginAttempt(ctx, &user.UserID, email, "failed", "email_not_verified", device, 0, "")
		return nil, nil, ErrEmailNotVerified
	}
	riskScore, suspiciousReason, _ := s.reputation.Score(ctx, device.IP)
	if user.IsTwoFactorEnabled {
		ticket, expiresAt, err := s.issueMFATicket(ctx, user.UserID, device)
		if err != nil {
			return nil, nil, err
		}
		_ = s.logLoginAttempt(ctx, &user.UserID, email, "mfa_pending", "", device, riskScore, suspiciousReason)
		return &LoginResult{
			MFARequired:  true,
			MFATicket:    ticket,
			MFAExpiresAt: &expiresAt,
			MFAMethods:   []string{"totp", "backup_code"},
		}, nil, nil
	}
	result, bundle, err := s.createSession(ctx, user, device, true)
	if err != nil {
		return nil, nil, err
	}
	_ = s.logLoginAttempt(ctx, &user.UserID, email, "success", "", device, riskScore, suspiciousReason)
	return result, bundle, nil
}
func (s *Service) CompleteMFALogin(ctx context.Context, input TwoFactorLoginInput, device DeviceContext) (*LoginResult, *TokenBundle, error) {
	if err := s.rateLimit(ctx, "mfa:ip", device.IP, s.cfg.MFALoginRateLimitIP); err != nil {
		return nil, nil, err
	}
	state, err := s.getMFATicketState(ctx, input.Ticket)
	if err != nil {
		return nil, nil, err
	}
	if err := s.rateLimit(ctx, "mfa:ticket", input.Ticket, s.cfg.MFALoginRateLimitTicket); err != nil {
		return nil, nil, err
	}
	if state.FingerprintHash != device.FingerprintHash || state.UserAgent != device.UserAgent {
		return nil, nil, ErrSuspiciousLogin
	}
	userID, err := uuid.Parse(state.UserID)
	if err != nil {
		return nil, nil, ErrInvalidToken
	}
	code, ok := normalizeMFALoginCode(input.Code)
	if !ok {
		return nil, nil, ErrInvalidMFACode
	}
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, nil, err
	}
	ok, err = s.verifySecondFactor(ctx, user, code)
	if err != nil {
		return nil, nil, err
	}
	if !ok {
		return nil, nil, ErrInvalidMFACode
	}
	_ = s.clearMFATicket(ctx, input.Ticket)
	return s.createSession(ctx, user, device, true)
}
func (s *Service) Refresh(ctx context.Context, refreshToken string, csrfToken string, device DeviceContext) (*LoginResult, *TokenBundle, error) {
	refreshToken = strings.TrimSpace(refreshToken)
	csrfToken = strings.TrimSpace(csrfToken)
	if refreshToken == "" || csrfToken == "" {
		return nil, nil, ErrInvalidRefreshToken
	}
	claims, err := s.parseToken(refreshToken, TokenTypeRefresh)
	if err != nil {
		return nil, nil, ErrInvalidRefreshToken
	}
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, nil, ErrInvalidRefreshToken
	}
	sessionID, err := uuid.Parse(claims.SessionID)
	if err != nil {
		return nil, nil, ErrInvalidRefreshToken
	}
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, nil, err
	}
	if user.BlockedAt != nil || user.Status == "blocked" {
		return nil, nil, ErrUserBlocked
	}
	if claims.SessionVersion != user.SessionVersion {
		return nil, nil, ErrUnauthorized
	}
	session, err := s.repo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return nil, nil, ErrInvalidRefreshToken
	}
	if session.RevokedAt != nil {
		_ = s.repo.RevokeSessionFamily(ctx, session.FamilyID, "refresh_reuse_detected", s.now())
		return nil, nil, ErrSessionCompromised
	}
	if session.ExpiresAt.Before(s.now()) {
		_ = s.repo.RevokeSessionFamily(ctx, session.FamilyID, "expired_session_refresh", s.now())
		return nil, nil, ErrInvalidRefreshToken
	}
	if session.RefreshJTI != claims.ID ||
		session.RefreshTokenHash != hashToken(refreshToken) ||
		session.CSRFTokenHash != hashToken(csrfToken) {
		_ = s.repo.RevokeSessionFamily(ctx, session.FamilyID, "refresh_reuse_detected", s.now())
		return nil, nil, ErrInvalidRefreshToken
	}
	nextSession := &models.AuthSession{
		SessionID:        uuid.New(),
		UserID:           user.UserID,
		FamilyID:         session.FamilyID,
		RefreshJTI:       "",
		RefreshTokenHash: "",
		CSRFTokenHash:    "",
		FingerprintHash:  device.FingerprintHash,
		DeviceLabel:      device.DeviceLabel,
		IP:               device.IP,
		UserAgent:        device.UserAgent,
		MFAVerified:      session.MFAVerified,
		CreatedAt:        s.now(),
		UpdatedAt:        s.now(),
		LastSeenAt:       s.now(),
		ExpiresAt:        s.now().Add(s.cfg.RefreshTTL),
		RevokeReason:     "",
	}
	bundle, err := s.issueTokens(user, nextSession.SessionID, session.FamilyID, device, session.MFAVerified)
	if err != nil {
		return nil, nil, err
	}
	nextSession.RefreshJTI = bundle.RefreshJTI
	nextSession.RefreshTokenHash = hashToken(bundle.RefreshToken)
	nextSession.CSRFTokenHash = hashToken(bundle.CSRFToken)
	nextSession.ExpiresAt = bundle.RefreshExpiresAt
	if err := s.repo.RotateSession(ctx, session.SessionID, nextSession, s.now()); err != nil {
		return nil, nil, err
	}
	result := buildLoginResult(user, nextSession, nil)
	return result, bundle, nil
}
func (s *Service) Logout(ctx context.Context, accessToken string) error {
	if strings.TrimSpace(accessToken) == "" {
		return nil
	}
	claims, err := s.parseToken(accessToken, TokenTypeAccess)
	if err != nil {
		return nil
	}
	sessionID, err := uuid.Parse(claims.SessionID)
	if err == nil {
		_ = s.repo.RevokeSession(ctx, sessionID, "logout", s.now())
	}
	_ = s.blacklistAccessToken(ctx, claims)
	return nil
}
func (s *Service) LogoutAll(ctx context.Context, userID uuid.UUID, accessClaims *Claims) error {
	if err := s.repo.IncrementSessionVersion(ctx, userID); err != nil {
		return err
	}
	if err := s.repo.RevokeAllSessionsForUser(ctx, userID, "logout_all", s.now()); err != nil {
		return err
	}
	if accessClaims != nil {
		_ = s.blacklistAccessToken(ctx, accessClaims)
	}
	return nil
}
func (s *Service) ForgotPassword(ctx context.Context, email string) error {
	normalized, err := normalizeEmail(email)
	if err != nil {
		return nil
	}
	if err := s.rateLimit(ctx, "forgot:mail", normalized, s.cfg.ResetRateLimit); err != nil {
		return err
	}
	user, err := s.repo.GetUserByEmail(ctx, normalized)
	if err != nil {
		return nil
	}
	return s.issueActionTokenEmail(ctx, user, "password_reset", s.cfg.ResetTTL, "/auth/reset-password")
}
func (s *Service) ResetPassword(ctx context.Context, input ResetPasswordInput) error {
	token := strings.TrimSpace(input.Token)
	if err := ValidatePasswordPolicy(input.NewPassword, s.cfg.PasswordMinLength, ""); err != nil {
		return err
	}
	action, err := s.repo.ConsumeActionToken(ctx, "password_reset", hashToken(token), s.now())
	if err != nil {
		return err
	}
	hash, err := HashPassword(input.NewPassword)
	if err != nil {
		return err
	}
	if err := s.repo.UpdatePassword(ctx, action.UserID, hash); err != nil {
		return err
	}
	if err := s.repo.IncrementSessionVersion(ctx, action.UserID); err != nil {
		return err
	}
	return s.repo.RevokeAllSessionsForUser(ctx, action.UserID, "password_reset", s.now())
}
func (s *Service) ChangePassword(ctx context.Context, user *models.User, input ChangePasswordInput, accessClaims *Claims) error {
	if !VerifyPassword(input.CurrentPassword, user.PasswordHash) {
		return ErrPasswordChangeDenied
	}
	if err := ValidatePasswordPolicy(input.NewPassword, s.cfg.PasswordMinLength, user.Email); err != nil {
		return err
	}
	hash, err := HashPassword(input.NewPassword)
	if err != nil {
		return err
	}
	if err := s.repo.UpdatePassword(ctx, user.UserID, hash); err != nil {
		return err
	}
	if err := s.repo.IncrementSessionVersion(ctx, user.UserID); err != nil {
		return err
	}
	if err := s.repo.RevokeAllSessionsForUser(ctx, user.UserID, "password_change", s.now()); err != nil {
		return err
	}
	if accessClaims != nil {
		_ = s.blacklistAccessToken(ctx, accessClaims)
	}
	return nil
}
func (s *Service) SetupTwoFactor(ctx context.Context, user *models.User) (*TwoFactorSetupResult, error) {
	if user.IsTwoFactorEnabled {
		_ = s.clearTwoFactorSetupCache(ctx, user.UserID)
		return nil, ErrTwoFactorAlreadyOn
	}
	if cached, err := s.getTwoFactorSetupCache(ctx, user.UserID); err == nil && cached != nil {
		if len(user.TOTPSecretCiphertext) == 0 {
			_ = s.clearTwoFactorSetupCache(ctx, user.UserID)
		} else {
			secret, decryptErr := decryptString(user.TOTPSecretCiphertext, s.cfg.EncryptionKey)
			if decryptErr != nil {
				return nil, decryptErr
			}
			return &TwoFactorSetupResult{
				Secret:      secret,
				OTPAuthURL:  cached.OTPAuthURL,
				QRDataURL:   cached.QRDataURL,
				BackupCodes: nil,
			}, nil
		}
	}
	if err := s.rateLimit(ctx, "2fa-setup:user", user.UserID.String(), s.cfg.TwoFactorSetupRateLimit); err != nil {
		return nil, err
	}
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      s.cfg.TOTPIssuer,
		AccountName: user.Email,
	})
	if err != nil {
		return nil, err
	}
	ciphertext, err := encryptString(key.Secret(), s.cfg.EncryptionKey)
	if err != nil {
		return nil, err
	}
	if err := s.repo.UpdateTOTPSecret(ctx, user.UserID, ciphertext); err != nil {
		return nil, err
	}
	qrDataURL, err := buildQRCodeDataURL(key.URL())
	if err != nil {
		return nil, err
	}
	result := &TwoFactorSetupResult{
		Secret:      key.Secret(),
		OTPAuthURL:  key.URL(),
		QRDataURL:   qrDataURL,
		BackupCodes: nil,
	}
	_ = s.storeTwoFactorSetupCache(ctx, user.UserID, result)
	return result, nil
}
func (s *Service) EnableTwoFactor(ctx context.Context, user *models.User, code string) ([]string, error) {
	if user.IsTwoFactorEnabled {
		return nil, ErrTwoFactorAlreadyOn
	}
	if len(user.TOTPSecretCiphertext) == 0 {
		return nil, ErrTwoFactorNotSetup
	}
	secret, err := decryptString(user.TOTPSecretCiphertext, s.cfg.EncryptionKey)
	if err != nil {
		return nil, err
	}
	ok, err := totp.ValidateCustom(strings.TrimSpace(code), secret, s.now(), totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})
	if err != nil || !ok {
		return nil, ErrInvalidMFACode
	}
	plainCodes, codeHashes, err := generateBackupCodes()
	if err != nil {
		return nil, err
	}
	if err := s.repo.ReplaceBackupCodes(ctx, user.UserID, codeHashes); err != nil {
		return nil, err
	}
	if err := s.repo.EnableTwoFactor(ctx, user.UserID); err != nil {
		return nil, err
	}
	_ = s.clearTwoFactorSetupCache(ctx, user.UserID)
	return plainCodes, nil
}
func (s *Service) DisableTwoFactor(ctx context.Context, user *models.User, input TwoFactorDisableInput) error {
	if !user.IsTwoFactorEnabled {
		return ErrTwoFactorNotEnabled
	}
	if !VerifyPassword(input.Password, user.PasswordHash) {
		return ErrInvalidCredentials
	}
	ok, err := s.verifySecondFactor(ctx, user, strings.TrimSpace(input.Code))
	if err != nil || !ok {
		return ErrInvalidMFACode
	}
	if err := s.repo.ReplaceBackupCodes(ctx, user.UserID, nil); err != nil {
		return err
	}
	_ = s.clearTwoFactorSetupCache(ctx, user.UserID)
	return s.repo.DisableTwoFactor(ctx, user.UserID)
}
func (s *Service) AuthenticateAccessToken(ctx context.Context, raw string, device DeviceContext) (*models.User, *Claims, error) {
	claims, err := s.parseToken(strings.TrimSpace(raw), TokenTypeAccess)
	if err != nil {
		return nil, nil, ErrUnauthorized
	}
	denied, err := s.redis.Exists(ctx, "auth:deny:access:"+claims.ID).Result()
	if err == nil && denied > 0 {
		return nil, nil, ErrUnauthorized
	}
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, nil, ErrUnauthorized
	}
	sessionID, err := uuid.Parse(claims.SessionID)
	if err != nil {
		return nil, nil, ErrUnauthorized
	}
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, nil, ErrUnauthorized
	}
	if user.BlockedAt != nil || user.Status == "blocked" {
		return nil, nil, ErrUserBlocked
	}
	if claims.SessionVersion != user.SessionVersion {
		return nil, nil, ErrUnauthorized
	}
	session, err := s.repo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return nil, nil, ErrUnauthorized
	}
	if session.UserID != user.UserID || session.RevokedAt != nil || session.ExpiresAt.Before(s.now()) {
		return nil, nil, ErrUnauthorized
	}
	_ = s.repo.TouchSession(ctx, sessionID, s.now(), device.IP, device.UserAgent)
	return user, claims, nil
}
func (s *Service) ListSessions(ctx context.Context, userID uuid.UUID, currentSessionID string) (*SessionListResult, error) {
	items, err := s.repo.ListSessionsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]SessionView, 0, len(items))
	for _, item := range items {
		out = append(out, SessionView{
			SessionID:    item.SessionID,
			DeviceLabel:  item.DeviceLabel,
			IP:           item.IP,
			UserAgent:    item.UserAgent,
			MFAVerified:  item.MFAVerified,
			LastSeenAt:   item.LastSeenAt,
			ExpiresAt:    item.ExpiresAt,
			RevokedAt:    item.RevokedAt,
			RevokeReason: item.RevokeReason,
			CreatedAt:    item.CreatedAt,
			IsCurrent:    item.SessionID.String() == currentSessionID,
		})
	}
	return &SessionListResult{Items: out}, nil
}
func (s *Service) RevokeUserSession(ctx context.Context, requester *models.User, sessionID uuid.UUID, currentAccess *Claims) error {
	session, err := s.repo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return err
	}
	if session.UserID != requester.UserID {
		return ErrForbidden
	}
	if err := s.repo.RevokeSession(ctx, sessionID, "user_revoked", s.now()); err != nil {
		return err
	}
	if currentAccess != nil && currentAccess.SessionID == sessionID.String() {
		_ = s.blacklistAccessToken(ctx, currentAccess)
	}
	return nil
}
func (s *Service) UpdateProfile(ctx context.Context, userID uuid.UUID, input UpdateProfileInput) (*UpdateProfileResult, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	displayName := user.DisplayName
	if input.DisplayName != nil {
		nextDisplayName := strings.TrimSpace(*input.DisplayName)
		if nextDisplayName != "" {
			displayName = nextDisplayName
		}
	}
	locale := user.Locale
	if input.Locale != nil {
		locale = fallbackString(*input.Locale, user.Locale)
	}
	timezone := user.Timezone
	if input.Timezone != nil {
		timezone = fallbackString(*input.Timezone, user.Timezone)
	}
	profile, err := mergeJSONObjects(user.Profile, input.Profile)
	if err != nil {
		return nil, ErrInvalidInput
	}
	settings, err := mergeJSONObjects(user.Settings, input.Settings)
	if err != nil {
		return nil, ErrInvalidInput
	}
	beforeAvatarKeys := extractAvatarImageKeys(user.Profile, user.Settings)
	afterAvatarKeys := extractAvatarImageKeys(profile, settings)
	if err := s.repo.UpdateProfile(ctx, userID, displayName, locale, timezone, profile, settings); err != nil {
		return nil, err
	}
	for _, key := range diffStringSets(beforeAvatarKeys, afterAvatarKeys) {
		if cleanupErr := media.DeleteIfUnreferenced(ctx, key); cleanupErr != nil {
			return nil, cleanupErr
		}
	}
	return &UpdateProfileResult{
		Changes: buildAuthUserPatch(user, displayName, locale, timezone, profile, settings, input),
	}, nil
}
func buildAuthUserPatch(
	current *models.User,
	displayName string,
	locale string,
	timezone string,
	profile json.RawMessage,
	settings json.RawMessage,
	input UpdateProfileInput,
) AuthUserPatch {
	patch := AuthUserPatch{}
	if input.DisplayName != nil && displayName != current.DisplayName {
		patch.DisplayName = &displayName
	}
	if input.Locale != nil && locale != current.Locale {
		patch.Locale = &locale
	}
	if input.Timezone != nil && timezone != current.Timezone {
		patch.Timezone = &timezone
	}
	if nextProfilePatch := buildJSONObjectPatch(current.Profile, profile, input.Profile); nextProfilePatch != nil {
		patch.Profile = nextProfilePatch
	}
	if nextSettingsPatch := buildJSONObjectPatch(current.Settings, settings, input.Settings); nextSettingsPatch != nil {
		patch.Settings = nextSettingsPatch
	}
	return patch
}
func buildJSONObjectPatch(current json.RawMessage, merged json.RawMessage, patch *json.RawMessage) map[string]any {
	if patch == nil {
		return nil
	}
	currentValue, err := decodeJSONObject(jsonOrDefault(current))
	if err != nil {
		return nil
	}
	mergedValue, err := decodeJSONObject(jsonOrDefault(merged))
	if err != nil {
		return nil
	}
	if reflect.DeepEqual(currentValue, mergedValue) {
		return nil
	}
	patchValue, err := decodeJSONObject(jsonOrDefault(*patch))
	if err != nil || len(patchValue) == 0 {
		return nil
	}
	return patchValue
}
func mergeJSONObjects(current json.RawMessage, patch *json.RawMessage) (json.RawMessage, error) {
	if patch == nil {
		return jsonOrDefault(current), nil
	}
	currentValue, err := decodeJSONObject(jsonOrDefault(current))
	if err != nil {
		return nil, err
	}
	patchValue, err := decodeJSONObject(jsonOrDefault(*patch))
	if err != nil {
		return nil, err
	}
	merged := mergeJSONObjectValue(currentValue, patchValue)
	raw, err := json.Marshal(merged)
	if err != nil {
		return nil, err
	}
	return raw, nil
}
func decodeJSONObject(raw json.RawMessage) (map[string]any, error) {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		return map[string]any{}, nil
	}
	var value map[string]any
	if err := json.Unmarshal(trimmed, &value); err != nil {
		return nil, err
	}
	if value == nil {
		return map[string]any{}, nil
	}
	return value, nil
}
func mergeJSONObjectValue(current map[string]any, patch map[string]any) map[string]any {
	merged := make(map[string]any, len(current))
	for key, value := range current {
		merged[key] = value
	}
	for key, value := range patch {
		if value == nil {
			delete(merged, key)
			continue
		}
		patchMap, patchIsMap := value.(map[string]any)
		currentMap, currentIsMap := merged[key].(map[string]any)
		if patchIsMap && currentIsMap {
			merged[key] = mergeJSONObjectValue(currentMap, patchMap)
			continue
		}
		merged[key] = value
	}
	return merged
}
func extractAvatarImageKeys(profile json.RawMessage, settings json.RawMessage) []string {
	keys := make(map[string]struct{})
	collectAvatarKeys(keys, profile)
	collectAvatarKeys(keys, settings)
	out := make([]string, 0, len(keys))
	for key := range keys {
		out = append(out, key)
	}
	return out
}
func collectAvatarKeys(keys map[string]struct{}, raw json.RawMessage) {
	if len(raw) == 0 || !json.Valid(raw) {
		return
	}
	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		return
	}
	avatar, ok := payload["avatar"].(map[string]any)
	if !ok {
		return
	}
	collectNestedAvatarKeys(keys, avatar)
}
func collectNestedAvatarKeys(keys map[string]struct{}, value any) {
	switch typed := value.(type) {
	case map[string]any:
		for _, field := range []string{"icon_image_key", "profile_image_key", "original_image_key"} {
			key := media.NormalizeStorageKey(fmt.Sprint(typed[field]))
			if key != "" {
				keys[key] = struct{}{}
			}
		}
		for _, nested := range typed {
			collectNestedAvatarKeys(keys, nested)
		}
	case []any:
		for _, item := range typed {
			collectNestedAvatarKeys(keys, item)
		}
	}
}
func diffStringSets(before []string, after []string) []string {
	afterSet := make(map[string]struct{}, len(after))
	for _, key := range after {
		afterSet[key] = struct{}{}
	}
	removed := make([]string, 0, len(before))
	for _, key := range before {
		if _, keep := afterSet[key]; keep {
			continue
		}
		removed = append(removed, key)
	}
	return removed
}
func (s *Service) ListSavedTraitSets(ctx context.Context, userID uuid.UUID) (*SavedTraitSetListResult, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	items, err := s.repo.ListSavedTraitSets(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]SavedTraitSetView, 0, len(items))
	for _, item := range items {
		exists, existsErr := s.ensureTraitTargetPersisted(ctx, item.TraitUUID)
		if existsErr != nil {
			return nil, existsErr
		}
		if !exists {
			if deleteErr := s.repo.DeleteSavedTraitSetByTraitUUID(ctx, userID, item.TraitUUID); deleteErr != nil {
				return nil, deleteErr
			}
			continue
		}
		out = append(out, ToSavedTraitSetView(item))
	}
	if user.PrimaryTraitUUID != nil {
		exists, existsErr := s.ensureTraitTargetPersisted(ctx, *user.PrimaryTraitUUID)
		if existsErr != nil {
			return nil, existsErr
		}
		if !exists {
			if err := s.repo.UpdatePrimaryTraitUUID(ctx, userID, nil); err != nil {
				return nil, err
			}
			user.PrimaryTraitUUID = nil
		}
	}
	return &SavedTraitSetListResult{
		PrimaryTraitUUID: user.PrimaryTraitUUID,
		Items:            out,
	}, nil
}
func (s *Service) SetPrimaryTrait(ctx context.Context, userID uuid.UUID, rawTraitUUID string) (*MeResult, error) {
	trimmed := strings.TrimSpace(rawTraitUUID)
	if trimmed == "" {
		if err := s.repo.UpdatePrimaryTraitUUID(ctx, userID, nil); err != nil {
			return nil, err
		}
		updated, err := s.repo.GetUserByID(ctx, userID)
		if err != nil {
			return nil, err
		}
		return &MeResult{User: ToAuthUser(updated)}, nil
	}
	traitUUID, err := uuid.Parse(trimmed)
	if err != nil {
		return nil, ErrInvalidInput
	}
	exists, err := s.ensureTraitTargetPersisted(ctx, traitUUID)
	if err != nil {
		return nil, err
	}
	if !exists {
		if err := s.repo.DeleteSavedTraitSetByTraitUUID(ctx, userID, traitUUID); err != nil {
			return nil, err
		}
		return nil, ErrTraitTargetNotFound
	}
	if err := s.repo.UpdatePrimaryTraitUUID(ctx, userID, &traitUUID); err != nil {
		return nil, err
	}
	updated, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &MeResult{User: ToAuthUser(updated)}, nil
}
func (s *Service) SaveTraitSet(ctx context.Context, userID uuid.UUID, input SaveTraitSetInput) (*SavedTraitSetView, error) {
	traitUUID, name, description, err := validateSavedTraitSetInput(input.TraitUUID, input.Name, input.Description)
	if err != nil {
		return nil, err
	}
	exists, err := s.ensureTraitTargetPersisted(ctx, traitUUID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrInvalidInput
	}
	item, err := s.repo.UpsertSavedTraitSet(ctx, &models.SavedTraitSet{
		SavedSetID:  uuid.New(),
		UserID:      userID,
		TraitUUID:   traitUUID,
		Name:        name,
		Description: description,
	})
	if err != nil {
		return nil, err
	}
	view := ToSavedTraitSetView(*item)
	return &view, nil
}
func (s *Service) ensureTraitTargetPersisted(ctx context.Context, id uuid.UUID) (bool, error) {
	exists, err := dbservices.PdbTraitTargetExists(ctx, id)
	if err != nil || exists {
		return exists, err
	}
	if _, ok := store.GetTrait(id); ok {
		if err := s.persistTraitFromMemory(ctx, id); err != nil {
			return false, err
		}
		return dbservices.PdbTraitTargetExists(ctx, id)
	}
	if _, ok := store.GetSet(id); ok {
		if err := s.persistSetFromMemory(ctx, id); err != nil {
			return false, err
		}
		return dbservices.PdbTraitTargetExists(ctx, id)
	}
	return false, nil
}
func (s *Service) persistTraitFromMemory(ctx context.Context, traitUUID uuid.UUID) error {
	trait, ok := store.GetTrait(traitUUID)
	if !ok || trait == nil {
		return nil
	}
	conn, err := db.PDB.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()
	keys, ok := store.GetKeys(trait.TKey)
	if !ok || len(keys) == 0 || keys[0] == nil {
		return fmt.Errorf("missing key for trait %s", traitUUID)
	}
	if err := dbservices.PdbInsertTraitKey(conn, keys[0]); err != nil {
		return err
	}
	return dbservices.PdbInsertTrait(conn, trait)
}
func (s *Service) persistSetFromMemory(ctx context.Context, setUUID uuid.UUID) error {
	setItem, ok := store.GetSet(setUUID)
	if !ok || setItem == nil {
		return nil
	}
	if len(setItem.SChilds) != 2 {
		return fmt.Errorf("invalid childs len for %s: %v", setItem.SUUID, setItem.SChilds)
	}
	for _, childUUID := range setItem.SChilds {
		childExists, err := s.ensureTraitTargetPersisted(ctx, childUUID)
		if err != nil {
			return err
		}
		if !childExists {
			return fmt.Errorf("missing child %s for set %s", childUUID, setUUID)
		}
	}
	return dbservices.PdbInsertSet(ctx, setItem)
}
func (s *Service) UpdateSavedTraitSet(ctx context.Context, userID uuid.UUID, savedSetID uuid.UUID, input UpdateSavedTraitSetInput) (*SavedTraitSetView, error) {
	name := normalizeSavedTraitSetName(input.Name)
	description := normalizeSavedTraitSetDescription(input.Description)
	if name == "" {
		return nil, ErrInvalidInput
	}
	item, err := s.repo.UpdateSavedTraitSet(ctx, userID, savedSetID, name, description)
	if err != nil {
		return nil, err
	}
	view := ToSavedTraitSetView(*item)
	return &view, nil
}
func (s *Service) DeleteSavedTraitSet(ctx context.Context, userID uuid.UUID, savedSetID uuid.UUID) error {
	return s.repo.DeleteSavedTraitSet(ctx, userID, savedSetID)
}
func (s *Service) ListLoginAttempts(ctx context.Context, user *models.User, limit int) (*LoginAttemptListResult, error) {
	items, err := s.repo.ListLoginAttempts(ctx, user.UserID, user.EmailNormalized, limit)
	if err != nil {
		return nil, err
	}
	out := make([]LoginAttemptView, 0, len(items))
	for _, item := range items {
		out = append(out, ToLoginAttemptView(item))
	}
	return &LoginAttemptListResult{Items: out}, nil
}
func normalizeSavedTraitSetName(value string) string {
	name := strings.TrimSpace(value)
	if name == "" {
		return ""
	}
	if len([]rune(name)) > 120 {
		return ""
	}
	return name
}
func normalizeSavedTraitSetDescription(value string) string {
	description := strings.TrimSpace(value)
	if len([]rune(description)) > 280 {
		return ""
	}
	return description
}
func validateSavedTraitSetInput(rawTraitUUID string, rawName string, rawDescription string) (uuid.UUID, string, string, error) {
	traitUUID, err := uuid.Parse(strings.TrimSpace(rawTraitUUID))
	if err != nil {
		return uuid.Nil, "", "", ErrInvalidInput
	}
	name := normalizeSavedTraitSetName(rawName)
	description := normalizeSavedTraitSetDescription(rawDescription)
	if name == "" {
		return uuid.Nil, "", "", ErrInvalidInput
	}
	if rawDescription != "" && description == "" {
		return uuid.Nil, "", "", ErrInvalidInput
	}
	return traitUUID, name, description, nil
}
func (s *Service) ListSecurityEvents(ctx context.Context, userID uuid.UUID, limit int) (*SecurityEventListResult, error) {
	items, err := s.repo.ListSecurityEvents(ctx, userID, limit)
	if err != nil {
		return nil, err
	}
	out := make([]SecurityEventView, 0, len(items))
	for _, item := range items {
		out = append(out, ToSecurityEventView(item))
	}
	return &SecurityEventListResult{Items: out}, nil
}
func (s *Service) BlockUser(ctx context.Context, actor *models.User, targetUserID uuid.UUID, reason string) error {
	if err := s.repo.SetBlockedState(ctx, targetUserID, true, strings.TrimSpace(reason), s.now()); err != nil {
		return err
	}
	if err := s.repo.IncrementSessionVersion(ctx, targetUserID); err != nil {
		return err
	}
	if err := s.repo.RevokeAllSessionsForUser(ctx, targetUserID, "blocked", s.now()); err != nil {
		return err
	}
	actorID := actor.UserID
	_ = s.logSecurityEvent(ctx, &targetUserID, &actorID, "admin", "auth.user.blocked", "warn", DeviceContext{}, map[string]any{
		"reason": reason,
	})
	return nil
}
func (s *Service) UnblockUser(ctx context.Context, actor *models.User, targetUserID uuid.UUID) error {
	if err := s.repo.SetBlockedState(ctx, targetUserID, false, "", s.now()); err != nil {
		return err
	}
	actorID := actor.UserID
	_ = s.logSecurityEvent(ctx, &targetUserID, &actorID, "admin", "auth.user.unblocked", "info", DeviceContext{}, nil)
	return nil
}
func (s *Service) ForceLogoutUser(ctx context.Context, actor *models.User, targetUserID uuid.UUID) error {
	if err := s.repo.IncrementSessionVersion(ctx, targetUserID); err != nil {
		return err
	}
	if err := s.repo.RevokeAllSessionsForUser(ctx, targetUserID, "admin_force_logout", s.now()); err != nil {
		return err
	}
	actorID := actor.UserID
	_ = s.logSecurityEvent(ctx, &targetUserID, &actorID, "admin", "auth.user.force_logout", "warn", DeviceContext{}, nil)
	return nil
}
func (s *Service) AdminDeleteUser(ctx context.Context, actor *models.User, targetUserID uuid.UUID) error {
	if actor != nil && actor.UserID == targetUserID {
		return ErrForbidden
	}
	target, err := s.repo.GetUserByID(ctx, targetUserID)
	if err != nil {
		if errors.Is(err, ErrUnauthorized) {
			return ErrUserNotFound
		}
		return err
	}
	if hasRole(target.Roles, "admin") {
		adminCount, countErr := s.repo.CountUsersByRole(ctx, "admin")
		if countErr != nil {
			return countErr
		}
		if adminCount <= 1 {
			return ErrForbidden
		}
	}
	if err := s.repo.DeleteUser(ctx, targetUserID); err != nil {
		return err
	}
	if actor != nil {
		actorID := actor.UserID
		_ = s.logSecurityEvent(ctx, nil, &actorID, "admin", "auth.user.deleted", "warn", DeviceContext{}, map[string]any{
			"target_user_id": targetUserID.String(),
			"target_email":   target.Email,
		})
	}
	return nil
}
func (s *Service) AdminListUsers(ctx context.Context, search string, status string, role string, limit int, offset int) (*AdminUserListResult, error) {
	items, err := s.repo.ListUsersForAdmin(ctx, strings.TrimSpace(search), strings.TrimSpace(status), strings.TrimSpace(role), limit, offset)
	if err != nil {
		return nil, err
	}
	total, err := s.repo.CountUsersForAdmin(ctx, strings.TrimSpace(search), strings.TrimSpace(status), strings.TrimSpace(role))
	if err != nil {
		return nil, err
	}
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}
	if offset < 0 {
		offset = 0
	}
	out := make([]AdminUserView, 0, len(items))
	for _, item := range items {
		out = append(out, AdminUserView{
			UserID:             item.UserID,
			Email:              item.Email,
			Status:             item.Status,
			DisplayName:        item.DisplayName,
			Locale:             item.Locale,
			Timezone:           item.Timezone,
			Roles:              append([]string(nil), item.Roles...),
			IsTwoFactorEnabled: item.IsTwoFactorEnabled,
			LastLoginAt:        item.LastLoginAt,
			BlockedReason:      item.BlockedReason,
			BlockedAt:          item.BlockedAt,
			CreatedAt:          item.CreatedAt,
		})
	}
	return &AdminUserListResult{
		Items:  out,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}
func (s *Service) AdminSummary(ctx context.Context, adminUserID uuid.UUID, currentSessionID string) (*AdminSummaryResult, error) {
	usersTotal, err := s.repo.CountUsersTotal(ctx)
	if err != nil {
		return nil, err
	}
	recipeStatusTotals, err := s.repo.CountKitchenRecipesByModerationStatus(ctx)
	if err != nil {
		return nil, err
	}
	pendingTotal := recipeStatusTotals["pending"]
	var sessionUUID *uuid.UUID
	if parsed, parseErr := uuid.Parse(strings.TrimSpace(currentSessionID)); parseErr == nil {
		sessionUUID = &parsed
	}
	lastLoginBeforeCurrent, err := s.repo.GetPreviousLoginAt(ctx, adminUserID, sessionUUID)
	if err != nil {
		return nil, err
	}
	notificationsReadAt, err := s.repo.GetAdminNotificationsReadAt(ctx, adminUserID)
	if err != nil {
		return nil, err
	}
	var notificationsSinceAt *time.Time
	if lastLoginBeforeCurrent != nil {
		copyTs := *lastLoginBeforeCurrent
		notificationsSinceAt = &copyTs
	}
	if notificationsReadAt != nil && (notificationsSinceAt == nil || notificationsReadAt.After(*notificationsSinceAt)) {
		copyTs := *notificationsReadAt
		notificationsSinceAt = &copyTs
	}
	newUsersSinceLastLogin := 0
	newPendingSinceLastLogin := 0
	if notificationsSinceAt != nil {
		newUsersSinceLastLogin, err = s.repo.CountUsersCreatedSince(ctx, *notificationsSinceAt)
		if err != nil {
			return nil, err
		}
		newPendingSinceLastLogin, err = s.repo.CountPendingKitchenRecipesCreatedSince(ctx, *notificationsSinceAt)
		if err != nil {
			return nil, err
		}
	}
	return &AdminSummaryResult{
		UsersTotal:                      usersTotal,
		PendingRecipesTotal:             pendingTotal,
		RecipeStatusTotals:              recipeStatusTotals,
		NewUsersSinceLastLogin:          newUsersSinceLastLogin,
		NewPendingRecipesSinceLastLogin: newPendingSinceLastLogin,
		LastLoginBeforeCurrentSessionAt: lastLoginBeforeCurrent,
		NotificationsReadAt:             notificationsReadAt,
		NotificationsSinceAt:            notificationsSinceAt,
		HasUnread:                       newUsersSinceLastLogin > 0 || newPendingSinceLastLogin > 0,
		CheckedAt:                       s.now(),
	}, nil
}
func (s *Service) UserSummary(ctx context.Context, userID uuid.UUID, currentSessionID string) (*UserSummaryResult, error) {
	var sessionUUID *uuid.UUID
	if parsed, parseErr := uuid.Parse(strings.TrimSpace(currentSessionID)); parseErr == nil {
		sessionUUID = &parsed
	}
	lastLoginBeforeCurrent, err := s.repo.GetPreviousLoginAt(ctx, userID, sessionUUID)
	if err != nil {
		return nil, err
	}
	notificationsReadAt, err := s.repo.GetUserNotificationsReadAt(ctx, userID)
	if err != nil {
		return nil, err
	}
	var notificationsSinceAt *time.Time
	if lastLoginBeforeCurrent != nil {
		copyTs := *lastLoginBeforeCurrent
		notificationsSinceAt = &copyTs
	}
	if notificationsReadAt != nil && (notificationsSinceAt == nil || notificationsReadAt.After(*notificationsSinceAt)) {
		copyTs := *notificationsReadAt
		notificationsSinceAt = &copyTs
	}
	newApprovedSinceLastLogin := 0
	newRejectedSinceLastLogin := 0
	if notificationsSinceAt != nil {
		eventTotals, countErr := s.repo.CountOwnedKitchenRecipeModerationEventsSince(ctx, userID, *notificationsSinceAt)
		if countErr != nil {
			return nil, countErr
		}
		newApprovedSinceLastLogin = eventTotals["approved"]
		newRejectedSinceLastLogin = eventTotals["rejected"]
	}
	return &UserSummaryResult{
		NewApprovedRecipesSinceLastLogin: newApprovedSinceLastLogin,
		NewRejectedRecipesSinceLastLogin: newRejectedSinceLastLogin,
		LastLoginBeforeCurrentSessionAt:  lastLoginBeforeCurrent,
		NotificationsReadAt:              notificationsReadAt,
		NotificationsSinceAt:             notificationsSinceAt,
		HasUnread:                        newApprovedSinceLastLogin > 0 || newRejectedSinceLastLogin > 0,
		CheckedAt:                        s.now(),
	}, nil
}
func (s *Service) MarkUserSummaryRead(ctx context.Context, userID uuid.UUID) error {
	return s.repo.SetUserNotificationsReadAt(ctx, userID, s.now())
}
func (s *Service) MarkAdminSummaryRead(ctx context.Context, adminUserID uuid.UUID) error {
	return s.repo.SetAdminNotificationsReadAt(ctx, adminUserID, s.now())
}
func (s *Service) AdminSetUserRole(ctx context.Context, actor *models.User, targetUserID uuid.UUID, role string) (*AuthUser, error) {
	nextRole := strings.ToLower(strings.TrimSpace(role))
	if nextRole != "user" && nextRole != "admin" {
		return nil, ErrInvalidInput
	}
	target, err := s.repo.GetUserByID(ctx, targetUserID)
	if err != nil {
		if errors.Is(err, ErrUnauthorized) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	if actor != nil && actor.UserID == targetUserID && nextRole != "admin" {
		return nil, ErrForbidden
	}
	currentIsAdmin := hasRole(target.Roles, "admin")
	if currentIsAdmin && nextRole != "admin" {
		adminCount, countErr := s.repo.CountUsersByRole(ctx, "admin")
		if countErr != nil {
			return nil, countErr
		}
		if adminCount <= 1 {
			return nil, ErrForbidden
		}
	}
	if err := s.repo.SetUserPrimaryRole(ctx, targetUserID, nextRole); err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	updated, err := s.repo.GetUserByID(ctx, targetUserID)
	if err != nil {
		return nil, err
	}
	if actor != nil {
		actorID := actor.UserID
		_ = s.logSecurityEvent(ctx, &targetUserID, &actorID, "admin", "auth.user.role_changed", "info", DeviceContext{}, map[string]any{
			"role": nextRole,
		})
	}
	item := ToAuthUser(updated)
	return &item, nil
}
func (s *Service) AdminSearchTraitKeys(ctx context.Context, query string, limit int) (*AdminTraitKeySearchResult, error) {
	items, err := s.repo.SearchTraitKeysForAdmin(ctx, strings.TrimSpace(query), limit)
	if err != nil {
		return nil, err
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	return &AdminTraitKeySearchResult{
		Items: items,
		Query: strings.TrimSpace(query),
		Limit: limit,
	}, nil
}
func (s *Service) AdminTraitsSetsAnalysis(ctx context.Context) (*AdminTraitsSetsAnalysis, error) {
	return s.repo.GetTraitsSetsAnalysisForAdmin(ctx)
}
func (s *Service) createSession(ctx context.Context, user *models.User, device DeviceContext, mfaVerified bool) (*LoginResult, *TokenBundle, error) {
	session := &models.AuthSession{
		SessionID:       uuid.New(),
		UserID:          user.UserID,
		FamilyID:        uuid.New(),
		CreatedAt:       s.now(),
		UpdatedAt:       s.now(),
		FingerprintHash: device.FingerprintHash,
		DeviceLabel:     device.DeviceLabel,
		IP:              device.IP,
		UserAgent:       device.UserAgent,
		MFAVerified:     mfaVerified,
		LastSeenAt:      s.now(),
		ExpiresAt:       s.now().Add(s.cfg.RefreshTTL),
		RevokeReason:    "",
	}
	bundle, err := s.issueTokens(user, session.SessionID, session.FamilyID, device, mfaVerified)
	if err != nil {
		return nil, nil, err
	}
	session.RefreshJTI = bundle.RefreshJTI
	session.RefreshTokenHash = hashToken(bundle.RefreshToken)
	session.CSRFTokenHash = hashToken(bundle.CSRFToken)
	session.ExpiresAt = bundle.RefreshExpiresAt
	if err := s.repo.CreateSession(ctx, session); err != nil {
		return nil, nil, err
	}
	_ = s.repo.UpdateUserLastLogin(ctx, user.UserID, s.now())
	seen, _ := s.repo.HasSeenFingerprint(ctx, user.UserID, device.FingerprintHash)
	eventType := "auth.login.success"
	if !seen {
		eventType = "auth.login.new_device"
	}
	_ = s.logSecurityEvent(ctx, &user.UserID, nil, "auth", eventType, "info", device, map[string]any{
		"session_id": session.SessionID.String(),
	})
	return buildLoginResult(user, session, bundle), bundle, nil
}
func buildLoginResult(user *models.User, session *models.AuthSession, bundle *TokenBundle) *LoginResult {
	result := &LoginResult{
		User: ptr(ToAuthUser(user)),
		Session: &SessionView{
			SessionID:    session.SessionID,
			DeviceLabel:  session.DeviceLabel,
			IP:           session.IP,
			UserAgent:    session.UserAgent,
			MFAVerified:  session.MFAVerified,
			LastSeenAt:   session.LastSeenAt,
			ExpiresAt:    session.ExpiresAt,
			RevokedAt:    session.RevokedAt,
			RevokeReason: session.RevokeReason,
			CreatedAt:    session.CreatedAt,
			IsCurrent:    true,
		},
	}
	if bundle != nil {
		result.CSRFToken = bundle.CSRFToken
	}
	return result
}
func normalizeMFALoginCode(raw string) (string, bool) {
	code := strings.TrimSpace(raw)
	if isSixDigitCode(code) {
		return code, true
	}
	code = strings.ToUpper(code)
	compact := strings.NewReplacer("-", "", " ", "").Replace(code)
	if len(compact) != 8 {
		return "", false
	}
	for _, ch := range compact {
		if (ch < '0' || ch > '9') && (ch < 'A' || ch > 'Z') {
			return "", false
		}
	}
	return compact[:4] + "-" + compact[4:], true
}
func isSixDigitCode(code string) bool {
	if len(code) != 6 {
		return false
	}
	for _, ch := range code {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}
func buildQRCodeDataURL(content string) (string, error) {
	qrCode, err := qr.Encode(content, qr.M, qr.Auto)
	if err != nil {
		return "", err
	}
	scaled, err := barcode.Scale(qrCode, 240, 240)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, scaled); err != nil {
		return "", err
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}
func (s *Service) verifySecondFactor(ctx context.Context, user *models.User, code string) (bool, error) {
	if len(user.TOTPSecretCiphertext) > 0 {
		secret, err := decryptString(user.TOTPSecretCiphertext, s.cfg.EncryptionKey)
		if err != nil {
			return false, err
		}
		ok, err := totp.ValidateCustom(code, secret, s.now(), totp.ValidateOpts{
			Period:    30,
			Skew:      1,
			Digits:    otp.DigitsSix,
			Algorithm: otp.AlgorithmSHA1,
		})
		if err == nil && ok {
			return true, nil
		}
	}
	codes, err := s.repo.ListBackupCodes(ctx, user.UserID)
	if err != nil {
		return false, err
	}
	for _, item := range codes {
		if VerifyPassword(code, item.CodeHash) {
			if err := s.repo.ConsumeBackupCode(ctx, item.BackupCodeID, s.now()); err != nil {
				return false, err
			}
			return true, nil
		}
	}
	return false, nil
}
func generateBackupCodes() ([]string, []string, error) {
	plain := make([]string, 0, 10)
	hashes := make([]string, 0, 10)
	for i := 0; i < 10; i++ {
		token, err := randomToken(6)
		if err != nil {
			return nil, nil, err
		}
		code := strings.ToUpper(strings.ReplaceAll(token[:8], "-", ""))
		if len(code) > 4 {
			code = fmt.Sprintf("%s-%s", code[:4], code[4:])
		}
		hash, err := HashPassword(code)
		if err != nil {
			return nil, nil, err
		}
		plain = append(plain, code)
		hashes = append(hashes, hash)
	}
	return plain, hashes, nil
}
func (s *Service) issueActionTokenEmail(ctx context.Context, user *models.User, purpose string, ttl time.Duration, path string) error {
	token, err := randomToken(32)
	if err != nil {
		return err
	}
	action := &models.ActionToken{
		TokenID:       uuid.New(),
		UserID:        user.UserID,
		Purpose:       purpose,
		TokenHash:     hashToken(token),
		ExpiresAt:     s.now().Add(ttl),
		DeliveryValue: user.Email,
		Meta:          []byte("{}"),
	}
	if err := s.repo.CreateActionToken(ctx, action); err != nil {
		return err
	}
	u, err := url.Parse(strings.TrimRight(s.cfg.PublicBaseURL, "/") + path)
	if err != nil {
		return err
	}
	q := u.Query()
	q.Set("token", token)
	u.RawQuery = q.Encode()
	subject, textBody, htmlBody := buildActionEmail(purpose, ttl, u.String(), s.cfg.PublicBaseURL)
	if err := s.mailer.Send(user.Email, subject, textBody, htmlBody); err != nil {
		_ = s.repo.DeleteActionToken(ctx, action.TokenID)
		return fmt.Errorf("%w: %v", ErrEmailDeliveryFailed, err)
	}
	return nil
}
func displayHostFromBaseURL(publicBaseURL string) string {
	value := strings.TrimSpace(publicBaseURL)
	if value == "" {
		return ""
	}
	if parsed, err := url.Parse(value); err == nil {
		host := strings.TrimSpace(parsed.Hostname())
		if host != "" {
			return host
		}
	}
	value = strings.TrimPrefix(value, "https://")
	value = strings.TrimPrefix(value, "http://")
	value = strings.TrimRight(value, "/")
	return strings.TrimSpace(value)
}
func buildEmailSubject(base string, projectHost string) string {
	if projectHost == "" {
		return base
	}
	return fmt.Sprintf("%s в %s", base, projectHost)
}
func withProjectHost(projectHost string, withHost string, withoutHost string, args ...any) string {
	if projectHost == "" {
		if len(args) == 0 {
			return withoutHost
		}
		return fmt.Sprintf(withoutHost, args...)
	}
	values := make([]any, 0, len(args)+1)
	values = append(values, projectHost)
	values = append(values, args...)
	return fmt.Sprintf(withHost, values...)
}
func buildActionEmail(purpose string, ttl time.Duration, actionURL string, publicBaseURL string) (string, string, string) {
	expiry := formatDurationRu(ttl)
	projectHost := displayHostFromBaseURL(publicBaseURL)
	switch strings.TrimSpace(purpose) {
	case "verify_email":
		subject := buildEmailSubject("Подтвердите email для входа", projectHost)
		textBody := fmt.Sprintf(
			"Здравствуйте.\n\n"+
				"%s\n\n"+
				"Чтобы подтвердить адрес и завершить регистрацию, откройте ссылку:\n%s\n\n"+
				"Ссылка действует %s.\n\n"+
				"%s\n",
			withProjectHost(
				projectHost,
				"Вы получили это письмо, потому что для адреса email был создан аккаунт на %s.",
				"Вы получили это письмо, потому что для этого адреса был создан аккаунт.",
			),
			actionURL,
			expiry,
			withProjectHost(
				projectHost,
				"Если вы не регистрировались на %s, просто проигнорируйте это письмо.",
				"Если вы не регистрировались, просто проигнорируйте это письмо.",
			),
		)
		htmlBody := actionEmailHTML(
			projectHost,
			"Подтвердите email",
			withProjectHost(
				projectHost,
				"Вы получили это письмо, потому что для этого адреса был создан аккаунт на %s. Подтвердите email, чтобы завершить регистрацию и войти.",
				"Вы получили это письмо, потому что для этого адреса был создан аккаунт. Подтвердите email, чтобы завершить регистрацию и войти.",
			),
			"Подтвердить email",
			actionURL,
			expiry,
			withProjectHost(
				projectHost,
				"Если вы не регистрировались на %s, просто проигнорируйте это письмо.",
				"Если вы не регистрировались, просто проигнорируйте это письмо.",
			),
		)
		return subject, textBody, htmlBody
	case "password_reset":
		subject := buildEmailSubject("Сброс пароля", projectHost)
		textBody := fmt.Sprintf(
			"Здравствуйте.\n\n"+
				"%s\n\n"+
				"Чтобы задать новый пароль, откройте ссылку:\n%s\n\n"+
				"Ссылка действует %s.\n\n"+
				"Если вы не запрашивали сброс пароля, просто проигнорируйте это письмо.\n",
			withProjectHost(
				projectHost,
				"Мы получили запрос на сброс пароля для вашего аккаунта %s.",
				"Мы получили запрос на сброс пароля для вашего аккаунта.",
			),
			actionURL,
			expiry,
		)
		htmlBody := actionEmailHTML(
			projectHost,
			"Сброс пароля",
			withProjectHost(
				projectHost,
				"Мы получили запрос на смену пароля для вашего аккаунта %s. Перейдите по ссылке, чтобы задать новый пароль.",
				"Мы получили запрос на смену пароля для вашего аккаунта. Перейдите по ссылке, чтобы задать новый пароль.",
			),
			"Сменить пароль",
			actionURL,
			expiry,
			"Если вы не запрашивали сброс пароля, просто проигнорируйте это письмо.",
		)
		return subject, textBody, htmlBody
	default:
		subject := buildEmailSubject("Требуется действие", projectHost)
		textBody := fmt.Sprintf(
			"Здравствуйте.\n\n"+
				"%s\n\n"+
				"Ссылка действует %s.\n",
			withProjectHost(
				projectHost,
				"Для вашего аккаунта %s требуется действие по ссылке:\n%s",
				"Для вашего аккаунта требуется действие по ссылке:\n%s",
				actionURL,
			),
			expiry,
		)
		htmlBody := actionEmailHTML(
			projectHost,
			"Требуется действие",
			withProjectHost(
				projectHost,
				"Для вашего аккаунта %s требуется подтверждение действия по ссылке ниже.",
				"Для вашего аккаунта требуется подтверждение действия по ссылке ниже.",
			),
			"Открыть ссылку",
			actionURL,
			expiry,
			"Если вы не запрашивали это действие, проигнорируйте письмо.",
		)
		return subject, textBody, htmlBody
	}
}
func formatDurationRu(ttl time.Duration) string {
	if ttl <= 0 {
		return "несколько минут"
	}
	hoursTotal := int(ttl.Hours())
	minutesTotal := int(ttl.Minutes())
	if hoursTotal > 0 && minutesTotal%60 == 0 {
		return fmt.Sprintf("%d %s", hoursTotal, pluralRu(hoursTotal, "час", "часа", "часов"))
	}
	if hoursTotal > 0 {
		minutes := minutesTotal % 60
		return fmt.Sprintf("%d %s %d %s",
			hoursTotal, pluralRu(hoursTotal, "час", "часа", "часов"),
			minutes, pluralRu(minutes, "минута", "минуты", "минут"))
	}
	return fmt.Sprintf("%d %s", minutesTotal, pluralRu(minutesTotal, "минута", "минуты", "минут"))
}
func pluralRu(n int, one string, few string, many string) string {
	n = n % 100
	if n >= 11 && n <= 14 {
		return many
	}
	last := n % 10
	if last == 1 {
		return one
	}
	if last >= 2 && last <= 4 {
		return few
	}
	return many
}
type mfaTicketState struct {
	UserID          string `json:"user_id"`
	FingerprintHash string `json:"fingerprint_hash"`
	IP              string `json:"ip"`
	UserAgent       string `json:"user_agent"`
	ExpiresAt       string `json:"expires_at"`
}
type twoFactorSetupCacheState struct {
	OTPAuthURL string `json:"otpauth_url"`
	QRDataURL  string `json:"qr_data_url"`
}
func (s *Service) issueMFATicket(ctx context.Context, userID uuid.UUID, device DeviceContext) (string, time.Time, error) {
	token, err := randomToken(32)
	if err != nil {
		return "", time.Time{}, err
	}
	expiresAt := s.now().Add(s.cfg.MFATicketTTL)
	state, _ := json.Marshal(mfaTicketState{
		UserID:          userID.String(),
		FingerprintHash: device.FingerprintHash,
		IP:              device.IP,
		UserAgent:       device.UserAgent,
		ExpiresAt:       expiresAt.UTC().Format(time.RFC3339Nano),
	})
	key := "auth:mfa:ticket:" + hashToken(token)
	storageTTL := s.cfg.MFATicketTTL + time.Minute
	if storageTTL < s.cfg.MFATicketTTL {
		storageTTL = s.cfg.MFATicketTTL
	}
	if err := s.redis.Set(ctx, key, state, storageTTL).Err(); err != nil {
		return "", time.Time{}, err
	}
	return token, expiresAt, nil
}
func (s *Service) getMFATicketState(ctx context.Context, token string) (*mfaTicketState, error) {
	key := "auth:mfa:ticket:" + hashToken(strings.TrimSpace(token))
	raw, err := s.redis.Get(ctx, key).Result()
	if err != nil {
		return nil, ErrInvalidToken
	}
	var state mfaTicketState
	if err := json.Unmarshal([]byte(raw), &state); err != nil {
		return nil, ErrInvalidToken
	}
	if state.ExpiresAt != "" {
		expiresAt, err := time.Parse(time.RFC3339Nano, state.ExpiresAt)
		if err != nil {
			return nil, ErrInvalidToken
		}
		if !expiresAt.After(s.now()) {
			_ = s.clearMFATicket(ctx, token)
			return nil, ErrMFATicketExpired
		}
	}
	return &state, nil
}
func (s *Service) clearMFATicket(ctx context.Context, token string) error {
	key := "auth:mfa:ticket:" + hashToken(strings.TrimSpace(token))
	return s.redis.Del(ctx, key).Err()
}
func twoFactorSetupCacheKey(userID uuid.UUID) string {
	return "auth:2fa:setup:" + userID.String()
}
func (s *Service) getTwoFactorSetupCache(ctx context.Context, userID uuid.UUID) (*twoFactorSetupCacheState, error) {
	if s.redis == nil {
		return nil, nil
	}
	raw, err := s.redis.Get(ctx, twoFactorSetupCacheKey(userID)).Result()
	if err != nil {
		return nil, err
	}
	var state twoFactorSetupCacheState
	if err := json.Unmarshal([]byte(raw), &state); err != nil {
		return nil, err
	}
	return &state, nil
}
func (s *Service) storeTwoFactorSetupCache(ctx context.Context, userID uuid.UUID, result *TwoFactorSetupResult) error {
	if s.redis == nil || result == nil {
		return nil
	}
	payload, err := json.Marshal(twoFactorSetupCacheState{
		OTPAuthURL: result.OTPAuthURL,
		QRDataURL:  result.QRDataURL,
	})
	if err != nil {
		return err
	}
	return s.redis.Set(ctx, twoFactorSetupCacheKey(userID), payload, s.cfg.TwoFactorSetupTTL).Err()
}
func (s *Service) clearTwoFactorSetupCache(ctx context.Context, userID uuid.UUID) error {
	if s.redis == nil {
		return nil
	}
	return s.redis.Del(ctx, twoFactorSetupCacheKey(userID)).Err()
}
func (s *Service) blacklistAccessToken(ctx context.Context, claims *Claims) error {
	if claims == nil || claims.ExpiresAt == nil || claims.ID == "" {
		return nil
	}
	ttl := time.Until(claims.ExpiresAt.Time)
	if ttl <= 0 {
		return nil
	}
	return s.redis.Set(ctx, "auth:deny:access:"+claims.ID, "1", ttl).Err()
}
func (s *Service) rateLimit(ctx context.Context, bucket string, identifier string, limit int) error {
	if limit <= 0 || s.redis == nil {
		return nil
	}
	id := strings.TrimSpace(identifier)
	if id == "" {
		id = "unknown"
	}
	window := s.now().UTC().Truncate(s.cfg.RateWindow).Format(time.RFC3339)
	key := fmt.Sprintf("auth:rate:%s:%s:%s", bucket, hashToken(id), window)
	count, err := s.redis.Incr(ctx, key).Result()
	if err != nil {
		return nil
	}
	if count == 1 {
		_ = s.redis.Expire(ctx, key, s.cfg.RateWindow).Err()
	}
	if count > int64(limit) {
		return ErrTooManyRequests
	}
	return nil
}
func (s *Service) logLoginAttempt(ctx context.Context, userID *uuid.UUID, email string, outcome string, reason string, device DeviceContext, riskScore int, suspiciousReason string) error {
	return s.repo.CreateLoginAttempt(ctx, &models.LoginAttempt{
		AttemptID:        uuid.New(),
		UserID:           userID,
		EmailNormalized:  email,
		IP:               device.IP,
		UserAgent:        device.UserAgent,
		Outcome:          outcome,
		FailureReason:    reason,
		RiskScore:        riskScore,
		FingerprintHash:  device.FingerprintHash,
		SuspiciousReason: suspiciousReason,
	})
}
func (s *Service) logSecurityEvent(ctx context.Context, userID *uuid.UUID, actorUserID *uuid.UUID, category string, eventType string, severity string, device DeviceContext, payload map[string]any) error {
	if payload == nil {
		payload = map[string]any{}
	}
	raw, _ := json.Marshal(payload)
	return s.repo.CreateSecurityEvent(ctx, &models.SecurityEvent{
		EventID:     uuid.New(),
		UserID:      userID,
		ActorUserID: actorUserID,
		Category:    category,
		EventType:   eventType,
		Severity:    severity,
		IP:          device.IP,
		UserAgent:   device.UserAgent,
		Payload:     raw,
	})
}
func normalizeEmail(raw string) (string, error) {
	addr, err := mail.ParseAddress(strings.TrimSpace(raw))
	if err != nil {
		return "", ErrInvalidCredentials
	}
	return strings.ToLower(strings.TrimSpace(addr.Address)), nil
}
func BuildDeviceContext(ip string, userAgent string, acceptLanguage string) DeviceContext {
	normalizedIP := strings.TrimSpace(ip)
	parsed := net.ParseIP(normalizedIP)
	if parsed != nil {
		if v4 := parsed.To4(); v4 != nil {
			normalizedIP = v4.String()
		} else {
			normalizedIP = parsed.String()
		}
	}
	normalizedUserAgent := strings.TrimSpace(strings.ToLower(userAgent))
	normalizedAcceptLanguage := normalizeAcceptLanguage(acceptLanguage)
	fingerprintSeed := strings.Join([]string{
		normalizedUserAgent,
		normalizedAcceptLanguage,
	}, "|")
	return DeviceContext{
		IP:              normalizedIP,
		UserAgent:       strings.TrimSpace(userAgent),
		AcceptLanguage:  strings.TrimSpace(acceptLanguage),
		FingerprintHash: hashToken(fingerprintSeed),
		DeviceLabel:     strings.TrimSpace(userAgent),
	}
}
func normalizeAcceptLanguage(raw string) string {
	if strings.TrimSpace(raw) == "" {
		return ""
	}
	parts := strings.Split(strings.ToLower(raw), ",")
	normalized := make([]string, 0, len(parts))
	for _, part := range parts {
		value := strings.TrimSpace(part)
		if value == "" {
			continue
		}
		if idx := strings.Index(value, ";"); idx >= 0 {
			value = strings.TrimSpace(value[:idx])
		}
		if value == "" {
			continue
		}
		normalized = append(normalized, value)
	}
	return strings.Join(normalized, ",")
}
func fallbackString(value string, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return strings.TrimSpace(value)
}
func hasRole(roles []string, role string) bool {
	target := strings.ToLower(strings.TrimSpace(role))
	if target == "" {
		return false
	}
	for _, item := range roles {
		if strings.ToLower(strings.TrimSpace(item)) == target {
			return true
		}
	}
	return false
}
func ptr[T any](v T) *T {
	return &v
}
func isUniqueViolation(err error) bool {
	return strings.Contains(strings.ToLower(err.Error()), "duplicate key")
}
