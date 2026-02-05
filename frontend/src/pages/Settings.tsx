import { Card, Form, Input, Button, Select, message, Divider } from 'antd'
import { useAuthStore } from '@/store/useAuthStore'

const { Option } = Select

const Settings = () => {
  const { user } = useAuthStore()
  const [form] = Form.useForm()

  const handleSave = async (values: any) => {
    try {
      // TODO: 实现保存设置的 API 调用
      console.log('保存设置:', values)
      message.success('设置已保存')
    } catch (error) {
      message.error('保存设置失败')
    }
  }

  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold mb-6">用户设置</h1>

      <Card title="基本信息" className="mb-6">
        <Form
          form={form}
          layout="vertical"
          initialValues={{
            username: user?.username,
            email: user?.email,
          }}
          onFinish={handleSave}
        >
          <Form.Item label="用户名" name="username">
            <Input disabled />
          </Form.Item>

          <Form.Item label="邮箱" name="email">
            <Input disabled />
          </Form.Item>

          <Form.Item label="角色" name="role">
            <Input disabled value={user?.role} />
          </Form.Item>
        </Form>
      </Card>

      <Card title="AI 配置" className="mb-6">
        <Form layout="vertical" onFinish={handleSave}>
          <Form.Item
            label="默认 AI 提供商"
            name="defaultProvider"
            initialValue="openai"
          >
            <Select>
              <Option value="openai">OpenAI</Option>
              <Option value="anthropic">Anthropic Claude</Option>
              <Option value="gemini">Google Gemini</Option>
            </Select>
          </Form.Item>

          <Form.Item
            label="默认模型"
            name="defaultModel"
            initialValue="gpt-4o"
          >
            <Select>
              <Option value="gpt-4o">GPT-4o</Option>
              <Option value="gpt-4-turbo">GPT-4 Turbo</Option>
              <Option value="gpt-3.5-turbo">GPT-3.5 Turbo</Option>
              <Option value="claude-3-opus">Claude 3 Opus</Option>
              <Option value="claude-3-sonnet">Claude 3 Sonnet</Option>
              <Option value="gemini-pro">Gemini Pro</Option>
            </Select>
          </Form.Item>

          <Form.Item
            label="OpenAI API Key"
            name="openaiApiKey"
            tooltip="留空则使用系统默认配置"
          >
            <Input.Password placeholder="sk-..." />
          </Form.Item>

          <Form.Item
            label="Anthropic API Key"
            name="anthropicApiKey"
            tooltip="留空则使用系统默认配置"
          >
            <Input.Password placeholder="sk-ant-..." />
          </Form.Item>

          <Form.Item
            label="Gemini API Key"
            name="geminiApiKey"
            tooltip="留空则使用系统默认配置"
          >
            <Input.Password />
          </Form.Item>

          <Divider />

          <Form.Item>
            <Button type="primary" htmlType="submit">
              保存设置
            </Button>
          </Form.Item>
        </Form>
      </Card>

      <Card title="编辑器设置" className="mb-6">
        <Form layout="vertical" onFinish={handleSave}>
          <Form.Item
            label="主题"
            name="theme"
            initialValue="light"
          >
            <Select>
              <Option value="light">浅色</Option>
              <Option value="dark">深色</Option>
              <Option value="auto">跟随系统</Option>
            </Select>
          </Form.Item>

          <Form.Item
            label="编辑器字体大小"
            name="editorFontSize"
            initialValue={16}
          >
            <Select>
              <Option value={12}>12px</Option>
              <Option value={14}>14px</Option>
              <Option value={16}>16px (默认)</Option>
              <Option value={18}>18px</Option>
              <Option value={20}>20px</Option>
            </Select>
          </Form.Item>

          <Form.Item
            label="语言"
            name="language"
            initialValue="zh-CN"
          >
            <Select>
              <Option value="zh-CN">简体中文</Option>
              <Option value="en-US">English</Option>
            </Select>
          </Form.Item>

          <Divider />

          <Form.Item>
            <Button type="primary" htmlType="submit">
              保存设置
            </Button>
          </Form.Item>
        </Form>
      </Card>

      <Card title="密码修改">
        <Form layout="vertical" onFinish={handleSave}>
          <Form.Item
            label="当前密码"
            name="currentPassword"
            rules={[{ required: true, message: '请输入当前密码' }]}
          >
            <Input.Password />
          </Form.Item>

          <Form.Item
            label="新密码"
            name="newPassword"
            rules={[
              { required: true, message: '请输入新密码' },
              { min: 8, message: '密码至少 8 个字符' },
            ]}
          >
            <Input.Password />
          </Form.Item>

          <Form.Item
            label="确认新密码"
            name="confirmPassword"
            dependencies={['newPassword']}
            rules={[
              { required: true, message: '请确认新密码' },
              ({ getFieldValue }) => ({
                validator(_, value) {
                  if (!value || getFieldValue('newPassword') === value) {
                    return Promise.resolve()
                  }
                  return Promise.reject(new Error('两次输入的密码不一致'))
                },
              }),
            ]}
          >
            <Input.Password />
          </Form.Item>

          <Divider />

          <Form.Item>
            <Button type="primary" htmlType="submit">
              修改密码
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </div>
  )
}

export default Settings
