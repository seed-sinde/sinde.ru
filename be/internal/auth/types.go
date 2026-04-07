package auth

import (
	"encoding/json"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"sinde.ru/internal/models"
	"time"
)

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

type Config struct {
	PublicBaseURL           string
	JWTIssuer               string
	JWTAudience             string
	JWTSecret               []byte
	EncryptionKey           []byte
	AccessTTL               time.Duration
	RefreshTTL              time.Duration
	VerifyTTL               time.Duration
	ResetTTL                time.Duration
	MFATicketTTL            time.Duration
	TwoFactorSetupTTL       time.Duration
	AccessCookieName        string
	RefreshCookieName       string
	CSRFCookieName          string
	SessionHintCookieName   string
	CookieDomain            string
	CookieSecure            bool
	CookieSameSite          string
	TOTPIssuer              string
	PasswordMinLength       int
	AllowedOrigins          []string
	RedisAddr               string
	RedisPassword           string
	RedisDB                 int
	RegisterRateLimit       int
	LoginRateLimitIP        int
	LoginRateLimitMail      int
	TwoFactorSetupRateLimit int
	MFALoginRateLimitIP     int
	MFALoginRateLimitTicket int
	VerifyRateLimit         int
	ResetRateLimit          int
	RateWindow              time.Duration
	MailFrom                string
	MailerDriver            string
	SMTPHost                string
	SMTPPort                int
	SMTPUsername            string
	SMTPPassword            string
	SMTPTLSMode             string
	SMTPHELO                string
}
type Claims struct {
	TokenType      string   `json:"typ"`
	UserID         string   `json:"uid"`
	SessionID      string   `json:"sid,omitempty"`
	Email          string   `json:"email,omitempty"`
	Roles          []string `json:"roles,omitempty"`
	MFAVerified    bool     `json:"mfa"`
	Fingerprint    string   `json:"fp,omitempty"`
	SessionVersion int64    `json:"sv"`
	jwt.RegisteredClaims
}
type DeviceContext struct {
	IP              string
	UserAgent       string
	AcceptLanguage  string
	FingerprintHash string
	DeviceLabel     string
}
type TokenBundle struct {
	AccessToken      string
	RefreshToken     string
	RefreshJTI       string
	CSRFToken        string
	AccessExpiresAt  time.Time
	RefreshExpiresAt time.Time
	Claims           *Claims
}
type RegisterInput struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	DisplayName  string `json:"display_name"`
	Locale       string `json:"locale"`
	Timezone     string `json:"timezone"`
	CaptchaToken string `json:"captcha_token"`
	Honeypot     string `json:"website"`
}
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type VerifyEmailInput struct {
	Token string `json:"token"`
}
type RequestEmailVerificationInput struct {
	Email string `json:"email"`
}
type RequestEmailChangeInput struct {
	Email string `json:"email"`
}
type RefreshInput struct {
	RefreshToken string
	CSRFToken    string
}
type LogoutInput struct {
	AccessToken string
}
type ForgotPasswordInput struct {
	Email string `json:"email"`
}
type ResetPasswordInput struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}
type ChangePasswordInput struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}
type UpdateProfileInput struct {
	DisplayName *string          `json:"display_name"`
	Locale      *string          `json:"locale"`
	Timezone    *string          `json:"timezone"`
	Profile     *json.RawMessage `json:"profile"`
	Settings    *json.RawMessage `json:"settings"`
}
type AuthUserPatch struct {
	DisplayName *string `json:"display_name,omitempty"`
	Locale      *string `json:"locale,omitempty"`
	Timezone    *string `json:"timezone,omitempty"`
	Profile     any     `json:"profile,omitempty"`
	Settings    any     `json:"settings,omitempty"`
}
type SetPrimaryTraitInput struct {
	SetUUID   string `json:"set_uuid"`
	TraitUUID string `json:"trait_uuid"`
}
type SaveTraitSetInput struct {
	SetUUID     string `json:"set_uuid"`
	TraitUUID   string `json:"trait_uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
type UpdateSavedTraitSetInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
type TwoFactorLoginInput struct {
	Ticket string `json:"ticket"`
	Code   string `json:"code"`
}
type TwoFactorEnableInput struct {
	Code string `json:"code"`
}
type TwoFactorDisableInput struct {
	Password string `json:"password"`
	Code     string `json:"code"`
}
type BlockUserInput struct {
	Reason string `json:"reason"`
}
type AdminSetUserRoleInput struct {
	Role string `json:"role"`
}
type AuthUser struct {
	UserID             uuid.UUID   `json:"user_id"`
	Email              string      `json:"email"`
	Status             string      `json:"status"`
	EmailVerifiedAt    *time.Time  `json:"email_verified_at,omitempty"`
	DisplayName        string      `json:"display_name"`
	Locale             string      `json:"locale"`
	Timezone           string      `json:"timezone"`
	Roles              []string    `json:"roles"`
	IsTwoFactorEnabled bool        `json:"is_two_factor_enabled"`
	PrimaryTraitUUID   *uuid.UUID  `json:"primary_trait_uuid,omitempty"`
	LastLoginAt        *time.Time  `json:"last_login_at,omitempty"`
	BlockedReason      string      `json:"blocked_reason,omitempty"`
	BlockedAt          *time.Time  `json:"blocked_at,omitempty"`
	CreatedAt          time.Time   `json:"created_at"`
	Profile            interface{} `json:"profile"`
	Settings           interface{} `json:"settings"`
}
type SessionView struct {
	SessionID    uuid.UUID  `json:"session_id"`
	DeviceLabel  string     `json:"device_label"`
	IP           string     `json:"ip"`
	UserAgent    string     `json:"user_agent"`
	MFAVerified  bool       `json:"mfa_verified"`
	LastSeenAt   time.Time  `json:"last_seen_at"`
	ExpiresAt    time.Time  `json:"expires_at"`
	RevokedAt    *time.Time `json:"revoked_at,omitempty"`
	RevokeReason string     `json:"revoke_reason,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	IsCurrent    bool       `json:"is_current"`
}
type RegisterResult struct {
	User            AuthUser `json:"user"`
	VerificationTTL string   `json:"verification_ttl"`
}
type LoginResult struct {
	User         *AuthUser    `json:"user,omitempty"`
	Session      *SessionView `json:"session,omitempty"`
	CSRFToken    string       `json:"csrf_token,omitempty"`
	MFARequired  bool         `json:"mfa_required"`
	MFATicket    string       `json:"mfa_ticket,omitempty"`
	MFAExpiresAt *time.Time   `json:"mfa_expires_at,omitempty"`
	MFAMethods   []string     `json:"mfa_methods,omitempty"`
}
type TwoFactorSetupResult struct {
	Secret      string   `json:"secret"`
	OTPAuthURL  string   `json:"otpauth_url"`
	QRDataURL   string   `json:"qr_data_url"`
	BackupCodes []string `json:"backup_codes"`
}
type MeResult struct {
	User AuthUser `json:"user"`
}
type VerifyEmailResult struct {
	Verified     bool   `json:"verified"`
	Action       string `json:"action"`
	Email        string `json:"email,omitempty"`
	SessionHints bool   `json:"session_hints,omitempty"`
}
type UpdateProfileResult struct {
	Changes AuthUserPatch `json:"changes"`
}
type SavedTraitSetView struct {
	SavedSetID  uuid.UUID `json:"saved_set_id"`
	SetUUID     uuid.UUID `json:"set_uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
type SavedTraitSetListResult struct {
	PrimaryTraitUUID *uuid.UUID          `json:"primary_trait_uuid,omitempty"`
	Items            []SavedTraitSetView `json:"items"`
}
type SessionListResult struct {
	Items []SessionView `json:"items"`
}
type LoginAttemptView struct {
	AttemptID        uuid.UUID `json:"attempt_id"`
	IP               string    `json:"ip"`
	UserAgent        string    `json:"user_agent"`
	Outcome          string    `json:"outcome"`
	FailureReason    string    `json:"failure_reason,omitempty"`
	RiskScore        int       `json:"risk_score"`
	SuspiciousReason string    `json:"suspicious_reason,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}
type SecurityEventView struct {
	EventID   uuid.UUID   `json:"event_id"`
	Category  string      `json:"category"`
	EventType string      `json:"event_type"`
	Severity  string      `json:"severity"`
	IP        string      `json:"ip"`
	UserAgent string      `json:"user_agent"`
	SessionID *uuid.UUID  `json:"session_id,omitempty"`
	Payload   interface{} `json:"payload"`
	CreatedAt time.Time   `json:"created_at"`
}
type LoginAttemptListResult struct {
	Items []LoginAttemptView `json:"items"`
}
type SecurityEventListResult struct {
	Items []SecurityEventView `json:"items"`
}
type AdminUserView struct {
	UserID             uuid.UUID  `json:"user_id"`
	Email              string     `json:"email"`
	Status             string     `json:"status"`
	DisplayName        string     `json:"display_name"`
	Locale             string     `json:"locale"`
	Timezone           string     `json:"timezone"`
	Roles              []string   `json:"roles"`
	IsTwoFactorEnabled bool       `json:"is_two_factor_enabled"`
	LastLoginAt        *time.Time `json:"last_login_at,omitempty"`
	BlockedReason      string     `json:"blocked_reason,omitempty"`
	BlockedAt          *time.Time `json:"blocked_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
}
type AdminUserListResult struct {
	Items  []AdminUserView `json:"items"`
	Total  int             `json:"total"`
	Limit  int             `json:"limit"`
	Offset int             `json:"offset"`
}
type UserSummaryResult struct {
	NewApprovedRecipesSinceLastLogin int                   `json:"new_approved_recipes_since_last_login"`
	NewRejectedRecipesSinceLastLogin int                   `json:"new_rejected_recipes_since_last_login"`
	LastLoginBeforeCurrentSessionAt  *time.Time            `json:"last_login_before_current_session_at,omitempty"`
	NotificationsReadAt              *time.Time            `json:"notifications_read_at,omitempty"`
	NotificationsSinceAt             *time.Time            `json:"notifications_since_at,omitempty"`
	Admin                            *UserSummaryAdminInfo `json:"admin,omitempty"`
	HasUnread                        bool                  `json:"has_unread"`
	CheckedAt                        time.Time             `json:"checked_at"`
}
type UserSummaryAdminInfo struct {
	UsersTotal                      int            `json:"users_total"`
	PendingRecipesTotal             int            `json:"pending_recipes_total"`
	RecipeStatusTotals              map[string]int `json:"recipe_status_totals"`
	NewUsersSinceLastLogin          int            `json:"new_users_since_last_login"`
	NewPendingRecipesSinceLastLogin int            `json:"new_pending_recipes_since_last_login"`
}
type AdminSummaryResult struct {
	UsersTotal                      int            `json:"users_total"`
	PendingRecipesTotal             int            `json:"pending_recipes_total"`
	RecipeStatusTotals              map[string]int `json:"recipe_status_totals"`
	NewUsersSinceLastLogin          int            `json:"new_users_since_last_login"`
	NewPendingRecipesSinceLastLogin int            `json:"new_pending_recipes_since_last_login"`
	LastLoginBeforeCurrentSessionAt *time.Time     `json:"last_login_before_current_session_at,omitempty"`
	NotificationsReadAt             *time.Time     `json:"notifications_read_at,omitempty"`
	NotificationsSinceAt            *time.Time     `json:"notifications_since_at,omitempty"`
	HasUnread                       bool           `json:"has_unread"`
	CheckedAt                       time.Time      `json:"checked_at"`
}
type AdminTraitKeySearchItem struct {
	KeyID       int64       `json:"key_id"`
	Syn         string      `json:"syn"`
	Meta        interface{} `json:"meta"`
	TraitCount  int64       `json:"trait_count"`
	CreatedFrom string      `json:"created_from,omitempty"`
}
type AdminTraitKeySearchResult struct {
	Items []AdminTraitKeySearchItem `json:"items"`
	Query string                    `json:"query"`
	Limit int                       `json:"limit"`
}
type AdminTraitsSetsAnalysisTopKey struct {
	KeyID      int64  `json:"key_id"`
	Syn        string `json:"syn"`
	TraitCount int64  `json:"trait_count"`
}
type AdminTraitsSetsAnalysis struct {
	TotalTraits             int64                           `json:"total_traits"`
	UniqueTraitPairs        int64                           `json:"unique_trait_pairs"`
	UniqueTraitKeys         int64                           `json:"unique_trait_keys"`
	TraitsReferencedInSets  int64                           `json:"traits_referenced_in_sets"`
	OrphanTraits            int64                           `json:"orphan_traits"`
	TotalSets               int64                           `json:"total_sets"`
	UniqueSetsByChildren    int64                           `json:"unique_sets_by_children"`
	DerivedSets             int64                           `json:"derived_sets"`
	DerivedSetRate          float64                         `json:"derived_set_rate"`
	SetUniquenessRate       float64                         `json:"set_uniqueness_rate"`
	TraitCoverageInSetsRate float64                         `json:"trait_coverage_in_sets_rate"`
	TopKeys                 []AdminTraitsSetsAnalysisTopKey `json:"top_keys"`
}

func ToAuthUser(user *models.User) AuthUser {
	return AuthUser{
		UserID:             user.UserID,
		Email:              user.Email,
		Status:             user.Status,
		EmailVerifiedAt:    user.EmailVerifiedAt,
		DisplayName:        user.DisplayName,
		Locale:             user.Locale,
		Timezone:           user.Timezone,
		Roles:              append([]string(nil), user.Roles...),
		IsTwoFactorEnabled: user.IsTwoFactorEnabled,
		PrimaryTraitUUID:   user.PrimaryTraitUUID,
		LastLoginAt:        user.LastLoginAt,
		BlockedReason:      user.BlockedReason,
		BlockedAt:          user.BlockedAt,
		CreatedAt:          user.CreatedAt,
		Profile:            jsonMap(user.Profile),
		Settings:           jsonMap(user.Settings),
	}
}
func ToSavedTraitSetView(item models.SavedTraitSet) SavedTraitSetView {
	return SavedTraitSetView{
		SavedSetID:  item.SavedSetID,
		SetUUID:     item.TraitUUID,
		Name:        item.Name,
		Description: item.Description,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}
}
func ToLoginAttemptView(item models.LoginAttempt) LoginAttemptView {
	return LoginAttemptView{
		AttemptID:        item.AttemptID,
		IP:               item.IP,
		UserAgent:        item.UserAgent,
		Outcome:          item.Outcome,
		FailureReason:    item.FailureReason,
		RiskScore:        item.RiskScore,
		SuspiciousReason: item.SuspiciousReason,
		CreatedAt:        item.CreatedAt,
	}
}
func ToSecurityEventView(item models.SecurityEvent) SecurityEventView {
	return SecurityEventView{
		EventID:   item.EventID,
		Category:  item.Category,
		EventType: item.EventType,
		Severity:  item.Severity,
		IP:        item.IP,
		UserAgent: item.UserAgent,
		SessionID: item.SessionID,
		Payload:   jsonMap(item.Payload),
		CreatedAt: item.CreatedAt,
	}
}
func jsonMap(raw []byte) any {
	if len(raw) == 0 {
		return map[string]any{}
	}
	var out any
	if err := json.Unmarshal(raw, &out); err != nil {
		return map[string]any{}
	}
	return out
}
