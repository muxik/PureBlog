// GET /sitemap.xml — a search-engine sitemap covering the static pages and
// every published post. Rendered on the fly from the backend post list.
export default defineEventHandler(async (event) => {
  const cfg = useRuntimeConfig(event)
  const siteUrl = cfg.public.siteUrl as string
  const posts = await fetchPublishedPosts(event)

  // Static, crawlable pages (the search page is intentionally excluded).
  const staticPaths = ['/', '/archive', '/tags', '/about']

  const urls: string[] = []
  for (const path of staticPaths) {
    urls.push(`  <url><loc>${escapeXml(absoluteUrl(siteUrl, path))}</loc></url>`)
  }
  for (const p of posts) {
    const lastmod = p.updatedAt || p.publishedAt || p.createdAt
    const loc = escapeXml(absoluteUrl(siteUrl, `/post/${p.slug}`))
    urls.push(
      lastmod
        ? `  <url><loc>${loc}</loc><lastmod>${new Date(lastmod).toISOString()}</lastmod></url>`
        : `  <url><loc>${loc}</loc></url>`,
    )
  }

  const xml =
    `<?xml version="1.0" encoding="UTF-8"?>\n` +
    `<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">\n` +
    urls.join('\n') +
    `\n</urlset>\n`

  setResponseHeader(event, 'Content-Type', 'application/xml; charset=utf-8')
  return xml
})
