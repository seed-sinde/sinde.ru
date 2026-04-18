<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    modelValue?: boolean
    items: MenuItem[]
  }>(),
  {
    modelValue: false
  }
)
const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
}>()
const { locale, key, load, t } = useI18nSection('nav')
await useAsyncData(key.value, load, { watch: [locale] })
const isOpen = computed(() => props.modelValue)
const close = () => {
  emit('update:modelValue', false)
}
</script>
<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition-opacity duration-200"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition-opacity duration-200"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div v-if="isOpen" class="fixed inset-0 z-50 lg:hidden">
        <button
          type="button"
          class="absolute inset-0 bg-[color-mix(in_srgb,var(--lab-bg-canvas)_72%,transparent)]"
          :aria-label="t('close_menu')"
          @click="close"
        />
        <Transition
          enter-active-class="transition-transform duration-200"
          enter-from-class="-translate-x-full"
          enter-to-class="translate-x-0"
          leave-active-class="transition-transform duration-200"
          leave-from-class="translate-x-0"
          leave-to-class="-translate-x-full"
        >
          <aside
            v-if="isOpen"
            class="relative flex h-full w-full max-w-fit flex-col border-r border-(--lab-border) bg-(--lab-bg-overlay)"
          >
            <div class="flex items-center justify-between pr-2">
              <LabNavHomeBtn @click="close" />
              <LabBaseButton
                icon="ic:round-close"
                icon-only
                icon-size="lg"
                size="sm"
                variant="ghost"
                :aria-label="t('close_menu')"
                @click="close"
              />
            </div>
            <LabNavSidebar :items="items" class="min-h-0 flex-1" @request-close="close" />
          </aside>
        </Transition>
      </div>
    </Transition>
  </Teleport>
</template>
