-- 分销商系统数据库迁移脚本
-- 创建日期: 2024-01-19

-- 1. 创建 distributors 表（分销商信息表）
CREATE TABLE IF NOT EXISTS `distributors` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '分销商ID',
  `user_id` BIGINT NOT NULL COMMENT '关联用户ID',
  `brand_id` BIGINT NOT NULL COMMENT '关联品牌ID',
  `level` INT NOT NULL DEFAULT 1 COMMENT '分销级别(1/2/3)',
  `parent_id` BIGINT DEFAULT NULL COMMENT '上级分销商ID',
  `status` VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '状态(pending/active/suspended)',
  `approved_by` BIGINT DEFAULT NULL COMMENT '审批人ID',
  `approved_at` DATETIME DEFAULT NULL COMMENT '审批时间',
  `total_earnings` DECIMAL(10,2) NOT NULL DEFAULT 0.00 COMMENT '累计收益',
  `subordinates_count` INT NOT NULL DEFAULT 0 COMMENT '下级人数',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_brand` (`user_id`, `brand_id`),
  KEY `idx_distributor_user` (`user_id`),
  KEY `idx_distributor_brand` (`brand_id`),
  KEY `idx_parent` (`parent_id`),
  KEY `idx_status` (`status`),
  KEY `idx_updated_at` (`updated_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='分销商信息表';

-- 2. 创建 distributor_applications 表（分销商申请表）
CREATE TABLE IF NOT EXISTS `distributor_applications` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '申请ID',
  `user_id` BIGINT NOT NULL COMMENT '申请用户ID',
  `brand_id` BIGINT NOT NULL COMMENT '申请品牌ID',
  `status` VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '状态(pending/approved/rejected)',
  `reason` TEXT COMMENT '申请理由',
  `reviewed_by` BIGINT DEFAULT NULL COMMENT '审核人ID',
  `reviewed_at` DATETIME DEFAULT NULL COMMENT '审核时间',
  `review_notes` TEXT COMMENT '审核备注',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_user` (`user_id`),
  KEY `idx_brand` (`brand_id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='分销商申请表';

-- 3. 创建 distributor_level_rewards 表（分销商级别奖励配置表）
CREATE TABLE IF NOT EXISTS `distributor_level_rewards` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '配置ID',
  `brand_id` BIGINT NOT NULL COMMENT '品牌ID',
  `level` INT NOT NULL COMMENT '级别(1/2/3)',
  `reward_percentage` DECIMAL(5,2) NOT NULL DEFAULT 0.00 COMMENT '奖励百分比',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_brand_level` (`brand_id`, `level`),
  KEY `idx_level_reward_brand` (`brand_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='分销商级别奖励配置表';

-- 4. 创建 distributor_rewards 表（分销商奖励记录表）
CREATE TABLE IF NOT EXISTS `distributor_rewards` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '奖励ID',
  `distributor_id` BIGINT NOT NULL COMMENT '分销商ID',
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `order_id` BIGINT NOT NULL COMMENT '订单ID',
  `campaign_id` BIGINT NOT NULL COMMENT '活动ID',
  `amount` DECIMAL(10,2) NOT NULL DEFAULT 0.00 COMMENT '奖励金额',
  `level` INT NOT NULL COMMENT '奖励级别(1/2/3)',
  `reward_rate` DECIMAL(5,2) NOT NULL COMMENT '奖励比例',
  `from_user_id` BIGINT DEFAULT NULL COMMENT '购买用户ID',
  `status` VARCHAR(20) NOT NULL DEFAULT 'settled' COMMENT '状态',
  `settled_at` DATETIME DEFAULT NULL COMMENT '结算时间',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_distributor` (`distributor_id`),
  KEY `idx_user` (`user_id`),
  KEY `idx_order` (`order_id`),
  KEY `idx_campaign` (`campaign_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='分销商奖励记录表';

-- 5. 创建 distributor_links 表（分销商推广链接表）
CREATE TABLE IF NOT EXISTS `distributor_links` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '链接ID',
  `distributor_id` BIGINT NOT NULL COMMENT '分销商ID',
  `campaign_id` BIGINT NOT NULL COMMENT '活动ID',
  `link_code` VARCHAR(50) NOT NULL COMMENT '推广码',
  `click_count` INT NOT NULL DEFAULT 0 COMMENT '点击次数',
  `order_count` INT NOT NULL DEFAULT 0 COMMENT '订单数',
  `status` VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '状态',
  `expires_at` DATETIME DEFAULT NULL COMMENT '过期时间',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_link_code` (`link_code`),
  KEY `idx_distributor` (`distributor_id`),
  KEY `idx_campaign` (`campaign_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='分销商推广链接表';

-- 6. 插入 distributor 角色到 roles 表
INSERT INTO `roles` (`name`, `code`, `description`, `created_at`, `updated_at`)
VALUES ('分销商', 'distributor', '具备推广资格的高级顾客角色', NOW(), NOW())
ON DUPLICATE KEY UPDATE `description` = '具备推广资格的高级顾客角色';

-- 7. 初始化默认的级别奖励配置（可选，品牌可以自己配置）
-- 这里为示例品牌1创建默认配置
INSERT IGNORE INTO `distributor_level_rewards` (`brand_id`, `level`, `reward_percentage`)
VALUES
  (1, 1, 5.0),  -- 一级分销商 5%
  (1, 2, 2.0),  -- 二级分销商 2%
  (1, 3, 1.0);  -- 三级分销商 1%

-- 8. 为分销商角色创建基础菜单（如果menus表存在）
INSERT INTO `menus` (`name`, `code`, `path`, `icon`, `sort`, `type`, `platform`, `status`, `created_at`, `updated_at`)
VALUES
  ('分销中心', 'distributor_center', '/distributor', 'gift', 50, 'menu', 'h5', 'active', NOW(), NOW()),
  ('分销商管理', 'distributor_management', '/admin/distributors', 'team', 50, 'menu', 'admin', 'active', NOW(), NOW()),
  ('分销商审批', 'distributor_approval', '/admin/distributor-approvals', 'audit', 51, 'menu', 'admin', 'active', NOW(), NOW())
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`);
