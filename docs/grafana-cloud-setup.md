# Grafana Cloud 接入 Railway 指标

两种方式将 `https://githubactions-production-0946.up.railway.app/metrics` 接入 Grafana Cloud。

---

## 方式一：Metrics Endpoint 集成（推荐，零运维）

Grafana Cloud 可直接抓取公网 Prometheus 兼容 URL，无需自建 Prometheus。

### 1. 注册 Grafana Cloud 免费版

- 访问：https://grafana.com/products/cloud/
- 点击 **Start for free**，注册并创建 Stack

### 2. 添加 Metrics Endpoint 集成

1. 登录 Grafana Cloud
2. 左侧菜单点击 **Connections** → **Add new connection**
3. 搜索并选择 **Metrics Endpoint**
4. 点击 **Install** / **Open**

### 3. 创建 Scrape Job

1. 在 **Configuration** 页面点击 **Create scrape job**
2. 填写：
   - **Name**：`railway-jmeter-api`（仅字母、数字、横线、下划线）
   - **URL**：`https://githubactions-production-0946.up.railway.app/metrics`
   - **Scrape interval**：`1m`（默认）
   - **Authentication**：选 **Basic** 或 **Bearer**

**重要**：Metrics Endpoint 要求目标 URL 必须带认证，否则会报错。

**为 /metrics 添加 Basic Auth**：在 Railway 或本地设置环境变量：

| 环境变量 | 说明 |
|----------|------|
| `METRICS_AUTH_USER` | Basic Auth 用户名（如 `grafana`） |
| `METRICS_AUTH_PASS` | Basic Auth 密码 |

两者同时设置后，访问 `/metrics` 需携带 `Authorization: Basic base64(user:pass)` 或浏览器弹出登录框。

### 4. 配置 Grafana 数据源

1. 左侧 **Connections** → **Data sources**
2. 选择 **Grafana Cloud Metrics**（内置，无需手动添加 Prometheus 数据源）
3. 在 Explore 或 Dashboard 中直接查询，数据会来自 Metrics Endpoint 抓取的结果

---

## 方式二：自建 Prometheus + Remote Write（无需对 /metrics 加鉴权）

本地或服务器运行 Prometheus，抓取 Railway 的 `/metrics`，再通过 remote_write 推送到 Grafana Cloud。

### 1. 获取 Grafana Cloud 凭证

1. Grafana Cloud 控制台 → **Cloud Portal**（或 **My Account**）
2. 在 **Prometheus** 卡片点击 **Details**
3. 复制：
   - **Remote Write endpoint**：`https://xxx.grafana.net/api/prom/push`
   - **Username**：Metrics 实例 ID
   - **Password**：Cloud Access Policy token

### 2. 创建 Prometheus 配置

新建 `prometheus.yml`：

```yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'railway-jmeter-api'
    static_configs:
      - targets: ['githubactions-production-0946.up.railway.app']
    scheme: https
    metrics_path: /metrics

remote_write:
  - url: <Your Grafana Cloud remote_write URL>
    basic_auth:
      username: <Your Metrics instance ID>
      password: <Your Cloud Access Policy token>
```

### 3. 运行 Prometheus

```bash
# Docker
docker run -d -p 9090:9090 \
  -v $(pwd)/prometheus.yml:/etc/prometheus/prometheus.yml \
  prom/prometheus:latest

# 或本地安装后
prometheus --config.file=prometheus.yml
```

### 4. 配置 Grafana 数据源

1. Grafana Cloud 中：**Connections** → **Data sources**
2. 选择 **Grafana Cloud Metrics** 作为数据源（内置）
3. 在 Explore / Dashboard 中查询指标（数据来自 remote_write）

---

## 总结

| 方式 | 优点 | 缺点 |
|------|------|------|
| **Metrics Endpoint** | 无需自建 Prometheus，配置简单 | 要求 `/metrics` 必须带认证 |
| **Prometheus + Remote Write** | 不需要对应用做任何改动 | 需要自己跑一个 Prometheus 实例 |

若暂时不想改代码，可先用方式二；若愿意给 `/metrics` 加 Basic Auth，则推荐方式一。
