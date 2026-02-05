import { create } from 'zustand'
import { persist } from 'zustand/middleware'

interface Project {
  id: string
  user_id: string
  title: string
  description?: string
  genre?: string
  status: string
  total_chapters: number
  total_words: number
  cover_url?: string
  metadata?: Record<string, unknown>
  is_deleted: boolean
  deleted_at?: string
  created_at: string
  updated_at: string
}

interface ProjectState {
  // 当前项目
  currentProject: Project | null

  // 项目列表
  projects: Project[]

  // 加载状态
  loading: boolean

  // 分页
  page: number
  pageSize: number
  total: number

  // 筛选和排序
  filters: {
    status?: string
    genre?: string
    search?: string
  }
  sort: string

  // 批量选择
  selectedIds: string[]

  // Actions
  setCurrentProject: (project: Project | null) => void
  setProjects: (projects: Project[]) => void
  addProject: (project: Project) => void
  updateProject: (id: string, updates: Partial<Project>) => void
  removeProject: (id: string) => void
  setLoading: (loading: boolean) => void
  setPage: (page: number) => void
  setPageSize: (pageSize: number) => void
  setTotal: (total: number) => void
  setFilters: (filters: Partial<ProjectState['filters']>) => void
  clearFilters: () => void
  setSort: (sort: string) => void
  toggleSelection: (id: string) => void
  selectAll: () => void
  clearSelection: () => void
  reset: () => void
}

export const useProjectStore = create<ProjectState>()(
  persist(
    (set) => ({
      currentProject: null,
      projects: [],
      loading: false,
      page: 1,
      pageSize: 20,
      total: 0,
      filters: {},
      sort: 'updated_at DESC',
      selectedIds: [],

      setCurrentProject: (project) => set({ currentProject: project }),

      setProjects: (projects) => set({ projects }),

      addProject: (project) =>
        set((state) => ({
          projects: [project, ...state.projects],
          total: state.total + 1,
        })),

      updateProject: (id, updates) =>
        set((state) => ({
          projects: state.projects.map((p) =>
            p.id === id ? { ...p, ...updates } : p
          ),
          currentProject:
            state.currentProject?.id === id
              ? { ...state.currentProject, ...updates }
              : state.currentProject,
        })),

      removeProject: (id) =>
        set((state) => ({
          projects: state.projects.filter((p) => p.id !== id),
          total: state.total - 1,
          currentProject:
            state.currentProject?.id === id ? null : state.currentProject,
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

      setSort: (sort) => set({ sort }),

      toggleSelection: (id) =>
        set((state) => ({
          selectedIds: state.selectedIds.includes(id)
            ? state.selectedIds.filter((selectedId) => selectedId !== id)
            : [...state.selectedIds, id],
        })),

      selectAll: () =>
        set((state) => ({
          selectedIds: state.projects.map((p) => p.id),
        })),

      clearSelection: () => set({ selectedIds: [] }),

      reset: () =>
        set({
          currentProject: null,
          projects: [],
          loading: false,
          page: 1,
          pageSize: 20,
          total: 0,
          filters: {},
          sort: 'updated_at DESC',
          selectedIds: [],
        }),
    }),
    {
      name: 'project-storage',
      // 只持久化当前项目和筛选条件
      partialize: (state) => ({
        currentProject: state.currentProject,
        filters: state.filters,
        sort: state.sort,
      }),
    }
  )
)
