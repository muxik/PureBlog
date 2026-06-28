import type { PostListResponse } from '@pureblog/api-types'

export interface PostListFilter {
  categorySlug?: string
  tagSlug?: string
}

/**
 * Shared composable for paginated post listing filtered by category or tag.
 * Page is synced to route.query.page for SSR + browser back/forward support.
 */
export async function usePostList(filter: PostListFilter) {
  // Collect all composables that need instance context before the first await.
  const route = useRoute()
  const apiBase = useApiBase()
  const page = computed(() => Number(route.query.page ?? 1))

  const { data } = await useFetch<PostListResponse>(`${apiBase}/posts`, {
    // Passing `page` (a ComputedRef) directly causes useFetch to watch it
    // and re-fetch whenever the route query changes (client-side navigation).
    query: { ...filter, page, pageSize: 10 },
  })

  const totalPages = computed(() => {
    const total = data.value?.total ?? 0
    const pageSize = data.value?.pageSize ?? 10
    return Math.max(1, Math.ceil(total / pageSize))
  })

  return { data, page, totalPages }
}
