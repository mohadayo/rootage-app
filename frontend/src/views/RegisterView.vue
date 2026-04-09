<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import RootageLogo from '@/components/RootageLogo.vue'

const auth = useAuthStore()
const router = useRouter()
const name = ref('')
const email = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

const allowedDomain = import.meta.env.VITE_ALLOWED_DOMAIN || ''

async function handleRegister() {
  error.value = ''
  if (!email.value.toLowerCase().endsWith('@' + allowedDomain)) {
    error.value = `@${allowedDomain} のメールアドレスのみ登録できます`
    return
  }
  loading.value = true
  try {
    await auth.register(email.value, password.value, name.value)
    router.push('/')
  } catch (e: any) {
    error.value = e.response?.data?.error || '登録に失敗しました'
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
      <p class="text-center text-rootage-muted text-sm mb-6">新規アカウント登録</p>

      <form @submit.prevent="handleRegister" class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-rootage-sub mb-1">名前</label>
          <input
            v-model="name"
            type="text"
            required
            class="w-full px-3 py-2 border border-rootage-rule rounded-md focus:outline-none focus:ring-2 focus:ring-rootage-accent bg-rootage-bg"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-rootage-sub mb-1">メールアドレス</label>
          <input
            v-model="email"
            type="email"
            required
            :placeholder="`example@${allowedDomain}`"
            class="w-full px-3 py-2 border border-rootage-rule rounded-md focus:outline-none focus:ring-2 focus:ring-rootage-accent bg-rootage-bg"
          />
          <p class="mt-1 text-xs text-rootage-muted">{{ `@${allowedDomain}` }} のアドレスのみ登録できます</p>
        </div>
        <div>
          <label class="block text-sm font-medium text-rootage-sub mb-1">パスワード</label>
          <input
            v-model="password"
            type="password"
            required
            minlength="6"
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
          {{ loading ? '登録中...' : '登録する' }}
        </button>
      </form>

      <p class="mt-4 text-center text-sm text-rootage-muted">
        既にアカウントをお持ちの方は
        <router-link to="/login" class="text-rootage-text hover:underline font-medium">ログイン</router-link>
      </p>
    </div>
  </div>
</template>
