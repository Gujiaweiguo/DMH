# Change: 实现 Feedback 反馈系统

## Why
当前 feedbacklogic.go 中有 9 个 TODO 待实现，包括用户反馈创建、列表查询、状态更新、满意度调查、FAQ 管理、使用统计等功能。实现完整的反馈系统可以收集用户意见、改进产品质量、提升用户体验。

## What Changes
- 创建 6 个数据模型：UserFeedback, FeatureSatisfactionSurvey, FAQItem, FeatureUsageStat, FeedbackTag, FeedbackTagRelation
- 实现 feedbacklogic.go 中的 9 个 TODO 函数
- 完善 types.go 中的响应类型（添加 CreatedAt 字段等）
- 实现权限控制（管理员 vs 普通用户）
- 实现数据统计和聚合查询

## Impact
- Affected specs: feedback-system (new capability)
- Affected code:
  - backend/model/ (新增 6 个模型文件)
  - backend/api/internal/logic/feedback/feedbacklogic.go (移除 TODO，实现逻辑)
  - backend/api/internal/types/types.go (完善类型定义)
- Dependencies: GORM (已有)
