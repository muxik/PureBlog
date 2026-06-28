/**
 * Theme utilities — shared contract with the blog via localStorage "muxi:theme".
 * Mirrors the design's admin.js theme mechanics exactly.
 */
import { computed, onMounted, onUnmounted, ref, type ComputedRef, type Ref } from 'vue'

const THEME_KEY = 'muxi:theme'

/** Returns true if the current effective theme is dark. */
export function effectiveDark(): boolean {
  const attr = document.documentElement.getAttribute('data-theme')
  if (attr === 'dark') return true
  if (attr === 'light') return false
  return !!(window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches)
}

/** Flips the theme, writes to data-theme attr and localStorage. */
export function toggleTheme(): void {
  const next = effectiveDark() ? 'light' : 'dark'
  document.documentElement.setAttribute('data-theme', next)
  try { localStorage.setItem(THEME_KEY, next) } catch (_) {}
}

/**
 * Reads the stored theme and applies it.
 * Called once at mount time (the head script also does this before first paint).
 */
export function applyStoredTheme(): void {
  try {
    const t = localStorage.getItem(THEME_KEY)
    if (t) document.documentElement.setAttribute('data-theme', t)
  } catch (_) {}
}

export interface ThemeComposable {
  isDark: Ref<boolean>
  toggle(): void
  /** The OPPOSITE label to display on the toggle button: dark→"浅色", light→"深色" */
  label: ComputedRef<string>
}

/** Reactive theme composable for use in Vue components. */
export function useTheme(): ThemeComposable {
  const isDark = ref(effectiveDark())

  function syncFromDOM() {
    isDark.value = effectiveDark()
  }

  function toggle() {
    toggleTheme()
    isDark.value = effectiveDark()
  }

  // Keep reactive when changed from another tab
  function onStorage(e: StorageEvent) {
    if (e.key === THEME_KEY) syncFromDOM()
  }

  onMounted(() => {
    window.addEventListener('storage', onStorage)
  })

  onUnmounted(() => {
    window.removeEventListener('storage', onStorage)
  })

  // The label shows the OPPOSITE of the current theme (what clicking will switch TO)
  const label = computed<string>(() => (isDark.value ? '浅色' : '深色'))

  return { isDark, toggle, label }
}
