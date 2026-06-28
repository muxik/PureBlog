<script setup lang="ts">
import type { PostListResponse } from '@pureblog/api-types'
import { decoratePost } from '~/composables/usePostView'

const route = useRoute()
const router = useRouter()
const apiBase = useApiBase()

// Normalise initial ?q= from URL (SSR-safe: route is available on the server).
const rawQ = route.query.q
const initialQ: string = Array.isArray(rawQ) ? (rawQ[0] ?? '') : ((rawQ as string) ?? '')

// `searchQ` drives the fetch; updated after the debounce fires.
// `inputText` is bound to the <input> directly so the field stays responsive.
const searchQ = ref<string>(initialQ)
const inputText = ref<string>(initialQ)

const settings = useSiteSettings()
const { format } = useDateFormat(
  (settings.value.defaultDateFormat as 'numeric' | 'lunar') ?? 'numeric',
)

// Nuxt watches reactive refs in the query object and re-fetches when they change.
// On SSR the initial value of `searchQ` is used, giving server-rendered results
// when the page is accessed with ?q= already set.
const { data: postsData } = await useFetch<PostListResponse>(`${apiBase}/posts`, {
  query: { q: searchQ, pageSize: 50 },
})

// Only surface results when there is an actual query term.
const posts = computed(() => {
  if (!searchQ.value.trim()) return []
  return (postsData.value?.items ?? []).map(p => decoratePost(p, format.value))
})

// Use the server-side total (accurate even when more than 50 results exist).
const heading = computed(() => {
  const q = searchQ.value.trim()
  if (!q) return ''
  const n = postsData.value?.total ?? posts.value.length
  return `「${q}」· 共 ${n} 篇`
})

// Sync URL ?q= back into state when navigating with browser back/forward.
watch(
  () => route.query.q,
  (newQ) => {
    const q = Array.isArray(newQ) ? (newQ[0] ?? '') : ((newQ as string) ?? '')
    if (q !== searchQ.value) {
      searchQ.value = q
      inputText.value = q
    }
  },
)

// Debounce: update the authoritative query and URL ~250 ms after the user stops typing.
let debounceTimer: ReturnType<typeof setTimeout> | null = null

function onInput(e: Event) {
  const value = (e.target as HTMLInputElement).value
  inputText.value = value
  if (debounceTimer !== null) clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => {
    searchQ.value = value
    router.replace({ query: value.trim() ? { q: value } : {} })
  }, 250)
}
</script>

<template>
  <section class="section--pad wrap">
    <h1 class="page-title" style="margin: 0 0 20px">搜索</h1>

    <div style="position: relative; margin-bottom: 8px">
      <input
        class="search-input"
        type="text"
        :value="inputText"
        placeholder="搜索标题、摘要或标签……"
        @input="onInput"
      />
    </div>

    <!-- Meta line: result count — only visible while a query is active -->
    <div v-if="searchQ.trim()" class="search-meta">
      <span class="search-heading">{{ heading }}</span>
    </div>

    <!-- Results -->
    <template v-if="searchQ.trim()">
      <template v-if="posts.length">
        <NuxtLink
          v-for="p in posts"
          :key="p.slug"
          :to="`/post/${p.slug}`"
          class="row post-summary"
          style="margin-top: 12px"
        >
          <h3 class="post-summary__title" data-title="" style="margin-bottom: 6px">
            {{ p.title }}
          </h3>
          <p class="post-summary__excerpt" style="white-space: normal">{{ p.excerpt }}</p>
          <div class="post-summary__meta">{{ p.metaLine }}</div>
        </NuxtLink>
      </template>
      <p v-else class="empty-note">没有找到相关的文章。换个词试试？</p>
    </template>

    <!-- No query yet -->
    <p v-else class="empty-note">输入关键词开始搜索。</p>
  </section>
</template>
