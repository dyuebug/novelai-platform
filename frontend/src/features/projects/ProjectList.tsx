import { useState, useEffect, useCallback } from 'react'
import { Card, List, Tag, Button, Dropdown, Empty, Spin, message } from 'antd'
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  MoreOutlined,
  BookOutlined,
  FileTextOutlined,
} from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { projectService, Project, ProjectListParams } from '@/services/project.service'

interface ProjectListProps {
  onCreateClick: () => void
}

const statusColors: Record<string, string> = {
  draft: 'default',
  writing: 'processing',
  completed: 'success',
  paused: 'warning',
}

const statusLabels: Record<string, string> = {
  draft: '草稿',
  writing: '写作中',
  completed: '已完成',
  paused: '已暂停',
}

const ProjectList = ({ onCreateClick }: ProjectListProps) => {
  const navigate = useNavigate()
  const [projects, setProjects] = useState<Project[]>([])
  const [loading, setLoading] = useState(true)
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 12,
    total: 0,
  })

  const fetchProjects = useCallback(async (params?: ProjectListParams) => {
    setLoading(true)
    try {
      const response = await projectService.list({
        page: params?.page || pagination.current,
        page_size: params?.page_size || pagination.pageSize,
        ...params,
      })
      setProjects(response.items || [])
      setPagination((prev) => ({
        ...prev,
        current: response.page,
        total: response.total,
      }))
    } catch (error) {
      message.error('获取项目列表失败')
      console.error(error)
    } finally {
      setLoading(false)
    }
  }, [pagination.current, pagination.pageSize])

  useEffect(() => {
    fetchProjects()
  }, [])

  const handleDelete = async (id: string) => {
    try {
      await projectService.delete(id)
      message.success('项目已删除')
      fetchProjects()
    } catch (error) {
      message.error('删除失败')
      console.error(error)
    }
  }

  const handlePageChange = (page: number, pageSize: number) => {
    fetchProjects({ page, page_size: pageSize })
  }

  const getDropdownItems = (project: Project) => [
    {
      key: 'edit',
      icon: <EditOutlined />,
      label: '编辑信息',
      onClick: () => navigate(`/projects/${project.id}/settings`),
    },
    {
      key: 'delete',
      icon: <DeleteOutlined />,
      label: '删除项目',
      danger: true,
      onClick: () => handleDelete(project.id),
    },
  ]

  if (loading && projects.length === 0) {
    return (
      <div className="flex justify-center items-center h-64">
        <Spin size="large" />
      </div>
    )
  }

  if (!loading && projects.length === 0) {
    return (
      <Card>
        <Empty
          description="暂无项目"
          image={Empty.PRESENTED_IMAGE_SIMPLE}
        >
          <Button type="primary" icon={<PlusOutlined />} onClick={onCreateClick}>
            创建第一个项目
          </Button>
        </Empty>
      </Card>
    )
  }

  return (
    <List
      grid={{
        gutter: 16,
        xs: 1,
        sm: 2,
        md: 2,
        lg: 3,
        xl: 4,
        xxl: 4,
      }}
      dataSource={projects}
      loading={loading}
      pagination={{
        ...pagination,
        onChange: handlePageChange,
        showSizeChanger: true,
        showTotal: (total) => `共 ${total} 个项目`,
      }}
      renderItem={(project) => (
        <List.Item>
          <Card
            hoverable
            className="h-full"
            onClick={() => navigate(`/projects/${project.id}`)}
            actions={[
              <Button
                key="open"
                type="link"
                icon={<BookOutlined />}
                onClick={(e) => {
                  e.stopPropagation()
                  navigate(`/projects/${project.id}`)
                }}
              >
                打开
              </Button>,
              <Dropdown
                key="more"
                menu={{ items: getDropdownItems(project) }}
                trigger={['click']}
              >
                <Button
                  type="link"
                  icon={<MoreOutlined />}
                  onClick={(e) => e.stopPropagation()}
                />
              </Dropdown>,
            ]}
          >
            <Card.Meta
              title={
                <div className="flex items-center justify-between">
                  <span className="truncate">{project.title}</span>
                  <Tag color={statusColors[project.status] || 'default'}>
                    {statusLabels[project.status] || project.status}
                  </Tag>
                </div>
              }
              description={
                <div className="space-y-2">
                  <p className="text-gray-500 line-clamp-2 h-10">
                    {project.description || '暂无简介'}
                  </p>
                  <div className="flex items-center gap-4 text-xs text-gray-400">
                    <span className="flex items-center gap-1">
                      <FileTextOutlined />
                      {project.total_chapters} 章
                    </span>
                    <span>{project.total_words.toLocaleString()} 字</span>
                  </div>
                  {project.genre && (
                    <Tag color="blue" className="mt-1">
                      {project.genre}
                    </Tag>
                  )}
                </div>
              }
            />
          </Card>
        </List.Item>
      )}
    />
  )
}

export default ProjectList
