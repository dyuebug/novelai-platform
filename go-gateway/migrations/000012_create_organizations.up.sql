-- 000012_create_organizations.up.sql
-- 创建组织表

CREATE TABLE organizations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(50),
    description TEXT,
    history TEXT,
    structure TEXT,
    goals TEXT,
    location_id UUID REFERENCES locations(id) ON DELETE SET NULL,
    logo_url VARCHAR(500),
    metadata JSONB DEFAULT '{}',
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_organizations_project ON organizations(project_id);
CREATE INDEX idx_organizations_type ON organizations(type);
CREATE INDEX idx_organizations_location ON organizations(location_id);
