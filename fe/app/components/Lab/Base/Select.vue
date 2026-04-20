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
      v-bind="triggerAttrs"
      @click="toggleOpen"
      @keydown="onTriggerKeydown"
    >
      <span class="flex min-w-0 items-center gap-2 text-left">
        <span
          v-if="selectedOption?.swatchColor"
          aria-hidden="true"
          class="h-2.5 w-2.5 shrink-0 rounded-full border border-(--lab-border)"
          :style="getSwatchStyle(selectedOption.swatchColor)"
        />
        <span class="block min-w-0 truncate">
          {{ selectedLabel }}
        </span>
      </span>
      <Icon
        name="ic:round-expand-more"
        class="pointer-events-none absolute top-1/2 right-3 -translate-y-1/2 text-base text-(--lab-text-soft) transition-transform duration-200"
        :class="isOpen ? 'rotate-180' : 'rotate-0'"
      />
    </button>
    <input v-if="name" type="hidden" :name="name" :value="resolvedValue" />
    <Teleport to="body">
      <div
        v-show="isOpen"
        ref="dropdownRef"
        class="fixed z-10000 overflow-hidden border border-(--lab-border)"
        :style="resolvedDropdownStyle"
      >
        <div
          v-if="enabledOptions.length || normalizedOptions.length"
          ref="listboxRef"
          class="max-h-72 overflow-auto"
          role="listbox"
          :aria-labelledby="id"
          @scroll.passive="updateScrollIndicators"
        >
          <button
            v-for="(option, index) in normalizedOptions"
            :key="option.key"
            :ref="getOptionRef(index)"
            type="button"
            class="block w-full px-3.5 py-2.5 text-left text-sm transition-colors disabled:pointer-events-none disabled:opacity-50"
            :class="
              index === highlightedIndex
                ? 'text-(--lab-text-primary)'
                : option.value === resolvedValue
                  ? 'text-(--lab-text-accent)'
                  : 'text-(--lab-text-secondary) hover:bg-zinc-800/10 hover:text-zinc-100'
            "
            :disabled="option.disabled"
            role="option"
            :aria-selected="option.value === resolvedValue ? 'true' : 'false'"
            @click="selectOption(option.value)"
          >
            <span class="flex min-w-0 items-center gap-2">
              <span
                v-if="option.swatchColor"
                aria-hidden="true"
                class="h-2.5 w-2.5 shrink-0 rounded-full border border-white/20"
                :style="getSwatchStyle(option.swatchColor)"
              />
              <span class="block min-w-0 truncate">{{ option.label }}</span>
            </span>
          </button>
        </div>
        <div
          v-if="isOpen && canScrollDown"
          aria-hidden="true"
          class="pointer-events-none absolute inset-x-0 bottom-1 flex justify-center"
        >
          <span class="inline-flex items-center text-xs text-(--lab-text-soft)">
            <Icon name="ic:round-keyboard-arrow-down" class="text-base" />
          </span>
        </div>
      </div>
    </Teleport>
  </div>
</template>
<script setup lang="ts">
import { twMerge } from 'tailwind-merge'
import type { ComponentPublicInstance, VNodeRef } from 'vue'
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
    matchTriggerWidth: false,
    offset: 0,
    crossAxisOffset: 0
  }
)
const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
  (e: 'change', event: Event): void
  (e: 'open' | 'close'): void
}>()
const attrs = useAttrs()
const rootRef = ref<HTMLElement | null>(null)
const triggerRef = ref<HTMLElement | null>(null)
const dropdownRef = ref<HTMLElement | null>(null)
const listboxRef = ref<HTMLElement | null>(null)
const isOpen = ref(false)
const highlightedIndex = ref(-1)
const canScrollDown = ref(false)
const optionRefs = ref<Array<HTMLElement | null>>([])
const { panelStyle, schedulePositionUpdate, resetPosition } = useFloatingPanelPosition({
  triggerRef,
  panelRef: dropdownRef,
  side: 'bottom',
  offset: computed(() => props.offset),
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
  return twMerge('lab-control lab-focus relative pr-10 text-left', stateClass, props.selectClass, externalClass.value)
})
const resolvedDropdownStyle = computed(() => ({
  ...panelStyle.value,
  minWidth: `${Math.round(triggerRef.value?.getBoundingClientRect().width || 0)}px`,
  width: props.matchTriggerWidth ? panelStyle.value.width : 'max-content',
  background: 'var(--lab-bg-canvas)',
  color: 'var(--lab-text-primary)'
}))
const getSwatchStyle = (color?: string) => (color ? { backgroundColor: color } : {})
const setOptionRef = (index: number, el: Element | ComponentPublicInstance | null) => {
  optionRefs.value[index] = el instanceof HTMLElement ? el : null
}
const getOptionRef =
  (index: number): VNodeRef =>
  el => {
    setOptionRef(index, el)
  }
const scrollHighlightedOptionIntoView = () => {
  if (!import.meta.client || !isOpen.value) return
  const listboxEl = listboxRef.value
  const optionEl = optionRefs.value[highlightedIndex.value]
  if (!listboxEl || !optionEl) return
  optionEl.scrollIntoView({ block: 'nearest' })
  updateScrollIndicators()
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
const findSelectableIndex = (startIndex: number, direction: 1 | -1) => {
  const options = normalizedOptions.value
  let nextIndex = startIndex
  while (nextIndex >= 0 && nextIndex < options.length) {
    const option = options[nextIndex]
    if (option && !option.disabled) return nextIndex
    nextIndex += direction
  }
  return -1
}
const selectAdjacentOption = (direction: 1 | -1) => {
  const selectedIndex = normalizedOptions.value.findIndex(
    option => option.value === resolvedValue.value && !option.disabled
  )
  const fallbackIndex = direction === 1 ? 0 : normalizedOptions.value.length - 1
  const nextIndex = findSelectableIndex(selectedIndex >= 0 ? selectedIndex + direction : fallbackIndex, direction)
  if (nextIndex < 0) return false
  highlightedIndex.value = nextIndex
  const nextOption = normalizedOptions.value[nextIndex]
  if (!nextOption) return false
  if (nextOption.value !== resolvedValue.value) {
    emit('update:modelValue', nextOption.value)
    emit('change', new Event('change'))
  }
  return true
}
const moveHighlight = (direction: 1 | -1) => {
  if (!normalizedOptions.value.length) return
  const fallbackIndex = direction === 1 ? 0 : normalizedOptions.value.length - 1
  const nextIndex = findSelectableIndex(
    highlightedIndex.value < 0 ? fallbackIndex : highlightedIndex.value + direction,
    direction
  )
  if (nextIndex >= 0) {
    highlightedIndex.value = nextIndex
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
        selectAdjacentOption(1)
        return
      }
      moveHighlight(1)
      break
    case 'ArrowUp':
      event.preventDefault()
      if (!isOpen.value) {
        selectAdjacentOption(-1)
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
  if (!isOpen.value) return
  schedulePositionUpdate()
  updateScrollIndicators()
}
const updateScrollIndicators = () => {
  const listboxEl = listboxRef.value
  if (!listboxEl) {
    canScrollDown.value = false
    return
  }
  const maxScrollTop = Math.max(0, listboxEl.scrollHeight - listboxEl.clientHeight)
  canScrollDown.value = listboxEl.scrollTop < maxScrollTop - 1
}
const onDocumentPointerDown = (event: MouseEvent | PointerEvent) => {
  const target = event.target as Node | null
  if (!target) return
  if (rootRef.value?.contains(target)) return
  if (dropdownRef.value?.contains(target)) return
  closeMenu()
}
onMounted(() => {
  document.addEventListener('click', onDocumentPointerDown, true)
  window.addEventListener('resize', onViewportChange, { passive: true })
  window.addEventListener('scroll', onViewportChange, { passive: true, capture: true })
})
onBeforeUnmount(() => {
  document.removeEventListener('click', onDocumentPointerDown, true)
  window.removeEventListener('resize', onViewportChange)
  window.removeEventListener('scroll', onViewportChange, true)
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
    optionRefs.value = []
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
  () => props.modelValue,
  () => {
    if (!isOpen.value) return
    syncHighlightWithValue()
  }
)
watch(isOpen, next => {
  if (!next) {
    optionRefs.value = []
    canScrollDown.value = false
    return
  }
  if (!import.meta.client) return
  schedulePositionUpdate()
  requestAnimationFrame(() => {
    scrollHighlightedOptionIntoView()
    updateScrollIndicators()
  })
})
watch(highlightedIndex, () => {
  if (!import.meta.client || !isOpen.value) return
  requestAnimationFrame(() => {
    scrollHighlightedOptionIntoView()
  })
})
</script>
