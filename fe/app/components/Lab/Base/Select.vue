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
          class="h-2.5 w-2.5 shrink-0 rounded-full border border-white/20"
          :style="getSwatchStyle(selectedOption.swatchColor)"
        />
        <span class="block min-w-0 truncate">
          {{ selectedLabel }}
        </span>
      </span>
      <Icon
        name="ic:round-expand-more"
        class="lab-text-soft pointer-events-none absolute top-1/2 right-3 -translate-y-1/2 text-base transition-transform duration-200"
        :class="isOpen ? 'rotate-180' : 'rotate-0'"
      />
    </button>
    <input v-if="name" type="hidden" :name="name" :value="resolvedValue" >
    <Teleport to="body">
      <div
        v-show="isOpen"
        ref="dropdownRef"
        class="fixed z-10000 overflow-hidden border"
        :style="resolvedDropdownStyle"
      >
        <div
          v-if="enabledOptions.length || normalizedOptions.length"
          ref="listboxRef"
          class="max-h-72 overflow-auto"
          role="listbox"
          :aria-labelledby="id"
        >
          <button
            v-for="(option, index) in normalizedOptions"
            :key="option.key"
            :ref="getOptionRef(index)"
            type="button"
            class="block w-full px-3.5 py-2.5 text-left text-sm transition-colors disabled:pointer-events-none disabled:opacity-50"
            :class="
              index === highlightedIndex
                ? 'lab-surface-subtle lab-text-primary'
                : option.value === resolvedValue
                  ? 'lab-text-accent'
                  : 'lab-text-secondary hover:bg-zinc-800/10 hover:text-zinc-100'
            "
            :disabled="option.disabled"
            role="option"
            :aria-selected="option.value === resolvedValue ? 'true' : 'false'"
            @click="selectOption(option.value)"
            @mouseenter="setHighlightedIndex(index)"
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
        <div v-else class="lab-text-soft px-3.5 py-2.5 text-sm">{{ emptyText }}</div>
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
const optionRefs = ref<Array<HTMLElement | null>>([])
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
  (props.options || []).map((option) => {
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
const enabledOptions = computed(() => normalizedOptions.value.filter((option) => !option.disabled))
const selectedOption = computed(() => {
  const value = resolvedValue.value
  return normalizedOptions.value.find((option) => option.value === value) || null
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
  background: 'var(--lab-bg-canvas)',
  color: 'var(--lab-text-primary)'
}))
const getSwatchStyle = (color?: string) => (color ? { backgroundColor: color } : {})
const setOptionRef = (index: number, el: Element | ComponentPublicInstance | null) => {
  optionRefs.value[index] = el instanceof HTMLElement ? el : null
}
const getOptionRef =
  (index: number): VNodeRef =>
  (el) => {
    setOptionRef(index, el)
  }
const setHighlightedIndex = (index: number) => {
  if (!normalizedOptions.value[index] || normalizedOptions.value[index].disabled) return
  highlightedIndex.value = index
}
const scrollHighlightedOptionIntoView = () => {
  if (!import.meta.client || !isOpen.value) return
  const listboxEl = listboxRef.value
  const optionEl = optionRefs.value[highlightedIndex.value]
  if (!listboxEl || !optionEl) return
  const bufferItems = 2
  const optionHeight = optionEl.offsetHeight || 0
  const buffer = optionHeight * bufferItems
  const optionTop = optionEl.offsetTop
  const optionBottom = optionTop + optionEl.offsetHeight
  const viewportTop = listboxEl.scrollTop
  const viewportBottom = viewportTop + listboxEl.clientHeight
  const nextTop = Math.max(0, optionTop - buffer)
  const maxTop = Math.max(0, listboxEl.scrollHeight - listboxEl.clientHeight)
  const nextBottom = Math.min(maxTop, optionBottom + buffer - listboxEl.clientHeight)
  if (optionTop < viewportTop + buffer) {
    listboxEl.scrollTo({ top: nextTop })
    return
  }
  if (optionBottom > viewportBottom - buffer) {
    listboxEl.scrollTo({ top: nextBottom })
  }
}
const syncHighlightWithValue = () => {
  const selectedIndex = normalizedOptions.value.findIndex(
    (option) => option.value === resolvedValue.value && !option.disabled
  )
  if (selectedIndex >= 0) {
    highlightedIndex.value = selectedIndex
    return
  }
  const firstEnabledIndex = normalizedOptions.value.findIndex((option) => !option.disabled)
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
  (next) => {
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
watch(isOpen, (next) => {
  if (!next) {
    optionRefs.value = []
    return
  }
  if (!import.meta.client) return
  schedulePositionUpdate()
})
watch(highlightedIndex, () => {
  if (!import.meta.client || !isOpen.value) return
  requestAnimationFrame(() => {
    scrollHighlightedOptionIntoView()
  })
})
</script>
