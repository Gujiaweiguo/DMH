# 📚 DMH项目GitHub上传指南

## 🎯 项目概述

**项目名称：** DMH Digital Marketing Hub  
**项目描述：** 数字营销中台系统，包含完整的RBAC权限管理、活动管理、用户管理等功能  
**技术栈：** Go + Vue.js + MySQL + JWT认证

## 📁 项目结构

```
DMH/
├── backend/                 # Go后端API服务
│   ├── api/                # API接口定义和实现
│   ├── model/              # 数据模型
│   ├── scripts/            # 数据库脚本
│   └── ...
├── frontend-admin/         # React管理后台
│   ├── index.tsx          # 主应用文件
│   ├── services/          # API服务
│   └── ...
├── frontend-h5/           # Vue.js H5前端
│   ├── src/               # 源代码
│   ├── public/            # 静态资源
│   └── ...
├── docs/                  # 项目文档
├── openspec/              # OpenSpec规范文档
└── *.html                 # 测试和演示页面
```

## 🚀 上传步骤

### 1. 初始化Git仓库

```bash
# 在项目根目录执行
cd /opt/code/DMH
git init
```

### 2. 创建.gitignore文件

```bash
# 创建.gitignore文件
cat > .gitignore << 'EOF'
# 依赖文件
node_modules/
vendor/

# 构建输出
dist/
build/
*.exe
*.dll
*.so
*.dylib

# 日志文件
*.log
logs/

# 环境配置
.env
.env.local
.env.development.local
.env.test.local
.env.production.local

# IDE配置
.vscode/
.idea/
*.swp
*.swo
*~

# 操作系统
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db

# 临时文件
*.tmp
*.temp
.cache/

# 数据库文件
*.db
*.sqlite
*.sqlite3

# Go特定
*.test
*.out
go.work
go.work.sum

# 前端特定
.nuxt/
.next/
.vuepress/dist/
.serverless/
.fusebox/
.dynamodb/
.tern-port

# 测试覆盖率
coverage/
*.cover
*.coverprofile

# 备份文件
*.bak
*.backup
*-backup.*
*-complete.*
*-fixed.*
EOF
```

### 3. 添加文件到Git

```bash
# 添加所有文件
git add .

# 查看状态
git status
```

### 4. 创建初始提交

```bash
# 创建初始提交
git commit -m "🎉 Initial commit: DMH Digital Marketing Hub

✨ Features:
- Complete RBAC permission system with 4 user roles
- JWT authentication with token refresh
- Campaign management with visual page designer
- User management and brand management
- H5 frontend for users and brand managers
- Admin dashboard for platform administrators
- MySQL database with complete schema
- API documentation and testing pages

🏗️ Architecture:
- Backend: Go + Gin + GORM + JWT
- Frontend: Vue.js 3 + Vant UI (H5) + React (Admin)
- Database: MySQL with RBAC tables
- Authentication: JWT with role-based access control

👥 User Roles:
- platform_admin: System administration
- brand_admin: Brand and campaign management  
- participant: Activity participation
- anonymous: Public access

🔧 Setup:
- Backend API: http://localhost:8888
- H5 Frontend: http://localhost:3100
- Admin Dashboard: http://localhost:3000

📋 Test Accounts:
- admin / 123456 (Platform Admin)
- brand_manager / 123456 (Brand Admin)
- user001 / 123456 (Participant)"
```

### 5. 在GitHub上创建仓库

1. 访问 https://github.com
2. 点击右上角的 "+" 按钮
3. 选择 "New repository"
4. 填写仓库信息：
   - **Repository name:** `dmh-digital-marketing-hub`
   - **Description:** `数字营销中台系统 - Digital Marketing Hub with RBAC, Campaign Management, and Multi-role Frontend`
   - **Visibility:** Public 或 Private (根据需要选择)
   - **不要**勾选 "Add a README file"、"Add .gitignore"、"Choose a license"
5. 点击 "Create repository"

### 6. 连接本地仓库到GitHub

```bash
# 添加远程仓库 (替换YOUR_USERNAME为你的GitHub用户名)
git remote add origin https://github.com/YOUR_USERNAME/dmh-digital-marketing-hub.git

# 设置主分支名称
git branch -M main

# 推送到GitHub
git push -u origin main
```

### 7. 创建README.md文件

```bash
cat > README.md << 'EOF'
# 🎯 DMH Digital Marketing Hub

数字营销中台系统，提供完整的营销活动管理、用户权限管理和数据分析功能。

## ✨ 主要功能

### 🔐 权限管理系统
- **4种用户角色**：平台管理员、品牌管理员、参与者、匿名用户
- **JWT认证**：Token自动刷新、登录状态管理
- **RBAC权限控制**：基于角色的访问控制
- **菜单权限**：动态菜单和按钮权限

### 🎨 活动管理
- **可视化页面设计器**：拖拽式组件设计
- **动态表单配置**：自定义报名字段
- **活动状态管理**：创建、编辑、启用、暂停
- **数据统计分析**：参与人数、转化率等

### 👥 用户管理
- **用户账号管理**：创建、编辑、禁用用户
- **品牌关系管理**：品牌管理员与品牌的绑定关系
- **密码策略**：安全密码规则和重置功能

## 🏗️ 技术架构

### 后端技术栈
- **Go 1.19+** - 主要开发语言
- **Gin** - Web框架
- **GORM** - ORM框架
- **JWT** - 身份认证
- **MySQL** - 数据库

### 前端技术栈
- **Vue.js 3** - H5前端框架
- **Vant UI** - 移动端UI组件库
- **React** - 管理后台框架
- **TypeScript** - 类型安全

## 🚀 快速开始

### 环境要求
- Go 1.19+
- Node.js 16+
- MySQL 8.0+

### 1. 克隆项目
```bash
git clone https://github.com/YOUR_USERNAME/dmh-digital-marketing-hub.git
cd dmh-digital-marketing-hub
```

### 2. 启动后端服务
```bash
cd backend
go mod tidy
go run api/dmh-api.go
```

### 3. 启动H5前端
```bash
cd frontend-h5
npm install
npm run dev
```

### 4. 启动管理后台
```bash
cd frontend-admin
npm install
npm run dev
```

### 5. 初始化数据库
```bash
# 执行数据库脚本
mysql -u root -p < backend/scripts/init.sql
```

## 🌐 访问地址

- **H5前端**: http://localhost:3100
- **管理后台**: http://localhost:3000  
- **后端API**: http://localhost:8888

## 👤 测试账号

| 角色 | 用户名 | 密码 | 功能 |
|------|--------|------|------|
| 平台管理员 | admin | 123456 | 系统管理、用户管理、品牌管理 |
| 品牌管理员 | brand_manager | 123456 | 活动管理、页面设计、数据查看 |
| 普通用户 | user001 | 123456 | 活动参与、报名提交 |

## 📱 用户角色说明

### 🔧 平台管理员 (platform_admin)
- 访问管理后台进行系统配置
- 管理用户账号和权限
- 管理品牌信息和关系
- 查看全局数据统计

### 🏢 品牌管理员 (brand_admin)  
- 登录H5前端的品牌管理功能
- 创建和管理营销活动
- 设计活动页面和表单
- 查看活动数据和报名信息

### 👤 普通用户 (participant)
- 浏览活动列表（无需登录）
- 参与活动报名
- 查看个人报名记录

## 📊 项目特色

### 🎨 可视化页面设计器
- **组件库**：横幅、文本、视频、倒计时等8种组件
- **实时预览**：所见即所得的设计体验
- **主题配置**：颜色、字体、布局自定义
- **表单设计**：动态字段配置和验证

### 🔒 完整的权限体系
- **数据级权限**：用户只能访问授权的数据
- **功能级权限**：基于角色的功能访问控制
- **API级权限**：所有接口都有权限验证
- **前端路由守卫**：页面访问权限控制

### 📈 数据分析功能
- **活动统计**：参与人数、转化率、收益分析
- **用户行为**：报名趋势、活跃度统计
- **品牌数据**：多品牌数据对比分析

## 🛠️ 开发指南

### API文档
- 接口定义：`backend/api/dmh.api`
- 在线文档：http://localhost:8888/swagger/

### 数据库设计
- 完整的RBAC权限表设计
- 活动和订单数据模型
- 用户和品牌关系管理

### 前端组件
- 可复用的UI组件库
- 统一的API调用封装
- 响应式设计适配

## 🤝 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 📞 联系方式

如有问题或建议，请通过以下方式联系：

- 提交 Issue
- 发送邮件
- 微信群讨论

---

⭐ 如果这个项目对你有帮助，请给个星标支持！
EOF
```

### 8. 提交README并推送

```bash
# 添加README文件
git add README.md

# 提交
git commit -m "📝 Add comprehensive README with setup guide and features"

# 推送到GitHub
git push origin main
```

## 🎯 推荐的GitHub仓库设置

### 1. 添加Topics标签
在GitHub仓库页面点击设置图标，添加以下标签：
- `digital-marketing`
- `rbac`
- `campaign-management`
- `vue3`
- `golang`
- `jwt-authentication`
- `mysql`
- `marketing-platform`

### 2. 创建Release
1. 在GitHub仓库页面点击 "Releases"
2. 点击 "Create a new release"
3. 填写版本信息：
   - **Tag version:** `v1.0.0`
   - **Release title:** `🎉 DMH v1.0.0 - Initial Release`
   - **Description:** 描述主要功能和特性

### 3. 设置分支保护
1. 进入仓库设置 → Branches
2. 添加分支保护规则
3. 保护 `main` 分支

## ⚠️ 注意事项

1. **敏感信息**：确保没有提交数据库密码、API密钥等敏感信息
2. **文件大小**：GitHub单个文件限制100MB，仓库建议不超过1GB
3. **许可证**：考虑添加适当的开源许可证
4. **文档**：保持README和代码注释的更新

## 🎉 完成后的效果

上传完成后，你的GitHub仓库将包含：
- ✅ 完整的项目代码
- ✅ 详细的README文档
- ✅ 合适的.gitignore配置
- ✅ 清晰的提交历史
- ✅ 项目标签和描述

这样其他开发者就能轻松理解和使用你的项目了！