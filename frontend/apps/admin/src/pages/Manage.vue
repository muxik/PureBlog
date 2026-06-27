<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import type { PostListResponse } from '@pureblog/api-types'
import { api } from '../api'
import { useAuthStore } from '../stores/auth'

const list = ref<PostListResponse | null>(null)
const router = useRouter()
const auth = useAuthStore()

async function load() {
  list.value = await api<PostListResponse>('/admin/posts?pageSize=50')
}

async function remove(id?: number) {
  if (!id) return
  if (!confirm('确认删除这篇文章?')) return
  await api(`/admin/posts/${id}`, { method: 'DELETE' })
  await load()
}

function logout() {
  auth.logout()
  router.push('/login')
}

onMounted(load)
</script>

<template>
  <div class="manage">
    <header class="bar">
      <h1>文章</h1>
      <div>
        <RouterLink to="/write"><button>写文章</button></RouterLink>
        <button class="ghost" @click="logout">退出</button>
      </div>
    </header>

    <table class="list">
      <tbody>
        <tr v-for="p in list?.items ?? []" :key="p.id">
          <td class="t">{{ p.title }}</td>
          <td class="s">{{ p.status === 'published' ? '已发布' : '草稿' }}</td>
          <td class="a">
            <RouterLink :to="`/write/${p.id}`">编辑</RouterLink>
            <button class="ghost" @click="remove(p.id)">删除</button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<style scoped>
.bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.list {
  width: 100%;
  border-collapse: collapse;
  margin-top: 1rem;
}
.list td {
  padding: 0.7rem 0;
  border-bottom: 1px solid var(--border, #e7e2d6);
}
.s {
  color: var(--ink-3, #8a857a);
  width: 6rem;
}
.a {
  width: 9rem;
  text-align: right;
}
</style>
