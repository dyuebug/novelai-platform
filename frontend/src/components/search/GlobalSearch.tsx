import { useState, useEffect, useCallback } from 'react'
import { Modal, Input, List, Tag, Empty, Spin } from 'antd'
import { SearchOutlined, FileTextOutlined, TeamOutlined, EnvironmentOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { useDebounce } from '@/hooks/useDebounce'
import type { FC } from 'react'

const { Search } = Input

interface SearchResult {
  id: string
  type: 'chapter' | 'character' | 'location' | 'world-setting'
  title: string
  content: string
  projectId: string
  projectTitle: string
  highlight?: string
}

interface GlobalSearchProps {
  visible: boolean
  onClose: () => void
}

const typeIcons = {
  chapter: <FileTextOutlined />,
  character: <TeamOutlined />,
  location: <EnvironmentOutlined />,
  'world-setting': <EnvironmentOutlined />,
}

const typeLabels = {
  chapter: '章节',
  character: '角色',
  location: '地点',
  'world-setting': '世界设定',
}

const typeColors = {
  chapter: 'blue',
  character: 'green',
  location: 'orange',
  'world-setting': 'purple',
}

export const GlobalSearch: FC<GlobalSearchProps> = ({ visible, onClose }) => {
  const navigate = useNavigate()
  const [searchText, setSearchText] = useState('')
  const [results, setResults] = useState<SearchResult[]>([])
  const [loading, setLoading] = useState(false)
  const [searchHistory, setSearchHistory] = useState<string[]>([])

  const debouncedSearchText = useDebounce(searchText, 300)

  // 加载搜索历史
  useEffect(() => {
    const history = localStorage.getItem('search-history')
    if (history) {
      setSearchHistory(JSON.parse(history))
    }
  }, [])

  // 保存搜索历史
  const saveSearchHistory = useCallback((query: string) => {
    const newHistory = [query, ...searchHistory.filter(h => h !== query)].slice(0, 10)
    setSearchHistory(newHistory)
    localStorage.setItem('search-history', JSON.stringify(newHistory))
  }, [searchHistory])

  // 执行搜索
  useEffect(() => {
    if (!debouncedSearchText) {
      setResults([])
      return
    }

    const performSearch = async () => {
      setLoading(true)
      try {
        // TODO: 调用实际API
        // const response = await searchService.globalSearch(debouncedSearchText)
        // setResults(response.data)

        // 模拟搜索结果
        await new Promise(resolve => setTimeout(resolve, 500))
        setResults([])
      } catch (error) {
        console.error('Search failed:', error)
      } finally {
        setLoading(false)
      }
    }

    performSearch()
  }, [debouncedSearchText])

  const handleResultClick = (result: SearchResult) => {
    saveSearchHistory(searchText)

    // 导航到结果详情页
    switch (result.type) {
      case 'chapter':
        navigate(`/projects/${result.projectId}/chapters/${result.id}`)
        break
      case 'character':
        navigate(`/projects/${result.projectId}/characters/${result.id}`)
        break
      case 'location':
        navigate(`/projects/${result.projectId}/locations/${result.id}`)
        break
      case 'world-setting':
        navigate(`/projects/${result.projectId}/world-settings/${result.id}`)
        break
    }

    onClose()
  }

  const handleHistoryClick = (query: string) => {
    setSearchText(query)
  }

  const clearHistory = () => {
    setSearchHistory([])
    localStorage.removeItem('search-history')
  }

  // 按类型分组结果
  const groupedResults = results.reduce((acc, result) => {
    if (!acc[result.type]) {
      acc[result.type] = []
    }
    acc[result.type].push(result)
    return acc
  }, {} as Record<string, SearchResult[]>)

  return (
    <Modal
      title={null}
      open={visible}
      onCancel={onClose}
      footer={null}
      width={700}
      bodyStyle={{ padding: 0 }}
      closeIcon={null}
    >
      <div className="p-4 border-b">
        <Search
          placeholder="搜索章节、角色、地点、世界设定..."
          prefix={<SearchOutlined />}
          value={searchText}
          onChange={(e) => setSearchText(e.target.value)}
          size="large"
          autoFocus
          allowClear
        />
      </div>

      <div className="max-h-[500px] overflow-auto">
        {loading && (
          <div className="flex justify-center py-8">
            <Spin />
          </div>
        )}

        {!loading && !searchText && searchHistory.length > 0 && (
          <div className="p-4">
            <div className="flex items-center justify-between mb-2">
              <span className="text-sm text-gray-500">搜索历史</span>
              <a onClick={clearHistory} className="text-xs">清除</a>
            </div>
            <div className="flex flex-wrap gap-2">
              {searchHistory.map((query, index) => (
                <Tag
                  key={index}
                  className="cursor-pointer"
                  onClick={() => handleHistoryClick(query)}
                >
                  {query}
                </Tag>
              ))}
            </div>
          </div>
        )}

        {!loading && searchText && results.length === 0 && (
          <Empty
            description="未找到匹配结果"
            className="py-8"
          />
        )}

        {!loading && results.length > 0 && (
          <div className="p-4 space-y-4">
            {Object.entries(groupedResults).map(([type, items]) => (
              <div key={type}>
                <div className="text-sm font-medium text-gray-500 mb-2">
                  {typeLabels[type as keyof typeof typeLabels]} ({items.length})
                </div>
                <List
                  size="small"
                  dataSource={items}
                  renderItem={(item) => (
                    <List.Item
                      className="cursor-pointer hover:bg-gray-50 px-3 py-2 rounded"
                      onClick={() => handleResultClick(item)}
                    >
                      <div className="flex-1">
                        <div className="flex items-center gap-2 mb-1">
                          {typeIcons[item.type]}
                          <span className="font-medium">{item.title}</span>
                          <Tag color={typeColors[item.type]} className="text-xs">
                            {item.projectTitle}
                          </Tag>
                        </div>
                        {item.highlight && (
                          <div
                            className="text-sm text-gray-600 line-clamp-2"
                            dangerouslySetInnerHTML={{ __html: item.highlight }}
                          />
                        )}
                      </div>
                    </List.Item>
                  )}
                />
              </div>
            ))}
          </div>
        )}
      </div>

      <div className="p-3 border-t bg-gray-50 text-xs text-gray-500 flex items-center justify-between">
        <span>按 Enter 跳转，Esc 关闭</span>
        <span>Ctrl+K 打开搜索</span>
      </div>
    </Modal>
  )
}
