# Фронтенд

Frontend приложения `sinde.ru` на Nuxt 4, pnpm, Vite, Nitro и Tailwind CSS v4.

## Назначение

Во frontend находятся:

- публичные страницы проекта;
- страницы аутентификации и аккаунта;
- прикладные разделы `traits`, `kitchen`, `chemistry`, `minerals`;
- SSR proxy до backend API;
- PWA-конфигурация и клиентская часть UI.

## Требования

- Node.js `24.x`
- pnpm `10.x`
- файл `fe/.env`
- запущенный backend на `http://127.0.0.1:3001`

Минимальная локальная подготовка:

```bash
cp .env.example .env
pnpm install
```

## Скрипты

- `pnpm dev` — локальный запуск Nuxt с HMR;
- `pnpm typecheck` — `nuxt typecheck`;
- `pnpm lint` — eslint-проверка;
- `pnpm build` — production build;
- `pnpm preview` — локальный просмотр production build;
- `pnpm start` — запуск `.output/server/index.mjs` через Node.js;
- `pnpm generate` — генерация статического вывода;
- `pnpm upd` — показать список обновлений зависимостей через `npm-check-updates`.

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
- для локальной проверки используй `pnpm build && pnpm preview`;
- сервис-воркер и manifest корректнее проверять на HTTPS-хосте, а не только на `localhost`.
