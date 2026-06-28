// PureBlog v3 — public site (SSR). https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-01-01',
  devtools: { enabled: false },
  css: [
    // 1. Design-system tokens (colors, typography, spacing, fonts + CDN @imports).
    '@pureblog/ui/tokens.css',
    // 2. Application component styles (verbatim from muxi-blog-design/styles/app.css).
    '~/assets/css/app.css',
  ],
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
      // Pre-paint theme script: apply the saved theme before first paint to
      // avoid a light→dark flash (FOUC). Verbatim from muxi-blog-design/index.html.
      script: [
        {
          innerHTML: `try{var t=localStorage.getItem("muxi:theme");if(t)document.documentElement.setAttribute("data-theme",t);}catch(e){}`,
        },
      ],
    },
  },
})
