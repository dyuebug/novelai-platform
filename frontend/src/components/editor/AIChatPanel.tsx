import { useState } from 'react'
import { Drawer, Tabs, Form, Select, Input, Button, Card, Space, message } from 'antd'
import { SendOutlined, StopOutlined, CopyOutlined, CheckOutlined } from '@ant-design/icons'
import type { FC } from 'react'

const { TextArea } = Input

interface AIChatPanelProps {
  visible: boolean
  onClose: () => void
  onInsert?: (text: string) => void
}

export const AIChatPanel: FC<AIChatPanelProps> = ({
  visible,
  onClose,
  onInsert,
}) => {
  const [activeTab, setActiveTab] = useState('generate')
  const [form] = Form.useForm()
  const [generating, setGenerating] = useState(false)
  const [generatedContent, setGeneratedContent] = useState('')

  const handleGenerate = async () => {
    try {
      await form.validateFields()
      setGenerating(true)
      setGeneratedContent('')

      // TODO: 实现SSE流式生成
      // const eventSource = new EventSource(`/api/ai/generate?prompt=${values.prompt}`)
      // eventSource.onmessage = (event) => {
      //   setGeneratedContent(prev => prev + event.data)
      // }

      // 模拟流式生成
      const mockContent = '这是AI生成的内容...'
      for (let i = 0; i < mockContent.length; i++) {
        await new Promise(resolve => setTimeout(resolve, 50))
        setGeneratedContent(mockContent.slice(0, i + 1))
      }
      setGenerating(false)
    } catch (error) {
      message.error('生成失败')
      setGenerating(false)
    }
  }

  const handleStop = () => {
    setGenerating(false)
    message.info('已停止生成')
  }

  const handleInsert = () => {
    if (generatedContent) {
      onInsert?.(generatedContent)
      onClose()
    }
  }

  const handleCopy = () => {
    navigator.clipboard.writeText(generatedContent)
    message.success('已复制到剪贴板')
  }

  const tabItems = [
    {
      key: 'generate',
      label: '生成',
      children: (
        <div className="space-y-4">
          <Form form={form} layout="vertical">
            <Form.Item label="AI提供商" name="provider" initialValue="openai">
              <Select
                options={[
                  { label: 'OpenAI', value: 'openai' },
                  { label: 'Anthropic', value: 'anthropic' },
                  { label: 'Gemini', value: 'gemini' },
                ]}
              />
            </Form.Item>

            <Form.Item label="模型" name="model" initialValue="gpt-4o">
              <Select
                options={[
                  { label: 'GPT-4o', value: 'gpt-4o' },
                  { label: 'Claude 3 Opus', value: 'claude-3-opus' },
                  { label: 'Gemini Pro', value: 'gemini-pro' },
                ]}
              />
            </Form.Item>

            <Form.Item
              label="提示词"
              name="prompt"
              rules={[{ required: true, message: '请输入提示词' }]}
            >
              <TextArea rows={6} placeholder="描述您想要生成的内容..." />
            </Form.Item>
          </Form>

          <Space className="w-full" direction="vertical">
            {!generating ? (
              <Button
                type="primary"
                icon={<SendOutlined />}
                onClick={handleGenerate}
                block
              >
                生成
              </Button>
            ) : (
              <Button danger icon={<StopOutlined />} onClick={handleStop} block>
                停止生成
              </Button>
            )}
          </Space>

          {generatedContent && (
            <Card
              title="生成结果"
              extra={
                <Space>
                  <Button size="small" icon={<CopyOutlined />} onClick={handleCopy}>
                    复制
                  </Button>
                  <Button
                    size="small"
                    type="primary"
                    icon={<CheckOutlined />}
                    onClick={handleInsert}
                  >
                    插入
                  </Button>
                </Space>
              }
            >
              <div className="whitespace-pre-wrap max-h-96 overflow-auto">
                {generatedContent}
              </div>
            </Card>
          )}
        </div>
      ),
    },
    {
      key: 'rewrite',
      label: '重写',
      children: (
        <div className="text-center text-gray-400 py-8">
          重写功能开发中...
        </div>
      ),
    },
    {
      key: 'polish',
      label: '润色',
      children: (
        <div className="text-center text-gray-400 py-8">
          润色功能开发中...
        </div>
      ),
    },
  ]

  return (
    <Drawer
      title="AI助手"
      placement="right"
      onClose={onClose}
      open={visible}
      width={480}
    >
      <Tabs activeKey={activeTab} onChange={setActiveTab} items={tabItems} />
    </Drawer>
  )
}

export default AIChatPanel
