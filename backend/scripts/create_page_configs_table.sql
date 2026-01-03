-- 创建页面配置表
CREATE TABLE IF NOT EXISTS page_configs (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    campaign_id BIGINT NOT NULL COMMENT '活动ID',
    components TEXT NOT NULL COMMENT '组件配置(JSON)',
    theme TEXT NOT NULL COMMENT '主题配置(JSON)',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL DEFAULT NULL COMMENT '删除时间',
    INDEX idx_campaign_id (campaign_id),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='页面配置表';
