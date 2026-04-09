<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useConfirm } from '@/composables/useConfirm'
import { getAdminQuestions, createQuestion, updateQuestion, deleteQuestion, getAdminCategories, importQuestions } from '@/api/admin'
import type { Question, ImportResult } from '@/api/admin'
import type { Category } from '@/api/quiz'

const questions = ref<Question[]>([])
const categories = ref<Category[]>([])
const loading = ref(true)
const showForm = ref(false)
const editingId = ref<string | null>(null)
const filterCategory = ref('')
const showImport = ref(false)
const importFile = ref<File | null>(null)
const importing = ref(false)
const importResult = ref<ImportResult | null>(null)
const { open: confirm } = useConfirm()

const form = ref({
  category_id: '',
  text: '',
  choices: ['', '', '', ''],
  correct_index: 0,
  explanation: '',
  difficulty: 'beginner',
})

onMounted(async () => {
  const [catRes] = await Promise.all([getAdminCategories()])
  categories.value = catRes.data || []
  await load()
})

async function load() {
  loading.value = true
  try {
    const { data } = await getAdminQuestions(filterCategory.value || undefined)
    questions.value = data || []
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editingId.value = null
  form.value = { category_id: categories.value[0]?.id || '', text: '', choices: ['', '', '', ''], correct_index: 0, explanation: '', difficulty: 'beginner' }
  showForm.value = true
}

function openEdit(q: Question) {
  editingId.value = q.id
  form.value = {
    category_id: q.category_id,
    text: q.text,
    choices: [...q.choices],
    correct_index: q.correct_index,
    explanation: q.explanation,
    difficulty: (q as any).difficulty || 'beginner',
  }
  showForm.value = true
}

async function save() {
  try {
    if (editingId.value) {
      await updateQuestion(editingId.value, {
        text: form.value.text,
        choices: form.value.choices,
        correct_index: form.value.correct_index,
        explanation: form.value.explanation,
        difficulty: form.value.difficulty,
      })
    } else {
      await createQuestion(form.value)
    }
    showForm.value = false
    await load()
  } catch (e: any) {
    alert(e.response?.data?.error || '保存に失敗しました')
  }
}

async function remove(id: string) {
  const ok = await confirm('問題の削除', 'この問題を削除しますか？この操作は取り消せません。')
  if (!ok) return
  try {
    await deleteQuestion(id)
    await load()
  } catch (e: any) {
    alert(e.response?.data?.error || '削除に失敗しました')
  }
}

function categoryName(id: string) {
  return categories.value.find(c => c.id === id)?.name || ''
}

function downloadTemplate() {
  const header = 'category_name,text,choice_a,choice_b,choice_c,choice_d,correct,explanation'
  const example = '技術基礎,HTTPステータスコード200は？,成功,失敗,リダイレクト,認証エラー,A,200はリクエスト成功を意味します'
  const csv = header + '\n' + example + '\n'
  const bom = '\uFEFF'
  const blob = new Blob([bom + csv], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = 'questions_template.csv'
  a.click()
  URL.revokeObjectURL(url)
}

function onFileSelect(e: Event) {
  const target = e.target as HTMLInputElement
  importFile.value = target.files?.[0] || null
}

async function handleImport() {
  if (!importFile.value) return
  importing.value = true
  importResult.value = null
  try {
    const { data } = await importQuestions(importFile.value)
    importResult.value = data
    if (data.imported > 0) {
      await load()
    }
  } catch (e: any) {
    alert(e.response?.data?.error || 'インポートに失敗しました')
  } finally {
    importing.value = false
  }
}

function closeImport() {
  showImport.value = false
  importFile.value = null
  importResult.value = null
}
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold text-rootage-text">問題管理</h1>
      <div class="flex gap-2">
        <button @click="showImport = true" class="px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 text-sm">
          CSV一括インポート
        </button>
        <button @click="openCreate" class="px-4 py-2 bg-rootage-dark text-white rounded-md hover:bg-rootage-accent-h text-sm">
          新規作成
        </button>
      </div>
    </div>

    <!-- Filter -->
    <div class="mb-4">
      <select v-model="filterCategory" @change="load" class="px-3 py-2 border border-rootage-rule rounded-md text-sm">
        <option value="">全カテゴリ</option>
        <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
      </select>
    </div>

    <!-- Import modal -->
    <div v-if="showImport" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
      <div class="bg-white rounded-lg p-6 w-full max-w-lg">
        <h2 class="text-lg font-bold mb-4">CSV一括インポート</h2>

        <div class="mb-4">
          <button @click="downloadTemplate" class="text-sm text-rootage-accent hover:underline">
            テンプレートCSVをダウンロード
          </button>
          <p class="text-xs text-rootage-sub mt-1">
            CSVフォーマット: category_name, text, choice_a, choice_b, choice_c, choice_d, correct(A/B/C/D), explanation
          </p>
        </div>

        <div class="mb-4">
          <input type="file" accept=".csv" @change="onFileSelect" class="w-full text-sm text-rootage-sub file:mr-4 file:py-2 file:px-4 file:rounded-md file:border-0 file:text-sm file:font-semibold file:bg-rootage-warm file:text-rootage-accent hover:file:bg-rootage-rule-faint" />
        </div>

        <div v-if="importResult" class="mb-4 space-y-2">
          <p class="text-sm font-medium text-green-700">{{ importResult.imported }}件 インポート成功</p>
          <div v-if="importResult.errors && importResult.errors.length > 0">
            <p class="text-sm font-medium text-red-600">{{ importResult.errors.length }}件 エラー:</p>
            <ul class="text-xs text-red-600 mt-1 max-h-32 overflow-y-auto space-y-1">
              <li v-for="(err, i) in importResult.errors" :key="i" class="bg-red-50 p-2 rounded">
                {{ err.row }}行目: {{ err.message }}
              </li>
            </ul>
          </div>
        </div>

        <div class="flex justify-end gap-2">
          <button @click="closeImport" class="px-4 py-2 text-gray-600 hover:bg-rootage-rule-faint rounded-md">閉じる</button>
          <button @click="handleImport" :disabled="!importFile || importing" class="px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 disabled:opacity-50">
            {{ importing ? 'インポート中...' : 'インポート実行' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Form modal -->
    <div v-if="showForm" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
      <div class="bg-white rounded-lg p-6 w-full max-w-lg max-h-[90vh] overflow-y-auto">
        <h2 class="text-lg font-bold mb-4">{{ editingId ? '問題編集' : '問題作成' }}</h2>
        <form @submit.prevent="save" class="space-y-4">
          <div v-if="!editingId">
            <label class="block text-sm font-medium text-gray-700 mb-1">カテゴリ</label>
            <select v-model="form.category_id" required class="w-full px-3 py-2 border border-rootage-rule rounded-md">
              <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">問題文</label>
            <textarea v-model="form.text" required rows="3" class="w-full px-3 py-2 border border-rootage-rule rounded-md"></textarea>
          </div>
          <div v-for="(_, i) in form.choices" :key="i">
            <label class="block text-sm font-medium text-gray-700 mb-1">
              選択肢{{ ['A', 'B', 'C', 'D'][i] }}
              <span v-if="i === form.correct_index" class="text-green-600">(正解)</span>
            </label>
            <input v-model="form.choices[i]" required class="w-full px-3 py-2 border border-rootage-rule rounded-md" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">正解</label>
            <select v-model="form.correct_index" class="w-full px-3 py-2 border border-rootage-rule rounded-md">
              <option :value="0">A</option>
              <option :value="1">B</option>
              <option :value="2">C</option>
              <option :value="3">D</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">難易度</label>
            <select v-model="form.difficulty" class="w-full px-3 py-2 border border-rootage-rule rounded-md">
              <option value="beginner">初級</option>
              <option value="intermediate">中級</option>
              <option value="advanced">上級</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">解説</label>
            <textarea v-model="form.explanation" required rows="3" class="w-full px-3 py-2 border border-rootage-rule rounded-md"></textarea>
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
      <div v-for="q in questions" :key="q.id" class="bg-white rounded-lg shadow-sm border border-rootage-rule p-4">
        <div class="flex items-start justify-between">
          <div class="flex-1">
            <span class="text-xs text-rootage-accent bg-rootage-warm px-2 py-0.5 rounded">{{ categoryName(q.category_id) }}</span>
            <p class="font-medium text-rootage-text mt-1">{{ q.text }}</p>
            <p class="text-sm text-rootage-sub mt-1">正解: {{ ['A', 'B', 'C', 'D'][q.correct_index] }}</p>
          </div>
          <div class="flex gap-2 ml-4">
            <button @click="openEdit(q)" class="text-sm text-rootage-accent hover:underline">編集</button>
            <button @click="remove(q.id)" class="text-sm text-red-600 hover:underline">削除</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
