package auth
import (
	"fmt"
	"time"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"sinde.ru/internal/models"
)
func (s *Service) issueTokens(user *models.User, sessionID uuid.UUID, familyID uuid.UUID, device DeviceContext, mfaVerified bool) (*TokenBundle, error) {
	now := s.now()
	accessJTI, err := randomToken(24)
	if err != nil {
		return nil, err
	}
	refreshJTI, err := randomToken(32)
	if err != nil {
		return nil, err
	}
	csrfToken, err := randomToken(24)
	if err != nil {
		return nil, err
	}
	accessExp := now.Add(s.cfg.AccessTTL)
	refreshExp := now.Add(s.cfg.RefreshTTL)
	accessClaims := &Claims{
		TokenType:      TokenTypeAccess,
		UserID:         user.UserID.String(),
		SessionID:      sessionID.String(),
		Email:          user.Email,
		Roles:          append([]string(nil), user.Roles...),
		MFAVerified:    mfaVerified,
		Fingerprint:    device.FingerprintHash,
		SessionVersion: user.SessionVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.cfg.JWTIssuer,
			Subject:   user.UserID.String(),
			Audience:  []string{s.cfg.JWTAudience},
			ExpiresAt: jwt.NewNumericDate(accessExp),
			NotBefore: jwt.NewNumericDate(now.Add(-10 * time.Second)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        accessJTI,
		},
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS512, accessClaims).SignedString(s.cfg.JWTSecret)
	if err != nil {
		return nil, err
	}
	refreshClaims := &Claims{
		TokenType:      TokenTypeRefresh,
		UserID:         user.UserID.String(),
		SessionID:      sessionID.String(),
		Email:          user.Email,
		Roles:          append([]string(nil), user.Roles...),
		MFAVerified:    mfaVerified,
		Fingerprint:    device.FingerprintHash,
		SessionVersion: user.SessionVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.cfg.JWTIssuer,
			Subject:   user.UserID.String(),
			Audience:  []string{s.cfg.JWTAudience},
			ExpiresAt: jwt.NewNumericDate(refreshExp),
			NotBefore: jwt.NewNumericDate(now.Add(-10 * time.Second)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        fmt.Sprintf("%s.%s", familyID.String(), refreshJTI),
		},
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS512, refreshClaims).SignedString(s.cfg.JWTSecret)
	if err != nil {
		return nil, err
	}
	return &TokenBundle{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		RefreshJTI:       refreshClaims.ID,
		CSRFToken:        csrfToken,
		AccessExpiresAt:  accessExp,
		RefreshExpiresAt: refreshExp,
		Claims:           accessClaims,
	}, nil
}
func (s *Service) parseToken(raw string, expectedType string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(raw, claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS512 {
			return nil, fmt.Errorf("unexpected jwt algorithm")
		}
		return s.cfg.JWTSecret, nil
	}, jwt.WithAudience(s.cfg.JWTAudience), jwt.WithIssuer(s.cfg.JWTIssuer))
	if err != nil {
		return nil, ErrInvalidToken
	}
	if !token.Valid {
		return nil, ErrInvalidToken
	}
	if claims.TokenType != expectedType {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
