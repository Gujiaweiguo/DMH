# DMH 部署指南

## 技术栈

### 后端
| 技术 | 版本 | 说明 |
|------|------|------|
| Go | 1.21+ | 编程语言 |
| go-zero | 1.6.0 | 微服务框架 |
| GORM | 1.25.5 | ORM 框架 |
| JWT | - | 身份认证 |
| bcrypt | - | 密码加密 |

### 前端
| 技术 | 版本 | 说明 |
|------|------|------|
| Vue | 3.x | 前端框架 |
| Vite | 5.x/6.x | 构建工具 |
| Element Plus | - | UI 组件库 (Admin) |
| Vant | - | UI 组件库 (H5) |
| Pinia | - | 状态管理 |
| Vue Router | 4.x | 路由管理 |
| Axios | - | HTTP 客户端 |

### 基础设施
| 组件 | 版本 | 说明 |
|------|------|------|
| Docker | 20.10+ | 容器运行环境 |
| MySQL | 8.0 | 关系型数据库 |
| Node.js | 20.x | 前端运行环境 |
| npm | 10.x | 包管理器 |

## 一、配置国内镜像源（重要）

### 1.1 Docker 镜像加速
```bash
sudo mkdir -p /etc/docker
sudo tee /etc/docker/daemon.json <<EOF
{
  "registry-mirrors": [
    "https://docker.1ms.run",
    "https://docker.xuanyuan.me",
    "https://hub.rat.dev"
  ]
}
EOF

# 重启 Docker
sudo systemctl restart docker
# WSL2 中使用
sudo service docker restart
```

### 1.2 Go 模块代理
```bash
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GOSUMDB=sum.golang.google.cn
```

### 1.3 npm 淘宝镜像
```bash
npm config set registry https://registry.npmmirror.com
```

## 二、安装 Docker

### Ubuntu/Debian
```bash
# 使用阿里云镜像安装
curl -fsSL https://get.docker.com | sh -s -- --mirror Aliyun

sudo systemctl start docker
sudo systemctl enable docker
sudo usermod -aG docker $USER
```

### WSL2
```bash
curl -fsSL https://get.docker.com | sh
sudo service docker start
```

## 三、安装 Go

```bash
# 下载 Go 1.21
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
# 或使用国内镜像
wget https://golang.google.cn/dl/go1.21.0.linux-amd64.tar.gz

# 解压安装
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz

# 配置环境变量
cat >> ~/.bashrc <<EOF
export PATH=\$PATH:/usr/local/go/bin
export GOPATH=\$HOME/go
export PATH=\$PATH:\$GOPATH/bin
export GOPROXY=https://goproxy.cn,direct
EOF
source ~/.bashrc

# 验证
go version
```

## 四、安装 Node.js

```bash
# 使用 nvm 安装（推荐）
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
source ~/.bashrc

# 安装 Node.js 20
nvm install 20
nvm use 20

# 配置淘宝镜像
npm config set registry https://registry.npmmirror.com

# 验证
node -v
npm -v
```

## 五、启动 MySQL 8

```bash
# 创建数据目录
sudo mkdir -p /opt/data/mysql

# 启动 MySQL 容器
sudo docker run -d \
  --name mysql8 \
  -p 3306:3306 \
  -e MYSQL_ROOT_PASSWORD='#Admin168' \
  -e MYSQL_DATABASE=dmh \
  -v /opt/data/mysql:/var/lib/mysql \
  --restart unless-stopped \
  mysql:8.0 \
  --character-set-server=utf8mb4 \
  --collation-server=utf8mb4_unicode_ci

# 等待启动（约 30 秒）
sleep 30
```

### MySQL 连接信息
| 参数 | 值 |
|------|------|
| Host | `127.0.0.1` 或 `172.17.0.1` (Docker 网关) |
| Port | `3306` |
| User | `root` |
| Password | `#Admin168` |
| Database | `dmh` |

## 六、初始化数据库

### SQL 脚本说明
| 文件 | 说明 | 必需 |
|------|------|------|
| `init.sql` | 主初始化（表结构+基础数据+测试用户） | ✅ 是 |
| `create_member_system_tables.sql` | 会员系统表结构 | ✅ 是 |
| `test_data.sql` | 品牌/活动测试数据 | 开发推荐 |
| `seed_member_campaign_data.sql` | 会员活动测试数据 | 开发推荐 |

### 执行初始化
```bash
# 1. 主初始化脚本（必需）
sudo docker exec -i mysql8 mysql -uroot -p'#Admin168' \
  --default-character-set=utf8mb4 < backend/scripts/init.sql

# 2. 会员系统表（必需）
sudo docker exec -i mysql8 mysql -uroot -p'#Admin168' \
  --default-character-set=utf8mb4 dmh < backend/scripts/create_member_system_tables.sql

# 3. 测试数据（开发环境推荐）
sudo docker exec -i mysql8 mysql -uroot -p'#Admin168' \
  --default-character-set=utf8mb4 dmh < backend/scripts/test_data.sql

sudo docker exec -i mysql8 mysql -uroot -p'#Admin168' \
  --default-character-set=utf8mb4 dmh < backend/scripts/seed_member_campaign_data.sql
```

## 七、部署后端

```bash
cd backend

# 下载依赖（go-zero, gorm 等自动下载）
go mod download

# 编译
go build -o dmh-api api/dmh.go

# 修改配置（如需要）
# vim api/etc/dmh-api.yaml

# 运行
./dmh-api -f api/etc/dmh-api.yaml
```

后端运行在 `http://localhost:8889`

## 八、部署前端

### 管理后台 (Vue 3 + Element Plus)
```bash
cd frontend-admin
npm install
npm run dev      # 开发模式
# npm run build  # 生产构建
```
运行在 `http://localhost:3000`

### H5 前端 (Vue 3 + Vant)
```bash
cd frontend-h5
npm install
npm run dev      # 开发模式
# npm run build  # 生产构建
```
运行在 `http://localhost:3100`

## 九、测试账号

| 用户名 | 密码 | 角色 |
|--------|------|------|
| admin | 123456 | 平台管理员 |
| brand_manager | 123456 | 品牌管理员 |
| user001 | 123456 | 普通用户 |

## 十、快速启动脚本

```bash
#!/bin/bash
# start.sh

# 启动 Docker
sudo service docker start

# 启动 MySQL
sudo docker start mysql8

# 等待 MySQL
sleep 5

# 启动后端
cd backend && ./dmh-api -f api/etc/dmh-api.yaml &

# 启动前端
cd frontend-admin && npm run dev &
cd frontend-h5 && npm run dev &

echo "服务已启动"
echo "  后端: http://localhost:8889"
echo "  管理后台: http://localhost:3000"
echo "  H5: http://localhost:3100"
```

## 十一、常见问题

### MySQL 连接失败
```bash
# WSL2 中需使用 Docker 网关 IP
# 修改 backend/api/etc/dmh-api.yaml
Mysql:
  DataSource: root:#Admin168@tcp(172.17.0.1:3306)/dmh?charset=utf8mb4&parseTime=true&loc=Local
```

### Go 编译失败
```bash
go mod tidy
go mod download
```

### npm 安装慢
```bash
npm config set registry https://registry.npmmirror.com
```
