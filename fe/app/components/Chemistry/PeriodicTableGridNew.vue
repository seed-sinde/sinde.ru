<script setup lang="ts">
  export const chemicalElementCategoryColors = [
    'oklab(0.72 0.00 0.00 / 0.24)', // nonmetals
    'oklab(0.78 -0.01 -0.10 / 0.26)', // noble gases
    'oklab(0.70 0.06 0.10 / 0.26)', // alkali metals
    'oklab(0.82 -0.02 0.10 / 0.26)', // alkaline earth metals
    'oklab(0.68 -0.04 0.04 / 0.26)', // metalloids
    'oklab(0.84 -0.06 0.08 / 0.26)', // halogens
    'oklab(0.62 0.02 -0.08 / 0.26)', // transition metals
    'oklab(0.76 0.03 0.03 / 0.26)', // post-transition metals
    'oklab(0.80 0.00 0.12 / 0.26)', // lanthanides
    'oklab(0.74 -0.03 0.12 / 0.26)' // actinides
  ] as const
</script>
<template>
  <section class="grid"></section>
</template>
