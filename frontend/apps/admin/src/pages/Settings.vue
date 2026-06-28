<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import type { SiteSettings, UpdateSettingsRequest } from '@pureblog/api-types'
import { api } from '../api'
import { short } from '../utils/lunar'

// ── Active tab ────────────────────────────────────────────────
const tab = ref<'site' | 'about' | 'social' | 'account'>('site')

// ── Loaded fields ─────────────────────────────────────────────
const siteName    = ref('')
const description = ref('')
const author      = ref('')   // preserved on save; not shown in UI
const aboutMd     = ref('')
const dateFormat  = ref('numeric')

// Social — fixed four keys
const socialEmail  = ref('')
const socialGithub = ref('')
const socialX      = ref('')
const socialWeibo  = ref('')

// Account tab — display only (no backend endpoint)
const acctUser = ref('')
const acctPass = ref('')

// ── Feedback ──────────────────────────────────────────────────
const saveError   = ref('')
const saveSuccess = ref('')
let successTimer: ReturnType<typeof setTimeout> | null = null

// ── Load ──────────────────────────────────────────────────────
async function load() {
  try {
    const s = await api<SiteSettings>('/settings', { auth: false })
    siteName.value    = s.siteName    ?? ''
    description.value = s.description ?? ''
    author.value      = s.author      ?? ''
    aboutMd.value     = s.aboutMd     ?? ''
    dateFormat.value  = s.defaultDateFormat ?? 'numeric'
    const soc = s.social ?? {}
    socialEmail.value  = soc['email']  ?? ''
    socialGithub.value = soc['github'] ?? ''
    socialX.value      = soc['x']      ?? ''
    socialWeibo.value  = soc['weibo']  ?? ''
  } catch (e) {
    saveError.value = e instanceof Error ? e.message : '加载失败'
  }
}

// ── Date format (writes to localStorage immediately) ──────────
function setDateFormat(fmt: string) {
  dateFormat.value = fmt
  try { localStorage.setItem('muxi:dateFormat', fmt) } catch { /* ignore */ }
}

const DEMO = '2026 · 03 · 14'
const fmtPreview = computed(() => {
  if (dateFormat.value === 'lunar') {
    const l = short(DEMO)
    if (l) return `示例：${l}（2026 年 3 月 14 日）`
  }
  return `示例：${DEMO}`
})

// ── Save (PUT /admin/settings with full merged body) ──────────
function showSuccess() {
  saveSuccess.value = '已保存 ✓'
  saveError.value   = ''
  if (successTimer) clearTimeout(successTimer)
  successTimer = setTimeout(() => { saveSuccess.value = '' }, 2500)
}

async function save() {
  saveError.value   = ''
  saveSuccess.value = ''

  const social: Record<string, string> = {}
  if (socialEmail.value.trim())  social['email']  = socialEmail.value.trim()
  if (socialGithub.value.trim()) social['github'] = socialGithub.value.trim()
  if (socialX.value.trim())      social['x']      = socialX.value.trim()
  if (socialWeibo.value.trim())  social['weibo']  = socialWeibo.value.trim()

  const body: UpdateSettingsRequest = {
    siteName:          siteName.value,
    description:       description.value,
    author:            author.value,
    aboutMd:           aboutMd.value,
    defaultDateFormat: dateFormat.value,
    social,
  }

  try {
    await api<SiteSettings>('/admin/settings', { method: 'PUT', body })
    showSuccess()
  } catch (e) {
    saveError.value = e instanceof Error ? e.message : '保存失败'
  }
}

onMounted(load)
</script>

<template>
  <section class="admin-section">
    <h1 class="settings-title">设置</h1>

    <div class="settings-grid">
      <!-- ── Left nav ───────────────────────────────────────── -->
      <nav class="settings-nav">
        <button
          class="set-tab"
          :class="{ 'is-active': tab === 'site' }"
          @click="tab = 'site'"
        >
          <span>站点</span>
          <span class="set-tab__hint">站名 · 简介</span>
        </button>
        <button
          class="set-tab"
          :class="{ 'is-active': tab === 'about' }"
          @click="tab = 'about'"
        >
          <span>关于页</span>
          <span class="set-tab__hint">自我介绍</span>
        </button>
        <button
          class="set-tab"
          :class="{ 'is-active': tab === 'social' }"
          @click="tab = 'social'"
        >
          <span>社交</span>
          <span class="set-tab__hint">联系方式</span>
        </button>
        <button
          class="set-tab"
          :class="{ 'is-active': tab === 'account' }"
          @click="tab = 'account'"
        >
          <span>账户</span>
          <span class="set-tab__hint">登录凭据</span>
        </button>
      </nav>

      <!-- ── Right panel ────────────────────────────────────── -->
      <div class="settings-panel">

        <!-- 站点 -->
        <div v-if="tab === 'site'" class="panel-narrow">
          <h2 class="panel-h">站点</h2>
          <p class="panel-desc">站名与首页的开场白。</p>
          <div class="field-stack">
            <label class="field">
              <span class="field-label">站名</span>
              <input
                v-model="siteName"
                class="admin-input field-name-input"
                type="text"
              />
            </label>
            <label class="field">
              <span class="field-label">首页简介</span>
              <textarea
                v-model="description"
                class="admin-body"
                style="min-height: 80px"
              />
            </label>
            <div>
              <span class="field-label">首页日期格式</span>
              <div class="seg-group" style="margin: 0">
                <button
                  class="seg seg--wide"
                  :class="{ 'is-active': dateFormat === 'numeric' }"
                  @click="setDateFormat('numeric')"
                >数字</button>
                <button
                  class="seg seg--wide"
                  :class="{ 'is-active': dateFormat === 'lunar' }"
                  @click="setDateFormat('lunar')"
                >农历</button>
              </div>
              <p class="fmt-preview">{{ fmtPreview }}</p>
            </div>
          </div>
          <div class="save-row">
            <button class="btn-solid--save" @click="save">保存</button>
            <span v-if="saveSuccess" class="feedback-ok">{{ saveSuccess }}</span>
            <span v-if="saveError"   class="feedback-err">{{ saveError }}</span>
          </div>
        </div>

        <!-- 关于页 -->
        <div v-else-if="tab === 'about'" class="panel-wide">
          <h2 class="panel-h">关于页</h2>
          <p class="panel-desc">这页讲你是谁。支持 Markdown。</p>
          <label class="field">
            <textarea
              v-model="aboutMd"
              class="admin-body"
              style="min-height: 280px; line-height: 1.85; padding: 14px"
            />
          </label>
          <div class="save-row">
            <button class="btn-solid--save" @click="save">保存</button>
            <span v-if="saveSuccess" class="feedback-ok">{{ saveSuccess }}</span>
            <span v-if="saveError"   class="feedback-err">{{ saveError }}</span>
          </div>
        </div>

        <!-- 社交 -->
        <div v-else-if="tab === 'social'" class="panel-narrow">
          <h2 class="panel-h">社交</h2>
          <p class="panel-desc">联系方式，会显示在博客页脚的落款处。留空则不显示。</p>
          <div class="field-stack field-stack--tight">
            <label class="field">
              <span class="field-label">邮箱</span>
              <input
                v-model="socialEmail"
                class="admin-input"
                type="text"
                placeholder="hi@example.com"
              />
            </label>
            <label class="field">
              <span class="field-label">GitHub</span>
              <input
                v-model="socialGithub"
                class="admin-input"
                type="text"
                placeholder="github.com/用户名"
              />
            </label>
            <label class="field">
              <span class="field-label">X · Twitter</span>
              <input
                v-model="socialX"
                class="admin-input"
                type="text"
                placeholder="@用户名"
              />
            </label>
            <label class="field">
              <span class="field-label">微博</span>
              <input
                v-model="socialWeibo"
                class="admin-input"
                type="text"
                placeholder="weibo.com/用户名"
              />
            </label>
          </div>
          <div class="save-row">
            <button class="btn-solid--save" @click="save">保存</button>
            <span v-if="saveSuccess" class="feedback-ok">{{ saveSuccess }}</span>
            <span v-if="saveError"   class="feedback-err">{{ saveError }}</span>
          </div>
        </div>

        <!-- 账户 -->
        <div v-else-if="tab === 'account'" class="panel-narrow">
          <h2 class="panel-h">账户</h2>
          <p class="panel-desc">登录写作后台的凭据。</p>
          <div class="field-stack field-stack--tight">
            <label class="field">
              <span class="field-label">用户名</span>
              <input
                v-model="acctUser"
                class="admin-input"
                type="text"
              />
            </label>
            <label class="field">
              <span class="field-label">新密码</span>
              <input
                v-model="acctPass"
                class="admin-input"
                type="password"
                placeholder="留空则不修改"
              />
            </label>
          </div>
          <div class="save-row">
            <button class="btn-ghost" disabled>更新账户</button>
            <p class="acct-notice">暂不支持修改账户（后端未提供接口）</p>
          </div>
        </div>

      </div><!-- /settings-panel -->
    </div><!-- /settings-grid -->
  </section>
</template>

<style scoped>
.feedback-ok {
  margin-left: 14px;
  font-family: var(--font-mono);
  font-size: var(--text-xs);
  color: var(--accent);
}
.feedback-err {
  margin-left: 14px;
  font-family: var(--font-body);
  font-size: var(--text-sm);
  color: #c0392b;
}
.acct-notice {
  font-family: var(--font-body);
  font-size: var(--text-sm);
  color: var(--text-faint);
  margin: 10px 0 0;
}
</style>
