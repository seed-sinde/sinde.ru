package models

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type User struct {
	UserID               uuid.UUID       `json:"user_id"`
	Email                string          `json:"email"`
	EmailNormalized      string          `json:"-"`
	PasswordHash         string          `json:"-"`
	Status               string          `json:"status"`
	EmailVerifiedAt      *time.Time      `json:"email_verified_at,omitempty"`
	DisplayName          string          `json:"display_name"`
	Locale               string          `json:"locale"`
	Timezone             string          `json:"timezone"`
	Roles                []string        `json:"roles"`
	IsTwoFactorEnabled   bool            `json:"is_two_factor_enabled"`
	TOTPSecretCiphertext []byte          `json:"-"`
	Profile              json.RawMessage `json:"profile"`
	Settings             json.RawMessage `json:"settings"`
	PrimaryTraitUUID     *uuid.UUID      `json:"primary_trait_uuid,omitempty"`
	SessionVersion       int64           `json:"session_version"`
	LastLoginAt          *time.Time      `json:"last_login_at,omitempty"`
	BlockedReason        string          `json:"blocked_reason,omitempty"`
	BlockedAt            *time.Time      `json:"blocked_at,omitempty"`
	CreatedAt            time.Time       `json:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at"`
}
type AuthSession struct {
	SessionID         uuid.UUID  `json:"session_id"`
	UserID            uuid.UUID  `json:"user_id"`
	FamilyID          uuid.UUID  `json:"family_id"`
	RefreshJTI        string     `json:"refresh_jti"`
	RefreshTokenHash  string     `json:"-"`
	CSRFTokenHash     string     `json:"-"`
	FingerprintHash   string     `json:"fingerprint_hash"`
	DeviceLabel       string     `json:"device_label"`
	IP                string     `json:"ip"`
	UserAgent         string     `json:"user_agent"`
	MFAVerified       bool       `json:"mfa_verified"`
	ReplacedBySession *uuid.UUID `json:"replaced_by_session_id,omitempty"`
	LastSeenAt        time.Time  `json:"last_seen_at"`
	ExpiresAt         time.Time  `json:"expires_at"`
	RevokedAt         *time.Time `json:"revoked_at,omitempty"`
	RevokeReason      string     `json:"revoke_reason,omitempty"`
	CompromisedAt     *time.Time `json:"compromised_at,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}
type ActionToken struct {
	TokenID       uuid.UUID       `json:"token_id"`
	UserID        uuid.UUID       `json:"user_id"`
	Purpose       string          `json:"purpose"`
	TokenHash     string          `json:"-"`
	ExpiresAt     time.Time       `json:"expires_at"`
	ConsumedAt    *time.Time      `json:"consumed_at,omitempty"`
	DeliveryValue string          `json:"delivery_value"`
	Meta          json.RawMessage `json:"meta"`
	CreatedAt     time.Time       `json:"created_at"`
}
type BackupCode struct {
	BackupCodeID uuid.UUID  `json:"backup_code_id"`
	UserID       uuid.UUID  `json:"user_id"`
	CodeHash     string     `json:"-"`
	ConsumedAt   *time.Time `json:"consumed_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}
type SavedTraitSet struct {
	SavedSetID  uuid.UUID `json:"saved_set_id"`
	UserID      uuid.UUID `json:"user_id"`
	TraitUUID   uuid.UUID `json:"trait_uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
type LoginAttempt struct {
	AttemptID        uuid.UUID  `json:"attempt_id"`
	UserID           *uuid.UUID `json:"user_id,omitempty"`
	EmailNormalized  string     `json:"email_normalized"`
	IP               string     `json:"ip"`
	UserAgent        string     `json:"user_agent"`
	Outcome          string     `json:"outcome"`
	FailureReason    string     `json:"failure_reason,omitempty"`
	RiskScore        int        `json:"risk_score"`
	FingerprintHash  string     `json:"fingerprint_hash,omitempty"`
	SuspiciousReason string     `json:"suspicious_reason,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
}
type SecurityEvent struct {
	EventID     uuid.UUID       `json:"event_id"`
	UserID      *uuid.UUID      `json:"user_id,omitempty"`
	ActorUserID *uuid.UUID      `json:"actor_user_id,omitempty"`
	SessionID   *uuid.UUID      `json:"session_id,omitempty"`
	Category    string          `json:"category"`
	EventType   string          `json:"event_type"`
	Severity    string          `json:"severity"`
	IP          string          `json:"ip"`
	UserAgent   string          `json:"user_agent"`
	Payload     json.RawMessage `json:"payload"`
	CreatedAt   time.Time       `json:"created_at"`
}
