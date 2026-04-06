package main

import (
	"log"
	"os"
	authsvc "sinde.ru/internal/auth"
	"sinde.ru/utils"
	"time"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("")
	utils.LoadEnv()
	authConfig, err := authsvc.LoadConfig()
	if err != nil {
		log.Fatalf("Auth config error: %v", err)
	}
	email := os.Getenv("MAIL_TEST_TO")
	if email == "" {
		email = authConfig.MailFrom
	}
	if email == "" {
		log.Fatal("missing MAIL_TEST_TO or AUTH_MAIL_FROM")
	}
	mailer := authsvc.NewMailer(authConfig)
	tests := []struct {
		purpose string
		ttl     time.Duration
		url     string
	}{
		{
			purpose: "verify_email",
			ttl:     authConfig.VerifyTTL,
			url:     authConfig.PublicBaseURL + "/auth/verify-email?token=test-verify-token",
		},
		{
			purpose: "password_reset",
			ttl:     authConfig.ResetTTL,
			url:     authConfig.PublicBaseURL + "/auth/reset-password?token=test-reset-token",
		},
		{
			purpose: "unknown",
			ttl:     time.Hour,
			url:     authConfig.PublicBaseURL + "/auth/action?token=test-action-token",
		},
	}
	for _, test := range tests {
		subject, textBody, htmlBody := authsvc.BuildActionEmailForTest(
			test.purpose,
			test.ttl,
			test.url,
			authConfig.PublicBaseURL,
		)
		if err := mailer.Send(email, subject, textBody, htmlBody); err != nil {
			log.Fatalf("send %s failed: %v", test.purpose, err)
		}
		log.Printf("ok: %s -> %s", test.purpose, email)
	}
}
