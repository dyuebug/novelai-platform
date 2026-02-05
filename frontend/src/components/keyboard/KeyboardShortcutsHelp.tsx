import { Modal, Input, List, Tag, Typography } from 'antd'
import { useState } from 'react'
import type { FC } from 'react'

const { Search } = Input
const { Text } = Typography

interface Shortcut {
  category: string
  shortcuts: Array<{
    keys: string
    description: string
  }>
}

const shortcuts: Shortcut[] = [
  {
    category: '编辑',
    shortcuts: [
      { keys: 'Ctrl+B', description: '加粗' },
      { keys: 'Ctrl+I', description: '斜体' },
      { keys: 'Ctrl+U', description: '下划线' },
      { keys: 'Ctrl+S', description: '保存' },
      { keys: 'Ctrl+Z', description: '撤销' },
      { keys: 'Ctrl+Y', description: '重做' },
    ],
  },
  {
    category: '导航',
    shortcuts: [
      { keys: 'Ctrl+K', description: '快速搜索' },
      { keys: 'Ctrl+←', description: '上一章节' },
      { keys: 'Ctrl+→', description: '下一章节' },
      { keys: 'Ctrl+Shift+C', description: '章节列表' },
    ],
  },
  {
    category: '面板控制',
    shortcuts: [
      { keys: 'Ctrl+Shift+A', description: '打开/关闭AI抽屉' },
      { keys: 'Ctrl+B', description: '折叠/展开侧边栏' },
      { keys: 'Ctrl+Shift+F', description: '沉浸模式' },
      { keys: 'F11', description: '全屏' },
      { keys: 'Esc', description: '退出沉浸模式/全屏' },
    ],
  },
  {
    category: 'AI功能',
    shortcuts: [
      { keys: 'Ctrl+Shift+G', description: '快速生成' },
      { keys: 'Ctrl+Shift+R', description: '重写选中文本' },
      { keys: 'Ctrl+Shift+P', description: '润色选中文本' },
    ],
  },
  {
    category: '其他',
    shortcuts: [
      { keys: 'Ctrl+/', description: '显示快捷键帮助' },
      { keys: '?', description: '显示快捷键帮助' },
    ],
  },
]

interface KeyboardShortcutsHelpProps {
  visible: boolean
  onClose: () => void
}

export const KeyboardShortcutsHelp: FC<KeyboardShortcutsHelpProps> = ({ visible, onClose }) => {
  const [searchText, setSearchText] = useState('')

  const filteredShortcuts = shortcuts
    .map((category) => ({
      ...category,
      shortcuts: category.shortcuts.filter(
        (shortcut) =>
          shortcut.keys.toLowerCase().includes(searchText.toLowerCase()) ||
          shortcut.description.toLowerCase().includes(searchText.toLowerCase())
      ),
    }))
    .filter((category) => category.shortcuts.length > 0)

  return (
    <Modal
      title="快捷键帮助"
      open={visible}
      onCancel={onClose}
      footer={null}
      width={600}
    >
      <div className="space-y-4">
        <Search
          placeholder="搜索快捷键..."
          value={searchText}
          onChange={(e) => setSearchText(e.target.value)}
          allowClear
        />

        <div className="max-h-[500px] overflow-auto">
          {filteredShortcuts.map((category) => (
            <div key={category.category} className="mb-6">
              <Text strong className="text-base mb-2 block">
                {category.category}
              </Text>
              <List
                size="small"
                dataSource={category.shortcuts}
                renderItem={(item) => (
                  <List.Item className="flex justify-between items-center">
                    <span>{item.description}</span>
                    <Tag className="font-mono">{item.keys}</Tag>
                  </List.Item>
                )}
              />
            </div>
          ))}

          {filteredShortcuts.length === 0 && (
            <div className="text-center text-gray-400 py-8">
              未找到匹配的快捷键
            </div>
          )}
        </div>
      </div>
    </Modal>
  )
}
