import api from '@/lib/axios';

// 伏笔状态
export type ForeshadowStatus = 'planted' | 'hinted' | 'resolved' | 'abandoned';

// 伏笔优先级
export type ForeshadowPriority = 'high' | 'medium' | 'low';

// 伏笔
export interface Foreshadow {
  id: string;
  project_id: string;
  title: string;
  description: string;
  content: string;
  priority: ForeshadowPriority;
  status: ForeshadowStatus;
  plant_chapter_id?: string;
  plant_chapter_num: number;
  plant_position: string;
  resolve_chapter_id?: string;
  resolve_chapter_num: number;
  resolve_content: string;
  resolved_at?: string;
  remind_chapter_num: number;
  remind_enabled: boolean;
  related_character_ids: string[];
  tags: string[];
  sort_order: number;
  created_at: string;
  updated_at: string;
}

// 伏笔暗示
export interface ForeshadowHint {
  id: string;
  foreshadow_id: string;
  chapter_id: string;
  chapter_num: number;
  content: string;
  hint_type: string;
  created_at: string;
}

// 伏笔提醒
export interface ForeshadowReminder {
  id: string;
  foreshadow_id: string;
  chapter_num: number;
  message: string;
  is_read: boolean;
  created_at: string;
  foreshadow?: Foreshadow;
}

// 伏笔统计
export interface ForeshadowStats {
  total: number;
  planted: number;
  hinted: number;
  resolved: number;
  abandoned: number;
  high_priority: number;
  medium_priority: number;
  low_priority: number;
  overdue_count: number;
}

// 创建伏笔请求
export interface CreateForeshadowRequest {
  title: string;
  description?: string;
  content?: string;
  priority?: ForeshadowPriority;
  plant_chapter_id?: string;
  plant_chapter_num?: number;
  plant_position?: string;
  remind_chapter_num?: number;
  related_character_ids?: string[];
  tags?: string[];
}

// 更新伏笔请求
export interface UpdateForeshadowRequest {
  title?: string;
  description?: string;
  content?: string;
  priority?: ForeshadowPriority;
  status?: ForeshadowStatus;
  plant_chapter_id?: string;
  plant_chapter_num?: number;
  plant_position?: string;
  resolve_chapter_id?: string;
  resolve_chapter_num?: number;
  resolve_content?: string;
  remind_chapter_num?: number;
  remind_enabled?: boolean;
  related_character_ids?: string[];
  tags?: string[];
}

// 回收伏笔请求
export interface ResolveForeshadowRequest {
  chapter_id?: string;
  chapter_num: number;
  content: string;
}

// 添加暗示请求
export interface AddHintRequest {
  chapter_id: string;
  chapter_num?: number;
  content: string;
  hint_type?: string;
}

// 创建提醒请求
export interface CreateReminderRequest {
  chapter_num: number;
  message: string;
}

// 列表响应
export interface ForeshadowListResponse {
  items: Foreshadow[];
  total: number;
  page: number;
  page_size: number;
}

// 伏笔服务
export const foreshadowService = {
  // 获取伏笔列表
  async list(
    projectId: string,
    params?: {
      status?: ForeshadowStatus;
      priority?: ForeshadowPriority;
      page?: number;
      page_size?: number;
    }
  ): Promise<ForeshadowListResponse> {
    const response = await api.get(`/projects/${projectId}/foreshadows`, { params });
    return response.data;
  },

  // 创建伏笔
  async create(projectId: string, data: CreateForeshadowRequest): Promise<Foreshadow> {
    const response = await api.post(`/projects/${projectId}/foreshadows`, data);
    return response.data;
  },

  // 获取伏笔详情
  async get(id: string): Promise<Foreshadow> {
    const response = await api.get(`/foreshadows/${id}`);
    return response.data;
  },

  // 更新伏笔
  async update(id: string, data: UpdateForeshadowRequest): Promise<Foreshadow> {
    const response = await api.put(`/foreshadows/${id}`, data);
    return response.data;
  },

  // 删除伏笔
  async delete(id: string): Promise<void> {
    await api.delete(`/foreshadows/${id}`);
  },

  // 回收伏笔
  async resolve(id: string, data: ResolveForeshadowRequest): Promise<Foreshadow> {
    const response = await api.post(`/foreshadows/${id}/resolve`, data);
    return response.data;
  },

  // 获取伏笔统计
  async getStats(projectId: string, currentChapter?: number): Promise<ForeshadowStats> {
    const response = await api.get(`/projects/${projectId}/foreshadows/stats`, {
      params: { current_chapter: currentChapter },
    });
    return response.data;
  },

  // 获取待提醒伏笔
  async getPending(projectId: string, currentChapter: number): Promise<Foreshadow[]> {
    const response = await api.get(`/projects/${projectId}/foreshadows/pending`, {
      params: { current_chapter: currentChapter },
    });
    return response.data;
  },

  // 获取超期伏笔
  async getOverdue(projectId: string, currentChapter: number): Promise<Foreshadow[]> {
    const response = await api.get(`/projects/${projectId}/foreshadows/overdue`, {
      params: { current_chapter: currentChapter },
    });
    return response.data;
  },

  // 检查并创建提醒
  async checkReminders(
    projectId: string,
    currentChapter: number
  ): Promise<{ created_reminders: ForeshadowReminder[]; count: number }> {
    const response = await api.post(`/projects/${projectId}/foreshadows/check-reminders`, null, {
      params: { current_chapter: currentChapter },
    });
    return response.data;
  },

  // 获取伏笔暗示列表
  async getHints(foreshadowId: string): Promise<ForeshadowHint[]> {
    const response = await api.get(`/foreshadows/${foreshadowId}/hints`);
    return response.data;
  },

  // 添加暗示
  async addHint(foreshadowId: string, data: AddHintRequest): Promise<ForeshadowHint> {
    const response = await api.post(`/foreshadows/${foreshadowId}/hints`, data);
    return response.data;
  },

  // 删除暗示
  async deleteHint(hintId: string): Promise<void> {
    await api.delete(`/hints/${hintId}`);
  },

  // 创建提醒
  async createReminder(foreshadowId: string, data: CreateReminderRequest): Promise<ForeshadowReminder> {
    const response = await api.post(`/foreshadows/${foreshadowId}/reminders`, data);
    return response.data;
  },

  // 获取未读提醒
  async getUnreadReminders(projectId: string): Promise<ForeshadowReminder[]> {
    const response = await api.get(`/projects/${projectId}/foreshadow-reminders/unread`);
    return response.data;
  },

  // 标记提醒为已读
  async markReminderAsRead(reminderId: string): Promise<void> {
    await api.post(`/reminders/${reminderId}/read`);
  },

  // 标记所有提醒为已读
  async markAllRemindersAsRead(projectId: string): Promise<void> {
    await api.post(`/projects/${projectId}/foreshadow-reminders/mark-all-read`);
  },
};

export default foreshadowService;
