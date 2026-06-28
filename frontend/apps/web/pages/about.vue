<script setup lang="ts">
const settings = useSiteSettings()

// ── Tiny inline Markdown → HTML converter ──────────────────────────────────
// Ported from the muxi-blog-design spirit: headings, blockquote, fenced code,
// unordered lists, bold, italic, inline code, links, paragraphs.
// aboutMd is admin-controlled content, so we escape user text before inserting
// our own HTML, which keeps XSS risk to the same level as the admin panel itself.

function esc(s: string): string {
  return s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')
}

function renderInline(s: string): string {
  return (
    esc(s)
      // links [text](url) — URL may contain & which was already escaped, href is fine
      .replace(/\[([^\]]+)\]\(([^)]+)\)/g, '<a href="$2">$1</a>')
      // bold **text**
      .replace(/\*\*(.+?)\*\*/g, '<strong>$1</strong>')
      // italic *text*
      .replace(/\*(.+?)\*/g, '<em>$1</em>')
      // inline code `code`
      .replace(/`([^`]+)`/g, '<code>$1</code>')
  )
}

function md(src: string): string {
  if (!src) return ''
  const lines = src.split('\n')
  const out: string[] = []
  let i = 0
  while (i < lines.length) {
    const line = lines[i]
    // Heading  # … ######
    const hm = line.match(/^(#{1,6}) (.+)/)
    if (hm) {
      const lvl = hm[1].length
      out.push(`<h${lvl}>${renderInline(hm[2])}</h${lvl}>`)
      i++
      continue
    }
    // Blockquote  > …
    if (line.startsWith('> ')) {
      out.push(`<blockquote><p>${renderInline(line.slice(2))}</p></blockquote>`)
      i++
      continue
    }
    // Fenced code block  ```
    if (line.startsWith('```')) {
      const acc: string[] = []
      i++
      while (i < lines.length && !lines[i].startsWith('```')) {
        acc.push(esc(lines[i]))
        i++
      }
      i++ // skip closing ```
      out.push(`<pre><code>${acc.join('\n')}</code></pre>`)
      continue
    }
    // Unordered list  - …  or  * …
    if (/^[-*] /.test(line)) {
      const items: string[] = []
      while (i < lines.length && /^[-*] /.test(lines[i])) {
        items.push(`<li>${renderInline(lines[i].slice(2))}</li>`)
        i++
      }
      out.push(`<ul class="muxi-list">${items.join('')}</ul>`)
      continue
    }
    // Blank line — skip
    if (!line.trim()) {
      i++
      continue
    }
    // Paragraph — collect consecutive non-special lines
    const p: string[] = []
    while (
      i < lines.length &&
      lines[i].trim() &&
      !lines[i].match(/^#{1,6} /) &&
      !lines[i].startsWith('> ') &&
      !lines[i].startsWith('```') &&
      !/^[-*] /.test(lines[i])
    ) {
      p.push(renderInline(lines[i]))
      i++
    }
    if (p.length) out.push(`<p>${p.join('\n')}</p>`)
  }
  return out.join('\n')
}

// ── Reactive data ───────────────────────────────────────────────────────────

const aboutHtml = computed(() => md(settings.value.aboutMd ?? ''))

const gridRows = computed(() => {
  const s = settings.value
  const rows: { k: string; v: string }[] = []
  if (s.author) rows.push({ k: '作者', v: s.author })
  if (s.siteName) rows.push({ k: '站点', v: s.siteName })
  if (s.description) rows.push({ k: '简介', v: s.description })
  if (s.social) {
    for (const [key, val] of Object.entries(s.social)) {
      if (val) rows.push({ k: key, v: String(val) })
    }
  }
  return rows
})

useHead(() => ({
  title: '关于 · ' + (settings.value.siteName ?? 'PureBlog'),
}))
</script>

<template>
  <section class="section--pad wrap">
    <p class="about-eyebrow">ABOUT</p>
    <h1 class="about-title">{{ settings.siteName || '关于' }}</h1>
    <!-- eslint-disable-next-line vue/no-v-html -->
    <div class="prose" v-html="aboutHtml" />
    <div v-if="gridRows.length" class="about-grid">
      <template v-for="row in gridRows" :key="row.k">
        <span class="about-grid__k">{{ row.k }}</span>
        <span class="about-grid__v">{{ row.v }}</span>
      </template>
    </div>
  </section>
</template>
