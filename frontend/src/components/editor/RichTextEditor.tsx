import { useEditor, EditorContent } from '@tiptap/react'
import StarterKit from '@tiptap/starter-kit'
import Placeholder from '@tiptap/extension-placeholder'
import CharacterCount from '@tiptap/extension-character-count'
import { useEffect } from 'react'
import type { FC } from 'react'
import './RichTextEditor.css'

interface RichTextEditorProps {
  content: string
  onChange: (content: string) => void
  placeholder?: string
  editable?: boolean
}

export const RichTextEditor: FC<RichTextEditorProps> = ({
  content,
  onChange,
  placeholder = '开始写作...',
  editable = true,
}) => {
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
    ],
    content,
    editable,
    onUpdate: ({ editor }) => {
      onChange(editor.getHTML())
    },
    editorProps: {
      attributes: {
        class: 'prose prose-lg max-w-none focus:outline-none min-h-[500px] p-4',
      },
    },
  })

  useEffect(() => {
    if (editor && content !== editor.getHTML()) {
      editor.commands.setContent(content)
    }
  }, [content, editor])

  if (!editor) {
    return null
  }

  const wordCount = editor.storage.characterCount.words()
  const charCount = editor.storage.characterCount.characters()

  return (
    <div className="rich-text-editor">
      <EditorContent editor={editor} />

      <div className="fixed bottom-4 right-4 bg-white shadow-lg rounded-lg px-4 py-2 text-sm text-gray-600">
        <span>{wordCount} 字</span>
        <span className="mx-2">|</span>
        <span>{charCount} 字符</span>
      </div>
    </div>
  )
}

export default RichTextEditor
