<template>
  <div v-if="hasContent" class="mt-2 space-y-0.5 text-xs leading-4 text-(--lab-text-muted)">
    <p class="line-clamp-2 space-x-2">
      <span v-if="title" class="text-sm font-semibold text-(--lab-text-primary) sm:text-base">
        {{ title }}
      </span>
      <span v-if="attribution">{{ attribution }}</span>
      <span v-if="author">© {{ author }}</span>
      <NuxtLink
        v-if="sourceUrl"
        :to="sourceUrl"
        class="text-(--lab-accent) underline underline-offset-2 decoration-[color-mix(in_srgb,var(--lab-accent)_48%,transparent)] transition hover:text-(--lab-accent-hover) hover:decoration-[color-mix(in_srgb,var(--lab-accent-hover)_56%,transparent)]"
        target="_blank"
        rel="noopener noreferrer">
        источник
      </NuxtLink>
      <span v-if="license">
        <NuxtLink
          v-if="licenseUrl"
          :to="licenseUrl"
          class="text-(--lab-accent) underline underline-offset-2 decoration-[color-mix(in_srgb,var(--lab-accent)_48%,transparent)] transition hover:text-(--lab-accent-hover) hover:decoration-[color-mix(in_srgb,var(--lab-accent-hover)_56%,transparent)]"
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
    author?: string | null | undefined
    attribution?: string | null | undefined
    sourceUrl?: string | null | undefined
    licenseUrl?: string | null | undefined
    license?: string | null | undefined
  }>()
  const hasContent = computed(() =>
    Boolean(props.title || props.author || props.attribution || props.sourceUrl || props.license)
  )
</script>
