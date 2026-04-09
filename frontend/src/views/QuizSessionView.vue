<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter, onBeforeRouteLeave } from 'vue-router'
import { useQuizStore } from '@/stores/quiz'
import { useConfirm } from '@/composables/useConfirm'

const router = useRouter()
const quizStore = useQuizStore()
const selectedIndex = ref<number | null>(null)
const answered = ref(false)
const lastAnswer = ref<{ is_correct: boolean; correct_index: number; explanation: string } | null>(null)
const submitting = ref(false)

const question = computed(() => quizStore.currentQuestion())
const progress = computed(() => `${quizStore.currentIndex + 1} / ${quizStore.questions.length}`)
const progressPercent = computed(() => ((quizStore.currentIndex + 1) / quizStore.questions.length) * 100)
const isLast = computed(() => quizStore.currentIndex >= quizStore.questions.length - 1)

function choiceClass(index: number) {
  if (!answered.value) {
    return index === selectedIndex.value
      ? 'border-rootage-accent bg-rootage-warm'
      : 'border-rootage-rule hover:border-rootage-rule'
  }
  if (index === lastAnswer.value?.correct_index) return 'border-green-500 bg-green-50'
  if (index === selectedIndex.value && !lastAnswer.value?.is_correct) return 'border-red-500 bg-red-50'
  return 'border-rootage-rule opacity-50'
}

async function submit() {
  if (selectedIndex.value === null || submitting.value) return
  submitting.value = true
  try {
    const result = await quizStore.submitAnswer(selectedIndex.value)
    lastAnswer.value = result
    answered.value = true
  } finally {
    submitting.value = false
  }
}

async function next() {
  if (isLast.value) {
    const result = await quizStore.finishQuiz()
    quizFinished.value = true
    router.replace(`/quiz/${quizStore.sessionId}/result`)
    return
  }
  quizStore.nextQuestion()
  selectedIndex.value = null
  answered.value = false
  lastAnswer.value = null
}

const quizFinished = ref(false)
const { open: confirm } = useConfirm()

onBeforeRouteLeave(async () => {
  if (quizFinished.value || !quizStore.sessionId) return true
  return await confirm('クイズを中断しますか？', '途中で離れると進捗は保存されません。', { label: '中断する', color: 'bg-rootage-dark hover:bg-rootage-accent-h' })
})

if (!quizStore.sessionId) {
  router.replace('/categories')
}
</script>

<template>
  <div v-if="question">
    <!-- Progress bar -->
    <div class="mb-6">
      <div class="flex justify-between text-sm text-rootage-sub mb-1">
        <span>{{ progress }}</span>
        <span>{{ Math.round(progressPercent) }}%</span>
      </div>
      <div class="w-full bg-gray-200 rounded-full h-2">
        <div class="bg-rootage-dark h-2 rounded-full transition-all" :style="{ width: progressPercent + '%' }"></div>
      </div>
    </div>

    <!-- Question -->
    <div class="bg-white rounded-lg shadow-sm border border-rootage-rule p-6 mb-4">
      <h2 class="text-lg font-semibold text-rootage-text mb-4">{{ question.text }}</h2>

      <div class="space-y-3">
        <button
          v-for="(choice, i) in question.choices"
          :key="i"
          @click="!answered && (selectedIndex = i)"
          :disabled="answered"
          :class="[
            'w-full text-left p-4 rounded-lg border-2 transition-colors',
            choiceClass(i),
          ]"
        >
          <span class="font-medium text-rootage-muted mr-2">{{ ['A', 'B', 'C', 'D'][i] }}.</span>
          {{ choice }}
        </button>
      </div>
    </div>

    <!-- Explanation -->
    <div v-if="answered && lastAnswer" class="bg-white rounded-lg shadow-sm border p-6 mb-4"
         :class="lastAnswer.is_correct ? 'border-green-300 bg-green-50' : 'border-red-300 bg-red-50'">
      <p class="font-semibold mb-2" :class="lastAnswer.is_correct ? 'text-green-700' : 'text-red-700'">
        {{ lastAnswer.is_correct ? '正解!' : '不正解' }}
      </p>
      <p class="text-gray-700 text-sm">{{ lastAnswer.explanation }}</p>
    </div>

    <!-- Action buttons -->
    <div class="flex justify-end gap-3">
      <button
        v-if="!answered"
        @click="submit"
        :disabled="selectedIndex === null || submitting"
        class="px-6 py-2 bg-rootage-dark text-white rounded-md hover:bg-rootage-accent-h disabled:opacity-50"
      >
        {{ submitting ? '送信中...' : '回答する' }}
      </button>
      <button
        v-else
        @click="next"
        class="px-6 py-2 bg-rootage-dark text-white rounded-md hover:bg-rootage-accent-h"
      >
        {{ isLast ? '結果を見る' : '次の問題' }}
      </button>
    </div>
  </div>
</template>
