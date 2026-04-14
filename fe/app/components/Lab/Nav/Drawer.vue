<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    modelValue?: boolean
    items: MenuItem[]
    homeTo?: string
  }>(),
  {
    modelValue: false,
    homeTo: '/'
  }
)
const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
}>()
const { t, faviconLightSrc, faviconDarkSrc } = useInterfacePreferences()
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
          :aria-label="t('nav.close_menu')"
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
            class="relative flex h-full w-full max-w-fit flex-col border-r bg-(--lab-bg-overlay) pt-1"
          >
            <div class="mb-1 flex items-center justify-between border-b px-1 pb-1">
              <NuxtLink
                :to="homeTo"
                class="inline-flex h-9 w-9 items-center justify-center text-(--lab-text-primary) transition-colors hover:bg-(--lab-bg-surface-hover) focus-visible:bg-(--lab-bg-surface-hover)"
                :aria-label="t('nav.home')"
                @click="close"
              >
                <picture>
                  <source :srcset="faviconDarkSrc" media="(prefers-color-scheme: dark)" >
                  <img :src="faviconLightSrc" alt="" class="h-4.5 w-4.5 object-contain" >
                </picture>
              </NuxtLink>
              <LabBaseButton
                icon="ic:round-close"
                icon-only
                size="sm"
                variant="ghost"
                :aria-label="t('nav.close_menu')"
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
