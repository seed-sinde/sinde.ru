package services

import (
	"strings"
	"testing"

	"sinde.ru/internal/models"
)

func TestBuildMineralsFilterSQLImageFilter(t *testing.T) {
	tests := []struct {
		name          string
		filter        models.MineralImageFilter
		wantFragment  string
		wantFragments []string
	}{
		{
			name:   "any omits image clause",
			filter: models.MineralImageFilterAny,
		},
		{
			name:         "with adds exists clause",
			filter:       models.MineralImageFilterWith,
			wantFragment: "EXISTS",
			wantFragments: []string{
				"usage.entity_type = 'chemistry_mineral'",
				"usage.usage_type = 'image'",
				"usage.field_name = 'gallery'",
			},
		},
		{
			name:         "without adds not exists clause",
			filter:       models.MineralImageFilterWithout,
			wantFragment: "NOT (",
			wantFragments: []string{
				"EXISTS",
				"usage.entity_type = 'chemistry_mineral'",
				"usage.usage_type = 'image'",
				"usage.field_name = 'gallery'",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sql, args := buildMineralsFilterSQL(models.MineralsListQuery{ImageFilter: test.filter})
			if len(args) != 0 {
				t.Fatalf("args length = %d, want 0", len(args))
			}
			if test.wantFragment == "" {
				if strings.TrimSpace(sql) != "" {
					t.Fatalf("sql = %q, want empty", sql)
				}
				return
			}
			if !strings.Contains(sql, test.wantFragment) {
				t.Fatalf("sql = %q, want fragment %q", sql, test.wantFragment)
			}
			for _, fragment := range test.wantFragments {
				if !strings.Contains(sql, fragment) {
					t.Fatalf("sql = %q, want fragment %q", sql, fragment)
				}
			}
		})
	}
}
