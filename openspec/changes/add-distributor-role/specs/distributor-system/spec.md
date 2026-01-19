# Spec: 分销商系统（Distributor System）

## ADDED Requirements

### Requirement: 分销商角色定义
系统 SHALL 定义 `distributor` 角色，作为高级顾客角色，具备分享推广并获得奖励的资格。

#### Scenario: 分销商角色创建
- **WHEN** 用户被批准成为分销商
- **THEN** 系统 SHALL 在 `user_roles` 表中为该用户添加 `distributor` 角色关联
- **AND** 在 `distributors` 表中创建分销商专属记录

#### Scenario: 分销商角色验证
- **WHEN** 用户访问分销商专属功能
- **THEN** 系统 SHALL 验证用户是否具有 `distributor` 角色
- **AND** 无该角色时 SHALL 返回 403 禁止访问错误

### Requirement: 分销商申请流程
系统 SHALL 支持普通顾客申请成为分销商，由品牌管理员审批。

#### Scenario: 顾客提交分销商申请
- **WHEN** 已登录顾客通过H5提交分销商申请
- **THEN** 系统 SHALL 创建申请记录，状态为 `pending`
- **AND** 关联申请顾客、品牌和申请理由
- **AND** 通知对应品牌管理员有待审批申请

#### Scenario: 重复申请检查
- **WHEN** 用户已存在待审批或已批准的分销商申请
- **THEN** 系统 SHALL 拒绝新的申请
- **AND** 提示用户当前申请状态

#### Scenario: 品牌管理员审批通过
- **WHEN** 品牌管理员批准分销商申请
- **THEN** 系统 SHALL 将用户角色添加 `distributor`
- **AND** 在 `distributors` 表创建记录，设置指定级别
- **AND** 记录审批人和审批时间
- **AND** 更新申请状态为 `approved`
- **AND** 通知用户审批结果

#### Scenario: 品牌管理员审批拒绝
- **WHEN** 品牌管理员拒绝分销商申请
- **THEN** 系统 SHALL 更新申请状态为 `rejected`
- **AND** 记录拒绝原因
- **AND** 通知用户审批结果

#### Scenario: 平台管理员查看全部分销商申请
- **WHEN** 平台管理员访问分销商申请列表
- **THEN** 系统 SHALL 返回所有品牌的待审批申请

#### Scenario: 品牌管理员查看本品牌分销商申请
- **WHEN** 品牌管理员访问分销商申请列表
- **THEN** 系统 SHALL 仅返回本品牌的待审批申请

### Requirement: 多级分销体系
系统 SHALL 支持最多3级分销体系，不同级别对应不同的奖励比例。

#### Scenario: 一级分销商奖励
- **WHEN** 新订单支付成功且有直接推荐人
- **THEN** 系统 SHALL 按一级奖励比例计算并分配奖励给直接推荐人
- **AND** 直接推荐人必须是分销商状态为 `active`

#### Scenario: 二级分销商奖励
- **WHEN** 一级分销商获得奖励且其有上级分销商
- **THEN** 系统 SHALL 按二级奖励比例计算并分配奖励给二级分销商

#### Scenario: 三级分销商奖励
- **WHEN** 二级分销商获得奖励且其有上级分销商
- **THEN** 系统 SHALL 按三级奖励比例计算并分配奖励给三级分销商

#### Scenario: 超过三级不分配奖励
- **WHEN** 推荐链超过三级
- **THEN** 系统 SHALL 只为前三级分销商分配奖励
- **AND** 第四级及之后不获得奖励

#### Scenario: 非分销商不获得推荐奖励
- **WHEN** 推荐人不是分销商或分销商状态非 `active`
- **THEN** 系统 SHALL 不为该推荐人分配奖励
- **AND** 继续向上查找是否有分销商推荐人

### Requirement: 分销商级别管理
系统 SHALL 支持品牌管理员或平台管理员调整分销商级别。

#### Scenario: 手动升级分销商级别
- **WHEN** 品牌管理员或平台管理员手动调整分销商级别
- **THEN** 系统 SHALL 更新 `distributors` 表中的 level 字段
- **AND** 记录级别变更日志

#### Scenario: 分销商级别变更影响后续奖励
- **WHEN** 分销商级别发生变更
- **THEN** 新级别 SHALL 适用于级别变更后的新订单
- **AND** 历史奖励保持不变

### Requirement: 分销商专属推广工具
系统 SHALL 为分销商提供专属推广链接和二维码生成功能。

#### Scenario: 生成专属推广链接
- **WHEN** 分销商请求生成推广链接
- **THEN** 系统 SHALL 生成包含分销商ID的专属链接
- **AND** 链接格式为 `{campaign_url}?distributor_id={id}`
- **AND** 记录链接生成时间和使用次数

#### Scenario: 生成推广二维码
- **WHEN** 分销商请求生成推广二维码
- **THEN** 系统 SHALL 基于专属推广链接生成二维码图片
- **AND** 返回二维码图片URL

#### Scenario: 推广链接访问追踪
- **WHEN** 用户通过分销商推广链接访问活动
- **THEN** 系统 SHALL 记录访问来源为该分销商
- **AND** 在用户最终下单时将分销商设为推荐人

### Requirement: 分销商数据查看
系统 SHALL 允许分销商查看自己的推广数据和奖励明细。

#### Scenario: 查看推广数据统计
- **WHEN** 分销商访问推广数据页面
- **THEN** 系统 SHALL 返回以下统计数据：
  - 累计订单数
  - 累计奖励金额
  - 下级分销商数量（仅一级）
  - 本月/本周新增订单

#### Scenario: 查看奖励明细
- **WHEN** 分销商访问奖励明细页面
- **THEN** 系统 SHALL 返回分页的奖励记录列表
- **AND** 每条记录包含：订单ID、奖励金额、奖励时间、来源订单

#### Scenario: 查看下级分销商列表
- **WHEN** 分销商访问下级列表页面
- **THEN** 系统 SHALL 返回其直接下级分销商列表（仅一级）
- **AND** 每个下级显示：姓名、级别、加入时间、累计订单数

#### Scenario: 数据隔离
- **WHEN** 分销商查看数据
- **THEN** 系统 SHALL 只返回该分销商自己的数据
- **AND** 不能查看其他分销商的数据

### Requirement: 分销商状态管理
系统 SHALL 支持分销商的激活、暂停和禁用状态。

#### Scenario: 暂停分销商资格
- **WHEN** 品牌管理员或平台管理员暂停分销商
- **THEN** 系统 SHALL 更新分销商状态为 `suspended`
- **AND** 暂停后的分销商不再获得新的推荐奖励
- **AND** 历史奖励保持不变

#### Scenario: 重新激活分销商
- **WHEN** 管理员重新激活被暂停的分销商
- **THEN** 系统 SHALL 更新分销商状态为 `active`
- **AND** 恢复后的分销商可继续获得推荐奖励

### Requirement: 分销商与品牌关联
系统 SHALL 支持分销商与特定品牌的关联关系。

#### Scenario: 申请时选择品牌
- **WHEN** 顾客申请成为分销商
- **THEN** 系统 SHALL 要求顾客选择申请成为哪个品牌的分销商

#### Scenario: 多品牌分销商
- **WHEN** 用户被批准为多个品牌的分销商
- **THEN** 系统 SHALL 为每个品牌创建独立的分销商记录
- **AND** 奖励按品牌分别计算

#### Scenario: 品牌管理员只能管理本品牌分销商
- **WHEN** 品牌管理员管理分销商
- **THEN** 系统 SHALL 只允许管理本品牌的分销商
- **AND** 不能操作其他品牌的分销商
