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
  Tabs,
  List,
  Empty,
  Avatar,
} from 'antd'
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  TeamOutlined,
  UserOutlined,
} from '@ant-design/icons'
import {
  organizationService,
  characterService,
  Organization,
  OrganizationMember,
  Character,
} from '@/services/world.service'

const { TextArea } = Input

const orgTypes = [
  { value: 'sect', label: '宗门' },
  { value: 'guild', label: '公会' },
  { value: 'kingdom', label: '王国' },
  { value: 'empire', label: '帝国' },
  { value: 'company', label: '公司' },
  { value: 'family', label: '家族' },
  { value: 'gang', label: '帮派' },
  { value: 'army', label: '军队' },
  { value: 'other', label: '其他' },
]

const OrganizationManager = () => {
  const { projectId } = useParams<{ projectId: string }>()
  const [organizations, setOrganizations] = useState<Organization[]>([])
  const [characters, setCharacters] = useState<Character[]>([])
  const [loading, setLoading] = useState(true)
  const [modalOpen, setModalOpen] = useState(false)
  const [memberModalOpen, setMemberModalOpen] = useState(false)
  const [editingOrg, setEditingOrg] = useState<Organization | null>(null)
  const [selectedOrg, setSelectedOrg] = useState<Organization | null>(null)
  const [members, setMembers] = useState<OrganizationMember[]>([])
  const [form] = Form.useForm()
  const [memberForm] = Form.useForm()

  const fetchOrganizations = useCallback(async () => {
    if (!projectId) return
    setLoading(true)
    try {
      const [orgResponse, charResponse] = await Promise.all([
        organizationService.list(projectId, { page_size: 100 }),
        characterService.list(projectId, { page_size: 200 }),
      ])
      setOrganizations(orgResponse.items || [])
      setCharacters(charResponse.items || [])
    } catch (error) {
      message.error('获取数据失败')
      console.error(error)
    } finally {
      setLoading(false)
    }
  }, [projectId])

  useEffect(() => {
    fetchOrganizations()
  }, [fetchOrganizations])

  const fetchMembers = async (orgId: string) => {
    try {
      const data = await organizationService.getMembers(orgId)
      setMembers(data || [])
    } catch (error) {
      console.error(error)
    }
  }

  const handleCreate = () => {
    setEditingOrg(null)
    form.resetFields()
    setModalOpen(true)
  }

  const handleEdit = (org: Organization) => {
    setEditingOrg(org)
    form.setFieldsValue(org)
    setModalOpen(true)
  }

  const handleDelete = async (id: string) => {
    try {
      await organizationService.delete(id)
      message.success('删除成功')
      setSelectedOrg(null)
      fetchOrganizations()
    } catch (error) {
      message.error('删除失败')
      console.error(error)
    }
  }

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      if (editingOrg) {
        await organizationService.update(editingOrg.id, values)
        message.success('更新成功')
      } else {
        await organizationService.create(projectId!, values)
        message.success('创建成功')
      }
      setModalOpen(false)
      fetchOrganizations()
    } catch (error) {
      console.error(error)
    }
  }

  const handleViewDetail = async (org: Organization) => {
    setSelectedOrg(org)
    await fetchMembers(org.id)
  }

  const handleAddMember = () => {
    memberForm.resetFields()
    setMemberModalOpen(true)
  }

  const handleMemberSubmit = async () => {
    if (!selectedOrg) return
    try {
      const values = await memberForm.validateFields()
      await organizationService.addMember(selectedOrg.id, values)
      message.success('添加成功')
      setMemberModalOpen(false)
      fetchMembers(selectedOrg.id)
    } catch (error) {
      console.error(error)
    }
  }

  const handleRemoveMember = async (memberId: string) => {
    try {
      await organizationService.removeMember(memberId)
      message.success('移除成功')
      if (selectedOrg) {
        fetchMembers(selectedOrg.id)
      }
    } catch (error) {
      message.error('移除失败')
      console.error(error)
    }
  }

  const columns = [
    {
      title: '组织',
      dataIndex: 'name',
      key: 'name',
      render: (name: string) => (
        <Space>
          <TeamOutlined />
          <span className="font-medium">{name}</span>
        </Space>
      ),
    },
    {
      title: '类型',
      dataIndex: 'type',
      key: 'type',
      width: 100,
      render: (type: string) => (
        <Tag>{orgTypes.find((t) => t.value === type)?.label || type || '未设置'}</Tag>
      ),
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
      ellipsis: true,
      render: (text: string) => text || '-',
    },
    {
      title: '操作',
      key: 'action',
      width: 180,
      render: (_: unknown, record: Organization) => (
        <Space>
          <Button type="link" size="small" onClick={() => handleViewDetail(record)}>
            详情
          </Button>
          <Button type="link" size="small" icon={<EditOutlined />} onClick={() => handleEdit(record)}>
            编辑
          </Button>
          <Popconfirm title="确定删除此组织？" onConfirm={() => handleDelete(record.id)}>
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
        title="组织管理"
        extra={
          <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
            新建组织
          </Button>
        }
      >
        <Table
          columns={columns}
          dataSource={organizations}
          rowKey="id"
          loading={loading}
          pagination={{ pageSize: 10 }}
        />
      </Card>

      {/* 创建/编辑弹窗 */}
      <Modal
        title={editingOrg ? '编辑组织' : '新建组织'}
        open={modalOpen}
        onOk={handleSubmit}
        onCancel={() => setModalOpen(false)}
        width={600}
        destroyOnClose
      >
        <Form form={form} layout="vertical" className="mt-4">
          <div className="grid grid-cols-2 gap-4">
            <Form.Item
              name="name"
              label="组织名称"
              rules={[{ required: true, message: '请输入组织名称' }]}
            >
              <Input placeholder="组织名称" />
            </Form.Item>
            <Form.Item name="type" label="组织类型">
              <Select placeholder="选择类型" options={orgTypes} allowClear />
            </Form.Item>
          </div>

          <Form.Item name="description" label="描述">
            <TextArea rows={2} placeholder="组织描述" />
          </Form.Item>

          <Form.Item name="history" label="历史">
            <TextArea rows={2} placeholder="组织历史" />
          </Form.Item>

          <Form.Item name="structure" label="组织结构">
            <TextArea rows={2} placeholder="组织结构说明" />
          </Form.Item>

          <Form.Item name="goals" label="目标/宗旨">
            <TextArea rows={2} placeholder="组织目标或宗旨" />
          </Form.Item>
        </Form>
      </Modal>

      {/* 详情弹窗 */}
      <Modal
        title={selectedOrg?.name}
        open={!!selectedOrg}
        onCancel={() => setSelectedOrg(null)}
        footer={null}
        width={800}
      >
        {selectedOrg && (
          <Tabs
            items={[
              {
                key: 'info',
                label: '基本信息',
                children: (
                  <div className="space-y-3">
                    <p><strong>类型：</strong>{orgTypes.find((t) => t.value === selectedOrg.type)?.label || selectedOrg.type || '-'}</p>
                    <p><strong>描述：</strong>{selectedOrg.description || '-'}</p>
                    <p><strong>历史：</strong>{selectedOrg.history || '-'}</p>
                    <p><strong>结构：</strong>{selectedOrg.structure || '-'}</p>
                    <p><strong>目标：</strong>{selectedOrg.goals || '-'}</p>
                  </div>
                ),
              },
              {
                key: 'members',
                label: `成员 (${members.length})`,
                children: (
                  <div>
                    <div className="mb-4">
                      <Button type="primary" icon={<PlusOutlined />} onClick={handleAddMember}>
                        添加成员
                      </Button>
                    </div>
                    {members.length > 0 ? (
                      <List
                        dataSource={members}
                        renderItem={(member) => (
                          <List.Item
                            actions={[
                              <Popconfirm
                                key="remove"
                                title="确定移除此成员？"
                                onConfirm={() => handleRemoveMember(member.id)}
                              >
                                <Button type="link" danger size="small">
                                  移除
                                </Button>
                              </Popconfirm>,
                            ]}
                          >
                            <List.Item.Meta
                              avatar={<Avatar icon={<UserOutlined />} />}
                              title={member.character?.name || '未知角色'}
                              description={
                                <Space>
                                  {member.position && <Tag color="blue">{member.position}</Tag>}
                                  {member.rank && <Tag>{member.rank}</Tag>}
                                </Space>
                              }
                            />
                          </List.Item>
                        )}
                      />
                    ) : (
                      <Empty description="暂无成员" />
                    )}
                  </div>
                ),
              },
            ]}
          />
        )}
      </Modal>

      {/* 添加成员弹窗 */}
      <Modal
        title="添加成员"
        open={memberModalOpen}
        onOk={handleMemberSubmit}
        onCancel={() => setMemberModalOpen(false)}
        destroyOnClose
      >
        <Form form={memberForm} layout="vertical" className="mt-4">
          <Form.Item
            name="character_id"
            label="选择角色"
            rules={[{ required: true, message: '请选择角色' }]}
          >
            <Select
              placeholder="选择角色"
              showSearch
              optionFilterProp="label"
              options={characters.map((c) => ({ value: c.id, label: c.name }))}
            />
          </Form.Item>
          <Form.Item name="position" label="职位">
            <Input placeholder="职位" />
          </Form.Item>
          <Form.Item name="rank" label="等级/阶级">
            <Input placeholder="等级或阶级" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default OrganizationManager
