<script setup lang="ts">
  definePageMeta({
    title: 'конвертер единиц измерения',
    description:
      'Конвертер единиц измерения с категориями, поиском, быстрым обменом направлений и архитектурой для дальнейшего расширения.'
  })
  const unitCatalog = {
    length: {
      key: 'length',
      kind: 'linear',
      title: 'Длина',
      description: 'Расстояния, размеры, высота, ширина и геометрические измерения.',
      baseUnitKey: 'meter',
      units: [
        { key: 'millimeter', label: 'Миллиметр', symbol: 'мм', aliases: ['mm'], toBase: 0.001 },
        { key: 'centimeter', label: 'Сантиметр', symbol: 'см', aliases: ['cm'], toBase: 0.01 },
        { key: 'meter', label: 'Метр', symbol: 'м', aliases: ['m'], toBase: 1 },
        { key: 'kilometer', label: 'Километр', symbol: 'км', aliases: ['km'], toBase: 1000 },
        { key: 'inch', label: 'Дюйм', symbol: 'in', aliases: ['inch'], toBase: 0.0254 },
        { key: 'foot', label: 'Фут', symbol: 'ft', aliases: ['foot'], toBase: 0.3048 },
        { key: 'yard', label: 'Ярд', symbol: 'yd', aliases: ['yard'], toBase: 0.9144 },
        { key: 'mile', label: 'Миля', symbol: 'mi', aliases: ['mile'], toBase: 1609.344 }
      ]
    },
    mass: {
      key: 'mass',
      kind: 'linear',
      title: 'Масса',
      description: 'Масса и вес для быта, науки и промышленности.',
      baseUnitKey: 'kilogram',
      units: [
        { key: 'milligram', label: 'Миллиграмм', symbol: 'мг', aliases: ['mg'], toBase: 0.000001 },
        { key: 'gram', label: 'Грамм', symbol: 'г', aliases: ['g'], toBase: 0.001 },
        { key: 'kilogram', label: 'Килограмм', symbol: 'кг', aliases: ['kg'], toBase: 1 },
        { key: 'tonne', label: 'Тонна', symbol: 'т', aliases: ['t'], toBase: 1000 },
        { key: 'ounce', label: 'Унция', symbol: 'oz', aliases: ['ounce'], toBase: 0.028349523125 },
        { key: 'pound', label: 'Фунт', symbol: 'lb', aliases: ['pound'], toBase: 0.45359237 }
      ]
    },
    temperature: {
      key: 'temperature',
      kind: 'temperature',
      title: 'Температура',
      description: 'Температурные шкалы с отдельными формулами преобразования.',
      baseUnitKey: 'kelvin',
      units: [
        {
          key: 'celsius',
          label: 'Цельсий',
          symbol: '°C',
          aliases: ['c'],
          toBase: value => value + 273.15,
          fromBase: value => value - 273.15
        },
        {
          key: 'fahrenheit',
          label: 'Фаренгейт',
          symbol: '°F',
          aliases: ['f'],
          toBase: value => ((value - 32) * 5) / 9 + 273.15,
          fromBase: value => ((value - 273.15) * 9) / 5 + 32
        },
        {
          key: 'kelvin',
          label: 'Кельвин',
          symbol: 'K',
          aliases: ['k'],
          toBase: value => value,
          fromBase: value => value
        }
      ]
    },
    area: {
      key: 'area',
      kind: 'linear',
      title: 'Площадь',
      description: 'Площадь поверхностей, участков и помещений.',
      baseUnitKey: 'square-meter',
      units: [
        { key: 'square-centimeter', label: 'Квадратный сантиметр', symbol: 'см²', aliases: ['cm2'], toBase: 0.0001 },
        { key: 'square-meter', label: 'Квадратный метр', symbol: 'м²', aliases: ['m2'], toBase: 1 },
        { key: 'square-kilometer', label: 'Квадратный километр', symbol: 'км²', aliases: ['km2'], toBase: 1000000 },
        { key: 'hectare', label: 'Гектар', symbol: 'га', aliases: ['ha'], toBase: 10000 },
        { key: 'acre', label: 'Акр', symbol: 'ac', aliases: ['acre'], toBase: 4046.8564224 }
      ]
    },
    volume: {
      key: 'volume',
      kind: 'linear',
      title: 'Объём',
      description: 'Объёмы жидкостей, газов и сыпучих веществ.',
      baseUnitKey: 'liter',
      units: [
        { key: 'milliliter', label: 'Миллилитр', symbol: 'мл', aliases: ['ml'], toBase: 0.001 },
        { key: 'liter', label: 'Литр', symbol: 'л', aliases: ['l'], toBase: 1 },
        { key: 'cubic-meter', label: 'Кубический метр', symbol: 'м³', aliases: ['m3'], toBase: 1000 },
        { key: 'teaspoon', label: 'Чайная ложка', symbol: 'tsp', aliases: ['teaspoon'], toBase: 0.00492892159375 },
        {
          key: 'tablespoon',
          label: 'Столовая ложка',
          symbol: 'tbsp',
          aliases: ['tablespoon'],
          toBase: 0.01478676478125
        },
        { key: 'cup', label: 'Чашка', symbol: 'cup', aliases: ['cup'], toBase: 0.2365882365 },
        { key: 'gallon-us', label: 'Галлон США', symbol: 'gal', aliases: ['gallon'], toBase: 3.785411784 }
      ]
    },
    speed: {
      key: 'speed',
      kind: 'linear',
      title: 'Скорость',
      description: 'Скорость движения транспорта, потоков и процессов.',
      baseUnitKey: 'meter-per-second',
      units: [
        { key: 'meter-per-second', label: 'Метр в секунду', symbol: 'м/с', aliases: ['m/s'], toBase: 1 },
        {
          key: 'kilometer-per-hour',
          label: 'Километр в час',
          symbol: 'км/ч',
          aliases: ['km/h'],
          toBase: 0.2777777777778
        },
        { key: 'mile-per-hour', label: 'Миля в час', symbol: 'mph', aliases: ['mile/hour'], toBase: 0.44704 },
        { key: 'knot', label: 'Узел', symbol: 'kn', aliases: ['knot'], toBase: 0.5144444444444 }
      ]
    },
    time: {
      key: 'time',
      kind: 'linear',
      title: 'Время',
      description: 'Продолжительность, интервалы и временные масштабы.',
      baseUnitKey: 'second',
      units: [
        { key: 'millisecond', label: 'Миллисекунда', symbol: 'мс', aliases: ['ms'], toBase: 0.001 },
        { key: 'second', label: 'Секунда', symbol: 'с', aliases: ['sec'], toBase: 1 },
        { key: 'minute', label: 'Минута', symbol: 'мин', aliases: ['min'], toBase: 60 },
        { key: 'hour', label: 'Час', symbol: 'ч', aliases: ['h'], toBase: 3600 },
        { key: 'day', label: 'Сутки', symbol: 'сут', aliases: ['day'], toBase: 86400 }
      ]
    },
    pressure: {
      key: 'pressure',
      kind: 'linear',
      title: 'Давление',
      description: 'Давление в технике, метеорологии и физике.',
      baseUnitKey: 'pascal',
      units: [
        { key: 'pascal', label: 'Паскаль', symbol: 'Па', aliases: ['pa'], toBase: 1 },
        { key: 'kilopascal', label: 'Килопаскаль', symbol: 'кПа', aliases: ['kpa'], toBase: 1000 },
        { key: 'bar', label: 'Бар', symbol: 'bar', aliases: ['bar'], toBase: 100000 },
        { key: 'atmosphere', label: 'Атмосфера', symbol: 'atm', aliases: ['atm'], toBase: 101325 },
        {
          key: 'mmhg',
          label: 'Миллиметр ртутного столба',
          symbol: 'мм рт. ст.',
          aliases: ['mmhg'],
          toBase: 133.322387415
        }
      ]
    },
    energy: {
      key: 'energy',
      kind: 'linear',
      title: 'Энергия',
      description: 'Энергия в физике, быту и электротехнике.',
      baseUnitKey: 'joule',
      units: [
        { key: 'joule', label: 'Джоуль', symbol: 'Дж', aliases: ['j'], toBase: 1 },
        { key: 'kilojoule', label: 'Килоджоуль', symbol: 'кДж', aliases: ['kj'], toBase: 1000 },
        { key: 'calorie', label: 'Калория', symbol: 'cal', aliases: ['cal'], toBase: 4.184 },
        { key: 'kilowatt-hour', label: 'Киловатт-час', symbol: 'кВт·ч', aliases: ['kwh'], toBase: 3600000 }
      ]
    },
    power: {
      key: 'power',
      kind: 'linear',
      title: 'Мощность',
      description: 'Мощность механизмов, двигателей и энергосистем.',
      baseUnitKey: 'watt',
      units: [
        { key: 'watt', label: 'Ватт', symbol: 'Вт', aliases: ['w'], toBase: 1 },
        { key: 'kilowatt', label: 'Киловатт', symbol: 'кВт', aliases: ['kw'], toBase: 1000 },
        { key: 'megawatt', label: 'Мегаватт', symbol: 'МВт', aliases: ['mw'], toBase: 1000000 },
        { key: 'horsepower-metric', label: 'Лошадиная сила', symbol: 'л.с.', aliases: ['hp'], toBase: 735.49875 }
      ]
    }
  } satisfies Record<MathUnitCategoryKey, MathUnitCategoryDef>
  const categoryOrder = Object.keys(unitCatalog) as MathUnitCategoryKey[]
  const selectedCategoryKey = ref<MathUnitCategoryKey>('length')
  const amountInput = ref('1')
  const searchQuery = ref('')
  const fromUnitKey = ref('meter')
  const toUnitKey = ref('kilometer')
  const decimalPlaces = ref(6)
  const showOnlyCommon = ref(false)
  const commonUnitKeys = new Set([
    'millimeter',
    'centimeter',
    'meter',
    'kilometer',
    'gram',
    'kilogram',
    'celsius',
    'fahrenheit',
    'kelvin',
    'liter',
    'milliliter',
    'kilometer-per-hour',
    'meter-per-second',
    'second',
    'minute',
    'hour',
    'pascal',
    'bar',
    'atmosphere',
    'joule',
    'kilowatt-hour',
    'watt',
    'kilowatt'
  ])
  const selectedCategory = computed(() => unitCatalog[selectedCategoryKey.value])
  const categoryUnits = computed(() => {
    const units = selectedCategory.value.units
    return units.filter(unit => {
      if (showOnlyCommon.value && !commonUnitKeys.has(unit.key)) return false
      if (!searchQuery.value.trim()) return true
      const needle = searchQuery.value.trim().toLowerCase()
      const haystack = [unit.label, unit.symbol, ...(unit.aliases || [])].join(' ').toLowerCase()
      return haystack.includes(needle)
    })
  })
  const fromUnit = computed(() => {
    const units = selectedCategory.value.units
    return units.find(unit => unit.key === fromUnitKey.value) ?? units[0] ?? null
  })
  const toUnit = computed(() => {
    const units = selectedCategory.value.units
    return units.find(unit => unit.key === toUnitKey.value) ?? units[1] ?? units[0] ?? null
  })
  watch(
    selectedCategoryKey,
    nextCategoryKey => {
      const nextCategory = unitCatalog[nextCategoryKey]
      fromUnitKey.value = nextCategory.units[0]?.key || ''
      toUnitKey.value = nextCategory.units[1]?.key || nextCategory.units[0]?.key || ''
      searchQuery.value = ''
    },
    { immediate: true }
  )
  const parsedAmount = computed(() => {
    const normalized = amountInput.value.replace(',', '.').trim()
    if (!normalized) return null
    const parsed = Number(normalized)
    return Number.isFinite(parsed) ? parsed : null
  })
  const convertLinearValue = (category: MathLinearCategoryDef, value: number, sourceKey: string, targetKey: string) => {
    const sourceUnit = category.units.find(unit => unit.key === sourceKey)
    const targetUnit = category.units.find(unit => unit.key === targetKey)
    if (!sourceUnit || !targetUnit) return null
    const baseValue = value * sourceUnit.toBase
    return baseValue / targetUnit.toBase
  }
  const convertTemperatureValue = (
    category: MathTemperatureCategoryDef,
    value: number,
    sourceKey: string,
    targetKey: string
  ) => {
    const sourceUnit = category.units.find(unit => unit.key === sourceKey)
    const targetUnit = category.units.find(unit => unit.key === targetKey)
    if (!sourceUnit || !targetUnit) return null
    const baseValue = sourceUnit.toBase(value)
    return targetUnit.fromBase(baseValue)
  }
  const convertValue = (category: MathUnitCategoryDef, value: number, sourceKey: string, targetKey: string) => {
    if (category.kind === 'linear') {
      return convertLinearValue(category, value, sourceKey, targetKey)
    }
    return convertTemperatureValue(category, value, sourceKey, targetKey)
  }
  const convertedValue = computed(() => {
    if (parsedAmount.value === null) return null
    return convertValue(selectedCategory.value, parsedAmount.value, fromUnitKey.value, toUnitKey.value)
  })
  const formatNumber = (value: number | null) => {
    if (value === null || Number.isNaN(value)) return '—'
    const normalized = Number(value.toFixed(decimalPlaces.value))
    return new Intl.NumberFormat('ru-RU', {
      maximumFractionDigits: decimalPlaces.value
    }).format(normalized)
  }
  const resultSummary = computed(() => {
    if (parsedAmount.value === null || convertedValue.value === null) {
      return 'Введите корректное число.'
    }
    if (!fromUnit.value || !toUnit.value) {
      return 'Единицы измерения недоступны.'
    }
    return `${formatNumber(parsedAmount.value)} ${fromUnit.value.symbol} = ${formatNumber(convertedValue.value)} ${toUnit.value.symbol}`
  })
  const quickTableRows = computed(() => {
    if (parsedAmount.value === null) return []
    return selectedCategory.value.units
      .filter(unit => unit.key !== fromUnitKey.value)
      .slice(0, 8)
      .map(unit => ({
        key: unit.key,
        label: unit.label,
        symbol: unit.symbol,
        value: convertValue(selectedCategory.value, parsedAmount.value as number, fromUnitKey.value, unit.key)
      }))
  })
  const swapUnits = () => {
    const currentFrom = fromUnitKey.value
    fromUnitKey.value = toUnitKey.value
    toUnitKey.value = currentFrom
  }
  const setCategory = (key: MathUnitCategoryKey) => {
    selectedCategoryKey.value = key
  }
  const applyPreset = (categoryKey: MathUnitCategoryKey, sourceKey: string, targetKey: string, value: string) => {
    selectedCategoryKey.value = categoryKey
    amountInput.value = value
    fromUnitKey.value = sourceKey
    toUnitKey.value = targetKey
  }
</script>
<template>
  <div class="min-h-0 overflow-hidden">
    <LabNavHeader
      title="Конвертер единиц измерения"
      :breadcrumb-items="[
        { label: 'Вики', to: '/edu' },
        { label: 'Математика', to: '/edu/math' },
        { label: 'Конвертер единиц', current: true }
      ]" />
    <section class="min-h-0 flex flex-wrap gap-3 overflow-hidden p-3">
      <aside class="min-h-0 overflow-hidden border border-white/10 bg-black/20">
        <div class="border-b border-white/10 px-3 py-2">
          <h3 class="text-sm font-semibold text-white">Категории</h3>
        </div>
        <ul class="h-full divide-y divide-zinc-800 overflow-auto">
          <li v-for="categoryKey in categoryOrder" :key="categoryKey">
            <button
              type="button"
              class="block w-full px-3 py-2 text-left text-sm transition"
              :class="
                selectedCategoryKey === categoryKey ? 'bg-white/10 text-white' : 'text-zinc-300 hover:bg-zinc-800'
              "
              @click="setCategory(categoryKey)">
              {{ unitCatalog[categoryKey].title }}
            </button>
          </li>
        </ul>
      </aside>
      <section class="grid min-h-0 gap-3 grid-rows-[auto_auto]">
        <div class="border border-white/10 bg-black/20 p-3">
          <div class="mb-1 text-[10px] uppercase tracking-[0.14em] text-zinc-500">Значение</div>
          <LabBaseInput v-model="amountInput" type="text" placeholder="Введите число" inputmode="decimal" />
        </div>
        <div class="border border-white/10 bg-black/20 p-3">
          <div class="mb-3 text-[10px] uppercase tracking-[0.14em] text-zinc-500">Выбор единиц</div>
          <div class="grid gap-3 grid-cols-[minmax(0,1fr)_42px_minmax(0,1fr)]">
            <div>
              <div class="mb-1 text-[10px] uppercase tracking-[0.14em] text-zinc-500">Из</div>
              <select
                v-model="fromUnitKey"
                class="h-9 w-full border border-white/10 bg-black/20 px-3 text-sm text-white outline-none transition focus:border-white/20">
                <option v-for="unit in selectedCategory.units" :key="unit.key" :value="unit.key">
                  {{ unit.label }} · {{ unit.symbol }}
                </option>
              </select>
            </div>
            <button
              type="button"
              class="mt-4.5 h-9 border border-white/10 bg-black/20 text-sm text-zinc-200 transition hover:bg-white/5"
              @click="swapUnits">
              ⇄
            </button>
            <div>
              <div class="mb-1 text-[10px] uppercase tracking-[0.14em] text-zinc-500">В</div>
              <select
                v-model="toUnitKey"
                class="h-9 w-full border border-white/10 bg-black/20 px-3 text-sm text-white outline-none transition focus:border-white/20">
                <option v-for="unit in selectedCategory.units" :key="unit.key" :value="unit.key">
                  {{ unit.label }} · {{ unit.symbol }}
                </option>
              </select>
            </div>
          </div>
        </div>
      </section>
      <section class="grid min-h-0 gap-3 grid-rows-[auto_minmax(0,1fr)]">
        <div class="border border-white/10 bg-black/20 p-3">
          <div class="mb-1 text-[10px] uppercase tracking-[0.14em] text-zinc-500">Результат</div>
          <div class="text-base font-semibold text-white">
            {{ resultSummary }}
          </div>
        </div>
        <section class="flex min-h-0 flex-col border border-white/10 bg-black/20">
          <div class="border-b border-white/10 px-3 py-2">
            <h3 class="text-sm font-semibold text-white">Быстрый пересчёт</h3>
          </div>
          <div class="min-h-0 flex-1 overflow-auto">
            <table class="min-w-full divide-y divide-white/10 text-sm">
              <thead class="sticky top-0 bg-zinc-900 text-zinc-400">
                <tr>
                  <th class="px-3 py-2 text-left font-medium">Единица</th>
                  <th class="px-3 py-2 text-right font-medium">Значение</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-white/10">
                <tr v-for="row in quickTableRows" :key="row.key" class="text-zinc-200">
                  <td class="px-3 py-2">
                    {{ row.label }}
                    <span class="text-zinc-500">{{ row.symbol }}</span>
                  </td>
                  <td class="px-3 py-2 text-right">{{ formatNumber(row.value) }} {{ row.symbol }}</td>
                </tr>
                <tr v-if="!quickTableRows.length">
                  <td colspan="2" class="px-3 py-4 text-sm text-zinc-500">Введите число.</td>
                </tr>
              </tbody>
            </table>
          </div>
        </section>
      </section>
    </section>
  </div>
</template>
