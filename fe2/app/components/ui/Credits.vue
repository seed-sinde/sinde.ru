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

<template>
  <div v-if="hasContent" class="mt-2 text-xs leading-4 text-zinc-700 dark:text-zinc-300">
    <slot v-bind="props">
      <p class="line-clamp-2 space-x-2">
        <span v-if="title" class="text-sm font-semibold text-(--text) sm:text-base">
          {{ title }}
        </span>
        <span v-if="attribution">{{ attribution }}</span>
        <span v-if="author">© {{ author }}</span>
        <NuxtLink
          v-if="sourceUrl"
          :to="sourceUrl"
          class="ui-focus text-(--accent) underline underline-offset-2 transition hover:ring-2 hover:ring-(--accent)"
          target="_blank"
          rel="noopener noreferrer"
        >
          источник
        </NuxtLink>
        <span v-if="license">
          <NuxtLink
            v-if="licenseUrl"
            :to="licenseUrl"
            class="ui-focus text-(--accent) underline underline-offset-2 transition hover:ring-2 hover:ring-(--accent)"
            target="_blank"
            rel="noopener noreferrer"
          >
            {{ license }}
          </NuxtLink>
          <span v-else>{{ license }}</span>
        </span>
      </p>
    </slot>
  </div>
</template>
