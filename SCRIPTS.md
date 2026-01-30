# DMH 管理脚本说明

本文档说明 DMH 项目的各种管理脚本及其用法。

---

## 快速启动/停止

### 启动服务

```bash
./start.sh
```

或

```bash
./dmh.sh start
```

### 停止服务

```bash
./stop.sh
```

或

```bash
./dmh.sh stop
```

---

## 主管理脚本 (dmh.sh)

`dmh.sh` 是主要的管理脚本，提供完整的服务管理功能。

### 用法

```bash
./dmh.sh [命令]
```

### 可用命令

| 命令 | 说明 | 使用场景 |
|------|------|----------|
| `init` | 初始化环境 | 首次运行项目时使用 |
| `start` | 启动所有服务 | 日常启动服务 |
| `stop` | 停止所有服务 | 停止所有运行的服务 |
| `restart` | 重启所有服务 | 修改配置后重启 |
| `status` | 查看服务状态 | 检查服务运行状态 |
| `logs` | 查看服务日志 | 调试和排查问题 |

### 详细说明

#### 1. init - 初始化环境

**首次运行项目时使用**

```bash
./dmh.sh init
```

执行内容：
- 检查/创建 MySQL Docker 容器
- 初始化数据库表结构
- 可选：导入测试数据
- 安装后端依赖（Go modules）
- 安装前端依赖（npm packages）

#### 2. start - 启动服务

**日常启动服务**

```bash
./dmh.sh start
```

执行内容：
- 检查 MySQL 容器状态（如未运行则启动）
- 检查端口占用（8889, 3000, 3100）
- 启动后端服务（Go API）
- 启动 H5 前端（Vue + Vite）
- 启动管理后台（React + Vite）

启动后的服务地址：
- 后端 API: http://localhost:8889
- H5 前端: http://localhost:3100
- 管理后台: http://localhost:3000

#### 3. stop - 停止服务

**停止所有运行的服务**

```bash
./dmh.sh stop
```

执行内容：
- 停止后端服务
- 停止 H5 前端
- 停止管理后台
- 清理残留进程
- 清理占用的端口

**注意**: MySQL 容器不会被停止

#### 4. restart - 重启服务

**修改配置后重启服务**

```bash
./dmh.sh restart
```

等同于：
```bash
./dmh.sh stop
./dmh.sh start
```

#### 5. status - 查看状态

**检查服务运行状态**

```bash
./dmh.sh status
```

显示内容：
- 后端服务状态（PID, Port）
- H5 前端状态（PID, Port）
- 管理后台状态（PID, Port）
- MySQL 容器状态

#### 6. logs - 查看日志

**查看服务日志**

```bash
./dmh.sh logs
```

可选择查看：
1. 后端日志
2. H5 前端日志
3. 管理后台日志
4. 全部日志

---

## 快捷脚本

### start.sh - 快速启动

```bash
./start.sh
```

简化的启动脚本，直接调用 `./dmh.sh start`

### stop.sh - 快速停止

```bash
./stop.sh
```

独立的停止脚本，功能包括：
- 停止所有服务（后端、H5、管理后台）
- 清理残留进程
- 清理占用的端口（8889, 3000, 3100）
- 显示详细的停止过程

**特点**：
- 更详细的输出信息
- 更彻底的进程清理
- 更友好的提示信息

---

## 日志文件

所有服务的日志和 PID 文件都存储在 `logs/` 目录：

| 文件 | 说明 |
|------|------|
| `logs/backend.log` | 后端服务日志 |
| `logs/backend.pid` | 后端服务 PID |
| `logs/h5.log` | H5 前端日志 |
| `logs/h5.pid` | H5 前端 PID |
| `logs/admin.log` | 管理后台日志 |
| `logs/admin.pid` | 管理后台 PID |
| `logs/dmh-api` | 后端编译后的可执行文件 |

---

## 常见使用场景

### 场景1: 首次运行项目

```bash
# 1. 初始化环境
./dmh.sh init

# 2. 启动服务
./start.sh

# 3. 查看状态
./dmh.sh status
```

### 场景2: 日常开发

```bash
# 启动服务
./start.sh

# 开发中...

# 停止服务
./stop.sh
```

### 场景3: 修改代码后重启

```bash
# 重启所有服务
./dmh.sh restart

# 或者只重启后端
./stop.sh
cd backend
go build -o ../logs/dmh-api api/dmh.go
nohup ../logs/dmh-api -f api/etc/dmh-api.yaml > ../logs/backend.log 2>&1 &
echo $! > ../logs/backend.pid
cd ..
```

### 场景4: 排查问题

```bash
# 查看服务状态
./dmh.sh status

# 查看日志
./dmh.sh logs

# 或直接查看日志文件
tail -f logs/backend.log
tail -f logs/h5.log
tail -f logs/admin.log
```

### 场景5: 清理环境

```bash
# 停止所有服务
./stop.sh

# 清理日志
rm -f logs/*.log

# 停止 MySQL 容器（可选）
docker stop mysql8

# 删除 MySQL 容器和数据（慎用！）
docker rm mysql8
docker volume rm mysql_data
```

---

## 环境要求

### 必需软件

- **Docker**: 用于运行 MySQL 容器
- **Go**: 1.19+ (后端开发)
- **Node.js**: 16+ (前端开发)
- **npm**: 8+ (包管理)

### 可选工具

- **lsof**: 用于检查端口占用（大多数 Linux 系统自带）

---

## 端口说明

| 端口 | 服务 | 说明 |
|------|------|------|
| 3306 | MySQL | 数据库服务 |
| 8889 | Backend API | Go 后端服务 |
| 3000 | Admin Frontend | React 管理后台 |
| 3100 | H5 Frontend | Vue H5 前端 |

---

## 故障排查

### 问题1: 端口被占用

**错误信息**: `✗ 端口 8889 已被占用`

**解决方法**:
```bash
# 查看占用端口的进程
lsof -i :8889

# 停止占用的进程
./stop.sh

# 或手动 kill
kill -9 <PID>
```

### 问题2: MySQL 容器未运行

**错误信息**: `✗ MySQL 容器不存在`

**解决方法**:
```bash
# 重新初始化
./dmh.sh init

# 或手动启动容器
docker start mysql8
```

### 问题3: 服务启动失败

**解决方法**:
```bash
# 查看日志
./dmh.sh logs

# 检查依赖
cd backend && go mod download
cd frontend-h5 && npm install
cd frontend-admin && npm install
```

### 问题4: 进程残留

**解决方法**:
```bash
# 使用 stop.sh 彻底清理
./stop.sh

# 手动清理
pkill -f dmh.go
pkill -f vite
```

---

## 脚本权限

确保所有脚本有执行权限：

```bash
chmod +x dmh.sh
chmod +x start.sh
chmod +x stop.sh
```

---

## 相关文档

- **系统架构**: `ARCHITECTURE.md`
- **启动指南**: `STARTUP.md`
- **开发指南**: `DEVELOPMENT.md`
- **登录信息**: `LOGIN_INFO.md`
- **角色说明**: `ROLES_AND_ENTRIES.md`

---

**维护者**: DMH Team  
**最后更新**: 2025-01-21
