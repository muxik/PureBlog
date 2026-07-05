<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from './stores/auth'
import { useTheme, applyStoredTheme } from './utils/theme'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const { label: themeLabel, toggle: toggleTheme } = useTheme()
const mobileNavOpen = ref(false)

// Sync stored theme on mount (the head script already ran before first paint;
// this ensures the reactive composable state is consistent after hydration).
applyStoredTheme()

function logout() {
  auth.logout()
  router.push('/login')
}

// Nav entries: existing routes + the design's three core areas
const navItems = [
  { to: '/write',      label: '写作' },
  { to: '/manage',     label: '文章' },
  { to: '/tags',       label: '标签' },
  { to: '/comments',   label: '评论' },
  { to: '/settings',   label: '设置' },
] as const

watch(() => route.fullPath, () => {
  mobileNavOpen.value = false
})
</script>

<template>
  <!-- Authenticated shell -->
  <div v-if="auth.isAuthed && route.path !== '/login'" class="admin">
    <header class="admin-header">
      <div class="admin-header__bar">
        <RouterLink class="admin-brand" to="/write">PureBlog</RouterLink>
        <button
          class="menu-button"
          type="button"
          :class="{ 'is-open': mobileNavOpen }"
          aria-label="打开导航菜单"
          :aria-expanded="mobileNavOpen"
          aria-controls="admin-nav-panel"
          @click="mobileNavOpen = !mobileNavOpen"
        >
          <span aria-hidden="true"></span>
          <span aria-hidden="true"></span>
          <span aria-hidden="true"></span>
        </button>
      </div>

      <div
        id="admin-nav-panel"
        class="admin-header__menu"
        :class="{ 'is-open': mobileNavOpen }"
      >
        <div class="admin-header__left">
          <nav class="admin-nav">
            <RouterLink
              v-for="item in navItems"
              :key="item.to"
              :to="item.to"
              class="admin-nav__btn"
              :class="{ 'is-active': route.path === item.to || route.path.startsWith(item.to + '/') }"
            >{{ item.label }}</RouterLink>
          </nav>
        </div>
        <div class="admin-header__right">
          <button class="chrome-btn" aria-label="切换深浅色" @click="toggleTheme">
            {{ themeLabel }}
          </button>
          <span class="vsep" />
          <button class="chrome-btn" @click="logout">退出</button>
        </div>
      </div>
    </header>
    <main class="admin-main">
      <RouterView />
    </main>
  </div>

  <!-- Unauthenticated: just render the login route -->
  <RouterView v-else />
</template>
