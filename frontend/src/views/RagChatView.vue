<script setup lang="ts">
import { ref, nextTick } from 'vue'
import { askQuestion, type RAGResponse, type RAGSource, type HistoryMessage } from '@/api/rag'

interface Message {
  role: 'user' | 'assistant'
  content: string
  sources?: RAGSource[]
}

const messages = ref<Message[]>([])
const question = ref('')
const loading = ref(false)
const chatContainer = ref<HTMLElement>()

const suggestedQuestions = [
  '有給休暇の申請方法は？',
  '研修プログラムについて教えて',
  '出退勤の打刻方法は？',
  '経費精算のやり方は？',
  '困ったときの相談先は？',
  '福利厚生にはどんなものがある？',
]

function selectSuggestion(q: string) {
  question.value = q
  send()
}

async function send() {
  if (!question.value.trim() || loading.value) return

  const q = question.value.trim()
  question.value = ''
  messages.value.push({ role: 'user', content: q })
  loading.value = true
  scrollToBottom()

  try {
    // 過去の会話履歴を構築（sourcesを除外してrole/contentのみ送る）
    const history: HistoryMessage[] = messages.value.slice(0, -1).map(m => ({
      role: m.role,
      content: m.content,
    }))
    const { data } = await askQuestion(q, history.length > 0 ? history : undefined)
    messages.value.push({
      role: 'assistant',
      content: data.answer,
      sources: data.sources,
    })
  } catch {
    messages.value.push({
      role: 'assistant',
      content: '申し訳ありません。回答の生成に失敗しました。',
    })
  } finally {
    loading.value = false
    scrollToBottom()
  }
}

async function scrollToBottom() {
  await nextTick()
  chatContainer.value?.scrollTo({ top: chatContainer.value.scrollHeight, behavior: 'smooth' })
}
</script>

<template>
  <div class="flex flex-col h-[calc(100vh-8rem)]">
    <h1 class="text-2xl font-bold text-rootage-text mb-4">社内検索AI</h1>

    <!-- Messages -->
    <div ref="chatContainer" class="flex-1 overflow-y-auto space-y-4 mb-4">
      <div v-if="messages.length === 0" class="text-center text-rootage-muted py-12">
        <p class="text-lg mb-2">ルーテイジについて分からないことがあれば何でも聞いてください</p>
        <p class="text-sm mb-6">よくある質問をタップして始められます</p>
        <div class="flex flex-wrap gap-2 justify-center max-w-lg mx-auto">
          <button
            v-for="q in suggestedQuestions"
            :key="q"
            @click="selectSuggestion(q)"
            class="px-4 py-2 bg-white border border-rootage-rule rounded-lg text-sm text-rootage-text hover:border-rootage-accent hover:text-rootage-accent-h transition-colors"
          >
            {{ q }}
          </button>
        </div>
      </div>

      <div v-for="(msg, i) in messages" :key="i" :class="msg.role === 'user' ? 'flex justify-end' : ''">
        <div :class="[
          'max-w-[85%] rounded-lg p-4',
          msg.role === 'user' ? 'bg-rootage-dark text-white' : 'bg-white border border-rootage-rule shadow-sm'
        ]">
          <p class="whitespace-pre-wrap">{{ msg.content }}</p>

          <!-- Sources -->
          <div v-if="msg.sources && msg.sources.length > 0" class="mt-3 pt-3 border-t border-rootage-rule">
            <p class="text-xs font-semibold text-rootage-sub mb-2">参照ソース</p>
            <div v-for="(src, j) in msg.sources" :key="j" class="text-xs bg-rootage-bg rounded p-2 mb-1">
              <p class="font-medium text-rootage-accent">{{ src.document_title }}</p>
              <p class="text-rootage-sub mt-1 line-clamp-2">{{ src.content }}</p>
            </div>
          </div>
        </div>
      </div>

      <div v-if="loading" class="flex">
        <div class="bg-white border border-rootage-rule rounded-lg p-4 shadow-sm">
          <p class="text-rootage-muted">回答を生成中...</p>
        </div>
      </div>
    </div>

    <!-- Input -->
    <form @submit.prevent="send" class="flex gap-2">
      <input
        v-model="question"
        type="text"
        placeholder="質問を入力..."
        :disabled="loading"
        class="flex-1 px-4 py-3 border border-rootage-rule rounded-lg focus:outline-none focus:ring-2 focus:ring-rootage-accent"
      />
      <button
        type="submit"
        :disabled="!question.trim() || loading"
        class="px-6 py-3 bg-rootage-dark text-white rounded-lg hover:bg-rootage-accent-h disabled:opacity-50"
      >
        送信
      </button>
    </form>
  </div>
</template>
