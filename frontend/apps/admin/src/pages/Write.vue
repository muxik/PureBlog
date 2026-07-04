<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import type {
  Post,
  SavePostRequest,
  Tag,
  TagListResponse,
  RenderResponse,
} from '@pureblog/api-types'
import { api, uploadFile } from '../api'
import MarkdownEditor from '../components/MarkdownEditor.vue'

const route = useRoute()
const router = useRouter()
const id = route.params.id ? Number(route.params.id) : null
type ViewMode = 'split' | 'source' | 'preview'
const PREVIEW_TYPING_DELAY_MS = 900
const PREVIEW_SETTLE_DELAY_MS = 220

/* ── Core editor state ───────────────────────────────────────────── */
const title = ref('')
const body = ref('')
const wordCount = ref(0)
const postStatus = ref<'draft' | 'published'>('draft')
const viewMode = ref<ViewMode>('split')
const previousMode = ref<ViewMode>('split')
const previewHtml = ref('')
const previewState = ref<'idle' | 'loading' | 'ready' | 'error'>('idle')
const drawerOpen = ref(false)
const showSourcePane = computed(() => viewMode.value !== 'preview')
const showPreviewPane = computed(() => viewMode.value !== 'source')

function isSourceView() {
  return viewMode.value === 'source'
}

/* ── Drawer / publish-settings state ─────────────────────────────── */
const summary = ref('')
const coverUrl = ref('')
const pinned = ref(false)
const tagInput = ref('')        // raw space-separated names shown in the input
const tagChips = ref<string[]>([])  // parsed chips shown below the input

/* ── DOM refs ─────────────────────────────────────────────────────── */
const editorEl = ref<InstanceType<typeof MarkdownEditor> | null>(null)
const fileInputEl = ref<HTMLInputElement | null>(null)
const uploading = ref(false)
let previewTimer: ReturnType<typeof setTimeout> | null = null
let previewSeq = 0
let previewRendering = false
let previewPending = false
let lastPreviewMarkdown = ''

/* ── Display date ─────────────────────────────────────────────────── */
function formatDate(d: Date): string {
  const pad = (n: number) => (n < 10 ? '0' + n : String(n))
  return `${d.getFullYear()} · ${pad(d.getMonth() + 1)} · ${pad(d.getDate())}`
}
const displayDate = ref(formatDate(new Date()))

function updateWordCount(val = body.value) {
  wordCount.value = val.replace(/\s/g, '').length
}

async function ensureEditorVisible() {
  if (viewMode.value === 'preview') {
    viewMode.value = 'split'
    await nextTick()
  }
  editorEl.value?.focus()
}

function syncBodyFromEditor() {
  const value = editorEl.value?.getValue()
  if (value !== undefined && value !== body.value) {
    body.value = value
  }
}

/* ── Toolbar helpers ──────────────────────────────────────────────── */
async function wrap(before: string, after: string) {
  await ensureEditorVisible()
  if (editorEl.value) editorEl.value.wrapSelection(before, after)
  else body.value = `${body.value}${before}${after}`
}

async function linePrefix(prefix: string) {
  await ensureEditorVisible()
  if (editorEl.value) editorEl.value.prefixLine(prefix)
  else body.value = `${prefix}${body.value}`
}

/* ── Image upload ─────────────────────────────────────────────────── */
async function insertAtCursor(text: string) {
  await ensureEditorVisible()
  if (editorEl.value) editorEl.value.insertMarkdown(text)
  else body.value = `${body.value}${text}`
}

// The "图" toolbar button opens the OS file picker.
function insertImage() {
  fileInputEl.value?.click()
}

// Upload one image file and insert its Markdown at the caret.
async function uploadAndInsert(file: File) {
  if (!file.type.startsWith('image/')) return
  uploading.value = true
  try {
    const { url } = await uploadFile(file)
    await insertAtCursor(`![](${url})`)
  } catch {
    await insertAtCursor(`\n<!-- 图片上传失败，请重试 -->\n`)
  } finally {
    uploading.value = false
  }
}

async function onFilePicked(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  if (file) await uploadAndInsert(file)
  input.value = '' // allow re-picking the same file
}

// Paste an image straight from the clipboard.
async function onPaste(e: ClipboardEvent) {
  const item = Array.from(e.clipboardData?.items ?? []).find((i) => i.type.startsWith('image/'))
  if (!item) return
  const file = item.getAsFile()
  if (!file) return
  e.preventDefault()
  await uploadAndInsert(file)
}

// Cover image: upload a file and set coverUrl to the returned URL.
const coverInputEl = ref<HTMLInputElement | null>(null)
async function onCoverPicked(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  if (file && file.type.startsWith('image/')) {
    uploading.value = true
    try {
      const { url } = await uploadFile(file)
      coverUrl.value = url
    } catch {
      /* leave coverUrl unchanged on failure */
    } finally {
      uploading.value = false
    }
  }
  input.value = ''
}

/* ── Live preview ─────────────────────────────────────────────────── */
function clearPreviewTimer() {
  if (!previewTimer) {
    return
  }

  clearTimeout(previewTimer)
  previewTimer = null
}

async function renderPreview(seq: number, markdown = body.value) {
  if (isSourceView()) {
    return
  }

  if (markdown === lastPreviewMarkdown && previewState.value === 'ready') {
    return
  }

  if (!markdown.trim()) {
    if (seq === previewSeq) {
      previewHtml.value = ''
      previewState.value = 'ready'
      lastPreviewMarkdown = markdown
    }
    return
  }

  previewRendering = true
  if (!previewHtml.value) previewState.value = 'loading'

  try {
    const res = await api<RenderResponse>('/admin/render', {
      method: 'POST',
      body: { markdown },
    })
    if (seq !== previewSeq) return
    previewHtml.value = res.html ?? ''
    previewState.value = 'ready'
    lastPreviewMarkdown = markdown
  } catch {
    if (seq !== previewSeq) return
    previewHtml.value = ''
    previewState.value = 'error'
  } finally {
    previewRendering = false
    if (previewPending && !isSourceView()) {
      previewPending = false
      queuePreview(PREVIEW_SETTLE_DELAY_MS)
    }
  }
}

function queuePreview(delay = PREVIEW_TYPING_DELAY_MS) {
  clearPreviewTimer()
  const seq = ++previewSeq
  if (isSourceView()) return

  if (previewRendering) {
    previewPending = true
    return
  }

  if (body.value === lastPreviewMarkdown && previewState.value === 'ready') return

  previewTimer = setTimeout(() => {
    previewTimer = null
    if (previewRendering) {
      previewPending = true
      return
    }
    void renderPreview(seq, body.value)
  }, delay)
}

function refreshPreview() {
  clearPreviewTimer()
  const seq = ++previewSeq
  if (previewRendering) {
    previewPending = true
    return
  }
  void renderPreview(seq, body.value)
}

async function setViewMode(mode: ViewMode) {
  syncBodyFromEditor()
  viewMode.value = mode
  if (mode !== 'source') refreshPreview()
  if (mode !== 'preview') {
    await nextTick()
    editorEl.value?.focus()
  }
}

async function toggleSourceMode() {
  syncBodyFromEditor()
  if (viewMode.value === 'source') {
    viewMode.value = previousMode.value === 'source' ? 'split' : previousMode.value
    refreshPreview()
  } else {
    previousMode.value = viewMode.value
    viewMode.value = 'source'
  }
  await nextTick()
  if (viewMode.value !== 'preview') editorEl.value?.focus()
}

function onWriteKeydown(e: KeyboardEvent) {
  if (e.defaultPrevented) return
  if (e.key !== '/' || (!e.metaKey && !e.ctrlKey)) return
  if ((e.target as Element | null)?.closest?.('.markdown-editor')) return
  e.preventDefault()
  void toggleSourceMode()
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
  window.addEventListener('keydown', onWriteKeydown)

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
  updateWordCount(initialBody)

  await nextTick()
  if (showPreviewPane.value) refreshPreview()
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
  syncBodyFromEditor()

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

watch(body, (value) => {
  updateWordCount(value)
  queuePreview()
})

onBeforeUnmount(() => {
  clearPreviewTimer()
  window.removeEventListener('keydown', onWriteKeydown)
})
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
      <button
        class="tool-btn"
        :title="uploading ? '上传中…' : '插入图片'"
        :disabled="uploading"
        @click="insertImage()"
      >{{ uploading ? '…' : '图' }}</button>
      <input
        ref="fileInputEl"
        type="file"
        accept="image/png,image/jpeg,image/gif,image/webp"
        style="display:none"
        @change="onFilePicked"
      />
      <span class="toolbar-spacer"></span>
      <button
        class="preview-btn"
        :class="{ 'is-on': viewMode === 'split' }"
        @click="setViewMode('split')"
      >双栏</button>
      <button
        class="preview-btn"
        :class="{ 'is-on': viewMode === 'source' }"
        @click="setViewMode('source')"
      >源码</button>
      <button
        class="preview-btn"
        :class="{ 'is-on': viewMode === 'preview' }"
        @click="setViewMode('preview')"
      >预览</button>
    </div>
  </div>

  <!-- ── Editor area ───────────────────────────────────────────────── -->
  <div class="editor" :class="[`editor--${viewMode}`]">
    <div class="editor__workspace">
      <section v-if="showSourcePane" class="editor__pane editor__pane--source">
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
        <MarkdownEditor
          ref="editorEl"
          v-model="body"
          class="editor__body"
          placeholder="开始写。Markdown 可用。慢一点，长一点。可粘贴图片。"
          :source-mode="viewMode === 'source'"
          @mod-slash="toggleSourceMode"
          @paste.capture="onPaste"
        />
      </section>

      <section v-if="showPreviewPane" class="editor__pane editor__pane--preview">
        <div class="preview-block">
          <div class="preview-label">
            {{ previewState === 'loading' ? '预览更新中' : '预览' }}
          </div>
          <h1 class="preview-title">{{ title || '无标题' }}</h1>
          <p v-if="previewState === 'error'" class="preview-empty">渲染失败，请检查网络或重试。</p>
          <p v-else-if="!body.trim()" class="preview-empty">正文预览会出现在这里。</p>
          <div v-else class="prose" style="margin:0" v-html="previewHtml"></div>
        </div>
      </section>
    </div>
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
          <button
            class="linkbtn"
            type="button"
            style="margin-top:8px"
            :disabled="uploading"
            @click="coverInputEl?.click()"
          >{{ uploading ? '上传中…' : '上传图片 ↑' }}</button>
          <input
            ref="coverInputEl"
            type="file"
            accept="image/png,image/jpeg,image/gif,image/webp"
            style="display:none"
            @change="onCoverPicked"
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
