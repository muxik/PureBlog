// MuxiLunar — minimal, verified Gregorian→Chinese-lunar converter.
// Scope: 2023–2028. Ported verbatim from the design prototype's lunar.js.

export type LunarDate = { year: number; month: number; day: number; isLeap: boolean }

const YEARS: Record<number, { base: string; info: number }> = {
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

const leapMonth = (info: number) => info & 0xf
const leapDays = (info: number) => (leapMonth(info) ? (info & 0x10000 ? 30 : 29) : 0)
const monthDays = (info: number, m: number) => (info & (0x10000 >> m) ? 30 : 29)
const ganzhi = (y: number) => GAN[(y - 4) % 10] + ZHI[(y - 4) % 12]

function dayCn(d: number): string {
  if (d === 10) return '初十'
  if (d === 20) return '二十'
  if (d === 30) return '三十'
  return D2[Math.floor(d / 10)] + D1[d % 10]
}

export function toLunar(date: Date): LunarDate | null {
  const keys = Object.keys(YEARS).map(Number).sort((a, b) => a - b)
  let ly: number | null = null
  for (const k of keys) {
    if (date.getTime() >= new Date(YEARS[k].base + 'T00:00:00').getTime()) ly = k
  }
  if (ly == null) return null

  const info = YEARS[ly].info
  const offset = Math.round(
    (date.getTime() - new Date(YEARS[ly].base + 'T00:00:00').getTime()) / 86400000,
  )
  const leap = leapMonth(info)
  const seq: { m: number; leap: boolean; days: number }[] = []
  for (let m = 1; m <= 12; m++) {
    seq.push({ m, leap: false, days: monthDays(info, m) })
    if (leap > 0 && m === leap) seq.push({ m, leap: true, days: leapDays(info) })
  }

  let day = offset
  for (const s of seq) {
    if (day < s.days) return { year: ly, month: s.m, day: day + 1, isLeap: s.leap }
    day -= s.days
  }
  return null
}

function parse(str: string): Date | null {
  const m = String(str).match(/(\d{4})\D+(\d{1,2})\D+(\d{1,2})/)
  if (!m) return null
  return new Date(Number(m[1]), Number(m[2]) - 1, Number(m[3]))
}

function asDate(d: Date | string): Date | null {
  return d instanceof Date ? d : parse(d)
}

/** Short form, e.g. 正月廿六 (no year). Accepts a Date or "YYYY · MM · DD". */
export function short(d: Date | string): string | null {
  const dt = asDate(d)
  if (!dt) return null
  const l = toLunar(dt)
  if (!l) return null
  return (l.isLeap ? '闰' : '') + MONTHS[l.month - 1] + '月' + dayCn(l.day)
}

/** Full form, e.g. 丙午年 正月廿六. Accepts a Date or "YYYY · MM · DD". */
export function full(d: Date | string): string | null {
  const dt = asDate(d)
  if (!dt) return null
  const l = toLunar(dt)
  if (!l) return null
  return ganzhi(l.year) + '年 ' + (l.isLeap ? '闰' : '') + MONTHS[l.month - 1] + '月' + dayCn(l.day)
}
