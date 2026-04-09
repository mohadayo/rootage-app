<script setup lang="ts">
import { ref } from 'vue'
import { forgotPassword } from '@/api/auth'
import RootageLogo from '@/components/RootageLogo.vue'

const email = ref('')
const loading = ref(false)
const sent = ref(false)
const error = ref('')

async function handleSubmit() {
  error.value = ''
  loading.value = true
  try {
    await forgotPassword(email.value)
    sent.value = true
  } catch {
    error.value = '送信に失敗しました。しばらく経ってから再度お試しください。'
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
      <p class="text-center text-rootage-muted text-sm mb-6">パスワードリセット</p>

      <template v-if="!sent">
        <p class="text-sm text-rootage-sub mb-4">登録したメールアドレスを入力してください。リセット用のリンクを送信します。</p>

        <form @submit.prevent="handleSubmit" class="space-y-4">
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

          <div v-if="error" class="text-red-600 text-sm bg-red-50 p-3 rounded">{{ error }}</div>

          <button
            type="submit"
            :disabled="loading"
            class="w-full py-2.5 px-4 bg-rootage-dark text-white rounded-md hover:bg-rootage-accent-h disabled:opacity-50 transition-colors"
          >
            {{ loading ? '送信中...' : 'リセットリンクを送信' }}
          </button>
        </form>
      </template>

      <template v-else>
        <div class="text-center py-4">
          <p class="text-green-600 font-medium mb-2">メールを送信しました</p>
          <p class="text-sm text-rootage-sub">登録されているメールアドレスの場合、リセット用のリンクが届きます。メールをご確認ください。</p>
        </div>
      </template>

      <p class="mt-4 text-center text-sm text-rootage-muted">
        <router-link to="/login" class="text-rootage-text hover:underline font-medium">ログイン画面に戻る</router-link>
      </p>
    </div>
  </div>
</template>
