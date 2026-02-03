import { useEditor, EditorContent } from '@tiptap/react'
import StarterKit from '@tiptap/starter-kit'
import Placeholder from '@tiptap/extension-placeholder'
import CharacterCount from '@tiptap/extension-character-count'
import Underline from '@tiptap/extension-underline'
import { useEffect, useCallback, useRef } from 'react'
import { Button, Tooltip, Divider } from 'antd'
import {
  BoldOutlined,
  ItalicOutlined,
  UnderlineOutlined,
  StrikethroughOutlined,
  OrderedListOutlined,
  UnorderedListOutlined,
  UndoOutlined,
  RedoOutlined,
} from '@ant-design/icons'

interface RichTextEditorProps {
  content: string
  onChange: (content: string) => void
  onWordCountChange?: (count: number) => void
  placeholder?: string
  editable?: boolean
  className?: string
}

const RichTextEditor = ({
  content,
  onChange,
  onWordCountChange,
  placeholder = '开始写作...',
  editable = true,
  className = '',
}: RichTextEditorProps) => {
  const isInternalUpdate = useRef(false)

  const editor = useEditor({
    extensions: [
      StarterKit.configure({
        heading: {
          levels: [1, 2, 3],
        },
      }),
      Placeholder.configure({
        placeholder,
      }),
      CharacterCount,
      Underline,
    ],
    content,
    editable,
    onUpdate: ({ editor }) => {
      if (!isInternalUpdate.current) {
        const html = editor.getHTML()
        onChange(html)

        // 计算字数（中文按字符计算）
        const text = editor.getText()
        const wordCount = text.length
        onWordCountChange?.(wordCount)
      }
    },
  })

  // 外部内容变化时更新编辑器
  useEffect(() => {
    if (editor && content !== editor.getHTML()) {
      isInternalUpdate.current = true
      editor.commands.setContent(content)
      isInternalUpdate.current = false
    }
  }, [content, editor])

  // 更新可编辑状态
  useEffect(() => {
    if (editor) {
      editor.setEditable(editable)
    }
  }, [editable, editor])

  const ToolbarButton = useCallback(
    ({
      icon,
      title,
      action,
      isActive = false,
      disabled = false,
    }: {
      icon: React.ReactNode
      title: string
      action: () => void
      isActive?: boolean
      disabled?: boolean
    }) => (
      <Tooltip title={title}>
        <Button
          type={isActive ? 'primary' : 'text'}
          icon={icon}
          onClick={action}
          disabled={disabled}
          size="small"
          className="mx-0.5"
        />
      </Tooltip>
    ),
    []
  )

  if (!editor) {
    return null
  }

  return (
    <div className={`border rounded-lg overflow-hidden ${className}`}>
      {/* 工具栏 */}
      <div className="bg-gray-50 border-b px-2 py-1 flex items-center flex-wrap">
        <ToolbarButton
          icon={<BoldOutlined />}
          title="粗体 (Ctrl+B)"
          action={() => editor.chain().focus().toggleBold().run()}
          isActive={editor.isActive('bold')}
        />
        <ToolbarButton
          icon={<ItalicOutlined />}
          title="斜体 (Ctrl+I)"
          action={() => editor.chain().focus().toggleItalic().run()}
          isActive={editor.isActive('italic')}
        />
        <ToolbarButton
          icon={<UnderlineOutlined />}
          title="下划线 (Ctrl+U)"
          action={() => editor.chain().focus().toggleUnderline().run()}
          isActive={editor.isActive('underline')}
        />
        <ToolbarButton
          icon={<StrikethroughOutlined />}
          title="删除线"
          action={() => editor.chain().focus().toggleStrike().run()}
          isActive={editor.isActive('strike')}
        />

        <Divider type="vertical" className="mx-1" />

        <ToolbarButton
          icon={<UnorderedListOutlined />}
          title="无序列表"
          action={() => editor.chain().focus().toggleBulletList().run()}
          isActive={editor.isActive('bulletList')}
        />
        <ToolbarButton
          icon={<OrderedListOutlined />}
          title="有序列表"
          action={() => editor.chain().focus().toggleOrderedList().run()}
          isActive={editor.isActive('orderedList')}
        />

        <Divider type="vertical" className="mx-1" />

        <ToolbarButton
          icon={<UndoOutlined />}
          title="撤销 (Ctrl+Z)"
          action={() => editor.chain().focus().undo().run()}
          disabled={!editor.can().undo()}
        />
        <ToolbarButton
          icon={<RedoOutlined />}
          title="重做 (Ctrl+Y)"
          action={() => editor.chain().focus().redo().run()}
          disabled={!editor.can().redo()}
        />

        <div className="flex-1" />

        <span className="text-xs text-gray-400 mr-2">
          {editor.storage.characterCount.characters()} 字
        </span>
      </div>

      {/* 编辑区域 */}
      <EditorContent
        editor={editor}
        className="prose prose-sm max-w-none p-4 min-h-[400px] focus:outline-none"
      />
    </div>
  )
}

export default RichTextEditor
