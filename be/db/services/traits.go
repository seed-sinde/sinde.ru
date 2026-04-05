package services
import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"sinde.ru/db"
	"sinde.ru/internal/models"
	"sinde.ru/internal/store"
	"sinde.ru/utils"
)
func PdbTraitExists(ctx context.Context, id uuid.UUID) (bool, error) {
	var exists bool
	err := db.PDB.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM traits_v WHERE t_uuid = $1)`, id).Scan(&exists)
	return exists, err
}
func PdbTraitTargetExists(ctx context.Context, id uuid.UUID) (bool, error) {
	const query = `
		SELECT EXISTS(
			SELECT 1 FROM traits_v WHERE t_uuid = $1
			UNION ALL
			SELECT 1 FROM sets WHERE s_uuid = $1
		)
	`
	var exists bool
	err := db.PDB.QueryRow(ctx, query, id).Scan(&exists)
	return exists, err
}
// Функция PdbInsertTrait сохраняет особенность в базе данных.
// Если пара (t_key, t_value) уже существует, из-за уникального ограничения это no-op.
func PdbInsertTrait(conn *pgxpool.Conn, trait *models.Trait) error {
	defer utils.Benchmark("PdbInsertTrait")()
	const query = `
		INSERT INTO traits_v (t_uuid, t_key, t_value)
		VALUES ($1, $2, $3)
		ON CONFLICT ON CONSTRAINT unique_t_key_value DO NOTHING;
	`
	_, err := conn.Exec(context.Background(), query, trait.TUUID, trait.TKey, trait.TValue)
	return err
}
// Функция PdbInsertTraitKey сохраняет ключ (id, syn, meta).
// Если id уже существует, это no-op.
func PdbInsertTraitKey(conn *pgxpool.Conn, key *models.Key) error {
	defer utils.Benchmark("PdbInsertTraitKey")()
	const query = `
		WITH syn_row AS (
			INSERT INTO key_syns (name)
			VALUES ($2)
			ON CONFLICT (name) DO NOTHING
			RETURNING id
		),
		syn_id AS (
			SELECT id FROM syn_row
			UNION ALL
			SELECT id FROM key_syns WHERE name = $2
			LIMIT 1
		)
		INSERT INTO traits_k (id, syn_id, meta)
		VALUES ($1, (SELECT id FROM syn_id), $3)
		ON CONFLICT (id) DO UPDATE SET syn_id = EXCLUDED.syn_id, meta = EXCLUDED.meta
		RETURNING syn_id;
	`
	var synID int64
	if err := conn.QueryRow(context.Background(), query, key.ID, key.Syn, key.Meta).Scan(&synID); err != nil {
		return err
	}
	if _, err := conn.Exec(context.Background(), `
		SELECT setval(
			pg_get_serial_sequence('traits_k', 'id'),
			GREATEST(COALESCE((SELECT MAX(id) FROM traits_k), 0), $1),
			TRUE
		)
	`, key.ID); err != nil {
		return err
	}
	key.SynID = synID
	return nil
}
// Функция PdbGetTraitByUUID возвращает особенность по её UUID.
func PdbGetTraitByUUID(uuid string) (*models.Trait, error) {
	defer utils.Benchmark("PdbGetTraitByUUID")()
	const query = `
		SELECT v.t_uuid, v.t_key, v.t_value
		FROM traits_v v
		WHERE v.t_uuid = $1;
	`
	var trait models.Trait
	err := db.PDB.QueryRow(context.Background(), query, uuid).
		Scan(&trait.TUUID, &trait.TKey, &trait.TValue)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("not found")
		}
		return nil, err
	}
	return &trait, nil
}
func PdbGetTraitByKeyValue(ctx context.Context, keyID int64, value string) (*models.Trait, error) {
	const query = `
		SELECT t_uuid, t_key, t_value
		FROM traits_v
		WHERE t_key = $1 AND t_value = $2
	`
	var trait models.Trait
	err := db.PDB.QueryRow(ctx, query, keyID, value).Scan(&trait.TUUID, &trait.TKey, &trait.TValue)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &trait, nil
}
func PdbGetOrCreateTraitByKeyValue(ctx context.Context, keyID int64, value string) (*models.Trait, error) {
	trait, err := PdbGetTraitByKeyValue(ctx, keyID, value)
	if err != nil || trait != nil {
		return trait, err
	}
	insertedUUID := uuid.New()
	const query = `
		INSERT INTO traits_v (t_uuid, t_key, t_value)
		VALUES ($1, $2, $3)
		ON CONFLICT ON CONSTRAINT unique_t_key_value DO NOTHING
		RETURNING t_uuid, t_key, t_value
	`
	var inserted models.Trait
	err = db.PDB.QueryRow(ctx, query, insertedUUID, keyID, value).Scan(&inserted.TUUID, &inserted.TKey, &inserted.TValue)
	if err == nil {
		return &inserted, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	return PdbGetTraitByKeyValue(ctx, keyID, value)
}
// Функция PdbLoadAllData загружает все сущности {ключи, особенности, наборы} из Postgres
// и инициализирует хранилище в памяти через store.LoadInitialData.
func PdbLoadAllData() error {
	defer utils.Benchmark("PdbLoadAllData")()
	ctx := context.Background()
	// --- Ключи ---
	var keys []*models.Key
	keyRows, err := db.PDB.Query(ctx, `
		SELECT k.id, k.syn_id, s.name, COALESCE(k.meta, '{}')
		FROM traits_k k
		JOIN key_syns s ON s.id = k.syn_id
	`)
	if err != nil {
		return err
	}
	defer keyRows.Close()
	for keyRows.Next() {
		var k models.Key
		if err := keyRows.Scan(&k.ID, &k.SynID, &k.Syn, &k.Meta); err != nil {
			return err
		}
		kk := k
		keys = append(keys, &kk)
	}
	// --- Особенности ---
	var traits []*models.Trait
	traitRows, err := db.PDB.Query(ctx, `SELECT t_uuid, t_key, t_value FROM traits_v`)
	if err != nil {
		return err
	}
	defer traitRows.Close()
	for traitRows.Next() {
		var t models.Trait
		if err := traitRows.Scan(&t.TUUID, &t.TKey, &t.TValue); err != nil {
			return err
		}
		tt := t
		traits = append(traits, &tt)
	}
	// --- Наборы ---
	var sets []*models.Set
	setRows, err := db.PDB.Query(ctx, `SELECT s_uuid, s_childs FROM sets`)
	if err != nil {
		return err
	}
	defer setRows.Close()
	for setRows.Next() {
		var s models.Set
		var raw []pgtype.UUID
		if err := setRows.Scan(&s.SUUID, &raw); err != nil {
			return err
		}
		for _, el := range raw {
			if !el.Valid {
				continue
			}
			child, err := uuid.FromBytes(el.Bytes[:])
			if err != nil {
				continue
			}
			s.SChilds = append(s.SChilds, child)
		}
		ss := s
		sets = append(sets, &ss)
	}
	// Атомарно инициализируем хранилище в памяти.
	store.LoadInitialData(traits, keys, sets)
	return nil
}
