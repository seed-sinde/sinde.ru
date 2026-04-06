package services

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"sinde.ru/db"
	"sinde.ru/internal/models"
)

var ErrStorageObjectNotFound = errors.New("storage object not found")

func normalizeJSONMetadata(raw json.RawMessage) json.RawMessage {
	trimmed := strings.TrimSpace(string(raw))
	if trimmed == "" || trimmed == "null" {
		return json.RawMessage(`{}`)
	}
	return raw
}

func scanStorageObject(row interface{ Scan(dest ...any) error }) (*models.StorageObject, error) {
	var item models.StorageObject
	var metadata []byte
	if err := row.Scan(
		&item.ObjectID,
		&item.StorageKey,
		&item.BucketName,
		&item.MediaFamily,
		&item.ContentType,
		&item.ByteSize,
		&item.FileHash,
		&item.SourceKind,
		&metadata,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrStorageObjectNotFound
		}
		return nil, err
	}
	item.Metadata = normalizeJSONMetadata(metadata)
	return &item, nil
}

func scanStorageObjectUsage(row interface{ Scan(dest ...any) error }) (*models.StorageObjectUsage, error) {
	var item models.StorageObjectUsage
	var usageMetadata []byte
	var object models.StorageObject
	var objectMetadata []byte
	if err := row.Scan(
		&item.UsageID,
		&item.ObjectID,
		&item.EntityType,
		&item.EntityID,
		&item.UsageType,
		&item.FieldName,
		&item.SortOrder,
		&item.IsPrimary,
		&usageMetadata,
		&item.CreatedAt,
		&object.ObjectID,
		&object.StorageKey,
		&object.BucketName,
		&object.MediaFamily,
		&object.ContentType,
		&object.ByteSize,
		&object.FileHash,
		&object.SourceKind,
		&objectMetadata,
		&object.CreatedAt,
		&object.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrStorageObjectNotFound
		}
		return nil, err
	}
	item.Metadata = normalizeJSONMetadata(usageMetadata)
	object.Metadata = normalizeJSONMetadata(objectMetadata)
	item.Object = &object
	return &item, nil
}

func PdbRegisterStorageObject(ctx context.Context, object *models.StorageObject) (*models.StorageObject, error) {
	if object == nil {
		return nil, errors.New("storage object is required")
	}
	row := db.PDB.QueryRow(ctx, `
		INSERT INTO storage_objects (
			storage_key,
			bucket_name,
			media_family,
			content_type,
			byte_size,
			file_hash,
			source_kind,
			metadata
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8::jsonb)
		ON CONFLICT (storage_key) DO UPDATE
		SET
			bucket_name = EXCLUDED.bucket_name,
			media_family = EXCLUDED.media_family,
			content_type = EXCLUDED.content_type,
			byte_size = EXCLUDED.byte_size,
			file_hash = EXCLUDED.file_hash,
			source_kind = EXCLUDED.source_kind,
			metadata = EXCLUDED.metadata,
			updated_at = now()
		RETURNING
			object_id,
			storage_key,
			bucket_name,
			media_family,
			content_type,
			byte_size,
			file_hash,
			source_kind,
			metadata,
			created_at,
			updated_at
	`,
		strings.TrimSpace(object.StorageKey),
		strings.TrimSpace(object.BucketName),
		strings.TrimSpace(object.MediaFamily),
		strings.TrimSpace(object.ContentType),
		object.ByteSize,
		strings.TrimSpace(object.FileHash),
		strings.TrimSpace(object.SourceKind),
		normalizeJSONMetadata(object.Metadata),
	)
	return scanStorageObject(row)
}

func PdbFindStorageObjectByStorageKey(ctx context.Context, storageKey string) (*models.StorageObject, error) {
	row := db.PDB.QueryRow(ctx, `
		SELECT
			object_id,
			storage_key,
			bucket_name,
			media_family,
			content_type,
			byte_size,
			file_hash,
			source_kind,
			metadata,
			created_at,
			updated_at
		FROM storage_objects
		WHERE storage_key = $1
		LIMIT 1
	`, strings.TrimSpace(storageKey))
	return scanStorageObject(row)
}

func PdbFindStorageObjectByHashAndFamily(ctx context.Context, fileHash, mediaFamily string) (*models.StorageObject, error) {
	row := db.PDB.QueryRow(ctx, `
		SELECT
			object_id,
			storage_key,
			bucket_name,
			media_family,
			content_type,
			byte_size,
			file_hash,
			source_kind,
			metadata,
			created_at,
			updated_at
		FROM storage_objects
		WHERE file_hash = $1 AND media_family = $2
		ORDER BY created_at DESC
		LIMIT 1
	`, strings.TrimSpace(fileHash), strings.TrimSpace(mediaFamily))
	return scanStorageObject(row)
}

func PdbDeleteStorageObjectByObjectID(ctx context.Context, objectID uuid.UUID) error {
	_, err := db.PDB.Exec(ctx, `DELETE FROM storage_objects WHERE object_id = $1`, objectID)
	return err
}

func PdbAttachStorageObjectUsage(ctx context.Context, usage *models.StorageObjectUsage) (*models.StorageObjectUsage, error) {
	if usage == nil {
		return nil, errors.New("storage object usage is required")
	}
	row := db.PDB.QueryRow(ctx, `
		INSERT INTO storage_object_usages (
			object_id,
			entity_type,
			entity_id,
			usage_type,
			field_name,
			sort_order,
			is_primary,
			metadata
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8::jsonb)
		ON CONFLICT (object_id, entity_type, entity_id, usage_type, field_name) DO UPDATE
		SET
			sort_order = EXCLUDED.sort_order,
			is_primary = EXCLUDED.is_primary,
			metadata = EXCLUDED.metadata
		RETURNING
			usage_id,
			object_id,
			entity_type,
			entity_id,
			usage_type,
			field_name,
			sort_order,
			is_primary,
			metadata,
			created_at,
			(
				SELECT object_id FROM storage_objects WHERE storage_objects.object_id = storage_object_usages.object_id
			),
			(
				SELECT storage_key FROM storage_objects WHERE storage_objects.object_id = storage_object_usages.object_id
			),
			(
				SELECT bucket_name FROM storage_objects WHERE storage_objects.object_id = storage_object_usages.object_id
			),
			(
				SELECT media_family FROM storage_objects WHERE storage_objects.object_id = storage_object_usages.object_id
			),
			(
				SELECT content_type FROM storage_objects WHERE storage_objects.object_id = storage_object_usages.object_id
			),
			(
				SELECT byte_size FROM storage_objects WHERE storage_objects.object_id = storage_object_usages.object_id
			),
			(
				SELECT file_hash FROM storage_objects WHERE storage_objects.object_id = storage_object_usages.object_id
			),
			(
				SELECT source_kind FROM storage_objects WHERE storage_objects.object_id = storage_object_usages.object_id
			),
			(
				SELECT metadata FROM storage_objects WHERE storage_objects.object_id = storage_object_usages.object_id
			),
			(
				SELECT created_at FROM storage_objects WHERE storage_objects.object_id = storage_object_usages.object_id
			),
			(
				SELECT updated_at FROM storage_objects WHERE storage_objects.object_id = storage_object_usages.object_id
			)
	`,
		usage.ObjectID,
		strings.TrimSpace(usage.EntityType),
		strings.TrimSpace(usage.EntityID),
		strings.TrimSpace(usage.UsageType),
		strings.TrimSpace(usage.FieldName),
		usage.SortOrder,
		usage.IsPrimary,
		normalizeJSONMetadata(usage.Metadata),
	)
	return scanStorageObjectUsage(row)
}

func listStorageObjectUsages(ctx context.Context, query string, args ...any) ([]*models.StorageObjectUsage, error) {
	rows, err := db.PDB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*models.StorageObjectUsage, 0, 8)
	for rows.Next() {
		item, scanErr := scanStorageObjectUsage(rows)
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

func PdbListStorageObjectUsagesByEntity(ctx context.Context, entityType, entityID string) ([]*models.StorageObjectUsage, error) {
	return listStorageObjectUsages(ctx, `
		SELECT
			u.usage_id,
			u.object_id,
			u.entity_type,
			u.entity_id,
			u.usage_type,
			u.field_name,
			u.sort_order,
			u.is_primary,
			u.metadata,
			u.created_at,
			o.object_id,
			o.storage_key,
			o.bucket_name,
			o.media_family,
			o.content_type,
			o.byte_size,
			o.file_hash,
			o.source_kind,
			o.metadata,
			o.created_at,
			o.updated_at
		FROM storage_object_usages u
		JOIN storage_objects o ON o.object_id = u.object_id
		WHERE u.entity_type = $1 AND u.entity_id = $2
		ORDER BY u.is_primary DESC, u.sort_order ASC, u.created_at ASC
	`, strings.TrimSpace(entityType), strings.TrimSpace(entityID))
}

func PdbListStorageObjectUsagesByObjectID(ctx context.Context, objectID uuid.UUID) ([]*models.StorageObjectUsage, error) {
	return listStorageObjectUsages(ctx, `
		SELECT
			u.usage_id,
			u.object_id,
			u.entity_type,
			u.entity_id,
			u.usage_type,
			u.field_name,
			u.sort_order,
			u.is_primary,
			u.metadata,
			u.created_at,
			o.object_id,
			o.storage_key,
			o.bucket_name,
			o.media_family,
			o.content_type,
			o.byte_size,
			o.file_hash,
			o.source_kind,
			o.metadata,
			o.created_at,
			o.updated_at
		FROM storage_object_usages u
		JOIN storage_objects o ON o.object_id = u.object_id
		WHERE u.object_id = $1
		ORDER BY u.entity_type ASC, u.entity_id ASC, u.sort_order ASC, u.created_at ASC
	`, objectID)
}

func PdbDeleteStorageObjectUsagesByEntity(ctx context.Context, entityType, entityID string) error {
	_, err := db.PDB.Exec(ctx, `
		DELETE FROM storage_object_usages
		WHERE entity_type = $1 AND entity_id = $2
	`, strings.TrimSpace(entityType), strings.TrimSpace(entityID))
	return err
}

func PdbDeleteStorageObjectUsagesBySelector(ctx context.Context, entityType, entityID, usageType, fieldName string) error {
	_, err := db.PDB.Exec(ctx, `
		DELETE FROM storage_object_usages
		WHERE entity_type = $1 AND entity_id = $2 AND usage_type = $3 AND field_name = $4
	`, strings.TrimSpace(entityType), strings.TrimSpace(entityID), strings.TrimSpace(usageType), strings.TrimSpace(fieldName))
	return err
}
