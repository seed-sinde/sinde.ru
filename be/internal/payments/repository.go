package payments

import (
	"context"
	"encoding/json"
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

const paymentOrderSelect = `
	order_id, public_token, user_id, user_email, user_display_name, provider,
	plan_code, base_amount, amount, tip_amount, currency, subscription_type, status,
	provider_status, provider_payment_id, provider_error_code, provider_message, provider_response,
	provider_fee_percent, fee_amount, net_amount, access_from, access_until,
	return_to, success_url, fail_url, payment_url, last_checked_at, notified_at,
	paid_at, failed_at, refunded_at, created_at, updated_at
`

func paymentJSONOrDefault(raw json.RawMessage) []byte {
	if len(strings.TrimSpace(string(raw))) == 0 {
		return []byte("{}")
	}
	return raw
}
func scanPaymentOrder(scanner interface{ Scan(dest ...any) error }) (*models.PaymentOrder, error) {
	item := &models.PaymentOrder{}
	var providerResponse []byte
	if err := scanner.Scan(
		&item.OrderID,
		&item.PublicToken,
		&item.UserID,
		&item.UserEmail,
		&item.UserDisplayName,
		&item.Provider,
		&item.PlanCode,
		&item.BaseAmount,
		&item.Amount,
		&item.TipAmount,
		&item.Currency,
		&item.SubscriptionType,
		&item.Status,
		&item.ProviderStatus,
		&item.ProviderPaymentID,
		&item.ProviderErrorCode,
		&item.ProviderMessage,
		&providerResponse,
		&item.ProviderFeePercent,
		&item.FeeAmount,
		&item.NetAmount,
		&item.AccessFrom,
		&item.AccessUntil,
		&item.ReturnTo,
		&item.SuccessURL,
		&item.FailURL,
		&item.PaymentURL,
		&item.LastCheckedAt,
		&item.NotifiedAt,
		&item.PaidAt,
		&item.FailedAt,
		&item.RefundedAt,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}
	item.ProviderResponse = paymentJSONOrDefault(providerResponse)
	return item, nil
}
func (r *Repository) GetUserSnapshot(ctx context.Context, userID uuid.UUID) (*UserSnapshot, error) {
	const query = `
		SELECT user_id, email, display_name, status, email_verified_at, roles
		FROM users
		WHERE user_id = $1
	`
	item := &UserSnapshot{}
	if err := r.db.QueryRow(ctx, query, userID).Scan(
		&item.UserID,
		&item.Email,
		&item.DisplayName,
		&item.Status,
		&item.EmailVerifiedAt,
		&item.Roles,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserUnavailable
		}
		return nil, err
	}
	return item, nil
}
func (r *Repository) CreateOrder(ctx context.Context, item *models.PaymentOrder) error {
	const query = `
		INSERT INTO payment_orders (
			order_id, public_token, user_id, user_email, user_display_name, provider,
			plan_code, base_amount, amount, tip_amount, currency, subscription_type, status,
			provider_status, provider_payment_id, provider_error_code, provider_message, provider_response,
			provider_fee_percent, fee_amount, net_amount, access_from, access_until, return_to, success_url,
			fail_url, payment_url, last_checked_at, notified_at, paid_at, failed_at, refunded_at
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7,	$8, $9, $10, $11, $12, $13,
			$14, $15, $16, $17, $18,
			$19, $20, $21, $22, $23, $24, $25,
			$26, $27, $28, $29, $30, $31, $32
		)
	`
	_, err := r.db.Exec(ctx, query,
		item.OrderID,
		item.PublicToken,
		item.UserID,
		item.UserEmail,
		item.UserDisplayName,
		item.Provider,
		item.PlanCode,
		item.BaseAmount,
		item.Amount,
		item.TipAmount,
		item.Currency,
		item.SubscriptionType,
		item.Status,
		item.ProviderStatus,
		item.ProviderPaymentID,
		item.ProviderErrorCode,
		item.ProviderMessage,
		paymentJSONOrDefault(item.ProviderResponse),
		item.ProviderFeePercent,
		item.FeeAmount,
		item.NetAmount,
		item.AccessFrom,
		item.AccessUntil,
		item.ReturnTo,
		item.SuccessURL,
		item.FailURL,
		item.PaymentURL,
		item.LastCheckedAt,
		item.NotifiedAt,
		item.PaidAt,
		item.FailedAt,
		item.RefundedAt,
	)
	return err
}
func (r *Repository) UpdateOrder(ctx context.Context, item *models.PaymentOrder) error {
	const query = `
	UPDATE payment_orders
	SET user_email = $2,
		user_display_name = $3,
		provider = $4,
		plan_code = $5,
		base_amount = $6,
		amount = $7,
		tip_amount = $8,
		currency = $9,
		subscription_type = $10,
		status = $11,
		provider_status = $12,
		provider_payment_id = $13,
		provider_error_code = $14,
		provider_message = $15,
		provider_response = $16,
		provider_fee_percent = $17,
		fee_amount = $18,
		net_amount = $19,
		access_from = $20,
		access_until = $21,
		return_to = $22,
		success_url = $23,
		fail_url = $24,
		payment_url = $25,
		last_checked_at = $26,
		notified_at = $27,
		paid_at = $28,
		failed_at = $29,
		refunded_at = $30
	WHERE order_id = $1
`
	tag, err := r.db.Exec(ctx, query,
		item.OrderID,
		item.UserEmail,
		item.UserDisplayName,
		item.Provider,
		item.PlanCode,
		item.BaseAmount,
		item.Amount,
		item.TipAmount,
		item.Currency,
		item.SubscriptionType,
		item.Status,
		item.ProviderStatus,
		item.ProviderPaymentID,
		item.ProviderErrorCode,
		item.ProviderMessage,
		paymentJSONOrDefault(item.ProviderResponse),
		item.ProviderFeePercent,
		item.FeeAmount,
		item.NetAmount,
		item.AccessFrom,
		item.AccessUntil,
		item.ReturnTo,
		item.SuccessURL,
		item.FailURL,
		item.PaymentURL,
		item.LastCheckedAt,
		item.NotifiedAt,
		item.PaidAt,
		item.FailedAt,
		item.RefundedAt,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrOrderNotFound
	}
	return nil
}
func (r *Repository) GetOrderByID(ctx context.Context, orderID uuid.UUID) (*models.PaymentOrder, error) {
	return scanPaymentOrder(r.db.QueryRow(ctx, `SELECT `+paymentOrderSelect+` FROM payment_orders WHERE order_id = $1`, orderID))
}
func (r *Repository) GetOrderByIDAndToken(ctx context.Context, orderID uuid.UUID, token string) (*models.PaymentOrder, error) {
	return scanPaymentOrder(r.db.QueryRow(ctx, `
		SELECT `+paymentOrderSelect+`
		FROM payment_orders
		WHERE order_id = $1 AND public_token = $2
	`, orderID, token))
}
func (r *Repository) GetOrderByIDAndUser(ctx context.Context, orderID uuid.UUID, userID uuid.UUID) (*models.PaymentOrder, error) {
	return scanPaymentOrder(r.db.QueryRow(ctx, `
		SELECT `+paymentOrderSelect+`
		FROM payment_orders
		WHERE order_id = $1 AND user_id = $2
	`, orderID, userID))
}
func (r *Repository) GetCurrentActiveAccessUntilByUser(ctx context.Context, userID uuid.UUID, at time.Time) (*time.Time, error) {
	var next *time.Time
	err := r.db.QueryRow(ctx, `
		SELECT MAX(access_until)
		FROM payment_orders
		WHERE user_id = $1
		  AND status = 'success'
		  AND access_until IS NOT NULL
		  AND access_until > $2
	`, userID, at).Scan(&next)
	return next, err
}
func (r *Repository) GetActiveOrderByUser(ctx context.Context, userID uuid.UUID, at time.Time) (*models.PaymentOrder, error) {
	return scanPaymentOrder(r.db.QueryRow(ctx, `
		SELECT `+paymentOrderSelect+`
		FROM payment_orders
		WHERE user_id = $1
		  AND status = 'success'
		  AND access_from IS NOT NULL
		  AND access_until IS NOT NULL
		  AND access_from <= $2
		  AND access_until > $2
		ORDER BY access_until DESC, created_at DESC
		LIMIT 1
	`, userID, at))
}
func (r *Repository) GetLatestOrderByUser(ctx context.Context, userID uuid.UUID) (*models.PaymentOrder, error) {
	return scanPaymentOrder(r.db.QueryRow(ctx, `
		SELECT `+paymentOrderSelect+`
		FROM payment_orders
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`, userID))
}

func (r *Repository) ListOrdersByUser(ctx context.Context, userID uuid.UUID, limit int) ([]models.PaymentOrder, error) {
	if limit <= 0 {
		limit = 100
	}
	if limit > 200 {
		limit = 200
	}
	rows, err := r.db.Query(ctx, `
		SELECT `+paymentOrderSelect+`
		FROM payment_orders
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]models.PaymentOrder, 0, limit)
	for rows.Next() {
		item, scanErr := scanPaymentOrder(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		items = append(items, *item)
	}
	return items, rows.Err()
}

func (r *Repository) ListOrdersForAdmin(ctx context.Context, search string, status string, planCode string, limit int, offset int) ([]models.PaymentOrder, error) {
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
		SELECT `+paymentOrderSelect+`
		FROM payment_orders
		WHERE
			($1 = '' OR lower(user_email) LIKE '%' || lower($1) || '%' OR lower(user_display_name) LIKE '%' || lower($1) || '%')
			AND ($2 = '' OR status = $2)
			AND ($3 = '' OR plan_code = $3)
		ORDER BY created_at DESC
		LIMIT $4 OFFSET $5
	`, search, status, planCode, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]models.PaymentOrder, 0, limit)
	for rows.Next() {
		item, scanErr := scanPaymentOrder(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		items = append(items, *item)
	}
	return items, rows.Err()
}
func (r *Repository) CountOrdersForAdmin(ctx context.Context, search string, status string, planCode string) (int, error) {
	var total int
	err := r.db.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM payment_orders
		WHERE
			($1 = '' OR lower(user_email) LIKE '%' || lower($1) || '%' OR lower(user_display_name) LIKE '%' || lower($1) || '%')
			AND ($2 = '' OR status = $2)
			AND ($3 = '' OR plan_code = $3)
	`, search, status, planCode).Scan(&total)
	return total, err
}
func (r *Repository) GetAdminSummary(ctx context.Context, now time.Time) (*AdminOrdersSummary, error) {
	periodStart := now.AddDate(0, 0, -30)
	summary := &AdminOrdersSummary{}
	err := r.db.QueryRow(ctx, `
		WITH status_counts AS (
			SELECT
				COUNT(*)::INT AS orders_total,
				COUNT(*) FILTER (WHERE status = 'success')::INT AS orders_success,
				COUNT(*) FILTER (WHERE status = 'pending')::INT AS orders_pending,
				COUNT(*) FILTER (WHERE status IN ('failed', 'canceled'))::INT AS orders_failed,
				COUNT(*) FILTER (WHERE status = 'refunded')::INT AS orders_refunded
			FROM payment_orders
		),
		success_orders AS (
			SELECT *
			FROM payment_orders
			WHERE status = 'success'
		),
		user_counts AS (
			SELECT
				COUNT(DISTINCT user_id)::INT AS paid_users_total,
				COUNT(DISTINCT user_id) FILTER (WHERE tip_amount > 0)::INT AS patron_users_total,
				COUNT(DISTINCT user_id) FILTER (WHERE access_from <= $1 AND access_until > $1)::INT AS active_access_users,
				COALESCE(SUM(amount), 0)::BIGINT AS gross_revenue,
				COALESCE(SUM(net_amount), 0)::BIGINT AS net_revenue,
				COALESCE(SUM(tip_amount), 0)::BIGINT AS tip_revenue,
				COALESCE(SUM(amount) FILTER (WHERE access_from <= $1 AND access_until > $1), 0)::BIGINT AS mrr,
				MAX(paid_at) AS last_successful_paid
			FROM success_orders
		),
		cohort AS (
			SELECT DISTINCT user_id
			FROM success_orders
			WHERE access_from <= $2
			  AND access_until > $2
		),
		churn AS (
			SELECT
				COUNT(*)::INT AS cohort_total,
				COUNT(*) FILTER (
					WHERE user_id NOT IN (
						SELECT user_id
						FROM success_orders
						WHERE access_until > $1
					)
				)::INT AS churned_total
			FROM cohort
		)
		SELECT
			sc.orders_total,
			sc.orders_success,
			sc.orders_pending,
			sc.orders_failed,
			sc.orders_refunded,
			uc.paid_users_total,
			uc.patron_users_total,
			uc.active_access_users,
			uc.gross_revenue,
			uc.net_revenue,
			uc.tip_revenue,
			uc.mrr,
			CASE
				WHEN ch.cohort_total = 0 THEN 0
				ELSE ch.churned_total::DOUBLE PRECISION / ch.cohort_total
			END AS churn_rate,
			CASE
				WHEN uc.paid_users_total = 0 THEN 0
				ELSE uc.patron_users_total::DOUBLE PRECISION / uc.paid_users_total
			END AS patron_share,
			uc.last_successful_paid
		FROM status_counts sc
		CROSS JOIN user_counts uc
		CROSS JOIN churn ch
	`, now, periodStart).Scan(
		&summary.OrdersTotal,
		&summary.OrdersSuccess,
		&summary.OrdersPending,
		&summary.OrdersFailed,
		&summary.OrdersRefunded,
		&summary.PaidUsersTotal,
		&summary.PatronUsersTotal,
		&summary.ActiveAccessUsers,
		&summary.GrossRevenue,
		&summary.NetRevenue,
		&summary.TipRevenue,
		&summary.MRR,
		&summary.ChurnRate,
		&summary.PatronShare,
		&summary.LastSuccessfulPaid,
	)
	if err != nil {
		return nil, err
	}
	return summary, nil
}
