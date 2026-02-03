import api from '@/lib/axios'

// 钩子分析
export interface HookAnalysis {
  hook_type: string
  hook_strength: string
  hook_content: string
  score: number
  suggestions: string[]
}

// 爽点项
export interface CoolpointItem {
  type: string
  content: string
  position: number
  intensity: number
}

// 微兑现项
export interface MicropayoffItem {
  goal: string
  payoff: string
  position: number
  satisfaction: number
}

// 追读力分析结果
export interface ReadingPowerResult {
  overall_score: number
  hook: HookAnalysis | null
  coolpoints: CoolpointItem[]
  micropayoffs: MicropayoffItem[]
  dimension_scores: Record<string, number>
  suggestions: string[]
}

// 一致性问题
export interface ConsistencyIssue {
  type: string
  severity: string
  description: string
  chapter_refs: number[]
  suggestion: string
}

// 一致性检查结果
export interface ConsistencyResult {
  overall_score: number
  issues: ConsistencyIssue[]
  character_score: number
  timeline_score: number
  worldview_score: number
}

// Agent审查结果
export interface AgentReview {
  agent_name: string
  role: string
  score: number
  strengths: string[]
  weaknesses: string[]
  suggestions: string[]
}

// 多Agent审查结果
export interface MultiAgentResult {
  overall_score: number
  reviews: AgentReview[]
  consensus_strengths: string[]
  consensus_weaknesses: string[]
  final_suggestions: string[]
}

// 综合质量评估结果
export interface QualityEvaluationResult {
  overall_score: number
  reading_power: ReadingPowerResult
  consistency: ConsistencyResult
  multi_agent: MultiAgentResult
  summary: {
    strengths: string[]
    weaknesses: string[]
    suggestions: string[]
  }
}

// 质量评估请求
export interface QualityRequest {
  content: string
  characters?: Array<{ name: string; personality?: string; background?: string }>
  world_settings?: Array<{ title: string; content?: string }>
  previous_summary?: string
  provider?: string
  model?: string
}

// 质量评估服务
export const qualityService = {
  // 追读力分析
  analyzeReadingPower: async (
    chapterId: string,
    data: QualityRequest
  ): Promise<ReadingPowerResult> => {
    const response = await api.post(`/chapters/${chapterId}/analyze-reading-power`, data)
    return response.data
  },

  // 一致性检查
  checkConsistency: async (
    chapterId: string,
    data: QualityRequest
  ): Promise<ConsistencyResult> => {
    const response = await api.post(`/chapters/${chapterId}/check-consistency`, data)
    return response.data
  },

  // 多Agent审查
  multiAgentReview: async (
    chapterId: string,
    data: QualityRequest
  ): Promise<MultiAgentResult> => {
    const response = await api.post(`/chapters/${chapterId}/multi-agent-review`, data)
    return response.data
  },

  // 综合质量评估
  evaluateQuality: async (
    chapterId: string,
    data: QualityRequest
  ): Promise<QualityEvaluationResult> => {
    const response = await api.post(`/chapters/${chapterId}/evaluate-quality`, data)
    return response.data
  },
}
