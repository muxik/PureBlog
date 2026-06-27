<script setup lang="ts">
import type { PostListResponse } from '@pureblog/api-types'

const { data } = await useFetch<PostListResponse>(`${useApiBase()}/posts`, {
  query: { pageSize: 20 },
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
</style>
