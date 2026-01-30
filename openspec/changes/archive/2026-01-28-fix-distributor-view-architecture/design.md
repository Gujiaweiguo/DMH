# Design: Fix Distributor View Architecture

## Overview

本设计文档描述了如何将 DMH 平台前端管理后台的 Vue 3 运行时从 CDN 迁移到 npm 本地安装，以修复分销监控页面的渲染问题。

## Problem Statement

### 当前架构问题

```
┌─────────────────────────────────────────┐
│         index.html                      │
│  ┌───────────────────────────────────┐  │
│  │  <script type="importmap">        │  │
│  │  {                                │  │
│  │    "vue": "https://esm.sh/vue"   │  │
│  │  }                                │  │
│  └───────────────────────────────────┘  │
└─────────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────┐
│    DistributorManagementView.tsx        │
│  ┌───────────────────────────────────┐  │
│  │  const { h } = window.Vue || {}   │  │
│  │  // h 可能为 undefined             │  │
│  └───────────────────────────────────┘  │
└─────────────────────────────────────────┘
                  │
                  ▼
         ❌ 渲染失败（空白页面）
```

**核心问题**：
1. Vue 从 CDN 加载，模块初始化时序不确定
2. 组件依赖 `window.Vue` 全局变量，可能未定义
3. CSP 策略可能阻止 CDN 脚本执行
4. 缺少类型支持，开发体验差

## Solution Architecture

### 目标架构

```
┌─────────────────────────────────────────┐
│         package.json                    │
│  ┌───────────────────────────────────┐  │
│  │  "dependencies": {                │  │
│  │    "vue": "^3.4.0",              │  │
│  │    "lucide-vue-next": "^0.263.0" │  │
│  │  }                                │  │
│  └───────────────────────────────────┘  │
└─────────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────┐
│         Vite Build System               │
│  - 解析 node_modules                    │
│  - 打包 Vue 运行时                      │
│  - 生成优化的 bundle                    │
└─────────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────┐
│    DistributorManagementView.tsx        │
│  ┌───────────────────────────────────┐  │
│  │  import { h } from 'vue'          │  │
│  │  // h 始终可用且类型安全           │  │
│  └───────────────────────────────────┘  │
└─────────────────────────────────────────┘
                  │
                  ▼
         ✅ 正常渲染（显示数据）
```

## Key Changes

### 1. 依赖管理

**Before (CDN)**:
```html
<script type="importmap">
{
  "imports": {
    "vue": "https://esm.sh/vue@3.3.4"
  }
}
</script>
```

**After (NPM)**:
```json
{
  "dependencies": {
    "vue": "^3.4.0",
    "lucide-vue-next": "^0.263.0"
  },
  "devDependencies": {
    "@vitejs/plugin-vue": "^5.0.0"
  }
}
```

### 2. 模块导入

**Before (Global)**:
```typescript
const { h } = (window as any).Vue || { h: () => null };
```

**After (ES6 Module)**:
```typescript
import { h, ref, computed, onMounted } from 'vue';
import { RefreshCw, Users } from 'lucide-vue-next';
```

### 3. 组件结构

**Before (Problematic)**:
```typescript
export const renderDistributorManagementView = (viewModel: any) => {
  const { h } = (window as any).Vue || { h: () => null };
  // 可能渲染失败
};
```

**After (Robust)**:
```typescript
import { h, ref, computed, onMounted } from 'vue';

export function DistributorManagementView() {
  // 响应式状态
  const distributors = ref<Distributor[]>([]);
  const loading = ref(false);
  const error = ref<string | null>(null);
  
  // 计算属性
  const stats = computed(() => ({
    total: distributors.value.length,
    normal: distributors.value.filter(d => d.status === 'normal').length,
    // ...
  }));
  
  // 数据加载
  const loadDistributors = async () => {
    loading.value = true;
    try {
      const response = await fetch('/api/v1/platform/distributors');
      distributors.value = await response.json();
    } catch (err) {
      error.value = err.message;
    } finally {
      loading.value = false;
    }
  };
  
  // 生命周期
  onMounted(() => {
    loadDistributors();
  });
  
  // 渲染函数
  return () => h('div', { class: 'p-6' }, [
    // 组件内容
  ]);
}
```

### 4. Vite 配置

**Before (Incomplete)**:
```typescript
export default defineConfig({
  server: { port: 3000 },
  resolve: {
    alias: { '@': path.resolve(__dirname, '.') }
  }
});
```

**After (Complete)**:
```typescript
import vue from '@vitejs/plugin-vue';

export default defineConfig({
  plugins: [vue()],  // 添加 Vue 插件
  server: { port: 3000 },
  resolve: {
    alias: {
      '@': path.resolve(__dirname, '.'),
      'vue': 'vue/dist/vue.esm-bundler.js'  // 确保使用正确的 Vue 构建版本
    }
  }
});
```

## Data Flow

### 组件数据流

```
┌─────────────────────────────────────────┐
│  DistributorManagementView              │
│                                         │
│  ┌─────────────────────────────────┐   │
│  │  State (ref)                    │   │
│  │  - distributors: []             │   │
│  │  - loading: false               │   │
│  │  - error: null                  │   │
│  │  - filters: { ... }             │   │
│  └─────────────────────────────────┘   │
│                │                        │
│                ▼                        │
│  ┌─────────────────────────────────┐   │
│  │  Computed Properties            │   │
│  │  - stats (统计数据)              │   │
│  │  - filteredDistributors (过滤)  │   │
│  └─────────────────────────────────┘   │
│                │                        │
│                ▼                        │
│  ┌─────────────────────────────────┐   │
│  │  Render Function (h)            │   │
│  │  - 统计卡片                      │   │
│  │  - 筛选栏                        │   │
│  │  - 数据表格                      │   │
│  └─────────────────────────────────┘   │
└─────────────────────────────────────────┘
```

### API 交互流程

```
Component                API                 Backend
    │                    │                     │
    │  loadDistributors  │                     │
    ├───────────────────>│                     │
    │                    │  GET /api/v1/...    │
    │                    ├────────────────────>│
    │                    │                     │
    │                    │  Response (JSON)    │
    │                    │<────────────────────┤
    │  Update State      │                     │
    │<───────────────────┤                     │
    │                    │                     │
    │  Re-render         │                     │
    │                    │                     │
```

## Component Interface

### TypeScript 类型定义

```typescript
// 分销商数据类型
interface Distributor {
  id: number;
  user_id: number;
  username: string;
  brand_name: string;
  level: number;
  parent_id: number | null;
  total_earnings: string;
  subordinate_count: number;
  status: 'normal' | 'paused' | 'pending';
  created_at: string;
  updated_at: string;
}

// 统计数据类型
interface DistributorStats {
  total: number;
  normal: number;
  paused: number;
  pending: number;
}

// 筛选状态类型
interface FilterState {
  status: string;
  level: string;
  search: string;
}

// 组件 Props 类型
interface DistributorManagementProps {
  readOnly?: boolean;
  isPlatformAdmin?: boolean;
}
```

### 组件 API

```typescript
// 导出的组件函数
export function DistributorManagementView(
  props?: DistributorManagementProps
): () => VNode;

// 内部函数
function loadDistributors(): Promise<void>;
function renderStatCard(title: string, value: number, icon: any, bgColor: string): VNode;
function renderTable(data: Distributor[]): VNode;
```

## Performance Considerations

### 优化策略

1. **懒加载**：
   - 使用 `defineAsyncComponent` 懒加载大型组件
   - 路由级别的代码分割

2. **计算属性缓存**：
   - 使用 `computed` 缓存计算结果
   - 避免在渲染函数中进行复杂计算

3. **防抖和节流**：
   - 搜索框输入使用防抖（300ms）
   - 滚动事件使用节流

4. **虚拟滚动**：
   - 如果数据量大（>100条），考虑使用虚拟滚动

### 性能指标

| 指标 | 目标 | 测量方法 |
|------|------|----------|
| 首次渲染 | < 2s | Performance API |
| 筛选响应 | < 500ms | Performance API |
| 刷新响应 | < 1s | Performance API |
| 包体积增长 | < 150% | Build output |

## Error Handling

### 错误类型和处理

```typescript
// 网络错误
try {
  const response = await fetch('/api/v1/platform/distributors');
  if (!response.ok) {
    throw new Error(`HTTP ${response.status}: ${response.statusText}`);
  }
} catch (err) {
  error.value = '网络错误，请稍后重试';
  console.error('加载失败:', err);
}

// 超时处理
const controller = new AbortController();
const timeoutId = setTimeout(() => controller.abort(), 2000);

try {
  const response = await fetch('/api/v1/platform/distributors', {
    signal: controller.signal
  });
} catch (err) {
  if (err.name === 'AbortError') {
    error.value = '请求超时，请重试';
  }
} finally {
  clearTimeout(timeoutId);
}

// 数据验证错误
if (!Array.isArray(data.distributors)) {
  error.value = '数据格式错误';
  console.error('Invalid data format:', data);
}
```

## Testing Strategy

### 单元测试

```typescript
// 测试组件初始化
test('should initialize with empty state', () => {
  const component = DistributorManagementView();
  expect(component.distributors.value).toEqual([]);
  expect(component.loading.value).toBe(false);
});

// 测试数据加载
test('should load distributors on mount', async () => {
  const component = DistributorManagementView();
  await component.loadDistributors();
  expect(component.distributors.value.length).toBeGreaterThan(0);
});

// 测试筛选功能
test('should filter by status', () => {
  const component = DistributorManagementView();
  component.filters.value.status = 'normal';
  expect(component.filteredDistributors.value.every(d => d.status === 'normal')).toBe(true);
});
```

### 集成测试

```typescript
// 测试完整流程
test('should display distributor data', async () => {
  render(<DistributorManagementView />);
  
  // 等待数据加载
  await waitFor(() => {
    expect(screen.getByText('总分销商')).toBeInTheDocument();
  });
  
  // 验证统计卡片
  expect(screen.getByText('8')).toBeInTheDocument();
  
  // 验证数据表格
  expect(screen.getAllByRole('row').length).toBeGreaterThan(1);
});
```

## Migration Path

### 迁移步骤

1. **Phase 1: 准备**
   - 备份当前代码
   - 更新依赖配置
   - 验证构建系统

2. **Phase 2: 重构**
   - 清理 HTML
   - 重构组件
   - 实现功能

3. **Phase 3: 测试**
   - 功能测试
   - 性能测试
   - 兼容性测试

4. **Phase 4: 部署**
   - 生产构建
   - 灰度发布
   - 全量发布

### 回滚策略

如果出现问题：
1. Git 回滚到上一个稳定版本
2. 使用测试页面作为临时方案
3. 分析问题并修复
4. 重新部署

## Security Considerations

### CSP 配置

移除 importmap 后，可以使用更严格的 CSP：

```html
<meta http-equiv="Content-Security-Policy" 
      content="default-src 'self'; 
               script-src 'self'; 
               style-src 'self' 'unsafe-inline'; 
               img-src 'self' data: https:;">
```

### 依赖安全

```bash
# 检查依赖漏洞
npm audit

# 修复漏洞
npm audit fix
```

## References

- [Vue 3 文档](https://vuejs.org/)
- [Vite 文档](https://vitejs.dev/)
- [TypeScript 文档](https://www.typescriptlang.org/)
- [Lucide Icons](https://lucide.dev/)

## Appendix

### 完整的组件代码示例

详细的组件实现请参考：
- `.kiro/specs/fix-distributor-view-architecture/design.md` - 完整设计文档
- `frontend-admin/distributor-test.html` - 参考实现

### 相关文件

- `openspec/changes/fix-distributor-view-architecture/proposal.md` - 变更提案
- `openspec/changes/fix-distributor-view-architecture/tasks.md` - 任务清单
- `.kiro/specs/fix-distributor-view-architecture/requirements.md` - 需求文档

