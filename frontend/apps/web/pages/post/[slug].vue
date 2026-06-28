<script setup lang="ts">
import type { Post, Comment, CommentListResponse, Category, CategoryListResponse } from '@pureblog/api-types'

const route = useRoute()
const slug = route.params.slug as string
const apiBase = useApiBase()

// SSR fetches
const { data: post } = await useFetch<Post>(`${apiBase}/posts/${slug}`)
const { data: categoriesData } = await useFetch<CategoryListResponse>(`${apiBase}/categories`)
const { data: commentsData } = await useFetch<CommentListResponse>(`${apiBase}/posts/${slug}/comments`)

useHead(() => ({ title: post.value?.title ?? 'PureBlog' }))

// Build id→category lookup
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

// Comment tree – flatten DFS so each node carries its depth for indentation
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

// Form reactive state
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
    submitError.value = fe?.data?.error ?? fe?.message ?? '提交失败,请稍后再试。'
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <article v-if="post" class="article">
    <h1 class="article-title">{{ post.title }}</h1>
    <!-- contentHtml is sanitised server-side by bluemonday -->
    <div class="prose" v-html="post.contentHtml" />

    <!-- Category + Tags row -->
    <div v-if="postCategory || post.tags?.length" class="meta-row">
      <NuxtLink
        v-if="postCategory && postCategory.slug"
        :to="`/categories/${postCategory.slug}`"
        class="chip chip-category"
      >
        {{ postCategory.name }}
      </NuxtLink>
      <template v-for="tag in post.tags ?? []" :key="tag.id ?? tag.slug ?? ''">
        <NuxtLink
          v-if="tag.slug && tag.name"
          :to="`/tags/${tag.slug}`"
          class="chip chip-tag"
        >
          #&nbsp;{{ tag.name }}
        </NuxtLink>
      </template>
    </div>

    <!-- Comments section -->
    <section class="comments">
      <h2 class="comments-heading">评论</h2>

      <!-- Comment list -->
      <div v-if="commentNodes.length" class="comment-list">
        <div
          v-for="node in commentNodes"
          :key="node.comment.id"
          class="comment"
          :style="node.depth > 0
            ? { paddingLeft: `${node.depth * 1.5}rem`, borderLeft: '2px solid var(--line, #e7e2d6)' }
            : {}"
        >
          <div class="comment-header">
            <span class="comment-author">{{ node.comment.authorName ?? '匿名' }}</span>
            <span class="comment-date">{{ fmtDate(node.comment.createdAt) }}</span>
          </div>
          <p class="comment-content">{{ node.comment.content }}</p>
          <button type="button" class="reply-btn" @click="startReply(node.comment)">回复</button>
        </div>
      </div>
      <p v-else class="comment-empty">还没有评论,来说第一句吧。</p>

      <!-- Comment form -->
      <div ref="formRef" class="comment-form">
        <h3 class="form-heading">
          <template v-if="replyingTo">
            回复 @{{ replyingTo.authorName }}
            <button type="button" class="cancel-btn" @click="cancelReply">取消</button>
          </template>
          <template v-else>留下评论</template>
        </h3>

        <p v-if="submitSuccess" class="form-notice">评论已提交,等待审核后显示。</p>
        <p v-if="submitError" class="form-error">{{ submitError }}</p>

        <form @submit.prevent="submitComment">
          <div class="field">
            <label class="label" for="cf-name">昵称 <span aria-hidden="true">*</span></label>
            <input id="cf-name" v-model="formAuthorName" type="text" class="input" required />
          </div>
          <div class="field">
            <label class="label" for="cf-email">邮箱（可选）</label>
            <input id="cf-email" v-model="formAuthorEmail" type="email" class="input" />
          </div>
          <div class="field">
            <label class="label" for="cf-content">内容 <span aria-hidden="true">*</span></label>
            <textarea id="cf-content" v-model="formContent" class="textarea" rows="4" required />
          </div>
          <button type="submit" class="submit-btn" :disabled="submitting">
            {{ submitting ? '提交中…' : '提交评论' }}
          </button>
        </form>
      </div>
    </section>
  </article>
  <p v-else>未找到这篇文章。</p>
</template>

<style scoped>
.article-title {
  font-family: var(--font-serif, serif);
  font-size: 1.7rem;
  margin: 0 0 1.5rem;
}
.prose :deep(p) {
  line-height: 1.85;
  margin: 1rem 0;
}
.prose :deep(a) {
  color: var(--accent, #235a73);
}
.prose :deep(pre) {
  background: var(--paper-2, #f1ece1);
  padding: 1rem;
  overflow: auto;
  font-family: var(--font-mono, monospace);
}

/* Category + Tags row */
.meta-row {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin-top: 1.5rem;
}
.chip {
  display: inline-block;
  text-decoration: none;
  font-size: 0.8rem;
  padding: 0.2rem 0.6rem;
  border: 1px solid var(--line, #e7e2d6);
  color: var(--ink-2, #5c574e);
}
.chip-category {
  background: var(--paper-subtle, #f7f3ec);
  color: var(--dai, #6b4e3d);
  border-color: var(--dai, #6b4e3d);
}
.chip:hover {
  color: var(--dai, #6b4e3d);
  border-color: var(--dai, #6b4e3d);
}

/* Comments section */
.comments {
  margin-top: 3rem;
  padding-top: 1.5rem;
  border-top: 1px solid var(--line, #e7e2d6);
}
.comments-heading {
  font-family: var(--font-serif, serif);
  font-size: 1.2rem;
  margin: 0 0 1.5rem;
}
.comment-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}
.comment {
  padding: 0.75rem 0;
}
.comment-header {
  display: flex;
  align-items: baseline;
  gap: 0.75rem;
  margin-bottom: 0.4rem;
}
.comment-author {
  font-weight: 600;
  font-size: 0.95rem;
  color: var(--ink, #1a1714);
}
.comment-date {
  font-size: 0.8rem;
  color: var(--ink-3, #8c877e);
}
.comment-content {
  margin: 0 0 0.4rem;
  line-height: var(--leading-body, 1.75);
  color: var(--ink, #1a1714);
}
.reply-btn {
  background: none;
  border: none;
  cursor: pointer;
  font-size: 0.8rem;
  color: var(--ink-3, #8c877e);
  padding: 0;
}
.reply-btn:hover {
  color: var(--dai, #6b4e3d);
}
.comment-empty {
  color: var(--ink-3, #8c877e);
  font-size: 0.95rem;
}

/* Comment form */
.comment-form {
  margin-top: 2.5rem;
  padding-top: 1.5rem;
  border-top: 1px solid var(--line, #e7e2d6);
}
.form-heading {
  font-family: var(--font-serif, serif);
  font-size: 1rem;
  margin: 0 0 1rem;
  display: flex;
  align-items: center;
  gap: 0.75rem;
}
.cancel-btn {
  background: none;
  border: none;
  cursor: pointer;
  font-size: 0.8rem;
  color: var(--ink-3, #8c877e);
  padding: 0;
  font-family: inherit;
}
.cancel-btn:hover {
  color: var(--dai, #6b4e3d);
}
.form-notice {
  color: var(--ink-2, #5c574e);
  background: var(--paper-subtle, #f7f3ec);
  padding: 0.6rem 0.9rem;
  font-size: 0.9rem;
  margin-bottom: 1rem;
  border-left: 3px solid var(--dai, #6b4e3d);
}
.form-error {
  color: #c0392b;
  font-size: 0.9rem;
  margin-bottom: 1rem;
}
.field {
  display: flex;
  flex-direction: column;
  gap: 0.3rem;
  margin-bottom: 1rem;
}
.label {
  font-size: 0.85rem;
  color: var(--ink-2, #5c574e);
}
.input,
.textarea {
  border: 1px solid var(--line, #e7e2d6);
  background: var(--paper, #fdf8f0);
  color: var(--ink, #1a1714);
  padding: 0.4rem 0.6rem;
  font-family: var(--font-body, inherit);
  font-size: 0.95rem;
  width: 100%;
  box-sizing: border-box;
}
.input:focus,
.textarea:focus {
  outline: 2px solid var(--dai, #6b4e3d);
  outline-offset: 1px;
}
.textarea {
  resize: vertical;
}
.submit-btn {
  background: var(--dai, #6b4e3d);
  color: var(--dai-ink, #fff);
  border: none;
  cursor: pointer;
  padding: 0.5rem 1.25rem;
  font-size: 0.9rem;
  font-family: var(--font-body, inherit);
}
.submit-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
.submit-btn:hover:not(:disabled) {
  opacity: 0.9;
}
</style>
