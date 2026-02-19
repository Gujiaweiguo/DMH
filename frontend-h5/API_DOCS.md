# 前端开发API文档

## 📅 Date: 2026-01-28

## 🎯 Purpose
为前端开发提供完整的API参考和开发指南。

---

## 📋 API端点列表

### 海报生成相关

#### 1. 生成活动海报
```
POST /api/v1/campaigns/:id/poster
Request Headers:
  Authorization: Bearer {token}
  Content-Type: application/json

Request Body:
  {
    "distributorId": number  // 分销商ID（可选）
  }

Response (200 OK):
  {
    "code": 0,
    "msg": "success",
    "data": {
      "posterUrl": "https://...",
      "qrcodeUrl": "https://...",
      "createdAt": "2026-01-28T10:00:00Z"
    }
  }

Response (400 Bad Request):
  {
    "code": 400,
    "msg": "活动不存在",
    "data": null
  }
```

#### 2. 获取海报文件
```
GET /api/v1/posters/:filename
Response (200 OK):
  {
    "code": 0,
    "msg": "success",
    "data": {
      "filename": "poster-123.png",
      "url": "https://.../poster-123.png"
    }
  }
```

---

### 订单核销相关

#### 1. 扫码获取订单信息
```
GET /api/v1/orders/scan/:code
Request Headers:
  Authorization: Bearer {token}

Request Params:
  code: string  // 核销码（二维码中的order_code）

Response (200 OK):
  {
    "code": 0,
    "msg": "success",
    "data": {
      "orderId": 123,
      "orderCode": "ABC123",
      "userId": 456,
      "userName": "张三",
      "userPhone": "138****8888",
      "campaignId": 789,
      "campaignName": "活动名称",
      "paymentStatus": "paid",  // paid/unpaid/refunded
      "verifyStatus": "unverified",  // unverified/verified
      "paymentAmount": 99.00,
      "paymentTime": "2026-01-28T10:00:00Z"
    }
  }

Response (404 Not Found):
  {
    "code": 404,
    "msg": "订单不存在或核销码无效",
    "data": null
  }
```

#### 2. 确认核销订单
```
POST /api/v1/orders/:id/verify
Request Headers:
  Authorization: Bearer {token}
  Content-Type: application/json

Request Body:
  {
    "notes": string  // 核销备注（可选）
  }

Response (200 OK):
  {
    "code": 0,
    "msg": "核销成功",
    "data": {
      "orderId": 123,
      "verifyStatus": "verified",
      "verifiedBy": "管理员姓名",
      "verifiedAt": "2026-01-28T10:00:00Z",
      "notes": "用户现场核销"
    }
  }

Response (400 Bad Request):
  {
    "code": 400,
    "msg": "订单已核销或状态不允许核销",
    "data": null
  }
```

#### 3. 取消核销订单
```
POST /api/v1/orders/:id/unverify
Request Headers:
  Authorization: Bearer {token}
  Content-Type: application/json

Request Body:
  {
    "reason": string  // 取消原因（必填）
  }

Response (200 OK):
  {
    "code": 0,
    "msg": "取消核销成功",
    "data": {
      "orderId": 123,
      "verifyStatus": "unverified",
      "cancelledBy": "管理员姓名",
      "cancelledAt": "2026-01-28T10:00:00Z",
      "reason": "误操作"
    }
  }
```

---

### 活动表单字段相关

#### FormField 数据结构
```typescript
interface FormField {
  type: string;           // 类型: text, email, phone, address, textarea, select, checkbox
  name: string;           // 字段名（用于数据存储）
  label: string;          // 显示标签
  required: boolean;       // 是否必填
  placeholder?: string;     // 占位提示
  options?: FormFieldOption[];  // 选项（select/checkbox类型）
}

interface FormFieldOption {
  label: string;          // 选项显示文本
  value: string;          // 选项值
}
```

---

## 🔧 开发环境配置

### 环境变量

后端服务：
```bash
export API_BASE_URL="http://localhost:8889"
export API_TOKEN="{your-test-token}"
```

前端开发：
```bash
# H5开发服务器
export H5_DEV_URL="http://localhost:3100"
```

### 后端服务启动

```bash
# 方式1：直接运行
cd /opt/code/dmh/backend/api
go run dmh.go -f etc/dmh-api.yaml

# 方式2：编译后运行
cd /opt/code/dmh/backend
go build -o dmh-api ./api/dmh.go
./dmh-api -f api/etc/dmh-api.yaml
```

### 前端开发服务器启动

```bash
cd /opt/code/dmh/frontend-h5
npm run dev

# 访问地址
# http://localhost:3100
```

---

## 📚 组件开发指南

### 1. PosterGenerator.vue 开发指南

#### 功能需求
- 支持活动专属海报生成
- 支持通用分销商海报生成
- 海报模板选择和预览
- 海报缩放和旋转
- 下载和分享功能

#### 关键技术点
- 使用html2canvas生成海报图片
- 使用qrcode生成二维码
- 使用Vue3 Composition API
- 使用Vant UI组件库

#### 开发步骤
1. 创建页面框架和路由
2. 实现模板选择器（Grid布局）
3. 实现海报预览（手势缩放、旋转）
4. 调用生成海报API
5. 实现下载功能（长按保存）
6. 实现分享功能（微信分享、复制链接）
7. 添加loading状态和错误处理

#### 依赖库
- `html2canvas@^1.4.1` - 海报图片生成
- `qrcode@^1.5.4` - 二维码生成
- `vue@^3.4.0` - Vue框架
- `vant@^4.8.0` - UI组件库

---

### 2. OrderVerification.vue 开发指南

#### 功能需求
- 扫描二维码获取订单信息
- 显示订单详情（用户信息、支付状态）
- 确认核销订单
- 取消核销订单
- 权限检查（仅品牌管理员）

#### 关键技术点
- 使用html5-qrcode扫描二维码
- 调用订单核销相关API
- 使用Vant Dialog组件
- 权限验证

#### 开发步骤
1. 创建页面框架和路由
2. 集成二维码扫描组件
3. 实现订单详情展示
4. 实现确认核销功能（带确认对话框）
5. 实现取消核销功能（带确认对话框）
6. 添加权限检查
7. 添加路由配置

#### 依赖库
- `html5-qrcode@^2.3.8` - 二维码扫描
- `vue@^3.4.0` - Vue框架
- `vant@^4.8.0` - UI组件库

---

## 🧪 测试指南

### 本地测试命令

```bash
# 1. 启动后端服务
cd /opt/code/dmh/backend/api
go run dmh.go -f etc/dmh-api.yaml

# 2. 启动前端开发服务器
cd /opt/code/dmh/frontend-h5
npm run dev

# 3. 测试海报生成
curl -X POST http://localhost:8889/api/v1/campaigns/1/poster \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"distributorId": 123}'

# 4. 测试订单扫码
curl http://localhost:8889/api/v1/orders/scan/ABC123 \
  -H "Authorization: Bearer YOUR_TOKEN"

# 5. 测试核销订单
curl -X POST http://localhost:8889/api/v1/orders/1/verify \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"notes": "现场核销"}'
```

### 端到端测试

```bash
# 1. 在手机上访问H5
# http://YOUR_IP:3100

# 2. 测试二维码扫描（使用实际手机摄像头）
# 3. 测试海报生成和下载
# 4. 测试订单核销流程

# 5. 跨浏览器测试
# Chrome, Safari, Firefox
```

---

## 📊 常见错误码

| 错误码 | 说明 | 解决方案 |
|--------|------|----------|
| 400 | 请求参数错误 | 检查请求体格式 |
| 401 | 未授权 | 检查token是否有效 |
| 403 | 无权限 | 检查用户角色 |
| 404 | 资源不存在 | 检查路径和ID |
| 500 | 服务器内部错误 | 查看后端日志 |

---

## 🔍 调试技巧

### 浏览器调试
```javascript
// 1. 开启Vue DevTools
// 2. 查看Network面板的API请求
// 3. 查看Console的错误日志
// 4. 使用Vue DevTools插件

// 5. 在代码中添加调试日志
console.log('API Request:', request);
console.log('API Response:', response);
console.error('Error:', error);
```

### 后端日志查看
```bash
# 查看实时日志
tail -f /var/log/dmh-api.log

# 查看错误日志
tail -f /var/log/dmh-error.log

# 搜索特定订单的日志
grep "orderId:123" /var/log/dmh-api.log
```

---

## ✅ 开发检查清单

### 提交前检查
- [ ] 代码通过ESLint检查
- [ ] 代码通过Prettier格式化
- [ ] 所有API调用都有错误处理
- [ ] 所有表单都有验证
- [ ] 所有异步操作都有loading状态
- [ ] 所有用户操作都有反馈（Toast/Dialog）

### 代码审查要点
- [ ] 组件命名清晰规范
- [ ] 函数单一职责
- [ ] 避免重复代码
- [ ] 注释充分但不冗余
- [ ] 性能优化（避免不必要的渲染）
- [ ] 错误边界处理完整

---

## 📝 备注

- 本文档会持续更新
- 如有问题请联系项目负责人
- 代码提交前请运行完整测试

---

**Document Version**: 1.0  
**Last Updated**: 2026-01-28
**Next Review**: 开发完成后更新
