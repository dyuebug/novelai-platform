-- 000013_create_organization_members.up.sql
-- 创建组织成员表

CREATE TABLE organization_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    character_id UUID NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    position VARCHAR(100),
    rank VARCHAR(50),
    join_chapter INTEGER,
    leave_chapter INTEGER,
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(organization_id, character_id)
);

CREATE INDEX idx_organization_members_org ON organization_members(organization_id);
CREATE INDEX idx_organization_members_char ON organization_members(character_id);
CREATE INDEX idx_organization_members_status ON organization_members(status);
