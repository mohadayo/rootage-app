<script setup lang="ts">
import { useQuizStore } from '@/stores/quiz'
import { useRouter } from 'vue-router'
import { computed } from 'vue'

const router = useRouter()
const quizStore = useQuizStore()

const result = computed(() => {
  const total = quizStore.questions.length
  const correct = quizStore.answers.filter(a => a.is_correct).length
  const percentage = total > 0 ? Math.round((correct / total) * 100) : 0
  return { total, correct, percentage }
})

const cleared = computed(() => result.value.percentage >= 80)

const scoreColor = computed(() => {
  if (result.value.percentage >= 80) return 'text-green-600'
  if (result.value.percentage >= 60) return 'text-yellow-600'
  return 'text-red-600'
})

const difficultyLabels: Record<string, string> = {
  beginner: '初級',
  intermediate: '中級',
  advanced: '上級',
}

const nextDifficulty = computed(() => {
  const order = ['beginner', 'intermediate', 'advanced']
  const idx = order.indexOf(quizStore.difficulty)
  if (idx < order.length - 1) return order[idx + 1]
  return null
})

function retry() {
  quizStore.startQuiz(quizStore.categoryId, quizStore.difficulty)
    .then(() => router.replace(`/quiz/${quizStore.sessionId}`))
}

function tryNextDifficulty() {
  if (!nextDifficulty.value) return
  quizStore.startQuiz(quizStore.categoryId, nextDifficulty.value)
    .then(() => router.replace(`/quiz/${quizStore.sessionId}`))
}

function goHome() {
  quizStore.reset()
  router.push('/categories')
}

if (!quizStore.sessionId) {
  router.replace('/categories')
}
</script>

<template>
  <div class="max-w-lg mx-auto">
    <div class="bg-white rounded-lg shadow-sm border border-rootage-rule p-8 text-center">
      <h1 class="text-2xl font-bold text-rootage-text mb-4">クイズ結果</h1>

      <div class="mb-6">
        <p :class="['text-5xl font-bold mb-2', scoreColor]">{{ result.percentage }}%</p>
        <p class="text-rootage-sub">{{ result.correct }} / {{ result.total }} 問正解</p>
        <p v-if="cleared" class="text-green-600 font-medium mt-2">クリア!</p>
        <p v-else class="text-rootage-muted text-sm mt-2">80%以上でクリアになります</p>
      </div>

      <div class="mb-6 space-y-2">
        <div
          v-for="(ans, i) in quizStore.answers"
          :key="i"
          class="flex items-center gap-2 text-sm"
        >
          <span :class="ans.is_correct ? 'text-green-500' : 'text-red-500'" class="font-bold">
            {{ ans.is_correct ? 'O' : 'X' }}
          </span>
          <span class="text-gray-600">問{{ i + 1 }}</span>
        </div>
      </div>

      <!-- 次のアクション -->
      <div v-if="cleared" class="space-y-2">
        <button
          v-if="nextDifficulty"
          @click="tryNextDifficulty"
          class="w-full py-2 bg-rootage-dark text-white rounded-md hover:bg-rootage-accent-h text-sm"
        >
          {{ difficultyLabels[nextDifficulty!] }}に挑戦する
        </button>
        <button
          @click="goHome"
          class="w-full py-2 bg-gray-100 text-gray-700 rounded-md hover:bg-gray-200 text-sm"
        >
          他のカテゴリに挑戦する
        </button>
      </div>
      <div v-else class="space-y-2">
        <button
          @click="retry"
          class="w-full py-2 bg-rootage-dark text-white rounded-md hover:bg-rootage-accent-h text-sm"
        >
          もう一度挑戦する
        </button>
        <router-link
          to="/quiz/review"
          class="block w-full py-2 bg-gray-100 text-gray-700 rounded-md hover:bg-gray-200 text-sm text-center"
        >
          間違えた問題を復習する
        </router-link>
        <button
          @click="goHome"
          class="w-full py-2 text-sm text-rootage-sub hover:text-rootage-text"
        >
          カテゴリ選択に戻る
        </button>
      </div>
    </div>
  </div>
</template>
