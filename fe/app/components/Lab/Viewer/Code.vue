<script setup lang="ts">
  const CODE_VIEWER_THEME_OPTIONS: Array<{ value: CodeViewerTheme; label: string }> = [
    { value: 'github-dark', label: 'GitHub Dark' },
    { value: 'github-dark-dimmed', label: 'GitHub Dark Dimmed' },
    { value: 'github', label: 'GitHub Light' },
    { value: 'atom-one-dark', label: 'Atom One Dark' },
    { value: 'atom-one-light', label: 'Atom One Light' },
    { value: 'base16/google-dark', label: 'Google Dark' },
    { value: 'base16/google-light', label: 'Google Light' },
    { value: 'base16/gruvbox-dark-hard', label: 'Gruvbox Dark Hard' },
    { value: 'base16/gruvbox-light-soft', label: 'Gruvbox Light Soft' },
    { value: 'default', label: 'Default' },
    { value: 'mono-blue', label: 'Mono Blue' },
    { value: 'monokai', label: 'Monokai' },
    { value: 'night-owl', label: 'Night Owl' },
    { value: 'nnfx-dark', label: 'NNFX Dark' },
    { value: 'stackoverflow-dark', label: 'StackOverflow Dark' },
    { value: 'stackoverflow-light', label: 'StackOverflow Light' },
    { value: 'tokyo-night-dark', label: 'Tokyo Night Dark' },
    { value: 'vs', label: 'VS' },
    { value: 'vs2015', label: 'VS 2015' }
  ]
  const CODE_VIEWER_THEME_IMPORTERS: Record<CodeViewerTheme, () => Promise<{ default: string }>> = {
    'atom-one-dark': () => import('highlight.js/styles/atom-one-dark.css?inline'),
    'atom-one-light': () => import('highlight.js/styles/atom-one-light.css?inline'),
    'base16/google-dark': () => import('highlight.js/styles/base16/google-dark.css?inline'),
    'base16/google-light': () => import('highlight.js/styles/base16/google-light.css?inline'),
    'base16/gruvbox-dark-hard': () => import('highlight.js/styles/base16/gruvbox-dark-hard.css?inline'),
    'base16/gruvbox-light-soft': () => import('highlight.js/styles/base16/gruvbox-light-soft.css?inline'),
    default: () => import('highlight.js/styles/default.css?inline'),
    github: () => import('highlight.js/styles/github.css?inline'),
    'github-dark': () => import('highlight.js/styles/github-dark.css?inline'),
    'github-dark-dimmed': () => import('highlight.js/styles/github-dark-dimmed.css?inline'),
    'mono-blue': () => import('highlight.js/styles/mono-blue.css?inline'),
    monokai: () => import('highlight.js/styles/monokai.css?inline'),
    'night-owl': () => import('highlight.js/styles/night-owl.css?inline'),
    'nnfx-dark': () => import('highlight.js/styles/nnfx-dark.css?inline'),
    'stackoverflow-dark': () => import('highlight.js/styles/stackoverflow-dark.css?inline'),
    'stackoverflow-light': () => import('highlight.js/styles/stackoverflow-light.css?inline'),
    'tokyo-night-dark': () => import('highlight.js/styles/tokyo-night-dark.css?inline'),
    vs: () => import('highlight.js/styles/vs.css?inline'),
    vs2015: () => import('highlight.js/styles/vs2015.css?inline')
  }
  /**
   * Prefixes highlight.js selectors so the chosen theme affects only the code viewer.
   */
  const scopeHighlightTheme = (css: string, scopeSelector: string): string => {
    return String(css || '')
      .replace(/\/\*[\s\S]*?\*\//g, '')
      .replace(/(^|})\s*([^@{}][^{}]*)\{/g, (_, prefix: string, selectorGroup: string) => {
        const scopedSelectors = String(selectorGroup)
          .split(',')
          .map(selector => selector.trim())
          .filter(Boolean)
          .map(selector => `${scopeSelector} ${selector}`)
          .join(',\n')
        return `${prefix}\n${scopedSelectors} {`
      })
      .trim()
  }
  /**
   * Loads the CSS source for one highlight.js theme.
   */
  const loadThemeCss = async (theme: CodeViewerTheme): Promise<string> => {
    const importer = CODE_VIEWER_THEME_IMPORTERS[theme]
    if (!importer) return ''
    const mod = await importer()
    return String(mod.default || '')
  }
  const props = withDefaults(
    defineProps<{
      code: string
      language?: string
    }>(),
    {
      language: 'text'
    }
  )
  const { copyFrom, copied } = useClipboard()
  const { effectiveTheme } = useInterfacePreferences()
  const uiPreferences = useUiPreferencesStore()
  const activeThemeCss = ref('')
  let themeLoadToken = 0
  const highlightedCode = computed(() => highlightCodeSnippet(props.code, props.language))
  const languageLabel = computed(() => highlightedCode.value.displayLanguage)
  const themeOptions = computed<SelectOptionInput[]>(() =>
    CODE_VIEWER_THEME_OPTIONS.map(option => ({
      value: option.value,
      label: option.label
    }))
  )
  const shouldSoftWrap = computed({
    get: () => uiPreferences.codeViewerSoftWrap,
    set: value => {
      uiPreferences.setCodeViewerSoftWrap(Boolean(value))
    }
  })
  const resolvedTheme = computed<CodeViewerTheme>(() => {
    if (uiPreferences.codeViewerThemeOverride) {
      return uiPreferences.codeViewerThemeOverride
    }
    return uiPreferences.defaultCodeViewerThemeFor(effectiveTheme.value)
  })
  const selectedTheme = computed<CodeViewerTheme>({
    get: () => resolvedTheme.value,
    set: value => {
      uiPreferences.setCodeViewerThemeOverride(value as CodeViewerTheme)
    }
  })
  watch(
    selectedTheme,
    async value => {
      const token = ++themeLoadToken
      const css = await loadThemeCss(value)
      if (token !== themeLoadToken) return
      activeThemeCss.value = scopeHighlightTheme(css, '.code-preview')
    },
    { immediate: true }
  )
  useHead(() => ({
    style: activeThemeCss.value
      ? [
          {
            key: 'lab-viewer-code-theme',
            innerHTML: activeThemeCss.value
          }
        ]
      : []
  }))
  /**
   * Copies the current raw code snippet to the clipboard.
   */
  const copyCode = async () => {
    await copyFrom(props.code)
  }
</script>
<template>
  <div class="code-preview">
    <div class="code-preview-head">
      <div class="code-preview-language">
        <Icon name="ic:round-code" class="h-4 w-4 shrink-0" />
        <span>{{ languageLabel }}</span>
      </div>
      <div class="code-preview-controls">
        <LabBaseSelect v-model="selectedTheme" :options="themeOptions" class="code-preview-theme-select" />
        <LabBaseCheckbox v-model="shouldSoftWrap" label="wrap" />
        <LabBaseButton
          :label="copied ? 'скопировано' : 'копировать'"
          :icon="copied ? 'ic:round-check' : 'ic:round-content-copy'"
          variant="ghost"
          size="xs"
          button-class="rounded"
          @click="copyCode" />
      </div>
    </div>
    <div
      class="code-preview-body"
      :class="{ 'code-preview-body-wrap': shouldSoftWrap }"
      v-html="highlightedCode.highlightedHtml"></div>
  </div>
</template>
<style scoped>
  @reference "~/assets/css/main.css";
  .code-preview {
    overflow: hidden;
    border: 1px solid color-mix(in srgb, var(--lab-border) 84%, transparent);
    background: color-mix(in srgb, var(--lab-bg-surface) 92%, transparent);
    color: var(--lab-text-primary);
    border-radius: 1rem;
  }
  .code-preview-head {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
    padding: 0.5rem 0.75rem;
    border-bottom: 1px solid color-mix(in srgb, var(--lab-border) 78%, transparent);
    background: color-mix(in srgb, var(--lab-bg-surface-muted) 88%, transparent);
  }
  .code-preview-language {
    display: flex;
    min-width: 0;
    align-items: center;
    gap: 0.5rem;
    font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
    font-size: 0.75rem;
    line-height: 1;
    text-transform: lowercase;
    letter-spacing: 0.04em;
    color: var(--lab-text-muted);
  }
  .code-preview-controls {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }
  .code-preview-theme-select {
    width: 100%;
    min-width: 0;
  }
  .code-preview-body :deep(pre) {
    overflow-x: auto;
    margin: 0;
    padding: 0;
    background: transparent;
  }
  .code-preview-body :deep(pre code.hljs) {
    display: block;
    padding: 1rem;
    background: transparent;
    font-size: 0.75rem;
    line-height: 1.55;
  }
  .code-preview-body-wrap :deep(pre) {
    overflow-x: visible;
  }
  .code-preview-body-wrap :deep(pre code.hljs) {
    white-space: pre-wrap;
    overflow-wrap: anywhere;
  }
  @media (min-width: 640px) {
    .code-preview-head {
      flex-direction: row;
      align-items: center;
      justify-content: space-between;
    }
    .code-preview-controls {
      flex-direction: row;
      align-items: center;
    }
    .code-preview-theme-select {
      width: 13rem;
    }
  }
</style>
