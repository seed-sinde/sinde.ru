<script setup lang="ts">
  usePageSeo({
    title: 'главная',
    description: 'Платформа прикладных инструментов и справочных знаний.'
  })
  const letters = ['S', 'I', 'N', 'D', 'E']
</script>
<template>
  <section class="bg-(--lab-bg-canvas)">
    <div class="flex w-3/7 justify-between overflow-hidden line-none px-4 py-2">
      <span v-for="char in letters" :key="char" class="inline-block text-2xl font-bold leading-none origin-center">
        {{ char }}
      </span>
    </div>
  </section>
</template>
