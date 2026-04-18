# AGENTS.md

## Stack

go, fiber, postgres, redis

## Principles & Code

- Flat & simple: short functions, early returns, no deep nesting.
- No over-engineering: avoid extra layers, structures, or abstractions.
- Dry: don't duplicate logic; make minimal, surgical changes.
- Performance: avoid unnecessary allocations and redundant data loops.

## Database & Redis

- DB: Idempotent changes, no breaking existing data, minimal migrations.
- Redis: Only for cache/state; don't duplicate Postgres logic.

## API & Errors

- API: Simple contracts, no extra layers or over-abstraction.
- Errors: Don't silence; keep error handling logic minimal and direct.
