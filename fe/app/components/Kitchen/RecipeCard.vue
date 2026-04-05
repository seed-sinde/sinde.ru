<script setup lang="ts">
  const props = withDefaults(
    defineProps<{
      recipe: KitchenRecipe
      to: string
      coverSrc?: string
      metaItems?: KitchenRecipeMetaItem[]
      canManage?: boolean
      showFavorite?: boolean
      favorite?: boolean
      favoritePending?: boolean
      showModeration?: boolean
      moderationText?: string
      moderationNote?: string
      deleteDisabled?: boolean
      editLabel?: string
    }>(),
    {
      coverSrc: '',
      metaItems: () => [],
      canManage: false,
      showFavorite: false,
      favorite: false,
      favoritePending: false,
      showModeration: false,
      moderationText: '',
      moderationNote: '',
      deleteDisabled: false,
      editLabel: 'Редактировать'
    }
  )
  const emit = defineEmits<{
    edit: []
    delete: []
    toggleFavorite: []
  }>()
</script>
<template>
  <article class="bg-zinc-900/70 p-4 space-y-3">
    <NuxtLink v-if="coverSrc" :to="to" class="group relative block">
      <img :src="coverSrc" alt="Фото готового блюда" class="h-36 w-full border border-zinc-700 object-cover" />
      <LabBaseButton
        v-if="showFavorite"
        :button-class="[
          'absolute left-2 top-2 inline-flex h-8 w-8 items-center justify-center overflow-hidden rounded-full border transition',
          favorite ? 'border-rose-300/90 text-white' : 'border-zinc-600/90 text-rose-200 hover:border-rose-400/80'
        ]"
        :aria-label="
          favorite ? `Убрать рецепт ${recipe.title} из избранного` : `Добавить рецепт ${recipe.title} в избранное`
        "
        :disabled="favoritePending"
        :icon="favorite ? 'ic:round-favorite' : 'ic:round-favorite-border'"
        icon-class="relative z-1 h-4.5 w-4.5"
        icon-only
        @click.stop.prevent="emit('toggleFavorite')">
        <span class="absolute inset-0 bg-zinc-950/80"></span>
        <span
          class="absolute inset-0 transition-opacity"
          :class="favorite ? 'bg-rose-600/75' : 'bg-rose-500/20 hover:bg-rose-500/35'"></span>
      </LabBaseButton>
    </NuxtLink>
    <div class="flex items-start gap-2">
      <LabBaseButton
        v-if="showFavorite && !coverSrc"
        :button-class="[
          'inline-flex h-8 w-8 shrink-0 items-center justify-center rounded-full border transition',
          favorite
            ? 'border-rose-300/90 bg-rose-600 text-white hover:bg-rose-500'
            : 'border-zinc-700 bg-zinc-900 text-rose-300 hover:border-rose-500/70 hover:bg-rose-500/10'
        ]"
        :aria-label="
          favorite ? `Убрать рецепт ${recipe.title} из избранного` : `Добавить рецепт ${recipe.title} в избранное`
        "
        :disabled="favoritePending"
        :icon="favorite ? 'ic:round-favorite' : 'ic:round-favorite-border'"
        icon-class="h-4.5 w-4.5"
        icon-only
        @click="emit('toggleFavorite')" />
      <div class="min-w-0 flex-1 space-y-1">
        <NuxtLink
          :to="to"
          class="flex min-w-0 items-start gap-1 text-base font-semibold text-zinc-100 hover:text-zinc-50">
          <span class="min-w-0 flex-1 truncate">{{ recipe.title }}</span>
          <Icon name="ic:round-open-in-new" class="h-4 w-4 shrink-0 text-zinc-400" />
        </NuxtLink>
        <p v-if="showModeration && moderationText" class="text-[11px] text-zinc-500">{{ moderationText }}</p>
        <p v-if="showModeration && moderationNote" class="text-[11px] text-amber-300/90">
          Причина отклонения: {{ moderationNote }}
        </p>
        <div v-if="canManage" class="flex flex-wrap gap-2">
          <LabBaseButton
            button-class="h-7 border border-zinc-700 bg-zinc-900 px-2.5 py-1 text-xs text-zinc-200 transition hover:bg-zinc-800"
            @click="emit('edit')">
            {{ editLabel }}
          </LabBaseButton>
          <LabConfirmActionButton
            label="Удалить"
            confirm-label="Подтвердить"
            tooltip="Подтвердить удаление рецепта?"
            button-class="h-7 w-[112px] px-2.5 py-1 text-xs"
            idle-class="border border-rose-500/50 bg-rose-500/10 text-rose-300 hover:bg-rose-500/20"
            confirm-class="border border-rose-300/90 bg-rose-600 text-white hover:bg-rose-500"
            progress-class="bg-rose-300/45"
            :disabled="deleteDisabled"
            @confirm="emit('delete')" />
        </div>
      </div>
    </div>
    <p class="text-sm text-zinc-400">{{ recipe.description || 'Без описания' }}</p>
    <div v-if="metaItems.length" class="flex flex-wrap gap-2 text-xs text-zinc-500">
      <span v-for="item in metaItems" :key="`${recipe.id}:${item.label}`">{{ item.label }}: {{ item.value }}</span>
    </div>
    <div class="flex flex-wrap gap-1">
      <span
        v-for="ing in recipe.ingredients.slice(0, 8)"
        :key="`${recipe.id}:${ing.name}`"
        class="border border-zinc-700 px-2 py-0.5 text-[11px] text-zinc-200">
        {{ ing.name }}
      </span>
    </div>
  </article>
</template>
