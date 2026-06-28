<script setup lang="ts">
import type { TagListResponse } from '@pureblog/api-types'

const route = useRoute()
const slug = route.params.slug as string

const { data: tags } = await useFetch<TagListResponse>(`${useApiBase()}/tags`)
const tag = computed(() => (tags.value?.items ?? []).find(t => t.slug === slug))

const { data: posts, page, totalPages } = await usePostList({ tagSlug: slug })

useHead(() => ({
  title: tag.value?.name ? `${tag.value.name} · PureBlog` : 'PureBlog',
}))
</script>

<template>
  <section>
    <header class="archive-header">
      <h1 class="archive-title"># {{ tag?.name ?? slug }}</h1>
    </header>

    <ul v-if="posts?.items?.length" class="post-list">
      <li v-for="p in posts.items" :key="p.id" class="post-row">
        <NuxtLink :to="`/post/${p.slug}`" class="post-link">
          <span class="post-title">{{ p.title }}</span>
          <span v-if="p.summary" class="post-summary">{{ p.summary }}</span>
        </NuxtLink>
      </li>
    </ul>
    <p v-else class="empty">该标签暂无文章。</p>

    <nav v-if="totalPages > 1" class="pagination">
      <NuxtLink
        v-if="page > 1"
        :to="{ query: { page: page - 1 } }"
        class="page-btn"
      >上一页</NuxtLink>
      <span class="page-info">第 {{ page }} / {{ totalPages }} 页</span>
      <NuxtLink
        v-if="page < totalPages"
        :to="{ query: { page: page + 1 } }"
        class="page-btn"
      >下一页</NuxtLink>
    </nav>
  </section>
</template>

<style scoped>
.archive-header {
  margin-bottom: 1.75rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid var(--border, #e6e0d3);
}
.archive-title {
  font-family: var(--font-serif, serif);
  font-size: 1.5rem;
  margin: 0;
}
.post-list {
  list-style: none;
  margin: 0;
  padding: 0;
}
.post-row {
  padding: 1.1rem 0;
  border-bottom: 1px solid var(--border, #e6e0d3);
}
.post-link {
  text-decoration: none;
  color: inherit;
  display: block;
}
.post-title {
  display: block;
  font-family: var(--font-serif, serif);
  font-size: 1.15rem;
}
.post-summary {
  display: block;
  margin-top: 0.35rem;
  color: var(--ink-2, #6b655c);
  font-size: 0.95rem;
}
.post-row:hover .post-title {
  color: var(--accent, #235a73);
}
.empty {
  padding: 1rem 0;
  color: var(--ink-3, #a39c8f);
}
.pagination {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-top: 2rem;
  font-size: 0.9rem;
}
.page-btn {
  color: var(--accent, #235a73);
  text-decoration: none;
}
.page-btn:hover {
  text-decoration: underline;
}
.page-info {
  color: var(--ink-2, #6b655c);
}
</style>
