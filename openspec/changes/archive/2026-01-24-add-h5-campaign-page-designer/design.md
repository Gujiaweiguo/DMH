# Design: H5 Campaign Page Designer

## 状态: ✅ 已完成实现

## Context
需要在H5端实现与Admin后台相同的活动页面设计功能，但需要适配移动端的交互方式和屏幕尺寸。

## Goals
- ✅ 提供完整的页面设计功能
- ✅ 移动端友好的交互体验
- ✅ 与Admin后台功能保持一致
- ✅ 支持实时预览

## Non-Goals
- 不实现复杂的拖拽动画（移动端性能考虑）
- 不支持自定义组件开发（使用预设组件）
- 不实现多人协作编辑

## 实现架构

### 1. 组件架构 (已实现)
```
CampaignPageDesigner.vue (主容器)
├── 组件库面板 (van-popup 底部弹出)
├── 设计画布 (中间区域)
│   └── 各类展示组件
├── 编辑面板 (van-popup 底部弹出)
│   └── 各类编辑器组件
└── 预览模式 (全屏预览)
```

### 2. 数据结构 (已实现)
```typescript
interface PageConfig {
  id: number
  campaignId: number
  components: ComponentConfig[]
  theme: ThemeConfig
  createdAt: string
  updatedAt: string
}

interface ComponentConfig {
  id: string
  type: 'poster' | 'title' | 'time' | 'location' | 'highlight' | 'detail' | 'divider' | 'button'
  data: Record<string, any>
}

interface ThemeConfig {
  primaryColor: string
  backgroundColor: string
}
```

### 3. 组件清单

#### 展示组件 (frontend-h5/src/components/designer/)
| 组件 | 文件 | 功能 |
|------|------|------|
| 海报 | PosterComponent.vue | 显示活动海报图片 |
| 标题 | TitleComponent.vue | 显示活动标题和副标题 |
| 时间 | TimeComponent.vue | 显示活动时间 |
| 地点 | LocationComponent.vue | 显示活动地点 |
| 亮点 | HighlightComponent.vue | 显示活动亮点列表 |
| 详情 | DetailComponent.vue | 显示活动详情富文本 |
| 分割线 | DividerComponent.vue | 显示分割线 |
| 按钮 | ButtonComponent.vue | 显示报名按钮 |

#### 编辑器组件 (frontend-h5/src/components/designer/editors/)
| 编辑器 | 文件 | 功能 |
|--------|------|------|
| 海报编辑 | PosterEditor.vue | 编辑海报图片URL |
| 标题编辑 | TitleEditor.vue | 编辑标题和副标题 |
| 时间编辑 | TimeEditor.vue | 编辑开始/结束时间 |
| 地点编辑 | LocationEditor.vue | 编辑地点名称和地址 |
| 亮点编辑 | HighlightEditor.vue | 编辑亮点列表 |
| 详情编辑 | DetailEditor.vue | 编辑详情内容 |
| 分割线编辑 | DividerEditor.vue | 编辑分割线样式 |
| 按钮编辑 | ButtonEditor.vue | 编辑按钮文字和颜色 |

### 4. API设计 (已实现)

#### 保存页面配置
```
POST /api/v1/campaign/page-config/:id
Request:
{
  "components": [...],
  "theme": {...}
}
Response:
{
  "id": 1,
  "campaignId": 1,
  "components": [...],
  "theme": {...},
  "createdAt": "...",
  "updatedAt": "..."
}
```

#### 获取页面配置
```
GET /api/v1/campaign/page-config/:id
Response:
{
  "id": 1,
  "campaignId": 1,
  "components": [...],
  "theme": {...},
  "createdAt": "...",
  "updatedAt": "..."
}
```

### 5. 数据库设计 (已实现)

```sql
CREATE TABLE page_configs (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  campaign_id BIGINT NOT NULL UNIQUE,
  components JSON,
  theme JSON,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_campaign_id (campaign_id)
);
```

### 6. 移动端适配策略 (已实现)

- ✅ 使用全屏布局
- ✅ 通过底部弹出层打开组件库 (van-popup)
- ✅ 使用底部弹出层编辑组件 (van-popup)
- ✅ 预览功能使用全屏模式
- ✅ 点击添加组件（不使用拖拽）
- ✅ 使用上下箭头调整组件顺序
- ✅ 使用Vant的Toast提供操作反馈

## 文件清单

### 前端文件
```
frontend-h5/src/
├── views/brand/
│   ├── CampaignPageDesigner.vue  # 主设计器页面
│   └── CampaignEditorVant.vue    # 活动编辑器(跳转入口)
├── components/designer/
│   ├── PosterComponent.vue
│   ├── TitleComponent.vue
│   ├── TimeComponent.vue
│   ├── LocationComponent.vue
│   ├── HighlightComponent.vue
│   ├── DetailComponent.vue
│   ├── DividerComponent.vue
│   ├── ButtonComponent.vue
│   └── editors/
│       ├── PosterEditor.vue
│       ├── TitleEditor.vue
│       ├── TimeEditor.vue
│       ├── LocationEditor.vue
│       ├── HighlightEditor.vue
│       ├── DetailEditor.vue
│       ├── DividerEditor.vue
│       └── ButtonEditor.vue
├── services/
│   └── brandApi.js               # API调用方法
├── router/
│   └── index.js                  # 路由配置
└── main.js                       # Vant组件注册
```

### 后端文件
```
backend/
├── api/
│   ├── dmh.api                   # API定义
│   └── internal/logic/campaign/
│       ├── save_page_config_logic.go
│       └── get_page_config_logic.go
├── model/
│   └── page_config.go            # 数据模型
└── scripts/
    └── create_page_configs_table.sql
```

## 访问路径
- 页面设计器: `/brand/campaigns/:id/page-design`
- 从活动编辑器点击"页面设计"按钮进入

## 技术栈
- Vue 3 + Composition API
- Vant 4 (DatePicker, TimePicker, RadioGroup, Radio 等)
- go-zero 后端框架
- MySQL 8 数据库

---

**实际开发时间**: 约3天
**状态**: ✅ 已完成
