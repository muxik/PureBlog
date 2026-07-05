/// <reference path="./solarlunar.d.ts" />

import solarLunar from 'solarlunar'
import type { SolarLunarResult } from 'solarlunar'

export type LunarDate = { year: number; month: number; day: number; isLeap: boolean }

const GAN = ['甲', '乙', '丙', '丁', '戊', '己', '庚', '辛', '壬', '癸']
const ZHI = ['子', '丑', '寅', '卯', '辰', '巳', '午', '未', '申', '酉', '戌', '亥']

function isValidDate(date: Date): boolean {
  return !Number.isNaN(date.getTime())
}

function ganzhi(year: number): string {
  return GAN[(year - 4) % 10] + ZHI[(year - 4) % 12]
}

export function toLunar(date: Date): LunarDate | null {
  if (!isValidDate(date)) return null

  const lunar = lunarFor(date)
  if (!lunar) return null

  return {
    year: lunar.lYear,
    month: lunar.lMonth,
    day: lunar.lDay,
    isLeap: lunar.isLeap,
  }
}

function parse(str: string): Date | null {
  const m = String(str).match(/(\d{4})\D+(\d{1,2})\D+(\d{1,2})/)
  if (!m) return null
  const date = new Date(Number(m[1]), Number(m[2]) - 1, Number(m[3]))
  return isValidDate(date) ? date : null
}

function asDate(d: Date | string): Date | null {
  if (d instanceof Date) return isValidDate(d) ? d : null
  return parse(d)
}

function lunarFor(d: Date | string): SolarLunarResult | null {
  const dt = asDate(d)
  if (!dt) return null

  const result = solarLunar.solar2lunar(dt.getFullYear(), dt.getMonth() + 1, dt.getDate())
  return result === -1 ? null : result
}

/** Short form, e.g. 正月廿六 (no year). Accepts a Date or "YYYY · MM · DD". */
export function short(d: Date | string): string | null {
  const lunar = lunarFor(d)
  if (!lunar) return null
  return `${lunar.monthCn}${lunar.dayCn}`
}

/** Full form, e.g. 丙午年 正月廿六. Accepts a Date or "YYYY · MM · DD". */
export function full(d: Date | string): string | null {
  const lunar = lunarFor(d)
  if (!lunar) return null
  return `${ganzhi(lunar.lYear)}年 ${lunar.monthCn}${lunar.dayCn}`
}
