import { useState, useCallback } from 'react'
import { Button } from 'antd'
import { PlusOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import ProjectList from '@/features/projects/ProjectList'
import CreateProjectModal from '@/features/projects/CreateProjectModal'
import { Project } from '@/services/project.service'

const Dashboard = () => {
  const navigate = useNavigate()
  const [createModalOpen, setCreateModalOpen] = useState(false)
  const [refreshKey, setRefreshKey] = useState(0)

  const handleCreateClick = () => {
    setCreateModalOpen(true)
  }

  const handleCreateSuccess = useCallback((project: Project) => {
    setCreateModalOpen(false)
    // 刷新列表
    setRefreshKey((prev) => prev + 1)
    // 跳转到项目详情
    navigate(`/projects/${project.id}`)
  }, [navigate])

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">我的项目</h1>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleCreateClick}>
          新建项目
        </Button>
      </div>

      <ProjectList key={refreshKey} onCreateClick={handleCreateClick} />

      <CreateProjectModal
        open={createModalOpen}
        onClose={() => setCreateModalOpen(false)}
        onSuccess={handleCreateSuccess}
      />
    </div>
  )
}

export default Dashboard
