package payments

import "errors"

var (
	ErrPaymentsDisabled   = errors.New("payments are disabled")
	ErrInvalidPlan        = errors.New("invalid payment plan")
	ErrInvalidAmount      = errors.New("invalid payment amount")
	ErrOrderNotFound      = errors.New("payment order not found")
	ErrOrderAccessDenied  = errors.New("payment order access denied")
	ErrRefundNotAllowed   = errors.New("payment refund is not allowed")
	ErrProviderRejected   = errors.New("payment provider rejected request")
	ErrNotificationToken  = errors.New("invalid payment notification token")
	ErrUnsupportedStatus  = errors.New("unsupported payment status")
	ErrUserUnavailable    = errors.New("payment user is unavailable")
	ErrRequestOriginEmpty = errors.New("payment origin is unavailable")
)
