<template>
  <div class="grid gap-2">
    <LabBaseField label="Категория блеска" :for-id="`${baseId}-gloss-category`">
      <LabBaseInput :id="`${baseId}-gloss-category`" v-model="glossCategory" placeholder="матовый, глянец" />
    </LabBaseField>
    <LabBaseField label="Число GU" :for-id="`${baseId}-gloss-gu`">
      <LabBaseInput :id="`${baseId}-gloss-gu`" v-model="glossGU" type="number" inputmode="decimal" placeholder="0" />
    </LabBaseField>
    <LabBaseField label="Тип рельефа" :for-id="`${baseId}-relief-type`">
      <LabBaseInput :id="`${baseId}-relief-type`" v-model="reliefType" placeholder="песок, кожа, сатин" />
    </LabBaseField>
    <LabBaseField label="Высота микрорельефа" :for-id="`${baseId}-micro`">
      <LabBaseInput
        :id="`${baseId}-micro`"
        v-model="microReliefHeight"
        type="number"
        inputmode="decimal"
        placeholder="0"
      />
    </LabBaseField>
  </div>
</template>

<script setup lang="ts">
const props = withDefaults(defineProps<{id?: string}>(), {id: ""})
const model = defineModel<TraitSurfaceModel>({required: true})
const baseId = computed(() => props.id || "TraitInputSurface")
const glossCategory = computed({
  get: () => String(model.value?.glossCategory || ""),
  set: next => {
    model.value = {...model.value, glossCategory: String(next || "").trim()}
  }
})
const glossGU = computed({
  get: () => String(model.value?.glossGU || ""),
  set: next => {
    model.value = {...model.value, glossGU: String(next || "").trim()}
  }
})
const reliefType = computed({
  get: () => String(model.value?.reliefType || ""),
  set: next => {
    model.value = {...model.value, reliefType: String(next || "").trim()}
  }
})
const microReliefHeight = computed({
  get: () => String(model.value?.microReliefHeight || ""),
  set: next => {
    model.value = {...model.value, microReliefHeight: String(next || "").trim()}
  }
})
</script>
