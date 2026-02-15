SET @has_helpul := (
  SELECT COUNT(*)
  FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME = 'faq_items'
    AND COLUMN_NAME = 'helpul_count'
);

SET @has_helpful := (
  SELECT COUNT(*)
  FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME = 'faq_items'
    AND COLUMN_NAME = 'helpful_count'
);

SET @sql_rename := IF(
  @has_helpul > 0 AND @has_helpful = 0,
  'ALTER TABLE faq_items CHANGE COLUMN helpul_count helpful_count INT DEFAULT 0 COMMENT \'有帮助次数\'' ,
  'SELECT 1'
);

PREPARE stmt FROM @sql_rename;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @has_not_helpful := (
  SELECT COUNT(*)
  FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME = 'faq_items'
    AND COLUMN_NAME = 'not_helpful_count'
);

SET @sql_add_not_helpful := IF(
  @has_not_helpful = 0,
  'ALTER TABLE faq_items ADD COLUMN not_helpful_count INT DEFAULT 0 COMMENT \'无帮助次数\'' ,
  'SELECT 1'
);

PREPARE stmt2 FROM @sql_add_not_helpful;
EXECUTE stmt2;
DEALLOCATE PREPARE stmt2;
