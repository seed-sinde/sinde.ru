<script setup lang="ts">
import {AmbientLight, Box3, DirectionalLight, Group, PerspectiveCamera, Scene, Vector3, WebGLRenderer, type Material, type Object3D} from 'three'
import {GLTFLoader, type GLTF} from 'three/examples/jsm/loaders/GLTFLoader.js'
defineOptions({ inheritAttrs: false })
const props = withDefaults(defineProps<{ src?: string | null; title?: string; downloadName?: string }>(), {
  src: '',
  title: '3D model',
  downloadName: 'model.glb'
})
const rootRef = ref<HTMLDivElement | null>(null)
const canvasRef = ref<HTMLCanvasElement | null>(null)
const loading = ref(false)
const error = ref('')
const isFullscreen = ref(false)
let renderer: WebGLRenderer | null = null
let scene: Scene | null = null
let camera: PerspectiveCamera | null = null
let resizeObserver: ResizeObserver | null = null
let disposeModel = () => {}
let frameId = 0
const drag = { active: false, x: 0, y: 0 }
const rotation = reactive({ x: -0.35, y: 0.55 })
const zoom = ref(2.4)
const src = computed(() => String(props.src || '').trim())
const rootClass = computed(() => [
  'overflow-hidden bg-(--lab-bg-canvas) select-none',
  isFullscreen.value ? 'fixed inset-0 z-60 h-screen w-screen' : 'relative h-full w-full'
])
const toolbarClass = computed(() => [
  'absolute top-2 right-2 z-10 flex items-center gap-1',
  isFullscreen.value && 'top-3 right-3'
])
const toggleFullscreen = () => {
  isFullscreen.value = !isFullscreen.value
}
const closeFullscreen = () => {
  isFullscreen.value = false
}
const onKeydown = (event: KeyboardEvent) => {
  if (event.key === 'Escape') closeFullscreen()
}
const animate = () => {
  if (!renderer || !scene || !camera) return
  if (!drag.active) {
    rotation.y += 0.004
    render()
  }
  frameId = requestAnimationFrame(animate)
}
const render = () => {
  if (!renderer || !scene || !camera) return
  camera.position.z = zoom.value
  scene.rotation.x = rotation.x
  scene.rotation.y = rotation.y
  renderer.render(scene, camera)
}
const resize = () => {
  if (!renderer || !camera || !rootRef.value) return
  const { width, height } = rootRef.value.getBoundingClientRect()
  if (!width || !height) return
  renderer.setSize(width, height, false)
  camera.aspect = width / height
  camera.updateProjectionMatrix()
  render()
}
const stop = () => {
  if (frameId) cancelAnimationFrame(frameId)
  frameId = 0
  resizeObserver?.disconnect()
  resizeObserver = null
  disposeModel()
  disposeModel = () => {}
  renderer?.dispose()
  renderer = null
  scene = null
  camera = null
}
const mountModel = async () => {
  stop()
  error.value = ''
  if (!import.meta.client || !rootRef.value || !canvasRef.value || !src.value) return
  loading.value = true
  const [{Mesh, MeshStandardMaterial}] = await Promise.all([
    import('three')
  ])
  scene = new Scene()
  camera = new PerspectiveCamera(45, 1, 0.01, 100)
  renderer = new WebGLRenderer({ antialias: true, alpha: true, canvas: canvasRef.value, powerPreference: 'low-power' })
  renderer.setPixelRatio(Math.min(window.devicePixelRatio || 1, 2))
  scene.add(new AmbientLight(0xffffff, 2.6))
  const light = new DirectionalLight(0xffffff, 2.4)
  light.position.set(3, 4, 6)
  scene.add(light)
  const pivot = new Group()
  scene.add(pivot)
  try {
    const gltf = await new Promise<GLTF>((resolve, reject) => {
      new GLTFLoader().load(src.value, resolve, undefined, reject)
    })
    const model = gltf.scene
    const box = new Box3().setFromObject(model)
    const size = box.getSize(new Vector3())
    const center = box.getCenter(new Vector3())
    const scale = 1.8 / Math.max(size.x || 0, size.y || 0, size.z || 0, 1)
    model.position.set(-center.x * scale, -center.y * scale, -center.z * scale)
    model.scale.setScalar(scale)
    model.traverse((node: Object3D) => {
      if (!(node instanceof Mesh) || Array.isArray(node.material)) return
      node.castShadow = false
      node.receiveShadow = false
      node.material = node.material instanceof MeshStandardMaterial ? node.material : new MeshStandardMaterial({ color: 0xd4d4d8 })
    })
    pivot.add(model)
    disposeModel = () => {
      pivot.remove(model)
      model.traverse((node: Object3D) => {
        if (!(node instanceof Mesh)) return
        node.geometry.dispose()
        ;(Array.isArray(node.material) ? node.material : [node.material]).forEach((material: Material) => material.dispose())
      })
    }
    resizeObserver = new ResizeObserver(resize)
    resizeObserver.observe(rootRef.value)
    resize()
    frameId = requestAnimationFrame(animate)
  } catch (cause) {
    error.value = cause instanceof Error ? cause.message : 'Не удалось загрузить модель'
    stop()
  } finally {
    loading.value = false
  }
}
const onPointerDown = (event: PointerEvent) => {
  drag.active = true
  drag.x = event.clientX
  drag.y = event.clientY
}
const onPointerMove = (event: PointerEvent) => {
  if (!drag.active) return
  rotation.y += (event.clientX - drag.x) * 0.01
  rotation.x = Math.max(-1.4, Math.min(1.4, rotation.x + (event.clientY - drag.y) * 0.01))
  drag.x = event.clientX
  drag.y = event.clientY
  render()
}
const onPointerUp = () => {
  drag.active = false
}
const onWheel = (event: WheelEvent) => {
  zoom.value = Math.max(1.2, Math.min(6, zoom.value + event.deltaY * 0.002))
  render()
}
watch(src, () => void mountModel(), { immediate: true })
watch(isFullscreen, () => nextTick(resize))
onMounted(() => {
  window.addEventListener('keydown', onKeydown)
  void mountModel()
})
onBeforeUnmount(() => {
  window.removeEventListener('keydown', onKeydown)
  stop()
})
</script>
<template>
  <div
    v-bind="$attrs"
    ref="rootRef"
    :class="rootClass"
    @pointerdown="onPointerDown"
    @pointermove="onPointerMove"
    @pointerup="onPointerUp"
    @pointerleave="onPointerUp"
    @wheel.prevent="onWheel">
    <canvas ref="canvasRef" class="h-full w-full touch-none" :aria-label="title" />
    <div v-if="src" :class="toolbarClass" @pointerdown.stop>
      <button
        v-if="isFullscreen"
        type="button"
        class="lab-focus inline-flex h-10 w-10 items-center justify-center bg-(--lab-bg-control) text-(--lab-text-primary)"
        :aria-label="`Закрыть ${title}`"
        @click.stop="closeFullscreen"
      >
        <Icon name="ic:round-close" class="text-xl" />
      </button>
      <button
        v-else
        type="button"
        class="lab-focus inline-flex h-8 w-8 items-center justify-center text-(--lab-text-primary)"
        :aria-label="`Открыть ${title} на весь экран`"
        @click.stop="toggleFullscreen"
      >
        <Icon name="ic:round-fullscreen" class="text-base" />
      </button>
      <a
        class="lab-focus inline-flex items-center justify-center text-(--lab-text-primary)"
        :class="isFullscreen ? 'h-10 w-10 bg-(--lab-bg-control)' : 'h-8 w-8'"
        :href="src"
        :download="downloadName"
        :aria-label="`Скачать ${title}`"
        @click.stop
      >
        <Icon name="ic:round-download" class="text-base" />
      </a>
    </div>
    <div v-if="loading || error || !src" class="absolute inset-0 grid place-items-center p-4 text-center text-sm text-(--lab-text-muted)">
      {{ !src ? 'Нет модели' : loading ? 'Загрузка модели…' : error }}
    </div>
  </div>
</template>
