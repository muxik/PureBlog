<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Vditor from 'vditor'
import 'vditor/dist/index.css'
import type { Post, SavePostRequest } from '@pureblog/api-types'
import { api } from '../api'

const route = useRoute()
const router = useRouter()
const id = route.params.id ? Number(route.params.id) : null

const title = ref('')
const slug = ref('')
const summary = ref('')
const editorEl = ref<HTMLDivElement | null>(null)
let editor: Vditor | null = null

onMounted(async () => {
  let initial = ''
  if (id) {
    const post = await api<Post>(`/admin/posts/${id}`)
    title.value = post.title ?? ''
    slug.value = post.slug ?? ''
    summary.value = post.summary ?? ''
    initial = post.contentMd ?? ''
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
</style>
