import { useState, useEffect, useCallback } from 'react'
import { useParams } from 'react-router-dom'
import {
  Card,
  Tree,
  Button,
  Space,
  Modal,
  Form,
  Input,
  Select,
  message,
  Popconfirm,
  Empty,
  Tabs,
  Typography,
} from 'antd'
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  BookOutlined,
  FileTextOutlined,
} from '@ant-design/icons'
import type { DataNode } from 'antd/es/tree'
import { worldSettingService, WorldSetting } from '@/services/world.service'

const { TextArea } = Input
const { Paragraph } = Typography

// 设定分类
const settingCategories = [
  { value: 'magic', label: '魔法/力量体系' },
  { value: 'technology', label: '科技体系' },
  { value: 'society', label: '社会制度' },
  { value: 'history', label: '历史事件' },
  { value: 'culture', label: '文化习俗' },
  { value: 'geography', label: '地理环境' },
  { value: 'species', label: '种族/物种' },
  { value: 'item', label: '物品/道具' },
  { value: 'rule', label: '世界规则' },
  { value: 'other', label: '其他' },
]

// 分类标签页（包含 key 属性）
const categoryTabs = [
  { key: 'all', label: '全部' },
  { key: 'magic', label: '魔法/力量体系' },
  { key: 'technology', label: '科技体系' },
  { key: 'society', label: '社会制度' },
  { key: 'history', label: '历史事件' },
  { key: 'culture', label: '文化习俗' },
  { key: 'geography', label: '地理环境' },
  { key: 'species', label: '种族/物种' },
  { key: 'item', label: '物品/道具' },
  { key: 'rule', label: '世界规则' },
  { key: 'other', label: '其他' },
]

const WorldSettingManager = () => {
  const { projectId } = useParams<{ projectId: string }>()
  const [settings, setSettings] = useState<WorldSetting[]>([])
  const [loading, setLoading] = useState(true)
  const [modalOpen, setModalOpen] = useState(false)
  const [editingSetting, setEditingSetting] = useState<WorldSetting | null>(null)
  const [selectedSetting, setSelectedSetting] = useState<WorldSetting | null>(null)
  const [parentId, setParentId] = useState<string | undefined>(undefined)
  const [activeCategory, setActiveCategory] = useState<string>('all')
  const [form] = Form.useForm()

  const fetchSettings = useCallback(async () => {
    if (!projectId) return
    setLoading(true)
    try {
      const response = await worldSettingService.list(projectId, { page_size: 500 })
      setSettings(response.items || [])
    } catch (error) {
      message.error('获取世界设定失败')
      console.error(error)
    } finally {
      setLoading(false)
    }
  }, [projectId])

  useEffect(() => {
    fetchSettings()
  }, [fetchSettings])

  // 按分类过滤设定
  const filteredSettings = activeCategory === 'all'
    ? settings
    : settings.filter((s) => s.category === activeCategory)

  // 构建树形数据
  const buildTreeData = (items: WorldSetting[], parentId?: string): DataNode[] => {
    return items
      .filter((item) => item.parent_id === parentId)
      .map((item) => ({
        key: item.id,
        title: (
          <Space>
            <FileTextOutlined />
            <span>{item.title}</span>
            <span className="text-gray-400 text-xs">
              ({settingCategories.find((c) => c.value === item.category)?.label || item.category})
            </span>
          </Space>
        ),
        children: buildTreeData(items, item.id),
        data: item,
      }))
  }

  const treeData = buildTreeData(filteredSettings, undefined)

  const handleCreate = (pId?: string, category?: string) => {
    setEditingSetting(null)
    setParentId(pId)
    form.resetFields()
    if (category && category !== 'all') {
      form.setFieldsValue({ category })
    }
    setModalOpen(true)
  }

  const handleEdit = (setting: WorldSetting) => {
    setEditingSetting(setting)
    setParentId(setting.parent_id)
    form.setFieldsValue(setting)
    setModalOpen(true)
  }

  const handleDelete = async (id: string) => {
    try {
      await worldSettingService.delete(id)
      message.success('删除成功')
      setSelectedSetting(null)
      fetchSettings()
    } catch (error) {
      message.error('删除失败')
      console.error(error)
    }
  }

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      const data = { ...values, parent_id: parentId }

      if (editingSetting) {
        await worldSettingService.update(editingSetting.id, data)
        message.success('更新成功')
      } else {
        await worldSettingService.create(projectId!, data)
        message.success('创建成功')
      }
      setModalOpen(false)
      fetchSettings()
    } catch (error) {
      console.error(error)
    }
  }

  const handleSelect = (selectedKeys: React.Key[]) => {
    if (selectedKeys.length > 0) {
      const setting = settings.find((s) => s.id === selectedKeys[0])
      setSelectedSetting(setting || null)
    } else {
      setSelectedSetting(null)
    }
  }

  // 分类标签页
  const tabItems = categoryTabs.map((cat) => ({
    key: cat.key,
    label: cat.label,
  }))

  return (
    <div className="flex gap-4">
      {/* 左侧分类和树形结构 */}
      <Card
        title={
          <Space>
            <BookOutlined />
            <span>世界设定</span>
          </Space>
        }
        className="w-1/3"
        extra={
          <Button
            type="primary"
            icon={<PlusOutlined />}
            size="small"
            onClick={() => handleCreate(undefined, activeCategory)}
          >
            新建
          </Button>
        }
        loading={loading}
      >
        <Tabs
          activeKey={activeCategory}
          onChange={setActiveCategory}
          size="small"
          items={tabItems}
        />

        {treeData.length > 0 ? (
          <Tree
            treeData={treeData}
            onSelect={handleSelect}
            defaultExpandAll
            showLine
          />
        ) : (
          <Empty description="暂无设定" />
        )}
      </Card>

      {/* 右侧详情 */}
      <Card title="设定详情" className="flex-1">
        {selectedSetting ? (
          <div>
            <div className="flex justify-end mb-4">
              <Space>
                <Button
                  icon={<PlusOutlined />}
                  onClick={() => handleCreate(selectedSetting.id, selectedSetting.category)}
                >
                  添加子设定
                </Button>
                <Button icon={<EditOutlined />} onClick={() => handleEdit(selectedSetting)}>
                  编辑
                </Button>
                <Popconfirm
                  title="确定删除此设定？"
                  onConfirm={() => handleDelete(selectedSetting.id)}
                >
                  <Button danger icon={<DeleteOutlined />}>
                    删除
                  </Button>
                </Popconfirm>
              </Space>
            </div>

            <div className="space-y-4">
              <div>
                <h3 className="text-lg font-medium mb-2">{selectedSetting.title}</h3>
                <span className="inline-block px-2 py-1 bg-blue-100 text-blue-800 rounded text-sm">
                  {settingCategories.find((c) => c.value === selectedSetting.category)?.label ||
                    selectedSetting.category}
                </span>
              </div>

              <div>
                <h4 className="text-sm font-medium text-gray-500 mb-1">内容</h4>
                <Paragraph className="whitespace-pre-wrap bg-gray-50 p-4 rounded">
                  {selectedSetting.content || '暂无内容'}
                </Paragraph>
              </div>

              <div className="text-sm text-gray-400">
                创建时间：{new Date(selectedSetting.created_at).toLocaleString()}
              </div>
            </div>
          </div>
        ) : (
          <Empty description="请选择一个设定查看详情" />
        )}
      </Card>

      {/* 创建/编辑弹窗 */}
      <Modal
        title={editingSetting ? '编辑设定' : '新建设定'}
        open={modalOpen}
        onOk={handleSubmit}
        onCancel={() => setModalOpen(false)}
        width={700}
        destroyOnClose
      >
        <Form form={form} layout="vertical" className="mt-4">
          <div className="grid grid-cols-2 gap-4">
            <Form.Item
              name="title"
              label="设定标题"
              rules={[{ required: true, message: '请输入设定标题' }]}
            >
              <Input placeholder="设定标题" />
            </Form.Item>

            <Form.Item
              name="category"
              label="分类"
              rules={[{ required: true, message: '请选择分类' }]}
            >
              <Select placeholder="选择分类" options={settingCategories} />
            </Form.Item>
          </div>

          <Form.Item name="content" label="内容">
            <TextArea rows={10} placeholder="详细描述这个世界设定..." />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default WorldSettingManager
