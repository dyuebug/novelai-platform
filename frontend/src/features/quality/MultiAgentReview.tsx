import { useState } from 'react'
import {
  Card,
  Button,
  Tag,
  List,
  Space,
  Spin,
  Empty,
  Avatar,
  Collapse,
  Rate,
} from 'antd'
import {
  TeamOutlined,
  UserOutlined,
  EditOutlined,
  EyeOutlined,
  TrophyOutlined,
  RiseOutlined,
  HeartOutlined,
  CheckOutlined,
  CloseOutlined,
  BulbOutlined,
} from '@ant-design/icons'
import {
  qualityService,
  MultiAgentResult,
  AgentReview,
  QualityRequest,
} from '@/services/quality.service'

interface MultiAgentReviewProps {
  chapterId: string
  content: string
  onReview?: (result: MultiAgentResult) => void
}

// Agent 角色配置
const agentConfig: Record<string, { label: string; icon: React.ReactNode; color: string }> = {
  reader: { label: '普通读者', icon: <UserOutlined />, color: '#1890ff' },
  editor: { label: '资深编辑', icon: <EditOutlined />, color: '#722ed1' },
  critic: { label: '文学评论家', icon: <EyeOutlined />, color: '#13c2c2' },
  writer: { label: '网文作者', icon: <TrophyOutlined />, color: '#fa8c16' },
  marketer: { label: '市场分析师', icon: <RiseOutlined />, color: '#52c41a' },
  psychologist: { label: '读者心理专家', icon: <HeartOutlined />, color: '#eb2f96' },
}

// 评分颜色
const getScoreColor = (score: number): string => {
  if (score >= 8) return '#52c41a'
  if (score >= 6) return '#faad14'
  return '#f5222d'
}

const MultiAgentReviewComponent: React.FC<MultiAgentReviewProps> = ({
  chapterId,
  content,
  onReview,
}) => {
  const [loading, setLoading] = useState(false)
  const [result, setResult] = useState<MultiAgentResult | null>(null)

  const handleReview = async () => {
    if (!content.trim()) return

    setLoading(true)
    try {
      const data: QualityRequest = { content }
      const res = await qualityService.multiAgentReview(chapterId, data)
      setResult(res)
      onReview?.(res)
    } catch (error) {
      console.error('多Agent审查失败', error)
    } finally {
      setLoading(false)
    }
  }

  const renderAgentReview = (review: AgentReview) => {
    return (
      <div className="space-y-3">
        {/* 评分 */}
        <div className="flex items-center justify-between">
          <span className="text-gray-500">评分</span>
          <Space>
            <Rate disabled value={review.score / 2} allowHalf />
            <span
              className="font-semibold"
              style={{ color: getScoreColor(review.score) }}
            >
              {review.score.toFixed(1)}
            </span>
          </Space>
        </div>

        {/* 优点 */}
        {review.strengths.length > 0 && (
          <div>
            <div className="text-sm text-gray-500 mb-1">
              <CheckOutlined className="text-green-500 mr-1" />
              优点
            </div>
            <List
              size="small"
              dataSource={review.strengths}
              renderItem={(item) => (
                <List.Item className="py-1 px-0 border-0">
                  <span className="text-sm text-green-600">• {item}</span>
                </List.Item>
              )}
            />
          </div>
        )}

        {/* 不足 */}
        {review.weaknesses.length > 0 && (
          <div>
            <div className="text-sm text-gray-500 mb-1">
              <CloseOutlined className="text-red-500 mr-1" />
              不足
            </div>
            <List
              size="small"
              dataSource={review.weaknesses}
              renderItem={(item) => (
                <List.Item className="py-1 px-0 border-0">
                  <span className="text-sm text-red-600">• {item}</span>
                </List.Item>
              )}
            />
          </div>
        )}

        {/* 建议 */}
        {review.suggestions.length > 0 && (
          <div>
            <div className="text-sm text-gray-500 mb-1">
              <BulbOutlined className="text-yellow-500 mr-1" />
              建议
            </div>
            <List
              size="small"
              dataSource={review.suggestions}
              renderItem={(item) => (
                <List.Item className="py-1 px-0 border-0">
                  <span className="text-sm">• {item}</span>
                </List.Item>
              )}
            />
          </div>
        )}
      </div>
    )
  }

  return (
    <Card
      title={
        <Space>
          <TeamOutlined />
          <span>多Agent审查</span>
        </Space>
      }
      extra={
        <Button
          type="primary"
          onClick={handleReview}
          loading={loading}
          disabled={!content.trim()}
        >
          开始审查
        </Button>
      }
    >
      {loading ? (
        <div className="flex items-center justify-center py-12">
          <Spin tip="6位专家正在审查中..." />
        </div>
      ) : result ? (
        <div className="space-y-6">
          {/* 综合评分 */}
          <div className="text-center">
            <div
              className="text-4xl font-bold"
              style={{ color: getScoreColor(result.overall_score) }}
            >
              {result.overall_score.toFixed(1)}
            </div>
            <div className="text-gray-500 mt-1">6位专家综合评分</div>
          </div>

          {/* 共识总结 */}
          <div className="grid grid-cols-2 gap-4">
            {/* 共识优点 */}
            <div className="bg-green-50 p-4 rounded-lg">
              <h4 className="font-medium text-green-700 mb-2">
                <CheckOutlined className="mr-1" />
                共识优点
              </h4>
              {result.consensus_strengths.length > 0 ? (
                <ul className="text-sm text-green-600 space-y-1">
                  {result.consensus_strengths.map((s, i) => (
                    <li key={i}>• {s}</li>
                  ))}
                </ul>
              ) : (
                <span className="text-sm text-gray-400">暂无共识</span>
              )}
            </div>

            {/* 共识不足 */}
            <div className="bg-red-50 p-4 rounded-lg">
              <h4 className="font-medium text-red-700 mb-2">
                <CloseOutlined className="mr-1" />
                共识不足
              </h4>
              {result.consensus_weaknesses.length > 0 ? (
                <ul className="text-sm text-red-600 space-y-1">
                  {result.consensus_weaknesses.map((w, i) => (
                    <li key={i}>• {w}</li>
                  ))}
                </ul>
              ) : (
                <span className="text-sm text-gray-400">暂无共识</span>
              )}
            </div>
          </div>

          {/* 最终建议 */}
          {result.final_suggestions.length > 0 && (
            <div className="bg-yellow-50 p-4 rounded-lg">
              <h4 className="font-medium text-yellow-700 mb-2">
                <BulbOutlined className="mr-1" />
                综合建议
              </h4>
              <ul className="text-sm space-y-1">
                {result.final_suggestions.map((s, i) => (
                  <li key={i}>• {s}</li>
                ))}
              </ul>
            </div>
          )}

          {/* 各Agent详细审查 */}
          <Collapse
            items={result.reviews.map((review) => {
              const config = agentConfig[review.role] || {
                label: review.agent_name,
                icon: <UserOutlined />,
                color: '#666',
              }
              return {
                key: review.role,
                label: (
                  <Space>
                    <Avatar
                      size="small"
                      icon={config.icon}
                      style={{ backgroundColor: config.color }}
                    />
                    <span>{config.label}</span>
                    <Tag color={getScoreColor(review.score)}>
                      {review.score.toFixed(1)}
                    </Tag>
                  </Space>
                ),
                children: renderAgentReview(review),
              }
            })}
          />
        </div>
      ) : (
        <Empty
          description={
            <div>
              <p>点击"开始审查"按钮启动多Agent审查</p>
              <p className="text-xs text-gray-400 mt-2">
                6位虚拟专家将从不同角度审查您的章节
              </p>
            </div>
          }
        />
      )}
    </Card>
  )
}

export default MultiAgentReviewComponent
