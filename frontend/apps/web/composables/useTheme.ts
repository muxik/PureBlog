/**
 * useTheme — SSR-safe dark/light theme composable.
 *
 * localStorage key : "muxi:theme"
 * values           : "dark" | "light" | (absent = auto/system)
 * data-theme attr  : set on document.documentElement
 *
 * Call once in app.vue to mount the toggle + cross-tab listener.
 * Other components can call it too; the shared `isDark` state is
 * kept in Nuxt's useState so it's reactive across the component tree.
 */

export function useTheme() {
  // Shared reactive state — defaults to false (light) for SSR.
  const isDark = useState<boolean>('theme:isDark', () => false)

  /** Read the actual effective dark state from the DOM + matchMedia. */
  function effectiveDark(): boolean {
    if (!import.meta.client) return false
    const attr = document.documentElement.getAttribute('data-theme')
    if (attr === 'dark') return true
    if (attr === 'light') return false
    return !!(window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches)
  }

  /** Apply a stored theme value ("dark"|"light"|null) and sync isDark. */
  function applyStored(value: string | null) {
    if (!import.meta.client) return
    if (value === 'dark' || value === 'light') {
      document.documentElement.setAttribute('data-theme', value)
    }
    isDark.value = effectiveDark()
  }

  /** Toggle between dark and light, persisting to localStorage. */
  function toggle() {
    if (!import.meta.client) return
    const next = effectiveDark() ? 'light' : 'dark'
    document.documentElement.setAttribute('data-theme', next)
    try {
      localStorage.setItem('muxi:theme', next)
    } catch {
      // ignore quota / private-browsing errors
    }
    isDark.value = effectiveDark()
  }

  /**
   * label — shows the word for the OPPOSITE mode so clicking toggles to it.
   * dark  → "浅色"  (click to go light)
   * light → "深色"  (click to go dark)
   */
  const label = computed<string>(() => (isDark.value ? '浅色' : '深色'))

  // Capture cleanup fns to run during onUnmounted (registered at setup-call time).
  const cleanupFns: (() => void)[] = []

  onMounted(() => {
    // Sync isDark to whatever the pre-paint script already set.
    isDark.value = effectiveDark()

    // Override with persisted value if present.
    try {
      const stored = localStorage.getItem('muxi:theme')
      if (stored === 'dark' || stored === 'light') {
        applyStored(stored)
      }
    } catch {
      // ignore
    }

    // Cross-tab sync: another tab (or the admin) changed the theme.
    const onStorage = (ev: StorageEvent) => {
      if (ev.key === 'muxi:theme') {
        applyStored(ev.newValue)
      }
    }
    window.addEventListener('storage', onStorage)
    cleanupFns.push(() => window.removeEventListener('storage', onStorage))

    // OS preference change while on "auto" (no explicit localStorage value).
    if (window.matchMedia) {
      const mq = window.matchMedia('(prefers-color-scheme: dark)')
      const onMq = () => {
        if (!document.documentElement.getAttribute('data-theme')) {
          isDark.value = effectiveDark()
        }
      }
      if (mq.addEventListener) {
        mq.addEventListener('change', onMq)
        cleanupFns.push(() => mq.removeEventListener('change', onMq))
      } else {
        // Safari < 14 fallback
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        ;(mq as any).addListener(onMq)
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        cleanupFns.push(() => (mq as any).removeListener(onMq))
      }
    }
  })

  onUnmounted(() => {
    cleanupFns.forEach((fn) => fn())
    cleanupFns.length = 0
  })

  return { isDark: isDark as Readonly<typeof isDark>, toggle, label }
}
