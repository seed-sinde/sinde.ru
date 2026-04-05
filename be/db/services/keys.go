package services
import (
	"context"
	"errors"
	"strings"
	"github.com/jackc/pgx/v5"
	"sinde.ru/db"
	"sinde.ru/internal/models"
)
// traits_k: id BIGINT PK, syn_id BIGINT FK -> key_syns, meta JSONB ...
func PdbGetKeyByID(ctx context.Context, id int64) (*models.Key, error) {
	const query = `
		SELECT k.id, k.syn_id, s.name, COALESCE(k.meta, '{}'::jsonb)::text
		FROM traits_k k
		JOIN key_syns s ON s.id = k.syn_id
		WHERE k.id = $1
	`
	var key models.Key
	err := db.PDB.QueryRow(ctx, query, id).Scan(&key.ID, &key.SynID, &key.Syn, &key.Meta)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &key, nil
}
func PdbGetCanonicalKeyBySyn(ctx context.Context, syn string) (*models.Key, error) {
	syn = strings.TrimSpace(syn)
	if syn == "" {
		return nil, nil
	}
	const query = `
		SELECT k.id, k.syn_id, s.name, COALESCE(k.meta, '{}'::jsonb)::text
		FROM traits_k k
		JOIN key_syns s ON s.id = k.syn_id
		WHERE lower(s.name) = lower($1)
		ORDER BY k.id
		LIMIT 1
	`
	var key models.Key
	err := db.PDB.QueryRow(ctx, query, syn).Scan(&key.ID, &key.SynID, &key.Syn, &key.Meta)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &key, nil
}
func PdbGetOrCreateCanonicalKeyBySyn(ctx context.Context, syn string) (*models.Key, error) {
	syn = strings.TrimSpace(syn)
	if syn == "" {
		return nil, nil
	}
	key, err := PdbGetCanonicalKeyBySyn(ctx, syn)
	if err != nil || key != nil {
		return key, err
	}
	const query = `
		WITH ensured_syn AS (
			INSERT INTO key_syns (name)
			VALUES ($1)
			ON CONFLICT (name) DO NOTHING
			RETURNING id
		),
		syn_row AS (
			SELECT id FROM ensured_syn
			UNION ALL
			SELECT id FROM key_syns WHERE name = $1
			LIMIT 1
		)
		INSERT INTO traits_k (syn_id, meta)
		VALUES ((SELECT id FROM syn_row), '{}'::jsonb)
		ON CONFLICT (syn_id, meta) DO NOTHING
	`
	if _, err := db.PDB.Exec(ctx, query, syn); err != nil {
		return nil, err
	}
	return PdbGetCanonicalKeyBySyn(ctx, syn)
}
// Обновление meta по ID.
func PDBUpdateKeyMetaByID(ctx context.Context, id int64, metaJSON string) (*models.Key, error) {
	const query = `
		UPDATE traits_k AS k
		SET meta = $1::jsonb
		FROM key_syns AS s
		WHERE k.id = $2
		  AND s.id = k.syn_id
		RETURNING k.id, k.syn_id, s.name, COALESCE(k.meta, '{}'::jsonb)::text
	`
	var key models.Key
	err := db.PDB.QueryRow(ctx, query, metaJSON, id).Scan(&key.ID, &key.SynID, &key.Syn, &key.Meta)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &key, nil
}
