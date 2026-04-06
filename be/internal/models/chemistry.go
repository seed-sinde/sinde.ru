package models

type ChemistryElementSample struct {
	File        string  `json:"file"`
	Title       *string `json:"title,omitempty"`
	Author      *string `json:"author,omitempty"`
	Attribution *string `json:"attribution,omitempty"`
	SourceURL   *string `json:"source_url,omitempty"`
	LicenseURL  *string `json:"license_url,omitempty"`
	License     *string `json:"license,omitempty"`
}

type ChemistryElementSampleFallback struct {
	Type        string `json:"type"`
	Label       string `json:"label"`
	Description string `json:"description"`
}

type ChemistryElement struct {
	Number                        int                             `json:"number"`
	Symbol                        string                          `json:"symbol"`
	Name                          string                          `json:"name"`
	RussianName                   string                          `json:"russian_name"`
	Appearance                    *string                         `json:"appearance,omitempty"`
	AtomicMass                    *float64                        `json:"atomic_mass,omitempty"`
	Boil                          *float64                        `json:"boil,omitempty"`
	Category                      string                          `json:"category"`
	Density                       *float64                        `json:"density,omitempty"`
	DiscoveredBy                  *string                         `json:"discovered_by,omitempty"`
	Melt                          *float64                        `json:"melt,omitempty"`
	MolarHeat                     *float64                        `json:"molar_heat,omitempty"`
	NamedBy                       *string                         `json:"named_by,omitempty"`
	Period                        *int                            `json:"period,omitempty"`
	Group                         *int                            `json:"group,omitempty"`
	Phase                         *string                         `json:"phase,omitempty"`
	Source                        *string                         `json:"source,omitempty"`
	BohrModelImage                *string                         `json:"bohr_model_image,omitempty"`
	BohrModel3D                   *string                         `json:"bohr_model_3d,omitempty"`
	SpectralImg                   *string                         `json:"spectral_img,omitempty"`
	Summary                       string                          `json:"summary"`
	XPos                          int                             `json:"xpos"`
	YPos                          int                             `json:"ypos"`
	WXPos                         *int                            `json:"wxpos,omitempty"`
	WYPos                         *int                            `json:"wypos,omitempty"`
	Shells                        []int                           `json:"shells,omitempty"`
	ElectronConfiguration         *string                         `json:"electron_configuration,omitempty"`
	ElectronConfigurationSemantic *string                         `json:"electron_configuration_semantic,omitempty"`
	ElectronAffinity              *float64                        `json:"electron_affinity,omitempty"`
	ElectronegativityPauling      *float64                        `json:"electronegativity_pauling,omitempty"`
	IonizationEnergies            []float64                       `json:"ionization_energies,omitempty"`
	CPKHex                        *string                         `json:"cpk_hex,omitempty"`
	Block                         *string                         `json:"block,omitempty"`
	Samples                       []ChemistryElementSample        `json:"samples,omitempty"`
	SampleFallback                *ChemistryElementSampleFallback `json:"sample_fallback,omitempty"`
}
