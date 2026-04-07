package models
import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)
type PaymentOrder struct {
	OrderID            uuid.UUID       `json:"order_id"`
	PublicToken        string          `json:"-"`
	UserID             uuid.UUID       `json:"user_id"`
	UserEmail          string          `json:"user_email"`
	UserDisplayName    string          `json:"user_display_name"`
	Provider           string          `json:"provider"`
	PlanCode           string          `json:"plan_code"`
	BaseAmount         int64           `json:"base_amount"`
	Amount             int64           `json:"amount"`
	TipAmount          int64           `json:"tip_amount"`
	Currency           string          `json:"currency"`
	SubscriptionType   string          `json:"subscription_type"`
	Status             string          `json:"status"`
	ProviderStatus     string          `json:"provider_status"`
	ProviderPaymentID  string          `json:"provider_payment_id"`
	ProviderErrorCode  string          `json:"provider_error_code"`
	ProviderMessage    string          `json:"provider_message"`
	ProviderResponse   json.RawMessage `json:"provider_response"`
	ProviderFeePercent float64         `json:"provider_fee_percent"`
	FeeAmount          int64           `json:"fee_amount"`
	NetAmount          int64           `json:"net_amount"`
	AccessFrom         *time.Time      `json:"access_from,omitempty"`
	AccessUntil        *time.Time      `json:"access_until,omitempty"`
	ReturnTo           string          `json:"return_to"`
	SuccessURL         string          `json:"success_url"`
	FailURL            string          `json:"fail_url"`
	PaymentURL         string          `json:"payment_url"`
	LastCheckedAt      *time.Time      `json:"last_checked_at,omitempty"`
	NotifiedAt         *time.Time      `json:"notified_at,omitempty"`
	PaidAt             *time.Time      `json:"paid_at,omitempty"`
	FailedAt           *time.Time      `json:"failed_at,omitempty"`
	RefundedAt         *time.Time      `json:"refunded_at,omitempty"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
}
