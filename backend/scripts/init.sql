-- DMH 数字营销中台数据库初始化脚本

-- 创建数据库
CREATE DATABASE IF NOT EXISTS dmh DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE dmh;

-- 营销活动表
CREATE TABLE IF NOT EXISTS campaigns (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '活动ID',
    name VARCHAR(200) NOT NULL COMMENT '活动名称',
    description TEXT COMMENT '活动描述',
    form_fields JSON COMMENT '动态表单字段配置',
    reward_rule DECIMAL(10,2) NOT NULL DEFAULT 0.00 COMMENT '奖励规则（比例或固定金额）',
    start_time DATETIME NOT NULL COMMENT '活动开始时间',
    end_time DATETIME NOT NULL COMMENT '活动结束时间',
    status VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '活动状态: active/paused/ended',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at DATETIME NULL DEFAULT NULL COMMENT '删除时间（软删除）',
    INDEX idx_status (status),
    INDEX idx_start_time (start_time),
    INDEX idx_end_time (end_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='营销活动表';

-- 订单表
CREATE TABLE IF NOT EXISTS orders (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '订单ID',
    campaign_id BIGINT NOT NULL COMMENT '活动ID',
    phone VARCHAR(20) NOT NULL COMMENT '用户手机号',
    form_data JSON COMMENT '表单数据',
    referrer_id BIGINT DEFAULT 0 COMMENT '推荐人ID',
    status VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '订单状态: pending/paid/cancelled',
    amount DECIMAL(10,2) NOT NULL DEFAULT 0.00 COMMENT '订单金额',
    pay_status VARCHAR(20) NOT NULL DEFAULT 'unpaid' COMMENT '支付状态: unpaid/paid/refunded',
    trade_no VARCHAR(100) DEFAULT '' COMMENT '交易流水号',
    sync_status VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '同步状态: pending/synced/failed',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at DATETIME NULL DEFAULT NULL COMMENT '删除时间（软删除）',
    INDEX idx_campaign_id (campaign_id),
    INDEX idx_phone (phone),
    INDEX idx_referrer_id (referrer_id),
    INDEX idx_status (status),
    INDEX idx_pay_status (pay_status),
    INDEX idx_sync_status (sync_status),
    UNIQUE KEY uk_campaign_phone (campaign_id, phone, deleted_at) COMMENT '同一活动同一手机号只能报名一次'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订单表';

-- 奖励记录表
CREATE TABLE IF NOT EXISTS rewards (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '奖励ID',
    user_id BIGINT NOT NULL COMMENT '用户ID（推荐人）',
    order_id BIGINT NOT NULL COMMENT '关联订单ID',
    campaign_id BIGINT NOT NULL COMMENT '活动ID',
    amount DECIMAL(10,2) NOT NULL DEFAULT 0.00 COMMENT '奖励金额',
    status VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '奖励状态: pending/settled/cancelled',
    settled_at DATETIME NULL COMMENT '结算时间',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_user_id (user_id),
    INDEX idx_order_id (order_id),
    INDEX idx_campaign_id (campaign_id),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='奖励记录表';

-- 用户余额表
CREATE TABLE IF NOT EXISTS user_balances (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '余额ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    balance DECIMAL(10,2) NOT NULL DEFAULT 0.00 COMMENT '当前余额',
    total_reward DECIMAL(10,2) NOT NULL DEFAULT 0.00 COMMENT '累计奖励',
    version BIGINT NOT NULL DEFAULT 0 COMMENT '版本号（乐观锁）',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_user_id (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户余额表';

-- 同步日志表
CREATE TABLE IF NOT EXISTS sync_logs (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '日志ID',
    order_id BIGINT NOT NULL COMMENT '订单ID',
    sync_type VARCHAR(20) NOT NULL COMMENT '同步类型: order/reward',
    sync_status VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '同步状态: pending/synced/failed',
    attempts INT NOT NULL DEFAULT 0 COMMENT '尝试次数',
    error_msg TEXT COMMENT '错误信息',
    synced_at DATETIME NULL COMMENT '同步成功时间',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_order_id (order_id),
    INDEX idx_sync_status (sync_status),
    INDEX idx_sync_type (sync_type),
    INDEX idx_synced_at (synced_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='同步日志表';

-- 外部数据库表（模拟外部系统的表结构）
-- 注意：实际使用时，这些表在外部数据库中创建

-- 外部订单表
CREATE TABLE IF NOT EXISTS external_orders (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    order_id BIGINT NOT NULL COMMENT 'DMH订单ID',
    campaign_id BIGINT NOT NULL COMMENT '活动ID',
    phone VARCHAR(20) NOT NULL COMMENT '手机号',
    form_data TEXT COMMENT '表单数据JSON',
    amount DECIMAL(10,2) NOT NULL DEFAULT 0.00 COMMENT '订单金额',
    pay_status VARCHAR(20) NOT NULL COMMENT '支付状态',
    created_at DATETIME NOT NULL COMMENT '创建时间',
    synced_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '同步时间',
    UNIQUE KEY uk_order_id (order_id),
    INDEX idx_phone (phone),
    INDEX idx_pay_status (pay_status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='外部订单表';

-- 外部奖励表
CREATE TABLE IF NOT EXISTS external_rewards (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    reward_id BIGINT NOT NULL COMMENT 'DMH奖励ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    order_id BIGINT NOT NULL COMMENT '订单ID',
    amount DECIMAL(10,2) NOT NULL DEFAULT 0.00 COMMENT '奖励金额',
    status VARCHAR(20) NOT NULL COMMENT '奖励状态',
    settled_at DATETIME NULL COMMENT '结算时间',
    synced_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '同步时间',
    UNIQUE KEY uk_reward_id (reward_id),
    INDEX idx_user_id (user_id),
    INDEX idx_order_id (order_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='外部奖励表';

-- ============================================
-- 权限管理相关表
-- ============================================

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '用户ID',
    username VARCHAR(50) NOT NULL COMMENT '用户名',
    password VARCHAR(255) NOT NULL COMMENT '密码（bcrypt加密）',
    phone VARCHAR(20) NOT NULL COMMENT '手机号',
    email VARCHAR(100) DEFAULT '' COMMENT '邮箱',
    avatar VARCHAR(255) DEFAULT '' COMMENT '头像URL',
    real_name VARCHAR(50) DEFAULT '' COMMENT '真实姓名',
    role VARCHAR(50) NOT NULL DEFAULT 'participant' COMMENT '用户角色: platform_admin/participant',
    status VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '用户状态: active/disabled/locked',
    login_attempts INT DEFAULT 0 COMMENT '登录尝试次数',
    locked_until DATETIME NULL COMMENT '锁定到期时间',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_username (username),
    UNIQUE KEY uk_phone (phone),
    INDEX idx_role (role),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 用户品牌关联表（用户与品牌的关联关系）
CREATE TABLE IF NOT EXISTS user_brands (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '关联ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    brand_id BIGINT NOT NULL COMMENT '品牌ID',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    UNIQUE KEY uk_user_brand (user_id, brand_id),
    INDEX idx_user_id (user_id),
    INDEX idx_brand_id (brand_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户品牌关联表';

-- 角色表
CREATE TABLE IF NOT EXISTS roles (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '角色ID',
    name VARCHAR(50) NOT NULL COMMENT '角色名称',
    code VARCHAR(50) NOT NULL COMMENT '角色编码',
    description VARCHAR(200) DEFAULT '' COMMENT '角色描述',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_code (code),
    INDEX idx_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表';

-- 用户角色关联表
CREATE TABLE IF NOT EXISTS user_roles (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '关联ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    role_id BIGINT NOT NULL COMMENT '角色ID',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    UNIQUE KEY uk_user_role (user_id, role_id),
    INDEX idx_user_id (user_id),
    INDEX idx_role_id (role_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户角色关联表';

-- 权限表
CREATE TABLE IF NOT EXISTS permissions (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '权限ID',
    name VARCHAR(100) NOT NULL COMMENT '权限名称',
    code VARCHAR(100) NOT NULL COMMENT '权限编码',
    resource VARCHAR(100) NOT NULL COMMENT '资源类型',
    action VARCHAR(50) NOT NULL COMMENT '操作类型',
    description VARCHAR(200) DEFAULT '' COMMENT '权限描述',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    UNIQUE KEY uk_code (code),
    INDEX idx_resource (resource),
    INDEX idx_action (action)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='权限表';

-- 角色权限关联表
CREATE TABLE IF NOT EXISTS role_permissions (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '关联ID',
    role_id BIGINT NOT NULL COMMENT '角色ID',
    permission_id BIGINT NOT NULL COMMENT '权限ID',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    UNIQUE KEY uk_role_permission (role_id, permission_id),
    INDEX idx_role_id (role_id),
    INDEX idx_permission_id (permission_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色权限关联表';

-- 菜单表
CREATE TABLE IF NOT EXISTS menus (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '菜单ID',
    name VARCHAR(100) NOT NULL COMMENT '菜单名称',
    code VARCHAR(100) NOT NULL COMMENT '菜单编码',
    path VARCHAR(200) DEFAULT '' COMMENT '菜单路径',
    icon VARCHAR(100) DEFAULT '' COMMENT '菜单图标',
    parent_id BIGINT DEFAULT NULL COMMENT '父菜单ID',
    sort INT DEFAULT 0 COMMENT '排序',
    type VARCHAR(20) NOT NULL DEFAULT 'menu' COMMENT '类型: menu/button',
    platform VARCHAR(20) NOT NULL DEFAULT 'admin' COMMENT '平台: admin/h5',
    status VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '状态: active/disabled',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_code (code),
    INDEX idx_parent_id (parent_id),
    INDEX idx_platform (platform),
    INDEX idx_status (status),
    INDEX idx_sort (sort)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='菜单表';

-- 角色菜单关联表
CREATE TABLE IF NOT EXISTS role_menus (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '关联ID',
    role_id BIGINT NOT NULL COMMENT '角色ID',
    menu_id BIGINT NOT NULL COMMENT '菜单ID',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    UNIQUE KEY uk_role_menu (role_id, menu_id),
    INDEX idx_role_id (role_id),
    INDEX idx_menu_id (menu_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色菜单关联表';

-- 品牌表
CREATE TABLE IF NOT EXISTS brands (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '品牌ID',
    name VARCHAR(100) NOT NULL COMMENT '品牌名称',
    logo VARCHAR(255) DEFAULT '' COMMENT '品牌Logo',
    description TEXT COMMENT '品牌描述',
    status VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '品牌状态: active/disabled',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_name (name),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='品牌表';

-- 品牌素材表
CREATE TABLE IF NOT EXISTS brand_assets (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '素材ID',
    brand_id BIGINT NOT NULL COMMENT '品牌ID',
    name VARCHAR(200) NOT NULL COMMENT '素材名称',
    type VARCHAR(50) NOT NULL COMMENT '素材类型: image/video/document',
    category VARCHAR(100) DEFAULT '' COMMENT '素材分类',
    tags VARCHAR(500) DEFAULT '' COMMENT '素材标签',
    file_url VARCHAR(500) NOT NULL COMMENT '文件URL',
    file_size BIGINT DEFAULT 0 COMMENT '文件大小（字节）',
    description TEXT COMMENT '素材描述',
    status VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '状态: active/disabled',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_brand_id (brand_id),
    INDEX idx_type (type),
    INDEX idx_category (category),
    INDEX idx_status (status),
    FOREIGN KEY (brand_id) REFERENCES brands(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='品牌素材表';

-- ============================================
-- 安全增强相关表
-- ============================================

-- 密码策略配置表
CREATE TABLE IF NOT EXISTS password_policies (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '策略ID',
    min_length INT NOT NULL DEFAULT 8 COMMENT '最小长度',
    require_uppercase BOOLEAN NOT NULL DEFAULT TRUE COMMENT '需要大写字母',
    require_lowercase BOOLEAN NOT NULL DEFAULT TRUE COMMENT '需要小写字母',
    require_numbers BOOLEAN NOT NULL DEFAULT TRUE COMMENT '需要数字',
    require_special_chars BOOLEAN NOT NULL DEFAULT TRUE COMMENT '需要特殊字符',
    max_age INT NOT NULL DEFAULT 90 COMMENT '密码最大有效期（天）',
    history_count INT NOT NULL DEFAULT 5 COMMENT '历史密码记录数量',
    max_login_attempts INT NOT NULL DEFAULT 5 COMMENT '最大登录尝试次数',
    lockout_duration INT NOT NULL DEFAULT 30 COMMENT '锁定时长（分钟）',
    session_timeout INT NOT NULL DEFAULT 480 COMMENT '会话超时时间（分钟）',
    max_concurrent_sessions INT NOT NULL DEFAULT 3 COMMENT '最大并发会话数',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='密码策略配置表';

-- 密码历史记录表
CREATE TABLE IF NOT EXISTS password_histories (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '历史ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    password_hash VARCHAR(255) NOT NULL COMMENT '密码哈希',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_user_id (user_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='密码历史记录表';

-- 登录尝试记录表
CREATE TABLE IF NOT EXISTS login_attempts (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '记录ID',
    user_id BIGINT DEFAULT NULL COMMENT '用户ID',
    username VARCHAR(50) NOT NULL COMMENT '尝试登录的用户名',
    client_ip VARCHAR(45) NOT NULL COMMENT '客户端IP',
    user_agent VARCHAR(500) DEFAULT '' COMMENT '用户代理',
    success BOOLEAN NOT NULL COMMENT '是否成功',
    fail_reason VARCHAR(200) DEFAULT '' COMMENT '失败原因',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_user_id (user_id),
    INDEX idx_username (username),
    INDEX idx_client_ip (client_ip),
    INDEX idx_success (success),
    INDEX idx_created_at (created_at),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='登录尝试记录表';

-- 用户会话记录表
CREATE TABLE IF NOT EXISTS user_sessions (
    id VARCHAR(64) PRIMARY KEY COMMENT '会话ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    client_ip VARCHAR(45) NOT NULL COMMENT '客户端IP',
    user_agent VARCHAR(500) DEFAULT '' COMMENT '用户代理',
    login_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '登录时间',
    last_active_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '最后活跃时间',
    expires_at DATETIME NOT NULL COMMENT '过期时间',
    status VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '状态: active/expired/revoked',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_user_id (user_id),
    INDEX idx_last_active_at (last_active_at),
    INDEX idx_expires_at (expires_at),
    INDEX idx_status (status),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户会话记录表';

-- 操作审计日志表
CREATE TABLE IF NOT EXISTS audit_logs (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '日志ID',
    user_id BIGINT DEFAULT NULL COMMENT '操作用户ID',
    username VARCHAR(50) DEFAULT '' COMMENT '操作用户名',
    action VARCHAR(100) NOT NULL COMMENT '操作类型',
    resource VARCHAR(100) DEFAULT '' COMMENT '操作资源',
    resource_id VARCHAR(100) DEFAULT '' COMMENT '资源ID',
    details TEXT COMMENT '操作详情',
    client_ip VARCHAR(45) DEFAULT '' COMMENT '客户端IP',
    user_agent VARCHAR(500) DEFAULT '' COMMENT '用户代理',
    status VARCHAR(20) NOT NULL DEFAULT 'success' COMMENT '操作状态: success/failed',
    error_msg VARCHAR(500) DEFAULT '' COMMENT '错误信息',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_user_id (user_id),
    INDEX idx_username (username),
    INDEX idx_action (action),
    INDEX idx_resource (resource),
    INDEX idx_resource_id (resource_id),
    INDEX idx_client_ip (client_ip),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='操作审计日志表';

-- 安全事件记录表
CREATE TABLE IF NOT EXISTS security_events (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '事件ID',
    event_type VARCHAR(50) NOT NULL COMMENT '事件类型',
    severity VARCHAR(20) NOT NULL COMMENT '严重程度: low/medium/high/critical',
    user_id BIGINT DEFAULT NULL COMMENT '相关用户ID',
    username VARCHAR(50) DEFAULT '' COMMENT '相关用户名',
    client_ip VARCHAR(45) NOT NULL COMMENT '客户端IP',
    user_agent VARCHAR(500) DEFAULT '' COMMENT '用户代理',
    description TEXT NOT NULL COMMENT '事件描述',
    details JSON COMMENT '事件详情',
    handled BOOLEAN NOT NULL DEFAULT FALSE COMMENT '是否已处理',
    handled_by BIGINT DEFAULT NULL COMMENT '处理人ID',
    handled_at DATETIME DEFAULT NULL COMMENT '处理时间',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_event_type (event_type),
    INDEX idx_severity (severity),
    INDEX idx_user_id (user_id),
    INDEX idx_username (username),
    INDEX idx_client_ip (client_ip),
    INDEX idx_handled (handled),
    INDEX idx_created_at (created_at),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (handled_by) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='安全事件记录表';

-- 提现申请表
CREATE TABLE IF NOT EXISTS withdrawals (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '提现ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    amount DECIMAL(10,2) NOT NULL COMMENT '提现金额',
    bank_name VARCHAR(100) DEFAULT '' COMMENT '银行名称',
    bank_account VARCHAR(50) DEFAULT '' COMMENT '银行账号',
    account_name VARCHAR(50) DEFAULT '' COMMENT '账户名称',
    status VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '提现状态: pending/approved/rejected',
    remark TEXT COMMENT '备注/拒绝原因',
    approved_by BIGINT DEFAULT NULL COMMENT '审核人ID',
    approved_at DATETIME NULL COMMENT '审核时间',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_user_id (user_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='提现申请表';

-- 添加 brand_id 到 campaigns 表
ALTER TABLE campaigns ADD COLUMN brand_id BIGINT DEFAULT 0 COMMENT '品牌ID' AFTER id;
ALTER TABLE campaigns ADD INDEX idx_brand_id (brand_id);

-- ============================================
-- 初始化基础数据
-- ============================================

-- 创建基础角色
INSERT INTO roles (name, code, description) VALUES
('平台管理员', 'platform_admin', '拥有系统所有权限的超级管理员'),
('参与者', 'participant', '普通用户，可参与活动'),
('匿名用户', 'anonymous', '未登录的访客用户');

-- 创建基础权限
INSERT INTO permissions (name, code, resource, action, description) VALUES
-- 用户管理权限
('用户查看', 'user:read', 'user', 'read', '查看用户信息'),
('用户创建', 'user:create', 'user', 'create', '创建新用户'),
('用户更新', 'user:update', 'user', 'update', '更新用户信息'),
('用户删除', 'user:delete', 'user', 'delete', '删除用户'),
-- 品牌管理权限
('品牌查看', 'brand:read', 'brand', 'read', '查看品牌信息'),
('品牌创建', 'brand:create', 'brand', 'create', '创建新品牌'),
('品牌更新', 'brand:update', 'brand', 'update', '更新品牌信息'),
('品牌删除', 'brand:delete', 'brand', 'delete', '删除品牌'),
-- 活动管理权限
('活动查看', 'campaign:read', 'campaign', 'read', '查看活动信息'),
('活动创建', 'campaign:create', 'campaign', 'create', '创建新活动'),
('活动更新', 'campaign:update', 'campaign', 'update', '更新活动信息'),
('活动删除', 'campaign:delete', 'campaign', 'delete', '删除活动'),
-- 订单管理权限
('订单查看', 'order:read', 'order', 'read', '查看订单信息'),
('订单创建', 'order:create', 'order', 'create', '创建新订单'),
('订单更新', 'order:update', 'order', 'update', '更新订单信息'),
-- 奖励管理权限
('奖励查看', 'reward:read', 'reward', 'read', '查看奖励信息'),
('奖励发放', 'reward:grant', 'reward', 'grant', '发放奖励'),
-- 提现管理权限
('提现申请', 'withdrawal:apply', 'withdrawal', 'apply', '申请提现'),
('提现审核', 'withdrawal:approve', 'withdrawal', 'approve', '审核提现申请'),
-- 角色权限管理
('角色查看', 'role:read', 'role', 'read', '查看角色信息'),
('权限配置', 'role:config', 'role', 'config', '配置角色权限'),
-- 菜单管理权限
('菜单查看', 'menu:read', 'menu', 'read', '查看菜单信息'),
('菜单创建', 'menu:create', 'menu', 'create', '创建新菜单'),
('菜单更新', 'menu:update', 'menu', 'update', '更新菜单信息'),
('菜单删除', 'menu:delete', 'menu', 'delete', '删除菜单');

-- 创建基础菜单（管理后台）
INSERT INTO menus (name, code, path, icon, parent_id, sort, type, platform, status) VALUES
-- 一级菜单
('仪表板', 'dashboard', '/dashboard', 'dashboard', NULL, 1, 'menu', 'admin', 'active'),
('用户管理', 'user-management', '/users', 'user', NULL, 2, 'menu', 'admin', 'active'),
('品牌管理', 'brand-management', '/brands', 'shop', NULL, 3, 'menu', 'admin', 'active'),
('活动管理', 'campaign-management', '/campaigns', 'calendar', NULL, 4, 'menu', 'admin', 'active'),
('订单管理', 'order-management', '/orders', 'shopping-cart', NULL, 5, 'menu', 'admin', 'active'),
('奖励管理', 'reward-management', '/rewards', 'gift', NULL, 6, 'menu', 'admin', 'active'),
('提现管理', 'withdrawal-management', '/withdrawals', 'money-collect', NULL, 7, 'menu', 'admin', 'active'),
('系统管理', 'system-management', '/system', 'setting', NULL, 8, 'menu', 'admin', 'active');

-- 二级菜单
INSERT INTO menus (name, code, path, icon, parent_id, sort, type, platform, status) VALUES
-- 用户管理子菜单
('用户列表', 'user-list', '/users/list', '', 2, 1, 'menu', 'admin', 'active'),
('创建用户', 'user-create', '/users/create', '', 2, 2, 'menu', 'admin', 'active'),
-- 品牌管理子菜单
('品牌列表', 'brand-list', '/brands/list', '', 3, 1, 'menu', 'admin', 'active'),
('创建品牌', 'brand-create', '/brands/create', '', 3, 2, 'menu', 'admin', 'active'),
('品牌关系管理', 'brand-relation', '/brands/relations', '', 3, 3, 'menu', 'admin', 'active'),
-- 活动管理子菜单
('活动列表', 'campaign-list', '/campaigns/list', '', 4, 1, 'menu', 'admin', 'active'),
('创建活动', 'campaign-create', '/campaigns/create', '', 4, 2, 'menu', 'admin', 'active'),
('页面配置', 'campaign-config', '/campaigns/config', '', 4, 3, 'menu', 'admin', 'active'),
-- 订单管理子菜单
('订单列表', 'order-list', '/orders/list', '', 5, 1, 'menu', 'admin', 'active'),
('订单统计', 'order-stats', '/orders/stats', '', 5, 2, 'menu', 'admin', 'active'),
-- 奖励管理子菜单
('奖励列表', 'reward-list', '/rewards/list', '', 6, 1, 'menu', 'admin', 'active'),
('奖励统计', 'reward-stats', '/rewards/stats', '', 6, 2, 'menu', 'admin', 'active'),
-- 提现管理子菜单
('提现申请', 'withdrawal-apply', '/withdrawals/apply', '', 7, 1, 'menu', 'admin', 'active'),
('提现审核', 'withdrawal-approve', '/withdrawals/approve', '', 7, 2, 'menu', 'admin', 'active'),
('提现记录', 'withdrawal-list', '/withdrawals/list', '', 7, 3, 'menu', 'admin', 'active'),
-- 系统管理子菜单
('角色管理', 'role-management', '/system/roles', '', 8, 1, 'menu', 'admin', 'active'),
('菜单管理', 'menu-management', '/system/menus', '', 8, 2, 'menu', 'admin', 'active'),
('权限配置', 'permission-config', '/system/permissions', '', 8, 3, 'menu', 'admin', 'active');

-- 配置角色权限关联
-- 平台管理员拥有所有权限
INSERT INTO role_permissions (role_id, permission_id)
SELECT 1, id FROM permissions;

-- 参与者权限
INSERT INTO role_permissions (role_id, permission_id)
SELECT 2, id FROM permissions WHERE code IN (
    'campaign:read',
    'order:create',
    'reward:read',
    'withdrawal:apply'
);

-- 配置角色菜单关联
-- 平台管理员拥有所有菜单
INSERT INTO role_menus (role_id, menu_id)
SELECT 1, id FROM menus WHERE platform = 'admin';

-- 参与者菜单（H5端，这里先创建admin端的基础菜单）
INSERT INTO role_menus (role_id, menu_id)
SELECT 2, id FROM menus WHERE code IN (
    'dashboard',
    'campaign-management', 'campaign-list',
    'order-management', 'order-list',
    'reward-management', 'reward-list',
    'withdrawal-management', 'withdrawal-apply', 'withdrawal-list'
);

-- 创建测试用户
-- 密码都是 123456 的bcrypt加密: $2a$10$iL5hmpD0wGKSkRDCY92TL.y8wGarBWmnqVoFYlRxLM7xr0eSCzPEm
INSERT INTO users (username, password, phone, email, real_name, role, status) VALUES
('admin', '$2a$10$iL5hmpD0wGKSkRDCY92TL.y8wGarBWmnqVoFYlRxLM7xr0eSCzPEm', '13800000001', 'admin@dmh.com', '系统管理员', 'platform_admin', 'active'),
('brand_manager', '$2a$10$iL5hmpD0wGKSkRDCY92TL.y8wGarBWmnqVoFYlRxLM7xr0eSCzPEm', '13800000002', 'brand@dmh.com', '品牌经理', 'participant', 'active'),
('user001', '$2a$10$iL5hmpD0wGKSkRDCY92TL.y8wGarBWmnqVoFYlRxLM7xr0eSCzPEm', '13800000003', 'user001@dmh.com', '张三', 'participant', 'active');

-- 创建测试品牌
INSERT INTO brands (name, logo, description, status) VALUES
('品牌A', 'https://via.placeholder.com/150', '这是品牌A的描述', 'active'),
('品牌B', 'https://via.placeholder.com/150', '这是品牌B的描述', 'active');

-- 将 brand_manager 分配给品牌A
INSERT INTO user_brands (user_id, brand_id) VALUES
(2, 1);

-- 分配用户角色
INSERT INTO user_roles (user_id, role_id) VALUES
(1, 1), -- admin -> platform_admin
(2, 2), -- brand_manager -> participant
(3, 2); -- user001 -> participant

-- 插入测试活动数据
INSERT INTO campaigns (brand_id, name, description, form_fields, reward_rule, start_time, end_time, status) VALUES
(1, '新年促销活动', '新年大促，推荐有礼', '["姓名", "手机号", "地址"]', 10.00, '2025-01-01 00:00:00', '2025-12-31 23:59:59', 'active');

-- 插入用户余额
INSERT INTO user_balances (user_id, balance, total_reward) VALUES
(1, 0.00, 0.00),
(2, 0.00, 0.00),
(3, 0.00, 0.00);

-- ============================================
-- 初始化安全配置数据
-- ============================================

-- 插入默认密码策略
INSERT INTO password_policies (
    min_length, require_uppercase, require_lowercase, require_numbers, require_special_chars,
    max_age, history_count, max_login_attempts, lockout_duration, session_timeout, max_concurrent_sessions
) VALUES (
    8, TRUE, TRUE, TRUE, TRUE,
    90, 5, 5, 30, 480, 3
);
