<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useConfirm } from '@/composables/useConfirm'
import { getDocuments, uploadDocument, deleteDocument, reindexDocument, createDocumentFromText } from '@/api/admin'
import type { Document } from '@/api/admin'

const documents = ref<Document[]>([])
const loading = ref(true)
const showUpload = ref(false)
const showTextInput = ref(false)
const uploading = ref(false)
const reindexing = ref<string | null>(null)
const title = ref('')
const textContent = ref('')
const textTitle = ref('')
const savingText = ref(false)
const fileInput = ref<HTMLInputElement>()
const selectedFile = ref<File | null>(null)
const { open: confirm } = useConfirm()

onMounted(() => load())

async function load() {
  loading.value = true
  try {
    const { data } = await getDocuments()
    documents.value = data || []
  } finally {
    loading.value = false
  }
}

function onFileChange(e: Event) {
  const target = e.target as HTMLInputElement
  selectedFile.value = target.files?.[0] || null
}

async function upload() {
  if (!title.value || !selectedFile.value) return
  uploading.value = true
  try {
    await uploadDocument(title.value, selectedFile.value)
    showUpload.value = false
    title.value = ''
    selectedFile.value = null
    await load()
  } catch (e: any) {
    alert(e.response?.data?.error || 'アップロードに失敗しました')
  } finally {
    uploading.value = false
  }
}

async function remove(id: string) {
  const ok = await confirm('文書の削除', 'この文書を削除しますか？この操作は取り消せません。')
  if (!ok) return
  try {
    await deleteDocument(id)
    await load()
  } catch (e: any) {
    alert(e.response?.data?.error || '削除に失敗しました')
  }
}

async function reindex(id: string) {
  reindexing.value = id
  try {
    await reindexDocument(id)
    alert('再インデックスが完了しました')
  } catch (e: any) {
    alert(e.response?.data?.error || '再インデックスに失敗しました')
  } finally {
    reindexing.value = null
  }
}

async function saveText() {
  if (!textTitle.value || !textContent.value) return
  savingText.value = true
  try {
    await createDocumentFromText(textTitle.value, textContent.value)
    showTextInput.value = false
    textTitle.value = ''
    textContent.value = ''
    await load()
  } catch (e: any) {
    alert(e.response?.data?.error || '保存に失敗しました')
  } finally {
    savingText.value = false
  }
}

function formatDate(d: string) {
  return new Date(d).toLocaleString('ja-JP')
}
</script>

<template>
  <div>
    <div class="mb-6">
      <div class="flex items-center justify-between mb-2">
        <h1 class="text-2xl font-bold text-rootage-text">社内ナレッジの登録</h1>
        <div class="flex gap-2">
          <button @click="showTextInput = true" class="px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 text-sm">
            テキスト入力
          </button>
          <button @click="showUpload = true" class="px-4 py-2 bg-rootage-dark text-white rounded-md hover:bg-rootage-accent-h text-sm">
            ファイルアップロード
          </button>
        </div>
      </div>
      <p class="text-sm text-rootage-muted">ここに登録した文書は「社内検索AI」でAIが回答に使います</p>
    </div>

    <!-- Text input modal -->
    <div v-if="showTextInput" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
      <div class="bg-white rounded-lg p-6 w-full max-w-2xl max-h-[90vh] overflow-y-auto">
        <h2 class="text-lg font-bold mb-4">テキスト入力で文書登録</h2>
        <p class="text-sm text-rootage-sub mb-4">自社の情報をテキストで直接入力できます。コピペもOK。</p>
        <form @submit.prevent="saveText" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">タイトル</label>
            <input v-model="textTitle" required placeholder="例: 会社概要、福利厚生、研修制度" class="w-full px-3 py-2 border border-rootage-rule rounded-md" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">内容</label>
            <textarea v-model="textContent" required rows="12" placeholder="自社の情報をここに入力またはコピペしてください..." class="w-full px-3 py-2 border border-rootage-rule rounded-md"></textarea>
          </div>
          <div class="flex justify-end gap-2">
            <button type="button" @click="showTextInput = false" class="px-4 py-2 text-gray-600 hover:bg-rootage-rule-faint rounded-md">キャンセル</button>
            <button type="submit" :disabled="savingText" class="px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 disabled:opacity-50">
              {{ savingText ? '保存中...' : '保存してインデックス作成' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Upload modal -->
    <div v-if="showUpload" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
      <div class="bg-white rounded-lg p-6 w-full max-w-md">
        <h2 class="text-lg font-bold mb-2">ファイルから登録</h2>
        <p class="text-sm text-rootage-muted mb-4">社内資料のテキストファイルをアップロードします</p>
        <form @submit.prevent="upload" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">タイトル</label>
            <input v-model="title" required placeholder="例: 就業規則、福利厚生一覧" class="w-full px-3 py-2 border border-rootage-rule rounded-md" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">ファイル (.txt, .md, .csv)</label>
            <input ref="fileInput" type="file" accept=".txt,.md,.csv" @change="onFileChange" required class="w-full text-sm" />
          </div>
          <div class="flex justify-end gap-2">
            <button type="button" @click="showUpload = false" class="px-4 py-2 text-gray-600 hover:bg-rootage-rule-faint rounded-md">キャンセル</button>
            <button type="submit" :disabled="uploading" class="px-4 py-2 bg-rootage-dark text-white rounded-md hover:bg-rootage-accent-h disabled:opacity-50">
              {{ uploading ? 'アップロード中...' : 'アップロード' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <div v-if="loading" class="text-center text-rootage-sub py-12">読み込み中...</div>

    <div v-else-if="documents.length === 0" class="text-center text-rootage-sub py-12">
      <p>文書がありません</p>
    </div>

    <div v-else class="space-y-3">
      <div v-for="doc in documents" :key="doc.id" class="bg-white rounded-lg shadow-sm border border-rootage-rule p-4 flex items-center justify-between">
        <div>
          <p class="font-semibold text-rootage-text">{{ doc.title }}</p>
          <p class="text-xs text-rootage-muted">{{ formatDate(doc.uploaded_at) }} 登録</p>
        </div>
        <div class="flex gap-2">
          <button
            @click="reindex(doc.id)"
            :disabled="reindexing === doc.id"
            class="text-sm text-rootage-accent hover:underline disabled:opacity-50"
          >
            {{ reindexing === doc.id ? '処理中...' : '再インデックス' }}
          </button>
          <button @click="remove(doc.id)" class="text-sm text-red-600 hover:underline">削除</button>
        </div>
      </div>
    </div>
  </div>
</template>
