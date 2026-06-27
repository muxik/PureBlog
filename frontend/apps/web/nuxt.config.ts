// PureBlog v3 — public site (SSR). https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-01-01',
  devtools: { enabled: false },
  css: ['@pureblog/ui/tokens.css'],
  runtimeConfig: {
    // server-side base (SSR → backend directly); override: NUXT_API_BASE_SERVER
    apiBaseServer: 'http://localhost:8080/api/v1',
    public: {
      // browser base; override: NUXT_PUBLIC_API_BASE
      apiBase: 'http://localhost:8080/api/v1',
    },
  },
  app: {
    head: {
      htmlAttrs: { lang: 'zh' },
      title: 'PureBlog',
      meta: [{ name: 'viewport', content: 'width=device-width, initial-scale=1' }],
    },
  },
})
