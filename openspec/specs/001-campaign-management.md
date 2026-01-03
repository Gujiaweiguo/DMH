# Spec: 营销活动管理模块

**Module**: Campaign Management  
**Priority**: P0  
**Status**: ✅ Approved  
**Related Proposal**: [001-dmh-mvp-core-features](../changes/001-dmh-mvp-core-features.md)

---

## 📋 模块概述

营销活动管理模块是 DMH 的核心模块，提供活动的完整生命周期管理，包括创建、配置、发布、监控和下线。支持动态表单配置，实现"零代码"上线新活动。

---

## 🎯 功能需求

### 1. 活动列表页

#### 功能描述
展示所有营销活动的列表视图，支持分页、搜索和筛选。

#### 界面要求
- 表格展示：活动名称、状态、时间范围、创建时间
- 状态标签：进行中（绿色）、已结束（灰色）、已禁用（红色）
- 操作按钮：查看、编辑、禁用/启用、删除
- 顶部工具栏：创建活动按钮、搜索框、状态筛选器

#### 交互要求
- 搜索框支持活动名称模糊搜索
- 状态筛选支持多选
- 分页默认每页 20 条
- 点击行可查看活动详情

### 2. 活动创建/编辑页

#### 基础信息配置
- **活动名称**: 必填，2-100 字符
- **活动描述**: 可选，最多 500 字符
- **活动时间**: 必填，开始时间 < 结束时间
- **活动主图**: 可选，支持上传，推荐尺寸 750x400

#### 动态表单配置
支持的字段类型：
1. **文本框** (text)
   - 字段名称、占位符、是否必填
   - 最小/最大长度限制
2. **手机号** (phone)
   - 自动校验 11 位手机号格式
   - 必填字段
3. **下拉选择** (select)
   - 配置选项列表
   - 支持单选

表单配置界面：
- 左侧：字段类型选择区
- 中间：表单预览区（实时预览）
- 右侧：字段属性配置区

#### 奖励规则配置
- **奖励类型**: 固定金额（MVP 阶段）
- **奖励金额**: 必填，0.01-999.99 元
- **奖励条件**: 推荐用户完成支付

#### 支付参数配置
- **报名费用**: 必填，0 元或固定金额
- **微信支付商户号**: 自动关联系统配置

### 3. 活动详情页

#### 显示内容
- 活动基本信息
- 动态表单配置（JSON 格式展示）
- 奖励规则
- 统计数据：
  - 总参与人数
  - 总支付金额
  - 总奖励发放
  - 同步状态统计

#### 操作区
- 编辑按钮
- 禁用/启用按钮
- 复制活动链接按钮
- 查看订单列表按钮

---

## 🔌 API 接口定义

### 1. 创建活动
```
POST /api/v1/campaigns
Content-Type: application/json

Request Body:
{
  "name": "新年促销活动",
  "description": "新年大促，推荐有礼",
  "formFields": [
    {
      "type": "text",
      "name": "name",
      "label": "姓名",
      "required": true,
      "placeholder": "请输入姓名"
    },
    {
      "type": "phone",
      "name": "phone",
      "label": "手机号",
      "required": true
    },
    {
      "type": "select",
      "name": "course",
      "label": "意向课程",
      "required": true,
      "options": ["前端开发", "后端开发", "全栈开发"]
    }
  ],
  "rewardRule": 10.00,
  "startTime": "2025-01-01T00:00:00Z",
  "endTime": "2025-12-31T23:59:59Z"
}

Response:
{
  "id": 1,
  "name": "新年促销活动",
  "description": "新年大促，推荐有礼",
  "formFields": [...],
  "rewardRule": 10.00,
  "startTime": "2025-01-01T00:00:00Z",
  "endTime": "2025-12-31T23:59:59Z",
  "status": "active",
  "createdAt": "2025-01-01T10:00:00Z"
}
```

### 2. 获取活动列表
```
GET /api/v1/campaigns?page=1&pageSize=20&status=active&keyword=促销

Response:
{
  "total": 100,
  "campaigns": [
    {
      "id": 1,
      "name": "新年促销活动",
      "status": "active",
      "startTime": "2025-01-01T00:00:00Z",
      "endTime": "2025-12-31T23:59:59Z",
      "createdAt": "2025-01-01T10:00:00Z"
    }
  ]
}
```

### 3. 获取活动详情
```
GET /api/v1/campaigns/:id

Response:
{
  "id": 1,
  "name": "新年促销活动",
  "description": "新年大促，推荐有礼",
  "formFields": [...],
  "rewardRule": 10.00,
  "startTime": "2025-01-01T00:00:00Z",
  "endTime": "2025-12-31T23:59:59Z",
  "status": "active",
  "createdAt": "2025-01-01T10:00:00Z",
  "statistics": {
    "totalParticipants": 1000,
    "totalAmount": 99990.00,
    "totalReward": 10000.00,
    "syncSuccessRate": 0.98
  }
}
```

### 4. 更新活动
```
PUT /api/v1/campaigns/:id
Content-Type: application/json

Request Body: (同创建活动)

Response: (同创建活动响应)
```

### 5. 删除活动（软删除）
```
DELETE /api/v1/campaigns/:id

Response:
{
  "message": "活动已删除"
}
```

---

## 💾 数据存储

### campaigns 表
```sql
CREATE TABLE campaigns (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(200) NOT NULL COMMENT '活动名称',
    description TEXT COMMENT '活动描述',
    form_fields JSON COMMENT '动态表单配置',
    reward_rule DECIMAL(10,2) NOT NULL DEFAULT 0.00 COMMENT '奖励金额',
    start_time DATETIME NOT NULL COMMENT '开始时间',
    end_time DATETIME NOT NULL COMMENT '结束时间',
    status VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '状态',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL DEFAULT NULL COMMENT '软删除',
    INDEX idx_status (status),
    INDEX idx_start_time (start_time),
    INDEX idx_end_time (end_time),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### form_fields JSON Schema
```json
{
  "type": "array",
  "items": {
    "type": "object",
    "properties": {
      "type": {"type": "string", "enum": ["text", "phone", "select"]},
      "name": {"type": "string"},
      "label": {"type": "string"},
      "required": {"type": "boolean"},
      "placeholder": {"type": "string"},
      "options": {"type": "array", "items": {"type": "string"}},
      "minLength": {"type": "integer"},
      "maxLength": {"type": "integer"}
    },
    "required": ["type", "name", "label"]
  }
}
```

---

## 🎨 前端实现

### 技术栈
- Vue 3 Composition API
- Element Plus UI 组件库
- Axios HTTP 客户端
- VueRouter 路由管理

### 页面路由
```
/campaigns          - 活动列表页
/campaigns/create   - 创建活动页
/campaigns/:id/edit - 编辑活动页
/campaigns/:id      - 活动详情页
```

### 核心组件
```
CampaignList.vue        - 活动列表
CampaignForm.vue        - 活动表单（创建/编辑）
  ├─ BasicInfo.vue      - 基础信息
  ├─ FormBuilder.vue    - 表单构建器
  ├─ RewardConfig.vue   - 奖励配置
  └─ PaymentConfig.vue  - 支付配置
CampaignDetail.vue      - 活动详情
FormPreview.vue         - 表单预览组件
```

---

## 🔧 后端实现

### Go-Zero 项目结构
```
backend/api/
├── internal/
│   ├── handler/
│   │   └── campaign/
│   │       ├── createcampaignhandler.go
│   │       ├── getcampaignshandler.go
│   │       ├── getcampaignhandler.go
│   │       ├── updatecampaignhandler.go
│   │       └── deletecampaignhandler.go
│   ├── logic/
│   │   └── campaign/
│   │       ├── createcampaignlogic.go
│   │       ├── getcampaignslogic.go
│   │       ├── getcampaignlogic.go
│   │       ├── updatecampaignlogic.go
│   │       └── deletecampaignlogic.go
│   └── svc/
│       └── servicecontext.go
└── model/
    └── campaignmodel.go
```

### 核心业务逻辑

#### CreateCampaignLogic
```go
1. 验证请求参数
2. 校验时间范围（startTime < endTime）
3. 验证 formFields JSON 格式
4. 插入 campaigns 表
5. 返回创建结果
```

#### GetCampaignsLogic
```go
1. 解析查询参数（page, pageSize, status, keyword）
2. 构建查询条件
3. 执行分页查询
4. 统计总数
5. 返回列表和总数
```

---

## ✅ 验收标准

### 功能验收
- [ ] 可以成功创建活动
- [ ] 活动列表正确显示
- [ ] 搜索和筛选功能正常
- [ ] 可以编辑活动信息
- [ ] 可以禁用/启用活动
- [ ] 可以删除活动（软删除）
- [ ] 动态表单配置正确保存

### 性能验收
- [ ] 列表页加载时间 < 1 秒
- [ ] 创建活动响应时间 < 500ms
- [ ] 支持 1000+ 活动数据量

### 用户体验验收
- [ ] 表单预览实时更新
- [ ] 操作有明确的成功/失败提示
- [ ] 界面响应式适配

---

## 🧪 测试用例

### 单元测试
1. 创建活动 - 正常流程
2. 创建活动 - 参数校验失败
3. 创建活动 - 时间范围错误
4. 查询活动列表 - 无筛选条件
5. 查询活动列表 - 按状态筛选
6. 查询活动列表 - 关键词搜索
7. 更新活动 - 正常流程
8. 删除活动 - 软删除验证

### 集成测试
1. 完整的活动创建流程
2. 创建后立即查询验证
3. 更新后数据一致性验证
4. 删除后列表不显示验证

---

## 📝 开发清单

### 后端开发
- [ ] 编写 API 定义文件
- [ ] 生成 handler 和 logic 代码
- [ ] 实现 CampaignModel
- [ ] 实现创建活动逻辑
- [ ] 实现查询列表逻辑
- [ ] 实现查询详情逻辑
- [ ] 实现更新活动逻辑
- [ ] 实现删除活动逻辑
- [ ] 编写单元测试
- [ ] API 文档生成

### 前端开发
- [ ] 创建活动列表页面
- [ ] 创建活动表单页面
- [ ] 实现表单构建器组件
- [ ] 实现表单预览组件
- [ ] 创建活动详情页面
- [ ] 集成 API 接口
- [ ] 表单验证逻辑
- [ ] 错误处理
- [ ] 页面联调测试

---

## 🔗 相关文档
- [Proposal: DMH MVP 核心功能](../changes/001-dmh-mvp-core-features.md)
- [API 定义](../../backend/api/dmh.api)
- [数据库设计](../../backend/scripts/init.sql)
