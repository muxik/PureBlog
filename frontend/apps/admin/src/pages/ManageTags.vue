<script setup lang="ts">
import { onMounted, ref } from 'vue'
import type { Tag, TagListResponse, SaveTagRequest } from '@pureblog/api-types'
import { api } from '../api'

const tags = ref<Tag[]>([])
const err = ref('')

const formName = ref('')
const formSlug = ref('')
const formErr = ref('')

async function load() {
  err.value = ''
  try {
    const res = await api<TagListResponse>('/tags', { auth: false })
    tags.value = res.items ?? []
  } catch (e) {
    err.value = e instanceof Error ? e.message : String(e)
  }
}

async function create() {
  formErr.value = ''
  if (!formName.value.trim()) {
    formErr.value = '名称不能为空'
    return
  }
  const body: SaveTagRequest = { name: formName.value.trim() }
  if (formSlug.value.trim()) body.slug = formSlug.value.trim()
  try {
    await api('/admin/tags', { method: 'POST', body })
    formName.value = ''
    formSlug.value = ''
    await load()
  } catch (e) {
    formErr.value = e instanceof Error ? e.message : String(e)
  }
}

async function remove(id?: number) {
  if (!id) return
  if (!confirm('确认删除这个标签?')) return
  err.value = ''
  try {
    await api(`/admin/tags/${id}`, { method: 'DELETE' })
    await load()
  } catch (e) {
    err.value = e instanceof Error ? e.message : String(e)
  }
}

onMounted(load)
</script>

<template>
  <div class="manage-tags">
    <RouterLink to="/manage" class="back">← 返回文章</RouterLink>
    <h1>标签管理</h1>

    <section class="form-section">
      <h2>新建标签</h2>
      <input v-model="formName" type="text" placeholder="名称（必填）" />
      <input v-model="formSlug" type="text" placeholder="slug（选填，留空自动生成）" />
      <p v-if="formErr" class="err">{{ formErr }}</p>
      <button @click="create">创建</button>
    </section>

    <p v-if="err" class="err">{{ err }}</p>

    <table class="list">
      <thead>
        <tr>
          <th class="th-name">名称</th>
          <th class="th-slug">slug</th>
          <th class="th-action">操作</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="tag in tags" :key="tag.id">
          <td>{{ tag.name }}</td>
          <td class="slug">{{ tag.slug }}</td>
          <td class="action">
            <button class="ghost" @click="remove(tag.id)">删除</button>
          </td>
        </tr>
        <tr v-if="tags.length === 0">
          <td colspan="3" class="empty">暂无标签</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<style scoped>
.back {
  display: inline-block;
  margin-bottom: 0.75rem;
  color: var(--accent, #235a73);
  text-decoration: none;
  font-size: 0.9rem;
}
.back:hover {
  text-decoration: underline;
}
h1 {
  margin-top: 0;
}
.form-section {
  margin-bottom: 2rem;
  padding: 1rem 1.25rem;
  border: 1px solid var(--border, #e7e2d6);
  border-radius: 6px;
  background: var(--paper, #faf8f2);
}
.form-section h2 {
  margin-top: 0;
  font-size: 1rem;
}
.list {
  width: 100%;
  border-collapse: collapse;
  margin-top: 1rem;
}
.list th,
.list td {
  padding: 0.7rem 0;
  border-bottom: 1px solid var(--border, #e7e2d6);
  text-align: left;
}
.list th {
  font-size: 0.85rem;
  color: var(--ink-2, #6a6560);
}
.th-slug,
.slug {
  color: var(--ink-2, #6a6560);
  width: 12rem;
}
.th-action,
.action {
  width: 6rem;
  text-align: right;
}
.empty {
  color: var(--ink-2, #6a6560);
  font-size: 0.9rem;
}
</style>
