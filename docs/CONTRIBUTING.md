# Как вносить изменения

## Общие правила

- изменения должны быть целевыми и минимальными;
- не добавляй секреты, реальные `.env` и production credentials;
- не дублируй существующую логику и типы без причины;
- если в проекте уже есть компонент, утилита или composable для задачи, используй его;
- пользовательские тексты в интерфейсе по умолчанию должны быть на русском, если задача не про локализацию.
- постоянные архитектурные правила описаны в [PROJECT_STANDARDS](PROJECT_STANDARDS.md).

## Frontend

- стек: Nuxt 4, Vue 3, pnpm, Tailwind CSS v4;
- использовать `script setup lang="ts"`;
- не писать ручные import для Nuxt auto-import сущностей;
- все запросы к <abbr title="Application Programming Interface">API</abbr> делать через `useAPI().json`, `useAPI().stream` или существующий доменный composable;
- новые иконки использовать из набора `ic:round-*`;
- проверка перед push:

```bash
cd fe
pnpm install
pnpm typecheck
pnpm lint
```

## Backend

- стек: Go `1.26.1`, PostgreSQL, Redis, MinIO;
- миграции именуются как `YYYYMMDDHHMMSS_description.sql`;
- локальная работа и обслуживание сидов идут через `be/dev.sh`;
- проверка перед push:

```bash
cd be
go test ./...
```

## Git

- основная ветка: `main`;
- коммиты должны быть понятными и по делу;
- если меняется публичное поведение проекта или документация, обновляй [docs/CHANGELOG.md](CHANGELOG.md).
