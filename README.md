# monitoring-stack

## 项目介绍
一个极简版的监控系统，整合了Prometheus（指标）、Promtail（日志采集）、Loki（日志）、Alertmanager（告警）和Grafana（可视化）。

## 启动验证
1. 启动webhook服务

项目地址：https://github.com/Light-Alex/webhook-server
```bash
# webhook-server
go run main.go
```

2. 启动所有监控组件
```bash
# 启动所有组件
docker-compose up -d
```

3. 检查服务状态：
```bash
# 检查所有服务状态
docker-compose ps
```
所有服务的状态应为 Up (healthy)。

4. 检查Loki日志告警是否正确配置
（1）检查Ruler是否正常运行
```bash
# 检查 Loki 容器日志
docker-compose logs loki | grep -i ruler
# 检查 Ruler API 是否可用
curl http://localhost:3100/loki/api/v1/rules
```


（2）验证规则文件是否正确加载
```bash
# 检查规则文件是否被加载
curl http://localhost:3100/loki/api/v1/rules
```

5. 生成一些流量（触发指标和日志）：
```bash
for i in {1..10}; do curl http://localhost:8015/api && sleep 2; done
```

6. 产生错误：
```bash
for i in {1..30}; do curl http://localhost:8015/error && sleep 2; done
```

## 组件功能演示
### Prometheus – 指标查询
- 访问 http://localhost:9090，进入 Graph 页面。
- 输入查询语句，例如 rate(http_requests_total[5m])，查看Go应用的请求速率。

### Loki – 日志查询
- 访问 http://localhost:3000，使用默认账号密码（admin/admin）登录。
- 在 Explore 页面，选择 Loki 数据源。
- 输入查询 {job="go-app"} |= "HTTP request processed"，查看Go应用产生的日志。

### Alertmanager – 告警管理
- 访问 http://localhost:9050，查看告警接收情况。
- 您可以在 config/prometheus/rules.yml 中定义告警规则（例如当请求错误率超过5%时触发告警），Prometheus会将告警发送到Alertmanager。

### Grafana – 统一可视化
- 在Grafana中，您可以将Prometheus的指标和Loki的日志整合到同一个仪表板中。
- 例如，创建一个面板显示请求速率（来自Prometheus），同时在下方表格显示最近的错误日志（来自Loki）。
