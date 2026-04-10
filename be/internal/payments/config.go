package payments

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	TBankAPIURL          string
	TBankTerminalKey     string
	TBankPassword        string
	TaxMode              string
	ReceiptEnabled       bool
	ReceiptTaxation      string
	ReceiptTax           string
	ReceiptPaymentMethod string
	ReceiptPaymentObject string
	ReceiptProTitle      string
	ReceiptDonationTitle string
	ReceiptSendOnCancel  bool
	PublicBaseURL        string
	Currency             string
	FeePercent           float64
	RequestTimeout       time.Duration
	NotificationPath     string
	SuccessPath          string
	FailPath             string
	AccessDurationMonths int
}

func LoadConfig() (Config, error) {
	get := func(key string) string {
		return strings.TrimSpace(os.Getenv(key))
	}
	cfg := Config{
		TBankAPIURL:          fallbackString(get("TBANK_API_URL"), "https://rest-api-test.tinkoff.ru/v2"),
		TBankTerminalKey:     get("TBANK_TERMINAL_KEY"),
		TBankPassword:        get("TBANK_PASSWORD"),
		TaxMode:              strings.ToLower(fallbackString(get("TBANK_TAX_MODE"), "generic")),
		ReceiptEnabled:       strings.EqualFold(get("TBANK_RECEIPT_ENABLED"), "true"),
		ReceiptTaxation:      get("TBANK_RECEIPT_TAXATION"),
		ReceiptTax:           fallbackString(get("TBANK_RECEIPT_TAX"), "none"),
		ReceiptPaymentMethod: fallbackString(get("TBANK_RECEIPT_PAYMENT_METHOD"), "full_payment"),
		ReceiptPaymentObject: fallbackString(get("TBANK_RECEIPT_PAYMENT_OBJECT"), "service"),
		ReceiptProTitle:      fallbackString(get("TBANK_RECEIPT_ITEM_NAME_PRO"), "Доступ sinde.ru Pro"),
		ReceiptDonationTitle: fallbackString(get("TBANK_RECEIPT_ITEM_NAME_DONATION"), "Поддержка проекта sinde.ru"),
		ReceiptSendOnCancel:  strings.EqualFold(get("TBANK_RECEIPT_SEND_ON_CANCEL"), "true"),
		PublicBaseURL:        fallbackString(get("PAYMENTS_PUBLIC_BASE_URL"), get("AUTH_PUBLIC_BASE_URL")),
		Currency:             fallbackString(get("TBANK_CURRENCY"), "RUB"),
		FeePercent:           0,
		RequestTimeout:       15 * time.Second,
		NotificationPath:     "/api/v1/payments/tbank/notify",
		SuccessPath:          "/payments/success",
		FailPath:             "/payments/fail",
		AccessDurationMonths: 1,
	}
	if raw := get("TBANK_FEE_PERCENT"); raw != "" {
		value, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return Config{}, fmt.Errorf("TBANK_FEE_PERCENT must be a number: %w", err)
		}
		cfg.FeePercent = value
	}
	if raw := get("TBANK_REQUEST_TIMEOUT_SECONDS"); raw != "" {
		value, err := strconv.Atoi(raw)
		if err != nil || value <= 0 {
			return Config{}, fmt.Errorf("TBANK_REQUEST_TIMEOUT_SECONDS must be a positive integer")
		}
		cfg.RequestTimeout = time.Duration(value) * time.Second
	}
	if raw := get("PAYMENTS_ACCESS_DURATION_MONTHS"); raw != "" {
		value, err := strconv.Atoi(raw)
		if err != nil || value <= 0 {
			return Config{}, fmt.Errorf("PAYMENTS_ACCESS_DURATION_MONTHS must be a positive integer")
		}
		cfg.AccessDurationMonths = value
	}
	if cfg.TBankAPIURL == "" {
		return Config{}, fmt.Errorf("TBANK_API_URL is required")
	}
	if cfg.UsesProviderReceipt() {
		if cfg.ReceiptTaxation == "" {
			return Config{}, fmt.Errorf("TBANK_RECEIPT_TAXATION is required when TBANK_RECEIPT_ENABLED=true")
		}
		if cfg.ReceiptTax == "" {
			return Config{}, fmt.Errorf("TBANK_RECEIPT_TAX is required when TBANK_RECEIPT_ENABLED=true")
		}
		if cfg.ReceiptPaymentMethod == "" {
			return Config{}, fmt.Errorf("TBANK_RECEIPT_PAYMENT_METHOD is required when TBANK_RECEIPT_ENABLED=true")
		}
		if cfg.ReceiptPaymentObject == "" {
			return Config{}, fmt.Errorf("TBANK_RECEIPT_PAYMENT_OBJECT is required when TBANK_RECEIPT_ENABLED=true")
		}
	}
	return cfg, nil
}

func (c Config) Enabled() bool {
	return c.TBankTerminalKey != "" && c.TBankPassword != ""
}

func (c Config) UsesProviderReceipt() bool {
	return c.ReceiptEnabled && c.TaxMode != "npd"
}
