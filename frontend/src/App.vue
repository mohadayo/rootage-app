<script setup lang="ts">
import { RouterView } from 'vue-router'
import AppHeader from '@/components/layout/AppHeader.vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
</script>

<template>
  <div class="min-h-screen bg-rootage-bg">
    <AppHeader v-if="auth.isLoggedIn" />
    <main class="max-w-4xl mx-auto px-4 py-6">
      <RouterView />
    </main>
    <ConfirmDialog />
  </div>
</template>
