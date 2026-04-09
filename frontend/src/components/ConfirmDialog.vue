<script setup lang="ts">
import { confirmState, useConfirm } from '@/composables/useConfirm'

const { close } = useConfirm()
</script>

<template>
  <Teleport to="body">
    <div v-if="confirmState.visible" class="fixed inset-0 z-50 flex items-center justify-center">
      <div class="absolute inset-0 bg-black/40 transition-opacity" @click="close(false)"></div>
      <div class="relative bg-white rounded-xl shadow-xl border border-rootage-rule w-full max-w-sm mx-4 p-6">
        <h3 class="text-base font-bold text-rootage-text mb-2">{{ confirmState.title }}</h3>
        <p class="text-sm text-rootage-sub mb-6">{{ confirmState.message }}</p>
        <div class="flex gap-3 justify-end">
          <button
            @click="close(false)"
            class="px-4 py-2 text-sm font-medium text-rootage-sub bg-gray-100 rounded-lg hover:bg-gray-200 transition-colors"
          >
            キャンセル
          </button>
          <button
            @click="close(true)"
            :class="['px-4 py-2 text-sm font-medium text-white rounded-lg transition-colors', confirmState.confirmColor]"
          >
            {{ confirmState.confirmLabel }}
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>
