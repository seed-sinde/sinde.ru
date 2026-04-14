<template>
  <div class="space-y-3 bg-transparent p-0">
    <div v-if="!lockedMode" class="flex flex-wrap gap-2 text-sm text-zinc-300">
      <label
        v-for="entry in modes"
        :key="entry"
        class="inline-flex cursor-pointer items-center gap-2 rounded-full border border-zinc-700/80 bg-zinc-950/75 px-3 py-1.5 transition hover:border-amber-400/60 hover:text-amber-200"
      >
        <input
          :id="`TraitColorMode-${entry}`"
          v-model="mode"
          name="TraitColorMode"
          type="radio"
          :value="entry"
          class="accent-amber-400"
        >
        <span>{{ modeLabels[entry] }}</span>
      </label>
    </div>
    <div v-if="mode === 'hex'" class="space-y-2">
      <LabField label="HEX" for-id="color-hex">
        <LabBaseInput
          id="color-hex"
          v-model="hex"
          name="color_hex"
          placeholder="#aabbcc"
          pattern="^#?[0-9a-fA-F]{6}$"
        />
      </LabField>
      <TraitsColorPicker v-model="hex" />
    </div>
    <div v-else-if="mode === 'lab'" class="space-y-2">
      <p class="text-xs font-medium tracking-[0.08em] text-zinc-500 uppercase">Слайдеры LAB</p>
      <div class="space-y-1">
        <div class="flex justify-between font-mono text-xs text-zinc-400">
          <span>L</span>
          <span>{{ lab.L.toFixed(0) }}</span>
        </div>
        <input
          v-model.number="lab.L"
          name="color_lab_l"
          type="range"
          min="0"
          max="100"
          step="1"
          class="h-2 w-full rounded-lg bg-zinc-800 accent-amber-400"
          :style="gradientL"
        >
      </div>
      <div class="space-y-1">
        <div class="flex justify-between font-mono text-xs text-zinc-400">
          <span>a</span>
          <span>{{ lab.a.toFixed(0) }}</span>
        </div>
        <input
          v-model.number="lab.a"
          name="color_lab_a"
          type="range"
          min="-128"
          max="127"
          step="1"
          class="h-2 w-full rounded-lg bg-zinc-800 accent-amber-400"
          :style="gradientA"
        >
      </div>
      <div class="space-y-1">
        <div class="flex justify-between font-mono text-xs text-zinc-400">
          <span>b</span>
          <span>{{ lab.b.toFixed(0) }}</span>
        </div>
        <input
          v-model.number="lab.b"
          name="color_lab_b"
          type="range"
          min="-128"
          max="127"
          step="1"
          class="h-2 w-full rounded-lg bg-zinc-800 accent-amber-400"
          :style="gradientB"
        >
      </div>
      <p class="pt-1 text-xs font-medium tracking-[0.08em] text-zinc-500 uppercase">Ручной ввод LAB</p>
      <div class="grid grid-cols-1 gap-2 sm:grid-cols-3">
        <LabField label="L" for-id="color-lab-l-input">
          <LabBaseInput
            id="color-lab-l-input"
            name="color_lab_l_input"
            :model-value="labInput.L"
            type="number"
            min="0"
            max="100"
            step="0.1"
            inputmode="decimal"
            @update:model-value="onLabBaseInput('L', $event)"
            @blur="onLabBaseInputBlur('L')"
          />
        </LabField>
        <LabField label="a" for-id="color-lab-a-input">
          <LabBaseInput
            id="color-lab-a-input"
            name="color_lab_a_input"
            :model-value="labInput.a"
            type="number"
            min="-128"
            max="127"
            step="0.1"
            inputmode="decimal"
            @update:model-value="onLabBaseInput('a', $event)"
            @blur="onLabBaseInputBlur('a')"
          />
        </LabField>
        <LabField label="b" for-id="color-lab-b-input">
          <LabBaseInput
            id="color-lab-b-input"
            name="color_lab_b_input"
            :model-value="labInput.b"
            type="number"
            min="-128"
            max="127"
            step="0.1"
            inputmode="decimal"
            @update:model-value="onLabBaseInput('b', $event)"
            @blur="onLabBaseInputBlur('b')"
          />
        </LabField>
      </div>
      <div class="flex items-center gap-3">
        <div class="h-10 w-14 rounded-xl border border-zinc-700/80" :style="{ background: labPreviewColor }" />
        <div class="font-mono text-xs text-zinc-300">
          {{ lab.L.toFixed(0) }},{{ lab.a.toFixed(0) }},{{ lab.b.toFixed(0) }}
        </div>
      </div>
    </div>
    <div v-else class="space-y-2">
      <LabField label="Категория цвета" for-id="color-spectrum">
        <LabBaseSelect id="color-spectrum" v-model="spectrum" name="color_spectrum" :options="spectrumSelectOptions" />
      </LabField>
      <div class="h-10 w-14 rounded-xl border border-zinc-700/80" :style="{ background: spectrumColor }" />
    </div>
  </div>
</template>
<script setup lang="ts">
import TraitsColorPicker from './ColorPicker.vue'
import { COLOR_MODE_LABELS, resolveColorMode } from '../../../utils/traitMeta'
import {
  COLOR_SPECTRUM_OPTIONS,
  COLOR_SPECTRUM_LABELS,
  COLOR_SPECTRUM_MAP,
  labToRgb,
  rgbToCss
} from '../../../utils/traitColor'
const props = defineProps<{ meta?: unknown }>()
const model = defineModel<Color>({ required: true })
const modes = ['hex', 'lab', 'spectrum'] as const
const modeLabels = COLOR_MODE_LABELS
const mode = ref<Color['mode']>(model.value.mode ?? 'hex')
const hex = ref(model.value.hex ?? '')
const lab = reactive({
  L: model.value.lab?.L ?? 50,
  a: model.value.lab?.a ?? 0,
  b: model.value.lab?.b ?? 0
})
const labInput = reactive({
  L: String(model.value.lab?.L ?? 50),
  a: String(model.value.lab?.a ?? 0),
  b: String(model.value.lab?.b ?? 0)
})
const spectrum = ref(model.value.spectrum ?? '')
let syncingLabBaseInput = false
const formatLabNumber = (value: number): string => {
  if (Number.isInteger(value)) return String(value)
  return String(Number(value.toFixed(3)))
}
const clamp = (value: number, min: number, max: number): number => {
  if (value < min) return min
  if (value > max) return max
  return value
}
const parseManualNumber = (raw: string): number | null => {
  const normalized = String(raw || '')
    .trim()
    .replace(',', '.')
  if (!normalized) return null
  const num = Number(normalized)
  return Number.isFinite(num) ? num : null
}
const onLabBaseInput = (channel: 'L' | 'a' | 'b', raw: string) => {
  labInput[channel] = String(raw ?? '')
  const parsed = parseManualNumber(raw)
  if (parsed === null) return
  syncingLabBaseInput = true
  if (channel === 'L') lab.L = clamp(parsed, 0, 100)
  if (channel === 'a') lab.a = clamp(parsed, -128, 127)
  if (channel === 'b') lab.b = clamp(parsed, -128, 127)
  syncingLabBaseInput = false
}
const onLabBaseInputBlur = (channel: 'L' | 'a' | 'b') => {
  if (channel === 'L') labInput.L = formatLabNumber(lab.L)
  if (channel === 'a') labInput.a = formatLabNumber(lab.a)
  if (channel === 'b') labInput.b = formatLabNumber(lab.b)
}
const gradientL = computed(() => {
  const left = rgbToCss(labToRgb(0, lab.a, lab.b))
  const right = rgbToCss(labToRgb(100, lab.a, lab.b))
  return { background: `linear-gradient(90deg, ${left}, ${right})` }
})
const gradientA = computed(() => {
  const left = rgbToCss(labToRgb(lab.L, -128, lab.b))
  const right = rgbToCss(labToRgb(lab.L, 127, lab.b))
  return { background: `linear-gradient(90deg, ${left}, ${right})` }
})
const gradientB = computed(() => {
  const left = rgbToCss(labToRgb(lab.L, lab.a, -128))
  const right = rgbToCss(labToRgb(lab.L, lab.a, 127))
  return { background: `linear-gradient(90deg, ${left}, ${right})` }
})
const labPreviewColor = computed(() => rgbToCss(labToRgb(lab.L, lab.a, lab.b)))
const spectrumOptions = COLOR_SPECTRUM_OPTIONS
const spectrumLabels = COLOR_SPECTRUM_LABELS
const spectrumSelectOptions = spectrumOptions.map((option) => ({
  value: option,
  label: spectrumLabels[option]
}))
const spectrumColor = computed(() => COLOR_SPECTRUM_MAP[spectrum.value] || '#6b7280')
const lockedMode = computed<Color['mode'] | ''>(() => {
  return resolveColorMode(props.meta && typeof props.meta === 'object' ? (props.meta as Record<string, unknown>) : null)
})
watchEffect(() => {
  if (lockedMode.value && mode.value !== lockedMode.value) {
    mode.value = lockedMode.value
  }
  model.value =
    mode.value === 'hex'
      ? { mode: 'hex', hex: hex.value }
      : mode.value === 'lab'
        ? { mode: 'lab', lab: { L: lab.L, a: lab.a, b: lab.b } }
        : { mode: 'spectrum', spectrum: spectrum.value }
})
watch(model, (next) => {
  mode.value = lockedMode.value || next.mode
  hex.value = next.hex ?? ''
  lab.L = next.lab?.L ?? 50
  lab.a = next.lab?.a ?? 0
  lab.b = next.lab?.b ?? 0
  labInput.L = formatLabNumber(lab.L)
  labInput.a = formatLabNumber(lab.a)
  labInput.b = formatLabNumber(lab.b)
  spectrum.value = next.spectrum ?? ''
})
watch(
  () => [lab.L, lab.a, lab.b] as const,
  ([nextL, nextA, nextB]) => {
    if (syncingLabBaseInput) return
    labInput.L = formatLabNumber(nextL)
    labInput.a = formatLabNumber(nextA)
    labInput.b = formatLabNumber(nextB)
  }
)
</script>
