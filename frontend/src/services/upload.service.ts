import api from '@/lib/axios'

// 文件类型
export interface File {
  id: string
  user_id: string
  file_name: string
  storage_path: string
  file_size: number
  mime_type: string
  file_type: string
  width?: number
  height?: number
  thumbnail_url?: string
  hash: string
  metadata: Record<string, unknown>
  created_at: string
  updated_at: string
}

// 文件列表响应
export interface FileListResponse {
  items: File[]
  total: number
  page: number
  page_size: number
  total_pages: number
}

// 存储统计
export interface StorageStats {
  file_count: number
  total_size: number
}

// 文件上传服务
export const uploadService = {
  // 上传文件
  upload: async (file: globalThis.File, fileType: string = 'image'): Promise<File> => {
    const formData = new FormData()
    formData.append('file', file)
    formData.append('file_type', fileType)

    const response = await api.post('/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    })
    return response.data.data
  },

  // 获取文件列表
  list: async (
    page: number = 1,
    pageSize: number = 20,
    fileType?: string,
    sort?: string
  ): Promise<FileListResponse> => {
    const response = await api.get('/files', {
      params: { page, page_size: pageSize, file_type: fileType, sort },
    })
    return response.data.data
  },

  // 获取文件元数据
  getFile: async (id: string): Promise<File> => {
    const response = await api.get(`/files/${id}`)
    return response.data.data
  },

  // 删除文件
  deleteFile: async (id: string): Promise<void> => {
    await api.delete(`/files/${id}`)
  },

  // 获取存储统计
  getStats: async (): Promise<StorageStats> => {
    const response = await api.get('/files/stats')
    return response.data.data
  },

  // 获取文件下载 URL
  getDownloadUrl: (id: string): string => {
    return `${api.defaults.baseURL}/files/${id}/download`
  },

  // 获取缩略图 URL
  getThumbnailUrl: (id: string): string => {
    return `${api.defaults.baseURL}/files/${id}/thumbnail`
  },
}
