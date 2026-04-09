<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getGuide, type Guide } from '@/api/guide'
import { marked } from 'marked'

const route = useRoute()
const guide = ref<Guide | null>(null)
const loading = ref(true)
const error = ref('')

const renderedContent = computed(() => {
  if (!guide.value) return ''
  return marked(guide.value.content)
})

onMounted(async () => {
  try {
    const { data } = await getGuide(route.params.id as string)
    guide.value = data
  } catch {
    error.value = 'ガイドが見つかりませんでした'
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div>
    <router-link to="/guides" class="inline-flex items-center gap-1 text-sm text-rootage-sub hover:text-rootage-text transition-colors mb-6">
      <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M15 19l-7-7 7-7" /></svg>
      ガイド一覧
    </router-link>

    <div v-if="loading" class="text-center text-rootage-sub py-12">読み込み中...</div>

    <div v-else-if="error" class="text-center text-red-600 py-12">{{ error }}</div>

    <article v-else-if="guide">
      <div class="mb-6">
        <span class="text-xs text-rootage-accent bg-rootage-warm px-2 py-0.5 rounded">{{ guide.category_name }}</span>
        <h1 class="text-2xl font-bold text-rootage-text mt-2 leading-snug">{{ guide.title }}</h1>
        <p class="text-xs text-rootage-muted mt-2">{{ new Date(guide.updated_at).toLocaleDateString('ja-JP') }} 更新</p>
      </div>

      <div class="bg-white rounded-lg border border-rootage-rule p-6 sm:p-8">
        <div class="prose prose-sm max-w-none" v-html="renderedContent"></div>
      </div>
    </article>
  </div>
</template>

<style scoped>
.prose :deep(h1) { font-size: 1.5rem; font-weight: 700; margin-top: 2rem; margin-bottom: 0.75rem; color: #1A1A1A; }
.prose :deep(h2) { font-size: 1.2rem; font-weight: 700; margin-top: 2rem; margin-bottom: 0.5rem; color: #1A1A1A; border-bottom: 1px solid #E5E2DD; padding-bottom: 0.375rem; }
.prose :deep(h3) { font-size: 1.05rem; font-weight: 600; margin-top: 1.5rem; margin-bottom: 0.5rem; color: #1A1A1A; }
.prose :deep(p) { margin-bottom: 0.875rem; line-height: 1.85; color: #404040; }
.prose :deep(ul), .prose :deep(ol) { margin-bottom: 0.875rem; padding-left: 1.5rem; }
.prose :deep(li) { margin-bottom: 0.375rem; line-height: 1.75; color: #404040; }
.prose :deep(ul) { list-style-type: disc; }
.prose :deep(ol) { list-style-type: decimal; }
.prose :deep(code) { background: #F3F1EE; padding: 0.15rem 0.4rem; border-radius: 0.25rem; font-size: 0.85em; }
.prose :deep(pre) { background: #1A1A1A; color: #E8E8E8; padding: 1rem 1.25rem; border-radius: 0.5rem; overflow-x: auto; margin-bottom: 1rem; font-size: 0.85rem; line-height: 1.6; }
.prose :deep(pre code) { background: none; padding: 0; color: inherit; }
.prose :deep(blockquote) { border-left: 3px solid #D8D8D8; padding-left: 1rem; color: #6B6B6B; margin-bottom: 0.875rem; font-style: italic; }
.prose :deep(table) { width: 100%; border-collapse: collapse; margin-bottom: 1rem; font-size: 0.9rem; }
.prose :deep(th), .prose :deep(td) { border: 1px solid #E5E2DD; padding: 0.5rem 0.75rem; text-align: left; }
.prose :deep(th) { background: #F8F7F5; font-weight: 600; }
.prose :deep(a) { color: #3d3a37; text-decoration: underline; text-underline-offset: 2px; }
.prose :deep(a:hover) { color: #1A1A1A; }
.prose :deep(strong) { font-weight: 700; }
.prose :deep(hr) { border: none; border-top: 1px solid #E5E2DD; margin: 1.5rem 0; }
</style>
