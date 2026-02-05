import { useState, useEffect } from 'react'
import { Layout, Menu, Button, Breadcrumb, Drawer, message } from 'antd'
import {
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  FileTextOutlined,
  TeamOutlined,
  EnvironmentOutlined,
  BankOutlined,
  GlobalOutlined,
  BulbOutlined,
  HomeOutlined,
} from '@ant-design/icons'
import { Outlet, useNavigate, useLocation, useParams } from 'react-router-dom'
import type { MenuProps } from 'antd'

const { Sider, Content } = Layout

type MenuItem = Required<MenuProps>['items'][number]

const menuItems: MenuItem[] = [
  {
    key: 'chapters',
    icon: <FileTextOutlined />,
    label: '章节',
  },
  {
    key: 'characters',
    icon: <TeamOutlined />,
    label: '角色',
  },
  {
    key: 'locations',
    icon: <EnvironmentOutlined />,
    label: '地点',
  },
  {
    key: 'organizations',
    icon: <BankOutlined />,
    label: '组织',
  },
  {
    key: 'world-settings',
    icon: <GlobalOutlined />,
    label: '世界设定',
  },
  {
    key: 'foreshadows',
    icon: <BulbOutlined />,
    label: '伏笔',
  },
]

export default function ProjectDetail() {
  const { projectId } = useParams<{ projectId: string }>()
  const navigate = useNavigate()
  const location = useLocation()
  const [collapsed, setCollapsed] = useState(false)
  const [mobileDrawerVisible, setMobileDrawerVisible] = useState(false)
  const [isMobile, setIsMobile] = useState(false)
  const [projectTitle, setProjectTitle] = useState('项目详情')

  // 检测移动端
  useEffect(() => {
    const checkMobile = () => {
      setIsMobile(window.innerWidth <= 768)
    }
    checkMobile()
    window.addEventListener('resize', checkMobile)
    return () => window.removeEventListener('resize', checkMobile)
  }, [])

  // 加载项目信息
  useEffect(() => {
    const loadProject = async () => {
      try {
        // TODO: 调用实际API
        // const project = await projectService.getProject(projectId!)
        // setProjectTitle(project.title)
        setProjectTitle('示例项目')
      } catch (error) {
        message.error('加载项目信息失败')
      }
    }
    if (projectId) {
      loadProject()
    }
  }, [projectId])

  // 获取当前选中的菜单项
  const selectedKey = location.pathname.split('/').pop() || 'chapters'

  const handleMenuClick = ({ key }: { key: string }) => {
    navigate(`/projects/${projectId}/${key}`)
    if (isMobile) {
      setMobileDrawerVisible(false)
    }
  }

  const renderMenu = () => (
    <Menu
      mode="inline"
      selectedKeys={[selectedKey]}
      items={menuItems}
      onClick={handleMenuClick}
      className="border-r-0"
    />
  )

  const renderBreadcrumb = () => (
    <Breadcrumb
      items={[
        {
          title: (
            <a onClick={() => navigate('/projects')}>
              <HomeOutlined /> 项目列表
            </a>
          ),
        },
        {
          title: projectTitle,
        },
        {
          title: (menuItems.find(item => item && 'key' in item && item.key === selectedKey) as any)?.label || '',
        },
      ]}
    />
  )

  return (
    <Layout className="min-h-screen">
      {/* 桌面端侧边栏 */}
      {!isMobile && (
        <Sider
          collapsible
          collapsed={collapsed}
          onCollapse={setCollapsed}
          width={220}
          collapsedWidth={80}
          className="bg-white border-r"
          trigger={null}
        >
          <div className="p-4 border-b">
            <h2 className={`font-bold truncate ${collapsed ? 'text-center' : ''}`}>
              {collapsed ? '项目' : projectTitle}
            </h2>
          </div>
          {renderMenu()}
        </Sider>
      )}

      {/* 移动端抽屉 */}
      {isMobile && (
        <Drawer
          title={projectTitle}
          placement="left"
          onClose={() => setMobileDrawerVisible(false)}
          open={mobileDrawerVisible}
          width={250}
        >
          {renderMenu()}
        </Drawer>
      )}

      <Layout>
        {/* 顶部栏 */}
        <div className="bg-white border-b px-6 py-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-4">
              <Button
                type="text"
                icon={collapsed || isMobile ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
                onClick={() => {
                  if (isMobile) {
                    setMobileDrawerVisible(true)
                  } else {
                    setCollapsed(!collapsed)
                  }
                }}
              />
              {renderBreadcrumb()}
            </div>
          </div>
        </div>

        {/* 主内容区 */}
        <Content className="p-6 bg-gray-50">
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  )
}
