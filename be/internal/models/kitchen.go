package models

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type KitchenIngredient struct {
	Name   string `json:"name"`
	Amount string `json:"amount,omitempty"`
	Unit   string `json:"unit,omitempty"`
	Note   string `json:"note,omitempty"`
}
type KitchenStep struct {
	Order    int    `json:"order"`
	Text     string `json:"text"`
	ImageKey string `json:"image_key,omitempty"`
}
type KitchenRecipeStep struct {
	StepID      uuid.UUID       `json:"step_id"`
	RecipeID    uuid.UUID       `json:"recipe_id"`
	StepOrder   int             `json:"step_order"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Metadata    json.RawMessage `json:"metadata,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}
type KitchenRecipe struct {
	ID            uuid.UUID  `json:"id"`
	OwnerUserID   *uuid.UUID `json:"owner_user_id,omitempty"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	CoverImageKey string     `json:"cover_image_key,omitempty"`
	Kcal          int        `json:"kcal"`
	PrepMinutes   int        `json:"prep_minutes"`
	CookMinutes   int        `json:"cook_minutes"`
	Servings      int        `json:"servings"`
	// Difficulty stores a catalog code from kitchen_filter_options(kind='difficulty'), not a fixed numeric level.
	Difficulty        string              `json:"difficulty"`
	MealType          string              `json:"meal_type"`
	CookingMethod     string              `json:"cooking_method,omitempty"`
	Cuisine           string              `json:"cuisine"`
	DietType          string              `json:"diet_type,omitempty"`
	Ingredients       []KitchenIngredient `json:"ingredients"`
	IngredientsSearch []string            `json:"ingredients_search,omitempty"`
	Steps             []KitchenStep       `json:"steps"`
	Tags              []string            `json:"tags"`
	IsPublic          bool                `json:"is_public"`
	ModerationStatus  string              `json:"moderation_status"`
	ModeratedByUserID *uuid.UUID          `json:"moderated_by_user_id,omitempty"`
	ModeratedAt       *time.Time          `json:"moderated_at,omitempty"`
	ModerationNote    string              `json:"moderation_note,omitempty"`
	CreatedAt         time.Time           `json:"created_at"`
	UpdatedAt         time.Time           `json:"updated_at"`
}
type KitchenUserIngredient struct {
	IngredientID   uuid.UUID `json:"ingredient_id"`
	UserID         uuid.UUID `json:"user_id"`
	Name           string    `json:"name"`
	Category       string    `json:"category"`
	NormalizedName string    `json:"normalized_name"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
type KitchenFavoriteIngredient struct {
	IngredientID   int64     `json:"ingredient_id"`
	UserID         uuid.UUID `json:"user_id"`
	ListType       string    `json:"list_type"`
	Name           string    `json:"name"`
	Category       string    `json:"category"`
	NormalizedName string    `json:"normalized_name"`
	CreatedAt      time.Time `json:"created_at"`
}
type KitchenFavoriteRecipe struct {
	UserID    uuid.UUID `json:"user_id"`
	RecipeID  uuid.UUID `json:"recipe_id"`
	CreatedAt time.Time `json:"created_at"`
}
type KitchenCategory struct {
	CategoryID  int64  `json:"category_id"`
	Slug        string `json:"slug"`
	Label       string `json:"label"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}
type KitchenCatalogIngredient struct {
	IngredientID   int64    `json:"ingredient_id"`
	Name           string   `json:"name"`
	NormalizedName string   `json:"normalized_name"`
	CategoryID     int64    `json:"category_id"`
	CategoryLabel  string   `json:"category_label"`
	Description    string   `json:"description"`
	ProteinG       *float64 `json:"protein_g,omitempty"`
	FatG           *float64 `json:"fat_g,omitempty"`
	CarbsG         *float64 `json:"carbs_g,omitempty"`
	Kcal           *float64 `json:"kcal,omitempty"`
	SortOrder      int      `json:"sort_order"`
}
type KitchenFilterOption struct {
	OptionID    int64  `json:"option_id"`
	Kind        string `json:"kind"`
	Code        string `json:"code"`
	Label       string `json:"label"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}
