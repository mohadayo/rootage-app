<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import RootageLogo from '@/components/RootageLogo.vue'

const auth = useAuthStore()
const router = useRouter()
const email = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

async function handleLogin() {
  error.value = ''
  loading.value = true
  try {
    await auth.login(email.value, password.value)
    router.push('/')
  } catch (e: any) {
    error.value = e.response?.data?.error || 'ログインに失敗しました'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-rootage-bg -mt-6">
    <div class="w-full max-w-md p-8 bg-white rounded-lg border border-rootage-rule">
      <div class="flex justify-center mb-4">
        <RootageLogo :size="36" textClass="text-xl text-rootage-text" />
      </div>
      <p class="text-center text-rootage-muted text-sm mb-6">新人エンジニア育成アプリ</p>

      <form @submit.prevent="handleLogin" class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-rootage-sub mb-1">メールアドレス</label>
          <input
            v-model="email"
            type="email"
            required
            class="w-full px-3 py-2 border border-rootage-rule rounded-md focus:outline-none focus:ring-2 focus:ring-rootage-accent bg-rootage-bg"
            placeholder="example@example.com"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-rootage-sub mb-1">パスワード</label>
          <input
            v-model="password"
            type="password"
            required
            class="w-full px-3 py-2 border border-rootage-rule rounded-md focus:outline-none focus:ring-2 focus:ring-rootage-accent bg-rootage-bg"
            placeholder="6文字以上"
          />
        </div>

        <div v-if="error" class="text-red-600 text-sm bg-red-50 p-3 rounded">{{ error }}</div>

        <button
          type="submit"
          :disabled="loading"
          class="w-full py-2.5 px-4 bg-rootage-dark text-white rounded-md hover:bg-rootage-accent-h disabled:opacity-50 transition-colors"
        >
          {{ loading ? 'ログイン中...' : 'ログイン' }}
        </button>
      </form>

      <p class="mt-4 text-center text-sm text-rootage-muted">
        アカウントをお持ちでない方は
        <router-link to="/register" class="text-rootage-text hover:underline font-medium">新規登録</router-link>
      </p>
      <p class="mt-2 text-center text-xs text-rootage-muted">
        パスワードを忘れた方は
        <router-link to="/forgot-password" class="text-rootage-text hover:underline">こちら</router-link>
      </p>
    </div>
  </div>
</template>
