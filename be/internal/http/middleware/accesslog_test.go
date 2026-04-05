package middleware
import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"github.com/gofiber/fiber/v3"
	"sinde.ru/internal/http/responses"
)
func TestFormatResponseSize(t *testing.T) {
	cases := []struct {
		name  string
		bytes int
		want  string
	}{
		{name: "stream", bytes: -1, want: "stream"},
		{name: "zero", bytes: 0, want: ""},
		{name: "bytes", bytes: 512, want: "512 B"},
		{name: "kilobytes", bytes: 1536, want: "1,500 KB"},
	}
	for _, tc := range cases {
		if got := formatResponseSize(tc.bytes); got != tc.want {
			t.Fatalf("%s: expected %q, got %q", tc.name, tc.want, got)
		}
	}
}
func TestStringifyLogDetails(t *testing.T) {
	if got := stringifyLogDetails(nil); got != "" {
		t.Fatalf("expected empty string for nil, got %q", got)
	}
	if got := stringifyLogDetails("plain"); got != "plain" {
		t.Fatalf("unexpected string details: %q", got)
	}
	if got := stringifyLogDetails(errors.New("boom")); got != "boom" {
		t.Fatalf("unexpected error details: %q", got)
	}
	if got := stringifyLogDetails(map[string]any{"field": "email"}); got != `{"field":"email"}` {
		t.Fatalf("unexpected json details: %q", got)
	}
}
func TestColorEnabledRespectsEnvironment(t *testing.T) {
	t.Setenv("TERM", "xterm-256color")
	t.Setenv("NO_COLOR", "")
	if !colorEnabled() {
		t.Fatal("expected color to be enabled")
	}
	t.Setenv("NO_COLOR", "1")
	if colorEnabled() {
		t.Fatal("expected color to be disabled when NO_COLOR is set")
	}
	t.Setenv("NO_COLOR", "")
	t.Setenv("TERM", "dumb")
	if colorEnabled() {
		t.Fatal("expected color to be disabled for dumb terminal")
	}
}
func TestAccessLogIncludesResponseMessageAndDetails(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	app := fiber.New()
	app.Use(AccessLog())
	app.Get("/err", func(c fiber.Ctx) error {
		return responses.Error(c, fiber.StatusForbidden, "forbidden", fiber.Map{"field": "email"})
	})
	var buf bytes.Buffer
	oldWriter := log.Writer()
	log.SetOutput(&buf)
	t.Cleanup(func() {
		log.SetOutput(oldWriter)
	})
	resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/err", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	resp.Body.Close()
	line := buf.String()
	if !strings.Contains(line, "403 - GET /err") {
		t.Fatalf("log line does not include status and route: %q", line)
	}
	if !strings.Contains(line, "Доступ запрещён.") {
		t.Fatalf("log line does not include localized error message: %q", line)
	}
	if !strings.Contains(line, `({"field":"email"})`) {
		t.Fatalf("log line does not include serialized details: %q", line)
	}
}
