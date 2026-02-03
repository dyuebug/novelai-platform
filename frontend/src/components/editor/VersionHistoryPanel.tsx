import { useState, useEffect } from 'react'
import { Drawer, List, Button, Tag, Spin, Empty, message, Modal } from 'antd'
import { HistoryOutlined, RollbackOutlined, EyeOutlined } from '@ant-design/icons'
import { chapterService, ChapterVersion } from '@/services/chapter.service'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'

dayjs.extend(relativeTime)
dayjs.locale('zh-cn')

interface VersionHistoryPanelProps {
  chapterId: string
  open: boolean
  onClose: () => void
  onRestore: () => void
}

const sourceLabels: Record<string, { label: string; color: string }> = {
  manual: { label: '手动保存', color: 'blue' },
  auto: { label: '自动保存', color: 'green' },
  ai: { label: 'AI 生成', color: 'purple' },
  restore: { label: '版本恢复', color: 'orange' },
}

const VersionHistoryPanel = ({
  chapterId,
  open,
  onClose,
  onRestore,
}: VersionHistoryPanelProps) => {
  const [versions, setVersions] = useState<ChapterVersion[]>([])
  const [loading, setLoading] = useState(false)
  const [previewVersion, setPreviewVersion] = useState<ChapterVersion | null>(null)
  const [previewOpen, setPreviewOpen] = useState(false)
  const [restoring, setRestoring] = useState(false)

  useEffect(() => {
    if (open && chapterId) {
      fetchVersions()
    }
  }, [open, chapterId])

  const fetchVersions = async () => {
    setLoading(true)
    try {
      const response = await chapterService.listVersions(chapterId, 1, 50)
      setVersions(response.items || [])
    } catch (error) {
      message.error('获取版本历史失败')
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  const handlePreview = async (version: ChapterVersion) => {
    setPreviewVersion(version)
    setPreviewOpen(true)
  }

  const handleRestore = async (versionNumber: number) => {
    Modal.confirm({
      title: '确认恢复',
      content: `确定要恢复到版本 ${versionNumber} 吗？当前内容将被替换。`,
      okText: '确认恢复',
      cancelText: '取消',
      onOk: async () => {
        setRestoring(true)
        try {
          await chapterService.restoreVersion(chapterId, versionNumber)
          message.success('版本恢复成功')
          onRestore()
          onClose()
        } catch (error) {
          message.error('版本恢复失败')
          console.error(error)
        } finally {
          setRestoring(false)
        }
      },
    })
  }

  return (
    <>
      <Drawer
        title={
          <span>
            <HistoryOutlined className="mr-2" />
            版本历史
          </span>
        }
        placement="right"
        width={400}
        open={open}
        onClose={onClose}
      >
        {loading ? (
          <div className="flex justify-center items-center h-40">
            <Spin />
          </div>
        ) : versions.length === 0 ? (
          <Empty description="暂无版本记录" />
        ) : (
          <List
            dataSource={versions}
            renderItem={(version) => {
              const sourceInfo = sourceLabels[version.source] || {
                label: version.source,
                color: 'default',
              }
              return (
                <List.Item
                  actions={[
                    <Button
                      key="preview"
                      type="link"
                      size="small"
                      icon={<EyeOutlined />}
                      onClick={() => handlePreview(version)}
                    >
                      预览
                    </Button>,
                    <Button
                      key="restore"
                      type="link"
                      size="small"
                      icon={<RollbackOutlined />}
                      onClick={() => handleRestore(version.version_number)}
                      loading={restoring}
                    >
                      恢复
                    </Button>,
                  ]}
                >
                  <List.Item.Meta
                    title={
                      <div className="flex items-center gap-2">
                        <span>版本 {version.version_number}</span>
                        <Tag color={sourceInfo.color}>{sourceInfo.label}</Tag>
                      </div>
                    }
                    description={
                      <div className="text-xs text-gray-400">
                        <div>{dayjs(version.created_at).fromNow()}</div>
                        <div>{version.word_count} 字</div>
                        {version.ai_model && (
                          <div className="text-purple-500">
                            {version.ai_provider} / {version.ai_model}
                          </div>
                        )}
                      </div>
                    }
                  />
                </List.Item>
              )
            }}
          />
        )}
      </Drawer>

      {/* 版本预览弹窗 */}
      <Modal
        title={`版本 ${previewVersion?.version_number} 预览`}
        open={previewOpen}
        onCancel={() => setPreviewOpen(false)}
        footer={[
          <Button key="close" onClick={() => setPreviewOpen(false)}>
            关闭
          </Button>,
          <Button
            key="restore"
            type="primary"
            icon={<RollbackOutlined />}
            onClick={() => {
              setPreviewOpen(false)
              if (previewVersion) {
                handleRestore(previewVersion.version_number)
              }
            }}
          >
            恢复此版本
          </Button>,
        ]}
        width={700}
      >
        {previewVersion && (
          <div className="max-h-96 overflow-auto">
            <div className="mb-2 text-sm text-gray-500">
              <span className="mr-4">字数: {previewVersion.word_count}</span>
              <span>保存时间: {dayjs(previewVersion.created_at).format('YYYY-MM-DD HH:mm:ss')}</span>
            </div>
            <div
              className="prose prose-sm max-w-none p-4 bg-gray-50 rounded"
              dangerouslySetInnerHTML={{ __html: previewVersion.content }}
            />
          </div>
        )}
      </Modal>
    </>
  )
}

export default VersionHistoryPanel
