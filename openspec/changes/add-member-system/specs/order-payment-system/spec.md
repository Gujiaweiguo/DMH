# 订单与支付系统规范变更（会员关联）

## MODIFIED Requirements

### Requirement: 订单与会员关联
系统 SHALL 在订单创建与支付流程中将订单与会员关联，以支持会员维度的支付分析。

#### Scenario: 创建订单时写入 member 关联
- **WHEN** 创建订单请求包含 `unionid`
- **THEN** 系统 SHALL 将订单关联到对应会员（memberId）
- **AND** 订单查询接口返回 memberId

#### Scenario: unionid 缺失时的订单创建约束
- **WHEN** 创建订单请求不包含 `unionid`
- **THEN** 系统 SHALL 拒绝创建订单
- **AND** 返回可理解的错误提示说明需要完成授权

#### Scenario: 支付回调后会员指标可汇总
- **WHEN** 订单支付成功并触发回调
- **THEN** 系统 SHALL 能基于 memberId 汇总支付金额与次数
