package utils
import (
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"github.com/joho/godotenv"
)
type ConfigT struct {
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	PostgresSSL      string
}
var Config *ConfigT
func loadFirstExisting(baseDir string, names ...string) bool {
	for _, name := range names {
		if err := godotenv.Load(filepath.Join(baseDir, name)); err == nil {
			return true
		}
	}
	return false
}
// Функция tryLoadDotenv пытается загрузить .env сначала из текущей рабочей директории, затем из директории бинарника.
// Ошибки загрузки не фатальны: переменные окружения могут быть заданы извне.
func tryLoadDotenv() {
	candidates := []string{".env.local", ".env.production", ".env"}
	// 1) Текущая рабочая директория
	if loadFirstExisting(".", candidates...) {
		return
	}
	// 2) Рядом с бинарником
	if exePath, err := os.Executable(); err == nil {
		exeDir := ""
		if realPath, err2 := filepath.EvalSymlinks(exePath); err2 == nil {
			exeDir = filepath.Dir(realPath)
		} else {
			exeDir = filepath.Dir(exePath)
		}
		_ = loadFirstExisting(exeDir, candidates...)
	}
}
// Функция GetEnvVar загружает конфигурацию из окружения (после попытки прочитать .env).
// Возвращает nil и пишет детали в лог, если чего-то не хватает.
func GetEnvVar() *ConfigT {
	tryLoadDotenv()
	get := func(k string) string {
		return strings.TrimSpace(os.Getenv(k))
	}
	values := map[string]string{
		"POSTGRES_HOST":     get("POSTGRES_HOST"),
		"POSTGRES_PORT":     get("POSTGRES_PORT"),
		"POSTGRES_USER":     get("POSTGRES_USER"),
		"POSTGRES_PASSWORD": get("POSTGRES_PASSWORD"),
		"POSTGRES_DB":       get("POSTGRES_DB"),
		"POSTGRES_SSL":      get("POSTGRES_SSL"),
	}
	missing := make([]string, 0, len(values))
	for key, value := range values {
		if value == "" {
			missing = append(missing, key)
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		log.Printf("Неполная конфигурация PostgreSQL. Отсутствуют: %s", strings.Join(missing, ", "))
		return nil
	}
	return &ConfigT{
		PostgresHost:     values["POSTGRES_HOST"],
		PostgresPort:     values["POSTGRES_PORT"],
		PostgresUser:     values["POSTGRES_USER"],
		PostgresPassword: values["POSTGRES_PASSWORD"],
		PostgresDB:       values["POSTGRES_DB"],
		PostgresSSL:      values["POSTGRES_SSL"],
	}
}
// Функция LoadEnv заполняет глобальную конфигурацию или завершает процесс с ошибкой.
func LoadEnv() {
	Config = GetEnvVar()
	if Config == nil {
		log.Fatal("Не удалось загрузить переменные среды.")
	}
}
