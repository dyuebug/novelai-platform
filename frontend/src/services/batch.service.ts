import api from '@/lib/axios'

// 批量操作结果
export interface BatchOperationResult {
  success: number
  failed: number
  total: number
  failed_ids?: string[]
  error_message?: string
}

// 批量更新请求
export interface BatchUpdateRequest {
  chapter_ids: string[]
  updates: Record<string, unknown>
}

// 批量删除请求
export interface BatchDeleteRequest {
  chapter_ids: string[]
}

// 批量状态更新请求
export interface BatchStatusUpdateRequest {
  chapter_ids: string[]
  status: string
}

// 重排序请求
export interface ReorderRequest {
  chapter_orders: Array<{
    chapter_id: string
    chapter_number: number
  }>
}

// 批量操作服务
export const batchService = {
  // 批量更新章节
  batchUpdate: async (
    projectId: string,
    data: BatchUpdateRequest
  ): Promise<BatchOperationResult> => {
    const response = await api.post(
      `/projects/${projectId}/chapters/batch-update`,
      data
    )
    return response.data.data
  },

  // 批量删除章节
  batchDelete: async (
    projectId: string,
    data: BatchDeleteRequest
  ): Promise<BatchOperationResult> => {
    const response = await api.post(
      `/projects/${projectId}/chapters/batch-delete`,
      data
    )
    return response.data.data
  },

  // 批量更新章节状态
  batchStatusUpdate: async (
    projectId: string,
    data: BatchStatusUpdateRequest
  ): Promise<BatchOperationResult> => {
    const response = await api.post(
      `/projects/${projectId}/chapters/batch-status`,
      data
    )
    return response.data.data
  },

  // 重排序章节
  reorder: async (projectId: string, data: ReorderRequest): Promise<void> => {
    await api.post(`/projects/${projectId}/chapters/reorder`, data)
  },
}
