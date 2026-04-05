package handlers
import (
	"errors"
	"strconv"
	"strings"
	"github.com/gofiber/fiber/v3"
	"sinde.ru/db/services"
	"sinde.ru/internal/http/responses"
	"sinde.ru/internal/media"
	"sinde.ru/internal/models"
)
const (
	defaultMineralsLimit = 30
	maxMineralsLimit     = 100
)
func parseMineralQueryList(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	seen := make(map[string]struct{}, 8)
	out := make([]string, 0, 8)
	for _, chunk := range strings.Split(raw, ",") {
		value := strings.TrimSpace(chunk)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}
func normalizeMineralChemistrySymbol(raw string) string {
	value := strings.TrimSpace(raw)
	if value == "" {
		return ""
	}
	if len(value) == 1 {
		return strings.ToUpper(value)
	}
	return strings.ToUpper(value[:1]) + strings.ToLower(value[1:])
}
func normalizeMineralChemistryQueryList(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	seen := make(map[string]struct{}, len(values))
	out := make([]string, 0, len(values))
	for _, value := range values {
		normalized := normalizeMineralChemistrySymbol(value)
		if normalized == "" {
			continue
		}
		if _, ok := seen[normalized]; ok {
			continue
		}
		seen[normalized] = struct{}{}
		out = append(out, normalized)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}
var allowedMineralCrystalSystems = map[string]struct{}{
	string(models.MineralCrystalSystemCubic):        {},
	string(models.MineralCrystalSystemHexagonal):    {},
	string(models.MineralCrystalSystemMonoclinic):   {},
	string(models.MineralCrystalSystemOrthorhombic): {},
	string(models.MineralCrystalSystemTetragonal):   {},
	string(models.MineralCrystalSystemTriclinic):    {},
	string(models.MineralCrystalSystemUnknown):      {},
}
func parseMineralsListQuery(c fiber.Ctx) (models.MineralsListQuery, error) {
	query := models.MineralsListQuery{
		Search:            strings.TrimSpace(c.Query("q")),
		Limit:             defaultMineralsLimit,
		Offset:            0,
		Sort:              models.MineralSortNameAsc,
		CrystalSystemMode: models.MineralCrystalSystemModeAny,
	}
	if rawLimit := strings.TrimSpace(c.Query("limit")); rawLimit != "" {
		limit, err := strconv.Atoi(rawLimit)
		if err != nil || limit <= 0 {
			return query, fiber.NewError(fiber.StatusBadRequest, "Некорректный запрос")
		}
		if limit > maxMineralsLimit {
			limit = maxMineralsLimit
		}
		query.Limit = limit
	}
	if rawOffset := strings.TrimSpace(c.Query("offset")); rawOffset != "" {
		offset, err := strconv.Atoi(rawOffset)
		if err != nil || offset < 0 {
			return query, fiber.NewError(fiber.StatusBadRequest, "Некорректный запрос")
		}
		query.Offset = offset
	}
	if rawSort := strings.TrimSpace(c.Query("sort")); rawSort != "" {
		sort := models.MineralSort(rawSort)
		if sort != models.MineralSortNameAsc && sort != models.MineralSortNameDesc {
			return query, fiber.NewError(fiber.StatusBadRequest, "Некорректный запрос")
		}
		query.Sort = sort
	}
	if rawOnlyWithImages := strings.TrimSpace(c.Query("onlyWithImages")); rawOnlyWithImages != "" {
		switch strings.ToLower(rawOnlyWithImages) {
		case "1", "true", "yes", "on":
			query.OnlyWithImages = true
		case "0", "false", "no", "off":
			query.OnlyWithImages = false
		default:
			return query, fiber.NewError(fiber.StatusBadRequest, "Некорректный запрос")
		}
	}
	query.ChemistryAll = normalizeMineralChemistryQueryList(parseMineralQueryList(c.Query("chemistryAll")))
	query.ChemistryAny = normalizeMineralChemistryQueryList(parseMineralQueryList(c.Query("chemistryAny")))
	query.ChemistryNone = normalizeMineralChemistryQueryList(parseMineralQueryList(c.Query("chemistryNone")))
	query.CrystalSystems = parseMineralQueryList(c.Query("crystalSystem"))
	if rawCrystalSystemMode := strings.TrimSpace(c.Query("crystalSystemMode")); rawCrystalSystemMode != "" {
		mode := models.MineralCrystalSystemMode(rawCrystalSystemMode)
		if mode != models.MineralCrystalSystemModeAny && mode != models.MineralCrystalSystemModeAll {
			return query, fiber.NewError(fiber.StatusBadRequest, "Некорректный запрос")
		}
		query.CrystalSystemMode = mode
	}
	for _, value := range query.CrystalSystems {
		if _, ok := allowedMineralCrystalSystems[value]; !ok {
			return query, fiber.NewError(fiber.StatusBadRequest, "Некорректный запрос")
		}
	}
	owners := make(map[string]string, len(query.ChemistryAll)+len(query.ChemistryAny)+len(query.ChemistryNone))
	for _, value := range query.ChemistryAll {
		owners[value] = "all"
	}
	for _, value := range query.ChemistryAny {
		if owners[value] != "" {
			return query, fiber.NewError(fiber.StatusBadRequest, "Некорректный запрос")
		}
		owners[value] = "any"
	}
	for _, value := range query.ChemistryNone {
		if owners[value] != "" {
			return query, fiber.NewError(fiber.StatusBadRequest, "Некорректный запрос")
		}
		owners[value] = "none"
	}
	return query, nil
}
func MineralGetByDatabaseIDHandler() fiber.Handler {
	return func(c fiber.Ctx) error {
		rawID := strings.TrimSpace(c.Params("database_id"))
		if rawID == "" {
			rawID = strings.TrimSpace(c.Params("id"))
		}
		if rawID == "" {
			return responses.Error(c, fiber.StatusBadRequest, "Некорректный запрос")
		}
		databaseID, err := strconv.ParseInt(rawID, 10, 64)
		if err != nil || databaseID <= 0 {
			return responses.Error(c, fiber.StatusBadRequest, "Некорректный запрос")
		}
		item, err := services.PdbGetMineralByDatabaseID(c.Context(), databaseID)
		if err != nil {
			if errors.Is(err, services.ErrMineralNotFound) {
				return responses.Error(c, fiber.StatusNotFound, "Минерал не найден")
			}
			return responses.Error(c, fiber.StatusInternalServerError, "Не удалось загрузить минерал", err.Error())
		}
		mediaItems := media.GetMineralImages(item.DatabaseID)
		if len(mediaItems) > 0 {
			item.Images = make([]models.MineralImage, 0, len(mediaItems))
			for _, image := range mediaItems {
				item.Images = append(item.Images, models.MineralImage{
					File:    image.File,
					RRUFFID: image.RRUFFID,
					Order:   image.Order,
				})
			}
		}
		return responses.Success(c, fiber.StatusOK, item)
	}
}
func MineralsListHandler() fiber.Handler {
	return func(c fiber.Ctx) error {
		query, err := parseMineralsListQuery(c)
		if err != nil {
			status := fiber.StatusBadRequest
			message := "Некорректный запрос"
			var fiberErr *fiber.Error
			if errors.As(err, &fiberErr) {
				status = fiberErr.Code
				if text := strings.TrimSpace(fiberErr.Message); text != "" {
					message = text
				}
			}
			return responses.Error(c, status, message)
		}
		result, err := services.PdbListMinerals(c.Context(), query)
		if err != nil {
			return responses.Error(c, fiber.StatusInternalServerError, "Не удалось загрузить минералы", err.Error())
		}
		return responses.Success(c, fiber.StatusOK, result)
	}
}
