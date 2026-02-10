# feedback-system Specification

## Purpose
实现用户反馈从提交、流转处理到统计分析的完整闭环能力，覆盖功能反馈、满意度调查、FAQ管理与功能使用行为统计，帮助产品与运营团队持续改进体验和质量。

## Requirements

### Requirement: 用户反馈提交
系统SHALL允许用户提交功能反馈，包含分类、评分、描述等信息。

#### Scenario: 用户提交功能反馈成功
- **WHEN** 已登录用户提交有效的反馈（包含标题、内容、类别）
- **THEN** 系统SHALL创建反馈记录并分配默认优先级(medium)
- **AND** 系统SHALL记录用户ID、设备信息、浏览器信息
- **AND** 系统SHALL返回创建成功的反馈ID

#### Scenario: 反馈参数验证
- **WHEN** 用户提交标题或内容为空的反馈
- **THEN** 系统SHALL返回400错误并提示必填字段
- **AND** 评分范围必须在1-5之间

#### Scenario: 反馈分类支持
- **WHEN** 用户提交反馈时选择分类
- **THEN** 系统SHALL支持以下分类: poster(海报), payment(支付), verification(核销), other(其他)
- **AND** 每个分类可包含可选的子分类

#### Scenario: 反馈优先级设置
- **WHEN** 管理员审核反馈
- **THEN** 系统SHALL支持优先级: low, medium, high
- **AND** 用户未指定时默认为medium

### Requirement: 反馈列表查询
系统SHALL根据用户角色和筛选条件返回反馈列表。

#### Scenario: 普通用户查看自己的反馈
- **WHEN** 普通用户请求反馈列表
- **THEN** 系统SHALL只返回该用户创建的反馈
- **AND** 支持按类别、状态、优先级筛选
- **AND** 支持分页查询

#### Scenario: 管理员查看所有反馈
- **WHEN** 平台管理员或品牌管理员请求反馈列表
- **THEN** 系统SHALL返回所有反馈
- **AND** 支持按用户ID、类别、状态、优先级筛选
- **AND** 按创建时间倒序排列

#### Scenario: 反馈状态筛选
- **WHEN** 用户筛选反馈时指定状态
- **THEN** 系统SHALL支持状态: pending(待处理), reviewing(审核中), resolved(已解决), closed(已关闭)
- **AND** 只返回匹配状态的反馈

### Requirement: 反馈详情查看
系统SHALL允许用户查看单个反馈的详细信息。

#### Scenario: 查看自己的反馈详情
- **WHEN** 用户请求查看自己创建的反馈
- **THEN** 系统SHALL返回完整的反馈信息
- **AND** 包含标题、内容、类别、评分、状态、处理回复等

#### Scenario: 管理员查看任意反馈
- **WHEN** 管理员请求查看任意反馈
- **THEN** 系统SHALL返回该反馈的完整信息
- **AND** 包含用户信息、设备信息、浏览器信息等元数据

#### Scenario: 反馈访问权限控制
- **WHEN** 普通用户尝试访问他人的反馈
- **THEN** 系统SHALL返回403禁止访问错误

### Requirement: 反馈状态更新
系统SHALL允许管理员更新反馈状态并添加处理回复。

#### Scenario: 管理员更新反馈状态
- **WHEN** 管理员更新反馈状态并指定处理人
- **THEN** 系统SHALL更新反馈的状态和指派信息
- **AND** 记录处理人ID和更新时间

#### Scenario: 管理员添加处理回复
- **WHEN** 管理员添加处理回复内容
- **THEN** 系统SHALL保存回复内容
- **AND** 更新反馈状态为"reviewing"或"resolved"

#### Scenario: 反馈解决时间记录
- **WHEN** 反馈状态更新为"resolved"
- **THEN** 系统SHALL记录解决时间(resolved_at)
- **AND** 用于计算平均解决时间统计

#### Scenario: 状态更新权限控制
- **WHEN** 非管理员尝试更新反馈状态
- **THEN** 系统SHALL返回403禁止访问错误

### Requirement: 满意度调查
系统SHALL允许用户对系统功能进行满意度评分。

#### Scenario: 提交满意度调查
- **WHEN** 已登录用户提交功能满意度评分
- **THEN** 系统SHALL记录用户ID、功能名称、各维度评分
- **AND** 评分维度包括: 易用性、性能、稳定性、整体满意度、推荐意愿

#### Scenario: 满意度评分验证
- **WHEN** 用户提交满意度调查
- **THEN** 系统SHALL验证所有评分在1-5范围
- **AND** 功能名称不能为空

#### Scenario: 记录详细反馈
- **WHEN** 用户在满意度调查中填写反馈
- **THEN** 系统SHALL记录最喜欢的功能、最不喜欢的功能、改进建议
- **AND** 记录希望增加的新功能需求

### Requirement: FAQ管理
系统SHALL提供常见问题(FAQ)的查询和反馈功能。

#### Scenario: 查询已发布的FAQ
- **WHEN** 用户请求FAQ列表
- **THEN** 系统SHALL返回所有已发布的FAQ
- **AND** 支持按分类筛选
- **AND** 支持关键词搜索（匹配问题或答案）

#### Scenario: FAQ点赞反馈
- **WHEN** 用户标记FAQ有帮助或无帮助
- **THEN** 系统SHALL增加对应的计数器
- **AND** 防止同一用户重复投票

#### Scenario: FAQ排序
- **WHEN** 返回FAQ列表
- **THEN** 系统SHALL按sort_order升序排列
- **AND** 同排序的按创建时间降序

### Requirement: 功能使用统计
系统SHALL记录用户功能使用情况用于分析和优化。

#### Scenario: 记录功能使用事件
- **WHEN** 用户使用系统功能（如生成海报、支付、核销等）
- **THEN** 系统SHALL记录功能名称、操作类型、执行结果
- **AND** 记录执行耗时（毫秒）
- **AND** 失败时记录错误信息

#### Scenario: 关联活动ID
- **WHEN** 使用记录与特定活动相关
- **THEN** 系统SHALL记录活动ID用于关联分析

#### Scenario: 异步记录使用数据
- **WHEN** 记录功能使用数据
- **THEN** 系统SHALL使用异步方式避免影响主流程性能

### Requirement: 反馈统计分析
系统SHALL为管理员提供反馈数据的统计分析功能。

#### Scenario: 反馈总数和分类统计
- **WHEN** 管理员请求反馈统计
- **THEN** 系统SHALL返回反馈总数
- **AND** 按类别统计反馈数量
- **AND** 按状态统计反馈数量
- **AND** 按优先级统计反馈数量

#### Scenario: 评分和解决率统计
- **WHEN** 管理员请求评分统计
- **THEN** 系统SHALL计算平均评分
- **AND** 计算解决率(已解决数/总数)
- **AND** 计算平均解决时间（小时）

#### Scenario: 评分分布统计
- **WHEN** 管理员请求评分分布
- **THEN** 系统SHALL返回1-5星的评分数量分布
- **AND** 用于绘制评分分布图

#### Scenario: 时间范围筛选
- **WHEN** 管理员指定统计时间范围
- **THEN** 系统SHALL只统计该时间范围内的反馈数据
- **AND** 支持日期格式YYYY-MM-DD

#### Scenario: 统计访问权限
- **WHEN** 非管理员请求统计数据
- **THEN** 系统SHALL返回403禁止访问错误

### Requirement: 反馈数据模型
系统SHALL使用以下数据模型存储反馈相关数据。

#### Scenario: UserFeedback模型
- **WHEN** 存储用户反馈
- **THEN** 系统SHALL使用UserFeedback表
- **AND** 包含字段: id, user_id, category, subcategory, rating, title, content
- **AND** 包含字段: feature_use_case, device_info, browser_info, priority, status
- **AND** 包含字段: assignee_id, response, resolved_at, created_at, updated_at

#### Scenario: FeatureSatisfactionSurvey模型
- **WHEN** 存储满意度调查
- **THEN** 系统SHALL使用FeatureSatisfactionSurvey表
- **AND** 包含字段: id, user_id, user_role, feature, ease_of_use, performance
- **AND** 包含字段: reliability, overall_satisfaction, would_recommend
- **AND** 包含字段: most_liked, least_liked, improvement_suggestions, would_like_more_features
- **AND** 包含字段: created_at

#### Scenario: FAQItem模型
- **WHEN** 存储FAQ
- **THEN** 系统SHALL使用FAQItem表
- **AND** 包含字段: id, category, question, answer, sort_order
- **AND** 包含字段: is_published, view_count, helpful_count, not_helpful_count
- **AND** 包含字段: created_at, updated_at

#### Scenario: FeatureUsageStat模型
- **WHEN** 存储功能使用统计
- **THEN** 系统SHALL使用FeatureUsageStat表
- **AND** 包含字段: id, user_id, user_role, feature, action, campaign_id
- **AND** 包含字段: success, duration_ms, error_message
- **AND** 包含字段: created_at

#### Scenario: FeedbackTag模型
- **WHEN** 存储反馈标签
- **THEN** 系统SHALL使用FeedbackTag表
- **AND** 包含字段: id, name, color, created_at

#### Scenario: FeedbackTagRelation模型
- **WHEN** 关联反馈和标签
- **THEN** 系统SHALL使用FeedbackTagRelation表
- **AND** 包含字段: feedback_id, tag_id
- **AND** 支持多对多关系
