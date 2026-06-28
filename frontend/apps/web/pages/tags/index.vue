<script setup lang="ts">
import type { TagListResponse, PostListResponse } from '@pureblog/api-types'

const apiBase = useApiBase()

// Fetch all tags and all posts (capped at 500) in parallel.
// The posts fetch gives us per-tag counts without a dedicated endpoint.
const [{ data: tagsData }, { data: postsData }] = await Promise.all([
  useFetch<TagListResponse>(`${apiBase}/tags`),
  useFetch<PostListResponse>(`${apiBase}/posts`, { query: { pageSize: 500 } }),
])

// Tally post count for each tag slug.
const tagCounts = computed<Record<string, number>>(() => {
  const counts: Record<string, number> = {}
  for (const post of postsData.value?.items ?? []) {
    for (const tag of post.tags ?? []) {
      if (tag.slug) {
        counts[tag.slug] = (counts[tag.slug] ?? 0) + 1
      }
    }
  }
  return counts
})

// Sort tags by count descending, then alphabetically within same count.
const sortedTags = computed(() => {
  const tags = [...(tagsData.value?.items ?? [])]
  return tags.sort((a, b) => {
    const ca = tagCounts.value[a.slug ?? ''] ?? 0
    const cb = tagCounts.value[b.slug ?? ''] ?? 0
    if (cb !== ca) return cb - ca
    return (a.name ?? '').localeCompare(b.name ?? '', 'zh')
  })
})
</script>

<template>
  <section class="section--pad wrap">
    <h1 class="page-title" style="margin: 0 0 6px">标签</h1>
    <p class="page-sub" style="margin: 0 0 28px">
      按主题翻阅。点一个标签，看看那一类里都写了些什么。
    </p>

    <div v-if="sortedTags.length" class="tags-grid">
      <NuxtLink
        v-for="tag in sortedTags"
        :key="tag.id"
        :to="`/tags/${tag.slug}`"
        class="chip"
      >
        <span class="chip__name">{{ tag.name }}</span>
        <span v-if="tagCounts[tag.slug ?? '']" class="chip__count">
          {{ tagCounts[tag.slug ?? ''] }}
        </span>
      </NuxtLink>
    </div>

    <p v-else class="empty-note">暂无标签。</p>
  </section>
</template>
