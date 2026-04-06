package auth

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"sinde.ru/internal/models"
	"strings"
	"time"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}
func (r *Repository) CreateUser(ctx context.Context, user *models.User) error {
	const query = `
		INSERT INTO users (
			user_id, email, email_normalized, password_hash, status,
			display_name, locale, timezone, roles, is_two_factor_enabled,
			profile, settings, session_version
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10,
			$11, $12, $13
		)
	`
	_, err := r.db.Exec(ctx, query,
		user.UserID, user.Email, user.EmailNormalized, user.PasswordHash, user.Status,
		user.DisplayName, user.Locale, user.Timezone, user.Roles, user.IsTwoFactorEnabled,
		jsonOrDefault(user.Profile), jsonOrDefault(user.Settings), user.SessionVersion,
	)
	return err
}
func (r *Repository) GetUserByEmail(ctx context.Context, emailNormalized string) (*models.User, error) {
	const query = `
		SELECT
			user_id, email, email_normalized, password_hash, status,
			email_verified_at, display_name, locale, timezone, roles,
			is_two_factor_enabled, totp_secret_ciphertext, profile, settings, primary_trait_uuid,
			session_version, last_login_at, blocked_reason, blocked_at,
			created_at, updated_at
		FROM users
		WHERE email_normalized = $1
	`
	return scanUser(r.db.QueryRow(ctx, query, emailNormalized))
}
func (r *Repository) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	const query = `
		SELECT
			user_id, email, email_normalized, password_hash, status,
			email_verified_at, display_name, locale, timezone, roles,
			is_two_factor_enabled, totp_secret_ciphertext, profile, settings, primary_trait_uuid,
			session_version, last_login_at, blocked_reason, blocked_at,
			created_at, updated_at
		FROM users
		WHERE user_id = $1
	`
	return scanUser(r.db.QueryRow(ctx, query, userID))
}
func (r *Repository) UpdateUserLastLogin(ctx context.Context, userID uuid.UUID, t time.Time) error {
	_, err := r.db.Exec(ctx, `UPDATE users SET last_login_at = $1 WHERE user_id = $2`, t, userID)
	return err
}
func (r *Repository) MarkUserEmailVerified(ctx context.Context, userID uuid.UUID, t time.Time) error {
	_, err := r.db.Exec(ctx, `
		UPDATE users
		SET status = 'active',
			email_verified_at = $1
		WHERE user_id = $2
	`, t, userID)
	return err
}
func (r *Repository) CreateActionToken(ctx context.Context, token *models.ActionToken) error {
	const query = `
		INSERT INTO auth_action_tokens (
			token_id, user_id, purpose, token_hash, expires_at,
			delivery_value, meta
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(ctx, query,
		token.TokenID, token.UserID, token.Purpose, token.TokenHash,
		token.ExpiresAt, token.DeliveryValue, jsonOrDefault(token.Meta),
	)
	return err
}
func (r *Repository) DeleteActionToken(ctx context.Context, tokenID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM auth_action_tokens WHERE token_id = $1`, tokenID)
	return err
}
func (r *Repository) ConsumeActionToken(ctx context.Context, purpose string, tokenHash string, now time.Time) (*models.ActionToken, error) {
	const query = `
		UPDATE auth_action_tokens
		SET consumed_at = $3
		WHERE purpose = $1 AND token_hash = $2 AND consumed_at IS NULL AND expires_at > $3
		RETURNING token_id, user_id, purpose, token_hash, expires_at, consumed_at, delivery_value, meta, created_at
	`
	return scanActionToken(r.db.QueryRow(ctx, query, purpose, tokenHash, now))
}
func (r *Repository) UpdateTOTPSecret(ctx context.Context, userID uuid.UUID, ciphertext []byte) error {
	_, err := r.db.Exec(ctx, `UPDATE users SET totp_secret_ciphertext = $1 WHERE user_id = $2`, ciphertext, userID)
	return err
}
func (r *Repository) EnableTwoFactor(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `UPDATE users SET is_two_factor_enabled = TRUE WHERE user_id = $1`, userID)
	return err
}
func (r *Repository) DisableTwoFactor(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
		UPDATE users
		SET is_two_factor_enabled = FALSE,
			totp_secret_ciphertext = NULL
		WHERE user_id = $1
	`, userID)
	return err
}
func (r *Repository) ReplaceBackupCodes(ctx context.Context, userID uuid.UUID, hashes []string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	if _, err := tx.Exec(ctx, `DELETE FROM auth_backup_codes WHERE user_id = $1`, userID); err != nil {
		return err
	}
	for _, hash := range hashes {
		if _, err := tx.Exec(ctx, `
			INSERT INTO auth_backup_codes (backup_code_id, user_id, code_hash)
			VALUES ($1, $2, $3)
		`, uuid.New(), userID, hash); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}
func (r *Repository) ListBackupCodes(ctx context.Context, userID uuid.UUID) ([]models.BackupCode, error) {
	rows, err := r.db.Query(ctx, `
		SELECT backup_code_id, user_id, code_hash, consumed_at, created_at
		FROM auth_backup_codes
		WHERE user_id = $1 AND consumed_at IS NULL
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]models.BackupCode, 0)
	for rows.Next() {
		var item models.BackupCode
		if err := rows.Scan(&item.BackupCodeID, &item.UserID, &item.CodeHash, &item.ConsumedAt, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
func (r *Repository) ConsumeBackupCode(ctx context.Context, backupCodeID uuid.UUID, now time.Time) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE auth_backup_codes
		SET consumed_at = $2
		WHERE backup_code_id = $1 AND consumed_at IS NULL
	`, backupCodeID, now)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrInvalidMFACode
	}
	return nil
}
func (r *Repository) CreateSession(ctx context.Context, session *models.AuthSession) error {
	const query = `
		INSERT INTO auth_sessions (
			session_id, user_id, family_id, refresh_jti, refresh_token_hash,
			csrf_token_hash, fingerprint_hash, device_label, ip, user_agent,
			mfa_verified, last_seen_at, expires_at, revoke_reason
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10,
			$11, $12, $13, $14
		)
	`
	_, err := r.db.Exec(ctx, query,
		session.SessionID, session.UserID, session.FamilyID, session.RefreshJTI, session.RefreshTokenHash,
		session.CSRFTokenHash, session.FingerprintHash, session.DeviceLabel, nullableINET(session.IP), session.UserAgent,
		session.MFAVerified, session.LastSeenAt, session.ExpiresAt, session.RevokeReason,
	)
	return err
}
func (r *Repository) GetSessionByID(ctx context.Context, sessionID uuid.UUID) (*models.AuthSession, error) {
	const query = `
		SELECT
			session_id, user_id, family_id, refresh_jti, refresh_token_hash,
			csrf_token_hash, fingerprint_hash, device_label, COALESCE(host(ip), ''), user_agent,
			mfa_verified, replaced_by_session_id, last_seen_at, expires_at, revoked_at,
			revoke_reason, compromised_at, created_at, updated_at
		FROM auth_sessions
		WHERE session_id = $1
	`
	return scanSession(r.db.QueryRow(ctx, query, sessionID))
}
func (r *Repository) RotateSession(ctx context.Context, currentSessionID uuid.UUID, next *models.AuthSession, revokeAt time.Time) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	if _, err := tx.Exec(ctx, `
		INSERT INTO auth_sessions (
			session_id, user_id, family_id, refresh_jti, refresh_token_hash,
			csrf_token_hash, fingerprint_hash, device_label, ip, user_agent,
			mfa_verified, last_seen_at, expires_at, revoke_reason
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10,
			$11, $12, $13, $14
		)
	`, next.SessionID, next.UserID, next.FamilyID, next.RefreshJTI, next.RefreshTokenHash,
		next.CSRFTokenHash, next.FingerprintHash, next.DeviceLabel, nullableINET(next.IP), next.UserAgent,
		next.MFAVerified, next.LastSeenAt, next.ExpiresAt, next.RevokeReason,
	); err != nil {
		return err
	}
	tag, err := tx.Exec(ctx, `
		UPDATE auth_sessions
		SET revoked_at = $2,
			revoke_reason = 'rotated',
			replaced_by_session_id = $3,
			last_seen_at = $2
		WHERE session_id = $1 AND revoked_at IS NULL
	`, currentSessionID, revokeAt, next.SessionID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrSessionNotFound
	}
	return tx.Commit(ctx)
}
func (r *Repository) TouchSession(ctx context.Context, sessionID uuid.UUID, t time.Time, ip string, userAgent string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE auth_sessions
		SET last_seen_at = $2,
			ip = $3,
			user_agent = $4
		WHERE session_id = $1
	`, sessionID, t, nullableINET(ip), userAgent)
	return err
}
func (r *Repository) RevokeSession(ctx context.Context, sessionID uuid.UUID, reason string, now time.Time) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE auth_sessions
		SET revoked_at = COALESCE(revoked_at, $3),
			revoke_reason = CASE WHEN revoke_reason = '' THEN $2 ELSE revoke_reason END
		WHERE session_id = $1
	`, sessionID, reason, now)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrSessionNotFound
	}
	return nil
}
func (r *Repository) RevokeSessionFamily(ctx context.Context, familyID uuid.UUID, reason string, now time.Time) error {
	_, err := r.db.Exec(ctx, `
		UPDATE auth_sessions
		SET revoked_at = COALESCE(revoked_at, $3),
			revoke_reason = CASE WHEN revoke_reason = '' THEN $2 ELSE revoke_reason END,
			compromised_at = $3
		WHERE family_id = $1
	`, familyID, reason, now)
	return err
}
func (r *Repository) RevokeAllSessionsForUser(ctx context.Context, userID uuid.UUID, reason string, now time.Time) error {
	_, err := r.db.Exec(ctx, `
		UPDATE auth_sessions
		SET revoked_at = COALESCE(revoked_at, $3),
			revoke_reason = CASE WHEN revoke_reason = '' THEN $2 ELSE revoke_reason END
		WHERE user_id = $1
	`, userID, reason, now)
	return err
}
func (r *Repository) ListSessionsByUser(ctx context.Context, userID uuid.UUID) ([]models.AuthSession, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			session_id, user_id, family_id, refresh_jti, refresh_token_hash,
			csrf_token_hash, fingerprint_hash, device_label, COALESCE(host(ip), ''), user_agent,
			mfa_verified, replaced_by_session_id, last_seen_at, expires_at, revoked_at,
			revoke_reason, compromised_at, created_at, updated_at
		FROM auth_sessions
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]models.AuthSession, 0)
	for rows.Next() {
		item, err := scanSession(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}
	return items, rows.Err()
}
func (r *Repository) HasSeenFingerprint(ctx context.Context, userID uuid.UUID, fingerprintHash string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM auth_sessions
			WHERE user_id = $1 AND fingerprint_hash = $2
		)
	`, userID, fingerprintHash).Scan(&exists)
	return exists, err
}
func (r *Repository) UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error {
	_, err := r.db.Exec(ctx, `UPDATE users SET password_hash = $1 WHERE user_id = $2`, passwordHash, userID)
	return err
}
func (r *Repository) UpdateProfile(ctx context.Context, userID uuid.UUID, displayName string, locale string, timezone string, profile []byte, settings []byte) error {
	_, err := r.db.Exec(ctx, `
		UPDATE users
		SET display_name = $2,
			locale = $3,
			timezone = $4,
			profile = $5,
			settings = $6
		WHERE user_id = $1
	`, userID, displayName, locale, timezone, jsonOrDefault(profile), jsonOrDefault(settings))
	return err
}
func (r *Repository) UpdatePrimaryTraitUUID(ctx context.Context, userID uuid.UUID, traitUUID *uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
		UPDATE users
		SET primary_trait_uuid = $2
		WHERE user_id = $1
	`, userID, traitUUID)
	return err
}
func (r *Repository) UpsertSavedTraitSet(ctx context.Context, item *models.SavedTraitSet) (*models.SavedTraitSet, error) {
	const query = `
		INSERT INTO auth_saved_trait_sets (
			saved_set_id, user_id, trait_uuid, name, description
		) VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, trait_uuid) DO UPDATE
		SET name = EXCLUDED.name,
			description = EXCLUDED.description,
			updated_at = now()
		RETURNING saved_set_id, user_id, trait_uuid, name, description, created_at, updated_at
	`
	return scanSavedTraitSet(r.db.QueryRow(ctx, query,
		item.SavedSetID, item.UserID, item.TraitUUID, item.Name, item.Description,
	))
}
func (r *Repository) ListSavedTraitSets(ctx context.Context, userID uuid.UUID) ([]models.SavedTraitSet, error) {
	rows, err := r.db.Query(ctx, `
		SELECT saved_set_id, user_id, trait_uuid, name, description, created_at, updated_at
		FROM auth_saved_trait_sets
		WHERE user_id = $1
		ORDER BY updated_at DESC, created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]models.SavedTraitSet, 0)
	for rows.Next() {
		item, scanErr := scanSavedTraitSet(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		items = append(items, *item)
	}
	return items, rows.Err()
}
func (r *Repository) UpdateSavedTraitSet(ctx context.Context, userID uuid.UUID, savedSetID uuid.UUID, name string, description string) (*models.SavedTraitSet, error) {
	const query = `
		UPDATE auth_saved_trait_sets
		SET name = $3,
			description = $4,
			updated_at = now()
		WHERE user_id = $1 AND saved_set_id = $2
		RETURNING saved_set_id, user_id, trait_uuid, name, description, created_at, updated_at
	`
	return scanSavedTraitSet(r.db.QueryRow(ctx, query, userID, savedSetID, name, description))
}
func (r *Repository) DeleteSavedTraitSet(ctx context.Context, userID uuid.UUID, savedSetID uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `
		DELETE FROM auth_saved_trait_sets
		WHERE user_id = $1 AND saved_set_id = $2
	`, userID, savedSetID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrSavedTraitSetNotFound
	}
	return nil
}
func (r *Repository) DeleteSavedTraitSetByTraitUUID(ctx context.Context, userID uuid.UUID, traitUUID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
		DELETE FROM auth_saved_trait_sets
		WHERE user_id = $1 AND trait_uuid = $2
	`, userID, traitUUID)
	return err
}
func (r *Repository) IncrementSessionVersion(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
		UPDATE users
		SET session_version = session_version + 1
		WHERE user_id = $1
	`, userID)
	return err
}
func (r *Repository) SetBlockedState(ctx context.Context, userID uuid.UUID, blocked bool, reason string, now time.Time) error {
	if blocked {
		_, err := r.db.Exec(ctx, `
			UPDATE users
			SET status = 'blocked',
				blocked_reason = $2,
				blocked_at = $3
			WHERE user_id = $1
		`, userID, reason, now)
		return err
	}
	_, err := r.db.Exec(ctx, `
		UPDATE users
		SET status = CASE WHEN email_verified_at IS NULL THEN 'pending_verification' ELSE 'active' END,
			blocked_reason = '',
			blocked_at = NULL
		WHERE user_id = $1
	`, userID)
	return err
}
func (r *Repository) CreateLoginAttempt(ctx context.Context, attempt *models.LoginAttempt) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO auth_login_attempts (
			attempt_id, user_id, email_normalized, ip, user_agent,
			outcome, failure_reason, risk_score, fingerprint_hash, suspicious_reason
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`, attempt.AttemptID, attempt.UserID, attempt.EmailNormalized, nullableINET(attempt.IP), attempt.UserAgent,
		attempt.Outcome, attempt.FailureReason, attempt.RiskScore, attempt.FingerprintHash, attempt.SuspiciousReason,
	)
	return err
}
func (r *Repository) CreateSecurityEvent(ctx context.Context, event *models.SecurityEvent) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO auth_security_events (
			event_id, user_id, actor_user_id, session_id, category,
			event_type, severity, ip, user_agent, payload
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`, event.EventID, event.UserID, event.ActorUserID, event.SessionID, event.Category,
		event.EventType, event.Severity, nullableINET(event.IP), event.UserAgent, jsonOrDefault(event.Payload),
	)
	return err
}
func nullableINET(raw string) any {
	value := strings.TrimSpace(raw)
	if value == "" {
		return nil
	}
	return value
}
func (r *Repository) ListLoginAttempts(ctx context.Context, userID uuid.UUID, emailNormalized string, limit int) ([]models.LoginAttempt, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := r.db.Query(ctx, `
		SELECT
			attempt_id, user_id, email_normalized, COALESCE(host(ip), ''), user_agent,
			outcome, failure_reason, risk_score, fingerprint_hash, suspicious_reason, created_at
		FROM auth_login_attempts
		WHERE user_id = $1 OR email_normalized = $2
		ORDER BY created_at DESC
		LIMIT $3
	`, userID, emailNormalized, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]models.LoginAttempt, 0)
	for rows.Next() {
		var item models.LoginAttempt
		if err := rows.Scan(
			&item.AttemptID, &item.UserID, &item.EmailNormalized, &item.IP, &item.UserAgent,
			&item.Outcome, &item.FailureReason, &item.RiskScore, &item.FingerprintHash,
			&item.SuspiciousReason, &item.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
func (r *Repository) ListSecurityEvents(ctx context.Context, userID uuid.UUID, limit int) ([]models.SecurityEvent, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := r.db.Query(ctx, `
		SELECT
			event_id, user_id, actor_user_id, session_id, category,
			event_type, severity, COALESCE(host(ip), ''), user_agent, payload, created_at
		FROM auth_security_events
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]models.SecurityEvent, 0)
	for rows.Next() {
		var item models.SecurityEvent
		if err := rows.Scan(
			&item.EventID, &item.UserID, &item.ActorUserID, &item.SessionID, &item.Category,
			&item.EventType, &item.Severity, &item.IP, &item.UserAgent, &item.Payload, &item.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
func (r *Repository) ListUsersForAdmin(ctx context.Context, search string, status string, role string, limit int, offset int) ([]models.User, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}
	if offset < 0 {
		offset = 0
	}
	rows, err := r.db.Query(ctx, `
		SELECT
			user_id, email, email_normalized, password_hash, status,
			email_verified_at, display_name, locale, timezone, roles,
			is_two_factor_enabled, totp_secret_ciphertext, profile, settings, primary_trait_uuid,
			session_version, last_login_at, blocked_reason, blocked_at,
			created_at, updated_at
		FROM users
		WHERE
			($1 = '' OR email_normalized LIKE '%' || lower($1) || '%' OR lower(display_name) LIKE '%' || lower($1) || '%')
			AND ($2 = '' OR status = $2)
			AND ($3 = '' OR $3 = ANY(roles))
		ORDER BY created_at DESC
		LIMIT $4 OFFSET $5
	`, search, status, role, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]models.User, 0, limit)
	for rows.Next() {
		item, scanErr := scanUser(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		items = append(items, *item)
	}
	return items, rows.Err()
}
func (r *Repository) CountUsersForAdmin(ctx context.Context, search string, status string, role string) (int, error) {
	var total int
	err := r.db.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM users
		WHERE
			($1 = '' OR email_normalized LIKE '%' || lower($1) || '%' OR lower(display_name) LIKE '%' || lower($1) || '%')
			AND ($2 = '' OR status = $2)
			AND ($3 = '' OR $3 = ANY(roles))
	`, search, status, role).Scan(&total)
	return total, err
}
func (r *Repository) CountUsersByRole(ctx context.Context, role string) (int, error) {
	var total int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM users WHERE $1 = ANY(roles)`, role).Scan(&total)
	return total, err
}
func (r *Repository) CountUsersTotal(ctx context.Context) (int, error) {
	var total int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&total)
	return total, err
}
func (r *Repository) CountUsersCreatedSince(ctx context.Context, since time.Time) (int, error) {
	var total int
	err := r.db.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM users
		WHERE created_at > $1
	`, since).Scan(&total)
	return total, err
}
func (r *Repository) CountPendingKitchenRecipes(ctx context.Context) (int, error) {
	var total int
	err := r.db.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM kitchen_recipes
		WHERE moderation_status = 'pending'
	`).Scan(&total)
	return total, err
}
func (r *Repository) CountPendingKitchenRecipesCreatedSince(ctx context.Context, since time.Time) (int, error) {
	var total int
	err := r.db.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM kitchen_recipes
		WHERE moderation_status = 'pending'
		  AND created_at > $1
	`, since).Scan(&total)
	return total, err
}
func (r *Repository) CountKitchenRecipesByModerationStatus(ctx context.Context) (map[string]int, error) {
	totals := map[string]int{
		"all":      0,
		"pending":  0,
		"approved": 0,
		"rejected": 0,
		"draft":    0,
	}
	rows, err := r.db.Query(ctx, `
		SELECT moderation_status, COUNT(*)::INT
		FROM kitchen_recipes
		GROUP BY moderation_status
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var status string
		var count int
		if scanErr := rows.Scan(&status, &count); scanErr != nil {
			return nil, scanErr
		}
		key := strings.ToLower(strings.TrimSpace(status))
		if key != "pending" && key != "approved" && key != "rejected" && key != "draft" {
			continue
		}
		totals[key] = count
		totals["all"] += count
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return totals, nil
}
func (r *Repository) GetAdminNotificationsReadAt(ctx context.Context, userID uuid.UUID) (*time.Time, error) {
	var raw string
	err := r.db.QueryRow(ctx, `
		SELECT COALESCE(settings->>'admin_notifications_read_at', '')
		FROM users
		WHERE user_id = $1
	`, userID).Scan(&raw)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	ts, parseErr := time.Parse(time.RFC3339Nano, raw)
	if parseErr != nil {
		return nil, nil
	}
	return &ts, nil
}
func (r *Repository) GetUserNotificationsReadAt(ctx context.Context, userID uuid.UUID) (*time.Time, error) {
	var raw string
	err := r.db.QueryRow(ctx, `
		SELECT COALESCE(settings->>'user_notifications_read_at', '')
		FROM users
		WHERE user_id = $1
	`, userID).Scan(&raw)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	ts, parseErr := time.Parse(time.RFC3339Nano, raw)
	if parseErr != nil {
		return nil, nil
	}
	return &ts, nil
}
func (r *Repository) SetAdminNotificationsReadAt(ctx context.Context, userID uuid.UUID, at time.Time) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE users
		SET settings = jsonb_set(
			COALESCE(settings, '{}'::jsonb),
			'{admin_notifications_read_at}',
			to_jsonb($2::text),
			true
		)
		WHERE user_id = $1
	`, userID, at.UTC().Format(time.RFC3339Nano))
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrUserNotFound
	}
	return nil
}
func (r *Repository) SetUserNotificationsReadAt(ctx context.Context, userID uuid.UUID, at time.Time) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE users
		SET settings = jsonb_set(
			COALESCE(settings, '{}'::jsonb),
			'{user_notifications_read_at}',
			to_jsonb($2::text),
			true
		)
		WHERE user_id = $1
	`, userID, at.UTC().Format(time.RFC3339Nano))
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrUserNotFound
	}
	return nil
}
func (r *Repository) CountOwnedKitchenRecipeModerationEventsSince(ctx context.Context, userID uuid.UUID, since time.Time) (map[string]int, error) {
	totals := map[string]int{
		"approved": 0,
		"rejected": 0,
	}
	rows, err := r.db.Query(ctx, `
		SELECT moderation_status, COUNT(*)::INT
		FROM kitchen_recipes
		WHERE owner_user_id = $1
			AND moderated_at IS NOT NULL
			AND moderated_at > $2
			AND moderation_status IN ('approved', 'rejected')
		GROUP BY moderation_status
	`, userID, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var status string
		var count int
		if scanErr := rows.Scan(&status, &count); scanErr != nil {
			return nil, scanErr
		}
		key := strings.ToLower(strings.TrimSpace(status))
		if key != "approved" && key != "rejected" {
			continue
		}
		totals[key] = count
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return totals, nil
}
func (r *Repository) GetPreviousLoginAt(ctx context.Context, userID uuid.UUID, currentSessionID *uuid.UUID) (*time.Time, error) {
	if currentSessionID != nil {
		var ts time.Time
		err := r.db.QueryRow(ctx, `
			SELECT created_at
			FROM auth_sessions
			WHERE user_id = $1
			  AND session_id <> $2
			ORDER BY created_at DESC
			LIMIT 1
		`, userID, *currentSessionID).Scan(&ts)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, nil
			}
			return nil, err
		}
		return &ts, nil
	}
	var ts time.Time
	err := r.db.QueryRow(ctx, `
		SELECT created_at
		FROM auth_sessions
		WHERE user_id = $1
		ORDER BY created_at DESC
		OFFSET 1
		LIMIT 1
	`, userID).Scan(&ts)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &ts, nil
}
func (r *Repository) SetUserPrimaryRole(ctx context.Context, userID uuid.UUID, role string) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE users
		SET roles = ARRAY[$2]::TEXT[]
		WHERE user_id = $1
	`, userID, role)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrUserNotFound
	}
	return nil
}
func (r *Repository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM users WHERE user_id = $1`, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrUserNotFound
	}
	return nil
}
func (r *Repository) SearchTraitKeysForAdmin(ctx context.Context, query string, limit int) ([]AdminTraitKeySearchItem, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	rows, err := r.db.Query(ctx, `
		SELECT
			k.id,
			s.name,
			COALESCE(k.meta, '{}'::jsonb),
			COUNT(v.t_uuid)::BIGINT
		FROM traits_k k
		JOIN key_syns s ON s.id = k.syn_id
		LEFT JOIN traits_v v ON v.t_key = k.id
		WHERE ($1 = '' OR lower(s.name) LIKE '%' || lower($1) || '%')
		GROUP BY k.id, s.name, k.meta
		ORDER BY
			CASE WHEN lower(s.name) = lower($1) THEN 0 ELSE 1 END,
			COUNT(v.t_uuid) DESC,
			s.name ASC
		LIMIT $2
	`, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]AdminTraitKeySearchItem, 0, limit)
	for rows.Next() {
		var item AdminTraitKeySearchItem
		var metaRaw []byte
		if err := rows.Scan(&item.KeyID, &item.Syn, &metaRaw, &item.TraitCount); err != nil {
			return nil, err
		}
		item.Meta = jsonMap(metaRaw)
		items = append(items, item)
	}
	return items, rows.Err()
}
func (r *Repository) GetTraitsSetsAnalysisForAdmin(ctx context.Context) (*AdminTraitsSetsAnalysis, error) {
	out := &AdminTraitsSetsAnalysis{}
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM traits_v`).Scan(&out.TotalTraits); err != nil {
		return nil, err
	}
	if err := r.db.QueryRow(ctx, `SELECT COUNT(DISTINCT (t_key, t_value)) FROM traits_v`).Scan(&out.UniqueTraitPairs); err != nil {
		return nil, err
	}
	if err := r.db.QueryRow(ctx, `SELECT COUNT(DISTINCT t_key) FROM traits_v`).Scan(&out.UniqueTraitKeys); err != nil {
		return nil, err
	}
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM sets`).Scan(&out.TotalSets); err != nil {
		return nil, err
	}
	if err := r.db.QueryRow(ctx, `SELECT COUNT(DISTINCT s_childs) FROM sets`).Scan(&out.UniqueSetsByChildren); err != nil {
		return nil, err
	}
	if err := r.db.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM sets s
		WHERE EXISTS (SELECT 1 FROM sets nested WHERE nested.s_uuid = s.s_childs[1])
	`).Scan(&out.DerivedSets); err != nil {
		return nil, err
	}
	if err := r.db.QueryRow(ctx, `
		SELECT COUNT(DISTINCT t.t_uuid)
		FROM traits_v t
		JOIN (
			SELECT unnest(s_childs) AS child_uuid
			FROM sets
		) refs ON refs.child_uuid = t.t_uuid
	`).Scan(&out.TraitsReferencedInSets); err != nil {
		return nil, err
	}
	if out.TotalTraits > 0 {
		out.OrphanTraits = out.TotalTraits - out.TraitsReferencedInSets
		out.TraitCoverageInSetsRate = float64(out.TraitsReferencedInSets) / float64(out.TotalTraits)
	}
	if out.TotalSets > 0 {
		out.DerivedSetRate = float64(out.DerivedSets) / float64(out.TotalSets)
		out.SetUniquenessRate = float64(out.UniqueSetsByChildren) / float64(out.TotalSets)
	}
	topRows, err := r.db.Query(ctx, `
		SELECT
			k.id,
			s.name,
			COUNT(v.t_uuid)::BIGINT AS trait_count
		FROM traits_k k
		JOIN key_syns s ON s.id = k.syn_id
		LEFT JOIN traits_v v ON v.t_key = k.id
		GROUP BY k.id, s.name
		ORDER BY trait_count DESC, s.name ASC
		LIMIT 10
	`)
	if err != nil {
		return nil, err
	}
	defer topRows.Close()
	top := make([]AdminTraitsSetsAnalysisTopKey, 0, 10)
	for topRows.Next() {
		var item AdminTraitsSetsAnalysisTopKey
		if err := topRows.Scan(&item.KeyID, &item.Syn, &item.TraitCount); err != nil {
			return nil, err
		}
		top = append(top, item)
	}
	if err := topRows.Err(); err != nil {
		return nil, err
	}
	out.TopKeys = top
	return out, nil
}
func scanUser(row pgx.Row) (*models.User, error) {
	var item models.User
	err := row.Scan(
		&item.UserID, &item.Email, &item.EmailNormalized, &item.PasswordHash, &item.Status,
		&item.EmailVerifiedAt, &item.DisplayName, &item.Locale, &item.Timezone, &item.Roles,
		&item.IsTwoFactorEnabled, &item.TOTPSecretCiphertext, &item.Profile, &item.Settings, &item.PrimaryTraitUUID,
		&item.SessionVersion, &item.LastLoginAt, &item.BlockedReason, &item.BlockedAt,
		&item.CreatedAt, &item.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUnauthorized
		}
		return nil, err
	}
	return &item, nil
}
func scanActionToken(row pgx.Row) (*models.ActionToken, error) {
	var item models.ActionToken
	err := row.Scan(
		&item.TokenID, &item.UserID, &item.Purpose, &item.TokenHash,
		&item.ExpiresAt, &item.ConsumedAt, &item.DeliveryValue, &item.Meta, &item.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrActionTokenNotFound
		}
		return nil, err
	}
	return &item, nil
}
func scanSession(row pgx.Row) (*models.AuthSession, error) {
	var item models.AuthSession
	err := row.Scan(
		&item.SessionID, &item.UserID, &item.FamilyID, &item.RefreshJTI, &item.RefreshTokenHash,
		&item.CSRFTokenHash, &item.FingerprintHash, &item.DeviceLabel, &item.IP, &item.UserAgent,
		&item.MFAVerified, &item.ReplacedBySession, &item.LastSeenAt, &item.ExpiresAt, &item.RevokedAt,
		&item.RevokeReason, &item.CompromisedAt, &item.CreatedAt, &item.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}
	return &item, nil
}
func jsonOrDefault(raw []byte) []byte {
	if len(raw) == 0 {
		return []byte("{}")
	}
	return raw
}
func scanSavedTraitSet(row pgx.Row) (*models.SavedTraitSet, error) {
	var item models.SavedTraitSet
	err := row.Scan(
		&item.SavedSetID,
		&item.UserID,
		&item.TraitUUID,
		&item.Name,
		&item.Description,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSavedTraitSetNotFound
		}
		return nil, err
	}
	return &item, nil
}
