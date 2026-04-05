package models
import "time"
type MineralSort string
const (
	MineralSortNameAsc  MineralSort = "name_asc"
	MineralSortNameDesc MineralSort = "name_desc"
)
type MineralCrystalSystem string
type MineralCrystalSystemMode string
const (
	MineralCrystalSystemCubic        MineralCrystalSystem = "cubic"
	MineralCrystalSystemHexagonal    MineralCrystalSystem = "hexagonal"
	MineralCrystalSystemMonoclinic   MineralCrystalSystem = "monoclinic"
	MineralCrystalSystemOrthorhombic MineralCrystalSystem = "orthorhombic"
	MineralCrystalSystemTetragonal   MineralCrystalSystem = "tetragonal"
	MineralCrystalSystemTriclinic    MineralCrystalSystem = "triclinic"
	MineralCrystalSystemUnknown      MineralCrystalSystem = "unknown"
	MineralCrystalSystemModeAny MineralCrystalSystemMode = "any"
	MineralCrystalSystemModeAll MineralCrystalSystemMode = "all"
)
type MineralImage struct {
	File    string `json:"file"`
	RRUFFID string `json:"rruff_id"`
	Order   int    `json:"order"`
}
type Mineral struct {
	ID                    int64     `json:"id"`
	DatabaseID            int64     `json:"database_id"`
	MineralName           string    `json:"mineral_name"`
	MineralNamePlain      string    `json:"mineral_name_plain"`
	ValenceChemistry      *string   `json:"valence_chemistry,omitempty"`
	ImaChemistry          *string   `json:"ima_chemistry,omitempty"`
	ChemistryElements     []string  `json:"chemistry_elements,omitempty"`
	ImaNumber             *string   `json:"ima_number,omitempty"`
	RRUFFIDs              []string  `json:"rruff_ids,omitempty"`
	CountryOfTypeLocality *string   `json:"country_of_type_locality,omitempty"`
	YearFirstPublished    *int      `json:"year_first_published,omitempty"`
	ImaStatus             *string   `json:"ima_status,omitempty"`
	StructuralGroupname   *string   `json:"structural_groupname,omitempty"`
	StatusNotes           *string   `json:"status_notes,omitempty"`
	CrystalSystems        []string  `json:"crystal_systems,omitempty"`
	SpaceGroups           []string  `json:"space_groups,omitempty"`
	OldestKnownAgeMa      *float64  `json:"oldest_known_age_ma,omitempty"`
	ValenceElements       []string  `json:"valence_elements,omitempty"`
	ImaMineralSymbol      *string   `json:"ima_mineral_symbol,omitempty"`
	ParageneticModes      []string  `json:"paragenetic_modes,omitempty"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
	MineralNameSearch string         `json:"-"`
	Images            []MineralImage `json:"images,omitempty"`
}
type MineralListItem struct {
	DatabaseID     int64    `json:"database_id"`
	MineralName    string   `json:"mineral_name"`
	CrystalSystems []string `json:"crystal_systems,omitempty"`
}
type MineralsListMeta struct {
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
	Total  int64 `json:"total"`
}
type MineralsListFacets struct {
	ChemistryElementsAvailable map[string]int64 `json:"chemistryElementsAvailable"`
}
type MineralsListResult struct {
	Items  []MineralListItem  `json:"items"`
	Meta   MineralsListMeta   `json:"meta"`
	Facets MineralsListFacets `json:"facets"`
}
type MineralsListQuery struct {
	Search            string
	Limit             int
	Offset            int
	Sort              MineralSort
	OnlyWithImages    bool
	CrystalSystems    []string
	CrystalSystemMode MineralCrystalSystemMode
	ChemistryAll      []string
	ChemistryAny      []string
	ChemistryNone     []string
}
type MineralImportResult struct {
	Processed int `json:"processed"`
	Inserted  int `json:"inserted"`
	Updated   int `json:"updated"`
}
