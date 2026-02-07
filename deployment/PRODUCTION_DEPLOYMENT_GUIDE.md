# 高级功能生产环境部署指南

## 📋 概述

本文档提供了 DMH 高级功能（海报生成、支付配置、表单增强、订单核销）的生产环境部署步骤。

**相关文档**:
- [OpenSpec 提案](../openspec/changes/add-campaign-advanced-features/proposal.md)
- [设计文档](../openspec/changes/add-campaign-advanced-features/design.md)
- [任务清单](../openspec/changes/add-campaign-advanced-features/tasks.md)

---

## 🎯 部署前检查清单

### 基础设施检查

- [ ] 服务器 CPU >= 4 核
- [ ] 服务器内存 >= 8 GB
- [ ] 磁盘空间 >= 50 GB（用于存储海报图片）
- [ ] MySQL 8.0+ 已安装并运行
- [ ] Redis 7.0+ 已安装并运行
- [ ] Go 1.23+ 已安装
- [ ] Node.js 20+ 已安装
- [ ] Nginx 1.25+ 已安装

### 配置文件检查

- [ ] 后端配置文件 `backend/api/etc/dmh-api.prod.yaml` 已准备
- [ ] 数据库连接配置正确
- [ ] Redis 连接配置正确
- [ ] JWT Secret 已设置为生产环境密钥
- [ ] 微信支付配置已更新为生产环境参数
- [ ] 日志目录 `/var/log/dmh-api` 已创建并设置正确权限

### 代码检查

- [ ] 后端代码已编译并测试通过
- [ ] H5 前端代码已构建
- [ ] 管理后台代码已构建
- [ ] 数据库迁移脚本已准备
- [ ] 备份脚本已准备

---

## 🗄️ 第一步：数据库迁移

### 1.1 备份现有数据库

```bash
# 创建备份目录
mkdir -p /data/backups/dmh/$(date +%Y%m%d)

# 备份数据库
mysqldump -h<host> -u<user> -p<password> \
  --single-transaction \
  --routines \
  --triggers \
  --events \
  dmh > /data/backups/dmh/$(date +%Y%m%d)/pre_deployment_backup.sql

# 验证备份文件
ls -lh /data/backups/dmh/$(date +%Y%m%d)/pre_deployment_backup.sql
```

### 1.2 执行迁移脚本

```bash
# 进入迁移脚本目录
cd backend/migrations

# 执行高级功能迁移
mysql -h<host> -u<user> -p<password> dmh < 20250124_add_advanced_features.sql

# 执行海报模板数据导入
mysql -h<host> -u<user> -p<password> dmh < insert_templates.sql

# 验证迁移结果
mysql -h<host> -u<user> -p<password> dmh -e "
  -- 验证 campaigns 表新字段
  DESC campaigns;

  -- 验证 orders 表新字段
  DESC orders;

  -- 验证 poster_templates 表
  SELECT COUNT(*) as template_count FROM poster_templates;

  -- 验证数据
  SELECT * FROM poster_templates LIMIT 3;
"
```

### 1.3 验证数据库迁移

```bash
# 执行验证脚本
mysql -h<host> -u<user> -p<password> dmh < verify_migration.sql

# 检查是否有错误
echo "迁移验证完成，请检查上述输出"
```

---

## 🔧 第二步：后端部署

### 2.1 准备后端配置

```bash
# 复制生产配置文件
cp backend/api/etc/dmh-api.yaml backend/api/etc/dmh-api.prod.yaml

# 编辑生产配置
vim backend/api/etc/dmh-api.prod.yaml
```

**关键配置项**:
```yaml
Name: dmh-api
Host: 0.0.0.0
Port: 8889

# 数据库配置（生产环境）
Mysql:
  DataSource: <prod_user>:<prod_password>@tcp(<prod_mysql_host>:3306)/dmh?charset=utf8mb4&parseTime=true&loc=Local

# Redis 配置（生产环境）
Redis:
  Host: <prod_redis_host>:6379
  Type: node
  Pass: "<prod_redis_password>"

# 日志配置（生产环境使用文件日志）
Log:
  ServiceName: dmh-api
  Mode: file
  Path: /var/log/dmh-api
  Level: info

# JWT 配置（生产环境使用强密钥）
Auth:
  AccessSecret: "<strong_secret_key_production>"
  AccessExpire: 86400

# 微信支付配置（生产环境）
WeChatPay:
  AppID: "<wechat_appid>"
  MchID: "<wechat_mchid>"
  APIKey: "<wechat_api_key>"
  APIKeyV3: "<wechat_api_key_v3>"
  APIClientCert: "/etc/dmh/wechat/apiclient_cert.pem"
  APIClientKey: "/etc/dmh/wechat/apiclient_key.pem"
  NotifyURL: "https://<your-domain>/api/v1/payment/wechat/notify"
  RefundNotifyURL: "https://<your-domain>/api/v1/payment/wechat/refund/notify"
  Sandbox: false  # 生产环境使用真实环境
  CacheTTL: 7200

# 频率限制配置（使用 Redis 存储）
RateLimit:
  PosterGenerate:
    MaxRequests: 5
    WindowDuration: 60
    Storage: redis
  Default:
    MaxRequests: 100
    WindowDuration: 60
    Storage: redis
```

### 2.2 编译后端二进制文件

```bash
# 进入后端目录
cd backend

# 编译生产版本
go build -ldflags "-s -w" -o /data/apps/dmh-api api/dmh.go

# 设置执行权限
chmod +x /data/apps/dmh-api

# 验证二进制文件
ls -lh /data/apps/dmh-api
/data/apps/dmh-api version
```

### 2.3 安装微信支付证书

```bash
# 创建证书目录
mkdir -p /etc/dmh/wechat

# 复制微信支付证书
cp <path_to_apiclient_cert.pem> /etc/dmh/wechat/apiclient_cert.pem
cp <path_to_apiclient_key.pem> /etc/dmh/wechat/apiclient_key.pem

# 设置证书权限
chmod 600 /etc/dmh/wechat/*.pem

# 验证证书
ls -la /etc/dmh/wechat/
```

### 2.4 创建海报存储目录

```bash
# 创建海报存储目录
mkdir -p /data/dmh/posters

# 设置目录权限
chown -R <dmh_user>:<dmh_group> /data/dmh/posters
chmod 755 /data/dmh/posters

# 验证目录
ls -la /data/dmh/
```

### 2.5 创建系统服务（systemd）

```bash
# 创建 systemd 服务文件
sudo vim /etc/systemd/system/dmh-api.service
```

**服务文件内容**:
```ini
[Unit]
Description=DMH API Service
After=network.target mysql.service redis.service

[Service]
Type=simple
User=dmh
Group=dmh
WorkingDirectory=/data/apps
Environment="GIN_MODE=release"
ExecStart=/data/apps/dmh-api -f /data/apps/dmh-api.prod.yaml
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal
SyslogIdentifier=dmh-api

[Install]
WantedBy=multi-user.target
```

**启动服务**:
```bash
# 重载 systemd 配置
sudo systemctl daemon-reload

# 启用服务（开机自启）
sudo systemctl enable dmh-api

# 启动服务
sudo systemctl start dmh-api

# 查看服务状态
sudo systemctl status dmh-api

# 查看服务日志
sudo journalctl -u dmh-api -f
```

---

## 🎨 第三步：前端部署

### 3.1 构建 H5 前端

```bash
# 进入 H5 前端目录
cd frontend-h5

# 安装依赖（如果需要）
npm ci --production=false

# 配置生产环境变量
cat > .env.production << 'EOF'
VITE_API_BASE_URL=https://<your-domain>/api
VITE_APP_TITLE=DMH 营销平台
EOF

# 构建生产版本
npm run build

# 验证构建产物
ls -la dist/
du -sh dist/
```

### 3.2 构建管理后台

```bash
# 进入管理后台目录
cd frontend-admin

# 安装依赖（如果需要）
npm ci --production=false

# 配置生产环境变量
cat > .env.production << 'EOF'
VITE_API_BASE_URL=https://<your-domain>/api
VITE_APP_TITLE=DMH 管理后台
EOF

# 构建生产版本
npm run build

# 验证构建产物
ls -la dist/
du -sh dist/
```

### 3.3 部署前端静态文件

```bash
# 创建前端目录
mkdir -p /data/www/dmh/h5
mkdir -p /data/www/dmh/admin

# 复制 H5 前端文件
cp -r frontend-h5/dist/* /data/www/dmh/h5/

# 复制管理后台文件
cp -r frontend-admin/dist/* /data/www/dmh/admin/

# 设置权限
chown -R nginx:nginx /data/www/dmh
chmod -R 755 /data/www/dmh

# 验证文件
ls -la /data/www/dmh/h5/
ls -la /data/www/dmh/admin/
```

### 3.4 配置 Nginx

```bash
# 创建 Nginx 配置文件
sudo vim /etc/nginx/conf.d/dmh.conf
```

**Nginx 配置**:
```nginx
# H5 前端 (端口 3100)
server {
    listen 3100;
    server_name _;
    root /data/www/dmh/h5;
    index index.html;

    # Gzip 压缩
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types text/plain text/css text/xml text/javascript
               application/x-javascript application/xml+rss
               application/json;

    # 静态资源缓存
    location ~* \.(jpg|jpeg|png|gif|ico|css|js|svg|woff2|woff)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    # Vue Router history 模式
    location / {
        try_files $uri $uri/ /index.html;
    }

    # 代理后端 API
    location /api/ {
        proxy_pass http://127.0.0.1:8889/api/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

# 管理后台 (端口 3000)
server {
    listen 3000;
    server_name _;
    root /data/www/dmh/admin;
    index index.html;

    # Gzip 压缩
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types text/plain text/css text/xml text/javascript
               application/x-javascript application/xml+rss
               application/json;

    # 静态资源缓存
    location ~* \.(jpg|jpeg|png|gif|ico|css|js|svg|woff2|woff)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    # Vue Router history 模式
    location / {
        try_files $uri $uri/ /index.html;
    }

    # 代理后端 API
    location /api/ {
        proxy_pass http://127.0.0.1:8889/api/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

**重启 Nginx**:
```bash
# 测试 Nginx 配置
sudo nginx -t

# 重新加载 Nginx
sudo nginx -s reload

# 查看 Nginx 状态
sudo systemctl status nginx
```

---

## ✅ 第四步：功能验证

### 4.1 健康检查

```bash
# 检查后端服务
curl -f http://localhost:8889/api/v1/health || echo "后端服务异常"

# 检查 H5 前端
curl -f http://localhost:3100 || echo "H5 前端异常"

# 检查管理后台
curl -f http://localhost:3000 || echo "管理后台异常"

# 检查数据库连接
mysql -h<host> -u<user> -p<password> -e "SELECT 1" dmh

# 检查 Redis 连接
redis-cli -h<host> -p<port> ping
```

### 4.2 登录验证

```bash
# 测试管理员登录
curl -X POST http://localhost:8889/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"<admin_username>","password":"<admin_password>"}'

# 测试品牌管理员登录
curl -X POST http://localhost:8889/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"<brand_admin_username>","password":"<brand_admin_password>"}'
```

### 4.3 高级功能验证

#### 4.3.1 海报生成功能

```bash
# 获取登录 token
TOKEN=$(curl -s -X POST http://localhost:8889/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"123456"}' | jq -r '.token')

# 创建测试活动
CAMPAIGN_ID=$(curl -s -X POST http://localhost:8889/api/v1/campaigns \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "brandId": 1,
    "name": "海报生成测试活动",
    "description": "用于测试海报生成功能",
    "rewardRule": 10.0,
    "startTime": "2026-02-01T10:00:00",
    "endTime": "2026-12-31T23:59:59",
    "formFields": [{"type":"text","name":"name","label":"姓名","required":true}]
  }' | jq -r '.id')

# 测试海报生成
echo "测试海报生成（活动ID: $CAMPAIGN_ID）"
curl -X POST http://localhost:8889/api/v1/campaigns/$CAMPAIGN_ID/poster \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"templateId":1}' | jq '.'

# 验证海报是否生成成功
ls -la /data/dmh/posters/ | tail -10
```

#### 4.3.2 支付二维码功能

```bash
# 测试支付二维码生成
echo "测试支付二维码生成"
curl -X GET http://localhost:8889/api/v1/campaigns/$CAMPAIGN_ID/payment-qrcode \
  -H "Authorization: Bearer $TOKEN" | jq '.'
```

#### 4.3.3 订单核销功能

```bash
# 创建测试订单
echo "创建测试订单"
ORDER_ID=$(curl -s -X POST http://localhost:8889/api/v1/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"campaignId\": $CAMPAIGN_ID,
    \"phone\": \"13800138000\",
    \"formData\": {\"name\": \"测试用户\"}
  }" | jq -r '.id')

# 获取订单详情（包含核销码）
echo "获取订单详情（订单ID: $ORDER_ID）"
ORDER_DETAIL=$(curl -s -X GET http://localhost:8889/api/v1/orders/$ORDER_ID \
  -H "Authorization: Bearer $TOKEN" | jq '.')

echo "$ORDER_DETAIL" | jq -r '.verificationCode'

# 测试核销订单
echo "测试核销订单"
curl -X POST http://localhost:8889/api/v1/orders/$ORDER_ID/verify \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"code":"'"$(echo "$ORDER_DETAIL" | jq -r '.verificationCode')"'}"' | jq '.'
```

### 4.4 浏览器验证

打开浏览器访问以下地址并测试各项功能：

- **H5 前端**: http://<your-server-ip>:3100
- **管理后台**: http://<your-server-ip>:3000

**测试清单**:
- [ ] 活动列表页面正常加载
- [ ] 创建活动成功
- [ ] 海报生成功能正常
- [ ] 支付二维码显示正常
- [ ] 表单字段配置正常
- [ ] 订单核销功能正常
- [ ] 数据统计显示正常

---

## 🔙 第五步：回滚准备

### 5.1 准备回滚脚本

```bash
# 创建回滚脚本目录
mkdir -p /data/scripts

# 创建数据库回滚脚本
cat > /data/scripts/rollback_advanced_features.sql << 'EOF'
-- 回滚高级功能数据库变更

-- 删除新增的字段（如果需要）
ALTER TABLE campaigns DROP COLUMN payment_config;
ALTER TABLE campaigns DROP COLUMN poster_template_id;
ALTER TABLE orders DROP COLUMN verification_status;
ALTER TABLE orders DROP COLUMN verified_at;
ALTER TABLE orders DROP COLUMN verified_by;
ALTER TABLE orders DROP COLUMN verification_code;

-- 删除新增的表（如果需要）
DROP TABLE IF EXISTS poster_templates;

-- 删除相关的索引
ALTER TABLE orders DROP INDEX IF EXISTS idx_verification_status;
ALTER TABLE orders DROP INDEX IF EXISTS idx_verified_at;
EOF
```

### 5.2 准备代码回滚

```bash
# 保存当前代码版本
cd /opt/code/DMH
git tag pre-deployment-$(date +%Y%m%d-%H%M%S)

# 创建回滚分支
git branch rollback/$(date +%Y%m%d-%H%M%S)

# 记录部署信息
cat > /data/deployment-info-$(date +%Y%m%d).txt << EOF
部署日期: $(date)
部署版本: $(git rev-parse HEAD)
数据库备份: /data/backups/dmh/$(date +%Y%m%d)/pre_deployment_backup.sql
回滚分支: rollback/$(date +%Y%m%d-%H%M%S)
EOF
```

### 5.3 回滚步骤

如果部署后发现问题，按以下步骤回滚：

```bash
# 1. 回滚数据库
mysql -h<host> -u<user> -p<password> dmh < /data/scripts/rollback_advanced_features.sql

# 2. 回滚后端代码
cd /opt/code/DMH
git checkout <previous_stable_commit_tag>

# 3. 重新编译和部署
go build -ldflags "-s -w" -o /data/apps/dmh-api api/dmh.go
sudo systemctl restart dmh-api

# 4. 回滚前端代码
git checkout <previous_stable_commit_tag>
cd frontend-h5 && npm run build
cd frontend-admin && npm run build

# 5. 重新部署前端文件
cp -r frontend-h5/dist/* /data/www/dmh/h5/
cp -r frontend-admin/dist/* /data/www/dmh/admin/
sudo nginx -s reload

# 6. 验证回滚
curl http://localhost:8889/api/v1/health
```

---

## 📊 第六步：监控配置

### 6.1 配置日志监控

```bash
# 创建日志目录
mkdir -p /var/log/dmh-api

# 设置日志轮转
sudo vim /etc/logrotate.d/dmh-api
```

**logrotate 配置**:
```
/var/log/dmh-api/*.log {
    daily
    rotate 30
    compress
    delaycompress
    missingok
    notifempty
    create 0644 dmh dmh
    sharedscripts
    postrotate
        systemctl reload dmh-api > /dev/null 2>&1 || true
    endscript
}
```

### 6.2 配置性能监控

```bash
# 安装监控工具
sudo apt-get install -y htop iotop sysstat

# 配置系统监控
sudo systemctl enable --now sysstat
```

---

## 🎯 部署后检查清单

- [ ] 数据库迁移成功，数据完整性验证通过
- [ ] 后端服务启动成功，无错误日志
- [ ] 前端页面加载正常
- [ ] 海报生成功能测试通过
- [ ] 支付二维码生成功能测试通过
- [ ] 订单核销功能测试通过
- [ ] 表单字段增强功能测试通过
- [ ] 日志输出正常
- [ ] 回滚方案已准备
- [ ] 监控已配置
- [ ] 用户培训材料已准备

---

## 📞 部署支持

### 常见问题

**Q: 数据库迁移失败怎么办？**
A: 立即停止部署，使用备份恢复数据库，检查迁移脚本错误。

**Q: 后端服务无法启动？**
A: 检查日志 `journalctl -u dmh-api -n 50`，排查配置文件和依赖问题。

**Q: 前端页面 404？**
A: 检查 Nginx 配置和文件权限，确保静态文件路径正确。

**Q: 海报生成失败？**
A: 检查 `/data/dmh/posters` 目录权限和磁盘空间。

### 联系方式

- **技术支持**: support@dmh.com
- **紧急联系**: +86-xxx-xxxx-xxxx
- **文档中心**: https://docs.dmh.com

---

**部署状态**: 待执行
**最后更新**: 2026-02-01
**版本**: v1.0.0
