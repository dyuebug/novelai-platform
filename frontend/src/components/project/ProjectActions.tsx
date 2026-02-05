import { Space, Button, Dropdown, Modal, message } from 'antd'
import {
  EditOutlined,
  DeleteOutlined,
  ExportOutlined,
  CopyOutlined,
  MoreOutlined,
  ExclamationCircleOutlined,
} from '@ant-design/icons'
import type { FC } from 'react'
import type { MenuProps } from 'antd'

const { confirm } = Modal

interface ProjectActionsProps {
  projectId: string
  projectTitle?: string
  onEdit?: (id: string) => void
  onDelete?: (id: string) => void
  onExport?: (id: string) => void
  onDuplicate?: (id: string) => void
  compact?: boolean
}

export const ProjectActions: FC<ProjectActionsProps> = ({
  projectId,
  projectTitle = '项目',
  onEdit,
  onDelete,
  onExport,
  onDuplicate,
  compact = false,
}) => {
  const handleEdit = () => {
    onEdit?.(projectId)
  }

  const handleDelete = () => {
    confirm({
      title: '确认删除',
      icon: <ExclamationCircleOutlined />,
      content: `确定要删除项目"${projectTitle}"吗？此操作可以在回收站中恢复。`,
      okText: '删除',
      okType: 'danger',
      cancelText: '取消',
      onOk: () => {
        onDelete?.(projectId)
        message.success('项目已删除')
      },
    })
  }

  const handleExport = () => {
    onExport?.(projectId)
    message.loading('正在导出项目...', 0)
  }

  const handleDuplicate = () => {
    onDuplicate?.(projectId)
    message.success('项目已复制')
  }

  const menuItems: MenuProps['items'] = [
    {
      key: 'edit',
      label: '编辑',
      icon: <EditOutlined />,
      onClick: handleEdit,
    },
    {
      key: 'duplicate',
      label: '复制',
      icon: <CopyOutlined />,
      onClick: handleDuplicate,
    },
    {
      key: 'export',
      label: '导出',
      icon: <ExportOutlined />,
      onClick: handleExport,
    },
    {
      type: 'divider',
    },
    {
      key: 'delete',
      label: '删除',
      icon: <DeleteOutlined />,
      danger: true,
      onClick: handleDelete,
    },
  ]

  if (compact) {
    return (
      <Dropdown menu={{ items: menuItems }} trigger={['click']}>
        <Button icon={<MoreOutlined />} />
      </Dropdown>
    )
  }

  return (
    <Space>
      <Button icon={<EditOutlined />} onClick={handleEdit}>
        编辑
      </Button>
      <Button icon={<ExportOutlined />} onClick={handleExport}>
        导出
      </Button>
      <Dropdown menu={{ items: menuItems }} trigger={['click']}>
        <Button icon={<MoreOutlined />} />
      </Dropdown>
    </Space>
  )
}
