# 监控系统搭建：Grafana + Prometheus 错误复盘

基于对话总结，整理监控系统（Go 应用 /metrics、Grafana Cloud Metrics Endpoint、Railway 部署）搭建过程中暴露的问题及解决方案。

---

## 错误1：MySQL 连接 Access denied

**背景**：本地启动 Go 应用时，需连接 MySQL 数据库。配置为 `root:root@127.0.0.1`。

**错误**：`Error 1045 (28000): Access denied for user 'root'@'localhost' (using password: YES)`，应用无法启动。

**原因**：
1. MySQL 8 默认使用 `caching_sha2_password`，部分 Go 驱动兼容性较差；
2. `root@localhost` 与 `root@127.0.0.1` 被视为不同用户，权限可能不同；
3. 本地 MySQL 实际密码或认证方式与配置不一致。

**方案**：
1. 将连接 host 改为 `localhost`（与 `mysql -u root -proot` 默认行为一致）；
2. 在 MySQL 中执行：
   ```sql
   ALTER USER 'root'@'localhost' IDENTIFIED WITH mysql_native_password BY 'root';
   ALTER USER 'root'@'127.0.0.1' IDENTIFIED WITH mysql_native_password BY 'root';
   FLUSH PRIVILEGES;
   ```
3. 增加「MySQL 连接失败时自动降级」逻辑，不阻塞应用启动，仍可提供 `/metrics`、`/health` 等接口。

---

## 错误2：/metrics 路径重复注册导致 panic

**背景**：在 Gin 中接入 go-gin-prometheus 和 promhttp，期望通过 `/metrics` 暴露 Prometheus 指标。

**错误**：`panic: handlers are already registered for path '/metrics'`。

**原因**：`p.Use(r)` 会同时注册中间件和 `/metrics` 路由，而我们又手动添加了 `r.GET("/metrics", ...)`，导致同一路径被注册两次。

**方案**：只使用 `r.Use(p.HandlerFunc())` 添加中间件，不调用 `p.Use(r)`；由我们单独注册 `/metrics`，并使用 `promhttp.Handler()` 以包含 process、cache、gorm 等全部指标。

---

## 错误3：ProcessCollector 重复注册导致 panic

**背景**：MySQL 连接失败时自动降级，应用继续启动。在降级场景下运行。

**错误**：`panic: duplicate metrics collector registration attempted`。

**原因**：`prometheus.MustRegister(collectors.NewProcessCollector(...))` 在与其他包已注册的 collector 冲突时直接 panic。

**方案**：改用 `prometheus.Register()`，并显式处理 `AlreadyRegisteredError`，遇到该错误时忽略，不 panic。

---

## 错误4：Railway 部署报 secret METRICS_AUTH_PASS: not found

**背景**：为满足 Grafana Cloud Metrics Endpoint 的认证要求，在 Railway 中配置 `METRICS_AUTH_USER` 和 `METRICS_AUTH_PASS`。

**错误**：构建日志中出现 `ERROR: failed to build: failed to solve: secret METRICS_AUTH_PASS: not found`，部署失败。

**原因**：将变量配置为「引用某个 Secret」，而非直接填写值。Railway 在构建时尝试解析该引用，但对应 Secret 不存在。

**方案**：以普通 Variable 形式添加，直接填写字符串值（如 `grafana`、`MyMetrics2024`），不要使用 `${{ secrets.xxx }}` 等引用不存在的 Secret。

---

## 错误5：Grafana Cloud Metrics Endpoint 测试连接失败

**背景**：在 Grafana Cloud 中创建 scrape job，配置 Basic Auth，目标为 Railway 的 `https://xxx.up.railway.app/metrics`。

**错误**：`The request to the provided endpoint with the credentials provided did not succeed`，Test Connection 失败。

**原因**：Grafana 中填写的 Basic 用户名/密码与 Railway 应用环境变量 `METRICS_AUTH_USER`、`METRICS_AUTH_PASS` 不一致（大小写、空格、特殊字符等）。

**方案**：
1. 在 Railway Variables 中确认 `METRICS_AUTH_USER`、`METRICS_AUTH_PASS` 的准确值；
2. 用 `curl -u user:pass https://xxx/metrics` 本地验证；
3. 在 Grafana scrape job 中填入完全相同的用户名和密码，注意不要多打空格。

---

## 错误6：Grafana Explore 查询返回 No data

**背景**：在 Grafana Explore 中查询 `rate(jmeter_api_requests_total[1m])` 或 `rate(process_cpu_seconds_total[1m])*100`。

**错误**：图表显示 `No data`。

**原因**：
1. **TPS 类指标**：无请求时计数器为 0，`rate()` 可能无有效数据；
2. **rate 时间窗口**：抓取间隔为 1 分钟时，`[1m]` 内样本不足，可能导致空结果；
3. **时间范围**：选择的时间范围内没有数据；
4. **未先验证原始指标**：未确认 `jmeter_api_requests_total` 或 `process_cpu_seconds_total` 是否存在。

**方案**：
1. 先用简单指标验证：如 `go_goroutines`、`process_cpu_seconds_total`；
2. 将时间范围调整为 Last 1 hour 或 Last 6 hours；
3. TPS 查询可尝试 `[5m]` 或 `[15m]`；
4. 用 JMeter 或 curl 对接口发一些请求后再查 TPS；
5. 使用 Metrics browser 查看实际存在的指标和标签。

---

## 错误7：process_cpu_seconds_total 有数据，但 rate() 查 CPU% 无数据

**背景**：`process_cpu_seconds_total` 原始指标有数据，但 `rate(process_cpu_seconds_total[1m])*100` 仍显示 No data。

**原因**：`rate()` 的 1 分钟窗口与 scrape 间隔（如 1 分钟）不匹配，可能无法计算出有效速率。

**方案**：改用更长窗口，如 `rate(process_cpu_seconds_total[5m])*100` 或 `rate(process_cpu_seconds_total[15m])*100`。

---

## 经验总结

| 类别 | 要点 |
|------|------|
| MySQL | localhost vs 127.0.0.1、mysql_native_password、失败时自动降级 |
| Prometheus/Gin | go-gin-prometheus 的 Use 会注册路由，避免重复；collector 冲突时用 Register + 忽略 AlreadyRegistered |
| Railway | 变量用直接值，不要引用不存在的 Secret |
| Grafana Cloud | Metrics Endpoint 必须认证；Basic 凭据要与应用端完全一致 |
| PromQL | 先查原始指标，再组合 rate；时间窗口需与 scrape 间隔适配 |
