import { create } from 'zustand'
import { persist } from 'zustand/middleware'

interface PanelState {
  visible: boolean
  width?: number
  height?: number
  collapsed?: boolean
}

interface EditorLayoutState {
  // 侧边栏状态
  sidebarCollapsed: boolean
  sidebarWidth: number

  // AI 抽屉状态
  aiDrawerOpen: boolean
  aiDrawerWidth: number

  // 工具栏状态
  toolbarVisible: boolean

  // 章节列表状态
  chapterListVisible: boolean

  // 沉浸模式
  immersiveMode: boolean

  // 面板状态
  panels: Record<string, PanelState>

  // Actions
  toggleSidebar: () => void
  setSidebarWidth: (width: number) => void
  toggleAiDrawer: () => void
  setAiDrawerWidth: (width: number) => void
  toggleToolbar: () => void
  toggleChapterList: () => void
  toggleImmersiveMode: () => void
  setPanelState: (panelId: string, state: Partial<PanelState>) => void
  resetLayout: () => void
}

const DEFAULT_SIDEBAR_WIDTH = 220
const DEFAULT_AI_DRAWER_WIDTH = 400
const MIN_SIDEBAR_WIDTH = 80
const MAX_SIDEBAR_WIDTH = 400
const MIN_AI_DRAWER_WIDTH = 280
const MAX_AI_DRAWER_WIDTH = 600

export const useEditorLayoutStore = create<EditorLayoutState>()(
  persist(
    (set) => ({
      sidebarCollapsed: false,
      sidebarWidth: DEFAULT_SIDEBAR_WIDTH,
      aiDrawerOpen: false,
      aiDrawerWidth: DEFAULT_AI_DRAWER_WIDTH,
      toolbarVisible: true,
      chapterListVisible: true,
      immersiveMode: false,
      panels: {},

      toggleSidebar: () =>
        set((state) => ({
          sidebarCollapsed: !state.sidebarCollapsed,
        })),

      setSidebarWidth: (width) =>
        set({
          sidebarWidth: Math.max(
            MIN_SIDEBAR_WIDTH,
            Math.min(MAX_SIDEBAR_WIDTH, width)
          ),
        }),

      toggleAiDrawer: () =>
        set((state) => ({
          aiDrawerOpen: !state.aiDrawerOpen,
        })),

      setAiDrawerWidth: (width) =>
        set({
          aiDrawerWidth: Math.max(
            MIN_AI_DRAWER_WIDTH,
            Math.min(MAX_AI_DRAWER_WIDTH, width)
          ),
        }),

      toggleToolbar: () =>
        set((state) => ({
          toolbarVisible: !state.toolbarVisible,
        })),

      toggleChapterList: () =>
        set((state) => ({
          chapterListVisible: !state.chapterListVisible,
        })),

      toggleImmersiveMode: () =>
        set((state) => {
          const newImmersiveMode = !state.immersiveMode
          return {
            immersiveMode: newImmersiveMode,
            // 进入沉浸模式时隐藏所有 UI
            toolbarVisible: !newImmersiveMode,
            sidebarCollapsed: newImmersiveMode,
            aiDrawerOpen: false,
          }
        }),

      setPanelState: (panelId, state) =>
        set((current) => ({
          panels: {
            ...current.panels,
            [panelId]: {
              ...current.panels[panelId],
              ...state,
            },
          },
        })),

      resetLayout: () =>
        set({
          sidebarCollapsed: false,
          sidebarWidth: DEFAULT_SIDEBAR_WIDTH,
          aiDrawerOpen: false,
          aiDrawerWidth: DEFAULT_AI_DRAWER_WIDTH,
          toolbarVisible: true,
          chapterListVisible: true,
          immersiveMode: false,
          panels: {},
        }),
    }),
    {
      name: 'editor-layout-storage',
    }
  )
)
