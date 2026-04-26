<script setup lang="ts">
definePageMeta({tabs: true})
const tabs = [
  {to: "/projects", label: "Проекты"},
  {to: "/new", label: "Новый проект"}
]
</script>
<template>
  <main class="p-4">
    <UiTabs :items="tabs">
      <template #panel-projects="{item}">
        <p class="text-sm leading-snug">{{ item.label }}: содержимое первой вкладки через слот.</p>
      </template>
      <template #panel-new="{item}">
        <p class="text-sm leading-snug">{{ item.label }}: отдельная панель второй вкладки.</p>
      </template>
    </UiTabs>
  </main>
</template>
