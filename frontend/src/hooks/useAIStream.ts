import { useState, useCallback, useRef } from 'react'

interface SSEEvent {
  type: 'content' | 'tool_call' | 'tool_result' | 'error' | 'done'
  data: string
}

interface UseAIStreamOptions {
  onChunk?: (chunk: string) => void
  onToolCall?: (tool: string, args: unknown) => void
  onComplete?: (fullContent: string) => void
  onError?: (error: Error) => void
}

export const useAIStream = (options: UseAIStreamOptions = {}) => {
  const [content, setContent] = useState('')
  const [isStreaming, setIsStreaming] = useState(false)
  const [error, setError] = useState<Error | null>(null)
  const abortControllerRef = useRef<AbortController | null>(null)

  const parseSSE = (chunk: string): SSEEvent[] => {
    const events: SSEEvent[] = []
    let currentEvent: Partial<SSEEvent> = {}

    for (const line of chunk.split('\n')) {
      if (line.startsWith('event: ')) {
        currentEvent.type = line.slice(7) as SSEEvent['type']
      } else if (line.startsWith('data: ')) {
        currentEvent.data = line.slice(6)
        if (currentEvent.type && currentEvent.data) {
          events.push(currentEvent as SSEEvent)
          currentEvent = {}
        }
      }
    }
    return events
  }

  const startStream = useCallback(
    async (url: string, body?: object) => {
      // 取消之前的请求
      if (abortControllerRef.current) {
        abortControllerRef.current.abort()
      }

      abortControllerRef.current = new AbortController()
      setContent('')
      setIsStreaming(true)
      setError(null)

      try {
        const token = localStorage.getItem('auth-storage')
        const authData = token ? JSON.parse(token) : null
        const accessToken = authData?.state?.accessToken

        const response = await fetch(url, {
          method: body ? 'POST' : 'GET',
          headers: {
            'Content-Type': 'application/json',
            ...(accessToken && { Authorization: `Bearer ${accessToken}` }),
          },
          body: body ? JSON.stringify(body) : undefined,
          signal: abortControllerRef.current.signal,
        })

        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`)
        }

        const reader = response.body?.getReader()
        const decoder = new TextDecoder()
        let fullContent = ''

        while (reader) {
          const { done, value } = await reader.read()
          if (done) break

          const chunk = decoder.decode(value, { stream: true })
          const events = parseSSE(chunk)

          for (const event of events) {
            switch (event.type) {
              case 'content':
                try {
                  const parsed = JSON.parse(event.data)
                  const text = parsed.content || parsed.text || ''
                  fullContent += text

                  // 节流更新 UI
                  requestAnimationFrame(() => {
                    setContent(fullContent)
                  })

                  options.onChunk?.(text)
                } catch {
                  // 非 JSON 数据，直接追加
                  fullContent += event.data
                  setContent(fullContent)
                }
                break

              case 'tool_call':
                try {
                  const toolData = JSON.parse(event.data)
                  options.onToolCall?.(toolData.tool, toolData.args)
                } catch {
                  // 忽略解析错误
                }
                break

              case 'error':
                throw new Error(event.data)

              case 'done':
                // 完成
                break
            }
          }
        }

        options.onComplete?.(fullContent)
      } catch (err) {
        if ((err as Error).name !== 'AbortError') {
          const error = err as Error
          setError(error)
          options.onError?.(error)
        }
      } finally {
        setIsStreaming(false)
      }
    },
    [options]
  )

  const stopStream = useCallback(() => {
    abortControllerRef.current?.abort()
    setIsStreaming(false)
  }, [])

  return {
    content,
    isStreaming,
    error,
    startStream,
    stopStream,
  }
}
