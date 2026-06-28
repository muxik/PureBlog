<script setup lang="ts">
import { nextTick, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import type {
  Post,
  SavePostRequest,
  Tag,
  TagListResponse,
  RenderResponse,
} from '@pureblog/api-types'
import { api } from '../api'

const route = useRoute()
const router = useRouter()
const id = route.params.id ? Number(route.params.id) : null

/* ── Core editor state ───────────────────────────────────────────── */
const title = ref('')
const body = ref('')        // synced from textarea on input / before save / before preview
const wordCount = ref(0)
const postStatus = ref<'draft' | 'published'>('draft')
const preview = ref(false)
const previewHtml = ref('')
const drawerOpen = ref(false)

/* ── Drawer / publish-settings state ─────────────────────────────── */
const summary = ref('')
const coverUrl = ref('')
const pinned = ref(false)
const tagInput = ref('')        // raw space-separated names shown in the input
const tagChips = ref<string[]>([])  // parsed chips shown below the input

/* ── DOM refs ─────────────────────────────────────────────────────── */
// We do NOT bind :value on the textarea — Vue never patches value during
// re-renders, so Chinese IME composition and caret position are preserved.
const textareaEl = ref<HTMLTextAreaElement | null>(null)

/* ── Display date ─────────────────────────────────────────────────── */
function formatDate(d: Date): string {
  const pad = (n: number) => (n < 10 ? '0' + n : String(n))
  return `${d.getFullYear()} · ${pad(d.getMonth() + 1)} · ${pad(d.getDate())}`
}
const displayDate = ref(formatDate(new Date()))

/* ── Body input handler (no structural re-render) ─────────────────── */
function onBodyInput(e: Event) {
  const val = (e.target as HTMLTextAreaElement).value
  body.value = val
  wordCount.value = val.replace(/\s/g, '').length
}

/* ── Toolbar helpers (operate on live textarea, no re-render) ──────── */
function wrap(before: string, after: string) {
  const t = textareaEl.value
  if (!t) {
    body.value = (body.value || '') + before + after
    return
  }
  const s = t.selectionStart
  const e = t.selectionEnd
  const v = t.value
  const sel = v.slice(s, e)
  t.value = v.slice(0, s) + before + sel + after + v.slice(e)
  body.value = t.value
  wordCount.value = t.value.replace(/\s/g, '').length
  t.focus()
  try { t.setSelectionRange(s + before.length, s + before.length + sel.length) } catch {}
}

function linePrefix(prefix: string) {
  const t = textareaEl.value
  if (!t) {
    body.value = prefix + (body.value || '')
    return
  }
  const s = t.selectionStart
  const v = t.value
  const lineStart = v.lastIndexOf('\n', s - 1) + 1
  t.value = v.slice(0, lineStart) + prefix + v.slice(lineStart)
  body.value = t.value
  wordCount.value = t.value.replace(/\s/g, '').length
  t.focus()
  try { t.setSelectionRange(s + prefix.length, s + prefix.length) } catch {}
}

function insertImage() {
  // No upload backend — insert a URL placeholder template
  wrap('![', '](https://)')
}

/* ── Preview toggle (structural — re-renders) ─────────────────────── */
async function togglePreview() {
  if (!preview.value) {
    // Sync body from textarea before switching away from it
    if (textareaEl.value) body.value = textareaEl.value.value
    try {
      const res = await api<RenderResponse>('/admin/render', {
        method: 'POST',
        body: { markdown: body.value },
      })
      previewHtml.value = res.html ?? ''
    } catch {
      previewHtml.value = '<p style="color:var(--accent)">渲染失败，请检查网络或重试。</p>'
    }
    preview.value = true
  } else {
    preview.value = false
    // After re-render the textarea is mounted again; restore its value
    await nextTick()
    if (textareaEl.value) textareaEl.value.value = body.value
  }
}

/* ── Tag input ────────────────────────────────────────────────────── */
function onTagInput(e: Event) {
  const val = (e.target as HTMLInputElement).value
  tagInput.value = val
  tagChips.value = val.split(/[\s,，、]+/).filter(Boolean)
}

/* ── Drawer ───────────────────────────────────────────────────────── */
function openDrawer() { drawerOpen.value = true }
function closeDrawer() { drawerOpen.value = false }

/* ── Mount: load data ─────────────────────────────────────────────── */
onMounted(async () => {
  const postData = id ? await api<Post>(`/admin/posts/${id}`) : null

  let initialBody = ''
  if (postData) {
    title.value = postData.title ?? ''
    summary.value = postData.summary ?? ''
    coverUrl.value = postData.coverUrl ?? ''
    pinned.value = postData.pinned ?? false
    postStatus.value = (postData.status as 'draft' | 'published') ?? 'draft'
    initialBody = postData.contentMd ?? ''

    // Populate tag state from post.tags (name-based)
    const postTagNames = (postData.tags ?? [])
      .map((t) => t.name)
      .filter((n): n is string => !!n)
    tagChips.value = postTagNames
    tagInput.value = postTagNames.join(' ')

    if (postData.publishedAt) {
      displayDate.value = formatDate(new Date(postData.publishedAt))
    }
  }

  body.value = initialBody
  wordCount.value = initialBody.replace(/\s/g, '').length

  // Set textarea value after Vue has mounted the element
  await nextTick()
  if (textareaEl.value) textareaEl.value.value = initialBody
})

/* ── Tag name → ID resolution ─────────────────────────────────────── */
// Matches each name against existing tags; creates missing ones via
// POST /admin/tags { name }, then collects all resolved IDs.
async function resolveTagIds(): Promise<number[]> {
  const names = tagChips.value
  if (!names.length) return []

  // Refresh the canonical tag list
  const res = await api<TagListResponse>('/tags', { auth: false })
  const current = res.items ?? []

  const ids: number[] = []
  for (const name of names) {
    const existing = current.find((t) => t.name === name)
    if (existing?.id != null) {
      ids.push(Number(existing.id))
    } else {
      // Create the missing tag, then add its new ID
      const created = await api<Tag>('/admin/tags', {
        method: 'POST',
        body: { name },
      })
      if (created?.id != null) ids.push(Number(created.id))
    }
  }
  return ids
}

/* ── Save / publish ───────────────────────────────────────────────── */
async function save(status: 'draft' | 'published') {
  // Always read the current textarea value before saving
  if (textareaEl.value) body.value = textareaEl.value.value

  let resolvedTagIds: number[] = []
  try {
    resolvedTagIds = await resolveTagIds()
  } catch {
    // Non-fatal: save without tag IDs rather than blocking the user
  }

  const req: SavePostRequest = {
    title: title.value,
    summary: summary.value,
    contentMd: body.value,
    coverUrl: coverUrl.value,
    status,
    pinned: pinned.value,
    tagIds: resolvedTagIds,
  }

  if (id) {
    await api(`/admin/posts/${id}`, { method: 'PUT', body: req })
  } else {
    await api('/admin/posts', { method: 'POST', body: req })
  }

  router.push('/manage')
}

async function drawerSave(status: 'draft' | 'published') {
  drawerOpen.value = false
  await save(status)
}
</script>

<template>
  <!-- ── Toolbar ───────────────────────────────────────────────────── -->
  <div class="toolbar-bar">
    <div class="toolbar-inner">
      <button class="tool-btn" title="标题" @click="linePrefix('## ')">H</button>
      <button class="tool-btn" title="粗体" @click="wrap('**', '**')">B</button>
      <button class="tool-btn" title="斜体" @click="wrap('*', '*')"><em>i</em></button>
      <button class="tool-btn" title="引用" @click="linePrefix('> ')">"</button>
      <button class="tool-btn" title="列表" @click="linePrefix('- ')">•</button>
      <button class="tool-btn" title="代码" @click="wrap('`', '`')">‹›</button>
      <button class="tool-btn" title="链接" @click="wrap('[', '](https://)')">链</button>
      <button class="tool-btn" title="插入图片" @click="insertImage()">图</button>
      <span class="toolbar-spacer"></span>
      <button
        class="preview-btn"
        :class="{ 'is-on': preview }"
        @click="togglePreview"
      >{{ preview ? '继续写 ✎' : '预览 ◹' }}</button>
    </div>
  </div>

  <!-- ── Editor area ───────────────────────────────────────────────── -->
  <div class="editor">
    <!-- Preview mode -->
    <template v-if="preview">
      <div class="preview-block">
        <div class="preview-label">预览</div>
        <h1 class="preview-title">{{ title || '无标题' }}</h1>
        <div class="prose" style="margin:0" v-html="previewHtml"></div>
      </div>
    </template>

    <!-- Edit mode -->
    <template v-else>
      <input
        v-model="title"
        class="editor__title"
        type="text"
        placeholder="无标题"
      />
      <div class="editor__meta">
        <span>{{ displayDate }}</span>
        <span>·</span>
        <span>{{ postStatus === 'published' ? '已发布' : '草稿' }}</span>
        <span>·</span>
        <span>{{ wordCount }} 字</span>
      </div>
      <!-- No :value binding — Vue never patches this element's value during
           re-renders, so Chinese IME composition and caret are never broken. -->
      <textarea
        ref="textareaEl"
        class="editor__body"
        placeholder="开始写。Markdown 可用。慢一点，长一点。"
        @input="onBodyInput"
      ></textarea>
    </template>
  </div>

  <!-- ── Edit bar (bottom sticky) ──────────────────────────────────── -->
  <div class="editbar">
    <div class="editbar-inner">
      <button class="linkbtn" @click="openDrawer">发布设置 ›</button>
      <span class="toolbar-spacer"></span>
      <button class="btn-ghost" @click="save('draft')">存草稿</button>
      <button class="btn-solid btn-solid--sm" @click="save('published')">
        {{ postStatus === 'published' ? '更新' : '发布' }}
      </button>
    </div>
  </div>

  <!-- ── 发布设置 drawer ────────────────────────────────────────────── -->
  <template v-if="drawerOpen">
    <div class="scrim" @click="closeDrawer"></div>
    <aside class="drawer">
      <div class="drawer-head">
        <span class="drawer-title">发布设置</span>
        <button class="drawer-close" aria-label="关闭" @click="closeDrawer">×</button>
      </div>

      <div class="drawer-body">
        <!-- 配图 / cover -->
        <label class="field">
          <span class="field-label">配图</span>
          <div class="cover-slot">
            <img v-if="coverUrl" :src="coverUrl" alt="封面预览" />
            <template v-else>封面图 URL</template>
          </div>
          <input
            v-model="coverUrl"
            class="admin-input"
            type="text"
            placeholder="https://example.com/cover.jpg"
            style="margin-top:8px"
          />
        </label>

        <!-- 标签 / tags -->
        <label class="field">
          <span class="field-label">标签</span>
          <input
            class="admin-input"
            type="text"
            placeholder="排版 中文 随笔（空格分隔）"
            :value="tagInput"
            @input="onTagInput"
          />
          <div class="tag-chips">
            <span v-for="chip in tagChips" :key="chip" class="tag-chip">{{ chip }}</span>
          </div>
        </label>

        <!-- 日期 / date (read-only; backend sets publishedAt) -->
        <label class="field">
          <span class="field-label">日期</span>
          <input
            class="admin-input"
            type="text"
            :value="displayDate"
            placeholder="发布时自动"
            disabled
            readonly
          />
        </label>

        <!-- 摘要 / summary -->
        <label class="field">
          <span class="field-label">摘要</span>
          <textarea
            v-model="summary"
            class="admin-body"
            style="min-height:84px;font-size:var(--text-sm)"
            placeholder="一句话，写给列表页。"
          ></textarea>
        </label>

        <!-- 置顶 / pinned -->
        <div class="switch-row">
          <button
            class="switch"
            :class="{ 'is-on': pinned }"
            role="switch"
            :aria-checked="pinned"
            @click="pinned = !pinned"
          >
            <span class="switch__knob"></span>
          </button>
          <span class="switch-row__label">置顶这篇</span>
        </div>
      </div>

      <div class="drawer-foot">
        <button class="btn-ghost" @click="drawerSave('draft')">存草稿</button>
        <button class="btn-solid--drawer" @click="drawerSave('published')">
          {{ postStatus === 'published' ? '更新' : '发布' }}
        </button>
      </div>
    </aside>
  </template>
</template>
