import { useState, useEffect } from 'react'
import { Layout, Spin, message } from 'antd'
import { useParams } from 'react-router-dom'
import { EditorToolbar } from '@/components/editor/EditorToolbar'
import { ChapterSidebar } from '@/components/editor/ChapterSidebar'
import { RichTextEditor } from '@/components/editor/RichTextEditor'
import { AIChatPanel } from '@/components/editor/AIChatPanel'
import type { FC } from 'react'

const { Sider, Content } = Layout

interface Chapter {
  id: string
  title: string
  content?: string
  sortOrder: number
  status?: string
  wordCount?: number
}

export const EditorLayout: FC = () => {
  const { projectId, chapterId } = useParams<{ projectId: string; chapterId: string }>()
  const [chapters, setChapters] = useState<Chapter[]>([])
  const [currentChapter, setCurrentChapter] = useState<Chapter | null>(null)
  const [loading, setLoading] = useState(false)
  const [saving, setSaving] = useState(false)
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false)
  const [aiDrawerVisible, setAiDrawerVisible] = useState(false)
  const [content, setContent] = useState('')

  // 加载章节列表
  useEffect(() => {
    const loadChapters = async () => {
      try {
        // TODO: 调用实际API
        // const response = await chapterService.getChapters(projectId!)
        // setChapters(response.data)
        setChapters([])
      } catch (error) {
        message.error('加载章节列表失败')
      }
    }
    if (projectId) {
      loadChapters()
    }
  }, [projectId])

  // 加载当前章节内容
  useEffect(() => {
    const loadChapter = async () => {
      if (!chapterId) return
      setLoading(true)
      try {
        // TODO: 调用实际API
        // const chapter = await chapterService.getChapter(chapterId)
        // setCurrentChapter(chapter)
        // setContent(chapter.content)
        setCurrentChapter(null)
        setContent('')
      } catch (error) {
        message.error('加载章节内容失败')
      } finally {
        setLoading(false)
      }
    }
    loadChapter()
  }, [chapterId])

  // 自动保存
  useEffect(() => {
    if (!currentChapter) return

    const timer = setTimeout(() => {
      handleSave()
    }, 3000)

    return () => clearTimeout(timer)
  }, [content])

  const handleSave = async () => {
    if (!currentChapter) return
    setSaving(true)
    try {
      // TODO: 调用实际API
      // await chapterService.updateChapter(currentChapter.id, { content })
      message.success('已保存')
    } catch (error) {
      message.error('保存失败')
    } finally {
      setSaving(false)
    }
  }

  const handleChapterSelect = (chapterId: string) => {
    // 导航到选中的章节
    window.location.href = `/projects/${projectId}/chapters/${chapterId}`
  }

  const handleChapterReorder = async (chapters: Chapter[]) => {
    try {
      // TODO: 调用实际API
      // await chapterService.reorderChapters(projectId!, chapters.map(c => c.id))
      setChapters(chapters)
      message.success('章节顺序已更新')
    } catch (error) {
      message.error('更新章节顺序失败')
    }
  }

  const handleUndo = () => {
    // TODO: 实现撤销功能
    message.info('撤销')
  }

  const handleRedo = () => {
    // TODO: 实现重做功能
    message.info('重做')
  }

  const handleOpenAI = () => {
    setAiDrawerVisible(true)
  }

  const handleOpenSettings = () => {
    // TODO: 打开编辑器设置
    message.info('打开设置')
  }

  return (
    <Layout className="h-screen">
      {/* 左侧章节导航栏 */}
      <Sider
        collapsible
        collapsed={sidebarCollapsed}
        onCollapse={setSidebarCollapsed}
        width={250}
        collapsedWidth={0}
        trigger={null}
        className="bg-white border-r"
      >
        <ChapterSidebar
          chapters={chapters}
          currentChapterId={chapterId}
          onSelect={handleChapterSelect}
          onReorder={handleChapterReorder}
        />
      </Sider>

      <Layout>
        {/* 顶部工具栏 */}
        <div className="bg-white border-b">
          <EditorToolbar
            onToggleSidebar={() => setSidebarCollapsed(!sidebarCollapsed)}
            sidebarCollapsed={sidebarCollapsed}
            onSave={handleSave}
            saving={saving}
            onUndo={handleUndo}
            onRedo={handleRedo}
            onOpenAI={handleOpenAI}
            onOpenSettings={handleOpenSettings}
          />
        </div>

        {/* 主编辑区 */}
        <Content className="overflow-auto">
          {loading ? (
            <div className="flex items-center justify-center h-full">
              <Spin size="large" tip="加载中..." />
            </div>
          ) : currentChapter ? (
            <div className="max-w-4xl mx-auto p-8">
              <RichTextEditor
                content={content}
                onChange={setContent}
                placeholder="开始写作..."
              />
            </div>
          ) : (
            <div className="flex items-center justify-center h-full text-gray-400">
              <div className="text-center">
                <p className="text-lg">请选择一个章节开始编辑</p>
                <p className="text-sm mt-2">或创建新章节</p>
              </div>
            </div>
          )}
        </Content>
      </Layout>

      {/* AI聊天抽屉 */}
      <AIChatPanel
        visible={aiDrawerVisible}
        onClose={() => setAiDrawerVisible(false)}
        onInsert={(text) => {
          setContent(content + text)
          message.success('已插入生成内容')
        }}
      />
    </Layout>
  )
}

export default EditorLayout
