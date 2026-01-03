# Implementation Tasks

## 1. 组件库设计
- [x] 1.1 创建组件库侧边栏组件
- [x] 1.2 实现活动海报组件 (PosterComponent.vue)
- [x] 1.3 实现活动标题组件 (TitleComponent.vue)
- [x] 1.4 实现活动时间组件 (TimeComponent.vue)
- [x] 1.5 实现活动地点组件 (LocationComponent.vue)
- [x] 1.6 实现活动亮点组件 (HighlightComponent.vue)
- [x] 1.7 实现活动详情组件 (DetailComponent.vue)
- [x] 1.8 实现分割线组件 (DividerComponent.vue)
- [x] 1.9 实现报名按钮组件 (ButtonComponent.vue)

## 2. 页面设计器核心功能
- [x] 2.1 创建主设计器页面组件 (CampaignPageDesigner.vue)
- [x] 2.2 实现组件点击添加功能
- [x] 2.3 实现组件编辑功能 (底部弹出编辑面板)
- [x] 2.4 实现组件删除功能
- [x] 2.5 实现组件排序功能 (上移/下移)
- [x] 2.6 实现页面配置保存
- [x] 2.7 实现页面配置加载

## 3. 移动端适配
- [x] 3.1 使用Vant组件优化移动端交互
- [x] 3.2 实现底部弹出式编辑面板 (van-popup)
- [x] 3.3 实现触摸友好的操作按钮
- [x] 3.4 优化小屏幕显示效果

## 4. 预览功能
- [x] 4.1 实现实时预览组件
- [x] 4.2 实现预览/编辑模式切换
- [x] 4.3 实现移动端预览样式
- [x] 4.4 实现预览数据绑定

## 5. 数据持久化
- [x] 5.1 设计页面配置数据结构 (components + theme)
- [x] 5.2 实现配置保存API (POST /api/v1/campaign/page-config/:id)
- [x] 5.3 实现配置加载API (GET /api/v1/campaign/page-config/:id)
- [x] 5.4 创建 page_configs 数据库表
- [x] 5.5 实现后端 Logic 层 (save_page_config_logic.go, get_page_config_logic.go)

## 6. 路由和导航
- [x] 6.1 添加页面设计器路由 (/brand/campaigns/:id/page-design)
- [x] 6.2 从活动编辑器跳转到页面设计器 (CampaignEditorVant.vue)
- [x] 6.3 实现返回和保存导航逻辑

## 7. 编辑器组件
- [x] 7.1 PosterEditor.vue - 海报编辑器
- [x] 7.2 TitleEditor.vue - 标题编辑器
- [x] 7.3 TimeEditor.vue - 时间编辑器
- [x] 7.4 LocationEditor.vue - 地点编辑器
- [x] 7.5 HighlightEditor.vue - 亮点编辑器
- [x] 7.6 DetailEditor.vue - 详情编辑器
- [x] 7.7 DividerEditor.vue - 分割线编辑器
- [x] 7.8 ButtonEditor.vue - 按钮编辑器

## 完成状态
✅ 所有核心功能已实现并测试通过

## 实现文件清单
### 前端 (frontend-h5)
- `src/views/brand/CampaignPageDesigner.vue` - 主设计器页面
- `src/views/brand/CampaignEditorVant.vue` - 活动编辑器(添加跳转按钮)
- `src/components/designer/*.vue` - 8个展示组件
- `src/components/designer/editors/*.vue` - 8个编辑器组件
- `src/services/brandApi.js` - API调用方法
- `src/router/index.js` - 路由配置
- `src/main.js` - Vant组件注册

### 后端 (backend)
- `api/dmh.api` - API定义
- `api/internal/logic/campaign/save_page_config_logic.go` - 保存逻辑
- `api/internal/logic/campaign/get_page_config_logic.go` - 获取逻辑
- `model/page_config.go` - 数据模型
- `scripts/create_page_configs_table.sql` - 数据库表脚本
