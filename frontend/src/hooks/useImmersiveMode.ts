import { useState, useEffect, useCallback } from 'react'

interface UseImmersiveModeOptions {
  onEnter?: () => void
  onExit?: () => void
  autoEnter?: boolean
}

export const useImmersiveMode = (options: UseImmersiveModeOptions = {}) => {
  const [isImmersive, setIsImmersive] = useState(false)
  const [isFullscreen, setIsFullscreen] = useState(false)

  // 检查是否支持全屏API
  const supportsFullscreen = typeof document !== 'undefined' &&
    (document.fullscreenEnabled ||
     (document as any).webkitFullscreenEnabled ||
     (document as any).mozFullScreenEnabled)

  // 进入沉浸模式
  const enter = useCallback(() => {
    setIsImmersive(true)
    localStorage.setItem('immersive-mode-active', 'true')
    options.onEnter?.()
  }, [options])

  // 退出沉浸模式
  const exit = useCallback(() => {
    setIsImmersive(false)
    localStorage.removeItem('immersive-mode-active')

    // 如果在全屏模式，也退出全屏
    if (isFullscreen) {
      exitFullscreen()
    }

    options.onExit?.()
  }, [isFullscreen, options])

  // 切换沉浸模式
  const toggle = useCallback(() => {
    if (isImmersive) {
      exit()
    } else {
      enter()
    }
  }, [isImmersive, enter, exit])

  // 进入全屏
  const enterFullscreen = useCallback(async () => {
    if (!supportsFullscreen) {
      console.warn('Fullscreen API not supported')
      return
    }

    try {
      const elem = document.documentElement
      if (elem.requestFullscreen) {
        await elem.requestFullscreen()
      } else if ((elem as any).webkitRequestFullscreen) {
        await (elem as any).webkitRequestFullscreen()
      } else if ((elem as any).mozRequestFullScreen) {
        await (elem as any).mozRequestFullScreen()
      }
      setIsFullscreen(true)
    } catch (error) {
      console.error('Failed to enter fullscreen:', error)
    }
  }, [supportsFullscreen])

  // 退出全屏
  const exitFullscreen = useCallback(async () => {
    try {
      if (document.exitFullscreen) {
        await document.exitFullscreen()
      } else if ((document as any).webkitExitFullscreen) {
        await (document as any).webkitExitFullscreen()
      } else if ((document as any).mozCancelFullScreen) {
        await (document as any).mozCancelFullScreen()
      }
      setIsFullscreen(false)
    } catch (error) {
      console.error('Failed to exit fullscreen:', error)
    }
  }, [])

  // 监听全屏变化
  useEffect(() => {
    const handleFullscreenChange = () => {
      const isCurrentlyFullscreen = !!(
        document.fullscreenElement ||
        (document as any).webkitFullscreenElement ||
        (document as any).mozFullScreenElement
      )
      setIsFullscreen(isCurrentlyFullscreen)
    }

    document.addEventListener('fullscreenchange', handleFullscreenChange)
    document.addEventListener('webkitfullscreenchange', handleFullscreenChange)
    document.addEventListener('mozfullscreenchange', handleFullscreenChange)

    return () => {
      document.removeEventListener('fullscreenchange', handleFullscreenChange)
      document.removeEventListener('webkitfullscreenchange', handleFullscreenChange)
      document.removeEventListener('mozfullscreenchange', handleFullscreenChange)
    }
  }, [])

  // 自动进入沉浸模式
  useEffect(() => {
    if (options.autoEnter) {
      const shouldAutoEnter = localStorage.getItem('immersive-mode-auto-enter') === 'true'
      if (shouldAutoEnter) {
        enter()
      }
    }
  }, [options.autoEnter, enter])

  // 恢复沉浸模式状态
  useEffect(() => {
    const wasImmersive = localStorage.getItem('immersive-mode-active') === 'true'
    if (wasImmersive) {
      setIsImmersive(true)
    }
  }, [])

  return {
    isImmersive,
    isFullscreen,
    enter,
    exit,
    toggle,
    enterFullscreen,
    exitFullscreen,
    supportsFullscreen,
  }
}

// 设置自动进入沉浸模式
export const setAutoEnterImmersiveMode = (enabled: boolean) => {
  if (enabled) {
    localStorage.setItem('immersive-mode-auto-enter', 'true')
  } else {
    localStorage.removeItem('immersive-mode-auto-enter')
  }
}

// 获取自动进入设置
export const getAutoEnterImmersiveMode = (): boolean => {
  return localStorage.getItem('immersive-mode-auto-enter') === 'true'
}
