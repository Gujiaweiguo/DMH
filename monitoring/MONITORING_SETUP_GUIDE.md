# DMH 高级功能监控和告警配置指南

## 📋 概述

本文档提供了 DMH 高级功能（海报生成、支付配置、表单增强、订单核销）的监控和告警配置方案。

**相关文档**:
- [性能测试报告](../backend/test/performance/PERFORMANCE_TEST_REPORT.md)
- [生产部署指南](../deployment/PRODUCTION_DEPLOYMENT_GUIDE.md)
- [任务清单](../openspec/changes/add-campaign-advanced-features/tasks.md)

---

## 🎯 监控目标

### 关键指标

| 监控项 | 目标值 | 告警阈值 | 说明 |
|--------|--------|----------|------|
| **海报生成时间** | < 3 秒 | > 5 秒 | 单次海报生成响应时间 |
| **二维码生成时间** | < 500ms | > 1 秒 | 支付二维码生成响应时间 |
| **核销接口响应时间** | < 500ms | > 1 秒 | 订单核销接口响应时间 |
| **海报生成成功率** | > 95% | < 90% | 海报生成请求的成功率 |
| **API 错误率** | < 0.1% | > 1% | 所有 API 的错误率 |
| **数据库连接数** | < 80% | > 90% | 数据库连接池使用率 |
| **Redis 缓存命中率** | > 90% | < 80% | Redis 缓存命中率 |
| **API QPS** | - | > 1000/s | API 每秒请求数 |
| **内存使用** | < 80% | > 90% | 应用内存使用率 |
| **CPU 使用率** | < 70% | > 85% | 应用 CPU 使用率 |
| **磁盘使用率** | < 80% | > 90% | 磁盘空间使用率 |

---

## 📊 监控系统架构

### 推荐监控栈

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  应用层     │────▶│  采集层     │────▶│  存储层     │
│  (DMH API)  │     │ (Prometheus) │     │(Prometheus) │
└─────────────┘     └─────────────┘     └─────────────┘
                                               │
                                               ▼
                                        ┌─────────────┐
                                        │  可视化层   │
                                        │ (Grafana)   │
                                        └─────────────┘
                                               │
                                               ▼
                                        ┌─────────────┐
                                        │  告警层     │
                                        │ (AlertManager)│
                                        └─────────────┘
```

### 系统组件

1. **Prometheus**: 时序数据库，负责指标采集和存储
2. **Grafana**: 可视化仪表盘，展示监控数据
3. **AlertManager**: 告警管理，处理告警规则和通知

---

## 🔧 第一部分：应用层监控配置

### 1.1 Prometheus 客户端集成

在 Go 后端应用中集成 Prometheus 客户端：

#### 1.1.1 添加依赖

```bash
cd backend
go get github.com/prometheus/client_golang/prometheus
go get github.com/prometheus/client_golang/prometheus/promhttp
```

#### 1.1.2 创建监控中间件

创建文件 `backend/api/internal/middleware/prometheus_middleware.go`:

```go
package middleware

import (
    "strconv"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    // HTTP 请求总数
    httpRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "dmh_http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "path", "status"},
    )

    // HTTP 请求耗时
    httpRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "dmh_http_request_duration_seconds",
            Help:    "HTTP request latency in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "path"},
    )

    // 海报生成耗时
    posterGenerationDuration = promauto.NewHistogram(
        prometheus.HistogramOpts{
            Name:    "dmh_poster_generation_duration_seconds",
            Help:    "Poster generation latency in seconds",
            Buckets: []float64{1, 2, 3, 5, 10},
        },
    )

    // 二维码生成耗时
    qrcodeGenerationDuration = promauto.NewHistogram(
        prometheus.HistogramOpts{
            Name:    "dmh_qrcode_generation_duration_seconds",
            Help:    "QRCode generation latency in seconds",
            Buckets: []float64{0.1, 0.2, 0.5, 1, 2},
        },
    )

    // 核销操作耗时
    orderVerifyDuration = promauto.NewHistogram(
        prometheus.HistogramOpts{
            Name:    "dmh_order_verify_duration_seconds",
            Help:    "Order verification latency in seconds",
            Buckets: []float64{0.1, 0.2, 0.5, 1, 2},
        },
    )

    // 海报生成总数
    posterGenerationTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "dmh_poster_generation_total",
            Help: "Total number of poster generations",
        },
        []string{"status"},
    )
)

// PrometheusMiddleware HTTP 请求监控中间件
func PrometheusMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()

        // 处理请求
        c.Next()

        // 记录指标
        duration := time.Since(start).Seconds()
        status := strconv.Itoa(c.Writer.Status())

        httpRequestsTotal.WithLabelValues(c.Request.Method, c.FullPath(), status).Inc()
        httpRequestDuration.WithLabelValues(c.Request.Method, c.FullPath()).Observe(duration)
    }
}

// RecordPosterGeneration 记录海报生成指标
func RecordPosterGeneration(duration time.Duration, status string) {
    posterGenerationDuration.Observe(duration.Seconds())
    posterGenerationTotal.WithLabelValues(status).Inc()
}

// RecordQRCodeGeneration 记录二维码生成指标
func RecordQRCodeGeneration(duration time.Duration) {
    qrcodeGenerationDuration.Observe(duration.Seconds())
}

// RecordOrderVerify 记录核销操作指标
func RecordOrderVerify(duration time.Duration) {
    orderVerifyDuration.Observe(duration.Seconds())
}
```

#### 1.1.3 在路由中集成监控中间件

修改 `backend/api/internal/handler/routes.go`:

```go
import (
    "dmh/api/internal/middleware"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

// 在路由注册时添加 Prometheus 中间件
func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
    // 添加监控中间件
    server.Use(middleware.PrometheusMiddleware())

    // 注册指标端点
    server.AddRoute(
        rest.Route{
            Method: http.MethodGet,
            Path:   "/metrics",
            Handler: gin.WrapH(promhttp.Handler()),
        },
    )

    // ... 其他路由配置
}
```

#### 1.1.4 在业务逻辑中记录指标

**海报生成逻辑修改**:

```go
// 在 backend/api/internal/logic/poster/generateCampaignPosterLogic.go 中
func (l *GenerateCampaignPosterLogic) GenerateCampaignPoster(req *types.GeneratePosterReq) (resp *types.GeneratePosterResp, err error) {
    startTime := time.Now()
    status := "success"

    defer func() {
        // 记录监控指标
        if err != nil {
            status = "error"
        }
        middleware.RecordPosterGeneration(time.Since(startTime), status)
    }()

    // ... 原有逻辑
}
```

**二维码生成逻辑修改**:

```go
// 在 backend/api/internal/logic/campaign/getPaymentQrcodeLogic.go 中
func (l *GetPaymentQrcodeLogic) GetPaymentQrcode(req *types.GetPaymentQrcodeReq) (resp *types.GetPaymentQrcodeResp, err error) {
    startTime := time.Now()

    defer func() {
        // 记录监控指标
        middleware.RecordQRCodeGeneration(time.Since(startTime))
    }()

    // ... 原有逻辑
}
```

**核销逻辑修改**:

```go
// 在 backend/api/internal/logic/order/verifyOrderLogic.go 中
func (l *VerifyOrderLogic) VerifyOrder(req *types.VerifyOrderReq) (resp *types.VerifyOrderResp, err error) {
    startTime := time.Now()

    defer func() {
        // 记录监控指标
        middleware.RecordOrderVerify(time.Since(startTime))
    }()

    // ... 原有逻辑
}
```

### 1.2 配置 Prometheus 端口

在配置文件 `backend/api/etc/dmh-api.prod.yaml` 中：

```yaml
# Prometheus 配置
Prometheus:
  Enabled: true
  Port: 9090  # 指标暴露端口
  Path: "/metrics"
```

---

## 🖥️ 第二部分：Prometheus 服务配置

### 2.1 安装 Prometheus

```bash
# 创建用户
sudo useradd --no-create-home --shell /bin/false prometheus

# 创建目录
sudo mkdir -p /etc/prometheus
sudo mkdir -p /var/lib/prometheus

# 下载 Prometheus
cd /tmp
wget https://github.com/prometheus/prometheus/releases/download/v2.45.0/prometheus-2.45.0.linux-amd64.tar.gz

# 解压
tar xvfz prometheus-2.45.0.linux-amd64.tar.gz
cd prometheus-2.45.0.linux-amd64

# 安装
sudo cp prometheus promtool /usr/local/bin/
sudo chown prometheus:prometheus /usr/local/bin/prometheus
sudo chown prometheus:prometheus /usr/local/bin/promtool

# 安装配置文件
sudo cp consoles/ console_libraries/ prometheus.yml /etc/prometheus/
sudo chown -R prometheus:prometheus /etc/prometheus
sudo chown -R prometheus:prometheus /var/lib/prometheus
```

### 2.2 配置 Prometheus

编辑 `/etc/prometheus/prometheus.yml`:

```yaml
# 全局配置
global:
  scrape_interval: 15s
  evaluation_interval: 15s
  external_labels:
    monitor: 'dmh-monitor'

# 告警管理器配置
alerting:
  alertmanagers:
    - static_configs:
        - targets: ['localhost:9093']

# 告警规则文件
rule_files:
  - '/etc/prometheus/rules/*.yml'

# 数据采集配置
scrape_configs:
  # DMH API 监控
  - job_name: 'dmh-api'
    static_configs:
      - targets: ['localhost:9090']
    metrics_path: '/metrics'
    scrape_interval: 10s

  # Node Exporter 监控（系统指标）
  - job_name: 'node-exporter'
    static_configs:
      - targets: ['localhost:9100']

  # MySQL Exporter 监控
  - job_name: 'mysql-exporter'
    static_configs:
      - targets: ['localhost:9104']

  # Redis Exporter 监控
  - job_name: 'redis-exporter'
    static_configs:
      - targets: ['localhost:9121']
```

### 2.3 配置告警规则

创建 `/etc/prometheus/rules/dmh-alerts.yml`:

```yaml
groups:
  - name: dmh_performance
    interval: 30s
    rules:
      # 海报生成时间告警
      - alert: PosterGenerationSlow
        expr: histogram_quantile(0.95, dmh_poster_generation_duration_seconds_bucket) > 5
        for: 5m
        labels:
          severity: warning
          component: poster
        annotations:
          summary: "海报生成时间过长"
          description: "95% 的海报生成请求耗时超过 5 秒（当前值: {{ $value }}s）"

      # 二维码生成时间告警
      - alert: QRCodeGenerationSlow
        expr: histogram_quantile(0.95, dmh_qrcode_generation_duration_seconds_bucket) > 1
        for: 5m
        labels:
          severity: warning
          component: qrcode
        annotations:
          summary: "二维码生成时间过长"
          description: "95% 的二维码生成请求耗时超过 1 秒（当前值: {{ $value }}s）"

      # 核销接口响应时间告警
      - alert: OrderVerifySlow
        expr: histogram_quantile(0.95, dmh_order_verify_duration_seconds_bucket) > 1
        for: 5m
        labels:
          severity: warning
          component: order
        annotations:
          summary: "核销接口响应时间过长"
          description: "95% 的核销请求耗时超过 1 秒（当前值: {{ $value }}s）"

      # 海报生成成功率告警
      - alert: PosterGenerationFailureRateHigh
        expr: |
          (
            sum(rate(dmh_poster_generation_total{status="error"}[5m]))
            /
            sum(rate(dmh_poster_generation_total[5m]))
          ) > 0.1
        for: 10m
        labels:
          severity: critical
          component: poster
        annotations:
          summary: "海报生成失败率过高"
          description: "海报生成失败率超过 10%（当前值: {{ $value | humanizePercentage }}）"

      # API 错误率告警
      - alert: APIErrorRateHigh
        expr: |
          (
            sum(rate(dmh_http_requests_total{status=~"5.."}[5m]))
            /
            sum(rate(dmh_http_requests_total[5m]))
          ) > 0.01
        for: 5m
        labels:
          severity: critical
          component: api
        annotations:
          summary: "API 错误率过高"
          description: "API 5xx 错误率超过 1%（当前值: {{ $value | humanizePercentage }}）"

  - name: system_health
    interval: 30s
    rules:
      # CPU 使用率告警
      - alert: HighCPUUsage
        expr: 100 - (avg by(instance) (rate(node_cpu_seconds_total{mode="idle"}[5m])) * 100) > 85
        for: 5m
        labels:
          severity: warning
          component: system
        annotations:
          summary: "CPU 使用率过高"
          description: "服务器 CPU 使用率超过 85%（当前值: {{ $value }}%）"

      # 内存使用率告警
      - alert: HighMemoryUsage
        expr: (1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)) * 100 > 90
        for: 5m
        labels:
          severity: critical
          component: system
        annotations:
          summary: "内存使用率过高"
          description: "服务器内存使用率超过 90%（当前值: {{ $value }}%）"

      # 磁盘使用率告警
      - alert: HighDiskUsage
        expr: (node_filesystem_avail_bytes{mountpoint="/"} / node_filesystem_size_bytes{mountpoint="/"}) * 100 < 10
        for: 5m
        labels:
          severity: warning
          component: system
        annotations:
          summary: "磁盘空间不足"
          description: "根分区剩余空间低于 10%（当前值: {{ $value }}%）"

      # 数据库连接数告警
      - alert: DatabaseConnectionsHigh
        expr: mysql_global_status_threads_connected / mysql_global_variables_max_connections * 100 > 90
        for: 5m
        labels:
          severity: warning
          component: database
        annotations:
          summary: "数据库连接数过高"
          description: "MySQL 连接数超过 90%（当前值: {{ $value }}%）"
```

### 2.4 创建 systemd 服务

创建 `/etc/systemd/system/prometheus.service`:

```ini
[Unit]
Description=Prometheus Monitoring System
After=network.target

[Service]
Type=simple
User=prometheus
Group=prometheus
WorkingDirectory=/etc/prometheus
ExecStart=/usr/local/bin/prometheus \
  --config.file=/etc/prometheus/prometheus.yml \
  --storage.tsdb.path=/var/lib/prometheus \
  --web.console.templates=/etc/prometheus/consoles \
  --web.console.libraries=/etc/prometheus/console_libraries \
  --web.listen-address=:9091 \
  --web.external-url=http://localhost:9091
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

**启动 Prometheus**:
```bash
# 重载 systemd
sudo systemctl daemon-reload

# 启用服务
sudo systemctl enable prometheus

# 启动服务
sudo systemctl start prometheus

# 查看状态
sudo systemctl status prometheus

# 访问 Prometheus UI
# http://localhost:9091
```

---

## 📈 第三部分：Grafana 配置

### 3.1 安装 Grafana

```bash
# 添加 Grafana APT 仓库
wget -q -O - https://packages.grafana.com/gpg.key | sudo apt-key add -
sudo add-apt-repository "deb https://packages.grafana.com/oss/deb stable main"

# 安装 Grafana
sudo apt-get update
sudo apt-get install -y grafana

# 启动 Grafana
sudo systemctl enable grafana-server
sudo systemctl start grafana-server

# 访问 Grafana
# http://localhost:3000
# 默认用户名/密码: admin/admin
```

### 3.2 配置 Prometheus 数据源

1. 登录 Grafana (http://localhost:3000)
2. 导航到 **Configuration** > **Data Sources**
3. 点击 **Add data source**
4. 选择 **Prometheus**
5. 配置数据源:
   - **Name**: Prometheus
   - **URL**: http://localhost:9091
   - **Access**: Server (default)
6. 点击 **Save & Test**

### 3.3 导入 DMH 监控仪表盘

#### 3.3.1 创建仪表盘 JSON

创建文件 `/opt/code/DMH/monitoring/grafana-dmh-dashboard.json`:

```json
{
  "dashboard": {
    "title": "DMH 高级功能监控",
    "uid": "dmh-advanced-features",
    "tags": ["dmh", "performance"],
    "timezone": "browser",
    "refresh": "30s",
    "panels": [
      {
        "id": 1,
        "title": "海报生成请求 QPS",
        "type": "graph",
        "targets": [
          {
            "expr": "sum(rate(dmh_poster_generation_total[1m]))",
            "legendFormat": "QPS"
          }
        ],
        "yaxes": [
          {
            "format": "ops"
          }
        ]
      },
      {
        "id": 2,
        "title": "海报生成 P95 延迟",
        "type": "graph",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, dmh_poster_generation_duration_seconds_bucket)",
            "legendFormat": "P95"
          }
        ],
        "yaxes": [
          {
            "format": "s"
          }
        ],
        "alert": {
          "conditions": [
            {
              "evaluator": {
                "params": [5],
                "type": "gt"
              },
              "operator": {
                "type": "and"
              },
              "query": {
                "params": ["A", "5m", "now"]
              },
              "reducer": {
                "params": [],
                "type": "avg"
              },
              "type": "query"
            }
          ]
        }
      },
      {
        "id": 3,
        "title": "二维码生成 P95 延迟",
        "type": "graph",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, dmh_qrcode_generation_duration_seconds_bucket)",
            "legendFormat": "P95"
          }
        ],
        "yaxes": [
          {
            "format": "s"
          }
        ]
      },
      {
        "id": 4,
        "title": "核销接口 P95 延迟",
        "type": "graph",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, dmh_order_verify_duration_seconds_bucket)",
            "legendFormat": "P95"
          }
        ],
        "yaxes": [
          {
            "format": "s"
          }
        ]
      },
      {
        "id": 5,
        "title": "API 错误率",
        "type": "graph",
        "targets": [
          {
            "expr": "sum(rate(dmh_http_requests_total{status=~\"5..\"}[5m])) / sum(rate(dmh_http_requests_total[5m])) * 100",
            "legendFormat": "5xx 错误率 %"
          }
        ],
        "yaxes": [
          {
            "format": "percent"
          }
        ]
      },
      {
        "id": 6,
        "title": "系统 CPU 使用率",
        "type": "graph",
        "targets": [
          {
            "expr": "100 - (avg by(instance) (rate(node_cpu_seconds_total{mode=\"idle\"}[5m])) * 100)",
            "legendFormat": "{{instance}}"
          }
        ],
        "yaxes": [
          {
            "format": "percent"
          }
        ]
      }
    ]
  }
}
```

#### 3.3.2 导入仪表盘

1. 导航到 **Dashboards** > **Import**
2. 点击 **Upload JSON file**
3. 选择 `grafana-dmh-dashboard.json`
4. 点击 **Load**
5. 点击 **Import**

---

## 🚨 第四部分：AlertManager 配置

### 4.1 安装 AlertManager

```bash
# 创建用户
sudo useradd --no-create-home --shell /bin/false alertmanager

# 创建目录
sudo mkdir -p /etc/alertmanager
sudo mkdir -p /var/lib/alertmanager

# 下载 AlertManager
cd /tmp
wget https://github.com/prometheus/alertmanager/releases/download/v0.26.0/alertmanager-0.26.0.linux-amd64.tar.gz

# 解压
tar xvfz alertmanager-0.26.0.linux-amd64.tar.gz
cd alertmanager-0.26.0.linux-amd64

# 安装
sudo cp alertmanager amtool /usr/local/bin/
sudo chown alertmanager:alertmanager /usr/local/bin/alertmanager
sudo chown alertmanager:alertmanager /usr/local/bin/amtool

# 安装配置文件
sudo cp alertmanager.yml /etc/alertmanager/
sudo chown -R alertmanager:alertmanager /etc/alertmanager
sudo chown -R alertmanager:alertmanager /var/lib/alertmanager
```

### 4.2 配置 AlertManager

编辑 `/etc/alertmanager/alertmanager.yml`:

```yaml
global:
  resolve_timeout: 5m

# 路由配置
route:
  group_by: ['alertname', 'severity']
  group_wait: 10s
  group_interval: 10s
  repeat_interval: 12h
  receiver: 'default'

  # 子路由
  routes:
    # Critical 告警 - 发送给运维团队
    - match:
        severity: critical
      receiver: 'critical-alerts'
      continue: true

    # Warning 告警 - 发送给开发团队
    - match:
        severity: warning
      receiver: 'warning-alerts'

# 接收器配置
receivers:
  # 默认接收器
  - name: 'default'
    email_configs:
      - to: 'ops@dmh.com'
        from: 'alertmanager@dmh.com'
        smarthost: 'smtp.example.com:587'
        auth_username: 'alertmanager@dmh.com'
        auth_password: 'your_password'
        headers:
          Subject: 'DMH 告警通知'

  # Critical 告警接收器
  - name: 'critical-alerts'
    email_configs:
      - to: 'ops@dmh.com,devops@dmh.com'
        from: 'alertmanager@dmh.com'
        smarthost: 'smtp.example.com:587'
        auth_username: 'alertmanager@dmh.com'
        auth_password: 'your_password'
        headers:
          Subject: '[CRITICAL] DMH 告警通知'

    webhook_configs:
      - url: 'https://hooks.slack.com/services/YOUR/WEBHOOK/URL'
        send_resolved: true

  # Warning 告警接收器
  - name: 'warning-alerts'
    email_configs:
      - to: 'dev@dmh.com'
        from: 'alertmanager@dmh.com'
        smarthost: 'smtp.example.com:587'
        auth_username: 'alertmanager@dmh.com'
        auth_password: 'your_password'
        headers:
          Subject: '[WARNING] DMH 告警通知'

# 告警抑制规则
inhibit_rules:
  # 如果 Critical 告警触发，抑制 Warning 告警
  - source_match:
      severity: 'critical'
    target_match:
      severity: 'warning'
    equal: ['alertname', 'instance']
```

### 4.3 创建 systemd 服务

创建 `/etc/systemd/system/alertmanager.service`:

```ini
[Unit]
Description=Alertmanager
After=network.target

[Service]
Type=simple
User=alertmanager
Group=alertmanager
WorkingDirectory=/etc/alertmanager
ExecStart=/usr/local/bin/alertmanager \
  --config.file=/etc/alertmanager/alertmanager.yml \
  --storage.path=/var/lib/alertmanager \
  --web.listen-address=:9093 \
  --web.external-url=http://localhost:9093
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

**启动 AlertManager**:
```bash
# 重载 systemd
sudo systemctl daemon-reload

# 启用服务
sudo systemctl enable alertmanager

# 启动服务
sudo systemctl start alertmanager

# 查看状态
sudo systemctl status alertmanager

# 访问 AlertManager UI
# http://localhost:9093
```

---

## 📧 第五部分：告警通知配置

### 5.1 邮件通知

修改 `/etc/alertmanager/alertmanager.yml` 中的邮件配置：

```yaml
email_configs:
  - to: 'ops@dmh.com'
    from: 'alertmanager@dmh.com'
    smarthost: 'smtp.gmail.com:587'
    auth_username: 'alertmanager@dmh.com'
    auth_password: 'your_app_password'  # 使用应用专用密码
    require_tls: true
    headers:
      Subject: '[{{ .Status | toUpper }}] {{ .CommonLabels.alertname }} - DMH 告警'
      X-Priority: '1'
```

### 5.2 Slack 通知

创建 Slack Incoming Webhook:

1. 访问 https://api.slack.com/apps
2. 创建 Incoming Webhook 应用
3. 获取 Webhook URL
4. 配置到 AlertManager:

```yaml
webhook_configs:
  - url: 'https://hooks.slack.com/services/YOUR/WEBHOOK/URL'
    send_resolved: true
```

### 5.3 钉钉通知（可选）

如果使用钉钉，可以使用 webhook 转发工具：

```bash
# 安装钉钉通知工具
git clone https://github.com/timonwong/prometheus-webhook-dingtalk.git
cd prometheus-webhook-dingtalk
go build

# 运行服务
./prometheus-webhook-dingtalk \
  --ding.profile=dmh \
  --ding.url=https://oapi.dingtalk.com/robot/send?access_token=YOUR_TOKEN \
  --ding.secret=YOUR_SECRET
```

配置 AlertManager:

```yaml
webhook_configs:
  - url: 'http://localhost:8060/dingtalk/dmh/profile'
    send_resolved: true
```

---

## ✅ 第六部分：监控验证

### 6.1 验证 Prometheus

```bash
# 检查 Prometheus 状态
curl http://localhost:9091/-/healthy

# 查询指标
curl 'http://localhost:9091/api/v1/query?query=up'

# 检查告警规则
curl 'http://localhost:9091/api/v1/rules'
```

### 6.2 验证 AlertManager

```bash
# 检查 AlertManager 状态
curl http://localhost:9093/-/healthy

# 查看告警
curl http://localhost:9093/api/v1/alerts

# 测试告警通知
curl -XPOST http://localhost:9093/api/v1/alerts \
  -d '[{
    "labels": {
      "alertname": "TestAlert",
      "severity": "warning"
    },
    "annotations": {
      "description": "This is a test alert"
    }
  }]'
```

### 6.3 验证 Grafana

1. 访问 Grafana Dashboard
2. 检查各面板数据是否正常显示
3. 确认数据刷新是否正常
4. 验证告警面板是否显示告警

---

## 📋 监控检查清单

- [ ] Prometheus 已安装并启动
- [ ] AlertManager 已安装并启动
- [ ] Grafana 已安装并启动
- [ ] Prometheus 数据源已配置
- [ ] 告警规则已创建
- [ ] 通知渠道已配置
- [ ] Grafana 仪表盘已导入
- [ ] 监控指标已正常采集
- [ ] 告警通知测试通过
- [ ] 监控访问权限已配置

---

## 📞 故障排查

### Prometheus 无法采集指标

**检查项**:
```bash
# 检查目标是否可达
curl http://localhost:9090/metrics

# 检查 Prometheus 配置
promtool check config /etc/prometheus/prometheus.yml

# 查看 Prometheus 日志
sudo journalctl -u prometheus -n 50
```

### 告警未触发

**检查项**:
```bash
# 检查告警规则
curl 'http://localhost:9091/api/v1/rules' | jq

# 检查告警状态
curl 'http://localhost:9091/api/v1/alerts' | jq

# 验证告警表达式
# 在 Prometheus UI 中测试表达式
```

### 通知未收到

**检查项**:
```bash
# 查看 AlertManager 日志
sudo journalctl -u alertmanager -n 50

# 检查 AlertManager 配置
amtool config validate /etc/alertmanager/alertmanager.yml

# 测试 SMTP 连接
telnet smtp.example.com 587
```

---

**配置状态**: 待执行
**最后更新**: 2026-02-01
**维护人员**: DevOps 团队
