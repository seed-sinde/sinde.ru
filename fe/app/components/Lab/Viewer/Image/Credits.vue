<template>
  <div v-if="hasContent" class="viewer-credits mt-2 space-y-0.5 text-[11px] leading-4">
    <p class="line-clamp-2 space-x-2">
      <span v-if="title" class="viewer-credits-title text-sm font-semibold sm:text-base">
        {{ title }}
      </span>
      <span v-if="attribution">{{ attribution }}</span>
      <span v-if="author">© {{ author }}</span>
      <NuxtLink
        v-if="sourceUrl"
        :to="sourceUrl"
        class="viewer-credits-link underline underline-offset-2"
        target="_blank"
        rel="noopener noreferrer">
        источник
      </NuxtLink>
      <span v-if="license">
        <NuxtLink
          v-if="licenseUrl"
          :to="licenseUrl"
          class="viewer-credits-link underline underline-offset-2"
          target="_blank"
          rel="noopener noreferrer">
          {{ license }}
        </NuxtLink>
        <span v-else>{{ license }}</span>
      </span>
    </p>
  </div>
</template>
<script setup lang="ts">
  const props = defineProps<{
    title?: string | null
    author?: string | null
    attribution?: string | null
    sourceUrl?: string | null
    licenseUrl?: string | null
    license?: string | null
  }>()
  const hasContent = computed(() =>
    Boolean(props.title || props.author || props.attribution || props.sourceUrl || props.license)
  )
</script>
<style scoped>
  .viewer-credits {
    color: var(--lab-text-secondary);
  }
  .viewer-credits-title {
    color: var(--lab-text-primary);
  }
  .viewer-credits-link {
    color: var(--lab-accent);
    text-decoration-color: color-mix(in srgb, var(--lab-accent) 48%, transparent);
    transition: color 0.2s ease, text-decoration-color 0.2s ease;
  }
  .viewer-credits-link:hover {
    color: var(--lab-accent-hover);
    text-decoration-color: color-mix(in srgb, var(--lab-accent-hover) 56%, transparent);
  }
</style>
