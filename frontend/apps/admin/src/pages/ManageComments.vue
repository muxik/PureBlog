<script setup lang="ts">
import { onMounted, ref } from 'vue'
import type { Comment, CommentListResponse } from '@pureblog/api-types'
import { api } from '../api'

type FilterStatus = 'pending' | 'approved' | 'all'

const comments = ref<Comment[]>([])
const filter = ref<FilterStatus>('pending')
const error = ref('')

async function load() {
  error.value = ''
  try {
    const path =
      filter.value === 'all'
        ? '/admin/comments'
        : `/admin/comments?status=${filter.value}`
    const res = await api<CommentListResponse>(path)
    comments.value = res.items ?? []
  } catch (e) {
    error.value = (e as Error).message
  }
}

async function moderate(id: number | undefined, status: 'approved' | 'pending') {
  if (!id) return
  error.value = ''
  try {
    await api<undefined>(`/admin/comments/${id}`, { method: 'PUT', body: { status } })
    await load()
  } catch (e) {
    error.value = (e as Error).message
  }
}

async function remove(id: number | undefined) {
  if (!id) return
  if (!confirm('确认删除这条评论?')) return
  error.value = ''
  try {
    await api<undefined>(`/admin/comments/${id}`, { method: 'DELETE' })
    await load()
  } catch (e) {
    error.value = (e as Error).message
  }
}

function setFilter(f: FilterStatus) {
  filter.value = f
  load()
}

onMounted(load)
</script>

<template>
  <div class="manage-comments">
    <RouterLink to="/manage" class="back">← 返回文章</RouterLink>
    <h1>评论审核</h1>

    <div class="filters">
      <button
        :class="{ active: filter === 'pending' }"
        @click="setFilter('pending')"
      >待审</button>
      <button
        :class="{ active: filter === 'approved' }"
        @click="setFilter('approved')"
      >已通过</button>
      <button
        :class="{ active: filter === 'all' }"
        @click="setFilter('all')"
      >全部</button>
    </div>

    <p v-if="error" class="err">{{ error }}</p>

    <p v-if="!comments.length && !error" class="empty">
      {{
        filter === 'pending'
          ? '没有待审评论'
          : filter === 'approved'
          ? '没有已通过评论'
          : '没有评论'
      }}
    </p>

    <div v-for="c in comments" :key="c.id" class="comment-card">
      <div class="comment-meta">
        <span class="author">{{ c.authorName }}</span>
        <span
          :class="['badge', c.status === 'approved' ? 'badge--approved' : 'badge--pending']"
        >{{ c.status === 'approved' ? '已通过' : '待审' }}</span>
        <span class="date">{{
          c.createdAt ? new Date(c.createdAt).toLocaleString('zh-CN') : ''
        }}</span>
        <span class="post-ref">文章 #{{ c.postId }}</span>
        <span v-if="c.parentId" class="reply-ref">↳ 回复 #{{ c.parentId }}</span>
      </div>
      <div class="comment-content">{{ c.content }}</div>
      <div class="comment-actions">
        <button v-if="c.status !== 'approved'" @click="moderate(c.id, 'approved')">通过</button>
        <button v-if="c.status === 'approved'" @click="moderate(c.id, 'pending')">撤回</button>
        <button class="ghost" @click="remove(c.id)">删除</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.back {
  display: inline-block;
  margin-bottom: 0.5rem;
  color: var(--accent, #235a73);
  text-decoration: none;
  font-size: 0.9rem;
}
.back:hover {
  text-decoration: underline;
}
.filters {
  display: flex;
  margin: 1rem 0;
  border: 1px solid var(--border, #e7e2d6);
  border-radius: 4px;
  overflow: hidden;
  width: fit-content;
}
.filters button {
  border: none;
  border-radius: 0;
  margin: 0;
  background: var(--paper, #faf8f2);
  color: var(--accent, #235a73);
  padding: 0.4rem 1rem;
}
.filters button + button {
  border-left: 1px solid var(--border, #e7e2d6);
}
.filters button.active {
  background: var(--accent, #235a73);
  color: #fff;
}
.empty {
  color: var(--ink-2, #8a857a);
  margin-top: 1.5rem;
}
.comment-card {
  border: 1px solid var(--border, #e7e2d6);
  border-radius: 6px;
  padding: 0.9rem 1rem;
  margin-top: 1rem;
  background: var(--paper-subtle, #f5f2ea);
}
.comment-meta {
  display: flex;
  align-items: center;
  gap: 0.7rem;
  flex-wrap: wrap;
  margin-bottom: 0.5rem;
  font-size: 0.85rem;
}
.author {
  font-weight: 600;
}
.badge {
  padding: 0.15rem 0.5rem;
  border-radius: 99px;
  font-size: 0.75rem;
}
.badge--pending {
  background: #fef3c7;
  color: #92400e;
}
.badge--approved {
  background: #d1fae5;
  color: #065f46;
}
.date {
  color: var(--ink-2, #8a857a);
}
.post-ref {
  color: var(--ink-2, #8a857a);
}
.reply-ref {
  color: var(--accent, #235a73);
  font-size: 0.8rem;
}
.comment-content {
  margin-bottom: 0.75rem;
  line-height: 1.5;
  white-space: pre-wrap;
  word-break: break-word;
}
.comment-actions {
  display: flex;
}
</style>
