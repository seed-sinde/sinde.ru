package auth

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"html"
	"log"
	"net"
	"net/mail"
	"net/smtp"
	"strconv"
	"strings"
	"time"
)

type Mailer interface {
	Send(to string, subject string, textBody string, htmlBody string) error
}

func NewMailer(cfg Config) Mailer {
	if strings.EqualFold(strings.TrimSpace(cfg.MailerDriver), "smtp") && strings.TrimSpace(cfg.SMTPHost) != "" {
		return SMTPMailer{
			Host:     cfg.SMTPHost,
			Port:     cfg.SMTPPort,
			Username: cfg.SMTPUsername,
			Password: cfg.SMTPPassword,
			From:     cfg.MailFrom,
			TLSMode:  cfg.SMTPTLSMode,
			HELO:     cfg.SMTPHELO,
		}
	}
	return LogMailer{}
}

type LogMailer struct{}

func (LogMailer) Send(to string, subject string, textBody string, htmlBody string) error {
	log.Printf("[auth-mail] получатель=%s тема=%q текст=%s html=%s", to, subject, textBody, htmlBody)
	return nil
}

type SMTPMailer struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	TLSMode  string
	HELO     string
}

const smtpDialTimeout = 10 * time.Second
const smtpSessionTimeout = 20 * time.Second

func (m SMTPMailer) Send(to string, subject string, textBody string, htmlBody string) error {
	from := strings.TrimSpace(m.From)
	if from == "" {
		return fmt.Errorf("smtp from is required")
	}
	parsedFrom, err := mail.ParseAddress(from)
	if err != nil {
		return fmt.Errorf("smtp from is invalid: %w", err)
	}
	if _, err := mail.ParseAddress(strings.TrimSpace(to)); err != nil {
		return fmt.Errorf("smtp recipient is invalid: %w", err)
	}
	addr := net.JoinHostPort(strings.TrimSpace(m.Host), strconv.Itoa(m.Port))
	message, err := buildSMTPMessage(parsedFrom, strings.TrimSpace(to), subject, textBody, htmlBody)
	if err != nil {
		return err
	}
	switch strings.ToLower(strings.TrimSpace(m.TLSMode)) {
	case "tls", "smtps":
		dialer := &net.Dialer{Timeout: smtpDialTimeout}
		conn, err := tls.DialWithDialer(dialer, "tcp", addr, &tls.Config{
			MinVersion: tls.VersionTLS12,
			ServerName: m.Host,
		})
		if err != nil {
			return err
		}
		_ = conn.SetDeadline(time.Now().Add(smtpSessionTimeout))
		client, err := smtp.NewClient(conn, m.Host)
		if err != nil {
			_ = conn.Close()
			return err
		}
		defer client.Close()
		if err := m.runSMTP(client, from, to, message, false); err != nil {
			return err
		}
		return client.Quit()
	default:
		conn, err := net.DialTimeout("tcp", addr, smtpDialTimeout)
		if err != nil {
			return err
		}
		_ = conn.SetDeadline(time.Now().Add(smtpSessionTimeout))
		client, err := smtp.NewClient(conn, m.Host)
		if err != nil {
			_ = conn.Close()
			return err
		}
		defer client.Close()
		startTLS := strings.ToLower(strings.TrimSpace(m.TLSMode)) != "plain"
		if err := m.runSMTP(client, from, to, message, startTLS); err != nil {
			return err
		}
		return client.Quit()
	}
}
func (m SMTPMailer) runSMTP(client *smtp.Client, from string, to string, message string, startTLS bool) error {
	if helo := strings.TrimSpace(m.HELO); helo != "" {
		if err := client.Hello(helo); err != nil {
			return err
		}
	}
	if startTLS {
		if ok, _ := client.Extension("STARTTLS"); ok {
			if err := client.StartTLS(&tls.Config{
				MinVersion: tls.VersionTLS12,
				ServerName: m.Host,
			}); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("smtp server does not support STARTTLS")
		}
	}
	if strings.TrimSpace(m.Username) != "" {
		auth := smtp.PlainAuth("", m.Username, m.Password, m.Host)
		if err := client.Auth(auth); err != nil {
			return err
		}
	}
	if err := client.Mail(from); err != nil {
		return err
	}
	if err := client.Rcpt(strings.TrimSpace(to)); err != nil {
		return err
	}
	writer, err := client.Data()
	if err != nil {
		return err
	}
	if _, err := writer.Write([]byte(message)); err != nil {
		_ = writer.Close()
		return err
	}
	return writer.Close()
}
func buildSMTPMessage(from *mail.Address, to string, subject string, textBody string, htmlBody string) (string, error) {
	fromHeader := (&mail.Address{Name: "sinde", Address: from.Address}).String()
	toHeader := to
	if parsedTo, err := mail.ParseAddress(to); err == nil {
		toHeader = parsedTo.String()
	}
	subjectHeader := encodeRFC2047(strings.TrimSpace(subject))
	textPart := base64UTF8(textBody)
	htmlPart := base64UTF8(htmlBody)
	boundary := fmt.Sprintf("sinde-alt-%d", time.Now().UnixNano())
	var buf bytes.Buffer
	headers := []string{
		"From: " + fromHeader,
		"To: " + toHeader,
		"Subject: " + subjectHeader,
		"MIME-Version: 1.0",
		"Content-Type: multipart/alternative; boundary=" + boundary,
		"",
	}
	buf.WriteString(strings.Join(headers, "\r\n"))
	buf.WriteString("\r\n--" + boundary + "\r\n")
	buf.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	buf.WriteString("Content-Transfer-Encoding: base64\r\n\r\n")
	buf.WriteString(textPart)
	buf.WriteString("\r\n--" + boundary + "\r\n")
	buf.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	buf.WriteString("Content-Transfer-Encoding: base64\r\n\r\n")
	buf.WriteString(htmlPart)
	buf.WriteString("\r\n--" + boundary + "--\r\n")
	return buf.String(), nil
}
func encodeRFC2047(value string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(strings.TrimSpace(value)))
	return "=?UTF-8?B?" + encoded + "?="
}
func base64UTF8(value string) string {
	raw := []byte(strings.ReplaceAll(value, "\r\n", "\n"))
	encoded := base64.StdEncoding.EncodeToString(raw)
	const lineLen = 76
	if len(encoded) <= lineLen {
		return encoded
	}
	var buf bytes.Buffer
	for start := 0; start < len(encoded); start += lineLen {
		end := start + lineLen
		if end > len(encoded) {
			end = len(encoded)
		}
		if start > 0 {
			buf.WriteString("\r\n")
		}
		buf.WriteString(encoded[start:end])
	}
	return buf.String()
}
func actionEmailHTML(projectHost string, title string, lead string, actionLabel string, actionURL string, expiry string, closing string) string {
	return fmt.Sprintf(`<!doctype html>
<html lang="ru">
  <body style="margin: 0; padding: 0; background: #0a0c10; font-family: Arial, sans-serif; color: #f6f2ea">
    <div style="margin: 0; padding: 32px 16px; background: #0a0c10">
      <div
        style="
          max-width: 640px;
          margin: 0 auto;
          background: #11141a;
          border: 1px solid #2a2f38;
          border-radius: 24px;
          overflow: hidden;
        ">
        <div style="padding: 28px 32px 24px; background: #0f1218; border-bottom: 1px solid #2a2f38">
          <div
            style="
              display: inline-block;
              padding: 8px 12px;
              border: 1px solid #6b4700;
              border-radius: 999px;
              background: #17120a;
              color: #e6b85c;
              font-size: 12px;
              line-height: 1;
              letter-spacing: 0.08em;
              font-weight: 700;
              text-transform: uppercase;
            ">
            %s • безопасность аккаунта
          </div>
          <h1 style="margin: 18px 0 0; font-size: 30px; line-height: 1.2; font-weight: 700; color: #f6f2ea">%s</h1>
        </div>
        <div style="padding: 32px">
          <p style="margin: 0 0 18px; font-size: 17px; line-height: 1.7; color: #f6f2ea">%s</p>
          <div
            style="
              margin: 0 0 24px;
              padding: 16px 18px;
              border: 1px solid #2a2f38;
              border-radius: 18px;
              background: #0d1015;
            ">
            <div
              style="
                margin: 0 0 6px;
                font-size: 12px;
                line-height: 1.4;
                letter-spacing: 0.08em;
                font-weight: 700;
                text-transform: uppercase;
                color: #89a9c7;
              ">
              срок действия
            </div>
            <div style="margin: 0; font-size: 15px; line-height: 1.6; color: #d8dde7">Ссылка действует %s.</div>
          </div>
          <p style="margin: 0 0 28px">
            <a
              href="%s"
              style="
                display: inline-block;
                padding: 14px 22px;
                border: 1px solid #8a5a08;
                border-radius: 999px;
                background: #131720;
                color: #f6f2ea;
                text-decoration: none;
                font-size: 15px;
                line-height: 1.2;
                font-weight: 700;
              ">
              %s
            </a>
          </p>
          <div
            style="
              margin: 0 0 24px;
              padding: 18px;
              border: 1px solid #2a2f38;
              border-radius: 18px;
              background: #0d1015;
            ">
            <p style="margin: 0 0 10px; font-size: 14px; line-height: 1.7; color: #c7d2e1">
              Если кнопка не открывается, используйте эту ссылку:
            </p>
            <p style="margin: 0; word-break: break-all; font-size: 14px; line-height: 1.7">
              <a href="%s" style="color: #afc3da; text-decoration: underline">%s</a>
            </p>
          </div>
          <p style="margin: 0; font-size: 14px; line-height: 1.8; color: #b6c2d2">%s</p>
        </div>
      </div>
    </div>
  </body>
</html>
`, html.EscapeString(projectHost),
		html.EscapeString(title),
		html.EscapeString(lead),
		html.EscapeString(expiry),
		html.EscapeString(actionURL),
		html.EscapeString(actionLabel),
		html.EscapeString(actionURL),
		html.EscapeString(actionURL),
		html.EscapeString(closing),
	)
}
