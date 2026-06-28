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
  <section class="admin-section">
    <div class="manage-head">
      <h1 class="section-title">评论审核</h1>
    </div>

    <!-- ── Status filter ── -->
    <div class="seg-group">
      <button
        class="seg"
        :class="{ 'is-active': filter === 'pending' }"
        @click="setFilter('pending')"
      >待审</button>
      <button
        class="seg"
        :class="{ 'is-active': filter === 'approved' }"
        @click="setFilter('approved')"
      >已通过</button>
      <button
        class="seg"
        :class="{ 'is-active': filter === 'all' }"
        @click="setFilter('all')"
      >全部</button>
    </div>

    <p v-if="error" class="cmt-err">{{ error }}</p>

    <!-- ── Empty state ── -->
    <p v-if="!comments.length && !error" class="manage-empty">
      {{
        filter === 'pending'
          ? '没有待审评论'
          : filter === 'approved'
          ? '没有已通过评论'
          : '没有评论'
      }}
    </p>

    <!-- ── Comment list ── -->
    <div class="manage-list">
      <div v-for="c in comments" :key="c.id" class="manage-row cmt-row">
        <div class="manage-row__main">
          <div class="manage-row__titlerow">
            <span class="manage-row__title cmt-author">{{ c.authorName }}</span>
            <span v-if="c.status !== 'approved'" class="badge-draft">待审</span>
            <span v-else class="badge-pin">已通过</span>
          </div>
          <div class="manage-row__tags">
            {{ c.createdAt ? new Date(c.createdAt).toLocaleString('zh-CN') : '' }}
            · 文章 #{{ c.postId }}
            <template v-if="c.parentId"> · 回复 #{{ c.parentId }}</template>
          </div>
          <div class="cmt-body">{{ c.content }}</div>
        </div>
        <div class="manage-row__acts">
          <button v-if="c.status !== 'approved'" class="act-btn" @click="moderate(c.id, 'approved')">通过</button>
          <button v-if="c.status === 'approved'" class="act-btn" @click="moderate(c.id, 'pending')">撤回</button>
          <button class="act-btn act-btn--del" @click="remove(c.id)">删除</button>
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped>
/* Comments have multi-line content so align to top, not baseline */
.cmt-row {
  align-items: flex-start;
}

/* Author name is display only — override pointer cursor from global manage-row__title */
.cmt-author {
  cursor: default;
  font-family: var(--font-body);
  font-size: var(--text-base);
}

/* Comment body text */
.cmt-body {
  margin-top: 8px;
  font-size: var(--text-base);
  line-height: 1.6;
  color: var(--text-body);
  white-space: pre-wrap;
  word-break: break-word;
}

.cmt-err {
  font-size: var(--text-sm);
  color: #b23b3b;
  margin: 12px 0 0;
}
</style>
