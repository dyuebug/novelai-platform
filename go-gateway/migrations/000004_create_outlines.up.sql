-- 000004_create_outlines.up.sql
-- 创建大纲表

CREATE TABLE outlines (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,

    title VARCHAR(200),
    content TEXT NOT NULL,
    outline_type VARCHAR(50) DEFAULT 'main' CHECK (outline_type IN ('main', 'volume', 'chapter')),
    parent_id UUID REFERENCES outlines(id) ON DELETE CASCADE,
    sort_order INTEGER DEFAULT 0,

    -- 关联章节范围
    start_chapter INTEGER,
    end_chapter INTEGER,

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_outlines_project ON outlines(project_id);
CREATE INDEX idx_outlines_parent ON outlines(parent_id);
CREATE INDEX idx_outlines_sort ON outlines(project_id, sort_order);

CREATE TRIGGER update_outlines_updated_at
    BEFORE UPDATE ON outlines
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
