import { reactive } from 'vue'

interface ConfirmState {
  visible: boolean
  title: string
  message: string
  confirmLabel: string
  confirmColor: string
  resolve: ((value: boolean) => void) | null
}

export const confirmState = reactive<ConfirmState>({
  visible: false,
  title: '',
  message: '',
  confirmLabel: '削除する',
  confirmColor: 'bg-red-600 hover:bg-red-700',
  resolve: null,
})

export function useConfirm() {
  function open(title: string, message: string, opts?: { label?: string; color?: string }): Promise<boolean> {
    confirmState.title = title
    confirmState.message = message
    confirmState.confirmLabel = opts?.label || '削除する'
    confirmState.confirmColor = opts?.color || 'bg-red-600 hover:bg-red-700'
    confirmState.visible = true
    return new Promise((resolve) => {
      confirmState.resolve = resolve
    })
  }

  function close(result: boolean) {
    confirmState.visible = false
    confirmState.resolve?.(result)
    confirmState.resolve = null
  }

  return { open, close }
}
