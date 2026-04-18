<script setup lang="ts">
usePageSeo({
  title: 'UI main.css preview',
  description: 'Preview of global UI classes from main.css.'
})
const text = ref('')
const disabled = ref(false)
const showFade = ref(true)
const tabs = ['one', 'two', 'three'] as const
const activeTab = ref<(typeof tabs)[number]>('one')
const badgeItems = [
  { cls: 'lab-badge-default', label: 'default' },
  { cls: 'lab-badge-muted', label: 'muted' },
  { cls: 'lab-badge-info', label: 'info' },
  { cls: 'lab-badge-success', label: 'success' },
  { cls: 'lab-badge-warning', label: 'warning' },
  { cls: 'lab-badge-danger', label: 'danger' }
]
</script>

<template>
  <div class="space-y-8 p-4 text-(--lab-text-secondary)">
    <section class="space-y-3">
      <h1 class="text-lg text-(--lab-text-primary)">main.css preview</h1>
      <p class="lab-field-label">Preview of global focus, control, dropdown, button, badge, grid and scroll classes.</p>
    </section>
    <section class="space-y-3">
      <h2 class="text-base text-(--lab-text-primary)">focus</h2>
      <div class="grid gap-3 min-[900px]:grid-cols-2">
        <button class="lab-button lab-button-secondary lab-focus min-h-11 border px-3" type="button">
          lab-focus button
        </button>
        <div class="lab-control lab-focus-manual">lab-focus-manual</div>
        <input v-model="text" class="lab-control lab-focus" placeholder="lab-control + lab-focus" />
        <input class="lab-control lab-focus lab-control-invalid" placeholder="lab-control-invalid" />
        <label class="flex min-h-11 items-center border border-(--lab-border) px-3">
          <input class="peer sr-only" type="checkbox" />
          <span class="lab-focus-peer flex min-h-11 w-full items-center px-3">peer + lab-focus-peer</span>
        </label>
        <button class="lab-button lab-tabs-focus min-h-11 border px-3" type="button">lab-tabs-focus</button>
      </div>
    </section>

    <section class="space-y-3">
      <h2 class="text-base text-(--lab-text-primary)">controls</h2>
      <div class="grid gap-3 min-[900px]:grid-cols-2">
        <input class="lab-control" placeholder="default control" />
        <input class="lab-control" :disabled="disabled" placeholder="disabled control" />
        <textarea class="lab-control min-h-24" placeholder="textarea control" />
        <select class="lab-control">
          <option>option one</option>
          <option>option two</option>
        </select>
      </div>
      <label class="inline-flex items-center gap-2 text-sm">
        <input v-model="disabled" type="checkbox" />
        disabled
      </label>
    </section>

    <section class="space-y-3">
      <h2 class="text-base text-(--lab-text-primary)">dropdown</h2>
      <div class="lab-dropdown-panel max-h-56 max-w-xl">
        <button class="lab-dropdown-option" type="button">
          <span>default option</span>
          <span class="lab-dropdown-option-meta">meta</span>
        </button>
        <button class="lab-dropdown-option lab-dropdown-option-active" type="button">
          <span>active option</span>
          <span class="lab-dropdown-option-meta">meta</span>
        </button>
        <button class="lab-dropdown-option" type="button" disabled>
          <span>disabled option</span>
          <span class="lab-dropdown-option-meta">meta</span>
        </button>
      </div>
    </section>

    <section class="space-y-3">
      <h2 class="text-base text-(--lab-text-primary)">buttons</h2>
      <div class="flex flex-wrap gap-3">
        <button class="lab-button lab-button-primary min-h-11 border px-3">primary</button>
        <button class="lab-button lab-button-secondary min-h-11 border px-3">secondary</button>
        <button class="lab-button lab-button-danger min-h-11 border px-3">danger</button>
        <button class="lab-button lab-button-success min-h-11 border px-3">success</button>
        <button class="lab-button lab-button-info min-h-11 border px-3">info</button>
        <button class="lab-button lab-button-ghost min-h-11 px-3">ghost</button>
        <button class="lab-button lab-button-plain min-h-11 px-3">plain</button>
        <button class="lab-button lab-button-secondary min-h-11 border px-3" disabled>disabled</button>
      </div>
    </section>

    <section class="space-y-3">
      <h2 class="text-base text-(--lab-text-primary)">badges</h2>
      <div class="flex flex-wrap gap-2">
        <span v-for="item in badgeItems" :key="item.cls" class="lab-badge border px-2 py-1 text-sm" :class="item.cls">
          {{ item.label }}
        </span>
      </div>
    </section>

    <section class="space-y-3">
      <h2 class="text-base text-(--lab-text-primary)">grid table</h2>
      <div class="lab-grid-table max-w-xl grid-cols-3">
        <div class="lab-grid-table-cell p-3">cell 1</div>
        <div class="lab-grid-table-cell p-3">cell 2</div>
        <div class="lab-grid-table-cell p-3">cell 3</div>
        <div class="lab-grid-table-cell p-3">cell 4</div>
        <div class="lab-grid-table-cell p-3">cell 5</div>
        <div class="lab-grid-table-cell p-3">cell 6</div>
      </div>
    </section>

    <section class="space-y-3">
      <h2 class="text-base text-(--lab-text-primary)">scroll hidden + fade</h2>
      <div class="relative max-w-full overflow-hidden border border-(--lab-border)">
        <div class="lab-scroll-hidden flex gap-2 overflow-x-auto px-6 py-3 whitespace-nowrap">
          <button
            v-for="n in 16"
            :key="n"
            class="lab-button lab-button-secondary min-h-10 shrink-0 border px-3"
            type="button"
          >
            item {{ n }}
          </button>
        </div>
        <div class="lab-scroll-fade lab-scroll-fade-visible lab-scroll-fade-x-left" :class="showFade ? '' : 'hidden'" />
        <div
          class="lab-scroll-fade lab-scroll-fade-visible lab-scroll-fade-x-right"
          :class="showFade ? '' : 'hidden'"
        />
      </div>
      <label class="inline-flex items-center gap-2 text-sm">
        <input v-model="showFade" type="checkbox" />
        show fade
      </label>
    </section>

    <section class="space-y-3">
      <h2 class="text-base text-(--lab-text-primary)">tabs focus preview</h2>
      <div class="flex flex-wrap gap-2 border-b border-(--lab-border)">
        <button
          v-for="tab in tabs"
          :key="tab"
          class="lab-tabs-focus min-h-10 border-b-2 border-transparent px-3 text-sm"
          :class="activeTab === tab ? 'text-(--lab-text-primary)' : 'text-(--lab-text-secondary)'"
          type="button"
          @click="activeTab = tab"
        >
          {{ tab }}
        </button>
      </div>
    </section>
  </div>
</template>
