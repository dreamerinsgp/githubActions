# Docker 镜像构建错误复盘

基于对话总结，整理在本地使用 Docker 构建 jmeter-test 镜像时暴露的问题及解决方案。

---

## 错误1：Dockerfile 不存在

**背景**：在项目根目录执行 `docker build -t jmeter-test:latest .`，准备构建 Go 应用的容器镜像。

**错误**：`ERROR: failed to build: failed to solve: failed to read dockerfile: open Dockerfile: no such file or directory`

**原因**：项目目录下尚未创建 Dockerfile，Docker 无法找到构建定义文件。Go 项目默认没有 Dockerfile，需手动添加。

**方案**：在项目根目录新增 Dockerfile，例如：
```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/main .

FROM alpine:latest
COPY --from=builder /app/main /main
EXPOSE 8080
CMD ["/main"]
```

---

## 错误2：无法拉取 Docker Hub 基础镜像

**背景**：Dockerfile 已创建，执行 `docker build -t jmeter-test:latest .` 时，在拉取 `golang:1.23-alpine` 和 `alpine:latest` 阶段失败。

**错误**：
```
ERROR: failed to build: failed to solve: failed to fetch anonymous token: Get "https://auth.docker.io/token?scope=repository%3Alibrary%2Falpine%3Apull&service=registry.docker.io": 
dial tcp 199.59.148.6:443: connectex: A connection attempt failed because the connected party did not properly respond 
after a period of time, or established connection failed because connected host has failed to respond.
```

**原因**：无法访问 Docker Hub（`auth.docker.io`、`registry.docker.io`），常见于网络限速、防火墙拦截或国内环境访问 Docker Hub 不稳定。

**方案**：
1. **配置镜像加速器**（推荐）：Docker Desktop → Settings → Docker Engine → 在 JSON 中添加 `registry-mirrors`：
   ```json
   {
     "registry-mirrors": [
       "https://docker.m.daocloud.io",
       "https://docker.nju.edu.cn"
     ]
   }
   ```
2. 阿里云用户：在容器镜像服务控制台获取专属加速地址 `https://你的ID.mirror.aliyuncs.com`。
3. 配置后点击 Apply & Restart 重启 Docker，再重新构建。

---

## 错误3：Go 版本与 go.mod 要求不匹配

**背景**：基础镜像已能拉取，构建在 `RUN go mod download` 阶段失败。

**错误**：`go: go.mod requires go >= 1.24.0 (running go 1.23.12; GOTOOLCHAIN=local)`

**原因**：Dockerfile 使用 `golang:1.23-alpine`，内置 Go 1.23.12，而项目的 `go.mod` 声明了 `go 1.24.0`，版本不满足要求，`go mod download` 拒绝执行。

**方案**：将 Dockerfile 中构建阶段的基础镜像升级为 `golang:1.24-alpine`：
```dockerfile
FROM golang:1.24-alpine AS builder
```
若 `golang:1.24-alpine` 不可用，可改为将 `go.mod` 中的最低版本降为 `go 1.23`，继续使用 `golang:1.23-alpine`。
