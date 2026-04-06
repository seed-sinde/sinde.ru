<script setup lang="ts">
  const title = 'Кодекс участника'
  usePageSeo({
    title,
    description: 'Соглашение между участником и платформой sinde.'
  })
  const principles = [
    {
      title: 'Контроль',
      text: 'Ты добавляешь только то, что считаешь уместным. Все действия и данные — твоя зона ответственности.',
      tone: 'text-cyan-400'
    },
    {
      title: 'Использование данных',
      text: 'Данные используются только внутри системы: хранение, отображение, связи и развитие.',
      tone: 'text-emerald-400'
    },
    {
      title: 'Минимизация',
      text: 'Ничего лишнего не собирается — только то, что ты явно вводишь или передаёшь.',
      tone: 'text-amber-400'
    },
    {
      title: 'Актуальность',
      text: 'Данные обновляются или удаляются тобой. Система не хранит лишнего.',
      tone: 'text-violet-400'
    },
    {
      title: 'Доступ',
      text: 'Доступ к данным ограничен логикой системы и твоими действиями.',
      tone: 'text-blue-400'
    },
    {
      title: 'Безопасность',
      text: 'Защита встроена в архитектуру и не требует отдельного контроля.',
      tone: 'text-green-400'
    },
    {
      title: 'Прозрачность',
      text: 'Все изменения системы доступны сразу. Ничего скрытого.',
      tone: 'text-pink-400'
    },
    {
      title: 'Согласие',
      text: 'Использование сервиса означает согласие с принципами: прозрачность и контроль у пользователя.',
      tone: 'text-orange-400'
    }
  ] as const
</script>
<template>
  <div>
    <LabNavHeader :title />
    <section class="space-y-4 p-4">
      <div class="grid gap-4 sm:grid-cols-2 xl:grid-cols-3">
        <article
          v-for="item in principles"
          :key="item.title"
          class="border border-zinc-800 bg-zinc-900/50 p-4 transition ring-1 ring-transparent hover:ring-zinc-700">
          <p class="mb-2 text-sm" :class="item.tone">{{ item.title }}</p>
          <p class="text-sm text-gray-300">{{ item.text }}</p>
        </article>
      </div>
    </section>
  </div>
</template>
