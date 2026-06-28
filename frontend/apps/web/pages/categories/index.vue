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
</script>

<template>
  <section>
    <h1 class="page-title">分类</h1>
    <ul v-if="tree.length" class="cat-list">
      <li v-for="cat in tree" :key="cat.id" class="cat-item">
        <NuxtLink :to="`/categories/${cat.slug}`" class="cat-link">{{ cat.name }}</NuxtLink>
        <p v-if="cat.description" class="cat-desc">{{ cat.description }}</p>
        <ul v-if="cat.children.length" class="cat-children">
          <li v-for="child in cat.children" :key="child.id" class="cat-item cat-item--child">
            <NuxtLink :to="`/categories/${child.slug}`" class="cat-link">{{ child.name }}</NuxtLink>
            <p v-if="child.description" class="cat-desc">{{ child.description }}</p>
          </li>
        </ul>
      </li>
    </ul>
    <p v-else class="empty">暂无分类。</p>
  </section>
</template>

<style scoped>
.page-title {
  font-family: var(--font-serif, serif);
  font-size: 1.5rem;
  margin: 0 0 1.5rem;
}
.cat-list {
  list-style: none;
  margin: 0;
  padding: 0;
}
.cat-item {
  padding: 0.75rem 0;
  border-bottom: 1px solid var(--border, #e6e0d3);
}
.cat-item--child {
  border-bottom: none;
  padding: 0.4rem 0;
}
.cat-link {
  text-decoration: none;
  color: var(--accent, #235a73);
  font-family: var(--font-serif, serif);
  font-size: 1.05rem;
}
.cat-link:hover {
  text-decoration: underline;
}
.cat-desc {
  margin: 0.25rem 0 0;
  color: var(--ink-2, #6b655c);
  font-size: 0.9rem;
}
.cat-children {
  list-style: none;
  margin: 0.5rem 0 0 1.25rem;
  padding: 0;
}
.empty {
  color: var(--ink-3, #a39c8f);
}
</style>
