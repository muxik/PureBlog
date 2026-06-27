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
    error.value = '用户名或密码错误'
  }
}
</script>

<template>
  <div class="login">
    <h1>PureBlog 后台</h1>
    <form @submit.prevent="submit">
      <input v-model="username" placeholder="用户名" autocomplete="username" />
      <input
        v-model="password"
        type="password"
        placeholder="密码"
        autocomplete="current-password"
      />
      <button type="submit">登录</button>
      <p v-if="error" class="err">{{ error }}</p>
    </form>
  </div>
</template>
