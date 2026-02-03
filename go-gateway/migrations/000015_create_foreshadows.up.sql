-- 000015_create_foreshadows.up.sql
-- 创建伏笔表

CREATE TABLE foreshadows (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    content TEXT,
    priority VARCHAR(20) DEFAULT 'medium',
    status VARCHAR(20) DEFAULT 'planted',

    -- 埋设信息
    plant_chapter_id UUID REFERENCES chapters(id) ON DELETE SET NULL,
    plant_chapter_num INTEGER DEFAULT 0,
    plant_position VARCHAR(500),

    -- 回收信息
    resolve_chapter_id UUID REFERENCES chapters(id) ON DELETE SET NULL,
    resolve_chapter_num INTEGER DEFAULT 0,
    resolve_content TEXT,
    resolved_at TIMESTAMP,

    -- 提醒设置
    remind_chapter_num INTEGER DEFAULT 0,
    remind_enabled BOOLEAN DEFAULT true,

    -- 关联
    related_character_ids TEXT DEFAULT '[]',
    tags TEXT DEFAULT '[]',

    -- 元数据
    metadata TEXT DEFAULT '{}',
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_foreshadows_project ON foreshadows(project_id);
CREATE INDEX idx_foreshadows_status ON foreshadows(status);
CREATE INDEX idx_foreshadows_priority ON foreshadows(priority);
CREATE INDEX idx_foreshadows_plant_chapter ON foreshadows(plant_chapter_num);
CREATE INDEX idx_foreshadows_remind ON foreshadows(remind_chapter_num) WHERE remind_enabled = true;
