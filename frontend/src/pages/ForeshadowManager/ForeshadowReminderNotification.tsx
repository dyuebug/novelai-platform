import React, { useState, useEffect, useCallback } from 'react';
import {
  Badge,
  Popover,
  List,
  Button,
  Empty,
  Tag,
  Space,
  message,
  Spin,
} from 'antd';
import {
  BellOutlined,
  CheckOutlined,
  ExclamationCircleOutlined,
} from '@ant-design/icons';
import foreshadowService, { ForeshadowReminder } from '../../services/foreshadow.service';

interface ForeshadowReminderNotificationProps {
  projectId: string;
  currentChapter?: number;
  onReminderClick?: (reminder: ForeshadowReminder) => void;
}

const ForeshadowReminderNotification: React.FC<ForeshadowReminderNotificationProps> = ({
  projectId,
  currentChapter = 0,
  onReminderClick,
}) => {
  const [reminders, setReminders] = useState<ForeshadowReminder[]>([]);
  const [loading, setLoading] = useState(false);
  const [visible, setVisible] = useState(false);

  // 加载未读提醒
  const loadReminders = useCallback(async () => {
    if (!projectId) return;
    setLoading(true);
    try {
      const data = await foreshadowService.getUnreadReminders(projectId);
      setReminders(data || []);
    } catch (error) {
      console.error('加载提醒失败:', error);
    } finally {
      setLoading(false);
    }
  }, [projectId]);

  useEffect(() => {
    loadReminders();
  }, [loadReminders]);

  // 检查并创建新提醒
  const checkReminders = async () => {
    if (!projectId || currentChapter <= 0) return;
    try {
      const result = await foreshadowService.checkReminders(projectId, currentChapter);
      if (result.count > 0) {
        message.info(`发现 ${result.count} 个新的伏笔提醒`);
        loadReminders();
      }
    } catch (error) {
      console.error('检查提醒失败:', error);
    }
  };

  // 标记单个为已读
  const markAsRead = async (reminderId: string, e: React.MouseEvent) => {
    e.stopPropagation();
    try {
      await foreshadowService.markReminderAsRead(reminderId);
      setReminders((prev) => prev.filter((r) => r.id !== reminderId));
      message.success('已标记为已读');
    } catch (error) {
      message.error('操作失败');
    }
  };

  // 标记全部为已读
  const markAllAsRead = async () => {
    try {
      await foreshadowService.markAllRemindersAsRead(projectId);
      setReminders([]);
      message.success('已全部标记为已读');
    } catch (error) {
      message.error('操作失败');
    }
  };

  const content = (
    <div style={{ width: 350 }}>
      <div
        style={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          marginBottom: 12,
          paddingBottom: 8,
          borderBottom: '1px solid #f0f0f0',
        }}
      >
        <span style={{ fontWeight: 500 }}>伏笔提醒</span>
        <Space>
          <Button size="small" onClick={checkReminders}>
            检查新提醒
          </Button>
          {reminders.length > 0 && (
            <Button size="small" type="link" onClick={markAllAsRead}>
              全部已读
            </Button>
          )}
        </Space>
      </div>

      {loading ? (
        <div style={{ textAlign: 'center', padding: 20 }}>
          <Spin />
        </div>
      ) : reminders.length === 0 ? (
        <Empty
          description="暂无未读提醒"
          image={Empty.PRESENTED_IMAGE_SIMPLE}
          style={{ padding: 20 }}
        />
      ) : (
        <List
          dataSource={reminders}
          style={{ maxHeight: 400, overflowY: 'auto' }}
          renderItem={(reminder) => (
            <List.Item
              style={{
                cursor: 'pointer',
                padding: '8px 0',
                borderBottom: '1px solid #f5f5f5',
              }}
              onClick={() => {
                onReminderClick?.(reminder);
                setVisible(false);
              }}
              actions={[
                <Button
                  key="read"
                  type="text"
                  size="small"
                  icon={<CheckOutlined />}
                  onClick={(e) => markAsRead(reminder.id, e)}
                />,
              ]}
            >
              <List.Item.Meta
                avatar={
                  <ExclamationCircleOutlined
                    style={{
                      fontSize: 20,
                      color: reminder.message.includes('超过') ? '#ff4d4f' : '#fa8c16',
                    }}
                  />
                }
                title={
                  <div>
                    {reminder.foreshadow?.title || '伏笔提醒'}
                    <Tag style={{ marginLeft: 8, fontSize: 12 }}>
                      第{reminder.chapter_num}章
                    </Tag>
                  </div>
                }
                description={
                  <div
                    style={{
                      fontSize: 12,
                      color: '#8c8c8c',
                      overflow: 'hidden',
                      textOverflow: 'ellipsis',
                      whiteSpace: 'nowrap',
                    }}
                  >
                    {reminder.message}
                  </div>
                }
              />
            </List.Item>
          )}
        />
      )}
    </div>
  );

  return (
    <Popover
      content={content}
      trigger="click"
      open={visible}
      onOpenChange={setVisible}
      placement="bottomRight"
    >
      <Badge count={reminders.length} size="small" offset={[-2, 2]}>
        <Button
          type="text"
          icon={<BellOutlined style={{ fontSize: 18 }} />}
          style={{ padding: '4px 8px' }}
        />
      </Badge>
    </Popover>
  );
};

export default ForeshadowReminderNotification;
