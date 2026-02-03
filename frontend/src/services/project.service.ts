import api from '@/lib/axios'

// 项目类型
export interface Project {
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

// 创建项目请求
export interface CreateProjectRequest {
  title: string
  description?: string
  genre?: string
  metadata?: Record<string, unknown>
}

// 更新项目请求
export interface UpdateProjectRequest {
  title?: string
  description?: string
  genre?: string
  status?: string
  cover_url?: string
}

// 项目列表响应
export interface ProjectListResponse {
  items: Project[]
  total: number
  page: number
  page_size: number
  total_pages: number
}

// 列表查询参数
export interface ProjectListParams {
  page?: number
  page_size?: number
  status?: string
  genre?: string
  sort?: string
}

export const projectService = {
  // 获取项目列表
  list: async (params?: ProjectListParams): Promise<ProjectListResponse> => {
    const response = await api.get('/projects', { params })
    return response.data
  },

  // 创建项目
  create: async (data: CreateProjectRequest): Promise<Project> => {
    const response = await api.post('/projects', data)
    return response.data
  },

  // 获取单个项目
  get: async (id: string): Promise<Project> => {
    const response = await api.get(`/projects/${id}`)
    return response.data
  },

  // 更新项目
  update: async (id: string, data: UpdateProjectRequest): Promise<Project> => {
    const response = await api.put(`/projects/${id}`, data)
    return response.data
  },

  // 删除项目
  delete: async (id: string): Promise<void> => {
    await api.delete(`/projects/${id}`)
  },

  // 恢复项目
  restore: async (id: string): Promise<Project> => {
    const response = await api.post(`/projects/${id}/restore`)
    return response.data
  },
}
