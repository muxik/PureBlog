// Returns the API base URL: the server-only base during SSR (so the Nuxt server
// talks to the backend container directly), the public base in the browser.
export function useApiBase(): string {
  const cfg = useRuntimeConfig()
  return import.meta.server ? cfg.apiBaseServer : cfg.public.apiBase
}
