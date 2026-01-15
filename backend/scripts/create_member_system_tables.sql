-- 会员系统相关表创建脚本

-- 1. 会员表
CREATE TABLE IF NOT EXISTS `members` (
  `id` BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `unionid` VARCHAR(100) NOT NULL UNIQUE COMMENT '微信 unionid，平台唯一',
  `nickname` VARCHAR(100) DEFAULT NULL COMMENT '昵称',
  `avatar` VARCHAR(500) DEFAULT NULL COMMENT '头像',
  `phone` VARCHAR(20) DEFAULT NULL COMMENT '手机号',
  `gender` INT DEFAULT 0 COMMENT '性别：0未知 1男 2女',
  `source` VARCHAR(50) DEFAULT NULL COMMENT '首次来源渠道',
  `status` VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '状态：active/disabled',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` TIMESTAMP NULL DEFAULT NULL,
  INDEX `idx_phone` (`phone`),
  INDEX `idx_source` (`source`),
  INDEX `idx_status` (`status`),
  INDEX `idx_created_at` (`created_at`),
  INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='会员表';

-- 2. 会员画像表
CREATE TABLE IF NOT EXISTS `member_profiles` (
  `id` BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `member_id` BIGINT NOT NULL UNIQUE COMMENT '会员ID',
  `total_orders` INT DEFAULT 0 COMMENT '累计订单数',
  `total_payment` DECIMAL(10,2) DEFAULT 0.00 COMMENT '累计支付金额',
  `total_reward` DECIMAL(10,2) DEFAULT 0.00 COMMENT '累计奖励金额',
  `first_order_at` TIMESTAMP NULL DEFAULT NULL COMMENT '首次下单时间',
  `last_order_at` TIMESTAMP NULL DEFAULT NULL COMMENT '最后下单时间',
  `first_payment_at` TIMESTAMP NULL DEFAULT NULL COMMENT '首次支付时间',
  `last_payment_at` TIMESTAMP NULL DEFAULT NULL COMMENT '最后支付时间',
  `participated_campaigns` INT DEFAULT 0 COMMENT '参与活动数',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (`member_id`) REFERENCES `members`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='会员画像表';

-- 3. 会员标签表
CREATE TABLE IF NOT EXISTS `member_tags` (
  `id` BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `name` VARCHAR(50) NOT NULL UNIQUE COMMENT '标签名称',
  `category` VARCHAR(50) DEFAULT NULL COMMENT '标签分类',
  `color` VARCHAR(20) DEFAULT NULL COMMENT '标签颜色',
  `description` VARCHAR(200) DEFAULT NULL COMMENT '标签描述',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX `idx_category` (`category`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='会员标签表';

-- 4. 会员标签关联表
CREATE TABLE IF NOT EXISTS `member_tag_links` (
  `id` BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `member_id` BIGINT NOT NULL COMMENT '会员ID',
  `tag_id` BIGINT NOT NULL COMMENT '标签ID',
  `created_by` BIGINT NOT NULL COMMENT '操作人ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY `uk_member_tag` (`member_id`, `tag_id`),
  INDEX `idx_member_id` (`member_id`),
  INDEX `idx_tag_id` (`tag_id`),
  FOREIGN KEY (`member_id`) REFERENCES `members`(`id`) ON DELETE CASCADE,
  FOREIGN KEY (`tag_id`) REFERENCES `member_tags`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='会员标签关联表';

-- 5. 会员品牌关联表
CREATE TABLE IF NOT EXISTS `member_brand_links` (
  `id` BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `member_id` BIGINT NOT NULL COMMENT '会员ID',
  `brand_id` BIGINT NOT NULL COMMENT '品牌ID',
  `first_campaign_id` BIGINT NOT NULL COMMENT '首次参与的活动ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY `uk_member_brand` (`member_id`, `brand_id`),
  INDEX `idx_member_id` (`member_id`),
  INDEX `idx_brand_id` (`brand_id`),
  FOREIGN KEY (`member_id`) REFERENCES `members`(`id`) ON DELETE CASCADE,
  FOREIGN KEY (`brand_id`) REFERENCES `brands`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='会员品牌关联表';

-- 6. 会员合并请求表
CREATE TABLE IF NOT EXISTS `member_merge_requests` (
  `id` BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `source_member_id` BIGINT NOT NULL COMMENT '被合并的会员ID',
  `target_member_id` BIGINT NOT NULL COMMENT '主会员ID（保留）',
  `status` VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '状态：pending/completed/failed',
  `reason` TEXT DEFAULT NULL COMMENT '合并原因',
  `conflict_info` JSON DEFAULT NULL COMMENT '冲突信息',
  `created_by` BIGINT NOT NULL COMMENT '操作人ID',
  `executed_at` TIMESTAMP NULL DEFAULT NULL COMMENT '执行时间',
  `error_msg` TEXT DEFAULT NULL COMMENT '错误信息',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX `idx_source_member_id` (`source_member_id`),
  INDEX `idx_target_member_id` (`target_member_id`),
  INDEX `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='会员合并请求表';

-- 7. 导出申请表
CREATE TABLE IF NOT EXISTS `export_requests` (
  `id` BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `brand_id` BIGINT NOT NULL COMMENT '申请导出的品牌ID',
  `requested_by` BIGINT NOT NULL COMMENT '申请人ID',
  `reason` TEXT NOT NULL COMMENT '导出原因',
  `filters` JSON DEFAULT NULL COMMENT '筛选条件',
  `status` VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '状态：pending/approved/rejected/completed',
  `approved_by` BIGINT DEFAULT NULL COMMENT '审批人ID',
  `approved_at` TIMESTAMP NULL DEFAULT NULL COMMENT '审批时间',
  `reject_reason` TEXT DEFAULT NULL COMMENT '驳回原因',
  `file_url` VARCHAR(500) DEFAULT NULL COMMENT '导出文件URL',
  `record_count` INT DEFAULT 0 COMMENT '导出记录数',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX `idx_brand_id` (`brand_id`),
  INDEX `idx_requested_by` (`requested_by`),
  INDEX `idx_status` (`status`),
  INDEX `idx_created_at` (`created_at`),
  FOREIGN KEY (`brand_id`) REFERENCES `brands`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='导出申请表';

-- 8. 为 orders 表添加会员关联字段
ALTER TABLE `orders` 
  ADD COLUMN `member_id` BIGINT DEFAULT NULL COMMENT '关联会员ID' AFTER `campaign_id`,
  ADD COLUMN `unionid` VARCHAR(100) DEFAULT NULL COMMENT '微信 unionid' AFTER `member_id`,
  ADD INDEX `idx_member_id` (`member_id`),
  ADD INDEX `idx_unionid` (`unionid`);

-- 9. 为 rewards 表添加会员关联字段
ALTER TABLE `rewards` 
  ADD COLUMN `member_id` BIGINT DEFAULT NULL COMMENT '关联会员ID' AFTER `user_id`,
  ADD INDEX `idx_member_id` (`member_id`);
