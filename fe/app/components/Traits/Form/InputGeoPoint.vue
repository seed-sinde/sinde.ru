<template>
  <LabField label="Координаты (широта, долгота)" :for-id="inputId">
    <LabBaseInput
      :id="inputId"
      name="TraitGeoPoint"
      v-model="raw"
      inputmode="decimal"
      placeholder="59.934280,30.335099"
      @blur="normalize" />
    <LabHint v-if="message" :text="message" hint-class="text-amber-400" />
  </LabField>
</template>
<script setup lang="ts">
  import { normalizeGeoPoint, parseGeoCoordinate, GEO_MIN_DECIMALS } from '../../../utils/traitValueCodec'
  const props = withDefaults(defineProps<{ id?: string; meta?: unknown }>(), { meta: undefined })
  const model = defineModel<TraitGeoPointModel>({ required: true })
  const inputId = computed(() => (props.id ? `${props.id}-geo` : 'TraitGeoPoint'))
  const raw = ref('')
  const message = ref('')
  const normalize = () => {
    message.value = ''
    const parts = raw.value
      .split(',')
      .map(part => part.trim())
      .filter(Boolean)
    if (parts.length !== 2) {
      message.value = 'Введите и широту, и долготу'
      return
    }
    const [latRaw, lngRaw] = parts as [string, string]
    const latParsed = parseGeoCoordinate(latRaw)
    const lngParsed = parseGeoCoordinate(lngRaw)
    if (!latParsed || !lngParsed) {
      message.value = 'Широта и долгота должны быть числами'
      return
    }
    const normalized = normalizeGeoPoint(latRaw, lngRaw)
    if (!normalized) {
      message.value = 'Широта и долгота должны быть числами'
      return
    }
    if (latParsed.decimals < GEO_MIN_DECIMALS || lngParsed.decimals < GEO_MIN_DECIMALS) {
      message.value = `Недостаточная точность (нужно не меньше ${GEO_MIN_DECIMALS} знаков после запятой)`
    }
    raw.value = `${normalized.lat},${normalized.lng}`
    model.value = { lat: normalized.lat, lng: normalized.lng }
  }
  watchEffect(() => {
    if (!model.value) return
    raw.value = `${model.value.lat},${model.value.lng}`
  })
</script>
