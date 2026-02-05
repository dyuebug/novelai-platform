import { Space, Button, Tooltip, Divider } from 'antd'
import {
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  SaveOutlined,
  UndoOutlined,
  RedoOutlined,
  RobotOutlined,
  SettingOutlined,
  LoadingOutlined,
} from '@ant-design/icons'
import type { FC } from 'react'

interface EditorToolbarProps {
  onToggleSidebar?: () => void
  sidebarCollapsed?: boolean
  onSave?: () => void
  saving?: boolean
  onUndo?: () => void
  onRedo?: () => void
  onOpenAI?: () => void
  onOpenSettings?: () => void
}

export const EditorToolbar: FC<EditorToolbarProps> = ({
  onToggleSidebar,
  sidebarCollapsed = false,
  onSave,
  saving = false,
  onUndo,
  onRedo,
  onOpenAI,
  onOpenSettings,
}) => {
  return (
    <div className="flex items-center justify-between px-4 py-2">
      <Space split={<Divider type="vertical" />}>
        {/* 侧边栏切换 */}
        <Tooltip title={sidebarCollapsed ? '展开侧边栏' : '折叠侧边栏'}>
          <Button
            type="text"
            icon={sidebarCollapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
            onClick={onToggleSidebar}
          />
        </Tooltip>

        {/* 保存 */}
        <Tooltip title="保存 (Ctrl+S)">
          <Button
            type="text"
            icon={saving ? <LoadingOutlined /> : <SaveOutlined />}
            onClick={onSave}
            disabled={saving}
          >
            {saving ? '保存中...' : '保存'}
          </Button>
        </Tooltip>

        {/* 撤销/重做 */}
        <Space>
          <Tooltip title="撤销 (Ctrl+Z)">
            <Button type="text" icon={<UndoOutlined />} onClick={onUndo} />
          </Tooltip>
          <Tooltip title="重做 (Ctrl+Y)">
            <Button type="text" icon={<RedoOutlined />} onClick={onRedo} />
          </Tooltip>
        </Space>
      </Space>

      <Space split={<Divider type="vertical" />}>
        {/* AI助手 */}
        <Tooltip title="AI助手 (Ctrl+Shift+A)">
          <Button type="primary" icon={<RobotOutlined />} onClick={onOpenAI}>
            AI助手
          </Button>
        </Tooltip>

        {/* 设置 */}
        <Tooltip title="编辑器设置">
          <Button type="text" icon={<SettingOutlined />} onClick={onOpenSettings} />
        </Tooltip>
      </Space>
    </div>
  )
}
