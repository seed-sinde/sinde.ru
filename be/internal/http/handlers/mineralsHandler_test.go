package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
)

func TestParseMineralsListQueryImageFilter(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		statusCode int
	}{
		{name: "default any", query: "/", statusCode: fiber.StatusNoContent},
		{name: "with", query: "/?image=with", statusCode: fiber.StatusNoContent},
		{name: "without", query: "/?image=without", statusCode: fiber.StatusNoContent},
		{name: "reject any", query: "/?image=any", statusCode: fiber.StatusBadRequest},
		{name: "reject garbage", query: "/?image=garbage", statusCode: fiber.StatusBadRequest},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			app := fiber.New()
			app.Get("/", func(c fiber.Ctx) error {
				_, err := parseMineralsListQuery(c)
				if err != nil {
					return err
				}
				return c.SendStatus(fiber.StatusNoContent)
			})

			req := httptest.NewRequest(http.MethodGet, test.query, nil)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("app.Test() error = %v", err)
			}
			if resp.StatusCode != test.statusCode {
				t.Fatalf("status = %d, want %d", resp.StatusCode, test.statusCode)
			}
		})
	}
}
