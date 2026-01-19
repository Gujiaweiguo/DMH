# Design: 分销商系统

## Context

DMH平台需要支持分销商（大C端）功能，这是一种高级顾客角色，具备分享推广并获得奖励的资格。分销商需要经过品牌管理员审批后才能生效。

系统当前已有：
- 普通顾客（participant）可参与活动、支付、获得优惠券
- 推荐人机制（referrer_id）已在订单中实现，用于单级推荐奖励

需要在此基础上扩展：
- 多级分销体系（1-3级）
- 分销商申请审批流程
- 分销商专属功能（推广工具、数据查看）

## Goals / Non-Goals

**Goals**:
- 实现分销商角色定义与权限管理
- 支持顾客申请成为分销商的审批流程
- 实现多级分销奖励计算（最多3级）
- 提供分销商推广工具（专属链接、二维码）
- 提供分销数据查看功能

**Non-Goals**:
- 分销商独立H5/小程序（复用现有H5）
- 复杂的团队管理系统
- 自动升级机制
- 分销商PC管理后台

## Decisions

### 1. 分销商角色设计

**Decision**: 分销商作为 participant 的扩展角色，而非完全独立的新角色

**Rationale**:
- 分销商本质上仍是顾客，只是具有额外的推广权限
- 复用现有的用户表和会员系统
- 降低系统复杂度

**Implementation**:
```
users 表保持不变
user_roles 表新增 distributor 角色关联
新增 distributors 表存储分销商专属信息（级别、上级等）
```

### 2. 多级分销层级限制

**Decision**: 最多支持3级分销

**Rationale**:
- 合规考虑：中国法规禁止超过3级的传销模式
- 业务平衡：足够激励推广，同时控制复杂度

**Implementation**:
```
一级分销：直接推荐人（referrer）
二级分销：推荐人的推荐人
三级分销：推荐人的推荐人的推荐人
```

### 3. 审批流程

**Decision**: 由品牌管理员审批，而非平台管理员

**Rationale**:
- 品牌对自己品牌的分销商质量负责
- 减轻平台管理员负担
- 品牌可以制定自己的分销商标准

### 4. 前端实现

**Decision**: 在现有H5基础上增加"分销中心"模块

**Rationale**:
- 用户无需安装额外应用
- 开发成本最低
- 与现有顾客功能无缝集成

## Architecture

### 数据模型

```
distributors 表：
- id: 分销商ID
- user_id: 关联用户ID
- brand_id: 关联品牌ID
- level: 分销级别(1/2/3)
- parent_id: 上级分销商ID
- status: 状态(pending/active/suspended)
- approved_by: 审批人ID
- approved_at: 审批时间
- total_earnings: 累计收益
- subordinates_count: 下级人数
```

### 奖励分配流程

```
订单支付成功 ->
  查询推荐链（最多3级）->
    计算各级奖励比例 ->
      创建奖励记录 ->
        更新分销商余额 ->
          发送通知
```

### 权限矩阵

| 功能 | participant | distributor | brand_admin | platform_admin |
|------|-------------|-------------|-------------|----------------|
| 参与活动 | ✓ | ✓ | ✓ | ✓ |
| 查看个人数据 | ✓ | ✓ | ✓ | ✓ |
| 申请成为分销商 | ✓ | - | - | - |
| 生成推广链接 | - | ✓ | - | - |
| 查看推广数据 | - | ✓ | - | - |
| 查看下级列表 | - | ✓(一级) | - | - |
| 审批分销商申请 | - | - | ✓ | ✓ |

## API Design

### 分销商申请
```
POST /api/v1/distributor/apply
Request: { brand_id, reason }
Response: { application_id, status }
```

### 审批申请
```
POST /api/v1/brands/:brandId/distributors/approve/:id
Request: { action: "approve|reject", level, reason }
```

### 生成推广链接
```
POST /api/v1/distributor/link/generate
Response: { link, qrcode_url }
```

### 查看推广数据
```
GET /api/v1/distributor/statistics
Response: {
  total_orders, total_earnings, subordinates_count,
  recent_orders[], recent_earnings[]
}
```

## Migration Plan

1. 创建 distributors 表
2. 修改 user_roles 表支持 distributor 角色
3. 部署后端API
4. 更新H5前端添加分销中心
5. 更新管理后台添加审批界面
6. 数据迁移：现有推荐人可选择升级为分销商

## Open Questions

1. ~~分销商级别是自动升级还是手动调整？~~ -> 手动调整
2. ~~分销商是否有有效期？~~ -> 无有效期，除非手动暂停
3. ~~是否支持分销商转让/更换上级？~~ -> 暂不支持
