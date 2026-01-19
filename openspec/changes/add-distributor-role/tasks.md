# Tasks: 分销商系统实现

## 1. 数据库设计与迁移
- [ ] 1.1 创建 `distributors` 表
  - id, user_id, brand_id, level, parent_id, status
  - approved_by, approved_at, total_earnings, subordinates_count
  - 索引：user_id, brand_id, parent_id, status
- [ ] 1.2 创建 `distributor_applications` 表
  - id, user_id, brand_id, status, reason
  - reviewed_by, reviewed_at, review_notes
  - 索引：user_id, brand_id, status
- [ ] 1.3 创建 `distributor_level_rewards` 表
  - brand_id, level(1/2/3), reward_percentage
  - 支持按品牌自定义各级奖励比例
- [ ] 1.4 修改 `roles` 表，插入 distributor 角色记录
- [ ] 1.5 创建数据库迁移脚本

## 2. 后端API实现
- [ ] 2.1 实现分销商申请API
  - POST /api/v1/distributor/apply - 提交申请
  - GET /api/v1/distributor/applications - 查询申请列表
  - GET /api/v1/distributor/applications/:id - 查询申请详情
- [ ] 2.2 实现分销商审批API
  - POST /api/v1/brands/:brandId/distributors/approve/:id - 审批申请
  - PUT /api/v1/distributors/:id/level - 调整级别
  - PUT /api/v1/distributors/:id/status - 暂停/激活
- [ ] 2.3 实现分销商推广工具API
  - POST /api/v1/distributor/link/generate - 生成推广链接
  - GET /api/v1/distributor/qrcode/:campaignId - 生成二维码
- [ ] 2.4 实现分销商数据查询API
  - GET /api/v1/distributor/statistics - 推广数据统计
  - GET /api/v1/distributor/rewards - 奖励明细
  - GET /api/v1/distributor/subordinates - 下级列表
- [ ] 2.5 实现管理后台分销商管理API
  - GET /api/v1/brands/:brandId/distributors - 分销商列表
  - GET /api/v1/distributors/:id - 分销商详情
- [ ] 2.6 修改奖励计算逻辑，支持多级分销

## 3. 前端H5实现（复用现有H5）
- [ ] 3.1 创建分销中心页面 (DistributorCenterView.vue)
  - 入口：个人中心 → 分销中心
  - 展示：当前级别、累计收益、下级数量
- [ ] 3.2 实现分销商申请页面
  - 申请表单：选择品牌、填写理由
  - 申请状态查看
- [ ] 3.3 实现推广工具页面
  - 活动选择器
  - 推广链接生成
  - 二维码展示和下载
- [ ] 3.4 实现推广数据统计页面
  - 数据卡片：订单数、收益、下级数
  - 奖励明细列表（分页）
  - 图表展示（可选）
- [ ] 3.5 实现下级列表页面
  - 一级下级分销商列表
  - 每个下级的基本信息

## 4. 前端管理后台实现
- [ ] 4.1 创建分销商申请审批页面 (DistributorApprovalView.tsx)
  - 待审批申请列表
  - 申请详情查看
  - 审批操作（批准/拒绝）
  - 设置级别
- [ ] 4.2 创建分销商管理页面 (DistributorManagementView.tsx)
  - 分销商列表（按品牌筛选）
  - 分销商详情查看
  - 级别调整
  - 状态管理（暂停/激活）
- [ ] 4.3 在品牌管理菜单中添加"分销商管理"入口
- [ ] 4.4 更新权限控制组件，支持 distributor 角色

## 5. 奖励系统改造
- [ ] 5.1 实现多级奖励计算逻辑
  - 订单支付成功时向上追溯最多3级
  - 按级别比例计算奖励
  - 只有 active 状态的分销商才获得奖励
- [ ] 5.2 更新奖励记录，关联分销商级别信息
- [ ] 5.3 确保奖励计算的幂等性和事务安全

## 6. 测试
- [ ] 6.1 单元测试
  - 分销商申请流程测试
  - 多级奖励计算测试
  - 权限控制测试
- [ ] 6.2 集成测试
  - 完整的分销商申请→审批→推广→奖励流程
  - 多级分销奖励分配验证
- [ ] 6.3 边界测试
  - 超过3级不分配奖励
  - 非分销商不获得奖励
  - 暂停状态不获得奖励

## 7. 文档和部署
- [ ] 7.1 更新API文档
- [ ] 7.2 编写分销商使用指南
- [ ] 7.3 准备数据库迁移脚本
- [ ] 7.4 部署验证

## 验收标准
- [ ] 顾客可以通过H5申请成为分销商
- [ ] 品牌管理员可以审批分销商申请
- [ ] 分销商可以生成专属推广链接和二维码
- [ ] 订单支付成功后，最多3级分销商自动获得奖励
- [ ] 分销商可以查看自己的推广数据和奖励明细
- [ ] 管理员可以管理分销商的级别和状态
- [ ] 数据隔离正确：分销商只能看到自己的数据
