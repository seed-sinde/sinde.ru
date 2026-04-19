# SINDE

Публичный репозиторий приложения `sinde.ru`.

В репозитории лежат:

- `be/` — backend на Go;
- `fe/` — frontend на Nuxt 4, pnpm и Tailwind CSS v4;
- `docker-compose.yml` — локальные зависимости для разработки: PostgreSQL, Redis, MinIO;
- `docs/` — рабочая документация по запуску, вкладу в проект и безопасности.

## Что умеет проект

- аккаунт и аутентификация: регистрация, вход, восстановление доступа, 2FA;
- кабинет пользователя и настройки интерфейса;
- traits / sets;
- kitchen recipes и связанные инструменты;
- wiki и учебные страницы;
- chemistry, minerals, astronomy и другие прикладные разделы;
- загрузка медиа для аватаров и рецептов.

## Требования

- Node.js `24.x`
- pnpm `10.x`
- Go `1.26.1`
- Docker и `docker compose`
- PostgreSQL, Redis и MinIO для локальной разработки
- `goose` для миграций backend

## Быстрый запуск

1. Скопировать шаблоны env:

```bash
cp .env.example .env
cp be/.env.example be/.env
cp fe/.env.example fe/.env
```

2. Поднять локальные сервисы:

```bash
docker compose up -d
```

3. Запустить backend:

```bash
cd be
./dev.sh
```

4. В другом терминале запустить frontend:

```bash
cd fe
pnpm install
pnpm dev
```

Приложение по умолчанию доступно на `http://127.0.0.1:3000`, backend на `http://127.0.0.1:3001`.

## Основные команды

Backend:

- `cd be && ./dev.sh` — локальный запуск backend;
- `cd be && ./dev.sh migrate` — применить миграции;
- `cd be && ./dev.sh seed` — импортировать сиды;
- `cd be && go test ./...` — тесты backend.

Frontend:

- `cd fe && pnpm dev` — локальный запуск frontend;
- `cd fe && pnpm typecheck` — проверка типов;
- `cd fe && pnpm lint` — frontend lint;
- `cd fe && pnpm build` — production build;

## Структура

- `be/cmd` — точки входа backend и служебных CLI;
- `be/db` — миграции и сиды;
- `be/internal` — бизнес-логика, HTTP-обработчики и сервисы;
- `fe/app` — страницы, компоненты, composables и data;
- `fe/public` — публичные статические ассеты;
- `docs` — документация для разработки и публикации.

## Документация

- [Обзор документации](docs/README.md)
- [Как вносить изменения](docs/CONTRIBUTING.md)
- [Безопасность](docs/SECURITY.md)
- [История изменений](docs/CHANGELOG.md)
- [README фронтенда](fe/README.md)

## Версия

Текущая версия продукта хранится в [VERSION](VERSION).
Изменение `VERSION` в `main` после успешного CI автоматически создаёт release tag `vX.Y.Z` и запускает production deploy через `seed-sinde/sinde.ru-ops`.

## Лицензия

Лицензия опубликована в [LICENSE](LICENSE).
