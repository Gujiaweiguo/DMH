# 实时奖励系统规范变更（会员关联）

## MODIFIED Requirements

### Requirement: 奖励与会员关联
系统 SHALL 将奖励记录与会员关联，用于会员维度的激励分析与生命周期分析。

#### Scenario: 创建奖励记录时写入 member 关联
- **WHEN** 系统为某订单结算奖励
- **THEN** 奖励记录 SHALL 关联到该订单对应会员（memberId）

#### Scenario: 会员维度奖励汇总
- **WHEN** 平台管理员查看会员详情
- **THEN** 系统 SHALL 展示该会员累计奖励、奖励次数、最近一次奖励时间

#### Scenario: 品牌隔离下的奖励查询
- **WHEN** 品牌管理员查看本品牌会员奖励汇总
- **THEN** 系统 SHALL 仅统计本品牌产生的奖励记录
