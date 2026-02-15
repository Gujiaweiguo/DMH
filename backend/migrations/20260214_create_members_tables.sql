-- Migration: Create members and member_profiles tables
-- Date: 2026-02-14
-- Purpose: Support member management integration tests

-- Members table (平台会员，以 unionid 为标识)
CREATE TABLE IF NOT EXISTS members (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  unionid VARCHAR(100) NOT NULL,
  nickname VARCHAR(100),
  avatar VARCHAR(500),
  phone VARCHAR(20),
  gender INT DEFAULT 0 COMMENT '0:未知 1:男 2:女',
  source VARCHAR(50) COMMENT '首次来源渠道',
  status VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT 'active/disabled',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at DATETIME,
  UNIQUE INDEX idx_unionid (unionid),
  INDEX idx_phone (phone),
  INDEX idx_source (source),
  INDEX idx_status (status),
  INDEX idx_created_at (created_at),
  INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='会员表';

-- Member profiles table (会员画像扩展)
CREATE TABLE IF NOT EXISTS member_profiles (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  member_id BIGINT NOT NULL,
  total_orders INT DEFAULT 0 COMMENT '累计订单数',
  total_payment DECIMAL(10,2) DEFAULT 0.00 COMMENT '累计支付金额',
  total_reward DECIMAL(10,2) DEFAULT 0.00 COMMENT '累计奖励金额',
  first_order_at DATETIME,
  last_order_at DATETIME,
  first_payment_at DATETIME,
  last_payment_at DATETIME,
  participated_campaigns INT DEFAULT 0 COMMENT '参与活动数',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE INDEX idx_member_id (member_id),
  FOREIGN KEY (member_id) REFERENCES members(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='会员画像扩展表';

-- Insert test member for integration tests
INSERT INTO members (unionid, nickname, phone, gender, source, status)
VALUES ('test_union_id_123', '测试会员', '13800138000', 0, 'wechat', 'active')
ON DUPLICATE KEY UPDATE nickname = '测试会员';

-- Insert corresponding member profile
INSERT INTO member_profiles (member_id, total_orders, total_payment, total_reward, participated_campaigns)
SELECT id, 0, 0.00, 0.00, 0 FROM members WHERE unionid = 'test_union_id_123'
ON DUPLICATE KEY UPDATE total_orders = 0;
