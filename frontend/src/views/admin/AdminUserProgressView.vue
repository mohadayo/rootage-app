<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getUserProgress, resetPassword, type UserProgressItem } from '@/api/admin'
import { useConfirm } from '@/composables/useConfirm'

const users = ref<UserProgressItem[]>([])
const loading = ref(true)
const expandedId = ref<string | null>(null)
const resetUserId = ref<string | null>(null)
const newPassword = ref('')
const resetting = ref(false)
const { open: confirm } = useConfirm()

async function handleResetPassword(userId: string, userName: string) {
  resetUserId.value = userId
  newPassword.value = ''
}

async function submitReset() {
  if (!resetUserId.value || !newPassword.value) return
  resetting.value = true
  try {
    await resetPassword(resetUserId.value, newPassword.value)
    alert('パスワードをリセットしました')
    resetUserId.value = null
    newPassword.value = ''
  } catch (e: any) {
    alert(e.response?.data?.error || 'リセットに失敗しました')
  } finally {
    resetting.value = false
  }
}

onMounted(async () => {
  try {
    const { data } = await getUserProgress()
    users.value = data || []
  } finally {
    loading.value = false
  }
})

function toggleExpand(id: string) {
  expandedId.value = expandedId.value === id ? null : id
}

function formatDate(d: string | null) {
  if (!d) return '-'
  return new Date(d).toLocaleDateString('ja-JP')
}

function pctColor(pct: number): string {
  if (pct >= 80) return 'text-green-600'
  if (pct >= 60) return 'text-yellow-600'
  return 'text-red-600'
}

function isAllPassed(user: UserProgressItem): boolean {
  if (!user.category_stats || user.category_stats.length < 4) return false
  return user.category_stats.every(cs => cs.percentage >= 80)
}
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold text-rootage-text mb-6">ユーザー進捗</h1>

    <div v-if="loading" class="text-center text-rootage-sub py-12">読み込み中...</div>

    <div v-else-if="users.length === 0" class="text-center text-rootage-muted py-12">
      登録ユーザーがいません
    </div>

    <!-- PC: テーブル表示 -->
    <div v-else class="hidden sm:block overflow-x-auto">
      <table class="w-full text-sm">
        <thead>
          <tr class="border-b border-rootage-rule">
            <th class="text-left py-3 px-2 font-medium text-rootage-sub">名前</th>
            <th class="text-left py-3 px-2 font-medium text-rootage-sub">メール</th>
            <th class="text-left py-3 px-2 font-medium text-rootage-sub">登録日</th>
            <th class="text-left py-3 px-2 font-medium text-rootage-sub">最終クイズ</th>
            <th class="text-right py-3 px-2 font-medium text-rootage-sub">総合正答率</th>
            <th class="text-right py-3 px-2 font-medium text-rootage-sub">操作</th>
          </tr>
        </thead>
        <tbody>
          <template v-for="user in users" :key="user.id">
            <tr
              @click="toggleExpand(user.id)"
              class="border-b border-rootage-rule-faint hover:bg-rootage-warm cursor-pointer transition-colors"
            >
              <td class="py-3 px-2 font-medium text-rootage-text">
                {{ user.name }}
                <span v-if="isAllPassed(user)" class="ml-1 text-xs font-bold text-green-700 bg-green-100 px-2 py-0.5 rounded-full">全合格</span>
              </td>
              <td class="py-3 px-2 text-rootage-sub">{{ user.email }}</td>
              <td class="py-3 px-2 text-rootage-sub">{{ formatDate(user.created_at) }}</td>
              <td class="py-3 px-2 text-rootage-sub">{{ formatDate(user.last_quiz_at) }}</td>
              <td class="py-3 px-2 text-right font-bold" :class="user.total_questions > 0 ? pctColor(user.percentage) : 'text-rootage-muted'">
                <template v-if="user.total_questions > 0">
                  {{ Math.round(user.percentage) }}%
                  <span class="font-normal text-rootage-muted">({{ user.total_correct }}/{{ user.total_questions }})</span>
                </template>
                <template v-else>未受験</template>
              </td>
              <td class="py-3 px-2 text-right">
                <button @click.stop="handleResetPassword(user.id, user.name)" class="text-xs text-rootage-sub hover:text-rootage-text">PW再設定</button>
              </td>
            </tr>
            <tr v-if="expandedId === user.id">
              <td :colspan="5" class="py-3 px-2 bg-rootage-bg">
                <div v-if="user.category_stats && user.category_stats.length > 0" class="flex flex-wrap gap-4 pl-4">
                  <div v-for="cs in user.category_stats" :key="cs.category_id" class="text-center">
                    <p class="text-lg font-bold" :class="pctColor(cs.percentage)">{{ Math.round(cs.percentage) }}%</p>
                    <p class="text-xs text-rootage-sub">{{ cs.category_name }}</p>
                    <p class="text-xs text-rootage-muted">{{ cs.correct }}/{{ cs.total_answers }}</p>
                    <p v-if="cs.percentage >= 80" class="text-xs font-bold text-green-700 mt-0.5">合格</p>
                  </div>
                </div>
                <p v-else class="text-rootage-muted text-xs pl-4">カテゴリ別データなし</p>
              </td>
            </tr>
          </template>
        </tbody>
      </table>
    </div>

    <!-- スマホ: カード表示 -->
    <div v-if="!loading && users.length > 0" class="sm:hidden space-y-3">
      <div
        v-for="user in users"
        :key="user.id"
        @click="toggleExpand(user.id)"
        class="bg-white rounded-lg border border-rootage-rule p-4 cursor-pointer"
      >
        <div class="flex items-center justify-between mb-2">
          <div>
            <p class="font-medium text-rootage-text">
              {{ user.name }}
              <span v-if="isAllPassed(user)" class="ml-1 text-xs font-bold text-green-700 bg-green-100 px-2 py-0.5 rounded-full">全合格</span>
            </p>
            <p class="text-xs text-rootage-muted">{{ user.email }}</p>
          </div>
          <p class="text-lg font-bold" :class="user.total_questions > 0 ? pctColor(user.percentage) : 'text-rootage-muted'">
            <template v-if="user.total_questions > 0">{{ Math.round(user.percentage) }}%</template>
            <template v-else>-</template>
          </p>
        </div>
        <div class="flex gap-4 text-xs text-rootage-sub">
          <span>登録: {{ formatDate(user.created_at) }}</span>
          <span>最終: {{ formatDate(user.last_quiz_at) }}</span>
        </div>

        <!-- カテゴリ別詳細 -->
        <div v-if="expandedId === user.id" class="mt-3 pt-3 border-t border-rootage-rule-faint">
          <div v-if="user.category_stats && user.category_stats.length > 0" class="grid grid-cols-2 gap-3">
            <div v-for="cs in user.category_stats" :key="cs.category_id" class="text-center bg-rootage-bg rounded p-2">
              <p class="text-lg font-bold" :class="pctColor(cs.percentage)">{{ Math.round(cs.percentage) }}%</p>
              <p class="text-xs text-rootage-sub">{{ cs.category_name }}</p>
              <p class="text-xs text-rootage-muted">{{ cs.correct }}/{{ cs.total_answers }}</p>
              <p v-if="cs.percentage >= 80" class="text-xs font-bold text-green-700 mt-0.5">合格</p>
            </div>
          </div>
          <p v-else class="text-rootage-muted text-xs">カテゴリ別データなし</p>
        </div>
      </div>
    </div>

    <!-- パスワードリセットモーダル -->
    <div v-if="resetUserId" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
      <div class="bg-white rounded-lg p-6 w-full max-w-sm">
        <h2 class="text-lg font-bold mb-4">パスワード再設定</h2>
        <form @submit.prevent="submitReset" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">新しいパスワード</label>
            <input v-model="newPassword" type="text" required minlength="6" placeholder="6文字以上" class="w-full px-3 py-2 border border-rootage-rule rounded-md" />
          </div>
          <div class="flex justify-end gap-2">
            <button type="button" @click="resetUserId = null" class="px-4 py-2 text-gray-600 hover:bg-rootage-rule-faint rounded-md">キャンセル</button>
            <button type="submit" :disabled="resetting" class="px-4 py-2 bg-rootage-dark text-white rounded-md hover:bg-rootage-accent-h disabled:opacity-50">
              {{ resetting ? '処理中...' : 'リセット' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>
