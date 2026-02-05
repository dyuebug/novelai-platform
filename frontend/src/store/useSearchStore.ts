import { create } from 'zustand'
import { persist } from 'zustand/middleware'

interface SearchResult {
  type: string
  id: string
  title: string
  snippet: string
  highlight: string
  metadata?: Record<string, unknown>
}

interface SearchState {
  // 搜索关键词
  query: string

  // 搜索结果
  results: SearchResult[]

  // 搜索历史（最多保存 10 条）
  history: string[]

  // 搜索状态
  loading: boolean

  // 分页
  page: number
  pageSize: number
  total: number

  // 筛选条件
  filters: {
    type?: string
    status?: string[]
    tags?: string[]
    dateFrom?: string
    dateTo?: string
    wordCountMin?: number
    wordCountMax?: number
  }

  // Actions
  setQuery: (query: string) => void
  setResults: (results: SearchResult[]) => void
  addToHistory: (query: string) => void
  clearHistory: () => void
  removeFromHistory: (query: string) => void
  setLoading: (loading: boolean) => void
  setPage: (page: number) => void
  setPageSize: (pageSize: number) => void
  setTotal: (total: number) => void
  setFilters: (filters: Partial<SearchState['filters']>) => void
  clearFilters: () => void
  reset: () => void
}

const MAX_HISTORY_SIZE = 10

export const useSearchStore = create<SearchState>()(
  persist(
    (set, get) => ({
      query: '',
      results: [],
      history: [],
      loading: false,
      page: 1,
      pageSize: 20,
      total: 0,
      filters: {},

      setQuery: (query) => set({ query }),

      setResults: (results) => set({ results }),

      addToHistory: (query) => {
        if (!query.trim()) return

        const { history } = get()

        // 移除重复项
        const newHistory = [
          query,
          ...history.filter((item) => item !== query),
        ].slice(0, MAX_HISTORY_SIZE)

        set({ history: newHistory })
      },

      clearHistory: () => set({ history: [] }),

      removeFromHistory: (query) =>
        set((state) => ({
          history: state.history.filter((item) => item !== query),
        })),

      setLoading: (loading) => set({ loading }),

      setPage: (page) => set({ page }),

      setPageSize: (pageSize) => set({ pageSize }),

      setTotal: (total) => set({ total }),

      setFilters: (filters) =>
        set((state) => ({
          filters: {
            ...state.filters,
            ...filters,
          },
        })),

      clearFilters: () => set({ filters: {} }),

      reset: () =>
        set({
          query: '',
          results: [],
          loading: false,
          page: 1,
          pageSize: 20,
          total: 0,
          filters: {},
        }),
    }),
    {
      name: 'search-storage',
      // 只持久化历史记录和筛选条件
      partialize: (state) => ({
        history: state.history,
        filters: state.filters,
      }),
    }
  )
)
