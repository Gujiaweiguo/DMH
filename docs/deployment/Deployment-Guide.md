# RBAC权限管理系统部署指南

## 概述

本指南详细说明了DMH数字营销中台RBAC权限管理系统的部署流程、配置要求和最佳实践。

## 推荐部署方式（生产环境）

生产环境推荐使用本仓库提供的 **Docker Compose + Nginx** 一键部署脚本，默认部署目录为 `/opt/dmh`。

### 前置条件

- 已安装 `docker` 与 `docker compose`
- 服务器开放 80/443（对外）与 3306/8080（按需，仅内部可关闭）

### 一键部署

在仓库根目录执行：

```bash
bash docs/deployment/scripts/deploy.sh production latest
```

首次运行会在 `/opt/dmh/.env` 生成模板，请至少配置：

- `MYSQL_ROOT_PASSWORD`
- `MYSQL_PASSWORD`（dmh_user 使用）
- `JWT_SECRET`（生产环境必须更换，建议 32 位以上随机串）

部署成功后访问：

- 管理后台：`http://<server>/`
- H5：`http://<server>/h5/`
- API：`http://<server>/api/v1/...`
- 健康检查：`http://<server>/health`

### 目录与配置说明

- `/opt/dmh/current`：代码仓库（脚本拉取/更新）
- `/opt/dmh/docker-compose.yml`：生产 compose（引用 `./current/*` 的相对路径）
- `/opt/dmh/config/*`：可持久化的配置（nginx/mysql/redis 与后端 `dmh-api.yaml`）
- `/opt/dmh/.env`：敏感变量（密码、JWT 等）

## 系统要求

### 硬件要求

#### 最小配置
- **CPU**: 2核心
- **内存**: 4GB RAM
- **存储**: 50GB SSD
- **网络**: 100Mbps

#### 推荐配置
- **CPU**: 4核心以上
- **内存**: 8GB RAM以上
- **存储**: 100GB SSD以上
- **网络**: 1Gbps

#### 生产环境配置
- **CPU**: 8核心以上
- **内存**: 16GB RAM以上
- **存储**: 200GB SSD以上（支持RAID）
- **网络**: 10Gbps
- **负载均衡**: 支持
- **高可用**: 多节点部署

### 软件要求

#### 操作系统
- **Linux**: Ubuntu 20.04+ / CentOS 8+ / RHEL 8+
- **容器**: Docker 20.10+ / Kubernetes 1.20+

#### 数据库
- **MySQL**: 8.0+
- **连接池**: 推荐使用连接池
- **备份**: 定期备份策略

#### 运行时环境
- **Go**: 1.23+（以 `backend/go.mod` 为准）
- **Node.js**: 20+（前端构建）
- **Nginx**: 1.20+ (反向代理)

## 部署架构

### 单机部署架构

```
┌─────────────────┐
│   Load Balancer │
│    (Nginx)      │
└─────────┬───────┘
          │
┌─────────▼───────┐
│   DMH Backend   │
│   (Go Service)  │
└─────────┬───────┘
          │
┌─────────▼───────┐
│   MySQL DB      │
│   (Data Store)  │
└─────────────────┘
```

### 高可用部署架构

```
┌─────────────────┐    ┌─────────────────┐
│ Load Balancer 1 │    │ Load Balancer 2 │
│    (Nginx)      │    │    (Nginx)      │
└─────────┬───────┘    └─────────┬───────┘
          │                      │
          └──────────┬───────────┘
                     │
    ┌────────────────┼────────────────┐
    │                │                │
┌───▼────┐    ┌──────▼──┐    ┌───────▼┐
│Backend1│    │Backend2 │    │Backend3│
│(Go App)│    │(Go App) │    │(Go App)│
└───┬────┘    └──────┬──┘    └───────┬┘
    │                │                │
    └────────────────┼────────────────┘
                     │
    ┌────────────────┼────────────────┐
    │                │                │
┌───▼────┐    ┌──────▼──┐    ┌───────▼┐
│MySQL   │    │MySQL    │    │MySQL   │
│Master  │◄──►│Slave1   │◄──►│Slave2  │
└────────┘    └─────────┘    └────────┘
```

## 部署步骤

### 1. 环境准备

#### 1.1 创建部署用户

```bash
# 创建部署用户
sudo useradd -m -s /bin/bash dmh
sudo usermod -aG sudo dmh

# 切换到部署用户
su - dmh
```

#### 1.2 创建目录结构

```bash
# 创建应用目录
mkdir -p /opt/dmh/{bin,config,logs,data,backup}
mkdir -p /opt/dmh/frontend/{admin,h5}

# 设置权限
sudo chown -R dmh:dmh /opt/dmh
chmod 755 /opt/dmh
```

#### 1.3 安装依赖软件

```bash
# 更新系统
sudo apt update && sudo apt upgrade -y

# 安装基础工具
sudo apt install -y curl wget git unzip nginx mysql-client

# 安装Go环境
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 安装Node.js
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs
```

### 2. 数据库部署

#### 2.1 安装MySQL

```bash
# 安装MySQL 8.0
sudo apt install -y mysql-server-8.0

# 启动MySQL服务
sudo systemctl start mysql
sudo systemctl enable mysql

# 安全配置
sudo mysql_secure_installation
```

#### 2.2 创建数据库和用户

```sql
-- 连接到MySQL
mysql -u root -p

-- 创建数据库
CREATE DATABASE dmh DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 创建用户
CREATE USER 'dmh_user'@'localhost' IDENTIFIED BY 'your_secure_password';
CREATE USER 'dmh_user'@'%' IDENTIFIED BY 'your_secure_password';

-- 授权
GRANT ALL PRIVILEGES ON dmh.* TO 'dmh_user'@'localhost';
GRANT ALL PRIVILEGES ON dmh.* TO 'dmh_user'@'%';
FLUSH PRIVILEGES;

-- 退出
EXIT;
```

#### 2.3 初始化数据库

```bash
# 执行初始化脚本
mysql -u dmh_user -p dmh < /opt/dmh/scripts/init.sql
```

#### 2.4 数据库配置优化

```bash
# 编辑MySQL配置
sudo nano /etc/mysql/mysql.conf.d/mysqld.cnf

# 添加优化配置
[mysqld]
# 基础配置
max_connections = 1000
max_connect_errors = 10000
table_open_cache = 2048
max_allowed_packet = 64M

# InnoDB配置
innodb_buffer_pool_size = 2G
innodb_log_file_size = 256M
innodb_log_buffer_size = 16M
innodb_flush_log_at_trx_commit = 2

# 查询缓存
query_cache_type = 1
query_cache_size = 256M

# 慢查询日志
slow_query_log = 1
slow_query_log_file = /var/log/mysql/slow.log
long_query_time = 2

# 重启MySQL
sudo systemctl restart mysql
```

### 3. 后端服务部署

#### 3.1 编译后端服务

```bash
# 克隆代码
cd /opt/dmh
git clone https://github.com/your-org/dmh-backend.git backend
cd backend

# 安装依赖
go mod download

# 编译服务
go build -o /opt/dmh/bin/dmh-api ./api/dmh.go

# 设置执行权限
chmod +x /opt/dmh/bin/dmh-api
```

#### 3.2 配置文件

创建配置文件 `/opt/dmh/config/config.yaml`:

```yaml
# 服务配置
Name: dmh-api
Host: 0.0.0.0
Port: 8080
Mode: prod

# 数据库配置
Mysql:
  DataSource: dmh_user:your_secure_password@tcp(localhost:3306)/dmh?charset=utf8mb4&parseTime=true&loc=Local

# JWT配置
Auth:
  AccessSecret: "your-jwt-secret-key-change-in-production"
  AccessExpire: 86400  # 24小时

# 日志配置
Log:
  ServiceName: dmh-api
  Mode: file
  Path: /opt/dmh/logs
  Level: info
  Compress: true
  KeepDays: 30

# Redis配置（可选）
Redis:
  Host: localhost:6379
  Type: node
  Pass: ""

# 安全配置
Security:
  # 密码策略
  PasswordPolicy:
    MinLength: 8
    RequireUppercase: true
    RequireLowercase: true
    RequireNumbers: true
    RequireSpecialChars: true
    MaxAge: 90
    HistoryCount: 5
    MaxLoginAttempts: 5
    LockoutDuration: 30
    SessionTimeout: 480
    MaxConcurrentSessions: 3
  
  # CORS配置
  CORS:
    AllowOrigins: ["https://admin.dmh.com", "https://h5.dmh.com"]
    AllowMethods: ["GET", "POST", "PUT", "DELETE", "OPTIONS"]
    AllowHeaders: ["Content-Type", "Authorization"]
    ExposeHeaders: ["Content-Length"]
    AllowCredentials: true
    MaxAge: 86400

# 文件上传配置
Upload:
  MaxSize: 10485760  # 10MB
  AllowedTypes: ["image/jpeg", "image/png", "image/gif", "application/pdf"]
  StoragePath: /opt/dmh/data/uploads
```

#### 3.3 创建系统服务

创建服务文件 `/etc/systemd/system/dmh-api.service`:

```ini
[Unit]
Description=DMH API Service
After=network.target mysql.service
Wants=mysql.service

[Service]
Type=simple
User=dmh
Group=dmh
WorkingDirectory=/opt/dmh
ExecStart=/opt/dmh/bin/dmh-api -f /opt/dmh/config/config.yaml
ExecReload=/bin/kill -HUP $MAINPID
KillMode=process
Restart=on-failure
RestartSec=5s
StandardOutput=journal
StandardError=journal
SyslogIdentifier=dmh-api

# 安全配置
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/opt/dmh/logs /opt/dmh/data

# 资源限制
LimitNOFILE=65536
LimitNPROC=32768

[Install]
WantedBy=multi-user.target
```

#### 3.4 启动后端服务

```bash
# 重载systemd配置
sudo systemctl daemon-reload

# 启动服务
sudo systemctl start dmh-api

# 设置开机自启
sudo systemctl enable dmh-api

# 检查服务状态
sudo systemctl status dmh-api

# 查看日志
sudo journalctl -u dmh-api -f
```

### 4. 前端部署

#### 4.1 构建前端项目

```bash
# 构建管理后台
cd /opt/dmh/frontend-admin
npm install
npm run build

# 构建H5前端
cd /opt/dmh/frontend-h5
npm install
npm run build

# 复制构建文件
cp -r /opt/dmh/frontend-admin/dist/* /opt/dmh/frontend/admin/
cp -r /opt/dmh/frontend-h5/dist/* /opt/dmh/frontend/h5/
```

#### 4.2 配置Nginx

创建Nginx配置文件 `/etc/nginx/sites-available/dmh`:

```nginx
# 管理后台配置
server {
    listen 80;
    server_name admin.dmh.com;
    
    # 重定向到HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name admin.dmh.com;
    
    # SSL证书配置
    ssl_certificate /etc/ssl/certs/dmh.crt;
    ssl_certificate_key /etc/ssl/private/dmh.key;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512:ECDHE-RSA-AES256-GCM-SHA384:DHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;
    
    # 安全头
    add_header X-Frame-Options DENY;
    add_header X-Content-Type-Options nosniff;
    add_header X-XSS-Protection "1; mode=block";
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    
    # 静态文件
    root /opt/dmh/frontend/admin;
    index index.html;
    
    # Gzip压缩
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types text/plain text/css text/xml text/javascript application/javascript application/xml+rss application/json;
    
    # 静态资源缓存
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }
    
    # API代理
    location /api/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # 超时配置
        proxy_connect_timeout 30s;
        proxy_send_timeout 30s;
        proxy_read_timeout 30s;
        
        # 缓冲配置
        proxy_buffering on;
        proxy_buffer_size 4k;
        proxy_buffers 8 4k;
    }
    
    # SPA路由支持
    location / {
        try_files $uri $uri/ /index.html;
    }
    
    # 访问日志
    access_log /var/log/nginx/dmh-admin.access.log;
    error_log /var/log/nginx/dmh-admin.error.log;
}

# H5前端配置
server {
    listen 80;
    server_name h5.dmh.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name h5.dmh.com;
    
    # SSL配置（同上）
    ssl_certificate /etc/ssl/certs/dmh.crt;
    ssl_certificate_key /etc/ssl/private/dmh.key;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512:ECDHE-RSA-AES256-GCM-SHA384:DHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;
    
    # 安全头
    add_header X-Frame-Options SAMEORIGIN;
    add_header X-Content-Type-Options nosniff;
    add_header X-XSS-Protection "1; mode=block";
    
    root /opt/dmh/frontend/h5;
    index index.html;
    
    # 静态资源配置（同上）
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }
    
    # API代理
    location /api/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
    
    location / {
        try_files $uri $uri/ /index.html;
    }
    
    access_log /var/log/nginx/dmh-h5.access.log;
    error_log /var/log/nginx/dmh-h5.error.log;
}

# API服务配置
server {
    listen 80;
    server_name api.dmh.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name api.dmh.com;
    
    # SSL配置
    ssl_certificate /etc/ssl/certs/dmh.crt;
    ssl_certificate_key /etc/ssl/private/dmh.key;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512:ECDHE-RSA-AES256-GCM-SHA384:DHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;
    
    # 限流配置
    limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;
    limit_req zone=api burst=20 nodelay;
    
    # API代理
    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # CORS配置
        add_header Access-Control-Allow-Origin $http_origin;
        add_header Access-Control-Allow-Methods "GET, POST, PUT, DELETE, OPTIONS";
        add_header Access-Control-Allow-Headers "Content-Type, Authorization";
        add_header Access-Control-Allow-Credentials true;
        
        if ($request_method = 'OPTIONS') {
            return 204;
        }
    }
    
    access_log /var/log/nginx/dmh-api.access.log;
    error_log /var/log/nginx/dmh-api.error.log;
}
```

#### 4.3 启用Nginx配置

```bash
# 启用站点配置
sudo ln -s /etc/nginx/sites-available/dmh /etc/nginx/sites-enabled/

# 测试配置
sudo nginx -t

# 重载Nginx
sudo systemctl reload nginx
```

### 5. SSL证书配置

#### 5.1 使用Let's Encrypt

```bash
# 安装Certbot
sudo apt install -y certbot python3-certbot-nginx

# 获取证书
sudo certbot --nginx -d admin.dmh.com -d h5.dmh.com -d api.dmh.com

# 设置自动续期
sudo crontab -e
# 添加以下行
0 12 * * * /usr/bin/certbot renew --quiet
```

#### 5.2 使用自签名证书（开发环境）

```bash
# 创建证书目录
sudo mkdir -p /etc/ssl/private

# 生成私钥
sudo openssl genrsa -out /etc/ssl/private/dmh.key 2048

# 生成证书
sudo openssl req -new -x509 -key /etc/ssl/private/dmh.key -out /etc/ssl/certs/dmh.crt -days 365

# 设置权限
sudo chmod 600 /etc/ssl/private/dmh.key
sudo chmod 644 /etc/ssl/certs/dmh.crt
```

### 6. 监控和日志

#### 6.1 日志配置

```bash
# 创建日志轮转配置
sudo nano /etc/logrotate.d/dmh

# 添加配置
/opt/dmh/logs/*.log {
    daily
    missingok
    rotate 30
    compress
    delaycompress
    notifempty
    create 644 dmh dmh
    postrotate
        systemctl reload dmh-api
    endscript
}

/var/log/nginx/dmh-*.log {
    daily
    missingok
    rotate 30
    compress
    delaycompress
    notifempty
    create 644 www-data www-data
    postrotate
        systemctl reload nginx
    endscript
}
```

#### 6.2 系统监控

```bash
# 安装监控工具
sudo apt install -y htop iotop nethogs

# 创建监控脚本
cat > /opt/dmh/scripts/monitor.sh << 'EOF'
#!/bin/bash

# 检查服务状态
check_service() {
    local service=$1
    if systemctl is-active --quiet $service; then
        echo "✓ $service is running"
    else
        echo "✗ $service is not running"
        systemctl status $service
    fi
}

# 检查端口
check_port() {
    local port=$1
    local service=$2
    if netstat -tuln | grep -q ":$port "; then
        echo "✓ $service (port $port) is listening"
    else
        echo "✗ $service (port $port) is not listening"
    fi
}

echo "=== DMH System Status ==="
echo "Time: $(date)"
echo

echo "--- Services ---"
check_service dmh-api
check_service nginx
check_service mysql

echo
echo "--- Ports ---"
check_port 8080 "DMH API"
check_port 80 "Nginx HTTP"
check_port 443 "Nginx HTTPS"
check_port 3306 "MySQL"

echo
echo "--- System Resources ---"
echo "CPU Usage: $(top -bn1 | grep "Cpu(s)" | awk '{print $2}' | cut -d'%' -f1)%"
echo "Memory Usage: $(free | grep Mem | awk '{printf("%.1f%%", $3/$2 * 100.0)}')"
echo "Disk Usage: $(df -h / | awk 'NR==2{printf "%s", $5}')"

echo
echo "--- Database ---"
mysql -u dmh_user -p'your_secure_password' -e "SELECT COUNT(*) as user_count FROM dmh.users;" 2>/dev/null || echo "Database connection failed"
EOF

chmod +x /opt/dmh/scripts/monitor.sh

# 设置定时监控
crontab -e
# 添加以下行
*/5 * * * * /opt/dmh/scripts/monitor.sh >> /opt/dmh/logs/monitor.log 2>&1
```

### 7. 备份策略

#### 7.1 数据库备份

```bash
# 创建备份脚本
cat > /opt/dmh/scripts/backup.sh << 'EOF'
#!/bin/bash

BACKUP_DIR="/opt/dmh/backup"
DATE=$(date +%Y%m%d_%H%M%S)
DB_NAME="dmh"
DB_USER="dmh_user"
DB_PASS="your_secure_password"

# 创建备份目录
mkdir -p $BACKUP_DIR/database
mkdir -p $BACKUP_DIR/files

# 数据库备份
mysqldump -u $DB_USER -p$DB_PASS $DB_NAME > $BACKUP_DIR/database/dmh_$DATE.sql

# 压缩备份文件
gzip $BACKUP_DIR/database/dmh_$DATE.sql

# 文件备份
tar -czf $BACKUP_DIR/files/files_$DATE.tar.gz /opt/dmh/data/uploads

# 清理7天前的备份
find $BACKUP_DIR -name "*.gz" -mtime +7 -delete

echo "Backup completed: $DATE"
EOF

chmod +x /opt/dmh/scripts/backup.sh

# 设置定时备份
crontab -e
# 添加以下行
0 2 * * * /opt/dmh/scripts/backup.sh >> /opt/dmh/logs/backup.log 2>&1
```

#### 7.2 配置文件备份

```bash
# 备份重要配置文件
mkdir -p /opt/dmh/backup/config
cp /opt/dmh/config/config.yaml /opt/dmh/backup/config/
cp /etc/nginx/sites-available/dmh /opt/dmh/backup/config/
cp /etc/systemd/system/dmh-api.service /opt/dmh/backup/config/
```

## 安全配置

### 1. 防火墙配置

```bash
# 安装UFW
sudo apt install -y ufw

# 默认策略
sudo ufw default deny incoming
sudo ufw default allow outgoing

# 允许SSH
sudo ufw allow ssh

# 允许HTTP/HTTPS
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# 允许MySQL（仅本地）
sudo ufw allow from 127.0.0.1 to any port 3306

# 启用防火墙
sudo ufw enable

# 查看状态
sudo ufw status verbose
```

### 2. 系统安全加固

```bash
# 更新系统
sudo apt update && sudo apt upgrade -y

# 安装安全工具
sudo apt install -y fail2ban unattended-upgrades

# 配置fail2ban
sudo cp /etc/fail2ban/jail.conf /etc/fail2ban/jail.local
sudo nano /etc/fail2ban/jail.local

# 添加DMH API保护
cat >> /etc/fail2ban/jail.local << 'EOF'

[dmh-api]
enabled = true
port = 8080
filter = dmh-api
logpath = /opt/dmh/logs/access.log
maxretry = 5
bantime = 3600
findtime = 600
EOF

# 创建过滤规则
cat > /etc/fail2ban/filter.d/dmh-api.conf << 'EOF'
[Definition]
failregex = ^.*"POST /api/v1/auth/login.*" 401.*$
ignoreregex =
EOF

# 重启fail2ban
sudo systemctl restart fail2ban
```

### 3. 数据库安全

```bash
# MySQL安全配置
mysql -u root -p << 'EOF'
-- 删除匿名用户
DELETE FROM mysql.user WHERE User='';

-- 删除测试数据库
DROP DATABASE IF EXISTS test;
DELETE FROM mysql.db WHERE Db='test' OR Db='test\\_%';

-- 禁用远程root登录
DELETE FROM mysql.user WHERE User='root' AND Host NOT IN ('localhost', '127.0.0.1', '::1');

-- 刷新权限
FLUSH PRIVILEGES;
EOF

# 配置SSL连接
sudo nano /etc/mysql/mysql.conf.d/mysqld.cnf
# 添加以下配置
[mysqld]
ssl-ca=/etc/mysql/ssl/ca-cert.pem
ssl-cert=/etc/mysql/ssl/server-cert.pem
ssl-key=/etc/mysql/ssl/server-key.pem
require_secure_transport=ON
```

## 性能优化

### 1. 系统优化

```bash
# 优化系统参数
sudo nano /etc/sysctl.conf

# 添加以下配置
# 网络优化
net.core.somaxconn = 65535
net.core.netdev_max_backlog = 5000
net.ipv4.tcp_max_syn_backlog = 65535
net.ipv4.tcp_fin_timeout = 30
net.ipv4.tcp_keepalive_time = 1200
net.ipv4.tcp_max_tw_buckets = 5000

# 文件描述符限制
fs.file-max = 65535

# 应用配置
sudo sysctl -p

# 设置用户限制
sudo nano /etc/security/limits.conf
# 添加以下行
dmh soft nofile 65535
dmh hard nofile 65535
dmh soft nproc 32768
dmh hard nproc 32768
```

### 2. 应用优化

```bash
# Go应用优化环境变量
cat >> /etc/systemd/system/dmh-api.service << 'EOF'

[Service]
Environment=GOGC=100
Environment=GOMAXPROCS=4
Environment=GOMEMLIMIT=2GiB
EOF

# 重载服务配置
sudo systemctl daemon-reload
sudo systemctl restart dmh-api
```

## 故障排除

### 1. 常见问题

#### 服务无法启动
```bash
# 检查服务状态
sudo systemctl status dmh-api

# 查看详细日志
sudo journalctl -u dmh-api -f

# 检查配置文件
/opt/dmh/bin/dmh-api -f /opt/dmh/config/config.yaml -test
```

#### 数据库连接失败
```bash
# 测试数据库连接
mysql -u dmh_user -p -h localhost dmh

# 检查MySQL状态
sudo systemctl status mysql

# 查看MySQL错误日志
sudo tail -f /var/log/mysql/error.log
```

#### 前端访问异常
```bash
# 检查Nginx状态
sudo systemctl status nginx

# 测试Nginx配置
sudo nginx -t

# 查看Nginx日志
sudo tail -f /var/log/nginx/error.log
```

### 2. 性能问题排查

```bash
# 检查系统资源
htop
iotop
nethogs

# 检查数据库性能
mysql -u dmh_user -p -e "SHOW PROCESSLIST;"
mysql -u dmh_user -p -e "SHOW ENGINE INNODB STATUS\G"

# 检查慢查询
sudo tail -f /var/log/mysql/slow.log
```

### 3. 日志分析

```bash
# 分析访问日志
sudo tail -f /var/log/nginx/dmh-admin.access.log | grep -E "(4[0-9]{2}|5[0-9]{2})"

# 分析应用日志
tail -f /opt/dmh/logs/access.log | grep ERROR

# 统计错误频率
grep ERROR /opt/dmh/logs/access.log | awk '{print $1}' | sort | uniq -c | sort -nr
```

## 维护指南

### 1. 日常维护

```bash
# 每日检查脚本
cat > /opt/dmh/scripts/daily_check.sh << 'EOF'
#!/bin/bash

echo "=== Daily Maintenance Check ==="
echo "Date: $(date)"

# 检查磁盘空间
echo "--- Disk Usage ---"
df -h

# 检查内存使用
echo "--- Memory Usage ---"
free -h

# 检查服务状态
echo "--- Service Status ---"
systemctl is-active dmh-api nginx mysql

# 检查错误日志
echo "--- Recent Errors ---"
grep ERROR /opt/dmh/logs/*.log | tail -10

# 检查数据库大小
echo "--- Database Size ---"
mysql -u dmh_user -p'your_secure_password' -e "
SELECT 
    table_schema AS 'Database',
    ROUND(SUM(data_length + index_length) / 1024 / 1024, 2) AS 'Size (MB)'
FROM information_schema.tables 
WHERE table_schema = 'dmh'
GROUP BY table_schema;
"

echo "=== Check Complete ==="
EOF

chmod +x /opt/dmh/scripts/daily_check.sh
```

### 2. 更新部署

```bash
# 创建更新脚本
cat > /opt/dmh/scripts/update.sh << 'EOF'
#!/bin/bash

set -e

echo "Starting DMH update process..."

# 备份当前版本
BACKUP_DIR="/opt/dmh/backup/$(date +%Y%m%d_%H%M%S)"
mkdir -p $BACKUP_DIR
cp /opt/dmh/bin/dmh-api $BACKUP_DIR/
cp -r /opt/dmh/config $BACKUP_DIR/

# 停止服务
sudo systemctl stop dmh-api

# 更新代码
cd /opt/dmh/backend
git pull origin main

# 编译新版本
go build -o /opt/dmh/bin/dmh-api-new ./api/dmh.go

# 替换二进制文件
mv /opt/dmh/bin/dmh-api /opt/dmh/bin/dmh-api-old
mv /opt/dmh/bin/dmh-api-new /opt/dmh/bin/dmh-api
chmod +x /opt/dmh/bin/dmh-api

# 数据库迁移（如果需要）
# mysql -u dmh_user -p'your_secure_password' dmh < migrations/latest.sql

# 启动服务
sudo systemctl start dmh-api

# 检查服务状态
sleep 5
if systemctl is-active --quiet dmh-api; then
    echo "Update successful!"
    rm /opt/dmh/bin/dmh-api-old
else
    echo "Update failed, rolling back..."
    sudo systemctl stop dmh-api
    mv /opt/dmh/bin/dmh-api-old /opt/dmh/bin/dmh-api
    sudo systemctl start dmh-api
    exit 1
fi
EOF

chmod +x /opt/dmh/scripts/update.sh
```

### 3. 健康检查

```bash
# 创建健康检查脚本
cat > /opt/dmh/scripts/health_check.sh << 'EOF'
#!/bin/bash

# 检查API健康状态
check_api() {
    local response=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/api/v1/health)
    if [ "$response" = "200" ]; then
        echo "✓ API is healthy"
        return 0
    else
        echo "✗ API health check failed (HTTP $response)"
        return 1
    fi
}

# 检查数据库连接
check_database() {
    if mysql -u dmh_user -p'your_secure_password' -e "SELECT 1" dmh >/dev/null 2>&1; then
        echo "✓ Database connection is healthy"
        return 0
    else
        echo "✗ Database connection failed"
        return 1
    fi
}

# 检查前端访问
check_frontend() {
    local response=$(curl -s -o /dev/null -w "%{http_code}" https://admin.dmh.com)
    if [ "$response" = "200" ]; then
        echo "✓ Frontend is accessible"
        return 0
    else
        echo "✗ Frontend access failed (HTTP $response)"
        return 1
    fi
}

# 执行所有检查
echo "=== Health Check ==="
echo "Time: $(date)"

check_api
api_status=$?

check_database
db_status=$?

check_frontend
frontend_status=$?

# 总体状态
if [ $api_status -eq 0 ] && [ $db_status -eq 0 ] && [ $frontend_status -eq 0 ]; then
    echo "✓ All systems healthy"
    exit 0
else
    echo "✗ Some systems are unhealthy"
    exit 1
fi
EOF

chmod +x /opt/dmh/scripts/health_check.sh

# 设置定时健康检查
crontab -e
# 添加以下行
*/5 * * * * /opt/dmh/scripts/health_check.sh >> /opt/dmh/logs/health.log 2>&1
```

## 总结

本部署指南涵盖了DMH RBAC权限管理系统的完整部署流程，包括：

1. **环境准备**: 系统要求、软件安装、目录结构
2. **数据库部署**: MySQL安装、配置、优化
3. **后端服务部署**: 编译、配置、系统服务
4. **前端部署**: 构建、Nginx配置、SSL证书
5. **安全配置**: 防火墙、系统加固、数据库安全
6. **性能优化**: 系统参数、应用优化
7. **监控维护**: 日志管理、备份策略、健康检查

遵循本指南可以确保系统的稳定、安全和高性能运行。在生产环境中，建议根据实际需求调整配置参数，并建立完善的监控和告警机制。
