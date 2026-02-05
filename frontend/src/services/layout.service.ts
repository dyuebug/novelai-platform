import api from '@/lib/axios'

// 布局配置类型
export interface LayoutConfig {
  id: string
  user_id: string
  layout_type: string
  panel_states: Record<string, unknown>
  panel_sizes: Record<string, unknown>
  created_at: string
  updated_at: string
}

// 保存布局配置请求
export interface SaveLayoutConfigRequest {
  panel_states: Record<string, unknown>
  panel_sizes: Record<string, unknown>
}

// 布局配置服务
export const layoutService = {
  // 获取所有布局配置
  getAllConfigs: async (): Promise<LayoutConfig[]> => {
    const response = await api.get('/user/layout-configs')
    return response.data.data
  },

  // 获取特定类型的布局配置
  getConfig: async (layoutType: string): Promise<LayoutConfig> => {
    const response = await api.get(`/user/layout-config/${layoutType}`)
    return response.data.data
  },

  // 保存布局配置
  saveConfig: async (
    layoutType: string,
    data: SaveLayoutConfigRequest
  ): Promise<LayoutConfig> => {
    const response = await api.post(`/user/layout-config/${layoutType}`, data)
    return response.data.data
  },

  // 删除布局配置
  deleteConfig: async (layoutType: string): Promise<void> => {
    await api.delete(`/user/layout-config/${layoutType}`)
  },

  // 重置布局配置
  resetConfig: async (layoutType: string): Promise<LayoutConfig> => {
    const response = await api.post(`/user/layout-config/${layoutType}/reset`)
    return response.data.data
  },
}
