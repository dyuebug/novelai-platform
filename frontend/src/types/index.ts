// API 响应类型
export interface ApiResponse<T> {
  code: number
  message: string
  data: T
}

export interface PaginatedResponse<T> {
  items: T[]
  total: number
  page: number
  page_size: number
  total_pages: number
}

// 用户相关类型
export interface User {
  id: string
  username: string
  email: string
  role: string
  created_at: string
  updated_at: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  username: string
  email: string
  password: string
}

export interface TokenResponse {
  access_token: string
  refresh_token: string
  token_type: string
  expires_in: number
}

// 项目相关类型
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

export interface CreateProjectRequest {
  title: string
  description?: string
  genre?: string
  metadata?: Record<string, unknown>
}

export interface UpdateProjectRequest {
  title?: string
  description?: string
  genre?: string
  status?: string
  cover_url?: string
}

// 章节相关类型
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

export interface CreateChapterRequest {
  title: string
  content?: string
  summary?: string
  chapter_number?: number
  pov_character?: string
  location?: string
  time_setting?: string
}

export interface UpdateChapterRequest {
  title?: string
  content?: string
  summary?: string
  status?: string
  pov_character?: string
  location?: string
  time_setting?: string
}

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

// 世界观相关类型
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

// 伏笔相关类型
export type ForeshadowStatus = 'planted' | 'hinted' | 'resolved' | 'abandoned'
export type ForeshadowPriority = 'high' | 'medium' | 'low'

export interface Foreshadow {
  id: string
  project_id: string
  title: string
  description: string
  content: string
  priority: ForeshadowPriority
  status: ForeshadowStatus
  plant_chapter_id?: string
  plant_chapter_num: number
  plant_position: string
  resolve_chapter_id?: string
  resolve_chapter_num: number
  resolve_content: string
  resolved_at?: string
  remind_chapter_num: number
  remind_enabled: boolean
  related_character_ids: string[]
  tags: string[]
  sort_order: number
  created_at: string
  updated_at: string
}

export interface ForeshadowHint {
  id: string
  foreshadow_id: string
  chapter_id: string
  chapter_num: number
  content: string
  hint_type: string
  created_at: string
}

export interface ForeshadowStats {
  total: number
  planted: number
  hinted: number
  resolved: number
  abandoned: number
  high_priority: number
  medium_priority: number
  low_priority: number
  overdue_count: number
}

// AI 相关类型
export interface GenerateChapterRequest {
  model?: string
  provider?: string
  instruction?: string
  context?: string
  temperature?: number
  max_tokens?: number
}

export interface PartialRegenerateRequest {
  selection: string
  instruction?: string
  model?: string
  provider?: string
  temperature?: number
}

export interface PolishRequest {
  model?: string
  provider?: string
  instruction?: string
  temperature?: number
}

// SSE 事件类型
export type SSEEventType = 'content' | 'tool_call' | 'tool_result' | 'error' | 'done'

export interface SSEEvent {
  type: SSEEventType
  data: string
}

export interface SSEContentEvent {
  text: string
}

export interface SSEDoneEvent {
  status: string
  word_count?: number
  original?: string
  replacement?: string
  word_count_delta?: number
}

export interface SSEErrorEvent {
  message: string
}

// 质量检查相关类型
export interface ConsistencyCheckResult {
  id: string
  chapter_id: string
  issues: ConsistencyIssue[]
  score: number
  created_at: string
}

export interface ConsistencyIssue {
  type: string
  severity: 'high' | 'medium' | 'low'
  description: string
  location: string
  suggestion?: string
}

export interface ReadingPowerAnalysis {
  id: string
  chapter_id: string
  readability_score: number
  engagement_score: number
  complexity_score: number
  suggestions: string[]
  created_at: string
}

// 通用类型
export interface PageParams {
  page?: number
  page_size?: number
  sort?: string
  order?: 'asc' | 'desc'
}

export interface ErrorResponse {
  code: number
  message: string
  details?: Record<string, unknown>
}
