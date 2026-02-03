-- 000014_create_world_settings.up.sql
-- 创建世界设定表

CREATE TABLE world_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    parent_id UUID REFERENCES world_settings(id) ON DELETE SET NULL,
    category VARCHAR(50) NOT NULL,
    title VARCHAR(200) NOT NULL,
    content TEXT,
    metadata JSONB DEFAULT '{}',
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_world_settings_project ON world_settings(project_id);
CREATE INDEX idx_world_settings_parent ON world_settings(parent_id);
CREATE INDEX idx_world_settings_category ON world_settings(category);
