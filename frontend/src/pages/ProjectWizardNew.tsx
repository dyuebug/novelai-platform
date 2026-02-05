import { useState, useEffect } from 'react'
import { Steps, Button, Form, Input, Select, Card, message, Modal } from 'antd'
import { useNavigate } from 'react-router-dom'
import type { FC } from 'react'

const { Step } = Steps
const { TextArea } = Input

interface WizardFormData {
  title: string
  description?: string
  genre?: string
  tags?: string[]
  aiPrompt?: string
  aiProvider?: string
  aiModel?: string
}

const genreOptions = [
  { label: '玄幻', value: 'xuanhuan' },
  { label: '都市', value: 'urban' },
  { label: '科幻', value: 'scifi' },
  { label: '武侠', value: 'wuxia' },
  { label: '仙侠', value: 'xianxia' },
  { label: '历史', value: 'history' },
  { label: '军事', value: 'military' },
  { label: '悬疑', value: 'mystery' },
  { label: '其他', value: 'other' },
]

const STORAGE_KEY = 'project_wizard_draft'

export const ProjectWizardNew: FC = () => {
  const navigate = useNavigate()
  const [currentStep, setCurrentStep] = useState(0)
  const [formData, setFormData] = useState<WizardFormData>({
    title: '',
    aiProvider: 'openai',
    aiModel: 'gpt-4o',
  })
  const [form] = Form.useForm()
  const [generating, setGenerating] = useState(false)
  const [generatedContent, setGeneratedContent] = useState('')

  // 恢复未完成的向导
  useEffect(() => {
    const savedDraft = localStorage.getItem(STORAGE_KEY)
    if (savedDraft) {
      Modal.confirm({
        title: '恢复未完成的项目',
        content: '检测到您有未完成的项目创建，是否继续？',
        okText: '继续',
        cancelText: '重新开始',
        onOk: () => {
          const draft = JSON.parse(savedDraft)
          setFormData(draft.formData)
          setCurrentStep(draft.currentStep)
          form.setFieldsValue(draft.formData)
        },
        onCancel: () => {
          localStorage.removeItem(STORAGE_KEY)
        },
      })
    }
  }, [])

  // 保存草稿
  const saveDraft = () => {
    localStorage.setItem(
      STORAGE_KEY,
      JSON.stringify({
        formData,
        currentStep,
        timestamp: Date.now(),
      })
    )
  }

  const handleNext = async () => {
    try {
      const values = await form.validateFields()
      const newFormData = { ...formData, ...values }
      setFormData(newFormData)
      setCurrentStep(currentStep + 1)
      saveDraft()
    } catch (error) {
      // 验证失败
    }
  }

  const handlePrev = () => {
    setCurrentStep(currentStep - 1)
  }

  const handleCancel = () => {
    Modal.confirm({
      title: '确认取消',
      content: '取消后将丢失所有已填写的信息，确定要取消吗？',
      okText: '确定',
      okType: 'danger',
      cancelText: '继续填写',
      onOk: () => {
        localStorage.removeItem(STORAGE_KEY)
        navigate('/projects')
      },
    })
  }

  const handleGenerate = async () => {
    try {
      await form.validateFields(['aiPrompt'])
      setGenerating(true)
      setGeneratedContent('')

      // TODO: 实现SSE流式生成
      // const values = await form.validateFields(['aiPrompt'])
      // const eventSource = new EventSource(`/api/wizard-stream/generate?prompt=${values.aiPrompt}`)
      // eventSource.onmessage = (event) => {
      //   setGeneratedContent(prev => prev + event.data)
      // }
      // eventSource.onerror = () => {
      //   eventSource.close()
      //   setGenerating(false)
      // }

      // 模拟流式生成
      const mockContent = '这是AI生成的项目设定内容...'
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

  const handleAbortGenerate = () => {
    // TODO: 中断SSE连接
    setGenerating(false)
    message.info('已停止生成')
  }

  const handleFinish = async () => {
    try {
      // TODO: 调用实际API创建项目
      // const project = await projectService.createProject({
      //   ...formData,
      //   generatedContent,
      // })
      message.success('项目创建成功')
      localStorage.removeItem(STORAGE_KEY)
      navigate('/projects')
    } catch (error) {
      message.error('创建项目失败')
    }
  }

  const steps = [
    {
      title: '基本信息',
      content: (
        <Form form={form} layout="vertical" initialValues={formData}>
          <Form.Item
            label="项目标题"
            name="title"
            rules={[
              { required: true, message: '请输入项目标题' },
              { min: 1, max: 200, message: '标题长度为1-200个字符' },
            ]}
          >
            <Input placeholder="请输入项目标题" />
          </Form.Item>

          <Form.Item label="项目描述" name="description">
            <TextArea rows={4} placeholder="请输入项目描述（可选）" maxLength={5000} showCount />
          </Form.Item>

          <Form.Item label="类型" name="genre">
            <Select placeholder="选择项目类型" options={genreOptions} />
          </Form.Item>

          <Form.Item label="标签" name="tags">
            <Select mode="tags" placeholder="添加标签（可选）" />
          </Form.Item>
        </Form>
      ),
    },
    {
      title: 'AI生成设定',
      content: (
        <Form form={form} layout="vertical" initialValues={formData}>
          <Form.Item label="AI提供商" name="aiProvider">
            <Select
              options={[
                { label: 'OpenAI', value: 'openai' },
                { label: 'Anthropic Claude', value: 'anthropic' },
                { label: 'Google Gemini', value: 'gemini' },
              ]}
            />
          </Form.Item>

          <Form.Item label="AI模型" name="aiModel">
            <Select
              options={[
                { label: 'GPT-4o', value: 'gpt-4o' },
                { label: 'GPT-4 Turbo', value: 'gpt-4-turbo' },
                { label: 'Claude 3 Opus', value: 'claude-3-opus' },
                { label: 'Gemini Pro', value: 'gemini-pro' },
              ]}
            />
          </Form.Item>

          <Form.Item
            label="生成提示词"
            name="aiPrompt"
            rules={[{ required: true, message: '请输入生成提示词' }]}
          >
            <TextArea
              rows={6}
              placeholder="描述您的小说世界观、角色、情节等，AI将根据您的描述生成详细设定"
            />
          </Form.Item>

          <div className="space-y-4">
            <Button
              type="primary"
              loading={generating}
              onClick={handleGenerate}
              block
            >
              {generating ? '生成中...' : 'AI生成设定'}
            </Button>

            {generating && (
              <Button danger onClick={handleAbortGenerate} block>
                停止生成
              </Button>
            )}

            {generatedContent && (
              <Card title="生成结果" className="mt-4">
                <div className="whitespace-pre-wrap">{generatedContent}</div>
              </Card>
            )}
          </div>
        </Form>
      ),
    },
  ]

  return (
    <div className="max-w-4xl mx-auto p-6">
      <Card>
        <Steps current={currentStep} className="mb-8">
          {steps.map((step) => (
            <Step key={step.title} title={step.title} />
          ))}
        </Steps>

        <div className="min-h-[400px]">{steps[currentStep].content}</div>

        <div className="flex justify-between mt-8 pt-6 border-t">
          <Button onClick={handleCancel}>取消</Button>
          <div className="space-x-2">
            {currentStep > 0 && (
              <Button onClick={handlePrev}>上一步</Button>
            )}
            {currentStep < steps.length - 1 && (
              <Button type="primary" onClick={handleNext}>
                下一步
              </Button>
            )}
            {currentStep === steps.length - 1 && (
              <Button type="primary" onClick={handleFinish}>
                创建项目
              </Button>
            )}
          </div>
        </div>
      </Card>
    </div>
  )
}

export default ProjectWizardNew
