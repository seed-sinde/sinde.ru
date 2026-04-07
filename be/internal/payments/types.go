package payments

import (
	"github.com/google/uuid"
	"time"
)

const (
	ProviderTBank           = "tbank"
	PlanPro                 = "pro"
	PlanDonation            = "donation"
	SubscriptionTypeOneTime = "one_time"
	StatusPending           = "pending"
	StatusSuccess           = "success"
	StatusFailed            = "failed"
	StatusRefunded          = "refunded"
	StatusCanceled          = "canceled"
	ProAmountKopecks  int64 = 39900
)

type UserSnapshot struct {
	UserID          uuid.UUID
	Email           string
	DisplayName     string
	Status          string
	EmailVerifiedAt *time.Time
	Roles           []string
}

type CreateOrderInput struct {
	PlanCode       string `json:"plan_code"`
	Amount         int64  `json:"amount"`
	ReturnTo       string `json:"return_to"`
	RequestBaseURL string `json:"-"`
}

type AccessSummary struct {
	HasActiveAccess bool        `json:"has_active_access"`
	PlanCode        string      `json:"plan_code"`
	Amount          int64       `json:"amount"`
	TipAmount       int64       `json:"tip_amount"`
	Currency        string      `json:"currency"`
	AccessFrom      *time.Time  `json:"access_from,omitempty"`
	AccessUntil     *time.Time  `json:"access_until,omitempty"`
	OrderID         *uuid.UUID  `json:"order_id,omitempty"`
	LatestOrder     *OrderView  `json:"latest_order,omitempty"`
}

type OrderView struct {
	OrderID            uuid.UUID  `json:"order_id"`
	UserID             uuid.UUID  `json:"user_id"`
	UserEmail          string     `json:"user_email"`
	UserDisplayName    string     `json:"user_display_name"`
	PlanCode           string     `json:"plan_code"`
	BaseAmount         int64      `json:"base_amount"`
	Amount             int64      `json:"amount"`
	TipAmount          int64      `json:"tip_amount"`
	Currency           string     `json:"currency"`
	SubscriptionType   string     `json:"subscription_type"`
	Status             string     `json:"status"`
	Provider           string     `json:"provider"`
	ProviderStatus     string     `json:"provider_status"`
	ProviderPaymentID  string     `json:"provider_payment_id"`
	ProviderErrorCode  string     `json:"provider_error_code"`
	ProviderMessage    string     `json:"provider_message"`
	ProviderFeePercent float64    `json:"provider_fee_percent"`
	FeeAmount          int64      `json:"fee_amount"`
	NetAmount          int64      `json:"net_amount"`
	AccessFrom         *time.Time `json:"access_from,omitempty"`
	AccessUntil        *time.Time `json:"access_until,omitempty"`
	ReturnTo           string     `json:"return_to"`
	PaymentURL         string     `json:"payment_url"`
	LastCheckedAt      *time.Time `json:"last_checked_at,omitempty"`
	NotifiedAt         *time.Time `json:"notified_at,omitempty"`
	PaidAt             *time.Time `json:"paid_at,omitempty"`
	FailedAt           *time.Time `json:"failed_at,omitempty"`
	RefundedAt         *time.Time `json:"refunded_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

type CreateOrderResult struct {
	Order      OrderView `json:"order"`
	PaymentURL string    `json:"payment_url"`
}

type PublicOrderLookupInput struct {
	OrderID   string `json:"order_id"`
	Token     string `json:"token"`
	SyncState bool   `json:"sync_state"`
}

type PublicOrderLookupResult struct {
	Order OrderView `json:"order"`
}

type AdminOrdersListResult struct {
	Items  []OrderView `json:"items"`
	Total  int         `json:"total"`
	Limit  int         `json:"limit"`
	Offset int         `json:"offset"`
}

type AdminOrdersSummary struct {
	OrdersTotal        int       `json:"orders_total"`
	OrdersSuccess      int       `json:"orders_success"`
	OrdersPending      int       `json:"orders_pending"`
	OrdersFailed       int       `json:"orders_failed"`
	OrdersRefunded     int       `json:"orders_refunded"`
	PaidUsersTotal     int       `json:"paid_users_total"`
	PatronUsersTotal   int       `json:"patron_users_total"`
	ActiveAccessUsers  int       `json:"active_access_users"`
	GrossRevenue       int64     `json:"gross_revenue"`
	NetRevenue         int64     `json:"net_revenue"`
	TipRevenue         int64     `json:"tip_revenue"`
	MRR                int64     `json:"mrr"`
	ChurnRate          float64   `json:"churn_rate"`
	PatronShare        float64   `json:"patron_share"`
	LastSuccessfulPaid *time.Time `json:"last_successful_paid,omitempty"`
}

type tBankInitResponse struct {
	Success     bool   `json:"Success"`
	ErrorCode   string `json:"ErrorCode"`
	Message     string `json:"Message"`
	TerminalKey string `json:"TerminalKey"`
	Status      string `json:"Status"`
	PaymentID   string `json:"PaymentId"`
	OrderID     string `json:"OrderId"`
	Amount      int64  `json:"Amount"`
	PaymentURL  string `json:"PaymentURL"`
}

type tBankStateResponse struct {
	Success     bool   `json:"Success"`
	ErrorCode   string `json:"ErrorCode"`
	Message     string `json:"Message"`
	TerminalKey string `json:"TerminalKey"`
	Status      string `json:"Status"`
	PaymentID   string `json:"PaymentId"`
	OrderID     string `json:"OrderId"`
	Amount      int64  `json:"Amount"`
}
