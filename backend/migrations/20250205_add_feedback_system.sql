-- 创建用户反馈表
CREATE TABLE IF NOT EXISTS user_feedback (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '反馈ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    user_name VARCHAR(100) DEFAULT '' COMMENT '用户姓名（可选）',
    user_role VARCHAR(50) DEFAULT 'participant' COMMENT '用户角色：platform_admin, brand_admin, distributor, participant',
    category VARCHAR(50) NOT NULL COMMENT '反馈类别：poster, payment, verification, other',
    subcategory VARCHAR(50) DEFAULT '' COMMENT '子类别',
    rating TINYINT DEFAULT NULL COMMENT '评分（1-5星）',
    title VARCHAR(200) NOT NULL COMMENT '反馈标题',
    content TEXT NOT NULL COMMENT '反馈内容',
    feature_use_case VARCHAR(500) DEFAULT '' COMMENT '使用场景描述',
    device_info VARCHAR(200) DEFAULT '' COMMENT '设备信息',
    browser_info VARCHAR(200) DEFAULT '' COMMENT '浏览器信息',
    error_message VARCHAR(500) COMMENT '错误信息',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_user_id (user_id),
    INDEX idx_feature (feature),
    INDEX idx_action (action),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='功能使用统计表';

-- 创建功能满意度调查表
CREATE TABLE IF NOT EXISTS feature_satisfaction_survey (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '调查ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    user_role VARCHAR(50) DEFAULT 'participant' COMMENT '用户角色',
    feature VARCHAR(50) NOT NULL COMMENT '功能名称',
    ease_of_use TINYINT DEFAULT NULL COMMENT '易用性（1-5）',
    performance TINYINT DEFAULT NULL COMMENT '性能满意度（1-5）',
    reliability TINYINT DEFAULT NULL COMMENT '稳定性满意度（1-5）',
    overall_satisfaction TINYINT DEFAULT NULL COMMENT '整体满意度（1-5）',
    would_recommend TINYINT DEFAULT NULL COMMENT '推荐意愿（1-5）',
    most_liked VARCHAR(500) DEFAULT '' COMMENT '最满意的方面',
    least_liked VARCHAR(500) DEFAULT '' COMMENT '最不满意的方面',
    improvement_suggestions TEXT COMMENT '改进建议',
    would_like_more_features VARCHAR(500) DEFAULT '' COMMENT '希望增加的功能',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_user_id (user_id),
    INDEX idx_feature (feature),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='功能满意度调查表';

-- 创建常见问题表
CREATE TABLE IF NOT EXISTS faq_items (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT 'FAQ ID',
    category VARCHAR(50) NOT NULL COMMENT '分类：poster, payment, verification, general',
    question VARCHAR(500) NOT NULL COMMENT '问题',
    answer TEXT NOT NULL COMMENT '答案',
    sort_order INT DEFAULT 0 COMMENT '排序序号',
    view_count INT DEFAULT 0 COMMENT '浏览次数',
    helpful_count INT DEFAULT 0 COMMENT '有帮助次数',
    not_helpful_count INT DEFAULT 0 COMMENT '无帮助次数',
    is_published BOOLEAN DEFAULT TRUE COMMENT '是否发布',
    created_by BIGINT DEFAULT NULL COMMENT '创建人ID',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_category (category),
    INDEX idx_is_published (is_published)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='常见问题表';

-- 创建反馈标签表
CREATE TABLE IF NOT EXISTS feedback_tags (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '标签ID',
    name VARCHAR(50) NOT NULL UNIQUE COMMENT '标签名称',
    description VARCHAR(200) DEFAULT '' COMMENT '标签描述',
    color VARCHAR(20) DEFAULT '#1890ff' COMMENT '标签颜色',
    usage_count INT DEFAULT 0 COMMENT '使用次数',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='反馈标签表';

-- 创建反馈-标签关联表
CREATE TABLE IF NOT EXISTS feedback_tag_relations (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '关联ID',
    feedback_id BIGINT NOT NULL COMMENT '反馈ID',
    tag_id BIGINT NOT NULL COMMENT '标签ID',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_feedback_id (feedback_id),
    INDEX idx_tag_id (tag_id),
    UNIQUE KEY uk_feedback_tag (feedback_id, tag_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='反馈标签关联表';

-- 插入默认标签
INSERT INTO feedback_tags (name, description, color) VALUES
('功能建议', '用户提出的改进建议', '#52c41a'),
('问题反馈', '使用过程中遇到的问题', '#f5222d'),
('体验优化', '用户体验相关建议', '#1890ff'),
('性能问题', '性能相关反馈', '#faad14'),
('界面问题', 'UI/UX相关反馈', '#722ed1')
ON DUPLICATE KEY UPDATE name=name;

-- 插入默认FAQ
INSERT INTO faq_items (category, question, answer, sort_order) VALUES
('poster', '海报生成失败怎么办？', '请检查：1. 网络连接是否正常；2. 后端服务是否运行（访问 http://localhost:8889/health）；3. 活动信息是否完整。如果问题依然存在，请联系技术支持。', 1),
('poster', '海报生成有频率限制吗？', '是的，为防止滥用，每个用户每分钟最多生成 5 张海报。超出限制会提示「请求过于频繁」，请等待 1 分钟后重试。', 2),
('payment', '支付二维码有效期多久？', '支付二维码缓存 2 小时，过期后会自动刷新。您也可以手动点击「刷新二维码」获取新的二维码。', 1),
('payment', '订金模式和全款模式有什么区别？', '订金模式：用户先支付订金预订名额，到现场支付余款后核销。适合高价值活动。全款模式：用户一次性支付全款，支付后可直接参与活动。适合标准化活动。', 2),
('verification', '核销码在哪里查看？', '用户登录 H5 前端后，进入「我的订单」页面，点击已支付的订单即可查看核销码。', 1),
('verification', '核销码的有效期是多久？', '核销码长期有效，建议在活动期间内使用。', 2),
('general', '如何联系技术支持？', '技术支持邮箱：support@example.com，工作时间：周一至周五 9:00-18:00。紧急问题请拨打：+86-xxx-xxxx-xxxx', 1)
ON DUPLICATE KEY UPDATE question=question;
