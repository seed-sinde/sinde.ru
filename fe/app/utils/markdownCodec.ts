import hljs from 'highlight.js/lib/common'
import stylus from 'highlight.js/lib/languages/stylus'
import highlightjsVueModule from 'highlightjs-vue'
type AlignMode = '' | 'left' | 'center' | 'right' | 'justify'
type EmojiIconData = {
  body: string
  left?: number
  top?: number
  width?: number
  height?: number
}
type EmojiAliasData = {
  parent: string
}
let emojiIcons: Record<string, EmojiIconData> = {}
let emojiAliases: Record<string, EmojiAliasData> = {}
let emojiDefaultWidth = 64
let emojiDefaultHeight = 64
let emojiCollectionLoaded = false
let emojiCollectionPromise: Promise<void> | null = null
const emojiDataUriCache = new Map<string, string>()
const DEFAULT_IMAGE_MAX_WIDTH_PX = 960
export const MIN_IMAGE_WIDTH_PX = 32
const MAX_IMAGE_WIDTH_PX = 1600
type HighlightVueRegistrar = (hljsInstance: typeof hljs) => void
const resolveVueLanguageRegistrar = (value: unknown): HighlightVueRegistrar | null => {
  if (typeof value === 'function') return value as HighlightVueRegistrar
  if (!value || typeof value !== 'object') return null
  const record = value as Record<string, unknown>
  if (typeof record.default === 'function') return record.default as HighlightVueRegistrar
  if (typeof record.registerVueLanguage === 'function') return record.registerVueLanguage as HighlightVueRegistrar
  return null
}
const registerVueLanguage = resolveVueLanguageRegistrar(highlightjsVueModule)
if (!hljs.getLanguage('stylus')) {
  hljs.registerLanguage('stylus', stylus)
}
if (!hljs.getLanguage('vue')) {
  registerVueLanguage?.(hljs)
}
const applyEmojiCollection = (collectionRaw: unknown) => {
  const collection = (collectionRaw || {}) as Record<string, any>
  emojiIcons = (collection.icons || {}) as Record<string, EmojiIconData>
  emojiAliases = (collection.aliases || {}) as Record<string, EmojiAliasData>
  emojiDefaultWidth = Number(collection.width || 64)
  emojiDefaultHeight = Number(collection.height || emojiDefaultWidth || 64)
}
export const ensureEmojiCollectionLoaded = async () => {
  if (emojiCollectionLoaded) return
  if (emojiCollectionPromise) {
    await emojiCollectionPromise
    return
  }
  emojiCollectionPromise = loadEmojiCollection()
    .then(collection => {
      applyEmojiCollection(collection)
    })
    .finally(() => {
      emojiCollectionLoaded = true
      emojiCollectionPromise = null
    })
  await emojiCollectionPromise
}
export const normalizeCssColor = (raw: string) => {
  const value = String(raw || '')
    .trim()
    .toLowerCase()
  if (!value) return ''
  if (/^#[0-9a-f]{3,8}$/i.test(value)) return value
  if (/^rgba?\(\s*[-\d.%\s,]+\)$/.test(value)) return value
  if (/^hsla?\(\s*[-\d.%\s,]+\)$/.test(value)) return value
  if (/^[a-z]{3,20}$/.test(value)) return value
  return ''
}
const sanitizeSpanColorStyle = (rawStyle: string) => {
  const entries = String(rawStyle || '')
    .split(';')
    .map(part => part.trim())
    .filter(Boolean)
  let color = ''
  let background = ''
  for (const entry of entries) {
    const [propRaw, ...valueParts] = entry.split(':')
    const prop = String(propRaw || '')
      .trim()
      .toLowerCase()
    const value = normalizeCssColor(valueParts.join(':'))
    if (!value) continue
    if (prop === 'color') {
      color = value
      continue
    }
    if (prop === 'background' || prop === 'background-color') {
      background = value
    }
  }
  const out: string[] = []
  if (color) out.push(`color: ${color}`)
  if (background) out.push(`background-color: ${background}`)
  return out.join('; ')
}
export const escapeHtml = (raw: string) =>
  String(raw || '')
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
export const escapeAttr = (raw: string) => String(raw || '').replace(/[^a-z0-9_+.-]/gi, '')
const escapeAttrValue = (raw: string) =>
  String(raw || '')
    .replaceAll('&', '&amp;')
    .replaceAll('"', '&quot;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
const CODE_LANGUAGE_ALIASES: Record<string, { highlight: string; display: string }> = {
  js: { highlight: 'javascript', display: 'javascript' },
  jsx: { highlight: 'javascript', display: 'javascript' },
  ts: { highlight: 'typescript', display: 'typescript' },
  tsx: { highlight: 'typescript', display: 'typescript' },
  yml: { highlight: 'yaml', display: 'yaml' },
  sh: { highlight: 'bash', display: 'bash' },
  shell: { highlight: 'bash', display: 'bash' },
  zsh: { highlight: 'bash', display: 'bash' },
  md: { highlight: 'markdown', display: 'markdown' },
  mdx: { highlight: 'markdown', display: 'markdown' },
  html: { highlight: 'xml', display: 'html' },
  svg: { highlight: 'xml', display: 'svg' },
  dotenv: { highlight: 'bash', display: 'env' },
  env: { highlight: 'bash', display: 'env' },
  dockerfile: { highlight: 'docker', display: 'docker' },
  postgres: { highlight: 'sql', display: 'postgresql' },
  postgressql: { highlight: 'sql', display: 'postgresql' },
  sql: { highlight: 'sql', display: 'sql' }
}
export type HighlightedCodeSnippet = {
  highlightLanguage: string
  displayLanguage: string
  highlightedHtml: string
}
export const normalizeImageDisplayWidth = (raw: string | number | null | undefined) => {
  const value = typeof raw === 'number' ? raw : Number.parseInt(String(raw || '').trim(), 10)
  if (!Number.isFinite(value)) return ''
  const normalized = Math.max(MIN_IMAGE_WIDTH_PX, Math.min(MAX_IMAGE_WIDTH_PX, Math.round(value)))
  return String(normalized)
}
export const normalizeImageDisplayAlign = (raw: string | null | undefined): ImageAlignMode => {
  const value = String(raw || '')
    .trim()
    .toLowerCase()
  if (value === 'left' || value === 'center' || value === 'right') return value
  return ''
}
export const normalizeImageTextWrap = (raw: unknown) => {
  if (typeof raw === 'boolean') return raw
  const value = String(raw ?? '')
    .trim()
    .toLowerCase()
  return value === '1' || value === 'true' || value === 'yes' || value === 'on'
}
const parseMarkdownImageOptions = (rawMeta?: string) => {
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
const buildImageMetaSuffix = (width?: string, align?: ImageAlignMode, wrap?: boolean) => {
  const normalizedWidth = normalizeImageDisplayWidth(width)
  const normalizedAlign = normalizeImageDisplayAlign(align)
  const normalizedWrap = normalizeImageTextWrap(wrap) && (normalizedAlign === 'left' || normalizedAlign === 'right')
  const parts: string[] = []
  if (normalizedWidth) parts.push(`width=${normalizedWidth}`)
  if (normalizedAlign) parts.push(`align=${normalizedAlign}`)
  if (normalizedWrap) parts.push('wrap=1')
  return parts.length > 0 ? `{${parts.join(' ')}}` : ''
}
export const buildImageStyleAttr = (width?: string, align?: ImageAlignMode, wrap?: boolean) => {
  const normalized = normalizeImageDisplayWidth(width)
  const normalizedAlign = normalizeImageDisplayAlign(align)
  const normalizedWrap = normalizeImageTextWrap(wrap) && (normalizedAlign === 'left' || normalizedAlign === 'right')
  const styles: string[] = []
  if (normalized) {
    styles.push(`width:min(${normalized}px, 100%)`, `max-width:${normalized}px`)
  } else {
    styles.push(`max-width:min(100%, ${DEFAULT_IMAGE_MAX_WIDTH_PX}px)`)
  }
  if (normalizedWrap) {
    styles.push(`float:${normalizedAlign}`)
    if (normalizedAlign === 'left') {
      styles.push('margin:0.25rem 0.9rem 0.45rem 0')
    } else {
      styles.push('margin:0.25rem 0 0.45rem 0.9rem')
    }
  } else {
    styles.push('float:none')
    if (normalizedAlign === 'center') {
      styles.push('display:block', 'margin-left:auto', 'margin-right:auto')
    } else if (normalizedAlign === 'right') {
      styles.push('display:block', 'margin-left:auto', 'margin-right:0')
    } else if (normalizedAlign === 'left') {
      styles.push('display:block', 'margin-left:0', 'margin-right:auto')
    }
  }
  return `${styles.join('; ')};`
}
const buildMarkdownImageToken = (src: string, alt: string, width?: string, align?: ImageAlignMode, wrap?: boolean) => {
  const suffix = buildImageMetaSuffix(width, align, wrap)
  return `![${alt || 'image'}](${src})${suffix}`
}
const MATH_SYMBOLS: Record<string, string> = {
  alpha: 'α',
  beta: 'β',
  gamma: 'γ',
  delta: 'δ',
  Delta: 'Δ',
  theta: 'θ',
  phi: 'φ',
  lambda: 'λ',
  mu: 'μ',
  nu: 'ν',
  omega: 'ω',
  Omega: 'Ω',
  pi: 'π',
  epsilon: 'ε',
  varepsilon: 'ε',
  tau: 'τ',
  sigma: 'σ',
  Sigma: 'Σ',
  eta: 'η',
  pm: '±',
  cdot: '·',
  times: '×',
  div: '÷',
  neq: '≠',
  le: '≤',
  ge: '≥',
  approx: '≈',
  to: '→',
  rightarrow: '→',
  leftarrow: '←',
  cup: '∪',
  cap: '∩',
  sim: '∼',
  infty: '∞',
  partial: '∂',
  int: '∫',
  sum: '∑',
  prod: '∏',
  degree: '°',
  mid: '|',
  ldots: '…',
  prime: '′',
  forall: '∀',
  exists: '∃',
  neg: '¬',
  land: '∧',
  lor: '∨',
  in: '∈',
  notin: '∉',
  subset: '⊂',
  subseteq: '⊆',
  emptyset: '∅',
  bot: '⊥',
  top: '⊤',
  Rightarrow: '⇒',
  Leftarrow: '⇐',
  Leftrightarrow: '⇔',
  implies: '⇒',
  iff: '⇔',
  lbrace: '{',
  rbrace: '}',
  langle: '⟨',
  rangle: '⟩'
}
const MATH_OPERATORS = new Set(['sin', 'cos', 'tan', 'cot', 'ln', 'log', 'det', 'max', 'min'])
const renderMathText = (raw: string) => {
  return escapeHtml(raw)
}
const wrapMathScript = (base: string, subscript: string, superscript: string) => {
  if (!subscript && !superscript) return base
  let out = `<span class="math-scripted">${base}`
  if (subscript) out += `<sub>${subscript}</sub>`
  if (superscript) out += `<sup>${superscript}</sup>`
  out += '</span>'
  return out
}
type MathParserState = {
  source: string
  index: number
}
const MATH_SPACING_COMMANDS: Record<string, string> = {
  ',': ' ',
  ':': ' ',
  ';': ' ',
  quad: '  ',
  qquad: '    '
}
const peekMath = (state: MathParserState) => state.source[state.index] || ''
const consumeMath = (state: MathParserState) => {
  const char = state.source[state.index] || ''
  state.index += 1
  return char
}
const parseMathExpression = (state: MathParserState, stopChar = ''): string => {
  const parts: string[] = []
  while (state.index < state.source.length) {
    const char = peekMath(state)
    if (stopChar && char === stopChar) break
    if (char === ' ') {
      consumeMath(state)
      parts.push(' ')
      continue
    }
    parts.push(parseMathPrimary(state))
  }
  return parts.join('')
}
const parseMathGroup = (state: MathParserState) => {
  if (peekMath(state) !== '{') return ''
  consumeMath(state)
  const inner = parseMathExpression(state, '}')
  if (peekMath(state) === '}') consumeMath(state)
  return inner
}
const parseMathCommandName = (state: MathParserState) => {
  let name = ''
  while (state.index < state.source.length) {
    const char = peekMath(state)
    if (!/[a-zA-Z]/.test(char)) break
    name += consumeMath(state)
  }
  return name
}
const parseMathRawGroup = (state: MathParserState) => {
  if (peekMath(state) !== '{') return ''
  consumeMath(state)
  let depth = 1
  let out = ''
  while (state.index < state.source.length && depth > 0) {
    const char = consumeMath(state)
    if (char === '{') {
      depth += 1
      out += char
      continue
    }
    if (char === '}') {
      depth -= 1
      if (depth === 0) break
      out += char
      continue
    }
    out += char
  }
  return out
}
const parseMathScript = (state: MathParserState) => {
  const char = peekMath(state)
  if (!char) return ''
  if (char === '{') return parseMathGroup(state)
  return parseMathPrimary(state)
}
const parseMathAtom = (state: MathParserState): string => {
  const char = peekMath(state)
  if (!char) return ''
  if (char === '{') {
    return parseMathGroup(state)
  }
  if (char === '\\') {
    consumeMath(state)
    const escapedChar = peekMath(state)
    if (escapedChar && !/[a-zA-Z]/.test(escapedChar)) {
      consumeMath(state)
      if (escapedChar === '\\') return '<br>'
      return MATH_SPACING_COMMANDS[escapedChar] || renderMathText(escapedChar)
    }
    const command = parseMathCommandName(state)
    if (!command) return renderMathText('\\')
    if (command === 'text') {
      return `<span class="math-text">${escapeHtml(parseMathRawGroup(state))}</span>`
    }
    if (command === 'left' || command === 'right') {
      return ''
    }
    if (MATH_SPACING_COMMANDS[command]) {
      return MATH_SPACING_COMMANDS[command]
    }
    if (command === 'frac') {
      const numerator = parseMathScript(state)
      const denominator = parseMathScript(state)
      return `<span class="math-frac"><span class="math-frac-num">${numerator}</span><span class="math-frac-den">${denominator}</span></span>`
    }
    if (command === 'sqrt') {
      const radicand = parseMathScript(state)
      return `<span class="math-radical"><span class="math-radical-sign">√</span><span class="math-radical-body">${radicand}</span></span>`
    }
    if (MATH_OPERATORS.has(command)) {
      return `<span class="math-op">${escapeHtml(command)}</span>`
    }
    if (MATH_SYMBOLS[command]) {
      return `<span class="math-symbol">${MATH_SYMBOLS[command]}</span>`
    }
    return `<span class="math-id">${escapeHtml(command)}</span>`
  }
  if (/\d/.test(char)) {
    let value = ''
    while (/\d/.test(peekMath(state))) {
      value += consumeMath(state)
    }
    return `<span class="math-num">${escapeHtml(value)}</span>`
  }
  if (/[a-zA-Z]/.test(char)) {
    let value = ''
    while (/[a-zA-Z]/.test(peekMath(state))) {
      value += consumeMath(state)
    }
    return `<span class="math-id">${escapeHtml(value)}</span>`
  }
  consumeMath(state)
  return `<span class="math-symbol">${renderMathText(char)}</span>`
}
const parseMathPrimary = (state: MathParserState): string => {
  let base = parseMathAtom(state)
  if (!base) return ''
  let subscript = ''
  let superscript = ''
  while (state.index < state.source.length) {
    const char = peekMath(state)
    if (char === '_') {
      consumeMath(state)
      subscript = parseMathScript(state)
      continue
    }
    if (char === '^') {
      consumeMath(state)
      superscript = parseMathScript(state)
      continue
    }
    break
  }
  return wrapMathScript(base, subscript, superscript)
}
export const renderMath = (raw: string, display = false) => {
  const source = String(raw || '')
    .replace(/\s+/g, ' ')
    .trim()
  if (!source) return ''
  const state: MathParserState = { source, index: 0 }
  const body = parseMathExpression(state)
  const className = display ? 'md-math md-math-block' : 'md-math'
  return `<span class="${className}">${body}</span>`
}
const injectMathTokens = (raw: string) => {
  const tokens: string[] = []
  const next = String(raw || '').replace(
    /\$\$([\s\S]+?)\$\$|\$([^$\n]+?)\$/g,
    (_, blockExpr: string, inlineExpr: string) => {
      const token = `@@MATH${tokens.length}@@`
      const expr = blockExpr ?? inlineExpr ?? ''
      tokens.push(renderMath(expr, blockExpr !== undefined))
      return token
    }
  )
  return { source: next, tokens }
}
const restoreMathTokens = (html: string, tokens: string[]) => {
  let out = html
  for (let index = 0; index < tokens.length; index += 1) {
    out = out.replace(`@@MATH${index}@@`, tokens[index] || '')
  }
  return out
}
export const getEmojiDataUri = (rawName: string) => {
  if (!emojiCollectionLoaded) return ''
  const name = String(rawName || '')
    .toLowerCase()
    .trim()
  if (!name) return ''
  if (emojiDataUriCache.has(name)) return emojiDataUriCache.get(name) || ''
  const resolvedName = emojiIcons[name] ? name : emojiAliases[name]?.parent || ''
  const icon = resolvedName ? emojiIcons[resolvedName] : null
  if (!icon?.body) {
    emojiDataUriCache.set(name, '')
    return ''
  }
  const left = Number(icon.left || 0)
  const top = Number(icon.top || 0)
  const width = Number(icon.width || emojiDefaultWidth)
  const height = Number(icon.height || emojiDefaultHeight)
  const svg = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="${left} ${top} ${width} ${height}">${icon.body}</svg>`
  const uri = `data:image/svg+xml;charset=UTF-8,${encodeURIComponent(svg)}`
  emojiDataUriCache.set(name, uri)
  return uri
}
const renderInline = (raw: string) => {
  const mathPrepared = injectMathTokens(raw)
  let out = escapeHtml(mathPrepared.source)
  out = out.replace(
    /&lt;span\s+style=&quot;([\s\S]*?)&quot;&gt;([\s\S]*?)&lt;\/span&gt;/gi,
    (_, style: string, inner: string) => {
      const safeStyle = sanitizeSpanColorStyle(style)
      if (!safeStyle) return inner
      return `<span style="${safeStyle}">${inner}</span>`
    }
  )
  out = out.replace(
    /!\[([^\]]*)\]\(((?:https?:\/\/|\/)[^\s)]+)\)(\{[^}]+\})?/g,
    (_match: string, alt: string, src: string, meta?: string) => {
      const parsed = parseMarkdownImageOptions(meta)
      const safeSrc = escapeAttrValue(src)
      const safeAlt = escapeAttrValue(alt)
      return `<img src="${safeSrc}" alt="${safeAlt}" loading="lazy" decoding="async" data-display-width="${parsed.width}" data-display-align="${parsed.align}" data-display-wrap="${parsed.wrap ? '1' : '0'}" style="${buildImageStyleAttr(parsed.width, parsed.align, parsed.wrap)}">`
    }
  )
  out = out.replace(
    /\[([^\]]+)\]\(((?:https?:\/\/|\/)[^)\s]+)\)/g,
    '<a href="$2" target="_blank" rel="noopener noreferrer">$1</a>'
  )
  out = out.replace(/`([^`]+)`/g, '<code>$1</code>')
  out = out.replace(/\*\*([^*]+)\*\*/g, '<strong>$1</strong>')
  out = out.replace(/~~([^~]+)~~/g, '<s>$1</s>')
  out = out.replace(/\*([^*]+)\*/g, '<em>$1</em>')
  out = out.replace(/&lt;u&gt;([\s\S]+?)&lt;\/u&gt;/g, '<u>$1</u>')
  out = out.replace(/:emojione:([a-z0-9-]+):/gi, (_, name: string) => {
    const src = getEmojiDataUri(name)
    if (!src) return `:emojione:${name}:`
    return `<img class="emoji-inline" src="${src}" alt="${escapeAttr(name)}" loading="lazy" decoding="async">`
  })
  return restoreMathTokens(out, mathPrepared.tokens)
}
export const highlightCodeSnippet = (code: string, language: string): HighlightedCodeSnippet => {
  const rawLanguage = escapeAttr(language.trim().toLowerCase()) || 'text'
  const alias = CODE_LANGUAGE_ALIASES[rawLanguage]
  const highlightLanguage = alias?.highlight || rawLanguage
  const displayLanguage = alias?.display || rawLanguage
  if (highlightLanguage && hljs.getLanguage(highlightLanguage)) {
    try {
      const highlighted = hljs.highlight(code, {
        language: highlightLanguage,
        ignoreIllegals: true
      }).value
      return {
        highlightLanguage,
        displayLanguage,
        highlightedHtml: `<pre><code class="hljs language-${highlightLanguage}">${highlighted}</code></pre>`
      }
    } catch {
      // fallback below
    }
  }
  const escaped = escapeHtml(code)
  return {
    highlightLanguage,
    displayLanguage,
    highlightedHtml: `<pre><code class="hljs language-${highlightLanguage}">${escaped}</code></pre>`
  }
}
export const renderCodeBlock = (code: string, language: string) => {
  const highlighted = highlightCodeSnippet(code, language)
  const safeLanguage = escapeAttr(highlighted.displayLanguage)
  const label = escapeHtml(highlighted.displayLanguage)
  return `<div class="md-code-block" data-language="${safeLanguage}"><div class="md-code-block-label">${label}</div>${highlighted.highlightedHtml}</div>`
}
const parseMarkdownTableRow = (raw: string) => {
  const line = String(raw || '').trim()
  if (!line || !line.includes('|')) return null
  let normalized = line
  if (normalized.startsWith('|')) normalized = normalized.slice(1)
  if (normalized.endsWith('|')) normalized = normalized.slice(0, -1)
  return normalized.split('|').map(cell => cell.trim())
}
const isMarkdownTableSeparatorCell = (raw: string) => {
  const cell = String(raw || '').replace(/\s+/g, '')
  return /^:?-{3,}:?$/.test(cell)
}
const toMarkdownTableColumnKey = (raw: string, index: number, used: Set<string>) => {
  const base =
    String(raw || '')
      .trim()
      .toLowerCase()
      .replace(/<[^>]+>/g, '')
      .replace(/&[a-z0-9#]+;/gi, '')
      .replace(/[^\p{L}\p{N}]+/gu, '_')
      .replace(/^_+|_+$/g, '') || `col_${index + 1}`
  let key = base
  let suffix = 2
  while (used.has(key)) {
    key = `${base}_${suffix}`
    suffix += 1
  }
  used.add(key)
  return key
}
const renderMarkdownTableHtml = (block: MarkdownViewerTableBlock) => {
  const head = block.columns.map(column => `<th>${escapeHtml(column.label)}</th>`).join('')
  const body =
    block.rows.length > 0
      ? block.rows
          .map(row => {
            const cells = block.columns.map(column => `<td>${row[column.key] || '&mdash;'}</td>`).join('')
            return `<tr>${cells}</tr>`
          })
          .join('')
      : `<tr><td colspan="${Math.max(block.columns.length, 1)}">&mdash;</td></tr>`
  return `<div class="md-table-html-wrap"><table class="md-table-html"><thead><tr>${head}</tr></thead><tbody>${body}</tbody></table></div>`
}
export const renderMarkdownToBlocks = (source: string): MarkdownViewerBlock[] => {
  const lines = String(source || '')
    .replace(/\r/g, '')
    .split('\n')
  const htmlBuffer: string[] = []
  const blocks: MarkdownViewerBlock[] = []
  let inCode = false
  let codeLanguage = ''
  const codeBuffer: string[] = []
  let inMathBlock = false
  const mathBuffer: string[] = []
  let inUl = false
  let inOl = false
  let inQuote = false
  let alignMode: AlignMode = ''
  let htmlBlockIndex = 0
  let tableBlockIndex = 0
  const pushHtml = (html: string) => {
    htmlBuffer.push(html)
  }
  const flushHtmlBuffer = () => {
    if (htmlBuffer.length === 0) return
    const html = htmlBuffer.join('\n').trim()
    htmlBuffer.length = 0
    if (!html) return
    blocks.push({
      type: 'html',
      key: `html-${htmlBlockIndex}`,
      html
    })
    htmlBlockIndex += 1
  }
  const closeLists = () => {
    if (inUl) {
      pushHtml('</ul>')
      inUl = false
    }
    if (inOl) {
      pushHtml('</ol>')
      inOl = false
    }
  }
  const closeQuote = () => {
    if (inQuote) {
      pushHtml('</blockquote>')
      inQuote = false
    }
  }
  const closeAlign = () => {
    if (!alignMode) return
    closeLists()
    closeQuote()
    pushHtml('</div>')
    alignMode = ''
  }
  const closeFlowBlocks = () => {
    closeLists()
    closeQuote()
  }
  const flushMathBlock = () => {
    const body = mathBuffer.join('\n')
    mathBuffer.length = 0
    inMathBlock = false
    if (!body.trim()) return
    pushHtml(`<p>${renderMath(body, true)}</p>`)
  }
  for (let lineIndex = 0; lineIndex < lines.length; lineIndex += 1) {
    const line = lines[lineIndex] || ''
    const trimmed = line.trim()
    if (inCode) {
      if (trimmed.startsWith('```')) {
        pushHtml(renderCodeBlock(codeBuffer.join('\n'), codeLanguage))
        inCode = false
        codeLanguage = ''
        codeBuffer.length = 0
      } else {
        codeBuffer.push(line)
      }
      continue
    }
    if (inMathBlock) {
      if (trimmed === '$$') {
        flushMathBlock()
      } else {
        mathBuffer.push(line)
      }
      continue
    }
    const fenceOpen = trimmed.match(/^```([a-zA-Z0-9_+.-]*)\s*$/)
    if (fenceOpen) {
      closeFlowBlocks()
      inCode = true
      codeLanguage = String(fenceOpen[1] || '').toLowerCase()
      continue
    }
    if (trimmed === '$$') {
      closeFlowBlocks()
      inMathBlock = true
      mathBuffer.length = 0
      continue
    }
    const singleLineMathBlock = trimmed.match(/^\$\$(.+)\$\$$/)
    if (singleLineMathBlock) {
      closeFlowBlocks()
      pushHtml(`<p>${renderMath(singleLineMathBlock[1] || '', true)}</p>`)
      continue
    }
    const alignDirectiveOpen = trimmed.match(/^:::\s*align-(left|center|right|justify)\s*$/i)
    if (alignDirectiveOpen) {
      closeFlowBlocks()
      closeAlign()
      alignMode = (alignDirectiveOpen[1] as AlignMode) || 'left'
      pushHtml(`<div class="md-align md-align-${alignMode}">`)
      continue
    }
    const alignOpen = trimmed.match(/^<div\s+align="(left|center|right|justify)">$/i)
    if (alignOpen) {
      closeFlowBlocks()
      closeAlign()
      alignMode = (alignOpen[1] as AlignMode) || 'left'
      pushHtml(`<div class="md-align md-align-${alignMode}">`)
      continue
    }
    if (trimmed === ':::' && alignMode) {
      closeFlowBlocks()
      pushHtml('</div>')
      alignMode = ''
      continue
    }
    if (trimmed.toLowerCase() === '</div>' && alignMode) {
      closeFlowBlocks()
      pushHtml('</div>')
      alignMode = ''
      continue
    }
    if (!trimmed) {
      closeFlowBlocks()
      continue
    }
    const headingMatch = trimmed.match(/^(#{1,6})\s+(.*)$/)
    if (headingMatch) {
      closeFlowBlocks()
      const level = (headingMatch[1] || '#').length
      pushHtml(`<h${level}>${renderInline(headingMatch[2] || '')}</h${level}>`)
      continue
    }
    const quoteMatch = trimmed.match(/^>\s?(.*)$/)
    if (quoteMatch) {
      closeLists()
      if (!inQuote) {
        pushHtml('<blockquote>')
        inQuote = true
      }
      pushHtml(`<p>${renderInline(quoteMatch[1] || '')}</p>`)
      continue
    }
    const headerCells = parseMarkdownTableRow(line)
    const separatorCells = parseMarkdownTableRow(lines[lineIndex + 1] || '')
    const isTableStart = Boolean(
      headerCells &&
      separatorCells &&
      headerCells.length > 0 &&
      headerCells.length === separatorCells.length &&
      separatorCells.every(isMarkdownTableSeparatorCell)
    )
    if (isTableStart && headerCells && separatorCells) {
      closeFlowBlocks()
      flushHtmlBuffer()
      const rows: string[][] = []
      let cursor = lineIndex + 2
      while (cursor < lines.length) {
        const rowCells = parseMarkdownTableRow(lines[cursor] || '')
        if (!rowCells || rowCells.length !== headerCells.length) break
        rows.push(rowCells)
        cursor += 1
      }
      const usedKeys = new Set<string>()
      const columns = headerCells.map((label, index) => ({
        key: toMarkdownTableColumnKey(label, index, usedKeys),
        label: label || `Колонка ${index + 1}`
      }))
      blocks.push({
        type: 'table',
        key: `table-${tableBlockIndex}`,
        columns,
        rows: rows.map((cells, rowIndex) =>
          columns.reduce<Record<string, string>>((acc, column, columnIndex) => {
            const rawCell = cells[columnIndex] || ''
            acc.id = `row_${tableBlockIndex}_${rowIndex + 1}`
            acc[column.key] = rawCell ? renderInline(rawCell) : '&mdash;'
            return acc
          }, {})
        )
      })
      tableBlockIndex += 1
      lineIndex = cursor - 1
      continue
    }
    const ulMatch = trimmed.match(/^[-*]\s+(.*)$/)
    if (ulMatch) {
      closeQuote()
      if (!inUl) {
        if (inOl) {
          pushHtml('</ol>')
          inOl = false
        }
        pushHtml('<ul>')
        inUl = true
      }
      pushHtml(`<li>${renderInline(ulMatch[1] || '')}</li>`)
      continue
    }
    const olMatch = trimmed.match(/^\d+\.\s+(.*)$/)
    if (olMatch) {
      closeQuote()
      if (!inOl) {
        if (inUl) {
          pushHtml('</ul>')
          inUl = false
        }
        pushHtml('<ol>')
        inOl = true
      }
      pushHtml(`<li>${renderInline(olMatch[1] || '')}</li>`)
      continue
    }
    closeFlowBlocks()
    pushHtml(`<p>${renderInline(trimmed)}</p>`)
  }
  if (inCode) {
    pushHtml(renderCodeBlock(codeBuffer.join('\n'), codeLanguage))
  }
  if (inMathBlock) {
    flushMathBlock()
  }
  closeFlowBlocks()
  closeAlign()
  flushHtmlBuffer()
  return blocks
}
export const renderMarkdownToHtml = (source: string) => {
  return renderMarkdownToBlocks(source)
    .map(block => (block.type === 'html' ? block.html : renderMarkdownTableHtml(block)))
    .join('\n')
}
const inlineNodeToMarkdown = (node: Node): string => {
  if (node.nodeType === Node.TEXT_NODE) {
    return node.textContent || ''
  }
  if (node.nodeType !== Node.ELEMENT_NODE) return ''
  const el = node as HTMLElement
  const tag = el.tagName.toLowerCase()
  const text = inlineChildrenToMarkdown(Array.from(el.childNodes))
  if (tag === 'br') return '\n'
  if (tag === 'strong' || tag === 'b') return `**${text}**`
  if (tag === 'em' || tag === 'i') return `*${text}*`
  if (tag === 'u') return `<u>${text}</u>`
  if (tag === 's' || tag === 'strike' || tag === 'del') return `~~${text}~~`
  if (tag === 'code') return `\`${(el.textContent || '').replace(/\n/g, ' ')}\``
  if (tag === 'font') {
    const color = normalizeCssColor(String(el.getAttribute('color') || ''))
    return color ? `<span style="color:${color}">${text}</span>` : text
  }
  if (tag === 'span') {
    const rawStyle = String(el.getAttribute('style') || '')
    const style = rawStyle.toLowerCase()
    let wrapped = text
    const safeStyle = sanitizeSpanColorStyle(rawStyle)
    if (/font-weight\s*:\s*(bold|[5-9]00)/.test(style)) {
      wrapped = `**${wrapped}**`
    }
    if (style.includes('font-style: italic')) {
      wrapped = `*${wrapped}*`
    }
    if (style.includes('line-through')) {
      wrapped = `~~${wrapped}~~`
    }
    if (style.includes('text-decoration') && style.includes('underline')) {
      wrapped = `<u>${wrapped}</u>`
    }
    if (safeStyle && /(^|;\s*)color\s*:/.test(safeStyle)) {
      wrapped = `<span style="${safeStyle}">${wrapped}</span>`
    }
    return wrapped
  }
  if (tag === 'a') {
    const href = String(el.getAttribute('href') || '').trim()
    const label = (text || href || 'ссылка').trim()
    return href ? `[${label}](${href})` : label
  }
  if (tag === 'img') {
    const alt = String(el.getAttribute('alt') || '').trim()
    const src = String(el.getAttribute('src') || '').trim()
    if (el.classList.contains('emoji-inline') && alt) {
      return `:emojione:${alt.toLowerCase()}:`
    }
    if (!src) return ''
    const explicitWidth = normalizeImageDisplayWidth(el.getAttribute('data-display-width') || el.style.width)
    const rawStyle = String(el.getAttribute('style') || '').toLowerCase()
    const styleFloat = rawStyle.match(/(?:^|;)\s*float\s*:\s*(left|right)\s*(?:;|$)/)?.[1] || ''
    const styleAlign = (() => {
      if (/margin-left\s*:\s*auto/.test(rawStyle) && /margin-right\s*:\s*auto/.test(rawStyle)) return 'center'
      if (/margin-left\s*:\s*auto/.test(rawStyle)) return 'right'
      if (/margin-right\s*:\s*auto/.test(rawStyle)) return 'left'
      return ''
    })()
    const explicitAlign = normalizeImageDisplayAlign(
      el.getAttribute('data-display-align') || el.getAttribute('align') || styleFloat || styleAlign
    )
    const explicitWrap = normalizeImageTextWrap(el.getAttribute('data-display-wrap')) || Boolean(styleFloat)
    return buildMarkdownImageToken(src, alt || 'image', explicitWidth, explicitAlign, explicitWrap)
  }
  return text
}
const inlineChildrenToMarkdown = (nodes: Node[]) => {
  return nodes.map(inlineNodeToMarkdown).join('')
}
const detectAlign = (el: HTMLElement): AlignMode => {
  const alignAttr = String(el.getAttribute('align') || '')
    .toLowerCase()
    .trim()
  if (['left', 'center', 'right', 'justify'].includes(alignAttr)) return alignAttr as AlignMode
  const className = String(el.className || '')
  if (className.includes('md-align-left')) return 'left'
  if (className.includes('md-align-center')) return 'center'
  if (className.includes('md-align-right')) return 'right'
  if (className.includes('md-align-justify')) return 'justify'
  const styleAlign = String(el.style?.textAlign || '')
    .toLowerCase()
    .trim()
  if (['left', 'center', 'right', 'justify'].includes(styleAlign)) return styleAlign as AlignMode
  return ''
}
const unwrapOuterAlignDirective = (raw: string, mode: Exclude<AlignMode, ''>) => {
  let text = String(raw || '').trim()
  if (!text) return ''
  const pattern = new RegExp(`^:::align-${mode}\\s*\\n([\\s\\S]*?)\\n:::\\s*$`, 'i')
  while (true) {
    const match = text.match(pattern)
    if (!match) break
    text = String(match[1] || '').trim()
  }
  return text
}
const wrapBlocksWithAlign = (lines: string[], align: Exclude<AlignMode, ''>) => {
  if (lines.length === 0) return []
  const body = lines
    .map(line => unwrapOuterAlignDirective(line, align))
    .join('\n\n')
    .trim()
  if (!body) return []
  return [`:::align-${align}\n${body}\n:::`]
}
const blockNodeToMarkdown = (node: Node): string[] => {
  if (node.nodeType === Node.TEXT_NODE) {
    const text = String(node.textContent || '').trim()
    return text ? [text] : []
  }
  if (node.nodeType !== Node.ELEMENT_NODE) return []
  const el = node as HTMLElement
  const tag = el.tagName.toLowerCase()
  const align = detectAlign(el)
  if (tag === 'h1' || tag === 'h2' || tag === 'h3' || tag === 'h4' || tag === 'h5' || tag === 'h6') {
    const level = Number(tag[1] || '1')
    const text = inlineChildrenToMarkdown(Array.from(el.childNodes)).trim()
    if (!text) return []
    const out = [`${'#'.repeat(level)} ${text}`]
    return align ? wrapBlocksWithAlign(out, align) : out
  }
  if (tag === 'p') {
    const text = inlineChildrenToMarkdown(Array.from(el.childNodes)).trim()
    if (!text) return []
    const out = [text]
    return align ? wrapBlocksWithAlign(out, align) : out
  }
  if (tag === 'blockquote') {
    const parts = Array.from(el.childNodes).flatMap(blockNodeToMarkdown)
    if (parts.length === 0) return []
    const text = parts
      .join('\n\n')
      .split('\n')
      .map(line => `> ${line}`.trimEnd())
      .join('\n')
    const out = [text]
    return align ? wrapBlocksWithAlign(out, align) : out
  }
  if (tag === 'ul' || tag === 'ol') {
    const lines = Array.from(el.children)
      .filter(item => item.tagName.toLowerCase() === 'li')
      .map((li, index) => {
        const body = inlineChildrenToMarkdown(Array.from(li.childNodes)).trim() || 'пункт'
        return tag === 'ol' ? `${index + 1}. ${body}` : `- ${body}`
      })
    if (lines.length === 0) return []
    const out = [lines.join('\n')]
    return align ? wrapBlocksWithAlign(out, align) : out
  }
  if (tag === 'pre') {
    const codeEl = el.querySelector('code')
    const className = String(codeEl?.getAttribute('class') || '')
    const language = (className.match(/language-([a-z0-9_+.-]+)/i)?.[1] || '').toLowerCase() || 'text'
    const raw = String(codeEl?.textContent || el.textContent || '')
    const out = [`\`\`\`${language}\n${raw}\n\`\`\``]
    return align ? wrapBlocksWithAlign(out, align) : out
  }
  if (tag === 'div') {
    const inner = Array.from(el.childNodes).flatMap(blockNodeToMarkdown)
    if (align) {
      if (inner.length === 0) return []
      return wrapBlocksWithAlign(inner, align)
    }
    return inner
  }
  const text = inlineChildrenToMarkdown(Array.from(el.childNodes)).trim()
  if (text) {
    const out = [text]
    return align ? wrapBlocksWithAlign(out, align) : out
  }
  const inner = Array.from(el.childNodes).flatMap(blockNodeToMarkdown)
  return align ? wrapBlocksWithAlign(inner, align) : inner
}
export const renderEditableHtmlToMarkdown = (html: string) => {
  if (!import.meta.client) return String(html || '')
  const holder = document.createElement('div')
  holder.innerHTML = String(html || '')
  const blocks = Array.from(holder.childNodes).flatMap(blockNodeToMarkdown)
  return blocks
    .join('\n\n')
    .replace(/\n{3,}/g, '\n\n')
    .trim()
}
export const buildAnchorHtml = (href: string, label: string) => {
  const safeHref = escapeAttrValue(href)
  const safeLabel = escapeHtml(label)
  return `<a href="${safeHref}" target="_blank" rel="noopener noreferrer">${safeLabel}</a>`
}
export const buildImageHtml = (src: string, alt: string, width?: string, align?: ImageAlignMode, wrap?: boolean) => {
  const safeSrc = escapeAttrValue(src)
  const safeAlt = escapeAttrValue(alt)
  const normalizedWidth = normalizeImageDisplayWidth(width)
  const normalizedAlign = normalizeImageDisplayAlign(align)
  const normalizedWrap = normalizeImageTextWrap(wrap) && (normalizedAlign === 'left' || normalizedAlign === 'right')
  return `<img src="${safeSrc}" alt="${safeAlt}" loading="lazy" decoding="async" data-display-width="${normalizedWidth}" data-display-align="${normalizedAlign}" data-display-wrap="${normalizedWrap ? '1' : '0'}" style="${buildImageStyleAttr(normalizedWidth, normalizedAlign, normalizedWrap)}">`
}
