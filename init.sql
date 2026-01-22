-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    "phoneNumber" VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    status INTEGER DEFAULT 1,
    created_at VARCHAR(50),
    updated_at VARCHAR(50)
);

-- 插入默认管理员用户
INSERT INTO users ("phoneNumber", password, status, created_at, updated_at) 
VALUES ('admin', 'admin123', 1, '2026-01-22 00:00:00', '2026-01-22 00:00:00')
ON CONFLICT ("phoneNumber") DO NOTHING;

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_users_phone_number ON users("phoneNumber");
