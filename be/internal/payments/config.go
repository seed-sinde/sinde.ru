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
	return cfg, nil
}

func (c Config) Enabled() bool {
	return c.TBankTerminalKey != "" && c.TBankPassword != ""
}
