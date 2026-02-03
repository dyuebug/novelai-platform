import api from '@/lib/axios'

// 角色类型
export interface Character {
  id: string
  project_id: string
  name: string
  alias?: string
  gender?: string
  age?: string
  birthday?: string
  appearance?: string
  personality?: string
  background?: string
  abilities?: string
  goals?: string
  role?: string
  status: string
  avatar_url?: string
  metadata?: Record<string, unknown>
  sort_order: number
  created_at: string
  updated_at: string
}

// 角色关系
export interface CharacterRelationship {
  id: string
  character_id: string
  target_id: string
  relation_type: string
  description?: string
  start_chapter?: number
  end_chapter?: number
  target?: Character
  created_at: string
}

// 角色经历
export interface CharacterExperience {
  id: string
  character_id: string
  title: string
  content?: string
  chapter_ref?: number
  time_point?: string
  sort_order: number
  created_at: string
}

// 地点类型
export interface Location {
  id: string
  project_id: string
  parent_id?: string
  name: string
  description?: string
  type?: string
  climate?: string
  features?: string
  image_url?: string
  metadata?: Record<string, unknown>
  sort_order: number
  children?: Location[]
  created_at: string
  updated_at: string
}

// 组织类型
export interface Organization {
  id: string
  project_id: string
  name: string
  type?: string
  description?: string
  history?: string
  structure?: string
  goals?: string
  location_id?: string
  logo_url?: string
  metadata?: Record<string, unknown>
  sort_order: number
  created_at: string
  updated_at: string
}

// 组织成员
export interface OrganizationMember {
  id: string
  organization_id: string
  character_id: string
  position?: string
  rank?: string
  join_chapter?: number
  leave_chapter?: number
  status: string
  character?: Character
  created_at: string
}

// 世界设定
export interface WorldSetting {
  id: string
  project_id: string
  parent_id?: string
  category: string
  title: string
  content?: string
  metadata?: Record<string, unknown>
  sort_order: number
  children?: WorldSetting[]
  created_at: string
  updated_at: string
}

// 列表响应
export interface ListResponse<T> {
  items: T[]
  total: number
  page: number
  page_size: number
  total_pages: number
}

// 请求类型
export interface CreateCharacterRequest {
  name: string
  alias?: string
  gender?: string
  age?: string
  birthday?: string
  appearance?: string
  personality?: string
  background?: string
  abilities?: string
  goals?: string
  role?: string
  avatar_url?: string
}

export interface CreateLocationRequest {
  parent_id?: string
  name: string
  description?: string
  type?: string
  climate?: string
  features?: string
  image_url?: string
}

export interface CreateOrganizationRequest {
  name: string
  type?: string
  description?: string
  history?: string
  structure?: string
  goals?: string
  location_id?: string
  logo_url?: string
}

export interface CreateWorldSettingRequest {
  parent_id?: string
  category: string
  title: string
  content?: string
}

// 角色服务
export const characterService = {
  list: async (projectId: string, params?: { page?: number; page_size?: number; role?: string }): Promise<ListResponse<Character>> => {
    const response = await api.get(`/projects/${projectId}/characters`, { params })
    return response.data
  },

  create: async (projectId: string, data: CreateCharacterRequest): Promise<Character> => {
    const response = await api.post(`/projects/${projectId}/characters`, data)
    return response.data
  },

  get: async (id: string): Promise<Character> => {
    const response = await api.get(`/characters/${id}`)
    return response.data
  },

  update: async (id: string, data: Partial<CreateCharacterRequest>): Promise<Character> => {
    const response = await api.put(`/characters/${id}`, data)
    return response.data
  },

  delete: async (id: string): Promise<void> => {
    await api.delete(`/characters/${id}`)
  },

  // 关系
  getRelationships: async (characterId: string): Promise<CharacterRelationship[]> => {
    const response = await api.get(`/characters/${characterId}/relationships`)
    return response.data
  },

  createRelationship: async (characterId: string, data: { target_id: string; relation_type: string; description?: string }): Promise<CharacterRelationship> => {
    const response = await api.post(`/characters/${characterId}/relationships`, data)
    return response.data
  },

  deleteRelationship: async (relationshipId: string): Promise<void> => {
    await api.delete(`/relationships/${relationshipId}`)
  },

  // 经历
  getExperiences: async (characterId: string): Promise<CharacterExperience[]> => {
    const response = await api.get(`/characters/${characterId}/experiences`)
    return response.data
  },

  createExperience: async (characterId: string, data: { title: string; content?: string; chapter_ref?: number; time_point?: string }): Promise<CharacterExperience> => {
    const response = await api.post(`/characters/${characterId}/experiences`, data)
    return response.data
  },

  deleteExperience: async (experienceId: string): Promise<void> => {
    await api.delete(`/experiences/${experienceId}`)
  },
}

// 地点服务
export const locationService = {
  list: async (projectId: string, params?: { page?: number; page_size?: number; type?: string }): Promise<ListResponse<Location>> => {
    const response = await api.get(`/projects/${projectId}/locations`, { params })
    return response.data
  },

  getTree: async (projectId: string): Promise<Location[]> => {
    const response = await api.get(`/projects/${projectId}/locations/tree`)
    return response.data
  },

  create: async (projectId: string, data: CreateLocationRequest): Promise<Location> => {
    const response = await api.post(`/projects/${projectId}/locations`, data)
    return response.data
  },

  get: async (id: string): Promise<Location> => {
    const response = await api.get(`/locations/${id}`)
    return response.data
  },

  update: async (id: string, data: Partial<CreateLocationRequest>): Promise<Location> => {
    const response = await api.put(`/locations/${id}`, data)
    return response.data
  },

  delete: async (id: string): Promise<void> => {
    await api.delete(`/locations/${id}`)
  },
}

// 组织服务
export const organizationService = {
  list: async (projectId: string, params?: { page?: number; page_size?: number; type?: string }): Promise<ListResponse<Organization>> => {
    const response = await api.get(`/projects/${projectId}/organizations`, { params })
    return response.data
  },

  create: async (projectId: string, data: CreateOrganizationRequest): Promise<Organization> => {
    const response = await api.post(`/projects/${projectId}/organizations`, data)
    return response.data
  },

  get: async (id: string): Promise<Organization> => {
    const response = await api.get(`/organizations/${id}`)
    return response.data
  },

  update: async (id: string, data: Partial<CreateOrganizationRequest>): Promise<Organization> => {
    const response = await api.put(`/organizations/${id}`, data)
    return response.data
  },

  delete: async (id: string): Promise<void> => {
    await api.delete(`/organizations/${id}`)
  },

  // 成员
  getMembers: async (orgId: string): Promise<OrganizationMember[]> => {
    const response = await api.get(`/organizations/${orgId}/members`)
    return response.data
  },

  addMember: async (orgId: string, data: { character_id: string; position?: string; rank?: string }): Promise<OrganizationMember> => {
    const response = await api.post(`/organizations/${orgId}/members`, data)
    return response.data
  },

  removeMember: async (memberId: string): Promise<void> => {
    await api.delete(`/members/${memberId}`)
  },
}

// 世界设定服务
export const worldSettingService = {
  list: async (projectId: string, params?: { page?: number; page_size?: number; category?: string }): Promise<ListResponse<WorldSetting>> => {
    const response = await api.get(`/projects/${projectId}/world-settings`, { params })
    return response.data
  },

  getTree: async (projectId: string, category?: string): Promise<WorldSetting[]> => {
    const response = await api.get(`/projects/${projectId}/world-settings/tree`, { params: { category } })
    return response.data
  },

  getByCategory: async (projectId: string, category: string): Promise<WorldSetting[]> => {
    const response = await api.get(`/projects/${projectId}/world-settings/category/${category}`)
    return response.data
  },

  create: async (projectId: string, data: CreateWorldSettingRequest): Promise<WorldSetting> => {
    const response = await api.post(`/projects/${projectId}/world-settings`, data)
    return response.data
  },

  get: async (id: string): Promise<WorldSetting> => {
    const response = await api.get(`/world-settings/${id}`)
    return response.data
  },

  update: async (id: string, data: Partial<CreateWorldSettingRequest>): Promise<WorldSetting> => {
    const response = await api.put(`/world-settings/${id}`, data)
    return response.data
  },

  delete: async (id: string): Promise<void> => {
    await api.delete(`/world-settings/${id}`)
  },
}
