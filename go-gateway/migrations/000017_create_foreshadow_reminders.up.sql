-- 000017_create_foreshadow_reminders.up.sql
-- 创建伏笔提醒表

CREATE TABLE foreshadow_reminders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    foreshadow_id UUID NOT NULL REFERENCES foreshadows(id) ON DELETE CASCADE,
    chapter_num INTEGER NOT NULL,
    message TEXT,
    is_read BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_foreshadow_reminders_foreshadow ON foreshadow_reminders(foreshadow_id);
CREATE INDEX idx_foreshadow_reminders_chapter ON foreshadow_reminders(chapter_num);
CREATE INDEX idx_foreshadow_reminders_unread ON foreshadow_reminders(is_read) WHERE is_read = false;
