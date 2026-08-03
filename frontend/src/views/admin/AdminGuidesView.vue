<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useConfirm } from '@/composables/useConfirm'
import {
  getAdminGuides, getAdminGuideCategories, createGuide, updateGuide, deleteGuide,
  createGuideCategory, updateGuideCategory, deleteGuideCategory,
  type Guide, type GuideCategory,
} from '@/api/admin'
import { renderMarkdown } from '@/markdown'

const guides = ref<Guide[]>([])
const categories = ref<GuideCategory[]>([])
const loading = ref(true)
const showForm = ref(false)
const showCategoryForm = ref(false)
const showPreview = ref(false)
const editingId = ref<string | null>(null)
const editingCategoryId = ref<string | null>(null)

const form = ref({
  category_id: '',
  title: '',
  content: '',
  is_published: false,
  sort_order: 0,
})

const categoryForm = ref({ name: '', sort_order: 0 })
const { open: confirm } = useConfirm()

const previewHtml = computed(() => renderMarkdown(form.value.content))

onMounted(async () => {
  await loadAll()
})

async function loadAll() {
  loading.value = true
  try {
    const [catRes, guideRes] = await Promise.all([getAdminGuideCategories(), getAdminGuides()])
    categories.value = catRes.data || []
    guides.value = guideRes.data || []
  } finally {
    loading.value = false
  }
}

// カテゴリ管理
function openCreateCategory() {
  editingCategoryId.value = null
  categoryForm.value = { name: '', sort_order: 0 }
  showCategoryForm.value = true
}

function openEditCategory(c: GuideCategory) {
  editingCategoryId.value = c.id
  categoryForm.value = { name: c.name, sort_order: c.sort_order }
  showCategoryForm.value = true
}

async function saveCategory() {
  try {
    if (editingCategoryId.value) {
      await updateGuideCategory(editingCategoryId.value, categoryForm.value)
    } else {
      await createGuideCategory(categoryForm.value)
    }
    showCategoryForm.value = false
    await loadAll()
  } catch (e: any) {
    alert(e.response?.data?.error || '保存に失敗しました')
  }
}

async function removeCategory(id: string) {
  const ok = await confirm('カテゴリの削除', 'このカテゴリとその中のガイドを全て削除しますか？')
  if (!ok) return
  try {
    await deleteGuideCategory(id)
    await loadAll()
  } catch (e: any) {
    alert(e.response?.data?.error || '削除に失敗しました')
  }
}

// ガイド管理
function openCreate() {
  editingId.value = null
  form.value = { category_id: categories.value[0]?.id || '', title: '', content: '', is_published: false, sort_order: 0 }
  showPreview.value = false
  showForm.value = true
}

function openEdit(g: Guide) {
  editingId.value = g.id
  form.value = {
    category_id: g.category_id,
    title: g.title,
    content: g.content,
    is_published: g.is_published,
    sort_order: g.sort_order,
  }
  showPreview.value = false
  showForm.value = true
}

async function save() {
  try {
    if (editingId.value) {
      await updateGuide(editingId.value, form.value)
    } else {
      await createGuide(form.value)
    }
    showForm.value = false
    await loadAll()
  } catch (e: any) {
    alert(e.response?.data?.error || '保存に失敗しました')
  }
}

async function remove(id: string) {
  const ok = await confirm('ガイドの削除', 'このガイドを削除しますか？この操作は取り消せません。')
  if (!ok) return
  try {
    await deleteGuide(id)
    await loadAll()
  } catch (e: any) {
    alert(e.response?.data?.error || '削除に失敗しました')
  }
}

function categoryName(id: string) {
  return categories.value.find(c => c.id === id)?.name || ''
}
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold text-rootage-text">ガイド管理</h1>
      <div class="flex gap-2">
        <button @click="openCreateCategory" class="px-4 py-2 bg-rootage-warm text-rootage-text rounded-md hover:bg-rootage-rule-faint text-sm border border-rootage-rule">
          カテゴリ管理
        </button>
        <button @click="openCreate" class="px-4 py-2 bg-rootage-dark text-white rounded-md hover:bg-rootage-accent-h text-sm">
          新規作成
        </button>
      </div>
    </div>

    <!-- カテゴリ一覧 -->
    <div class="mb-6">
      <h2 class="text-sm font-semibold text-rootage-sub mb-2">カテゴリ</h2>
      <div class="flex flex-wrap gap-2">
        <div v-for="cat in categories" :key="cat.id" class="flex items-center gap-1 bg-rootage-warm px-3 py-1 rounded-full text-sm">
          <span class="text-rootage-text">{{ cat.name }}</span>
          <button @click="openEditCategory(cat)" class="text-rootage-muted hover:text-rootage-text ml-1">...</button>
          <button @click="removeCategory(cat.id)" class="text-rootage-muted hover:text-red-600">x</button>
        </div>
        <button v-if="categories.length === 0" @click="openCreateCategory" class="text-sm text-rootage-accent hover:underline">
          + カテゴリを作成
        </button>
      </div>
    </div>

    <!-- カテゴリフォームモーダル -->
    <div v-if="showCategoryForm" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
      <div class="bg-white rounded-lg p-6 w-full max-w-md">
        <h2 class="text-lg font-bold mb-4">{{ editingCategoryId ? 'カテゴリ編集' : 'カテゴリ作成' }}</h2>
        <form @submit.prevent="saveCategory" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">カテゴリ名</label>
            <input v-model="categoryForm.name" required class="w-full px-3 py-2 border border-rootage-rule rounded-md" placeholder="例: ロードマップ" />
          </div>
          <div class="flex justify-end gap-2">
            <button type="button" @click="showCategoryForm = false" class="px-4 py-2 text-gray-600 hover:bg-rootage-rule-faint rounded-md">キャンセル</button>
            <button type="submit" class="px-4 py-2 bg-rootage-dark text-white rounded-md hover:bg-rootage-accent-h">保存</button>
          </div>
        </form>
      </div>
    </div>

    <!-- ガイド作成・編集モーダル -->
    <div v-if="showForm" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
      <div class="bg-white rounded-lg p-6 w-full max-w-4xl max-h-[90vh] overflow-y-auto">
        <h2 class="text-lg font-bold mb-4">{{ editingId ? 'ガイド編集' : 'ガイド作成' }}</h2>
        <form @submit.prevent="save" class="space-y-4">
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">カテゴリ</label>
              <select v-model="form.category_id" required class="w-full px-3 py-2 border border-rootage-rule rounded-md">
                <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
              </select>
            </div>
            <div class="flex items-end gap-4">
              <label class="flex items-center gap-2 pb-2 cursor-pointer">
                <input v-model="form.is_published" type="checkbox" class="w-4 h-4" />
                <span class="text-sm font-medium" :class="form.is_published ? 'text-green-600' : 'text-rootage-muted'">
                  {{ form.is_published ? '公開' : '下書き' }}
                </span>
              </label>
            </div>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">タイトル</label>
            <input v-model="form.title" required class="w-full px-3 py-2 border border-rootage-rule rounded-md" placeholder="例: 現場初日の過ごし方" />
          </div>
          <div>
            <div class="flex items-center justify-between mb-1">
              <label class="block text-sm font-medium text-gray-700">内容 (Markdown)</label>
              <button type="button" @click="showPreview = !showPreview" class="text-xs text-rootage-accent hover:underline">
                {{ showPreview ? '編集に戻る' : 'プレビュー' }}
              </button>
            </div>
            <textarea
              v-if="!showPreview"
              v-model="form.content"
              required
              rows="16"
              class="w-full px-3 py-2 border border-rootage-rule rounded-md font-mono text-sm"
              placeholder="Markdown形式で記事を書いてください..."
            ></textarea>
            <div
              v-else
              class="w-full px-4 py-3 border border-rootage-rule rounded-md bg-rootage-bg min-h-[24rem] overflow-y-auto prose prose-sm max-w-none"
              v-html="previewHtml"
            ></div>
          </div>
          <div class="flex justify-end gap-2">
            <button type="button" @click="showForm = false" class="px-4 py-2 text-gray-600 hover:bg-rootage-rule-faint rounded-md">キャンセル</button>
            <button type="submit" class="px-4 py-2 bg-rootage-dark text-white rounded-md hover:bg-rootage-accent-h">
              {{ form.is_published ? '公開する' : '下書き保存' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- ガイド一覧 -->
    <div v-if="loading" class="text-center text-rootage-sub py-12">読み込み中...</div>
    <div v-else-if="guides.length === 0" class="text-center text-rootage-muted py-12">
      ガイドがありません。「新規作成」から記事を書きましょう。
    </div>
    <div v-else class="space-y-3">
      <div v-for="g in guides" :key="g.id" class="bg-white rounded-lg shadow-sm border border-rootage-rule p-4">
        <div class="flex items-start justify-between">
          <div class="flex-1">
            <div class="flex items-center gap-2">
              <span class="text-xs text-rootage-accent bg-rootage-warm px-2 py-0.5 rounded">{{ categoryName(g.category_id) }}</span>
              <span v-if="g.is_published" class="text-xs text-green-600 bg-green-50 px-2 py-0.5 rounded">公開</span>
              <span v-else class="text-xs text-rootage-muted bg-gray-100 px-2 py-0.5 rounded">下書き</span>
            </div>
            <p class="font-medium text-rootage-text mt-1">{{ g.title }}</p>
            <p class="text-xs text-rootage-muted mt-1">{{ new Date(g.updated_at).toLocaleDateString('ja-JP') }} 更新</p>
          </div>
          <div class="flex gap-2 ml-4">
            <button @click="openEdit(g)" class="text-sm text-rootage-accent hover:underline">編集</button>
            <button @click="remove(g.id)" class="text-sm text-red-600 hover:underline">削除</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
