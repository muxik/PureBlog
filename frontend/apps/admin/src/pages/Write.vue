<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Vditor from 'vditor'
import 'vditor/dist/index.css'
import type {
  Post,
  SavePostRequest,
  Category,
  CategoryListResponse,
  Tag,
  TagListResponse,
} from '@pureblog/api-types'
import { api } from '../api'

const route = useRoute()
const router = useRouter()
const id = route.params.id ? Number(route.params.id) : null

const title = ref('')
const slug = ref('')
const summary = ref('')
const categoryId = ref<number | null>(null)
const tagIds = ref<number[]>([])
const categories = ref<Category[]>([])
const allTags = ref<Tag[]>([])
const editorEl = ref<HTMLDivElement | null>(null)
let editor: Vditor | null = null

function isTagSelected(tag: Tag): boolean {
  return tag.id !== undefined && tagIds.value.includes(tag.id)
}

function handleTagToggle(tag: Tag): void {
  if (tag.id === undefined) return
  const idx = tagIds.value.indexOf(tag.id)
  if (idx === -1) tagIds.value.push(tag.id)
  else tagIds.value.splice(idx, 1)
}

onMounted(async () => {
  let initial = ''

  // Fetch categories and tags (public) in parallel with optional post load
  const [catRes, tagRes, postData] = await Promise.all([
    api<CategoryListResponse>('/categories', { auth: false }),
    api<TagListResponse>('/tags', { auth: false }),
    id ? api<Post>(`/admin/posts/${id}`) : Promise.resolve(null as Post | null),
  ])

  categories.value = catRes.items ?? []
  allTags.value = tagRes.items ?? []

  if (postData) {
    title.value = postData.title ?? ''
    slug.value = postData.slug ?? ''
    summary.value = postData.summary ?? ''
    initial = postData.contentMd ?? ''
    categoryId.value = postData.categoryId ?? null
    tagIds.value = (postData.tags ?? [])
      .map((t) => t.id)
      .filter((n): n is number => n !== undefined)
  }

  editor = new Vditor(editorEl.value as HTMLElement, {
    height: 480,
    placeholder: '在此写作……',
    cache: { enable: false },
    after: () => editor?.setValue(initial),
  })
})

onBeforeUnmount(() => editor?.destroy())

async function save(status: 'draft' | 'published') {
  const body: SavePostRequest = {
    title: title.value,
    slug: slug.value,
    summary: summary.value,
    contentMd: editor?.getValue() ?? '',
    status,
    categoryId: categoryId.value ?? undefined,
    tagIds: tagIds.value,
  }
  if (id) await api(`/admin/posts/${id}`, { method: 'PUT', body })
  else await api('/admin/posts', { method: 'POST', body })
  router.push('/manage')
}
</script>

<template>
  <div class="write">
    <input v-model="title" class="title" placeholder="标题" />
    <input v-model="slug" placeholder="slug(留空自动生成)" />
    <input v-model="summary" placeholder="摘要" />

    <div class="field">
      <label class="field-label">分类</label>
      <select v-model="categoryId" class="field-select">
        <option :value="null">无分类</option>
        <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
      </select>
    </div>

    <div class="field">
      <label class="field-label">标签</label>
      <div class="tags">
        <button
          v-for="tag in allTags"
          :key="tag.id"
          type="button"
          class="chip"
          :class="{ selected: isTagSelected(tag) }"
          @click="handleTagToggle(tag)"
        >{{ tag.name }}</button>
        <span v-if="allTags.length === 0" class="no-tags">暂无标签</span>
      </div>
    </div>

    <div ref="editorEl" class="editor"></div>
    <div class="actions">
      <button class="ghost" @click="save('draft')">存草稿</button>
      <button @click="save('published')">发布</button>
    </div>
  </div>
</template>

<style scoped>
.title {
  font-size: 1.2rem;
}
.editor {
  margin: 0.8rem 0;
}
.actions {
  margin-top: 1rem;
}
.field {
  margin: 0.4rem 0;
}
.field-label {
  display: block;
  font-size: 0.85rem;
  color: var(--ink-2, #7a7060);
  margin-bottom: 0.25rem;
}
.field-select {
  display: block;
  width: 100%;
  padding: 0.5rem 0.7rem;
  border: 1px solid var(--border, #e7e2d6);
  border-radius: 4px;
  background: var(--paper, #faf8f2);
  color: inherit;
  font: inherit;
}
.tags {
  display: flex;
  flex-wrap: wrap;
  gap: 0.4rem;
  padding: 0.4rem 0;
}
.chip {
  padding: 0.25rem 0.75rem;
  margin: 0;
  border: 1px solid var(--border, #e7e2d6);
  border-radius: 99px;
  background: transparent;
  color: var(--ink-2, #7a7060);
  cursor: pointer;
  font: inherit;
  font-size: 0.85rem;
  line-height: 1.4;
  transition: background 0.15s, color 0.15s, border-color 0.15s;
}
.chip.selected {
  background: var(--accent, #235a73);
  border-color: var(--accent, #235a73);
  color: #fff;
}
.chip:hover:not(.selected) {
  border-color: var(--accent, #235a73);
  color: var(--accent, #235a73);
}
.no-tags {
  font-size: 0.85rem;
  color: var(--ink-2, #7a7060);
}
</style>
