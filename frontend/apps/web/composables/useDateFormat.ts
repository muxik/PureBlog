import { computed, onMounted, onUnmounted, toValue, watch } from 'vue'
import type { MaybeRefOrGetter } from 'vue'
import { useSiteSettings } from './useSiteSettings'

export type DateFormat = 'numeric' | 'lunar'

function toDateFormat(value: unknown): DateFormat | undefined {
  return value === 'numeric' || value === 'lunar' ? value : undefined
}

function readStoredFormat(): DateFormat | undefined {
  try {
    return toDateFormat(localStorage.getItem('muxi:dateFormat'))
  } catch {
    return undefined
  }
}

function writeStoredFormat(value: DateFormat): void {
  try {
    localStorage.setItem('muxi:dateFormat', value)
  } catch {
    // ignore
  }
}

/**
 * Reactive date-format preference.
 *
 * When a backend setting is provided, it is authoritative and is mirrored into
 * localStorage for old design-system listeners. localStorage is only used as a
 * fallback when no backend value is available yet.
 */
export function useDateFormat(fallback?: MaybeRefOrGetter<DateFormat | string | undefined>) {
  const backendFormat = computed(() => toDateFormat(toValue(fallback)))
  const format = useState<DateFormat>('dateFormat', () => backendFormat.value ?? 'numeric')

  const cleanupFns: (() => void)[] = []

  watch(
    backendFormat,
    (value) => {
      if (!value) return
      format.value = value
      if (import.meta.client) writeStoredFormat(value)
    },
    { immediate: true },
  )

  onMounted(() => {
    if (backendFormat.value) {
      writeStoredFormat(backendFormat.value)
    } else {
      const stored = readStoredFormat()
      if (stored) format.value = stored
    }

    // Cross-tab sync (admin panel or another tab changed the format).
    const onStorage = (ev: StorageEvent) => {
      if (ev.key === 'muxi:dateFormat') {
        const v = toDateFormat(ev.newValue)
        if (v) format.value = v
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

export function useSiteDateFormat() {
  const settings = useSiteSettings()
  return useDateFormat(computed(() => settings.value.defaultDateFormat))
}
