import { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { Form, Input, Button, Card, message } from 'antd'
import { LockOutlined, MailOutlined } from '@ant-design/icons'
import { authService } from '@/services/auth.service'
import { useAuthStore } from '@/store/useAuthStore'
import OAuthButtons from '@/components/auth/OAuthButtons'

interface LoginFormData {
  email: string
  password: string
}

const Login = () => {
  const [loading, setLoading] = useState(false)
  const navigate = useNavigate()
  const { setAuth } = useAuthStore()

  const onFinish = async (values: LoginFormData) => {
    setLoading(true)
    try {
      const tokens = await authService.login(values)

      // 获取用户信息
      // 简化处理：从 token 解析用户信息
      const user = {
        id: 'temp-id',
        username: values.email.split('@')[0],
        email: values.email,
        role: 'user',
      }

      setAuth(user, tokens.access_token, tokens.refresh_token)
      message.success('登录成功')
      navigate('/dashboard')
    } catch (error: unknown) {
      const err = error as { response?: { data?: { message?: string } } }
      message.error(err.response?.data?.message || '登录失败')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <Card className="w-full max-w-md shadow-lg">
        <div className="text-center mb-8">
          <h1 className="text-2xl font-bold text-gray-800">NovelAI</h1>
          <p className="text-gray-500 mt-2">智能小说创作助手</p>
        </div>

        <Form
          name="login"
          onFinish={onFinish}
          autoComplete="off"
          layout="vertical"
          size="large"
        >
          <Form.Item
            name="email"
            rules={[
              { required: true, message: '请输入邮箱' },
              { type: 'email', message: '请输入有效的邮箱地址' },
            ]}
          >
            <Input
              prefix={<MailOutlined className="text-gray-400" />}
              placeholder="邮箱"
            />
          </Form.Item>

          <Form.Item
            name="password"
            rules={[{ required: true, message: '请输入密码' }]}
          >
            <Input.Password
              prefix={<LockOutlined className="text-gray-400" />}
              placeholder="密码"
            />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading} block>
              登录
            </Button>
          </Form.Item>
        </Form>

        <OAuthButtons disabled={loading} />

        <div className="text-center mt-6">
          <span className="text-gray-500">还没有账号？</span>
          <Link to="/register" className="text-blue-500 ml-1">
            立即注册
          </Link>
        </div>
      </Card>
    </div>
  )
}

export default Login
