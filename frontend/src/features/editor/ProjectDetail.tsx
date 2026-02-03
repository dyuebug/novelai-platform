import { useState, useEffect, useCallback } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import {
  Card,
  Button,
  Space,
  Spin,
  message,
  List,
  Empty,
  Modal,
  Form,
  Input,
  InputNumber,
  Popconfirm,
  Descriptions,
} from 'antd'
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  FileTextOutlined,
  GlobalOutlined,
  ArrowLeftOutlined,
} from '@ant-design/icons'
import { projectService, Project } from '@/services/project.service'
import { chapterService, Chapter } from '@/services/chapter.service'

const { TextArea } = Input

const ProjectDetail = () => {
  const { projectId } = useParams<{ projectId: string }>()
  const navigate = useNavigate()
  const [project, setProject] = useState<Project | null>(null)
  const [chapters, setChapters] = useState<Chapter[]>([])
  const [loading, setLoading] = useState(true)
  const [chapterModalOpen, setChapterModalOpen] = useState(false)
  const [editingChapter, setEditingChapter] = useState<Chapter | null>(null)
  const [form] = Form.useForm()

  const fetchProject = useCallback(async () => {
    if (!projectId) return
    try {
      const data = await projectService.get(projectId)
      setProject(data)
    } catch (error) {
      message.error('获取项目失败')
      console.error(error)
    }
  }, [projectId])

  const fetchChapters = useCallback(async () => {
    if (!projectId) return
    try {
      const response = await chapterService.list(projectId, { page_size: 100 })
      setChapters(response.items || [])
    } catch (error) {
      message.error('获取章节列表失败')
      console.error(error)
    }
  }, [projectId])

  useEffect(() => {
    const loadData = async () => {
      setLoading(true)
      await Promise.all([fetchProject(), fetchChapters()])
      setLoading(false)
    }
    loadData()
  }, [fetchProject, fetchChapters])

  const handleCreateChapter = () => {
    setEditingChapter(null)
    form.resetFields()
    form.setFieldsValue({ chapter_number: chapters.length + 1 })
    setChapterModalOpen(true)
  }

  const handleEditChapter = (chapter: Chapter) => {
    setEditingChapter(chapter)
    form.setFieldsValue(chapter)
    setChapterModalOpen(true)
  }

  const handleDeleteChapter = async (id: string) => {
    try {
      await chapterService.delete(id)
      message.success('删除成功')
      fetchChapters()
    } catch (error) {
      message.error('删除失败')
      console.error(error)
    }
  }

  const handleChapterSubmit = async () => {
    try {
      const values = await form.validateFields()
      if (editingChapter) {
        await chapterService.update(editingChapter.id, values)
        message.success('更新成功')
      } else {
        await chapterService.create(projectId!, values)
        message.success('创建成功')
      }
      setChapterModalOpen(false)
      fetchChapters()
    } catch (error) {
      console.error(error)
    }
  }

  const handleOpenChapter = (chapter: Chapter) => {
    navigate(`/projects/${projectId}/chapters/${chapter.id}`)
  }

  const handleOpenWorld = () => {
    navigate(`/projects/${projectId}/world`)
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Spin size="large" />
      </div>
    )
  }

  if (!project) {
    return (
      <div className="p-6">
        <Empty description="项目不存在" />
      </div>
    )
  }

  return (
    <div className="p-6">
      {/* 顶部导航 */}
      <div className="mb-4">
        <Button icon={<ArrowLeftOutlined />} onClick={() => navigate('/dashboard')}>
          返回项目列表
        </Button>
      </div>

      {/* 项目信息卡片 */}
      <Card className="mb-6">
        <div className="flex justify-between items-start">
          <div>
            <h1 className="text-2xl font-bold mb-2">{project.title}</h1>
            <p className="text-gray-500 mb-4">{project.description || '暂无描述'}</p>
          </div>
          <Space>
            <Button icon={<GlobalOutlined />} onClick={handleOpenWorld}>
              世界观管理
            </Button>
          </Space>
        </div>

        <Descriptions size="small" column={4}>
          <Descriptions.Item label="类型">{project.genre || '-'}</Descriptions.Item>
          <Descriptions.Item label="状态">{project.status}</Descriptions.Item>
          <Descriptions.Item label="章节数">{chapters.length}</Descriptions.Item>
          <Descriptions.Item label="创建时间">
            {new Date(project.created_at).toLocaleDateString()}
          </Descriptions.Item>
        </Descriptions>
      </Card>

      {/* 章节列表 */}
      <Card
        title="章节列表"
        extra={
          <Button type="primary" icon={<PlusOutlined />} onClick={handleCreateChapter}>
            新建章节
          </Button>
        }
      >
        {chapters.length > 0 ? (
          <List
            dataSource={chapters}
            renderItem={(chapter) => (
              <List.Item
                actions={[
                  <Button
                    key="edit"
                    type="link"
                    icon={<EditOutlined />}
                    onClick={() => handleEditChapter(chapter)}
                  >
                    编辑
                  </Button>,
                  <Popconfirm
                    key="delete"
                    title="确定删除此章节？"
                    onConfirm={() => handleDeleteChapter(chapter.id)}
                  >
                    <Button type="link" danger icon={<DeleteOutlined />}>
                      删除
                    </Button>
                  </Popconfirm>,
                ]}
              >
                <List.Item.Meta
                  avatar={<FileTextOutlined className="text-2xl text-blue-500" />}
                  title={
                    <a onClick={() => handleOpenChapter(chapter)}>
                      第 {chapter.chapter_number} 章：{chapter.title}
                    </a>
                  }
                  description={
                    <Space>
                      <span>字数：{chapter.word_count || 0}</span>
                      <span>状态：{chapter.status}</span>
                      <span>
                        更新：{new Date(chapter.updated_at).toLocaleDateString()}
                      </span>
                    </Space>
                  }
                />
              </List.Item>
            )}
          />
        ) : (
          <Empty description="暂无章节，点击上方按钮创建第一章" />
        )}
      </Card>

      {/* 章节编辑弹窗 */}
      <Modal
        title={editingChapter ? '编辑章节' : '新建章节'}
        open={chapterModalOpen}
        onOk={handleChapterSubmit}
        onCancel={() => setChapterModalOpen(false)}
        destroyOnClose
      >
        <Form form={form} layout="vertical" className="mt-4">
          <Form.Item
            name="title"
            label="章节标题"
            rules={[{ required: true, message: '请输入章节标题' }]}
          >
            <Input placeholder="章节标题" />
          </Form.Item>

          <Form.Item
            name="chapter_number"
            label="章节序号"
            rules={[{ required: true, message: '请输入章节序号' }]}
          >
            <InputNumber min={1} className="w-full" />
          </Form.Item>

          <Form.Item name="summary" label="章节概要">
            <TextArea rows={3} placeholder="章节概要" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default ProjectDetail
