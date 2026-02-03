import { useState, useEffect } from 'react'
import { Button, Space, message, Spin } from 'antd'
import { GithubOutlined, GoogleOutlined, LoginOutlined } from '@ant-design/icons'
import { authService } from '@/services/auth.service'

interface OAuthButtonsProps {
  onSuccess?: () => void
  disabled?: boolean
}

// OAuth 提供商配置
const providerConfig: Record<string, { icon: React.ReactNode; label: string; color: string }> = {
  github: {
    icon: <GithubOutlined />,
    label: 'GitHub',
    color: '#24292e',
  },
  google: {
    icon: <GoogleOutlined />,
    label: 'Google',
    color: '#4285f4',
  },
  linuxdo: {
    icon: <LoginOutlined />,
    label: 'LinuxDo',
    color: '#1890ff',
  },
}

const OAuthButtons: React.FC<OAuthButtonsProps> = ({ disabled = false }) => {
  const [providers, setProviders] = useState<string[]>([])
  const [loading, setLoading] = useState(true)
  const [authLoading, setAuthLoading] = useState<string | null>(null)

  useEffect(() => {
    loadProviders()
  }, [])

  const loadProviders = async () => {
    try {
      const enabledProviders = await authService.getOAuthProviders()
      setProviders(enabledProviders)
    } catch (error) {
      console.error('Failed to load OAuth providers:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleOAuthLogin = async (provider: string) => {
    setAuthLoading(provider)
    try {
      // 获取当前页面 URL 作为回调地址
      const redirectUri = `${window.location.origin}/oauth/callback`
      const authUrl = await authService.getOAuthAuthorizationUrl(provider, redirectUri)

      // 跳转到 OAuth 授权页面
      window.location.href = authUrl
    } catch (error) {
      console.error('OAuth login failed:', error)
      message.error('获取授权链接失败，请稍后重试')
      setAuthLoading(null)
    }
  }

  if (loading) {
    return (
      <div className="text-center py-4">
        <Spin size="small" />
      </div>
    )
  }

  if (providers.length === 0) {
    return null
  }

  return (
    <div className="oauth-buttons">
      <div className="relative my-6">
        <div className="absolute inset-0 flex items-center">
          <div className="w-full border-t border-gray-200" />
        </div>
        <div className="relative flex justify-center text-sm">
          <span className="px-4 bg-white text-gray-500">或使用以下方式登录</span>
        </div>
      </div>

      <Space direction="vertical" className="w-full" size="middle">
        {providers.map((provider) => {
          const config = providerConfig[provider]
          if (!config) return null

          return (
            <Button
              key={provider}
              icon={config.icon}
              size="large"
              block
              loading={authLoading === provider}
              disabled={disabled || authLoading !== null}
              onClick={() => handleOAuthLogin(provider)}
              style={{
                backgroundColor: 'white',
                borderColor: config.color,
                color: config.color,
              }}
              className="hover:opacity-80 transition-opacity"
            >
              使用 {config.label} 登录
            </Button>
          )
        })}
      </Space>
    </div>
  )
}

export default OAuthButtons
