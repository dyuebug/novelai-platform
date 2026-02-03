import { useState, useEffect, useCallback } from 'react'
import { useParams } from 'react-router-dom'
import {
  Card,
  Table,
  Button,
  Space,
  Modal,
  Form,
  Input,
  Select,
  Tag,
  message,
  Popconfirm,
  Avatar,
  Tabs,
  List,
  Empty,
} from 'antd'
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  UserOutlined,
} from '@ant-design/icons'
import {
  characterService,
  Character,
  CharacterRelationship,
  CharacterExperience,
} from '@/services/world.service'

const { TextArea } = Input

const roleOptions = [
  { value: 'protagonist', label: '主角' },
  { value: 'antagonist', label: '反派' },
  { value: 'supporting', label: '配角' },
  { value: 'minor', label: '龙套' },
]

const roleColors: Record<string, string> = {
  protagonist: 'gold',
  antagonist: 'red',
  supporting: 'blue',
  minor: 'default',
}

const relationTypes = [
  { value: 'family', label: '家人' },
  { value: 'friend', label: '朋友' },
  { value: 'enemy', label: '敌人' },
  { value: 'lover', label: '恋人' },
  { value: 'mentor', label: '师徒' },
  { value: 'colleague', label: '同事' },
  { value: 'rival', label: '对手' },
]

const CharacterManager = () => {
  const { projectId } = useParams<{ projectId: string }>()
  const [characters, setCharacters] = useState<Character[]>([])
  const [loading, setLoading] = useState(true)
  const [modalOpen, setModalOpen] = useState(false)
  const [editingCharacter, setEditingCharacter] = useState<Character | null>(null)
  const [detailCharacter, setDetailCharacter] = useState<Character | null>(null)
  const [relationships, setRelationships] = useState<CharacterRelationship[]>([])
  const [experiences, setExperiences] = useState<CharacterExperience[]>([])
  const [form] = Form.useForm()

  const fetchCharacters = useCallback(async () => {
    if (!projectId) return
    setLoading(true)
    try {
      const response = await characterService.list(projectId, { page_size: 100 })
      setCharacters(response.items || [])
    } catch (error) {
      message.error('获取角色列表失败')
      console.error(error)
    } finally {
      setLoading(false)
    }
  }, [projectId])

  useEffect(() => {
    fetchCharacters()
  }, [fetchCharacters])

  const handleCreate = () => {
    setEditingCharacter(null)
    form.resetFields()
    setModalOpen(true)
  }

  const handleEdit = (character: Character) => {
    setEditingCharacter(character)
    form.setFieldsValue(character)
    setModalOpen(true)
  }

  const handleDelete = async (id: string) => {
    try {
      await characterService.delete(id)
      message.success('删除成功')
      fetchCharacters()
    } catch (error) {
      message.error('删除失败')
      console.error(error)
    }
  }

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      if (editingCharacter) {
        await characterService.update(editingCharacter.id, values)
        message.success('更新成功')
      } else {
        await characterService.create(projectId!, values)
        message.success('创建成功')
      }
      setModalOpen(false)
      fetchCharacters()
    } catch (error) {
      console.error(error)
    }
  }

  const handleViewDetail = async (character: Character) => {
    setDetailCharacter(character)
    try {
      const [rels, exps] = await Promise.all([
        characterService.getRelationships(character.id),
        characterService.getExperiences(character.id),
      ])
      setRelationships(rels || [])
      setExperiences(exps || [])
    } catch (error) {
      console.error(error)
    }
  }

  const columns = [
    {
      title: '角色',
      dataIndex: 'name',
      key: 'name',
      render: (name: string, record: Character) => (
        <Space>
          <Avatar icon={<UserOutlined />} src={record.avatar_url} />
          <span className="font-medium">{name}</span>
          {record.alias && <span className="text-gray-400">({record.alias})</span>}
        </Space>
      ),
    },
    {
      title: '类型',
      dataIndex: 'role',
      key: 'role',
      width: 100,
      render: (role: string) => (
        <Tag color={roleColors[role] || 'default'}>
          {roleOptions.find((r) => r.value === role)?.label || role || '未设置'}
        </Tag>
      ),
    },
    {
      title: '性别',
      dataIndex: 'gender',
      key: 'gender',
      width: 80,
    },
    {
      title: '年龄',
      dataIndex: 'age',
      key: 'age',
      width: 80,
    },
    {
      title: '简介',
      dataIndex: 'background',
      key: 'background',
      ellipsis: true,
      render: (text: string) => text || '-',
    },
    {
      title: '操作',
      key: 'action',
      width: 180,
      render: (_: unknown, record: Character) => (
        <Space>
          <Button type="link" size="small" onClick={() => handleViewDetail(record)}>
            详情
          </Button>
          <Button type="link" size="small" icon={<EditOutlined />} onClick={() => handleEdit(record)}>
            编辑
          </Button>
          <Popconfirm title="确定删除此角色？" onConfirm={() => handleDelete(record.id)}>
            <Button type="link" size="small" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ]

  return (
    <div>
      <Card
        title="角色管理"
        extra={
          <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
            新建角色
          </Button>
        }
      >
        <Table
          columns={columns}
          dataSource={characters}
          rowKey="id"
          loading={loading}
          pagination={{ pageSize: 10 }}
        />
      </Card>

      {/* 创建/编辑弹窗 */}
      <Modal
        title={editingCharacter ? '编辑角色' : '新建角色'}
        open={modalOpen}
        onOk={handleSubmit}
        onCancel={() => setModalOpen(false)}
        width={700}
        destroyOnClose
      >
        <Form form={form} layout="vertical" className="mt-4">
          <div className="grid grid-cols-2 gap-4">
            <Form.Item
              name="name"
              label="角色名"
              rules={[{ required: true, message: '请输入角色名' }]}
            >
              <Input placeholder="角色名称" />
            </Form.Item>
            <Form.Item name="alias" label="别名/外号">
              <Input placeholder="别名或外号" />
            </Form.Item>
          </div>

          <div className="grid grid-cols-3 gap-4">
            <Form.Item name="gender" label="性别">
              <Select placeholder="选择性别" allowClear>
                <Select.Option value="男">男</Select.Option>
                <Select.Option value="女">女</Select.Option>
                <Select.Option value="其他">其他</Select.Option>
              </Select>
            </Form.Item>
            <Form.Item name="age" label="年龄">
              <Input placeholder="年龄" />
            </Form.Item>
            <Form.Item name="role" label="角色类型">
              <Select placeholder="选择类型" options={roleOptions} allowClear />
            </Form.Item>
          </div>

          <Form.Item name="appearance" label="外貌描述">
            <TextArea rows={2} placeholder="外貌特征描述" />
          </Form.Item>

          <Form.Item name="personality" label="性格特点">
            <TextArea rows={2} placeholder="性格特点描述" />
          </Form.Item>

          <Form.Item name="background" label="背景故事">
            <TextArea rows={3} placeholder="角色背景故事" />
          </Form.Item>

          <Form.Item name="abilities" label="能力/技能">
            <TextArea rows={2} placeholder="角色能力或技能" />
          </Form.Item>

          <Form.Item name="goals" label="目标/动机">
            <TextArea rows={2} placeholder="角色的目标或动机" />
          </Form.Item>
        </Form>
      </Modal>

      {/* 详情弹窗 */}
      <Modal
        title={detailCharacter?.name}
        open={!!detailCharacter}
        onCancel={() => setDetailCharacter(null)}
        footer={null}
        width={800}
      >
        {detailCharacter && (
          <Tabs
            items={[
              {
                key: 'info',
                label: '基本信息',
                children: (
                  <div className="space-y-3">
                    <p><strong>别名：</strong>{detailCharacter.alias || '-'}</p>
                    <p><strong>性别：</strong>{detailCharacter.gender || '-'}</p>
                    <p><strong>年龄：</strong>{detailCharacter.age || '-'}</p>
                    <p><strong>外貌：</strong>{detailCharacter.appearance || '-'}</p>
                    <p><strong>性格：</strong>{detailCharacter.personality || '-'}</p>
                    <p><strong>背景：</strong>{detailCharacter.background || '-'}</p>
                    <p><strong>能力：</strong>{detailCharacter.abilities || '-'}</p>
                    <p><strong>目标：</strong>{detailCharacter.goals || '-'}</p>
                  </div>
                ),
              },
              {
                key: 'relationships',
                label: `关系 (${relationships.length})`,
                children: relationships.length > 0 ? (
                  <List
                    dataSource={relationships}
                    renderItem={(rel) => (
                      <List.Item>
                        <Space>
                          <Tag>{relationTypes.find((t) => t.value === rel.relation_type)?.label || rel.relation_type}</Tag>
                          <span>{rel.target?.name || '未知角色'}</span>
                          {rel.description && <span className="text-gray-400">- {rel.description}</span>}
                        </Space>
                      </List.Item>
                    )}
                  />
                ) : (
                  <Empty description="暂无关系" />
                ),
              },
              {
                key: 'experiences',
                label: `经历 (${experiences.length})`,
                children: experiences.length > 0 ? (
                  <List
                    dataSource={experiences}
                    renderItem={(exp) => (
                      <List.Item>
                        <List.Item.Meta
                          title={exp.title}
                          description={
                            <div>
                              {exp.time_point && <Tag>{exp.time_point}</Tag>}
                              {exp.chapter_ref && <Tag>第{exp.chapter_ref}章</Tag>}
                              <p className="mt-1">{exp.content}</p>
                            </div>
                          }
                        />
                      </List.Item>
                    )}
                  />
                ) : (
                  <Empty description="暂无经历" />
                ),
              },
            ]}
          />
        )}
      </Modal>
    </div>
  )
}

export default CharacterManager
