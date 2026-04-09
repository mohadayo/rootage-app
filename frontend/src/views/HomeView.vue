<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { getPublishedGuides, type Guide } from '@/api/guide'
import { enableRAG } from '@/config'

const auth = useAuthStore()
const recentGuides = ref<Guide[]>([])

onMounted(async () => {
  try {
    const { data } = await getPublishedGuides()
    recentGuides.value = (data || [])
      .sort((a, b) => new Date(b.updated_at).getTime() - new Date(a.updated_at).getTime())
      .slice(0, 3)
  } catch {}
})
</script>

<template>
  <div>
    <div class="mb-8">
      <h1 class="text-2xl font-bold text-rootage-text">
        {{ auth.user?.name }}さん、こんにちは
      </h1>
      <p class="text-rootage-muted mt-1">学習を始めましょう</p>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
      <!-- クイズ -->
      <router-link
        to="/categories"
        class="block p-6 bg-white rounded-lg border border-rootage-rule hover:border-rootage-accent transition-colors group"
      >
        <div class="flex items-center gap-3 mb-2">
          <span class="w-8 h-8 rounded-md flex items-center justify-center bg-rootage-dark text-white">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
          </span>
          <h2 class="text-lg font-semibold text-rootage-text group-hover:text-rootage-accent-h transition-colors">クイズ</h2>
        </div>
        <p class="text-rootage-sub text-sm">面談対策・技術知識を4択クイズで学ぶ</p>
      </router-link>

      <!-- 実践ガイド -->
      <router-link
        to="/guides"
        class="block p-6 bg-white rounded-lg border border-rootage-rule hover:border-rootage-accent transition-colors group"
      >
        <div class="flex items-center gap-3 mb-2">
          <span class="w-8 h-8 rounded-md flex items-center justify-center bg-rootage-dark text-white">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" /></svg>
          </span>
          <h2 class="text-lg font-semibold text-rootage-text group-hover:text-rootage-accent-h transition-colors">実践ガイド</h2>
        </div>
        <p class="text-rootage-sub text-sm">現場で困ったときに見返せるナレッジ集</p>
      </router-link>

      <!-- 社内検索AI (RAG無効時は非表示) -->
      <template v-if="enableRAG">
      <router-link
        to="/rag"
        class="block p-6 bg-white rounded-lg border border-rootage-rule hover:border-rootage-accent transition-colors group"
      >
        <div class="flex items-center gap-3 mb-2">
          <span class="w-8 h-8 rounded-md flex items-center justify-center bg-rootage-dark text-white">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" /></svg>
          </span>
          <h2 class="text-lg font-semibold text-rootage-text group-hover:text-rootage-accent-h transition-colors">社内検索AI</h2>
        </div>
        <p class="text-rootage-sub text-sm">ルーテイジのことをAIに何でも聞ける</p>
      </router-link>
      </template>

      <!-- 管理画面 -->
      <router-link
        v-if="auth.isAdmin"
        to="/admin"
        class="block p-6 bg-rootage-rule-faint rounded-lg border border-rootage-rule hover:border-rootage-accent transition-colors group"
      >
        <div class="flex items-center gap-3 mb-2">
          <span class="w-8 h-8 rounded-md flex items-center justify-center bg-rootage-dark text-white">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.066 2.573c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.573 1.066c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.066-2.573c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" /><path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" /></svg>
          </span>
          <h2 class="text-lg font-semibold text-rootage-text group-hover:text-rootage-accent-h transition-colors">管理画面</h2>
        </div>
        <p class="text-rootage-sub text-sm">問題・ガイド・文書の管理</p>
      </router-link>
    </div>

    <!-- おすすめガイド -->
    <div v-if="recentGuides.length > 0" class="mt-8">
      <div class="flex items-center justify-between mb-3">
        <h2 class="text-lg font-semibold text-rootage-text">おすすめガイド</h2>
        <router-link to="/guides" class="text-sm text-rootage-sub hover:text-rootage-text">すべて見る →</router-link>
      </div>
      <div class="space-y-2">
        <router-link
          v-for="guide in recentGuides"
          :key="guide.id"
          :to="`/guides/${guide.id}`"
          class="flex items-center justify-between p-4 bg-white rounded-lg border border-rootage-rule hover:border-rootage-accent transition-colors group"
        >
          <div>
            <span class="text-xs text-rootage-accent bg-rootage-warm px-2 py-0.5 rounded">{{ guide.category_name }}</span>
            <h3 class="font-medium text-sm text-rootage-text group-hover:text-rootage-accent-h transition-colors mt-1">{{ guide.title }}</h3>
          </div>
          <svg class="w-4 h-4 text-rootage-muted flex-shrink-0" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7" />
          </svg>
        </router-link>
      </div>
    </div>
  </div>
</template>
