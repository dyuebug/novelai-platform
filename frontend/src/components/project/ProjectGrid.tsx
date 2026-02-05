import { Row, Col, Empty, Spin } from 'antd'
import { ProjectCard } from './ProjectCard'
import type { FC } from 'react'

interface Project {
  id: string
  title: string
  description?: string
  coverImageUrl?: string
  tags?: string[]
  totalChapters?: number
  totalWords?: number
  lastUpdateTime?: string
  status?: 'draft' | 'in_progress' | 'completed'
}

interface ProjectGridProps {
  projects: Project[]
  loading?: boolean
  onEdit?: (id: string) => void
  onDelete?: (id: string) => void
  onSelect?: (id: string, selected: boolean) => void
  selectedIds?: string[]
  selectionMode?: boolean
  emptyText?: string
}

export const ProjectGrid: FC<ProjectGridProps> = ({
  projects,
  loading = false,
  onEdit,
  onDelete,
  onSelect,
  selectedIds = [],
  selectionMode = false,
  emptyText = '暂无项目',
}) => {
  if (loading) {
    return (
      <div className="flex items-center justify-center py-20">
        <Spin size="large" tip="加载中..." />
      </div>
    )
  }

  if (projects.length === 0) {
    return (
      <div className="py-20">
        <Empty description={emptyText} />
      </div>
    )
  }

  return (
    <Row gutter={[16, 16]}>
      {projects.map((project) => (
        <Col
          key={project.id}
          xs={24}
          sm={24}
          md={12}
          lg={8}
          xl={6}
          xxl={6}
        >
          <ProjectCard
            {...project}
            onEdit={onEdit}
            onDelete={onDelete}
            onSelect={onSelect}
            selected={selectedIds.includes(project.id)}
            selectionMode={selectionMode}
          />
        </Col>
      ))}
    </Row>
  )
}
