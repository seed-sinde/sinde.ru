<script setup lang="ts">
  import type { WikiTool } from '~/constants/wikiTools'
  const props = withDefaults(
    defineProps<{
      tool: WikiTool
      compact?: boolean
    }>(),
    {
      compact: false
    }
  )
  const cardStyle = computed(() => ({
    '--wiki-accent': props.tool.accentColor,
    '--wiki-accent-soft': props.tool.accentSoftColor
  }))
</script>
<template>
  <NuxtLink
    :to="tool.to"
    class="wiki-tool-card group flex h-full flex-col border transition-colors"
    :class="compact ? 'wiki-tool-card-compact' : 'wiki-tool-card-regular'"
    :style="cardStyle">
    <div class="wiki-tool-card__body flex min-w-0 flex-1 flex-col gap-4">
      <div class="flex items-start justify-between gap-3">
        <div class="min-w-0 space-y-3">
          <div class="wiki-tool-card__badge inline-flex items-center gap-1.5 px-2.5 py-1 text-xs font-medium">
            <span class="wiki-tool-card__badge-dot h-2 w-2 shrink-0 rounded-full"></span>
            <span class="truncate">{{ tool.categoryLabel }}</span>
          </div>
          <div class="flex items-start gap-3">
            <span class="wiki-tool-card__icon inline-flex h-11 w-11 shrink-0 items-center justify-center">
              <Icon :name="tool.icon" class="h-5 w-5" />
            </span>
            <div class="min-w-0 space-y-2">
              <h3 class="text-lg font-semibold leading-tight lab-text-primary">{{ tool.title }}</h3>
              <p class="wrap-break-word text-sm leading-6 lab-text-secondary">{{ tool.text }}</p>
            </div>
          </div>
        </div>
        <div class="wiki-tool-card__art shrink-0" :class="`wiki-tool-card__art-${tool.art}`" aria-hidden="true">
          <svg v-if="tool.art === 'elements'" viewBox="0 0 120 120" class="h-full w-full">
            <rect x="14" y="18" width="34" height="34" rx="10" class="wiki-tool-card__shape-main" />
            <rect x="56" y="18" width="24" height="24" rx="8" class="wiki-tool-card__shape-soft" />
            <rect x="84" y="18" width="22" height="22" rx="8" class="wiki-tool-card__shape-stroke" />
            <rect x="28" y="60" width="26" height="26" rx="8" class="wiki-tool-card__shape-soft" />
            <rect x="60" y="50" width="44" height="44" rx="12" class="wiki-tool-card__shape-main" />
            <path d="M18 92h88" class="wiki-tool-card__line" />
          </svg>
          <svg v-else-if="tool.art === 'crystal'" viewBox="0 0 120 120" class="h-full w-full">
            <path d="M60 14 86 34 76 88 44 88 34 34Z" class="wiki-tool-card__shape-main" />
            <path d="M34 34 60 14 48 56Z" class="wiki-tool-card__shape-soft" />
            <path d="M60 14 86 34 72 56Z" class="wiki-tool-card__shape-stroke" />
            <path d="M48 56h24L60 88Z" class="wiki-tool-card__shape-soft" />
            <circle cx="92" cy="86" r="12" class="wiki-tool-card__shape-stroke" />
          </svg>
          <svg v-else-if="tool.art === 'compounds'" viewBox="0 0 120 120" class="h-full w-full">
            <circle cx="26" cy="60" r="12" class="wiki-tool-card__shape-main" />
            <circle cx="60" cy="30" r="10" class="wiki-tool-card__shape-soft" />
            <circle cx="92" cy="60" r="13" class="wiki-tool-card__shape-main" />
            <circle cx="60" cy="88" r="9" class="wiki-tool-card__shape-stroke" />
            <path d="M34 54 52 36M68 36 84 53M33 67 53 83M67 82 84 67" class="wiki-tool-card__line" />
          </svg>
          <svg v-else-if="tool.art === 'isotopes'" viewBox="0 0 120 120" class="h-full w-full">
            <circle cx="60" cy="60" r="16" class="wiki-tool-card__shape-main" />
            <circle cx="60" cy="60" r="34" class="wiki-tool-card__ring" />
            <circle cx="60" cy="60" r="48" class="wiki-tool-card__ring-soft" />
            <circle cx="60" cy="26" r="6" class="wiki-tool-card__shape-soft" />
            <circle cx="28" cy="76" r="7" class="wiki-tool-card__shape-stroke" />
            <circle cx="92" cy="76" r="5" class="wiki-tool-card__shape-soft" />
          </svg>
          <svg v-else-if="tool.art === 'materials'" viewBox="0 0 120 120" class="h-full w-full">
            <rect x="18" y="20" width="32" height="32" rx="8" class="wiki-tool-card__shape-main" />
            <path d="M66 18h30v30H66z" class="wiki-tool-card__shape-stroke" />
            <path d="M20 74h36L38 102H20z" class="wiki-tool-card__shape-soft" />
            <circle cx="86" cy="84" r="16" class="wiki-tool-card__shape-main" />
            <path d="M70 84h32M86 68v32" class="wiki-tool-card__line" />
          </svg>
          <svg v-else-if="tool.art === 'products'" viewBox="0 0 120 120" class="h-full w-full">
            <path d="M26 34h68l-8 52H34Z" class="wiki-tool-card__shape-main" />
            <path d="M40 34c0-11 8-18 20-18s20 7 20 18" class="wiki-tool-card__shape-stroke" />
            <path d="M42 56h36M42 68h28M42 80h22" class="wiki-tool-card__line" />
            <rect x="82" y="58" width="14" height="14" rx="4" class="wiki-tool-card__shape-soft" />
          </svg>
          <svg v-else-if="tool.art === 'production'" viewBox="0 0 120 120" class="h-full w-full">
            <circle cx="38" cy="60" r="15" class="wiki-tool-card__shape-main" />
            <circle cx="82" cy="60" r="11" class="wiki-tool-card__shape-soft" />
            <path
              d="M38 37v10M38 73v10M15 60h10M51 60h10M24 46l7 7M45 67l7 7M24 74l7-7M45 53l7-7M82 43v8M82 69v8M65 60h8M91 60h8"
              class="wiki-tool-card__line" />
          </svg>
          <svg v-else viewBox="0 0 120 120" class="h-full w-full">
            <circle cx="60" cy="60" r="16" class="wiki-tool-card__shape-main" />
            <circle cx="60" cy="60" r="34" class="wiki-tool-card__ring" />
            <circle cx="60" cy="60" r="50" class="wiki-tool-card__ring-soft" />
            <circle cx="94" cy="44" r="6" class="wiki-tool-card__shape-soft" />
            <circle cx="26" cy="74" r="8" class="wiki-tool-card__shape-stroke" />
            <path d="M18 28h10M92 92h10M26 16l6 6" class="wiki-tool-card__line" />
          </svg>
        </div>
      </div>
      <div class="wiki-tool-card__footer mt-auto inline-flex items-center gap-1.5 text-sm font-medium">
        <span>Открыть раздел</span>
        <Icon name="ic:round-arrow-forward" class="h-4 w-4" />
      </div>
    </div>
  </NuxtLink>
</template>
<style scoped>
  .wiki-tool-card {
    border-color: color-mix(in srgb, var(--wiki-accent) 22%, var(--lab-border));
    background:
      radial-gradient(
        circle at top right,
        color-mix(in srgb, var(--wiki-accent-soft) 26%, transparent) 0,
        transparent 54%
      ),
      linear-gradient(
        160deg,
        color-mix(in srgb, var(--wiki-accent) 10%, var(--lab-bg-surface)) 0%,
        color-mix(in srgb, var(--wiki-accent-soft) 6%, var(--lab-bg-surface)) 100%
      );
    overflow: hidden;
  }
  .wiki-tool-card:hover {
    border-color: color-mix(in srgb, var(--wiki-accent) 44%, var(--lab-border-strong));
    background:
      radial-gradient(
        circle at top right,
        color-mix(in srgb, var(--wiki-accent-soft) 32%, transparent) 0,
        transparent 58%
      ),
      linear-gradient(
        160deg,
        color-mix(in srgb, var(--wiki-accent) 14%, var(--lab-bg-surface-hover)) 0%,
        color-mix(in srgb, var(--wiki-accent-soft) 9%, var(--lab-bg-surface-hover)) 100%
      );
  }
  .wiki-tool-card-regular .wiki-tool-card__body {
    padding: 1.25rem;
  }
  .wiki-tool-card-compact .wiki-tool-card__body {
    padding: 1rem;
  }
  .wiki-tool-card__badge {
    border: 1px solid color-mix(in srgb, var(--wiki-accent) 28%, var(--lab-border));
    background: color-mix(in srgb, var(--wiki-accent) 12%, var(--lab-bg-control));
    color: color-mix(in srgb, var(--wiki-accent) 72%, var(--lab-text-primary));
  }
  .wiki-tool-card__badge-dot {
    background: var(--wiki-accent);
  }
  .wiki-tool-card__icon {
    border: 1px solid color-mix(in srgb, var(--wiki-accent) 28%, var(--lab-border));
    background:
      radial-gradient(circle at top, color-mix(in srgb, var(--wiki-accent-soft) 34%, transparent) 0, transparent 62%),
      color-mix(in srgb, var(--wiki-accent) 12%, var(--lab-bg-surface-muted));
    color: color-mix(in srgb, var(--wiki-accent) 76%, var(--lab-text-primary));
  }
  .wiki-tool-card__art {
    width: 6.5rem;
    height: 6.5rem;
    color: var(--wiki-accent);
    opacity: 0.96;
  }
  .wiki-tool-card__footer {
    color: color-mix(in srgb, var(--wiki-accent) 82%, var(--lab-text-primary));
  }
  .wiki-tool-card__shape-main,
  .wiki-tool-card__shape-soft,
  .wiki-tool-card__shape-stroke,
  .wiki-tool-card__ring,
  .wiki-tool-card__ring-soft,
  .wiki-tool-card__line {
    vector-effect: non-scaling-stroke;
  }
  .wiki-tool-card__shape-main {
    fill: color-mix(in srgb, var(--wiki-accent) 40%, transparent);
    stroke: color-mix(in srgb, var(--wiki-accent) 70%, transparent);
    stroke-width: 4;
  }
  .wiki-tool-card__shape-soft {
    fill: color-mix(in srgb, var(--wiki-accent-soft) 34%, transparent);
    stroke: color-mix(in srgb, var(--wiki-accent-soft) 68%, transparent);
    stroke-width: 3;
  }
  .wiki-tool-card__shape-stroke {
    fill: transparent;
    stroke: color-mix(in srgb, var(--wiki-accent) 58%, transparent);
    stroke-width: 3;
  }
  .wiki-tool-card__ring {
    fill: none;
    stroke: color-mix(in srgb, var(--wiki-accent) 50%, transparent);
    stroke-width: 4;
    stroke-dasharray: 10 9;
  }
  .wiki-tool-card__ring-soft {
    fill: none;
    stroke: color-mix(in srgb, var(--wiki-accent-soft) 44%, transparent);
    stroke-width: 3;
    stroke-dasharray: 6 8;
  }
  .wiki-tool-card__line {
    fill: none;
    stroke: color-mix(in srgb, var(--wiki-accent-soft) 58%, transparent);
    stroke-linecap: round;
    stroke-width: 3;
  }
  :global(.dark) .wiki-tool-card {
    background:
      radial-gradient(
        circle at top right,
        color-mix(in srgb, var(--wiki-accent-soft) 16%, transparent) 0,
        transparent 54%
      ),
      linear-gradient(
        160deg,
        color-mix(in srgb, var(--wiki-accent) 9%, var(--lab-bg-surface)) 0%,
        color-mix(in srgb, var(--wiki-accent-soft) 5%, var(--lab-bg-surface)) 100%
      );
  }
  :global(.dark) .wiki-tool-card:hover {
    background:
      radial-gradient(
        circle at top right,
        color-mix(in srgb, var(--wiki-accent-soft) 20%, transparent) 0,
        transparent 58%
      ),
      linear-gradient(
        160deg,
        color-mix(in srgb, var(--wiki-accent) 11%, var(--lab-bg-surface-hover)) 0%,
        color-mix(in srgb, var(--wiki-accent-soft) 6%, var(--lab-bg-surface-hover)) 100%
      );
  }
  @media (max-width: 639px) {
    .wiki-tool-card__art {
      width: 5rem;
      height: 5rem;
    }
    .wiki-tool-card-regular .wiki-tool-card__body,
    .wiki-tool-card-compact .wiki-tool-card__body {
      padding: 1rem;
    }
  }
</style>
