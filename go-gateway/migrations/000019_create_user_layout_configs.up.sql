-- 创建用户布局配置表
CREATE TABLE IF NOT EXISTS user_layout_configs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    layout_type VARCHAR(50) NOT NULL,
    panel_states JSONB DEFAULT '{}',
    panel_sizes JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, layout_type)
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_user_layout_configs_user_id ON user_layout_configs(user_id);
CREATE INDEX IF NOT EXISTS idx_user_layout_configs_layout_type ON user_layout_configs(layout_type);
