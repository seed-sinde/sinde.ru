include .env
GOOSE_DIR=./db/migrations
strip_quotes = $(patsubst "%,%,$(patsubst %",%,$(strip $(1))))
PG_HOST := $(call strip_quotes,$(POSTGRES_HOST))
PG_PORT := $(call strip_quotes,$(POSTGRES_PORT))
PG_USER := $(call strip_quotes,$(POSTGRES_USER))
PG_PASSWORD := $(call strip_quotes,$(POSTGRES_PASSWORD))
PG_DB := $(call strip_quotes,$(POSTGRES_DB))
PG_SSL := $(call strip_quotes,$(POSTGRES_SSL))
DB_DSN := postgres://$(PG_USER):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$(PG_DB)?sslmode=$(PG_SSL)

.PHONY: help goose-up goose-down goose-status goose-create goose-redo goose-reset build-seeds-import

help:
	@echo "Available commands:"
	@grep -E '^[a-zA-Z0-9_-]+:.*##' Makefile | \
		awk 'BEGIN {FS=":.*##"} {printf "  %-20s %s\n", $$1, $$2}'

goose-up: ## Apply all migrations
	goose -dir $(GOOSE_DIR) postgres "$(DB_DSN)" up

goose-down: ## Rollback last migration
	goose -dir $(GOOSE_DIR) postgres "$(DB_DSN)" down

goose-status: ## Show migration status
	goose -dir $(GOOSE_DIR) postgres "$(DB_DSN)" status

goose-create: ## Create new migration (use name=...)
	goose -dir $(GOOSE_DIR) create $(name) sql

goose-redo: ## Redo last migration
	goose -dir $(GOOSE_DIR) postgres "$(DB_DSN)" redo

goose-reset: ## Drop all tables and re-run migrations
	goose -dir $(GOOSE_DIR) postgres "$(DB_DSN)" reset
build-seeds-import: ## Build seeds import binary
	go build -buildvcs=true -trimpath -o ./seeds-import ./cmd/seeds-import
