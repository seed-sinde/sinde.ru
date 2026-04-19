# Стандарты развития проекта

Документ фиксирует текущие правила для `sinde.ru`. Он описывает уже принятые паттерны проекта и нужен, чтобы новые изменения не расползались по разным стилям.

## Общий порядок изменений

- Сначала ищи существующий модуль, composable, handler, service или тип.
- Новую сущность добавляй только если она убирает реальное дублирование или закрепляет уже повторяющийся паттерн.
- Не размещай business logic одновременно в странице, composable и utility. Выбери один уровень владения.
- Удаляй generated artifacts и cache-файлы из git; `*.tsbuildinfo`, бинарники и runtime cache не коммитятся.
- Не добавляй зависимости, пока задача не упирается в устойчивую проектную потребность.

## Frontend

- Стек: Nuxt 4, Vue 3, TypeScript, <abbr title="Server-Side Rendering">SSR</abbr>, Tailwind CSS v4, pnpm.
- Все новые компоненты пишутся через `script setup lang="ts"`.
- Auto-import считается доступным для `components/`, `composables/`, `shared/types`, `shared/utils`, `app/data`, `app/constants`, `app/utils`.
- Не импортируй вручную auto-import сущности, если проект уже поднимает их сам.
- Для <abbr title="Application Programming Interface">API</abbr> используй `useAPI().json` и `useAPI().stream`; прямой `$fetch`/`fetch` допустим только внутри самой API-обертки или server middleware/proxy.
- SSR-запросы должны идти через front proxy `/api/proxy`, чтобы сохранялись cookies, csrf, locale и версия backend API.
- Ошибки <abbr title="Application Programming Interface">API</abbr> в <abbr title="User Interface">UI</abbr> нормализуй через `extractApiErrorMessage`, если нет более узкого доменного helper.
- Route/query/params для tabs нормализуй через `normalizeTabRouteValue` и связанные helpers из `app/utils/tabRoute.ts`.
- <abbr title="Search Engine Optimization">SEO</abbr> для страниц задавай через `usePageSeo`; `useHead` оставляй для layout/canonical или нестандартных meta-структур.
- Суммы платежей форматируй через `paymentPresentation.ts`, не создавай локальные `Intl.NumberFormat`-копии.
- Для loading/error/empty используй Lab-компоненты (`LabNotify`, `LabErrorMessage`, `LabLoader`, `LabDataTable`) там, где они уже подходят.
- Избегай watcher/computed-цепочек, которые просто копируют значение в другое значение. Derived state должен быть `computed`, побочные эффекты должны быть в узком `watch`.
- Для client-only поведения проверяй `import.meta.client`; не допускай разных <abbr title="Server-Side Rendering">SSR</abbr>/<abbr title="Client-Side Rendering">CSR</abbr> initial values без явного fallback.
- Tailwind: не добавляй `rounded*`, `shadow*`, `backdrop-blur`, gradient utilities, `break-words`, `text-inherit` с явным цветом. Исключение: `rounded-full` для круглых icon-only controls.

## Backend

- Стек: Go, Fiber, PostgreSQL, Redis, MinIO.
- `internal/http/routes.go` владеет регистрацией `/api/vN` маршрутов.
- `internal/http/handlers/*` владеют transport-слоем: чтение params/body, auth context, status code, вызов service/store, `responses.Success`/`responses.Error`.
- `internal/*/service.go` владеет бизнес-операциями и не должен знать presentation JSON-форму ответа.
- `internal/*/repository.go` или `db/services/*.go` владеют доступом к хранилищам и не должны возвращать handler-specific payloads.
- `internal/models/*.go` содержит общие модели и view-структуры, если они переиспользуются между handlers/services.
- Новые <abbr title="Data Transfer Object">DTO</abbr> для одного transport endpoint держи рядом с handler/service, пока они не стали общими.
- Ошибки <abbr title="Hypertext Transfer Protocol">HTTP</abbr> возвращай через `internal/http/responses`; прямой `c.Status(...).JSON(...)` оставляй только для внешних webhook contracts, где нужен нестандартный plain-text ответ.
- В handler не добавляй долгую бизнес-логику. Если функция начинает валидировать несколько доменных правил или вызывает несколько хранилищ, перенеси это в service.

## <abbr title="Application Programming Interface">API</abbr> Versioning

- Текущая публичная backend версия регистрируется как `/api/v1`; frontend выбирает ее через `API_VERSION`.
- Breaking change: переименование/удаление поля ответа, изменение типа поля, изменение обязательности request field, изменение статуса или error contract для уже используемого endpoint.
- Non-breaking change: новое optional поле, новый endpoint, расширение enum с безопасным fallback, исправление внутренней реализации без изменения contract.
- Новую `/api/vN` создавай только для breaking changes. Внутреннюю business logic переиспользуй через service/repository, а differences держи в version-specific handlers/mappers.
- Frontend не должен hardcode backend origin или version в компонентах/страницах; все идет через `/api/proxy`.

## Migration Guide

1. Новая frontend-фича: добавь доменный composable рядом с существующими, используй `useAPI().json`/`stream`, типы бери из `shared/types` или ближайшего доменного файла.
2. Новая страница: route/query нормализуй helper-ами, SEO ставь через `usePageSeo`, loading/error/empty собирай на Lab-компонентах.
3. Новая backend-фича: зарегистрируй route в `/api/v1`, handler оставь тонким, доменную логику держи в service, storage code в repository/db service.
4. Новый response contract: сначала добавь/переиспользуй тип, затем проверь frontend consumer и proxy path.
5. Удаление старого кода: проверь `rg` по export/component/composable, удали generated artifacts, запусти lint/typecheck/build и backend tests.
6. Разделять файл стоит, когда в нем появились независимые состояния или владельцы. Объединять стоит, когда несколько файлов описывают одну маленькую сущность и заставляют прыгать по проекту без пользы.
