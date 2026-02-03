import api from '@/lib/axios'

// 章节类型
export interface Chapter {
  id: string
  project_id: string
  chapter_number: number
  title: string
  content?: string
  summary?: string
  status: string
  word_count: number
  pov_character?: string
  location?: string
  time_setting?: string
  created_at: string
  updated_at: string
  published_at?: string
}

// 创建章节请求
export interface CreateChapterRequest {
  title: string
  content?: string
  summary?: string
  chapter_number?: number
  pov_character?: string
  location?: string
  time_setting?: string
}

// 更新章节请求
export interface UpdateChapterRequest {
  title?: string
  content?: string
  summary?: string
  status?: string
  pov_character?: string
  location?: string
  time_setting?: string
}

// 章节列表响应
export interface ChapterListResponse {
  items: Chapter[]
  total: number
  page: number
  page_size: number
  total_pages: number
}

// 章节版本
export interface ChapterVersion {
  id: string
  chapter_id: string
  version_number: number
  content: string
  word_count: number
  source: string
  ai_provider?: string
  ai_model?: string
  created_at: string
}

// 版本列表响应
export interface VersionListResponse {
  items: ChapterVersion[]
  total: number
  page: number
  page_size: number
  total_pages: number
}

// 版本对比响应
export interface DiffResponse {
  version1: number
  version2: number
  content1: string
  content2: string
  additions: number
  deletions: number
  diff_lines: DiffLine[]
}

export interface DiffLine {
  type: 'add' | 'delete' | 'equal'
  content: string
  line_num?: number
}

// 列表查询参数
export interface ChapterListParams {
  page?: number
  page_size?: number
  status?: string
  sort?: string
}

export const chapterService = {
  // 获取章节列表
  list: async (projectId: string, params?: ChapterListParams): Promise<ChapterListResponse> => {
    const response = await api.get(`/projects/${projectId}/chapters`, { params })
    return response.data
  },

  // 创建章节
  create: async (projectId: string, data: CreateChapterRequest): Promise<Chapter> => {
    const response = await api.post(`/projects/${projectId}/chapters`, data)
    return response.data
  },

  // 获取单个章节
  get: async (id: string): Promise<Chapter> => {
    const response = await api.get(`/chapters/${id}`)
    return response.data
  },

  // 更新章节
  update: async (id: string, data: UpdateChapterRequest): Promise<Chapter> => {
    const response = await api.put(`/chapters/${id}`, data)
    return response.data
  },

  // 删除章节
  delete: async (id: string): Promise<void> => {
    await api.delete(`/chapters/${id}`)
  },

  // 重排序章节
  reorder: async (projectId: string, chapterIds: string[]): Promise<void> => {
    await api.post(`/projects/${projectId}/chapters/reorder`, {
      chapter_ids: chapterIds,
    })
  },

  // 获取版本列表
  listVersions: async (chapterId: string, page = 1, pageSize = 20): Promise<VersionListResponse> => {
    const response = await api.get(`/chapters/${chapterId}/versions`, {
      params: { page, page_size: pageSize },
    })
    return response.data
  },

  // 获取指定版本
  getVersion: async (chapterId: string, versionNumber: number): Promise<ChapterVersion> => {
    const response = await api.get(`/chapters/${chapterId}/versions/${versionNumber}`)
    return response.data
  },

  // 恢复到指定版本
  restoreVersion: async (chapterId: string, versionNumber: number): Promise<Chapter> => {
    const response = await api.post(`/chapters/${chapterId}/versions/${versionNumber}/restore`)
    return response.data
  },

  // 版本对比
  diffVersions: async (chapterId: string, v1: number, v2: number): Promise<DiffResponse> => {
    const response = await api.get(`/chapters/${chapterId}/versions/diff`, {
      params: { v1, v2 },
    })
    return response.data
  },
}
