<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import type { Category, CategoryListResponse, SaveCategoryRequest } from '@pureblog/api-types'
import { api } from '../api'

// ─── list state ───────────────────────────────────────────────────────────────
const items = ref<Category[]>([])
const error = ref<string | null>(null)
const submitting = ref(false)

async function load() {
  try {
    const res = await api<CategoryListResponse>('/categories', { auth: false })
    items.value = res.items ?? []
  } catch (e) {
    error.value = (e as Error).message
  }
}

// ─── tree builder ─────────────────────────────────────────────────────────────
interface TreeNode {
  cat: Category
  level: number
}

function buildTreeNodes(flat: Category[]): TreeNode[] {
  const childMap = new Map<number, Category[]>()
  const roots: Category[] = []

  for (const cat of flat) {
    if (cat.parentId != null) {
      const list = childMap.get(cat.parentId) ?? []
      list.push(cat)
      childMap.set(cat.parentId, list)
    } else {
      roots.push(cat)
    }
  }

  const bySort = (a: Category, b: Category) => (a.sort ?? 0) - (b.sort ?? 0)
  const result: TreeNode[] = []

  function traverse(cats: Category[], level: number): void {
    for (const cat of [...cats].sort(bySort)) {
      result.push({ cat, level })
      if (cat.id != null) {
        traverse(childMap.get(cat.id) ?? [], level + 1)
      }
    }
  }

  traverse([...roots].sort(bySort), 0)
  return result
}

const tree = computed(() => buildTreeNodes(items.value))

// ─── form state ───────────────────────────────────────────────────────────────
const editingId = ref<number | null>(null)
const formName = ref('')
const formSlug = ref('')
const formDescription = ref('')
const formParentId = ref('') // '' = top-level; '123' = parent id
const formSort = ref<number>(0)

function resetForm() {
  editingId.value = null
  formName.value = ''
  formSlug.value = ''
  formDescription.value = ''
  formParentId.value = ''
  formSort.value = 0
  error.value = null
}

function startEdit(cat: Category) {
  editingId.value = cat.id ?? null
  formName.value = cat.name ?? ''
  formSlug.value = cat.slug ?? ''
  formDescription.value = cat.description ?? ''
  formParentId.value = cat.parentId != null ? String(cat.parentId) : ''
  formSort.value = cat.sort ?? 0
  error.value = null
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

// ─── submit ───────────────────────────────────────────────────────────────────
async function submit() {
  if (!formName.value.trim()) {
    error.value = '名称不能为空'
    return
  }
  error.value = null
  submitting.value = true
  try {
    const parsedSort = Number(formSort.value)
    const body: SaveCategoryRequest = {
      name: formName.value.trim(),
      slug: formSlug.value.trim() || undefined,
      description: formDescription.value.trim() || undefined,
      parentId: formParentId.value !== '' ? Number(formParentId.value) : undefined,
      sort: Number.isFinite(parsedSort) ? parsedSort : undefined,
    }
    if (editingId.value !== null) {
      await api(`/admin/categories/${editingId.value}`, { method: 'PUT', body })
    } else {
      await api('/admin/categories', { method: 'POST', body })
    }
    resetForm()
    await load()
  } catch (e) {
    error.value = (e as Error).message
  } finally {
    submitting.value = false
  }
}

// ─── delete ───────────────────────────────────────────────────────────────────
async function remove(id?: number) {
  if (id == null) return
  if (!confirm('确认删除该分类？若分类下有文章或子分类，删除将失败。')) return
  error.value = null
  try {
    await api(`/admin/categories/${id}`, { method: 'DELETE' })
    // if the deleted category was being edited, reset the form
    if (editingId.value === id) resetForm()
    await load()
  } catch (e) {
    error.value = (e as Error).message
  }
}

onMounted(load)
</script>

<template>
  <div class="manage-cat">
    <RouterLink to="/manage" class="back">← 返回文章</RouterLink>
    <h1>分类管理</h1>

    <p v-if="error" class="err">{{ error }}</p>

    <!-- ── form ── -->
    <section class="form-box">
      <h2 class="form-title">{{ editingId !== null ? '编辑分类' : '新建分类' }}</h2>
      <form @submit.prevent="submit">
        <label class="lbl">名称 <span class="req">*</span></label>
        <input v-model="formName" placeholder="分类名称" />

        <label class="lbl">Slug <span class="hint">（留空自动生成）</span></label>
        <input v-model="formSlug" placeholder="url-friendly-slug" />

        <label class="lbl">描述</label>
        <input v-model="formDescription" placeholder="可选描述" />

        <label class="lbl">父分类</label>
        <select v-model="formParentId" class="sel">
          <option value="">（无 / 顶级）</option>
          <option
            v-for="cat in items"
            :key="cat.id"
            :value="String(cat.id)"
            :disabled="cat.id === editingId"
          >
            {{ cat.name }}
          </option>
        </select>

        <label class="lbl">排序</label>
        <input v-model.number="formSort" type="number" placeholder="0" />

        <div class="form-actions">
          <button type="submit" :disabled="submitting">
            {{ editingId !== null ? '保存' : '创建' }}
          </button>
          <button type="button" class="ghost" @click="resetForm">
            {{ editingId !== null ? '取消' : '清空' }}
          </button>
        </div>
      </form>
    </section>

    <!-- ── list ── -->
    <table class="list">
      <tbody>
        <tr v-for="{ cat, level } in tree" :key="cat.id">
          <td class="t">
            <span :style="{ paddingLeft: `${level * 1.5}rem` }">
              <span v-if="level > 0" class="indent-mark">└ </span>{{ cat.name }}
            </span>
          </td>
          <td class="slug">{{ cat.slug }}</td>
          <td class="a">
            <button class="ghost" @click="startEdit(cat)">编辑</button>
            <button class="ghost" @click="remove(cat.id)">删除</button>
          </td>
        </tr>
        <tr v-if="tree.length === 0">
          <td colspan="3" class="empty">暂无分类</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<style scoped>
.back {
  display: inline-block;
  font-size: 0.9rem;
  color: var(--accent, #235a73);
  text-decoration: none;
  margin-bottom: 0.25rem;
}
.back:hover {
  text-decoration: underline;
}
h1 {
  margin-top: 0.4rem;
  margin-bottom: 1.25rem;
}

/* form */
.form-box {
  background: var(--paper-subtle, #f5f2eb);
  border: 1px solid var(--border, #e7e2d6);
  border-radius: 6px;
  padding: 1rem 1.25rem 1.25rem;
  margin-bottom: 1.75rem;
}
.form-title {
  margin: 0 0 0.75rem;
  font-size: 1rem;
  font-weight: 600;
}
.lbl {
  display: block;
  font-size: 0.83rem;
  color: var(--ink-2, #5c5750);
  margin-top: 0.65rem;
}
.req {
  color: #b23b3b;
}
.hint {
  font-weight: normal;
  color: var(--ink-2, #5c5750);
}
.sel {
  display: block;
  width: 100%;
  padding: 0.5rem 0.7rem;
  margin: 0.4rem 0;
  border: 1px solid var(--border, #e7e2d6);
  border-radius: 4px;
  background: var(--paper, #faf8f2);
  color: inherit;
  font: inherit;
}
.form-actions {
  margin-top: 1rem;
}

/* table */
.list {
  width: 100%;
  border-collapse: collapse;
  margin-top: 0.25rem;
}
.list td {
  padding: 0.7rem 0;
  border-bottom: 1px solid var(--border, #e7e2d6);
}
.slug {
  font-size: 0.85rem;
  color: var(--ink-2, #5c5750);
  width: 12rem;
}
.a {
  width: 9rem;
  text-align: right;
}
.indent-mark {
  color: var(--ink-2, #5c5750);
}
.empty {
  color: var(--ink-2, #5c5750);
  text-align: center;
  padding: 2rem 0;
  font-size: 0.9rem;
}
.err {
  color: #b23b3b;
  margin-bottom: 1rem;
}
</style>
