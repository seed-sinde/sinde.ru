package middleware
import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"github.com/gofiber/fiber/v3"
	"sinde.ru/internal/http/responses"
)
const (
	ansiReset         = "\033[0m"
	ansiRed           = "\033[31m"
	ansiYellow        = "\033[33m"
	ansiGreen         = "\033[32m"
	ansiCyan          = "\033[36m"
	ansiGray          = "\033[90m"
	ansiBrightMagenta = "\033[95m"
)
func AccessLog() fiber.Handler {
	return func(c fiber.Ctx) error {
		err := c.Next()
		status := c.Response().StatusCode()
		statusPart := fmt.Sprintf("%d", status)
		methodPart := c.Method()
		pathPart := c.Path()
		sizePart := formatResponseSize(responseSizeBytes(c))
		if colorEnabled() {
			statusPart = statusColor(status) + statusPart + ansiReset
			methodPart = ansiCyan + methodPart + ansiReset
			pathPart = ansiGray + pathPart + ansiReset
			if sizePart != "" {
				sizePart = ansiBrightMagenta + sizePart + ansiReset
			}
		}
		line := fmt.Sprintf("%s %s - %s %s", time.Now().Format("02.01 15:04:05"), statusPart, methodPart, pathPart)
		if sizePart != "" {
			line += " " + sizePart
		}
		message, _ := c.Locals(responses.ErrorMessageLocalKey).(string)
		if message != "" {
			line += " " + message
		}
		if details := stringifyLogDetails(c.Locals(responses.ErrorDetailsLocalKey)); details != "" {
			line += " (" + details + ")"
		}
		log.Println(line)
		return err
	}
}
func responseSizeBytes(c fiber.Ctx) int {
	if c.Response().IsBodyStream() {
		return -1
	}
	if n := c.Response().Header.ContentLength(); n >= 0 {
		return n
	}
	if bodyLen := len(c.Response().Body()); bodyLen > 0 {
		return bodyLen
	}
	return -1
}
func formatResponseSize(bytes int) string {
	if bytes < 0 {
		return "stream"
	}
	if bytes == 0 {
		return ""
	}
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	}
	kb := float64(bytes) / 1024.0
	text := fmt.Sprintf("%.3f", kb)
	text = strings.Replace(text, ".", ",", 1)
	return text + " KB"
}
func statusColor(status int) string {
	switch {
	case status >= 500:
		return ansiRed
	case status >= 400:
		return ansiYellow
	case status >= 200 && status < 300:
		return ansiGreen
	default:
		return ansiCyan
	}
}
func colorEnabled() bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	if os.Getenv("TERM") == "dumb" {
		return false
	}
	return true
}
func stringifyLogDetails(details any) string {
	if details == nil {
		return ""
	}
	switch value := details.(type) {
	case string:
		return value
	case error:
		return value.Error()
	default:
		encoded, err := json.Marshal(value)
		if err != nil {
			return fmt.Sprintf("%v", value)
		}
		return string(encoded)
	}
}
