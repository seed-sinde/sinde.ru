package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"sinde.ru/db"
	"sinde.ru/internal/models"
	"strings"
)

var ErrKitchenRecipeNotFound = errors.New("kitchen recipe not found")
var ErrKitchenRecipeForbidden = errors.New("kitchen recipe forbidden")
var ErrKitchenRecipeOwnerNotFound = errors.New("kitchen recipe owner not found")

type KitchenSearchParams struct {
	Query                  string
	IngredientTags         []string
	ExcludedIngredientTags []string
	RequireAll             bool
	Difficulty             string
	MealType               string
	CookingMethod          string
	Cuisine                string
	DietType               string
	KcalMin                int
	KcalMax                int
	MaxTotalMinutes        int
	ServingsMin            int
	ServingsMax            int
	Limit                  int
	Offset                 int
}
type kitchenScanner interface {
	Scan(dest ...any) error
}
type KitchenCatalogPayload struct {
	Categories    []*models.KitchenCategory          `json:"categories"`
	Ingredients   []*models.KitchenCatalogIngredient `json:"ingredients"`
	FilterOptions []*models.KitchenFilterOption      `json:"filter_options"`
}

const kitchenRecipeSelect = `
	SELECT
		id,
		owner_user_id,
		title,
		description,
		cover_image_key,
		kcal,
		prep_minutes,
		cook_minutes,
		servings,
		difficulty,
		meal_type,
		cooking_method,
		cuisine,
		diet_type,
		ingredients,
		ingredients_search,
	steps,
	tags,
	is_public,
	moderation_status,
	moderated_by_user_id,
	moderated_at,
	moderation_note,
	created_at,
	updated_at
	FROM kitchen_recipes
`

func scanKitchenRecipe(row kitchenScanner) (*models.KitchenRecipe, error) {
	var recipe models.KitchenRecipe
	var ownerUserID pgtype.UUID
	var moderatedByUserID pgtype.UUID
	var moderatedAt pgtype.Timestamptz
	var ingredientsRaw []byte
	var stepsRaw []byte
	if err := row.Scan(
		&recipe.ID,
		&ownerUserID,
		&recipe.Title,
		&recipe.Description,
		&recipe.CoverImageKey,
		&recipe.Kcal,
		&recipe.PrepMinutes,
		&recipe.CookMinutes,
		&recipe.Servings,
		&recipe.Difficulty,
		&recipe.MealType,
		&recipe.CookingMethod,
		&recipe.Cuisine,
		&recipe.DietType,
		&ingredientsRaw,
		&recipe.IngredientsSearch,
		&stepsRaw,
		&recipe.Tags,
		&recipe.IsPublic,
		&recipe.ModerationStatus,
		&moderatedByUserID,
		&moderatedAt,
		&recipe.ModerationNote,
		&recipe.CreatedAt,
		&recipe.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrKitchenRecipeNotFound
		}
		return nil, err
	}
	if ownerUserID.Valid {
		id := uuid.UUID(ownerUserID.Bytes)
		recipe.OwnerUserID = &id
	}
	if moderatedByUserID.Valid {
		id := uuid.UUID(moderatedByUserID.Bytes)
		recipe.ModeratedByUserID = &id
	}
	if moderatedAt.Valid {
		ts := moderatedAt.Time
		recipe.ModeratedAt = &ts
	}
	if len(ingredientsRaw) > 0 {
		if err := json.Unmarshal(ingredientsRaw, &recipe.Ingredients); err != nil {
			return nil, err
		}
	}
	if recipe.Tags == nil {
		recipe.Tags = []string{}
	}
	if recipe.IngredientsSearch == nil {
		recipe.IngredientsSearch = []string{}
	}
	if recipe.Ingredients == nil {
		recipe.Ingredients = []models.KitchenIngredient{}
	}
	recipe.CoverImageKey = ""
	recipe.Steps = []models.KitchenStep{}
	return &recipe, nil
}
func scanKitchenRecipeCollect(row pgx.CollectableRow) (*models.KitchenRecipe, error) {
	return scanKitchenRecipe(row)
}
func firstStorageKeyForUsage(usages []*models.StorageObjectUsage, usageType string, fieldName string) string {
	for _, usage := range usages {
		if usage == nil || usage.Object == nil {
			continue
		}
		if usage.UsageType != usageType || usage.FieldName != fieldName {
			continue
		}
		return usage.Object.StorageKey
	}
	return ""
}
func attachKitchenRecipeCoverUsage(ctx context.Context, recipeID uuid.UUID, storageKey string) error {
	if err := PdbDeleteStorageObjectUsagesBySelector(ctx, "kitchen_recipe", recipeID.String(), "image", "cover"); err != nil {
		return err
	}
	if strings.TrimSpace(storageKey) == "" {
		return nil
	}
	object, err := PdbFindStorageObjectByStorageKey(ctx, storageKey)
	if err != nil {
		return err
	}
	_, err = PdbAttachStorageObjectUsage(ctx, &models.StorageObjectUsage{
		ObjectID:   object.ObjectID,
		EntityType: "kitchen_recipe",
		EntityID:   recipeID.String(),
		UsageType:  "image",
		FieldName:  "cover",
		SortOrder:  0,
		IsPrimary:  true,
		Metadata:   json.RawMessage(`{}`),
	})
	return err
}
func syncKitchenRecipeSteps(ctx context.Context, recipe *models.KitchenRecipe) error {
	existingSteps, err := PdbListKitchenRecipeStepsByRecipeID(ctx, recipe.ID)
	if err != nil {
		return err
	}
	for _, step := range existingSteps {
		if step == nil {
			continue
		}
		if err := PdbDeleteStorageObjectUsagesByEntity(ctx, "kitchen_recipe_step", step.StepID.String()); err != nil {
			return err
		}
	}
	if err := PdbDeleteKitchenRecipeStepsByRecipeID(ctx, recipe.ID); err != nil {
		return err
	}
	nextSteps := make([]models.KitchenStep, 0, len(recipe.Steps))
	for index, step := range recipe.Steps {
		record, createErr := PdbCreateKitchenRecipeStep(ctx, &models.KitchenRecipeStep{
			StepID:      uuid.New(),
			RecipeID:    recipe.ID,
			StepOrder:   index + 1,
			Title:       "",
			Description: step.Text,
			Metadata:    json.RawMessage(`{}`),
		})
		if createErr != nil {
			return createErr
		}
		imageKey := strings.TrimSpace(step.ImageKey)
		if imageKey != "" {
			object, findErr := PdbFindStorageObjectByStorageKey(ctx, imageKey)
			if findErr != nil {
				return findErr
			}
			if _, attachErr := PdbAttachStorageObjectUsage(ctx, &models.StorageObjectUsage{
				ObjectID:   object.ObjectID,
				EntityType: "kitchen_recipe_step",
				EntityID:   record.StepID.String(),
				UsageType:  "image",
				FieldName:  "content",
				SortOrder:  0,
				IsPrimary:  true,
				Metadata:   json.RawMessage(`{}`),
			}); attachErr != nil {
				return attachErr
			}
		}
		nextSteps = append(nextSteps, models.KitchenStep{
			Order:    record.StepOrder,
			Text:     record.Description,
			ImageKey: imageKey,
		})
	}
	recipe.Steps = nextSteps
	return nil
}
func hydrateKitchenRecipe(ctx context.Context, recipe *models.KitchenRecipe) error {
	if recipe == nil {
		return nil
	}
	if usages, err := PdbListStorageObjectUsagesByEntity(ctx, "kitchen_recipe", recipe.ID.String()); err != nil {
		return err
	} else if coverKey := firstStorageKeyForUsage(usages, "image", "cover"); coverKey != "" {
		recipe.CoverImageKey = coverKey
	}
	stepRecords, err := PdbListKitchenRecipeStepsByRecipeID(ctx, recipe.ID)
	if err != nil {
		return err
	}
	if len(stepRecords) == 0 {
		return nil
	}
	steps := make([]models.KitchenStep, 0, len(stepRecords))
	for _, step := range stepRecords {
		if step == nil {
			continue
		}
		usages, usageErr := PdbListStorageObjectUsagesByEntity(ctx, "kitchen_recipe_step", step.StepID.String())
		if usageErr != nil {
			return usageErr
		}
		text := strings.TrimSpace(step.Description)
		if text == "" {
			text = strings.TrimSpace(step.Title)
		}
		steps = append(steps, models.KitchenStep{
			Order:    step.StepOrder,
			Text:     text,
			ImageKey: firstStorageKeyForUsage(usages, "image", "content"),
		})
	}
	recipe.Steps = steps
	return nil
}
func hydrateKitchenRecipes(ctx context.Context, recipes []*models.KitchenRecipe) error {
	for _, recipe := range recipes {
		if err := hydrateKitchenRecipe(ctx, recipe); err != nil {
			return err
		}
	}
	return nil
}
func PdbKitchenCreateRecipe(ctx context.Context, recipe *models.KitchenRecipe) (*models.KitchenRecipe, error) {
	ingredientsRaw, err := json.Marshal(recipe.Ingredients)
	if err != nil {
		return nil, err
	}
	insertQuery := `
		INSERT INTO kitchen_recipes (
			id, owner_user_id, title, description, cover_image_key, kcal, prep_minutes, cook_minutes, servings,
			difficulty, meal_type, cooking_method, cuisine, diet_type,
			ingredients, ingredients_search, steps, tags, is_public,
			moderation_status, moderated_by_user_id, moderated_at, moderation_note
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9,
			$10, $11, $12, $13, $14,
			$15, $16, $17, $18, $19,
			$20, $21, $22, $23
		)
		RETURNING
			id,
			owner_user_id,
			title,
			description,
			cover_image_key,
			kcal,
			prep_minutes,
			cook_minutes,
			servings,
			difficulty,
			meal_type,
			cooking_method,
			cuisine,
			diet_type,
			ingredients,
			ingredients_search,
			steps,
			tags,
			is_public,
			moderation_status,
			moderated_by_user_id,
			moderated_at,
			moderation_note,
			created_at,
			updated_at
	`
	row := db.PDB.QueryRow(
		ctx,
		insertQuery,
		recipe.ID,
		recipe.OwnerUserID,
		recipe.Title,
		recipe.Description,
		"",
		recipe.Kcal,
		recipe.PrepMinutes,
		recipe.CookMinutes,
		recipe.Servings,
		recipe.Difficulty,
		recipe.MealType,
		recipe.CookingMethod,
		recipe.Cuisine,
		recipe.DietType,
		ingredientsRaw,
		recipe.IngredientsSearch,
		json.RawMessage(`[]`),
		recipe.Tags,
		recipe.IsPublic,
		recipe.ModerationStatus,
		recipe.ModeratedByUserID,
		recipe.ModeratedAt,
		recipe.ModerationNote,
	)
	item, err := scanKitchenRecipe(row)
	if err != nil {
		return nil, err
	}
	if err := attachKitchenRecipeCoverUsage(ctx, item.ID, recipe.CoverImageKey); err != nil {
		return nil, err
	}
	if err := syncKitchenRecipeSteps(ctx, recipe); err != nil {
		return nil, err
	}
	if err := hydrateKitchenRecipe(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}
func PdbKitchenUpdateRecipe(ctx context.Context, recipe *models.KitchenRecipe, ownerUserID uuid.UUID) (*models.KitchenRecipe, error) {
	ingredientsRaw, err := json.Marshal(recipe.Ingredients)
	if err != nil {
		return nil, err
	}
	query := `
		UPDATE kitchen_recipes
		SET
			title = $3,
			description = $4,
			cover_image_key = $5,
			kcal = $6,
			prep_minutes = $7,
			cook_minutes = $8,
			servings = $9,
			difficulty = $10,
			meal_type = $11,
			cooking_method = $12,
			cuisine = $13,
			diet_type = $14,
			ingredients = $15,
			ingredients_search = $16,
			steps = $17,
			tags = $18,
			is_public = $19,
			moderation_status = $20,
			moderated_by_user_id = $21,
			moderated_at = $22,
			moderation_note = $23
		WHERE id = $1 AND owner_user_id = $2
		RETURNING
			id,
			owner_user_id,
			title,
			description,
			cover_image_key,
			kcal,
			prep_minutes,
			cook_minutes,
			servings,
			difficulty,
			meal_type,
			cooking_method,
			cuisine,
			diet_type,
			ingredients,
			ingredients_search,
			steps,
			tags,
			is_public,
			moderation_status,
			moderated_by_user_id,
			moderated_at,
			moderation_note,
			created_at,
			updated_at
	`
	row := db.PDB.QueryRow(
		ctx,
		query,
		recipe.ID,
		ownerUserID,
		recipe.Title,
		recipe.Description,
		"",
		recipe.Kcal,
		recipe.PrepMinutes,
		recipe.CookMinutes,
		recipe.Servings,
		recipe.Difficulty,
		recipe.MealType,
		recipe.CookingMethod,
		recipe.Cuisine,
		recipe.DietType,
		ingredientsRaw,
		recipe.IngredientsSearch,
		json.RawMessage(`[]`),
		recipe.Tags,
		recipe.IsPublic,
		recipe.ModerationStatus,
		recipe.ModeratedByUserID,
		recipe.ModeratedAt,
		recipe.ModerationNote,
	)
	item, err := scanKitchenRecipe(row)
	if err != nil {
		if errors.Is(err, ErrKitchenRecipeNotFound) {
			return nil, ErrKitchenRecipeForbidden
		}
		return nil, err
	}
	if err := attachKitchenRecipeCoverUsage(ctx, item.ID, recipe.CoverImageKey); err != nil {
		return nil, err
	}
	if err := syncKitchenRecipeSteps(ctx, recipe); err != nil {
		return nil, err
	}
	if err := hydrateKitchenRecipe(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}
func PdbKitchenDeleteRecipe(ctx context.Context, recipeID uuid.UUID, ownerUserID uuid.UUID) error {
	stepRecords, err := PdbListKitchenRecipeStepsByRecipeID(ctx, recipeID)
	if err != nil {
		return err
	}
	tag, err := db.PDB.Exec(ctx, `DELETE FROM kitchen_recipes WHERE id = $1 AND owner_user_id = $2`, recipeID, ownerUserID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrKitchenRecipeForbidden
	}
	for _, step := range stepRecords {
		if step == nil {
			continue
		}
		if err := PdbDeleteStorageObjectUsagesByEntity(ctx, "kitchen_recipe_step", step.StepID.String()); err != nil {
			return err
		}
	}
	if err := PdbDeleteStorageObjectUsagesByEntity(ctx, "kitchen_recipe", recipeID.String()); err != nil {
		return err
	}
	return nil
}
func PdbKitchenUpdateRecipeAsAdmin(ctx context.Context, recipe *models.KitchenRecipe) (*models.KitchenRecipe, error) {
	ingredientsRaw, err := json.Marshal(recipe.Ingredients)
	if err != nil {
		return nil, err
	}
	query := `
		UPDATE kitchen_recipes
		SET
			title = $2,
			description = $3,
			cover_image_key = $4,
			kcal = $5,
			prep_minutes = $6,
			cook_minutes = $7,
			servings = $8,
			difficulty = $9,
			meal_type = $10,
			cooking_method = $11,
			cuisine = $12,
			diet_type = $13,
			ingredients = $14,
			ingredients_search = $15,
			steps = $16,
			tags = $17,
			is_public = $18,
			moderation_status = $19,
			moderated_by_user_id = $20,
			moderated_at = $21,
			moderation_note = $22
		WHERE id = $1
		RETURNING
			id,
			owner_user_id,
			title,
			description,
			cover_image_key,
			kcal,
			prep_minutes,
			cook_minutes,
			servings,
			difficulty,
			meal_type,
			cooking_method,
			cuisine,
			diet_type,
			ingredients,
			ingredients_search,
			steps,
			tags,
			is_public,
			moderation_status,
			moderated_by_user_id,
			moderated_at,
			moderation_note,
			created_at,
			updated_at
	`
	row := db.PDB.QueryRow(
		ctx,
		query,
		recipe.ID,
		recipe.Title,
		recipe.Description,
		"",
		recipe.Kcal,
		recipe.PrepMinutes,
		recipe.CookMinutes,
		recipe.Servings,
		recipe.Difficulty,
		recipe.MealType,
		recipe.CookingMethod,
		recipe.Cuisine,
		recipe.DietType,
		ingredientsRaw,
		recipe.IngredientsSearch,
		json.RawMessage(`[]`),
		recipe.Tags,
		recipe.IsPublic,
		recipe.ModerationStatus,
		recipe.ModeratedByUserID,
		recipe.ModeratedAt,
		recipe.ModerationNote,
	)
	item, err := scanKitchenRecipe(row)
	if err != nil {
		if errors.Is(err, ErrKitchenRecipeNotFound) {
			return nil, ErrKitchenRecipeNotFound
		}
		return nil, err
	}
	if err := attachKitchenRecipeCoverUsage(ctx, item.ID, recipe.CoverImageKey); err != nil {
		return nil, err
	}
	if err := syncKitchenRecipeSteps(ctx, recipe); err != nil {
		return nil, err
	}
	if err := hydrateKitchenRecipe(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}
func PdbKitchenDeleteRecipeAsAdmin(ctx context.Context, recipeID uuid.UUID) error {
	stepRecords, err := PdbListKitchenRecipeStepsByRecipeID(ctx, recipeID)
	if err != nil {
		return err
	}
	tag, err := db.PDB.Exec(ctx, `DELETE FROM kitchen_recipes WHERE id = $1`, recipeID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrKitchenRecipeNotFound
	}
	for _, step := range stepRecords {
		if step == nil {
			continue
		}
		if err := PdbDeleteStorageObjectUsagesByEntity(ctx, "kitchen_recipe_step", step.StepID.String()); err != nil {
			return err
		}
	}
	if err := PdbDeleteStorageObjectUsagesByEntity(ctx, "kitchen_recipe", recipeID.String()); err != nil {
		return err
	}
	return nil
}
func PdbKitchenGetRecipeByID(ctx context.Context, id uuid.UUID) (*models.KitchenRecipe, error) {
	query := kitchenRecipeSelect + `
		WHERE id = $1 AND is_public = TRUE
	`
	row := db.PDB.QueryRow(ctx, query, id)
	item, err := scanKitchenRecipe(row)
	if err != nil {
		return nil, err
	}
	if err := hydrateKitchenRecipe(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}
func PdbKitchenGetRecipeByIDAny(ctx context.Context, id uuid.UUID) (*models.KitchenRecipe, error) {
	query := kitchenRecipeSelect + `
		WHERE id = $1
	`
	row := db.PDB.QueryRow(ctx, query, id)
	item, err := scanKitchenRecipe(row)
	if err != nil {
		return nil, err
	}
	if err := hydrateKitchenRecipe(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}
func PdbKitchenListRecipesByOwner(ctx context.Context, ownerUserID uuid.UUID, limit int) ([]*models.KitchenRecipe, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}
	query := kitchenRecipeSelect + `
		WHERE owner_user_id = $1
		ORDER BY updated_at DESC, created_at DESC
		LIMIT $2
	`
	rows, err := db.PDB.Query(ctx, query, ownerUserID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items, err := pgx.CollectRows(rows, scanKitchenRecipeCollect)
	if err != nil {
		return nil, err
	}
	if items == nil {
		return []*models.KitchenRecipe{}, nil
	}
	if err := hydrateKitchenRecipes(ctx, items); err != nil {
		return nil, err
	}
	return items, nil
}
func PdbKitchenLatestRecipes(ctx context.Context, limit int) ([]*models.KitchenRecipe, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	query := kitchenRecipeSelect + `
		WHERE is_public = TRUE
		ORDER BY created_at DESC
		LIMIT $1
	`
	rows, err := db.PDB.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items, err := pgx.CollectRows(rows, scanKitchenRecipeCollect)
	if err != nil {
		return nil, err
	}
	if items == nil {
		return []*models.KitchenRecipe{}, nil
	}
	if err := hydrateKitchenRecipes(ctx, items); err != nil {
		return nil, err
	}
	return items, nil
}
func PdbKitchenPublicRecipesCount(ctx context.Context) (int64, error) {
	var total int64
	if err := db.PDB.QueryRow(ctx, `SELECT COUNT(*)::BIGINT FROM kitchen_recipes WHERE is_public = TRUE`).Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}
func NormalizeKitchenModerationStatus(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "all":
		return "all"
	case "draft":
		return "draft"
	case "pending":
		return "pending"
	case "approved":
		return "approved"
	case "rejected":
		return "rejected"
	default:
		return ""
	}
}
func PdbKitchenListRecipesForModeration(ctx context.Context, status string, limit, offset int) ([]*models.KitchenRecipe, int64, error) {
	if limit <= 0 {
		limit = 30
	}
	if limit > 200 {
		limit = 200
	}
	if offset < 0 {
		offset = 0
	}
	status = NormalizeKitchenModerationStatus(status)
	if status == "" {
		status = "pending"
	}
	baseWhere := " WHERE moderation_status = 'pending'"
	countWhere := baseWhere
	args := []any{limit, offset}
	if status == "all" {
		baseWhere = ""
		countWhere = ""
	} else if status != "pending" {
		baseWhere = " WHERE moderation_status = $3"
		countWhere = " WHERE moderation_status = $1"
		args = append(args, status)
	}
	query := kitchenRecipeSelect + baseWhere + `
		ORDER BY created_at ASC
		LIMIT $1 OFFSET $2
	`
	rows, err := db.PDB.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items, err := pgx.CollectRows(rows, scanKitchenRecipeCollect)
	if err != nil {
		return nil, 0, err
	}
	if items == nil {
		items = []*models.KitchenRecipe{}
	}
	if err := hydrateKitchenRecipes(ctx, items); err != nil {
		return nil, 0, err
	}
	countQuery := `SELECT COUNT(*)::BIGINT FROM kitchen_recipes` + countWhere
	var total int64
	countArgs := []any{}
	if status != "all" && status != "pending" {
		countArgs = append(countArgs, status)
	}
	if err := db.PDB.QueryRow(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}
	return items, total, nil
}
func PdbKitchenModerationStatusTotals(ctx context.Context) (map[string]int64, error) {
	totals := map[string]int64{
		"all":      0,
		"pending":  0,
		"approved": 0,
		"rejected": 0,
		"draft":    0,
	}
	rows, err := db.PDB.Query(ctx, `
		SELECT moderation_status, COUNT(*)::BIGINT
		FROM kitchen_recipes
		GROUP BY moderation_status
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var status string
		var count int64
		if scanErr := rows.Scan(&status, &count); scanErr != nil {
			return nil, scanErr
		}
		normalized := NormalizeKitchenModerationStatus(status)
		if normalized == "" || normalized == "all" {
			continue
		}
		totals[normalized] = count
		totals["all"] += count
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return totals, nil
}
func PdbKitchenModerateRecipe(ctx context.Context, recipeID uuid.UUID, approve bool, adminUserID uuid.UUID, note string) (*models.KitchenRecipe, error) {
	status := "rejected"
	isPublic := false
	normalizedNote := strings.TrimSpace(note)
	if approve {
		status = "approved"
		isPublic = true
		normalizedNote = ""
	}
	row := db.PDB.QueryRow(ctx, `
		UPDATE kitchen_recipes
		SET
			is_public = $2,
			moderation_status = $3,
			moderated_by_user_id = $4,
			moderated_at = now(),
			moderation_note = $5
		WHERE id = $1
		RETURNING
			id,
			owner_user_id,
			title,
			description,
			cover_image_key,
			kcal,
			prep_minutes,
			cook_minutes,
			servings,
			difficulty,
			meal_type,
			cooking_method,
			cuisine,
			diet_type,
			ingredients,
			ingredients_search,
			steps,
			tags,
			is_public,
			moderation_status,
			moderated_by_user_id,
			moderated_at,
			moderation_note,
			created_at,
			updated_at
	`, recipeID, isPublic, status, adminUserID, normalizedNote)
	item, err := scanKitchenRecipe(row)
	if err != nil {
		if errors.Is(err, ErrKitchenRecipeNotFound) {
			return nil, ErrKitchenRecipeNotFound
		}
		return nil, err
	}
	if err := hydrateKitchenRecipe(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}
func PdbKitchenChangeRecipeOwner(ctx context.Context, recipeID uuid.UUID, ownerUserID uuid.UUID) (*models.KitchenRecipe, error) {
	row := db.PDB.QueryRow(ctx, `
		UPDATE kitchen_recipes
		SET
			owner_user_id = $2,
			updated_at = now()
		WHERE id = $1
		RETURNING
			id,
			owner_user_id,
			title,
			description,
			cover_image_key,
			kcal,
			prep_minutes,
			cook_minutes,
			servings,
			difficulty,
			meal_type,
			cooking_method,
			cuisine,
			diet_type,
			ingredients,
			ingredients_search,
			steps,
			tags,
			is_public,
			moderation_status,
			moderated_by_user_id,
			moderated_at,
			moderation_note,
			created_at,
			updated_at
	`, recipeID, ownerUserID)
	item, err := scanKitchenRecipe(row)
	if err != nil {
		if errors.Is(err, ErrKitchenRecipeNotFound) {
			return nil, ErrKitchenRecipeNotFound
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return nil, ErrKitchenRecipeOwnerNotFound
		}
		return nil, err
	}
	if err := hydrateKitchenRecipe(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}
func normalizeKitchenIngredientName(value string) string {
	return strings.ToLower(strings.Join(strings.Fields(strings.TrimSpace(value)), " "))
}
func PdbKitchenListUserIngredients(ctx context.Context, userID uuid.UUID) ([]*models.KitchenUserIngredient, error) {
	rows, err := db.PDB.Query(ctx, `
		SELECT ingredient_id, user_id, name, category, normalized_name, created_at, updated_at
		FROM kitchen_user_ingredients
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*models.KitchenUserIngredient, error) {
		var item models.KitchenUserIngredient
		if err := row.Scan(
			&item.IngredientID,
			&item.UserID,
			&item.Name,
			&item.Category,
			&item.NormalizedName,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		return &item, nil
	})
	if err != nil {
		return nil, err
	}
	if items == nil {
		return []*models.KitchenUserIngredient{}, nil
	}
	return items, nil
}
func PdbKitchenDeleteUserIngredient(ctx context.Context, userID uuid.UUID, ingredientID uuid.UUID) error {
	row := db.PDB.QueryRow(ctx, `
		DELETE FROM kitchen_user_ingredients
		WHERE user_id = $1 AND ingredient_id = $2
		RETURNING normalized_name
	`, userID, ingredientID)
	var normalized string
	if err := row.Scan(&normalized); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgx.ErrNoRows
		}
		return err
	}
	if strings.TrimSpace(normalized) != "" {
		if _, err := db.PDB.Exec(ctx, `
			DELETE FROM kitchen_favorite_ingredients
			WHERE user_id = $1 AND normalized_name = $2
		`, userID, normalized); err != nil {
			return err
		}
	}
	return nil
}
func PdbKitchenCreateUserIngredient(ctx context.Context, userID uuid.UUID, name string, category string) (*models.KitchenUserIngredient, error) {
	normalized := normalizeKitchenIngredientName(name)
	query := `
		INSERT INTO kitchen_user_ingredients (
			ingredient_id, user_id, name, category, normalized_name
		)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, normalized_name) DO UPDATE
		SET
			name = EXCLUDED.name,
			category = EXCLUDED.category
		RETURNING ingredient_id, user_id, name, category, normalized_name, created_at, updated_at
	`
	row := db.PDB.QueryRow(ctx, query, uuid.New(), userID, strings.TrimSpace(name), strings.TrimSpace(category), normalized)
	var item models.KitchenUserIngredient
	if err := row.Scan(
		&item.IngredientID,
		&item.UserID,
		&item.Name,
		&item.Category,
		&item.NormalizedName,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &item, nil
}
func normalizeKitchenFavoriteListType(raw string) string {
	listType := strings.ToLower(strings.TrimSpace(raw))
	if listType == "exclude" {
		return "exclude"
	}
	return "include"
}
func PdbKitchenListFavoriteIngredients(ctx context.Context, userID uuid.UUID, listType string) ([]*models.KitchenFavoriteIngredient, error) {
	scope := normalizeKitchenFavoriteListType(listType)
	rows, err := db.PDB.Query(ctx, `
		SELECT ingredient_id, user_id, list_type, name, category, normalized_name, created_at
		FROM kitchen_favorite_ingredients
		WHERE user_id = $1
		  AND list_type = $2
		ORDER BY created_at DESC
	`, userID, scope)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*models.KitchenFavoriteIngredient, error) {
		var item models.KitchenFavoriteIngredient
		if err := row.Scan(&item.IngredientID, &item.UserID, &item.ListType, &item.Name, &item.Category, &item.NormalizedName, &item.CreatedAt); err != nil {
			return nil, err
		}
		return &item, nil
	})
	if err != nil {
		return nil, err
	}
	if items == nil {
		return []*models.KitchenFavoriteIngredient{}, nil
	}
	return items, nil
}
func PdbKitchenFavoriteIngredient(ctx context.Context, userID uuid.UUID, ingredientID int64, listType string) (*models.KitchenFavoriteIngredient, error) {
	scope := normalizeKitchenFavoriteListType(listType)
	row := db.PDB.QueryRow(ctx, `
		INSERT INTO kitchen_favorite_ingredients (user_id, ingredient_id, list_type, name, category, normalized_name)
		SELECT $1, i.ingredient_id, $3, i.name, c.label, i.normalized_name
		FROM kitchen_ingredients_catalog i
		JOIN kitchen_categories c ON c.category_id = i.category_id
		WHERE i.ingredient_id = $2
		ON CONFLICT (user_id, ingredient_id, list_type) DO UPDATE
		SET
			name = EXCLUDED.name,
			category = EXCLUDED.category,
			normalized_name = EXCLUDED.normalized_name
		RETURNING ingredient_id, user_id, list_type, name, category, normalized_name, created_at
	`, userID, ingredientID, scope)
	var item models.KitchenFavoriteIngredient
	if err := row.Scan(&item.IngredientID, &item.UserID, &item.ListType, &item.Name, &item.Category, &item.NormalizedName, &item.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}
		return nil, err
	}
	return &item, nil
}
func PdbKitchenUnfavoriteIngredient(ctx context.Context, userID uuid.UUID, ingredientID int64, listType string) error {
	scope := normalizeKitchenFavoriteListType(listType)
	tag, err := db.PDB.Exec(ctx, `
		DELETE FROM kitchen_favorite_ingredients
		WHERE user_id = $1 AND ingredient_id = $2 AND list_type = $3
	`, userID, ingredientID, scope)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
func PdbKitchenFavoriteRecipe(ctx context.Context, userID, recipeID uuid.UUID) error {
	_, err := db.PDB.Exec(ctx, `
		INSERT INTO kitchen_favorite_recipes (user_id, recipe_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, recipe_id) DO NOTHING
	`, userID, recipeID)
	return err
}
func PdbKitchenUnfavoriteRecipe(ctx context.Context, userID, recipeID uuid.UUID) error {
	tag, err := db.PDB.Exec(ctx, `
		DELETE FROM kitchen_favorite_recipes
		WHERE user_id = $1 AND recipe_id = $2
	`, userID, recipeID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
func PdbKitchenRecipeIsFavorite(ctx context.Context, userID, recipeID uuid.UUID) (bool, error) {
	var exists bool
	if err := db.PDB.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1
			FROM kitchen_favorite_recipes
			WHERE user_id = $1 AND recipe_id = $2
		)
	`, userID, recipeID).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}
func PdbKitchenListCategories(ctx context.Context) ([]*models.KitchenCategory, error) {
	rows, err := db.PDB.Query(ctx, `
		SELECT category_id, slug, label, description, sort_order
		FROM kitchen_categories
		ORDER BY sort_order ASC, category_id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*models.KitchenCategory, error) {
		var item models.KitchenCategory
		if err := row.Scan(
			&item.CategoryID,
			&item.Slug,
			&item.Label,
			&item.Description,
			&item.SortOrder,
		); err != nil {
			return nil, err
		}
		return &item, nil
	})
	if err != nil {
		return nil, err
	}
	if items == nil {
		return []*models.KitchenCategory{}, nil
	}
	return items, nil
}
func PdbKitchenListCatalogIngredients(ctx context.Context) ([]*models.KitchenCatalogIngredient, error) {
	rows, err := db.PDB.Query(ctx, `
		SELECT
			i.ingredient_id,
			i.name,
			i.normalized_name,
			i.category_id,
			c.label AS category_label,
			i.description,
			i.protein_g::float8,
			i.fat_g::float8,
			i.carbs_g::float8,
			i.kcal::float8,
			i.sort_order
		FROM kitchen_ingredients_catalog i
		JOIN kitchen_categories c ON c.category_id = i.category_id
		ORDER BY c.sort_order ASC, i.sort_order ASC, i.ingredient_id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*models.KitchenCatalogIngredient, error) {
		var item models.KitchenCatalogIngredient
		if err := row.Scan(
			&item.IngredientID,
			&item.Name,
			&item.NormalizedName,
			&item.CategoryID,
			&item.CategoryLabel,
			&item.Description,
			&item.ProteinG,
			&item.FatG,
			&item.CarbsG,
			&item.Kcal,
			&item.SortOrder,
		); err != nil {
			return nil, err
		}
		return &item, nil
	})
	if err != nil {
		return nil, err
	}
	if items == nil {
		return []*models.KitchenCatalogIngredient{}, nil
	}
	return items, nil
}
func PdbKitchenListFilterOptions(ctx context.Context) ([]*models.KitchenFilterOption, error) {
	rows, err := db.PDB.Query(ctx, `
		SELECT option_id, kind, code, label, description, sort_order
		FROM kitchen_filter_options
		ORDER BY kind ASC, sort_order ASC, option_id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*models.KitchenFilterOption, error) {
		var item models.KitchenFilterOption
		if err := row.Scan(
			&item.OptionID,
			&item.Kind,
			&item.Code,
			&item.Label,
			&item.Description,
			&item.SortOrder,
		); err != nil {
			return nil, err
		}
		return &item, nil
	})
	if err != nil {
		return nil, err
	}
	if items == nil {
		return []*models.KitchenFilterOption{}, nil
	}
	return items, nil
}
func PdbKitchenCatalog(ctx context.Context) (*KitchenCatalogPayload, error) {
	categories, err := PdbKitchenListCategories(ctx)
	if err != nil {
		return nil, err
	}
	ingredients, err := PdbKitchenListCatalogIngredients(ctx)
	if err != nil {
		return nil, err
	}
	filterOptions, err := PdbKitchenListFilterOptions(ctx)
	if err != nil {
		return nil, err
	}
	return &KitchenCatalogPayload{
		Categories:    categories,
		Ingredients:   ingredients,
		FilterOptions: filterOptions,
	}, nil
}
func PdbKitchenSearchRecipes(ctx context.Context, params KitchenSearchParams) ([]*models.KitchenRecipe, error) {
	if params.Limit <= 0 {
		params.Limit = 20
	}
	if params.Limit > 100 {
		params.Limit = 100
	}
	if params.Offset < 0 {
		params.Offset = 0
	}
	if params.MaxTotalMinutes < 0 {
		params.MaxTotalMinutes = 0
	}
	if params.KcalMin < 0 {
		params.KcalMin = 0
	}
	if params.KcalMax < 0 {
		params.KcalMax = 0
	}
	if params.KcalMin > 0 && params.KcalMax > 0 && params.KcalMax < params.KcalMin {
		params.KcalMin, params.KcalMax = params.KcalMax, params.KcalMin
	}
	if params.ServingsMin < 0 {
		params.ServingsMin = 0
	}
	if params.ServingsMax < 0 {
		params.ServingsMax = 0
	}
	if params.ServingsMin > 0 && params.ServingsMax > 0 && params.ServingsMax < params.ServingsMin {
		params.ServingsMin, params.ServingsMax = params.ServingsMax, params.ServingsMin
	}
	query := strings.Builder{}
	query.WriteString(kitchenRecipeSelect)
	query.WriteString(" WHERE is_public = TRUE")
	args := make([]any, 0, 4)
	idx := 1
	if q := strings.TrimSpace(params.Query); q != "" {
		args = append(args, "%"+q+"%")
		query.WriteString(fmt.Sprintf(" AND (title ILIKE $%d OR description ILIKE $%d)", idx, idx))
		idx++
	}
	if len(params.IngredientTags) > 0 {
		args = append(args, params.IngredientTags)
		if params.RequireAll {
			query.WriteString(fmt.Sprintf(" AND $%d::text[] <@ ingredients_search", idx))
		} else {
			query.WriteString(fmt.Sprintf(" AND ingredients_search && $%d::text[]", idx))
		}
		idx++
	}
	if len(params.ExcludedIngredientTags) > 0 {
		args = append(args, params.ExcludedIngredientTags)
		query.WriteString(fmt.Sprintf(" AND NOT (ingredients_search && $%d::text[])", idx))
		idx++
	}
	if v := strings.TrimSpace(params.Difficulty); v != "" {
		args = append(args, v)
		query.WriteString(fmt.Sprintf(" AND difficulty = $%d", idx))
		idx++
	}
	if v := strings.TrimSpace(params.MealType); v != "" {
		args = append(args, v)
		query.WriteString(fmt.Sprintf(" AND meal_type = $%d", idx))
		idx++
	}
	if v := strings.TrimSpace(params.CookingMethod); v != "" {
		args = append(args, v)
		query.WriteString(fmt.Sprintf(" AND cooking_method = $%d", idx))
		idx++
	}
	if v := strings.TrimSpace(params.Cuisine); v != "" {
		args = append(args, "%"+v+"%")
		query.WriteString(fmt.Sprintf(" AND cuisine ILIKE $%d", idx))
		idx++
	}
	if v := strings.TrimSpace(params.DietType); v != "" {
		args = append(args, v)
		query.WriteString(fmt.Sprintf(" AND diet_type = $%d", idx))
		idx++
	}
	if params.KcalMin > 0 {
		args = append(args, params.KcalMin)
		query.WriteString(fmt.Sprintf(" AND kcal >= $%d", idx))
		idx++
	}
	if params.KcalMax > 0 {
		args = append(args, params.KcalMax)
		query.WriteString(fmt.Sprintf(" AND kcal <= $%d", idx))
		idx++
	}
	if params.MaxTotalMinutes > 0 {
		args = append(args, params.MaxTotalMinutes)
		query.WriteString(fmt.Sprintf(" AND (prep_minutes + cook_minutes) <= $%d", idx))
		idx++
	}
	if params.ServingsMin > 0 {
		args = append(args, params.ServingsMin)
		query.WriteString(fmt.Sprintf(" AND servings >= $%d", idx))
		idx++
	}
	if params.ServingsMax > 0 {
		args = append(args, params.ServingsMax)
		query.WriteString(fmt.Sprintf(" AND servings <= $%d", idx))
		idx++
	}
	args = append(args, params.Limit, params.Offset)
	query.WriteString(fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", idx, idx+1))
	rows, err := db.PDB.Query(ctx, query.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items, err := pgx.CollectRows(rows, scanKitchenRecipeCollect)
	if err != nil {
		return nil, err
	}
	if items == nil {
		return []*models.KitchenRecipe{}, nil
	}
	if err := hydrateKitchenRecipes(ctx, items); err != nil {
		return nil, err
	}
	return items, nil
}
