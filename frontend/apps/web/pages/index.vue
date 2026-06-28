<script setup lang="ts">
import type { PostListResponse, CategoryListResponse } from '@pureblog/api-types'

const { data } = await useFetch<PostListResponse>(`${useApiBase()}/posts`, {
  query: { pageSize: 20 },
})
const { data: cats } = await useFetch<CategoryListResponse>(`${useApiBase()}/categories`)

// id → { name, slug } map for resolving post.categoryId to a link
const catMap = computed(() => {
  const m = new Map<number, { name: string; slug: string }>()
  for (const c of cats.value?.items ?? []) {
    if (c.id != null) m.set(c.id, { name: c.name ?? '', slug: c.slug ?? '' })
  }
  return m
})
</script>

<template>
  <section>
    <ul class="post-list">
      <li v-for="p in data?.items ?? []" :key="p.id" class="post-row">
        <NuxtLink :to="`/post/${p.slug}`" class="post-link">
          <span class="post-title">{{ p.title }}</span>
          <span v-if="p.summary" class="post-summary">{{ p.summary }}</span>
        </NuxtLink>
        <div v-if="(p.categoryId != null && catMap.get(p.categoryId)) || p.tags?.length" class="post-meta">
          <NuxtLink
            v-if="p.categoryId != null && catMap.get(p.categoryId)"
            :to="`/categories/${catMap.get(p.categoryId)?.slug}`"
            class="meta-cat"
          >{{ catMap.get(p.categoryId)?.name }}</NuxtLink>
          <NuxtLink
            v-for="t in p.tags ?? []"
            :key="t.id"
            :to="`/tags/${t.slug}`"
            class="meta-tag"
          >{{ t.name }}</NuxtLink>
        </div>
      </li>
    </ul>
  </section>
</template>

<style scoped>
.post-list {
  list-style: none;
  margin: 0;
  padding: 0;
}
.post-row {
  padding: 1.1rem 0;
  border-bottom: 1px solid var(--border, #e7e2d6);
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
  color: var(--ink-2, #5c574e);
  font-size: 0.95rem;
}
.post-row:hover .post-title {
  color: var(--accent, #235a73);
}
.post-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 0.4rem;
  margin-top: 0.5rem;
}
.meta-cat {
  display: inline-block;
  padding: 0.15rem 0.5rem;
  background: color-mix(in srgb, var(--accent, #235a73) 10%, transparent);
  color: var(--accent, #235a73);
  border-radius: 2px;
  font-size: 0.8rem;
  text-decoration: none;
  border-bottom: none;
  line-height: 1.5;
}
.meta-cat:hover {
  background: color-mix(in srgb, var(--accent, #235a73) 18%, transparent);
}
.meta-tag {
  display: inline-block;
  padding: 0.15rem 0.5rem;
  background: var(--paper-subtle, #f1ece0);
  color: var(--ink-2, #6b655c);
  border-radius: 2px;
  font-size: 0.8rem;
  text-decoration: none;
  border-bottom: none;
  line-height: 1.5;
}
.meta-tag:hover {
  color: var(--accent, #235a73);
}
</style>
