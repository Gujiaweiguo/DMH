-- ============================================
-- 填充测试数据：核销记录和海报记录
-- 创建时间: 2026-01-29
-- ============================================

USE dmh;

-- 填充核销记录测试数据（20条）
INSERT INTO verification_records (order_id, verification_status, verified_at, verified_by, verification_code, verification_method, remark, created_at) VALUES
(1, 'verified', '2026-01-22 10:30:00', 1, 'VR001', 'manual', '现场核销', '2026-01-22 10:25:00'),
(2, 'verified', '2026-01-23 14:15:00', 2, 'VR002', 'qrcode', '扫码核销', '2026-01-23 14:10:00'),
(3, 'pending', NULL, NULL, 'VR003', 'manual', '', '2026-01-24 09:20:00'),
(4, 'verified', '2026-01-24 16:45:00', 1, 'VR004', 'auto', '系统自动核销', '2026-01-24 16:40:00'),
(5, 'cancelled', '2026-01-25 11:30:00', 2, 'VR005', 'manual', '用户取消', '2026-01-25 11:25:00'),
(6, 'verified', '2026-01-25 15:20:00', 1, 'VR006', 'qrcode', '扫码核销', '2026-01-25 15:15:00'),
(7, 'pending', NULL, NULL, 'VR007', 'manual', '', '2026-01-26 08:30:00'),
(8, 'verified', '2026-01-26 10:45:00', 2, 'VR008', 'manual', '现场核销', '2026-01-26 10:40:00'),
(9, 'verified', '2026-01-26 14:00:00', 1, 'VR009', 'auto', '系统自动核销', '2026-01-26 13:55:00'),
(10, 'pending', NULL, NULL, 'VR010', 'manual', '', '2026-01-27 09:00:00'),
(11, 'verified', '2026-01-27 11:30:00', 2, 'VR011', 'qrcode', '扫码核销', '2026-01-27 11:25:00'),
(12, 'verified', '2026-01-27 16:00:00', 1, 'VR012', 'manual', '现场核销', '2026-01-27 15:55:00'),
(13, 'cancelled', '2026-01-28 09:20:00', 1, 'VR013', 'manual', '用户取消', '2026-01-28 09:15:00'),
(14, 'pending', NULL, NULL, 'VR014', 'manual', '', '2026-01-28 10:30:00'),
(15, 'verified', '2026-01-28 13:45:00', 2, 'VR015', 'qrcode', '扫码核销', '2026-01-28 13:40:00'),
(16, 'verified', '2026-01-28 15:30:00', 1, 'VR016', 'auto', '系统自动核销', '2026-01-28 15:25:00'),
(17, 'pending', NULL, NULL, 'VR017', 'manual', '', '2026-01-29 08:00:00'),
(18, 'pending', NULL, NULL, 'VR018', 'manual', '', '2026-01-29 09:30:00'),
(19, 'verified', '2026-01-29 11:15:00', 2, 'VR019', 'qrcode', '扫码核销', '2026-01-29 11:10:00'),
(20, 'pending', NULL, NULL, 'VR020', 'manual', '', '2026-01-29 13:00:00');

-- 填充海报生成记录测试数据（20条）
INSERT INTO poster_records (record_type, campaign_id, distributor_id, template_name, poster_url, thumbnail_url, file_size, generation_time, download_count, share_count, generated_by, status, created_at) VALUES
('personal', 1, 3, '红色主题模板', 'https://example.com/posters/poster_001.jpg', 'https://example.com/posters/thumb_001.jpg', '2.5MB', 1200, 5, 3, 3, 'active', '2026-01-22 09:00:00'),
('personal', 1, 4, '蓝色主题模板', 'https://example.com/posters/poster_002.jpg', 'https://example.com/posters/thumb_002.jpg', '3.1MB', 1500, 8, 5, 4, 'active', '2026-01-22 10:30:00'),
('brand', 1, 0, '品牌官方海报', 'https://example.com/posters/poster_003.jpg', 'https://example.com/posters/thumb_003.jpg', '4.2MB', 1800, 15, 12, 1, 'active', '2026-01-22 14:00:00'),
('personal', 2, 3, '新年促销模板', 'https://example.com/posters/poster_004.jpg', 'https://example.com/posters/thumb_004.jpg', '2.8MB', 1100, 3, 2, 3, 'active', '2026-01-23 08:45:00'),
('personal', 2, 5, '节日特惠模板', 'https://example.com/posters/poster_005.jpg', 'https://example.com/posters/thumb_005.jpg', '3.5MB', 1600, 7, 4, 5, 'active', '2026-01-23 11:20:00'),
('brand', 2, 0, '品牌活动海报', 'https://example.com/posters/poster_006.jpg', 'https://example.com/posters/thumb_006.jpg', '3.8MB', 1700, 20, 18, 1, 'active', '2026-01-23 15:30:00'),
('personal', 3, 4, '春季主题模板', 'https://example.com/posters/poster_007.jpg', 'https://example.com/posters/thumb_007.jpg', '2.6MB', 1300, 4, 3, 4, 'active', '2026-01-24 09:15:00'),
('personal', 3, 3, '夏日促销模板', 'https://example.com/posters/poster_008.jpg', 'https://example.com/posters/thumb_008.jpg', '2.9MB', 1400, 6, 4, 3, 'active', '2026-01-24 13:45:00'),
('brand', 3, 0, '品牌宣传海报', 'https://example.com/posters/poster_009.jpg', 'https://example.com/posters/thumb_009.jpg', '4.0MB', 1900, 12, 10, 1, 'active', '2026-01-24 16:20:00'),
('personal', 4, 5, '秋季优惠模板', 'https://example.com/posters/poster_010.jpg', 'https://example.com/posters/thumb_010.jpg', '3.2MB', 1550, 5, 3, 5, 'active', '2026-01-25 10:00:00'),
('personal', 4, 4, '冬季清仓模板', 'https://example.com/posters/poster_011.jpg', 'https://example.com/posters/thumb_011.jpg', '3.4MB', 1650, 9, 6, 4, 'active', '2026-01-25 14:30:00'),
('brand', 4, 0, '品牌新品海报', 'https://example.com/posters/poster_012.jpg', 'https://example.com/posters/thumb_012.jpg', '3.7MB', 1750, 18, 15, 1, 'active', '2026-01-25 17:00:00'),
('personal', 5, 3, '会员专享模板', 'https://example.com/posters/poster_013.jpg', 'https://example.com/posters/thumb_013.jpg', '2.7MB', 1250, 2, 1, 3, 'active', '2026-01-26 08:30:00'),
('personal', 5, 5, '限时抢购模板', 'https://example.com/posters/poster_014.jpg', 'https://example.com/posters/thumb_014.jpg', '3.0MB', 1450, 8, 5, 5, 'active', '2026-01-26 12:15:00'),
('brand', 5, 0, '品牌周年海报', 'https://example.com/posters/poster_015.jpg', 'https://example.com/posters/thumb_015.jpg', '4.1MB', 1850, 22, 19, 1, 'active', '2026-01-26 15:45:00'),
('personal', 1, 4, '红色主题模板', 'https://example.com/posters/poster_016.jpg', 'https://example.com/posters/thumb_016.jpg', '2.4MB', 1150, 6, 4, 4, 'active', '2026-01-27 09:45:00'),
('personal', 2, 3, '蓝色主题模板', 'https://example.com/posters/poster_017.jpg', 'https://example.com/posters/thumb_017.jpg', '3.3MB', 1580, 10, 7, 3, 'active', '2026-01-27 13:20:00'),
('brand', 1, 0, '品牌活动海报', 'https://example.com/posters/poster_018.jpg', 'https://example.com/posters/thumb_018.jpg', '3.9MB', 1780, 14, 11, 1, 'active', '2026-01-27 16:50:00'),
('personal', 3, 5, '新年促销模板', 'https://example.com/posters/poster_019.jpg', 'https://example.com/posters/thumb_019.jpg', '2.9MB', 1390, 3, 2, 5, 'active', '2026-01-28 10:10:00'),
('personal', 4, 4, '节日特惠模板', 'https://example.com/posters/poster_020.jpg', 'https://example.com/posters/thumb_020.jpg', '3.6MB', 1620, 7, 5, 4, 'active', '2026-01-28 14:40:00');
