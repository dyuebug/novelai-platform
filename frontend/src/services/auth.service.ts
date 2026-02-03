import api from '@/lib/axios'

export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  username: string
  email: string
  password: string
}

export interface TokenResponse {
  access_token: string
  refresh_token: string
  token_type: string
  expires_in: number
}

export interface User {
  id: string
  username: string
  email: string
  role: string
  created_at: string
}

export interface OAuthProvider {
  name: string
  displayName: string
  icon: string
}

export const authService = {
  login: async (data: LoginRequest): Promise<TokenResponse> => {
    const response = await api.post('/auth/login', data)
    return response.data
  },

  register: async (data: RegisterRequest): Promise<User> => {
    const response = await api.post('/auth/register', data)
    return response.data
  },

  refresh: async (refreshToken: string): Promise<TokenResponse> => {
    const response = await api.post('/auth/refresh', {
      refresh_token: refreshToken,
    })
    return response.data
  },

  getCurrentUser: async (): Promise<User> => {
    const response = await api.get('/auth/me')
    return response.data
  },

  requestPasswordReset: async (email: string): Promise<void> => {
    await api.post('/auth/password-reset/request', { email })
  },

  resetPassword: async (token: string, newPassword: string): Promise<void> => {
    await api.post('/auth/password-reset/verify', {
      token,
      new_password: newPassword,
    })
  },

  // OAuth 相关方法
  getOAuthProviders: async (): Promise<string[]> => {
    const response = await api.get('/oauth/providers')
    return response.data.providers || []
  },

  getOAuthAuthorizationUrl: async (
    provider: string,
    redirectUri?: string
  ): Promise<string> => {
    const params = redirectUri ? `?redirect_uri=${encodeURIComponent(redirectUri)}` : ''
    const response = await api.get(`/oauth/${provider}/authorize${params}`)
    return response.data.authorization_url
  },

  handleOAuthCallback: async (
    provider: string,
    code: string,
    state: string
  ): Promise<TokenResponse> => {
    const response = await api.post(`/oauth/${provider}/callback`, { code, state })
    return response.data
  },
}
