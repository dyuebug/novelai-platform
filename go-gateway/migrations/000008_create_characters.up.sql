-- 000008_create_characters.up.sql
-- 创建角色表

CREATE TABLE characters (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    alias VARCHAR(200),
    gender VARCHAR(20),
    age VARCHAR(50),
    birthday VARCHAR(50),
    appearance TEXT,
    personality TEXT,
    background TEXT,
    abilities TEXT,
    goals TEXT,
    role VARCHAR(50),
    status VARCHAR(20) DEFAULT 'active',
    avatar_url VARCHAR(500),
    metadata JSONB DEFAULT '{}',
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_characters_project ON characters(project_id);
CREATE INDEX idx_characters_name ON characters(project_id, name);
CREATE INDEX idx_characters_role ON characters(role);
