package responses

import (
	"encoding/json"
	"github.com/gofiber/fiber/v3"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestErrorLocalizesMessageAndStoresLocals(t *testing.T) {
	app := fiber.New()
	var gotMessage string
	var gotDetails any
	app.Get("/err", func(c fiber.Ctx) error {
		details := fiber.Map{"field": "email"}
		err := Error(c, fiber.StatusForbidden, "forbidden", details)
		gotMessage, _ = c.Locals(ErrorMessageLocalKey).(string)
		gotDetails = c.Locals(ErrorDetailsLocalKey)
		return err
	})
	resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/err", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusForbidden {
		t.Fatalf("unexpected status: %d", resp.StatusCode)
	}
	var body map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if ok, _ := body["ok"].(bool); ok {
		t.Fatalf("expected ok=false, got %#v", body["ok"])
	}
	if body["message"] != "Доступ запрещён." {
		t.Fatalf("unexpected localized message: %#v", body["message"])
	}
	wantDetails := fiber.Map{"field": "email"}
	if !reflect.DeepEqual(body["details"], map[string]any{"field": "email"}) {
		t.Fatalf("unexpected response details: %#v", body["details"])
	}
	if gotMessage != "Доступ запрещён." {
		t.Fatalf("unexpected locals message: %q", gotMessage)
	}
	if !reflect.DeepEqual(gotDetails, wantDetails) {
		t.Fatalf("unexpected locals details: %#v", gotDetails)
	}
}
func TestErrorRejectsNonErrorStatusCode(t *testing.T) {
	app := fiber.New()
	app.Get("/err", func(c fiber.Ctx) error {
		return Error(c, fiber.StatusOK, "forbidden")
	})
	resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/err", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusInternalServerError {
		t.Fatalf("unexpected status: %d", resp.StatusCode)
	}
	var body map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if body["message"] != "Некорректный код статуса для ошибки (должен быть 4xx или 5xx)" {
		t.Fatalf("unexpected validation message: %#v", body["message"])
	}
}
func TestSuccessReturnsPayloadAndRejectsInvalidStatus(t *testing.T) {
	app := fiber.New()
	app.Get("/ok", func(c fiber.Ctx) error {
		return Success(c, fiber.StatusCreated, fiber.Map{"id": "42"})
	})
	app.Get("/bad", func(c fiber.Ctx) error {
		return Success(c, fiber.StatusBadRequest, fiber.Map{"id": "42"})
	})
	resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/ok", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusCreated {
		t.Fatalf("unexpected status: %d", resp.StatusCode)
	}
	var okBody map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&okBody); err != nil {
		t.Fatalf("failed to decode success response: %v", err)
	}
	if ok, _ := okBody["ok"].(bool); !ok {
		t.Fatalf("expected ok=true, got %#v", okBody["ok"])
	}
	if !reflect.DeepEqual(okBody["data"], map[string]any{"id": "42"}) {
		t.Fatalf("unexpected payload: %#v", okBody["data"])
	}
	badResp, err := app.Test(httptest.NewRequest(http.MethodGet, "/bad", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer badResp.Body.Close()
	if badResp.StatusCode != fiber.StatusInternalServerError {
		t.Fatalf("unexpected status: %d", badResp.StatusCode)
	}
}
