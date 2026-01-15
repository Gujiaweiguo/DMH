# 移动端落地页规范变更（UnionID 会员关联）

## MODIFIED Requirements

### Requirement: 获取 unionid 并用于参与/下单
移动端落地页 SHALL 在用户参与/下单前获取 `unionid` 并随请求提交，用于会员关联。

#### Scenario: 微信环境获取 unionid 成功
- **GIVEN** 用户在微信环境打开活动落地页并完成授权
- **WHEN** 页面发起参与/下单
- **THEN** 请求参数 SHALL 包含 `unionid`
- **AND** 后端可据此创建/关联会员

#### Scenario: unionid 获取失败阻止参与
- **GIVEN** 页面无法获取 `unionid`
- **WHEN** 用户尝试提交参与/下单
- **THEN** 页面 SHALL 阻止提交
- **AND** 提示用户完成授权或在支持的环境中打开

#### Scenario: unionid 缓存与复用
- **GIVEN** 用户已成功获取过 `unionid`
- **WHEN** 用户在同一会话再次参与其他活动
- **THEN** 页面 SHALL 复用已获取的 `unionid` 以减少重复授权
