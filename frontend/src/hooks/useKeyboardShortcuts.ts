import { useEffect, useCallback } from 'react'

export type ShortcutHandler = () => void

interface Shortcut {
  key: string
  ctrl?: boolean
  shift?: boolean
  alt?: boolean
  meta?: boolean
  handler: ShortcutHandler
  description?: string
}

interface UseKeyboardShortcutsOptions {
  shortcuts: Shortcut[]
  enabled?: boolean
}

const isMac = typeof window !== 'undefined' && /Mac|iPod|iPhone|iPad/.test(navigator.platform)

export const useKeyboardShortcuts = ({ shortcuts, enabled = true }: UseKeyboardShortcutsOptions) => {
  const handleKeyDown = useCallback(
    (event: KeyboardEvent) => {
      if (!enabled) return

      // 忽略在输入框中的快捷键（除了特定的保存等操作）
      const target = event.target as HTMLElement
      const isInput = ['INPUT', 'TEXTAREA'].includes(target.tagName) || target.isContentEditable

      for (const shortcut of shortcuts) {
        const ctrlKey = isMac ? event.metaKey : event.ctrlKey
        const metaKey = isMac ? event.metaKey : event.ctrlKey

        const keyMatch = event.key.toLowerCase() === shortcut.key.toLowerCase()
        const ctrlMatch = shortcut.ctrl ? ctrlKey : !ctrlKey
        const shiftMatch = shortcut.shift ? event.shiftKey : !event.shiftKey
        const altMatch = shortcut.alt ? event.altKey : !event.altKey
        const metaMatch = shortcut.meta ? metaKey : true

        if (keyMatch && ctrlMatch && shiftMatch && altMatch && metaMatch) {
          // 某些快捷键在输入框中也应该工作
          const allowInInput = ['s'].includes(shortcut.key.toLowerCase()) && shortcut.ctrl

          if (!isInput || allowInInput) {
            event.preventDefault()
            shortcut.handler()
            break
          }
        }
      }
    },
    [shortcuts, enabled]
  )

  useEffect(() => {
    window.addEventListener('keydown', handleKeyDown)
    return () => window.removeEventListener('keydown', handleKeyDown)
  }, [handleKeyDown])
}

// 格式化快捷键显示
export const formatShortcut = (shortcut: Shortcut): string => {
  const parts: string[] = []

  if (shortcut.ctrl) parts.push(isMac ? '⌘' : 'Ctrl')
  if (shortcut.shift) parts.push(isMac ? '⇧' : 'Shift')
  if (shortcut.alt) parts.push(isMac ? '⌥' : 'Alt')

  parts.push(shortcut.key.toUpperCase())

  return parts.join(isMac ? '' : '+')
}
