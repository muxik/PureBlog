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
  <section class="admin-section">
    <div class="manage-head">
      <h1 class="section-title">分类管理</h1>
    </div>

    <!-- ── Create / Edit Form ── -->
    <div class="cat-form">
      <h2 class="panel-h">{{ editingId !== null ? '编辑分类' : '新建分类' }}</h2>
      <p v-if="error" class="cat-err">{{ error }}</p>
      <form @submit.prevent="submit">
        <div class="field-stack field-stack--tight">
          <label class="field">
            <span class="field-label">名称<span class="cat-req"> *</span></span>
            <input class="admin-input" v-model="formName" placeholder="分类名称" />
          </label>
          <label class="field">
            <span class="field-label">Slug（留空自动生成）</span>
            <input class="admin-input" v-model="formSlug" placeholder="url-friendly-slug" />
          </label>
          <label class="field">
            <span class="field-label">描述</span>
            <input class="admin-input" v-model="formDescription" placeholder="可选描述" />
          </label>
          <label class="field">
            <span class="field-label">父分类</span>
            <select class="admin-input cat-select" v-model="formParentId">
              <option value="">（无 / 顶级）</option>
              <option
                v-for="cat in items"
                :key="cat.id"
                :value="String(cat.id)"
                :disabled="cat.id === editingId"
              >{{ cat.name }}</option>
            </select>
          </label>
          <label class="field">
            <span class="field-label">排序</span>
            <input class="admin-input" v-model.number="formSort" type="number" placeholder="0" />
          </label>
        </div>
        <div class="save-row cat-actions">
          <button class="btn-solid--save" type="submit" :disabled="submitting">
            {{ editingId !== null ? '保存' : '创建' }}
          </button>
          <button class="btn-ghost" type="button" @click="resetForm">
            {{ editingId !== null ? '取消' : '清空' }}
          </button>
        </div>
      </form>
    </div>

    <!-- ── Category Tree List ── -->
    <div class="manage-list">
      <div v-for="{ cat, level } in tree" :key="cat.id" class="manage-row">
        <div class="manage-row__main">
          <div class="manage-row__titlerow">
            <button
              class="manage-row__title"
              type="button"
              :style="level > 0 ? { paddingLeft: `${level * 1.5}rem` } : {}"
              @click="startEdit(cat)"
            >
              <span v-if="level > 0" class="cat-indent">└ </span>{{ cat.name }}
            </button>
          </div>
          <div class="manage-row__tags">
            {{ cat.slug }}<template v-if="cat.parentId != null"> · 父级 #{{ cat.parentId }}</template>
          </div>
        </div>
        <div class="manage-row__acts">
          <button class="act-btn" type="button" @click="startEdit(cat)">编辑</button>
          <button class="act-btn act-btn--del" type="button" @click="remove(cat.id)">删除</button>
        </div>
      </div>
      <p v-if="tree.length === 0" class="manage-empty">暂无分类</p>
    </div>
  </section>
</template>

<style scoped>
.cat-form {
  background: var(--bg-subtle);
  border: 1px solid var(--line);
  border-radius: var(--radius-sm, 4px);
  padding: 20px 24px 24px;
  margin: 24px 0 32px;
  max-width: 540px;
}

.cat-err {
  font-size: var(--text-sm);
  color: #b23b3b;
  margin: 0 0 12px;
}

.cat-req {
  color: #b23b3b;
}

.cat-indent {
  color: var(--ink-2, #5c5750);
}

.cat-actions {
  display: flex;
  gap: 10px;
  align-items: center;
}

.cat-select {
  cursor: pointer;
}
</style>
