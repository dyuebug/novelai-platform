import { Tabs } from 'antd'
import { useSearchParams } from 'react-router-dom'
import {
  UserOutlined,
  EnvironmentOutlined,
  TeamOutlined,
  BookOutlined,
  ApartmentOutlined,
} from '@ant-design/icons'
import CharacterManager from './CharacterManager'
import LocationManager from './LocationManager'
import OrganizationManager from './OrganizationManager'
import WorldSettingManager from './WorldSettingManager'
import RelationshipGraph from './RelationshipGraph'

const WorldManager = () => {
  const [searchParams, setSearchParams] = useSearchParams()
  const activeTab = searchParams.get('tab') || 'characters'

  const handleTabChange = (key: string) => {
    setSearchParams({ tab: key })
  }

  const tabItems = [
    {
      key: 'characters',
      label: (
        <span>
          <UserOutlined />
          角色
        </span>
      ),
      children: <CharacterManager />,
    },
    {
      key: 'relationships',
      label: (
        <span>
          <ApartmentOutlined />
          关系图
        </span>
      ),
      children: <RelationshipGraph />,
    },
    {
      key: 'locations',
      label: (
        <span>
          <EnvironmentOutlined />
          地点
        </span>
      ),
      children: <LocationManager />,
    },
    {
      key: 'organizations',
      label: (
        <span>
          <TeamOutlined />
          组织
        </span>
      ),
      children: <OrganizationManager />,
    },
    {
      key: 'settings',
      label: (
        <span>
          <BookOutlined />
          世界设定
        </span>
      ),
      children: <WorldSettingManager />,
    },
  ]

  return (
    <div className="p-6">
      <Tabs
        activeKey={activeTab}
        onChange={handleTabChange}
        items={tabItems}
        size="large"
      />
    </div>
  )
}

export default WorldManager
