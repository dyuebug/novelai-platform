import { createBrowserRouter, Navigate } from 'react-router-dom'
import { lazy, Suspense } from 'react'
import { Spin } from 'antd'
import App from './App'

// 懒加载页面
const Login = lazy(() => import('@/features/auth/Login'))
const Register = lazy(() => import('@/features/auth/Register'))
const OAuthCallback = lazy(() => import('@/features/auth/OAuthCallback'))
const Dashboard = lazy(() => import('@/features/dashboard/Dashboard'))
const ProjectDetail = lazy(() => import('@/features/editor/ProjectDetail'))
const ChapterEditor = lazy(() => import('@/features/editor/ChapterEditor'))
const WorldManager = lazy(() => import('@/features/world'))
const ForeshadowManager = lazy(() => import('@/pages/ForeshadowManager'))
const Settings = lazy(() => import('@/pages/Settings'))

// 布局组件
const AppLayout = lazy(() => import('@/components/layout/AppLayout'))
const ProtectedRoute = lazy(() => import('@/features/auth/ProtectedRoute'))

// 加载中组件
const Loading = () => (
  <div className="flex items-center justify-center h-screen">
    <Spin size="large" />
  </div>
)

export const router = createBrowserRouter([
  {
    path: '/',
    element: <App />,
    children: [
      // 公开路由
      {
        path: 'login',
        element: (
          <Suspense fallback={<Loading />}>
            <Login />
          </Suspense>
        ),
      },
      {
        path: 'register',
        element: (
          <Suspense fallback={<Loading />}>
            <Register />
          </Suspense>
        ),
      },
      {
        path: 'oauth/callback',
        element: (
          <Suspense fallback={<Loading />}>
            <OAuthCallback />
          </Suspense>
        ),
      },
      {
        path: 'oauth/:provider/callback',
        element: (
          <Suspense fallback={<Loading />}>
            <OAuthCallback />
          </Suspense>
        ),
      },

      // 受保护路由
      {
        path: '',
        element: (
          <Suspense fallback={<Loading />}>
            <ProtectedRoute>
              <AppLayout />
            </ProtectedRoute>
          </Suspense>
        ),
        children: [
          {
            index: true,
            element: <Navigate to="/dashboard" replace />,
          },
          {
            path: 'dashboard',
            element: (
              <Suspense fallback={<Loading />}>
                <Dashboard />
              </Suspense>
            ),
          },
          {
            path: 'projects/:projectId',
            element: (
              <Suspense fallback={<Loading />}>
                <ProjectDetail />
              </Suspense>
            ),
          },
          {
            path: 'projects/:projectId/chapters/:chapterId',
            element: (
              <Suspense fallback={<Loading />}>
                <ChapterEditor />
              </Suspense>
            ),
          },
          {
            path: 'projects/:projectId/world',
            element: (
              <Suspense fallback={<Loading />}>
                <WorldManager />
              </Suspense>
            ),
          },
          {
            path: 'projects/:projectId/foreshadows',
            element: (
              <Suspense fallback={<Loading />}>
                <ForeshadowManager />
              </Suspense>
            ),
          },
          {
            path: 'settings',
            element: (
              <Suspense fallback={<Loading />}>
                <Settings />
              </Suspense>
            ),
          },
        ],
      },
    ],
  },
])
