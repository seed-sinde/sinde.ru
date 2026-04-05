<script setup lang="ts">
  const title = 'Влияние'
  usePageSeo({
    title,
    description:
      'Положительные практики влияния человека на экологию и качество жизни в обществе: бережное потребление, сохранение природы и развитие устойчивых привычек.'
  })
  const sectionEntries = computed(() => {
    const grouped: Record<string, typeof impact.docs> = {}
    for (const doc of impact.docs) {
      const key = doc.section
      if (!grouped[key]) grouped[key] = []
      grouped[key]!.push(doc)
    }
    return Object.entries(grouped)
  })
</script>
<template>
  <div class="space-y-6">
    <LabNavHeader :title />
    <LabSectionDocGrid
      v-for="[section, docs] in sectionEntries"
      :key="section"
      :section="section"
      :docs="docs"
      base-path="/impact" />
  </div>
</template>
