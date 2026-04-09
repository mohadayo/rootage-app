<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { getPublishedGuides, type Guide } from '@/api/guide'

const guides = ref<Guide[]>([])
const loading = ref(true)
const searchQuery = ref('')

onMounted(async () => {
  try {
    const { data } = await getPublishedGuides()
    guides.value = data || []
  } finally {
    loading.value = false
  }
})

const filteredGuides = computed(() => {
  const q = searchQuery.value.trim().toLowerCase()
  if (!q) return guides.value
  return guides.value.filter(g =>
    g.title.toLowerCase().includes(q) ||
    g.category_name.toLowerCase().includes(q) ||
    g.content.toLowerCase().includes(q)
  )
})

const groupedGuides = computed(() => {
  const groups: Record<string, { name: string; guides: Guide[] }> = {}
  for (const g of filteredGuides.value) {
    if (!groups[g.category_id]) {
      groups[g.category_id] = { name: g.category_name, guides: [] }
    }
    groups[g.category_id]!.guides.push(g)
  }
  return Object.values(groups)
})

function excerpt(content: string, maxLen = 60): string {
  const text = content.replace(/[#*`>\-|[\]()]/g, '').replace(/\n+/g, ' ').trim()
  if (text.length <= maxLen) return text
  return text.slice(0, maxLen) + '...'
}
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold text-rootage-text mb-1">実践ガイド</h1>
    <p class="text-rootage-muted text-sm mb-5">現場で困ったときに見返せるナレッジ集</p>

    <!-- 検索 -->
    <div v-if="guides.length > 0" class="relative mb-6">
      <input
        v-model="searchQuery"
        type="text"
        placeholder="キーワードで検索..."
        class="w-full px-4 py-2.5 pl-10 border border-rootage-rule rounded-lg focus:outline-none focus:ring-2 focus:ring-rootage-accent bg-white text-sm"
      />
      <svg class="w-4 h-4 text-rootage-muted absolute left-3.5 top-3" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
      </svg>
    </div>

    <div v-if="loading" class="text-center text-rootage-sub py-12">読み込み中...</div>

    <div v-else-if="guides.length === 0" class="text-center text-rootage-muted py-12">
      <p class="mb-2">まだガイドが登録されていません</p>
      <p class="text-sm">管理者がガイド記事を公開すると、ここに表示されます</p>
    </div>

    <div v-else-if="filteredGuides.length === 0" class="text-center text-rootage-muted py-12">
      <p class="mb-2">「{{ searchQuery }}」に一致するガイドが見つかりません</p>
      <button @click="searchQuery = ''" class="text-rootage-accent hover:underline text-sm">検索をクリア</button>
    </div>

    <div v-else class="space-y-6">
      <section v-for="group in groupedGuides" :key="group.name">
        <h2 class="text-lg font-bold text-rootage-text mb-3 flex items-center gap-2">
          <span class="w-2 h-2 bg-rootage-dark"></span>
          {{ group.name }}
        </h2>
        <div class="space-y-2">
          <router-link
            v-for="guide in group.guides"
            :key="guide.id"
            :to="`/guides/${guide.id}`"
            class="flex items-center justify-between p-4 bg-white rounded-lg border border-rootage-rule hover:border-rootage-accent transition-colors group"
          >
            <div class="min-w-0">
              <h3 class="font-medium text-sm text-rootage-text group-hover:text-rootage-accent-h transition-colors">{{ guide.title }}</h3>
              <p class="text-xs text-rootage-muted mt-1 truncate">{{ excerpt(guide.content) }}</p>
            </div>
            <svg class="w-4 h-4 text-rootage-muted flex-shrink-0 ml-3" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7" />
            </svg>
          </router-link>
        </div>
      </section>
    </div>
  </div>
</template>
