-- ============================================
-- 核销记录表和海报生成记录表
-- 创建时间: 2026-01-29
-- 说明: 支持核销记录和海报生成记录的查询功能
-- ============================================

USE dmh;

-- 核销记录表
CREATE TABLE IF NOT EXISTS verification_records (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '核销记录ID',
    order_id BIGINT NOT NULL COMMENT '关联订单ID',
    verification_status VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '核销状态: pending/verified/cancelled',
    verified_at DATETIME NULL COMMENT '核销时间',
    verified_by BIGINT DEFAULT NULL COMMENT '核销人ID',
    verification_code VARCHAR(50) DEFAULT '' COMMENT '核销码',
    verification_method VARCHAR(20) DEFAULT 'manual' COMMENT '核销方式: manual/auto/qrcode',
    remark VARCHAR(500) DEFAULT '' COMMENT '备注说明',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_order_id (order_id),
    INDEX idx_verification_status (verification_status),
    INDEX idx_verified_at (verified_at),
    INDEX idx_verified_by (verified_by),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='核销记录表';

-- 海报生成记录表
CREATE TABLE IF NOT EXISTS poster_records (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '海报记录ID',
    record_type VARCHAR(20) NOT NULL DEFAULT 'personal' COMMENT '记录类型: personal/brand',
    campaign_id BIGINT NOT NULL COMMENT '关联活动ID',
    distributor_id BIGINT DEFAULT 0 COMMENT '推广人ID（个人海报）',
    template_name VARCHAR(100) NOT NULL DEFAULT '' COMMENT '模板名称',
    poster_url VARCHAR(500) NOT NULL DEFAULT '' COMMENT '海报URL',
    thumbnail_url VARCHAR(500) DEFAULT '' COMMENT '缩略图URL',
    file_size VARCHAR(50) DEFAULT '' COMMENT '文件大小（如 2.5MB）',
    generation_time INT DEFAULT 0 COMMENT '生成耗时（毫秒）',
    download_count INT DEFAULT 0 COMMENT '下载次数',
    share_count INT DEFAULT 0 COMMENT '分享次数',
    generated_by BIGINT DEFAULT NULL COMMENT '生成人ID',
    status VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '状态: active/deleted',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_record_type (record_type),
    INDEX idx_campaign_id (campaign_id),
    INDEX idx_distributor_id (distributor_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='海报生成记录表';
