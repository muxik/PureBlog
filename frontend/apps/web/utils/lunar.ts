// Thin re-export so page agents can import from '~/utils/lunar'
// The implementation lives in @pureblog/ui/src/lunar.ts (verbatim port of
// muxi-blog-design/lunar.js, verified for 2023–2028).
export { toLunar, short, full } from '@pureblog/ui/lunar'
export type { LunarDate } from '@pureblog/ui/lunar'
