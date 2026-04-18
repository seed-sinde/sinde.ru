<template>
  <div class="space-y-2">
    <template v-if="geoType === 'point'">
      <LabField label="Широта" :for-id="`${baseId}-lat`">
        <LabBaseInput :id="`${baseId}-lat`" v-model="pointLat" inputmode="decimal" placeholder="59.934280" />
      </LabField>
      <LabField label="Долгота" :for-id="`${baseId}-lng`">
        <LabBaseInput :id="`${baseId}-lng`" v-model="pointLng" inputmode="decimal" placeholder="30.335099" />
      </LabField>
      <LabField label="Высота" :for-id="`${baseId}-height`">
        <LabBaseInput :id="`${baseId}-height`" v-model="height" inputmode="decimal" placeholder="0" />
      </LabField>
    </template>
    <template v-else-if="geoType === 'zone'">
      <LabField label="Широта" :for-id="`${baseId}-lat`">
        <LabBaseInput :id="`${baseId}-lat`" v-model="pointLat" inputmode="decimal" placeholder="59.934280" />
      </LabField>
      <LabField label="Долгота" :for-id="`${baseId}-lng`">
        <LabBaseInput :id="`${baseId}-lng`" v-model="pointLng" inputmode="decimal" placeholder="30.335099" />
      </LabField>
      <LabField label="Радиус" :for-id="`${baseId}-radius`">
        <LabBaseInput :id="`${baseId}-radius`" v-model="radius" inputmode="decimal" placeholder="10" />
      </LabField>
    </template>
    <LabField v-else label="Точки полигона" :for-id="`${baseId}-points`" hint="По одной точке на строку: широта,долгота">
      <LabBaseTextarea :id="`${baseId}-points`" v-model="rawPoints" rows="4" placeholder="59.93,30.33" />
    </LabField>
  </div>
</template>

<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    id?: string
    meta?: Pick<KeyMeta, "geoType" | "heightUnit" | "radiusUnit">
  }>(),
  {
    id: ""
  }
)
const model = defineModel<TraitGeoModel>({required: true})
const baseId = computed(() => props.id || "TraitInputGeo")
const geoType = computed(() => resolveGeoType(props.meta?.geoType))
const pointLat = computed({
  get: () => String(model.value?.lat || ""),
  set: next => {
    model.value = {...model.value, type: geoType.value, lat: String(next || "").trim()}
  }
})
const pointLng = computed({
  get: () => String(model.value?.lng || ""),
  set: next => {
    model.value = {...model.value, type: geoType.value, lng: String(next || "").trim()}
  }
})
const height = computed({
  get: () => String(model.value?.height || ""),
  set: next => {
    model.value = {...model.value, type: geoType.value, height: String(next || "").trim(), heightUnit: props.meta?.heightUnit || ""}
  }
})
const radius = computed({
  get: () => String(model.value?.radius || ""),
  set: next => {
    model.value = {...model.value, type: geoType.value, radius: String(next || "").trim(), radiusUnit: props.meta?.radiusUnit || ""}
  }
})
const rawPoints = computed({
  get: () =>
    Array.isArray(model.value?.points)
      ? model.value.points.map(point => `${String(point?.lat || "").trim()},${String(point?.lng || "").trim()}`).join("\n")
      : "",
  set: next => {
    model.value = {
      type: geoType.value,
      points: String(next || "")
        .split("\n")
        .map(line => line.trim())
        .filter(Boolean)
        .map(line => {
          const [lat = "", lng = ""] = line.split(",").map(part => part.trim())
          return {lat, lng}
        })
    }
  }
})
watch(
  geoType,
  next => {
    if (model.value?.type === next) return
    model.value = {
      type: next,
      ...(next === "polygon" ? {points: []} : {})
    }
  },
  {immediate: true}
)
</script>
