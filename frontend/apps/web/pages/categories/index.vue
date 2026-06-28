<script setup lang="ts">
import type { Category, CategoryListResponse } from '@pureblog/api-types'

const { data } = await useFetch<CategoryListResponse>(`${useApiBase()}/categories`)

type CatNode = Category & { children: CatNode[] }

const tree = computed<CatNode[]>(() => {
  const items = [...(data.value?.items ?? [])].sort(
    (a, b) => (a.sort ?? 0) - (b.sort ?? 0),
  )
  const map = new Map<number, CatNode>()
  for (const c of items) {
    if (c.id != null) map.set(c.id, { ...c, children: [] })
  }
  const roots: CatNode[] = []
  for (const c of items) {
    if (c.id == null) continue
    const node = map.get(c.id)!
    if (c.parentId != null && map.has(c.parentId)) {
      map.get(c.parentId)!.children.push(node)
    } else {
      roots.push(node)
    }
  }
  return roots
})

useHead({ title: '分类 · PureBlog' })
</script>

<template>
  <section class="section--pad wrap">
    <h1 class="page-title">分类</h1>
    <p class="page-sub" style="margin: 8px 0 32px">按类别浏览全部文章。</p>

    <div v-if="tree.length">
      <div v-for="cat in tree" :key="cat.id" class="cat-group">
        <div class="tags-grid">
          <NuxtLink :to="`/categories/${cat.slug}`" class="chip cat-chip">
            <span class="chip__name">{{ cat.name }}</span>
          </NuxtLink>
        </div>
        <p v-if="cat.description" class="page-sub cat-desc">{{ cat.description }}</p>
        <div v-if="cat.children.length" class="tags-grid cat-children">
          <NuxtLink
            v-for="child in cat.children"
            :key="child.id"
            :to="`/categories/${child.slug}`"
            class="chip cat-chip"
          >
            <span class="chip__name">{{ child.name }}</span>
          </NuxtLink>
        </div>
      </div>
    </div>
    <p v-else class="empty-note">暂无分类。</p>
  </section>
</template>

<style scoped>
.cat-group {
  margin-bottom: 28px;
}
.cat-desc {
  margin: 6px 0 0;
  font-size: var(--text-sm);
}
.cat-children {
  margin-top: 10px;
  padding-left: 16px;
}
/* Reset anchor decoration so .chip looks like its button counterpart */
.cat-chip {
  text-decoration: none;
}
</style>
