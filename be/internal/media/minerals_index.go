package media
import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"sync"
)
type MineralImageMeta struct {
	File    string `json:"file"`
	RRUFFID string `json:"rruff_id"`
	Order   int    `json:"order"`
}
var (
	mineralImagesMu    sync.RWMutex
	mineralImagesIndex = map[int64][]MineralImageMeta{}
)
func mineralsMetaPath() string {
	return filepath.Join(".", "storage", "media", "minerals", "meta.json")
}
func LoadMineralImagesIndex() error {
	raw, err := ReadObjectBytes(context.Background(), "minerals/meta.json")
	if err != nil {
		path := mineralsMetaPath()
		raw, err = os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				mineralImagesMu.Lock()
				mineralImagesIndex = map[int64][]MineralImageMeta{}
				mineralImagesMu.Unlock()
				return nil
			}
			return err
		}
	}
	var decoded map[string][]MineralImageMeta
	if err := json.Unmarshal(raw, &decoded); err != nil {
		return err
	}
	next := make(map[int64][]MineralImageMeta, len(decoded))
	for rawID, items := range decoded {
		id, err := strconv.ParseInt(rawID, 10, 64)
		if err != nil || id <= 0 {
			continue
		}
		clean := make([]MineralImageMeta, 0, len(items))
		for _, item := range items {
			if item.File == "" {
				continue
			}
			clean = append(clean, item)
		}
		sort.SliceStable(clean, func(i, j int) bool {
			if clean[i].Order != clean[j].Order {
				return clean[i].Order < clean[j].Order
			}
			if clean[i].RRUFFID != clean[j].RRUFFID {
				return clean[i].RRUFFID < clean[j].RRUFFID
			}
			return clean[i].File < clean[j].File
		})
		if len(clean) > 0 {
			next[id] = clean
		}
	}
	mineralImagesMu.Lock()
	mineralImagesIndex = next
	mineralImagesMu.Unlock()
	return nil
}
func GetMineralImages(databaseID int64) []MineralImageMeta {
	if databaseID <= 0 {
		return nil
	}
	mineralImagesMu.RLock()
	items := mineralImagesIndex[databaseID]
	mineralImagesMu.RUnlock()
	if len(items) == 0 {
		return nil
	}
	out := make([]MineralImageMeta, len(items))
	copy(out, items)
	return out
}
func GetMineralImageDatabaseIDs() []int64 {
	mineralImagesMu.RLock()
	ids := make([]int64, 0, len(mineralImagesIndex))
	for databaseID := range mineralImagesIndex {
		ids = append(ids, databaseID)
	}
	mineralImagesMu.RUnlock()
	sort.Slice(ids, func(i, j int) bool {
		return ids[i] < ids[j]
	})
	return ids
}
