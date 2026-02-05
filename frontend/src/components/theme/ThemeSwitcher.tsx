import { Dropdown, Button } from 'antd'
import { BulbOutlined, BulbFilled } from '@ant-design/icons'
import { useThemeStore } from '@/store/useThemeStore'
import type { FC } from 'react'
import type { MenuProps } from 'antd'

export const ThemeSwitcher: FC = () => {
  const { mode, setMode } = useThemeStore()

  const menuItems: MenuProps['items'] = [
    {
      key: 'light',
      label: '亮色模式',
      icon: <BulbOutlined />,
    },
    {
      key: 'dark',
      label: '暗色模式',
      icon: <BulbFilled />,
    },
    {
      key: 'auto',
      label: '跟随系统',
      icon: <BulbOutlined />,
    },
  ]

  const handleMenuClick: MenuProps['onClick'] = ({ key }) => {
    setMode(key as 'light' | 'dark' | 'auto')
  }

  const getIcon = () => {
    switch (mode) {
      case 'dark':
        return <BulbFilled />
      case 'light':
        return <BulbOutlined />
      case 'auto':
        return <BulbOutlined />
      default:
        return <BulbOutlined />
    }
  }

  const getLabel = () => {
    switch (mode) {
      case 'dark':
        return '暗色'
      case 'light':
        return '亮色'
      case 'auto':
        return '自动'
      default:
        return '主题'
    }
  }

  return (
    <Dropdown
      menu={{ items: menuItems, onClick: handleMenuClick, selectedKeys: [mode] }}
      trigger={['click']}
    >
      <Button icon={getIcon()}>
        {getLabel()}
      </Button>
    </Dropdown>
  )
}
