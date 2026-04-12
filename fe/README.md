# Фронтенд

Frontend приложения `sinde.ru` на Nuxt 4, Bun, Vite, Nitro и Tailwind CSS v4.

## Назначение

Во frontend находятся:

- публичные страницы проекта;
- страницы аутентификации и аккаунта;
- прикладные разделы `traits`, `kitchen`, `chemistry`, `minerals`, `astronomy`;
- SSR proxy до backend API;
- PWA-конфигурация и клиентская часть UI.

## Требования

- Bun `1.3.11`
- файл `fe/.env`
- запущенный backend на `http://127.0.0.1:3001`

Минимальная локальная подготовка:

```bash
cp .env.example .env
bun install
```

## Скрипты

- `bun run dev` — локальный запуск Nuxt с HMR;
- `bun run typecheck` — `nuxt typecheck`;
- `bun run build` — production build;
- `bun run preview` — локальный просмотр production build;
- `bun run start` — запуск `.output/server/index.mjs` через Bun;
- `bun run generate` — генерация статического вывода;
- `bun run upd` — показать список обновлений зависимостей через `npm-check-updates`.

## Переменные окружения

Основные переменные в [.env.example](.env.example):

- `API_INTERNAL_URL`
- `NUXT_PUBLIC_API_URL`
- `NUXT_PUBLIC_BASE_URL`
- `NUXT_PUBLIC_MEDIA_IMAGE_MAX_BYTES`
- `NITRO_HOST`
- `NITRO_PORT`

## PWA

- PWA включается для production build;
- для локальной проверки используй `bun run build && bun run preview`;
- сервис-воркер и manifest корректнее проверять на HTTPS-хосте, а не только на `localhost`.
