# Grafana 指标查看指南

Samples Received 只是「收到的样本数」，要看到 CPU、内存、接口 TPS 等，需要新建面板并配置 PromQL。

---

## 一、在 Explore 中查看单指标

1. 左侧 **Explore** → 数据源选 **Grafana Cloud Metrics**
2. 输入下面的 PromQL，点击 **Run query**

---

## 二、常用 PromQL

### CPU 使用率（累计 → 每秒增量）

```promql
rate(process_cpu_seconds_total[1m]) * 100
```

单位：% CPU。折线图即可。

---

### 内存（RSS）

```promql
process_resident_memory_bytes / 1024 / 1024
```

单位：MB。

---

### 接口 TPS（每秒请求数）

```promql
sum(rate(jmeter_api_requests_total[1m])) by (handler)
```

或按 path：

```promql
sum(rate(jmeter_api_requests_total[1m])) by (url)
```

---

### 接口 P99 延迟（秒）

```promql
histogram_quantile(0.99, sum(rate(jmeter_api_request_duration_seconds_bucket[5m])) by (le, url))
```

---

### goroutine 数量

```promql
go_goroutines
```

---

### 缓存命中率（需先有请求）

```promql
jmeter_api_item_count_cache_hits_total / (jmeter_api_item_count_cache_hits_total + jmeter_api_item_count_cache_misses_total) * 100
```

---

## 三、硬盘

当前 `/metrics` **未暴露**磁盘 I/O 指标；`process_*` 只有内存、CPU、文件描述符等。  
如需磁盘监控，可接入 node_exporter 或类似组件单独抓取。

---

## 四、创建看板面板

1. 左侧 **Dashboards** → **New** → **New dashboard**
2. **Add visualization**
3. 数据源选 **Grafana Cloud Metrics**
4. 在 Query 中填入上述 PromQL
5. 选择 Visualization：**Time series**（折线）或 **Stat**（单值）
6. 保存面板

---

## 五、推荐面板配置

| 面板名 | PromQL | 图表类型 |
|--------|--------|----------|
| CPU % | `rate(process_cpu_seconds_total[1m]) * 100` | Time series |
| 内存 MB | `process_resident_memory_bytes / 1024 / 1024` | Time series |
| TPS | `sum(rate(jmeter_api_requests_total[1m]))` | Time series |
| P99 延迟 | `histogram_quantile(0.99, sum(rate(jmeter_api_request_duration_seconds_bucket[5m])) by (le, url))` | Time series |
| goroutine | `go_goroutines` | Stat |
| 错误率 | `sum(rate(jmeter_api_requests_total{code=~"5.."}[5m])) / sum(rate(jmeter_api_requests_total[5m])) * 100` | Stat |

---

## 六、导入现成看板

项目已提供 `docs/grafana-jmeter-api-dashboard.json`：

1. Grafana 左侧 **Dashboards** → **New** → **Import**
2. **Upload JSON file**，选择 `grafana-jmeter-api-dashboard.json`
3. **Prometheus 数据源** 选 **Grafana Cloud Metrics**（或你配置的 Prometheus 数据源）
4. 点击 **Import**

看板包含：CPU、内存、接口 TPS、P99 延迟、goroutine、总 TPS。

若导入后无数据，在每块面板的 **Edit** 中检查数据源 UID 是否与实际一致（可在 **Connections → Data sources** 查看）。

---

## 七、注意事项

1. **job 标签**：若 Grafana Cloud 对 scrape job 加了 `job="xxx"`，PromQL 中需要加：`{job="railway-jmeter-api"}`。
2. **无数据**：检查 Explore 是否有数据；时间范围选 **Last 1 hour** 或更长。
3. **压测才有 TPS**：无请求时 TPS 为 0，可用 JMeter 或 curl 压测后再看。
