package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"sinde.ru/utils"
	"time"
)

var PDB *pgxpool.Pool

func Setup() (*pgxpool.Pool, error) {
	config := utils.Config
	if config == nil {
		log.Fatal("Environment configuration is not initilized.")
	}
	dbConfig, err := pgxpool.ParseConfig(fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.PostgresHost,
		config.PostgresPort,
		config.PostgresUser,
		config.PostgresPassword,
		config.PostgresDB,
		config.PostgresSSL,
	))
	if err != nil {
		log.Fatalf("Unable to parse DATABASE_URL: %v\n", err)
	}
	dbConfig.MaxConns = 10
	dbConfig.MinConns = 2
	dbConfig.MaxConnLifetime = time.Hour
	dbConfig.MaxConnIdleTime = 30 * time.Minute
	PDB, err = pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}
	err = PDB.Ping(context.Background())
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	return PDB, nil
}
func Init() (err error) {
	PDB, err = Setup()
	if err != nil {
		log.Fatalf("Connection error PostgreSQL: %v", err)
	}
	return nil
}
