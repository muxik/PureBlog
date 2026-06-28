// GET /feed.xml — an RSS 2.0 feed of the latest published posts, rendered from
// the backend post list and site settings.
import type { SiteSettings } from '@pureblog/api-types'

export default defineEventHandler(async (event) => {
  const cfg = useRuntimeConfig(event)
  const siteUrl = cfg.public.siteUrl as string

  const [posts, settings] = await Promise.all([
    fetchPublishedPosts(event),
    $fetch<SiteSettings>(`${cfg.apiBaseServer}/settings`).catch(() => ({}) as SiteSettings),
  ])

  const siteName = settings.siteName || 'PureBlog'
  const description = settings.description || ''
  const self = escapeXml(absoluteUrl(siteUrl, '/feed.xml'))
  const home = escapeXml(siteUrl.replace(/\/$/, ''))

  // Newest 30 posts, sorted by publish/create time descending.
  const items = posts
    .slice()
    .sort((a, b) => {
      const ta = new Date(a.publishedAt || a.createdAt || 0).getTime()
      const tb = new Date(b.publishedAt || b.createdAt || 0).getTime()
      return tb - ta
    })
    .slice(0, 30)
    .map((p) => {
      const link = escapeXml(absoluteUrl(siteUrl, `/post/${p.slug}`))
      const pub = p.publishedAt || p.createdAt
      const date = pub ? new Date(pub).toUTCString() : ''
      return (
        `  <item>\n` +
        `    <title>${escapeXml(p.title || '')}</title>\n` +
        `    <link>${link}</link>\n` +
        `    <guid isPermaLink="true">${link}</guid>\n` +
        (date ? `    <pubDate>${date}</pubDate>\n` : '') +
        `    <description>${escapeXml(p.summary || '')}</description>\n` +
        `  </item>`
      )
    })
    .join('\n')

  const xml =
    `<?xml version="1.0" encoding="UTF-8"?>\n` +
    `<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">\n` +
    `  <channel>\n` +
    `    <title>${escapeXml(siteName)}</title>\n` +
    `    <link>${home}</link>\n` +
    `    <description>${escapeXml(description)}</description>\n` +
    `    <language>zh-CN</language>\n` +
    `    <atom:link href="${self}" rel="self" type="application/rss+xml"/>\n` +
    items +
    `\n  </channel>\n</rss>\n`

  setResponseHeader(event, 'Content-Type', 'application/rss+xml; charset=utf-8')
  return xml
})
