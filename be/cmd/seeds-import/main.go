package main
import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"sinde.ru/db"
	"sinde.ru/utils"
)
const DefaultSeedsRoot = "./db/seeds"
var identRE = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
type Manifest struct {
	Version int            `json:"version"`
	Files   []ManifestFile `json:"files"`
}
type ManifestFile struct {
	File      string   `json:"file"`
	Table     string   `json:"table"`
	Order     int      `json:"order"`
	Mode      string   `json:"mode,omitempty"`
	Columns   []string `json:"columns,omitempty"`
	Delimiter string   `json:"delimiter,omitempty"`
	Header    *bool    `json:"header,omitempty"`
	Null      string   `json:"null,omitempty"`
}
type ImportPlan struct {
	Name         string
	ManifestPath string
	SectionDir   string
	Entry        ManifestFile
}
type ImportOptions struct {
	SeedsRoot     string
	TruncateFirst bool
	ContinueOnErr bool
	DryRun        bool
	Logger        func(format string, args ...any)
}
type ImportResult struct {
	Name         string
	ManifestPath string
	CSVPath      string
	Table        string
	Rows         int64
}
type VerifyResult struct {
	Name         string
	CSVPath      string
	Table        string
	ExpectedRows int64
	ActualRows   int64
	AutoSeeded   bool
}
type Importer struct {
	pool *pgxpool.Pool
	opts ImportOptions
}
func main() {
	if err := Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
func Run(ctx context.Context) error {
	utils.LoadEnv()
	log.SetFlags(0)
	log.SetPrefix("")
	if err := db.Init(); err != nil {
		return fmt.Errorf("DB initialization error: %w", err)
	}
	defer db.PDB.Close()
	importer := New(db.PDB, ImportOptions{
		Logger: log.Printf,
	})
	command, names, autoSeedZero, err := parseCLI(os.Args[1:])
	if err != nil {
		return err
	}
	switch command {
	case "list":
		plans, err := DiscoverPlans(importer.opts.SeedsRoot)
		if err != nil {
			return err
		}
		for _, plan := range plans {
			log.Printf("%s", plan.Name)
		}
		return nil
	case "reset":
		results, err := importer.Reset(ctx, names)
		if err != nil {
			return fmt.Errorf("seeds reset failed: %w", err)
		}
		for _, item := range results {
			log.Printf("reset: %s -> %s", item.Name, item.Table)
		}
		return nil
	case "verify":
		results, err := importer.Verify(ctx, names, autoSeedZero)
		if err != nil {
			return fmt.Errorf("seeds verify failed: %w", err)
		}
		for _, item := range results {
			switch {
			case item.ExpectedRows == item.ActualRows:
				continue
			case item.AutoSeeded:
				log.Printf(
					"auto-seeded: %s -> %s (expected=%d, actual=%d)",
					item.Name,
					item.Table,
					item.ExpectedRows,
					item.ActualRows,
				)
			default:
				log.Printf(
					"mismatch: %s -> %s (expected=%d, actual=%d)",
					item.Name,
					item.Table,
					item.ExpectedRows,
					item.ActualRows,
				)
			}
		}
		return nil
	default:
		if _, err := importer.Run(ctx, names); err != nil {
			return fmt.Errorf("seeds import failed: %w", err)
		}
	}
	return nil
}
func parseCLI(args []string) (command string, names []string, autoSeedZero bool, err error) {
	command = "import"
	if len(args) == 0 {
		return command, nil, false, nil
	}
	switch strings.ToLower(strings.TrimSpace(args[0])) {
	case "import":
		return "import", normalizePlanNames(args[1:]), false, nil
	case "reset":
		return "reset", normalizePlanNames(args[1:]), false, nil
	case "verify":
		names = make([]string, 0, len(args)-1)
		for _, arg := range args[1:] {
			switch strings.TrimSpace(arg) {
			case "--auto-seed-zero":
				autoSeedZero = true
			case "":
			default:
				names = append(names, strings.TrimSpace(arg))
			}
		}
		return "verify", normalizePlanNames(names), autoSeedZero, nil
	case "list":
		return "list", nil, false, nil
	default:
		return "import", normalizePlanNames(args), false, nil
	}
}
func New(pool *pgxpool.Pool, opts ImportOptions) *Importer {
	if opts.SeedsRoot == "" {
		opts.SeedsRoot = resolveSeedsRoot()
	}
	if opts.Logger == nil {
		opts.Logger = func(string, ...any) {}
	}
	return &Importer{
		pool: pool,
		opts: opts,
	}
}
func resolveSeedsRoot() string {
	if value := strings.TrimSpace(os.Getenv("SEEDS_ROOT")); value != "" {
		return value
	}
	return DefaultSeedsRoot
}
func (i *Importer) Run(ctx context.Context, names []string) ([]ImportResult, error) {
	plans, err := DiscoverPlans(i.opts.SeedsRoot)
	if err != nil {
		return nil, err
	}
	plans, err = filterPlans(plans, names)
	if err != nil {
		return nil, err
	}
	if len(plans) == 0 {
		return nil, fmt.Errorf("no import plans found in %s", i.opts.SeedsRoot)
	}
	results := make([]ImportResult, 0, len(plans))
	for _, plan := range plans {
		result, err := i.importPlan(ctx, plan)
		if err != nil {
			if i.opts.ContinueOnErr {
				i.opts.Logger("seed import failed: %v", err)
				continue
			}
			return results, err
		}
		results = append(results, result)
	}
	return results, nil
}
func (i *Importer) Reset(ctx context.Context, names []string) ([]ImportResult, error) {
	plans, err := DiscoverPlans(i.opts.SeedsRoot)
	if err != nil {
		return nil, err
	}
	plans, err = filterPlans(plans, names)
	if err != nil {
		return nil, err
	}
	if len(plans) == 0 {
		return nil, fmt.Errorf("no import plans found in %s", i.opts.SeedsRoot)
	}
	results := make([]ImportResult, 0, len(plans))
	for _, plan := range plans {
		if err := i.resetPlan(ctx, plan); err != nil {
			return results, err
		}
		results = append(results, ImportResult{
			Name:         plan.Name,
			ManifestPath: plan.ManifestPath,
			CSVPath:      filepath.Join(plan.SectionDir, plan.Entry.File),
			Table:        plan.Entry.Table,
		})
	}
	return results, nil
}
func (i *Importer) Verify(ctx context.Context, names []string, autoSeedZero bool) ([]VerifyResult, error) {
	plans, err := DiscoverPlans(i.opts.SeedsRoot)
	if err != nil {
		return nil, err
	}
	plans, err = filterPlans(plans, names)
	if err != nil {
		return nil, err
	}
	if len(plans) == 0 {
		return nil, fmt.Errorf("no import plans found in %s", i.opts.SeedsRoot)
	}
	results := make([]VerifyResult, 0, len(plans))
	for _, plan := range plans {
		csvPath := filepath.Join(plan.SectionDir, plan.Entry.File)
		expectedRows, err := countCSVRows(csvPath, plan.Entry)
		if err != nil {
			return results, err
		}
		actualRows, err := i.tableRowCount(ctx, plan.Entry.Table)
		if err != nil {
			return results, err
		}
		result := VerifyResult{
			Name:         plan.Name,
			CSVPath:      csvPath,
			Table:        plan.Entry.Table,
			ExpectedRows: expectedRows,
			ActualRows:   actualRows,
		}
		if autoSeedZero && expectedRows > 0 && actualRows == 0 {
			importResult, err := i.importPlan(ctx, plan)
			if err != nil {
				return results, err
			}
			result.ActualRows = importResult.Rows
			result.AutoSeeded = true
		}
		results = append(results, result)
	}
	return results, nil
}
func DiscoverPlans(seedsRoot string) ([]ImportPlan, error) {
	manifestPaths := make([]string, 0, 16)
	err := filepath.WalkDir(seedsRoot, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		if d.Name() == "_manifest.json" {
			manifestPaths = append(manifestPaths, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk seeds root: %w", err)
	}
	sort.Strings(manifestPaths)
	plans := make([]ImportPlan, 0, 32)
	planNames := make(map[string]string)
	for _, manifestPath := range manifestPaths {
		manifest, err := ReadManifest(manifestPath)
		if err != nil {
			return nil, err
		}
		if manifest.Version != 1 {
			return nil, fmt.Errorf("unsupported manifest version %d in %s", manifest.Version, manifestPath)
		}
		sectionDir := filepath.Dir(manifestPath)
		entries := append([]ManifestFile(nil), manifest.Files...)
		sort.SliceStable(entries, func(a, b int) bool {
			if entries[a].Order == entries[b].Order {
				return entries[a].File < entries[b].File
			}
			return entries[a].Order < entries[b].Order
		})
		for _, entry := range entries {
			if err := validateManifestEntry(sectionDir, entry); err != nil {
				return nil, fmt.Errorf("manifest %s: %w", manifestPath, err)
			}
			name := strings.TrimSuffix(filepath.Base(entry.File), filepath.Ext(entry.File))
			if existingPath, exists := planNames[name]; exists {
				return nil, fmt.Errorf("duplicate import name %q in %s and %s", name, existingPath, manifestPath)
			}
			planNames[name] = manifestPath
			plans = append(plans, ImportPlan{
				Name:         name,
				ManifestPath: manifestPath,
				SectionDir:   sectionDir,
				Entry:        entry,
			})
		}
	}
	return plans, nil
}
func ReadManifest(path string) (Manifest, error) {
	var manifest Manifest
	data, err := os.ReadFile(path)
	if err != nil {
		return manifest, fmt.Errorf("read manifest %s: %w", path, err)
	}
	if err := jsonUnmarshalStrict(data, &manifest); err != nil {
		return manifest, fmt.Errorf("parse manifest %s: %w", path, err)
	}
	if len(manifest.Files) == 0 {
		return manifest, fmt.Errorf("manifest %s has no files", path)
	}
	return manifest, nil
}
func (i *Importer) importPlan(ctx context.Context, plan ImportPlan) (ImportResult, error) {
	csvPath := filepath.Join(plan.SectionDir, plan.Entry.File)
	mode := normalizeMode(plan.Entry.Mode)
	hasHeader := resolveHeader(plan.Entry.Header)
	columns, err := resolveColumns(csvPath, plan.Entry)
	if err != nil {
		return ImportResult{}, err
	}
	if err := validateIdentifier(plan.Entry.Table); err != nil {
		return ImportResult{}, fmt.Errorf("csv %s: invalid table name: %w", csvPath, err)
	}
	if err := validateIdentifiers(columns); err != nil {
		return ImportResult{}, fmt.Errorf("csv %s: invalid columns: %w", csvPath, err)
	}
	if i.opts.DryRun {
		i.opts.Logger("dry-run %s -> %s (%s)", csvPath, plan.Entry.Table, mode)
		return ImportResult{
			Name:         plan.Name,
			ManifestPath: plan.ManifestPath,
			CSVPath:      csvPath,
			Table:        plan.Entry.Table,
			Rows:         0,
		}, nil
	}
	tx, err := i.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return ImportResult{}, fmt.Errorf("begin tx for %s: %w", csvPath, err)
	}
	defer tx.Rollback(ctx)
	quotedTable := quoteIdent(plan.Entry.Table)
	if i.opts.TruncateFirst || mode == "truncate" || mode == "replace" {
		if _, err := tx.Exec(ctx, "TRUNCATE TABLE "+quotedTable+" RESTART IDENTITY CASCADE"); err != nil {
			return ImportResult{}, fmt.Errorf("truncate %s: %w", plan.Entry.Table, err)
		}
	}
	copySQL := buildCopySQL(plan.Entry.Table, columns, hasHeader, plan.Entry.Delimiter, plan.Entry.Null)
	copyFile, err := os.Open(csvPath)
	if err != nil {
		return ImportResult{}, fmt.Errorf("open csv %s: %w", csvPath, err)
	}
	defer copyFile.Close()
	commandTag, err := tx.Conn().PgConn().CopyFrom(ctx, copyFile, copySQL)
	if err != nil {
		return ImportResult{}, fmt.Errorf("copy from %s into %s: %w", csvPath, plan.Entry.Table, err)
	}
	if err := tx.Commit(ctx); err != nil {
		return ImportResult{}, fmt.Errorf("commit %s: %w", csvPath, err)
	}
	rows := commandTag.RowsAffected()
	i.opts.Logger("imported %s -> %s (%d rows)", csvPath, plan.Entry.Table, rows)
	return ImportResult{
		Name:         plan.Name,
		ManifestPath: plan.ManifestPath,
		CSVPath:      csvPath,
		Table:        plan.Entry.Table,
		Rows:         rows,
	}, nil
}
func (i *Importer) resetPlan(ctx context.Context, plan ImportPlan) error {
	if i.opts.DryRun {
		i.opts.Logger("dry-run reset %s -> %s", plan.Name, plan.Entry.Table)
		return nil
	}
	quotedTable := quoteIdent(plan.Entry.Table)
	if _, err := i.pool.Exec(ctx, "TRUNCATE TABLE "+quotedTable+" RESTART IDENTITY CASCADE"); err != nil {
		return fmt.Errorf("truncate %s: %w", plan.Entry.Table, err)
	}
	return nil
}
func (i *Importer) tableRowCount(ctx context.Context, table string) (int64, error) {
	var rows int64
	if err := i.pool.QueryRow(ctx, "SELECT COUNT(*)::BIGINT FROM "+quoteIdent(table)).Scan(&rows); err != nil {
		return 0, fmt.Errorf("count %s: %w", table, err)
	}
	return rows, nil
}
func filterPlans(plans []ImportPlan, names []string) ([]ImportPlan, error) {
	if len(names) == 0 {
		return plans, nil
	}
	byName := make(map[string]ImportPlan, len(plans))
	for _, plan := range plans {
		byName[plan.Name] = plan
	}
	filtered := make([]ImportPlan, 0, len(names))
	for _, name := range names {
		plan, ok := byName[name]
		if !ok {
			return nil, fmt.Errorf("unknown import %q", name)
		}
		filtered = append(filtered, plan)
	}
	return filtered, nil
}
func normalizePlanNames(values []string) []string {
	out := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		name := strings.TrimSpace(value)
		if name == "" {
			continue
		}
		if _, exists := seen[name]; exists {
			continue
		}
		seen[name] = struct{}{}
		out = append(out, name)
	}
	return out
}
func countCSVRows(csvPath string, entry ManifestFile) (int64, error) {
	file, err := os.Open(csvPath)
	if err != nil {
		return 0, fmt.Errorf("open csv %s: %w", csvPath, err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.ReuseRecord = false
	reader.FieldsPerRecord = -1
	reader.TrimLeadingSpace = false
	if delim := strings.TrimSpace(entry.Delimiter); delim != "" {
		runes := []rune(delim)
		if len(runes) != 1 {
			return 0, fmt.Errorf("csv %s: delimiter must be exactly one rune", csvPath)
		}
		reader.Comma = runes[0]
	}
	var rows int64
	firstRow := true
	for {
		_, err := reader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return 0, fmt.Errorf("read csv %s: %w", csvPath, err)
		}
		if firstRow && resolveHeader(entry.Header) {
			firstRow = false
			continue
		}
		firstRow = false
		rows += 1
	}
	return rows, nil
}
func resolveColumns(csvPath string, entry ManifestFile) ([]string, error) {
	if len(entry.Columns) > 0 {
		return append([]string(nil), entry.Columns...), nil
	}
	if !resolveHeader(entry.Header) {
		return nil, fmt.Errorf("csv %s: columns are required when header=false", csvPath)
	}
	file, err := os.Open(csvPath)
	if err != nil {
		return nil, fmt.Errorf("open csv %s: %w", csvPath, err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.ReuseRecord = false
	reader.FieldsPerRecord = -1
	reader.TrimLeadingSpace = false
	if delim := strings.TrimSpace(entry.Delimiter); delim != "" {
		runes := []rune(delim)
		if len(runes) != 1 {
			return nil, fmt.Errorf("csv %s: delimiter must be exactly one rune", csvPath)
		}
		reader.Comma = runes[0]
	}
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("read first row %s: %w", csvPath, err)
	}
	return normalizeColumns(header), nil
}
func buildCopySQL(table string, columns []string, hasHeader bool, delimiter string, nullToken string) string {
	options := []string{"FORMAT csv"}
	if hasHeader {
		options = append(options, "HEADER true")
	}
	if delim := strings.TrimSpace(delimiter); delim != "" {
		options = append(options, "DELIMITER '"+escapeCopyLiteral(delim)+"'")
	}
	if nullValue := strings.TrimSpace(nullToken); nullValue != "" {
		options = append(options, "NULL '"+escapeCopyLiteral(nullValue)+"'")
	}
	return "COPY " + quoteIdent(table) +
		" (" + strings.Join(quoteIdents(columns), ", ") + ")" +
		" FROM STDIN WITH (" + strings.Join(options, ", ") + ")"
}
func normalizeMode(value string) string {
	mode := strings.ToLower(strings.TrimSpace(value))
	if mode == "" {
		return "append"
	}
	return mode
}
func resolveHeader(value *bool) bool {
	if value == nil {
		return true
	}
	return *value
}
func validateManifestEntry(sectionDir string, entry ManifestFile) error {
	if strings.TrimSpace(entry.File) == "" {
		return errors.New("file is required")
	}
	if strings.TrimSpace(entry.Table) == "" {
		return errors.New("table is required")
	}
	if filepath.IsAbs(entry.File) {
		return fmt.Errorf("absolute file path is forbidden: %s", entry.File)
	}
	cleaned := filepath.Clean(entry.File)
	if strings.HasPrefix(cleaned, "..") {
		return fmt.Errorf("file escapes section directory: %s", entry.File)
	}
	csvPath := filepath.Join(sectionDir, cleaned)
	if filepath.Ext(csvPath) != ".csv" {
		return fmt.Errorf("file must be .csv: %s", entry.File)
	}
	if _, err := os.Stat(csvPath); err != nil {
		return fmt.Errorf("csv not found: %s", entry.File)
	}
	if err := validateIdentifier(entry.Table); err != nil {
		return fmt.Errorf("invalid table %q: %w", entry.Table, err)
	}
	if len(entry.Columns) > 0 {
		if err := validateIdentifiers(entry.Columns); err != nil {
			return err
		}
	}
	return nil
}
func normalizeColumns(header []string) []string {
	out := make([]string, 0, len(header))
	for _, col := range header {
		col = strings.TrimSpace(col)
		col = strings.TrimPrefix(col, "\ufeff")
		out = append(out, col)
	}
	return out
}
func validateIdentifiers(values []string) error {
	for _, value := range values {
		if err := validateIdentifier(value); err != nil {
			return err
		}
	}
	return nil
}
func validateIdentifier(value string) error {
	if !identRE.MatchString(value) {
		return fmt.Errorf("%q does not match postgres-safe identifier policy", value)
	}
	return nil
}
func quoteIdents(values []string) []string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		out = append(out, quoteIdent(value))
	}
	return out
}
func quoteIdent(value string) string {
	return `"` + strings.ReplaceAll(value, `"`, `""`) + `"`
}
func escapeCopyLiteral(value string) string {
	return strings.ReplaceAll(value, `'`, `''`)
}
func jsonUnmarshalStrict(data []byte, dst any) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return err
	}
	var extra any
	if err := dec.Decode(&extra); err == nil {
		return errors.New("manifest must contain a single JSON object")
	}
	return nil
}
