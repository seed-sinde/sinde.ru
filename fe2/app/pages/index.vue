<script setup lang="ts">
import type {CheckboxValue} from "~/components/ui/Checkbox.vue"
import IcBaselineArrowRight from "~icons/ic/baseline-arrow-right"
definePageMeta({tabs: true})
const details = ref(false)
const details_2 = ref(false)
const sinde = `SINDE — System for Integrated Networked Data Exchange`
const email = ref("")
const hasEmailError = ref(true)
const hasPassError = ref(true)
const email_error = ref("Ошибка ввода почты")
const password = ref("")
const password_error = ref("Нет доступа")
const submitting = ref(false)
const toggleState = ref(false)
const disabledToggleState = ref(true)
const headerClass = `inline-block border-b-2 border-dotted border-transparent text-(--accent) hover:border-(--border-color)`
const checkbox_1 = ref<CheckboxValue>(false)
const checkbox_2 = ref<CheckboxValue>(true)
const checkbox_3 = ref<CheckboxValue>("partial")
const selectOptions = ref([
  {value: "option 1", label: "Опция 1"},
  {value: "option 2", label: "Опция 2"},
  {value: "option 3", label: "Опция 3 disabled", disabled: true}
])
const scaleValue = ref("middle")
const scaleOptions = [
  {value: "easy", label: "Легко", activeColor: "#16a34a"},
  {value: "medium", label: "Средне", activeColor: "#eab308"},
  {value: "hard", label: "Сложно", activeColor: "#f97316"},
  {value: "insane", label: "Тяжело", activeColor: "#dc2626"}
]
const tabs = [
  {to: "/tab1", label: "Вкладка 1"},
  {to: "/tab2", label: "Вкладка 2"},
  {to: "/tab3", label: "Вкладка 3"},
  {to: "/tab4", label: "Вкладка 4"},
  {to: "/tab5", label: "Вкладка 5"}
]
const result = ref<{email: string; password: string} | null>(null)
const viewerOpen = ref(false)
const viewerIndex = ref(0)
const viewerImages = [
  {
    src: "/pwa-icon-source.svg",
    thumbnailSrc: "/pwa-icon-source.svg",
    title: "PWA icon source",
    alt: "PWA icon source"
  },
  {
    src: "/pwa-maskable-source.svg",
    thumbnailSrc: "/pwa-maskable-source.svg",
    title: "PWA maskable source",
    alt: "PWA maskable source"
  }
]
const validate = () => {
  hasEmailError.value = !/^\S+@\S+\.\S+$/.test(email.value)
  hasPassError.value = password.value.length < 6
  email_error.value = hasEmailError.value ? "Некорректная почта" : ""
  password_error.value = hasPassError.value ? "Минимум 6 символов" : ""
  return !hasEmailError.value && !hasPassError.value
}
const onSubmit = () => {
  result.value = null
  if (!validate()) return
  submitting.value = true
  setTimeout(() => {
    submitting.value = false
    result.value = {
      email: email.value,
      password: password.value
    }
  }, 2000)
}
const onLogout = () => {
  submitting.value = true
  setTimeout(() => {
    submitting.value = false
    result.value = null
    email.value = ""
    password.value = ""
    hasEmailError.value = true
    hasPassError.value = true
    email_error.value = "Ошибка ввода почты"
    password_error.value = "Нет доступа"
  }, 1000)
}
const loading = ref(true)
let timer: ReturnType<typeof setTimeout> | null = null
const stop = () => {
  loading.value = false
  timer = null
}
const reloadLoader = () => {
  if (timer) clearTimeout(timer)
  loading.value = true
  timer = setTimeout(stop, 3000)
}
const openImage = (i = 0) => {
  viewerIndex.value = i
  viewerOpen.value = true
}
onMounted(reloadLoader)
useSeoMeta({
  title: "Главная",
  description: "Платформа прикладны"
})
</script>
<template>
  <UiStack>
    <div class="space-y-2">
      <h2 :class="headerClass" v-text="'<UiTabs />'" />
      <UiTabs :items="tabs">
        <template #panel-tab1="{item}">
          <p class="text-sm leading-snug">
            {{ item.label }}: содержимое первой вкладки через слот.
          </p>
        </template>
        <template #panel-tab2="{item}">
          <p class="text-sm leading-snug">{{ item.label }}: отдельная панель второй вкладки.</p>
        </template>
        <template #panel-tab3="{item}">
          <p class="text-sm leading-snug">{{ item.label }}: можно передавать любой шаблон.</p>
        </template>
        <template #panel-tab4="{item}">
          <p class="text-sm leading-snug">
            {{ item.label }}: слот получает данные активной вкладки.
          </p>
        </template>
        <template #panel-tab5="{item}">
          <p class="text-sm leading-snug">{{ item.label }}: запасной текст</p>
        </template>
      </UiTabs>
    </div>
    <div class="space-y-2">
      <h2 :class="headerClass" v-text="'<UiLoader />'" />
      <div class="flex gap-2"><UiLoader /> Загрузка...</div>
    </div>
    <div class="space-y-2">
      <h2 :class="headerClass" v-text="'<UiSpoiler>'" />
      <h3>Кнопка показывает/скрывает содержимое по клику</h3>
      <div class="grid grid-cols-2 gap-2">
        <div>
          <UiSpoiler v-model="details" label="Блок с иконкой слева">
            <p>
              Иконка меняет свой стиль и текст подсказки, при клике раскрывает содержимое визуально
              скрытого блока.
            </p>
          </UiSpoiler>
        </div>
        <div>
          <UiSpoiler v-model="details_2" label="Блок с иконкой справа" icon-position="right">
            <p>
              Иконка меняет свой стиль и текст подсказки, при клике раскрывает содержимое визуально
              скрытого блока.
            </p>
          </UiSpoiler>
        </div>
      </div>
    </div>
    <div class="space-y-2">
      <h2 :class="headerClass" v-text="'<UiTooltip>'" />
      <h3>Показать всплывающую подсказку на тексте</h3>
      <p>
        Welcome to
        <UiTooltip :text="sinde">
          <span class="border-b-2 border-dotted text-(--accent)">SINDE</span>
        </UiTooltip>
        and Happy New Year!
      </p>
    </div>
    <div class="flex flex-col items-start gap-2">
      <h2 :class="headerClass" v-text="'<UiScale>'" />
      <UiScale
        v-model="scaleValue"
        name="difficulty"
        label="Шкала сложности"
        :options="scaleOptions"
      />
    </div>
    <div class="flex flex-col items-start gap-2">
      <h2 :class="headerClass" v-text="'<UiCheckbox>'" />
      <h3>выбор из списка / множественный выбор / согласие с условиями / "отметить/не отметить"</h3>
      <div class="flex flex-wrap items-start gap-2">
        <UiCheckbox v-model="checkbox_1">false</UiCheckbox>
        <UiCheckbox v-model="checkbox_2">true</UiCheckbox>
        <UiCheckbox v-model="checkbox_3" partial>partial</UiCheckbox>
        <UiCheckbox disabled>disabled</UiCheckbox>
      </div>
    </div>
    <div class="space-y-2">
      <h2 :class="headerClass" v-text="'<UiButton>'" />
      <h3>Несколько примеров использования кнопок с иконками</h3>
      <div class="flex flex-wrap gap-2">
        <UiButton>Variant default</UiButton>
        <UiButton variant="ghost">Variant Ghost</UiButton>
        <UiButton>
          <template #left>
            <IcBaselineArrowRight />
          </template>
          icon left
        </UiButton>
        <UiButton>
          <template #right>
            <IcBaselineArrowRight />
          </template>
          icon right
        </UiButton>
        <UiButton>
          <template #left>
            <IcBaselineArrowRight />
          </template>
          icons both
          <template #right>
            <IcBaselineArrowRight />
          </template>
        </UiButton>
        <UiButton :loading="loading" @click="reloadLoader">Загрузить</UiButton>
        <UiButton disabled>disabled</UiButton>
      </div>
    </div>
    <div class="space-y-2">
      <h2 :class="headerClass" v-text="'<UiToggle>'" />
      <div class="flex flex-wrap items-start gap-2">
        <UiToggle v-model="toggleState">Активный переключатель</UiToggle>
        <UiToggle v-model="disabledToggleState" disabled>Выключенный переключатель</UiToggle>
      </div>
    </div>
    <div class="flex gap-2"></div>
    <div class="flex gap-2">
      <form v-if="!result" class="flex w-fit flex-col gap-1.5" @submit.prevent="onSubmit">
        <UiField label="Почта" hint="Введите почту" :error="email_error" required>
          <UiInput
            v-model="email"
            type="email"
            placeholder="seed@sinde.ru"
            :invalid="hasEmailError"
            autocomplete="username"
          />
        </UiField>
        <UiField
          label="Пароль"
          hint="Введите пароль не менее 12 символов"
          :error="password_error"
          required
        >
          <UiInput
            v-model="password"
            type="password"
            placeholder="Пароль"
            :invalid="hasPassError"
            autocomplete="current-password"
          />
        </UiField>
        <UiButton type="submit" :loading="submitting">Войти</UiButton>
      </form>
      <div v-if="result" class="mt-2 rounded-xl bg-(--elevated) p-2 text-sm">
        <div><b>Email:</b> {{ result.email }}</div>
        <div><b>Password:</b> {{ result.password }}</div>
        <UiButton :loading="submitting" @click="onLogout">Выйти</UiButton>
      </div>
      <div class="flex items-start gap-2">
        <UiTextarea placeholder="Текстовое поле ввода..." />
        <UiSelect :options="selectOptions" />
      </div>
    </div>
    <div class="flex gap-2">
      <UiImage
        v-model="viewerOpen"
        :items="viewerImages"
        :initial-index="viewerIndex"
        @active-index-change="viewerIndex = $event"
      >
        <button
          type="button"
          class="group relative min-h-80 overflow-hidden bg-(--elevated) transition focus:outline-none focus-visible:ring-2 focus-visible:ring-(--accent) md:min-h-120"
          @click="openImage(0)"
        >
          <img
            src="/pwa-icon-source.svg"
            alt="PWA icon source"
            class="h-full w-full object-contain p-8"
          />
          <span
            class="absolute inset-x-0 bottom-0 bg-zinc-950/70 p-3 text-sm font-medium text-white opacity-0 transition-opacity group-hover:opacity-100 group-focus-visible:opacity-100"
          >
            Посмотреть
          </span>
        </button>
      </UiImage>
      <UiThumbnails :items="viewerImages" :active-index="viewerIndex" @select="openImage" />
    </div>
  </UiStack>
</template>
