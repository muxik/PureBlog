// Shared helpers for the SEO server routes (sitemap.xml, feed.xml).
// Auto-imported by Nitro into server/routes handlers.
import type { H3Event } from 'h3'
import type { Post, PostListResponse } from '@pureblog/api-types'

/** Escape the five XML-significant characters for safe inclusion in markup. */
export function escapeXml(s: string): string {
  return s
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&apos;')
}

/** Join an origin and a path into one absolute URL, collapsing the seam slash. */
export function absoluteUrl(siteUrl: string, path: string): string {
  return `${siteUrl.replace(/\/$/, '')}/${path.replace(/^\//, '')}`
}

/**
 * Fetch all published posts from the backend during SSR. Uses the server-side
 * API base (direct to the backend container) and a generous page size so the
 * sitemap/feed cover the whole archive in one request.
 */
export async function fetchPublishedPosts(event: H3Event): Promise<Post[]> {
  const cfg = useRuntimeConfig(event)
  const res = await $fetch<PostListResponse>(`${cfg.apiBaseServer}/posts`, {
    query: { pageSize: 1000 },
  })
  return res.items ?? []
}
