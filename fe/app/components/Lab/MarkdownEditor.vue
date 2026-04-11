<script setup lang="ts">
  const { ensureMarkdownCodec } = useMarkdownCodecLoader()
  const props = withDefaults(
    defineProps<{
      modelValue: string
      rows?: number
    }>(),
    {
      rows: 12
    }
  )
  const emit = defineEmits<{
    (e: 'update:modelValue', value: string): void
  }>()
  const textareaRef = ref<HTMLTextAreaElement | null>(null)
  const previewEditorRef = ref<HTMLDivElement | null>(null)
  const toolbarRef = ref<HTMLElement | null>(null)
  const colorToolRef = ref<HTMLElement | null>(null)
  const linkToolRef = ref<HTMLElement | null>(null)
  const imageToolRef = ref<HTMLElement | null>(null)
  const menu = ref<ToolMenu>('none')
  const hideMenuTimer = ref<ReturnType<typeof setTimeout> | null>(null)
  const emojiScrollRaf = ref<number | null>(null)
  const linkUrlInput = ref('https://')
  const linkError = ref<string | null>(null)
  const linkSelectionState = ref<SelectionState | null>(null)
  const linkSelectionRange = ref<SavedRange>(null)
  const colorValueInput = ref('#22c55e')
  const colorHexInput = ref('#22c55e')
  const colorError = ref<string | null>(null)
  const colorSelectionState = ref<SelectionState | null>(null)
  const colorSelectionRange = ref<SavedRange>(null)
  const imageUrlInput = ref('https://')
  const imageAltInput = ref('image')
  const imageWidthInput = ref('')
  const imageAlignInput = ref<ImageAlignMode>('')
  const imageWrapInput = ref(false)
  const imageError = ref<string | null>(null)
  const imageSelectionState = ref<SelectionState | null>(null)
  const imageSelectionRange = ref<SavedRange>(null)
  const savedPreviewRange = ref<SavedRange>(null)
  const skipNextPreviewModelSync = ref(false)
  const languageFilter = ref('')
  const languageBatch = ref(20)
  const emojiFilter = ref('')
  const emojiBatch = ref(120)
  const emojiNames = ref<string[]>([])
  const emojiLoading = ref(false)
  const emojiLoadError = ref<string | null>(null)
  const previewMode = ref<PreviewMode>('preview')
  const min_image_width_px = MIN_IMAGE_WIDTH_PX
  const emojiDataUriResolver = ref<(name: string) => string>(() => '')
  const emoji_data_uri = (name: string) => emojiDataUriResolver.value(name)
  const EMOJI_CACHE_KEY = 'lab-emojione-face-names-v2'
  const EMOJI_BATCH_STEP = 180
  const MENU_HIDE_DELAY_MS = 120
  const EMOJI_TOKEN_RE = /:emojione:[a-z0-9-]+:/i
  const HEX_COLOR_RE = /^#[0-9a-f]{6}$/i
  const SHORT_HEX_COLOR_RE = /^#[0-9a-f]{3}$/i
  const PREVIEW_MODE_ORDER: PreviewMode[] = ['preview', 'edit']
  const PREVIEW_MODE_META: Record<PreviewMode, { label: string; shortLabel: string; icon: string }> = {
    edit: {
      label: 'Только код',
      shortLabel: 'C',
      icon: 'ic:round-code'
    },
    preview: {
      label: 'WYSIWYG',
      shortLabel: 'W',
      icon: 'ic:round-remove-red-eye'
    }
  }
  const STYLE_OPTIONS: StyleOption[] = [
    {
      label: 'Жирный',
      icon: 'ic:round-format-bold',
      command: 'bold',
      before: '**',
      after: '**',
      placeholder: 'жирный'
    },
    {
      label: 'Курсив',
      icon: 'ic:round-format-italic',
      command: 'italic',
      before: '*',
      after: '*',
      placeholder: 'курсив'
    },
    {
      label: 'Подчеркнутый',
      icon: 'ic:round-format-underlined',
      command: 'underline',
      before: '<u>',
      after: '</u>',
      placeholder: 'подчеркнутый'
    },
    {
      label: 'Перечеркнутый',
      icon: 'ic:round-format-strikethrough',
      command: 'strikeThrough',
      before: '~~',
      after: '~~',
      placeholder: 'перечеркнутый'
    }
  ]
  const HEADING_OPTIONS: HeadingOption[] = [
    { level: 1, size: '1.4rem' },
    { level: 2, size: '1.25rem' },
    { level: 3, size: '1.12rem' },
    { level: 4, size: '1rem' },
    { level: 5, size: '0.94rem' },
    { level: 6, size: '0.88rem' }
  ]
  const ALIGN_OPTIONS: Array<{ mode: 'left' | 'center' | 'right' | 'justify'; label: string; icon: string }> = [
    { mode: 'left', label: 'Влево', icon: 'ic:round-format-align-left' },
    { mode: 'center', label: 'По центру', icon: 'ic:round-format-align-center' },
    { mode: 'right', label: 'Вправо', icon: 'ic:round-format-align-right' },
    { mode: 'justify', label: 'По ширине', icon: 'ic:round-format-align-justify' }
  ]
  const ALIGN_CLASS_LIST = ['md-align-left', 'md-align-center', 'md-align-right', 'md-align-justify'] as const
  const ALIGNABLE_BLOCK_TAGS = new Set(['p', 'h1', 'h2', 'h3', 'h4', 'h5', 'h6', 'li', 'blockquote', 'pre'])
  const ALIGNABLE_BLOCK_SELECTOR = [
    'p',
    'h1',
    'h2',
    'h3',
    'h4',
    'h5',
    'h6',
    'li',
    'blockquote',
    'pre',
    '[align]',
    '[style*="text-align"]',
    '.md-align',
    '.md-align-left',
    '.md-align-center',
    '.md-align-right',
    '.md-align-justify'
  ].join(', ')
  const COLOR_PRESETS = ['#22c55e', '#38bdf8', '#f59e0b', '#f43f5e', '#a78bfa', '#f4f4f5'] as const
  const LANGUAGE_OPTIONS: LanguageOption[] = [
    { key: 'javascript', label: 'JavaScript', category: 'Web / Frontend' },
    { key: 'typescript', label: 'TypeScript', category: 'Web / Frontend' },
    { key: 'html', label: 'HTML', category: 'Web / Frontend' },
    { key: 'css', label: 'CSS', category: 'Web / Frontend' },
    { key: 'scss', label: 'SCSS', category: 'Web / Frontend' },
    { key: 'vue', label: 'Vue', category: 'Web / Frontend' },
    { key: 'jsx', label: 'JSX', category: 'Web / Frontend' },
    { key: 'tsx', label: 'TSX', category: 'Web / Frontend' },
    { key: 'go', label: 'Go', category: 'Web / Backend' },
    { key: 'python', label: 'Python', category: 'Web / Backend' },
    { key: 'php', label: 'PHP', category: 'Web / Backend' },
    { key: 'ruby', label: 'Ruby', category: 'Web / Backend' },
    { key: 'java', label: 'Java', category: 'Web / Backend' },
    { key: 'csharp', label: 'C#', category: 'Web / Backend' },
    { key: 'rust', label: 'Rust', category: 'Web / Backend' },
    { key: 'sql', label: 'SQL', category: 'Web / Backend' },
    { key: 'c', label: 'C', category: 'Низкоуровневые' },
    { key: 'cpp', label: 'C++', category: 'Низкоуровневые' },
    { key: 'asm', label: 'Assembly', category: 'Низкоуровневые' },
    { key: 'zig', label: 'Zig', category: 'Низкоуровневые' },
    { key: 'rust', label: 'Rust', category: 'Низкоуровневые' },
    { key: 'python', label: 'Python', category: 'Высокоуровневые' },
    { key: 'javascript', label: 'JavaScript', category: 'Высокоуровневые' },
    { key: 'typescript', label: 'TypeScript', category: 'Высокоуровневые' },
    { key: 'kotlin', label: 'Kotlin', category: 'Высокоуровневые' },
    { key: 'swift', label: 'Swift', category: 'Высокоуровневые' },
    { key: 'dart', label: 'Dart', category: 'Высокоуровневые' },
    { key: 'bash', label: 'Bash', category: 'Скрипты / DevOps' },
    { key: 'powershell', label: 'PowerShell', category: 'Скрипты / DevOps' },
    { key: 'dockerfile', label: 'Dockerfile', category: 'Скрипты / DevOps' },
    { key: 'yaml', label: 'YAML', category: 'Скрипты / DevOps' },
    { key: 'json', label: 'JSON', category: 'Скрипты / DevOps' },
    { key: 'toml', label: 'TOML', category: 'Скрипты / DevOps' },
    { key: 'ini', label: 'INI', category: 'Скрипты / DevOps' },
    { key: 'markdown', label: 'Markdown', category: 'Скрипты / DevOps' },
    { key: 'text', label: 'Text', category: 'Скрипты / DevOps' }
  ]
  const EMOJI_FALLBACK = Array.from(new Set([...UI_EMOJI_NAMES, ...EMOJI_PICKER_FALLBACK]))
  const currentPreviewMeta = computed(() => PREVIEW_MODE_META[previewMode.value])
  const isStickyInputMenu = computed(() => menu.value === 'link' || menu.value === 'image' || menu.value === 'color')
  const filteredLanguages = computed(() => {
    const query = languageFilter.value.trim().toLowerCase()
    if (!query) return LANGUAGE_OPTIONS
    return LANGUAGE_OPTIONS.filter(item => {
      return (
        item.label.toLowerCase().includes(query) ||
        item.key.toLowerCase().includes(query) ||
        item.category.toLowerCase().includes(query)
      )
    })
  })
  const visibleLanguages = computed(() => {
    return filteredLanguages.value.slice(0, languageBatch.value)
  })
  const filteredEmojiNames = computed(() => {
    const query = emojiFilter.value.trim().toLowerCase()
    if (!query) return emojiNames.value
    return emojiNames.value.filter(name => name.includes(query))
  })
  const visibleEmojiNames = computed(() => {
    return filteredEmojiNames.value.slice(0, emojiBatch.value)
  })
  const isPreviewEditing = () => previewMode.value === 'preview'
  const cloneRange = (range: SavedRange): SavedRange => {
    if (!range) return null
    try {
      return range.cloneRange()
    } catch {
      return null
    }
  }
  const activeTextarea = () => {
    return textareaRef.value
  }
  const readEmojiNamesCache = () => {
    if (!import.meta.client) return []
    try {
      const raw = window.localStorage.getItem(EMOJI_CACHE_KEY)
      if (!raw) return []
      const parsed = JSON.parse(raw)
      if (!Array.isArray(parsed)) return []
      return parsed.filter(item => typeof item === 'string')
    } catch {
      return []
    }
  }
  const ensureEmojiCollectionForModel = async (source: string) => {
    if (!EMOJI_TOKEN_RE.test(String(source || ''))) return
    const markdownCodec = await ensureMarkdownCodec()
    emojiDataUriResolver.value = markdownCodec.getEmojiDataUri
    await markdownCodec.ensureEmojiCollectionLoaded()
    await syncPreviewFromModel(true)
  }
  const writeEmojiNamesCache = (names: string[]) => {
    if (!import.meta.client || names.length === 0) return
    try {
      window.localStorage.setItem(EMOJI_CACHE_KEY, JSON.stringify(names))
    } catch {
      // ignore storage errors
    }
  }
  watch(languageFilter, () => {
    languageBatch.value = 20
  })
  watch(emojiFilter, () => {
    emojiBatch.value = 120
  })
  watch(colorValueInput, next => {
    syncColorHexFromTextInput(next)
  })
  watch(menu, async next => {
    if (next !== 'emoji' || emojiNames.value.length > 0 || emojiLoading.value) return
    const cached = readEmojiNamesCache()
    if (cached.length > 0) {
      emojiNames.value = cached
      return
    }
    emojiLoading.value = true
    emojiLoadError.value = null
    try {
      const markdownCodec = await ensureMarkdownCodec()
      emojiDataUriResolver.value = markdownCodec.getEmojiDataUri
      await markdownCodec.ensureEmojiCollectionLoaded()
      const names = await loadEmojiPickerNames()
      emojiNames.value = names.length > 0 ? names.sort() : EMOJI_FALLBACK
      writeEmojiNamesCache(emojiNames.value)
    } catch {
      emojiNames.value = EMOJI_FALLBACK
      emojiLoadError.value = 'Не удалось загрузить face-набор emoji, показан базовый список.'
    } finally {
      emojiLoading.value = false
    }
  })
  watch(
    () => props.modelValue,
    () => {
      void ensureEmojiCollectionForModel(props.modelValue)
      if (!isPreviewEditing()) return
      nextTick(() => {
        if (skipNextPreviewModelSync.value) {
          skipNextPreviewModelSync.value = false
          return
        }
        void syncPreviewFromModel(false)
      })
    }
  )
  watch(previewMode, nextMode => {
    if (nextMode !== 'preview') return
    nextTick(() => {
      void syncPreviewFromModel(true)
    })
  })
  const clearHideMenuTimer = () => {
    if (!hideMenuTimer.value) return
    clearTimeout(hideMenuTimer.value)
    hideMenuTimer.value = null
  }
  const closeMenu = () => {
    clearHideMenuTimer()
    menu.value = 'none'
  }
  const toggleMenu = (target: Exclude<typeof menu.value, 'none'>) => {
    clearHideMenuTimer()
    menu.value = menu.value === target ? 'none' : target
  }
  const openMenuOnHover = (target: Exclude<typeof menu.value, 'none'>) => {
    if (isStickyInputMenu.value && target !== 'link' && target !== 'image') return
    clearHideMenuTimer()
    menu.value = target
  }
  const onGroupLeave = (target: Exclude<typeof menu.value, 'none'>) => {
    if (target === 'link' || target === 'image' || target === 'color') return
    if (menu.value !== target) return
    clearHideMenuTimer()
    hideMenuTimer.value = setTimeout(() => {
      if (menu.value === target) {
        menu.value = 'none'
      }
      hideMenuTimer.value = null
    }, MENU_HIDE_DELAY_MS)
  }
  const onDocumentClick = (event: MouseEvent) => {
    const node = event.target as Node | null
    if (!node) return
    if (menu.value === 'link') {
      if (linkToolRef.value?.contains(node)) return
      closeMenu()
      return
    }
    if (menu.value === 'color') {
      if (colorToolRef.value?.contains(node)) return
      closeMenu()
      return
    }
    if (menu.value === 'image') {
      if (imageToolRef.value?.contains(node)) return
      closeMenu()
      return
    }
    if (!toolbarRef.value) return
    if (toolbarRef.value.contains(node)) return
    closeMenu()
  }
  const onToolbarMouseDown = (event: MouseEvent) => {
    const target = event.target as HTMLElement | null
    if (!target) return
    if (target.closest('input, textarea')) return
    if (target.closest('button')) {
      event.preventDefault()
    }
  }
  const onInput = (event: Event) => {
    emit('update:modelValue', (event.target as HTMLTextAreaElement).value)
  }
  const cyclePreviewMode = () => {
    const idx = PREVIEW_MODE_ORDER.indexOf(previewMode.value)
    const nextIndex = idx < 0 ? 0 : (idx + 1) % PREVIEW_MODE_ORDER.length
    previewMode.value = PREVIEW_MODE_ORDER[nextIndex] || 'edit'
  }
  const placePreviewCaretAtEnd = () => {
    if (!import.meta.client) return
    const editor = previewEditorRef.value
    if (!editor) return
    const selection = window.getSelection()
    if (!selection) return
    const range = document.createRange()
    range.selectNodeContents(editor)
    range.collapse(false)
    selection.removeAllRanges()
    selection.addRange(range)
    savedPreviewRange.value = cloneRange(range)
  }
  const capturePreviewRange = () => {
    if (!import.meta.client || !isPreviewEditing()) return
    const editor = previewEditorRef.value
    if (!editor) return
    const selection = window.getSelection()
    if (!selection || selection.rangeCount === 0) return
    const range = selection.getRangeAt(0)
    if (!editor.contains(range.commonAncestorContainer)) return
    savedPreviewRange.value = cloneRange(range)
  }
  const isPreviewSelectionInsideEditor = () => {
    if (!import.meta.client || !isPreviewEditing()) return false
    const editor = previewEditorRef.value
    if (!editor) return false
    const selection = window.getSelection()
    if (!selection || selection.rangeCount === 0) {
      return document.activeElement === editor
    }
    const anchor = selection.anchorNode
    const focus = selection.focusNode
    if (anchor && editor.contains(anchor)) return true
    if (focus && editor.contains(focus)) return true
    try {
      return editor.contains(selection.getRangeAt(0).commonAncestorContainer)
    } catch {
      return document.activeElement === editor
    }
  }
  const restorePreviewRange = (preferredRange?: SavedRange) => {
    if (!import.meta.client) return false
    const editor = previewEditorRef.value
    if (!editor) return false
    const selection = window.getSelection()
    if (!selection) return false
    const candidate = cloneRange(preferredRange || savedPreviewRange.value)
    if (candidate) {
      try {
        if (editor.contains(candidate.commonAncestorContainer)) {
          selection.removeAllRanges()
          selection.addRange(candidate)
          savedPreviewRange.value = cloneRange(candidate)
          return true
        }
      } catch {
        // ignored
      }
    }
    placePreviewCaretAtEnd()
    return true
  }
  const syncPreviewFromModel = async (force = false) => {
    if (!import.meta.client) return
    const editor = previewEditorRef.value
    if (!editor) return
    if (isPreviewSelectionInsideEditor()) return
    if (!force && document.activeElement === editor) return
    const markdownCodec = await ensureMarkdownCodec()
    emojiDataUriResolver.value = markdownCodec.getEmojiDataUri
    const nextHtml = markdownCodec.renderMarkdownToHtml(props.modelValue)
    if (editor.innerHTML !== nextHtml) {
      editor.innerHTML = nextHtml
    }
  }
  const syncModelFromPreview = async () => {
    if (!import.meta.client || !isPreviewEditing()) return
    const editor = previewEditorRef.value
    if (!editor) return
    const markdownCodec = await ensureMarkdownCodec()
    emojiDataUriResolver.value = markdownCodec.getEmojiDataUri
    const nextMarkdown = markdownCodec.renderEditableHtmlToMarkdown(editor.innerHTML)
    if (nextMarkdown !== props.modelValue) {
      if (isPreviewSelectionInsideEditor()) {
        skipNextPreviewModelSync.value = true
      }
      emit('update:modelValue', nextMarkdown)
    }
  }
  const onPreviewInput = () => {
    capturePreviewRange()
    void syncModelFromPreview()
  }
  const onPreviewKeydown = (event: KeyboardEvent) => {
    if (!import.meta.client || !isPreviewEditing()) return
    if (event.key !== 'Enter' || event.shiftKey || event.ctrlKey || event.metaKey || event.altKey) return
    const editor = previewEditorRef.value
    if (!editor) return
    const selection = window.getSelection()
    if (!selection || selection.rangeCount === 0) return
    const activeRange = selection.getRangeAt(0)
    if (!editor.contains(activeRange.commonAncestorContainer)) return
    const getHeadingHost = (node: Node | null) => {
      let current: HTMLElement | null =
        node?.nodeType === Node.ELEMENT_NODE ? (node as HTMLElement) : node?.parentElement || null
      while (current && current !== editor) {
        const tag = current.tagName.toLowerCase()
        if (/^h[1-6]$/.test(tag)) return current
        current = current.parentElement
      }
      return null
    }
    const heading = getHeadingHost(activeRange.startContainer)
    if (!heading) return
    event.preventDefault()
    const paragraph = document.createElement('p')
    paragraph.appendChild(document.createElement('br'))
    heading.insertAdjacentElement('afterend', paragraph)
    const nextRange = document.createRange()
    nextRange.setStart(paragraph, 0)
    nextRange.collapse(true)
    editor.focus()
    selection.removeAllRanges()
    selection.addRange(nextRange)
    savedPreviewRange.value = cloneRange(nextRange)
    void syncModelFromPreview()
  }
  const onPreviewFocus = () => {
    capturePreviewRange()
  }
  const onPreviewBlur = () => {
    capturePreviewRange()
  }
  const runPreviewCommand = (command: string, value = '', preferredRange?: SavedRange) => {
    if (!import.meta.client || !isPreviewEditing()) return false
    const editor = previewEditorRef.value
    if (!editor) return false
    editor.focus()
    restorePreviewRange(preferredRange)
    try {
      document.execCommand('styleWithCSS', false, 'false')
      document.execCommand(command, false, value)
    } catch {
      return false
    }
    capturePreviewRange()
    void syncModelFromPreview()
    return true
  }
  const insertPreviewTextFallback = (text: string) => {
    if (!import.meta.client) return false
    const selection = window.getSelection()
    if (!selection || selection.rangeCount === 0) return false
    const range = selection.getRangeAt(0)
    range.deleteContents()
    const textNode = document.createTextNode(text)
    range.insertNode(textNode)
    range.setStartAfter(textNode)
    range.collapse(true)
    selection.removeAllRanges()
    selection.addRange(range)
    return true
  }
  const insertPreviewText = (text: string, preferredRange?: SavedRange) => {
    if (!import.meta.client || !isPreviewEditing()) return false
    const editor = previewEditorRef.value
    if (!editor) return false
    editor.focus()
    restorePreviewRange(preferredRange)
    let inserted = false
    try {
      inserted = document.execCommand('insertText', false, text)
    } catch {
      inserted = false
    }
    if (!inserted) {
      inserted = insertPreviewTextFallback(text)
    }
    capturePreviewRange()
    void syncModelFromPreview()
    return inserted
  }
  const insertPreviewHtmlFallback = (html: string) => {
    if (!import.meta.client) return false
    const selection = window.getSelection()
    if (!selection || selection.rangeCount === 0) return false
    const range = selection.getRangeAt(0)
    range.deleteContents()
    const holder = document.createElement('div')
    holder.innerHTML = html
    const fragment = document.createDocumentFragment()
    let lastNode: ChildNode | null = null
    while (holder.firstChild) {
      lastNode = holder.firstChild
      fragment.appendChild(holder.firstChild)
    }
    range.insertNode(fragment)
    if (lastNode) {
      range.setStartAfter(lastNode)
      range.collapse(true)
      selection.removeAllRanges()
      selection.addRange(range)
    }
    return true
  }
  const insertPreviewHtml = (html: string, preferredRange?: SavedRange, preferRangeInsertion = false) => {
    if (!import.meta.client || !isPreviewEditing()) return false
    const editor = previewEditorRef.value
    if (!editor) return false
    editor.focus()
    restorePreviewRange(preferredRange)
    let inserted = false
    const before = editor.innerHTML
    if (!preferRangeInsertion) {
      try {
        inserted = document.execCommand('insertHTML', false, html)
      } catch {
        inserted = false
      }
      if (inserted && editor.innerHTML === before) {
        inserted = false
      }
    }
    if (!inserted) {
      inserted = insertPreviewHtmlFallback(html)
    }
    capturePreviewRange()
    void syncModelFromPreview()
    return inserted
  }
  const readSelection = (): SelectionState => {
    const value = String(props.modelValue || '')
    const area = activeTextarea()
    if (!area) {
      return {
        value,
        start: value.length,
        end: value.length,
        selected: ''
      }
    }
    const start = area.selectionStart ?? 0
    const end = area.selectionEnd ?? start
    return {
      value,
      start,
      end,
      selected: value.slice(start, end)
    }
  }
  const parseMarkdownImageMeta = (rawMeta?: string) => {
    const normalized = String(rawMeta || '').trim()
    const body = normalized.startsWith('{') && normalized.endsWith('}') ? normalized.slice(1, -1).trim() : ''
    const out: { width: string; align: ImageAlignMode; wrap: boolean } = {
      width: '',
      align: '',
      wrap: false
    }
    if (!body) return out
    const parts = body.split(/\s+/).filter(Boolean)
    for (const part of parts) {
      const [keyRaw, ...valueParts] = part.split('=')
      const key = String(keyRaw || '')
        .trim()
        .toLowerCase()
      const value = String(valueParts.join('=') || '').trim()
      if (!key || !value) continue
      if (key === 'width') {
        out.width = normalizeImageDisplayWidth(value)
        continue
      }
      if (key === 'align') {
        out.align = normalizeImageDisplayAlign(value)
        continue
      }
      if (key === 'wrap') {
        out.wrap = normalizeImageTextWrap(value)
      }
    }
    if (out.wrap && out.align !== 'left' && out.align !== 'right') {
      out.wrap = false
    }
    return out
  }
  const buildMarkdownImageMeta = (width: string, align: ImageAlignMode, wrap: boolean) => {
    const normalizedWidth = normalizeImageDisplayWidth(width)
    const normalizedAlign = normalizeImageDisplayAlign(align)
    const normalizedWrap = normalizeImageTextWrap(wrap) && (normalizedAlign === 'left' || normalizedAlign === 'right')
    const parts: string[] = []
    if (normalizedWidth) parts.push(`width=${normalizedWidth}`)
    if (normalizedAlign) parts.push(`align=${normalizedAlign}`)
    if (normalizedWrap) parts.push('wrap=1')
    return parts.length > 0 ? `{${parts.join(' ')}}` : ''
  }
  const buildMarkdownImageToken = (alt: string, url: string, width: string, align: ImageAlignMode, wrap: boolean) => {
    const meta = buildMarkdownImageMeta(width, align, wrap)
    return `![${alt || 'image'}](${url})${meta}`
  }
  const findMarkdownTokenAtSelection = (state: SelectionState, kind: UrlTargetKind): MarkdownInlineToken | null => {
    const value = String(state.value || '')
    if (!value) return null
    const left = Math.max(0, Math.min(state.start, state.end))
    const right = Math.max(state.start, state.end)
    const pattern = /(!?)\[([^\]]*)\]\(([^)\s]+)\)(\{[^}]+\})?/g
    let match: RegExpExecArray | null
    while ((match = pattern.exec(value)) !== null) {
      const isImage = match[1] === '!'
      if ((kind === 'image') !== isImage) continue
      const tokenStart = match.index
      const tokenEnd = tokenStart + match[0].length
      const touchesToken =
        left === right ? left >= tokenStart && left <= tokenEnd : right > tokenStart && left < tokenEnd
      if (!touchesToken) continue
      const imageMeta = isImage
        ? parseMarkdownImageMeta(match[4] || '')
        : { width: '', align: '' as ImageAlignMode, wrap: false }
      return {
        start: tokenStart,
        end: tokenEnd,
        isImage,
        label: String(match[2] || ''),
        url: String(match[3] || ''),
        width: imageMeta.width,
        align: imageMeta.align,
        wrap: imageMeta.wrap
      }
    }
    return null
  }
  const findClosestEditorTag = (node: Node | null, editor: HTMLElement, tagName: 'a' | 'img'): HTMLElement | null => {
    let current: HTMLElement | null = null
    if (node?.nodeType === Node.ELEMENT_NODE) {
      current = node as HTMLElement
    } else {
      current = node?.parentElement || null
    }
    while (current && current !== editor) {
      if (current.tagName.toLowerCase() === tagName) return current
      current = current.parentElement
    }
    return null
  }
  const findClosestColorSpan = (node: Node | null, editor: HTMLElement) => {
    let current: HTMLElement | null = null
    if (node?.nodeType === Node.ELEMENT_NODE) {
      current = node as HTMLElement
    } else {
      current = node?.parentElement || null
    }
    while (current && current !== editor) {
      if (
        current.tagName.toLowerCase() === 'span' &&
        /(^|;\s*)color\s*:/.test(String(current.getAttribute('style') || '').toLowerCase())
      ) {
        return current
      }
      current = current.parentElement
    }
    return null
  }
  const elementFromPreviewNode = (node: Node | null, editor: HTMLElement, kind: UrlTargetKind): HTMLElement | null => {
    if (kind === 'link') {
      return findClosestEditorTag(node, editor, 'a')
    }
    const image = findClosestEditorTag(node, editor, 'img')
    if (image) return image
    const anchor = findClosestEditorTag(node, editor, 'a')
    const nestedImage = anchor?.querySelector('img')
    return nestedImage && editor.contains(nestedImage) ? (nestedImage as HTMLElement) : null
  }
  const urlFromPreviewNode = (node: Node | null, editor: HTMLElement, kind: UrlTargetKind) => {
    const element = elementFromPreviewNode(node, editor, kind)
    if (!element) return ''
    if (kind === 'link') return String(element.getAttribute('href') || '').trim()
    if (kind === 'image') return String(element.getAttribute('src') || '').trim()
    return ''
  }
  const findPreviewElementAtRange = (range: SavedRange, kind: UrlTargetKind) => {
    if (!import.meta.client) return null
    const editor = previewEditorRef.value
    if (!editor || !range) return null
    const byStart = elementFromPreviewNode(range.startContainer, editor, kind)
    if (byStart) return byStart
    const byEnd = elementFromPreviewNode(range.endContainer, editor, kind)
    if (byEnd) return byEnd
    const edgeTargets: Array<{ container: Node; offset: number }> = [
      { container: range.startContainer, offset: range.startOffset },
      { container: range.endContainer, offset: range.endOffset }
    ]
    for (const edge of edgeTargets) {
      if (edge.container.nodeType !== Node.ELEMENT_NODE) continue
      const host = edge.container as HTMLElement
      const near = [host.childNodes[edge.offset], host.childNodes[edge.offset - 1]]
      for (const candidate of near) {
        if (!candidate || !editor.contains(candidate)) continue
        const found = elementFromPreviewNode(candidate, editor, kind)
        if (found) return found
      }
    }
    return null
  }
  const findPreviewColorSpanAtRange = (range: SavedRange) => {
    if (!import.meta.client) return null
    const editor = previewEditorRef.value
    if (!editor || !range) return null
    const byStart = findClosestColorSpan(range.startContainer, editor)
    if (byStart) return byStart
    const byEnd = findClosestColorSpan(range.endContainer, editor)
    if (byEnd) return byEnd
    return null
  }
  const readPreviewUrlAtRange = (range: SavedRange, kind: UrlTargetKind) => {
    if (!import.meta.client) return ''
    const editor = previewEditorRef.value
    if (!editor || !range) return ''
    const byStart = urlFromPreviewNode(range.startContainer, editor, kind)
    if (byStart) return byStart
    const byEnd = urlFromPreviewNode(range.endContainer, editor, kind)
    if (byEnd) return byEnd
    const edgeTargets: Array<{ container: Node; offset: number }> = [
      { container: range.startContainer, offset: range.startOffset },
      { container: range.endContainer, offset: range.endOffset }
    ]
    for (const edge of edgeTargets) {
      if (edge.container.nodeType !== Node.ELEMENT_NODE) continue
      const host = edge.container as HTMLElement
      const near = [host.childNodes[edge.offset], host.childNodes[edge.offset - 1]]
      for (const candidate of near) {
        if (!candidate || !editor.contains(candidate)) continue
        const found = urlFromPreviewNode(candidate, editor, kind)
        if (found) return found
      }
    }
    return ''
  }
  const readMarkdownUrlAtSelection = (state: SelectionState, kind: UrlTargetKind) => {
    const token = findMarkdownTokenAtSelection(state, kind)
    return String(token?.url || '').trim()
  }
  const replaceMarkdownImageTokenByState = (
    state: SelectionState,
    url: string,
    alt: string,
    width: string,
    align: ImageAlignMode,
    wrap: boolean
  ) => {
    const token = findMarkdownTokenAtSelection(state, 'image')
    if (!token) return false
    const replacement = buildMarkdownImageToken(alt, url, width, align, wrap)
    const next = `${state.value.slice(0, token.start)}${replacement}${state.value.slice(token.end)}`
    const caretPos = token.start + replacement.length
    applyNextValue(next, caretPos)
    return true
  }
  const findMarkdownColorSpanAtSelection = (state: SelectionState) => {
    const value = String(state.value || '')
    if (!value) return null
    const left = Math.max(0, Math.min(state.start, state.end))
    const right = Math.max(state.start, state.end)
    const pattern = /<span\s+style="([^"]*color\s*:[^"]+)"\s*>([\s\S]*?)<\/span>/gi
    let match: RegExpExecArray | null
    while ((match = pattern.exec(value)) !== null) {
      const tokenStart = match.index
      const tokenEnd = tokenStart + match[0].length
      const touchesToken =
        left === right ? left >= tokenStart && left <= tokenEnd : right > tokenStart && left < tokenEnd
      if (!touchesToken) continue
      return {
        start: tokenStart,
        end: tokenEnd,
        text: String(match[2] || ''),
        color: normalizeCssColor(
          String(match[1] || '')
            .split(';')
            .find(part => part.toLowerCase().includes('color'))
            ?.split(':')
            .slice(1)
            .join(':') || ''
        )
      }
    }
    return null
  }
  const replaceMarkdownTokenUrlByState = (state: SelectionState, kind: UrlTargetKind, url: string) => {
    const token = findMarkdownTokenAtSelection(state, kind)
    if (!token) return false
    if (token.isImage) {
      const replacement = buildMarkdownImageToken(token.label, url, token.width, token.align, token.wrap)
      const next = `${state.value.slice(0, token.start)}${replacement}${state.value.slice(token.end)}`
      const caretPos = token.start + replacement.length
      applyNextValue(next, caretPos)
      return true
    }
    const prefix = token.isImage ? '!' : ''
    const replacement = `${prefix}[${token.label}](${url})`
    const next = `${state.value.slice(0, token.start)}${replacement}${state.value.slice(token.end)}`
    const caretPos = token.start + replacement.length
    applyNextValue(next, caretPos)
    return true
  }
  const stripInlineFormatting = (raw: string) => {
    let out = String(raw || '')
    let prev = ''
    while (out !== prev) {
      prev = out
      out = out.replace(/!\[([^\]]*)\]\(([^)\s]+)\)(?:\{[^}]+\})?/g, (_, alt: string, src: string) =>
        (alt || src || '').trim()
      )
      out = out.replace(/\[([^\]]+)\]\(([^)\s]+)\)/g, '$1')
      out = out.replace(/<u>([\s\S]*?)<\/u>/gi, '$1')
      out = out.replace(/`([^`]+)`/g, '$1')
      out = out.replace(/\*\*([^*]+)\*\*/g, '$1')
      out = out.replace(/__([^_]+)__/g, '$1')
      out = out.replace(/~~([^~]+)~~/g, '$1')
      out = out.replace(/\*([^*\n]+)\*/g, '$1')
      out = out.replace(/_([^_\n]+)_/g, '$1')
      out = out.replace(/:emojione:([a-z0-9-]+):/gi, '$1')
      out = out.replace(/<\/?u>/gi, '')
    }
    return out.replace(/<[^>]+>/g, '')
  }
  const stripMarkdownFormatting = (raw: string) => {
    const lines = String(raw || '')
      .replace(/\r/g, '')
      .split('\n')
    const cleaned = lines.map(line => {
      const trimmed = line.trim()
      if (/^```/.test(trimmed)) return ''
      if (/^:::\s*align-(left|center|right|justify)\s*$/i.test(trimmed)) return ''
      if (trimmed === ':::') return ''
      if (/^<div\s+align="(left|center|right|justify)"\s*>$/i.test(trimmed)) return ''
      if (/^<\/div>\s*$/i.test(trimmed)) return ''
      let out = line
      out = out.replace(/^\s{0,3}#{1,6}\s+/, '')
      out = out.replace(/^\s*>\s?/, '')
      out = out.replace(/^\s*([-*+]|\d+\.)\s+/, '')
      out = stripInlineFormatting(out)
      return out
    })
    return cleaned.join('\n').replace(/[ \t]+$/gm, '')
  }
  const clearFormatting = () => {
    if (isPreviewEditing()) {
      const selectedRange = cloneRange(savedPreviewRange.value)
      if (selectedRange && !selectedRange.collapsed) {
        const plain = selectedRange.toString()
        insertPreviewText(plain, selectedRange)
        closeMenu()
        return
      }
      const plain = stripMarkdownFormatting(props.modelValue)
      emit('update:modelValue', plain)
      nextTick(() => {
        void syncPreviewFromModel(true)
        placePreviewCaretAtEnd()
      })
      closeMenu()
      return
    }
    const state = readSelection()
    if (state.selected) {
      const plainSelected = stripMarkdownFormatting(state.selected)
      replaceSelectionByState(state, plainSelected, true)
    } else {
      const plainAll = stripMarkdownFormatting(state.value)
      applyNextValue(plainAll, Math.min(state.start, plainAll.length))
    }
    closeMenu()
  }
  const replaceSelectionByState = (state: SelectionState, replacement: string, selectWhole = false) => {
    const next = `${state.value.slice(0, state.start)}${replacement}${state.value.slice(state.end)}`
    if (selectWhole) {
      applyNextValue(next, state.start, state.start + replacement.length)
      return
    }
    applyNextValue(next, state.start + replacement.length)
  }
  const applyNextValue = (value: string, selectionStart: number, selectionEnd = selectionStart) => {
    emit('update:modelValue', value)
    nextTick(() => {
      const area = activeTextarea()
      if (!area) return
      area.focus()
      area.setSelectionRange(selectionStart, selectionEnd)
    })
  }
  const wrapSelection = (before: string, after: string, placeholder = 'текст') => {
    const state = readSelection()
    const picked = state.selected || placeholder
    const replacement = `${before}${picked}${after}`
    const next = `${state.value.slice(0, state.start)}${replacement}${state.value.slice(state.end)}`
    const selStart = state.start + before.length
    const selEnd = selStart + picked.length
    applyNextValue(next, selStart, selEnd)
  }
  const replaceSelection = (replacement: string, selectWhole = false) => {
    const state = readSelection()
    replaceSelectionByState(state, replacement, selectWhole)
  }
  const applyInlineStyle = (option: StyleOption) => {
    if (isPreviewEditing()) {
      runPreviewCommand(option.command)
      closeMenu()
      return
    }
    wrapSelection(option.before, option.after, option.placeholder)
    closeMenu()
  }
  const applyHeading = (level: number) => {
    if (isPreviewEditing()) {
      runPreviewCommand('formatBlock', `H${level}`)
      closeMenu()
      return
    }
    const prefix = `${'#'.repeat(level)} `
    const state = readSelection()
    const text = (state.selected || `Заголовок H${level}`)
      .split('\n')
      .map(line => line.replace(/^#{1,6}\s+/, '').trim() || `Заголовок H${level}`)
      .map(line => `${prefix}${line}`)
      .join('\n')
    replaceSelection(text, true)
    closeMenu()
  }
  const applyList = (ordered: boolean) => {
    if (isPreviewEditing()) {
      runPreviewCommand(ordered ? 'insertOrderedList' : 'insertUnorderedList')
      return
    }
    const marker = ordered ? '1. ' : '- '
    const state = readSelection()
    if (!state.selected) {
      replaceSelection(marker)
      return
    }
    const transformed = state.selected
      .split('\n')
      .map((line, index) => {
        const clean = line.replace(/^\s*([-*+]|\d+\.)\s+/, '').trim() || 'пункт'
        return ordered ? `${index + 1}. ${clean}` : `- ${clean}`
      })
      .join('\n')
    const beforeChar = state.start > 0 ? state.value[state.start - 1] : '\n'
    const afterChar = state.end < state.value.length ? state.value[state.end] : '\n'
    const prefix = beforeChar === '\n' ? '' : '\n'
    const suffix = afterChar === '\n' ? '' : '\n'
    const replacement = `${prefix}${transformed}${suffix}`
    const next = `${state.value.slice(0, state.start)}${replacement}${state.value.slice(state.end)}`
    const selStart = state.start + prefix.length
    const selEnd = selStart + transformed.length
    applyNextValue(next, selStart, selEnd)
  }
  const readElementAlignMode = (el: HTMLElement | null): '' | 'left' | 'center' | 'right' | 'justify' => {
    if (!el) return ''
    const byAttr = String(el.getAttribute('align') || '')
      .toLowerCase()
      .trim()
    if (['left', 'center', 'right', 'justify'].includes(byAttr)) {
      return byAttr as 'left' | 'center' | 'right' | 'justify'
    }
    for (const cls of ALIGN_CLASS_LIST) {
      if (el.classList.contains(cls)) {
        return cls.replace('md-align-', '') as 'left' | 'center' | 'right' | 'justify'
      }
    }
    const byStyle = String(el.style.textAlign || '')
      .toLowerCase()
      .trim()
    if (['left', 'center', 'right', 'justify'].includes(byStyle)) {
      return byStyle as 'left' | 'center' | 'right' | 'justify'
    }
    return ''
  }
  const getClosestAlignElement = (node: Node | null, editor: HTMLElement): HTMLElement | null => {
    let current: Node | null = node
    while (current && current !== editor) {
      if (current.nodeType === Node.ELEMENT_NODE) {
        const el = current as HTMLElement
        if (readElementAlignMode(el)) return el
      }
      current = current.parentNode
    }
    return null
  }
  const isAlignableBlockElement = (el: HTMLElement) => {
    const tag = el.tagName.toLowerCase()
    if (ALIGNABLE_BLOCK_TAGS.has(tag)) return true
    return readElementAlignMode(el) !== ''
  }
  const getClosestAlignableBlock = (node: Node | null, editor: HTMLElement): HTMLElement | null => {
    let current: Node | null = node
    while (current && current !== editor) {
      if (current.nodeType === Node.ELEMENT_NODE) {
        const el = current as HTMLElement
        if (isAlignableBlockElement(el)) return el
      }
      current = current.parentNode
    }
    return null
  }
  const setElementAlignMode = (el: HTMLElement, mode: 'left' | 'center' | 'right' | 'justify') => {
    const hasMdAlign = el.classList.contains('md-align') || ALIGN_CLASS_LIST.some(cls => el.classList.contains(cls))
    el.classList.remove(...ALIGN_CLASS_LIST)
    if (hasMdAlign) {
      el.classList.add('md-align', `md-align-${mode}`)
      el.removeAttribute('align')
      el.style.removeProperty('text-align')
      return
    }
    if (el.hasAttribute('align')) {
      el.setAttribute('align', mode)
      el.style.removeProperty('text-align')
      return
    }
    el.style.textAlign = mode
  }
  const collectPreviewAlignmentTargets = (range: Range, editor: HTMLElement) => {
    const out: HTMLElement[] = []
    const seen = new Set<HTMLElement>()
    const pushTarget = (candidate: HTMLElement | null) => {
      if (!candidate || candidate === editor) return
      if (!editor.contains(candidate)) return
      if (!isAlignableBlockElement(candidate)) return
      if (seen.has(candidate)) return
      seen.add(candidate)
      out.push(candidate)
    }
    pushTarget(getClosestAlignElement(range.startContainer, editor))
    pushTarget(getClosestAlignElement(range.endContainer, editor))
    pushTarget(getClosestAlignableBlock(range.startContainer, editor))
    pushTarget(getClosestAlignableBlock(range.endContainer, editor))
    if (!range.collapsed) {
      const candidates = Array.from(editor.querySelectorAll<HTMLElement>(ALIGNABLE_BLOCK_SELECTOR))
      for (const candidate of candidates) {
        try {
          if (!range.intersectsNode(candidate)) continue
        } catch {
          continue
        }
        pushTarget(candidate)
      }
    }
    const pruned = out.filter(node => {
      return !out.some(other => other !== node && other.contains(node) && readElementAlignMode(other))
    })
    return pruned.length > 0 ? pruned : out
  }
  const flattenNestedPreviewAlignments = () => {
    const editor = previewEditorRef.value
    if (!editor) return
    const alignNodes = Array.from(
      editor.querySelectorAll<HTMLElement>(
        '[align], [style*="text-align"], .md-align-left, .md-align-center, .md-align-right, .md-align-justify, .md-align'
      )
    )
    for (const node of alignNodes) {
      const parent = node.parentElement
      if (!parent || !editor.contains(parent)) continue
      const styleRules = String(node.getAttribute('style') || '')
        .toLowerCase()
        .split(';')
        .map(item => item.trim())
        .filter(Boolean)
      const isPureDivStyleAlign =
        node.tagName.toLowerCase() === 'div' && styleRules.length === 1 && styleRules[0]?.startsWith('text-align:')
      if (
        !node.classList.contains('md-align') &&
        !ALIGN_CLASS_LIST.some(cls => node.classList.contains(cls)) &&
        !node.hasAttribute('align') &&
        !isPureDivStyleAlign
      ) {
        continue
      }
      const ownMode = readElementAlignMode(node)
      const parentMode = readElementAlignMode(parent)
      if (!ownMode || !parentMode || ownMode !== parentMode) continue
      while (node.firstChild) {
        parent.insertBefore(node.firstChild, node)
      }
      parent.removeChild(node)
    }
  }
  const tryUpdateExistingPreviewAlignment = (mode: 'left' | 'center' | 'right' | 'justify') => {
    if (!import.meta.client || !isPreviewEditing()) return false
    const editor = previewEditorRef.value
    if (!editor) return false
    const active = cloneRange(savedPreviewRange.value)
    if (!active) return false
    const startAlign = getClosestAlignElement(active.startContainer, editor)
    const endAlign = getClosestAlignElement(active.endContainer, editor)
    if (!startAlign || startAlign !== endAlign) return false
    setElementAlignMode(startAlign, mode)
    flattenNestedPreviewAlignments()
    capturePreviewRange()
    void syncModelFromPreview()
    return true
  }
  const applyPreviewAlignmentToSelection = (mode: 'left' | 'center' | 'right' | 'justify') => {
    if (!import.meta.client || !isPreviewEditing()) return false
    const editor = previewEditorRef.value
    if (!editor) return false
    capturePreviewRange()
    const active = cloneRange(savedPreviewRange.value)
    if (!active) return false
    const targets = collectPreviewAlignmentTargets(active, editor)
    if (targets.length === 0) return false
    for (const target of targets) {
      setElementAlignMode(target, mode)
    }
    flattenNestedPreviewAlignments()
    capturePreviewRange()
    void syncModelFromPreview()
    return true
  }
  const unwrapOuterAlignBlock = (raw: string) => {
    let value = String(raw || '').trim()
    if (!value) return ''
    while (true) {
      const mdMatch = value.match(/^:::align-(left|center|right|justify)\s*\n([\s\S]*?)\n:::\s*$/i)
      if (mdMatch) {
        value = String(mdMatch[2] || '').trim()
        continue
      }
      const divMatch = value.match(/^<div\s+align="(left|center|right|justify)"\s*>\s*([\s\S]*?)\s*<\/div>\s*$/i)
      if (divMatch) {
        value = String(divMatch[2] || '').trim()
        continue
      }
      break
    }
    return value
  }
  const applyAlignment = (mode: 'left' | 'center' | 'right' | 'justify') => {
    if (isPreviewEditing()) {
      if (tryUpdateExistingPreviewAlignment(mode) || applyPreviewAlignmentToSelection(mode)) {
        closeMenu()
        return
      }
      closeMenu()
      return
    }
    const state = readSelection()
    const base = unwrapOuterAlignBlock(state.selected || 'текст') || 'текст'
    replaceSelectionByState(state, `:::align-${mode}\n${base}\n:::`, true)
    closeMenu()
  }
  const normalizeInsertUrl = (raw: string) => String(raw || '').trim()
  const extractColorFromStyle = (rawStyle: string) => {
    const entries = String(rawStyle || '').split(';')
    for (const entry of entries) {
      const [propRaw, ...valueParts] = entry.split(':')
      if (
        String(propRaw || '')
          .trim()
          .toLowerCase() !== 'color'
      )
        continue
      return normalizeCssColor(valueParts.join(':'))
    }
    return ''
  }
  const normalizeHexColor = (raw: string) => {
    const value = String(raw || '')
      .trim()
      .toLowerCase()
    if (HEX_COLOR_RE.test(value)) return value
    if (SHORT_HEX_COLOR_RE.test(value)) {
      return `#${value[1]}${value[1]}${value[2]}${value[2]}${value[3]}${value[3]}`
    }
    return ''
  }
  const syncColorHexFromTextInput = (raw: string) => {
    const hex = normalizeHexColor(raw)
    if (!hex || colorHexInput.value === hex) return
    colorHexInput.value = hex
  }
  const onColorHexInput = (event: Event) => {
    const next = normalizeHexColor((event.target as HTMLInputElement).value) || '#22c55e'
    colorHexInput.value = next
    colorValueInput.value = next
    colorError.value = null
  }
  const openColorInput = () => {
    colorError.value = null
    let prefillColor = ''
    if (isPreviewEditing()) {
      capturePreviewRange()
      const range = cloneRange(savedPreviewRange.value)
      colorSelectionRange.value = range
      const currentSpan = findPreviewColorSpanAtRange(range)
      prefillColor = extractColorFromStyle(String(currentSpan?.getAttribute('style') || ''))
    } else {
      const state = readSelection()
      colorSelectionState.value = state
      const token = findMarkdownColorSpanAtSelection(state)
      prefillColor = token?.color || ''
    }
    colorValueInput.value = prefillColor || '#22c55e'
    colorHexInput.value = normalizeHexColor(prefillColor) || '#22c55e'
    toggleMenu('color')
  }
  const applyColor = (rawColor: string) => {
    const color = normalizeCssColor(rawColor)
    if (!color) {
      colorError.value = 'Введите корректный CSS color.'
      return
    }
    colorError.value = null
    if (isPreviewEditing()) {
      const selectedRange = cloneRange(colorSelectionRange.value || savedPreviewRange.value)
      const currentSpan = findPreviewColorSpanAtRange(selectedRange)
      if (currentSpan) {
        currentSpan.style.color = color
        savedPreviewRange.value = cloneRange(selectedRange)
        capturePreviewRange()
        void syncModelFromPreview()
        closeMenu()
        return
      }
      runPreviewCommand('foreColor', color, selectedRange)
      closeMenu()
      return
    }
    const state = colorSelectionState.value || readSelection()
    const existing = findMarkdownColorSpanAtSelection(state)
    if (existing) {
      const text = existing.text || 'text'
      const prefix = `<span style="color:${color}">`
      const replacement = `${prefix}${text}</span>`
      const next = `${state.value.slice(0, existing.start)}${replacement}${state.value.slice(existing.end)}`
      const caretPos = existing.start + prefix.length
      applyNextValue(next, caretPos)
      closeMenu()
      return
    }
    const text = state.selected || 'text'
    const prefix = `<span style="color:${color}">`
    const replacement = `${prefix}${text}</span>`
    const next = `${state.value.slice(0, state.start)}${replacement}${state.value.slice(state.end)}`
    const caretPos = state.start + prefix.length
    applyNextValue(next, caretPos)
    closeMenu()
  }
  const openLinkInput = () => {
    linkError.value = null
    let prefillUrl = ''
    if (isPreviewEditing()) {
      capturePreviewRange()
      const range = cloneRange(savedPreviewRange.value)
      linkSelectionRange.value = range
      prefillUrl = readPreviewUrlAtRange(range, 'link')
    } else {
      const state = readSelection()
      linkSelectionState.value = state
      prefillUrl = readMarkdownUrlAtSelection(state, 'link')
    }
    linkUrlInput.value = prefillUrl || 'https://'
    toggleMenu('link')
  }
  const openImageInput = () => {
    imageError.value = null
    let prefillUrl = ''
    imageWidthInput.value = ''
    imageAlignInput.value = ''
    imageWrapInput.value = false
    if (isPreviewEditing()) {
      capturePreviewRange()
      const range = cloneRange(savedPreviewRange.value)
      imageSelectionRange.value = range
      prefillUrl = readPreviewUrlAtRange(range, 'image')
      const currentImage = findPreviewElementAtRange(range, 'image') as HTMLImageElement | null
      if (currentImage) {
        fillImageInputsFromPreviewElement(currentImage)
      }
    } else {
      const state = readSelection()
      imageSelectionState.value = state
      const token = findMarkdownTokenAtSelection(state, 'image')
      prefillUrl = String(token?.url || '').trim()
      if (token) {
        imageAltInput.value = token.label || 'image'
        imageWidthInput.value = token.width || ''
        imageAlignInput.value = token.align || ''
        imageWrapInput.value = token.wrap
      }
    }
    imageUrlInput.value = prefillUrl || 'https://'
    toggleMenu('image')
  }
  const fillImageInputsFromPreviewElement = (image: HTMLImageElement) => {
    const rawStyle = String(image.getAttribute('style') || '').toLowerCase()
    const styleFloat = rawStyle.match(/(?:^|;)\s*float\s*:\s*(left|right)\s*(?:;|$)/)?.[1] || ''
    const styleAlign = (() => {
      if (/margin-left\s*:\s*auto/.test(rawStyle) && /margin-right\s*:\s*auto/.test(rawStyle)) return 'center'
      if (/margin-left\s*:\s*auto/.test(rawStyle)) return 'right'
      if (/margin-right\s*:\s*auto/.test(rawStyle)) return 'left'
      return ''
    })()
    imageUrlInput.value = String(image.getAttribute('src') || '').trim() || 'https://'
    imageAltInput.value = String(image.getAttribute('alt') || '').trim() || 'image'
    imageWidthInput.value = normalizeImageDisplayWidth(image.getAttribute('data-display-width') || image.style.width)
    imageAlignInput.value = normalizeImageDisplayAlign(
      image.getAttribute('data-display-align') || image.getAttribute('align') || styleFloat || styleAlign
    )
    imageWrapInput.value =
      (normalizeImageTextWrap(image.getAttribute('data-display-wrap')) || Boolean(styleFloat)) &&
      (imageAlignInput.value === 'left' || imageAlignInput.value === 'right')
  }
  const onPreviewClick = (event: MouseEvent) => {
    if (!import.meta.client || !isPreviewEditing()) return
    const editor = previewEditorRef.value
    if (!editor) return
    const image = elementFromPreviewNode(event.target as Node | null, editor, 'image') as HTMLImageElement | null
    if (!image) {
      capturePreviewRange()
      return
    }
    event.preventDefault()
    event.stopPropagation()
    const range = document.createRange()
    range.selectNode(image)
    const selection = window.getSelection()
    if (selection) {
      selection.removeAllRanges()
      selection.addRange(range)
    }
    const clonedRange = cloneRange(range)
    savedPreviewRange.value = clonedRange
    imageSelectionRange.value = clonedRange
    imageError.value = null
    fillImageInputsFromPreviewElement(image)
    menu.value = 'image'
  }
  const submitLink = () => {
    const url = normalizeInsertUrl(linkUrlInput.value)
    if (!url) {
      linkError.value = 'Введите ссылку.'
      return
    }
    linkError.value = null
    if (isPreviewEditing()) {
      const selectedRange = cloneRange(linkSelectionRange.value || savedPreviewRange.value)
      const currentAnchor = findPreviewElementAtRange(selectedRange, 'link') as HTMLAnchorElement | null
      if (currentAnchor) {
        currentAnchor.setAttribute('href', url)
        capturePreviewRange()
        void syncModelFromPreview()
        closeMenu()
        return
      }
      const selectedText = selectedRange?.toString() || url
      insertPreviewHtml(buildAnchorHtml(url, selectedText), selectedRange)
      closeMenu()
      return
    }
    const state = linkSelectionState.value || readSelection()
    if (replaceMarkdownTokenUrlByState(state, 'link', url)) {
      closeMenu()
      return
    }
    const selectedText = state.selected || url
    replaceSelectionByState(state, `[${selectedText}](${url})`, true)
    closeMenu()
  }
  const submitImage = () => {
    const url = normalizeInsertUrl(imageUrlInput.value)
    if (!url) {
      imageError.value = 'Введите ссылку.'
      return
    }
    imageError.value = null
    const alt = imageAltInput.value.trim() || 'image'
    const width = normalizeImageDisplayWidth(imageWidthInput.value)
    const align = normalizeImageDisplayAlign(imageAlignInput.value)
    const wrap = normalizeImageTextWrap(imageWrapInput.value) && (align === 'left' || align === 'right')
    if (isPreviewEditing()) {
      const selectedRange = cloneRange(imageSelectionRange.value || savedPreviewRange.value)
      const currentImage = findPreviewElementAtRange(selectedRange, 'image') as HTMLImageElement | null
      if (currentImage) {
        currentImage.setAttribute('src', url)
        currentImage.setAttribute('alt', alt)
        currentImage.setAttribute('data-display-width', width)
        currentImage.setAttribute('data-display-align', align)
        currentImage.setAttribute('data-display-wrap', wrap ? '1' : '0')
        currentImage.setAttribute('style', buildImageStyleAttr(width, align, wrap))
        capturePreviewRange()
        void syncModelFromPreview()
        closeMenu()
        return
      }
      insertPreviewHtml(buildImageHtml(url, alt, width, align, wrap), selectedRange)
      closeMenu()
      return
    }
    const state = imageSelectionState.value || readSelection()
    if (replaceMarkdownImageTokenByState(state, url, alt, width, align, wrap)) {
      closeMenu()
      return
    }
    replaceSelectionByState(state, buildMarkdownImageToken(alt, url, width, align, wrap))
    closeMenu()
  }
  const normalizeCodeIndentation = (raw: string) => {
    const lines = String(raw || '')
      .replace(/\r/g, '')
      .split('\n')
    while (lines.length > 0) {
      const first = lines[0]
      if (!first || first.trim()) break
      lines.shift()
    }
    while (lines.length > 0) {
      const last = lines[lines.length - 1]
      if (!last || last.trim()) break
      lines.pop()
    }
    if (lines.length === 0) return ''
    const indents = lines
      .filter(line => line.trim())
      .map(line => {
        const leading = (line.match(/^[\t ]*/) || [''])[0].replace(/\t/g, '  ')
        return leading.length
      })
    const minIndent = indents.length > 0 ? Math.min(...indents) : 0
    return lines
      .map(line => {
        const expanded = line.replace(/\t/g, '  ')
        return expanded.slice(Math.min(minIndent, expanded.length))
      })
      .join('\n')
  }
  const insertCodeBlock = (language: string) => {
    const lang = language.trim().toLowerCase() || 'text'
    if (isPreviewEditing()) {
      const selectedRange = cloneRange(savedPreviewRange.value)
      const selectedText = selectedRange?.toString() || ''
      const normalized = normalizeCodeIndentation(selectedText)
      const selected = normalized || (lang === 'text' ? 'text' : `// ${lang} code`)
      const langClass = ` class="language-${escapeAttr(lang)}"`
      const blockHtml = `<pre spellcheck="false"><code${langClass}>${escapeHtml(selected)}</code></pre><p><br></p>`
      insertPreviewHtml(blockHtml, selectedRange, true)
      closeMenu()
      return
    }
    const state = readSelection()
    const normalized = normalizeCodeIndentation(state.selected)
    const selected = normalized || (lang === 'text' ? 'text' : `// ${lang} code`)
    const block = `\`\`\`${lang}\n${selected}\n\`\`\`\n`
    replaceSelection(block)
    closeMenu()
  }
  const insertEmoji = (name: string) => {
    const token = `:emojione:${name}:`
    if (isPreviewEditing()) {
      const src = emojiDataUriResolver.value(name)
      if (src) {
        const html = `<img class="emoji-inline" src="${src}" alt="${escapeAttr(name)}" loading="lazy" decoding="async">`
        insertPreviewHtml(html)
      } else {
        insertPreviewText(token)
      }
      closeMenu()
      return
    }
    replaceSelection(token)
    closeMenu()
  }
  const onLanguageScroll = (event: Event) => {
    const target = event.target as HTMLElement
    if (target.scrollTop + target.clientHeight < target.scrollHeight - 20) return
    languageBatch.value = Math.min(languageBatch.value + 20, filteredLanguages.value.length)
  }
  const onEmojiScroll = (event: Event) => {
    const target = event.target as HTMLElement
    if (target.scrollTop + target.clientHeight < target.scrollHeight - 20) return
    if (emojiScrollRaf.value && import.meta.client) return
    if (!import.meta.client) {
      emojiBatch.value = Math.min(emojiBatch.value + EMOJI_BATCH_STEP, filteredEmojiNames.value.length)
      return
    }
    emojiScrollRaf.value = window.requestAnimationFrame(() => {
      emojiBatch.value = Math.min(emojiBatch.value + EMOJI_BATCH_STEP, filteredEmojiNames.value.length)
      emojiScrollRaf.value = null
    })
  }
  onMounted(() => {
    void ensureEmojiCollectionForModel(props.modelValue)
    const cached = readEmojiNamesCache()
    if (cached.length > 0) {
      emojiNames.value = cached
    }
    document.addEventListener('click', onDocumentClick, true)
    nextTick(() => {
        void syncPreviewFromModel(true)
    })
  })
  onBeforeUnmount(() => {
    document.removeEventListener('click', onDocumentClick, true)
    clearHideMenuTimer()
    if (emojiScrollRaf.value && import.meta.client) {
      window.cancelAnimationFrame(emojiScrollRaf.value)
      emojiScrollRaf.value = null
    }
  })
</script>
<template>
  <div class="space-y-2">
    <div ref="toolbarRef" class="toolbar-shell" @mousedown="onToolbarMouseDown">
      <div class="toolbar-row">
        <div class="tool-group" @mouseenter="openMenuOnHover('textStyle')" @mouseleave="onGroupLeave('textStyle')">
          <button
            type="button"
            class="tool-btn"
            :class="{ 'tool-btn-active': menu === 'textStyle' }"
            title="Стили текста"
            @mouseenter="openMenuOnHover('textStyle')"
            @click.stop="toggleMenu('textStyle')">
            <Icon name="ic:round-format-bold" class="tool-icon" />
          </button>
          <div v-if="menu === 'textStyle'" class="tool-pop">
            <button
              v-for="option in STYLE_OPTIONS"
              :key="option.label"
              type="button"
              class="tool-option"
              @click="applyInlineStyle(option)">
              <Icon :name="option.icon" class="text-sm" />
              <span>{{ option.label }}</span>
            </button>
          </div>
        </div>
        <div class="tool-group" @mouseenter="openMenuOnHover('heading')" @mouseleave="onGroupLeave('heading')">
          <button
            type="button"
            class="tool-btn"
            :class="{ 'tool-btn-active': menu === 'heading' }"
            title="Заголовки"
            @mouseenter="openMenuOnHover('heading')"
            @click.stop="toggleMenu('heading')">
            <Icon name="ic:baseline-format-size" class="tool-icon" />
          </button>
          <div v-if="menu === 'heading'" class="tool-pop tool-pop-wide">
            <button
              v-for="option in HEADING_OPTIONS"
              :key="`h-${option.level}`"
              type="button"
              class="tool-option-heading"
              @click="applyHeading(option.level)">
              <span class="heading-preview" :style="{ fontSize: option.size }">Образец</span>
              <span class="heading-level">H{{ option.level }}</span>
            </button>
          </div>
        </div>
        <div class="tool-group" @mouseenter="openMenuOnHover('align')" @mouseleave="onGroupLeave('align')">
          <button
            type="button"
            class="tool-btn"
            :class="{ 'tool-btn-active': menu === 'align' }"
            title="Выравнивание"
            @mouseenter="openMenuOnHover('align')"
            @click.stop="toggleMenu('align')">
            <Icon name="ic:round-format-align-left" class="tool-icon" />
          </button>
          <div v-if="menu === 'align'" class="tool-pop">
            <button
              v-for="option in ALIGN_OPTIONS"
              :key="option.mode"
              type="button"
              class="tool-option"
              @click="applyAlignment(option.mode)">
              <Icon :name="option.icon" class="text-sm" />
              <span>{{ option.label }}</span>
            </button>
          </div>
        </div>
        <div ref="colorToolRef" class="tool-group" @mouseleave="onGroupLeave('color')">
          <button
            type="button"
            class="tool-btn"
            :class="{ 'tool-btn-active': menu === 'color' }"
            title="Цвет текста"
            @click.stop="openColorInput">
            <Icon name="ic:round-format-color-text" class="tool-icon" />
          </button>
          <div v-if="menu === 'color'" class="tool-pop tool-pop-inline-form tool-pop-color">
            <div class="tool-color-grid">
              <button
                v-for="preset in COLOR_PRESETS"
                :key="preset"
                type="button"
                class="tool-color-chip"
                :title="preset"
                :style="{ backgroundColor: preset }"
                @click="applyColor(preset)" />
            </div>
            <div class="tool-inline-row tool-inline-row-color">
              <input
                v-model="colorHexInput"
                name="md_text_color_hex_input"
                type="color"
                class="tool-color-picker"
                aria-label="HEX color picker"
                @input="onColorHexInput" />
              <span class="tool-color-hex">{{ colorHexInput.toUpperCase() }}</span>
            </div>
            <div class="tool-inline-row">
              <input
                v-model="colorValueInput"
                name="md_text_color_input"
                type="text"
                class="tool-search tool-search-compact"
                placeholder="#22c55e или tomato"
                @keydown.enter.prevent="applyColor(colorValueInput)" />
              <button type="button" class="tool-ok-btn" @click="applyColor(colorValueInput)">OK</button>
            </div>
            <LabNotify :text="colorError" tone="error" size="xs" class-name="tool-error" />
          </div>
        </div>
        <div ref="linkToolRef" class="tool-group" @mouseleave="onGroupLeave('link')">
          <button
            type="button"
            class="tool-btn"
            :class="{ 'tool-btn-active': menu === 'link' }"
            title="Ссылка"
            @click.stop="openLinkInput">
            <Icon name="ic:round-link" class="tool-icon" />
          </button>
          <div v-if="menu === 'link'" class="tool-pop tool-pop-inline-form">
            <div class="tool-inline-row">
              <input
                v-model="linkUrlInput"
                name="md_link_url_input"
                type="url"
                class="tool-search tool-search-compact"
                placeholder="https://example.com"
                @keydown.enter.prevent="submitLink" />
              <button type="button" class="tool-ok-btn" @click="submitLink">OK</button>
            </div>
            <LabNotify :text="linkError" tone="error" size="xs" class-name="tool-error" />
          </div>
        </div>
        <div ref="imageToolRef" class="tool-group" @mouseleave="onGroupLeave('image')">
          <button
            type="button"
            class="tool-btn"
            :class="{ 'tool-btn-active': menu === 'image' }"
            title="Изображение"
            @click.stop="openImageInput">
            <Icon name="ic:round-image" class="tool-icon" />
          </button>
          <div v-if="menu === 'image'" class="tool-pop tool-pop-inline-form">
            <div class="tool-inline-row">
              <input
                v-model="imageUrlInput"
                name="md_image_url_input"
                type="url"
                class="tool-search tool-search-compact"
                placeholder="https://site/image.png"
                @keydown.enter.prevent="submitImage" />
              <button type="button" class="tool-ok-btn" @click="submitImage">OK</button>
            </div>
            <input
              v-model="imageAltInput"
              name="md_image_alt_input"
              type="text"
              class="tool-search tool-search-compact tool-search-alt"
              placeholder="Alt" />
            <input
              v-model="imageWidthInput"
              name="md_image_width_input"
              type="number"
              :min="min_image_width_px"
              max="1600"
              class="tool-search tool-search-compact tool-search-alt"
              placeholder="Width px, например 320" />
            <select
              v-model="imageAlignInput"
              name="md_image_align_input"
              class="tool-search tool-search-compact tool-search-alt">
              <option value="">Без выравнивания</option>
              <option value="left">Влево</option>
              <option value="center">По центру</option>
              <option value="right">Вправо</option>
            </select>
            <label class="tool-check-row">
              <input
                v-model="imageWrapInput"
                name="md_image_wrap_input"
                type="checkbox"
                :disabled="imageAlignInput !== 'left' && imageAlignInput !== 'right'" />
              <span>Обтекание текстом</span>
            </label>
            <LabNotify :text="imageError" tone="error" size="xs" class-name="tool-error" />
          </div>
        </div>
        <button type="button" class="tool-btn" title="Маркированный список" @click="applyList(false)">
          <Icon name="ic:round-format-list-bulleted" class="tool-icon" />
        </button>
        <button type="button" class="tool-btn" title="Нумерованный список" @click="applyList(true)">
          <Icon name="ic:round-format-list-numbered" class="tool-icon" />
        </button>
        <button type="button" class="tool-btn" title="Очистить форматирование" @click="clearFormatting">
          <Icon name="ic:round-format-clear" class="tool-icon" />
        </button>
        <div class="tool-group" @mouseenter="openMenuOnHover('language')" @mouseleave="onGroupLeave('language')">
          <button
            type="button"
            class="tool-btn"
            :class="{ 'tool-btn-active': menu === 'language' }"
            title="Блок кода"
            @mouseenter="openMenuOnHover('language')"
            @click.stop="toggleMenu('language')">
            <Icon name="ic:round-code" class="tool-icon" />
          </button>
          <div v-if="menu === 'language'" class="tool-pop tool-pop-wide">
            <input
              v-model="languageFilter"
              name="md_language_filter"
              type="text"
              class="tool-search"
              placeholder="Поиск языка (python, js, go...)" />
            <div class="tool-scroll" @scroll="onLanguageScroll">
              <button
                v-for="lang in visibleLanguages"
                :key="`${lang.key}-${lang.category}`"
                type="button"
                class="tool-option tool-option-language"
                @click="insertCodeBlock(lang.key)">
                <span class="text-xs text-zinc-100">{{ lang.label }}</span>
                <span class="tool-meta">{{ lang.category }}</span>
              </button>
              <button type="button" class="tool-option tool-option-language" @click="insertCodeBlock('text')">
                <span class="text-xs text-zinc-100">Text</span>
              </button>
              <p class="tool-info">Показано {{ visibleLanguages.length }} из {{ filteredLanguages.length }}</p>
            </div>
          </div>
        </div>
        <div class="tool-group" @mouseenter="openMenuOnHover('emoji')" @mouseleave="onGroupLeave('emoji')">
          <button
            type="button"
            class="tool-btn"
            :class="{ 'tool-btn-active': menu === 'emoji' }"
            title="Emoji"
            @mouseenter="openMenuOnHover('emoji')"
            @click.stop="toggleMenu('emoji')">
            <Icon name="ic:round-emoji-emotions" class="tool-icon" />
          </button>
          <div v-if="menu === 'emoji'" class="tool-pop tool-pop-wide tool-pop-right">
            <input
              v-model="emojiFilter"
              name="md_emoji_filter"
              type="text"
              class="tool-search"
              placeholder="Поиск face-emoji (grinning-face, thinking-face...)" />
            <LabNotify :text="emojiLoadError" tone="warning" size="xs" class-name="tool-warning" />
            <LabNotify
              v-if="emojiLoading"
              text="Загрузка face-набора emoji..."
              tone="info"
              size="xs"
              as="div"
              class-name="tool-info" />
            <div v-else class="tool-scroll tool-scroll-emoji" @scroll="onEmojiScroll">
              <div class="emoji-grid">
                <button
                  v-for="name in visibleEmojiNames"
                  :key="name"
                  type="button"
                  class="emoji-btn"
                  :title="name"
                  @click="insertEmoji(name)">
                  <img
                    v-if="emoji_data_uri(name)"
                    :src="emoji_data_uri(name)"
                    :alt="name"
                    class="mx-auto h-5 w-5"
                    loading="lazy"
                    decoding="async" />
                  <span v-else class="block truncate text-[9px] font-medium leading-tight text-zinc-200">
                    {{ name }}
                  </span>
                </button>
              </div>
              <p class="tool-info">Показано {{ visibleEmojiNames.length }} из {{ filteredEmojiNames.length }}</p>
            </div>
          </div>
        </div>
        <button
          type="button"
          class="tool-btn tool-btn-preview"
          :title="`Режим: ${currentPreviewMeta.label}. Нажмите для переключения`"
          @click="cyclePreviewMode">
          <Icon :name="currentPreviewMeta.icon" class="tool-icon" />
          <span class="tool-hint">{{ currentPreviewMeta.shortLabel }}</span>
        </button>
      </div>
    </div>
    <div class="editor-layout" :class="`editor-layout-${previewMode}`">
      <div v-if="previewMode !== 'preview'" class="editor-pane editor-pane-input">
        <textarea
          ref="textareaRef"
          name="md_editor_input"
          :value="props.modelValue"
          :rows="rows"
          class="editor-input"
          placeholder="Полная статья (Markdown поддерживается)"
          @input="onInput" />
      </div>
      <div v-if="previewMode !== 'edit'" class="editor-pane editor-pane-preview">
        <div class="preview-box">
          <div
            ref="previewEditorRef"
            class="preview-editable md-view"
            contenteditable="true"
            spellcheck="true"
            data-placeholder="Редактируйте прямо в опубликованном виде"
            @input="onPreviewInput"
            @keydown="onPreviewKeydown"
            @keyup="capturePreviewRange"
            @mouseup="capturePreviewRange"
            @focus="onPreviewFocus"
            @blur="onPreviewBlur"
            @click="onPreviewClick" />
        </div>
      </div>
    </div>
  </div>
</template>
<style scoped>
  .toolbar-shell {
    position: relative;
    z-index: 25;
    overflow: visible;
    max-width: 100%;
  }
  .toolbar-row {
    display: flex;
    align-items: center;
    gap: 0.22rem;
    flex-wrap: wrap;
    overflow: visible;
    padding-bottom: 0;
    row-gap: 0.3rem;
  }
  .tool-group {
    position: relative;
    flex-shrink: 0;
  }
  .tool-btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 0.18rem;
    border: 1px solid rgb(63 63 70);
    background: rgb(39 39 42);
    border-radius: 0.42rem;
    color: rgb(228 228 231);
    width: 2rem;
    height: 2rem;
    flex-shrink: 0;
  }
  .tool-btn:hover {
    background: rgb(63 63 70);
  }
  .tool-btn-active {
    border-color: rgb(99 102 241);
    background: rgba(67, 56, 202, 0.35);
  }
  .tool-btn-preview {
    width: auto;
    padding-inline: 0.35rem;
  }
  .tool-icon {
    font-size: 1rem;
    line-height: 1;
  }
  .tool-hint {
    font-size: 0.58rem;
    font-weight: 600;
    color: rgb(161 161 170);
  }
  .tool-pop {
    position: absolute;
    top: calc(100% + 0.35rem);
    left: 0;
    z-index: 60;
    min-width: 11rem;
    background: rgb(9 9 11);
    border: 1px solid rgb(63 63 70);
    border-radius: 0.55rem;
    padding: 0.45rem;
  }
  .tool-pop-wide {
    min-width: 19rem;
  }
  .tool-pop-inline-form {
    min-width: 20rem;
  }
  .tool-pop-color {
    min-width: 14rem;
  }
  .tool-pop-right {
    right: 0;
    left: auto;
  }
  .tool-option,
  .tool-option-heading {
    width: 100%;
    display: flex;
    align-items: center;
    gap: 0.45rem;
    border: 1px solid rgb(63 63 70);
    background: rgb(24 24 27);
    color: rgb(228 228 231);
    border-radius: 0.38rem;
    padding: 0.3rem 0.45rem;
    font-size: 0.72rem;
    line-height: 1.2;
    text-align: left;
  }
  .tool-option + .tool-option,
  .tool-option-heading + .tool-option-heading {
    margin-top: 0.28rem;
  }
  .tool-option:hover,
  .tool-option-heading:hover {
    background: rgb(39 39 42);
  }
  .tool-option-language {
    justify-content: space-between;
  }
  .heading-preview {
    color: rgb(244 244 245);
    font-weight: 600;
    line-height: 1.15;
  }
  .heading-level {
    margin-left: auto;
    font-size: 0.62rem;
    color: rgb(161 161 170);
  }
  .tool-meta {
    margin-left: auto;
    color: rgb(161 161 170);
    font-size: 0.66rem;
  }
  .tool-search {
    width: 100%;
    border: 1px solid rgb(63 63 70);
    background: rgb(24 24 27);
    border-radius: 0.35rem;
    color: rgb(228 228 231);
    padding: 0.34rem 0.45rem;
    font-size: 0.72rem;
    margin-bottom: 0.4rem;
  }
  .tool-search-compact {
    margin-bottom: 0;
  }
  .tool-search-alt {
    margin-top: 0.35rem;
  }
  .tool-check-row {
    margin-top: 0.35rem;
    display: flex;
    align-items: center;
    gap: 0.4rem;
    font-size: 0.7rem;
    color: rgb(212 212 216);
  }
  .tool-check-row input[type='checkbox'] {
    accent-color: rgb(99 102 241);
  }
  .tool-check-row input[type='checkbox']:disabled + span {
    opacity: 0.5;
  }
  .tool-inline-row {
    display: flex;
    align-items: center;
    gap: 0.35rem;
  }
  .tool-inline-row-color {
    margin-bottom: 0.35rem;
  }
  .tool-color-picker {
    width: 2.2rem;
    height: 1.9rem;
    border: 1px solid rgb(63 63 70);
    border-radius: 0.35rem;
    background: rgb(24 24 27);
    padding: 0.12rem;
  }
  .tool-color-hex {
    font-size: 0.68rem;
    color: rgb(161 161 170);
    letter-spacing: 0.03em;
    font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
  }
  .tool-color-grid {
    display: grid;
    grid-template-columns: repeat(6, minmax(0, 1fr));
    gap: 0.35rem;
    margin-bottom: 0.4rem;
  }
  .tool-color-chip {
    width: 100%;
    aspect-ratio: 1;
    border-radius: 999px;
    border: 1px solid rgb(82 82 91);
  }
  .tool-color-chip:hover {
    transform: scale(1.04);
  }
  .tool-ok-btn {
    border: 1px solid rgb(63 63 70);
    background: rgb(63 63 70);
    color: rgb(244 244 245);
    border-radius: 0.35rem;
    height: 1.9rem;
    padding: 0 0.6rem;
    font-size: 0.7rem;
    font-weight: 600;
    white-space: nowrap;
  }
  .tool-ok-btn:hover {
    background: rgb(82 82 91);
  }
  .tool-ok-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
  .tool-scroll {
    max-height: 12rem;
    overflow-y: auto;
    padding-right: 0.2rem;
  }
  .tool-scroll-emoji {
    max-height: 14rem;
  }
  .tool-info {
    margin-top: 0.38rem;
    font-size: 0.64rem;
    color: rgb(113 113 122);
  }
  .tool-warning {
    margin-bottom: 0.3rem;
    font-size: 0.66rem;
    color: rgb(252 211 77);
  }
  .tool-error {
    margin-top: 0.35rem;
    font-size: 0.66rem;
    color: rgb(251 113 133);
  }
  .emoji-grid {
    display: grid;
    grid-template-columns: repeat(8, minmax(0, 1fr));
    gap: 0.25rem;
  }
  .emoji-btn {
    border: 1px solid rgb(63 63 70);
    background: rgb(24 24 27);
    border-radius: 0.35rem;
    padding: 0.2rem;
    color: rgb(228 228 231);
  }
  .emoji-btn:hover {
    background: rgb(39 39 42);
  }
  .editor-layout {
    display: grid;
    gap: 0.45rem;
    min-height: 360px;
  }
  .editor-layout-edit,
  .editor-layout-preview {
    grid-template-columns: minmax(0, 1fr);
  }
  .editor-input {
    width: 100%;
    min-height: 320px;
    border: 1px solid rgb(63 63 70);
    background: rgb(24 24 27);
    border-radius: 0.65rem;
    padding: 0.6rem 0.75rem;
    color: rgb(228 228 231);
    font-size: 0.88rem;
    resize: vertical;
    overflow: auto;
  }
  .editor-pane {
    min-height: 320px;
  }
  .editor-pane-input {
    overflow: visible;
  }
  .editor-pane-preview {
    overflow: visible;
  }
  .preview-box {
    min-height: 100%;
    border: 1px solid rgb(63 63 70);
    background: rgba(24, 24, 27, 0.88);
    border-radius: 0.65rem;
    padding: 0.5rem 0.65rem;
  }
  .preview-editable {
    width: 100%;
    min-height: 290px;
    border: 0;
    background: transparent;
    border-radius: 0;
    padding: 0.2rem 0.25rem;
    color: rgb(212 212 216);
    font-size: 0.92rem;
    line-height: 1.6;
    resize: vertical;
    overflow: auto;
    outline: none;
  }
  .preview-editable:focus {
    outline: 1px solid rgba(99, 102, 241, 0.7);
    border-radius: 0.35rem;
  }
  .preview-editable:empty::before {
    content: attr(data-placeholder);
    color: rgb(113 113 122);
  }
  .preview-editable :deep(h1),
  .preview-editable :deep(h2),
  .preview-editable :deep(h3),
  .preview-editable :deep(h4),
  .preview-editable :deep(h5),
  .preview-editable :deep(h6) {
    color: rgb(244 244 245);
    font-weight: 600;
    margin-top: 1.1rem;
    margin-bottom: 0.5rem;
  }
  .preview-editable :deep(h1) {
    font-size: 1.55rem;
  }
  .preview-editable :deep(h2) {
    font-size: 1.25rem;
  }
  .preview-editable :deep(h3) {
    font-size: 1.1rem;
  }
  .preview-editable :deep(h4),
  .preview-editable :deep(h5),
  .preview-editable :deep(h6) {
    font-size: 1rem;
  }
  .preview-editable :deep(p) {
    margin-bottom: 0.55rem;
  }
  .preview-editable :deep(ul),
  .preview-editable :deep(ol) {
    margin-bottom: 0.65rem;
    padding-left: 1.6rem;
    margin-left: 0.2rem;
  }
  .preview-editable :deep(ul) {
    list-style: disc;
  }
  .preview-editable :deep(ol) {
    list-style: decimal;
  }
  .preview-editable :deep(li) {
    margin-bottom: 0.3rem;
  }
  .preview-editable :deep(code) {
    background: rgba(24, 24, 27, 0.95);
    border: 1px solid rgba(63, 63, 70, 0.9);
    border-radius: 0.35rem;
    color: rgb(209 213 219);
    padding: 0.1rem 0.35rem;
    font-size: 0.9em;
    font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
  }
  .preview-editable :deep(.md-code-block) {
    margin-bottom: 0.75rem;
    border: 1px solid rgba(63, 63, 70, 0.9);
    border-radius: 0.6rem;
    overflow: hidden;
    background: rgba(24, 24, 27, 0.95);
  }
  .preview-editable :deep(.md-code-block-label) {
    padding: 0.38rem 0.7rem;
    border-bottom: 1px solid rgba(63, 63, 70, 0.9);
    background: rgba(39, 39, 42, 0.95);
    color: rgb(161 161 170);
    font-size: 0.7rem;
    line-height: 1;
    font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
    text-transform: lowercase;
    letter-spacing: 0.04em;
  }
  .preview-editable :deep(pre) {
    background: transparent;
    border: 0;
    border-radius: 0;
    padding: 0.75rem;
    overflow-x: auto;
    margin-bottom: 0;
  }
  .preview-editable :deep(pre code) {
    background: transparent;
    border: 0;
    padding: 0;
  }
  .preview-editable :deep(blockquote) {
    border-left: 3px solid rgba(99, 102, 241, 0.7);
    padding-left: 0.75rem;
    color: rgb(161 161 170);
    margin-bottom: 0.7rem;
  }
  .preview-editable :deep(a) {
    color: rgb(129 140 248);
    text-decoration: underline;
  }
  .preview-editable :deep(img) {
    max-width: min(100%, 60rem);
    width: auto;
    border-radius: 0.45rem;
    border: 1px solid rgba(63, 63, 70, 0.8);
  }
  .preview-editable :deep(img.emoji-inline) {
    width: 1.2em;
    height: 1.2em;
    vertical-align: -0.2em;
    border: 0;
    border-radius: 0;
    display: inline-block;
  }
  .preview-editable :deep(.md-align-left) {
    text-align: left;
  }
  .preview-editable :deep(.md-align-center) {
    text-align: center;
  }
  .preview-editable :deep(.md-align-right) {
    text-align: right;
  }
  .preview-editable :deep(.md-align-justify) {
    text-align: justify;
  }
  @media (min-width: 640px) {
    .toolbar-row {
      flex-wrap: nowrap;
      min-width: max-content;
    }
    .editor-layout {
      min-height: 420px;
    }
    .editor-input,
    .editor-pane {
      min-height: 420px;
    }
    .preview-editable {
      min-height: 390px;
    }
  }
</style>
