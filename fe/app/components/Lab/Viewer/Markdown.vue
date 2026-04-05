<script setup lang="ts">
  const props = defineProps<{
    source: string
  }>()
  const EMOJI_TOKEN_RE = /:emojione:[a-z0-9-]+:/i
  const emojiCollectionReady = ref(false)
  const ensureEmojiCollectionForSource = async (source: string) => {
    if (emojiCollectionReady.value) return
    if (!EMOJI_TOKEN_RE.test(String(source || ''))) return
    await ensureEmojiCollectionLoaded()
    emojiCollectionReady.value = true
  }
  await ensureEmojiCollectionForSource(props.source)
  watch(
    () => props.source,
    async next => {
      await ensureEmojiCollectionForSource(next)
    }
  )
  const blocks = computed<MarkdownViewerBlock[]>(() => {
    emojiCollectionReady.value
    return renderMarkdownToBlocks(props.source)
  })
</script>
<template>
  <article class="md-content">
    <template v-for="block in blocks" :key="block.key">
      <div v-if="block.type === 'html'" v-html="block.html"></div>
      <div v-else class="md-table-wrap">
        <LabDataTable
          :columns="block.columns"
          :rows="block.rows"
          row-key="id"
          empty-text="Таблица пуста."
          table-class="min-w-full text-sm"
          thead-class="sticky top-0 z-10 bg-zinc-950/95"
          tbody-class="bg-transparent"
          row-class="border-b border-zinc-800 align-top"
          max-height-class="max-h-none"
          container-class="bg-transparent">
          <template #cell="{ value }">
            <span class="inline-block min-w-full" v-html="String(value || '&mdash;')"></span>
          </template>
        </LabDataTable>
      </div>
    </template>
  </article>
</template>
<style scoped>
  @reference "~/assets/css/main.css";
  .md-content {
    @apply pb-20 text-zinc-300;
    line-height: 1.6;
  }
  .md-content :deep(.md-align-left) {
    @apply text-left;
  }
  .md-content :deep(.md-align-center) {
    @apply text-center;
  }
  .md-content :deep(.md-align-right) {
    @apply text-right;
  }
  .md-content :deep(.md-align-justify) {
    @apply text-justify;
  }
  .md-content :deep(h1) {
    @apply mt-4 mb-2 text-2xl font-semibold text-zinc-100;
  }
  .md-content :deep(h2) {
    @apply mt-4 mb-2 text-xl font-semibold text-zinc-100;
  }
  .md-content :deep(h3) {
    @apply mt-4 mb-2 text-lg font-semibold text-zinc-100;
  }
  .md-content :deep(h4),
  .md-content :deep(h5),
  .md-content :deep(h6) {
    @apply mt-4 mb-2 text-base font-semibold text-zinc-100;
  }
  .md-content :deep(p) {
    @apply mb-2;
  }
  .md-content :deep(ul) {
    @apply mb-3 list-disc pl-6;
  }
  .md-content :deep(ol) {
    @apply mb-3 list-decimal pl-6;
  }
  .md-content :deep(li) {
    @apply mb-1;
  }
  .md-content :deep(a) {
    @apply text-indigo-400 underline decoration-from-font;
  }
  .md-content :deep(blockquote) {
    @apply mb-3 border-l-[3px] border-indigo-500/70 pl-3 text-zinc-400;
  }
  .md-content :deep(code) {
    @apply border border-zinc-700/90 bg-zinc-900/95 px-1.5 py-px font-mono text-[0.9em] text-zinc-300;
  }
  .md-content :deep(pre) {
    @apply mb-0 overflow-x-auto border-0 bg-transparent p-3;
  }
  .md-content :deep(pre code) {
    @apply border-0 bg-transparent p-0;
  }
  .md-content :deep(.md-code-block) {
    @apply mb-3 overflow-hidden border border-zinc-700/90 bg-zinc-900/95;
  }
  .md-content :deep(.md-code-block-label) {
    @apply border-b border-zinc-700/90 bg-zinc-800/95 px-3 py-1.5 font-mono text-[0.7rem] leading-none tracking-[0.04em] text-zinc-400 lowercase;
  }
  .md-content :deep(img:not(.emoji-inline)) {
    @apply w-auto max-w-[min(100%,60rem)] border border-zinc-700/80;
  }
  .md-content :deep(img.emoji-inline) {
    display: inline-block;
    width: 1.2em;
    min-width: 1.2em;
    max-width: 1.2em;
    height: 1.2em;
    min-height: 1.2em;
    max-height: 1.2em;
    vertical-align: -0.2em;
    border: 0;
  }
  .md-table-wrap {
    @apply my-4 overflow-hidden border border-zinc-800 bg-zinc-950/70;
  }
  .md-content :deep(.md-math-block) {
    @apply my-1 inline-flex;
  }
  .md-content :deep(.md-math) {
    @apply inline-flex flex-wrap items-baseline gap-[0.04em] text-[1.02em] leading-[1.2] text-zinc-100;
    font-family: Cambria, 'Times New Roman', serif;
  }
  .md-content :deep(.math-id) {
    @apply italic;
  }
  .md-content :deep(.math-op),
  .md-content :deep(.math-num),
  .md-content :deep(.math-symbol),
  .md-content :deep(.math-text) {
    @apply not-italic;
  }
  .md-content :deep(.math-frac) {
    @apply mx-[0.12em] inline-flex flex-col items-stretch align-middle leading-[1.05];
  }
  .md-content :deep(.math-frac-num) {
    @apply block border-b border-zinc-200/90 px-[0.18em] pb-[0.06em] text-center;
  }
  .md-content :deep(.math-frac-den) {
    @apply block px-[0.18em] pt-[0.06em] text-center;
  }
  .md-content :deep(.math-radical) {
    @apply mx-[0.06em] inline-flex items-start;
  }
  .md-content :deep(.math-radical-sign) {
    @apply pr-[0.04em] text-[1.16em] leading-none;
  }
  .md-content :deep(.math-radical-body) {
    @apply border-t border-zinc-200/90 px-[0.14em] pt-[0.08em] pr-[0.14em] pl-[0.1em];
  }
  .md-content :deep(.math-scripted) {
    @apply relative inline-flex items-baseline;
  }
  .md-content :deep(.math-scripted sup) {
    @apply relative text-[0.72em] leading-none;
    top: -0.45em;
  }
  .md-content :deep(.math-scripted sub) {
    @apply relative text-[0.72em] leading-none;
    top: 0.28em;
  }
</style>
