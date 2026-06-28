<script setup lang="ts">
import type { PostListResponse } from '@pureblog/api-types'
import { decoratePost, groupByYear } from '~/composables/usePostView'

const { data } = await useFetch<PostListResponse>(`${useApiBase()}/posts`, {
  query: { pageSize: 500 },
})

const settings = useSiteSettings()
const dfmt = settings.value.defaultDateFormat === 'lunar' ? 'lunar' : 'numeric'
const { format } = useDateFormat(dfmt)

const yearGroups = computed(() => {
  const posts = data.value?.items ?? []
  const views = posts.map((p) => decoratePost(p, format.value))
  return groupByYear(views)
})

const totalCount = computed(() => data.value?.items?.length ?? 0)

useHead({
  title: computed(() => '归档 · ' + (settings.value.siteName || 'PureBlog')),
})
</script>

<template>
  <section class="section--pad wrap">
    <div style="margin-bottom: 8px">
      <h1 class="page-title">归档</h1>
      <p class="page-sub" style="margin: 8px 0 0">
        共 {{ totalCount }} 篇，跨 {{ yearGroups.length }} 个年份。
      </p>
    </div>

    <template v-for="g in yearGroups" :key="g.year">
      <div class="year-rule year-rule--archive">
        <h2 class="year-rule__year">{{ g.year }}</h2>
        <span class="year-rule__count">{{ g.posts.length }} 篇</span>
        <span class="year-rule__line" />
      </div>

      <NuxtLink
        v-for="pv in g.posts"
        :key="pv.slug"
        :to="`/post/${pv.slug}`"
        class="row post-row post-row--archive"
      >
        <time class="post-row__time">{{ pv.md }}</time>
        <span class="post-row__title post-row__title--song" data-title="">{{ pv.title }}</span>
        <span v-if="pv.pinned" class="badge-pin">置顶</span>
      </NuxtLink>
    </template>
  </section>
</template>
