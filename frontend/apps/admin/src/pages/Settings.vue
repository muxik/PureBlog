<script setup lang="ts">
import { onMounted, ref } from 'vue'
import type { SiteSettings, UpdateSettingsRequest } from '@pureblog/api-types'
import { api } from '../api'

interface SocialRow {
  key: string
  value: string
}

const siteName = ref('')
const author = ref('')
const description = ref('')
const aboutMd = ref('')
const defaultDateFormat = ref('numeric')
const socialRows = ref<SocialRow[]>([])

const error = ref('')
const success = ref('')

async function load() {
  try {
    const settings = await api<SiteSettings>('/settings', { auth: false })
    siteName.value = settings.siteName ?? ''
    author.value = settings.author ?? ''
    description.value = settings.description ?? ''
    aboutMd.value = settings.aboutMd ?? ''
    defaultDateFormat.value = settings.defaultDateFormat ?? 'numeric'
    const social = settings.social ?? {}
    socialRows.value = Object.entries(social).map(([key, value]) => ({ key, value }))
  } catch (e) {
    error.value = e instanceof Error ? e.message : '加载失败'
  }
}

function addSocialRow() {
  socialRows.value.push({ key: '', value: '' })
}

function removeSocialRow(index: number) {
  socialRows.value.splice(index, 1)
}

async function save() {
  error.value = ''
  success.value = ''
  const social: Record<string, string> = {}
  for (const row of socialRows.value) {
    const k = row.key.trim()
    const v = row.value.trim()
    if (k && v) {
      social[k] = v
    }
  }
  const body: UpdateSettingsRequest = {
    siteName: siteName.value,
    author: author.value,
    description: description.value,
    aboutMd: aboutMd.value,
    defaultDateFormat: defaultDateFormat.value,
    social,
  }
  try {
    await api<SiteSettings>('/admin/settings', { method: 'PUT', body })
    success.value = '已保存'
  } catch (e) {
    error.value = e instanceof Error ? e.message : '保存失败'
  }
}

onMounted(load)
</script>

<template>
  <div class="settings">
    <RouterLink to="/manage" class="back">← 返回文章</RouterLink>
    <h1>站点设置</h1>

    <form @submit.prevent="save">
      <label class="field-label">站名</label>
      <input v-model="siteName" type="text" placeholder="站名" />

      <label class="field-label">作者</label>
      <input v-model="author" type="text" placeholder="作者" />

      <label class="field-label">简介</label>
      <textarea v-model="description" rows="3" placeholder="站点简介" />

      <label class="field-label">关于页 Markdown</label>
      <textarea v-model="aboutMd" rows="8" placeholder="关于页内容（Markdown 源码）" class="mono" />

      <label class="field-label">日期格式</label>
      <select v-model="defaultDateFormat">
        <option value="numeric">公历(数字)</option>
        <option value="lunar">农历</option>
      </select>

      <div class="social-section">
        <div class="social-header">
          <span class="field-label">社交链接</span>
          <button type="button" class="ghost" @click="addSocialRow">添加</button>
        </div>
        <div v-for="(row, i) in socialRows" :key="i" class="social-row">
          <input v-model="row.key" type="text" placeholder="平台 (e.g. github)" class="social-key" />
          <input v-model="row.value" type="text" placeholder="链接 / 账号" class="social-val" />
          <button type="button" class="ghost del" @click="removeSocialRow(i)">删除</button>
        </div>
      </div>

      <div class="actions">
        <button type="submit">保存</button>
      </div>

      <p v-if="error" class="err">{{ error }}</p>
      <p v-if="success" class="ok">{{ success }}</p>
    </form>
  </div>
</template>

<style scoped>
.back {
  display: inline-block;
  margin-bottom: 1rem;
  color: var(--accent, #235a73);
  text-decoration: none;
}
.back:hover {
  text-decoration: underline;
}
.field-label {
  display: block;
  margin-top: 1rem;
  margin-bottom: 0.2rem;
  font-size: 0.875rem;
  color: var(--ink-2, #5a5650);
}
textarea,
select {
  display: block;
  width: 100%;
  padding: 0.5rem 0.7rem;
  margin: 0.4rem 0;
  border: 1px solid var(--border, #e7e2d6);
  border-radius: 4px;
  background: var(--paper, #faf8f2);
  color: inherit;
  font: inherit;
  box-sizing: border-box;
}
textarea {
  resize: vertical;
}
.mono {
  font-family: ui-monospace, 'Cascadia Code', 'Source Code Pro', monospace;
  font-size: 0.875rem;
}
.social-section {
  margin-top: 1rem;
}
.social-header {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-bottom: 0.4rem;
}
.social-header .field-label {
  margin: 0;
}
.social-row {
  display: flex;
  gap: 0.5rem;
  align-items: center;
  margin-bottom: 0.4rem;
}
.social-key {
  width: 9rem;
  flex-shrink: 0;
}
.social-val {
  flex: 1;
}
.del {
  flex-shrink: 0;
  color: #b23b3b;
  border-color: #b23b3b;
}
.actions {
  margin-top: 1.5rem;
}
.ok {
  color: var(--accent, #235a73);
  margin-top: 0.5rem;
}
</style>
