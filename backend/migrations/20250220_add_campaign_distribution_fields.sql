SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;
SET collation_connection = 'utf8mb4_unicode_ci';

-- ============================================
-- 为 campaigns 表添加分销相关字段
-- ============================================

-- 检查并添加 enable_distribution 字段
SET @col_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS 
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'campaigns' AND COLUMN_NAME = 'enable_distribution');
SET @sql = IF(@col_exists = 0, 
    'ALTER TABLE campaigns ADD COLUMN enable_distribution TINYINT(1) NOT NULL DEFAULT 0 COMMENT ''是否启用分销''',
    'SELECT 1');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 检查并添加 distribution_level 字段
SET @col_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS 
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'campaigns' AND COLUMN_NAME = 'distribution_level');
SET @sql = IF(@col_exists = 0, 
    'ALTER TABLE campaigns ADD COLUMN distribution_level INT NOT NULL DEFAULT 0 COMMENT ''分销层级''',
    'SELECT 1');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 检查并添加 distribution_rewards 字段
SET @col_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS 
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'campaigns' AND COLUMN_NAME = 'distribution_rewards');
SET @sql = IF(@col_exists = 0, 
    'ALTER TABLE campaigns ADD COLUMN distribution_rewards TEXT COMMENT ''分销奖励配置(JSON)''',
    'SELECT 1');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 检查并添加 payment_config 字段
SET @col_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS 
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'campaigns' AND COLUMN_NAME = 'payment_config');
SET @sql = IF(@col_exists = 0, 
    'ALTER TABLE campaigns ADD COLUMN payment_config JSON COMMENT ''支付配置(JSON)''',
    'SELECT 1');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 检查并添加 poster_template_id 字段
SET @col_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS 
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'campaigns' AND COLUMN_NAME = 'poster_template_id');
SET @sql = IF(@col_exists = 0, 
    'ALTER TABLE campaigns ADD COLUMN poster_template_id BIGINT DEFAULT NULL COMMENT ''海报模板ID''',
    'SELECT 1');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
