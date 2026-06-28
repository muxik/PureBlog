/**
 * usePostView — data adapter: converts raw API Post → display-ready PostView.
 *
 * Two pure functions (no composable boilerplate needed — import them directly):
 *   decoratePost(p, fmt)   → PostView
 *   groupByYear(views)     → { year, posts }[]
 */

import type { Post, Tag } from '@pureblog/api-types'
import { short as lunarShort, full as lunarFull } from '~/utils/lunar'
import type { DateFormat } from '~/composables/useDateFormat'

// ── Types ────────────────────────────────────────────────────────────────────

export interface PostView {
  /** URL slug, e.g. "hello-world" */
  slug: string
  /** Post title */
  title: string
  /** 4-digit year string, e.g. "2026" */
  year: string
  /**
   * Full date string.
   * numeric → "2026 · 03 · 14"
   * lunar   → "丙午年 正月廿六" (falls back to numeric if out-of-range)
   */
  date: string
  /**
   * Short month·day string for timeline variant A.
   * numeric → "03 · 14"
   * lunar   → "正月廿六" (falls back to numeric if out-of-range)
   */
  md: string
  /** Raw tag objects from the API */
  tags: Tag[]
  /** "排版 · 随笔" */
  tagStr: string
  /** "2026 · 03 · 14　·　排版 · 随笔" (full-width spaces around middle dot) */
  metaLine: string
  /** Post summary / excerpt */
  excerpt: string
  /** Whether the post is pinned to the top */
  pinned: boolean
  /**
   * Rough read-time estimate in minutes (0 = unknown).
   * Derived from summary length as a proxy (real content not always available).
   */
  readMin: number
}

// ── Helpers ──────────────────────────────────────────────────────────────────

function pad2(n: number): string {
  return String(n).padStart(2, '0')
}

// ── decoratePost ─────────────────────────────────────────────────────────────

/**
 * Convert a raw API Post into a display-ready PostView.
 *
 * @param p    Post from GET /posts or GET /posts/:slug
 * @param fmt  Current date format preference ("numeric" | "lunar")
 */
export function decoratePost(p: Post, fmt: DateFormat): PostView {
  // Base date: prefer publishedAt, fall back to createdAt, fall back to now.
  const baseStr = p.publishedAt ?? p.createdAt ?? new Date().toISOString()
  const base = new Date(baseStr)

  const y = base.getFullYear()
  const mo = base.getMonth() + 1
  const d = base.getDate()

  const numericDate = `${y} · ${pad2(mo)} · ${pad2(d)}`
  const numericMd = `${pad2(mo)} · ${pad2(d)}`

  let date: string
  let md: string

  if (fmt === 'lunar') {
    date = lunarFull(base) ?? numericDate
    md = lunarShort(base) ?? numericMd
  } else {
    date = numericDate
    md = numericMd
  }

  const tags: Tag[] = p.tags ?? []
  const tagStr = tags
    .map((t) => t.name ?? '')
    .filter(Boolean)
    .join(' · ')

  // Full-width spaces (U+3000) around the middle dot separator, matching the design.
  const metaLine = `${date}　·　${tagStr}`

  // Rough read estimate: ~400 Chinese characters per minute from the summary.
  const readMin = p.summary ? Math.max(1, Math.ceil(p.summary.length / 400)) : 0

  return {
    slug: p.slug ?? '',
    title: p.title ?? '',
    year: String(y),
    date,
    md,
    tags,
    tagStr,
    metaLine,
    excerpt: p.summary ?? '',
    pinned: p.pinned ?? false,
    readMin,
  }
}

// ── groupByYear ───────────────────────────────────────────────────────────────

/**
 * Group a flat PostView array by year, preserving input order.
 * Each group: { year: "2026", posts: PostView[] }
 */
export function groupByYear(views: PostView[]): { year: string; posts: PostView[] }[] {
  const seen: string[] = []
  for (const v of views) {
    if (!seen.includes(v.year)) seen.push(v.year)
  }
  return seen.map((year) => ({
    year,
    posts: views.filter((v) => v.year === year),
  }))
}
