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
  Descriptions,
} from 'antd'
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  EnvironmentOutlined,
} from '@ant-design/icons'
import type { DataNode } from 'antd/es/tree'
import { locationService, Location } from '@/services/world.service'

const { TextArea } = Input

const locationTypes = [
  { value: 'world', label: '世界' },
  { value: 'continent', label: '大陆' },
  { value: 'country', label: '国家' },
  { value: 'region', label: '地区' },
  { value: 'city', label: '城市' },
  { value: 'building', label: '建筑' },
  { value: 'room', label: '房间' },
  { value: 'other', label: '其他' },
]

const LocationManager = () => {
  const { projectId } = useParams<{ projectId: string }>()
  const [locations, setLocations] = useState<Location[]>([])
  const [loading, setLoading] = useState(true)
  const [modalOpen, setModalOpen] = useState(false)
  const [editingLocation, setEditingLocation] = useState<Location | null>(null)
  const [selectedLocation, setSelectedLocation] = useState<Location | null>(null)
  const [parentId, setParentId] = useState<string | undefined>(undefined)
  const [form] = Form.useForm()

  const fetchLocations = useCallback(async () => {
    if (!projectId) return
    setLoading(true)
    try {
      const response = await locationService.list(projectId, { page_size: 200 })
      setLocations(response.items || [])
    } catch (error) {
      message.error('获取地点列表失败')
      console.error(error)
    } finally {
      setLoading(false)
    }
  }, [projectId])

  useEffect(() => {
    fetchLocations()
  }, [fetchLocations])

  // 构建树形数据
  const buildTreeData = (items: Location[], parentId?: string): DataNode[] => {
    return items
      .filter((item) => item.parent_id === parentId)
      .map((item) => ({
        key: item.id,
        title: (
          <Space>
            <EnvironmentOutlined />
            <span>{item.name}</span>
            {item.type && <span className="text-gray-400 text-xs">({locationTypes.find((t) => t.value === item.type)?.label || item.type})</span>}
          </Space>
        ),
        children: buildTreeData(items, item.id),
        data: item,
      }))
  }

  const treeData = buildTreeData(locations, undefined)

  const handleCreate = (pId?: string) => {
    setEditingLocation(null)
    setParentId(pId)
    form.resetFields()
    setModalOpen(true)
  }

  const handleEdit = (location: Location) => {
    setEditingLocation(location)
    setParentId(location.parent_id)
    form.setFieldsValue(location)
    setModalOpen(true)
  }

  const handleDelete = async (id: string) => {
    try {
      await locationService.delete(id)
      message.success('删除成功')
      setSelectedLocation(null)
      fetchLocations()
    } catch (error) {
      message.error('删除失败')
      console.error(error)
    }
  }

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      const data = { ...values, parent_id: parentId }

      if (editingLocation) {
        await locationService.update(editingLocation.id, data)
        message.success('更新成功')
      } else {
        await locationService.create(projectId!, data)
        message.success('创建成功')
      }
      setModalOpen(false)
      fetchLocations()
    } catch (error) {
      console.error(error)
    }
  }

  const handleSelect = (selectedKeys: React.Key[]) => {
    if (selectedKeys.length > 0) {
      const location = locations.find((l) => l.id === selectedKeys[0])
      setSelectedLocation(location || null)
    } else {
      setSelectedLocation(null)
    }
  }

  return (
    <div className="flex gap-4">
      {/* 左侧树形结构 */}
      <Card
        title="地点列表"
        className="w-1/3"
        extra={
          <Button type="primary" icon={<PlusOutlined />} size="small" onClick={() => handleCreate()}>
            新建
          </Button>
        }
        loading={loading}
      >
        {treeData.length > 0 ? (
          <Tree
            treeData={treeData}
            onSelect={handleSelect}
            defaultExpandAll
            showLine
          />
        ) : (
          <Empty description="暂无地点" />
        )}
      </Card>

      {/* 右侧详情 */}
      <Card title="地点详情" className="flex-1">
        {selectedLocation ? (
          <div>
            <div className="flex justify-end mb-4">
              <Space>
                <Button icon={<PlusOutlined />} onClick={() => handleCreate(selectedLocation.id)}>
                  添加子地点
                </Button>
                <Button icon={<EditOutlined />} onClick={() => handleEdit(selectedLocation)}>
                  编辑
                </Button>
                <Popconfirm title="确定删除此地点？" onConfirm={() => handleDelete(selectedLocation.id)}>
                  <Button danger icon={<DeleteOutlined />}>
                    删除
                  </Button>
                </Popconfirm>
              </Space>
            </div>

            <Descriptions column={2} bordered>
              <Descriptions.Item label="名称">{selectedLocation.name}</Descriptions.Item>
              <Descriptions.Item label="类型">
                {locationTypes.find((t) => t.value === selectedLocation.type)?.label || selectedLocation.type || '-'}
              </Descriptions.Item>
              <Descriptions.Item label="气候">{selectedLocation.climate || '-'}</Descriptions.Item>
              <Descriptions.Item label="创建时间">
                {new Date(selectedLocation.created_at).toLocaleDateString()}
              </Descriptions.Item>
              <Descriptions.Item label="描述" span={2}>
                {selectedLocation.description || '-'}
              </Descriptions.Item>
              <Descriptions.Item label="特征" span={2}>
                {selectedLocation.features || '-'}
              </Descriptions.Item>
            </Descriptions>
          </div>
        ) : (
          <Empty description="请选择一个地点查看详情" />
        )}
      </Card>

      {/* 创建/编辑弹窗 */}
      <Modal
        title={editingLocation ? '编辑地点' : '新建地点'}
        open={modalOpen}
        onOk={handleSubmit}
        onCancel={() => setModalOpen(false)}
        destroyOnClose
      >
        <Form form={form} layout="vertical" className="mt-4">
          <Form.Item
            name="name"
            label="地点名称"
            rules={[{ required: true, message: '请输入地点名称' }]}
          >
            <Input placeholder="地点名称" />
          </Form.Item>

          <Form.Item name="type" label="地点类型">
            <Select placeholder="选择类型" options={locationTypes} allowClear />
          </Form.Item>

          <Form.Item name="climate" label="气候">
            <Input placeholder="气候特征" />
          </Form.Item>

          <Form.Item name="description" label="描述">
            <TextArea rows={3} placeholder="地点描述" />
          </Form.Item>

          <Form.Item name="features" label="特征">
            <TextArea rows={3} placeholder="地点特征" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default LocationManager
