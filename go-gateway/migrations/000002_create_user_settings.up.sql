-- 000002_create_user_settings.up.sql
-- 创建用户设置表

CREATE TABLE user_settings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- AI 提供商设置
    default_provider VARCHAR(50) DEFAULT 'openai',
    default_model VARCHAR(100),
    openai_api_key TEXT,  -- 加密存储
    openai_base_url VARCHAR(500),
    anthropic_api_key TEXT,
    gemini_api_key TEXT,

    -- 生成参数
    default_temperature DECIMAL(3,2) DEFAULT 0.7,
    default_max_tokens INTEGER DEFAULT 4096,

    -- UI 偏好
    theme VARCHAR(20) DEFAULT 'light',
    language VARCHAR(10) DEFAULT 'zh-CN',
    editor_font_size INTEGER DEFAULT 16,

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_user_settings_user ON user_settings(user_id);

CREATE TRIGGER update_user_settings_updated_at
    BEFORE UPDATE ON user_settings
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
