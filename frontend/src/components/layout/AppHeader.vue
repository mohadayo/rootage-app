<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import { ref } from 'vue'
import RootageLogo from '@/components/RootageLogo.vue'
import { enableRAG } from '@/config'

const auth = useAuthStore()
const router = useRouter()
const menuOpen = ref(false)

function logout() {
  auth.logout()
  router.push('/login')
}
</script>

<template>
  <header class="bg-white border-b border-rootage-rule">
    <div class="max-w-4xl mx-auto px-4">
      <div class="flex items-center justify-between h-14">
        <router-link to="/" class="hover:opacity-70 transition-opacity">
          <RootageLogo :size="22" />
        </router-link>

        <!-- Mobile menu button -->
        <button
          @click="menuOpen = !menuOpen"
          class="md:hidden p-2 rounded text-rootage-sub hover:bg-rootage-rule-faint"
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path v-if="!menuOpen" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
            <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>

        <!-- Desktop nav -->
        <nav class="hidden md:flex items-center gap-5 text-sm">
          <router-link to="/" class="text-rootage-sub hover:text-rootage-text transition-colors">ホーム</router-link>
          <router-link to="/categories" class="text-rootage-sub hover:text-rootage-text transition-colors">クイズ</router-link>
          <router-link to="/guides" class="text-rootage-sub hover:text-rootage-text transition-colors">実践ガイド</router-link>
          <router-link v-if="enableRAG" to="/rag" class="text-rootage-sub hover:text-rootage-text transition-colors">社内検索AI</router-link>
          <router-link v-if="auth.isAdmin" to="/admin" class="text-rootage-sub hover:text-rootage-text transition-colors">管理画面</router-link>

          <span class="text-rootage-rule">|</span>
          <span class="text-rootage-muted text-xs">{{ auth.user?.name }}</span>
          <button @click="logout" class="text-rootage-muted hover:text-rootage-text text-xs transition-colors">ログアウト</button>
        </nav>
      </div>

      <!-- Mobile nav -->
      <nav v-if="menuOpen" class="md:hidden pb-4 space-y-2 text-sm">
        <router-link to="/" class="block py-2 text-rootage-sub" @click="menuOpen = false">ホーム</router-link>
        <router-link to="/categories" class="block py-2 text-rootage-sub" @click="menuOpen = false">クイズ</router-link>
        <router-link to="/guides" class="block py-2 text-rootage-sub" @click="menuOpen = false">実践ガイド</router-link>
        <router-link to="/rag" class="block py-2 text-rootage-sub" @click="menuOpen = false">社内検索AI</router-link>
        <router-link v-if="auth.isAdmin" to="/admin" class="block py-2 text-rootage-sub" @click="menuOpen = false">管理画面</router-link>
        <button @click="logout" class="block py-2 text-rootage-muted">ログアウト</button>
      </nav>
    </div>
  </header>
</template>
