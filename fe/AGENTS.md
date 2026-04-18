# AGENTS.md

## Stack

node 24, pnpm, nuxt, vite, nitro, vue, ts, tailwind v4

## Rules & Code

- script setup, ts, auto-import, no-any, Nuxt naming.
- Arrow fns, expression-body, ternaries, early returns.
- Minimal variables, no single-use temps, no duplicate calcs.
- pnpm: format, typecheck, lint. No SSR breaking.

## Template & UI

- Clean DOM: minimal containers/blocks, local v-if.
- Style: minimalism, high contrast, readable, no visual noise.
- Forbidden: rounded*, shadow*, backdrop-blur, gradients, text-inherit (if color set).
- Exception: rounded-full (icon-only buttons), lab-scroll-fade\*.

## Tailwind v4 & CSS

- Use utility: bg-(--var), scale, wrap-break-word, aspect-4/3. No custom helpers.
- Colors: strictly via theme tokens or `color-mix()` in component CSS.
- Files: `assets/css/*.css` in kebab-case by domain.
- Focus: `lab-focus` (base), `lab-focus-manual`, `lab-focus-peer`, `lab-tabs-focus` (bottom accent).
