<template>
  <div class="space-y-1">
    <LabField label="Список" :for-id="inputId" :hint="hint">
      <LabBaseTextarea
        :id="inputId"
        v-model="raw"
        name="TraitInputList"
        rows="4"
        placeholder="Каждое значение с новой строки"
      />
    </LabField>
    <LabErrorMessage v-if="showMinItemsError" :text="`Нужно не меньше ${minItems} значений.`" />
  </div>
</template>

<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    id?: string
    meta?: Pick<KeyMeta, "minItems" | "listOrdered" | "listUnique" | "listValueType">
  }>(),
  {
    id: ""
  }
)
const model = defineModel<TraitListModel>({required: true})
const inputId = computed(() => props.id || "TraitInputList")
const minItems = computed(() => Math.max(Number(props.meta?.minItems || 0), 0))
const hint = computed(() => {
  const parts = [
    props.meta?.listValueType ? `Тип: ${props.meta.listValueType}` : "",
    props.meta?.listOrdered ? "Упорядоченный" : "Неупорядоченный",
    props.meta?.listUnique ? "Только уникальные" : ""
  ].filter(Boolean)
  return parts.join(" · ")
})
const normalizeItems = (items: string[]) => {
  if (!props.meta?.listUnique) return items
  return Array.from(new Set(items))
}
const raw = computed<string>({
  get: () => (Array.isArray(model.value?.items) ? model.value.items : []).join("\n"),
  set: next => {
    model.value = {
      items: normalizeItems(
        String(next || "")
          .split("\n")
          .map(item => item.trim())
          .filter(Boolean)
      )
    }
  }
})
const showMinItemsError = computed(() => minItems.value > 0 && model.value.items.length > 0 && model.value.items.length < minItems.value)
</script>
