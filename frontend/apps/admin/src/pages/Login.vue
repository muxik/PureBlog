<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const username = ref('')
const password = ref('')
const error = ref('')
const auth = useAuthStore()
const router = useRouter()

async function submit() {
  error.value = ''
  try {
    await auth.login({ username: username.value, password: password.value })
    router.push('/manage')
  } catch {
    error.value = '用户名或密码不对，再试试。'
  }
}
</script>

<template>
  <div class="login-wrap">
    <div class="login-card">
      <div class="login-head">
        <span class="seal" style="width:40px;height:40px;font-size:24px;">博</span>
        <h1 class="login-title">PureBlog 后台</h1>
        <p class="login-sub">写作后台 · 请登录</p>
      </div>
      <form class="login-form" @submit.prevent="submit">
        <input
          v-model="username"
          class="admin-input"
          type="text"
          placeholder="用户名"
          autocomplete="username"
        />
        <input
          v-model="password"
          class="admin-input"
          type="password"
          placeholder="密码"
          autocomplete="current-password"
        />
        <p v-if="error" class="login-error">{{ error }}</p>
        <button class="btn-solid" type="submit" style="margin-top:4px;">登录</button>
      </form>
      <p class="login-hint">PureBlog · 写作后台</p>
    </div>
  </div>
</template>
