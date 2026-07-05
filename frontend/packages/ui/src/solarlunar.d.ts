declare module 'solarlunar' {
  export interface SolarLunarResult {
    lYear: number
    lMonth: number
    lDay: number
    monthCn: string
    dayCn: string
    isLeap: boolean
  }

  const solarLunar: {
    solar2lunar(year: number, month: number, day: number): SolarLunarResult | -1
  }

  export default solarLunar
}
