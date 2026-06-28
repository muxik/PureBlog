<script setup lang="ts">
import type { Post, PostListResponse } from '@pureblog/api-types'
import { decoratePost, groupByYear } from '~/composables/usePostView'

// ── Settings & date format ──────────────────────────────────────────────────
const settings = useSiteSettings()
const dfmt = settings.value.defaultDateFormat
const { format } = useDateFormat(
  dfmt === 'numeric' || dfmt === 'lunar' ? dfmt : undefined,
)

// ── Page title ─────────────────────────────────────────────────────────────
useHead({ title: settings.value.siteName || 'PureBlog' })

// ── Layout variant ─────────────────────────────────────────────────────────
const { variant, set } = useLayoutVariant()

// ── Fetch posts (SSR + generous pageSize so SSR renders real content) ───────
const { data } = await useFetch<PostListResponse>(`${useApiBase()}/posts`, {
  query: { pageSize: 50 },
})

// ── Decorate & sort: pinned first, then preserve API order (date desc) ──────
const allPosts = computed(() => {
  const raw: Post[] = data.value?.items ?? []
  const sorted = [...raw].sort((a, b) => {
    if (a.pinned && !b.pinned) return -1
    if (!a.pinned && b.pinned) return 1
    return 0
  })
  return sorted.map((p) => decoratePost(p, format.value))
})

// ── Infinite-scroll pagination ─────────────────────────────────────────────
const PAGE_SIZE = 6
const page = ref(1)

// Reset page when the user switches layout variants
watch(variant, () => {
  page.value = 1
})

const pagedPosts = computed(() => allPosts.value.slice(0, page.value * PAGE_SIZE))
const hasMore = computed(() => pagedPosts.value.length < allPosts.value.length)

// Year-grouped view for variant A — derived from the already-paged slice
const groupedPosts = computed(() => groupByYear(pagedPosts.value))

function loadMore() {
  if (hasMore.value) page.value++
}

// ── Sentinel + IntersectionObserver (client-only) ───────────────────────────
const sentinelEl = ref<HTMLElement | null>(null)

onMounted(() => {
  if (!import.meta.client || !('IntersectionObserver' in window)) return

  const io = new IntersectionObserver(
    (entries) => {
      if (entries.some((e) => e.isIntersecting)) loadMore()
    },
    { rootMargin: '240px 0px' },
  )

  // Follow the sentinel element as it appears / disappears from the DOM
  watch(
    sentinelEl,
    (el, prev) => {
      if (prev) io.unobserve(prev)
      if (el) io.observe(el)
    },
    { immediate: true },
  )

  onUnmounted(() => io.disconnect())
})
</script>

<template>
  <!-- Intro lede -->
  <section class="intro wrap">
    <p class="intro__lede">
      {{ settings.description || '这里记录一些思考与文字。慢一点，长一点。' }}
    </p>
  </section>

  <!-- Posts section -->
  <section class="section wrap">
    <!-- Section header: title + segmented control -->
    <div class="section-head">
      <h2 class="section-head__title">文章</h2>
      <div class="segmented">
        <button
          class="seg"
          :class="{ 'is-active': variant === 'A' }"
          @click="set('A')"
        >年表</button>
        <button
          class="seg"
          :class="{ 'is-active': variant === 'B' }"
          @click="set('B')"
        >摘要</button>
        <button
          class="seg"
          :class="{ 'is-active': variant === 'C' }"
          @click="set('C')"
        >双栏</button>
      </div>
    </div>

    <!-- Variant A: 年表 (Timeline grouped by year) -->
    <template v-if="variant === 'A'">
      <template v-for="group in groupedPosts" :key="group.year">
        <div class="year-rule year-rule--timeline">
          <span class="year-rule__label">{{ group.year }}</span>
          <span class="year-rule__line"></span>
        </div>
        <NuxtLink
          v-for="pv in group.posts"
          :key="pv.slug"
          :to="`/post/${pv.slug}`"
          class="row post-row"
        >
          <time class="post-row__time">{{ pv.md }}</time>
          <span class="post-row__title" data-title="">{{ pv.title }}</span>
          <span v-if="pv.pinned" class="badge-pin">置顶</span>
        </NuxtLink>
      </template>
    </template>

    <!-- Variant B: 摘要 (Summary cards) -->
    <div v-else-if="variant === 'B'" class="list-pad">
      <NuxtLink
        v-for="pv in pagedPosts"
        :key="pv.slug"
        :to="`/post/${pv.slug}`"
        class="row post-summary"
      >
        <div class="post-summary__head">
          <h3 class="post-summary__title" data-title="">{{ pv.title }}</h3>
          <span v-if="pv.pinned" class="badge-pin">置顶</span>
        </div>
        <p class="post-summary__excerpt">{{ pv.excerpt }}</p>
        <div class="post-summary__meta">{{ pv.metaLine }}</div>
      </NuxtLink>
    </div>

    <!-- Variant C: 双栏 (Two-column date + title) -->
    <div v-else class="list-pad">
      <NuxtLink
        v-for="pv in pagedPosts"
        :key="pv.slug"
        :to="`/post/${pv.slug}`"
        class="row post-cols"
      >
        <div class="post-cols__date">{{ pv.date }}</div>
        <div>
          <div class="post-cols__head">
            <span class="post-cols__title" data-title="">{{ pv.title }}</span>
            <span v-if="pv.pinned" class="badge-pin">置顶</span>
          </div>
          <div class="post-cols__tags">{{ pv.tagStr }}</div>
        </div>
      </NuxtLink>
    </div>

    <!-- Infinite-scroll sentinel: visible while more posts remain -->
    <div
      v-if="hasMore"
      id="sentinel"
      ref="sentinelEl"
      class="sentinel"
      @click="loadMore"
    >
      <span class="sentinel__dot"></span>
      <span class="sentinel__label">载入中……</span>
    </div>
  </section>
</template>
