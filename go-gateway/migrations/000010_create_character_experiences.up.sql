-- 000010_create_character_experiences.up.sql
-- 创建角色经历表

CREATE TABLE character_experiences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    character_id UUID NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    title VARCHAR(200) NOT NULL,
    content TEXT,
    chapter_ref INTEGER,
    time_point VARCHAR(100),
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_character_experiences_character ON character_experiences(character_id);
CREATE INDEX idx_character_experiences_chapter ON character_experiences(chapter_ref);
