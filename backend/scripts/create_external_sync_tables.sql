-- 外部数据库同步表结构
-- 用于 Sync Adapter 将订单和奖励数据同步到外部数据库

-- 1. 外部订单表
CREATE TABLE IF NOT EXISTS `external_orders` (
  `id` BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `order_id` BIGINT NOT NULL UNIQUE COMMENT '订单ID（来自主系统）',
  `campaign_id` BIGINT NOT NULL COMMENT '活动ID',
  `member_id` BIGINT DEFAULT NULL COMMENT '会员ID',
  `unionid` VARCHAR(100) DEFAULT NULL COMMENT '微信 unionid',
  `phone` VARCHAR(20) DEFAULT NULL COMMENT '手机号',
  `form_data` JSON DEFAULT NULL COMMENT '表单数据',
  `amount` DECIMAL(10,2) DEFAULT 0.00 COMMENT '订单金额',
  `pay_status` VARCHAR(20) DEFAULT NULL COMMENT '支付状态',
  `created_at` TIMESTAMP NOT NULL COMMENT '订单创建时间',
  `synced_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '同步时间',
  INDEX `idx_campaign_id` (`campaign_id`),
  INDEX `idx_member_id` (`member_id`),
  INDEX `idx_unionid` (`unionid`),
  INDEX `idx_phone` (`phone`),
  INDEX `idx_created_at` (`created_at`),
  INDEX `idx_synced_at` (`synced_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='外部订单同步表';

-- 2. 外部奖励表
CREATE TABLE IF NOT EXISTS `external_rewards` (
  `id` BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `reward_id` BIGINT NOT NULL UNIQUE COMMENT '奖励ID（来自主系统）',
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `member_id` BIGINT DEFAULT NULL COMMENT '会员ID',
  `order_id` BIGINT NOT NULL COMMENT '关联订单ID',
  `amount` DECIMAL(10,2) DEFAULT 0.00 COMMENT '奖励金额',
  `status` VARCHAR(20) DEFAULT NULL COMMENT '奖励状态',
  `settled_at` TIMESTAMP NULL DEFAULT NULL COMMENT '结算时间',
  `synced_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '同步时间',
  INDEX `idx_user_id` (`user_id`),
  INDEX `idx_member_id` (`member_id`),
  INDEX `idx_order_id` (`order_id`),
  INDEX `idx_status` (`status`),
  INDEX `idx_synced_at` (`synced_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='外部奖励同步表';

-- 说明：
-- 1. 这些表应该在外部数据库中创建（MySQL/Oracle/SQL Server）
-- 2. member_id 和 unionid 字段用于会员系统关联
-- 3. 同步适配器会自动处理数据的插入和更新（幂等性）
-- 4. 根据实际外部数据库类型，可能需要调整数据类型和语法
