/**
 * useLayoutVariant — home-list layout variant preference.
 *
 * localStorage key : "muxi:variant"
 * values           : "A" (年表, default) | "B" (摘要) | "C" (双栏)
 *
 * SSR-safe: always returns "A" on the server.
 */

export type LayoutVariant = 'A' | 'B' | 'C'

export function useLayoutVariant() {
  const variant = useState<LayoutVariant>('layoutVariant', () => 'A')

  function set(v: LayoutVariant) {
    variant.value = v
    if (import.meta.client) {
      try {
        localStorage.setItem('muxi:variant', v)
      } catch {
        // ignore
      }
    }
  }

  onMounted(() => {
    try {
      const stored = localStorage.getItem('muxi:variant')
      if (stored === 'A' || stored === 'B' || stored === 'C') {
        variant.value = stored
      }
    } catch {
      // ignore
    }
  })

  return { variant, set }
}
