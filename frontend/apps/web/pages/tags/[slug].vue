<script setup lang="ts">
import type { TagListResponse, PostListResponse } from '@pureblog/api-types'
import { decoratePost } from '~/composables/usePostView'

const route = useRoute()
const slug = route.params.slug as string
const apiBase = useApiBase()

const settings = useSiteSettings()
const { format } = useDateFormat(
  (settings.value.defaultDateFormat as 'numeric' | 'lunar') ?? 'numeric',
)

// Fetch tags (to resolve the display name) and tag-filtered posts in parallel.
const [{ data: tagsData }, { data: postsData }] = await Promise.all([
  useFetch<TagListResponse>(`${apiBase}/tags`),
  useFetch<PostListResponse>(`${apiBase}/posts`, {
    query: { tagSlug: slug, pageSize: 200 },
  }),
])

// Resolve display name from the tag list.
const tag = computed(() => (tagsData.value?.items ?? []).find(t => t.slug === slug))

// Decorate posts; recomputes automatically when `format` changes.
const posts = computed(() =>
  (postsData.value?.items ?? []).map(p => decoratePost(p, format.value)),
)

useHead(() => ({
  title: tag.value?.name ? `${tag.value.name} · PureBlog` : 'PureBlog',
}))
</script>

<template>
  <section class="section--pad wrap">
    <h1 class="page-title" style="margin: 0 0 12px">
      {{ tag?.name ?? slug }}
    </h1>
    <NuxtLink to="/tags" class="tag-clear">← 全部标签</NuxtLink>

    <template v-if="posts.length">
      <NuxtLink
        v-for="p in posts"
        :key="p.slug"
        :to="`/post/${p.slug}`"
        class="row post-row post-row--archive"
      >
        <time class="post-row__time">{{ p.md }}</time>
        <span class="post-row__title post-row__title--song" data-title="">{{ p.title }}</span>
      </NuxtLink>
    </template>

    <p v-else class="empty-note">该标签暂无文章。</p>
  </section>
</template>
