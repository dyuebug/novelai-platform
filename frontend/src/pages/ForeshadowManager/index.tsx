import React, { useState, useEffect, useCallback } from 'react';
import { useParams } from 'react-router-dom';
import {
  Card,
  Table,
  Button,
  Space,
  Tag,
  Modal,
  Form,
  Input,
  Select,
  InputNumber,
  message,
  Popconfirm,
  Statistic,
  Row,
  Col,
  Tooltip,
  Drawer,
  Timeline,
  Empty,
  Switch,
  Tabs,
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  CheckCircleOutlined,
  ExclamationCircleOutlined,
  EyeOutlined,
  HistoryOutlined,
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import foreshadowService, {
  Foreshadow,
  ForeshadowStats,
  ForeshadowHint,
  ForeshadowReminder,
  ForeshadowStatus,
  ForeshadowPriority,
  CreateForeshadowRequest,
  UpdateForeshadowRequest,
} from '../../services/foreshadow.service';
import ForeshadowTimeline from './ForeshadowTimeline';
import ForeshadowReminderNotification from './ForeshadowReminderNotification';

interface ForeshadowManagerProps {
  currentChapter?: number;
}

const statusConfig: Record<ForeshadowStatus, { color: string; label: string }> = {
  planted: { color: 'blue', label: '已埋设' },
  hinted: { color: 'orange', label: '已暗示' },
  resolved: { color: 'green', label: '已回收' },
  abandoned: { color: 'default', label: '已放弃' },
};

const priorityConfig: Record<ForeshadowPriority, { color: string; label: string }> = {
  high: { color: 'red', label: '高优先级' },
  medium: { color: 'orange', label: '中优先级' },
  low: { color: 'blue', label: '低优先级' },
};

const ForeshadowManager: React.FC<ForeshadowManagerProps> = ({
  currentChapter = 0,
}) => {
  const { projectId } = useParams<{ projectId: string }>();
  const [foreshadows, setForeshadows] = useState<Foreshadow[]>([]);
  const [stats, setStats] = useState<ForeshadowStats | null>(null);
  const [loading, setLoading] = useState(false);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const [statusFilter, setStatusFilter] = useState<ForeshadowStatus | undefined>();
  const [priorityFilter, setPriorityFilter] = useState<ForeshadowPriority | undefined>();

  // 模态框状态
  const [createModalVisible, setCreateModalVisible] = useState(false);
  const [editModalVisible, setEditModalVisible] = useState(false);
  const [resolveModalVisible, setResolveModalVisible] = useState(false);
  const [detailDrawerVisible, setDetailDrawerVisible] = useState(false);
  const [selectedForeshadow, setSelectedForeshadow] = useState<Foreshadow | null>(null);
  const [hints, setHints] = useState<ForeshadowHint[]>([]);

  const [createForm] = Form.useForm();
  const [editForm] = Form.useForm();
  const [resolveForm] = Form.useForm();

  // 加载伏笔列表
  const loadForeshadows = useCallback(async () => {
    if (!projectId) return;
    setLoading(true);
    try {
      const response = await foreshadowService.list(projectId, {
        status: statusFilter,
        priority: priorityFilter,
        page,
        page_size: pageSize,
      });
      setForeshadows(response.items || []);
      setTotal(response.total);
    } catch (error) {
      message.error('加载伏笔列表失败');
    } finally {
      setLoading(false);
    }
  }, [projectId, statusFilter, priorityFilter, page, pageSize]);

  // 加载统计数据
  const loadStats = useCallback(async () => {
    if (!projectId) return;
    try {
      const data = await foreshadowService.getStats(projectId, currentChapter);
      setStats(data);
    } catch (error) {
      console.error('加载统计失败:', error);
    }
  }, [projectId, currentChapter]);

  useEffect(() => {
    loadForeshadows();
    loadStats();
  }, [loadForeshadows, loadStats]);

  // 创建伏笔
  const handleCreate = async (values: CreateForeshadowRequest) => {
    if (!projectId) return;
    try {
      await foreshadowService.create(projectId, values);
      message.success('创建成功');
      setCreateModalVisible(false);
      createForm.resetFields();
      loadForeshadows();
      loadStats();
    } catch (error) {
      message.error('创建失败');
    }
  };

  // 更新伏笔
  const handleUpdate = async (values: UpdateForeshadowRequest) => {
    if (!selectedForeshadow) return;
    try {
      await foreshadowService.update(selectedForeshadow.id, values);
      message.success('更新成功');
      setEditModalVisible(false);
      editForm.resetFields();
      loadForeshadows();
      loadStats();
    } catch (error) {
      message.error('更新失败');
    }
  };

  // 删除伏笔
  const handleDelete = async (id: string) => {
    try {
      await foreshadowService.delete(id);
      message.success('删除成功');
      loadForeshadows();
      loadStats();
    } catch (error) {
      message.error('删除失败');
    }
  };

  // 回收伏笔
  const handleResolve = async (values: { chapter_num: number; content: string }) => {
    if (!selectedForeshadow) return;
    try {
      await foreshadowService.resolve(selectedForeshadow.id, values);
      message.success('回收成功');
      setResolveModalVisible(false);
      resolveForm.resetFields();
      loadForeshadows();
      loadStats();
    } catch (error) {
      message.error('回收失败');
    }
  };

  // 查看详情
  const handleViewDetail = async (foreshadow: Foreshadow) => {
    setSelectedForeshadow(foreshadow);
    setDetailDrawerVisible(true);
    try {
      const hintsData = await foreshadowService.getHints(foreshadow.id);
      setHints(hintsData || []);
    } catch (error) {
      console.error('加载暗示失败:', error);
    }
  };

  // 打开编辑模态框
  const openEditModal = (foreshadow: Foreshadow) => {
    setSelectedForeshadow(foreshadow);
    editForm.setFieldsValue(foreshadow);
    setEditModalVisible(true);
  };

  // 打开回收模态框
  const openResolveModal = (foreshadow: Foreshadow) => {
    setSelectedForeshadow(foreshadow);
    resolveForm.setFieldsValue({ chapter_num: currentChapter });
    setResolveModalVisible(true);
  };

  // 表格列定义
  const columns: ColumnsType<Foreshadow> = [
    {
      title: '标题',
      dataIndex: 'title',
      key: 'title',
      width: 200,
      ellipsis: true,
      render: (text, record) => (
        <a onClick={() => handleViewDetail(record)}>{text}</a>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (status: ForeshadowStatus) => (
        <Tag color={statusConfig[status]?.color}>{statusConfig[status]?.label}</Tag>
      ),
    },
    {
      title: '优先级',
      dataIndex: 'priority',
      key: 'priority',
      width: 100,
      render: (priority: ForeshadowPriority) => (
        <Tag color={priorityConfig[priority]?.color}>{priorityConfig[priority]?.label}</Tag>
      ),
    },
    {
      title: '埋设章节',
      dataIndex: 'plant_chapter_num',
      key: 'plant_chapter_num',
      width: 100,
      render: (num: number) => (num > 0 ? `第${num}章` : '-'),
    },
    {
      title: '提醒章节',
      dataIndex: 'remind_chapter_num',
      key: 'remind_chapter_num',
      width: 100,
      render: (num: number, record) => {
        if (num <= 0) return '-';
        const isOverdue = currentChapter > num && record.status === 'planted';
        return (
          <span style={{ color: isOverdue ? '#ff4d4f' : undefined }}>
            第{num}章
            {isOverdue && <ExclamationCircleOutlined style={{ marginLeft: 4 }} />}
          </span>
        );
      },
    },
    {
      title: '回收章节',
      dataIndex: 'resolve_chapter_num',
      key: 'resolve_chapter_num',
      width: 100,
      render: (num: number) => (num > 0 ? `第${num}章` : '-'),
    },
    {
      title: '标签',
      dataIndex: 'tags',
      key: 'tags',
      width: 150,
      render: (tags: string[]) =>
        tags?.slice(0, 2).map((tag) => (
          <Tag key={tag} style={{ marginBottom: 2 }}>
            {tag}
          </Tag>
        )),
    },
    {
      title: '操作',
      key: 'action',
      width: 200,
      fixed: 'right',
      render: (_, record) => (
        <Space size="small">
          <Tooltip title="查看详情">
            <Button
              type="text"
              size="small"
              icon={<EyeOutlined />}
              onClick={() => handleViewDetail(record)}
            />
          </Tooltip>
          <Tooltip title="编辑">
            <Button
              type="text"
              size="small"
              icon={<EditOutlined />}
              onClick={() => openEditModal(record)}
            />
          </Tooltip>
          {record.status === 'planted' || record.status === 'hinted' ? (
            <Tooltip title="回收">
              <Button
                type="text"
                size="small"
                icon={<CheckCircleOutlined />}
                onClick={() => openResolveModal(record)}
              />
            </Tooltip>
          ) : null}
          <Popconfirm
            title="确定删除此伏笔?"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Tooltip title="删除">
              <Button type="text" size="small" danger icon={<DeleteOutlined />} />
            </Tooltip>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  // 处理提醒点击
  const handleReminderClick = (reminder: ForeshadowReminder) => {
    if (reminder.foreshadow) {
      handleViewDetail(reminder.foreshadow);
    }
  };

  // 处理时间线伏笔点击
  const handleTimelineForeshadowClick = (foreshadow: Foreshadow) => {
    handleViewDetail(foreshadow);
  };

  return (
    <div style={{ padding: 24 }}>
      {/* 页面标题和提醒通知 */}
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
        <h2 style={{ margin: 0 }}>伏笔管理</h2>
        {projectId && (
          <ForeshadowReminderNotification
            projectId={projectId}
            currentChapter={currentChapter}
            onReminderClick={handleReminderClick}
          />
        )}
      </div>

      {/* 统计卡片 */}
      {stats && (
        <Row gutter={16} style={{ marginBottom: 24 }}>
          <Col span={3}>
            <Card size="small">
              <Statistic title="总数" value={stats.total} />
            </Card>
          </Col>
          <Col span={3}>
            <Card size="small">
              <Statistic
                title="已埋设"
                value={stats.planted}
                valueStyle={{ color: '#1890ff' }}
              />
            </Card>
          </Col>
          <Col span={3}>
            <Card size="small">
              <Statistic
                title="已暗示"
                value={stats.hinted}
                valueStyle={{ color: '#fa8c16' }}
              />
            </Card>
          </Col>
          <Col span={3}>
            <Card size="small">
              <Statistic
                title="已回收"
                value={stats.resolved}
                valueStyle={{ color: '#52c41a' }}
              />
            </Card>
          </Col>
          <Col span={3}>
            <Card size="small">
              <Statistic title="已放弃" value={stats.abandoned} />
            </Card>
          </Col>
          <Col span={3}>
            <Card size="small">
              <Statistic
                title="高优先级"
                value={stats.high_priority}
                valueStyle={{ color: '#ff4d4f' }}
              />
            </Card>
          </Col>
          <Col span={3}>
            <Card size="small">
              <Statistic
                title="超期未回收"
                value={stats.overdue_count}
                valueStyle={{ color: stats.overdue_count > 0 ? '#ff4d4f' : undefined }}
                prefix={stats.overdue_count > 0 ? <ExclamationCircleOutlined /> : null}
              />
            </Card>
          </Col>
        </Row>
      )}

      {/* 视图切换 Tabs */}
      <Tabs
        defaultActiveKey="list"
        items={[
          {
            key: 'list',
            label: '列表视图',
            children: (
              <Card>
                <Space style={{ marginBottom: 16 }}>
                  <Button
                    type="primary"
                    icon={<PlusOutlined />}
                    onClick={() => setCreateModalVisible(true)}
                  >
                    新建伏笔
                  </Button>
                  <Select
                    placeholder="状态筛选"
                    allowClear
                    style={{ width: 120 }}
                    value={statusFilter}
                    onChange={setStatusFilter}
                  >
                    <Select.Option value="planted">已埋设</Select.Option>
                    <Select.Option value="hinted">已暗示</Select.Option>
                    <Select.Option value="resolved">已回收</Select.Option>
                    <Select.Option value="abandoned">已放弃</Select.Option>
                  </Select>
                  <Select
                    placeholder="优先级筛选"
                    allowClear
                    style={{ width: 120 }}
                    value={priorityFilter}
                    onChange={setPriorityFilter}
                  >
                    <Select.Option value="high">高优先级</Select.Option>
                    <Select.Option value="medium">中优先级</Select.Option>
                    <Select.Option value="low">低优先级</Select.Option>
                  </Select>
                </Space>

                {/* 伏笔列表 */}
                <Table
                  columns={columns}
                  dataSource={foreshadows}
                  rowKey="id"
                  loading={loading}
                  scroll={{ x: 1200 }}
                  pagination={{
                    current: page,
                    pageSize,
                    total,
                    showSizeChanger: true,
                    showQuickJumper: true,
                    showTotal: (t) => `共 ${t} 条`,
                    onChange: (p, ps) => {
                      setPage(p);
                      setPageSize(ps);
                    },
                  }}
                />
              </Card>
            ),
          },
          {
            key: 'timeline',
            label: '时间线视图',
            children: projectId ? (
              <ForeshadowTimeline
                projectId={projectId}
                currentChapter={currentChapter}
                maxChapter={100}
                onForeshadowClick={handleTimelineForeshadowClick}
              />
            ) : null,
          },
        ]}
      />

      {/* 创建伏笔模态框 */}
      <Modal
        title="新建伏笔"
        open={createModalVisible}
        onCancel={() => {
          setCreateModalVisible(false);
          createForm.resetFields();
        }}
        onOk={() => createForm.submit()}
        width={600}
      >
        <Form form={createForm} layout="vertical" onFinish={handleCreate}>
          <Form.Item
            name="title"
            label="标题"
            rules={[{ required: true, message: '请输入标题' }]}
          >
            <Input placeholder="伏笔标题" />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <Input.TextArea rows={2} placeholder="伏笔描述" />
          </Form.Item>
          <Form.Item name="content" label="内容/原文">
            <Input.TextArea rows={3} placeholder="伏笔内容或原文引用" />
          </Form.Item>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="priority" label="优先级" initialValue="medium">
                <Select>
                  <Select.Option value="high">高优先级 (主线)</Select.Option>
                  <Select.Option value="medium">中优先级 (支线)</Select.Option>
                  <Select.Option value="low">低优先级 (细节)</Select.Option>
                </Select>
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="plant_chapter_num" label="埋设章节">
                <InputNumber min={0} style={{ width: '100%' }} placeholder="章节号" />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="remind_chapter_num" label="提醒章节">
                <InputNumber min={0} style={{ width: '100%' }} placeholder="到达此章节时提醒" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="plant_position" label="埋设位置">
                <Input placeholder="位置描述" />
              </Form.Item>
            </Col>
          </Row>
          <Form.Item name="tags" label="标签">
            <Select mode="tags" placeholder="输入标签后回车" />
          </Form.Item>
        </Form>
      </Modal>

      {/* 编辑伏笔模态框 */}
      <Modal
        title="编辑伏笔"
        open={editModalVisible}
        onCancel={() => {
          setEditModalVisible(false);
          editForm.resetFields();
        }}
        onOk={() => editForm.submit()}
        width={600}
      >
        <Form form={editForm} layout="vertical" onFinish={handleUpdate}>
          <Form.Item
            name="title"
            label="标题"
            rules={[{ required: true, message: '请输入标题' }]}
          >
            <Input placeholder="伏笔标题" />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <Input.TextArea rows={2} placeholder="伏笔描述" />
          </Form.Item>
          <Form.Item name="content" label="内容/原文">
            <Input.TextArea rows={3} placeholder="伏笔内容或原文引用" />
          </Form.Item>
          <Row gutter={16}>
            <Col span={8}>
              <Form.Item name="priority" label="优先级">
                <Select>
                  <Select.Option value="high">高优先级</Select.Option>
                  <Select.Option value="medium">中优先级</Select.Option>
                  <Select.Option value="low">低优先级</Select.Option>
                </Select>
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item name="status" label="状态">
                <Select>
                  <Select.Option value="planted">已埋设</Select.Option>
                  <Select.Option value="hinted">已暗示</Select.Option>
                  <Select.Option value="resolved">已回收</Select.Option>
                  <Select.Option value="abandoned">已放弃</Select.Option>
                </Select>
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item name="remind_enabled" label="启用提醒" valuePropName="checked">
                <Switch />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="plant_chapter_num" label="埋设章节">
                <InputNumber min={0} style={{ width: '100%' }} />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="remind_chapter_num" label="提醒章节">
                <InputNumber min={0} style={{ width: '100%' }} />
              </Form.Item>
            </Col>
          </Row>
          <Form.Item name="tags" label="标签">
            <Select mode="tags" placeholder="输入标签后回车" />
          </Form.Item>
        </Form>
      </Modal>

      {/* 回收伏笔模态框 */}
      <Modal
        title="回收伏笔"
        open={resolveModalVisible}
        onCancel={() => {
          setResolveModalVisible(false);
          resolveForm.resetFields();
        }}
        onOk={() => resolveForm.submit()}
      >
        <Form form={resolveForm} layout="vertical" onFinish={handleResolve}>
          <Form.Item
            name="chapter_num"
            label="回收章节"
            rules={[{ required: true, message: '请输入回收章节' }]}
          >
            <InputNumber min={1} style={{ width: '100%' }} placeholder="章节号" />
          </Form.Item>
          <Form.Item
            name="content"
            label="回收内容"
            rules={[{ required: true, message: '请输入回收内容' }]}
          >
            <Input.TextArea rows={4} placeholder="描述伏笔如何被回收" />
          </Form.Item>
        </Form>
      </Modal>

      {/* 详情抽屉 */}
      <Drawer
        title={selectedForeshadow?.title}
        open={detailDrawerVisible}
        onClose={() => setDetailDrawerVisible(false)}
        width={500}
      >
        {selectedForeshadow && (
          <div>
            <div style={{ marginBottom: 16 }}>
              <Tag color={statusConfig[selectedForeshadow.status]?.color}>
                {statusConfig[selectedForeshadow.status]?.label}
              </Tag>
              <Tag color={priorityConfig[selectedForeshadow.priority]?.color}>
                {priorityConfig[selectedForeshadow.priority]?.label}
              </Tag>
            </div>

            {selectedForeshadow.description && (
              <div style={{ marginBottom: 16 }}>
                <h4>描述</h4>
                <p>{selectedForeshadow.description}</p>
              </div>
            )}

            {selectedForeshadow.content && (
              <div style={{ marginBottom: 16 }}>
                <h4>内容/原文</h4>
                <p style={{ background: '#f5f5f5', padding: 12, borderRadius: 4 }}>
                  {selectedForeshadow.content}
                </p>
              </div>
            )}

            <Row gutter={16} style={{ marginBottom: 16 }}>
              <Col span={12}>
                <Statistic
                  title="埋设章节"
                  value={
                    selectedForeshadow.plant_chapter_num > 0
                      ? `第${selectedForeshadow.plant_chapter_num}章`
                      : '-'
                  }
                />
              </Col>
              <Col span={12}>
                <Statistic
                  title="提醒章节"
                  value={
                    selectedForeshadow.remind_chapter_num > 0
                      ? `第${selectedForeshadow.remind_chapter_num}章`
                      : '-'
                  }
                />
              </Col>
            </Row>

            {selectedForeshadow.status === 'resolved' && (
              <div style={{ marginBottom: 16 }}>
                <h4>回收信息</h4>
                <p>
                  <strong>回收章节:</strong> 第{selectedForeshadow.resolve_chapter_num}章
                </p>
                {selectedForeshadow.resolve_content && (
                  <p style={{ background: '#f6ffed', padding: 12, borderRadius: 4 }}>
                    {selectedForeshadow.resolve_content}
                  </p>
                )}
              </div>
            )}

            {selectedForeshadow.tags?.length > 0 && (
              <div style={{ marginBottom: 16 }}>
                <h4>标签</h4>
                {selectedForeshadow.tags.map((tag) => (
                  <Tag key={tag}>{tag}</Tag>
                ))}
              </div>
            )}

            <div>
              <h4>
                <HistoryOutlined style={{ marginRight: 8 }} />
                暗示记录
              </h4>
              {hints.length > 0 ? (
                <Timeline
                  items={hints.map((hint) => ({
                    children: (
                      <div>
                        <div style={{ fontWeight: 500 }}>
                          第{hint.chapter_num}章
                          {hint.hint_type && (
                            <Tag style={{ marginLeft: 8, fontSize: 12 }}>
                              {hint.hint_type}
                            </Tag>
                          )}
                        </div>
                        <div style={{ color: '#666' }}>{hint.content}</div>
                      </div>
                    ),
                  }))}
                />
              ) : (
                <Empty description="暂无暗示记录" image={Empty.PRESENTED_IMAGE_SIMPLE} />
              )}
            </div>
          </div>
        )}
      </Drawer>
    </div>
  );
};

export default ForeshadowManager;
