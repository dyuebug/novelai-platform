-- 000009_create_character_relationships.up.sql
-- 创建角色关系表

CREATE TABLE character_relationships (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    character_id UUID NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    target_id UUID NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    relation_type VARCHAR(50) NOT NULL,
    description TEXT,
    start_chapter INTEGER,
    end_chapter INTEGER,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(character_id, target_id, relation_type)
);

CREATE INDEX idx_character_relationships_character ON character_relationships(character_id);
CREATE INDEX idx_character_relationships_target ON character_relationships(target_id);
