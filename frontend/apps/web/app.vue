<script setup lang="ts">
// Site chrome: header + main + footer.
// Page agents own pages/*.vue; this file owns only the shared shell.

const route = useRoute()
const { toggle, label: themeLabel } = useTheme()
const settings = useSiteSettings()

const siteName = computed(() => settings.value?.siteName ?? 'PureBlog')
const currentYear = computed(() => new Date().getFullYear())

/**
 * Whether a nav link is "active" for the current route.
 * Home ("/") uses exact match; all other paths use startsWith.
 */
function isActive(path: string): boolean {
  if (path === '/') return route.path === '/'
  return route.path.startsWith(path)
}
</script>

<template>
  <div class="app">
    <header class="site-header wrap">
      <NuxtLink class="masthead" to="/" aria-label="回到首页">
        <span class="masthead__name">{{ siteName }}</span>
      </NuxtLink>

      <nav class="nav">
        <NuxtLink
          class="nav__link"
          :class="{ 'is-active': isActive('/') }"
          to="/"
        >文章</NuxtLink>

        <NuxtLink
          class="nav__link"
          :class="{ 'is-active': isActive('/archive') }"
          to="/archive"
        >归档</NuxtLink>

        <NuxtLink
          class="nav__link"
          :class="{ 'is-active': isActive('/tags') }"
          to="/tags"
        >标签</NuxtLink>

        <NuxtLink
          class="nav__link"
          :class="{ 'is-active': isActive('/about') }"
          to="/about"
        >关于</NuxtLink>

        <NuxtLink
          class="nav__link"
          :class="{ 'is-active': isActive('/search') }"
          to="/search"
        >搜索</NuxtLink>

        <span class="nav__sep" aria-hidden="true"></span>

        <button
          class="nav__toggle"
          aria-label="切换深浅色"
          @click="toggle"
        >{{ themeLabel }}</button>
      </nav>
    </header>

    <main class="site-main">
      <NuxtPage />
    </main>

    <footer class="site-footer wrap">
      <span class="site-footer__dot" aria-hidden="true"></span>
      <span>© {{ currentYear }} {{ siteName }} · 保留所有权利</span>
    </footer>
  </div>
</template>
