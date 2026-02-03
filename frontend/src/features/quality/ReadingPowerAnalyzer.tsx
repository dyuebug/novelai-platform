import { useState } from 'react'
import {
  Card,
  Button,
  Progress,
  Tag,
  List,
  Space,
  Spin,
  Empty,
  Tooltip,
  Collapse,
} from 'antd'
import {
  ThunderboltOutlined,
  FireOutlined,
  GiftOutlined,
  BulbOutlined,
} from '@ant-design/icons'
import {
  qualityService,
  ReadingPowerResult,
  QualityRequest,
} from '@/services/quality.service'

interface ReadingPowerAnalyzerProps {
  chapterId: string
  content: string
  onAnalyze?: (result: ReadingPowerResult) => void
}

// 钩子类型映射
const hookTypeLabels: Record<string, string> = {
  crisis: '危机钩子',
  mystery: '悬念钩子',
  emotion: '情感钩子',
  choice: '抉择钩子',
  desire: '欲望钩子',
}

// 爽点类型映射
const coolpointTypeLabels: Record<string, string> = {
  achievement: '成就爽点',
  revenge: '复仇爽点',
  romance: '情感爽点',
  power_up: '升级爽点',
  recognition: '认可爽点',
}

// 评分颜色
const getScoreColor = (score: number): string => {
  if (score >= 8) return '#52c41a'
  if (score >= 6) return '#faad14'
  return '#f5222d'
}

const ReadingPowerAnalyzer: React.FC<ReadingPowerAnalyzerProps> = ({
  chapterId,
  content,
  onAnalyze,
}) => {
  const [loading, setLoading] = useState(false)
  const [result, setResult] = useState<ReadingPowerResult | null>(null)

  const handleAnalyze = async () => {
    if (!content.trim()) return

    setLoading(true)
    try {
      const data: QualityRequest = { content }
      const res = await qualityService.analyzeReadingPower(chapterId, data)
      setResult(res)
      onAnalyze?.(res)
    } catch (error) {
      console.error('追读力分析失败', error)
    } finally {
      setLoading(false)
    }
  }

  return (
    <Card
      title={
        <Space>
          <ThunderboltOutlined />
          <span>追读力分析</span>
        </Space>
      }
      extra={
        <Button
          type="primary"
          onClick={handleAnalyze}
          loading={loading}
          disabled={!content.trim()}
        >
          开始分析
        </Button>
      }
    >
      {loading ? (
        <div className="flex items-center justify-center py-12">
          <Spin tip="正在分析追读力..." />
        </div>
      ) : result ? (
        <div className="space-y-6">
          {/* 综合评分 */}
          <div className="text-center">
            <div className="text-4xl font-bold" style={{ color: getScoreColor(result.overall_score) }}>
              {result.overall_score.toFixed(1)}
            </div>
            <div className="text-gray-500 mt-1">综合追读力评分</div>
          </div>

          {/* 维度评分 */}
          <div className="grid grid-cols-3 gap-4">
            <div className="text-center p-4 bg-gray-50 rounded-lg">
              <Tooltip title="章节开头的吸引力">
                <div className="text-2xl font-semibold" style={{ color: getScoreColor(result.dimension_scores.hook || 0) }}>
                  {(result.dimension_scores.hook || 0).toFixed(1)}
                </div>
              </Tooltip>
              <div className="text-gray-500 text-sm mt-1">钩子</div>
            </div>
            <div className="text-center p-4 bg-gray-50 rounded-lg">
              <Tooltip title="让读者感到爽快的情节">
                <div className="text-2xl font-semibold" style={{ color: getScoreColor(result.dimension_scores.coolpoint || 0) }}>
                  {(result.dimension_scores.coolpoint || 0).toFixed(1)}
                </div>
              </Tooltip>
              <div className="text-gray-500 text-sm mt-1">爽点</div>
            </div>
            <div className="text-center p-4 bg-gray-50 rounded-lg">
              <Tooltip title="小目标的及时满足">
                <div className="text-2xl font-semibold" style={{ color: getScoreColor(result.dimension_scores.micropayoff || 0) }}>
                  {(result.dimension_scores.micropayoff || 0).toFixed(1)}
                </div>
              </Tooltip>
              <div className="text-gray-500 text-sm mt-1">微兑现</div>
            </div>
          </div>

          <Collapse
            items={[
              // 钩子分析
              {
                key: 'hook',
                label: (
                  <Space>
                    <BulbOutlined />
                    <span>钩子分析</span>
                    {result.hook && (
                      <Tag color={result.hook.hook_strength === 'strong' ? 'green' : result.hook.hook_strength === 'medium' ? 'orange' : 'red'}>
                        {result.hook.hook_strength === 'strong' ? '强' : result.hook.hook_strength === 'medium' ? '中' : '弱'}
                      </Tag>
                    )}
                  </Space>
                ),
                children: result.hook ? (
                  <div className="space-y-3">
                    <div>
                      <span className="text-gray-500">类型：</span>
                      <Tag color="blue">{hookTypeLabels[result.hook.hook_type] || result.hook.hook_type}</Tag>
                    </div>
                    <div>
                      <span className="text-gray-500">内容：</span>
                      <span className="italic">"{result.hook.hook_content}"</span>
                    </div>
                    {result.hook.suggestions.length > 0 && (
                      <div>
                        <span className="text-gray-500">建议：</span>
                        <ul className="list-disc list-inside mt-1">
                          {result.hook.suggestions.map((s, i) => (
                            <li key={i} className="text-sm">{s}</li>
                          ))}
                        </ul>
                      </div>
                    )}
                  </div>
                ) : (
                  <Empty description="未检测到钩子" />
                ),
              },
              // 爽点分析
              {
                key: 'coolpoints',
                label: (
                  <Space>
                    <FireOutlined />
                    <span>爽点分析</span>
                    <Tag>{result.coolpoints.length} 个</Tag>
                  </Space>
                ),
                children: result.coolpoints.length > 0 ? (
                  <List
                    size="small"
                    dataSource={result.coolpoints}
                    renderItem={(cp) => (
                      <List.Item>
                        <div className="w-full">
                          <div className="flex justify-between items-center">
                            <Space>
                              <Tag color="orange">{coolpointTypeLabels[cp.type] || cp.type}</Tag>
                              <span className="text-sm">{cp.content}</span>
                            </Space>
                            <span className="text-gray-400 text-xs">位置 {cp.position}%</span>
                          </div>
                          <Progress
                            percent={cp.intensity * 10}
                            size="small"
                            strokeColor={getScoreColor(cp.intensity)}
                            format={() => `${cp.intensity.toFixed(1)}`}
                          />
                        </div>
                      </List.Item>
                    )}
                  />
                ) : (
                  <Empty description="未检测到爽点" />
                ),
              },
              // 微兑现分析
              {
                key: 'micropayoffs',
                label: (
                  <Space>
                    <GiftOutlined />
                    <span>微兑现分析</span>
                    <Tag>{result.micropayoffs.length} 个</Tag>
                  </Space>
                ),
                children: result.micropayoffs.length > 0 ? (
                  <List
                    size="small"
                    dataSource={result.micropayoffs}
                    renderItem={(mp) => (
                      <List.Item>
                        <div className="w-full">
                          <div className="flex justify-between items-center mb-1">
                            <span className="text-gray-500 text-sm">目标：{mp.goal}</span>
                            <span className="text-gray-400 text-xs">位置 {mp.position}%</span>
                          </div>
                          <div className="text-sm">兑现：{mp.payoff}</div>
                          <Progress
                            percent={mp.satisfaction * 10}
                            size="small"
                            strokeColor={getScoreColor(mp.satisfaction)}
                            format={() => `满足度 ${mp.satisfaction.toFixed(1)}`}
                          />
                        </div>
                      </List.Item>
                    )}
                  />
                ) : (
                  <Empty description="未检测到微兑现" />
                ),
              },
            ]}
            defaultActiveKey={['hook']}
          />

          {/* 改进建议 */}
          {result.suggestions.length > 0 && (
            <div>
              <h4 className="font-medium mb-2">改进建议</h4>
              <List
                size="small"
                dataSource={result.suggestions}
                renderItem={(suggestion) => (
                  <List.Item>
                    <BulbOutlined className="text-yellow-500 mr-2" />
                    {suggestion}
                  </List.Item>
                )}
              />
            </div>
          )}
        </div>
      ) : (
        <Empty description='点击"开始分析"按钮分析章节追读力' />
      )}
    </Card>
  )
}

export default ReadingPowerAnalyzer
