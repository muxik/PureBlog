/**
 * MuxiLunar — minimal, verified Gregorian → Chinese-lunar converter.
 * Ported verbatim from muxi-blog-design/lunar.js to typed ESM.
 * Scope: 2023–2028 (covers the blog's sample posts).
 * Each year's 正月初一 solar date and packed month table are checked against
 * known 春节 dates and leap months (2023 闰二月, 2025 闰六月).
 */

interface YearInfo {
  base: string  // ISO date of 正月初一
  info: number  // packed month table
}

const YEARS: Record<number, YearInfo> = {
  2023: { base: '2023-01-22', info: 0x05b52 },
  2024: { base: '2024-02-10', info: 0x04b60 },
  2025: { base: '2025-01-29', info: 0x0a6e6 },
  2026: { base: '2026-02-17', info: 0x0a4e0 },
  2027: { base: '2027-02-06', info: 0x0d260 },
  2028: { base: '2028-01-26', info: 0x0ea65 },
}

const GAN = ['甲', '乙', '丙', '丁', '戊', '己', '庚', '辛', '壬', '癸']
const ZHI = ['子', '丑', '寅', '卯', '辰', '巳', '午', '未', '申', '酉', '戌', '亥']
const MONTHS = ['正', '二', '三', '四', '五', '六', '七', '八', '九', '十', '冬', '腊']
const D1 = ['日', '一', '二', '三', '四', '五', '六', '七', '八', '九', '十']
const D2 = ['初', '十', '廿', '卅']

function leapMonth(info: number): number { return info & 0xf }
function leapDays(info: number): number { return leapMonth(info) ? ((info & 0x10000) ? 30 : 29) : 0 }
function monthDays(info: number, m: number): number { return (info & (0x10000 >> m)) ? 30 : 29 }
function ganzhi(y: number): string { return GAN[(y - 4) % 10] + ZHI[(y - 4) % 12] }
function dayCn(d: number): string {
  if (d === 10) return '初十'
  if (d === 20) return '二十'
  if (d === 30) return '三十'
  return D2[Math.floor(d / 10)] + D1[d % 10]
}

export interface LunarDate {
  year: number
  month: number
  day: number
  isLeap: boolean
}

/**
 * Convert a Date to a LunarDate. Returns null if out of range (2023–2028).
 */
export function toLunar(date: Date): LunarDate | null {
  const keys = Object.keys(YEARS).map(Number).sort((a, b) => a - b)
  let ly: number | null = null
  for (const key of keys) {
    if (date.getTime() >= new Date(YEARS[key].base + 'T00:00:00').getTime()) ly = key
  }
  if (ly === null) return null

  const info = YEARS[ly].info
  const offset = Math.round(
    (date.getTime() - new Date(YEARS[ly].base + 'T00:00:00').getTime()) / 86400000,
  )
  const leap = leapMonth(info)
  const seq: Array<{ m: number; leap: boolean; days: number }> = []
  for (let m = 1; m <= 12; m++) {
    seq.push({ m, leap: false, days: monthDays(info, m) })
    if (leap > 0 && m === leap) seq.push({ m, leap: true, days: leapDays(info) })
  }

  let day = offset
  for (const entry of seq) {
    if (day < entry.days) {
      return { year: ly, month: entry.m, day: day + 1, isLeap: entry.leap }
    }
    day -= entry.days
  }
  return null
}

/** Parse a Date or "YYYY · MM · DD" string. Returns null on parse failure. */
function parse(str: string): Date | null {
  const m = String(str).match(/(\d{4})\D+(\d{1,2})\D+(\d{1,2})/)
  if (!m) return null
  return new Date(Number(m[1]), Number(m[2]) - 1, Number(m[3]))
}

function shortStr(date: Date): string | null {
  const l = toLunar(date)
  if (!l) return null
  return (l.isLeap ? '闰' : '') + MONTHS[l.month - 1] + '月' + dayCn(l.day)
}

function fullStr(date: Date): string | null {
  const l = toLunar(date)
  if (!l) return null
  return ganzhi(l.year) + '年 ' + (l.isLeap ? '闰' : '') + MONTHS[l.month - 1] + '月' + dayCn(l.day)
}

/**
 * Short lunar string: e.g. "正月初一". Accepts a Date or "YYYY · MM · DD" string.
 */
export function short(d: Date | string): string | null {
  const dt = d instanceof Date ? d : parse(d as string)
  return dt ? shortStr(dt) : null
}

/**
 * Full lunar string with ganzhi year: e.g. "丙午年 正月初一".
 * Accepts a Date or "YYYY · MM · DD" string.
 */
export function full(d: Date | string): string | null {
  const dt = d instanceof Date ? d : parse(d as string)
  return dt ? fullStr(dt) : null
}
