// Thin typed fetch wrapper. Reads the access token straight from localStorage to
// avoid a circular import with the auth store.
const BASE = import.meta.env.VITE_API_BASE ?? 'http://localhost:8080/api/v1'

export const ACCESS_KEY = 'pb:access'
export const REFRESH_KEY = 'pb:refresh'

type Options = { method?: string; body?: unknown; auth?: boolean }

export async function api<T>(path: string, opts: Options = {}): Promise<T> {
  const headers: Record<string, string> = { 'Content-Type': 'application/json' }
  if (opts.auth !== false) {
    const token = localStorage.getItem(ACCESS_KEY)
    if (token) headers.Authorization = `Bearer ${token}`
  }
  const res = await fetch(BASE + path, {
    method: opts.method ?? 'GET',
    headers,
    body: opts.body !== undefined ? JSON.stringify(opts.body) : undefined,
  })
  if (!res.ok) throw new Error(`API ${res.status}`)
  if (res.status === 204) return undefined as T
  return (await res.json()) as T
}
