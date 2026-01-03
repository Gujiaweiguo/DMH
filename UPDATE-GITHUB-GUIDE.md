# 🔄 更新现有GitHub项目指南

## 📋 项目信息
- **GitHub仓库**: https://github.com/Gujiaweiguo/DMH.git
- **本地路径**: /opt/code/DMH
- **更新内容**: 完整的RBAC权限系统、活动管理、页面设计器等功能

## 🚀 更新步骤

### 1. 检查当前Git状态

```bash
cd /opt/code/DMH

# 检查是否已经是Git仓库
git status

# 如果不是Git仓库，需要初始化并连接到远程仓库
git init
git remote add origin https://github.com/Gujiaweiguo/DMH.git
```

### 2. 拉取远程仓库最新代码

```bash
# 拉取远程仓库的最新代码
git fetch origin

# 如果远程仓库有内容，先合并
git pull origin main --allow-unrelated-histories
```

### 3. 创建.gitignore文件（如果没有）

```bash
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

### 4. 更新README.md

```bash
cat > README.md << 'EOF'
# 🎯 DMH Digital Marketing Hub

数字营销中台系统，提供完整的营销活动管理、用户权限管理和数据分析功能。

![DMH Logo](https://img.shields.io/badge/DMH-Digital%20Marketing%20Hub-blue)
![Go Version](https://img.shields.io/badge/Go-1.19+-00ADD8)
![Vue Version](https://img.shields.io/badge/Vue.js-3.0+-4FC08D)
![License](https://img.shields.io/badge/License-MIT-green)

## ✨ 主要功能

### 🔐 完整的RBAC权限系统
- **4种用户角色**：平台管理员、品牌管理员、参与者、匿名用户
- **JWT认证**：Token自动刷新、登录状态管理、安全验证
- **权限控制**：API级别、页面级别、数据级别的全方位权限管理
- **菜单权限**：动态菜单生成和按钮权限控制

### 🎨 可视化活动管理
- **页面设计器**：拖拽式组件设计，包含8种常用组件
- **动态表单**：自定义报名字段，支持文本、手机号、邮箱、选择等类型
- **实时预览**：所见即所得的页面设计体验
- **主题配置**：颜色、字体、布局完全自定义

### 👥 用户和品牌管理
- **用户管理**：创建、编辑、禁用用户账号
- **品牌管理**：品牌信息管理和品牌管理员关系绑定
- **数据统计**：用户行为分析、活动参与统计

## 🏗️ 技术架构

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   H5 Frontend   │    │  Admin Frontend │    │   Backend API   │
│   (Vue.js 3)    │    │    (React)      │    │     (Go)        │
│   Port: 3100    │    │   Port: 3000    │    │   Port: 8888    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────┐
                    │     MySQL       │
                    │   Database      │
                    └─────────────────┘
```

### 后端技术栈
- **Go 1.19+** - 高性能后端服务
- **Gin** - 轻量级Web框架
- **GORM** - 强大的ORM框架
- **JWT** - 安全的身份认证
- **MySQL 8.0+** - 可靠的数据存储

### 前端技术栈
- **Vue.js 3** - 现代化的H5前端
- **Vant UI** - 优秀的移动端组件库
- **React 18** - 功能丰富的管理后台
- **TypeScript** - 类型安全的开发体验

## 🚀 快速开始

### 环境要求
- Go 1.19+
- Node.js 16+
- MySQL 8.0+

### 1. 克隆项目
```bash
git clone https://github.com/Gujiaweiguo/DMH.git
cd DMH
```

### 2. 数据库初始化
```bash
# 创建数据库
mysql -u root -p -e "CREATE DATABASE dmh_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 导入数据库结构和初始数据
mysql -u root -p dmh_db < backend/scripts/init.sql
```

### 3. 启动后端服务
```bash
cd backend
go mod tidy
go run api/dmh-api.go
```
后端服务将在 http://localhost:8888 启动

### 4. 启动H5前端
```bash
cd frontend-h5
npm install
npm run dev
```
H5前端将在 http://localhost:3100 启动

### 5. 启动管理后台
```bash
cd frontend-admin
npm install
npm run dev
```
管理后台将在 http://localhost:3000 启动

## 🌐 访问地址和测试账号

| 服务 | 地址 | 用户类型 | 用户名 | 密码 | 功能描述 |
|------|------|----------|--------|------|----------|
| H5前端 | http://localhost:3100 | 普通用户 | - | - | 浏览活动、参与报名 |
| H5前端 | http://localhost:3100/brand/login | 品牌管理员 | brand_manager | 123456 | 活动管理、页面设计 |
| 管理后台 | http://localhost:3000 | 平台管理员 | admin | 123456 | 系统管理、用户管理 |
| 后端API | http://localhost:8888 | - | - | - | RESTful API服务 |

## 📱 用户角色详解

### 🔧 平台管理员 (platform_admin)
**访问方式**: http://localhost:3000
- ✅ 用户账号管理（创建、编辑、禁用、重置密码）
- ✅ 品牌信息管理（创建、编辑品牌）
- ✅ 权限配置管理（角色权限、菜单权限）
- ✅ 系统设置和全局数据查看
- ✅ 活动管理（查看、编辑所有品牌的活动）

### 🏢 品牌管理员 (brand_admin)
**访问方式**: http://localhost:3100/brand/login
- ✅ 营销活动管理（创建、编辑、启用、暂停）
- ✅ 可视化页面设计器（8种组件、主题配置）
- ✅ 动态表单设计（自定义字段、验证规则）
- ✅ 活动数据分析（参与统计、转化率）
- ✅ 报名信息管理（查看、导出报名数据）
- ✅ 素材管理（上传、管理活动素材）

### 👤 普通用户 (participant)
**访问方式**: http://localhost:3100
- ✅ 浏览活动列表（无需登录）
- ✅ 查看活动详情和页面
- ✅ 填写报名表单参与活动
- ✅ 查看个人报名记录
- ✅ 活动筛选和搜索

## 🎨 核心功能展示

### 可视化页面设计器
```
📦 组件库                    ⚙️ 组件配置                   👁️ 实时预览
├── 🖼️ 横幅图片              ├── 图片URL设置               ├── 页面标题
├── 📝 文本内容              ├── 文本内容编辑             ├── 活动描述  
├── 🎬 视频播放              ├── 字体大小/对齐            ├── 组件预览
├── ⏰ 倒计时                ├── 视频URL配置               ├── 表单字段
├── 💬 用户评价              ├── 倒计时设置               ├── 报名按钮
├── ❓ 常见问题              ├── 评价内容管理             └── 实时更新
├── 📞 联系方式              ├── 问答列表编辑
└── 🔗 社交媒体              └── 联系信息配置
```

### RBAC权限体系
```
用户 (User)
├── 拥有角色 (Has Roles)
│   ├── platform_admin (平台管理员)
│   ├── brand_admin (品牌管理员)
│   ├── participant (参与者)
│   └── anonymous (匿名用户)
│
├── 角色权限 (Role Permissions)
│   ├── 资源权限 (Resource Permissions)
│   ├── 操作权限 (Action Permissions)
│   └── 数据权限 (Data Permissions)
│
└── 菜单权限 (Menu Permissions)
    ├── 页面访问权限
    ├── 按钮操作权限
    └── 功能模块权限
```

## 📊 数据库设计

### 核心数据表
- **users** - 用户基础信息
- **roles** - 角色定义
- **permissions** - 权限定义
- **user_roles** - 用户角色关系
- **role_permissions** - 角色权限关系
- **menus** - 菜单定义
- **role_menus** - 角色菜单权限
- **brands** - 品牌信息
- **campaigns** - 营销活动
- **orders** - 报名订单
- **audit_logs** - 操作审计日志

## 🛠️ 开发指南

### API文档
- **接口定义**: `backend/api/dmh.api`
- **在线文档**: http://localhost:8888/swagger/
- **Postman集合**: 导入 `docs/api/DMH-API.postman_collection.json`

### 前端开发
```bash
# H5前端开发
cd frontend-h5
npm run dev

# 管理后台开发
cd frontend-admin  
npm run dev

# 构建生产版本
npm run build
```

### 后端开发
```bash
# 运行开发服务器
cd backend
go run api/dmh-api.go

# 运行测试
go test ./...

# 构建生产版本
go build -o dmh-api api/dmh-api.go
```

## 🔧 配置说明

### 数据库配置
```go
// backend/api/internal/config/config.go
type Config struct {
    MySQL struct {
        Host     string `json:"host"`
        Port     int    `json:"port"`
        Username string `json:"username"`
        Password string `json:"password"`
        Database string `json:"database"`
    } `json:"mysql"`
}
```

### JWT配置
```go
// JWT密钥和过期时间配置
const (
    JWTSecret = "your-secret-key"
    TokenExpire = 24 * time.Hour
    RefreshExpire = 7 * 24 * time.Hour
)
```

## 🚀 部署指南

### Docker部署
```bash
# 构建镜像
docker-compose build

# 启动服务
docker-compose up -d
```

### 生产环境部署
1. 构建前端静态文件
2. 编译Go后端程序
3. 配置Nginx反向代理
4. 设置MySQL数据库
5. 配置SSL证书

## 🤝 贡献指南

1. Fork 本仓库
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

### 代码规范
- Go代码遵循 `gofmt` 格式
- 前端代码使用 ESLint + Prettier
- 提交信息使用 [Conventional Commits](https://conventionalcommits.org/)

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 🙏 致谢

感谢所有为这个项目做出贡献的开发者！

## 📞 联系方式

- **GitHub Issues**: [提交问题](https://github.com/Gujiaweiguo/DMH/issues)
- **Email**: 项目相关问题咨询
- **微信群**: 扫码加入开发者交流群

---

⭐ **如果这个项目对你有帮助，请给个星标支持！**

🔗 **项目链接**: https://github.com/Gujiaweiguo/DMH
EOF
```

### 5. 添加所有更改

```bash
# 查看当前状态
git status

# 添加所有文件
git add .

# 查看将要提交的更改
git diff --cached
```

### 6. 创建提交

```bash
git commit -m "🚀 Major Update: Complete RBAC System & Campaign Management

✨ New Features:
- Complete RBAC permission system with 4 user roles
- Visual campaign page designer with 8 component types
- Dynamic form builder with custom field validation
- JWT authentication with token refresh mechanism
- User management and brand relationship management
- Menu-based permission control system
- Campaign data analytics and reporting
- Multi-role frontend interfaces (H5 + Admin)

🏗️ Architecture Improvements:
- Migrated from raw SQL to GORM ORM
- Enhanced security with bcrypt password hashing
- Implemented comprehensive error handling
- Added audit logging for all operations
- Optimized database schema with proper indexing

🎨 Frontend Enhancements:
- Vue.js 3 H5 frontend for users and brand managers
- React admin dashboard for platform administrators
- Responsive design with mobile-first approach
- Real-time preview in page designer
- Intuitive user interface with modern UI components

🔧 Technical Stack:
- Backend: Go + Gin + GORM + JWT + MySQL
- Frontend: Vue.js 3 + Vant UI + React + TypeScript
- Database: MySQL 8.0 with complete RBAC schema
- Authentication: JWT with role-based access control

📋 Test Accounts:
- Platform Admin: admin / 123456
- Brand Manager: brand_manager / 123456  
- Participant: user001 / 123456

🌐 Access URLs:
- H5 Frontend: http://localhost:3100
- Admin Dashboard: http://localhost:3000
- Backend API: http://localhost:8888

This update transforms DMH into a production-ready digital marketing platform with enterprise-level features and security."
```

### 7. 推送到GitHub

```bash
# 推送到远程仓库
git push origin main

# 如果遇到冲突，先拉取并合并
git pull origin main --rebase
git push origin main
```

## 🎯 推送后的建议操作

### 1. 更新仓库描述
在GitHub仓库页面：
- 点击仓库名称旁的设置图标
- 更新描述：`数字营销中台系统 - 完整的RBAC权限管理、活动管理、页面设计器`
- 添加网站链接（如果有演示地址）

### 2. 添加Topics标签
添加以下标签：
```
digital-marketing, rbac, campaign-management, vue3, golang, 
jwt-authentication, mysql, marketing-platform, page-designer, 
user-management, brand-management, data-analytics
```

### 3. 创建Release版本
- 点击 "Releases" → "Create a new release"
- Tag version: `v2.0.0`
- Release title: `🎉 DMH v2.0.0 - Complete RBAC & Campaign Management System`
- 描述主要功能和更新内容

### 4. 设置GitHub Pages（可选）
如果想展示文档或演示页面：
- 进入 Settings → Pages
- 选择源分支和文件夹
- 可以展示你创建的HTML演示页面

## ✅ 完成检查清单

- [ ] 代码已推送到GitHub
- [ ] README.md已更新
- [ ] .gitignore已配置
- [ ] 仓库描述已更新
- [ ] Topics标签已添加
- [ ] Release版本已创建
- [ ] 项目文档完整

执行这些步骤后，你的GitHub项目将展现为一个专业的开源项目！