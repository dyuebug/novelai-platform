import api from '@/lib/axios'

// 搜索结果类型
export interface SearchResult {
  type: string
  id: string
  title: string
  snippet: string
  highlight: string
  metadata?: Record<string, unknown>
}

// 搜索响应
export interface SearchResponse {
  results: SearchResult[]
  total: number
  page: number
  page_size: number
  total_pages: number
}

// 全局搜索请求
export interface GlobalSearchRequest {
  q: string
  page?: number
  page_size?: number
}

// 项目内搜索请求
export interface ProjectSearchRequest {
  q: string
  type?: string // chapter, character, location, all
  page?: number
  page_size?: number
}

// 高级筛选请求
export interface AdvancedFilterRequest {
  tags?: string[]
  status?: string[]
  date_from?: string
  date_to?: string
  word_count_min?: number
  word_count_max?: number
  page?: number
  page_size?: number
}

// 搜索服务
export const searchService = {
  // 全局搜索
  globalSearch: async (params: GlobalSearchRequest): Promise<SearchResponse> => {
    const response = await api.get('/search', { params })
    return response.data.data
  },

  // 项目内搜索
  projectSearch: async (
    projectId: string,
    params: ProjectSearchRequest
  ): Promise<SearchResponse> => {
    const response = await api.get(`/projects/${projectId}/search`, { params })
    return response.data.data
  },

  // 高级筛选
  advancedFilter: async (
    projectId: string,
    data: AdvancedFilterRequest
  ): Promise<SearchResponse> => {
    const response = await api.post(`/projects/${projectId}/advanced-filter`, data)
    return response.data.data
  },
}
