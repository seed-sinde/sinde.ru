package utils

import (
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
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
	RedisAddr        string
	RedisPassword    string
	RedisDB          int
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

func tryLoadDotenv() {
	candidates := []string{".env.local", ".env.production", ".env"}
	if loadFirstExisting(".", candidates...) {
		return
	}
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

func GetEnvVar() *ConfigT {
	tryLoadDotenv()

	get := func(k string) string { return strings.TrimSpace(os.Getenv(k)) }

	values := map[string]string{
		"POSTGRES_HOST":     get("POSTGRES_HOST"),
		"POSTGRES_PORT":     get("POSTGRES_PORT"),
		"POSTGRES_USER":     get("POSTGRES_USER"),
		"POSTGRES_PASSWORD": get("POSTGRES_PASSWORD"),
		"POSTGRES_DB":       get("POSTGRES_DB"),
		"POSTGRES_SSL":      get("POSTGRES_SSL"),

		"REDIS_ADDR":     get("REDIS_ADDR"),
		"REDIS_PASSWORD": get("REDIS_PASSWORD"),
		"REDIS_DB":       get("REDIS_DB"),
	}

	missing := make([]string, 0, 6)
	for _, key := range []string{
		"POSTGRES_HOST",
		"POSTGRES_PORT",
		"POSTGRES_USER",
		"POSTGRES_PASSWORD",
		"POSTGRES_DB",
		"POSTGRES_SSL",
		"REDIS_ADDR",
		"REDIS_PASSWORD",
		"REDIS_DB",
	} {
		if values[key] == "" {
			missing = append(missing, key)
		}
	}

	if len(missing) > 0 {
		sort.Strings(missing)
		log.Printf("Неполная конфигурация PostgreSQL, Redis. Отсутствуют: %s", strings.Join(missing, ", "))
		return nil
	}

	redisDB := 0
	if v := values["REDIS_DB"]; v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			redisDB = n
		} else {
			log.Printf("REDIS_DB должен быть числом, получено: %s", v)
		}
	}

	return &ConfigT{
		PostgresHost:     values["POSTGRES_HOST"],
		PostgresPort:     values["POSTGRES_PORT"],
		PostgresUser:     values["POSTGRES_USER"],
		PostgresPassword: values["POSTGRES_PASSWORD"],
		PostgresDB:       values["POSTGRES_DB"],
		PostgresSSL:      values["POSTGRES_SSL"],

		RedisAddr:     values["REDIS_ADDR"],
		RedisPassword: values["REDIS_PASSWORD"],
		RedisDB:       redisDB,
	}
}

func LoadEnv() {
	Config = GetEnvVar()
	if Config == nil {
		log.Fatal("Не удалось загрузить переменные среды.")
	}
}
