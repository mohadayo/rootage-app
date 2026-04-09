import { defineStore } from 'pinia'
import { ref } from 'vue'
import * as quizApi from '@/api/quiz'
import type { QuizQuestion, AnswerResponse } from '@/api/quiz'

export const useQuizStore = defineStore('quiz', () => {
  const sessionId = ref('')
  const categoryId = ref('')
  const difficulty = ref('')
  const questions = ref<QuizQuestion[]>([])
  const currentIndex = ref(0)
  const answers = ref<(AnswerResponse & { selectedIndex: number })[]>([])
  const isFinished = ref(false)

  const currentQuestion = () => questions.value[currentIndex.value]

  async function startQuiz(catId: string, diff?: string) {
    const { data } = await quizApi.startQuiz(catId, diff)
    sessionId.value = data.session_id
    categoryId.value = catId
    difficulty.value = diff || 'beginner'
    questions.value = data.questions
    currentIndex.value = 0
    answers.value = []
    isFinished.value = false
  }

  async function submitAnswer(selectedIndex: number) {
    const question = questions.value[currentIndex.value]!
    const { data } = await quizApi.submitAnswer(sessionId.value, question.id, selectedIndex)
    answers.value.push({ ...data, selectedIndex })
    return data
  }

  function nextQuestion() {
    if (currentIndex.value < questions.value.length - 1) {
      currentIndex.value++
    }
  }

  async function finishQuiz() {
    const { data } = await quizApi.finishQuiz(sessionId.value)
    isFinished.value = true
    return data
  }

  function reset() {
    sessionId.value = ''
    categoryId.value = ''
    difficulty.value = ''
    questions.value = []
    currentIndex.value = 0
    answers.value = []
    isFinished.value = false
  }

  return {
    sessionId, categoryId, difficulty, questions, currentIndex, answers, isFinished,
    currentQuestion, startQuiz, submitAnswer, nextQuestion, finishQuiz, reset,
  }
})
