<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { getAdminSummary, type AdminSummary } from '@/api/admin'
import { enableRAG } from '@/config'

const summary = ref<AdminSummary | null>(null)

const allLinks = [
  {
    title: '新人の学習状況',
    description: '学習進捗・クイズ成績を確認',
    to: '/admin/user-progress',
    icon: '<path stroke-linecap="round" stroke-linejoin="round" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z" />',
  },
  {
    title: 'クイズ問題の編集',
    description: '問題の作成・編集・CSV一括インポート',
    to: '/admin/questions',
    icon: '<path stroke-linecap="round" stroke-linejoin="round" d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />',
  },
  {
    title: 'クイズカテゴリの編集',
    description: 'カテゴリの追加・変更・削除',
    to: '/admin/categories',
    icon: '<path stroke-linecap="round" stroke-linejoin="round" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />',
  },
  {
    title: '実践ガイドの編集',
    description: 'ガイド記事の作成・編集・公開',
    to: '/admin/guides',
    icon: '<path stroke-linecap="round" stroke-linejoin="round" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />',
  },
  {
    title: '社内ナレッジの登録',
    description: 'AIが検索に使う社内文書の管理',
    to: '/admin/documents',
    icon: '<path stroke-linecap="round" stroke-linejoin="round" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />',
    ragOnly: true,
  },
]

const links = computed(() => enableRAG ? allLinks : allLinks.filter(l => !l.ragOnly))

onMounted(async () => {
  try {
    const { data } = await getAdminSummary()
    summary.value = data
  } catch {}
})
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold text-rootage-text mb-6">管理画面</h1>

    <!-- サマリー + メニューを統合 -->
    <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
      <!-- サマリーカード -->
      <div v-if="summary" class="sm:col-span-2 bg-white rounded-lg border border-rootage-rule p-5">
        <div class="grid grid-cols-4 gap-4">
          <div class="text-center">
            <p class="text-3xl font-bold text-rootage-text">{{ summary.total_users }}</p>
            <p class="text-xs text-rootage-muted mt-1">ユーザー</p>
          </div>
          <div class="text-center">
            <p class="text-3xl font-bold text-rootage-text">{{ summary.total_questions }}</p>
            <p class="text-xs text-rootage-muted mt-1">クイズ問題</p>
          </div>
          <div class="text-center">
            <p class="text-3xl font-bold text-rootage-text">{{ summary.total_guides }}</p>
            <p class="text-xs text-rootage-muted mt-1">ガイド記事</p>
          </div>
          <div class="text-center">
            <p class="text-3xl font-bold text-rootage-text">{{ summary.total_documents }}</p>
            <p class="text-xs text-rootage-muted mt-1">RAG文書</p>
          </div>
        </div>
      </div>

      <!-- メニューカード -->
      <router-link
        v-for="link in links"
        :key="link.to"
        :to="link.to"
        class="flex items-center gap-4 p-5 bg-white rounded-lg border border-rootage-rule hover:border-rootage-accent transition-colors group"
      >
        <span class="w-10 h-10 rounded-lg flex-shrink-0 flex items-center justify-center bg-rootage-rule-faint text-rootage-sub transition-colors">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24" v-html="link.icon"></svg>
        </span>
        <div>
          <h2 class="text-sm font-semibold text-rootage-text group-hover:text-rootage-accent-h transition-colors">{{ link.title }}</h2>
          <p class="text-xs text-rootage-muted">{{ link.description }}</p>
        </div>
      </router-link>
    </div>
  </div>
</template>
