package handlers

import (
	"github.com/gofiber/fiber/v3"

	"sinde.ru/db/services"
	"sinde.ru/internal/http/responses"
)

func ChemistryElementsListHandler() fiber.Handler {
	return func(c fiber.Ctx) error {
		items, err := services.PdbListChemistryElements(c.Context())
		if err != nil {
			return responses.Error(c, fiber.StatusInternalServerError, "Не удалось загрузить элементы", err.Error())
		}
		return responses.Success(c, fiber.StatusOK, items)
	}
}
