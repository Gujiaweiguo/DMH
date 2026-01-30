# RBAC 权限管理规范变更（会员与导出审批）

## ADDED Requirements

### Requirement: 会员数据访问权限
系统 SHALL 为会员相关功能定义可分配的权限点，并按角色限制访问范围。

#### Scenario: 平台管理员访问全量会员
- **WHEN** platform_admin 访问会员列表/详情
- **THEN** 系统 SHALL 允许访问全量会员

#### Scenario: 品牌管理员访问本品牌会员
- **WHEN** brand_admin 访问会员列表/详情
- **THEN** 系统 SHALL 仅允许访问其品牌关联会员

#### Scenario: 非授权角色访问
- **WHEN** participant 或 anonymous 访问会员管理接口
- **THEN** 系统 SHALL 返回 403 Forbidden

### Requirement: 会员合并权限
系统 SHALL 限制会员合并为平台管理员专属能力。

#### Scenario: brand_admin 尝试合并
- **WHEN** brand_admin 发起会员合并
- **THEN** 系统 SHALL 拒绝并返回 403

#### Scenario: platform_admin 合并
- **WHEN** platform_admin 发起会员合并
- **THEN** 系统 SHALL 允许执行并记录审计日志

### Requirement: 导出审批权限
系统 SHALL 支持“申请导出”与“审批导出”分离权限。

#### Scenario: brand_admin 申请导出
- **WHEN** brand_admin 发起会员导出
- **THEN** 系统 SHALL 允许创建导出申请并进入待审批

#### Scenario: platform_admin 审批导出
- **WHEN** platform_admin 审批导出申请
- **THEN** 系统 SHALL 允许审批通过/驳回并记录审批信息

#### Scenario: 未审批前不可下载
- **WHEN** 导出申请处于待审批或已驳回
- **THEN** 系统 SHALL 不提供下载链接
