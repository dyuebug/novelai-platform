import { useState, useEffect } from 'react'
import { Layout, Menu, Avatar, Dropdown, Button, Drawer } from 'antd'
import { Outlet, useNavigate, useLocation } from 'react-router-dom'
import {
  HomeOutlined,
  FileTextOutlined,
  SettingOutlined,
  LogoutOutlined,
  UserOutlined,
  MenuOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
} from '@ant-design/icons'
import { useAuthStore } from '@/store/useAuthStore'

const { Header, Sider, Content } = Layout

// 响应式断点
const BREAKPOINTS = {
  mobile: 768,
  tablet: 1024,
}

// 自定义 Hook: 检测屏幕尺寸
const useResponsive = () => {
  const [windowWidth, setWindowWidth] = useState(
    typeof window !== 'undefined' ? window.innerWidth : 1200
  )

  useEffect(() => {
    const handleResize = () => setWindowWidth(window.innerWidth)
    window.addEventListener('resize', handleResize)
    return () => window.removeEventListener('resize', handleResize)
  }, [])

  return {
    isMobile: windowWidth < BREAKPOINTS.mobile,
    isTablet: windowWidth >= BREAKPOINTS.mobile && windowWidth < BREAKPOINTS.tablet,
    isDesktop: windowWidth >= BREAKPOINTS.tablet,
    windowWidth,
  }
}

const AppLayout = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const { user, clearAuth } = useAuthStore()
  const { isMobile, isTablet } = useResponsive()

  const [collapsed, setCollapsed] = useState(false)
  const [drawerVisible, setDrawerVisible] = useState(false)

  // 在平板模式下默认折叠侧边栏
  useEffect(() => {
    setCollapsed(isTablet)
  }, [isTablet])

  // 路由变化时关闭移动端抽屉
  useEffect(() => {
    setDrawerVisible(false)
  }, [location.pathname])

  const menuItems = [
    {
      key: '/dashboard',
      icon: <HomeOutlined />,
      label: '仪表盘',
    },
    {
      key: '/projects',
      icon: <FileTextOutlined />,
      label: '我的项目',
    },
    {
      key: '/settings',
      icon: <SettingOutlined />,
      label: '设置',
    },
  ]

  const userMenuItems = [
    {
      key: 'profile',
      icon: <UserOutlined />,
      label: '个人资料',
    },
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: '退出登录',
      danger: true,
    },
  ]

  const handleMenuClick = ({ key }: { key: string }) => {
    navigate(key)
  }

  const handleUserMenuClick = ({ key }: { key: string }) => {
    if (key === 'logout') {
      clearAuth()
      navigate('/login')
    } else if (key === 'profile') {
      navigate('/settings')
    }
  }

  // 侧边栏菜单内容
  const SidebarContent = () => (
    <>
      <div className="h-16 flex items-center justify-center border-b">
        <h1 className={`font-bold text-blue-600 ${collapsed ? 'text-sm' : 'text-xl'}`}>
          {collapsed ? 'N' : 'NovelAI'}
        </h1>
      </div>
      <Menu
        mode="inline"
        selectedKeys={[location.pathname]}
        items={menuItems}
        onClick={handleMenuClick}
        className="border-none"
        inlineCollapsed={collapsed}
      />
    </>
  )

  // 移动端布局
  if (isMobile) {
    return (
      <Layout className="min-h-screen">
        <Header className="bg-white border-b px-4 flex items-center justify-between fixed top-0 left-0 right-0 z-50">
          <div className="flex items-center">
            <Button
              type="text"
              icon={<MenuOutlined />}
              onClick={() => setDrawerVisible(true)}
              className="mr-2"
            />
            <h1 className="text-lg font-bold text-blue-600 m-0">NovelAI</h1>
          </div>
          <Dropdown
            menu={{
              items: userMenuItems,
              onClick: handleUserMenuClick,
            }}
            placement="bottomRight"
          >
            <Avatar icon={<UserOutlined />} className="bg-blue-500 cursor-pointer" />
          </Dropdown>
        </Header>

        <Drawer
          title="菜单"
          placement="left"
          onClose={() => setDrawerVisible(false)}
          open={drawerVisible}
          width={250}
          styles={{ body: { padding: 0 } }}
        >
          <Menu
            mode="inline"
            selectedKeys={[location.pathname]}
            items={menuItems}
            onClick={handleMenuClick}
            className="border-none"
          />
        </Drawer>

        <Content className="p-4 bg-gray-50 mt-16">
          <Outlet />
        </Content>
      </Layout>
    )
  }

  // 桌面端/平板端布局
  return (
    <Layout className="min-h-screen">
      <Sider
        theme="light"
        className="border-r"
        width={220}
        collapsedWidth={80}
        collapsed={collapsed}
        trigger={null}
      >
        <SidebarContent />
      </Sider>

      <Layout>
        <Header className="bg-white border-b px-6 flex items-center justify-between">
          <Button
            type="text"
            icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
            onClick={() => setCollapsed(!collapsed)}
          />
          <Dropdown
            menu={{
              items: userMenuItems,
              onClick: handleUserMenuClick,
            }}
            placement="bottomRight"
          >
            <div className="flex items-center cursor-pointer">
              <Avatar icon={<UserOutlined />} className="bg-blue-500" />
              <span className="ml-2 text-gray-700">{user?.username || '用户'}</span>
            </div>
          </Dropdown>
        </Header>

        <Content className="p-6 bg-gray-50">
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  )
}

export default AppLayout
