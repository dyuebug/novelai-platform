import { useState } from 'react'
import { Modal, Form, Input, Select, message } from 'antd'
import { projectService, CreateProjectRequest, Project } from '@/services/project.service'

const { TextArea } = Input

interface CreateProjectModalProps {
  open: boolean
  onClose: () => void
  onSuccess: (project: Project) => void
}

const genreOptions = [
  { value: '玄幻', label: '玄幻' },
  { value: '仙侠', label: '仙侠' },
  { value: '都市', label: '都市' },
  { value: '历史', label: '历史' },
  { value: '科幻', label: '科幻' },
  { value: '游戏', label: '游戏' },
  { value: '悬疑', label: '悬疑' },
  { value: '奇幻', label: '奇幻' },
  { value: '武侠', label: '武侠' },
  { value: '其他', label: '其他' },
]

const CreateProjectModal = ({ open, onClose, onSuccess }: CreateProjectModalProps) => {
  const [form] = Form.useForm()
  const [loading, setLoading] = useState(false)

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      setLoading(true)

      const data: CreateProjectRequest = {
        title: values.title,
        description: values.description,
        genre: values.genre,
      }

      const project = await projectService.create(data)
      message.success('项目创建成功')
      form.resetFields()
      onSuccess(project)
    } catch (error) {
      if (error instanceof Error) {
        message.error(error.message || '创建失败')
      }
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  const handleCancel = () => {
    form.resetFields()
    onClose()
  }

  return (
    <Modal
      title="新建项目"
      open={open}
      onOk={handleSubmit}
      onCancel={handleCancel}
      confirmLoading={loading}
      okText="创建"
      cancelText="取消"
      destroyOnClose
    >
      <Form
        form={form}
        layout="vertical"
        className="mt-4"
      >
        <Form.Item
          name="title"
          label="项目名称"
          rules={[
            { required: true, message: '请输入项目名称' },
            { min: 1, max: 200, message: '名称长度应在 1-200 字符之间' },
          ]}
        >
          <Input placeholder="请输入小说名称" maxLength={200} showCount />
        </Form.Item>

        <Form.Item
          name="genre"
          label="类型"
        >
          <Select
            placeholder="请选择类型"
            options={genreOptions}
            allowClear
          />
        </Form.Item>

        <Form.Item
          name="description"
          label="简介"
          rules={[
            { max: 5000, message: '简介不能超过 5000 字符' },
          ]}
        >
          <TextArea
            placeholder="请输入小说简介（选填）"
            rows={4}
            maxLength={5000}
            showCount
          />
        </Form.Item>
      </Form>
    </Modal>
  )
}

export default CreateProjectModal
