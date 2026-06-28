<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import type { Post, PostListResponse, SavePostRequest } from '@pureblog/api-types'
import { api } from '../api'

const router = useRouter()
const filter = ref<'all' | 'published' | 'draft'>('all')
const posts = ref<Post[]>([])

async function load() {
  const qs = filter.value === 'all' ? 'pageSize=200' : `status=${filter.value}&pageSize=200`
  const res = await api<PostListResponse>(`/admin/posts?${qs}`)
  posts.value = res.items ?? []
}

const sorted = computed(() =>
  posts.value.slice().sort((a, b) => (b.pinned ? 1 : 0) - (a.pinned ? 1 : 0))
)

function formatDate(p: Post): string {
  const raw = p.publishedAt || p.createdAt || ''
  if (!raw) return '—'
  const d = new Date(raw)
  if (isNaN(d.getTime())) return raw
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${y} · ${m} · ${day}`
}

async function setFilter(f: 'all' | 'published' | 'draft') {
  filter.value = f
  await load()
}

async function remove(id?: number) {
  if (!id) return
  if (!confirm('确认删除这篇文章?')) return
  await api(`/admin/posts/${id}`, { method: 'DELETE' })
  await load()
}

function edit(id?: number) {
  if (!id) return
  router.push(`/write/${id}`)
}

async function togglePin(id?: number) {
  if (!id) return
  const post = await api<Post>(`/admin/posts/${id}`)
  const body: SavePostRequest = {
    title: post.title ?? '',
    slug: post.slug,
    summary: post.summary,
    contentMd: post.contentMd,
    coverUrl: post.coverUrl,
    status: post.status,
    tagIds: (post.tags ?? []).map((t) => t.id).filter((id): id is number => id !== undefined),
    pinned: !post.pinned,
  }
  await api(`/admin/posts/${id}`, { method: 'PUT', body })
  await load()
}

onMounted(load)
</script>

<template>
  <section class="admin-section">
    <div class="manage-head">
      <h1 class="section-title">文章</h1>
      <button class="add-link" @click="router.push('/write')">＋ 写新文章</button>
    </div>

    <div class="seg-group">
      <button
        class="seg"
        :class="{ 'is-active': filter === 'all' }"
        @click="setFilter('all')"
      >全部</button>
      <button
        class="seg"
        :class="{ 'is-active': filter === 'published' }"
        @click="setFilter('published')"
      >已发</button>
      <button
        class="seg"
        :class="{ 'is-active': filter === 'draft' }"
        @click="setFilter('draft')"
      >草稿</button>
    </div>

    <p v-if="sorted.length === 0" class="manage-empty">
      这里还没有文章。<button @click="router.push('/write')">写第一篇 ›</button>
    </p>

    <div v-else class="manage-list">
      <div v-for="p in sorted" :key="p.id" class="manage-row">
        <time class="manage-row__date">{{ formatDate(p) }}</time>
        <div class="manage-row__main">
          <div class="manage-row__titlerow">
            <button class="manage-row__title" @click="edit(p.id)">{{ p.title || '无标题' }}</button>
            <span v-if="p.status === 'draft'" class="badge-draft">草稿</span>
            <span v-if="p.pinned" class="badge-pin">置顶</span>
          </div>
          <div class="manage-row__tags">
            {{ (p.tags ?? []).map((t) => t.name).filter(Boolean).join(' · ') || '—' }}
          </div>
        </div>
        <div class="manage-row__acts">
          <button class="act-btn" @click="togglePin(p.id)">{{ p.pinned ? '取消置顶' : '置顶' }}</button>
          <button class="act-btn" @click="edit(p.id)">编辑</button>
          <button class="act-btn act-btn--del" @click="remove(p.id)">删除</button>
        </div>
      </div>
    </div>
  </section>
</template>
