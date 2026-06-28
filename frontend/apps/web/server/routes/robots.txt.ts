// GET /robots.txt — allow all crawlers and advertise the sitemap.
export default defineEventHandler((event) => {
  const cfg = useRuntimeConfig(event)
  const siteUrl = (cfg.public.siteUrl as string).replace(/\/$/, '')
  setResponseHeader(event, 'Content-Type', 'text/plain; charset=utf-8')
  return `User-agent: *\nAllow: /\n\nSitemap: ${siteUrl}/sitemap.xml\n`
})
