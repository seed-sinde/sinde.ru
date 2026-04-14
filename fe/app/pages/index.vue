<script setup lang="ts">
const brand = 'SINDE'
const title = 'главная'
const letters = brand.split('')
usePageSeo({ title, description: 'Платформа прикладных инструментов и справочных знаний.' })
</script>

<template>
  <section class="relative flex min-h-[calc(100svh-48px)] items-center justify-center overflow-hidden">
    <h1 class="flex gap-x-6 text-[clamp(2rem,6vw,4rem)] leading-none font-light tracking-[0.45em] select-none">
      <span
        v-for="letter in letters"
        :key="letter"
        class="bg-linear-to-b from-fuchsia-950 via-cyan-300 to-fuchsia-400 bg-clip-text text-transparent"
      >
        {{ letter }}
      </span>
    </h1>
  </section>
</template>
