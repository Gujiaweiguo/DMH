-- 添加活动高级功能数据库迁移脚本
-- 日期: 2025-01-24
-- 变更: add-campaign-advanced-features

USE dmh;

-- ============================================
-- 1. 为 campaigns 表添加支付配置和海报模板字段
-- ============================================

-- 添加 payment_config 字段（支付配置）
ALTER TABLE campaigns 
ADD COLUMN payment_config JSON COMMENT '支付配置（订金、全款、商户号等）' AFTER reward_rule;

-- 添加 poster_template_id 字段（海报模板ID）
ALTER TABLE campaigns 
ADD COLUMN poster_template_id INT DEFAULT 1 COMMENT '海报模板ID' AFTER payment_config;

-- ============================================
-- 2. 为 orders 表添加核销相关字段
-- ============================================

-- 添加核销状态字段
ALTER TABLE orders
ADD COLUMN verification_status VARCHAR(20) DEFAULT 'unverified' COMMENT '核销状态: unverified/verified/cancelled' AFTER pay_status;

-- 添加核销时间字段
ALTER TABLE orders
ADD COLUMN verified_at DATETIME NULL COMMENT '核销时间' AFTER verification_status;

-- 添加核销人ID字段
ALTER TABLE orders
ADD COLUMN verified_by BIGINT NULL COMMENT '核销人用户ID' AFTER verified_at;

-- 添加核销码字段
ALTER TABLE orders
ADD COLUMN verification_code VARCHAR(50) NULL COMMENT '核销码（包含签名）' AFTER verified_by;

-- 添加索引
ALTER TABLE orders
ADD INDEX idx_verification_status (verification_status);

ALTER TABLE orders
ADD INDEX idx_verified_at (verified_at);

-- ============================================
-- 3. 创建海报模板配置表
-- ============================================
-- 注意：poster_templates 表已存在用于存储生成的海报记录
-- 这里创建 poster_template_configs 表用于存储模板配置

CREATE TABLE IF NOT EXISTS poster_template_configs (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '模板ID',
    name VARCHAR(100) NOT NULL COMMENT '模板名称',
    preview_image VARCHAR(255) COMMENT '预览图URL',
    config JSON NOT NULL COMMENT '模板配置（元素位置、样式等）',
    status VARCHAR(20) DEFAULT 'active' COMMENT '状态: active/inactive',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='海报模板配置表';

-- ============================================
-- 4. 插入默认海报模板数据
-- ============================================

INSERT INTO poster_template_configs (id, name, preview_image, config, status) VALUES
(1, '经典模板', '/templates/classic.jpg', '{
  "width": 750,
  "height": 1334,
  "background": "#FFFFFF",
  "elements": [
    {
      "type": "text",
      "content": "{{campaignName}}",
      "x": 50,
      "y": 100,
      "fontSize": 32,
      "color": "#333333",
      "fontWeight": "bold",
      "maxWidth": 650
    },
    {
      "type": "text",
      "content": "{{campaignDescription}}",
      "x": 50,
      "y": 160,
      "fontSize": 18,
      "color": "#666666",
      "maxWidth": 650
    },
    {
      "type": "qrcode",
      "content": "{{distributorLink}}",
      "x": 275,
      "y": 1000,
      "size": 200
    },
    {
      "type": "text",
      "content": "扫码参与活动",
      "x": 375,
      "y": 1220,
      "fontSize": 16,
      "color": "#999999",
      "align": "center"
    }
  ]
}', 'active'),

(2, '简约模板', '/templates/simple.jpg', '{
  "width": 750,
  "height": 1334,
  "background": "#F5F5F5",
  "elements": [
    {
      "type": "text",
      "content": "{{campaignName}}",
      "x": 375,
      "y": 200,
      "fontSize": 36,
      "color": "#000000",
      "fontWeight": "bold",
      "align": "center",
      "maxWidth": 700
    },
    {
      "type": "qrcode",
      "content": "{{distributorLink}}",
      "x": 275,
      "y": 500,
      "size": 200
    },
    {
      "type": "text",
      "content": "长按识别二维码",
      "x": 375,
      "y": 750,
      "fontSize": 18,
      "color": "#666666",
      "align": "center"
    }
  ]
}', 'active'),

(3, '时尚模板', '/templates/modern.jpg', '{
  "width": 750,
  "height": 1334,
  "background": "linear-gradient(135deg, #667eea 0%, #764ba2 100%)",
  "elements": [
    {
      "type": "text",
      "content": "{{campaignName}}",
      "x": 50,
      "y": 150,
      "fontSize": 40,
      "color": "#FFFFFF",
      "fontWeight": "bold",
      "maxWidth": 650
    },
    {
      "type": "text",
      "content": "{{campaignDescription}}",
      "x": 50,
      "y": 230,
      "fontSize": 20,
      "color": "#FFFFFF",
      "maxWidth": 650,
      "opacity": 0.9
    },
    {
      "type": "rect",
      "x": 225,
      "y": 900,
      "width": 300,
      "height": 300,
      "fill": "#FFFFFF",
      "radius": 20
    },
    {
      "type": "qrcode",
      "content": "{{distributorLink}}",
      "x": 275,
      "y": 950,
      "size": 200
    },
    {
      "type": "text",
      "content": "扫码立即参与",
      "x": 375,
      "y": 1230,
      "fontSize": 18,
      "color": "#FFFFFF",
      "align": "center"
    }
  ]
}', 'active');

-- ============================================
-- 5. 验证迁移结果
-- ============================================

-- 查看 campaigns 表结构
SHOW COLUMNS FROM campaigns LIKE 'payment_config';
SHOW COLUMNS FROM campaigns LIKE 'poster_template_id';

-- 查看 orders 表结构
SHOW COLUMNS FROM orders LIKE 'verification_status';
SHOW COLUMNS FROM orders LIKE 'verified_at';
SHOW COLUMNS FROM orders LIKE 'verified_by';
SHOW COLUMNS FROM orders LIKE 'verification_code';

-- 查看 poster_template_configs 表
SELECT COUNT(*) as template_count FROM poster_template_configs;

-- 显示迁移完成信息
SELECT '数据库迁移完成！' as status, 
       '已添加支付配置、海报模板、订单核销功能' as description,
       NOW() as migration_time;
