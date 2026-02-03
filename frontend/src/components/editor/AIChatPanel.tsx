import { useState, useRef, useCallback } from 'react'
import {
  Drawer,
  Input,
  Button,
  Select,
  Tabs,
  Space,
  message,
  Spin,
  Typography,
} from 'antd'
import {
  SendOutlined,
  StopOutlined,
  RobotOutlined,
  EditOutlined,
  HighlightOutlined,
} from '@ant-design/icons'
import { aiService, SSEDoneEvent } from '@/services/ai.service'

const { TextArea } = Input
const { Text } = Typography

interface AIChatPanelProps {
  chapterId: string
  open: boolean
  onClose: () => void
  onContentGenerated?: (content: string) => void
  selectedText?: string
  onReplaceSelection?: (original: string, replacement: string) => void
}

const modelOptions = [
  { value: 'gpt-4o', label: 'GPT-4o' },
  { value: 'gpt-4o-mini', label: 'GPT-4o Mini' },
  { value: 'claude-3-5-sonnet-20241022', label: 'Claude 3.5 Sonnet' },
  { value: 'gemini-2.0-flash', label: 'Gemini 2.0 Flash' },
]

const providerOptions = [
  { value: 'openai', label: 'OpenAI' },
  { value: 'anthropic', label: 'Anthropic' },
  { value: 'gemini', label: 'Google' },
]

type TabKey = 'generate' | 'rewrite' | 'polish'

const AIChatPanel = ({
  chapterId,
  open,
  onClose,
  onContentGenerated,
  selectedText,
  onReplaceSelection,
}: AIChatPanelProps) => {
  const [activeTab, setActiveTab] = useState<TabKey>('generate')
  const [model, setModel] = useState('gpt-4o')
  const [provider, setProvider] = useState('openai')
  const [instruction, setInstruction] = useState('')
  const [generating, setGenerating] = useState(false)
  const [generatedContent, setGeneratedContent] = useState('')

  const abortControllerRef = useRef<AbortController | null>(null)

  // 处理生成完成
  const handleDone = useCallback(
    (data: SSEDoneEvent) => {
      setGenerating(false)
      if (data.word_count) {
        message.success(`生成完成，共 ${data.word_count} 字`)
      }
    },
    []
  )

  // 处理错误
  const handleError = useCallback((error: string) => {
    setGenerating(false)
    message.error(`生成失败: ${error}`)
  }, [])

  // 停止生成
  const handleStop = useCallback(() => {
    if (abortControllerRef.current) {
      abortControllerRef.current.abort()
      abortControllerRef.current = null
    }
    setGenerating(false)
  }, [])

  // 生成章节内容
  const handleGenerate = useCallback(() => {
    if (!chapterId) return

    setGenerating(true)
    setGeneratedContent('')

    abortControllerRef.current = new AbortController()

    aiService.generateChapter(
      chapterId,
      {
        model,
        provider,
        instruction: instruction || undefined,
      },
      (text) => {
        setGeneratedContent((prev) => prev + text)
      },
      handleDone,
      handleError,
      abortControllerRef.current.signal
    )
  }, [chapterId, model, provider, instruction, handleDone, handleError])

  // 局部重写
  const handleRewrite = useCallback(() => {
    if (!chapterId || !selectedText) {
      message.warning('请先选择要重写的文本')
      return
    }

    setGenerating(true)
    setGeneratedContent('')

    abortControllerRef.current = new AbortController()

    aiService.partialRegenerate(
      chapterId,
      {
        selection: selectedText,
        instruction: instruction || undefined,
        model,
        provider,
      },
      (text) => {
        setGeneratedContent((prev) => prev + text)
      },
      (data) => {
        setGenerating(false)
        if (data.replacement && onReplaceSelection) {
          // 提供替换选项
          message.success('重写完成，点击"应用"替换原文')
        }
      },
      handleError,
      abortControllerRef.current.signal
    )
  }, [chapterId, selectedText, model, provider, instruction, handleError, onReplaceSelection])

  // 润色
  const handlePolish = useCallback(() => {
    if (!chapterId) return

    setGenerating(true)
    setGeneratedContent('')

    abortControllerRef.current = new AbortController()

    aiService.polish(
      chapterId,
      {
        model,
        provider,
        instruction: instruction || undefined,
      },
      (text) => {
        setGeneratedContent((prev) => prev + text)
      },
      handleDone,
      handleError,
      abortControllerRef.current.signal
    )
  }, [chapterId, model, provider, instruction, handleDone, handleError])

  // 应用生成的内容
  const handleApply = useCallback(() => {
    if (activeTab === 'rewrite' && selectedText && onReplaceSelection) {
      onReplaceSelection(selectedText, generatedContent)
      message.success('已替换选中内容')
    } else if (onContentGenerated) {
      onContentGenerated(generatedContent)
      message.success('已应用生成内容')
    }
    setGeneratedContent('')
    setInstruction('')
  }, [activeTab, selectedText, generatedContent, onContentGenerated, onReplaceSelection])

  const tabItems = [
    {
      key: 'generate',
      label: (
        <span>
          <RobotOutlined />
          生成
        </span>
      ),
      children: (
        <div className="space-y-4">
          <TextArea
            value={instruction}
            onChange={(e) => setInstruction(e.target.value)}
            placeholder="输入写作指令，例如：写一段主角觉醒的场景..."
            rows={3}
          />
          <Button
            type="primary"
            icon={generating ? <StopOutlined /> : <SendOutlined />}
            onClick={generating ? handleStop : handleGenerate}
            block
          >
            {generating ? '停止生成' : '开始生成'}
          </Button>
        </div>
      ),
    },
    {
      key: 'rewrite',
      label: (
        <span>
          <EditOutlined />
          重写
        </span>
      ),
      children: (
        <div className="space-y-4">
          {selectedText ? (
            <div className="p-3 bg-gray-50 rounded text-sm max-h-24 overflow-auto">
              <Text type="secondary">选中的文本：</Text>
              <div className="mt-1">{selectedText}</div>
            </div>
          ) : (
            <div className="p-3 bg-yellow-50 rounded text-sm">
              <Text type="warning">请在编辑器中选择要重写的文本</Text>
            </div>
          )}
          <TextArea
            value={instruction}
            onChange={(e) => setInstruction(e.target.value)}
            placeholder="输入重写指令，例如：让这段更有张力..."
            rows={2}
          />
          <Button
            type="primary"
            icon={generating ? <StopOutlined /> : <SendOutlined />}
            onClick={generating ? handleStop : handleRewrite}
            disabled={!selectedText}
            block
          >
            {generating ? '停止生成' : '重写选中内容'}
          </Button>
        </div>
      ),
    },
    {
      key: 'polish',
      label: (
        <span>
          <HighlightOutlined />
          润色
        </span>
      ),
      children: (
        <div className="space-y-4">
          <TextArea
            value={instruction}
            onChange={(e) => setInstruction(e.target.value)}
            placeholder="输入润色要求，例如：增强对话的情感表达..."
            rows={2}
          />
          <Button
            type="primary"
            icon={generating ? <StopOutlined /> : <SendOutlined />}
            onClick={generating ? handleStop : handlePolish}
            block
          >
            {generating ? '停止生成' : '开始润色'}
          </Button>
        </div>
      ),
    },
  ]

  return (
    <Drawer
      title={
        <span>
          <RobotOutlined className="mr-2" />
          AI 助手
        </span>
      }
      placement="right"
      width={450}
      open={open}
      onClose={onClose}
      extra={
        <Space>
          <Select
            value={provider}
            onChange={setProvider}
            options={providerOptions}
            style={{ width: 100 }}
            size="small"
          />
          <Select
            value={model}
            onChange={setModel}
            options={modelOptions}
            style={{ width: 140 }}
            size="small"
          />
        </Space>
      }
    >
      <div className="flex flex-col h-full">
        <Tabs
          activeKey={activeTab}
          onChange={(key) => setActiveTab(key as TabKey)}
          items={tabItems}
        />

        {/* 生成结果区域 */}
        <div className="flex-1 mt-4 overflow-hidden">
          <div className="text-sm text-gray-500 mb-2">生成结果：</div>
          <div className="h-64 p-3 bg-gray-50 rounded overflow-auto">
            {generating && generatedContent === '' ? (
              <div className="flex items-center justify-center h-full">
                <Spin tip="正在生成..." />
              </div>
            ) : generatedContent ? (
              <div className="whitespace-pre-wrap text-sm">{generatedContent}</div>
            ) : (
              <div className="flex items-center justify-center h-full text-gray-400">
                生成的内容将显示在这里
              </div>
            )}
          </div>
        </div>

        {/* 应用按钮 */}
        {generatedContent && !generating && (
          <div className="mt-4">
            <Button type="primary" onClick={handleApply} block>
              应用到编辑器
            </Button>
          </div>
        )}
      </div>
    </Drawer>
  )
}

export default AIChatPanel
