<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import type { Post, PostListResponse, SavePostRequest } from '@pureblog/api-types'
import { api, downloadExport, importBundle, type ImportSummary } from '../api'

const router = useRouter()
const filter = ref<'all' | 'published' | 'draft'>('all')
const posts = ref<Post[]>([])

async function load() {
  const qs = filter.value === 'all' ? 'pageSize=100' : `status=${filter.value}&pageSize=100`
  const res = await api<PostListResponse>(`/admin/posts?${qs}`)
  posts.value = res.items ?? []
  // Drop any selections that are no longer in the list.
  selected.value = selected.value.filter((id) => posts.value.some((p) => p.id === id))
}

/* ── Multi-select ─────────────────────────────────────────────────── */
const selected = ref<number[]>([])
const busy = ref(false)

function isSelected(id?: number): boolean {
  return id != null && selected.value.includes(id)
}
function toggleSelect(id?: number) {
  if (id == null) return
  selected.value = isSelected(id)
    ? selected.value.filter((x) => x !== id)
    : [...selected.value, id]
}
const allSelected = computed(
  () => sorted.value.length > 0 && sorted.value.every((p) => isSelected(p.id)),
)
function toggleAll() {
  selected.value = allSelected.value
    ? []
    : sorted.value.map((p) => p.id).filter((id): id is number => id != null)
}

/* ── Batch actions ────────────────────────────────────────────────── */
async function batch(action: 'delete' | 'publish' | 'draft') {
  if (!selected.value.length) return
  if (action === 'delete' && !confirm(`确认删除选中的 ${selected.value.length} 篇文章?`)) return
  busy.value = true
  try {
    const res = await api<{ affected: number; failed: { id: number; reason: string }[] }>(
      '/admin/posts/batch',
      { method: 'POST', body: { action, ids: selected.value } },
    )
    selected.value = []
    await load()
    if (res?.failed?.length) {
      alert(`已处理 ${res.affected} 篇，${res.failed.length} 篇失败。`)
    }
  } catch {
    alert('批量操作失败，请重试。')
  } finally {
    busy.value = false
  }
}

/* ── Export / Import ──────────────────────────────────────────────── */
const importing = ref(false)
const importResult = ref<ImportSummary | null>(null)
const importInputEl = ref<HTMLInputElement | null>(null)

async function exportAll() {
  busy.value = true
  try {
    await downloadExport()
  } catch {
    alert('导出失败，请重试。')
  } finally {
    busy.value = false
  }
}
async function exportSelected() {
  if (!selected.value.length) return
  const slugs = posts.value
    .filter((p) => p.id != null && selected.value.includes(p.id))
    .map((p) => p.slug)
    .filter((s): s is string => !!s)
  busy.value = true
  try {
    await downloadExport(slugs)
  } catch {
    alert('导出失败，请重试。')
  } finally {
    busy.value = false
  }
}
async function onImportPicked(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  input.value = ''
  if (!file) return
  importing.value = true
  importResult.value = null
  try {
    importResult.value = await importBundle(file)
    await load()
  } catch {
    alert('导入失败：请确认上传的是导出的 .zip 包。')
  } finally {
    importing.value = false
  }
}

const sorted = computed(() =>
  posts.value.slice().sort((a, b) => (b.pinned ? 1 : 0) - (a.pinned ? 1 : 0))
)

function formatDateParts(p: Post): { full: string; label: string } {
  const raw = p.publishedAt || p.createdAt || ''
  if (!raw) return { full: '—', label: '—' }
  const d = new Date(raw)
  if (isNaN(d.getTime())) return { full: raw, label: raw }
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return {
    full: `${y}-${m}-${day}`,
    label: `${y}.${m}.${day}`,
  }
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
      <div class="manage-head__acts">
        <button class="act-btn" :disabled="busy" @click="exportAll()">导出全部 ↓</button>
        <button class="act-btn" :disabled="importing" @click="importInputEl?.click()">
          {{ importing ? '导入中…' : '导入 ↑' }}
        </button>
        <button class="add-link" @click="router.push('/write')">＋ 写新文章</button>
        <input
          ref="importInputEl"
          type="file"
          accept=".zip,application/zip"
          style="display:none"
          @change="onImportPicked"
        />
      </div>
    </div>

    <!-- Import result feedback -->
    <div v-if="importResult" class="import-note">
      导入完成：新增 {{ importResult.created }} · 更新 {{ importResult.updated }} · 图片 {{ importResult.images }}
      <template v-if="importResult.errors.length">
        · <span class="import-note__err">{{ importResult.errors.length }} 个失败</span>
        <ul class="import-note__list">
          <li v-for="(e, i) in importResult.errors" :key="i">{{ e.file }}：{{ e.reason }}</li>
        </ul>
      </template>
      <button class="import-note__close" @click="importResult = null">×</button>
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

    <!-- Select-all + batch action bar -->
    <div v-if="sorted.length" class="batch-bar">
      <label class="batch-bar__all">
        <input type="checkbox" :checked="allSelected" @change="toggleAll" />
        <span>{{ selected.length ? `已选 ${selected.length} 篇` : '全选' }}</span>
      </label>
      <template v-if="selected.length">
        <span class="batch-bar__spacer"></span>
        <button class="act-btn" :disabled="busy" @click="batch('publish')">发布</button>
        <button class="act-btn" :disabled="busy" @click="batch('draft')">转草稿</button>
        <button class="act-btn" :disabled="busy" @click="exportSelected()">导出选中 ↓</button>
        <button class="act-btn act-btn--del" :disabled="busy" @click="batch('delete')">删除</button>
      </template>
    </div>

    <div v-if="sorted.length" class="manage-list">
      <div v-for="p in sorted" :key="p.id" class="manage-row" :class="{ 'is-selected': isSelected(p.id) }">
        <input
          class="manage-row__check"
          type="checkbox"
          :checked="isSelected(p.id)"
          @change="toggleSelect(p.id)"
        />
        <div class="manage-row__main">
          <div class="manage-row__titlerow">
            <button class="manage-row__title" @click="edit(p.id)">{{ p.title || '无标题' }}</button>
            <span v-if="p.status === 'draft'" class="badge-draft">草稿</span>
            <span v-if="p.pinned" class="badge-pin">置顶</span>
          </div>
          <div class="manage-row__meta">
            <time
              class="manage-row__meta-date"
              :datetime="p.publishedAt || p.createdAt || undefined"
              :title="formatDateParts(p).full"
            >
              {{ formatDateParts(p).label }}
            </time>
            <span class="manage-row__tags">
              {{ (p.tags ?? []).map((t) => t.name).filter(Boolean).join(' · ') || '—' }}
            </span>
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
