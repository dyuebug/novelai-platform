-- 000016_create_foreshadow_hints.up.sql
-- 创建伏笔暗示表

CREATE TABLE foreshadow_hints (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    foreshadow_id UUID NOT NULL REFERENCES foreshadows(id) ON DELETE CASCADE,
    chapter_id UUID NOT NULL REFERENCES chapters(id) ON DELETE CASCADE,
    chapter_num INTEGER DEFAULT 0,
    content TEXT,
    hint_type VARCHAR(50),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_foreshadow_hints_foreshadow ON foreshadow_hints(foreshadow_id);
CREATE INDEX idx_foreshadow_hints_chapter ON foreshadow_hints(chapter_id);
