import { useState, useEffect, useCallback, useMemo } from 'react'
import { useParams } from 'react-router-dom'
import { Card, Spin, Empty, Select, Space, Tag, message } from 'antd'
import { characterService, Character, CharacterRelationship } from '@/services/world.service'

// 关系类型配置
const relationTypes: Record<string, { label: string; color: string }> = {
  family: { label: '家人', color: '#52c41a' },
  friend: { label: '朋友', color: '#1890ff' },
  enemy: { label: '敌人', color: '#f5222d' },
  lover: { label: '恋人', color: '#eb2f96' },
  mentor: { label: '师徒', color: '#722ed1' },
  colleague: { label: '同事', color: '#13c2c2' },
  rival: { label: '对手', color: '#fa8c16' },
}

// 角色类型颜色
const roleColors: Record<string, string> = {
  protagonist: '#faad14',
  antagonist: '#f5222d',
  supporting: '#1890ff',
  minor: '#8c8c8c',
}

interface GraphNode {
  id: string
  label: string
  role?: string
  style?: {
    fill?: string
    stroke?: string
  }
}

interface GraphEdge {
  source: string
  target: string
  label?: string
  style?: {
    stroke?: string
    lineWidth?: number
  }
}

const RelationshipGraph = () => {
  const { projectId } = useParams<{ projectId: string }>()
  const [characters, setCharacters] = useState<Character[]>([])
  const [allRelationships, setAllRelationships] = useState<Map<string, CharacterRelationship[]>>(new Map())
  const [loading, setLoading] = useState(true)
  const [selectedTypes, setSelectedTypes] = useState<string[]>(Object.keys(relationTypes))

  const fetchData = useCallback(async () => {
    if (!projectId) return
    setLoading(true)
    try {
      // 获取所有角色
      const charResponse = await characterService.list(projectId, { page_size: 200 })
      const chars = charResponse.items || []
      setCharacters(chars)

      // 获取每个角色的关系
      const relationshipMap = new Map<string, CharacterRelationship[]>()
      await Promise.all(
        chars.map(async (char) => {
          try {
            const rels = await characterService.getRelationships(char.id)
            relationshipMap.set(char.id, rels || [])
          } catch {
            relationshipMap.set(char.id, [])
          }
        })
      )
      setAllRelationships(relationshipMap)
    } catch (error) {
      message.error('获取数据失败')
      console.error(error)
    } finally {
      setLoading(false)
    }
  }, [projectId])

  useEffect(() => {
    fetchData()
  }, [fetchData])

  // 构建图数据
  const graphData = useMemo(() => {
    const nodes: GraphNode[] = characters.map((char) => ({
      id: char.id,
      label: char.name,
      role: char.role,
      style: {
        fill: roleColors[char.role || ''] || '#d9d9d9',
        stroke: '#fff',
      },
    }))

    const edges: GraphEdge[] = []
    const addedEdges = new Set<string>()

    allRelationships.forEach((rels, charId) => {
      rels.forEach((rel) => {
        // 过滤关系类型
        if (!selectedTypes.includes(rel.relation_type)) return

        // 避免重复边
        const edgeKey = [charId, rel.target_id].sort().join('-')
        if (addedEdges.has(edgeKey)) return
        addedEdges.add(edgeKey)

        const relConfig = relationTypes[rel.relation_type]
        edges.push({
          source: charId,
          target: rel.target_id,
          label: relConfig?.label || rel.relation_type,
          style: {
            stroke: relConfig?.color || '#999',
            lineWidth: 2,
          },
        })
      })
    })

    return { nodes, edges }
  }, [characters, allRelationships, selectedTypes])

  // 简单的 SVG 关系图渲染
  const renderSimpleGraph = () => {
    const { nodes, edges } = graphData
    if (nodes.length === 0) {
      return <Empty description="暂无角色数据" />
    }

    // 计算节点位置（圆形布局）
    const centerX = 400
    const centerY = 300
    const radius = Math.min(250, 50 + nodes.length * 20)
    const nodePositions = new Map<string, { x: number; y: number }>()

    nodes.forEach((node, index) => {
      const angle = (2 * Math.PI * index) / nodes.length - Math.PI / 2
      nodePositions.set(node.id, {
        x: centerX + radius * Math.cos(angle),
        y: centerY + radius * Math.sin(angle),
      })
    })

    return (
      <svg width="100%" height="600" viewBox="0 0 800 600">
        {/* 绘制边 */}
        {edges.map((edge, index) => {
          const source = nodePositions.get(edge.source)
          const target = nodePositions.get(edge.target)
          if (!source || !target) return null

          const midX = (source.x + target.x) / 2
          const midY = (source.y + target.y) / 2

          return (
            <g key={`edge-${index}`}>
              <line
                x1={source.x}
                y1={source.y}
                x2={target.x}
                y2={target.y}
                stroke={edge.style?.stroke || '#999'}
                strokeWidth={edge.style?.lineWidth || 1}
              />
              {edge.label && (
                <text
                  x={midX}
                  y={midY}
                  textAnchor="middle"
                  fontSize="10"
                  fill={edge.style?.stroke || '#666'}
                  dy="-5"
                >
                  {edge.label}
                </text>
              )}
            </g>
          )
        })}

        {/* 绘制节点 */}
        {nodes.map((node) => {
          const pos = nodePositions.get(node.id)
          if (!pos) return null

          return (
            <g key={node.id}>
              <circle
                cx={pos.x}
                cy={pos.y}
                r="30"
                fill={node.style?.fill || '#d9d9d9'}
                stroke={node.style?.stroke || '#fff'}
                strokeWidth="2"
              />
              <text
                x={pos.x}
                y={pos.y}
                textAnchor="middle"
                dominantBaseline="middle"
                fontSize="12"
                fill="#fff"
                fontWeight="bold"
              >
                {node.label.length > 4 ? node.label.slice(0, 4) + '...' : node.label}
              </text>
            </g>
          )
        })}
      </svg>
    )
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-96">
        <Spin size="large" />
      </div>
    )
  }

  return (
    <Card
      title="角色关系图"
      extra={
        <Space>
          <span>显示关系：</span>
          <Select
            mode="multiple"
            value={selectedTypes}
            onChange={setSelectedTypes}
            style={{ minWidth: 300 }}
            placeholder="选择要显示的关系类型"
            options={Object.entries(relationTypes).map(([value, config]) => ({
              value,
              label: config.label,
            }))}
          />
        </Space>
      }
    >
      {/* 图例 */}
      <div className="mb-4 flex flex-wrap gap-4">
        <div className="flex items-center gap-2">
          <span className="text-gray-500">角色类型：</span>
          <Tag color="gold">主角</Tag>
          <Tag color="red">反派</Tag>
          <Tag color="blue">配角</Tag>
          <Tag color="default">龙套</Tag>
        </div>
        <div className="flex items-center gap-2">
          <span className="text-gray-500">关系类型：</span>
          {Object.entries(relationTypes).map(([key, config]) => (
            <Tag key={key} color={config.color}>
              {config.label}
            </Tag>
          ))}
        </div>
      </div>

      {/* 关系图 */}
      <div className="border rounded-lg bg-gray-50 overflow-hidden">
        {graphData.nodes.length > 0 ? (
          renderSimpleGraph()
        ) : (
          <Empty description="暂无角色或关系数据" className="py-20" />
        )}
      </div>

      {/* 统计信息 */}
      <div className="mt-4 text-gray-500 text-sm">
        共 {graphData.nodes.length} 个角色，{graphData.edges.length} 条关系
      </div>
    </Card>
  )
}

export default RelationshipGraph
