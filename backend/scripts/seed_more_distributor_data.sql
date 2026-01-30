-- ============================================
-- 生成更多分销商测试数据
-- 用于分销监控页面展示
-- ============================================

SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

USE dmh;

-- ============================================
-- 1. 创建更多普通用户（作为订单客户）
-- ============================================
INSERT IGNORE INTO users (id, username, password, phone, email, real_name, role, status) VALUES
(20, 'customer01', '$2a$10$iL5hmpD0wGKSkRDCY92TL.y8wGarBWmnqVoFYlRxLM7xr0eSCzPEm', '13800000020', 'customer01@test.com', '客户01', 'participant', 'active'),
(21, 'customer02', '$2a$10$iL5hmpD0wGKSkRDCY92TL.y8wGarBWmnqVoFYlRxLM7xr0eSCzPEm', '13800000021', 'customer02@test.com', '客户02', 'participant', 'active'),
(22, 'customer03', '$2a$10$iL5hmpD0wGKSkRDCY92TL.y8wGarBWmnqVoFYlRxLM7xr0eSCzPEm', '13800000022', 'customer03@test.com', '客户03', 'participant', 'active'),
(23, 'customer04', '$2a$10$iL5hmpD0wGKSkRDCY92TL.y8wGarBWmnqVoFYlRxLM7xr0eSCzPEm', '13800000023', 'customer04@test.com', '客户04', 'participant', 'active'),
(24, 'customer05', '$2a$10$iL5hmpD0wGKSkRDCY92TL.y8wGarBWmnqVoFYlRxLM7xr0eSCzPEm', '13800000024', 'customer05@test.com', '客户05', 'participant', 'active'),
(25, 'customer06', '$2a$10$iL5hmpD0wGKSkRDCY92TL.y8wGarBWmnqVoFYlRxLM7xr0eSCzPEm', '13800000025', 'customer06@test.com', '客户06', 'participant', 'active'),
(26, 'customer07', '$2a$10$iL5hmpD0wGKSkRDCY92TL.y8wGarBWmnqVoFYlRxLM7xr0eSCzPEm', '13800000026', 'customer07@test.com', '客户07', 'participant', 'active'),
(27, 'customer08', '$2a$10$iL5hmpD0wGKSkRDCY92TL.y8wGarBWmnqVoFYlRxLM7xr0eSCzPEm', '13800000027', 'customer08@test.com', '客户08', 'participant', 'active'),
(28, 'customer09', '$2a$10$iL5hmpD0wGKSkRDCY92TL.y8wGarBWmnqVoFYlRxLM7xr0eSCzPEm', '13800000028', 'customer09@test.com', '客户09', 'participant', 'active'),
(29, 'customer10', '$2a$10$iL5hmpD0wGKSkRDCY92TL.y8wGarBWmnqVoFYlRxLM7xr0eSCzPEm', '13800000029', 'customer10@test.com', '客户10', 'participant', 'active'),
(30, 'customer11', '$2a$10$iL5hmpD0wGKSkRDCY92TL.y8wGarBWmnqVoFYlRxLM7xr0eSCzPEm', '13800000030', 'customer11@test.com', '客户11', 'participant', 'active'),
(31, 'customer12', '$2a$10$iL5hmpD0wGKSkRDCY92TL.y8wGarBWmnqVoFYlRxLM7xr0eSCzPEm', '13800000031', 'customer12@test.com', '客户12', 'participant', 'active'),
(32, 'customer13', '$2a$10$iL5hmpD0wGKSkRDCY92TL.y8wGarBWmnqVoFYlRxLM7xr0eSCzPEm', '13800000032', 'customer13@test.com', '客户13', 'participant', 'active'),
(33, 'customer14', '$2a$10$iL5hmpD0wGKSkRDCY92TL.y8wGarBWmnqVoFYlRxLM7xr0eSCzPEm', '13800000033', 'customer14@test.com', '客户14', 'participant', 'active'),
(34, 'customer15', '$2a$10$iL5hmpD0wGKSkRDCY92TL.y8wGarBWmnqVoFYlRxLM7xr0eSCzPEm', '13800000034', 'customer15@test.com', '客户15', 'participant', 'active'),
(35, 'customer16', '$2a$10$iL5hmpD0wGKSkRDCY92TL.y8wGarBWmnqVoFYlRxLM7xr0eSCzPEm', '13800000035', 'customer16@test.com', '客户16', 'participant', 'active'),
(36, 'customer17', '$2a$10$iL5hmpD0wGKSkRDCY92TL.y8wGarBWmnqVoFYlRxLM7xr0eSCzPEm', '13800000036', 'customer17@test.com', '客户17', 'participant', 'active'),
(37, 'customer18', '$2a$10$iL5hmpD0wGKSkRDCY92TL.y8wGarBWmnqVoFYlRxLM7xr0eSCzPEm', '13800000037', 'customer18@test.com', '客户18', 'participant', 'active'),
(38, 'customer19', '$2a$10$iL5hmpD0wGKSkRDCY92TL.y8wGarBWmnqVoFYlRxLM7xr0eSCzPEm', '13800000038', 'customer19@test.com', '客户19', 'participant', 'active'),
(39, 'customer20', '$2a$10$iL5hmpD0wGKSkRDCY92TL.y8wGarBWmnqVoFYlRxLM7xr0eSCzPEm', '13800000039', 'customer20@test.com', '客户20', 'participant', 'active');

-- ============================================
-- 2. 获取活动ID
-- ============================================
SET @campaign_id := (SELECT id FROM campaigns WHERE name = '新年促销活动' AND brand_id = 1 LIMIT 1);

-- 如果没有找到，使用第一个启用分销的活动
SET @campaign_id := IFNULL(@campaign_id, (SELECT id FROM campaigns WHERE enable_distribution = TRUE AND brand_id = 1 LIMIT 1));

-- ============================================
-- 3. 创建更多订单（通过不同分销商推荐）
-- ============================================

-- distributor1 (ID=4) 的订单 - 一级分销
INSERT IGNORE INTO orders (campaign_id, phone, form_data, referrer_id, distributor_path, status, pay_status, amount, paid_at, created_at) VALUES
(@campaign_id, '13800000020', '{"name":"客户01","phone":"13800000020"}', 4, '1', 'paid', 'paid', 150.00, DATE_SUB(NOW(), INTERVAL 1 DAY), DATE_SUB(NOW(), INTERVAL 1 DAY)),
(@campaign_id, '13800000021', '{"name":"客户02","phone":"13800000021"}', 4, '1', 'paid', 'paid', 200.00, DATE_SUB(NOW(), INTERVAL 2 DAY), DATE_SUB(NOW(), INTERVAL 2 DAY)),
(@campaign_id, '13800000022', '{"name":"客户03","phone":"13800000022"}', 4, '1', 'paid', 'paid', 180.00, DATE_SUB(NOW(), INTERVAL 3 DAY), DATE_SUB(NOW(), INTERVAL 3 DAY)),
(@campaign_id, '13800000023', '{"name":"客户04","phone":"13800000023"}', 4, '1', 'paid', 'paid', 220.00, DATE_SUB(NOW(), INTERVAL 4 DAY), DATE_SUB(NOW(), INTERVAL 4 DAY)),
(@campaign_id, '13800000024', '{"name":"客户05","phone":"13800000024"}', 4, '1', 'paid', 'paid', 160.00, DATE_SUB(NOW(), INTERVAL 5 DAY), DATE_SUB(NOW(), INTERVAL 5 DAY));

-- distributor2 (ID=5) 的订单 - 二级分销 (1,2)
INSERT IGNORE INTO orders (campaign_id, phone, form_data, referrer_id, distributor_path, status, pay_status, amount, paid_at, created_at) VALUES
(@campaign_id, '13800000025', '{"name":"客户06","phone":"13800000025"}', 5, '1,2', 'paid', 'paid', 190.00, DATE_SUB(NOW(), INTERVAL 1 DAY), DATE_SUB(NOW(), INTERVAL 1 DAY)),
(@campaign_id, '13800000026', '{"name":"客户07","phone":"13800000026"}', 5, '1,2', 'paid', 'paid', 210.00, DATE_SUB(NOW(), INTERVAL 2 DAY), DATE_SUB(NOW(), INTERVAL 2 DAY)),
(@campaign_id, '13800000027', '{"name":"客户08","phone":"13800000027"}', 5, '1,2', 'paid', 'paid', 170.00, DATE_SUB(NOW(), INTERVAL 3 DAY), DATE_SUB(NOW(), INTERVAL 3 DAY)),
(@campaign_id, '13800000028', '{"name":"客户09","phone":"13800000028"}', 5, '1,2', 'paid', 'paid', 230.00, DATE_SUB(NOW(), INTERVAL 4 DAY), DATE_SUB(NOW(), INTERVAL 4 DAY));

-- distributor3 (ID=6) 的订单 - 三级分销 (1,2,3)
INSERT IGNORE INTO orders (campaign_id, phone, form_data, referrer_id, distributor_path, status, pay_status, amount, paid_at, created_at) VALUES
(@campaign_id, '13800000029', '{"name":"客户10","phone":"13800000029"}', 6, '1,2,3', 'paid', 'paid', 240.00, DATE_SUB(NOW(), INTERVAL 1 DAY), DATE_SUB(NOW(), INTERVAL 1 DAY)),
(@campaign_id, '13800000030', '{"name":"客户11","phone":"13800000030"}', 6, '1,2,3', 'paid', 'paid', 260.00, DATE_SUB(NOW(), INTERVAL 2 DAY), DATE_SUB(NOW(), INTERVAL 2 DAY)),
(@campaign_id, '13800000031', '{"name":"客户12","phone":"13800000031"}', 6, '1,2,3', 'paid', 'paid', 280.00, DATE_SUB(NOW(), INTERVAL 3 DAY), DATE_SUB(NOW(), INTERVAL 3 DAY));

-- 一些待支付的订单
INSERT IGNORE INTO orders (campaign_id, phone, form_data, referrer_id, distributor_path, status, pay_status, amount, created_at) VALUES
(@campaign_id, '13800000032', '{"name":"客户13","phone":"13800000032"}', 4, '1', 'pending', 'unpaid', 150.00, NOW()),
(@campaign_id, '13800000033', '{"name":"客户14","phone":"13800000033"}', 5, '1,2', 'pending', 'unpaid', 180.00, NOW()),
(@campaign_id, '13800000034', '{"name":"客户15","phone":"13800000034"}', 6, '1,2,3', 'pending', 'unpaid', 200.00, NOW());

-- 一些已取消的订单
INSERT IGNORE INTO orders (campaign_id, phone, form_data, referrer_id, distributor_path, status, pay_status, amount, created_at) VALUES
(@campaign_id, '13800000035', '{"name":"客户16","phone":"13800000035"}', 4, '1', 'cancelled', 'unpaid', 120.00, DATE_SUB(NOW(), INTERVAL 7 DAY)),
(@campaign_id, '13800000036', '{"name":"客户17","phone":"13800000036"}', 5, '1,2', 'cancelled', 'unpaid', 140.00, DATE_SUB(NOW(), INTERVAL 8 DAY));

-- ============================================
-- 4. 为新订单创建分销奖励记录
-- ============================================

-- 获取新创建的订单ID
SET @order_id_1 := (SELECT id FROM orders WHERE phone = '13800000020' LIMIT 1);
SET @order_id_2 := (SELECT id FROM orders WHERE phone = '13800000021' LIMIT 1);
SET @order_id_3 := (SELECT id FROM orders WHERE phone = '13800000022' LIMIT 1);
SET @order_id_4 := (SELECT id FROM orders WHERE phone = '13800000023' LIMIT 1);
SET @order_id_5 := (SELECT id FROM orders WHERE phone = '13800000024' LIMIT 1);
SET @order_id_6 := (SELECT id FROM orders WHERE phone = '13800000025' LIMIT 1);
SET @order_id_7 := (SELECT id FROM orders WHERE phone = '13800000026' LIMIT 1);
SET @order_id_8 := (SELECT id FROM orders WHERE phone = '13800000027' LIMIT 1);
SET @order_id_9 := (SELECT id FROM orders WHERE phone = '13800000028' LIMIT 1);
SET @order_id_10 := (SELECT id FROM orders WHERE phone = '13800000029' LIMIT 1);
SET @order_id_11 := (SELECT id FROM orders WHERE phone = '13800000030' LIMIT 1);
SET @order_id_12 := (SELECT id FROM orders WHERE phone = '13800000031' LIMIT 1);

-- distributor1 的一级奖励 (10%)
INSERT IGNORE INTO distributor_rewards (distributor_id, user_id, order_id, campaign_id, amount, distributor_level, reward_percentage, status, settled_at) VALUES
(1, 4, @order_id_1, @campaign_id, 15.00, 1, 10.0, 'settled', DATE_SUB(NOW(), INTERVAL 1 DAY)),
(1, 4, @order_id_2, @campaign_id, 20.00, 1, 10.0, 'settled', DATE_SUB(NOW(), INTERVAL 2 DAY)),
(1, 4, @order_id_3, @campaign_id, 18.00, 1, 10.0, 'settled', DATE_SUB(NOW(), INTERVAL 3 DAY)),
(1, 4, @order_id_4, @campaign_id, 22.00, 1, 10.0, 'settled', DATE_SUB(NOW(), INTERVAL 4 DAY)),
(1, 4, @order_id_5, @campaign_id, 16.00, 1, 10.0, 'settled', DATE_SUB(NOW(), INTERVAL 5 DAY));

-- distributor2 的二级分销奖励
-- distributor1 获得一级奖励 (10%)
INSERT IGNORE INTO distributor_rewards (distributor_id, user_id, order_id, campaign_id, amount, distributor_level, reward_percentage, status, settled_at) VALUES
(1, 4, @order_id_6, @campaign_id, 19.00, 1, 10.0, 'settled', DATE_SUB(NOW(), INTERVAL 1 DAY)),
(1, 4, @order_id_7, @campaign_id, 21.00, 1, 10.0, 'settled', DATE_SUB(NOW(), INTERVAL 2 DAY)),
(1, 4, @order_id_8, @campaign_id, 17.00, 1, 10.0, 'settled', DATE_SUB(NOW(), INTERVAL 3 DAY)),
(1, 4, @order_id_9, @campaign_id, 23.00, 1, 10.0, 'settled', DATE_SUB(NOW(), INTERVAL 4 DAY));

-- distributor2 获得二级奖励 (5%)
INSERT IGNORE INTO distributor_rewards (distributor_id, user_id, order_id, campaign_id, amount, distributor_level, reward_percentage, status, settled_at) VALUES
(2, 5, @order_id_6, @campaign_id, 9.50, 2, 5.0, 'settled', DATE_SUB(NOW(), INTERVAL 1 DAY)),
(2, 5, @order_id_7, @campaign_id, 10.50, 2, 5.0, 'settled', DATE_SUB(NOW(), INTERVAL 2 DAY)),
(2, 5, @order_id_8, @campaign_id, 8.50, 2, 5.0, 'settled', DATE_SUB(NOW(), INTERVAL 3 DAY)),
(2, 5, @order_id_9, @campaign_id, 11.50, 2, 5.0, 'settled', DATE_SUB(NOW(), INTERVAL 4 DAY));

-- distributor3 的三级分销奖励
-- distributor1 获得一级奖励 (10%)
INSERT IGNORE INTO distributor_rewards (distributor_id, user_id, order_id, campaign_id, amount, distributor_level, reward_percentage, status, settled_at) VALUES
(1, 4, @order_id_10, @campaign_id, 24.00, 1, 10.0, 'settled', DATE_SUB(NOW(), INTERVAL 1 DAY)),
(1, 4, @order_id_11, @campaign_id, 26.00, 1, 10.0, 'settled', DATE_SUB(NOW(), INTERVAL 2 DAY)),
(1, 4, @order_id_12, @campaign_id, 28.00, 1, 10.0, 'settled', DATE_SUB(NOW(), INTERVAL 3 DAY));

-- distributor2 获得二级奖励 (5%)
INSERT IGNORE INTO distributor_rewards (distributor_id, user_id, order_id, campaign_id, amount, distributor_level, reward_percentage, status, settled_at) VALUES
(2, 5, @order_id_10, @campaign_id, 12.00, 2, 5.0, 'settled', DATE_SUB(NOW(), INTERVAL 1 DAY)),
(2, 5, @order_id_11, @campaign_id, 13.00, 2, 5.0, 'settled', DATE_SUB(NOW(), INTERVAL 2 DAY)),
(2, 5, @order_id_12, @campaign_id, 14.00, 2, 5.0, 'settled', DATE_SUB(NOW(), INTERVAL 3 DAY));

-- distributor3 获得三级奖励 (3%)
INSERT IGNORE INTO distributor_rewards (distributor_id, user_id, order_id, campaign_id, amount, distributor_level, reward_percentage, status, settled_at) VALUES
(3, 6, @order_id_10, @campaign_id, 7.20, 3, 3.0, 'settled', DATE_SUB(NOW(), INTERVAL 1 DAY)),
(3, 6, @order_id_11, @campaign_id, 7.80, 3, 3.0, 'settled', DATE_SUB(NOW(), INTERVAL 2 DAY)),
(3, 6, @order_id_12, @campaign_id, 8.40, 3, 3.0, 'settled', DATE_SUB(NOW(), INTERVAL 3 DAY));

-- ============================================
-- 5. 更新分销商累计收益
-- ============================================
UPDATE distributors d
SET total_earnings = (
    SELECT COALESCE(SUM(amount), 0)
    FROM distributor_rewards dr
    WHERE dr.distributor_id = d.id AND dr.status = 'settled'
)
WHERE d.id IN (1, 2, 3);

-- ============================================
-- 6. 更新用户余额
-- ============================================
INSERT INTO user_balances (user_id, balance, total_reward, version)
SELECT 
    u.id,
    COALESCE(SUM(dr.amount), 0),
    COALESCE(SUM(dr.amount), 0),
    0
FROM users u
LEFT JOIN distributor_rewards dr ON u.id = dr.user_id AND dr.status = 'settled'
WHERE u.id IN (4, 5, 6)
GROUP BY u.id
ON DUPLICATE KEY UPDATE
    balance = VALUES(balance),
    total_reward = VALUES(total_reward),
    version = version + 1;

-- ============================================
-- 7. 创建更多推广链接
-- ============================================
INSERT IGNORE INTO distributor_links (distributor_id, campaign_id, link_code, click_count, order_count, status, created_at) VALUES
(1, @campaign_id, 'LINK_D1_NEW1', 150, 15, 'active', DATE_SUB(NOW(), INTERVAL 10 DAY)),
(1, @campaign_id, 'LINK_D1_NEW2', 120, 12, 'active', DATE_SUB(NOW(), INTERVAL 8 DAY)),
(2, @campaign_id, 'LINK_D2_NEW1', 80, 8, 'active', DATE_SUB(NOW(), INTERVAL 7 DAY)),
(2, @campaign_id, 'LINK_D2_NEW2', 60, 6, 'active', DATE_SUB(NOW(), INTERVAL 5 DAY)),
(3, @campaign_id, 'LINK_D3_NEW1', 40, 4, 'active', DATE_SUB(NOW(), INTERVAL 4 DAY));

-- ============================================
-- 8. 显示统计结果
-- ============================================
SELECT '✓ 更多分销测试数据生成完成' AS result;

SELECT '========================================' AS '';
SELECT '数据统计' AS '';
SELECT '========================================' AS '';

SELECT 
    '分销商数量' AS item,
    COUNT(*) AS count
FROM distributors
WHERE brand_id = 1

UNION ALL

SELECT 
    '订单总数',
    COUNT(*)
FROM orders
WHERE campaign_id = @campaign_id

UNION ALL

SELECT 
    '已支付订单',
    COUNT(*)
FROM orders
WHERE campaign_id = @campaign_id AND pay_status = 'paid'

UNION ALL

SELECT 
    '奖励记录',
    COUNT(*)
FROM distributor_rewards
WHERE campaign_id = @campaign_id

UNION ALL

SELECT 
    '已结算奖励',
    COUNT(*)
FROM distributor_rewards
WHERE campaign_id = @campaign_id AND status = 'settled';

-- 显示分销商收益
SELECT '========================================' AS '';
SELECT '分销商收益统计' AS '';
SELECT '========================================' AS '';

SELECT 
    u.username AS '用户名',
    u.real_name AS '姓名',
    d.level AS '级别',
    d.total_earnings AS '累计收益',
    d.subordinates_count AS '下级数量',
    d.status AS '状态'
FROM distributors d
JOIN users u ON d.user_id = u.id
WHERE d.brand_id = 1
ORDER BY d.id;
