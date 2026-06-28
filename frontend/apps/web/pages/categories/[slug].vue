<script setup lang="ts">
import type { CategoryListResponse, PostListResponse } from '@pureblog/api-types'
import { decoratePost } from '~/composables/usePostView'

const route = useRoute()
const slug = route.params.slug as string

const settings = useSiteSettings()

// Resolve the date format preference (falls back to 'numeric' if unset)
const df = settings.value.defaultDateFormat
const { format } = useDateFormat(
  df === 'numeric' || df === 'lunar' ? df : undefined,
)

// Fetch category list to resolve name + description from the slug
const { data: cats } = await useFetch<CategoryListResponse>(`${useApiBase()}/categories`)
const category = computed(() => (cats.value?.items ?? []).find(c => c.slug === slug))

// Fetch all posts for this category in one shot (design: no pagination)
const { data: postsData } = await useFetch<PostListResponse>(`${useApiBase()}/posts`, {
  query: { categorySlug: slug, pageSize: 200 },
})

const postViews = computed(() =>
  (postsData.value?.items ?? []).map(p => decoratePost(p, format.value)),
)

useHead(() => ({
  title: category.value?.name
    ? `${category.value.name} · ${settings.value.siteName ?? 'PureBlog'}`
    : `PureBlog`,
}))
</script>

<template>
  <section class="section--pad wrap">
    <!-- Back link -->
    <NuxtLink to="/categories" class="tag-clear back-link">← 所有分类</NuxtLink>

    <h1 class="page-title" style="margin-top: 20px">{{ category?.name ?? slug }}</h1>
    <p v-if="category?.description" class="page-sub" style="margin: 8px 0 0">
      {{ category.description }}
    </p>

    <div style="margin-top: 28px">
      <template v-if="postViews.length">
        <NuxtLink
          v-for="p in postViews"
          :key="p.slug"
          :to="`/post/${p.slug}`"
          class="row post-row post-row--archive"
        >
          <time class="post-row__time">{{ p.md }}</time>
          <span class="post-row__title" data-title>{{ p.title }}</span>
        </NuxtLink>
      </template>
      <p v-else class="empty-note">该分类暂无文章。</p>
    </div>
  </section>
</template>

<style scoped>
/* .tag-clear is a button style globally; reset anchor defaults */
.back-link {
  text-decoration: none;
  display: inline-block;
}
</style>
