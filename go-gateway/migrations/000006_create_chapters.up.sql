-- 000006_create_chapters.up.sql
-- 创建章节表

CREATE TABLE chapters (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    chapter_number INTEGER NOT NULL,
    title VARCHAR(200) NOT NULL,
    content TEXT,
    summary TEXT,
    status VARCHAR(20) DEFAULT 'draft',
    word_count INTEGER DEFAULT 0,
    pov_character VARCHAR(100),
    location VARCHAR(200),
    time_setting VARCHAR(200),
    hook_type VARCHAR(50),
    hook_strength VARCHAR(20),
    coolpoint_patterns JSONB DEFAULT '[]',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    published_at TIMESTAMP,
    UNIQUE(project_id, chapter_number)
);

CREATE INDEX idx_chapters_project ON chapters(project_id);
CREATE INDEX idx_chapters_status ON chapters(status);
CREATE INDEX idx_chapters_number ON chapters(project_id, chapter_number);
