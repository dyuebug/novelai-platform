import { useState, useEffect, useCallback, useRef } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import {
  Card,
  Input,
  Button,
  Spin,
  message,
  Breadcrumb,
  Space,
  Dropdown,
  Tag,
  Tooltip,
} from 'antd'
import {
  SaveOutlined,
  HistoryOutlined,
  ArrowLeftOutlined,
  MoreOutlined,
  CheckCircleOutlined,
  SyncOutlined,
  RobotOutlined,
} from '@ant-design/icons'
import { chapterService, Chapter } from '@/services/chapter.service'
import { projectService, Project } from '@/services/project.service'
import RichTextEditor from '@/components/editor/RichTextEditor'
import VersionHistoryPanel from '@/components/editor/VersionHistoryPanel'
import AIChatPanel from '@/components/editor/AIChatPanel'

const AUTOSAVE_DELAY = 2000 // 2秒自动保存

const ChapterEditor = () => {
  const { projectId, chapterId } = useParams<{ projectId: string; chapterId: string }>()
  const navigate = useNavigate()

  const [project, setProject] = useState<Project | null>(null)
  const [chapter, setChapter] = useState<Chapter | null>(null)
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [saveStatus, setSaveStatus] = useState<'saved' | 'saving' | 'unsaved'>('saved')

  const [title, setTitle] = useState('')
  const [content, setContent] = useState('')
  const [wordCount, setWordCount] = useState(0)

  const [historyOpen, setHistoryOpen] = useState(false)
  const [aiPanelOpen, setAiPanelOpen] = useState(false)
  const [selectedText, setSelectedText] = useState('')

  // 自动保存定时器
  const autoSaveTimer = useRef<ReturnType<typeof setTimeout> | null>(null)
  const lastSavedContent = useRef<string>('')
  const lastSavedTitle = useRef<string>('')

  // 加载数据
  useEffect(() => {
    const fetchData = async () => {
      if (!projectId || !chapterId) return

      setLoading(true)
      try {
        const [projectData, chapterData] = await Promise.all([
          projectService.get(projectId),
          chapterService.get(chapterId),
        ])

        setProject(projectData)
        setChapter(chapterData)
        setTitle(chapterData.title)
        setContent(chapterData.content || '')
        setWordCount(chapterData.word_count)

        lastSavedContent.current = chapterData.content || ''
        lastSavedTitle.current = chapterData.title
      } catch (error) {
        message.error('加载章节失败')
        console.error(error)
      } finally {
        setLoading(false)
      }
    }

    fetchData()
  }, [projectId, chapterId])

  // 保存章节
  const saveChapter = useCallback(async (isAuto = false) => {
    if (!chapterId) return

    // 检查是否有变化
    if (content === lastSavedContent.current && title === lastSavedTitle.current) {
      return
    }

    setSaving(true)
    setSaveStatus('saving')

    try {
      await chapterService.update(chapterId, {
        title,
        content,
      })

      lastSavedContent.current = content
      lastSavedTitle.current = title
      setSaveStatus('saved')

      if (!isAuto) {
        message.success('保存成功')
      }
    } catch (error) {
      setSaveStatus('unsaved')
      if (!isAuto) {
        message.error('保存失败')
      }
      console.error(error)
    } finally {
      setSaving(false)
    }
  }, [chapterId, title, content])

  // 内容变化时触发自动保存
  const handleContentChange = useCallback((newContent: string) => {
    setContent(newContent)
    setSaveStatus('unsaved')

    // 清除之前的定时器
    if (autoSaveTimer.current) {
      clearTimeout(autoSaveTimer.current)
    }

    // 设置新的自动保存定时器
    autoSaveTimer.current = setTimeout(() => {
      saveChapter(true)
    }, AUTOSAVE_DELAY)
  }, [saveChapter])

  // 标题变化
  const handleTitleChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    setTitle(e.target.value)
    setSaveStatus('unsaved')

    // 清除之前的定时器
    if (autoSaveTimer.current) {
      clearTimeout(autoSaveTimer.current)
    }

    // 设置新的自动保存定时器
    autoSaveTimer.current = setTimeout(() => {
      saveChapter(true)
    }, AUTOSAVE_DELAY)
  }, [saveChapter])

  // 组件卸载时清理定时器并保存
  useEffect(() => {
    return () => {
      if (autoSaveTimer.current) {
        clearTimeout(autoSaveTimer.current)
      }
    }
  }, [])

  // 版本恢复后刷新
  const handleVersionRestore = useCallback(async () => {
    if (!chapterId) return

    try {
      const chapterData = await chapterService.get(chapterId)
      setChapter(chapterData)
      setTitle(chapterData.title)
      setContent(chapterData.content || '')
      setWordCount(chapterData.word_count)

      lastSavedContent.current = chapterData.content || ''
      lastSavedTitle.current = chapterData.title
      setSaveStatus('saved')
    } catch (error) {
      console.error(error)
    }
  }, [chapterId])

  // AI 生成内容后应用
  const handleAIContentGenerated = useCallback((generatedContent: string) => {
    setContent(generatedContent)
    setSaveStatus('unsaved')
    // 触发自动保存
    if (autoSaveTimer.current) {
      clearTimeout(autoSaveTimer.current)
    }
    autoSaveTimer.current = setTimeout(() => {
      saveChapter(true)
    }, AUTOSAVE_DELAY)
  }, [saveChapter])

  // AI 替换选中内容
  const handleReplaceSelection = useCallback((original: string, replacement: string) => {
    const newContent = content.replace(original, replacement)
    setContent(newContent)
    setSaveStatus('unsaved')
    setSelectedText('')
    // 触发自动保存
    if (autoSaveTimer.current) {
      clearTimeout(autoSaveTimer.current)
    }
    autoSaveTimer.current = setTimeout(() => {
      saveChapter(true)
    }, AUTOSAVE_DELAY)
  }, [content, saveChapter])

  // 保存状态图标
  const SaveStatusIcon = () => {
    switch (saveStatus) {
      case 'saving':
        return (
          <Tooltip title="保存中...">
            <SyncOutlined spin className="text-blue-500" />
          </Tooltip>
        )
      case 'saved':
        return (
          <Tooltip title="已保存">
            <CheckCircleOutlined className="text-green-500" />
          </Tooltip>
        )
      case 'unsaved':
        return (
          <Tooltip title="未保存">
            <span className="w-2 h-2 bg-orange-500 rounded-full inline-block" />
          </Tooltip>
        )
    }
  }

  const moreMenuItems = [
    {
      key: 'ai',
      icon: <RobotOutlined />,
      label: 'AI 助手',
      onClick: () => setAiPanelOpen(true),
    },
    {
      key: 'history',
      icon: <HistoryOutlined />,
      label: '版本历史',
      onClick: () => setHistoryOpen(true),
    },
  ]

  if (loading) {
    return (
      <div className="flex justify-center items-center h-96">
        <Spin size="large" />
      </div>
    )
  }

  if (!chapter || !project) {
    return (
      <Card>
        <div className="text-center py-8">
          <p className="text-gray-500">章节不存在</p>
          <Button type="link" onClick={() => navigate(-1)}>
            返回
          </Button>
        </div>
      </Card>
    )
  }

  return (
    <div className="h-full flex flex-col">
      {/* 顶部导航 */}
      <div className="mb-4">
        <Breadcrumb
          items={[
            { title: <a onClick={() => navigate('/dashboard')}>我的项目</a> },
            { title: <a onClick={() => navigate(`/projects/${projectId}`)}>{project.title}</a> },
            { title: `第 ${chapter.chapter_number} 章` },
          ]}
        />
      </div>

      {/* 工具栏 */}
      <Card className="mb-4" bodyStyle={{ padding: '12px 16px' }}>
        <div className="flex items-center justify-between">
          <Space>
            <Button
              icon={<ArrowLeftOutlined />}
              onClick={() => navigate(`/projects/${projectId}`)}
            >
              返回
            </Button>
            <Tag color="blue">第 {chapter.chapter_number} 章</Tag>
            <SaveStatusIcon />
          </Space>

          <Space>
            <span className="text-sm text-gray-500">{wordCount} 字</span>
            <Button
              type="primary"
              icon={<SaveOutlined />}
              onClick={() => saveChapter(false)}
              loading={saving}
            >
              保存
            </Button>
            <Button
              icon={<RobotOutlined />}
              onClick={() => setAiPanelOpen(true)}
            >
              AI
            </Button>
            <Button
              icon={<HistoryOutlined />}
              onClick={() => setHistoryOpen(true)}
            >
              历史
            </Button>
            <Dropdown menu={{ items: moreMenuItems }} trigger={['click']}>
              <Button icon={<MoreOutlined />} />
            </Dropdown>
          </Space>
        </div>
      </Card>

      {/* 编辑区域 */}
      <Card className="flex-1" bodyStyle={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
        <Input
          value={title}
          onChange={handleTitleChange}
          placeholder="章节标题"
          className="text-xl font-bold mb-4"
          variant="borderless"
          style={{ fontSize: '1.25rem', fontWeight: 'bold' }}
        />

        <div className="flex-1">
          <RichTextEditor
            content={content}
            onChange={handleContentChange}
            onWordCountChange={setWordCount}
            placeholder="开始写作..."
            className="h-full"
          />
        </div>
      </Card>

      {/* 版本历史面板 */}
      <VersionHistoryPanel
        chapterId={chapterId || ''}
        open={historyOpen}
        onClose={() => setHistoryOpen(false)}
        onRestore={handleVersionRestore}
      />

      {/* AI 助手面板 */}
      <AIChatPanel
        chapterId={chapterId || ''}
        open={aiPanelOpen}
        onClose={() => setAiPanelOpen(false)}
        onContentGenerated={handleAIContentGenerated}
        selectedText={selectedText}
        onReplaceSelection={handleReplaceSelection}
      />
    </div>
  )
}

export default ChapterEditor
