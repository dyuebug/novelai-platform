import { Button, Typography, Tag, Empty } from 'antd'
import { PlusOutlined, FileTextOutlined } from '@ant-design/icons'
import { DragDropContext, Droppable, Draggable, DropResult, DroppableProvided, DraggableProvided, DraggableStateSnapshot } from 'react-beautiful-dnd'
import type { FC } from 'react'

const { Text } = Typography

interface Chapter {
  id: string
  title: string
  content?: string
  sortOrder: number
  status?: string
  wordCount?: number
}

interface ChapterSidebarProps {
  chapters: Chapter[]
  currentChapterId?: string
  onSelect?: (chapterId: string) => void
  onReorder?: (chapters: Chapter[]) => void | Promise<void>
  onAdd?: () => void
}

const statusColors = {
  draft: 'default',
  in_progress: 'processing',
  completed: 'success',
} as const

const statusLabels = {
  draft: '草稿',
  in_progress: '进行中',
  completed: '已完成',
} as const

export const ChapterSidebar: FC<ChapterSidebarProps> = ({
  chapters,
  currentChapterId,
  onSelect,
  onReorder,
  onAdd,
}) => {
  const handleDragEnd = (result: DropResult) => {
    if (!result.destination) return

    const items = Array.from(chapters)
    const [reorderedItem] = items.splice(result.source.index, 1)
    items.splice(result.destination.index, 0, reorderedItem)

    // 更新sortOrder
    const updatedItems = items.map((item, index) => ({
      ...item,
      sortOrder: index,
    }))

    onReorder?.(updatedItems)
  }

  const formatWordCount = (count?: number) => {
    if (!count) return '0字'
    if (count >= 10000) {
      return `${(count / 10000).toFixed(1)}万字`
    }
    return `${count}字`
  }

  if (chapters.length === 0) {
    return (
      <div className="p-4">
        <div className="mb-4">
          <Button type="primary" icon={<PlusOutlined />} onClick={onAdd} block>
            新建章节
          </Button>
        </div>
        <Empty
          image={<FileTextOutlined className="text-4xl text-gray-300" />}
          description="暂无章节"
          imageStyle={{ height: 60 }}
        />
      </div>
    )
  }

  return (
    <div className="h-full flex flex-col">
      <div className="p-4 border-b">
        <Button type="primary" icon={<PlusOutlined />} onClick={onAdd} block>
          新建章节
        </Button>
      </div>

      <div className="flex-1 overflow-auto">
        <DragDropContext onDragEnd={handleDragEnd}>
          <Droppable droppableId="chapters">
            {(provided: DroppableProvided) => (
              <div {...provided.droppableProps} ref={provided.innerRef}>
                {chapters.map((chapter, index) => (
                  <Draggable key={chapter.id} draggableId={chapter.id} index={index}>
                    {(provided: DraggableProvided, snapshot: DraggableStateSnapshot) => (
                      <div
                        ref={provided.innerRef}
                        {...provided.draggableProps}
                        {...provided.dragHandleProps}
                        className={`
                          px-4 py-3 border-b cursor-pointer transition-colors
                          ${currentChapterId === chapter.id ? 'bg-blue-50 border-l-4 border-l-blue-500' : 'hover:bg-gray-50'}
                          ${snapshot.isDragging ? 'shadow-lg bg-white' : ''}
                        `}
                        onClick={() => onSelect?.(chapter.id)}
                      >
                        <div className="space-y-1">
                          <div className="flex items-start justify-between">
                            <Text
                              strong={currentChapterId === chapter.id}
                              className="flex-1 line-clamp-2"
                            >
                              {chapter.title || '未命名章节'}
                            </Text>
                          </div>

                          <div className="flex items-center justify-between text-xs text-gray-500">
                            <span>{formatWordCount(chapter.wordCount)}</span>
                            {chapter.status && (
                              <Tag
                                color={statusColors[chapter.status as keyof typeof statusColors]}
                                className="text-xs"
                              >
                                {statusLabels[chapter.status as keyof typeof statusLabels]}
                              </Tag>
                            )}
                          </div>
                        </div>
                      </div>
                    )}
                  </Draggable>
                ))}
                {provided.placeholder}
              </div>
            )}
          </Droppable>
        </DragDropContext>
      </div>
    </div>
  )
}
