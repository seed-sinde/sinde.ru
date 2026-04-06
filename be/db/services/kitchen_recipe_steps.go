package services

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"sinde.ru/db"
	"sinde.ru/internal/models"
)

var ErrKitchenRecipeStepNotFound = errors.New("kitchen recipe step not found")

func scanKitchenRecipeStep(row interface{ Scan(dest ...any) error }) (*models.KitchenRecipeStep, error) {
	var item models.KitchenRecipeStep
	var metadata []byte
	if err := row.Scan(
		&item.StepID,
		&item.RecipeID,
		&item.StepOrder,
		&item.Title,
		&item.Description,
		&metadata,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrKitchenRecipeStepNotFound
		}
		return nil, err
	}
	item.Metadata = normalizeJSONMetadata(metadata)
	return &item, nil
}

func PdbListKitchenRecipeStepsByRecipeID(ctx context.Context, recipeID uuid.UUID) ([]*models.KitchenRecipeStep, error) {
	rows, err := db.PDB.Query(ctx, `
		SELECT
			step_id,
			recipe_id,
			step_order,
			title,
			description,
			metadata,
			created_at,
			updated_at
		FROM kitchen_recipe_steps
		WHERE recipe_id = $1
		ORDER BY step_order ASC, created_at ASC
	`, recipeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*models.KitchenRecipeStep, 0, 8)
	for rows.Next() {
		item, scanErr := scanKitchenRecipeStep(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func PdbCreateKitchenRecipeStep(ctx context.Context, step *models.KitchenRecipeStep) (*models.KitchenRecipeStep, error) {
	if step == nil {
		return nil, errors.New("kitchen recipe step is required")
	}
	row := db.PDB.QueryRow(ctx, `
		INSERT INTO kitchen_recipe_steps (
			step_id,
			recipe_id,
			step_order,
			title,
			description,
			metadata
		)
		VALUES ($1, $2, $3, $4, $5, $6::jsonb)
		RETURNING
			step_id,
			recipe_id,
			step_order,
			title,
			description,
			metadata,
			created_at,
			updated_at
	`,
		step.StepID,
		step.RecipeID,
		step.StepOrder,
		step.Title,
		step.Description,
		normalizeJSONMetadata(step.Metadata),
	)
	return scanKitchenRecipeStep(row)
}

func PdbUpdateKitchenRecipeStep(ctx context.Context, step *models.KitchenRecipeStep) (*models.KitchenRecipeStep, error) {
	if step == nil {
		return nil, errors.New("kitchen recipe step is required")
	}
	row := db.PDB.QueryRow(ctx, `
		UPDATE kitchen_recipe_steps
		SET
			step_order = $2,
			title = $3,
			description = $4,
			metadata = $5::jsonb
		WHERE step_id = $1
		RETURNING
			step_id,
			recipe_id,
			step_order,
			title,
			description,
			metadata,
			created_at,
			updated_at
	`,
		step.StepID,
		step.StepOrder,
		step.Title,
		step.Description,
		normalizeJSONMetadata(step.Metadata),
	)
	return scanKitchenRecipeStep(row)
}

func PdbDeleteKitchenRecipeStep(ctx context.Context, stepID uuid.UUID) error {
	tag, err := db.PDB.Exec(ctx, `DELETE FROM kitchen_recipe_steps WHERE step_id = $1`, stepID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrKitchenRecipeStepNotFound
	}
	return nil
}

func PdbDeleteKitchenRecipeStepsByRecipeID(ctx context.Context, recipeID uuid.UUID) error {
	_, err := db.PDB.Exec(ctx, `DELETE FROM kitchen_recipe_steps WHERE recipe_id = $1`, recipeID)
	return err
}

func BuildKitchenRecipeStepMetadata(imageKey string) json.RawMessage {
	if imageKey == "" {
		return json.RawMessage(`{}`)
	}
	return json.RawMessage(`{"legacy_image_key":` + string(mustMarshalJSONString(imageKey)) + `}`)
}

func mustMarshalJSONString(value string) []byte {
	data, _ := json.Marshal(value)
	return data
}
