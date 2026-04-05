<script setup lang="ts">
  defineOptions({
    inheritAttrs: false
  })
  const props = withDefaults(
    defineProps<{
      source?: string | null
      formula?: string | null
      display?: boolean | null
    }>(),
    {
      source: '',
      formula: '',
      display: null
    }
  )
  const normalizeFormulaInput = (raw: string) => {
    const value = String(raw || '').trim()
    if (!value) return { source: '', display: null as boolean | null }
    if (value.startsWith('$$') && value.endsWith('$$') && value.length >= 4) {
      return {
        source: value.slice(2, -2).trim(),
        display: true
      }
    }
    if (value.startsWith('$') && value.endsWith('$') && value.length >= 2) {
      return {
        source: value.slice(1, -1).trim(),
        display: false
      }
    }
    return {
      source: value,
      display: null
    }
  }
  const normalizedFormula = computed(() => {
    const rawSource = String(props.source || '').trim()
    const rawFormula = String(props.formula || '').trim()
    return normalizeFormulaInput(rawSource || rawFormula)
  })
  const isDisplayMode = computed(() => {
    if (typeof props.display === 'boolean') return props.display
    return normalizedFormula.value.display ?? false
  })
  const renderedHtml = computed(() => {
    if (!normalizedFormula.value.source) return ''
    return renderMath(normalizedFormula.value.source, isDisplayMode.value)
  })
  const rootTag = computed(() => {
    return isDisplayMode.value ? 'div' : 'span'
  })
</script>
<template>
  <component
    :is="rootTag"
    v-if="renderedHtml"
    v-bind="$attrs"
    class="latex-viewer"
    :class="isDisplayMode ? 'latex-viewer-block' : 'latex-viewer-inline'">
    <span class="latex-viewer-surface" v-html="renderedHtml"></span>
  </component>
</template>
<style scoped>
  @reference "~/assets/css/main.css";
  .latex-viewer {
    @apply max-w-full text-zinc-100;
  }
  .latex-viewer-inline {
    @apply inline;
  }
  .latex-viewer-block {
    @apply block overflow-x-auto py-2;
  }
  .latex-viewer-surface {
    @apply inline-block max-w-full;
  }
  .latex-viewer-block .latex-viewer-surface {
    min-width: max-content;
  }
  .latex-viewer :deep(.md-math) {
    display: inline-flex;
    flex-wrap: wrap;
    align-items: baseline;
    gap: 0.04em;
    color: rgb(244 244 245);
    font-family: Cambria, 'Times New Roman', serif;
    font-size: 1.02em;
    line-height: 1.2;
  }
  .latex-viewer :deep(.md-math-block) {
    display: inline-flex;
    margin: 0;
  }
  .latex-viewer :deep(.math-id) {
    font-style: italic;
  }
  .latex-viewer :deep(.math-op),
  .latex-viewer :deep(.math-num),
  .latex-viewer :deep(.math-symbol),
  .latex-viewer :deep(.math-text) {
    font-style: normal;
  }
  .latex-viewer :deep(.math-frac) {
    display: inline-flex;
    flex-direction: column;
    align-items: stretch;
    margin-inline: 0.12em;
    vertical-align: middle;
    line-height: 1.05;
  }
  .latex-viewer :deep(.math-frac-num) {
    display: block;
    padding-inline: 0.18em;
    padding-bottom: 0.06em;
    border-bottom: 1px solid rgb(228 228 231 / 0.9);
    text-align: center;
  }
  .latex-viewer :deep(.math-frac-den) {
    display: block;
    padding-top: 0.06em;
    padding-inline: 0.18em;
    text-align: center;
  }
  .latex-viewer :deep(.math-radical) {
    display: inline-flex;
    align-items: flex-start;
    margin-inline: 0.06em;
  }
  .latex-viewer :deep(.math-radical-sign) {
    padding-right: 0.04em;
    font-size: 1.16em;
    line-height: 1;
  }
  .latex-viewer :deep(.math-radical-body) {
    padding-top: 0.08em;
    padding-right: 0.14em;
    padding-left: 0.1em;
    border-top: 1px solid rgb(228 228 231 / 0.9);
  }
  .latex-viewer :deep(.math-scripted) {
    display: inline-flex;
    position: relative;
    align-items: baseline;
  }
  .latex-viewer :deep(.math-scripted sup) {
    position: relative;
    top: -0.45em;
    font-size: 0.72em;
    line-height: 1;
  }
  .latex-viewer :deep(.math-scripted sub) {
    position: relative;
    top: 0.28em;
    font-size: 0.72em;
    line-height: 1;
  }
</style>
