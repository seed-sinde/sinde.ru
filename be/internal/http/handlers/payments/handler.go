package paymentshandlers

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"sinde.ru/internal/http/middleware"
	"sinde.ru/internal/http/responses"
	"sinde.ru/internal/payments"
)

type Handler struct {
	service *payments.Service
}

func New(service *payments.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateOrder() fiber.Handler {
	return func(c fiber.Ctx) error {
		user := middleware.CurrentUser(c)
		if user == nil {
			return responses.Error(c, fiber.StatusUnauthorized, "требуется аутентификация")
		}
		var input payments.CreateOrderInput
		if err := c.Bind().Body(&input); err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "некорректный запрос", err.Error())
		}
		input.RequestBaseURL = resolveRequestBaseURL(c)
		result, err := h.service.CreateOrder(c.Context(), user.UserID, input)
		if err != nil {
			return paymentError(c, err)
		}
		return responses.Success(c, fiber.StatusCreated, result)
	}
}

func (h *Handler) Access() fiber.Handler {
	return func(c fiber.Ctx) error {
		user := middleware.CurrentUser(c)
		if user == nil {
			return responses.Error(c, fiber.StatusUnauthorized, "требуется аутентификация")
		}
		result, err := h.service.GetAccessSummary(c.Context(), user.UserID)
		if err != nil {
			return paymentError(c, err)
		}
		return responses.Success(c, fiber.StatusOK, result)
	}
}

func (h *Handler) History() fiber.Handler {
	return func(c fiber.Ctx) error {
		user := middleware.CurrentUser(c)
		if user == nil {
			return responses.Error(c, fiber.StatusUnauthorized, "требуется аутентификация")
		}
		limit := 100
		if raw := strings.TrimSpace(c.Query("limit")); raw != "" {
			if value, err := strconv.Atoi(raw); err == nil {
				limit = value
			}
		}
		result, err := h.service.ListUserOrders(c.Context(), user.UserID, limit)
		if err != nil {
			return paymentError(c, err)
		}
		return responses.Success(c, fiber.StatusOK, result)
	}
}

func (h *Handler) PublicLookup() fiber.Handler {
	return func(c fiber.Ctx) error {
		var input payments.PublicOrderLookupInput
		if err := c.Bind().Body(&input); err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "некорректный запрос", err.Error())
		}
		result, err := h.service.LookupPublicOrder(c.Context(), input)
		if err != nil {
			return paymentError(c, err)
		}
		return responses.Success(c, fiber.StatusOK, result)
	}
}

func (h *Handler) Refund() fiber.Handler {
	return func(c fiber.Ctx) error {
		user := middleware.CurrentUser(c)
		if user == nil {
			return responses.Error(c, fiber.StatusUnauthorized, "требуется аутентификация")
		}
		orderID, err := uuid.Parse(strings.TrimSpace(c.Params("orderId")))
		if err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "некорректный идентификатор заказа")
		}
		result, err := h.service.RefundOrder(c.Context(), user.UserID, orderID)
		if err != nil {
			return paymentError(c, err)
		}
		return responses.Success(c, fiber.StatusOK, result)
	}
}

func (h *Handler) TBankNotification() fiber.Handler {
	return func(c fiber.Ctx) error {
		var payload map[string]any
		rawBody := c.Body()
		if err := json.Unmarshal(rawBody, &payload); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("invalid json")
		}
		if _, err := h.service.HandleNotification(c.Context(), payload, rawBody); err != nil {
			switch {
			case errors.Is(err, payments.ErrNotificationToken), errors.Is(err, payments.ErrOrderNotFound):
				return c.Status(fiber.StatusBadRequest).SendString("invalid notification")
			default:
				return c.Status(fiber.StatusInternalServerError).SendString("notification error")
			}
		}
		c.Status(fiber.StatusOK)
		return c.SendString("OK")
	}
}

func (h *Handler) AdminListOrders() fiber.Handler {
	return func(c fiber.Ctx) error {
		limit := 50
		if raw := strings.TrimSpace(c.Query("limit")); raw != "" {
			if value, err := strconv.Atoi(raw); err == nil {
				limit = value
			}
		}
		offset := 0
		if raw := strings.TrimSpace(c.Query("offset")); raw != "" {
			if value, err := strconv.Atoi(raw); err == nil {
				offset = value
			}
		}
		result, err := h.service.AdminListOrders(c.Context(), c.Query("q"), c.Query("status"), c.Query("plan"), limit, offset)
		if err != nil {
			return paymentError(c, err)
		}
		return responses.Success(c, fiber.StatusOK, result)
	}
}

func (h *Handler) AdminSummary() fiber.Handler {
	return func(c fiber.Ctx) error {
		result, err := h.service.AdminSummary(c.Context())
		if err != nil {
			return paymentError(c, err)
		}
		return responses.Success(c, fiber.StatusOK, result)
	}
}

func resolveRequestBaseURL(c fiber.Ctx) string {
	if origin := strings.TrimSpace(c.Get("Origin")); origin != "" {
		return strings.TrimRight(origin, "/")
	}
	if baseURL := strings.TrimSpace(c.BaseURL()); baseURL != "" {
		return strings.TrimRight(baseURL, "/")
	}
	return ""
}

func paymentError(c fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, payments.ErrPaymentsDisabled):
		return responses.Error(c, fiber.StatusServiceUnavailable, "приём платежей временно недоступен")
	case errors.Is(err, payments.ErrInvalidPlan), errors.Is(err, payments.ErrInvalidAmount):
		return responses.Error(c, fiber.StatusBadRequest, "некорректные параметры платежа")
	case errors.Is(err, payments.ErrOrderNotFound):
		return responses.Error(c, fiber.StatusNotFound, "заказ не найден")
	case errors.Is(err, payments.ErrOrderAccessDenied):
		return responses.Error(c, fiber.StatusForbidden, "доступ к заказу запрещен")
	case errors.Is(err, payments.ErrRefundNotAllowed):
		return responses.Error(c, fiber.StatusConflict, "возврат для этого платежа недоступен")
	case errors.Is(err, payments.ErrUserUnavailable):
		return responses.Error(c, fiber.StatusForbidden, "аккаунт не готов к оплате")
	case errors.Is(err, payments.ErrRequestOriginEmpty):
		return responses.Error(c, fiber.StatusBadRequest, "не удалось определить базовый адрес приложения для возврата платежа")
	case errors.Is(err, payments.ErrProviderRejected):
		return responses.Error(c, fiber.StatusBadGateway, "платёжный провайдер отклонил запрос")
	default:
		return responses.Error(c, fiber.StatusInternalServerError, "ошибка платёжного процесса", err.Error())
	}
}
