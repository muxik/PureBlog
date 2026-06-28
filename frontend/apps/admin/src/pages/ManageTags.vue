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
  <section class="admin-section">
    <div class="manage-head">
      <h1 class="section-title">标签管理</h1>
    </div>

    <!-- ── Create Form ── -->
    <div class="tag-form">
      <h2 class="panel-h">新建标签</h2>
      <p v-if="formErr" class="tag-err">{{ formErr }}</p>
      <div class="field-stack field-stack--tight">
        <label class="field">
          <span class="field-label">名称<span class="tag-req"> *</span></span>
          <input class="admin-input" v-model="formName" type="text" placeholder="名称（必填）" />
        </label>
        <label class="field">
          <span class="field-label">Slug（留空自动生成）</span>
          <input class="admin-input" v-model="formSlug" type="text" placeholder="url-friendly-slug" />
        </label>
      </div>
      <div class="save-row">
        <button class="btn-solid--save" @click="create">创建</button>
      </div>
    </div>

    <!-- List-level error -->
    <p v-if="err" class="tag-err tag-err--list">{{ err }}</p>

    <!-- ── Tag List ── -->
    <div class="manage-list">
      <div v-for="tag in tags" :key="tag.id" class="manage-row">
        <div class="manage-row__main">
          <div class="manage-row__titlerow">
            <span class="manage-row__title tag-name">{{ tag.name }}</span>
          </div>
          <div class="manage-row__tags">{{ tag.slug }}</div>
        </div>
        <div class="manage-row__acts">
          <button class="act-btn act-btn--del" @click="remove(tag.id)">删除</button>
        </div>
      </div>
      <p v-if="tags.length === 0" class="manage-empty">暂无标签</p>
    </div>
  </section>
</template>

<style scoped>
.tag-form {
  background: var(--bg-subtle);
  border: 1px solid var(--line);
  border-radius: var(--radius-sm, 4px);
  padding: 20px 24px 24px;
  margin: 24px 0 32px;
  max-width: 480px;
}

.tag-err {
  font-size: var(--text-sm);
  color: #b23b3b;
  margin: 0 0 12px;
}

.tag-err--list {
  margin: 0 0 16px;
}

.tag-req {
  color: #b23b3b;
}

/* Tag names are not clickable — override the pointer cursor from the global class */
.tag-name {
  cursor: default;
  font-family: var(--font-body);
}
</style>
