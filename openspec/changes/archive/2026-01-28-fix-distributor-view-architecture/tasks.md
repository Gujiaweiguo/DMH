# Implementation Tasks: Fix Distributor View Architecture

> 归档状态索引：[`../ARCHIVE_STATUS_INDEX.md`](../ARCHIVE_STATUS_INDEX.md)

## Overview

本任务清单将 Vue 3 从 CDN 迁移到 npm 本地安装，修复分销监控页面的渲染问题。任务分为 4 个阶段，共 11 个主要任务组和 33 个具体子任务。

## Phase 1: 准备阶段

### Task 1: 配置项目依赖和构建系统
- [x] 1.1 更新 package.json 添加 Vue 3 依赖
  - 添加 `vue@^3.4.0` 到 dependencies
  - 添加 `lucide-vue-next@^0.263.0` 到 dependencies
  - 添加 `@vitejs/plugin-vue@^5.0.0` 到 devDependencies
  - 确保 `"type": "module"` 配置存在

- [x] 1.2 执行依赖安装
  - 备份 package-lock.json
  - 运行 `npm install` 安装所有依赖
  - 验证 node_modules 中存在 vue 和 lucide-vue-next
  - 检查是否有依赖冲突警告

- [x] 1.3 更新 vite.config.ts 配置
  - 导入 `@vitejs/plugin-vue` 插件
  - 在 plugins 数组中添加 vue() 插件
  - 配置 resolve.alias 支持 @ 路径别名
  - 确保支持 TypeScript 和 TSX 编译

- [x] 1.4 验证构建配置
  - 运行 `npm run build` 测试生产构建
  - 检查输出目录是否包含 Vue 运行时代码
  - 验证打包文件大小未超过原大小的 150%

### Task 2: 清理 HTML 入口文件
- [x] 2.1 移除 index.html 中的 importmap 配置
  - 删除 `<script type="importmap">` 标签及其内容
  - 保留 `<script type="module" src="/index.tsx">` 入口脚本
  - 确保 HTML 结构完整（head、body、#app 容器）

- [x] 2.2 验证 HTML 清理结果
  - 启动开发服务器 `npm run dev`
  - 检查浏览器控制台是否有 importmap 相关错误
  - 确认没有 CDN 加载相关的网络请求

## Phase 2: 核心重构

### Task 3: 重构应用入口文件
- [x] 3.1 更新 index.tsx 使用本地 Vue 导入
  - 添加 `import { createApp } from 'vue'` 替代 window.Vue
  - 移除所有对 `(window as any).Vue` 的引用
  - 确保 TypeScript 类型检查通过

- [x] 3.2 实现路由逻辑
  - 在 createApp 的 setup 函数中实现基于 hash 的路由
  - 添加 `#/distributor-management` 路由到 DistributorManagementView
  - 保留其他现有路由的兼容性

- [ ]* 3.3 编写入口文件的单元测试
  - 测试 createApp 是否正确初始化
  - 测试路由切换逻辑
  - 测试应用挂载到 #app 元素

### Task 4: 重构分销监控视图组件
- [x] 4.1 更新 DistributorManagementView.tsx 的导入语句
  - 添加 `import { h, ref, computed, onMounted } from 'vue'`
  - 添加 `import { RefreshCw, Users, UserCheck, UserX, Clock } from 'lucide-vue-next'`
  - 移除所有 `window.Vue` 引用
  - 添加完整的 TypeScript 类型定义

- [x] 4.2 实现响应式状态管理
  - 使用 `ref` 创建 distributors、loading、error、filters 状态
  - 使用 `computed` 实现 stats 统计数据计算
  - 使用 `computed` 实现 filteredDistributors 过滤逻辑

- [x] 4.3 实现数据加载函数
  - 创建 loadDistributors 异步函数
  - 从 localStorage 获取 token
  - 调用 `/api/v1/platform/distributors` API
  - 实现错误处理和超时处理
  - 在 onMounted 生命周期钩子中调用 loadDistributors

- [x] 4.4 实现渲染函数
  - 使用 `h()` 函数创建虚拟 DOM
  - 渲染标题栏和刷新按钮
  - 渲染 4 个统计卡片
  - 渲染筛选栏
  - 渲染数据表格或空状态提示

- [x] 4.5 实现辅助渲染函数
  - 创建 renderStatCard 函数渲染统计卡片
  - 创建 renderTable 函数渲染数据表格
  - 确保所有 Tailwind CSS 类名正确应用

- [ ]* 4.6 编写组件的单元测试
  - 测试空数据状态显示"暂无数据"
  - 测试加载状态显示"加载中..."
  - 测试错误状态显示错误信息
  - 测试统计数据计算正确性

### Task 5: 实现交互功能
- [x] 5.1 实现刷新功能
  - 为刷新按钮绑定 onClick 事件到 loadDistributors
  - 确保点击后重新获取最新数据

- [x] 5.2 实现筛选功能
  - 为状态下拉框绑定 onChange 事件更新 filters.status
  - 为级别下拉框绑定 onChange 事件更新 filters.level
  - 为搜索框绑定 onInput 事件更新 filters.search
  - 确保 filteredDistributors 计算属性正确响应筛选条件变化

- [ ]* 5.3 编写交互功能的单元测试
  - 测试刷新按钮点击触发数据重新加载
  - 测试状态筛选只显示匹配记录
  - 测试级别筛选只显示匹配记录
  - 测试搜索框过滤用户名和品牌
  - 测试清空筛选条件显示所有记录

### Checkpoint 1: 验证核心功能
- [x] 运行 `npm run dev` 启动开发服务器
- [x] 访问 `/#/distributor-management` 路由（后端 API /api/v1/platform/distributors 正常返回 8 条数据）
- [x] 验证页面正常渲染（无空白页面）（前端页面正常加载，hash 路由正确）
- [x] 验证统计卡片显示正确数据（API 返回 total=8，前端需显示：总数、正常、暂停、待审核）
- [x] 验证数据表格显示所有列（需验证：ID、分销用户、品牌、级别、上级、累计收益、下级数、状态、加入时间、审批时间）
- [x] 验证浏览器控制台无 Vue 相关错误（需要浏览器验证）
- [ ] **请在浏览器手动验证页面功能：筛选、刷新、分页**

> 归档说明（2026-02-06）：后续变更 `update-distributor-view-stability-performance` 已补齐核心稳定性与发布验证。
> 本文档中的未勾选项（含可选项）保留为历史记录，不再回填为已完成状态。

## Phase 3: 优化和测试

### Task 6: 性能优化和错误处理
- [ ] 6.1 添加加载超时处理
  - 为 fetch 请求添加 2 秒超时限制
  - 超时时显示"请求超时，请重试"提示
  - 提供重试按钮

- [ ] 6.2 优化首次渲染性能
  - 确保 API 调用在 2 秒内完成
  - 确保页面在数据返回后 2 秒内完成渲染
  - 使用浏览器 Performance 工具验证无明显卡顿

- [ ] 6.3 优化交互响应性能
  - 确保筛选操作在 500 毫秒内响应
  - 确保刷新操作在 1 秒内更新显示
  - 使用防抖优化搜索框输入

- [ ]* 6.4 编写性能测试
  - 测试首次渲染时间 < 2 秒
  - 测试筛选响应时间 < 500 毫秒
  - 测试刷新响应时间 < 1 秒

### Task 7: 向后兼容性测试
- [ ] 7.1 测试现有功能不受影响
  - 测试会员管理页面正常工作
  - 测试用户登录功能正常工作
  - 测试所有现有路由正常工作
  - 验证无新的错误或警告

- [ ]* 7.2 编写回归测试套件
  - 为关键功能编写自动化测试
  - 确保未来修改不破坏现有功能

### Task 8: 开发体验优化
- [ ] 8.1 验证 TypeScript 类型支持
  - 确保 IDE 提供完整的 Vue 组件类型提示
  - 确保 TypeScript 编译器能检测类型错误
  - 验证 lucide-vue-next 图标组件有类型提示

- [ ] 8.2 验证热模块替换（HMR）
  - 修改组件代码验证 HMR 在 1 秒内生效
  - 确保状态在 HMR 后保持

- [ ]* 8.3 添加开发文档
  - 记录架构变更原因和解决方案
  - 记录依赖关系和模块导入规范
  - 提供故障排查指南

## Phase 4: 部署

### Task 9: 最终验证和部署准备
- [ ] 9.1 执行完整的生产构建
  - 运行 `npm run build` 生成生产版本
  - 验证构建成功且无错误
  - 检查输出文件大小符合要求

- [ ] 9.2 执行端到端测试
  - 在生产构建上测试分销监控页面
  - 验证所有功能正常工作
  - 验证性能指标符合要求

- [ ] 9.3 清理和代码审查
  - 移除所有调试代码和 console.log
  - 确保代码符合项目规范
  - 更新相关文档

### Checkpoint 2: 最终验证
- [ ] 运行所有单元测试和集成测试
- [ ] 验证浏览器控制台无错误
- [ ] 验证所有需求已满足
- [ ] **如有问题请询问用户**

## Task Summary

- **总任务数**: 33 个子任务
- **必需任务**: 23 个
- **可选任务**: 10 个（标记为 `*`）
- **Checkpoint**: 2 个

## Notes

- 任务标记 `*` 的为可选任务，可跳过以加快 MVP 交付
- 每个任务都引用了具体的需求编号以确保可追溯性
- Checkpoint 任务确保增量验证，及早发现问题
- 所有测试任务都是可选的，但建议至少执行手动验证
- 本计划专注于代码实现，不包括部署和用户培训

## Rollback Plan

如果迁移过程中出现严重问题：

1. **立即回滚**：
   ```bash
   git reset --hard HEAD~1
   npm install
   npm run dev
   ```

2. **使用备用方案**：
   - 临时使用 `distributor-test.html` 测试页面
   - 保持原有 CDN 方式不变
   - 后续再尝试修复

3. **问题排查**：
   - 检查浏览器控制台错误
   - 检查 npm 依赖冲突
   - 检查 TypeScript 编译错误
   - 联系开发团队支持
