# 会员管理前端路由集成指南

## 已创建的文件

### 服务层
- `services/memberApi.ts` - 会员管理 API 服务

### 视图层
- `views/MemberListView.tsx` - 会员列表页面
- `views/MemberDetailView.tsx` - 会员详情页面
- `views/MemberMergeView.tsx` - 会员合并页面
- `views/MemberExportView.tsx` - 导出管理页面

### 样式
- `styles/member.css` - 会员管理样式

## 路由配置

在主应用文件（`index.tsx` 或路由配置文件）中添加以下路由：

```typescript
import MemberListView from './views/MemberListView';
import MemberDetailView from './views/MemberDetailView';
import MemberMergeView from './views/MemberMergeView';
import MemberExportView from './views/MemberExportView';

// 在路由配置中添加
const routes = {
  '/members': MemberListView,
  '/members/:id': MemberDetailView,
  '/members/merge': MemberMergeView,
  '/members/export': MemberExportView,
  // ... 其他路由
};
```

## 菜单配置

在侧边栏菜单中添加会员管理入口：

```typescript
const menuItems = [
  // ... 其他菜单项
  {
    id: 'members',
    label: '会员管理',
    icon: 'Users',
    path: '#/members',
    permission: 'member:read',
    children: [
      {
        id: 'member-list',
        label: '会员列表',
        path: '#/members',
      },
      {
        id: 'member-export',
        label: '导出管理',
        path: '#/members/export',
      },
    ],
  },
];
```

## 样式引入

在主 HTML 文件或主应用文件中引入样式：

```html
<link rel="stylesheet" href="/styles/member.css">
```

或在 TypeScript/JavaScript 中：

```typescript
import './styles/member.css';
```

## 权限控制

会员管理功能需要以下权限：

- `member:read` - 查看会员列表和详情
- `member:merge` - 合并会员（仅平台管理员）
- `member:export:create` - 创建导出申请
- `member:export:approve` - 审批导出申请（仅平台管理员）
- `member:tag:create` - 创建标签（仅平台管理员）
- `member:tag:manage` - 管理会员标签（仅平台管理员）

## 使用示例

### 1. 从其他页面跳转到会员列表

```typescript
// 跳转到会员列表
window.location.hash = '#/members';

// 跳转到会员详情
window.location.hash = `#/members/${memberId}`;

// 跳转到会员合并（带参数）
window.location.hash = `#/members/merge?source=${sourceId}&target=${targetId}`;

// 跳转到导出管理
window.location.hash = '#/members/export';
```

### 2. 在订单详情页添加"查看会员"按钮

```typescript
// 在订单详情中
if (order.memberId) {
  return h('button', {
    onClick: () => {
      window.location.hash = `#/members/${order.memberId}`;
    },
    class: 'btn btn-link'
  }, '查看会员');
}
```

### 3. 在活动统计中添加会员维度

```typescript
// 在活动详情中显示会员统计
const memberStats = {
  totalMembers: 150,
  newMembers: 30,
  returningMembers: 120,
};

// 点击查看会员列表（筛选该活动的会员）
const viewCampaignMembers = (campaignId: number) => {
  window.location.hash = `#/members?campaignId=${campaignId}`;
};
```

## API 调用示例

### 查询会员列表

```typescript
import { memberApi } from './services/memberApi';

// 基础查询
const members = await memberApi.getMembers({
  page: 1,
  pageSize: 20,
});

// 带筛选条件
const filteredMembers = await memberApi.getMembers({
  page: 1,
  pageSize: 20,
  keyword: '张三',
  status: 'active',
  brandId: 1,
  startDate: '2024-01-01',
  endDate: '2024-12-31',
});
```

### 查询会员详情

```typescript
const member = await memberApi.getMember(memberId);
console.log(member.totalOrders, member.totalPayment);
```

### 会员合并

```typescript
// 预览合并
const preview = await memberApi.previewMerge({
  sourceMemberId: 2,
  targetMemberId: 1,
  reason: '重复会员',
});

if (preview.canMerge) {
  // 执行合并
  await memberApi.mergeMember({
    sourceMemberId: 2,
    targetMemberId: 1,
    reason: '重复会员',
  });
}
```

### 导出申请

```typescript
// 创建导出申请
await memberApi.createExportRequest({
  brandId: 1,
  reason: '市场分析需要',
  filters: JSON.stringify({ status: 'active' }),
});

// 审批导出
await memberApi.approveExportRequest(requestId, {
  approve: true,
});
```

## 注意事项

1. **权限检查**：在显示菜单和按钮前检查用户权限
2. **错误处理**：所有 API 调用都应该有 try-catch 错误处理
3. **加载状态**：在数据加载时显示 loading 状态
4. **数据刷新**：操作成功后及时刷新列表数据
5. **用户反馈**：操作成功或失败后给用户明确的提示

## 测试清单

- [ ] 会员列表加载正常
- [ ] 分页功能正常
- [ ] 筛选功能正常
- [ ] 会员详情显示完整
- [ ] 标签添加功能正常
- [ ] 会员合并预览正常
- [ ] 会员合并执行成功
- [ ] 导出申请创建成功
- [ ] 导出审批流程正常
- [ ] 权限控制生效
- [ ] 错误提示友好
- [ ] 响应式布局正常

## 下一步

1. 将路由配置集成到主应用
2. 添加菜单项
3. 配置权限
4. 测试所有功能
5. 优化用户体验
