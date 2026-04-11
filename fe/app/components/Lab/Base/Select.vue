<template>
  <div ref="rootRef" class="relative" :class="disabled ? 'opacity-60' : ''">
    <button
      :id="id"
      ref="triggerRef"
      type="button"
      :disabled="disabled"
      :aria-label="ariaLabel"
      :aria-expanded="isOpen ? 'true' : 'false'"
      aria-haspopup="listbox"
      :class="selectClassList"
      :style="triggerStyle"
      v-bind="triggerAttrs"
      @click="toggleOpen"
      @keydown="onTriggerKeydown">
      <span class="flex min-w-0 items-center gap-2 text-left">
        <span
          v-if="selectedOption?.swatchColor"
          aria-hidden="true"
          class="h-2.5 w-2.5 shrink-0 rounded-full border border-white/20"
          :style="getSwatchStyle(selectedOption.swatchColor)"></span>
        <span class="block min-w-0 truncate">
          {{ selectedLabel }}
        </span>
      </span>
      <Icon
        name="ic:round-expand-more"
        class="pointer-events-none absolute right-3 top-1/2 -translate-y-1/2 text-base lab-text-soft transition-transform duration-200"
        :class="isOpen ? 'rotate-180' : 'rotate-0'" />
    </button>
    <input v-if="name" type="hidden" :name="name" :value="resolvedValue" />
    <Teleport to="body">
      <div
        v-show="isOpen"
        ref="dropdownRef"
        class="fixed z-10000 overflow-hidden border"
        :style="resolvedDropdownStyle">
        <div
          v-if="enabledOptions.length || normalizedOptions.length"
          class="max-h-72 overflow-auto"
          role="listbox"
          :aria-labelledby="id">
          <button
            v-for="(option, index) in normalizedOptions"
            :key="option.key"
            type="button"
            class="block w-full px-3.5 py-2.5 text-left text-sm transition-colors disabled:pointer-events-none disabled:opacity-50"
            :class="
              index === highlightedIndex
                ? 'lab-surface-subtle lab-text-primary'
                : option.value === resolvedValue
                  ? ' lab-text-accent'
                  : 'lab-text-secondary hover:bg-zinc-800/10 hover:text-zinc-100'
            "
            :disabled="option.disabled"
            role="option"
            :aria-selected="option.value === resolvedValue ? 'true' : 'false'"
            @click="selectOption(option.value)"
            @mouseenter="setHighlightedIndex(index)">
            <span class="flex min-w-0 items-center gap-2">
              <span
                v-if="option.swatchColor"
                aria-hidden="true"
                class="h-2.5 w-2.5 shrink-0 rounded-full border border-white/20"
                :style="getSwatchStyle(option.swatchColor)"></span>
              <span class="block min-w-0 truncate">{{ option.label }}</span>
            </span>
          </button>
        </div>
        <div v-else class="px-3.5 py-2.5 text-sm lab-text-soft">{{ emptyText }}</div>
      </div>
    </Teleport>
  </div>
</template>
<script setup lang="ts">
  import { twMerge } from 'tailwind-merge'
  defineOptions({ inheritAttrs: false })
  const props = withDefaults(
    defineProps<{
      modelValue?: string | number | null | undefined
      options?: SelectOptionInput[]
      id?: string
      name?: string
      disabled?: boolean
      ariaLabel?: string
      invalid?: boolean
      selectClass?: string
      placeholder?: string
      emptyText?: string
      matchTriggerWidth?: boolean
      offset?: number
      crossAxisOffset?: number
    }>(),
    {
      modelValue: '',
      options: () => [],
      id: '',
      name: '',
      disabled: false,
      ariaLabel: '',
      invalid: false,
      selectClass: '',
      placeholder: 'Выберите значение',
      emptyText: 'Нет доступных вариантов',
      matchTriggerWidth: true,
      offset: 6,
      crossAxisOffset: 10
    }
  )
  const emit = defineEmits<{
    (e: 'update:modelValue', value: string): void
    (e: 'change', event: Event): void
    (e: 'open'): void
    (e: 'close'): void
  }>()
  const attrs = useAttrs()
  const rootRef = ref<HTMLElement | null>(null)
  const triggerRef = ref<HTMLElement | null>(null)
  const dropdownRef = ref<HTMLElement | null>(null)
  const isOpen = ref(false)
  const highlightedIndex = ref(-1)
  const measuredMinWidthPx = ref<number>(0)
  const { panelStyle, schedulePositionUpdate, cancelScheduledUpdate, resetPosition } = useFloatingPanelPosition({
    triggerRef,
    panelRef: dropdownRef,
    side: 'bottom',
    align: 'left',
    offset: computed(() => props.offset),
    crossAxisOffset: computed(() => props.crossAxisOffset),
    viewportPadding: 12,
    matchTriggerWidth: computed(() => props.matchTriggerWidth)
  })
  const resolvedValue = computed(() => {
    if (props.modelValue === null || props.modelValue === undefined) return ''
    return String(props.modelValue)
  })
  const normalizedOptions = computed<NormalizedSelectOption[]>(() =>
    (props.options || []).map(option => {
      const normalizedValue = option.value === null || option.value === undefined ? '' : String(option.value)
      const normalizedOption: NormalizedSelectOption = {
        key: normalizedValue || option.label,
        value: normalizedValue,
        label: option.label,
        disabled: Boolean(option.disabled)
      }
      if (option.swatchColor) {
        normalizedOption.swatchColor = option.swatchColor
      }
      return normalizedOption
    })
  )
  const enabledOptions = computed(() => normalizedOptions.value.filter(option => !option.disabled))
  const selectedOption = computed(() => {
    const value = resolvedValue.value
    return normalizedOptions.value.find(option => option.value === value) || null
  })
  const selectedLabel = computed(() => selectedOption.value?.label || props.placeholder)
  const highlightedOption = computed(() => {
    const index = highlightedIndex.value
    return index >= 0 ? normalizedOptions.value[index] || null : null
  })
  const triggerAttrs = computed(() => {
    const out: Record<string, unknown> = {}
    for (const [key, value] of Object.entries(attrs)) {
      if (key === 'class') continue
      if (key === 'style') continue
      out[key] = value
    }
    return out
  })
  const externalClass = computed(() => String(attrs.class || '').trim())
  const selectClassList = computed(() => {
    const stateClass = props.invalid ? 'lab-control-invalid' : ''
    return twMerge('lab-control relative pr-10 text-left', stateClass, props.selectClass, externalClass.value)
  })
  const triggerStyle = computed(() => {
    if (measuredMinWidthPx.value <= 0) return undefined
    return {
      minWidth: `${measuredMinWidthPx.value}px`
    }
  })
  const resolvedDropdownStyle = computed(() => ({
    ...panelStyle.value,
    background: 'var(--lab-bg-canvas)',
    color: 'var(--lab-text-primary)'
  }))
  const getSwatchStyle = (color?: string) => (color ? { backgroundColor: color } : {})
  const measureMinWidth = () => {
    if (!import.meta.client) return
    const triggerEl = triggerRef.value
    if (!triggerEl) return
    const style = window.getComputedStyle(triggerEl)
    const labels = normalizedOptions.value.map(option => option.label).filter(Boolean)
    labels.push(props.placeholder)
    const font =
      style.font ||
      `${style.fontStyle} ${style.fontVariant} ${style.fontWeight} ${style.fontSize} / ${style.lineHeight} ${style.fontFamily}`
    const canvas = document.createElement('canvas')
    const context = canvas.getContext('2d')
    if (!context) return
    context.font = font
    const widestLabelWidth = labels.reduce((maxWidth, label) => {
      const nextWidth = Math.ceil(context.measureText(label).width)
      return Math.max(maxWidth, nextWidth)
    }, 0)
    const horizontalPadding =
      Number.parseFloat(style.paddingLeft || '0') +
      Number.parseFloat(style.paddingRight || '0') +
      Number.parseFloat(style.borderLeftWidth || '0') +
      Number.parseFloat(style.borderRightWidth || '0')
    const hasSwatch = normalizedOptions.value.some(option => option.swatchColor)
    const swatchReserve = hasSwatch ? 18 : 0
    const caretReserve = 4
    measuredMinWidthPx.value = Math.ceil(widestLabelWidth + horizontalPadding + swatchReserve + caretReserve)
  }
  const scheduleMinWidthMeasurement = () => {
    if (!import.meta.client) return
    requestAnimationFrame(() => {
      measureMinWidth()
    })
  }
  const setHighlightedIndex = (index: number) => {
    if (!normalizedOptions.value[index] || normalizedOptions.value[index].disabled) return
    highlightedIndex.value = index
  }
  const syncHighlightWithValue = () => {
    const selectedIndex = normalizedOptions.value.findIndex(
      option => option.value === resolvedValue.value && !option.disabled
    )
    if (selectedIndex >= 0) {
      highlightedIndex.value = selectedIndex
      return
    }
    const firstEnabledIndex = normalizedOptions.value.findIndex(option => !option.disabled)
    highlightedIndex.value = firstEnabledIndex
  }
  const openMenu = () => {
    if (props.disabled || !normalizedOptions.value.length) return
    if (isOpen.value) return
    isOpen.value = true
    syncHighlightWithValue()
    schedulePositionUpdate()
    emit('open')
  }
  const closeMenu = () => {
    if (!isOpen.value) return
    isOpen.value = false
    highlightedIndex.value = -1
    resetPosition()
    emit('close')
  }
  const toggleOpen = () => {
    if (props.disabled) return
    if (isOpen.value) closeMenu()
    else openMenu()
  }
  const selectOption = (value: string) => {
    if (value === resolvedValue.value) {
      closeMenu()
      return
    }
    emit('update:modelValue', value)
    emit('change', new Event('change'))
    closeMenu()
  }
  const moveHighlight = (direction: 1 | -1) => {
    const options = normalizedOptions.value
    if (!options.length) return
    let nextIndex = highlightedIndex.value
    if (nextIndex < 0) {
      nextIndex = direction === 1 ? 0 : options.length - 1
    } else {
      nextIndex += direction
    }
    while (nextIndex >= 0 && nextIndex < options.length) {
      const option = options[nextIndex]
      if (option && !option.disabled) {
        highlightedIndex.value = nextIndex
        return
      }
      nextIndex += direction
    }
  }
  const selectHighlightedOption = () => {
    const option = highlightedOption.value
    if (!option || option.disabled) return
    selectOption(option.value)
  }
  const onTriggerKeydown = (event: KeyboardEvent) => {
    if (props.disabled) return
    switch (event.key) {
      case 'ArrowDown':
        event.preventDefault()
        if (!isOpen.value) {
          openMenu()
          return
        }
        moveHighlight(1)
        break
      case 'ArrowUp':
        event.preventDefault()
        if (!isOpen.value) {
          openMenu()
          return
        }
        moveHighlight(-1)
        break
      case 'Enter':
      case ' ':
        event.preventDefault()
        if (!isOpen.value) {
          openMenu()
          return
        }
        selectHighlightedOption()
        break
      case 'Escape':
        if (!isOpen.value) return
        event.preventDefault()
        closeMenu()
        break
      case 'Tab':
        closeMenu()
        break
    }
  }
  const onViewportChange = () => {
    scheduleMinWidthMeasurement()
    if (!isOpen.value) return
    schedulePositionUpdate()
  }
  const onDocumentPointerDown = (event: MouseEvent | PointerEvent) => {
    const target = event.target as Node | null
    if (!target) return
    if (rootRef.value?.contains(target)) return
    if (dropdownRef.value?.contains(target)) return
    closeMenu()
  }
  onMounted(() => {
    scheduleMinWidthMeasurement()
    if (typeof document !== 'undefined' && 'fonts' in document) {
      void document.fonts.ready.then(() => {
        scheduleMinWidthMeasurement()
      })
    }
    document.addEventListener('click', onDocumentPointerDown, true)
    window.addEventListener('resize', onViewportChange, { passive: true })
    window.addEventListener('scroll', onViewportChange, { passive: true, capture: true })
  })
  onBeforeUnmount(() => {
    document.removeEventListener('click', onDocumentPointerDown, true)
    window.removeEventListener('resize', onViewportChange)
    window.removeEventListener('scroll', onViewportChange, true)
    cancelScheduledUpdate()
  })
  watch(
    () => props.disabled,
    next => {
      if (next) closeMenu()
    }
  )
  watch(
    () => props.options,
    () => {
      scheduleMinWidthMeasurement()
      if (!normalizedOptions.value.length) {
        closeMenu()
        return
      }
      if (!isOpen.value) return
      syncHighlightWithValue()
      schedulePositionUpdate()
    },
    { deep: true }
  )
  watch(
    () => props.placeholder,
    () => {
      scheduleMinWidthMeasurement()
    }
  )
  watch(
    () => props.modelValue,
    () => {
      if (!isOpen.value) return
      syncHighlightWithValue()
    }
  )
  watch(isOpen, next => {
    if (!import.meta.client || !next) return
    schedulePositionUpdate()
  })
</script>
