/**
 * useDateFormat — reactive date-format preference composable.
 *
 * localStorage key : "muxi:dateFormat"
 * values           : "numeric" (default) | "lunar"
 *
 * SSR-safe: always returns "numeric" on the server.
 * Cross-tab sync: listens to the `storage` event for "muxi:dateFormat".
 *
 * @param fallback  Optional default when no localStorage value exists.
 *                  Typically from SiteSettings.defaultDateFormat.
 *                  Falls back to "numeric" if not provided.
 */

export type DateFormat = 'numeric' | 'lunar'

export function useDateFormat(fallback?: DateFormat) {
  const format = useState<DateFormat>('dateFormat', () => fallback ?? 'numeric')

  const cleanupFns: (() => void)[] = []

  onMounted(() => {
    // Read persisted value on client.
    try {
      const stored = localStorage.getItem('muxi:dateFormat')
      if (stored === 'numeric' || stored === 'lunar') {
        format.value = stored
      } else if (fallback) {
        format.value = fallback
      }
    } catch {
      // ignore
    }

    // Cross-tab sync (admin panel or another tab changed the format).
    const onStorage = (ev: StorageEvent) => {
      if (ev.key === 'muxi:dateFormat') {
        const v = ev.newValue
        if (v === 'numeric' || v === 'lunar') {
          format.value = v
        }
      }
    }
    window.addEventListener('storage', onStorage)
    cleanupFns.push(() => window.removeEventListener('storage', onStorage))
  })

  onUnmounted(() => {
    cleanupFns.forEach((fn) => fn())
    cleanupFns.length = 0
  })

  return { format }
}
