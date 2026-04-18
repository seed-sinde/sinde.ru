package authhandlers

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"sinde.ru/internal/auth"
	"sinde.ru/internal/http/middleware"
	"sinde.ru/internal/http/responses"
	"sinde.ru/internal/models"
)

type Handler struct {
	service *auth.Service
}

func New(service *auth.Service) *Handler {
	return &Handler{service: service}
}
func (h *Handler) Service() *auth.Service {
	return h.service
}

const (
	summaryStreamHeartbeatInterval = 15 * time.Second
	summaryStreamRefreshInterval   = 20 * time.Second
)

func isAdminUser(user *models.User) bool {
	if user == nil {
		return false
	}
	for _, role := range user.Roles {
		if strings.EqualFold(strings.TrimSpace(role), "admin") {
			return true
		}
	}
	return false
}
func (h *Handler) buildAdminSummary(ctx context.Context, adminUserID uuid.UUID, sessionID string) (*auth.AdminSummaryResult, error) {
	return h.service.AdminSummary(ctx, adminUserID, sessionID)
}
func (h *Handler) buildUserSummary(ctx context.Context, user *models.User, sessionID string) (*auth.UserSummaryResult, error) {
	if user == nil {
		return nil, auth.ErrUnauthorized
	}
	result, err := h.service.UserSummary(ctx, user.UserID, sessionID)
	if err != nil {
		return nil, err
	}
	if !isAdminUser(user) {
		return result, nil
	}
	adminResult, err := h.service.AdminSummary(ctx, user.UserID, sessionID)
	if err != nil {
		return nil, err
	}
	result.Admin = &auth.UserSummaryAdminInfo{
		UsersTotal:                      adminResult.UsersTotal,
		PendingRecipesTotal:             adminResult.PendingRecipesTotal,
		RecipeStatusTotals:              adminResult.RecipeStatusTotals,
		NewUsersSinceLastLogin:          adminResult.NewUsersSinceLastLogin,
		NewPendingRecipesSinceLastLogin: adminResult.NewPendingRecipesSinceLastLogin,
	}
	return result, nil
}
func writeSSEEvent(w *bufio.Writer, event string, payload any) bool {
	data, err := json.Marshal(payload)
	if err != nil {
		return false
	}
	if _, err := fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event, data); err != nil {
		return false
	}
	return w.Flush() == nil
}
func writeSSEComment(w *bufio.Writer, comment string) bool {
	if _, err := fmt.Fprintf(w, ": %s\n\n", comment); err != nil {
		return false
	}
	return w.Flush() == nil
}
func (h *Handler) Register() fiber.Handler {
	return func(c fiber.Ctx) error {
		if !h.validateOrigin(c) {
			return responses.Error(c, fiber.StatusForbidden, "источник запроса не разрешен")
		}
		var input auth.RegisterInput
		if err := c.Bind().Body(&input); err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "некорректный запрос", err.Error())
		}
		result, err := h.service.Register(c.Context(), input, deviceFromCtx(c))
		if err != nil {
			return authError(c, err)
		}
		return responses.Success(c, fiber.StatusCreated, result)
	}
}
func (h *Handler) RequestEmailVerification() fiber.Handler {
	return func(c fiber.Ctx) error {
		if !h.validateOrigin(c) {
			return responses.Error(c, fiber.StatusForbidden, "источник запроса не разрешен")
		}
		var input auth.RequestEmailVerificationInput
		if err := c.Bind().Body(&input); err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "некорректный запрос", err.Error())
		}
		if err := h.service.RequestEmailVerification(c.Context(), input.Email); err != nil {
			return authError(c, err)
		}
		return responses.Success(c, fiber.StatusAccepted, fiber.Map{"queued": true})
	}
}
func (h *Handler) RequestEmailChange() fiber.Handler {
	return func(c fiber.Ctx) error {
		user := middleware.CurrentUser(c)
		if user == nil {
			return responses.Error(c, fiber.StatusUnauthorized, "требуется аутентификация")
		}
		if !h.validateOrigin(c) {
			return responses.Error(c, fiber.StatusForbidden, "источник запроса не разрешен")
		}
		var input auth.RequestEmailChangeInput
		if err := c.Bind().Body(&input); err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "некорректный запрос", err.Error())
		}
		if err := h.service.RequestEmailChange(c.Context(), user, input.Email, deviceFromCtx(c)); err != nil {
			return authError(c, err)
		}
		return responses.Success(c, fiber.StatusAccepted, fiber.Map{"queued": true})
	}
}
func (h *Handler) VerifyEmail() fiber.Handler {
	return func(c fiber.Ctx) error {
		if !h.validateOrigin(c) {
			return responses.Error(c, fiber.StatusForbidden, "источник запроса не разрешен")
		}
		var input auth.VerifyEmailInput
		if err := c.Bind().Body(&input); err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "некорректный запрос", err.Error())
		}
		result, err := h.service.VerifyEmail(c.Context(), input.Token, deviceFromCtx(c))
		if err != nil {
			return authError(c, err)
		}
		return responses.Success(c, fiber.StatusOK, result)
	}
}
func (h *Handler) Login() fiber.Handler {
	return func(c fiber.Ctx) error {
		if !h.validateOrigin(c) {
			return responses.Error(c, fiber.StatusForbidden, "источник запроса не разрешен")
		}
		var input auth.LoginInput
		if err := c.Bind().Body(&input); err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "некорректный запрос", err.Error())
		}
		result, bundle, err := h.service.Login(c.Context(), input, deviceFromCtx(c))
		if err != nil {
			return authError(c, err)
		}
		if bundle != nil {
			h.service.SetSessionCookies(c, bundle)
		}
		return responses.Success(c, fiber.StatusOK, result)
	}
}
func (h *Handler) Login2FA() fiber.Handler {
	return func(c fiber.Ctx) error {
		if !h.validateOrigin(c) {
			return responses.Error(c, fiber.StatusForbidden, "источник запроса не разрешен")
		}
		var input auth.TwoFactorLoginInput
		if err := c.Bind().Body(&input); err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "некорректный запрос", err.Error())
		}
		result, bundle, err := h.service.CompleteMFALogin(c.Context(), input, deviceFromCtx(c))
		if err != nil {
			return authError(c, err)
		}
		h.service.SetSessionCookies(c, bundle)
		return responses.Success(c, fiber.StatusOK, result)
	}
}
func (h *Handler) Refresh() fiber.Handler {
	return func(c fiber.Ctx) error {
		if !h.validateOrigin(c) {
			return responses.Error(c, fiber.StatusForbidden, "источник запроса не разрешен")
		}
		refreshToken := strings.TrimSpace(c.Cookies(h.service.RefreshCookieName()))
		csrfToken := strings.TrimSpace(c.Get("X-CSRF-Token"))
		result, bundle, err := h.service.Refresh(c.Context(), refreshToken, csrfToken, deviceFromCtx(c))
		if err != nil {
			return authError(c, err)
		}
		h.service.SetSessionCookies(c, bundle)
		return responses.Success(c, fiber.StatusOK, result)
	}
}
func (h *Handler) Logout() fiber.Handler {
	return func(c fiber.Ctx) error {
		if !h.validateOrigin(c) {
			return responses.Error(c, fiber.StatusForbidden, "источник запроса не разрешен")
		}
		accessToken := strings.TrimSpace(c.Cookies(h.service.AccessCookieName()))
		_ = h.service.Logout(c.Context(), accessToken)
		h.service.ClearSessionCookies(c)
		return responses.Success(c, fiber.StatusOK, fiber.Map{"logged_out": true})
	}
}
func (h *Handler) LogoutAll() fiber.Handler {
	return func(c fiber.Ctx) error {
		user := middleware.CurrentUser(c)
		claims := middleware.CurrentClaims(c)
		if user == nil {
			return responses.Error(c, fiber.StatusUnauthorized, "требуется аутентификация")
		}
		if err := h.service.LogoutAll(c.Context(), user.UserID, claims); err != nil {
			return authError(c, err)
		}
		h.service.ClearSessionCookies(c)
		return responses.Success(c, fiber.StatusOK, fiber.Map{"logged_out": true, "all_devices": true})
	}
}
func (h *Handler) Me() fiber.Handler {
	return func(c fiber.Ctx) error {
		user := middleware.CurrentUser(c)
		if user == nil {
			return responses.Error(c, fiber.StatusUnauthorized, "требуется аутентификация")
		}
		return responses.Success(c, fiber.StatusOK, auth.MeResult{
			User: auth.ToAuthUser(user),
		})
	}
}
func (h *Handler) UpdateMe() fiber.Handler {
	return func(c fiber.Ctx) error {
		user := middleware.CurrentUser(c)
		if user == nil {
			return responses.Error(c, fiber.StatusUnauthorized, "требуется аутентификация")
		}
		var input auth.UpdateProfileInput
		if err := c.Bind().Body(&input); err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "некорректный запрос", err.Error())
		}
		result, err := h.service.UpdateProfile(c.Context(), user.UserID, input)
		if err != nil {
			return authError(c, err)
		}
		return responses.Success(c, fiber.StatusOK, result)
	}
}
func (h *Handler) ListSavedTraitSets() fiber.Handler {
	return func(c fiber.Ctx) error {
		user := middleware.CurrentUser(c)
		if user == nil {
			return responses.Error(c, fiber.StatusUnauthorized, "требуется аутентификация")
		}
		result, err := h.service.ListSavedTraitSets(c.Context(), user.UserID)
		if err != nil {
			return authError(c, err)
		}
		return responses.Success(c, fiber.StatusOK, result)
	}
}
func (h *Handler) SetPrimaryTrait() fiber.Handler {
	return func(c fiber.Ctx) error {
		user := middleware.CurrentUser(c)
		if user == nil {
			return responses.Error(c, fiber.StatusUnauthorized, "требуется аутентификация")
		}
		var input auth.SetPrimaryTraitInput
		if err := c.Bind().Body(&input); err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "некорректный запрос", err.Error())
		}
		setUUID := strings.TrimSpace(input.SetUUID)
		if setUUID == "" {
			setUUID = strings.TrimSpace(input.TraitUUID)
		}
		result, err := h.service.SetPrimaryTrait(c.Context(), user.UserID, setUUID)
		if err != nil {
			return authError(c, err)
		}
		return responses.Success(c, fiber.StatusOK, result)
	}
}
func (h *Handler) SaveTraitSet() fiber.Handler {
	return func(c fiber.Ctx) error {
		user := middleware.CurrentUser(c)
		if user == nil {
			return responses.Error(c, fiber.StatusUnauthorized, "требуется аутентификация")
		}
		var input auth.SaveTraitSetInput
		if err := c.Bind().Body(&input); err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "некорректный запрос", err.Error())
		}
		if strings.TrimSpace(input.SetUUID) != "" && strings.TrimSpace(input.TraitUUID) == "" {
			input.TraitUUID = strings.TrimSpace(input.SetUUID)
		}
		result, err := h.service.SaveTraitSet(c.Context(), user.UserID, input)
		if err != nil {
			return authError(c, err)
		}
		return responses.Success(c, fiber.StatusCreated, result)
	}
}
func (h *Handler) UpdateSavedTraitSet() fiber.Handler {
	return func(c fiber.Ctx) error {
		user := middleware.CurrentUser(c)
		if user == nil {
			return responses.Error(c, fiber.StatusUnauthorized, "требуется аутентификация")
		}
		savedSetID, err := uuid.Parse(strings.TrimSpace(c.Params("id")))
		if err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "ID сохраненного набора должен быть UUID")
		}
		var input auth.UpdateSavedTraitSetInput
		if err := c.Bind().Body(&input); err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "некорректный запрос", err.Error())
		}
		result, err := h.service.UpdateSavedTraitSet(c.Context(), user.UserID, savedSetID, input)
		if err != nil {
			return authError(c, err)
		}
		return responses.Success(c, fiber.StatusOK, result)
	}
}
func (h *Handler) DeleteSavedTraitSet() fiber.Handler {
	return func(c fiber.Ctx) error {
		user := middleware.CurrentUser(c)
		if user == nil {
			return responses.Error(c, fiber.StatusUnauthorized, "требуется аутентификация")
		}
		savedSetID, err := uuid.Parse(strings.TrimSpace(c.Params("id")))
		if err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "ID сохраненного набора должен быть UUID")
		}
		if err := h.service.DeleteSavedTraitSet(c.Context(), user.UserID, savedSetID); err != nil {
			return authError(c, err)
		}
		return responses.Success(c, fiber.StatusOK, fiber.Map{"deleted": true})
	}
}
func (h *Handler) ListSessions() fiber.Handler {
	return func(c fiber.Ctx) error {
		user := middleware.CurrentUser(c)
		claims := middleware.CurrentClaims(c)
		if user == nil || claims == nil {
			return responses.Error(c, fiber.StatusUnauthorized, "требуется аутентификация")
		}
		result, err := h.service.ListSessions(c.Context(), user.UserID, claims.SessionID)
		if err != nil {
			return authError(c, err)
		}
		return responses.Success(c, fiber.StatusOK, result)
	}
}
func (h *Handler) PublicUserProfile() fiber.Handler {
	return func(c fiber.Ctx) error {
		targetID, err := uuid.Parse(strings.TrimSpace(c.Params("id")))
		if err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "ID пользователя должен быть UUID")
		}
		result, err := h.service.PublicUserProfile(c.Context(), targetID)
		if err != nil {
			return authError(c, err)
		}
		return responses.Success(c, fiber.StatusOK, result)
	}
}
func (h *Handler) ListLoginAttempts() fiber.Handler {
	return func(c fiber.Ctx) error {
		user := middleware.CurrentUser(c)
		if user == nil {
			return responses.Error(c, fiber.StatusUnauthorized, "требуется аутентификация")
		}
		result, err := h.service.ListLoginAttempts(c.Context(), user, 50)
		if err != nil {
			return authError(c, err)
		}
		return responses.Success(c, fiber.StatusOK, result)
	}
}
func (h *Handler) ListSecurityEvents() fiber.Handler {
	return func(c fiber.Ctx) error {
		user := middleware.CurrentUser(c)
		if user == nil {
			return responses.Error(c, fiber.StatusUnauthorized, "требуется аутентификация")
		}
		result, err := h.service.ListSecurityEvents(c.Context(), user.UserID, 50)
		if err != nil {
			return authError(c, err)
		}
		return responses.Success(c, fiber.StatusOK, result)
	}
}
func (h *Handler) RevokeSession() fiber.Handler {
	return func(c fiber.Ctx) error {
		user := middleware.CurrentUser(c)
		claims := middleware.CurrentClaims(c)
		if user == nil {
			return responses.Error(c, fiber.StatusUnauthorized, "требуется аутентификация")
		}
		sessionID, err := uuid.Parse(strings.TrimSpace(c.Params("id")))
		if err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "ID сессии должен быть UUID")
		}
		if err := h.service.RevokeUserSession(c.Context(), user, sessionID, claims); err != nil {
			return authError(c, err)
		}
		if claims != nil && claims.SessionID == sessionID.String() {
			h.service.ClearSessionCookies(c)
		}
		return responses.Success(c, fiber.StatusOK, fiber.Map{"revoked": true})
	}
}
func (h *Handler) ForgotPassword() fiber.Handler {
	return func(c fiber.Ctx) error {
		if !h.validateOrigin(c) {
			return responses.Error(c, fiber.StatusForbidden, "источник запроса не разрешен")
		}
		var input auth.ForgotPasswordInput
		if err := c.Bind().Body(&input); err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "некорректный запрос", err.Error())
		}
		if err := h.service.ForgotPassword(c.Context(), input.Email); err != nil {
			return authError(c, err)
		}
		return responses.Success(c, fiber.StatusAccepted, fiber.Map{"queued": true})
	}
}
func (h *Handler) ResetPassword() fiber.Handler {
	return func(c fiber.Ctx) error {
		if !h.validateOrigin(c) {
			return responses.Error(c, fiber.StatusForbidden, "источник запроса не разрешен")
		}
		var input auth.ResetPasswordInput
		if err := c.Bind().Body(&input); err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "некорректный запрос", err.Error())
		}
		if err := h.service.ResetPassword(c.Context(), input); err != nil {
			return authError(c, err)
		}
		h.service.ClearSessionCookies(c)
		return responses.Success(c, fiber.StatusOK, fiber.Map{"password_reset": true})
	}
}
func (h *Handler) ChangePassword() fiber.Handler {
	return func(c fiber.Ctx) error {
		var input auth.ChangePasswordInput
		if err := c.Bind().Body(&input); err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "некорректный запрос", err.Error())
		}
		user := middleware.CurrentUser(c)
		claims := middleware.CurrentClaims(c)
		if user == nil {
			return responses.Error(c, fiber.StatusUnauthorized, "требуется аутентификация")
		}
		if err := h.service.ChangePassword(c.Context(), user, input, claims); err != nil {
			return authError(c, err)
		}
		h.service.ClearSessionCookies(c)
		return responses.Success(c, fiber.StatusOK, fiber.Map{"password_changed": true})
	}
}
func (h *Handler) SetupTwoFactor() fiber.Handler {
	return func(c fiber.Ctx) error {
		user := middleware.CurrentUser(c)
		if user == nil {
			return responses.Error(c, fiber.StatusUnauthorized, "требуется аутентификация")
		}
		result, err := h.service.SetupTwoFactor(c.Context(), user)
		if err != nil {
			return authError(c, err)
		}
		return responses.Success(c, fiber.StatusOK, result)
	}
}
func (h *Handler) EnableTwoFactor() fiber.Handler {
	return func(c fiber.Ctx) error {
		user := middleware.CurrentUser(c)
		if user == nil {
			return responses.Error(c, fiber.StatusUnauthorized, "требуется аутентификация")
		}
		var input auth.TwoFactorEnableInput
		if err := c.Bind().Body(&input); err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "некорректный запрос", err.Error())
		}
		codes, err := h.service.EnableTwoFactor(c.Context(), user, input.Code)
		if err != nil {
			return authError(c, err)
		}
		return responses.Success(c, fiber.StatusOK, fiber.Map{
			"enabled":      true,
			"backup_codes": codes,
		})
	}
}
func (h *Handler) DisableTwoFactor() fiber.Handler {
	return func(c fiber.Ctx) error {
		user := middleware.CurrentUser(c)
		if user == nil {
			return responses.Error(c, fiber.StatusUnauthorized, "требуется аутентификация")
		}
		var input auth.TwoFactorDisableInput
		if err := c.Bind().Body(&input); err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "некорректный запрос", err.Error())
		}
		if err := h.service.DisableTwoFactor(c.Context(), user, input); err != nil {
			return authError(c, err)
		}
		return responses.Success(c, fiber.StatusOK, fiber.Map{"disabled": true})
	}
}
func (h *Handler) AdminBlockUser() fiber.Handler {
	return func(c fiber.Ctx) error {
		admin := middleware.CurrentUser(c)
		if admin == nil {
			return responses.Error(c, fiber.StatusUnauthorized, "требуется аутентификация")
		}
		targetID, err := uuid.Parse(strings.TrimSpace(c.Params("id")))
		if err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "ID пользователя должен быть UUID")
		}
		var input auth.BlockUserInput
		if len(c.Body()) > 0 {
			if err := c.Bind().Body(&input); err != nil {
				return responses.Error(c, fiber.StatusBadRequest, "некорректный запрос", err.Error())
			}
		}
		if err := h.service.BlockUser(c.Context(), admin, targetID, input.Reason); err != nil {
			return authError(c, err)
		}
		return responses.Success(c, fiber.StatusOK, fiber.Map{"blocked": true})
	}
}
func (h *Handler) AdminUnblockUser() fiber.Handler {
	return func(c fiber.Ctx) error {
		admin := middleware.CurrentUser(c)
		if admin == nil {
			return responses.Error(c, fiber.StatusUnauthorized, "требуется аутентификация")
		}
		targetID, err := uuid.Parse(strings.TrimSpace(c.Params("id")))
		if err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "ID пользователя должен быть UUID")
		}
		if err := h.service.UnblockUser(c.Context(), admin, targetID); err != nil {
			return authError(c, err)
		}
		return responses.Success(c, fiber.StatusOK, fiber.Map{"unblocked": true})
	}
}
func (h *Handler) AdminForceLogoutUser() fiber.Handler {
	return func(c fiber.Ctx) error {
		admin := middleware.CurrentUser(c)
		if admin == nil {
			return responses.Error(c, fiber.StatusUnauthorized, "требуется аутентификация")
		}
		targetID, err := uuid.Parse(strings.TrimSpace(c.Params("id")))
		if err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "ID пользователя должен быть UUID")
		}
		if err := h.service.ForceLogoutUser(c.Context(), admin, targetID); err != nil {
			return authError(c, err)
		}
		return responses.Success(c, fiber.StatusOK, fiber.Map{"forced_logout": true})
	}
}
func (h *Handler) AdminDeleteUser() fiber.Handler {
	return func(c fiber.Ctx) error {
		admin := middleware.CurrentUser(c)
		if admin == nil {
			return responses.Error(c, fiber.StatusUnauthorized, "требуется аутентификация")
		}
		targetID, err := uuid.Parse(strings.TrimSpace(c.Params("id")))
		if err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "ID пользователя должен быть UUID")
		}
		if err := h.service.AdminDeleteUser(c.Context(), admin, targetID); err != nil {
			return authError(c, err)
		}
		return responses.Success(c, fiber.StatusOK, fiber.Map{"deleted": true})
	}
}
func (h *Handler) AdminListUsers() fiber.Handler {
	return func(c fiber.Ctx) error {
		limit := 50
		if raw := strings.TrimSpace(c.Query("limit")); raw != "" {
			if n, err := strconv.Atoi(raw); err == nil {
				limit = n
			}
		}
		offset := 0
		if raw := strings.TrimSpace(c.Query("offset")); raw != "" {
			if n, err := strconv.Atoi(raw); err == nil {
				offset = n
			}
		}
		result, err := h.service.AdminListUsers(
			c.Context(),
			c.Query("q"),
			c.Query("status"),
			c.Query("role"),
			limit,
			offset,
		)
		if err != nil {
			return authError(c, err)
		}
		return responses.Success(c, fiber.StatusOK, result)
	}
}
func (h *Handler) AdminUserDetail() fiber.Handler {
	return func(c fiber.Ctx) error {
		targetID, err := uuid.Parse(strings.TrimSpace(c.Params("id")))
		if err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "ID пользователя должен быть UUID")
		}
		result, err := h.service.AdminUserDetail(c.Context(), targetID)
		if err != nil {
			return authError(c, err)
		}
		return responses.Success(c, fiber.StatusOK, result)
	}
}
func (h *Handler) AdminSummary() fiber.Handler {
	return func(c fiber.Ctx) error {
		admin := middleware.CurrentUser(c)
		if admin == nil {
			return responses.Error(c, fiber.StatusUnauthorized, "требуется аутентификация")
		}
		claims := middleware.CurrentClaims(c)
		sessionID := ""
		if claims != nil {
			sessionID = claims.SessionID
		}
		result, err := h.buildAdminSummary(c.Context(), admin.UserID, sessionID)
		if err != nil {
			return authError(c, err)
		}
		return responses.Success(c, fiber.StatusOK, result)
	}
}
func (h *Handler) UserSummary() fiber.Handler {
	return func(c fiber.Ctx) error {
		user := middleware.CurrentUser(c)
		if user == nil {
			return responses.Error(c, fiber.StatusUnauthorized, "требуется аутентификация")
		}
		claims := middleware.CurrentClaims(c)
		sessionID := ""
		if claims != nil {
			sessionID = claims.SessionID
		}
		result, err := h.buildUserSummary(c.Context(), user, sessionID)
		if err != nil {
			return authError(c, err)
		}
		return responses.Success(c, fiber.StatusOK, result)
	}
}
func (h *Handler) UserSummaryStream() fiber.Handler {
	return func(c fiber.Ctx) error {
		user := middleware.CurrentUser(c)
		if user == nil {
			return responses.Error(c, fiber.StatusUnauthorized, "требуется аутентификация")
		}
		claims := middleware.CurrentClaims(c)
		sessionID := ""
		if claims != nil {
			sessionID = claims.SessionID
		}
		var lastPayload string
		c.Set(fiber.HeaderContentType, "text/event-stream")
		c.Set(fiber.HeaderCacheControl, "no-store")
		c.Set("X-Accel-Buffering", "no")
		c.Status(fiber.StatusOK)
		c.RequestCtx().Response.ImmediateHeaderFlush = true
		return c.SendStreamWriter(func(w *bufio.Writer) {
			summaryTicker := time.NewTicker(summaryStreamRefreshInterval)
			heartbeatTicker := time.NewTicker(summaryStreamHeartbeatInterval)
			defer summaryTicker.Stop()
			defer heartbeatTicker.Stop()
			sendSnapshot := func() bool {
				result, err := h.buildUserSummary(context.Background(), user, sessionID)
				if err != nil {
					return false
				}
				data, err := json.Marshal(result)
				if err != nil {
					return false
				}
				payload := string(data)
				if payload == lastPayload {
					return true
				}
				lastPayload = payload
				return writeSSEEvent(w, "summary", result)
			}
			if !writeSSEComment(w, "connected") {
				return
			}
			if !sendSnapshot() {
				return
			}
			for {
				select {
				case <-summaryTicker.C:
					if !sendSnapshot() {
						return
					}
				case <-heartbeatTicker.C:
					if !writeSSEComment(w, "keepalive") {
						return
					}
				}
			}
		})
	}
}
func (h *Handler) AdminSummaryStream() fiber.Handler {
	return func(c fiber.Ctx) error {
		admin := middleware.CurrentUser(c)
		if admin == nil {
			return responses.Error(c, fiber.StatusUnauthorized, "требуется аутентификация")
		}
		claims := middleware.CurrentClaims(c)
		sessionID := ""
		if claims != nil {
			sessionID = claims.SessionID
		}
		var lastPayload string
		c.Set(fiber.HeaderContentType, "text/event-stream")
		c.Set(fiber.HeaderCacheControl, "no-store")
		c.Set("X-Accel-Buffering", "no")
		c.Status(fiber.StatusOK)
		c.RequestCtx().Response.ImmediateHeaderFlush = true
		return c.SendStreamWriter(func(w *bufio.Writer) {
			summaryTicker := time.NewTicker(summaryStreamRefreshInterval)
			heartbeatTicker := time.NewTicker(summaryStreamHeartbeatInterval)
			defer summaryTicker.Stop()
			defer heartbeatTicker.Stop()
			sendSnapshot := func() bool {
				result, err := h.buildAdminSummary(context.Background(), admin.UserID, sessionID)
				if err != nil {
					return false
				}
				data, err := json.Marshal(result)
				if err != nil {
					return false
				}
				payload := string(data)
				if payload == lastPayload {
					return true
				}
				lastPayload = payload
				return writeSSEEvent(w, "summary", result)
			}
			if !writeSSEComment(w, "connected") {
				return
			}
			if !sendSnapshot() {
				return
			}
			for {
				select {
				case <-summaryTicker.C:
					if !sendSnapshot() {
						return
					}
				case <-heartbeatTicker.C:
					if !writeSSEComment(w, "keepalive") {
						return
					}
				}
			}
		})
	}
}
func (h *Handler) UserSummaryRead() fiber.Handler {
	return func(c fiber.Ctx) error {
		user := middleware.CurrentUser(c)
		if user == nil {
			return responses.Error(c, fiber.StatusUnauthorized, "требуется аутентификация")
		}
		if err := h.service.MarkUserSummaryRead(c.Context(), user.UserID); err != nil {
			return authError(c, err)
		}
		return responses.Success(c, fiber.StatusOK, fiber.Map{"read": true})
	}
}
func (h *Handler) AdminSummaryRead() fiber.Handler {
	return func(c fiber.Ctx) error {
		admin := middleware.CurrentUser(c)
		if admin == nil {
			return responses.Error(c, fiber.StatusUnauthorized, "требуется аутентификация")
		}
		if err := h.service.MarkAdminSummaryRead(c.Context(), admin.UserID); err != nil {
			return authError(c, err)
		}
		return responses.Success(c, fiber.StatusOK, fiber.Map{"read": true})
	}
}
func (h *Handler) AdminSetUserRole() fiber.Handler {
	return func(c fiber.Ctx) error {
		admin := middleware.CurrentUser(c)
		if admin == nil {
			return responses.Error(c, fiber.StatusUnauthorized, "требуется аутентификация")
		}
		targetID, err := uuid.Parse(strings.TrimSpace(c.Params("id")))
		if err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "ID пользователя должен быть UUID")
		}
		var input auth.AdminSetUserRoleInput
		if err := c.Bind().Body(&input); err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "некорректный запрос", err.Error())
		}
		updated, err := h.service.AdminSetUserRole(c.Context(), admin, targetID, input.Role)
		if err != nil {
			return authError(c, err)
		}
		return responses.Success(c, fiber.StatusOK, fiber.Map{"user": updated})
	}
}
func (h *Handler) AdminSearchTraitKeys() fiber.Handler {
	return func(c fiber.Ctx) error {
		limit := 20
		if raw := strings.TrimSpace(c.Query("limit")); raw != "" {
			if n, err := strconv.Atoi(raw); err == nil {
				limit = n
			}
		}
		result, err := h.service.AdminSearchTraitKeys(c.Context(), c.Query("q"), limit)
		if err != nil {
			return authError(c, err)
		}
		return responses.Success(c, fiber.StatusOK, result)
	}
}
func (h *Handler) AdminTraitsSetsAnalysis() fiber.Handler {
	return func(c fiber.Ctx) error {
		result, err := h.service.AdminTraitsSetsAnalysis(c.Context())
		if err != nil {
			return authError(c, err)
		}
		return responses.Success(c, fiber.StatusOK, result)
	}
}
func deviceFromCtx(c fiber.Ctx) auth.DeviceContext {
	return auth.BuildDeviceContext(c.IP(), c.Get("User-Agent"), c.Get("Accept-Language"))
}
func (h *Handler) validateOrigin(c fiber.Ctx) bool {
	origin := strings.TrimSpace(c.Get("Origin"))
	if origin != "" && !h.service.IsAllowedOrigin(origin) {
		return false
	}
	referer := strings.TrimSpace(c.Get("Referer"))
	if referer != "" {
		parsed, err := url.Parse(referer)
		if err == nil && parsed.Scheme != "" && parsed.Host != "" {
			refOrigin := parsed.Scheme + "://" + parsed.Host
			if !h.service.IsAllowedOrigin(refOrigin) {
				return false
			}
		}
	}
	return true
}
func authError(c fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, auth.ErrEmailAlreadyExists):
		return responses.Error(c, fiber.StatusConflict, "email уже зарегистрирован")
	case errors.Is(err, auth.ErrInvalidInput):
		return responses.Error(c, fiber.StatusBadRequest, "некорректные данные")
	case errors.Is(err, auth.ErrInvalidCredentials):
		return responses.Error(c, fiber.StatusUnauthorized, "неверный email или пароль")
	case errors.Is(err, auth.ErrEmailDeliveryFailed):
		return responses.Error(c, fiber.StatusServiceUnavailable, "не удалось отправить письмо подтверждения, попробуйте позже")
	case errors.Is(err, auth.ErrEmailNotVerified):
		return responses.Error(c, fiber.StatusForbidden, "email не подтвержден")
	case errors.Is(err, auth.ErrTooManyRequests):
		return responses.Error(c, fiber.StatusTooManyRequests, "слишком много запросов, попробуйте позже")
	case errors.Is(err, auth.ErrWeakPassword):
		return responses.Error(c, fiber.StatusBadRequest, "пароль не соответствует требованиям безопасности")
	case errors.Is(err, auth.ErrInvalidToken), errors.Is(err, auth.ErrActionTokenNotFound), errors.Is(err, auth.ErrInvalidRefreshToken):
		return responses.Error(c, fiber.StatusUnauthorized, "некорректный или просроченный токен")
	case errors.Is(err, auth.ErrUserBlocked):
		return responses.Error(c, fiber.StatusForbidden, "аккаунт заблокирован")
	case errors.Is(err, auth.ErrTwoFactorAlreadyOn):
		return responses.Error(c, fiber.StatusConflict, "2FA уже включена")
	case errors.Is(err, auth.ErrTwoFactorNotSetup):
		return responses.Error(c, fiber.StatusBadRequest, "сначала получите секрет для настройки 2FA")
	case errors.Is(err, auth.ErrTwoFactorNotEnabled):
		return responses.Error(c, fiber.StatusBadRequest, "2FA сейчас отключена")
	case errors.Is(err, auth.ErrInvalidMFACode):
		return responses.Error(c, fiber.StatusUnauthorized, "неверный код 2FA")
	case errors.Is(err, auth.ErrMFATicketExpired):
		return responses.Error(c, fiber.StatusUnauthorized, "Время подтверждения истекло. Войдите заново.")
	case errors.Is(err, auth.ErrForbidden):
		return responses.Error(c, fiber.StatusForbidden, "доступ запрещен")
	case errors.Is(err, auth.ErrUserNotFound):
		return responses.Error(c, fiber.StatusNotFound, "пользователь не найден")
	case errors.Is(err, auth.ErrSavedTraitSetNotFound):
		return responses.Error(c, fiber.StatusNotFound, "сохранённый набор не найден")
	case errors.Is(err, auth.ErrTraitTargetNotFound):
		return responses.Error(c, fiber.StatusNotFound, "набор или особенность больше не существуют")
	case errors.Is(err, auth.ErrRegistrationBlocked):
		return responses.Error(c, fiber.StatusForbidden, "регистрация временно недоступна")
	case errors.Is(err, auth.ErrSuspiciousLogin):
		return responses.Error(c, fiber.StatusUnauthorized, "подозрительная попытка входа")
	case errors.Is(err, auth.ErrSessionCompromised):
		return responses.Error(c, fiber.StatusUnauthorized, "сессия скомпрометирована")
	case errors.Is(err, auth.ErrPasswordChangeDenied):
		return responses.Error(c, fiber.StatusUnauthorized, "текущий пароль введен неверно")
	default:
		return responses.Error(c, fiber.StatusInternalServerError, "ошибка процесса аутентификации", err.Error())
	}
}
