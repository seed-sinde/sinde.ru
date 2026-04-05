package handlers
import (
	"context"
	"errors"
	"strconv"
	"strings"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"sinde.ru/db/services"
	"sinde.ru/internal/http/middleware"
	"sinde.ru/internal/http/responses"
	"sinde.ru/internal/media"
	"sinde.ru/internal/models"
)
type kitchenRecipeIngredientInput struct {
	Name   string `json:"name"`
	Amount string `json:"amount"`
	Unit   string `json:"unit"`
	Note   string `json:"note"`
}
type kitchenRecipeStepInput struct {
	Text     string `json:"text"`
	ImageKey string `json:"image_key"`
}
type kitchenRecipeCreateRequest struct {
	RecipeID      string                         `json:"recipe_id"`
	Title         string                         `json:"title"`
	Description   string                         `json:"description"`
	CoverImageKey string                         `json:"cover_image_key"`
	Kcal          *int                           `json:"kcal"`
	PrepMinutes   int                            `json:"prep_minutes"`
	CookMinutes   int                            `json:"cook_minutes"`
	Servings      int                            `json:"servings"`
	Difficulty    string                         `json:"difficulty"`
	MealType      string                         `json:"meal_type"`
	CookingMethod string                         `json:"cooking_method"`
	Cuisine       string                         `json:"cuisine"`
	DietType      string                         `json:"diet_type"`
	Ingredients   []kitchenRecipeIngredientInput `json:"ingredients"`
	Steps         []kitchenRecipeStepInput       `json:"steps"`
	Tags          []string                       `json:"tags"`
	IsPublic      *bool                          `json:"is_public"`
}
type kitchenRecipeSearchRequest struct {
	Query                 string   `json:"query"`
	Ingredients           []string `json:"ingredients"`
	ExcludedIngredients   []string `json:"excludedIngredients"`
	RequireAllIngredients bool     `json:"requireAllIngredients"`
	Difficulty            string   `json:"difficulty"`
	MealType              string   `json:"mealType"`
	CookingMethod         string   `json:"cookingMethod"`
	Cuisine               string   `json:"cuisine"`
	DietType              string   `json:"dietType"`
	KcalMin               *int     `json:"kcalMin"`
	KcalMax               *int     `json:"kcalMax"`
	MaxTotalMinutes       *int     `json:"maxTotalMinutes"`
	ServingsMin           *int     `json:"servingsMin"`
	ServingsMax           *int     `json:"servingsMax"`
	Limit                 int      `json:"limit"`
	Offset                int      `json:"offset"`
}
type kitchenIngredientMutationRequest struct {
	IngredientID int64  `json:"ingredient_id"`
	Name         string `json:"name"`
	Category     string `json:"category"`
	ListType     string `json:"list_type"`
}
type kitchenRecipeModerationRequest struct {
	Approve bool   `json:"approve"`
	Note    string `json:"note"`
}
type kitchenRecipeOwnerChangeRequest struct {
	OwnerUserID string `json:"owner_user_id"`
}
func normalizeKitchenStepImageKey(raw string) string {
	return media.NormalizeStorageKey(raw)
}
func kitchenFavoriteListType(raw string) string {
	value := strings.ToLower(strings.TrimSpace(raw))
	if value == "exclude" {
		return "exclude"
	}
	return "include"
}
func normalizeKitchenToken(v string) string {
	return strings.ToLower(strings.TrimSpace(v))
}
func normalizeKitchenDietType(v string) string {
	return normalizeKitchenToken(v)
}
func normalizeKitchenCookingMethod(v string) string {
	return normalizeKitchenToken(v)
}
func uniqueNormalized(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	res := make([]string, 0, len(values))
	for _, item := range values {
		n := normalizeKitchenToken(item)
		if n == "" {
			continue
		}
		if _, ok := seen[n]; ok {
			continue
		}
		seen[n] = struct{}{}
		res = append(res, n)
	}
	return res
}
func kitchenIntOrZero(v *int) int {
	if v == nil {
		return 0
	}
	return *v
}
func KitchenIngredientsHandler() fiber.Handler {
	return func(c fiber.Ctx) error {
		items, err := services.PdbKitchenListCategories(c)
		if err != nil {
			return responses.Error(c, fiber.StatusInternalServerError, "не удалось загрузить категории ингредиентов", err.Error())
		}
		categories := make([]string, 0, len(items))
		for _, item := range items {
			if item == nil {
				continue
			}
			categories = append(categories, item.Label)
		}
		return responses.Success(c, fiber.StatusOK, fiber.Map{
			"categories": categories,
		})
	}
}
func KitchenCatalogHandler() fiber.Handler {
	return func(c fiber.Ctx) error {
		payload, err := services.PdbKitchenCatalog(c)
		if err != nil {
			return responses.Error(c, fiber.StatusInternalServerError, "не удалось загрузить каталог кухни", err.Error())
		}
		return responses.Success(c, fiber.StatusOK, payload)
	}
}
func currentKitchenUser(c fiber.Ctx) (*models.User, error) {
	user := middleware.CurrentUser(c)
	if user == nil {
		return nil, responses.Error(c, fiber.StatusUnauthorized, "требуется аутентификация")
	}
	return user, nil
}
func ensureRecipeUploadAllowed(c fiber.Ctx, user *models.User, recipeID uuid.UUID) error {
	recipe, err := services.PdbKitchenGetRecipeByIDAny(c.Context(), recipeID)
	if err != nil {
		if errors.Is(err, services.ErrKitchenRecipeNotFound) {
			return nil
		}
		return responses.Error(c, fiber.StatusInternalServerError, "не удалось проверить загрузку рецепта", err.Error())
	}
	if kitchenUserHasRole(user, "admin") {
		return nil
	}
	if recipe.OwnerUserID != nil && *recipe.OwnerUserID == user.UserID {
		return nil
	}
	return responses.Error(c, fiber.StatusForbidden, "можно загружать изображения только для своих рецептов")
}
func kitchenRecipeImageKeys(recipe *models.KitchenRecipe) []string {
	if recipe == nil {
		return nil
	}
	keysMap := make(map[string]struct{})
	addKey := func(raw string) {
		key := media.NormalizeStorageKey(raw)
		if key == "" {
			return
		}
		keysMap[key] = struct{}{}
	}
	addKey(recipe.CoverImageKey)
	for _, step := range recipe.Steps {
		addKey(step.ImageKey)
	}
	keys := make([]string, 0, len(keysMap))
	for key := range keysMap {
		keys = append(keys, key)
	}
	return keys
}
func cleanupRemovedRecipeImages(ctx context.Context, before *models.KitchenRecipe, after *models.KitchenRecipe) error {
	beforeKeys := kitchenRecipeImageKeys(before)
	if len(beforeKeys) == 0 {
		return nil
	}
	afterSet := make(map[string]struct{})
	for _, key := range kitchenRecipeImageKeys(after) {
		afterSet[key] = struct{}{}
	}
	for _, key := range beforeKeys {
		if _, keep := afterSet[key]; keep {
			continue
		}
		if err := media.DeleteIfUnreferenced(ctx, key); err != nil {
			return err
		}
	}
	return nil
}
func kitchenUserHasRole(user *models.User, role string) bool {
	if user == nil {
		return false
	}
	target := strings.ToLower(strings.TrimSpace(role))
	if target == "" {
		return false
	}
	for _, item := range user.Roles {
		if strings.ToLower(strings.TrimSpace(item)) == target {
			return true
		}
	}
	return false
}
func buildKitchenRecipeFromRequest(req kitchenRecipeCreateRequest, ownerUserID *uuid.UUID, isAdmin bool) (*models.KitchenRecipe, error) {
	recipeID := uuid.New()
	if raw := strings.TrimSpace(req.RecipeID); raw != "" {
		parsed, err := uuid.Parse(raw)
		if err != nil {
			return nil, errors.New("ID рецепта некорректен")
		}
		recipeID = parsed
	}
	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, errors.New("title is required")
	}
	if req.PrepMinutes < 0 || req.CookMinutes < 0 {
		return nil, errors.New("time values must be non-negative")
	}
	kcal := 0
	if req.Kcal != nil {
		if *req.Kcal < 0 {
			return nil, errors.New("kcal must be non-negative")
		}
		kcal = *req.Kcal
	}
	ingredients := make([]models.KitchenIngredient, 0, len(req.Ingredients))
	ingredientsSearchRaw := make([]string, 0, len(req.Ingredients))
	for _, item := range req.Ingredients {
		name := strings.TrimSpace(item.Name)
		if name == "" {
			continue
		}
		ingredients = append(ingredients, models.KitchenIngredient{
			Name:   name,
			Amount: strings.TrimSpace(item.Amount),
			Unit:   strings.TrimSpace(item.Unit),
			Note:   strings.TrimSpace(item.Note),
		})
		ingredientsSearchRaw = append(ingredientsSearchRaw, name)
	}
	if len(ingredients) == 0 {
		return nil, errors.New("at least one ingredient is required")
	}
	steps := make([]models.KitchenStep, 0, len(req.Steps))
	for _, step := range req.Steps {
		text := strings.TrimSpace(step.Text)
		if text == "" {
			continue
		}
		imageKey := ""
		if step.ImageKey != "" {
			imageKey = normalizeKitchenStepImageKey(step.ImageKey)
			if imageKey == "" {
				return nil, errors.New("step image key is invalid")
			}
		}
		steps = append(steps, models.KitchenStep{
			Order:    len(steps) + 1,
			Text:     text,
			ImageKey: imageKey,
		})
	}
	if len(steps) == 0 {
		return nil, errors.New("at least one step is required")
	}
	servings := req.Servings
	if servings <= 0 {
		servings = 1
	}
	difficulty := strings.TrimSpace(req.Difficulty)
	if difficulty == "" {
		difficulty = "easy"
	}
	mealType := strings.TrimSpace(req.MealType)
	if mealType == "" {
		mealType = "other"
	}
	cookingMethod := normalizeKitchenCookingMethod(req.CookingMethod)
	dietType := normalizeKitchenDietType(req.DietType)
	coverImageKey := ""
	if req.CoverImageKey != "" {
		coverImageKey = normalizeKitchenStepImageKey(req.CoverImageKey)
		if coverImageKey == "" {
			return nil, errors.New("cover image key is invalid")
		}
	}
	isPublic := true
	if req.IsPublic != nil {
		isPublic = *req.IsPublic
	}
	moderationStatus := "approved"
	if isAdmin {
		if !isPublic {
			moderationStatus = "draft"
		}
	} else {
		isPublic = false
		moderationStatus = "pending"
	}
	return &models.KitchenRecipe{
		ID:                recipeID,
		OwnerUserID:       ownerUserID,
		Title:             title,
		Description:       strings.TrimSpace(req.Description),
		CoverImageKey:     coverImageKey,
		Kcal:              kcal,
		PrepMinutes:       req.PrepMinutes,
		CookMinutes:       req.CookMinutes,
		Servings:          servings,
		Difficulty:        difficulty,
		MealType:          mealType,
		CookingMethod:     cookingMethod,
		Cuisine:           strings.TrimSpace(req.Cuisine),
		DietType:          dietType,
		Ingredients:       ingredients,
		IngredientsSearch: uniqueNormalized(ingredientsSearchRaw),
		Steps:             steps,
		Tags:              uniqueNormalized(req.Tags),
		IsPublic:          isPublic,
		ModerationStatus:  moderationStatus,
		ModeratedByUserID: nil,
		ModeratedAt:       nil,
		ModerationNote:    "",
	}, nil
}
func KitchenCreateRecipeHandler() fiber.Handler {
	return func(c fiber.Ctx) error {
		user, authErr := currentKitchenUser(c)
		if authErr != nil {
			return authErr
		}
		var req kitchenRecipeCreateRequest
		if err := c.Bind().Body(&req); err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "некорректный запрос", err.Error())
		}
		recipe, err := buildKitchenRecipeFromRequest(req, &user.UserID, kitchenUserHasRole(user, "admin"))
		if err != nil {
			return responses.Error(c, fiber.StatusBadRequest, err.Error())
		}
		created, err := services.PdbKitchenCreateRecipe(c, recipe)
		if err != nil {
			return responses.Error(c, fiber.StatusInternalServerError, "не удалось создать рецепт", err.Error())
		}
		return responses.Success(c, fiber.StatusCreated, created)
	}
}
func KitchenUpdateRecipeHandler() fiber.Handler {
	return func(c fiber.Ctx) error {
		user, authErr := currentKitchenUser(c)
		if authErr != nil {
			return authErr
		}
		recipeID, err := uuid.Parse(strings.TrimSpace(c.Params("id")))
		if err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "ID рецепта некорректен")
		}
		var req kitchenRecipeCreateRequest
		if err := c.Bind().Body(&req); err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "некорректный запрос", err.Error())
		}
		isAdmin := kitchenUserHasRole(user, "admin")
		existing, err := services.PdbKitchenGetRecipeByIDAny(c, recipeID)
		if err != nil {
			if errors.Is(err, services.ErrKitchenRecipeNotFound) {
				return responses.Error(c, fiber.StatusNotFound, "рецепт не найден")
			}
			return responses.Error(c, fiber.StatusInternalServerError, "не удалось загрузить рецепт", err.Error())
		}
		if !isAdmin && (existing.OwnerUserID == nil || *existing.OwnerUserID != user.UserID) {
			return responses.Error(c, fiber.StatusForbidden, "можно редактировать только свои рецепты")
		}
		recipe, err := buildKitchenRecipeFromRequest(req, &user.UserID, isAdmin)
		if err != nil {
			return responses.Error(c, fiber.StatusBadRequest, err.Error())
		}
		recipe.ID = recipeID
		if !isAdmin {
			recipe.IsPublic = false
			recipe.ModerationStatus = "pending"
			recipe.ModeratedByUserID = nil
			recipe.ModeratedAt = nil
			recipe.ModerationNote = ""
		}
		var updated *models.KitchenRecipe
		if isAdmin {
			updated, err = services.PdbKitchenUpdateRecipeAsAdmin(c, recipe)
		} else {
			updated, err = services.PdbKitchenUpdateRecipe(c, recipe, user.UserID)
		}
		if err != nil {
			if errors.Is(err, services.ErrKitchenRecipeForbidden) {
				return responses.Error(c, fiber.StatusForbidden, "можно редактировать только свои рецепты")
			}
			if errors.Is(err, services.ErrKitchenRecipeNotFound) {
				return responses.Error(c, fiber.StatusNotFound, "рецепт не найден")
			}
			return responses.Error(c, fiber.StatusInternalServerError, "не удалось обновить рецепт", err.Error())
		}
		if err := cleanupRemovedRecipeImages(c.Context(), existing, updated); err != nil {
			return responses.Error(c, fiber.StatusInternalServerError, "рецепт обновлен, но не удалось удалить удаленные изображения", err.Error())
		}
		return responses.Success(c, fiber.StatusOK, updated)
	}
}
func KitchenDeleteRecipeHandler() fiber.Handler {
	return func(c fiber.Ctx) error {
		user, authErr := currentKitchenUser(c)
		if authErr != nil {
			return authErr
		}
		recipeID, err := uuid.Parse(strings.TrimSpace(c.Params("id")))
		if err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "ID рецепта некорректен")
		}
		existing, err := services.PdbKitchenGetRecipeByIDAny(c, recipeID)
		if err != nil {
			if errors.Is(err, services.ErrKitchenRecipeNotFound) {
				return responses.Error(c, fiber.StatusNotFound, "рецепт не найден")
			}
			return responses.Error(c, fiber.StatusInternalServerError, "не удалось загрузить рецепт", err.Error())
		}
		var deleteErr error
		if kitchenUserHasRole(user, "admin") {
			deleteErr = services.PdbKitchenDeleteRecipeAsAdmin(c, recipeID)
		} else {
			deleteErr = services.PdbKitchenDeleteRecipe(c, recipeID, user.UserID)
		}
		if deleteErr != nil {
			if errors.Is(deleteErr, services.ErrKitchenRecipeForbidden) {
				return responses.Error(c, fiber.StatusForbidden, "можно удалять только свои рецепты")
			}
			if errors.Is(deleteErr, services.ErrKitchenRecipeNotFound) {
				return responses.Error(c, fiber.StatusNotFound, "рецепт не найден")
			}
			return responses.Error(c, fiber.StatusInternalServerError, "не удалось удалить рецепт", deleteErr.Error())
		}
		if err := cleanupRemovedRecipeImages(c.Context(), existing, nil); err != nil {
			return responses.Error(c, fiber.StatusInternalServerError, "рецепт удален, но не удалось очистить изображения", err.Error())
		}
		return responses.Success(c, fiber.StatusOK, fiber.Map{"deleted": true})
	}
}
func KitchenMyRecipesHandler() fiber.Handler {
	return func(c fiber.Ctx) error {
		user, authErr := currentKitchenUser(c)
		if authErr != nil {
			return authErr
		}
		items, err := services.PdbKitchenListRecipesByOwner(c, user.UserID, 100)
		if err != nil {
			return responses.Error(c, fiber.StatusInternalServerError, "не удалось загрузить ваши рецепты", err.Error())
		}
		return responses.Success(c, fiber.StatusOK, fiber.Map{"items": items})
	}
}
func KitchenAdminModerationRecipesHandler() fiber.Handler {
	return func(c fiber.Ctx) error {
		user, authErr := currentKitchenUser(c)
		if authErr != nil {
			return authErr
		}
		if !kitchenUserHasRole(user, "admin") {
			return responses.Error(c, fiber.StatusForbidden, "требуется роль администратора")
		}
		limit := 30
		if raw := strings.TrimSpace(c.Query("limit")); raw != "" {
			if n, err := strconv.Atoi(raw); err == nil {
				limit = n
			}
		}
		offset := 0
		if raw := strings.TrimSpace(c.Query("offset")); raw != "" {
			if n, err := strconv.Atoi(raw); err == nil {
				offset = n
			}
		}
		status := services.NormalizeKitchenModerationStatus(c.Query("status"))
		if status == "" {
			status = "pending"
		}
		items, total, err := services.PdbKitchenListRecipesForModeration(c, status, limit, offset)
		if err != nil {
			return responses.Error(c, fiber.StatusInternalServerError, "не удалось загрузить рецепты на модерации", err.Error())
		}
		statusTotals, err := services.PdbKitchenModerationStatusTotals(c)
		if err != nil {
			return responses.Error(c, fiber.StatusInternalServerError, "не удалось загрузить рецепты на модерации", err.Error())
		}
		return responses.Success(c, fiber.StatusOK, fiber.Map{
			"items":         items,
			"total":         total,
			"limit":         limit,
			"offset":        offset,
			"status":        status,
			"status_totals": statusTotals,
		})
	}
}
func KitchenAdminModerateRecipeHandler() fiber.Handler {
	return func(c fiber.Ctx) error {
		user, authErr := currentKitchenUser(c)
		if authErr != nil {
			return authErr
		}
		if !kitchenUserHasRole(user, "admin") {
			return responses.Error(c, fiber.StatusForbidden, "требуется роль администратора")
		}
		recipeID, err := uuid.Parse(strings.TrimSpace(c.Params("id")))
		if err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "ID рецепта некорректен")
		}
		var req kitchenRecipeModerationRequest
		if err := c.Bind().Body(&req); err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "некорректный запрос", err.Error())
		}
		req.Note = strings.TrimSpace(req.Note)
		if !req.Approve && req.Note == "" {
			return responses.Error(c, fiber.StatusBadRequest, "заметка модератора обязательна")
		}
		item, err := services.PdbKitchenModerateRecipe(c, recipeID, req.Approve, user.UserID, req.Note)
		if err != nil {
			if errors.Is(err, services.ErrKitchenRecipeNotFound) {
				return responses.Error(c, fiber.StatusNotFound, "рецепт не найден")
			}
			return responses.Error(c, fiber.StatusInternalServerError, "не удалось модерировать рецепт", err.Error())
		}
		return responses.Success(c, fiber.StatusOK, item)
	}
}
func KitchenAdminChangeRecipeOwnerHandler() fiber.Handler {
	return func(c fiber.Ctx) error {
		user, authErr := currentKitchenUser(c)
		if authErr != nil {
			return authErr
		}
		if !kitchenUserHasRole(user, "admin") {
			return responses.Error(c, fiber.StatusForbidden, "требуется роль администратора")
		}
		recipeID, err := uuid.Parse(strings.TrimSpace(c.Params("id")))
		if err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "ID рецепта некорректен")
		}
		var req kitchenRecipeOwnerChangeRequest
		if err := c.Bind().Body(&req); err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "некорректный запрос", err.Error())
		}
		ownerUserIDRaw := strings.TrimSpace(req.OwnerUserID)
		if ownerUserIDRaw == "" {
			return responses.Error(c, fiber.StatusBadRequest, "ID владельца рецепта обязателен")
		}
		ownerUserID, err := uuid.Parse(ownerUserIDRaw)
		if err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "ID владельца рецепта некорректен")
		}
		item, err := services.PdbKitchenChangeRecipeOwner(c, recipeID, ownerUserID)
		if err != nil {
			if errors.Is(err, services.ErrKitchenRecipeNotFound) {
				return responses.Error(c, fiber.StatusNotFound, "рецепт не найден")
			}
			if errors.Is(err, services.ErrKitchenRecipeOwnerNotFound) {
				return responses.Error(c, fiber.StatusNotFound, "пользователь-владелец рецепта не найден")
			}
			return responses.Error(c, fiber.StatusInternalServerError, "не удалось сменить владельца рецепта", err.Error())
		}
		return responses.Success(c, fiber.StatusOK, item)
	}
}
func KitchenAccountIngredientsHandler() fiber.Handler {
	return func(c fiber.Ctx) error {
		user, authErr := currentKitchenUser(c)
		if authErr != nil {
			return authErr
		}
		custom, err := services.PdbKitchenListUserIngredients(c, user.UserID)
		if err != nil {
			return responses.Error(c, fiber.StatusInternalServerError, "не удалось загрузить пользовательские ингредиенты", err.Error())
		}
		favorites, err := services.PdbKitchenListFavoriteIngredients(c, user.UserID, "include")
		if err != nil {
			return responses.Error(c, fiber.StatusInternalServerError, "не удалось загрузить избранные ингредиенты", err.Error())
		}
		excludedFavorites, err := services.PdbKitchenListFavoriteIngredients(c, user.UserID, "exclude")
		if err != nil {
			return responses.Error(c, fiber.StatusInternalServerError, "не удалось загрузить исключенные избранные ингредиенты", err.Error())
		}
		return responses.Success(c, fiber.StatusOK, fiber.Map{
			"custom":            custom,
			"favorites":         favorites,
			"favorites_include": favorites,
			"favorites_exclude": excludedFavorites,
		})
	}
}
func KitchenCreateCustomIngredientHandler() fiber.Handler {
	return func(c fiber.Ctx) error {
		user, authErr := currentKitchenUser(c)
		if authErr != nil {
			return authErr
		}
		var req kitchenIngredientMutationRequest
		if err := c.Bind().Body(&req); err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "некорректный запрос", err.Error())
		}
		name := strings.TrimSpace(req.Name)
		if name == "" {
			return responses.Error(c, fiber.StatusBadRequest, "название ингредиента обязательно")
		}
		category := strings.TrimSpace(req.Category)
		if category == "" {
			category = "другое"
		}
		item, err := services.PdbKitchenCreateUserIngredient(c, user.UserID, name, category)
		if err != nil {
			return responses.Error(c, fiber.StatusInternalServerError, "не удалось сохранить пользовательский ингредиент", err.Error())
		}
		return responses.Success(c, fiber.StatusCreated, item)
	}
}
func KitchenDeleteCustomIngredientHandler() fiber.Handler {
	return func(c fiber.Ctx) error {
		user, authErr := currentKitchenUser(c)
		if authErr != nil {
			return authErr
		}
		ingredientID, err := uuid.Parse(strings.TrimSpace(c.Params("id")))
		if err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "ID ингредиента некорректен")
		}
		if err := services.PdbKitchenDeleteUserIngredient(c, user.UserID, ingredientID); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return responses.Error(c, fiber.StatusNotFound, "пользовательский ингредиент не найден")
			}
			return responses.Error(c, fiber.StatusInternalServerError, "не удалось удалить пользовательский ингредиент", err.Error())
		}
		return responses.Success(c, fiber.StatusOK, fiber.Map{"deleted": true})
	}
}
func KitchenFavoriteIngredientHandler() fiber.Handler {
	return func(c fiber.Ctx) error {
		user, authErr := currentKitchenUser(c)
		if authErr != nil {
			return authErr
		}
		var req kitchenIngredientMutationRequest
		if err := c.Bind().Body(&req); err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "некорректный запрос", err.Error())
		}
		if req.IngredientID <= 0 {
			return responses.Error(c, fiber.StatusBadRequest, "ID ингредиента обязателен")
		}
		item, err := services.PdbKitchenFavoriteIngredient(c, user.UserID, req.IngredientID, kitchenFavoriteListType(req.ListType))
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return responses.Error(c, fiber.StatusNotFound, "ингредиент не найден")
			}
			return responses.Error(c, fiber.StatusInternalServerError, "не удалось добавить ингредиент в избранное", err.Error())
		}
		return responses.Success(c, fiber.StatusCreated, item)
	}
}
func KitchenUnfavoriteIngredientHandler() fiber.Handler {
	return func(c fiber.Ctx) error {
		user, authErr := currentKitchenUser(c)
		if authErr != nil {
			return authErr
		}
		rawID := strings.TrimSpace(c.Params("id"))
		ingredientID, err := strconv.ParseInt(rawID, 10, 64)
		if err != nil || ingredientID <= 0 {
			return responses.Error(c, fiber.StatusBadRequest, "ID ингредиента некорректен")
		}
		listType := kitchenFavoriteListType(c.Query("list_type"))
		if err := services.PdbKitchenUnfavoriteIngredient(c, user.UserID, ingredientID, listType); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return responses.Error(c, fiber.StatusNotFound, "избранный ингредиент не найден")
			}
			return responses.Error(c, fiber.StatusInternalServerError, "не удалось удалить ингредиент из избранного", err.Error())
		}
		return responses.Success(c, fiber.StatusOK, fiber.Map{"deleted": true})
	}
}
func KitchenFavoriteRecipeHandler() fiber.Handler {
	return func(c fiber.Ctx) error {
		user, authErr := currentKitchenUser(c)
		if authErr != nil {
			return authErr
		}
		recipeID, err := uuid.Parse(strings.TrimSpace(c.Params("id")))
		if err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "ID рецепта некорректен")
		}
		if _, err := services.PdbKitchenGetRecipeByID(c, recipeID); err != nil {
			if errors.Is(err, services.ErrKitchenRecipeNotFound) {
				return responses.Error(c, fiber.StatusNotFound, "рецепт не найден")
			}
			return responses.Error(c, fiber.StatusInternalServerError, "не удалось загрузить рецепт", err.Error())
		}
		if err := services.PdbKitchenFavoriteRecipe(c, user.UserID, recipeID); err != nil {
			return responses.Error(c, fiber.StatusInternalServerError, "не удалось добавить рецепт в избранное", err.Error())
		}
		return responses.Success(c, fiber.StatusOK, fiber.Map{"favorited": true})
	}
}
func KitchenUnfavoriteRecipeHandler() fiber.Handler {
	return func(c fiber.Ctx) error {
		user, authErr := currentKitchenUser(c)
		if authErr != nil {
			return authErr
		}
		recipeID, err := uuid.Parse(strings.TrimSpace(c.Params("id")))
		if err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "ID рецепта некорректен")
		}
		if err := services.PdbKitchenUnfavoriteRecipe(c, user.UserID, recipeID); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return responses.Error(c, fiber.StatusNotFound, "избранный рецепт не найден")
			}
			return responses.Error(c, fiber.StatusInternalServerError, "не удалось убрать рецепт из избранного", err.Error())
		}
		return responses.Success(c, fiber.StatusOK, fiber.Map{"favorited": false})
	}
}
func KitchenRecipeFavoriteStatusHandler() fiber.Handler {
	return func(c fiber.Ctx) error {
		user, authErr := currentKitchenUser(c)
		if authErr != nil {
			return authErr
		}
		recipeID, err := uuid.Parse(strings.TrimSpace(c.Params("id")))
		if err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "ID рецепта некорректен")
		}
		favorited, err := services.PdbKitchenRecipeIsFavorite(c, user.UserID, recipeID)
		if err != nil {
			return responses.Error(c, fiber.StatusInternalServerError, "не удалось загрузить статус избранного", err.Error())
		}
		return responses.Success(c, fiber.StatusOK, fiber.Map{"favorited": favorited})
	}
}
func KitchenLatestRecipesHandler() fiber.Handler {
	return func(c fiber.Ctx) error {
		limit := 20
		if raw := strings.TrimSpace(c.Query("limit")); raw != "" {
			if n, err := strconv.Atoi(raw); err == nil {
				limit = n
			}
		}
		items, err := services.PdbKitchenLatestRecipes(c, limit)
		if err != nil {
			return responses.Error(c, fiber.StatusInternalServerError, "не удалось получить последние рецепты", err.Error())
		}
		totalPublic, err := services.PdbKitchenPublicRecipesCount(c)
		if err != nil {
			return responses.Error(c, fiber.StatusInternalServerError, "не удалось посчитать публичные рецепты", err.Error())
		}
		return responses.Success(c, fiber.StatusOK, fiber.Map{
			"items":        items,
			"total_public": totalPublic,
		})
	}
}
func KitchenSearchRecipesHandler() fiber.Handler {
	return func(c fiber.Ctx) error {
		var req kitchenRecipeSearchRequest
		if err := c.Bind().Body(&req); err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "некорректный запрос", err.Error())
		}
		params := services.KitchenSearchParams{
			Query:                  strings.TrimSpace(req.Query),
			IngredientTags:         uniqueNormalized(req.Ingredients),
			ExcludedIngredientTags: uniqueNormalized(req.ExcludedIngredients),
			RequireAll:             req.RequireAllIngredients,
			Difficulty:             strings.TrimSpace(req.Difficulty),
			MealType:               strings.TrimSpace(req.MealType),
			CookingMethod:          normalizeKitchenCookingMethod(req.CookingMethod),
			Cuisine:                strings.TrimSpace(req.Cuisine),
			KcalMin:                kitchenIntOrZero(req.KcalMin),
			KcalMax:                kitchenIntOrZero(req.KcalMax),
			MaxTotalMinutes:        kitchenIntOrZero(req.MaxTotalMinutes),
			ServingsMin:            kitchenIntOrZero(req.ServingsMin),
			ServingsMax:            kitchenIntOrZero(req.ServingsMax),
			Limit:                  req.Limit,
			Offset:                 req.Offset,
		}
		params.DietType = normalizeKitchenDietType(req.DietType)
		items, err := services.PdbKitchenSearchRecipes(c, params)
		if err != nil {
			return responses.Error(c, fiber.StatusInternalServerError, "не удалось выполнить поиск рецептов", err.Error())
		}
		mode := "any"
		if params.RequireAll {
			mode = "all"
		}
		return responses.Success(c, fiber.StatusOK, fiber.Map{
			"items":       items,
			"mode":        mode,
			"ingredients": params.IngredientTags,
		})
	}
}
func KitchenGetRecipeHandler() fiber.Handler {
	return func(c fiber.Ctx) error {
		rawID := strings.TrimSpace(c.Params("id"))
		if rawID == "" {
			return responses.Error(c, fiber.StatusBadRequest, "ID обязателен")
		}
		id, err := uuid.Parse(rawID)
		if err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "ID должен быть корректным UUID")
		}
		recipe, err := services.PdbKitchenGetRecipeByID(c, id)
		if err != nil {
			if errors.Is(err, services.ErrKitchenRecipeNotFound) {
				return responses.Error(c, fiber.StatusNotFound, "рецепт не найден")
			}
			return responses.Error(c, fiber.StatusInternalServerError, "не удалось получить рецепт", err.Error())
		}
		return responses.Success(c, fiber.StatusOK, recipe)
	}
}
func KitchenGetManageRecipeHandler() fiber.Handler {
	return func(c fiber.Ctx) error {
		user, authErr := currentKitchenUser(c)
		if authErr != nil {
			return authErr
		}
		rawID := strings.TrimSpace(c.Params("id"))
		if rawID == "" {
			return responses.Error(c, fiber.StatusBadRequest, "ID обязателен")
		}
		id, err := uuid.Parse(rawID)
		if err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "ID должен быть корректным UUID")
		}
		recipe, err := services.PdbKitchenGetRecipeByIDAny(c, id)
		if err != nil {
			if errors.Is(err, services.ErrKitchenRecipeNotFound) {
				return responses.Error(c, fiber.StatusNotFound, "рецепт не найден")
			}
			return responses.Error(c, fiber.StatusInternalServerError, "не удалось загрузить рецепт", err.Error())
		}
		if kitchenUserHasRole(user, "admin") {
			return responses.Success(c, fiber.StatusOK, recipe)
		}
		if recipe.OwnerUserID != nil && *recipe.OwnerUserID == user.UserID {
			return responses.Success(c, fiber.StatusOK, recipe)
		}
		return responses.Error(c, fiber.StatusForbidden, "можно просматривать только свои рецепты")
	}
}
