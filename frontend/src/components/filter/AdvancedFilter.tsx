import { Drawer, Form, Select, DatePicker, InputNumber, Button, Space, Divider } from 'antd'
import { FilterOutlined } from '@ant-design/icons'
import type { FC } from 'react'
import type { Dayjs } from 'dayjs'

const { RangePicker } = DatePicker

interface FilterValues {
  tags?: string[]
  status?: string[]
  dateRange?: [Dayjs | null, Dayjs | null] | null
  wordCountMin?: number
  wordCountMax?: number
}

interface AdvancedFilterProps {
  visible: boolean
  onClose: () => void
  onApply: (values: FilterValues) => void
  availableTags?: string[]
  initialValues?: FilterValues
}

const statusOptions = [
  { label: '草稿', value: 'draft' },
  { label: '进行中', value: 'in_progress' },
  { label: '已完成', value: 'completed' },
]

export const AdvancedFilter: FC<AdvancedFilterProps> = ({
  visible,
  onClose,
  onApply,
  availableTags = [],
  initialValues = {},
}) => {
  const [form] = Form.useForm()

  const handleApply = async () => {
    try {
      const values = await form.validateFields()
      onApply(values)
      onClose()
    } catch (error) {
      // 验证失败
    }
  }

  const handleReset = () => {
    form.resetFields()
  }

  const handleValuesChange = () => {
    // 表单值变化处理
  }

  return (
    <Drawer
      title={
        <div className="flex items-center gap-2">
          <FilterOutlined />
          <span>高级筛选</span>
        </div>
      }
      placement="right"
      onClose={onClose}
      open={visible}
      width={400}
      footer={
        <div className="flex justify-between">
          <Button onClick={handleReset}>重置</Button>
          <Space>
            <Button onClick={onClose}>取消</Button>
            <Button type="primary" onClick={handleApply}>
              应用筛选
            </Button>
          </Space>
        </div>
      }
    >
      <Form
        form={form}
        layout="vertical"
        initialValues={initialValues}
        onValuesChange={handleValuesChange}
      >
        <Form.Item label="标签" name="tags">
          <Select
            mode="multiple"
            placeholder="选择标签"
            options={availableTags.map(tag => ({ label: tag, value: tag }))}
            maxTagCount="responsive"
            allowClear
          />
        </Form.Item>

        <Form.Item label="状态" name="status">
          <Select
            mode="multiple"
            placeholder="选择状态"
            options={statusOptions}
            maxTagCount="responsive"
            allowClear
          />
        </Form.Item>

        <Divider />

        <Form.Item label="日期范围" name="dateRange">
          <RangePicker
            placeholder={['开始日期', '结束日期']}
            className="w-full"
          />
        </Form.Item>

        <Divider />

        <div className="space-y-4">
          <div className="text-sm font-medium">字数范围</div>
          <Form.Item label="最小字数" name="wordCountMin">
            <InputNumber
              placeholder="最小字数"
              min={0}
              className="w-full"
              formatter={(value) => `${value}`.replace(/\B(?=(\d{3})+(?!\d))/g, ',')}
              parser={(value) => Number(value?.replace(/\$\s?|(,*)/g, '') || '0') as any}
            />
          </Form.Item>

          <Form.Item label="最大字数" name="wordCountMax">
            <InputNumber
              placeholder="最大字数"
              min={0}
              className="w-full"
              formatter={(value) => `${value}`.replace(/\B(?=(\d{3})+(?!\d))/g, ',')}
              parser={(value) => Number(value?.replace(/\$\s?|(,*)/g, '') || '0') as any}
            />
          </Form.Item>
        </div>

        <Divider />

        <div className="text-xs text-gray-500">
          <p>筛选条件将以 AND 逻辑组合</p>
          <p className="mt-1">即：同时满足所有选中的条件</p>
        </div>
      </Form>
    </Drawer>
  )
}
