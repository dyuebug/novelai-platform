import { useState } from 'react'
import {
  Card,
  Button,
  Tag,
  List,
  Space,
  Spin,
  Empty,
  Alert,
  Progress,
} from 'antd'
import {
  CheckCircleOutlined,
  WarningOutlined,
  CloseCircleOutlined,
  UserOutlined,
  ClockCircleOutlined,
  GlobalOutlined,
} from '@ant-design/icons'
import {
  qualityService,
  ConsistencyResult,
  ConsistencyIssue,
  QualityRequest,
} from '@/services/quality.service'
import { Character, WorldSetting } from '@/services/world.service'

interface ConsistencyCheckerProps {
  chapterId: string
  content: string
  characters?: Character[]
  worldSettings?: WorldSetting[]
  previousSummary?: string
  onCheck?: (result: ConsistencyResult) => void
}

// 问题类型映射
const issueTypeLabels: Record<string, { label: string; icon: React.ReactNode }> = {
  character: { label: '角色一致性', icon: <UserOutlined /> },
  timeline: { label: '时间线一致性', icon: <ClockCircleOutlined /> },
  worldview: { label: '世界观一致性', icon: <GlobalOutlined /> },
}

// 严重程度映射
const severityConfig: Record<string, { color: string; label: string }> = {
  critical: { color: 'red', label: '严重' },
  major: { color: 'orange', label: '重要' },
  minor: { color: 'blue', label: '轻微' },
}

// 评分颜色
const getScoreColor = (score: number): string => {
  if (score >= 8) return '#52c41a'
  if (score >= 6) return '#faad14'
  return '#f5222d'
}

const ConsistencyChecker: React.FC<ConsistencyCheckerProps> = ({
  chapterId,
  content,
  characters = [],
  worldSettings = [],
  previousSummary = '',
  onCheck,
}) => {
  const [loading, setLoading] = useState(false)
  const [result, setResult] = useState<ConsistencyResult | null>(null)

  const handleCheck = async () => {
    if (!content.trim()) return

    setLoading(true)
    try {
      const data: QualityRequest = {
        content,
        characters: characters.map((c) => ({
          name: c.name,
          personality: c.personality,
          background: c.background,
        })),
        world_settings: worldSettings.map((s) => ({
          title: s.title,
          content: s.content,
        })),
        previous_summary: previousSummary,
      }
      const res = await qualityService.checkConsistency(chapterId, data)
      setResult(res)
      onCheck?.(res)
    } catch (error) {
      console.error('一致性检查失败', error)
    } finally {
      setLoading(false)
    }
  }

  const renderIssue = (issue: ConsistencyIssue) => {
    const typeConfig = issueTypeLabels[issue.type] || { label: issue.type, icon: null }
    const severity = severityConfig[issue.severity] || severityConfig.minor

    return (
      <List.Item>
        <div className="w-full">
          <div className="flex items-center justify-between mb-2">
            <Space>
              {typeConfig.icon}
              <Tag>{typeConfig.label}</Tag>
              <Tag color={severity.color}>{severity.label}</Tag>
            </Space>
          </div>
          <div className="text-sm mb-2">{issue.description}</div>
          {issue.suggestion && (
            <Alert
              type="info"
              message={issue.suggestion}
              showIcon
              className="text-xs"
            />
          )}
        </div>
      </List.Item>
    )
  }

  const getOverallStatus = () => {
    if (!result) return null
    if (result.overall_score >= 8) {
      return { icon: <CheckCircleOutlined />, color: '#52c41a', text: '一致性良好' }
    }
    if (result.overall_score >= 6) {
      return { icon: <WarningOutlined />, color: '#faad14', text: '存在一些问题' }
    }
    return { icon: <CloseCircleOutlined />, color: '#f5222d', text: '一致性较差' }
  }

  const status = getOverallStatus()

  return (
    <Card
      title={
        <Space>
          <CheckCircleOutlined />
          <span>一致性检查</span>
        </Space>
      }
      extra={
        <Button
          type="primary"
          onClick={handleCheck}
          loading={loading}
          disabled={!content.trim()}
        >
          开始检查
        </Button>
      }
    >
      {loading ? (
        <div className="flex items-center justify-center py-12">
          <Spin tip="正在检查一致性..." />
        </div>
      ) : result ? (
        <div className="space-y-6">
          {/* 综合评分 */}
          <div className="text-center">
            <div className="text-4xl font-bold" style={{ color: status?.color }}>
              {result.overall_score.toFixed(1)}
            </div>
            <div className="flex items-center justify-center gap-2 mt-2">
              <span style={{ color: status?.color }}>{status?.icon}</span>
              <span className="text-gray-500">{status?.text}</span>
            </div>
          </div>

          {/* 维度评分 */}
          <div className="space-y-3">
            <div>
              <div className="flex justify-between text-sm mb-1">
                <Space>
                  <UserOutlined />
                  <span>角色一致性</span>
                </Space>
                <span style={{ color: getScoreColor(result.character_score) }}>
                  {result.character_score.toFixed(1)}
                </span>
              </div>
              <Progress
                percent={result.character_score * 10}
                strokeColor={getScoreColor(result.character_score)}
                showInfo={false}
                size="small"
              />
            </div>

            <div>
              <div className="flex justify-between text-sm mb-1">
                <Space>
                  <ClockCircleOutlined />
                  <span>时间线一致性</span>
                </Space>
                <span style={{ color: getScoreColor(result.timeline_score) }}>
                  {result.timeline_score.toFixed(1)}
                </span>
              </div>
              <Progress
                percent={result.timeline_score * 10}
                strokeColor={getScoreColor(result.timeline_score)}
                showInfo={false}
                size="small"
              />
            </div>

            <div>
              <div className="flex justify-between text-sm mb-1">
                <Space>
                  <GlobalOutlined />
                  <span>世界观一致性</span>
                </Space>
                <span style={{ color: getScoreColor(result.worldview_score) }}>
                  {result.worldview_score.toFixed(1)}
                </span>
              </div>
              <Progress
                percent={result.worldview_score * 10}
                strokeColor={getScoreColor(result.worldview_score)}
                showInfo={false}
                size="small"
              />
            </div>
          </div>

          {/* 问题列表 */}
          {result.issues.length > 0 ? (
            <div>
              <h4 className="font-medium mb-2">
                发现 {result.issues.length} 个问题
              </h4>
              <List
                size="small"
                dataSource={result.issues}
                renderItem={renderIssue}
              />
            </div>
          ) : (
            <Alert
              type="success"
              message="未发现一致性问题"
              description="章节内容与角色设定、时间线、世界观保持一致"
              showIcon
            />
          )}
        </div>
      ) : (
        <Empty
          description={
            <div>
              <p>点击"开始检查"按钮检查章节一致性</p>
              <p className="text-xs text-gray-400 mt-2">
                提示：添加角色和世界设定可以获得更准确的检查结果
              </p>
            </div>
          }
        />
      )}
    </Card>
  )
}

export default ConsistencyChecker
