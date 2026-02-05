import { Empty, Button } from 'antd'
import { PlusOutlined, FileTextOutlined } from '@ant-design/icons'
import type { FC, ReactNode } from 'react'

interface EmptyStateProps {
  title?: string
  description?: string
  icon?: ReactNode
  actionText?: string
  onAction?: () => void
  showAction?: boolean
}

export const EmptyState: FC<EmptyStateProps> = ({
  title = '暂无项目',
  description = '创建您的第一个小说项目，开始创作之旅',
  icon = <FileTextOutlined className="text-6xl text-gray-300" />,
  actionText = '创建新项目',
  onAction,
  showAction = true,
}) => {
  return (
    <div className="flex flex-col items-center justify-center py-20">
      <Empty
        image={<div className="mb-4">{icon}</div>}
        description={
          <div className="space-y-2">
            <div className="text-lg font-medium text-gray-700">{title}</div>
            <div className="text-sm text-gray-500">{description}</div>
          </div>
        }
      >
        {showAction && onAction && (
          <Button
            type="primary"
            size="large"
            icon={<PlusOutlined />}
            onClick={onAction}
            className="mt-4"
          >
            {actionText}
          </Button>
        )}
      </Empty>
    </div>
  )
}
