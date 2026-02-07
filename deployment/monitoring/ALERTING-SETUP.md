# DMH 错误告警配置指南

## 概述

本文档说明如何配置 DMH 系统的错误告警，确保及时发现和响应系统异常。

---

## 1. 告警渠道配置

### 1.1 邮件告警

编辑 `deployment/monitoring/alertmanager.yml`：

```yaml
global:
  resolve_timeout: 5m
  smtp_smarthost: 'smtp.example.com:587'
  smtp_from: 'alertmanager@example.com'
  smtp_auth_username: 'alertmanager@example.com'
  smtp_auth_password: 'password'
  smtp_require_tls: true

route:
  group_by: ['alertname', 'severity']
  group_wait: 10s
  group_interval: 10s
  repeat_interval: 12h
  receiver: 'default-receiver'

receivers:
  - name: 'default-receiver'
    email_configs:
      - to: 'team@example.com'
        headers:
          Subject: '[DMH 告警] {{ .GroupLabels.alertname }} - {{ .Status | toUpper }}'
```

### 1.2 微信企业号告警

```yaml
receivers:
  - name: 'wechat-receiver'
    wechat_configs:
      - corp_id: 'wwwwwwwwwwwwwwwwww'
        agent_id: '1000001'
        api_secret: 'xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx'
        to_user: '@all'
        message: |
          {{ range .Alerts }}
          告警: {{ .Labels.alertname }}
          级别: {{ .Labels.severity }}
          描述: {{ .Annotations.description }}
          时间: {{ .StartsAt.Format "2006-01-02 15:04:05" }}
          {{ end }}
```

### 1.3 钉钉告警

```yaml
receivers:
  - name: 'dingtalk-receiver'
    webhook_configs:
      - url: 'https://oapi.dingtalk.com/robot/send?access_token=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx'
        send_resolved: true
```

---

## 2. 告警规则分类

### 2.1 严重级别说明

| 级别 | 响应时间 | 通知方式 | 说明 |
|------|---------|---------|------|
| Critical | 立即 | 电话 + 微信 + 邮件 | 系统不可用或严重影响 |
| Warning | 10分钟内 | 微信 + 邮件 | 性能下降或潜在问题 |
| Info | 1小时内 | 邮件 | 信息性通知 |

### 2.2 告警分组

| 分组 | 说明 | 典型告警 |
|------|------|---------|
| API | API 层面问题 | 服务不可用、响应时间高 |
| Database | 数据库层面问题 | 连接池耗尽、慢查询 |
| Cache | 缓存层面问题 | Redis 连接不足、缓存失效 |
| Infrastructure | 基础设施问题 | CPU/内存/磁盘问题 |

---

## 3. 静默规则配置

### 3.1 维护窗口静默

编辑 `deployment/monitoring/alertmanager.yml`：

```yaml
route:
  group_by: ['alertname', 'severity']
  group_wait: 10s
  group_interval: 10s
  repeat_interval: 12h
  receiver: 'default-receiver'
  
  # 维护窗口路由
  routes:
    - match:
        severity: warning
      receiver: 'maintenance-receiver'
      mute_time_intervals:
        - 'maintenance-window'

receivers:
  - name: 'default-receiver'
    email_configs:
      - to: 'team@example.com'
  
  - name: 'maintenance-receiver'
    # 只发送邮件，不发送微信/钉钉
    email_configs:
      - to: 'maintenance@example.com'

inhibit_rules:
  # 在维护窗口内静默非严重告警
  - source_match:
      severity: 'warning'
    target_match:
      alertname: 'Maintenance'
    equal: ['maintenance']
```

### 3.2 静默时间配置

```yaml
mute_time_intervals:
  - name: 'maintenance-window'
    time_intervals:
      - name: 'daily-maintenance'
        start_time: '00:00'
        end_time: '06:00'
      - name: 'weekend-maintenance'
        times:
          - - start_time: '00:00'
              weekdays: ['Saturday', 'Sunday']
            - start_time: '23:59'
              weekdays: ['Friday']
```

---

## 4. 测试告警

### 4.1 手动触发告警

访问 Prometheus UI 并执行以下 PromQL：

```promql
# 触发测试告警
vector(1)
```

### 4.2 验证告警发送

1. 检查 Alertmanager 日志
```bash
docker logs dmh-alertmanager
```

2. 检查邮件/微信/钉钉是否收到告警

---

## 5. 告警通知模板

### 5.1 邮件模板

```
主题: [{{ .Status | toUpper }}] [DMH] {{ .GroupLabels.alertname }}

{{ range .Alerts }}
告警名称: {{ .Labels.alertname }}
严重级别: {{ .Labels.severity }}
告警组件: {{ .Labels.component }}
发生时间: {{ .StartsAt.Format "2006-01-02 15:04:05" }}
描述: {{ .Annotations.description }}

详细信息:
  - 实例: {{ .Labels.instance }}
  - 值: {{ .Labels.value }}
{{ if .Labels.mountpoint }}
  - 挂载点: {{ .Labels.mountpoint }}
{{ end }}

建议操作:
  - 立即检查服务状态
  - 查看 Grafana 仪表板: http://your-grafana:3001
  - 查看 Prometheus: http://your-prometheus:9090
{{ end }}
```

### 5.2 微信模板

```
{{ range .Alerts }}
🚨 DMH 系统告警

告警: {{ .Labels.alertname }}
级别: {{ .Labels.severity }}
组件: {{ .Labels.component }}
时间: {{ .StartsAt.Format "2006-01-02 15:04:05" }}

{{ .Annotations.description }}

【建议操作】
1. 检查服务状态: docker ps
2. 查看日志: docker logs {{ .Labels.instance }}
3. 查看监控: http://your-grafana:3001
{{ end }}
```

---

## 6. 故障响应流程

### 6.1 Critical 告警响应流程

```
收到 Critical 告警
    ↓
立即响应（5分钟内）
    ↓
确认问题类型
    ├─ 服务不可用 → 检查服务进程
    ├─ 性能问题 → 查看监控数据
    └─ 错误率过高 → 查看应用日志
    ↓
执行修复
    ├─ 重启服务
    ├─ 扩容资源
    └─ 修复代码
    ↓
验证恢复
    ↓
更新知识库
```

### 6.2 Warning 告警响应流程

```
收到 Warning 告警
    ↓
10分钟内响应
    ↓
评估问题严重性
    ├─ 可延迟处理 → 记录到待办事项
    └─ 需要处理 → 按Critical流程处理
    ↓
定期监控
    ↓
问题解决后关闭告警
```

---

## 7. 告警管理

### 7.1 查看活跃告警

访问 Alertmanager UI：`http://localhost:9093`

### 7.2 静默告警

在 Alertmanager UI 中：
1. 点击 "Silence" 按钮
2. 设置静默时间
3. 添加评论说明原因

### 7.3 删除已静默告警

```bash
# 使用 API 删除静默
curl -X POST http://localhost:9093/api/v1/silences/<silence-id>
```

---

## 8. 最佳实践

### 8.1 告警规则设计原则

1. **避免告警风暴**
   - 使用合理的告警阈值
   - 设置合适的告警间隔

2. **告警信息清晰**
   - 提供准确的描述
   - 给出明确的操作建议

3. **合理分组**
   - 相关告警使用同一组
   - 避免重复通知

4. **定期审查**
   - 每月审查告警规则
   - 调整不合理的阈值
   - 删除无用的告警

### 8.2 监控仪表板配置

在 Grafana 中配置以下仪表板：

1. **API 性能仪表板**
   - 请求速率
   - 响应时间
   - 错误率
   - 活跃连接数

2. **系统资源仪表板**
   - CPU 使用率
   - 内存使用率
   - 磁盘空间
   - 网络流量

3. **高级功能仪表板**
   - 海报生成性能
   - 二维码生成性能
   - 订单核销性能

---

## 9. 常见问题

### 9.1 告警重复发送

**问题**：同一告警重复发送多次

**原因**：
- repeat_interval 设置过短
- 告警规则匹配条件过于宽松

**解决方案**：
```yaml
route:
  repeat_interval: 12h  # 增加重复间隔
```

### 9.2 告警未发送

**问题**：告警触发但未收到通知

**原因**：
- 告警通知配置错误
- 网络问题
- 微信/钉钉配置错误

**解决方案**：
1. 检查 Alertmanager 日志
2. 验证邮件/微信/钉钉配置
3. 测试通知渠道

### 9.3 告警误报

**问题**：正常情况下触发告警

**原因**：
- 告警阈值设置过低
- 监控指标异常

**解决方案**：
1. 调整告警阈值
2. 检查监控指标是否正常
3. 添加告警抑制规则

---

## 附录

### A. 告警规则列表

见 `deployment/monitoring/alerts/*.yml`

### B. 告警模板

见 Alertmanager 配置文件

### C. 联系人列表

| 角色 | 姓名 | 邮件 | 微信 | 电话 |
|------|------|------|------|------|
| 技术负责人 | - | - | - | - |
| 运维负责人 | - | - | - | - |
| 产品负责人 | - | - | - | - |

---

**文档版本**：1.0  
**最后更新**：2025-02-01
