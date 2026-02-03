import { useAuthStore } from '@/store/useAuthStore'

const API_BASE = '/api/v1'

// AI 生成请求
export interface GenerateChapterRequest {
  model?: string
  provider?: string
  instruction?: string
  context?: string
  temperature?: number
  max_tokens?: number
}

// 局部重写请求
export interface PartialRegenerateRequest {
  selection: string
  instruction?: string
  model?: string
  provider?: string
  temperature?: number
}

// 润色请求
export interface PolishRequest {
  model?: string
  provider?: string
  instruction?: string
  temperature?: number
}

// SSE 事件类型
export type SSEEventType = 'content' | 'done' | 'error'

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

// 创建 SSE 连接的通用函数
export function createSSEConnection(
  url: string,
  body: object,
  onContent: (text: string) => void,
  onDone: (data: SSEDoneEvent) => void,
  onError: (error: string) => void,
  signal?: AbortSignal
): void {
  const { accessToken } = useAuthStore.getState()

  fetch(`${API_BASE}${url}`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${accessToken}`,
    },
    body: JSON.stringify(body),
    signal,
  })
    .then((response) => {
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      if (!response.body) {
        throw new Error('No response body')
      }

      const reader = response.body.getReader()
      const decoder = new TextDecoder()
      let buffer = ''

      function processChunk(): Promise<void> {
        return reader.read().then(({ done, value }) => {
          if (done) {
            return
          }

          buffer += decoder.decode(value, { stream: true })
          const lines = buffer.split('\n')
          buffer = lines.pop() || ''

          for (const line of lines) {
            if (line.startsWith('event:')) {
              // 解析事件类型
              continue
            }
            if (line.startsWith('data:')) {
              const data = line.slice(5).trim()
              if (data) {
                try {
                  const parsed = JSON.parse(data)

                  // 根据数据结构判断事件类型
                  if (parsed.text !== undefined) {
                    onContent(parsed.text)
                  } else if (parsed.status === 'completed') {
                    onDone(parsed)
                  } else if (parsed.message !== undefined) {
                    onError(parsed.message)
                  }
                } catch (e) {
                  console.error('Failed to parse SSE data:', e)
                }
              }
            }
          }

          return processChunk()
        })
      }

      return processChunk()
    })
    .catch((error) => {
      if (error.name !== 'AbortError') {
        onError(error.message || 'Connection failed')
      }
    })
}

export const aiService = {
  // 生成章节内容
  generateChapter: (
    chapterId: string,
    request: GenerateChapterRequest,
    onContent: (text: string) => void,
    onDone: (data: SSEDoneEvent) => void,
    onError: (error: string) => void,
    signal?: AbortSignal
  ) => {
    createSSEConnection(
      `/chapters/${chapterId}/generate-stream`,
      request,
      onContent,
      onDone,
      onError,
      signal
    )
  },

  // 局部重写
  partialRegenerate: (
    chapterId: string,
    request: PartialRegenerateRequest,
    onContent: (text: string) => void,
    onDone: (data: SSEDoneEvent) => void,
    onError: (error: string) => void,
    signal?: AbortSignal
  ) => {
    createSSEConnection(
      `/chapters/${chapterId}/partial-regenerate-stream`,
      request,
      onContent,
      onDone,
      onError,
      signal
    )
  },

  // 润色
  polish: (
    chapterId: string,
    request: PolishRequest,
    onContent: (text: string) => void,
    onDone: (data: SSEDoneEvent) => void,
    onError: (error: string) => void,
    signal?: AbortSignal
  ) => {
    createSSEConnection(
      `/chapters/${chapterId}/polish-stream`,
      request,
      onContent,
      onDone,
      onError,
      signal
    )
  },
}
