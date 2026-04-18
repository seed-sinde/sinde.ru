<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    items?: MobileHeaderAction[]
  }>(),
  {
    items: () => []
  }
)
const { pasteUuidAndNavigate } = usePasteUuid()
const onPasteClick = async (action: Extract<MobileHeaderAction, { kind: 'traits-paste-uuid' }>) => {
  await pasteUuidAndNavigate(uuid => (action.mode === 'trait' ? `/traits/trait/${uuid}` : `/traits/${uuid}`))
}
</script>
<template>
  <div class="flex shrink-0 items-center gap-1">
    <template v-for="item in props.items" :key="item.kind === 'traits-copy-uuid' ? item.uuid : item.mode">
      <TraitsUuidButton v-if="item.kind === 'traits-copy-uuid'" action="copy" :uuid="item.uuid" :compact="true" />
      <TraitsUuidButton v-else action="paste" :compact="true" @click="onPasteClick(item)" />
    </template>
  </div>
</template>
