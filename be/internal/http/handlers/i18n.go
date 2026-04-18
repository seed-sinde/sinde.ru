package handlers

import (
	"errors"
	"strings"

	"sinde.ru/internal/http/responses"
	"sinde.ru/internal/store"

	"github.com/gofiber/fiber/v3"
)

func GetI18nHandler(i18nStore *store.I18nStore) fiber.Handler {
	return func(c fiber.Ctx) error {
		locale := store.LocaleCode(strings.TrimSpace(c.Params("locale")))
		namespace := strings.TrimSpace(c.Params("namespace"))

		if locale == "" || namespace == "" {
			return responses.Error(c, 400, "invalid i18n params")
		}

		switch locale {
		case store.LocaleEN, store.LocaleRU, store.LocaleZH, store.LocaleJA:
		default:
			return responses.Error(c, 400, "unsupported locale", locale)
		}

		data, err := i18nStore.ResolveSection(c.Context(), locale, namespace)
		if err != nil {
			if errors.Is(err, store.ErrI18nNamespaceNotFound) {
				return responses.Error(c, 404, "i18n namespace not found", namespace)
			}
			return responses.Error(c, 500, "failed to load i18n section")
		}

		c.Set("Cache-Control", "public, max-age=300")

		return responses.Success(c, 200, data)
	}
}
