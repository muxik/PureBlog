import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import type { LoginRequest, TokenResponse } from '@pureblog/api-types'
import { ACCESS_KEY, REFRESH_KEY, api } from '../api'

export const useAuthStore = defineStore('auth', () => {
  const accessToken = ref<string | null>(localStorage.getItem(ACCESS_KEY))
  const refreshToken = ref<string | null>(localStorage.getItem(REFRESH_KEY))
  const isAuthed = computed(() => !!accessToken.value)

  function persist() {
    if (accessToken.value) localStorage.setItem(ACCESS_KEY, accessToken.value)
    else localStorage.removeItem(ACCESS_KEY)
    if (refreshToken.value) localStorage.setItem(REFRESH_KEY, refreshToken.value)
    else localStorage.removeItem(REFRESH_KEY)
  }

  async function login(payload: LoginRequest) {
    const res = await api<TokenResponse>('/auth/login', {
      method: 'POST',
      body: payload,
      auth: false,
    })
    accessToken.value = res.accessToken ?? null
    refreshToken.value = res.refreshToken ?? null
    persist()
  }

  function logout() {
    accessToken.value = null
    refreshToken.value = null
    persist()
  }

  return { accessToken, refreshToken, isAuthed, login, logout }
})
