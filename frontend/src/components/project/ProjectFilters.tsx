import { Input, Select, DatePicker, Button, Tag } from 'antd'
import { SearchOutlined, FilterOutlined, CloseCircleOutlined } from '@ant-design/icons'
import type { FC } from 'react'
import type { Dayjs } from 'dayjs'

const { RangePicker } = DatePicker

interface FilterValues {
  search?: string
  tags?: string[]
  status?: string[]
  dateRange?: [Dayjs | null, Dayjs | null] | null
}

interface ProjectFiltersProps {
  value?: FilterValues
  onChange?: (values: FilterValues) => void
  availableTags?: string[]
  onReset?: () => void
}

const statusOptions = [
  { label: '草稿', value: 'draft' },
  { label: '进行中', value: 'in_progress' },
  { label: '已完成', value: 'completed' },
]

export const ProjectFilters: FC<ProjectFiltersProps> = ({
  value = {},
  onChange,
  availableTags = [],
  onReset,
}) => {
  const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    onChange?.({ ...value, search: e.target.value })
  }

  const handleTagsChange = (tags: string[]) => {
    onChange?.({ ...value, tags })
  }

  const handleStatusChange = (status: string[]) => {
    onChange?.({ ...value, status })
  }

  const handleDateRangeChange = (dates: [Dayjs | null, Dayjs | null] | null) => {
    onChange?.({ ...value, dateRange: dates })
  }

  const handleReset = () => {
    onChange?.({})
    onReset?.()
  }

  const hasFilters =
    value.search ||
    (value.tags && value.tags.length > 0) ||
    (value.status && value.status.length > 0) ||
    value.dateRange

  return (
    <div className="space-y-4">
      <div className="flex flex-wrap gap-4">
        <Input
          placeholder="搜索项目标题..."
          prefix={<SearchOutlined />}
          value={value.search}
          onChange={handleSearchChange}
          allowClear
          className="flex-1 min-w-[200px]"
        />

        <Select
          mode="multiple"
          placeholder="选择标签"
          value={value.tags}
          onChange={handleTagsChange}
          options={availableTags.map(tag => ({ label: tag, value: tag }))}
          className="min-w-[200px]"
          maxTagCount="responsive"
          allowClear
        />

        <Select
          mode="multiple"
          placeholder="选择状态"
          value={value.status}
          onChange={handleStatusChange}
          options={statusOptions}
          className="min-w-[150px]"
          maxTagCount="responsive"
          allowClear
        />

        <RangePicker
          placeholder={['开始日期', '结束日期']}
          value={value.dateRange}
          onChange={handleDateRangeChange}
          className="min-w-[250px]"
        />

        {hasFilters && (
          <Button
            icon={<CloseCircleOutlined />}
            onClick={handleReset}
          >
            清除筛选
          </Button>
        )}
      </div>

      {hasFilters && (
        <div className="flex flex-wrap gap-2 items-center">
          <span className="text-sm text-gray-500">
            <FilterOutlined /> 当前筛选：
          </span>

          {value.search && (
            <Tag closable onClose={() => onChange?.({ ...value, search: undefined })}>
              搜索: {value.search}
            </Tag>
          )}

          {value.tags?.map(tag => (
            <Tag
              key={tag}
              closable
              onClose={() => {
                const newTags = value.tags?.filter(t => t !== tag)
                onChange?.({ ...value, tags: newTags })
              }}
            >
              标签: {tag}
            </Tag>
          ))}

          {value.status?.map(status => (
            <Tag
              key={status}
              closable
              onClose={() => {
                const newStatus = value.status?.filter(s => s !== status)
                onChange?.({ ...value, status: newStatus })
              }}
            >
              状态: {statusOptions.find(opt => opt.value === status)?.label}
            </Tag>
          ))}

          {value.dateRange && (
            <Tag closable onClose={() => onChange?.({ ...value, dateRange: null })}>
              日期: {value.dateRange[0]?.format('YYYY-MM-DD')} ~ {value.dateRange[1]?.format('YYYY-MM-DD')}
            </Tag>
          )}
        </div>
      )}
    </div>
  )
}
