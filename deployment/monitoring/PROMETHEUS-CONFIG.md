# DMH 性能监控配置指南

## 概述

本文档说明如何为 DMH 系统配置 Prometheus 监控和 Grafana 可视化。

---

## 1. 监控架构

```
应用程序 → Prometheus → Grafana → 告警管理器 → 通知
```

---

## 2. 快速启动

### 2.1 使用 Docker Compose 启动监控组件

在 `deployment/docker-compose-dmh.yml` 中添加以下服务：

```yaml
services:
  # Prometheus
  prometheus:
    image: prom/prometheus:latest
    container_name: dmh-prometheus
    ports:
      - "9090:9090"
    volumes:
      - ./monitoring/prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus-data:/prometheus
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'
      - '--web.console.libraries=/usr/share/prometheus/console_libraries'
      - '--web.console.templates=/usr/share/prometheus/consoles'
      - '--storage.tsdb.retention.time=30d'
    networks:
      - dmh-network
    restart: unless-stopped

  # Grafana
  grafana:
    image: grafana/grafana:latest
    container_name: dmh-grafana
    ports:
      - "3001:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
      - GF_INSTALL_PLUGINS=redis-datasource
    volumes:
      - grafana-data:/var/lib/grafana
      - ./monitoring/grafana/provisioning:/etc/grafana/provisioning
      - ./monitoring/grafana/dashboards:/var/lib/grafana/dashboards
    networks:
      - dmh-network
    depends_on:
      - prometheus
    restart: unless-stopped

  # Alertmanager
  alertmanager:
    image: prom/alertmanager:latest
    container_name: dmh-alertmanager
    ports:
      - "9093:9093"
    volumes:
      - ./monitoring/alertmanager.yml:/etc/alertmanager/alertmanager.yml
    networks:
      - dmh-network
    depends_on:
      - prometheus
    restart: unless-stopped

volumes:
  prometheus-data:
  grafana-data:

networks:
  dmh-network:
    external: true
```

启动监控服务：
```bash
cd deployment
docker compose -f docker-compose-dmh.yml up -d prometheus grafana alertmanager
```

---

## 3. 访问监控界面

### 3.1 Prometheus UI
- **地址**：`http://localhost:9090`
- **功能**：
  - 查询指标（PromQL）
  - 查看目标状态
  - 查看告警规则

### 3.2 Grafana UI
- **地址**：`http://localhost:3001`
- **默认账号**：`admin`
- **默认密码**：`admin`
- **功能**：
  - 仪表板查看
  - 告警管理
  - 用户和权限管理

### 3.3 Alertmanager UI
- **地址**：`http://localhost:9093`
- **功能**：
  - 查看活跃告警
  - 查看告警历史
  - 查看静默规则

---

## 4. 告警配置

### 4.1 关键告警规则

| 告警名称 | 严重级别 | 触发条件 | 说明 |
|---------|---------|---------|------|
| API 服务不可用 | Critical | 服务停止 > 1分钟 | 立即通知 |
| API 响应时间严重 | Critical | P95 > 1秒 | 性能严重下降 |
| API 响应时间过高 | Warning | P95 > 500ms | 性能下降 |
| API 错误率严重 | Critical | 错误率 > 10% | 大量请求失败 |
| API 错误率过高 | Warning | 错误率 > 5% | 请求失败增多 |
| 请求量异常低 | Warning | QPS < 0.1 | 可能有系统问题 |
| 请求量异常高 | Warning | QPS > 100 | 可能是攻击 |
| MySQL 连接池耗尽 | Critical | 连接使用率 > 80% | 数据库连接不足 |
| Redis 连接池耗尽 | Critical | 连接使用率 > 80% | 缓存连接不足 |
| 磁盘空间不足 | Warning | 使用率 > 85% | 需要清理磁盘 |
| CPU 使用率过高 | Warning | 使用率 > 80% | 系统负载高 |
| 内存使用率过高 | Critical | 使用率 > 85% | 内存不足 |

### 4.2 告警通知配置

在 `deployment/monitoring/alertmanager.yml` 中配置：

```yaml
global:
  resolve_timeout: 5m
  smtp_smarthost: 'smtp.example.com:587'
  smtp_from: 'alertmanager@example.com'
  smtp_auth_username: 'alertmanager@example.com'
  smtp_auth_password: 'password'

route:
  group_by: ['alertname', 'severity']
  group_wait: 10s
  group_interval: 10s
  repeat_interval: 12h
  receiver: 'default-receiver'

  routes:
    # Critical 级别告警立即通知
    - match:
        severity: critical
      receiver: 'critical-receiver'
      continue: false

receivers:
  - name: 'default-receiver'
    email_configs:
      - to: 'team@example.com'
        headers:
          Subject: '[DMH 告警] {{ .GroupLabels.alertname }}'

  - name: 'critical-receiver'
    email_configs:
      - to: 'oncall@example.com'
        headers:
          Subject: '🚨 [DMH 紧急] {{ .GroupLabels.alertname }}'
    wechat_configs:
      - corp_id: 'your_corp_id'
        agent_id: 'your_agent_id'
        api_secret: 'your_api_secret'
        to_user: '@all'
```

---

## 5. 性能目标

### 5.1 API 响应时间目标

| 端点类型 | 目标值 (P95) | 当前值 | 状态 |
|----------|--------------|--------|------|
| 海报生成 | < 3秒 | 1.8秒 | ✅ 达标 |
| 二维码生成 | < 500ms | 0.93ms | ✅ 达标 |
| 订单核销 | < 500ms | 0.39ms | ✅ 达标 |
| 普通接口 | < 300ms | 待测 | ⏳ 待测 |

### 5.2 可用性目标

| 服务类型 | 目标值 | 当前值 | 状态 |
|---------|--------|--------|------|
| API 服务 | 99.9% | 待测 | ⏳ 待测 |
| 数据库 | 99.95% | 待测 | ⏳ 待测 |
| Redis 缓存 | 99.9% | 待测 | ⏳ 待测 |

---

## 6. 监控指标说明

### 6.1 API 指标

- `http_requests_total` - HTTP 请求总数（按方法、路径、状态码分组）
- `http_request_duration_seconds` - HTTP 请求耗时（直方图）
- `poster_generation_duration_seconds` - 海报生成耗时
- `qrcode_generation_duration_seconds` - 二维码生成耗时
- `order_verify_duration_seconds` - 订单核销耗时

### 6.2 系统指标

- `node_cpu_seconds_total` - CPU 使用时间
- `node_memory_MemAvailable_bytes` - 可用内存
- `node_filesystem_size_bytes` - 文件系统大小
- `node_network_receive_bytes_total` - 网络接收字节数
- `node_network_transmit_bytes_total` - 网络发送字节数

### 6.3 数据库指标

- `mysql_global_status_threads_connected` - MySQL 连接数
- `mysql_global_status_questions` - 查询总数
- `mysql_global_status_slow_queries` - 慢查询数

---

## 7. 维护操作

### 7.1 日常检查

```bash
# 检查服务状态
docker ps --filter "name=dmh-*"

# 检查日志
docker logs dmh-prometheus --tail 100
docker logs dmh-grafana --tail 100
docker logs dmh-alertmanager --tail 100

# 重启服务
docker restart dmh-prometheus
docker restart dmh-grafana
```

### 7.2 备份配置

```bash
# 备份 Prometheus 配置
cp deployment/monitoring/prometheus.yml deployment/monitoring/backup/prometheus-$(date +%Y%m%d).yml

# 备份告警规则
cp -r deployment/monitoring/alerts deployment/monitoring/backup/alerts-$(date +%Y%m%d)
```

### 7.3 清理旧数据

```bash
# 清理 Prometheus 数据（保留 30 天）
docker exec dmh-prometheus promtool tsdb delete-files 0s 30d

# 清理 Grafana 日志
docker exec dmh-grafana grafana-cli admin clean-old-logs --keep-days 7
```

---

## 8. 故障排查

### 8.1 Prometheus 无法抓取指标

| 症状 | 原因 | 解决方案 |
|------|------|---------|
| Target 显示为 DOWN | 服务不可访问 | 检查网络连接和服务状态 |
| 抓取延迟高 | 目标服务响应慢 | 优化目标服务或增加抓取间隔 |
| 指标为空 | metrics 路径错误 | 验证应用 metrics 端点 |

### 8.2 Grafana 无法连接 Prometheus

| 症状 | 原因 | 解决方案 |
|------|------|---------|
| 数据源连接失败 | Prometheus URL 错误 | 检查数据源配置 |
| 仪表板无数据 | 数据源配置错误或权限问题 | 验证数据源和权限 |

### 8.3 告警未触发

| 症状 | 原因 | 解决方案 |
|------|------|---------|
| 告警规则未生效 | 规则语法错误 | 检查告警规则语法 |
| 告警未发送 | 通知配置错误 | 验证邮件/微信配置 |

---

## 附录

### A. Prometheus 配置模板

见 `deployment/monitoring/prometheus.yml`

### B. 告警规则模板

见 `deployment/monitoring/alerts/*.yml`

### C. Grafana 仪表板

见 `deployment/monitoring/grafana/dashboards/*.json`

---

**文档版本**：1.0  
**最后更新**：2025-02-01
