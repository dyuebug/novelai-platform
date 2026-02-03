import { Tabs } from 'antd'
import {
  ThunderboltOutlined,
  CheckCircleOutlined,
  TeamOutlined,
} from '@ant-design/icons'
import ReadingPowerAnalyzer from './ReadingPowerAnalyzer'
import ConsistencyChecker from './ConsistencyChecker'
import MultiAgentReview from './MultiAgentReview'
import { Character, WorldSetting } from '@/services/world.service'

interface QualityPanelProps {
  chapterId: string
  content: string
  characters?: Character[]
  worldSettings?: WorldSetting[]
  previousSummary?: string
}

const QualityPanel: React.FC<QualityPanelProps> = ({
  chapterId,
  content,
  characters = [],
  worldSettings = [],
  previousSummary = '',
}) => {
  const tabItems = [
    {
      key: 'reading-power',
      label: (
        <span>
          <ThunderboltOutlined />
          追读力分析
        </span>
      ),
      children: (
        <ReadingPowerAnalyzer
          chapterId={chapterId}
          content={content}
        />
      ),
    },
    {
      key: 'consistency',
      label: (
        <span>
          <CheckCircleOutlined />
          一致性检查
        </span>
      ),
      children: (
        <ConsistencyChecker
          chapterId={chapterId}
          content={content}
          characters={characters}
          worldSettings={worldSettings}
          previousSummary={previousSummary}
        />
      ),
    },
    {
      key: 'multi-agent',
      label: (
        <span>
          <TeamOutlined />
          多Agent审查
        </span>
      ),
      children: (
        <MultiAgentReview
          chapterId={chapterId}
          content={content}
        />
      ),
    },
  ]

  return (
    <div className="quality-panel">
      <Tabs items={tabItems} />
    </div>
  )
}

export default QualityPanel

// 导出所有组件
export { ReadingPowerAnalyzer, ConsistencyChecker, MultiAgentReview }
