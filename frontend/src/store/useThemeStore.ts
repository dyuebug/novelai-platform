import { create } from 'zustand'
import { persist } from 'zustand/middleware'

export type ThemeMode = 'light' | 'dark' | 'auto'

interface ThemeState {
  mode: ThemeMode
  primaryColor: string
  setMode: (mode: ThemeMode) => void
  setPrimaryColor: (color: string) => void
  resetTheme: () => void
}

const DEFAULT_PRIMARY_COLOR = '#1890ff'

export const useThemeStore = create<ThemeState>()(
  persist(
    (set) => ({
      mode: 'light',
      primaryColor: DEFAULT_PRIMARY_COLOR,

      setMode: (mode) => set({ mode }),

      setPrimaryColor: (color) => set({ primaryColor: color }),

      resetTheme: () =>
        set({
          mode: 'light',
          primaryColor: DEFAULT_PRIMARY_COLOR,
        }),
    }),
    {
      name: 'theme-storage',
    }
  )
)
