package payments

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"sinde.ru/internal/models"
)

type Service struct {
	cfg    Config
	repo   *Repository
	client *http.Client
	now    func() time.Time
}

type Dependencies struct {
	Client *http.Client
	Now    func() time.Time
}

type receiptItem struct {
	Name   string
	Amount int64
}

func NewService(repo *Repository, cfg Config, deps Dependencies) *Service {
	client := deps.Client
	if client == nil {
		client = &http.Client{Timeout: cfg.RequestTimeout}
	}
	nowFn := deps.Now
	if nowFn == nil {
		nowFn = time.Now
	}
	return &Service{
		cfg:    cfg,
		repo:   repo,
		client: client,
		now:    nowFn,
	}
}

func fallbackString(value string, fallback string) string {
	value = strings.TrimSpace(value)
	if value != "" {
		return value
	}
	return strings.TrimSpace(fallback)
}

func normalizePlanCode(raw string) string {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case PlanDonation:
		return PlanDonation
	default:
		return PlanPro
	}
}

func hasRole(roles []string, role string) bool {
	needle := strings.ToLower(strings.TrimSpace(role))
	if needle == "" {
		return false
	}
	for _, item := range roles {
		if strings.ToLower(strings.TrimSpace(item)) == needle {
			return true
		}
	}
	return false
}

func resolveOrderAmounts(planCode string, requestedAmount int64, isAdmin bool) (string, int64, int64, int64, error) {
	switch normalizePlanCode(planCode) {
	case PlanDonation:
		if requestedAmount < MinPaymentAmountKopecks {
			return "", 0, 0, 0, ErrInvalidAmount
		}
		if requestedAmount < ProAmountKopecks {
			if !isAdmin {
				return "", 0, 0, 0, ErrInvalidAmount
			}
			return PlanPro, requestedAmount, requestedAmount, 0, nil
		}
		return PlanDonation, ProAmountKopecks, requestedAmount, requestedAmount - ProAmountKopecks, nil
	default:
		if isAdmin && requestedAmount > 0 && requestedAmount < ProAmountKopecks {
			if requestedAmount < MinPaymentAmountKopecks {
				return "", 0, 0, 0, ErrInvalidAmount
			}
			return PlanPro, requestedAmount, requestedAmount, 0, nil
		}
		return PlanPro, ProAmountKopecks, ProAmountKopecks, 0, nil
	}
}

func normalizeReturnTo(raw string) string {
	value := strings.TrimSpace(raw)
	if value == "" {
		return ""
	}
	if strings.Contains(value, "://") || strings.HasPrefix(value, "//") {
		return ""
	}
	if !strings.HasPrefix(value, "/") {
		value = "/" + value
	}
	return value
}

func randomPublicToken() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

func hashTBankToken(payload map[string]any, password string) string {
	pairs := make(map[string]string)
	for key, value := range payload {
		if strings.EqualFold(strings.TrimSpace(key), "Token") {
			continue
		}
		if value == nil {
			continue
		}
		kind := reflect.TypeOf(value).Kind()
		switch kind {
		case reflect.Map, reflect.Slice, reflect.Array, reflect.Struct:
			continue
		}
		pairs[key] = fmt.Sprint(value)
	}
	pairs["Password"] = password
	keys := make([]string, 0, len(pairs))
	for key := range pairs {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var builder strings.Builder
	for _, key := range keys {
		builder.WriteString(pairs[key])
	}
	sum := sha256.Sum256([]byte(builder.String()))
	return hex.EncodeToString(sum[:])
}

func buildURL(baseURL string, path string, values url.Values) string {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	path = "/" + strings.TrimLeft(strings.TrimSpace(path), "/")
	if baseURL == "" {
		return ""
	}
	parsed, err := url.Parse(baseURL + path)
	if err != nil {
		return ""
	}
	if len(values) > 0 {
		parsed.RawQuery = values.Encode()
	}
	return parsed.String()
}

func redactOrderLookupToken(raw string) string {
	value := strings.TrimSpace(raw)
	if value == "" {
		return ""
	}
	parsed, err := url.Parse(value)
	if err != nil {
		return value
	}
	query := parsed.Query()
	if query.Has("token") {
		query.Set("token", "[redacted]")
		parsed.RawQuery = query.Encode()
	}
	return parsed.String()
}

func redactReceiptForLog(value any) any {
	receipt, ok := value.(map[string]any)
	if !ok || receipt == nil {
		return value
	}
	clone := make(map[string]any, len(receipt))
	for key, item := range receipt {
		switch strings.TrimSpace(key) {
		case "Email":
			clone[key] = "[redacted]"
		default:
			clone[key] = item
		}
	}
	return clone
}

func debugPayloadForLog(payload map[string]any) map[string]any {
	if payload == nil {
		return map[string]any{}
	}
	debugPayload := make(map[string]any, len(payload))
	for key, value := range payload {
		switch strings.TrimSpace(key) {
		case "Token":
			debugPayload[key] = "[redacted]"
		case "SuccessURL", "FailURL", "NotificationURL":
			debugPayload[key] = redactOrderLookupToken(fmt.Sprint(value))
		case "Receipt":
			debugPayload[key] = redactReceiptForLog(value)
		default:
			debugPayload[key] = value
		}
	}
	return debugPayload
}

func buildRedirectURL(baseURL string, path string, orderID uuid.UUID, token string, paymentID string, status string, returnTo string) string {
	values := url.Values{}
	values.Set("order_id", orderID.String())
	values.Set("token", token)
	if paymentID != "" {
		values.Set("payment_id", paymentID)
	}
	if status != "" {
		values.Set("status", status)
	}
	if next := normalizeReturnTo(returnTo); next != "" {
		values.Set("next", next)
	}
	// T-Bank expects these placeholders literally in the return URL query.
	return strings.NewReplacer(
		"%24%7BPaymentId%7D", "${PaymentId}",
		"%24%7BStatus%7D", "${Status}",
	).Replace(buildURL(baseURL, path, values))
}

func normalizeReceiptEmail(email string) string {
	return strings.TrimSpace(email)
}

func (s *Service) buildReceiptItems(order *models.PaymentOrder) []receiptItem {
	if order == nil || order.Amount <= 0 {
		return nil
	}
	items := make([]receiptItem, 0, 2)
	baseAmount := order.Amount
	if order.BaseAmount > 0 && order.BaseAmount < baseAmount {
		baseAmount = order.BaseAmount
	}
	if baseAmount > 0 {
		items = append(items, receiptItem{
			Name:   s.cfg.ReceiptProTitle,
			Amount: baseAmount,
		})
	}
	if order.TipAmount > 0 {
		items = append(items, receiptItem{
			Name:   s.cfg.ReceiptDonationTitle,
			Amount: order.TipAmount,
		})
	}
	if len(items) == 0 {
		items = append(items, receiptItem{
			Name:   s.cfg.ReceiptProTitle,
			Amount: order.Amount,
		})
	}
	return items
}

func (s *Service) buildReceipt(order *models.PaymentOrder, email string) map[string]any {
	if !s.cfg.UsesProviderReceipt() || order == nil {
		return nil
	}
	email = normalizeReceiptEmail(email)
	if email == "" {
		return nil
	}
	sourceItems := s.buildReceiptItems(order)
	if len(sourceItems) == 0 {
		return nil
	}
	items := make([]map[string]any, 0, len(sourceItems))
	for _, item := range sourceItems {
		if item.Amount <= 0 {
			continue
		}
		items = append(items, map[string]any{
			"Name":          item.Name,
			"Price":         item.Amount,
			"Quantity":      1,
			"Amount":        item.Amount,
			"Tax":           s.cfg.ReceiptTax,
			"PaymentMethod": s.cfg.ReceiptPaymentMethod,
			"PaymentObject": s.cfg.ReceiptPaymentObject,
		})
	}
	if len(items) == 0 {
		return nil
	}
	return map[string]any{
		"Email":    email,
		"Taxation": s.cfg.ReceiptTaxation,
		"Items":    items,
		"Payments": map[string]any{"Electronic": order.Amount},
	}
}

func internalStatusFromProvider(providerStatus string, success bool) string {
	status := strings.ToUpper(strings.TrimSpace(providerStatus))
	switch status {
	case "CONFIRMED":
		return StatusSuccess
	case "REFUNDED", "PARTIAL_REFUNDED":
		return StatusRefunded
	case "REVERSED", "CANCELED":
		return StatusCanceled
	case "AUTH_FAIL", "REJECTED", "DEADLINE_EXPIRED":
		return StatusFailed
	case "AUTHORIZED", "NEW", "FORM_SHOWED", "3DS_CHECKING", "CHECKING", "COMPLETING":
		return StatusPending
	default:
		if success {
			return StatusPending
		}
		return StatusFailed
	}
}

func (s *Service) callTBank(ctx context.Context, method string, payload map[string]any, out any) ([]byte, error) {
	if !s.cfg.Enabled() {
		return nil, ErrPaymentsDisabled
	}
	payload["TerminalKey"] = s.cfg.TBankTerminalKey
	payload["Token"] = hashTBankToken(payload, s.cfg.TBankPassword)
	debugBody, debugErr := json.Marshal(debugPayloadForLog(payload))
	if debugErr != nil {
		log.Printf("[payments:tbank] request method=%s payload_log_error=%v", method, debugErr)
	} else {
		log.Printf("[payments:tbank] request method=%s url=%s payload=%s", method, strings.TrimRight(s.cfg.TBankAPIURL, "/")+"/"+strings.TrimLeft(method, "/"), string(debugBody))
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(s.cfg.TBankAPIURL, "/")+"/"+strings.TrimLeft(method, "/"), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := s.client.Do(req)
	if err != nil {
		log.Printf("[payments:tbank] transport_error method=%s err=%v", method, err)
		return nil, err
	}
	defer res.Body.Close()
	raw, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("[payments:tbank] read_error method=%s err=%v", method, err)
		return nil, err
	}
	log.Printf("[payments:tbank] response method=%s status_code=%d body=%s", method, res.StatusCode, strings.TrimSpace(string(raw)))
	if out != nil {
		if err := json.Unmarshal(raw, out); err != nil {
			log.Printf("[payments:tbank] decode_error method=%s err=%v", method, err)
			return raw, err
		}
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return raw, fmt.Errorf("%w: %s", ErrProviderRejected, strings.TrimSpace(string(raw)))
	}
	return raw, nil
}

func (s *Service) buildOrderView(item *models.PaymentOrder) OrderView {
	return OrderView{
		OrderID:            item.OrderID,
		UserID:             item.UserID,
		UserEmail:          item.UserEmail,
		UserDisplayName:    item.UserDisplayName,
		PlanCode:           item.PlanCode,
		BaseAmount:         item.BaseAmount,
		Amount:             item.Amount,
		TipAmount:          item.TipAmount,
		Currency:           item.Currency,
		SubscriptionType:   item.SubscriptionType,
		Status:             item.Status,
		Provider:           item.Provider,
		ProviderStatus:     item.ProviderStatus,
		ProviderPaymentID:  item.ProviderPaymentID,
		ProviderErrorCode:  item.ProviderErrorCode,
		ProviderMessage:    item.ProviderMessage,
		ProviderFeePercent: item.ProviderFeePercent,
		FeeAmount:          item.FeeAmount,
		NetAmount:          item.NetAmount,
		AccessFrom:         item.AccessFrom,
		AccessUntil:        item.AccessUntil,
		ReturnTo:           item.ReturnTo,
		PaymentURL:         item.PaymentURL,
		LastCheckedAt:      item.LastCheckedAt,
		NotifiedAt:         item.NotifiedAt,
		PaidAt:             item.PaidAt,
		FailedAt:           item.FailedAt,
		RefundedAt:         item.RefundedAt,
		CanRefund:          s.canRefundOrder(item),
		CreatedAt:          item.CreatedAt,
		UpdatedAt:          item.UpdatedAt,
	}
}

func (s *Service) canRefundOrder(item *models.PaymentOrder) bool {
	if item == nil {
		return false
	}
	return item.Provider == ProviderTBank &&
		strings.TrimSpace(item.ProviderPaymentID) != "" &&
		strings.EqualFold(strings.TrimSpace(item.Status), StatusSuccess)
}

func (s *Service) CreateOrder(ctx context.Context, userID uuid.UUID, input CreateOrderInput) (*CreateOrderResult, error) {
	if !s.cfg.Enabled() {
		return nil, ErrPaymentsDisabled
	}
	requestBaseURL := strings.TrimRight(strings.TrimSpace(input.RequestBaseURL), "/")
	if requestBaseURL == "" {
		requestBaseURL = strings.TrimRight(strings.TrimSpace(s.cfg.PublicBaseURL), "/")
	}
	if requestBaseURL == "" {
		return nil, ErrRequestOriginEmpty
	}
	user, err := s.repo.GetUserSnapshot(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user.Email == "" || user.Status != "active" || user.EmailVerifiedAt == nil {
		return nil, ErrUserUnavailable
	}
	planCode, baseAmount, amount, tipAmount, err := resolveOrderAmounts(input.PlanCode, input.Amount, hasRole(user.Roles, "admin"))
	if err != nil {
		return nil, err
	}
	now := s.now().UTC()
	orderID := uuid.New()
	publicToken, err := randomPublicToken()
	if err != nil {
		return nil, err
	}
	providerFee := s.cfg.FeePercent
	feeAmount := int64(float64(amount) * providerFee / 100)
	netAmount := amount - feeAmount
	returnTo := normalizeReturnTo(input.ReturnTo)
	successURL := buildRedirectURL(requestBaseURL, s.cfg.SuccessPath, orderID, publicToken, "${PaymentId}", "${Status}", returnTo)
	failURL := buildRedirectURL(requestBaseURL, s.cfg.FailPath, orderID, publicToken, "${PaymentId}", "${Status}", returnTo)
	notificationURL := buildURL(s.cfg.PublicBaseURL, s.cfg.NotificationPath, nil)
	order := &models.PaymentOrder{
		OrderID:            orderID,
		PublicToken:        publicToken,
		UserID:             user.UserID,
		UserEmail:          user.Email,
		UserDisplayName:    user.DisplayName,
		Provider:           ProviderTBank,
		PlanCode:           planCode,
		BaseAmount:         baseAmount,
		Amount:             amount,
		TipAmount:          tipAmount,
		Currency:           s.cfg.Currency,
		SubscriptionType:   SubscriptionTypeOneTime,
		Status:             StatusPending,
		ProviderStatus:     "NEW",
		ProviderResponse:   []byte("{}"),
		ProviderFeePercent: providerFee,
		FeeAmount:          feeAmount,
		NetAmount:          netAmount,
		ReturnTo:           returnTo,
		SuccessURL:         successURL,
		FailURL:            failURL,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	if err := s.repo.CreateOrder(ctx, order); err != nil {
		return nil, err
	}
	description := "sinde.ru Pro"
	if planCode == PlanDonation {
		description = "sinde.ru Donation"
	}
	payload := map[string]any{
		"Amount":      amount,
		"OrderId":     orderID.String(),
		"Description": description,
		"CustomerKey": user.UserID.String(),
		"Language":    "ru",
		"SuccessURL":  successURL,
		"FailURL":     failURL,
		"DATA": map[string]string{
			"OperationInitiatorType": "0",
			"plan_code":              planCode,
			"user_id":                user.UserID.String(),
		},
	}
	if notificationURL != "" {
		payload["NotificationURL"] = notificationURL
	}
	if receipt := s.buildReceipt(order, user.Email); receipt != nil {
		payload["Receipt"] = receipt
	}
	var response tBankInitResponse
	rawResponse, err := s.callTBank(ctx, "Init", payload, &response)
	order.ProviderResponse = rawResponse
	order.ProviderStatus = strings.ToUpper(strings.TrimSpace(response.Status))
	order.ProviderPaymentID = strings.TrimSpace(response.PaymentID)
	order.ProviderErrorCode = strings.TrimSpace(response.ErrorCode)
	order.ProviderMessage = strings.TrimSpace(response.Message)
	order.PaymentURL = strings.TrimSpace(response.PaymentURL)
	order.LastCheckedAt = ptrTime(now)
	if err != nil {
		log.Printf("[payments:init] rejected order_id=%s user_id=%s plan=%s amount=%d err=%v provider_status=%s provider_error_code=%s provider_message=%q payment_url=%s",
			order.OrderID.String(),
			order.UserID.String(),
			order.PlanCode,
			order.Amount,
			err,
			order.ProviderStatus,
			order.ProviderErrorCode,
			order.ProviderMessage,
			strings.TrimSpace(order.PaymentURL),
		)
		order.Status = StatusFailed
		order.FailedAt = ptrTime(now)
		_ = s.repo.UpdateOrder(ctx, order)
		return nil, err
	}
	order.Status = internalStatusFromProvider(response.Status, response.Success)
	if order.Status == StatusFailed {
		order.FailedAt = ptrTime(now)
	}
	if updateErr := s.repo.UpdateOrder(ctx, order); updateErr != nil {
		return nil, updateErr
	}
	if !response.Success || strings.TrimSpace(response.ErrorCode) != "" && strings.TrimSpace(response.ErrorCode) != "0" {
		log.Printf("[payments:init] provider_declined order_id=%s user_id=%s plan=%s amount=%d provider_status=%s provider_error_code=%s provider_message=%q success=%t payment_id=%s payment_url=%s",
			order.OrderID.String(),
			order.UserID.String(),
			order.PlanCode,
			order.Amount,
			order.ProviderStatus,
			order.ProviderErrorCode,
			order.ProviderMessage,
			response.Success,
			order.ProviderPaymentID,
			strings.TrimSpace(order.PaymentURL),
		)
		return nil, ErrProviderRejected
	}
	return &CreateOrderResult{
		Order:      s.buildOrderView(order),
		PaymentURL: order.PaymentURL,
	}, nil
}

func ptrTime(value time.Time) *time.Time {
	next := value.UTC()
	return &next
}

func (s *Service) syncOrderState(ctx context.Context, order *models.PaymentOrder) (*models.PaymentOrder, error) {
	if order == nil {
		return nil, ErrOrderNotFound
	}
	if !s.cfg.Enabled() {
		return order, ErrPaymentsDisabled
	}
	now := s.now().UTC()
	payload := map[string]any{
		"OrderId": order.OrderID.String(),
	}
	method := "CheckOrder"
	if strings.TrimSpace(order.ProviderPaymentID) != "" {
		method = "GetState"
		payload = map[string]any{
			"PaymentId": order.ProviderPaymentID,
		}
	}
	var response tBankStateResponse
	rawResponse, err := s.callTBank(ctx, method, payload, &response)
	if err != nil {
		return nil, err
	}
	order.ProviderResponse = rawResponse
	order.ProviderStatus = strings.ToUpper(strings.TrimSpace(response.Status))
	order.ProviderPaymentID = fallbackString(strings.TrimSpace(response.PaymentID), order.ProviderPaymentID)
	order.ProviderErrorCode = strings.TrimSpace(response.ErrorCode)
	order.ProviderMessage = strings.TrimSpace(response.Message)
	order.LastCheckedAt = ptrTime(now)
	nextStatus := internalStatusFromProvider(response.Status, response.Success)
	order.Status = nextStatus
	switch nextStatus {
	case StatusSuccess:
		if order.PaidAt == nil {
			order.PaidAt = ptrTime(now)
		}
		order.FailedAt = nil
		if order.AccessFrom == nil || order.AccessUntil == nil {
			accessStart := now
			currentUntil, currentErr := s.repo.GetCurrentActiveAccessUntilByUser(ctx, order.UserID, now)
			if currentErr != nil {
				return nil, currentErr
			}
			if currentUntil != nil && currentUntil.After(accessStart) {
				accessStart = currentUntil.UTC()
			}
			accessUntil := accessStart.AddDate(0, s.cfg.AccessDurationMonths, 0)
			order.AccessFrom = ptrTime(accessStart)
			order.AccessUntil = ptrTime(accessUntil)
			// TODO(payments): when access logging is introduced, record activation IP, activation timestamp,
			// and subsequent paid-feature usage from this access-grant point.
		}
	case StatusRefunded:
		if order.RefundedAt == nil {
			order.RefundedAt = ptrTime(now)
		}
	case StatusFailed, StatusCanceled:
		if order.FailedAt == nil {
			order.FailedAt = ptrTime(now)
		}
	}
	if err := s.repo.UpdateOrder(ctx, order); err != nil {
		return nil, err
	}
	return order, nil
}

func (s *Service) GetAccessSummary(ctx context.Context, userID uuid.UUID) (*AccessSummary, error) {
	now := s.now().UTC()
	summary := &AccessSummary{}
	activeOrder, err := s.repo.GetActiveOrderByUser(ctx, userID, now)
	if err == nil {
		summary.HasActiveAccess = true
		summary.PlanCode = activeOrder.PlanCode
		summary.Amount = activeOrder.Amount
		summary.TipAmount = activeOrder.TipAmount
		summary.Currency = activeOrder.Currency
		summary.AccessFrom = activeOrder.AccessFrom
		summary.AccessUntil = activeOrder.AccessUntil
		summary.OrderID = &activeOrder.OrderID
	}
	if err != nil && !errors.Is(err, ErrOrderNotFound) {
		return nil, err
	}
	latestOrder, err := s.repo.GetLatestOrderByUser(ctx, userID)
	if err == nil {
		view := s.buildOrderView(latestOrder)
		summary.LatestOrder = &view
	}
	if err != nil && !errors.Is(err, ErrOrderNotFound) {
		return nil, err
	}
	return summary, nil
}

func (s *Service) ListUserOrders(ctx context.Context, userID uuid.UUID, limit int) (*UserOrdersListResult, error) {
	items, err := s.repo.ListOrdersByUser(ctx, userID, limit)
	if err != nil {
		return nil, err
	}
	out := make([]OrderView, 0, len(items))
	for _, item := range items {
		copyItem := item
		out = append(out, s.buildOrderView(&copyItem))
	}
	return &UserOrdersListResult{
		Items: out,
	}, nil
}

func (s *Service) LookupPublicOrder(ctx context.Context, input PublicOrderLookupInput) (*PublicOrderLookupResult, error) {
	orderID, err := uuid.Parse(strings.TrimSpace(input.OrderID))
	if err != nil {
		return nil, ErrOrderNotFound
	}
	order, err := s.repo.GetOrderByIDAndToken(ctx, orderID, strings.TrimSpace(input.Token))
	if err != nil {
		return nil, err
	}
	if input.SyncState && order.Status == StatusPending {
		order, err = s.syncOrderState(ctx, order)
		if err != nil && !errors.Is(err, ErrPaymentsDisabled) {
			return nil, err
		}
	}
	return &PublicOrderLookupResult{
		Order: s.buildOrderView(order),
	}, nil
}

func (s *Service) RefundOrder(ctx context.Context, userID uuid.UUID, orderID uuid.UUID) (*RefundOrderResult, error) {
	if !s.cfg.Enabled() {
		return nil, ErrPaymentsDisabled
	}
	order, err := s.repo.GetOrderByIDAndUser(ctx, orderID, userID)
	if err != nil {
		return nil, err
	}
	if !s.canRefundOrder(order) {
		return nil, ErrRefundNotAllowed
	}
	now := s.now().UTC()
	payload := map[string]any{
		"PaymentId": order.ProviderPaymentID,
		"Amount":    order.Amount,
	}
	if s.cfg.UsesProviderReceipt() && s.cfg.ReceiptSendOnCancel {
		if receipt := s.buildReceipt(order, order.UserEmail); receipt != nil {
			payload["Receipt"] = receipt
		}
	}
	var response tBankCancelResponse
	rawResponse, err := s.callTBank(ctx, "Cancel", payload, &response)
	order.ProviderResponse = rawResponse
	order.ProviderStatus = strings.ToUpper(strings.TrimSpace(response.Status))
	order.ProviderPaymentID = fallbackString(strings.TrimSpace(response.PaymentID), order.ProviderPaymentID)
	order.ProviderErrorCode = strings.TrimSpace(response.ErrorCode)
	order.ProviderMessage = fallbackString(strings.TrimSpace(response.Message), strings.TrimSpace(response.Details))
	order.LastCheckedAt = ptrTime(now)
	if err != nil {
		_ = s.repo.UpdateOrder(ctx, order)
		return nil, err
	}
	if updateErr := s.repo.UpdateOrder(ctx, order); updateErr != nil {
		return nil, updateErr
	}
	if !response.Success || order.ProviderErrorCode != "" && order.ProviderErrorCode != "0" {
		return nil, ErrProviderRejected
	}
	syncedOrder, syncErr := s.syncOrderState(ctx, order)
	if syncErr != nil {
		return nil, syncErr
	}
	return &RefundOrderResult{
		Order: s.buildOrderView(syncedOrder),
	}, nil
}

func (s *Service) verifyNotificationToken(payload map[string]any) bool {
	expected := strings.TrimSpace(fmt.Sprint(payload["Token"]))
	if expected == "" {
		return false
	}
	return strings.EqualFold(expected, hashTBankToken(payload, s.cfg.TBankPassword))
}

func (s *Service) HandleNotification(ctx context.Context, payload map[string]any, rawBody []byte) (*models.PaymentOrder, error) {
	if !s.cfg.Enabled() {
		return nil, ErrPaymentsDisabled
	}
	if !s.verifyNotificationToken(payload) {
		return nil, ErrNotificationToken
	}
	orderID, err := uuid.Parse(strings.TrimSpace(fmt.Sprint(payload["OrderId"])))
	if err != nil {
		return nil, ErrOrderNotFound
	}
	order, err := s.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	now := s.now().UTC()
	order.ProviderResponse = rawBody
	order.ProviderStatus = strings.ToUpper(strings.TrimSpace(fmt.Sprint(payload["Status"])))
	order.ProviderPaymentID = fallbackString(strings.TrimSpace(fmt.Sprint(payload["PaymentId"])), order.ProviderPaymentID)
	order.ProviderErrorCode = strings.TrimSpace(fmt.Sprint(payload["ErrorCode"]))
	order.ProviderMessage = fallbackString(strings.TrimSpace(fmt.Sprint(payload["Message"])), strings.TrimSpace(fmt.Sprint(payload["Details"])))
	order.NotifiedAt = ptrTime(now)
	order.LastCheckedAt = ptrTime(now)
	order.Status = internalStatusFromProvider(order.ProviderStatus, strings.EqualFold(fmt.Sprint(payload["Success"]), "true"))
	switch order.Status {
	case StatusSuccess:
		if order.PaidAt == nil {
			order.PaidAt = ptrTime(now)
		}
		if order.AccessFrom == nil || order.AccessUntil == nil {
			accessStart := now
			currentUntil, currentErr := s.repo.GetCurrentActiveAccessUntilByUser(ctx, order.UserID, now)
			if currentErr != nil {
				return nil, currentErr
			}
			if currentUntil != nil && currentUntil.After(accessStart) {
				accessStart = currentUntil.UTC()
			}
			accessUntil := accessStart.AddDate(0, s.cfg.AccessDurationMonths, 0)
			order.AccessFrom = ptrTime(accessStart)
			order.AccessUntil = ptrTime(accessUntil)
			// TODO(payments): when access logging is introduced, record activation IP, activation timestamp,
			// and subsequent paid-feature usage from this access-grant point.
		}
	case StatusRefunded:
		if order.RefundedAt == nil {
			order.RefundedAt = ptrTime(now)
		}
	case StatusFailed, StatusCanceled:
		if order.FailedAt == nil {
			order.FailedAt = ptrTime(now)
		}
	}
	if err := s.repo.UpdateOrder(ctx, order); err != nil {
		return nil, err
	}
	return order, nil
}

func (s *Service) AdminListOrders(ctx context.Context, search string, status string, planCode string, limit int, offset int) (*AdminOrdersListResult, error) {
	items, err := s.repo.ListOrdersForAdmin(ctx, strings.TrimSpace(search), strings.TrimSpace(status), strings.TrimSpace(planCode), limit, offset)
	if err != nil {
		return nil, err
	}
	total, err := s.repo.CountOrdersForAdmin(ctx, strings.TrimSpace(search), strings.TrimSpace(status), strings.TrimSpace(planCode))
	if err != nil {
		return nil, err
	}
	out := make([]OrderView, 0, len(items))
	for _, item := range items {
		copyItem := item
		out = append(out, s.buildOrderView(&copyItem))
	}
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}
	if offset < 0 {
		offset = 0
	}
	return &AdminOrdersListResult{
		Items:  out,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}

func (s *Service) AdminSummary(ctx context.Context) (*AdminOrdersSummary, error) {
	return s.repo.GetAdminSummary(ctx, s.now().UTC())
}
