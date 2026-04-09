<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getReview, type ReviewItem } from '@/api/quiz'

const items = ref<ReviewItem[]>([])
const loading = ref(true)
const expandedId = ref<string | null>(null)

onMounted(async () => {
  try {
    const { data } = await getReview()
    items.value = data || []
  } finally {
    loading.value = false
  }
})

function toggle(id: string) {
  expandedId.value = expandedId.value === id ? null : id
}
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold text-rootage-text mb-6">復習 - 不正解問題</h1>

    <div v-if="loading" class="text-center text-rootage-sub py-12">読み込み中...</div>

    <div v-else-if="items.length === 0" class="text-center text-rootage-sub py-12">
      <p>復習する問題がありません</p>
      <router-link to="/categories" class="text-rootage-accent hover:underline mt-2 inline-block">
        クイズを始める
      </router-link>
    </div>

    <div v-else class="space-y-4">
      <div
        v-for="item in items"
        :key="item.question_id"
        class="bg-white rounded-lg shadow-sm border border-rootage-rule overflow-hidden"
      >
        <button
          @click="toggle(item.question_id)"
          class="w-full p-4 text-left flex items-start justify-between"
        >
          <div>
            <span class="text-xs text-rootage-accent bg-rootage-warm px-2 py-0.5 rounded">{{ item.category_name }}</span>
            <p class="font-medium text-rootage-text mt-1">{{ item.question_text }}</p>
          </div>
          <span class="text-rootage-muted ml-2">{{ expandedId === item.question_id ? '-' : '+' }}</span>
        </button>

        <div v-if="expandedId === item.question_id" class="px-4 pb-4 border-t border-gray-100 pt-3">
          <div class="space-y-2 mb-3">
            <div
              v-for="(choice, i) in item.choices"
              :key="i"
              :class="[
                'p-2 rounded text-sm',
                i === item.correct_index ? 'bg-green-50 border border-green-300 text-green-800' : 'bg-rootage-bg text-gray-600'
              ]"
            >
              <span class="font-medium mr-1">{{ ['A', 'B', 'C', 'D'][i] }}.</span> {{ choice }}
              <span v-if="i === item.correct_index" class="ml-1 text-green-600 font-medium">(正解)</span>
            </div>
          </div>
          <div class="bg-blue-50 p-3 rounded text-sm text-blue-800">
            <p class="font-medium mb-1">解説</p>
            <p>{{ item.explanation }}</p>
          </div>
        </div>
      </div>

    </div>
  </div>
</template>
