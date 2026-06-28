<script setup lang="ts">
import type { TagListResponse } from '@pureblog/api-types'

const { data } = await useFetch<TagListResponse>(`${useApiBase()}/tags`)
</script>

<template>
  <section>
    <h1 class="page-title">标签</h1>
    <div v-if="data?.items?.length" class="tag-cloud">
      <NuxtLink
        v-for="tag in data.items"
        :key="tag.id"
        :to="`/tags/${tag.slug}`"
        class="tag-chip"
      >{{ tag.name }}</NuxtLink>
    </div>
    <p v-else class="empty">暂无标签。</p>
  </section>
</template>

<style scoped>
.page-title {
  font-family: var(--font-serif, serif);
  font-size: 1.5rem;
  margin: 0 0 1.5rem;
}
.tag-cloud {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}
.tag-chip {
  display: inline-block;
  padding: 0.25rem 0.65rem;
  background: var(--paper-subtle, #f1ece0);
  color: var(--ink-2, #6b655c);
  border-radius: 2px;
  font-size: 0.875rem;
  text-decoration: none;
  border-bottom: none;
  line-height: 1.5;
}
.tag-chip:hover {
  color: var(--accent, #235a73);
}
.empty {
  color: var(--ink-3, #a39c8f);
}
</style>
