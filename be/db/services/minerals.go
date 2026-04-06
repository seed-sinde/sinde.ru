package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"sinde.ru/db"
	"sinde.ru/internal/models"
	"strconv"
	"strings"
)

var ErrMineralNotFound = errors.New("mineral not found")
var allChemicalElementSymbols = []string{
	"H", "He", "Li", "Be", "B", "C", "N", "O", "F", "Ne",
	"Na", "Mg", "Al", "Si", "P", "S", "Cl", "Ar",
	"K", "Ca", "Sc", "Ti", "V", "Cr", "Mn", "Fe", "Co", "Ni", "Cu", "Zn",
	"Ga", "Ge", "As", "Se", "Br", "Kr",
	"Rb", "Sr", "Y", "Zr", "Nb", "Mo", "Tc", "Ru", "Rh", "Pd", "Ag", "Cd",
	"In", "Sn", "Sb", "Te", "I", "Xe",
	"Cs", "Ba", "La", "Ce", "Pr", "Nd", "Pm", "Sm", "Eu", "Gd", "Tb", "Dy",
	"Ho", "Er", "Tm", "Yb", "Lu",
	"Hf", "Ta", "W", "Re", "Os", "Ir", "Pt", "Au", "Hg",
	"Tl", "Pb", "Bi", "Po", "At", "Rn",
	"Fr", "Ra", "Ac", "Th", "Pa", "U", "Np", "Pu", "Am", "Cm", "Bk", "Cf",
	"Es", "Fm", "Md", "No", "Lr",
	"Rf", "Db", "Sg", "Bh", "Hs", "Mt", "Ds", "Rg", "Cn",
	"Nh", "Fl", "Mc", "Lv", "Ts", "Og",
}

type mineralScanner interface {
	Scan(dest ...any) error
}

func nullableStringPtr(value pgtype.Text) *string {
	if !value.Valid {
		return nil
	}
	text := strings.TrimSpace(value.String)
	if text == "" {
		return nil
	}
	return &text
}
func nullableIntPtr(value pgtype.Int4) *int {
	if !value.Valid {
		return nil
	}
	number := int(value.Int32)
	return &number
}
func nullableFloatPtr(value pgtype.Numeric) (*float64, error) {
	if !value.Valid {
		return nil, nil
	}
	floatValue, err := value.Float64Value()
	if err != nil {
		return nil, err
	}
	if !floatValue.Valid {
		return nil, nil
	}
	number := floatValue.Float64
	return &number, nil
}
func nonEmptyStrings(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	out := make([]string, 0, len(values))
	for _, value := range values {
		text := strings.TrimSpace(value)
		if text == "" {
			continue
		}
		out = append(out, text)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}
func buildMineralsSearchPattern(raw string) string {
	search := strings.TrimSpace(raw)
	if search == "" {
		return ""
	}
	return "%" + search + "%"
}
func scanMineral(row mineralScanner) (*models.Mineral, error) {
	var item models.Mineral
	var valenceChemistry pgtype.Text
	var imaChemistry pgtype.Text
	var imaNumber pgtype.Text
	var imaMineralSymbol pgtype.Text
	var imaStatus pgtype.Text
	var structuralGroupname pgtype.Text
	var statusNotes pgtype.Text
	var countryOfTypeLocality pgtype.Text
	var yearFirstPublished pgtype.Int4
	var oldestKnownAgeMa pgtype.Numeric
	if err := row.Scan(
		&item.ID,
		&item.DatabaseID,
		&item.MineralName,
		&item.MineralNamePlain,
		&item.MineralNameSearch,
		&valenceChemistry,
		&imaChemistry,
		&item.ChemistryElements,
		&imaNumber,
		&item.RRUFFIDs,
		&countryOfTypeLocality,
		&yearFirstPublished,
		&imaStatus,
		&structuralGroupname,
		&statusNotes,
		&item.CrystalSystems,
		&item.SpaceGroups,
		&oldestKnownAgeMa,
		&item.ValenceElements,
		&imaMineralSymbol,
		&item.ParageneticModes,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrMineralNotFound
		}
		return nil, err
	}
	item.ValenceChemistry = nullableStringPtr(valenceChemistry)
	item.ImaChemistry = nullableStringPtr(imaChemistry)
	item.ImaNumber = nullableStringPtr(imaNumber)
	item.ImaMineralSymbol = nullableStringPtr(imaMineralSymbol)
	item.ImaStatus = nullableStringPtr(imaStatus)
	item.StructuralGroupname = nullableStringPtr(structuralGroupname)
	item.StatusNotes = nullableStringPtr(statusNotes)
	item.CountryOfTypeLocality = nullableStringPtr(countryOfTypeLocality)
	item.YearFirstPublished = nullableIntPtr(yearFirstPublished)
	age, err := nullableFloatPtr(oldestKnownAgeMa)
	if err != nil {
		return nil, err
	}
	item.OldestKnownAgeMa = age
	item.ChemistryElements = nonEmptyStrings(item.ChemistryElements)
	item.RRUFFIDs = nonEmptyStrings(item.RRUFFIDs)
	item.CrystalSystems = nonEmptyStrings(item.CrystalSystems)
	item.SpaceGroups = nonEmptyStrings(item.SpaceGroups)
	item.ValenceElements = nonEmptyStrings(item.ValenceElements)
	item.ParageneticModes = nonEmptyStrings(item.ParageneticModes)
	return &item, nil
}
func buildMineralsFilterSQL(params models.MineralsListQuery) (string, []any) {
	clauses := make([]string, 0, 4)
	args := make([]any, 0, 8)
	if params.OnlyWithImages {
		clauses = append(clauses, `
			EXISTS (
				SELECT 1
				FROM storage_object_usages usage
				WHERE
					usage.entity_type = 'chemistry_mineral'
					AND usage.entity_id = chemistry_minerals.database_id::text
					AND usage.usage_type = 'image'
					AND usage.field_name = 'gallery'
			)
		`)
	}
	if pattern := buildMineralsSearchPattern(params.Search); pattern != "" {
		args = append(args, pattern)
		placeholder := fmt.Sprintf("$%d", len(args))
		args = append(args, strings.TrimSpace(params.Search))
		exactPlaceholder := fmt.Sprintf("$%d", len(args))
		clauses = append(clauses, fmt.Sprintf(
			`(
				mineral_name ILIKE %s OR
				mineral_name_plain ILIKE %s OR
				mineral_name_search ILIKE %s OR
				COALESCE(ima_number, '') ILIKE %s OR
				COALESCE(ima_mineral_symbol, '') ILIKE %s OR
				COALESCE(ima_chemistry, '') ILIKE %s OR
				database_id::text = %s
			)`,
			placeholder,
			placeholder,
			placeholder,
			placeholder,
			placeholder,
			placeholder,
			exactPlaceholder,
		))
	}
	if len(params.ChemistryAll) > 0 {
		args = append(args, params.ChemistryAll)
		clauses = append(clauses, fmt.Sprintf(
			"COALESCE(chemistry_elements, ARRAY[]::text[]) @> $%d::text[]",
			len(args),
		))
	}
	if len(params.ChemistryAny) > 0 {
		args = append(args, params.ChemistryAny)
		clauses = append(clauses, fmt.Sprintf(
			"COALESCE(chemistry_elements, ARRAY[]::text[]) && $%d::text[]",
			len(args),
		))
	}
	if len(params.ChemistryNone) > 0 {
		args = append(args, params.ChemistryNone)
		clauses = append(clauses, fmt.Sprintf(
			"NOT (COALESCE(chemistry_elements, ARRAY[]::text[]) && $%d::text[])",
			len(args),
		))
	}
	if len(params.CrystalSystems) > 0 {
		args = append(args, params.CrystalSystems)
		if params.CrystalSystemMode == models.MineralCrystalSystemModeAll {
			clauses = append(clauses, fmt.Sprintf(
				"COALESCE(crystal_systems, ARRAY[]::text[]) @> $%d::text[]",
				len(args),
			))
		} else {
			clauses = append(clauses, fmt.Sprintf(
				"COALESCE(crystal_systems, ARRAY[]::text[]) && $%d::text[]",
				len(args),
			))
		}
	}
	if len(clauses) == 0 {
		return "", args
	}
	return "WHERE " + strings.Join(clauses, " AND "), args
}
func PdbGetMineralByDatabaseID(ctx context.Context, databaseID int64) (*models.Mineral, error) {
	row := db.PDB.QueryRow(ctx, `
		SELECT
			id,
			database_id,
			mineral_name,
			mineral_name_plain,
			mineral_name_search,
			valence_chemistry,
			ima_chemistry,
			chemistry_elements,
			ima_number,
			rruff_ids,
			country_of_type_locality,
			year_first_published,
			ima_status,
			structural_groupname,
			status_notes,
			crystal_systems,
			space_groups,
			oldest_known_age_ma,
			valence_elements,
			ima_mineral_symbol,
			paragenetic_modes,
			created_at,
			updated_at
		FROM chemistry_minerals
		WHERE database_id = $1
		LIMIT 1
	`, databaseID)
	return scanMineral(row)
}
func PdbGetMineralByDatabaseIDString(ctx context.Context, raw string) (*models.Mineral, error) {
	databaseID, err := strconv.ParseInt(strings.TrimSpace(raw), 10, 64)
	if err != nil || databaseID <= 0 {
		return nil, ErrMineralNotFound
	}
	return PdbGetMineralByDatabaseID(ctx, databaseID)
}
func PdbListMineralImages(ctx context.Context, databaseID int64) ([]models.MineralImage, error) {
	rows, err := db.PDB.Query(ctx, `
		SELECT metadata
		FROM storage_object_usages
		WHERE
			entity_type = 'chemistry_mineral'
			AND entity_id = $1
			AND usage_type = 'image'
			AND field_name = 'gallery'
		ORDER BY sort_order ASC, created_at ASC
	`, strconv.FormatInt(databaseID, 10))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]models.MineralImage, 0, 8)
	for rows.Next() {
		var rawMetadata []byte
		if scanErr := rows.Scan(&rawMetadata); scanErr != nil {
			return nil, scanErr
		}
		var payload struct {
			File    string `json:"file"`
			RRUFFID string `json:"rruff_id"`
			Order   int    `json:"order"`
		}
		if err := json.Unmarshal(rawMetadata, &payload); err != nil {
			return nil, err
		}
		file := strings.TrimSpace(payload.File)
		if file == "" {
			continue
		}
		items = append(items, models.MineralImage{
			File:    file,
			RRUFFID: strings.TrimSpace(payload.RRUFFID),
			Order:   payload.Order,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
func PdbListMinerals(ctx context.Context, params models.MineralsListQuery) (*models.MineralsListResult, error) {
	whereSQL, filterArgs := buildMineralsFilterSQL(params)
	result := &models.MineralsListResult{
		Items: []models.MineralListItem{},
		Meta: models.MineralsListMeta{
			Limit:  params.Limit,
			Offset: params.Offset,
			Total:  0,
		},
		Facets: models.MineralsListFacets{
			ChemistryElementsAvailable: map[string]int64{},
		},
	}
	for _, symbol := range allChemicalElementSymbols {
		result.Facets.ChemistryElementsAvailable[symbol] = 0
	}
	if err := db.PDB.QueryRow(
		ctx,
		fmt.Sprintf(`SELECT COUNT(*)::BIGINT FROM chemistry_minerals %s`, whereSQL),
		filterArgs...,
	).Scan(&result.Meta.Total); err != nil {
		return nil, err
	}
	sortClause := "mineral_name ASC, database_id ASC"
	if params.Sort == models.MineralSortNameDesc {
		sortClause = "mineral_name DESC, database_id DESC"
	}
	listArgs := append([]any{}, filterArgs...)
	listArgs = append(listArgs, params.Limit, params.Offset)
	rows, err := db.PDB.Query(ctx, fmt.Sprintf(`
		SELECT
			database_id,
			mineral_name,
			crystal_systems
		FROM chemistry_minerals
		%s
		ORDER BY %s
		LIMIT $%d OFFSET $%d
	`, whereSQL, sortClause, len(listArgs)-1, len(listArgs)), listArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var item models.MineralListItem
		if err := rows.Scan(
			&item.DatabaseID,
			&item.MineralName,
			&item.CrystalSystems,
		); err != nil {
			return nil, err
		}
		item.CrystalSystems = nonEmptyStrings(item.CrystalSystems)
		result.Items = append(result.Items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	facetRows, err := db.PDB.Query(ctx, fmt.Sprintf(`
		SELECT element, COUNT(*)::BIGINT
		FROM (
			SELECT DISTINCT chemistry_minerals.id, element
			FROM chemistry_minerals
			CROSS JOIN LATERAL UNNEST(COALESCE(chemistry_elements, ARRAY[]::text[])) AS element
			%s
		) AS filtered_elements
		GROUP BY element
	`, whereSQL), filterArgs...)
	if err != nil {
		return nil, err
	}
	defer facetRows.Close()
	for facetRows.Next() {
		var symbol string
		var count int64
		if err := facetRows.Scan(&symbol, &count); err != nil {
			return nil, err
		}
		result.Facets.ChemistryElementsAvailable[symbol] = count
	}
	if err := facetRows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
