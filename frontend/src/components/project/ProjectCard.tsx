import { Card, Tag, Space, Button, Tooltip, Typography } from 'antd'
import { EditOutlined, DeleteOutlined, EyeOutlined, BookOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import type { FC } from 'react'

const { Text, Paragraph } = Typography

interface ProjectCardProps {
  id: string
  title: string
  description?: string
  coverImageUrl?: string
  tags?: string[]
  totalChapters?: number
  totalWords?: number
  lastUpdateTime?: string
  status?: 'draft' | 'in_progress' | 'completed'
  onEdit?: (id: string) => void
  onDelete?: (id: string) => void
  onSelect?: (id: string, selected: boolean) => void
  selected?: boolean
  selectionMode?: boolean
}

const statusColors = {
  draft: 'default',
  in_progress: 'processing',
  completed: 'success',
} as const

const statusLabels = {
  draft: '草稿',
  in_progress: '进行中',
  completed: '已完成',
} as const

export const ProjectCard: FC<ProjectCardProps> = ({
  id,
  title,
  description,
  coverImageUrl,
  tags = [],
  totalChapters = 0,
  totalWords = 0,
  lastUpdateTime,
  status = 'draft',
  onEdit,
  onDelete,
  onSelect,
  selected = false,
  selectionMode = false,
}) => {
  const navigate = useNavigate()

  const handleView = () => {
    navigate(`/projects/${id}`)
  }

  const handleEdit = (e: React.MouseEvent) => {
    e.stopPropagation()
    onEdit?.(id)
  }

  const handleDelete = (e: React.MouseEvent) => {
    e.stopPropagation()
    onDelete?.(id)
  }

  const handleSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    e.stopPropagation()
    onSelect?.(id, e.target.checked)
  }

  const formatNumber = (num: number) => {
    if (num >= 10000) {
      return `${(num / 10000).toFixed(1)}万`
    }
    return num.toString()
  }

  const formatDate = (dateStr?: string) => {
    if (!dateStr) return '未知'
    const date = new Date(dateStr)
    const now = new Date()
    const diff = now.getTime() - date.getTime()
    const days = Math.floor(diff / (1000 * 60 * 60 * 24))

    if (days === 0) return '今天'
    if (days === 1) return '昨天'
    if (days < 7) return `${days}天前`
    if (days < 30) return `${Math.floor(days / 7)}周前`
    if (days < 365) return `${Math.floor(days / 30)}个月前`
    return `${Math.floor(days / 365)}年前`
  }

  return (
    <Card
      hoverable
      className={`project-card ${selected ? 'selected' : ''}`}
      onClick={handleView}
      cover={
        coverImageUrl ? (
          <div className="relative h-48 overflow-hidden">
            <img
              alt={title}
              src={coverImageUrl}
              className="w-full h-full object-cover"
            />
            {selectionMode && (
              <div className="absolute top-2 left-2">
                <input
                  type="checkbox"
                  checked={selected}
                  onChange={handleSelect}
                  className="w-5 h-5 cursor-pointer"
                  onClick={(e) => e.stopPropagation()}
                />
              </div>
            )}
          </div>
        ) : (
          <div className="relative h-48 bg-gray-100 flex items-center justify-center">
            <BookOutlined className="text-6xl text-gray-400" />
            {selectionMode && (
              <div className="absolute top-2 left-2">
                <input
                  type="checkbox"
                  checked={selected}
                  onChange={handleSelect}
                  className="w-5 h-5 cursor-pointer"
                  onClick={(e) => e.stopPropagation()}
                />
              </div>
            )}
          </div>
        )
      }
      actions={[
        <Tooltip title="查看" key="view">
          <Button type="text" icon={<EyeOutlined />} onClick={handleView} />
        </Tooltip>,
        <Tooltip title="编辑" key="edit">
          <Button type="text" icon={<EditOutlined />} onClick={handleEdit} />
        </Tooltip>,
        <Tooltip title="删除" key="delete">
          <Button type="text" danger icon={<DeleteOutlined />} onClick={handleDelete} />
        </Tooltip>,
      ]}
    >
      <div className="space-y-2">
        <div className="flex items-start justify-between">
          <Text strong className="text-lg line-clamp-1">
            {title}
          </Text>
          <Tag color={statusColors[status]}>{statusLabels[status]}</Tag>
        </div>

        {description && (
          <Paragraph
            ellipsis={{ rows: 2 }}
            className="text-gray-600 text-sm mb-2"
          >
            {description}
          </Paragraph>
        )}

        {tags.length > 0 && (
          <div className="flex flex-wrap gap-1">
            {tags.slice(0, 3).map((tag) => (
              <Tag key={tag} className="text-xs">
                {tag}
              </Tag>
            ))}
            {tags.length > 3 && (
              <Tag className="text-xs">+{tags.length - 3}</Tag>
            )}
          </div>
        )}

        <div className="flex items-center justify-between text-sm text-gray-500 pt-2 border-t">
          <Space size="large">
            <span>{totalChapters} 章节</span>
            <span>{formatNumber(totalWords)} 字</span>
          </Space>
          <span>{formatDate(lastUpdateTime)}</span>
        </div>
      </div>
    </Card>
  )
}
