<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { resetPassword } from '@/api/auth'
import RootageLogo from '@/components/RootageLogo.vue'

const route = useRoute()
const router = useRouter()
const password = ref('')
const confirmPassword = ref('')
const loading = ref(false)
const error = ref('')
const success = ref(false)
const token = ref('')

onMounted(() => {
  token.value = (route.query.token as string) || ''
  if (!token.value) {
    error.value = 'リセットリンクが無効です'
  }
})

async function handleSubmit() {
  error.value = ''

  if (password.value !== confirmPassword.value) {
    error.value = 'パスワードが一致しません'
    return
  }

  loading.value = true
  try {
    await resetPassword(token.value, password.value)
    success.value = true
    setTimeout(() => router.push('/login'), 3000)
  } catch (e: any) {
    error.value = e.response?.data?.error || 'パスワードの再設定に失敗しました'
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
      <p class="text-center text-rootage-muted text-sm mb-6">パスワード再設定</p>

      <template v-if="success">
        <div class="text-center py-4">
          <p class="text-green-600 font-medium mb-2">パスワードを再設定しました</p>
          <p class="text-sm text-rootage-sub">3秒後にログイン画面に移動します...</p>
          <router-link to="/login" class="mt-4 inline-block text-rootage-text hover:underline font-medium text-sm">
            すぐにログインする
          </router-link>
        </div>
      </template>

      <template v-else-if="!token">
        <div class="text-center py-4">
          <p class="text-red-600 mb-2">リセットリンクが無効です</p>
          <router-link to="/forgot-password" class="text-rootage-text hover:underline font-medium text-sm">
            もう一度リセットリンクを送信する
          </router-link>
        </div>
      </template>

      <template v-else>
        <form @submit.prevent="handleSubmit" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-rootage-sub mb-1">新しいパスワード</label>
            <input
              v-model="password"
              type="password"
              required
              minlength="6"
              class="w-full px-3 py-2 border border-rootage-rule rounded-md focus:outline-none focus:ring-2 focus:ring-rootage-accent bg-rootage-bg"
              placeholder="6文字以上"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-rootage-sub mb-1">パスワード確認</label>
            <input
              v-model="confirmPassword"
              type="password"
              required
              minlength="6"
              class="w-full px-3 py-2 border border-rootage-rule rounded-md focus:outline-none focus:ring-2 focus:ring-rootage-accent bg-rootage-bg"
              placeholder="もう一度入力"
            />
          </div>

          <div v-if="error" class="text-red-600 text-sm bg-red-50 p-3 rounded">{{ error }}</div>

          <button
            type="submit"
            :disabled="loading"
            class="w-full py-2.5 px-4 bg-rootage-dark text-white rounded-md hover:bg-rootage-accent-h disabled:opacity-50 transition-colors"
          >
            {{ loading ? '処理中...' : 'パスワードを再設定する' }}
          </button>
        </form>
      </template>

      <p class="mt-4 text-center text-sm text-rootage-muted">
        <router-link to="/login" class="text-rootage-text hover:underline font-medium">ログイン画面に戻る</router-link>
      </p>
    </div>
  </div>
</template>
