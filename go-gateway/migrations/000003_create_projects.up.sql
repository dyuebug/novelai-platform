-- 000003_create_projects.up.sql
-- 创建项目表

CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    title VARCHAR(200) NOT NULL,
    description TEXT,
    genre VARCHAR(50),  -- '玄幻', '都市', '科幻', '历史', '言情'
    status VARCHAR(20) DEFAULT 'draft' CHECK (status IN ('draft', 'writing', 'completed', 'archived')),

    -- 统计信息
    total_chapters INTEGER DEFAULT 0,
    total_words INTEGER DEFAULT 0,

    -- 封面
    cover_url VARCHAR(500),

    -- 元数据
    metadata JSONB DEFAULT '{}',

    -- 软删除
    deleted_at TIMESTAMP,
    is_deleted BOOLEAN DEFAULT FALSE,

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_projects_user ON projects(user_id);
CREATE INDEX idx_projects_status ON projects(status);
CREATE INDEX idx_projects_deleted ON projects(is_deleted);
CREATE INDEX idx_projects_user_active ON projects(user_id) WHERE is_deleted = FALSE;

CREATE TRIGGER update_projects_updated_at
    BEFORE UPDATE ON projects
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
