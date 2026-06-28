/**
 * useSiteSettings — fetch and cache public site settings from GET /settings.
 *
 * Fetched once per SSR request (deduped by key 'site-settings') and
 * rehydrated on the client from the Nuxt payload — no extra round-trip.
 *
 * Returns a Ref<SiteSettings> (never null; falls back to empty object
 * so consumers can safely read .siteName ?? 'PureBlog').
 */

import type { SiteSettings } from '@pureblog/api-types'

export function useSiteSettings(): Ref<SiteSettings> {
  const apiBase = useApiBase()
  const { data } = useAsyncData<SiteSettings>(
    'site-settings',
    () => $fetch<SiteSettings>(`${apiBase}/settings`),
    { default: () => ({} as SiteSettings) },
  )
  // data is Ref<SiteSettings> because of the default factory above.
  return data as Ref<SiteSettings>
}
