package handlers
import (
	"strings"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"sinde.ru/internal/http/responses"
	"sinde.ru/internal/media"
)
func MediaUploadHandler() fiber.Handler {
	return func(c fiber.Ctx) error {
		user, authErr := currentKitchenUser(c)
		if authErr != nil {
			return authErr
		}
		fileHeader, err := c.FormFile("file")
		if err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "файл изображения обязателен")
		}
		section := media.Section(strings.ToLower(strings.TrimSpace(c.FormValue("section"))))
		if section == "" {
			return responses.Error(c, fiber.StatusBadRequest, "секция media обязательна")
		}
		collection := media.Collection(strings.ToLower(strings.TrimSpace(c.FormValue("collection"))))
		if collection == "" {
			return responses.Error(c, fiber.StatusBadRequest, "коллекция media обязательна")
		}
		target := media.UploadTarget{
			Section:    section,
			Collection: collection,
			UserID:     user.UserID,
		}
		switch {
		case section == media.SectionUsers && collection == media.CollectionAvatars:
		case section == media.SectionKitchen && collection == media.CollectionRecipes:
			recipeIDRaw := strings.TrimSpace(c.FormValue("recipe_id"))
			if recipeIDRaw != "" {
				recipeID, parseErr := uuid.Parse(recipeIDRaw)
				if parseErr != nil {
					return responses.Error(c, fiber.StatusBadRequest, "ID рецепта некорректен")
				}
				if err := ensureRecipeUploadAllowed(c, user, recipeID); err != nil {
					return err
				}
				target.RecipeID = &recipeID
			}
		default:
			return responses.Error(c, fiber.StatusBadRequest, "цель media некорректна")
		}
		stored, err := media.SaveUploadedFile(c.Context(), fileHeader, target)
		if err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "не удалось сохранить изображение", err.Error())
		}
		return responses.Success(c, fiber.StatusCreated, fiber.Map{
			"image_key":  stored.ImageKey,
			"image_url":  stored.ImageURL,
			"reused":     stored.Reused,
			"file_hash":  stored.FileHash,
			"section":    string(target.Section),
			"collection": string(target.Collection),
		})
	}
}
func MediaGetFileHandler() fiber.Handler {
	return func(c fiber.Ctx) error {
		storageKey := media.NormalizeStorageKey(c.Params("*"))
		if storageKey == "" {
			return responses.Error(c, fiber.StatusNotFound, "изображение не найдено")
		}
		object, err := media.OpenObject(c.Context(), storageKey)
		if err != nil {
			return responses.Error(c, fiber.StatusNotFound, "изображение не найдено")
		}
		c.Set("Cache-Control", "public, max-age=31536000, immutable")
		if object.ContentType != "" {
			c.Set("Content-Type", object.ContentType)
		}
		if contentLength := media.ContentLengthHeader(object.Size); contentLength != "" {
			c.Set("Content-Length", contentLength)
			return c.SendStream(object.Body, int(object.Size))
		}
		return c.SendStream(object.Body)
	}
}
