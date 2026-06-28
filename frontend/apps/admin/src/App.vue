<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from './stores/auth'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

function logout() {
  auth.logout()
  router.push('/login')
}
</script>

<template>
  <div class="admin">
    <nav v-if="auth.isAuthed && route.path !== '/login'" class="nav">
      <div class="nav-links">
        <RouterLink to="/manage">文章</RouterLink>
        <RouterLink to="/categories">分类</RouterLink>
        <RouterLink to="/tags">标签</RouterLink>
        <RouterLink to="/comments">评论</RouterLink>
        <RouterLink to="/settings">设置</RouterLink>
      </div>
      <button class="ghost logout" @click="logout">退出</button>
    </nav>
    <RouterView />
  </div>
</template>

<style>
.admin {
  max-width: 48rem;
  margin: 0 auto;
  padding: 2rem 1.25rem;
}
.admin input,
.admin button {
  font: inherit;
}
.admin input {
  display: block;
  width: 100%;
  padding: 0.5rem 0.7rem;
  margin: 0.4rem 0;
  border: 1px solid var(--border, #e7e2d6);
  border-radius: 4px;
  background: var(--paper, #faf8f2);
  color: inherit;
}
.admin button {
  padding: 0.45rem 1rem;
  margin-right: 0.5rem;
  border: 1px solid var(--accent, #235a73);
  border-radius: 4px;
  background: var(--accent, #235a73);
  color: #fff;
  cursor: pointer;
}
.admin button.ghost {
  background: transparent;
  color: var(--accent, #235a73);
}
.err {
  color: #b23b3b;
}
.admin .nav {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  margin-bottom: 1.75rem;
  padding-bottom: 0.75rem;
  border-bottom: 1px solid var(--border, #e7e2d6);
}
.admin .nav-links {
  display: flex;
  gap: 1.25rem;
}
.admin .nav-links a {
  color: var(--ink-2, #6b655c);
  text-decoration: none;
}
.admin .nav-links a:hover {
  color: var(--accent, #235a73);
}
.admin .nav-links a.router-link-active {
  color: var(--accent, #235a73);
}
.admin .nav .logout {
  margin-right: 0;
  padding: 0.25rem 0.7rem;
}
</style>
