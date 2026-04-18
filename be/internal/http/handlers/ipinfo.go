package handlers

import (
	"github.com/gofiber/fiber/v3"
	"sinde.ru/db/services"
	"sinde.ru/internal/http/responses"
)

func GetIPInfo() fiber.Handler {
	return func(c fiber.Ctx) error {
		result, err := services.GetIPInfo(c.Context(), c.IP())
		if err != nil {
			return responses.Error(c, fiber.StatusInternalServerError, "failed to get ip info")
		}
		return responses.Success(c, fiber.StatusOK, result)
	}
}
