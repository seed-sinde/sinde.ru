<template>
  <LabField label="Координаты (широта, долгота)" :for-id="inputId">
    <LabBaseInput
      :id="inputId"
      v-model="raw"
      name="TraitGeoPoint"
      inputmode="decimal"
      placeholder="59.934280,30.335099"
      @input="onInput"
      @blur="normalize"
    />
    <LabHint v-if="message" :text="message" hint-class="text-amber-400" />
  </LabField>
</template>
<script setup lang="ts">
import { normalizeGeoPoint, parseGeoCoordinate, GEO_MIN_DECIMALS } from '../../../utils/traitValueCodec'
const props = defineProps<{ id?: string; meta?: unknown }>()
const model = defineModel<TraitGeoPointModel>({ required: true })
const inputId = computed(() => (props.id ? `${props.id}-geo` : 'TraitGeoPoint'))
const raw = ref('')
const message = ref('')
const emptyPoint = (): TraitGeoPointModel => ({ lat: '', lng: '' })
const syncFromRaw = (options: { normalizeRaw: boolean; showErrors: boolean }) => {
  if (!options.showErrors) {
    message.value = ''
  }
  const parts = raw.value
    .split(',')
    .map(part => part.trim())
    .filter(Boolean)
  if (parts.length !== 2) {
    model.value = emptyPoint()
    if (options.showErrors) {
      message.value = 'Введите и широту, и долготу'
    }
    return
  }
  const [latRaw, lngRaw] = parts as [string, string]
  const latParsed = parseGeoCoordinate(latRaw)
  const lngParsed = parseGeoCoordinate(lngRaw)
  if (!latParsed || !lngParsed) {
    model.value = emptyPoint()
    if (options.showErrors) {
      message.value = 'Широта и долгота должны быть числами'
    }
    return
  }
  const normalized = normalizeGeoPoint(latRaw, lngRaw)
  if (!normalized) {
    model.value = emptyPoint()
    if (options.showErrors) {
      message.value = 'Широта и долгота должны быть числами'
    }
    return
  }
  if (options.normalizeRaw) {
    raw.value = `${normalized.lat},${normalized.lng}`
  }
  model.value = { lat: normalized.lat, lng: normalized.lng }
  if (options.showErrors) {
    message.value = ''
    if (latParsed.decimals < GEO_MIN_DECIMALS || lngParsed.decimals < GEO_MIN_DECIMALS) {
      message.value = `Недостаточная точность (нужно не меньше ${GEO_MIN_DECIMALS} знаков после запятой)`
    }
  }
}
const onInput = () => {
  syncFromRaw({ normalizeRaw: false, showErrors: false })
}
const normalize = () => {
  syncFromRaw({ normalizeRaw: true, showErrors: true })
}
watchEffect(() => {
  if (!model.value) return
  const lat = String(model.value.lat || '').trim()
  const lng = String(model.value.lng || '').trim()
  const nextRaw = lat && lng ? `${lat},${lng}` : ''
  if (raw.value !== nextRaw) {
    raw.value = nextRaw
  }
})
</script>
