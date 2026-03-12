# Go 性能测试练习项目

简单的 Golang REST API 项目，集成 MySQL，提供多种 HTTP 接口供 JMeter 进行性能测试练习。

## 环境准备

- Go 1.21+
- MySQL 8.x（已配置：user=root, pass=root）

## 创建数据库

首次运行前，需要手动创建数据库：

```sql
CREATE DATABASE jmeter_test CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

若出现 `Access denied for user 'root'@'localhost'`，通常由以下原因引起：

1. **密码不正确**：MySQL 中 root 实际密码与代码中的 `root` 不一致  
2. **认证插件**：MySQL 8 默认用 `caching_sha2_password`，部分客户端不兼容  
3. **用户区分**：`root@localhost` 与 `root@127.0.0.1` 被视为不同用户  

**解决办法**（以 root 身份登录 MySQL 后执行）：

```sql
-- 修改 root 用户认证和密码
ALTER USER 'root'@'localhost' IDENTIFIED WITH mysql_native_password BY 'root';
ALTER USER 'root'@'127.0.0.1' IDENTIFIED WITH mysql_native_password BY 'root';
FLUSH PRIVILEGES;
```

**临时方式**：若 MySQL 暂不可用，应用会在连接失败时自动降级为仅指标模式（`/metrics`、`/health`、`/items/slow`），可先验证接口。

## 启动服务

```bash
# 安装依赖（国内网络可设置：go env -w GOPROXY=https://goproxy.cn,direct）
go mod tidy

# 启动
go run main.go
```

服务默认监听 `http://localhost:8080`。

### MySQL 不可用时（SKIP_DB 模式）

若 MySQL 连接失败（如 Access denied），可设置 `SKIP_DB=1` 跳过数据库，仅启动 `/metrics`、`/health`、`/api/items/slow`，用于验证 Prometheus 指标或压测：

```powershell
$env:SKIP_DB="1"
go run main.go
```

### 配置 MySQL 连接

默认使用 `root:root@127.0.0.1:3306/jmeter_test`。若密码不同，可通过环境变量覆盖：

```bash
# Windows PowerShell - 仅修改密码
$env:MYSQL_PASS="你的实际密码"
go run main.go

# 或使用完整连接串
$env:MYSQL_DSN="root:你的密码@tcp(127.0.0.1:3306)/jmeter_test?charset=utf8mb4&parseTime=True"
go run main.go
```

支持的环境变量：`MYSQL_USER`、`MYSQL_PASS`、`MYSQL_HOST`、`MYSQL_PORT`、`MYSQL_DB`、`MYSQL_DSN`、`PORT`

## SQL 优化说明

项目已做以下数据库访问优化：

| 优化项 | 说明 |
|--------|------|
| 连接池 | MaxIdleConns=20, MaxOpenConns=100，高并发下避免连接耗尽 |
| PrepareStmt | GORM 预编译语句复用，减少 SQL 解析开销 |
| UpdateItem | 用 `Updates(map)` 按字段更新，避免 SELECT+Save 全量更新 |
| ListItems | 限制 maxOffset=10000，防止超大 offset 全表扫描；显式 `Order("id")` 利用主键索引 |
| 内存缓存 | 当 items 总数 > 10 万时，缓存 `COUNT(*)` 结果 60 秒；Create/Update/Delete 时自动失效 |

## GitHub Actions 冒烟测试

| Workflow | 用途 | 触发 |
|----------|------|------|
| Smoke Test (CI) | 临时 MySQL + 临时 Go，验证构建与基本连通 | push、手动 |
| Smoke Test (Production) | 针对真实部署环境（Railway） | 手动、每日 8:00 UTC |

**Production 需配置**：仓库 Settings → Secrets → 新增 `API_BASE_URL` = `https://你的应用.up.railway.app`

## API 接口

| 方法   | 路径                     | 描述           |
| ------ | ------------------------ | -------------- |
| GET    | /api/health              | 健康检查       |
| GET    | /api/items               | 列表（支持 page, page_size） |
| GET    | /api/items/:id           | 按 ID 查询     |
| POST   | /api/items               | 创建（JSON: name, description） |
| PUT    | /api/items/:id           | 更新           |
| DELETE | /api/items/:id           | 删除           |
| GET    | /api/items/slow?ms=100   | 模拟延迟（ms 参数单位毫秒） |

## JMeter 配置示例

### 1. 添加线程组

- 右键 Test Plan → Add → Threads (Users) → Thread Group
- 建议：线程数 10，循环次数 100，Ramp-up 时间 1 秒

### 2. 添加 HTTP 请求

- 右键 Thread Group → Add → Sampler → HTTP Request
- Server Name: `localhost`
- Port: `8080`
- Path: 根据测试场景填写，例如：
  - `/api/health` - 基础吞吐量
  - `/api/items` - 列表读
  - `/api/items/1` - 单条查询
  - `/api/items/slow?ms=200` - 响应时间

### 3. POST 创建请求

- Method: POST
- Path: `/api/items`
- Body Data:
  ```json
  {"name": "test-item", "description": "jmeter test"}
  ```
- 添加 HTTP Header Manager：`Content-Type: application/json`

### 4. 添加监听器

- 右键 Thread Group → Add → Listener → Aggregate Report（聚合报告）
- 可再添加 View Results Tree 查看请求详情

### 5. 常用参数说明

| 参数        | 说明                     |
| ----------- | ------------------------ |
| 线程数      | 并发用户数               |
| 循环次数    | 每个用户执行的请求次数   |
| Ramp-up     | 所有线程启动所需时间(秒) |

### 6. 使用现成压测脚本（推荐）

项目已提供现成 JMeter 脚本，可直接导入使用：

1. 启动 JMeter → **File** → **Open** → 选择 `jmeter/load_test.jmx`
2. 确认 Go 服务已启动（`go run main.go`）
3. 点击绿色 **Start** 运行
4. 在 **聚合报告** 查看 TPS、平均响应时间、错误率
5. 在 **查看结果树** 查看每次请求详情

脚本默认配置：10 线程、50 次循环、5 秒 ramp-up，测试 `/api/health`。
# githubActions
