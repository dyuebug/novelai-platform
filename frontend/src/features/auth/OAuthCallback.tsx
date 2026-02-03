import { useEffect, useState } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { Spin, Result, Button } from 'antd'
import { authService } from '@/services/auth.service'
import { useAuthStore } from '@/store/useAuthStore'

const OAuthCallback = () => {
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()
  const { setAuth } = useAuthStore()
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    handleCallback()
  }, [])

  const handleCallback = async () => {
    // 检查 URL 参数中是否直接包含 token（后端重定向模式）
    const accessToken = searchParams.get('access_token')
    const refreshToken = searchParams.get('refresh_token')

    if (accessToken && refreshToken) {
      // 直接从 URL 获取 token
      const user = {
        id: 'oauth-user',
        username: 'OAuth User',
        email: '',
        role: 'user',
      }
      setAuth(user, accessToken, refreshToken)
      navigate('/dashboard', { replace: true })
      return
    }

    // 检查是否有错误
    const errorParam = searchParams.get('error')
    if (errorParam) {
      const errorDesc = searchParams.get('error_description') || '授权失败'
      setError(`${errorParam}: ${errorDesc}`)
      return
    }

    // 获取授权码和 state
    const code = searchParams.get('code')
    const state = searchParams.get('state')

    if (!code || !state) {
      setError('缺少授权参数')
      return
    }

    // 从 URL 路径获取 provider（如果有的话）
    // 或者从 state 中解析（需要后端支持）
    const pathParts = window.location.pathname.split('/')
    const providerIndex = pathParts.indexOf('oauth')
    const provider = providerIndex >= 0 ? pathParts[providerIndex + 1] : null

    if (!provider || provider === 'callback') {
      // 如果 URL 中没有 provider，尝试从 localStorage 获取
      const savedProvider = localStorage.getItem('oauth_provider')
      if (savedProvider) {
        localStorage.removeItem('oauth_provider')
        await processCallback(savedProvider, code, state)
      } else {
        setError('无法确定 OAuth 提供商')
      }
      return
    }

    await processCallback(provider, code, state)
  }

  const processCallback = async (provider: string, code: string, state: string) => {
    try {
      const tokens = await authService.handleOAuthCallback(provider, code, state)

      // 简化处理：创建临时用户对象
      const user = {
        id: 'oauth-user',
        username: 'OAuth User',
        email: '',
        role: 'user',
      }

      setAuth(user, tokens.access_token, tokens.refresh_token)
      navigate('/dashboard', { replace: true })
    } catch (err) {
      console.error('OAuth callback failed:', err)
      setError('登录失败，请重试')
    }
  }

  if (error) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <Result
          status="error"
          title="登录失败"
          subTitle={error}
          extra={[
            <Button type="primary" key="login" onClick={() => navigate('/login')}>
              返回登录
            </Button>,
          ]}
        />
      </div>
    )
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="text-center">
        <Spin size="large" />
        <p className="mt-4 text-gray-500">正在处理登录...</p>
      </div>
    </div>
  )
}

export default OAuthCallback
