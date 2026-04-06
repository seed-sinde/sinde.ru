package media

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"sinde.ru/db"
	"strconv"
	"strings"
)

type Section string
type Collection string

const (
	SectionUsers      Section    = "users"
	SectionKitchen    Section    = "kitchen"
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
	objectID, err := findObjectIDByStorageKey(ctx, key)
	if err != nil {
		return err
	}
	if err := DeleteObject(ctx, key); err != nil && !os.IsNotExist(err) {
		return err
	}
	if objectID == uuid.Nil {
		return nil
	}
	_, err = db.PDB.Exec(ctx, `DELETE FROM storage_objects WHERE object_id = $1`, objectID)
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
func mediaFamilyForTarget(target UploadTarget) (string, error) {
	switch {
	case target.Section == SectionUsers && target.Collection == CollectionAvatars:
		return "users-avatars", nil
	case target.Section == SectionKitchen && target.Collection == CollectionRecipes:
		return "kitchen-recipes", nil
	default:
		return "", errors.New("media target is invalid")
	}
}
func findByHash(ctx context.Context, fileHash string, section Section, collection Collection) (string, error) {
	mediaFamily, err := mediaFamilyForTarget(UploadTarget{
		Section:    section,
		Collection: collection,
	})
	if err != nil {
		return "", err
	}
	var storageKey string
	err = db.PDB.QueryRow(ctx, `
		SELECT storage_key
		FROM storage_objects
		WHERE file_hash = $1
			AND media_family = $2
		ORDER BY created_at DESC
		LIMIT 1
	`, fileHash, mediaFamily).Scan(&storageKey)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return NormalizeStorageKey(storageKey), nil
}
func deleteMetadataByHash(ctx context.Context, fileHash string, section Section, collection Collection) error {
	mediaFamily, err := mediaFamilyForTarget(UploadTarget{
		Section:    section,
		Collection: collection,
	})
	if err != nil {
		return err
	}
	_, err = db.PDB.Exec(ctx, `
		DELETE FROM storage_objects
		WHERE file_hash = $1
			AND media_family = $2
	`, fileHash, mediaFamily)
	return err
}
func upsertMetadata(ctx context.Context, imageKey string, fileHash string, target UploadTarget, ext string, byteSize int64) (string, error) {
	mediaFamily, err := mediaFamilyForTarget(target)
	if err != nil {
		return "", err
	}
	var storedKey string
	err = db.PDB.QueryRow(ctx, `
		INSERT INTO storage_objects (
			storage_key,
			file_hash,
			bucket_name,
			media_family,
			content_type,
			byte_size,
			source_kind,
			metadata
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			'runtime',
			jsonb_build_object(
				'storage_ext', $7,
				'section', $8,
				'collection', $9,
				'user_id', NULLIF($10, ''),
				'recipe_id', NULLIF($11, '')
			)
		)
		ON CONFLICT (storage_key) DO UPDATE
		SET
			file_hash = EXCLUDED.file_hash,
			bucket_name = EXCLUDED.bucket_name,
			media_family = EXCLUDED.media_family,
			content_type = EXCLUDED.content_type,
			byte_size = EXCLUDED.byte_size,
			source_kind = EXCLUDED.source_kind,
			metadata = EXCLUDED.metadata,
			updated_at = now()
		RETURNING storage_key
	`,
		imageKey,
		fileHash,
		MinIOBucket(),
		mediaFamily,
		detectContentType(imageKey, nil),
		byteSize,
		ext,
		string(target.Section),
		string(target.Collection),
		target.UserID.String(),
		func() string {
			if target.RecipeID == nil {
				return ""
			}
			return target.RecipeID.String()
		}(),
	).Scan(&storedKey)
	if err != nil {
		return "", err
	}
	return NormalizeStorageKey(storedKey), nil
}
func hasReferences(ctx context.Context, storageKey string) (bool, error) {
	objectID, err := findObjectIDByStorageKey(ctx, storageKey)
	if err != nil {
		return false, err
	}
	if objectID != uuid.Nil {
		var referencedByUsage bool
		if err := db.PDB.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1
				FROM storage_object_usages
				WHERE object_id = $1
			)
		`, objectID).Scan(&referencedByUsage); err != nil {
			return false, err
		}
		if referencedByUsage {
			return true, nil
		}
	}
	var referenced bool
	err = db.PDB.QueryRow(ctx, `
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
func findObjectIDByStorageKey(ctx context.Context, storageKey string) (uuid.UUID, error) {
	var objectID uuid.UUID
	err := db.PDB.QueryRow(ctx, `
		SELECT object_id
		FROM storage_objects
		WHERE storage_key = $1
		LIMIT 1
	`, storageKey).Scan(&objectID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, nil
		}
		return uuid.Nil, err
	}
	return objectID, nil
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
