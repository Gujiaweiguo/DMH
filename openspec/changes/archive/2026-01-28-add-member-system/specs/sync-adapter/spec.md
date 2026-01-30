# 外网数据同步适配器规范变更（会员字段）

## MODIFIED Requirements

### Requirement: 订单与奖励同步包含会员字段
当启用外网同步时，系统 SHALL 在订单与奖励同步数据中包含 `memberId` 与 `unionid` 字段，以便外部系统进行会员维度分析。

#### Scenario: 同步订单包含 memberId/unionid
- **WHEN** 同步一笔已关联会员的订单数据
- **THEN** 同步 payload SHALL 包含 memberId 与 unionid

#### Scenario: 同步奖励包含 memberId/unionid
- **WHEN** 同步一笔奖励数据
- **THEN** 同步 payload SHALL 包含 memberId 与 unionid

#### Scenario: 兼容缺失会员关联的数据
- **WHEN** 同步一笔未关联会员的数据（例如历史数据）
- **THEN** memberId/unionid 字段 SHALL 允许为空
- **AND** 同步仍可成功完成
