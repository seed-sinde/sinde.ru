<script setup lang="ts">
usePageSeo({
  title: 'UI lab.css preview',
  description: 'Preview of global UI classes moved into ui-lab.css.'
})

const text = ref('Лабораторный ввод')
const note = ref('Многострочный пример')
const selectValue = ref('beta')
const tabValue = ref<'one' | 'two' | 'three'>('one')

const buttonVariants = [
  { variant: 'default', label: 'default' },
  { variant: 'secondary', label: 'secondary' },
  { variant: 'primary', label: 'primary' },
  { variant: 'ghost', label: 'ghost' },
  { variant: 'plain', label: 'plain' },
  { variant: 'danger', label: 'danger' },
  { variant: 'success', label: 'success' },
  { variant: 'info', label: 'info' }
] as const

const badgeVariants = [
  { variant: 'default', label: 'default' },
  { variant: 'muted', label: 'muted' },
  { variant: 'info', label: 'info' },
  { variant: 'success', label: 'success' },
  { variant: 'warning', label: 'warning' },
  { variant: 'danger', label: 'danger' }
] as const

const selectOptions = [
  { value: 'alpha', label: 'alpha' },
  { value: 'beta', label: 'beta' },
  { value: 'gamma', label: 'gamma' }
]

const tableColumns = [
  { key: 'name', label: 'name' },
  { key: 'state', label: 'state', nowrap: true },
  { key: 'value', label: 'value', nowrap: true }
]

const tableRows = [
  { id: '1', name: 'north line', state: 'ready', value: '12' },
  { id: '2', name: 'south line', state: 'queued', value: '7' },
  { id: '3', name: 'west line', state: 'paused', value: '3' }
]

const tabItems = [
  { value: 'one', label: 'one' },
  { value: 'two', label: 'two' },
  { value: 'three', label: 'three' }
]

const scrollItems = [
  'grain flow',
  'amber field',
  'quiet branch',
  'signal path',
  'iron pulse',
  'north vector',
  'slow drift',
  'clean edge',
  'plain text',
  'open lane',
  'soft grid',
  'dry channel'
]
</script>

<template>
  <div class="space-y-8 p-4 text-(--lab-text-secondary)">
    <section class="space-y-3">
      <h1 class="text-lg text-(--lab-text-primary)">ui-lab.css preview</h1>
      <p class="lab-field-label">
        Global controls, buttons, badges, tables, tabs and scroll helpers now live in ui-lab.css.
      </p>
    </section>

    <section class="space-y-3">
      <h2 class="text-base text-(--lab-text-primary)">controls</h2>
      <div class="grid gap-3 min-[900px]:grid-cols-2">
        <LabBaseInput v-model="text" placeholder="lab-control" />
        <LabBaseInput v-model="text" invalid placeholder="lab-control-invalid" />
        <LabBaseTextarea v-model="note" :rows="3" placeholder="textarea" />
        <LabBaseSelect v-model="selectValue" :options="selectOptions" placeholder="select value" />
      </div>
    </section>

    <section class="space-y-3">
      <h2 class="text-base text-(--lab-text-primary)">buttons</h2>
      <div class="flex flex-wrap gap-3">
        <LabBaseButton v-for="item in buttonVariants" :key="item.variant" :variant="item.variant" class="min-h-11 px-3">
          {{ item.label }}
        </LabBaseButton>
        <LabBaseButton variant="secondary" class="min-h-11 px-3">manual focus</LabBaseButton>
      </div>
    </section>

    <section class="space-y-3">
      <h2 class="text-base text-(--lab-text-primary)">badges</h2>
      <div class="flex flex-wrap gap-2">
        <LabBaseBadge v-for="item in badgeVariants" :key="item.variant" :variant="item.variant" size="md">
          {{ item.label }}
        </LabBaseBadge>
      </div>
    </section>

    <section class="space-y-3">
      <h2 class="text-base text-(--lab-text-primary)">table</h2>
      <LabDataTable :columns="tableColumns" :rows="tableRows" empty-text="no rows">
        <template #cell-state="{ value }">
          <LabBaseBadge :variant="value === 'ready' ? 'success' : value === 'queued' ? 'warning' : 'muted'" size="md">
            {{ value }}
          </LabBaseBadge>
        </template>
      </LabDataTable>
    </section>

    <section class="space-y-3">
      <h2 class="text-base text-(--lab-text-primary)">tabs</h2>
      <LabNavTabs v-model="tabValue" :items="tabItems">
        <template #panel-one>
          <p class="text-sm text-(--lab-text-muted)">one panel</p>
        </template>
        <template #panel-two>
          <p class="text-sm text-(--lab-text-muted)">two panel</p>
        </template>
        <template #panel-three>
          <p class="text-sm text-(--lab-text-muted)">three panel</p>
        </template>
      </LabNavTabs>
    </section>

    <section class="space-y-3">
      <h2 class="text-base text-(--lab-text-primary)">scroll hidden + fade</h2>
      <div class="relative max-w-full overflow-hidden border border-(--lab-border)">
        <div class="lab-scroll-hidden flex gap-3 overflow-x-auto px-6 py-3 whitespace-nowrap">
          <span v-for="item in scrollItems" :key="item" class="shrink-0 text-sm text-(--lab-text-secondary)">
            {{ item }}
          </span>
        </div>
        <div class="lab-scroll-fade lab-scroll-fade-visible lab-scroll-fade-x-left" aria-hidden="true" />
        <div class="lab-scroll-fade lab-scroll-fade-visible lab-scroll-fade-x-right" aria-hidden="true" />
      </div>
    </section>
  </div>
</template>
