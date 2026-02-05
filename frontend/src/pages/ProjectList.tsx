import { useState, useEffect, useMemo } from 'react'
import { Button, Space, Pagination, Modal, Form, Input, Select, message } from 'antd'
import { PlusOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { ProjectGrid } from '@/components/project/ProjectGrid'
import { ProjectFilters } from '@/components/project/ProjectFilters'
import { EmptyState } from '@/components/project/EmptyState'
import type { Dayjs } from 'dayjs'

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
  createdAt?: string
}

interface FilterValues {
  search?: string
  tags?: string[]
  status?: string[]
  dateRange?: [Dayjs | null, Dayjs | null] | null
}

type SortField = 'createdAt' | 'lastUpdateTime' | 'title'
type SortOrder = 'asc' | 'desc'

export default function ProjectList() {
  const navigate = useNavigate()
  const [projects, setProjects] = useState<Project[]>([])
  const [loading, setLoading] = useState(false)
  const [filters, setFilters] = useState<FilterValues>({})
  const [sortField, setSortField] = useState<SortField>('lastUpdateTime')
  const [sortOrder, setSortOrder] = useState<SortOrder>('desc')
  const [currentPage, setCurrentPage] = useState(1)
  const [pageSize, setPageSize] = useState(20)
  const [total, setTotal] = useState(0)
  const [selectedIds, setSelectedIds] = useState<string[]>([])
  const [selectionMode, setSelectionMode] = useState(false)
  const [editModalVisible, setEditModalVisible] = useState(false)
  const [form] = Form.useForm()

  // 获取所有可用标签
  const availableTags = useMemo(() => {
    const tagsSet = new Set<string>()
    projects.forEach(project => {
      project.tags?.forEach(tag => tagsSet.add(tag))
    })
    return Array.from(tagsSet)
  }, [projects])

  // 加载项目列表
  const loadProjects = async () => {
    setLoading(true)
    try {
      // TODO: 调用实际API
      // const response = await projectService.getProjects({
      //   page: currentPage,
      //   pageSize,
      //   ...filters,
      //   sortField,
      //   sortOrder,
      // })
      // setProjects(response.data)
      // setTotal(response.total)

      // 模拟数据
      await new Promise(resolve => setTimeout(resolve, 500))
      setProjects([])
      setTotal(0)
    } catch (error) {
      message.error('加载项目列表失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadProjects()
  }, [currentPage, pageSize, filters, sortField, sortOrder])

  const handleCreateProject = () => {
    navigate('/projects/new')
  }

  const handleEditProject = (id: string) => {
    const project = projects.find(p => p.id === id)
    if (project) {
      form.setFieldsValue(project)
      setEditModalVisible(true)
    }
  }

  const handleDeleteProject = async (_id: string) => {
    try {
      // TODO: 调用实际API
      // await projectService.deleteProject(_id)
      message.success('项目已删除')
      loadProjects()
    } catch (error) {
      message.error('删除项目失败')
    }
  }

  const handleSaveEdit = async () => {
    try {
      await form.validateFields()
      // TODO: 调用实际API
      // const values = await form.validateFields()
      // await projectService.updateProject(editingProject!.id, values)
      message.success('项目已更新')
      setEditModalVisible(false)
      loadProjects()
    } catch (error) {
      message.error('更新项目失败')
    }
  }

  const handleSelect = (id: string, selected: boolean) => {
    if (selected) {
      setSelectedIds([...selectedIds, id])
    } else {
      setSelectedIds(selectedIds.filter(sid => sid !== id))
    }
  }

  const handleSelectAll = () => {
    if (selectedIds.length === projects.length) {
      setSelectedIds([])
    } else {
      setSelectedIds(projects.map(p => p.id))
    }
  }

  const handleBatchDelete = () => {
    Modal.confirm({
      title: '批量删除',
      content: `确定要删除选中的 ${selectedIds.length} 个项目吗？`,
      okText: '删除',
      okType: 'danger',
      cancelText: '取消',
      onOk: async () => {
        try {
          // TODO: 调用实际API
          // await projectService.batchDelete(selectedIds)
          message.success(`已删除 ${selectedIds.length} 个项目`)
          setSelectedIds([])
          setSelectionMode(false)
          loadProjects()
        } catch (error) {
          message.error('批量删除失败')
        }
      },
    })
  }

  const handleFiltersChange = (newFilters: FilterValues) => {
    setFilters(newFilters)
    setCurrentPage(1)
  }

  const handleSortChange = (field: SortField) => {
    if (sortField === field) {
      setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc')
    } else {
      setSortField(field)
      setSortOrder('desc')
    }
  }

  return (
    <div className="p-6 space-y-6">
      {/* 顶部操作栏 */}
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">我的项目</h1>
        <Space>
          {selectionMode && (
            <>
              <Button onClick={handleSelectAll}>
                {selectedIds.length === projects.length ? '取消全选' : '全选'}
              </Button>
              <Button danger disabled={selectedIds.length === 0} onClick={handleBatchDelete}>
                删除选中 ({selectedIds.length})
              </Button>
              <Button onClick={() => { setSelectionMode(false); setSelectedIds([]) }}>
                退出批量操作
              </Button>
            </>
          )}
          {!selectionMode && (
            <>
              <Button onClick={() => setSelectionMode(true)}>批量操作</Button>
              <Select
                value={sortField}
                onChange={handleSortChange}
                style={{ width: 150 }}
                options={[
                  { label: '最后更新', value: 'lastUpdateTime' },
                  { label: '创建时间', value: 'createdAt' },
                  { label: '标题', value: 'title' },
                ]}
              />
              <Button type="primary" icon={<PlusOutlined />} onClick={handleCreateProject}>
                创建项目
              </Button>
            </>
          )}
        </Space>
      </div>

      {/* 筛选器 */}
      <ProjectFilters
        value={filters}
        onChange={handleFiltersChange}
        availableTags={availableTags}
        onReset={() => setFilters({})}
      />

      {/* 项目网格 */}
      {projects.length === 0 && !loading ? (
        <EmptyState onAction={handleCreateProject} />
      ) : (
        <>
          <ProjectGrid
            projects={projects}
            loading={loading}
            onEdit={handleEditProject}
            onDelete={handleDeleteProject}
            onSelect={handleSelect}
            selectedIds={selectedIds}
            selectionMode={selectionMode}
          />

          {/* 分页 */}
          {total > 0 && (
            <div className="flex justify-center mt-6">
              <Pagination
                current={currentPage}
                pageSize={pageSize}
                total={total}
                onChange={(page, size) => {
                  setCurrentPage(page)
                  setPageSize(size)
                }}
                showSizeChanger
                showQuickJumper
                showTotal={(total) => `共 ${total} 个项目`}
              />
            </div>
          )}
        </>
      )}

      {/* 编辑项目对话框 */}
      <Modal
        title="编辑项目"
        open={editModalVisible}
        onOk={handleSaveEdit}
        onCancel={() => setEditModalVisible(false)}
        okText="保存"
        cancelText="取消"
      >
        <Form form={form} layout="vertical">
          <Form.Item
            label="项目标题"
            name="title"
            rules={[{ required: true, message: '请输入项目标题' }]}
          >
            <Input placeholder="请输入项目标题" />
          </Form.Item>
          <Form.Item label="项目描述" name="description">
            <Input.TextArea rows={4} placeholder="请输入项目描述" />
          </Form.Item>
          <Form.Item label="标签" name="tags">
            <Select mode="tags" placeholder="添加标签" />
          </Form.Item>
          <Form.Item label="状态" name="status">
            <Select
              options={[
                { label: '草稿', value: 'draft' },
                { label: '进行中', value: 'in_progress' },
                { label: '已完成', value: 'completed' },
              ]}
            />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
