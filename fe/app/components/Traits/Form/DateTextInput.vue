<template>
  <LabBaseInput
    :id="id"
    :name="name"
    :model-value="displayValue"
    :placeholder="resolvedPlaceholder"
    :aria-label="ariaLabel"
    :invalid="invalid"
    :maxlength="maxLength"
    autocomplete="off"
    spellcheck="false"
    inputmode="numeric"
    @update:model-value="onInputValue"
    @focus="onFocus"
    @blur="onBlur"
    @keydown="onKeydown"
    @paste="onPaste" />
</template>
<script setup lang="ts">
  type InputDateTextMode = 'date' | 'time' | 'datetime'
  const model = defineModel<string>({ required: true })
  const props = withDefaults(
    defineProps<{
      mode?: InputDateTextMode
      id?: string
      name?: string
      placeholder?: string
      ariaLabel?: string
      invalid?: boolean
    }>(),
    {
      mode: 'datetime',
      id: '',
      name: '',
      placeholder: '',
      ariaLabel: '',
      invalid: false
    }
  )
  const resolvedMode = computed<InputDateTextMode>(() => props.mode || 'datetime')
  const resolvedPlaceholder = computed(() => props.placeholder || datePlaceholder(resolvedMode.value))
  const maxLength = computed(() => dateTextMaxLength(resolvedMode.value))
  const isFocused = ref(false)
  const displayValue = ref(maskDateText(model.value || '', resolvedMode.value))
  const maskSpec = (mode: InputDateTextMode): DateTextMaskSpec => {
    if (mode === 'time') return { template: '__:__', slots: [0, 1, 3, 4] }
    if (mode === 'date') return { template: '__.__.____', slots: [0, 1, 3, 4, 6, 7, 8, 9] }
    return { template: '__.__.____ __:__', slots: [0, 1, 3, 4, 6, 7, 8, 9, 11, 12, 14, 15] }
  }
  const toSlots = (value: string, mode: InputDateTextMode): string[] => {
    const digits = String(value || '').replace(/\D/g, '')
    const { slots } = maskSpec(mode)
    return slots.map((_, idx) => digits[idx] || '')
  }
  const fromSlots = (slotsValue: string[]): string => {
    return slotsValue.join('')
  }
  const renderSlots = (slotsValue: string[], mode: InputDateTextMode): string => {
    const spec = maskSpec(mode)
    const chars = spec.template.split('')
    spec.slots.forEach((pos, idx) => {
      chars[pos] = slotsValue[idx] || '_'
    })
    return chars.join('')
  }
  const findSlotAtOrAfter = (caret: number, mode: InputDateTextMode): number => {
    const positions = maskSpec(mode).slots
    const index = positions.findIndex(pos => pos >= caret)
    return index >= 0 ? index : positions.length - 1
  }
  const findSlotAtOrBefore = (caret: number, mode: InputDateTextMode): number => {
    const positions = maskSpec(mode).slots
    for (let idx = positions.length - 1; idx >= 0; idx -= 1) {
      const pos = positions[idx]
      if (pos !== undefined && pos <= caret) return idx
    }
    return 0
  }
  const setCaret = (input: HTMLInputElement, caret: number) => {
    requestAnimationFrame(() => {
      try {
        input.setSelectionRange(caret, caret)
      } catch {
        // ignore unsupported input types
      }
    })
  }
  const setCaretBySlot = (input: HTMLInputElement, slotIndex: number, mode: InputDateTextMode) => {
    const positions = maskSpec(mode).slots
    const templateLen = maskSpec(mode).template.length
    if (slotIndex >= positions.length) {
      setCaret(input, templateLen)
      return
    }
    const safeIndex = Math.max(0, slotIndex)
    const caret = positions[safeIndex] ?? positions[positions.length - 1] ?? 0
    setCaret(input, caret)
  }
  const syncModelFromSlots = (slotsValue: string[], mode: InputDateTextMode) => {
    model.value = maskDateText(fromSlots(slotsValue), mode)
  }
  const onInputValue = (raw: string) => {
    if (!isFocused.value) {
      model.value = maskDateText(raw, resolvedMode.value)
      displayValue.value = model.value
    }
  }
  const onFocus = (event: FocusEvent) => {
    isFocused.value = true
    const mode = resolvedMode.value
    const input = event.target as HTMLInputElement
    const slotsValue = toSlots(model.value || '', mode)
    displayValue.value = renderSlots(slotsValue, mode)
    const firstEmpty = slotsValue.findIndex(item => !item)
    const targetSlot = firstEmpty >= 0 ? firstEmpty : slotsValue.length
    setCaretBySlot(input, targetSlot, mode)
  }
  const onBlur = () => {
    isFocused.value = false
    const mode = resolvedMode.value
    const normalized = maskDateText(model.value || '', mode)
    model.value = normalized
    displayValue.value = normalized
  }
  const onKeydown = (event: KeyboardEvent) => {
    const mode = resolvedMode.value
    const input = event.target as HTMLInputElement
    if (!input) return
    if (event.ctrlKey || event.metaKey || event.altKey) return
    const slotsValue = toSlots(displayValue.value, mode)
    const key = event.key
    const caret = input.selectionStart ?? 0
    const selectionEnd = input.selectionEnd ?? caret
    const clearSelection = () => {
      if (selectionEnd <= caret) return
      const startSlot = findSlotAtOrAfter(caret, mode)
      const endSlot = findSlotAtOrBefore(selectionEnd - 1, mode)
      for (let idx = startSlot; idx <= endSlot; idx += 1) slotsValue[idx] = ''
    }
    if (/^\d$/.test(key)) {
      event.preventDefault()
      clearSelection()
      const slotIndex = findSlotAtOrAfter(caret, mode)
      slotsValue[slotIndex] = key
      displayValue.value = renderSlots(slotsValue, mode)
      syncModelFromSlots(slotsValue, mode)
      setCaretBySlot(input, slotIndex + 1, mode)
      return
    }
    if (key === 'Backspace') {
      event.preventDefault()
      if (selectionEnd > caret) {
        clearSelection()
        displayValue.value = renderSlots(slotsValue, mode)
        syncModelFromSlots(slotsValue, mode)
        setCaretBySlot(input, findSlotAtOrAfter(caret, mode), mode)
        return
      }
      const slotIndex = findSlotAtOrBefore(caret - 1, mode)
      slotsValue[slotIndex] = ''
      displayValue.value = renderSlots(slotsValue, mode)
      syncModelFromSlots(slotsValue, mode)
      setCaretBySlot(input, slotIndex, mode)
      return
    }
    if (key === 'Delete') {
      event.preventDefault()
      if (selectionEnd > caret) {
        clearSelection()
        displayValue.value = renderSlots(slotsValue, mode)
        syncModelFromSlots(slotsValue, mode)
        setCaretBySlot(input, findSlotAtOrAfter(caret, mode), mode)
        return
      }
      const slotIndex = findSlotAtOrAfter(caret, mode)
      slotsValue[slotIndex] = ''
      displayValue.value = renderSlots(slotsValue, mode)
      syncModelFromSlots(slotsValue, mode)
      setCaretBySlot(input, slotIndex, mode)
    }
  }
  const onPaste = (event: ClipboardEvent) => {
    event.preventDefault()
    const mode = resolvedMode.value
    const input = event.target as HTMLInputElement
    if (!input) return
    const pastedDigits = String(event.clipboardData?.getData('text') || '').replace(/\D/g, '')
    if (!pastedDigits) return
    const slotsValue = toSlots(displayValue.value, mode)
    let slotIndex = findSlotAtOrAfter(input.selectionStart ?? 0, mode)
    for (const digit of pastedDigits) {
      if (slotIndex >= slotsValue.length) break
      slotsValue[slotIndex] = digit
      slotIndex += 1
    }
    displayValue.value = renderSlots(slotsValue, mode)
    syncModelFromSlots(slotsValue, mode)
    setCaretBySlot(input, slotIndex, mode)
  }
  watch(
    [() => model.value, resolvedMode],
    ([nextModel, mode]) => {
      if (isFocused.value) return
      displayValue.value = maskDateText(String(nextModel || ''), mode)
    },
    { immediate: true }
  )
</script>
