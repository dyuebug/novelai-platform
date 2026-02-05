import { useEffect, useMemo } from 'react'
import { ConfigProvider, theme as antdTheme } from 'antd'
import { useThemeStore } from '@/store/useThemeStore'
import type { FC, ReactNode } from 'react'

interface ThemeProviderProps {
  children: ReactNode
}

const lightTheme = {
  colorPrimary: '#1890ff',
  colorBgBase: '#ffffff',
  colorTextBase: '#000000',
  colorBorder: '#d9d9d9',
  borderRadius: 6,
}

const darkTheme = {
  colorPrimary: '#1890ff',
  colorBgBase: '#141414',
  colorTextBase: '#ffffff',
  colorBorder: '#434343',
  borderRadius: 6,
}

export const ThemeProvider: FC<ThemeProviderProps> = ({ children }) => {
  const { mode, primaryColor } = useThemeStore()

  // 检测系统主题
  const systemTheme = useMemo(() => {
    if (typeof window === 'undefined') return 'light'
    return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
  }, [])

  // 监听系统主题变化
  useEffect(() => {
    if (mode !== 'auto') return

    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
    const handleChange = () => {
      applyTheme(mediaQuery.matches ? 'dark' : 'light')
    }

    mediaQuery.addEventListener('change', handleChange)
    return () => mediaQuery.removeEventListener('change', handleChange)
  }, [mode])

  // 应用主题
  const applyTheme = (themeMode: 'light' | 'dark') => {
    const root = document.documentElement

    if (themeMode === 'dark') {
      root.classList.add('dark')
      root.style.setProperty('--color-bg-base', darkTheme.colorBgBase)
      root.style.setProperty('--color-text-base', darkTheme.colorTextBase)
      root.style.setProperty('--color-border', darkTheme.colorBorder)
    } else {
      root.classList.remove('dark')
      root.style.setProperty('--color-bg-base', lightTheme.colorBgBase)
      root.style.setProperty('--color-text-base', lightTheme.colorTextBase)
      root.style.setProperty('--color-border', lightTheme.colorBorder)
    }

    root.style.setProperty('--color-primary', primaryColor)
  }

  // 确定当前主题
  const currentTheme = mode === 'auto' ? systemTheme : mode

  // 应用主题
  useEffect(() => {
    applyTheme(currentTheme)
  }, [currentTheme, primaryColor])

  // Ant Design 主题配置
  const antdThemeConfig = useMemo(() => {
    const baseTheme = currentTheme === 'dark' ? darkTheme : lightTheme

    return {
      token: {
        ...baseTheme,
        colorPrimary: primaryColor,
      },
      algorithm: currentTheme === 'dark' ? antdTheme.darkAlgorithm : antdTheme.defaultAlgorithm,
    }
  }, [currentTheme, primaryColor])

  return (
    <ConfigProvider theme={antdThemeConfig}>
      <div className="transition-colors duration-200">
        {children}
      </div>
    </ConfigProvider>
  )
}
