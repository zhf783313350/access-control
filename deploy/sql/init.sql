-- PostgreSQL 初始化脚本
-- 创建数据库（如果不存在）
-- CREATE DATABASE access_control;

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    phone VARCHAR(20) NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uk_phone UNIQUE (phone)
);

COMMENT ON TABLE users IS '用户表';

-- 角色表
CREATE TABLE IF NOT EXISTS roles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    code VARCHAR(50) NOT NULL,
    description VARCHAR(255),
    status SMALLINT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    CONSTRAINT uk_role_code UNIQUE (code)
);

COMMENT ON TABLE roles IS '角色表';

-- 权限表
CREATE TABLE IF NOT EXISTS permissions (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    code VARCHAR(100) NOT NULL,
    description VARCHAR(255),
    type SMALLINT DEFAULT 1,
    parent_id BIGINT DEFAULT 0,
    path VARCHAR(255),
    method VARCHAR(10),
    sort INT DEFAULT 0,
    status SMALLINT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    CONSTRAINT uk_permission_code UNIQUE (code)
);

COMMENT ON TABLE permissions IS '权限表';
COMMENT ON COLUMN permissions.type IS '类型: 1-菜单, 2-按钮, 3-接口';

-- 用户角色关联表
CREATE TABLE IF NOT EXISTS user_roles (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    role_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uk_user_role UNIQUE (user_id, role_id)
);

CREATE INDEX IF NOT EXISTS idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_role_id ON user_roles(role_id);

COMMENT ON TABLE user_roles IS '用户角色关联表';

-- 角色权限关联表
CREATE TABLE IF NOT EXISTS role_permissions (
    id BIGSERIAL PRIMARY KEY,
    role_id BIGINT NOT NULL,
    permission_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uk_role_permission UNIQUE (role_id, permission_id)
);

CREATE INDEX IF NOT EXISTS idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX IF NOT EXISTS idx_role_permissions_permission_id ON role_permissions(permission_id);

COMMENT ON TABLE role_permissions IS '角色权限关联表';

-- 插入初始数据
-- 密码 123456 的 SHA256 值: 8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92
INSERT INTO users (phone, password) VALUES
('18888888888', '8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92')
ON CONFLICT (phone) DO NOTHING;

INSERT INTO roles (name, code, description) VALUES
('超级管理员', 'super_admin', '拥有所有权限'),
('管理员', 'admin', '系统管理员'),
('普通用户', 'user', '普通用户')
ON CONFLICT (code) DO NOTHING;

INSERT INTO permissions (name, code, description, type, path, method) VALUES
('用户管理', 'user:manage', '用户增删改查', 1, '/api/users', 'GET'),
('用户创建', 'user:create', '创建用户', 2, '/api/users', 'POST'),
('用户编辑', 'user:update', '编辑用户', 2, '/api/users/:id', 'PUT'),
('用户删除', 'user:delete', '删除用户', 2, '/api/users/:id', 'DELETE'),
('角色管理', 'role:manage', '角色增删改查', 1, '/api/roles', 'GET'),
('权限管理', 'permission:manage', '权限管理', 1, '/api/permissions', 'GET')
ON CONFLICT (code) DO NOTHING;

INSERT INTO user_roles (user_id, role_id) VALUES (1, 1)
ON CONFLICT (user_id, role_id) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id) VALUES
(1, 1), (1, 2), (1, 3), (1, 4), (1, 5), (1, 6),
(2, 1), (2, 5),
(3, 1)
ON CONFLICT (role_id, permission_id) DO NOTHING;
