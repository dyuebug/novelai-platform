import React, { useState, useEffect, useRef } from 'react';
import { Card, Spin, Empty, Tooltip, Tag, Select, Space } from 'antd';
import {
  FlagOutlined,
  BulbOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
} from '@ant-design/icons';
import foreshadowService, {
  Foreshadow,
  ForeshadowStatus,
  ForeshadowPriority,
} from '../../services/foreshadow.service';

interface ForeshadowTimelineProps {
  projectId: string;
  currentChapter?: number;
  maxChapter?: number;
  onForeshadowClick?: (foreshadow: Foreshadow) => void;
}

const statusIcons: Record<ForeshadowStatus, React.ReactNode> = {
  planted: <FlagOutlined />,
  hinted: <BulbOutlined />,
  resolved: <CheckCircleOutlined />,
  abandoned: <CloseCircleOutlined />,
};

const statusColors: Record<ForeshadowStatus, string> = {
  planted: '#1890ff',
  hinted: '#fa8c16',
  resolved: '#52c41a',
  abandoned: '#8c8c8c',
};

const priorityColors: Record<ForeshadowPriority, string> = {
  high: '#ff4d4f',
  medium: '#fa8c16',
  low: '#1890ff',
};

const ForeshadowTimeline: React.FC<ForeshadowTimelineProps> = ({
  projectId,
  currentChapter = 0,
  maxChapter = 100,
  onForeshadowClick,
}) => {
  const [foreshadows, setForeshadows] = useState<Foreshadow[]>([]);
  const [loading, setLoading] = useState(false);
  const [filterStatus, setFilterStatus] = useState<ForeshadowStatus | 'all'>('all');
  const [filterPriority, setFilterPriority] = useState<ForeshadowPriority | 'all'>('all');
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const loadForeshadows = async () => {
      setLoading(true);
      try {
        const response = await foreshadowService.list(projectId, {
          page: 1,
          page_size: 200,
        });
        setForeshadows(response.items || []);
      } catch (error) {
        console.error('加载伏笔失败:', error);
      } finally {
        setLoading(false);
      }
    };
    loadForeshadows();
  }, [projectId]);

  // 过滤伏笔
  const filteredForeshadows = foreshadows.filter((f) => {
    if (filterStatus !== 'all' && f.status !== filterStatus) return false;
    if (filterPriority !== 'all' && f.priority !== filterPriority) return false;
    return true;
  });

  // 计算章节范围
  const minChapterNum = Math.min(
    ...filteredForeshadows
      .map((f) => f.plant_chapter_num)
      .filter((n) => n > 0),
    currentChapter || 1
  );
  const maxChapterNum = Math.max(
    ...filteredForeshadows
      .map((f) => Math.max(f.plant_chapter_num, f.resolve_chapter_num, f.remind_chapter_num))
      .filter((n) => n > 0),
    currentChapter || 1,
    maxChapter
  );

  const chapterRange = maxChapterNum - minChapterNum + 1;
  const chapterWidth = 60; // 每章节宽度
  const rowHeight = 40; // 每行高度
  const headerHeight = 40;
  const padding = 20;

  // 按行分配伏笔（避免重叠）
  const assignRows = (items: Foreshadow[]): Map<string, number> => {
    const rowMap = new Map<string, number>();
    const rowEndChapters: number[] = [];

    const sortedItems = [...items].sort((a, b) => a.plant_chapter_num - b.plant_chapter_num);

    for (const item of sortedItems) {
      const startChapter = item.plant_chapter_num || 1;
      const endChapter = item.resolve_chapter_num || item.remind_chapter_num || startChapter + 10;

      let assignedRow = -1;
      for (let i = 0; i < rowEndChapters.length; i++) {
        if (rowEndChapters[i] < startChapter - 1) {
          assignedRow = i;
          rowEndChapters[i] = endChapter;
          break;
        }
      }

      if (assignedRow === -1) {
        assignedRow = rowEndChapters.length;
        rowEndChapters.push(endChapter);
      }

      rowMap.set(item.id, assignedRow);
    }

    return rowMap;
  };

  const rowAssignments = assignRows(filteredForeshadows);
  const totalRows = Math.max(...Array.from(rowAssignments.values()), 0) + 1;
  const svgHeight = headerHeight + totalRows * rowHeight + padding * 2;
  const svgWidth = chapterRange * chapterWidth + padding * 2;

  // 章节号转X坐标
  const chapterToX = (chapter: number) => padding + (chapter - minChapterNum) * chapterWidth;

  if (loading) {
    return (
      <Card title="伏笔时间线">
        <div style={{ textAlign: 'center', padding: 40 }}>
          <Spin />
        </div>
      </Card>
    );
  }

  if (foreshadows.length === 0) {
    return (
      <Card title="伏笔时间线">
        <Empty description="暂无伏笔数据" />
      </Card>
    );
  }

  return (
    <Card
      title="伏笔时间线"
      extra={
        <Space>
          <Select
            value={filterStatus}
            onChange={setFilterStatus}
            style={{ width: 100 }}
            size="small"
          >
            <Select.Option value="all">全部状态</Select.Option>
            <Select.Option value="planted">已埋设</Select.Option>
            <Select.Option value="hinted">已暗示</Select.Option>
            <Select.Option value="resolved">已回收</Select.Option>
            <Select.Option value="abandoned">已放弃</Select.Option>
          </Select>
          <Select
            value={filterPriority}
            onChange={setFilterPriority}
            style={{ width: 100 }}
            size="small"
          >
            <Select.Option value="all">全部优先级</Select.Option>
            <Select.Option value="high">高优先级</Select.Option>
            <Select.Option value="medium">中优先级</Select.Option>
            <Select.Option value="low">低优先级</Select.Option>
          </Select>
        </Space>
      }
    >
      <div
        ref={containerRef}
        style={{
          overflowX: 'auto',
          overflowY: 'hidden',
        }}
      >
        <svg width={svgWidth} height={svgHeight}>
          {/* 背景网格 */}
          <defs>
            <pattern
              id="grid"
              width={chapterWidth}
              height={rowHeight}
              patternUnits="userSpaceOnUse"
            >
              <path
                d={`M ${chapterWidth} 0 L 0 0 0 ${rowHeight}`}
                fill="none"
                stroke="#f0f0f0"
                strokeWidth="1"
              />
            </pattern>
          </defs>
          <rect
            x={padding}
            y={headerHeight}
            width={chapterRange * chapterWidth}
            height={totalRows * rowHeight}
            fill="url(#grid)"
          />

          {/* 章节刻度 */}
          {Array.from({ length: chapterRange }, (_, i) => {
            const chapter = minChapterNum + i;
            const x = chapterToX(chapter);
            const isCurrentChapter = chapter === currentChapter;
            return (
              <g key={`chapter-${chapter}`}>
                <text
                  x={x + chapterWidth / 2}
                  y={headerHeight - 10}
                  textAnchor="middle"
                  fontSize={10}
                  fill={isCurrentChapter ? '#1890ff' : '#8c8c8c'}
                  fontWeight={isCurrentChapter ? 'bold' : 'normal'}
                >
                  {chapter}
                </text>
                {isCurrentChapter && (
                  <line
                    x1={x + chapterWidth / 2}
                    y1={headerHeight}
                    x2={x + chapterWidth / 2}
                    y2={svgHeight - padding}
                    stroke="#1890ff"
                    strokeWidth={2}
                    strokeDasharray="4,4"
                  />
                )}
              </g>
            );
          })}

          {/* 伏笔条 */}
          {filteredForeshadows.map((foreshadow) => {
            const row = rowAssignments.get(foreshadow.id) || 0;
            const startChapter = foreshadow.plant_chapter_num || minChapterNum;
            const endChapter =
              foreshadow.resolve_chapter_num ||
              foreshadow.remind_chapter_num ||
              startChapter + 5;

            const x = chapterToX(startChapter);
            const y = headerHeight + row * rowHeight + 8;
            const width = Math.max((endChapter - startChapter + 1) * chapterWidth - 4, chapterWidth - 4);
            const height = rowHeight - 16;

            const barColor = statusColors[foreshadow.status];
            const borderColor = priorityColors[foreshadow.priority];

            return (
              <Tooltip
                key={foreshadow.id}
                title={
                  <div>
                    <div style={{ fontWeight: 'bold' }}>{foreshadow.title}</div>
                    <div>状态: {foreshadow.status}</div>
                    <div>优先级: {foreshadow.priority}</div>
                    <div>埋设: 第{startChapter}章</div>
                    {foreshadow.resolve_chapter_num > 0 && (
                      <div>回收: 第{foreshadow.resolve_chapter_num}章</div>
                    )}
                    {foreshadow.remind_chapter_num > 0 && (
                      <div>提醒: 第{foreshadow.remind_chapter_num}章</div>
                    )}
                  </div>
                }
              >
                <g
                  style={{ cursor: 'pointer' }}
                  onClick={() => onForeshadowClick?.(foreshadow)}
                >
                  {/* 伏笔条背景 */}
                  <rect
                    x={x}
                    y={y}
                    width={width}
                    height={height}
                    rx={4}
                    ry={4}
                    fill={barColor}
                    fillOpacity={0.2}
                    stroke={borderColor}
                    strokeWidth={2}
                  />
                  {/* 伏笔标题 */}
                  <text
                    x={x + 8}
                    y={y + height / 2 + 4}
                    fontSize={11}
                    fill="#333"
                    style={{ pointerEvents: 'none' }}
                  >
                    {foreshadow.title.length > 10
                      ? foreshadow.title.slice(0, 10) + '...'
                      : foreshadow.title}
                  </text>
                  {/* 状态图标 */}
                  <foreignObject x={x + width - 20} y={y + 2} width={16} height={16}>
                    <div style={{ color: barColor, fontSize: 12 }}>
                      {statusIcons[foreshadow.status]}
                    </div>
                  </foreignObject>
                </g>
              </Tooltip>
            );
          })}
        </svg>
      </div>

      {/* 图例 */}
      <div style={{ marginTop: 16, display: 'flex', gap: 16, flexWrap: 'wrap' }}>
        <span style={{ color: '#8c8c8c' }}>状态:</span>
        <Tag color="blue" icon={<FlagOutlined />}>
          已埋设
        </Tag>
        <Tag color="orange" icon={<BulbOutlined />}>
          已暗示
        </Tag>
        <Tag color="green" icon={<CheckCircleOutlined />}>
          已回收
        </Tag>
        <Tag icon={<CloseCircleOutlined />}>已放弃</Tag>
        <span style={{ marginLeft: 16, color: '#8c8c8c' }}>优先级边框:</span>
        <span style={{ color: '#ff4d4f' }}>■ 高</span>
        <span style={{ color: '#fa8c16' }}>■ 中</span>
        <span style={{ color: '#1890ff' }}>■ 低</span>
      </div>
    </Card>
  );
};

export default ForeshadowTimeline;
