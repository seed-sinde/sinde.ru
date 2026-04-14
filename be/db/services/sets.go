package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"sinde.ru/db"
	"sinde.ru/internal/models"
	"sinde.ru/utils"
)

func PdbSetExists(ctx context.Context, setUUID uuid.UUID) (bool, error) {
	var exists bool
	err := db.PDB.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM sets WHERE s_uuid = $1)`, setUUID).Scan(&exists)
	return exists, err
}

// PdbInsertSet вставляет набор без изменения порядка потомков.
//
// Параметры:
//
//	ctx — контекст операции.
//	set — набор для сохранения.
//
// Ограничения:
//
//	Допустимы только [t_uuid, t_uuid] и [s_uuid, t_uuid] (проверяется триггером БД).
//
// Возвращает:
//
//	Ошибку БД либо nil.
func PdbInsertSet(ctx context.Context, set *models.Set) error {
	defer utils.Benchmark("PdbInsertSet")()
	if len(set.SChilds) != 2 {
		return fmt.Errorf("invalid childs len for %s: %v", set.SUUID, set.SChilds)
	}
	query := `
		INSERT INTO sets (s_uuid, s_childs)
		VALUES ($1, ARRAY[$2::uuid, $3::uuid])
		ON CONFLICT (s_uuid) DO NOTHING;
	`
	_, err := db.PDB.Exec(ctx, query, set.SUUID, set.SChilds[0], set.SChilds[1])
	if err != nil {
		return fmt.Errorf("failed to insert set %s: %w", set.SUUID, err)
	}
	return nil
}
func PdbGetSetByChildren(ctx context.Context, children []uuid.UUID) (*models.Set, error) {
	if len(children) != 2 {
		return nil, fmt.Errorf("invalid childs len: %v", children)
	}
	const query = `
		SELECT s_uuid, s_childs
		FROM sets
		WHERE s_childs = ARRAY[$1::uuid, $2::uuid]
		LIMIT 1
	`
	var setItem models.Set
	var raw []pgtype.UUID
	err := db.PDB.QueryRow(ctx, query, children[0], children[1]).Scan(&setItem.SUUID, &raw)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	for _, item := range raw {
		if !item.Valid {
			continue
		}
		child, err := uuid.FromBytes(item.Bytes[:])
		if err != nil {
			continue
		}
		setItem.SChilds = append(setItem.SChilds, child)
	}
	return &setItem, nil
}
func PdbGetOrCreateSet(ctx context.Context, children []uuid.UUID) (*models.Set, error) {
	setItem, err := PdbGetSetByChildren(ctx, children)
	if err != nil || setItem != nil {
		return setItem, err
	}
	insertedUUID := uuid.New()
	const query = `
		INSERT INTO sets (s_uuid, s_childs)
		VALUES ($1, ARRAY[$2::uuid, $3::uuid])
		ON CONFLICT (s_childs) DO NOTHING
		RETURNING s_uuid, s_childs
	`
	var inserted models.Set
	var raw []pgtype.UUID
	err = db.PDB.QueryRow(ctx, query, insertedUUID, children[0], children[1]).Scan(&inserted.SUUID, &raw)
	if err == nil {
		for _, item := range raw {
			if !item.Valid {
				continue
			}
			child, err := uuid.FromBytes(item.Bytes[:])
			if err != nil {
				continue
			}
			inserted.SChilds = append(inserted.SChilds, child)
		}
		return &inserted, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	return PdbGetSetByChildren(ctx, children)
}

// Функция PdbGetSetTraits
func PdbGetSetTraits(ctx context.Context, setUUID uuid.UUID) ([]models.Trait, error) {
	defer utils.Benchmark("PdbGetSetTraits")()
	query := `
		WITH RECURSIVE set_hierarchy AS (
			SELECT unnest(s_childs) as child_uuid
			FROM sets
			WHERE s_uuid = $1
			UNION ALL
			SELECT unnest(s.s_childs)
			FROM sets s
			INNER JOIN set_hierarchy sh ON s.s_uuid = sh.child_uuid
		)
		SELECT v.t_uuid, v.t_key, v.t_value
		FROM traits_v v
		WHERE v.t_uuid IN (SELECT child_uuid FROM set_hierarchy);
	`
	rows, err := db.PDB.Query(ctx, query, setUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute recursive query for set %s: %w", setUUID, err)
	}
	defer rows.Close()
	result, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (models.Trait, error) {
		var t models.Trait
		err := row.Scan(&t.TUUID, &t.TKey, &t.TValue)
		return t, err
	})
	if err != nil {
		return nil, fmt.Errorf("failed to scan traits for set %s: %w", setUUID, err)
	}
	if len(result) == 0 {
		var exists bool
		err := db.PDB.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM sets WHERE s_uuid = $1)", setUUID).Scan(&exists)
		if err != nil {
			return nil, fmt.Errorf("failed to check existence for set %s: %w", setUUID, err)
		}
		if !exists {
			return nil, errors.New("not found")
		}
	}
	return result, nil
}
