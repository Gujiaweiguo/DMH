-- ============================================
-- DMH 完整测试数据初始化脚本
-- 包含：会员数据 + 分销商数据 + 订单数据
-- ============================================

SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

USE dmh;

-- ============================================
-- 第一部分：会员数据
-- ============================================

-- 1. 创建会员标签
INSERT IGNORE INTO member_tags (id, name, category, color, description) VALUES
(1, 'VIP会员', '等级', '#FFD700', '高价值会员'),
(2, '活跃用户', '行为', '#4CAF50', '经常参与活动'),
(3, '新用户', '等级', '#2196F3', '新注册用户'),
(4, '沉睡用户', '行为', '#9E9E9E', '长时间未活跃'),
(5, '高转化', '转化', '#FF5722', '转化率高的用户'),
(6, '推广达人', '行为', '#9C27B0', '推荐人数多');

-- 2. 创建会员数据（30个会员）
INSERT IGNORE INTO members (id, unionid, nickname, avatar, phone, gender, source, status, created_at) VALUES
(1, 'wx_union_0001', '张小明', 'https://api.dicebear.com/7.x/avataaars/svg?seed=1', '13900000101', 1, 'campaign_1', 'active', DATE_SUB(NOW(), INTERVAL 30 DAY)),
(2, 'wx_union_0002', '李小红', 'https://api.dicebear.com/7.x/avataaars/svg?seed=2', '13900000102', 2, 'campaign_1', 'active', DATE_SUB(NOW(), INTERVAL 28 DAY)),
(3, 'wx_union_0003', '王大力', 'https://api.dicebear.com/7.x/avataaars/svg?seed=3', '13900000103', 1, 'campaign_1', 'active', DATE_SUB(NOW(), INTERVAL 25 DAY)),
(4, 'wx_union_0004', '赵美丽', 'https://api.dicebear.com/7.x/avataaars/svg?seed=4', '13900000104', 2, 'campaign_2', 'active', DATE_SUB(NOW(), INTERVAL 22 DAY)),
(5, 'wx_union_0005', '刘强', 'https://api.dicebear.com/7.x/avataaars/svg?seed=5', '13900000105', 1, 'campaign_2', 'active', DATE_SUB(NOW(), INTERVAL 20 DAY)),
(6, 'wx_union_0006', '陈静', 'https://api.dicebear.com/7.x/avataaars/svg?seed=6', '13900000106', 2, 'campaign_1', 'active', DATE_SUB(NOW(), INTERVAL 18 DAY)),
(7, 'wx_union_0007', '杨帆', 'https://api.dicebear.com/7.x/avataaars/svg?seed=7', '13900000107', 1, 'campaign_3', 'active', DATE_SUB(NOW(), INTERVAL 15 DAY)),
(8, 'wx_union_0008', '周敏', 'https://api.dicebear.com/7.x/avataaars/svg?seed=8', '13900000108', 2, 'campaign_3', 'active', DATE_SUB(NOW(), INTERVAL 12 DAY)),
(9, 'wx_union_0009', '吴刚', 'https://api.dicebear.com/7.x/avataaars/svg?seed=9', '13900000109', 1, 'campaign_1', 'active', DATE_SUB(NOW(), INTERVAL 10 DAY)),
(10, 'wx_union_0010', '郑丽', 'https://api.dicebear.com/7.x/avataaars/svg?seed=10', '13900000110', 2, 'campaign_2', 'active', DATE_SUB(NOW(), INTERVAL 8 DAY)),
(11, 'wx_union_0011', '孙浩', 'https://api.dicebear.com/7.x/avataaars/svg?seed=11', '13900000111', 1, 'campaign_1', 'active', DATE_SUB(NOW(), INTERVAL 7 DAY)),
(12, 'wx_union_0012', '朱婷', 'https://api.dicebear.com/7.x/avataaars/svg?seed=12', '13900000112', 2, 'campaign_2', 'active', DATE_SUB(NOW(), INTERVAL 6 DAY)),
(13, 'wx_union_0013', '胡军', 'https://api.dicebear.com/7.x/avataaars/svg?seed=13', '13900000113', 1, 'campaign_3', 'active', DATE_SUB(NOW(), INTERVAL 5 DAY)),
(14, 'wx_union_0014', '高雪', 'https://api.dicebear.com/7.x/avataaars/svg?seed=14', '13900000114', 2, 'campaign_1', 'active', DATE_SUB(NOW(), INTERVAL 4 DAY)),
(15, 'wx_union_0015', '林峰', 'https://api.dicebear.com/7.x/avataaars/svg?seed=15', '13900000115', 1, 'campaign_2', 'active', DATE_SUB(NOW(), INTERVAL 3 DAY)),
(16, 'wx_union_0016', '何琳', 'https://api.dicebear.com/7.x/avataaars/svg?seed=16', '13900000116', 2, 'campaign_3', 'active', DATE_SUB(NOW(), INTERVAL 2 DAY)),
(17, 'wx_union_0017', '罗伟', 'https://api.dicebear.com/7.x/avataaars/svg?seed=17', '13900000117', 1, 'campaign_1', 'active', DATE_SUB(NOW(), INTERVAL 1 DAY)),
(18, 'wx_union_0018', '宋娟', 'https://api.dicebear.com/7.x/avataaars/svg?seed=18', '13900000118', 2, 'campaign_2', 'active', NOW()),
(19, 'wx_union_0019', '梁超', 'https://api.dicebear.com/7.x/avataaars/svg?seed=19', '13900000119', 1, 'campaign_1', 'active', NOW()),
(20, 'wx_union_0020', '韩梅', 'https://api.dicebear.com/7.x/avataaars/svg?seed=20', '13900000120', 2, 'campaign_3', 'active', NOW()),
(21, 'wx_union_0021', '冯涛', 'https://api.dicebear.com/7.x/avataaars/svg?seed=21', '13900000121', 1, 'campaign_1', 'disabled', DATE_SUB(NOW(), INTERVAL 60 DAY)),
(22, 'wx_union_0022', '曹艳', 'https://api.dicebear.com/7.x/avataaars/svg?seed=22', '13900000122', 2, 'campaign_2', 'active', DATE_SUB(NOW(), INTERVAL 14 DAY)),
(23, 'wx_union_0023', '彭飞', 'https://api.dicebear.com/7.x/avataaars/svg?seed=23', '13900000123', 1, 'campaign_3', 'active', DATE_SUB(NOW(), INTERVAL 11 DAY)),
(24, 'wx_union_0024', '董娜', 'https://api.dicebear.com/7.x/avataaars/svg?seed=24', '13900000124', 2, 'campaign_1', 'active', DATE_SUB(NOW(), INTERVAL 9 DAY)),
(25, 'wx_union_0025', '袁斌', 'https://api.dicebear.com/7.x/avataaars/svg?seed=25', '13900000125', 1, 'campaign_2', 'active', DATE_SUB(NOW(), INTERVAL 13 DAY)),
(26, 'wx_union_0026', '蒋莉', 'https://api.dicebear.com/7.x/avataaars/svg?seed=26', '13900000126', 2, 'campaign_3', 'active', DATE_SUB(NOW(), INTERVAL 16 DAY)),
(27, 'wx_union_0027', '薛鹏', 'https://api.dicebear.com/7.x/avataaars/svg?seed=27', '13900000127', 1, 'campaign_1', 'active', DATE_SUB(NOW(), INTERVAL 19 DAY)),
(28, 'wx_union_0028', '谭芳', 'https://api.dicebear.com/7.x/avataaars/svg?seed=28', '13900000128', 2, 'campaign_2', 'active', DATE_SUB(NOW(), INTERVAL 21 DAY)),
(29, 'wx_union_0029', '邹杰', 'https://api.dicebear.com/7.x/avataaars/svg?seed=29', '13900000129', 1, 'campaign_3', 'active', DATE_SUB(NOW(), INTERVAL 24 DAY)),
(30, 'wx_union_0030', '熊慧', 'https://api.dicebear.com/7.x/avataaars/svg?seed=30', '13900000130', 2, 'campaign_1', 'active', DATE_SUB(NOW(), INTERVAL 27 DAY));

-- 3. 创建会员画像数据
INSERT IGNORE INTO member_profiles (member_id, total_orders, total_payment, total_reward, first_order_at, last_order_at, first_payment_at, last_payment_at, participated_campaigns) VALUES
(1, 5, 1250.00, 125.00, DATE_SUB(NOW(), INTERVAL 29 DAY), DATE_SUB(NOW(), INTERVAL 1 DAY), DATE_SUB(NOW(), INTERVAL 29 DAY), DATE_SUB(NOW(), INTERVAL 1 DAY), 3),
(2, 3, 750.00, 75.00, DATE_SUB(NOW(), INTERVAL 27 DAY), DATE_SUB(NOW(), INTERVAL 2 DAY), DATE_SUB(NOW(), INTERVAL 27 DAY), DATE_SUB(NOW(), INTERVAL 2 DAY), 2),
(3, 4, 1000.00, 100.00, DATE_SUB(NOW(), INTERVAL 24 DAY), DATE_SUB(NOW(), INTERVAL 3 DAY), DATE_SUB(NOW(), INTERVAL 24 DAY), DATE_SUB(NOW(), INTERVAL 3 DAY), 2),
(4, 2, 500.00, 50.00, DATE_SUB(NOW(), INTERVAL 21 DAY), DATE_SUB(NOW(), INTERVAL 4 DAY), DATE_SUB(NOW(), INTERVAL 21 DAY), DATE_SUB(NOW(), INTERVAL 4 DAY), 1),
(5, 6, 1500.00, 150.00, DATE_SUB(NOW(), INTERVAL 19 DAY), DATE_SUB(NOW(), INTERVAL 1 DAY), DATE_SUB(NOW(), INTERVAL 19 DAY), DATE_SUB(NOW(), INTERVAL 1 DAY), 3),
(6, 3, 750.00, 75.00, DATE_SUB(NOW(), INTERVAL 17 DAY), DATE_SUB(NOW(), INTERVAL 5 DAY), DATE_SUB(NOW(), INTERVAL 17 DAY), DATE_SUB(NOW(), INTERVAL 5 DAY), 2),
(7, 2, 500.00, 50.00, DATE_SUB(NOW(), INTERVAL 14 DAY), DATE_SUB(NOW(), INTERVAL 6 DAY), DATE_SUB(NOW(), INTERVAL 14 DAY), DATE_SUB(NOW(), INTERVAL 6 DAY), 1),
(8, 4, 1000.00, 100.00, DATE_SUB(NOW(), INTERVAL 11 DAY), DATE_SUB(NOW(), INTERVAL 2 DAY), DATE_SUB(NOW(), INTERVAL 11 DAY), DATE_SUB(NOW(), INTERVAL 2 DAY), 2),
(9, 3, 750.00, 75.00, DATE_SUB(NOW(), INTERVAL 9 DAY), DATE_SUB(NOW(), INTERVAL 3 DAY), DATE_SUB(NOW(), INTERVAL 9 DAY), DATE_SUB(NOW(), INTERVAL 3 DAY), 2),
(10, 5, 1250.00, 125.00, DATE_SUB(NOW(), INTERVAL 7 DAY), DATE_SUB(NOW(), INTERVAL 1 DAY), DATE_SUB(NOW(), INTERVAL 7 DAY), DATE_SUB(NOW(), INTERVAL 1 DAY), 3),
(11, 2, 500.00, 50.00, DATE_SUB(NOW(), INTERVAL 6 DAY), DATE_SUB(NOW(), INTERVAL 4 DAY), DATE_SUB(NOW(), INTERVAL 6 DAY), DATE_SUB(NOW(), INTERVAL 4 DAY), 1),
(12, 3, 750.00, 75.00, DATE_SUB(NOW(), INTERVAL 5 DAY), DATE_SUB(NOW(), INTERVAL 2 DAY), DATE_SUB(NOW(), INTERVAL 5 DAY), DATE_SUB(NOW(), INTERVAL 2 DAY), 2),
(13, 4, 1000.00, 100.00, DATE_SUB(NOW(), INTERVAL 4 DAY), DATE_SUB(NOW(), INTERVAL 1 DAY), DATE_SUB(NOW(), INTERVAL 4 DAY), DATE_SUB(NOW(), INTERVAL 1 DAY), 2),
(14, 2, 500.00, 50.00, DATE_SUB(NOW(), INTERVAL 3 DAY), DATE_SUB(NOW(), INTERVAL 1 DAY), DATE_SUB(NOW(), INTERVAL 3 DAY), DATE_SUB(NOW(), INTERVAL 1 DAY), 1),
(15, 3, 750.00, 75.00, DATE_SUB(NOW(), INTERVAL 2 DAY), NOW(), DATE_SUB(NOW(), INTERVAL 2 DAY), NOW(), 2),
(16, 1, 250.00, 25.00, DATE_SUB(NOW(), INTERVAL 1 DAY), DATE_SUB(NOW(), INTERVAL 1 DAY), DATE_SUB(NOW(), INTERVAL 1 DAY), DATE_SUB(NOW(), INTERVAL 1 DAY), 1),
(17, 2, 500.00, 50.00, NOW(), NOW(), NOW(), NOW(), 1),
(18, 1, 250.00, 25.00, NOW(), NOW(), NOW(), NOW(), 1),
(19, 1, 250.00, 25.00, NOW(), NOW(), NOW(), NOW(), 1),
(20, 1, 250.00, 25.00, NOW(), NOW(), NOW(), NOW(), 1),
(21, 0, 0.00, 0.00, NULL, NULL, NULL, NULL, 0),
(22, 2, 500.00, 50.00, DATE_SUB(NOW(), INTERVAL 13 DAY), DATE_SUB(NOW(), INTERVAL 5 DAY), DATE_SUB(NOW(), INTERVAL 13 DAY), DATE_SUB(NOW(), INTERVAL 5 DAY), 1),
(23, 3, 750.00, 75.00, DATE_SUB(NOW(), INTERVAL 10 DAY), DATE_SUB(NOW(), INTERVAL 3 DAY), DATE_SUB(NOW(), INTERVAL 10 DAY), DATE_SUB(NOW(), INTERVAL 3 DAY), 2),
(24, 2, 500.00, 50.00, DATE_SUB(NOW(), INTERVAL 8 DAY), DATE_SUB(NOW(), INTERVAL 4 DAY), DATE_SUB(NOW(), INTERVAL 8 DAY), DATE_SUB(NOW(), INTERVAL 4 DAY), 1),
(25, 4, 1000.00, 100.00, DATE_SUB(NOW(), INTERVAL 12 DAY), DATE_SUB(NOW(), INTERVAL 2 DAY), DATE_SUB(NOW(), INTERVAL 12 DAY), DATE_SUB(NOW(), INTERVAL 2 DAY), 2),
(26, 3, 750.00, 75.00, DATE_SUB(NOW(), INTERVAL 15 DAY), DATE_SUB(NOW(), INTERVAL 6 DAY), DATE_SUB(NOW(), INTERVAL 15 DAY), DATE_SUB(NOW(), INTERVAL 6 DAY), 2),
(27, 5, 1250.00, 125.00, DATE_SUB(NOW(), INTERVAL 18 DAY), DATE_SUB(NOW(), INTERVAL 3 DAY), DATE_SUB(NOW(), INTERVAL 18 DAY), DATE_SUB(NOW(), INTERVAL 3 DAY), 3),
(28, 2, 500.00, 50.00, DATE_SUB(NOW(), INTERVAL 20 DAY), DATE_SUB(NOW(), INTERVAL 7 DAY), DATE_SUB(NOW(), INTERVAL 20 DAY), DATE_SUB(NOW(), INTERVAL 7 DAY), 1),
(29, 3, 750.00, 75.00, DATE_SUB(NOW(), INTERVAL 23 DAY), DATE_SUB(NOW(), INTERVAL 5 DAY), DATE_SUB(NOW(), INTERVAL 23 DAY), DATE_SUB(NOW(), INTERVAL 5 DAY), 2),
(30, 4, 1000.00, 100.00, DATE_SUB(NOW(), INTERVAL 26 DAY), DATE_SUB(NOW(), INTERVAL 4 DAY), DATE_SUB(NOW(), INTERVAL 26 DAY), DATE_SUB(NOW(), INTERVAL 4 DAY), 2);

-- 4. 创建会员标签关联
INSERT IGNORE INTO member_tag_links (member_id, tag_id, created_by) VALUES
(1, 1, 1), (1, 2, 1), (1, 5, 1),
(2, 2, 1), (2, 5, 1),
(3, 2, 1),
(4, 3, 1),
(5, 1, 1), (5, 2, 1), (5, 6, 1),
(6, 2, 1),
(7, 3, 1),
(8, 2, 1), (8, 5, 1),
(9, 2, 1),
(10, 1, 1), (10, 2, 1),
(11, 3, 1),
(12, 2, 1),
(13, 2, 1), (13, 5, 1),
(14, 3, 1),
(15, 2, 1),
(16, 3, 1),
(17, 3, 1),
(18, 3, 1),
(19, 3, 1),
(20, 3, 1),
(21, 4, 1),
(22, 2, 1),
(23, 2, 1),
(24, 3, 1),
(25, 2, 1), (25, 5, 1),
(26, 2, 1),
(27, 1, 1), (27, 2, 1),
(28, 3, 1),
(29, 2, 1),
(30, 2, 1);

-- 5. 创建会员品牌关联
INSERT IGNORE INTO member_brand_links (member_id, brand_id, first_campaign_id) VALUES
(1, 1, 1), (2, 1, 1), (3, 1, 1), (4, 1, 2), (5, 1, 2),
(6, 1, 1), (7, 1, 3), (8, 1, 3), (9, 1, 1), (10, 1, 2),
(11, 1, 1), (12, 1, 2), (13, 1, 3), (14, 1, 1), (15, 1, 2),
(16, 1, 3), (17, 1, 1), (18, 1, 2), (19, 1, 1), (20, 1, 3),
(22, 1, 2), (23, 1, 3), (24, 1, 1), (25, 1, 2), (26, 1, 3),
(27, 1, 1), (28, 1, 2), (29, 1, 3), (30, 1, 1);

-- ============================================
-- 第二部分：确保分销商数据存在
-- ============================================

-- 如果分销商数据不存在，重新导入
INSERT IGNORE INTO distributors (id, user_id, brand_id, level, parent_id, status, total_earnings, subordinates_count, created_at) VALUES
(1, 4, 1, 1, NULL, 'active', 339.00, 2, DATE_SUB(NOW(), INTERVAL 60 DAY)),
(2, 5, 1, 2, 1, 'active', 114.00, 1, DATE_SUB(NOW(), INTERVAL 55 DAY)),
(3, 6, 1, 3, 2, 'active', 35.40, 0, DATE_SUB(NOW(), INTERVAL 50 DAY)),
(101, 9, 2, 1, NULL, 'active', 54.40, 2, DATE_SUB(NOW(), INTERVAL 45 DAY)),
(102, 10, 2, 2, 101, 'pending', 17.60, 1, DATE_SUB(NOW(), INTERVAL 40 DAY)),
(103, 11, 2, 3, 102, 'suspended', 0.00, 0, DATE_SUB(NOW(), INTERVAL 35 DAY)),
(104, 14, 1, 1, NULL, 'pending', 0.00, 0, DATE_SUB(NOW(), INTERVAL 30 DAY)),
(105, 15, 1, 2, 1, 'suspended', 0.00, 0, DATE_SUB(NOW(), INTERVAL 25 DAY));

-- ============================================
-- 第三部分：统计信息
-- ============================================

SELECT '✓ 完整测试数据初始化完成' AS result;

SELECT '========================================' AS '';
SELECT '数据统计' AS '';
SELECT '========================================' AS '';

SELECT 
    '会员总数' AS item,
    COUNT(*) AS count
FROM members

UNION ALL

SELECT 
    '活跃会员',
    COUNT(*)
FROM members
WHERE status = 'active'

UNION ALL

SELECT 
    '会员画像',
    COUNT(*)
FROM member_profiles

UNION ALL

SELECT 
    '会员标签',
    COUNT(*)
FROM member_tags

UNION ALL

SELECT 
    '分销商总数',
    COUNT(*)
FROM distributors

UNION ALL

SELECT 
    '活跃分销商',
    COUNT(*)
FROM distributors
WHERE status = 'active'

UNION ALL

SELECT 
    '订单总数',
    COUNT(*)
FROM orders

UNION ALL

SELECT 
    '已支付订单',
    COUNT(*)
FROM orders
WHERE pay_status = 'paid';

-- 显示会员统计
SELECT '========================================' AS '';
SELECT '会员数据预览' AS '';
SELECT '========================================' AS '';

SELECT 
    m.id AS 'ID',
    m.nickname AS '昵称',
    m.phone AS '手机号',
    m.status AS '状态',
    mp.total_orders AS '订单数',
    mp.total_payment AS '累计支付',
    DATE_FORMAT(m.created_at, '%Y-%m-%d') AS '注册时间'
FROM members m
LEFT JOIN member_profiles mp ON m.id = mp.member_id
ORDER BY m.id
LIMIT 10;

-- 显示分销商统计
SELECT '========================================' AS '';
SELECT '分销商数据预览' AS '';
SELECT '========================================' AS '';

SELECT 
    d.id AS 'ID',
    u.username AS '用户名',
    d.level AS '级别',
    d.status AS '状态',
    d.total_earnings AS '累计收益',
    d.subordinates_count AS '下级数'
FROM distributors d
JOIN users u ON d.user_id = u.id
WHERE d.brand_id = 1
ORDER BY d.id;
