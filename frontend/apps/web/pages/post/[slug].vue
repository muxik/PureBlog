<script setup lang="ts">
import type { Post } from '@pureblog/api-types'

const route = useRoute()
const { data: post } = await useFetch<Post>(`${useApiBase()}/posts/${route.params.slug}`)

useHead(() => ({ title: post.value?.title ?? 'PureBlog' }))
</script>

<template>
  <article v-if="post" class="article">
    <h1 class="article-title">{{ post.title }}</h1>
    <!-- contentHtml is sanitised server-side by bluemonday -->
    <div class="prose" v-html="post.contentHtml" />
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
</style>
