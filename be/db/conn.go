package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"sinde.ru/utils"
)

var PDB *pgxpool.Pool
var RDB *redis.Client

func Setup() (*pgxpool.Pool, *redis.Client, error) {
	cfg := utils.Config
	if cfg == nil {
		log.Fatal("Environment configuration is not initialized.")
	}
	pgCfg, err := pgxpool.ParseConfig(fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDB,
		cfg.PostgresSSL,
	))
	if err != nil {
		return nil, nil, err
	}
	pgCfg.MaxConns = 10
	pgCfg.MinConns = 2
	pgCfg.MaxConnLifetime = time.Hour
	pgCfg.MaxConnIdleTime = 30 * time.Minute

	pdb, err := pgxpool.NewWithConfig(context.Background(), pgCfg)
	if err != nil {
		return nil, nil, err
	}
	if err = pdb.Ping(context.Background()); err != nil {
		return nil, nil, err
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = rdb.Ping(ctx).Err(); err != nil {
		return nil, nil, err
	}
	PDB = pdb
	RDB = rdb
	return pdb, rdb, nil
}

func Init() error {
	_, _, err := Setup()
	if err != nil {
		return err
	}
	return nil
}
