# DMH权限管理系统故障排除指南

## 概述

本指南提供了DMH数字营销中台RBAC权限管理系统常见问题的诊断和解决方案，帮助运维人员快速定位和解决系统故障。

## 目录

1. [系统监控和诊断](#系统监控和诊断)
2. [服务启动问题](#服务启动问题)
3. [数据库问题](#数据库问题)
4. [网络和连接问题](#网络和连接问题)
5. [权限和认证问题](#权限和认证问题)
6. [性能问题](#性能问题)
7. [前端问题](#前端问题)
8. [日志分析](#日志分析)
9. [应急处理](#应急处理)

## 系统监控和诊断

### 快速健康检查

使用以下脚本快速检查系统状态：

```bash
#!/bin/bash
# 系统健康检查脚本

echo "=== DMH系统健康检查 ==="
echo "时间: $(date)"
echo

# 检查服务状态
echo "--- 服务状态 ---"
if docker compose version >/dev/null 2>&1; then
  docker compose ps
else
    systemctl status dmh-api nginx mysql
fi

echo
echo "--- 端口监听 ---"
netstat -tuln | grep -E ':(80|443|3306|8080|6379) '

echo
echo "--- 系统资源 ---"
echo "CPU使用率: $(top -bn1 | grep "Cpu(s)" | awk '{print $2}' | cut -d'%' -f1)%"
echo "内存使用率: $(free | grep Mem | awk '{printf("%.1f%%", $3/$2 * 100.0)}')"
echo "磁盘使用率: $(df -h / | awk 'NR==2{printf "%s", $5}')"

echo
echo "--- API健康检查 ---"
if curl -f http://localhost:8080/health &> /dev/null; then
    echo "✓ API服务正常"
else
    echo "✗ API服务异常"
fi

echo
echo "--- 数据库连接 ---"
if mysql -u dmh_user -p'your_password' -e "SELECT 1" dmh &> /dev/null; then
    echo "✓ 数据库连接正常"
else
    echo "✗ 数据库连接异常"
fi
```

### 系统指标监控

#### CPU和内存监控

```bash
# 实时监控系统资源
htop

# 查看进程资源使用
ps aux --sort=-%cpu | head -10  # CPU使用率最高的进程
ps aux --sort=-%mem | head -10  # 内存使用率最高的进程

# 查看系统负载
uptime
cat /proc/loadavg
```

#### 磁盘空间监控

```bash
# 查看磁盘使用情况
df -h

# 查看目录大小
du -sh /opt/dmh/*
du -sh /var/log/*

# 查找大文件
find /opt/dmh -type f -size +100M -exec ls -lh {} \;
```

#### 网络监控

```bash
# 查看网络连接
netstat -tuln
ss -tuln

# 查看网络流量
iftop
nethogs

# 测试网络连接
ping google.com
curl -I http://localhost:8080/health
```

## 服务启动问题

### DMH API服务无法启动

#### 症状
- 服务启动失败
- 端口8080无法访问
- 系统日志显示启动错误

#### 诊断步骤

1. **检查服务状态**
```bash
# Docker环境
docker compose ps
docker compose logs dmh-api

# 系统服务环境
systemctl status dmh-api
journalctl -u dmh-api -f
```

2. **检查配置文件**
```bash
# 验证配置文件语法
/opt/dmh/bin/dmh-api -f /opt/dmh/config/config.yaml -test

# 检查配置文件权限
ls -la /opt/dmh/config/config.yaml
```

3. **检查端口占用**
```bash
netstat -tuln | grep 8080
lsof -i :8080
```

#### 常见解决方案

**配置文件错误**
```bash
# 检查YAML语法
python3 -c "import yaml; yaml.safe_load(open('/opt/dmh/config/config.yaml'))"

# 重置配置文件
cp /opt/dmh/config/config.yaml.backup /opt/dmh/config/config.yaml
```

**端口冲突**
```bash
# 杀死占用端口的进程
sudo kill -9 $(lsof -t -i:8080)

# 或修改配置文件中的端口
sed -i 's/Port: 8080/Port: 8081/' /opt/dmh/config/config.yaml
```

**权限问题**
```bash
# 修复文件权限
sudo chown -R dmh:dmh /opt/dmh
chmod +x /opt/dmh/bin/dmh-api
```

### MySQL数据库启动问题

#### 症状
- 数据库连接失败
- 端口3306无法访问
- 应用无法连接数据库

#### 诊断步骤

1. **检查MySQL状态**
```bash
# Docker环境
docker compose logs mysql

# 系统服务环境
systemctl status mysql
journalctl -u mysql -f
```

2. **检查数据库文件**
```bash
# 检查数据目录权限
ls -la /var/lib/mysql/

# 检查磁盘空间
df -h /var/lib/mysql
```

3. **测试数据库连接**
```bash
mysql -u root -p -h localhost
mysql -u dmh_user -p -h localhost dmh
```

#### 常见解决方案

**数据库无法启动**
```bash
# 检查错误日志
tail -f /var/log/mysql/error.log

# 修复数据库
sudo mysqld --user=mysql --skip-grant-tables --skip-networking &
mysql -u root
FLUSH PRIVILEGES;
```

**连接被拒绝**
```bash
# 检查MySQL配置
sudo nano /etc/mysql/mysql.conf.d/mysqld.cnf

# 确保bind-address设置正确
bind-address = 0.0.0.0

# 重启MySQL
sudo systemctl restart mysql
```

**权限问题**
```bash
# 重置用户权限
mysql -u root -p
DROP USER 'dmh_user'@'localhost';
CREATE USER 'dmh_user'@'localhost' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON dmh.* TO 'dmh_user'@'localhost';
FLUSH PRIVILEGES;
```

### Nginx启动问题

#### 症状
- 前端页面无法访问
- 502 Bad Gateway错误
- SSL证书错误

#### 诊断步骤

1. **检查Nginx状态**
```bash
systemctl status nginx
nginx -t  # 测试配置文件语法
```

2. **检查配置文件**
```bash
nginx -T  # 显示完整配置
cat /etc/nginx/sites-enabled/dmh
```

3. **检查日志**
```bash
tail -f /var/log/nginx/error.log
tail -f /var/log/nginx/dmh-admin.error.log
```

#### 常见解决方案

**配置语法错误**
```bash
# 测试配置
nginx -t

# 重载配置
nginx -s reload
```

**SSL证书问题**
```bash
# 检查证书文件
ls -la /etc/ssl/certs/dmh.crt /etc/ssl/private/dmh.key

# 测试证书
openssl x509 -in /etc/ssl/certs/dmh.crt -text -noout
```

**上游服务连接失败**
```bash
# 测试后端服务
curl http://localhost:8080/api/v1/health

# 检查防火墙
sudo ufw status
```

## 数据库问题

### 连接池耗尽

#### 症状
- 应用报告数据库连接超时
- 大量"too many connections"错误
- 系统响应缓慢

#### 诊断步骤

```bash
# 查看当前连接数
mysql -u root -p -e "SHOW PROCESSLIST;"
mysql -u root -p -e "SHOW STATUS LIKE 'Threads_connected';"

# 查看最大连接数
mysql -u root -p -e "SHOW VARIABLES LIKE 'max_connections';"
```

#### 解决方案

```bash
# 临时增加最大连接数
mysql -u root -p -e "SET GLOBAL max_connections = 1000;"

# 永久修改配置
echo "max_connections = 1000" >> /etc/mysql/mysql.conf.d/mysqld.cnf
systemctl restart mysql

# 优化连接池配置
# 在应用配置中调整连接池参数
```

### 慢查询问题

#### 症状
- 数据库响应缓慢
- 应用超时错误
- CPU使用率高

#### 诊断步骤

```bash
# 启用慢查询日志
mysql -u root -p -e "SET GLOBAL slow_query_log = 'ON';"
mysql -u root -p -e "SET GLOBAL long_query_time = 2;"

# 查看慢查询
tail -f /var/log/mysql/slow.log

# 分析慢查询
mysqldumpslow /var/log/mysql/slow.log
```

#### 解决方案

```bash
# 分析查询执行计划
mysql -u root -p dmh
EXPLAIN SELECT * FROM users WHERE username = 'admin';

# 添加索引
ALTER TABLE users ADD INDEX idx_username (username);

# 优化查询语句
# 避免SELECT *，使用具体字段
# 合理使用WHERE条件
# 使用LIMIT限制结果集
```

### 数据库锁问题

#### 症状
- 事务等待超时
- 死锁错误
- 数据更新失败

#### 诊断步骤

```bash
# 查看锁等待情况
mysql -u root -p -e "SHOW ENGINE INNODB STATUS\G" | grep -A 20 "LATEST DETECTED DEADLOCK"

# 查看当前锁
mysql -u root -p -e "SELECT * FROM information_schema.INNODB_LOCKS;"
mysql -u root -p -e "SELECT * FROM information_schema.INNODB_LOCK_WAITS;"
```

#### 解决方案

```bash
# 杀死长时间运行的事务
mysql -u root -p -e "SHOW PROCESSLIST;"
mysql -u root -p -e "KILL [process_id];"

# 优化事务
# 减少事务持有时间
# 按相同顺序访问资源
# 使用合适的隔离级别
```

## 网络和连接问题

### API接口超时

#### 症状
- 前端请求超时
- 502/504错误
- 连接重置

#### 诊断步骤

```bash
# 测试API响应时间
curl -w "@curl-format.txt" -o /dev/null -s http://localhost:8080/api/v1/health

# curl-format.txt内容:
#      time_namelookup:  %{time_namelookup}\n
#         time_connect:  %{time_connect}\n
#      time_appconnect:  %{time_appconnect}\n
#     time_pretransfer:  %{time_pretransfer}\n
#        time_redirect:  %{time_redirect}\n
#   time_starttransfer:  %{time_starttransfer}\n
#                     ----------\n
#           time_total:  %{time_total}\n

# 检查网络延迟
ping localhost
traceroute localhost
```

#### 解决方案

```bash
# 调整Nginx超时设置
proxy_connect_timeout 30s;
proxy_send_timeout 30s;
proxy_read_timeout 30s;

# 调整应用超时设置
# 在config.yaml中增加超时配置
Timeout: 30s

# 优化数据库查询
# 添加索引，优化慢查询
```

### SSL/TLS问题

#### 症状
- HTTPS访问失败
- 证书错误警告
- 握手失败

#### 诊断步骤

```bash
# 测试SSL连接
openssl s_client -connect admin.dmh.com:443

# 检查证书有效期
openssl x509 -in /etc/ssl/certs/dmh.crt -noout -dates

# 测试SSL配置
curl -I https://admin.dmh.com
```

#### 解决方案

```bash
# 更新证书
certbot renew

# 检查证书链
openssl verify -CAfile /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/dmh.crt

# 更新SSL配置
ssl_protocols TLSv1.2 TLSv1.3;
ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512;
```

## 权限和认证问题

### JWT Token问题

#### 症状
- 用户无法登录
- Token验证失败
- 权限检查错误

#### 诊断步骤

```bash
# 检查JWT密钥配置
grep JWT_SECRET /opt/dmh/config/config.yaml

# 测试Token生成
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password"}'

# 验证Token
jwt_token="your_token_here"
curl -H "Authorization: Bearer $jwt_token" \
  http://localhost:8080/api/v1/auth/userinfo
```

#### 解决方案

```bash
# 重新生成JWT密钥
openssl rand -base64 32

# 更新配置文件
sed -i 's/AccessSecret: .*/AccessSecret: "new_secret_key"/' /opt/dmh/config/config.yaml

# 重启服务
systemctl restart dmh-api
```

### 权限检查失败

#### 症状
- 403权限不足错误
- 用户无法访问功能
- 权限配置不生效

#### 诊断步骤

```bash
# 检查用户权限
mysql -u dmh_user -p dmh -e "
SELECT u.username, r.name as role, p.code as permission 
FROM users u 
JOIN user_roles ur ON u.id = ur.user_id 
JOIN roles r ON ur.role_id = r.id 
JOIN role_permissions rp ON r.id = rp.role_id 
JOIN permissions p ON rp.permission_id = p.id 
WHERE u.username = 'admin';
"

# 检查品牌权限
mysql -u dmh_user -p dmh -e "
SELECT u.username, b.name as brand 
FROM users u 
JOIN brand_admins ba ON u.id = ba.user_id 
JOIN brands b ON ba.brand_id = b.id 
WHERE u.username = 'brand_admin';
"
```

#### 解决方案

```bash
# 重新分配权限
mysql -u dmh_user -p dmh -e "
INSERT INTO role_permissions (role_id, permission_id) 
SELECT r.id, p.id 
FROM roles r, permissions p 
WHERE r.code = 'brand_admin' AND p.code = 'brand:manage';
"

# 清理权限缓存
redis-cli FLUSHDB
```

## 性能问题

### 高CPU使用率

#### 症状
- 系统响应缓慢
- CPU使用率持续高于80%
- 负载平均值过高

#### 诊断步骤

```bash
# 查看CPU使用情况
top -p $(pgrep dmh-api)
htop

# 分析进程
perf top -p $(pgrep dmh-api)

# 查看系统调用
strace -p $(pgrep dmh-api)
```

#### 解决方案

```bash
# 优化Go应用
export GOGC=100
export GOMAXPROCS=4

# 添加缓存
# 在应用中实现Redis缓存
# 缓存频繁查询的数据

# 数据库优化
# 添加索引
# 优化查询语句
```

### 高内存使用率

#### 症状
- 内存使用率持续高于90%
- 系统出现OOM错误
- 应用响应缓慢

#### 诊断步骤

```bash
# 查看内存使用
free -h
cat /proc/meminfo

# 分析进程内存
ps aux --sort=-%mem | head -10
pmap -d $(pgrep dmh-api)

# 检查内存泄漏
valgrind --tool=memcheck --leak-check=full ./dmh-api
```

#### 解决方案

```bash
# 调整Go垃圾回收
export GOGC=50  # 更频繁的GC

# 限制容器内存
# 在docker-compose.yml中添加
mem_limit: 2g

# 优化数据库缓存
# 调整MySQL缓存大小
innodb_buffer_pool_size = 1G
```

### 磁盘I/O问题

#### 症状
- 磁盘使用率100%
- 系统响应缓慢
- 数据库查询超时

#### 诊断步骤

```bash
# 查看磁盘I/O
iostat -x 1
iotop

# 查看磁盘使用
df -h
du -sh /opt/dmh/* | sort -hr

# 分析I/O模式
pidstat -d 1
```

#### 解决方案

```bash
# 清理日志文件
find /opt/dmh/logs -name "*.log" -mtime +7 -delete
logrotate -f /etc/logrotate.d/dmh

# 优化数据库
# 调整innodb设置
innodb_flush_log_at_trx_commit = 2
innodb_flush_method = O_DIRECT

# 使用SSD存储
# 迁移数据到SSD
# 调整文件系统参数
```

## 前端问题

### 页面加载失败

#### 症状
- 白屏或404错误
- 静态资源加载失败
- JavaScript错误

#### 诊断步骤

```bash
# 检查静态文件
ls -la /usr/share/nginx/html/admin/
ls -la /opt/dmh/frontend/admin/

# 检查Nginx配置
nginx -t
cat /etc/nginx/sites-enabled/dmh

# 查看浏览器控制台错误
# 使用F12开发者工具
```

#### 解决方案

```bash
# 重新构建前端
cd /opt/dmh/frontend-admin
npm run build

# 检查文件权限
chmod -R 644 /usr/share/nginx/html/admin/
chown -R www-data:www-data /usr/share/nginx/html/admin/

# 清理浏览器缓存
# 或添加缓存控制头
add_header Cache-Control "no-cache, must-revalidate";
```

### API调用失败

#### 症状
- 前端显示网络错误
- CORS错误
- 认证失败

#### 诊断步骤

```bash
# 测试API连接
curl -I http://localhost:8080/api/v1/health

# 检查CORS配置
curl -H "Origin: https://admin.dmh.com" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: X-Requested-With" \
  -X OPTIONS \
  http://localhost:8080/api/v1/auth/login
```

#### 解决方案

```bash
# 配置CORS
# 在config.yaml中添加
CORS:
  AllowOrigins: ["https://admin.dmh.com"]
  AllowMethods: ["GET", "POST", "PUT", "DELETE"]
  AllowHeaders: ["Content-Type", "Authorization"]

# 检查代理配置
# 在nginx配置中
proxy_set_header Host $host;
proxy_set_header X-Real-IP $remote_addr;
```

## 日志分析

### 应用日志分析

```bash
# 查看实时日志
tail -f /opt/dmh/logs/access.log
tail -f /opt/dmh/logs/error.log

# 分析错误日志
grep ERROR /opt/dmh/logs/*.log | tail -20

# 统计错误类型
grep ERROR /opt/dmh/logs/error.log | awk '{print $4}' | sort | uniq -c | sort -nr

# 分析访问模式
awk '{print $1}' /opt/dmh/logs/access.log | sort | uniq -c | sort -nr | head -10
```

### 系统日志分析

```bash
# 查看系统日志
journalctl -u dmh-api -f
journalctl -u nginx -f
journalctl -u mysql -f

# 分析启动问题
journalctl -u dmh-api --since "1 hour ago"

# 查看内核日志
dmesg | tail -20
```

### 数据库日志分析

```bash
# 查看MySQL错误日志
tail -f /var/log/mysql/error.log

# 分析慢查询日志
mysqldumpslow -s t -t 10 /var/log/mysql/slow.log

# 查看二进制日志
mysqlbinlog /var/log/mysql/mysql-bin.000001
```

## 应急处理

### 服务快速恢复

```bash
#!/bin/bash
# 应急恢复脚本

echo "开始应急恢复..."

# 停止所有服务
systemctl stop dmh-api nginx

# 恢复配置文件
cp /opt/dmh/backup/config/config.yaml /opt/dmh/config/
cp /opt/dmh/backup/config/nginx.conf /etc/nginx/sites-available/dmh

# 恢复数据库
mysql -u root -p dmh < /opt/dmh/backup/database.sql

# 重启服务
systemctl start mysql
systemctl start dmh-api
systemctl start nginx

# 检查服务状态
systemctl status dmh-api nginx mysql

echo "应急恢复完成"
```

### 数据备份恢复

```bash
# 紧急数据备份
mysqldump -u root -p --single-transaction --routines --triggers dmh > emergency_backup_$(date +%Y%m%d_%H%M%S).sql

# 恢复数据
mysql -u root -p dmh < emergency_backup_20240101_120000.sql

# 验证数据完整性
mysql -u root -p dmh -e "SELECT COUNT(*) FROM users;"
```

### 流量切换

```bash
# 临时维护页面
cat > /usr/share/nginx/html/maintenance.html << 'EOF'
<!DOCTYPE html>
<html>
<head>
    <title>系统维护中</title>
</head>
<body>
    <h1>系统维护中</h1>
    <p>系统正在进行维护，请稍后再试。</p>
</body>
</html>
EOF

# 修改Nginx配置
cat > /etc/nginx/sites-available/maintenance << 'EOF'
server {
    listen 80 default_server;
    listen 443 ssl default_server;
    
    root /usr/share/nginx/html;
    index maintenance.html;
    
    location / {
        try_files /maintenance.html =503;
    }
}
EOF

# 启用维护模式
ln -sf /etc/nginx/sites-available/maintenance /etc/nginx/sites-enabled/default
nginx -s reload
```

## 监控和告警

### 设置监控脚本

```bash
#!/bin/bash
# 监控脚本

ALERT_EMAIL="admin@dmh.com"
LOG_FILE="/opt/dmh/logs/monitor.log"

check_service() {
    local service=$1
    if ! systemctl is-active --quiet $service; then
        echo "$(date): $service is down" >> $LOG_FILE
        echo "$service服务异常" | mail -s "DMH系统告警" $ALERT_EMAIL
        return 1
    fi
    return 0
}

check_disk_space() {
    local usage=$(df / | awk 'NR==2{print $5}' | cut -d'%' -f1)
    if [ $usage -gt 90 ]; then
        echo "$(date): Disk usage is $usage%" >> $LOG_FILE
        echo "磁盘空间不足: $usage%" | mail -s "DMH系统告警" $ALERT_EMAIL
    fi
}

check_memory() {
    local usage=$(free | grep Mem | awk '{printf("%.0f", $3/$2 * 100.0)}')
    if [ $usage -gt 90 ]; then
        echo "$(date): Memory usage is $usage%" >> $LOG_FILE
        echo "内存使用率过高: $usage%" | mail -s "DMH系统告警" $ALERT_EMAIL
    fi
}

# 执行检查
check_service dmh-api
check_service nginx
check_service mysql
check_disk_space
check_memory
```

### 设置定时任务

```bash
# 添加到crontab
crontab -e

# 每5分钟检查一次
*/5 * * * * /opt/dmh/scripts/monitor.sh

# 每天备份数据库
0 2 * * * /opt/dmh/scripts/backup.sh

# 每周清理日志
0 0 * * 0 find /opt/dmh/logs -name "*.log" -mtime +7 -delete
```

## 联系支持

如果以上方法无法解决问题，请联系技术支持：

- **邮箱**: support@dmh.com
- **电话**: 400-xxx-xxxx
- **紧急热线**: 138-xxxx-xxxx（24小时）

提供以下信息以便快速定位问题：

1. 问题描述和复现步骤
2. 错误日志和截图
3. 系统环境信息
4. 最近的变更记录

---

*本指南将根据实际运维经验持续更新，建议定期查看最新版本。*
