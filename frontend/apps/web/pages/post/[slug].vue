<script setup lang="ts">
import type { Post, Comment, CommentListResponse, Category, CategoryListResponse, PostListResponse } from '@pureblog/api-types'
import { decoratePost } from '~/composables/usePostView'

const route = useRoute()
const slug = route.params.slug as string
const apiBase = useApiBase()

// Site settings + date format
const settings = useSiteSettings()
const initialFmt: 'numeric' | 'lunar' = settings.value.defaultDateFormat === 'lunar' ? 'lunar' : 'numeric'
const { format } = useDateFormat(initialFmt)

// SSR fetches
const { data: post } = await useFetch<Post>(`${apiBase}/posts/${slug}`)
const { data: categoriesData } = await useFetch<CategoryListResponse>(`${apiBase}/categories`)
const { data: commentsData } = await useFetch<CommentListResponse>(`${apiBase}/posts/${slug}/comments`)
const { data: postsListData } = await useFetch<PostListResponse>(`${apiBase}/posts`, {
  query: { pageSize: 200 },
})

useHead(() => ({ title: post.value?.title ?? 'PureBlog' }))

// Build id → category lookup
const categoryMap = computed((): Map<number, Category> => {
  const m = new Map<number, Category>()
  for (const c of categoriesData.value?.items ?? []) {
    if (c.id != null) m.set(c.id, c)
  }
  return m
})

const postCategory = computed((): Category | null => {
  const id = post.value?.categoryId
  return id != null ? (categoryMap.value.get(id) ?? null) : null
})

// Decorated display data (date, tagStr, readMin, …)
const decorated = computed(() => {
  if (!post.value) return null
  return decoratePost(post.value, format.value)
})

// Filtered tags that have both slug and name
const displayTags = computed(() =>
  (decorated.value?.tags ?? []).filter((t) => t.slug && t.name),
)

// Pager — match the design's prev/next logic on the sorted list
const pagerPrev = computed((): Post | null => {
  const items = postsListData.value?.items ?? []
  const idx = items.findIndex((p) => p.slug === slug)
  return idx > 0 ? (items[idx - 1] ?? null) : null
})
const pagerNext = computed((): Post | null => {
  const items = postsListData.value?.items ?? []
  const idx = items.findIndex((p) => p.slug === slug)
  return idx >= 0 && idx < items.length - 1 ? (items[idx + 1] ?? null) : null
})

// Seal character — first glyph of site name, or 木
const sealChar = computed(() => settings.value.siteName?.[0] ?? '木')

// Comment tree — flatten DFS so each node carries its depth for indentation
interface CommentNode {
  comment: Comment
  depth: number
}

function flattenComments(comments: Comment[]): CommentNode[] {
  const byParent = new Map<number, Comment[]>()
  byParent.set(0, [])
  for (const c of comments) {
    const pid = c.parentId ?? 0
    if (!byParent.has(pid)) byParent.set(pid, [])
    byParent.get(pid)!.push(c)
  }
  const result: CommentNode[] = []
  const walk = (pid: number, depth: number): void => {
    for (const c of byParent.get(pid) ?? []) {
      result.push({ comment: c, depth })
      if (c.id != null) walk(c.id, depth + 1)
    }
  }
  walk(0, 0)
  return result
}

const commentNodes = computed(() => flattenComments(commentsData.value?.items ?? []))

function fmtDate(iso: string | undefined): string {
  if (!iso) return ''
  return new Date(iso).toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' })
}

// Comment form reactive state
const formAuthorName = ref('')
const formAuthorEmail = ref('')
const formContent = ref('')
const replyingTo = ref<{ id: number; authorName: string } | null>(null)
const submitting = ref(false)
const submitSuccess = ref(false)
const submitError = ref('')
const formRef = ref<HTMLElement | null>(null)

function startReply(comment: Comment): void {
  if (comment.id == null) return
  replyingTo.value = { id: comment.id, authorName: comment.authorName ?? '匿名' }
  nextTick(() => {
    formRef.value?.scrollIntoView({ behavior: 'smooth', block: 'start' })
  })
}

function cancelReply(): void {
  replyingTo.value = null
}

async function submitComment(): Promise<void> {
  if (!formAuthorName.value.trim() || !formContent.value.trim()) return
  submitting.value = true
  submitSuccess.value = false
  submitError.value = ''
  try {
    const body: { authorName: string; content: string; parentId?: number; authorEmail?: string } = {
      authorName: formAuthorName.value.trim(),
      content: formContent.value.trim(),
    }
    if (formAuthorEmail.value.trim()) body.authorEmail = formAuthorEmail.value.trim()
    if (replyingTo.value != null) body.parentId = replyingTo.value.id
    await $fetch(`${apiBase}/posts/${slug}/comments`, { method: 'POST', body })
    formAuthorName.value = ''
    formAuthorEmail.value = ''
    formContent.value = ''
    replyingTo.value = null
    submitSuccess.value = true
  } catch (err: unknown) {
    const fe = err as { data?: { error?: string }; message?: string }
    submitError.value = fe?.data?.error ?? fe?.message ?? '提交失败，请稍后再试。'
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <article v-if="post && decorated" class="article wrap">
    <!-- Back link -->
    <NuxtLink class="backlink" to="/">← 返回</NuxtLink>

    <!-- Article header -->
    <header class="article__header">
      <h1 class="article__title">{{ post.title }}</h1>
      <div class="article__meta">
        <time class="article__date">{{ decorated.date }}</time>

        <!-- Tag links -->
        <span v-if="displayTags.length" class="article__tags">
          <template v-for="(tag, i) in displayTags" :key="tag.id ?? tag.slug">
            <NuxtLink :to="`/tags/${tag.slug}`" class="tag-link">{{ tag.name }}</NuxtLink>
            <span v-if="i < displayTags.length - 1" aria-hidden="true"> · </span>
          </template>
        </span>

        <!-- Estimated read time -->
        <span v-if="decorated.readMin" class="article__read">约 {{ decorated.readMin }} 分钟</span>

        <!-- Category chip -->
        <NuxtLink
          v-if="postCategory?.slug && postCategory?.name"
          :to="`/categories/${postCategory.slug}`"
          class="chip cat-chip"
        >{{ postCategory.name }}</NuxtLink>
      </div>
    </header>

    <!-- Post body — contentHtml is sanitised server-side by bluemonday -->
    <div class="prose" v-html="post.contentHtml" />

    <!-- Pager -->
    <div class="pager-wrap">
      <nav class="pager">
        <NuxtLink
          v-if="pagerPrev?.slug"
          :to="`/post/${pagerPrev.slug}`"
          class="row pager__cell"
        >
          <span class="pager__dir">← 上一篇</span>
          <span class="pager__title" data-title="">{{ pagerPrev.title }}</span>
        </NuxtLink>
        <span v-else class="pager__cell--empty" />

        <NuxtLink
          v-if="pagerNext?.slug"
          :to="`/post/${pagerNext.slug}`"
          class="row pager__cell pager__cell--next"
        >
          <span class="pager__dir">下一篇 →</span>
          <span class="pager__title" data-title="">{{ pagerNext.title }}</span>
        </NuxtLink>
        <span v-else class="pager__cell--empty" />
      </nav>
    </div>

    <!-- Seal (落款印) -->
    <div class="seal-row">
      <span class="seal" style="width:30px;height:30px;font-size:18px">{{ sealChar }}</span>
    </div>

    <!-- ── Comments ──────────────────────────────────────────────── -->
    <section class="cmt">
      <h2 class="cmt__heading">评论</h2>

      <!-- Approved comment list -->
      <div v-if="commentNodes.length" class="cmt__list">
        <div
          v-for="node in commentNodes"
          :key="node.comment.id"
          class="cmt__item"
          :class="{ 'cmt__item--reply': node.depth > 0 }"
          :style="node.depth > 0 ? { paddingLeft: `${node.depth * 1.25}rem` } : {}"
        >
          <div class="cmt__row">
            <span class="cmt__author">{{ node.comment.authorName ?? '匿名' }}</span>
            <span class="cmt__date">{{ fmtDate(node.comment.createdAt) }}</span>
          </div>
          <p class="cmt__body">{{ node.comment.content }}</p>
          <button type="button" class="cmt__reply-btn" @click="startReply(node.comment)">回复</button>
        </div>
      </div>
      <p v-else class="cmt__empty">还没有评论，来说第一句吧。</p>

      <!-- Comment form -->
      <div ref="formRef" class="cmt__form">
        <h3 class="cmt__form-title">
          <template v-if="replyingTo">
            回复 @{{ replyingTo.authorName }}
            <button type="button" class="cmt__cancel" @click="cancelReply">取消</button>
          </template>
          <template v-else>留下评论</template>
        </h3>

        <p v-if="submitSuccess" class="cmt__notice">评论已提交，等待审核后显示。</p>
        <p v-if="submitError" class="cmt__error">{{ submitError }}</p>

        <form @submit.prevent="submitComment">
          <div class="cmt__field">
            <label class="cmt__label" for="cf-name">昵称 <span aria-hidden="true">*</span></label>
            <input id="cf-name" v-model="formAuthorName" type="text" class="cmt__input" required />
          </div>
          <div class="cmt__field">
            <label class="cmt__label" for="cf-email">邮箱（可选）</label>
            <input id="cf-email" v-model="formAuthorEmail" type="email" class="cmt__input" />
          </div>
          <div class="cmt__field">
            <label class="cmt__label" for="cf-content">内容 <span aria-hidden="true">*</span></label>
            <textarea id="cf-content" v-model="formContent" class="cmt__textarea" rows="4" required />
          </div>
          <button type="submit" class="cmt__submit" :disabled="submitting">
            {{ submitting ? '提交中…' : '提交评论' }}
          </button>
        </form>
      </div>
    </section>
  </article>

  <p v-else class="wrap" style="padding:40px 20px;color:var(--text-muted)">未找到这篇文章。</p>
</template>

<style scoped>
/* ── Prose content ────────────────────────────────────────────
   Global app.css has no .prose deep-rules; add them here.
   ─────────────────────────────────────────────────────────── */
.prose :deep(p) {
  line-height: var(--leading-body);
  margin: 1rem 0;
  color: var(--text-body);
}
.prose :deep(a) {
  color: var(--accent);
  text-decoration: underline;
  text-underline-offset: 2px;
}
.prose :deep(a:hover) { opacity: 0.8; }
.prose :deep(h2) {
  font-family: var(--font-display);
  font-size: var(--text-h2);
  font-weight: 700;
  color: var(--text-body);
  margin: 2rem 0 0.75rem;
  letter-spacing: 0;
}
.prose :deep(h3) {
  font-family: var(--font-display);
  font-size: var(--text-h3);
  font-weight: 600;
  color: var(--text-body);
  margin: 1.5rem 0 0.5rem;
  letter-spacing: 0;
}
.prose :deep(blockquote) {
  margin: 1.5rem 0;
  padding: 0 0 0 1.1rem;
  border-left: 3px solid var(--accent);
  color: var(--text-muted);
  font-style: italic;
}
.prose :deep(blockquote p) { margin: 0; }
.prose :deep(pre) {
  background: var(--bg-subtle);
  border: 1px solid var(--line);
  padding: 1rem 1.2rem;
  overflow-x: auto;
  font-family: var(--font-mono);
  font-size: var(--text-sm);
  line-height: 1.6;
  border-radius: var(--radius-sm);
  margin: 1.5rem 0;
}
.prose :deep(code) {
  font-family: var(--font-mono);
  font-size: 0.9em;
  background: var(--bg-subtle);
  padding: 0.1em 0.3em;
  border-radius: 2px;
}
.prose :deep(pre code) {
  background: none;
  padding: 0;
  font-size: inherit;
}
.prose :deep(ul),
.prose :deep(ol) {
  padding-left: 1.4em;
  margin: 1rem 0;
}
.prose :deep(li) {
  margin: 0.3em 0;
  line-height: var(--leading-body);
}
.prose :deep(img) {
  max-width: 100%;
  height: auto;
}

/* ── Tag links in article meta ───────────────────────────────── */
.tag-link {
  color: var(--text-muted);
  text-decoration: none;
  font-family: var(--font-body);
  font-size: var(--text-sm);
  transition: color var(--dur-fast) var(--ease);
}
.tag-link:hover { color: var(--accent); }

/* ── Category chip — compact variant for inline meta ─────────── */
.cat-chip {
  text-decoration: none;
  font-size: var(--text-xs);
  padding: 2px 8px;
  line-height: 1.6;
}

/* ── Comments section ────────────────────────────────────────── */
.cmt {
  margin-top: 64px;
  padding-top: 24px;
  border-top: 1px solid var(--line);
}
.cmt__heading {
  font-family: var(--font-display);
  font-size: var(--text-h3);
  font-weight: 700;
  color: var(--text-body);
  margin: 0 0 1.5rem;
  letter-spacing: 0;
}

/* Comment list */
.cmt__list { display: flex; flex-direction: column; }
.cmt__item {
  padding: 1rem 0;
  border-bottom: 1px solid var(--line);
}
.cmt__item--reply { border-left: 2px solid var(--line); }
.cmt__row {
  display: flex;
  align-items: baseline;
  gap: 0.75rem;
  margin-bottom: 0.4rem;
}
.cmt__author {
  font-family: var(--font-body);
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--text-body);
}
.cmt__date {
  font-family: var(--font-mono);
  font-size: var(--text-xs);
  color: var(--text-faint);
}
.cmt__body {
  margin: 0 0 0.4rem;
  font-family: var(--font-body);
  font-size: var(--text-base);
  line-height: var(--leading-body);
  color: var(--text-body);
}
.cmt__reply-btn {
  background: none;
  border: none;
  cursor: pointer;
  font-family: var(--font-mono);
  font-size: var(--text-xs);
  color: var(--text-faint);
  padding: 0;
  letter-spacing: 0.04em;
  transition: color var(--dur-fast) var(--ease);
}
.cmt__reply-btn:hover { color: var(--accent); }
.cmt__empty {
  font-family: var(--font-body);
  font-size: var(--text-base);
  color: var(--text-faint);
  margin: 0 0 1rem;
}

/* Comment form */
.cmt__form {
  margin-top: 2.5rem;
  padding-top: 1.5rem;
  border-top: 1px solid var(--line);
}
.cmt__form-title {
  font-family: var(--font-display);
  font-size: 1rem;
  font-weight: 600;
  color: var(--text-body);
  margin: 0 0 1rem;
  display: flex;
  align-items: center;
  gap: 0.75rem;
  letter-spacing: 0;
}
.cmt__cancel {
  background: none;
  border: none;
  cursor: pointer;
  font-family: var(--font-mono);
  font-size: var(--text-xs);
  color: var(--text-faint);
  padding: 0;
  letter-spacing: 0.04em;
  transition: color var(--dur-fast) var(--ease);
}
.cmt__cancel:hover { color: var(--accent); }
.cmt__notice {
  font-family: var(--font-body);
  font-size: var(--text-sm);
  color: var(--text-muted);
  background: var(--bg-subtle);
  padding: 0.6rem 0.9rem;
  border-left: 3px solid var(--accent);
  margin-bottom: 1rem;
}
.cmt__error {
  font-family: var(--font-body);
  font-size: var(--text-sm);
  color: #c0392b;
  margin-bottom: 1rem;
}
.cmt__field {
  display: flex;
  flex-direction: column;
  gap: 0.3rem;
  margin-bottom: 1rem;
}
.cmt__label {
  font-family: var(--font-mono);
  font-size: var(--text-xs);
  color: var(--text-muted);
  letter-spacing: 0.04em;
}
.cmt__input,
.cmt__textarea {
  font-family: var(--font-body);
  font-size: var(--text-base);
  color: var(--text-body);
  background: var(--bg-subtle);
  border: 1px solid var(--line);
  border-radius: var(--radius-sm);
  padding: 0.5rem 0.75rem;
  width: 100%;
  box-sizing: border-box;
  outline: none;
  transition: border-color var(--dur-fast) var(--ease);
}
.cmt__input:focus,
.cmt__textarea:focus { border-color: var(--accent); }
.cmt__textarea {
  resize: vertical;
  line-height: var(--leading-body);
}
.cmt__submit {
  background: var(--accent);
  color: var(--accent-on);
  border: none;
  cursor: pointer;
  font-family: var(--font-body);
  font-size: var(--text-sm);
  padding: 0.5rem 1.25rem;
  border-radius: var(--radius-sm);
  transition: opacity var(--dur-fast) var(--ease);
}
.cmt__submit:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
.cmt__submit:hover:not(:disabled) { opacity: 0.85; }
</style>
