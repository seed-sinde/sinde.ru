package media
import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"sinde.ru/db"
)
type Section string
type Collection string
const (
	SectionUsers   Section = "users"
	SectionKitchen Section = "kitchen"
	CollectionAvatars Collection = "avatars"
	CollectionRecipes Collection = "recipes"
)
type UploadTarget struct {
	Section    Section
	Collection Collection
	UserID     uuid.UUID
	RecipeID   *uuid.UUID
}
type StoredFile struct {
	ImageKey string
	ImageURL string
	FileHash string
	Reused   bool
}
var storageKeyPattern = regexp.MustCompile(`^[a-zA-Z0-9/_\.-]+$`)
func RootDir() string {
	if raw := strings.TrimSpace(os.Getenv("MEDIA_FILES_DIR")); raw != "" {
		return raw
	}
	return "./storage"
}
func MaxBytes() int64 {
	const defaultLimit int64 = 8 * 1024 * 1024
	raw := strings.TrimSpace(os.Getenv("MEDIA_IMAGE_MAX_BYTES"))
	if raw == "" {
		return defaultLimit
	}
	n, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || n <= 0 {
		return defaultLimit
	}
	return n
}
func NormalizeStorageKey(raw string) string {
	key := strings.TrimSpace(strings.ReplaceAll(raw, "\\", "/"))
	key = strings.TrimPrefix(key, "/")
	if key == "" {
		return ""
	}
	key = filepath.ToSlash(filepath.Clean(key))
	if key == "." || strings.HasPrefix(key, "../") || strings.Contains(key, "/../") || strings.HasPrefix(key, "/") {
		return ""
	}
	if !storageKeyPattern.MatchString(key) {
		return ""
	}
	return key
}
func MediaProxyURL(storageKey string) string {
	key := NormalizeStorageKey(storageKey)
	if key == "" {
		return ""
	}
	return "/api/proxy/media/files/" + key
}
func AbsPath(storageKey string) (string, error) {
	root, err := filepath.Abs(RootDir())
	if err != nil {
		return "", err
	}
	key := NormalizeStorageKey(storageKey)
	if key == "" {
		return "", errors.New("invalid media key")
	}
	target := filepath.Join(root, filepath.FromSlash(key))
	rel, err := filepath.Rel(root, target)
	if err != nil {
		return "", err
	}
	rel = filepath.ToSlash(rel)
	if rel == ".." || strings.HasPrefix(rel, "../") {
		return "", errors.New("invalid media path")
	}
	return target, nil
}
func ResolveFilePath(storageKey string) (string, error) {
	candidates, err := pathCandidates(storageKey)
	if err != nil {
		return "", err
	}
	for _, candidate := range candidates {
		if _, statErr := os.Stat(candidate); statErr == nil {
			return candidate, nil
		}
	}
	return candidates[0], nil
}
func SaveUploadedFile(ctx context.Context, fileHeader *multipart.FileHeader, target UploadTarget) (*StoredFile, error) {
	if fileHeader == nil {
		return nil, errors.New("file is required")
	}
	if fileHeader.Size <= 0 {
		return nil, errors.New("file is empty")
	}
	maxBytes := MaxBytes()
	if fileHeader.Size > maxBytes {
		return nil, fmt.Errorf("file too large, max %d bytes", maxBytes)
	}
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()
	limited := io.LimitReader(file, maxBytes+1)
	data, err := io.ReadAll(limited)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errors.New("file is empty")
	}
	if int64(len(data)) > maxBytes {
		return nil, fmt.Errorf("file too large, max %d bytes", maxBytes)
	}
	normalizedData, ext, err := normalizeImageForStorage(data)
	if err != nil {
		return nil, errors.New("unsupported image format")
	}
	if int64(len(normalizedData)) > maxBytes {
		return nil, fmt.Errorf("file too large after optimization, max %d bytes", maxBytes)
	}
	fileHash := sha256Hex(normalizedData)
	if existingKey, err := findByHash(ctx, fileHash, target.Section, target.Collection); err != nil {
		return nil, err
	} else if existingKey != "" {
		exists, existsErr := ObjectExists(ctx, existingKey)
		if existsErr == nil && exists {
			return &StoredFile{
				ImageKey: existingKey,
				ImageURL: MediaProxyURL(existingKey),
				FileHash: fileHash,
				Reused:   true,
			}, nil
		}
		if existsErr == nil && !exists {
			_ = deleteMetadataByHash(ctx, fileHash, target.Section, target.Collection)
		}
	}
	imageKey, err := buildStorageKey(target, ext)
	if err != nil {
		return nil, err
	}
	if err := PutObjectBytes(ctx, imageKey, normalizedData, "image/webp"); err != nil {
		return nil, err
	}
	storedKey, err := upsertMetadata(ctx, imageKey, fileHash, target, ext, int64(len(normalizedData)))
	if err != nil {
		return nil, err
	}
	if storedKey != imageKey {
		_ = DeleteObject(ctx, imageKey)
		return &StoredFile{
			ImageKey: storedKey,
			ImageURL: MediaProxyURL(storedKey),
			FileHash: fileHash,
			Reused:   true,
		}, nil
	}
	return &StoredFile{
		ImageKey: imageKey,
		ImageURL: MediaProxyURL(imageKey),
		FileHash: fileHash,
		Reused:   false,
	}, nil
}
func DeleteIfUnreferenced(ctx context.Context, storageKey string) error {
	key := NormalizeStorageKey(storageKey)
	if key == "" {
		return nil
	}
	referenced, err := hasReferences(ctx, key)
	if err != nil {
		return err
	}
	if referenced {
		return nil
	}
	if err := DeleteObject(ctx, key); err != nil && !os.IsNotExist(err) {
		return err
	}
	_, err = db.PDB.Exec(ctx, `DELETE FROM media_files WHERE image_key = $1`, key)
	return err
}
func buildStorageKey(target UploadTarget, ext string) (string, error) {
	switch {
	case target.Section == SectionUsers && target.Collection == CollectionAvatars:
		if target.UserID == uuid.Nil {
			return "", errors.New("user id is required")
		}
		return fmt.Sprintf("users/avatars/%s/%s%s", target.UserID.String(), uuid.NewString(), ext), nil
	case target.Section == SectionKitchen && target.Collection == CollectionRecipes:
		if target.RecipeID == nil || *target.RecipeID == uuid.Nil {
			if target.UserID == uuid.Nil {
				return "", errors.New("user id is required")
			}
			return fmt.Sprintf("kitchen/recipes/drafts/%s/%s%s", target.UserID.String(), uuid.NewString(), ext), nil
		}
		return fmt.Sprintf("kitchen/recipes/%s/%s%s", target.RecipeID.String(), uuid.NewString(), ext), nil
	default:
		return "", errors.New("media target is invalid")
	}
}
func sha256Hex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
func findByHash(ctx context.Context, fileHash string, section Section, collection Collection) (string, error) {
	var imageKey string
	err := db.PDB.QueryRow(ctx, `
		SELECT image_key
		FROM media_files
		WHERE file_hash = $1
			AND storage_section = $2
			AND storage_collection = $3
		LIMIT 1
	`, fileHash, string(section), string(collection)).Scan(&imageKey)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return NormalizeStorageKey(imageKey), nil
}
func deleteMetadataByHash(ctx context.Context, fileHash string, section Section, collection Collection) error {
	_, err := db.PDB.Exec(ctx, `
		DELETE FROM media_files
		WHERE file_hash = $1
			AND storage_section = $2
			AND storage_collection = $3
	`, fileHash, string(section), string(collection))
	return err
}
func upsertMetadata(ctx context.Context, imageKey string, fileHash string, target UploadTarget, ext string, byteSize int64) (string, error) {
	var ownerUserID any
	if target.UserID != uuid.Nil {
		ownerUserID = target.UserID
	}
	var recipeID any
	if target.RecipeID != nil && *target.RecipeID != uuid.Nil {
		recipeID = *target.RecipeID
	}
	var mediaKind string
	switch {
	case target.Section == SectionUsers && target.Collection == CollectionAvatars:
		mediaKind = "avatar"
	case target.Section == SectionKitchen && target.Collection == CollectionRecipes:
		mediaKind = "recipe"
	default:
		return "", errors.New("media target is invalid")
	}
	var storedKey string
	err := db.PDB.QueryRow(ctx, `
		INSERT INTO media_files (
			image_key,
			file_hash,
			storage_section,
			storage_collection,
			media_kind,
			owner_user_id,
			recipe_id,
			content_ext,
			byte_size
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (file_hash, storage_section, storage_collection) DO UPDATE
		SET updated_at = now()
		RETURNING image_key
	`, imageKey, fileHash, string(target.Section), string(target.Collection), mediaKind, ownerUserID, recipeID, ext, byteSize).Scan(&storedKey)
	if err != nil {
		return "", err
	}
	return NormalizeStorageKey(storedKey), nil
}
func hasReferences(ctx context.Context, storageKey string) (bool, error) {
	var referenced bool
	err := db.PDB.QueryRow(ctx, `
		SELECT
			EXISTS (
				SELECT 1
				FROM users
				WHERE
					profile::text LIKE '%"' || $1 || '"%'
					OR settings::text LIKE '%"' || $1 || '"%'
			)
			OR EXISTS (
				SELECT 1
				FROM kitchen_recipes
				WHERE
					cover_image_key = $1
					OR EXISTS (
						SELECT 1
						FROM jsonb_array_elements(steps) AS step
						WHERE step->>'image_key' = $1
					)
			)
	`, storageKey).Scan(&referenced)
	return referenced, err
}
func removeFileAndEmptyParents(storageKey string) error {
	path, err := ResolveFilePath(storageKey)
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil {
		return err
	}
	root, err := filepath.Abs(RootDir())
	if err != nil {
		return err
	}
	dir := filepath.Dir(path)
	for dir != root && dir != "." && dir != string(filepath.Separator) {
		if removeErr := os.Remove(dir); removeErr != nil {
			if errors.Is(removeErr, os.ErrNotExist) {
				dir = filepath.Dir(dir)
				continue
			}
			break
		}
		dir = filepath.Dir(dir)
	}
	return nil
}
func pathCandidates(storageKey string) ([]string, error) {
	key := NormalizeStorageKey(storageKey)
	if key == "" {
		return nil, errors.New("invalid media key")
	}
	path, err := AbsPath(key)
	if err != nil {
		return nil, err
	}
	return []string{path}, nil
}
