<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getCategories, getUserDashboard, getReview, type Category, type UserDashboardResponse, type ReviewItem } from '@/api/quiz'
import { useQuizStore } from '@/stores/quiz'

const router = useRouter()
const quizStore = useQuizStore()
const categories = ref<Category[]>([])
const dashboard = ref<UserDashboardResponse | null>(null)
const reviewItems = ref<ReviewItem[]>([])
const loading = ref(true)
const loadError = ref(false)
const starting = ref('')
const expandedReview = ref<string | null>(null)
const selectedDifficulty = ref<Record<string, string>>({})

const colors = ['bg-blue-500', 'bg-green-500', 'bg-yellow-500', 'bg-purple-500']

const difficultyLabels: Record<string, string> = {
  beginner: '初級',
  intermediate: '中級',
  advanced: '上級',
}

onMounted(async () => {
  try {
    const [catRes, dashRes, revRes] = await Promise.all([
      getCategories(),
      getUserDashboard(),
      getReview(),
    ])
    categories.value = catRes.data
    dashboard.value = dashRes.data
    reviewItems.value = (revRes.data || []).slice(0, 5)
  } catch {
    loadError.value = true
  } finally {
    loading.value = false
  }
})

const difficultyOrder = ['beginner', 'intermediate', 'advanced']

// カテゴリ一覧からアチーブメントグリッドを構築
const achievementGrid = computed(() => {
  return categories.value.map(cat => {
    const rows = difficultyOrder.map(diff => {
      const a = dashboard.value?.achievements?.find(
        item => item.category_id === cat.id && item.difficulty === diff
      )
      return {
        difficulty: diff,
        label: difficultyLabels[diff],
        bestScore: a?.best_score ?? null,
        total: a?.total ?? null,
        percentage: a?.percentage ?? null,
        cleared: a?.cleared ?? false,
        attempted: !!a,
      }
    })
    const allCleared = rows.every(r => r.cleared)
    return { id: cat.id, name: cat.name, difficulties: rows, allCleared }
  })
})

const totalCleared = computed(() => {
  return achievementGrid.value.reduce((sum, cat) => sum + cat.difficulties.filter(d => d.cleared).length, 0)
})

const totalSlots = computed(() => achievementGrid.value.length * 3)

function getDifficulty(categoryId: string): string {
  return selectedDifficulty.value[categoryId] || 'beginner'
}

async function start(categoryId: string) {
  starting.value = categoryId
  try {
    await quizStore.startQuiz(categoryId, getDifficulty(categoryId))
    router.push(`/quiz/${quizStore.sessionId}`)
  } catch (e: any) {
    alert(e.response?.data?.error || 'クイズの開始に失敗しました')
  } finally {
    starting.value = ''
  }
}

function getCategoryPct(categoryId: string): number | null {
  const diff = getDifficulty(categoryId)
  const a = dashboard.value?.achievements?.find(
    item => item.category_id === categoryId && item.difficulty === diff
  )
  return a ? a.percentage : null
}

function pctColor(pct: number): string {
  if (pct >= 80) return 'text-green-600'
  if (pct >= 60) return 'text-yellow-600'
  return 'text-red-600'
}

function barColor(pct: number): string {
  if (pct >= 80) return 'bg-green-500'
  if (pct >= 60) return 'bg-yellow-500'
  return 'bg-red-400'
}

</script>

<template>
  <div>
    <h1 class="text-2xl font-bold text-rootage-text mb-4">カテゴリを選択</h1>

    <div v-if="loading" class="text-center text-rootage-sub py-12">読み込み中...</div>

    <div v-else-if="loadError" class="text-center text-red-600 py-12">
      <p>データの読み込みに失敗しました</p>
      <button @click="$router.go(0)" class="mt-2 text-sm text-rootage-accent hover:underline">再読み込み</button>
    </div>

    <div v-else class="grid grid-cols-1 sm:grid-cols-2 gap-4">
      <div
        v-for="(cat, i) in categories"
        :key="cat.id"
        class="p-6 bg-white rounded-lg shadow-sm border border-rootage-rule"
      >
        <div class="flex items-center justify-between mb-3">
          <div class="flex items-center gap-3">
            <span :class="[colors[i % colors.length], 'w-3 h-3 rounded-full']"></span>
            <h2 class="text-lg font-semibold text-rootage-text">{{ cat.name }}</h2>
          </div>
          <div class="flex items-center gap-2">
            <span v-if="getCategoryPct(cat.id) !== null && getCategoryPct(cat.id)! >= 80" class="text-xs font-bold text-green-700 bg-green-100 px-2 py-0.5 rounded-full">合格</span>
            <span
              v-if="getCategoryPct(cat.id) !== null"
              class="text-sm font-bold"
              :class="pctColor(getCategoryPct(cat.id)!)"
            >
              {{ Math.round(getCategoryPct(cat.id)!) }}%
            </span>
          </div>
        </div>
        <p class="text-rootage-sub text-sm mb-3">{{ cat.description }}</p>
        <div class="flex items-center gap-1.5 mb-3">
          <button
            v-for="(label, key) in difficultyLabels"
            :key="key"
            @click="selectedDifficulty[cat.id] = key"
            :class="[
              'px-2.5 py-1 rounded-full text-xs font-medium transition-colors',
              getDifficulty(cat.id) === key
                ? 'bg-rootage-dark text-white'
                : 'bg-gray-100 text-gray-500 hover:bg-gray-200'
            ]"
          >
            {{ label }}
          </button>
        </div>
        <button
          @click="start(cat.id)"
          :disabled="starting === cat.id"
          class="w-full py-2 bg-rootage-dark text-white rounded-md hover:bg-rootage-accent-h disabled:opacity-50 text-sm"
        >
          {{ starting === cat.id ? '準備中...' : `${difficultyLabels[getDifficulty(cat.id)]}でクイズを始める (10問)` }}
        </button>
      </div>
    </div>

    <!-- マイ進捗（アチーブメント） -->
    <template v-if="dashboard && categories.length > 0">
      <div class="flex items-center justify-between mt-10 mb-2 pb-2 border-b-2 border-rootage-rule">
        <h2 class="text-lg font-bold text-rootage-text">マイ進捗</h2>
        <span class="text-sm text-rootage-muted">{{ totalCleared }}/{{ totalSlots }} クリア</span>
      </div>
      <p class="text-xs text-rootage-muted mb-4">各カテゴリ・難易度で80%以上取るとクリアになります</p>

      <div class="space-y-3">
        <div
          v-for="cat in achievementGrid"
          :key="cat.id"
          class="bg-white rounded-lg border border-rootage-rule p-4"
        >
          <div class="flex items-center gap-2 mb-3">
            <span class="font-semibold text-sm text-rootage-text">{{ cat.name }}</span>
            <span v-if="cat.allCleared" class="text-xs font-bold text-green-700 bg-green-100 px-2 py-0.5 rounded-full">全制覇</span>
          </div>
          <div class="grid grid-cols-3 gap-2">
            <div
              v-for="d in cat.difficulties"
              :key="d.difficulty"
              :class="[
                'rounded-lg px-3 py-2 text-center transition-all',
                d.cleared
                  ? 'bg-green-50 border border-green-300'
                  : d.attempted
                    ? 'bg-rootage-warm border border-rootage-rule'
                    : 'bg-rootage-rule-faint border border-rootage-rule-faint'
              ]"
            >
              <p class="text-xs font-medium" :class="d.cleared ? 'text-green-700' : 'text-rootage-sub'">{{ d.label }}</p>
              <template v-if="d.attempted">
                <p class="text-base font-bold" :class="pctColor(d.percentage!)">{{ Math.round(d.percentage!) }}%</p>
                <p class="text-xs text-rootage-muted">{{ d.bestScore }}/{{ d.total }}</p>
              </template>
              <template v-else>
                <p class="text-base text-rootage-rule">--</p>
                <p class="text-xs text-rootage-muted">未挑戦</p>
              </template>
            </div>
          </div>
        </div>
      </div>
    </template>

    <!-- 復習 -->
    <h2 class="text-lg font-bold text-rootage-text mt-10 mb-4 pb-2 border-b-2 border-rootage-rule">復習（間違えた問題）</h2>

    <div v-if="reviewItems.length > 0" class="space-y-2">
      <div
        v-for="item in reviewItems"
        :key="item.question_id"
        class="bg-white rounded-lg border border-rootage-rule overflow-hidden"
      >
        <button
          @click="expandedReview = expandedReview === item.question_id ? null : item.question_id"
          class="w-full p-4 text-left flex items-start justify-between"
        >
          <div>
            <span class="text-xs text-rootage-accent bg-rootage-warm px-2 py-0.5 rounded">{{ item.category_name }}</span>
            <p class="text-sm text-rootage-text mt-1">{{ item.question_text }}</p>
          </div>
          <span class="text-rootage-muted ml-2 flex-shrink-0">{{ expandedReview === item.question_id ? '-' : '+' }}</span>
        </button>
        <div v-if="expandedReview === item.question_id" class="px-4 pb-4 border-t border-rootage-rule-faint pt-3">
          <div class="space-y-2 mb-3">
            <div
              v-for="(choice, ci) in item.choices"
              :key="ci"
              :class="[
                'p-2 rounded text-sm',
                ci === item.correct_index ? 'bg-green-50 border border-green-300 text-green-800' : 'bg-rootage-bg text-gray-600'
              ]"
            >
              <span class="font-medium mr-1">{{ ['A', 'B', 'C', 'D'][ci] }}.</span> {{ choice }}
              <span v-if="ci === item.correct_index" class="ml-1 text-green-600 font-medium">(正解)</span>
            </div>
          </div>
          <div class="bg-blue-50 p-3 rounded text-sm text-blue-800">
            <p class="font-medium mb-1">解説</p>
            <p>{{ item.explanation }}</p>
          </div>
        </div>
      </div>
      <div class="text-right mt-2">
        <router-link to="/quiz/review" class="text-sm text-rootage-sub hover:text-rootage-text">
          すべての復習問題を見る →
        </router-link>
      </div>
    </div>
    <p v-else class="text-sm text-rootage-muted">クイズで間違えた問題がここに表示されます</p>

  </div>
</template>
