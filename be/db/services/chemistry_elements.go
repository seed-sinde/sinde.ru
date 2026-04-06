package services

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"

	"sinde.ru/db"
	"sinde.ru/internal/models"
)

type chemistryElementMediaUsageSet struct {
	bohrModel  *string
	bohrModel3 bool
	spectral   *string
	model      *string
	samples    []models.ChemistryElementSample
}

func nullableFloat64Ptr(value pgtype.Float8) *float64 {
	if !value.Valid {
		return nil
	}
	number := value.Float64
	return &number
}

func nullableChemistryFallback(raw []byte) (*models.ChemistryElementSampleFallback, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return nil, nil
	}
	var item models.ChemistryElementSampleFallback
	if err := json.Unmarshal(raw, &item); err != nil {
		return nil, err
	}
	item.Type = strings.TrimSpace(item.Type)
	item.Label = strings.TrimSpace(item.Label)
	item.Description = strings.TrimSpace(item.Description)
	if item.Type == "" || item.Label == "" || item.Description == "" {
		return nil, nil
	}
	return &item, nil
}

func listChemistryElementMediaUsages(ctx context.Context) (map[int]*chemistryElementMediaUsageSet, error) {
	rows, err := db.PDB.Query(ctx, `
		SELECT
			usage.entity_id,
			usage.field_name,
			object.storage_key,
			usage.metadata
		FROM storage_object_usages usage
		JOIN storage_objects object ON object.object_id = usage.object_id
		WHERE usage.entity_type = 'chemistry_element'
		ORDER BY usage.entity_id::integer ASC, usage.sort_order ASC, usage.created_at ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make(map[int]*chemistryElementMediaUsageSet, 128)
	for rows.Next() {
		var entityID string
		var fieldName string
		var storageKey string
		var rawMetadata []byte
		if err := rows.Scan(&entityID, &fieldName, &storageKey, &rawMetadata); err != nil {
			return nil, err
		}
		number, err := strconv.Atoi(strings.TrimSpace(entityID))
		if err != nil || number <= 0 {
			continue
		}
		item := items[number]
		if item == nil {
			item = &chemistryElementMediaUsageSet{}
			items[number] = item
		}
		key := strings.TrimSpace(storageKey)
		switch strings.TrimSpace(fieldName) {
		case "bohr_model":
			if key != "" {
				item.bohrModel = &key
			}
		case "spectral":
			if key != "" {
				item.spectral = &key
			}
		case "bohr_model_3d":
			if key != "" {
				item.model = &key
			}
		case "sample":
			if key == "" {
				continue
			}
			payload := struct {
				Title       *string `json:"title"`
				Author      *string `json:"author"`
				Attribution *string `json:"attribution"`
				SourceURL   *string `json:"source_url"`
				LicenseURL  *string `json:"license_url"`
				License     *string `json:"license"`
			}{}
			if len(rawMetadata) > 0 {
				if err := json.Unmarshal(rawMetadata, &payload); err != nil {
					return nil, err
				}
			}
			item.samples = append(item.samples, models.ChemistryElementSample{
				File:        pathBase(key),
				Title:       trimStringPtr(payload.Title),
				Author:      trimStringPtr(payload.Author),
				Attribution: trimStringPtr(payload.Attribution),
				SourceURL:   trimStringPtr(payload.SourceURL),
				LicenseURL:  trimStringPtr(payload.LicenseURL),
				License:     trimStringPtr(payload.License),
			})
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func pathBase(value string) string {
	normalized := strings.TrimSpace(value)
	if normalized == "" {
		return ""
	}
	parts := strings.Split(normalized, "/")
	return parts[len(parts)-1]
}

func trimStringPtr(value *string) *string {
	if value == nil {
		return nil
	}
	text := strings.TrimSpace(*value)
	if text == "" {
		return nil
	}
	return &text
}

func PdbListChemistryElements(ctx context.Context) ([]models.ChemistryElement, error) {
	usageMap, err := listChemistryElementMediaUsages(ctx)
	if err != nil {
		return nil, err
	}
	rows, err := db.PDB.Query(ctx, `
		SELECT
			number,
			symbol,
			name,
			russian_name,
			appearance,
			atomic_mass::double precision,
			boil::double precision,
			category,
			density::double precision,
			discovered_by,
			melt::double precision,
			molar_heat::double precision,
			named_by,
			period,
			group_number,
			phase,
			source,
			summary,
			xpos,
			ypos,
			wxpos,
			wypos,
			COALESCE(shells, ARRAY[]::integer[]),
			electron_configuration,
			electron_configuration_semantic,
			electron_affinity::double precision,
			electronegativity_pauling::double precision,
			COALESCE(ionization_energies, ARRAY[]::double precision[]),
			cpk_hex,
			block,
			sample_fallback
		FROM chemistry_elements
		ORDER BY number ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]models.ChemistryElement, 0, 118)
	for rows.Next() {
		item := models.ChemistryElement{}
		var appearance pgtype.Text
		var atomicMass pgtype.Float8
		var boil pgtype.Float8
		var density pgtype.Float8
		var discoveredBy pgtype.Text
		var melt pgtype.Float8
		var molarHeat pgtype.Float8
		var namedBy pgtype.Text
		var period pgtype.Int4
		var group pgtype.Int4
		var phase pgtype.Text
		var source pgtype.Text
		var wxpos pgtype.Int4
		var wypos pgtype.Int4
		var shells []int32
		var electronConfiguration pgtype.Text
		var electronConfigurationSemantic pgtype.Text
		var electronAffinity pgtype.Float8
		var electronegativityPauling pgtype.Float8
		var ionizationEnergies []float64
		var cpkHex pgtype.Text
		var block pgtype.Text
		var sampleFallbackRaw []byte
		if err := rows.Scan(
			&item.Number,
			&item.Symbol,
			&item.Name,
			&item.RussianName,
			&appearance,
			&atomicMass,
			&boil,
			&item.Category,
			&density,
			&discoveredBy,
			&melt,
			&molarHeat,
			&namedBy,
			&period,
			&group,
			&phase,
			&source,
			&item.Summary,
			&item.XPos,
			&item.YPos,
			&wxpos,
			&wypos,
			&shells,
			&electronConfiguration,
			&electronConfigurationSemantic,
			&electronAffinity,
			&electronegativityPauling,
			&ionizationEnergies,
			&cpkHex,
			&block,
			&sampleFallbackRaw,
		); err != nil {
			return nil, err
		}
		item.Appearance = nullableStringPtr(appearance)
		item.AtomicMass = nullableFloat64Ptr(atomicMass)
		item.Boil = nullableFloat64Ptr(boil)
		item.Density = nullableFloat64Ptr(density)
		item.DiscoveredBy = nullableStringPtr(discoveredBy)
		item.Melt = nullableFloat64Ptr(melt)
		item.MolarHeat = nullableFloat64Ptr(molarHeat)
		item.NamedBy = nullableStringPtr(namedBy)
		item.Period = nullableIntPtr(period)
		item.Group = nullableIntPtr(group)
		item.Phase = nullableStringPtr(phase)
		item.Source = nullableStringPtr(source)
		item.WXPos = nullableIntPtr(wxpos)
		item.WYPos = nullableIntPtr(wypos)
		item.Shells = make([]int, 0, len(shells))
		for _, shell := range shells {
			item.Shells = append(item.Shells, int(shell))
		}
		item.ElectronConfiguration = nullableStringPtr(electronConfiguration)
		item.ElectronConfigurationSemantic = nullableStringPtr(electronConfigurationSemantic)
		item.ElectronAffinity = nullableFloat64Ptr(electronAffinity)
		item.ElectronegativityPauling = nullableFloat64Ptr(electronegativityPauling)
		item.IonizationEnergies = ionizationEnergies
		item.CPKHex = nullableStringPtr(cpkHex)
		item.Block = nullableStringPtr(block)
		item.SampleFallback, err = nullableChemistryFallback(sampleFallbackRaw)
		if err != nil {
			return nil, err
		}
		if usage := usageMap[item.Number]; usage != nil {
			item.BohrModelImage = usage.bohrModel
			item.BohrModel3D = usage.model
			item.SpectralImg = usage.spectral
			item.Samples = usage.samples
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
