<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useConfirm } from '@/composables/useConfirm'
import { getAdminCategories, createCategory, updateCategory, deleteCategory } from '@/api/admin'
import type { Category } from '@/api/quiz'

const categories = ref<Category[]>([])
const loading = ref(true)
const showForm = ref(false)
const editingId = ref<string | null>(null)
const form = ref({ name: '', description: '' })
const { open: confirm } = useConfirm()

onMounted(() => load())

async function load() {
  loading.value = true
  try {
    const { data } = await getAdminCategories()
    categories.value = data || []
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editingId.value = null
  form.value = { name: '', description: '' }
  showForm.value = true
}

function openEdit(cat: Category) {
  editingId.value = cat.id
  form.value = { name: cat.name, description: cat.description }
  showForm.value = true
}

async function save() {
  try {
    if (editingId.value) {
      await updateCategory(editingId.value, form.value)
    } else {
      await createCategory(form.value)
    }
    showForm.value = false
    await load()
  } catch (e: any) {
    alert(e.response?.data?.error || '保存に失敗しました')
  }
}

async function remove(id: string) {
  const ok = await confirm('カテゴリの削除', 'このカテゴリを削除しますか？関連する問題もすべて削除されます。')
  if (!ok) return
  try {
    await deleteCategory(id)
    await load()
  } catch (e: any) {
    alert(e.response?.data?.error || '削除に失敗しました')
  }
}
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold text-rootage-text">カテゴリ管理</h1>
      <button @click="openCreate" class="px-4 py-2 bg-rootage-dark text-white rounded-md hover:bg-rootage-accent-h text-sm">
        新規作成
      </button>
    </div>

    <!-- Form modal -->
    <div v-if="showForm" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
      <div class="bg-white rounded-lg p-6 w-full max-w-md">
        <h2 class="text-lg font-bold mb-4">{{ editingId ? 'カテゴリ編集' : 'カテゴリ作成' }}</h2>
        <form @submit.prevent="save" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">名前</label>
            <input v-model="form.name" required class="w-full px-3 py-2 border border-rootage-rule rounded-md" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">説明</label>
            <textarea v-model="form.description" rows="3" class="w-full px-3 py-2 border border-rootage-rule rounded-md"></textarea>
          </div>
          <div class="flex justify-end gap-2">
            <button type="button" @click="showForm = false" class="px-4 py-2 text-gray-600 hover:bg-rootage-rule-faint rounded-md">キャンセル</button>
            <button type="submit" class="px-4 py-2 bg-rootage-dark text-white rounded-md hover:bg-rootage-accent-h">保存</button>
          </div>
        </form>
      </div>
    </div>

    <div v-if="loading" class="text-center text-rootage-sub py-12">読み込み中...</div>
    <div v-else class="space-y-3">
      <div v-for="cat in categories" :key="cat.id" class="bg-white rounded-lg shadow-sm border border-rootage-rule p-4 flex items-center justify-between">
        <div>
          <p class="font-semibold text-rootage-text">{{ cat.name }}</p>
          <p class="text-sm text-rootage-sub">{{ cat.description }}</p>
        </div>
        <div class="flex gap-2">
          <button @click="openEdit(cat)" class="text-sm text-rootage-accent hover:underline">編集</button>
          <button @click="remove(cat.id)" class="text-sm text-red-600 hover:underline">削除</button>
        </div>
      </div>
    </div>
  </div>
</template>
