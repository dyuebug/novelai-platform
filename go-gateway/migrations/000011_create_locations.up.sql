-- 000011_create_locations.up.sql
-- 创建地点表

CREATE TABLE locations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    parent_id UUID REFERENCES locations(id) ON DELETE SET NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    type VARCHAR(50),
    climate VARCHAR(100),
    features TEXT,
    image_url VARCHAR(500),
    metadata JSONB DEFAULT '{}',
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_locations_project ON locations(project_id);
CREATE INDEX idx_locations_parent ON locations(parent_id);
CREATE INDEX idx_locations_type ON locations(type);
