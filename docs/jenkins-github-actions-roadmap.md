# Jenkins 安装管理及接入 GitHub Actions 路线图

从零搭建 Jenkins，并实现与 GitHub Actions 的混合 CI/CD 流程。

---

## 阶段一：Jenkins 安装

### 1.1 Docker 方式（推荐，跨平台）

```bash
docker run -d \
  --name jenkins \
  -p 8080:8080 \
  -p 50000:50000 \
  -v jenkins_home:/var/jenkins_home \
  jenkins/jenkins:lts
```

**说明**：
- `8080`：Web 管理界面
- `50000`：Agent 通信端口
- `jenkins_home`：数据持久化

### 1.2 获取初始密码

```bash
docker exec jenkins cat /var/jenkins_home/secrets/initialAdminPassword
```

### 1.3 首次访问

1. 浏览器打开 `http://localhost:8080`
2. 输入上述密码
3. 选择「安装推荐的插件」
4. 创建管理员用户

---

## 阶段二：Jenkins 基础配置

### 2.1 启用「触发远程构建」

默认已支持，需为 Job 配置 `token` 参数。

### 2.2 安装常用插件（可选）

- **Git**：拉取代码
- **Docker Pipeline**：构建/推送镜像
- **Kubernetes**：部署到 K8s（若需）
- **Generic Webhook Trigger**：更灵活的自定义 Webhook（可选）

通过 **Manage Jenkins → Plugins** 安装。

---

## 阶段三：创建参数化 Job

### 3.1 新建 Job

1. **New Item**
2. 输入名称：`jmeter-test-deploy`（或自定义）
3. 选择 **Freestyle project**
4. 确定

### 3.2 配置「参数化构建」

1. 勾选 **This project is parameterized**
2. 添加 **String Parameter**：
   - Name: `DOCKER_TAG`
   - Default: `latest`（可选）
3. 添加 **String Parameter**（安全令牌）：
   - Name: `token`
   - Default: 留空，由 GitHub Actions 传入

### 3.3 配置 Token 校验（可选但推荐）

1. **Build Triggers** → 勾选 **Trigger builds remotely**
2. 在 **Authentication Token** 中填写：`your-secret-token`（自行生成，与 GitHub Secrets 一致）

### 3.4 配置构建步骤

**方式 A：仅拉取镜像并重启（无 K8s）**

```bash
# 拉取镜像
docker pull 你的用户名/jmeter-test:${DOCKER_TAG}

# 若本地已有运行中的容器，重启
docker stop jmeter-test 2>/dev/null || true
docker rm jmeter-test 2>/dev/null || true
docker run -d --name jmeter-test -p 8080:8080 你的用户名/jmeter-test:${DOCKER_TAG}
```

**方式 B：部署到 K8s**

```bash
kubectl set image deployment/jmeter-test-api api=你的用户名/jmeter-test:${DOCKER_TAG} -n default
kubectl rollout status deployment/jmeter-test-api -n default
```

### 3.5 记录 Webhook URL

- 格式：`http://你的Jenkins地址/job/jmeter-test-deploy/buildWithParameters`
- 示例：`http://localhost:8080/job/jmeter-test-deploy/buildWithParameters`

若需从外网访问，需用公网 IP 或域名，并配置端口转发/反向代理。

---

## 阶段四：GitHub Actions 集成

### 4.1 配置 GitHub Secrets

仓库 **Settings → Secrets and variables → Actions** 中新增：

| Secret 名称 | 说明 |
|-------------|------|
| `DOCKER_USERNAME` | Docker Hub 用户名 |
| `DOCKER_ACCESS_TOKEN` | Docker Hub Access Token |
| `JENKINS_URL` | 如 `http://your-jenkins:8080` 或 `https://jenkins.example.com` |
| `JENKINS_JOB_PATH` | 如 `/job/jmeter-test-deploy/buildWithParameters` |
| `JENKINS_USERNAME` | Jenkins 用户名（用于 Basic Auth） |
| `JENKINS_PASSWORD` | Jenkins 密码或 API Token |
| `JENKINS_TOKEN` | 与 Job 中 Authentication Token 一致 |

### 4.2 创建 Workflow 文件

在 `.github/workflows/` 下新建 `docker-build-push.yml`：

```yaml
name: Docker Build, Push & Deploy

on:
  push:
    tags:
      - 'jmeter-test/v*'

env:
  APP_NAME: jmeter-test

jobs:
  build-and-push:
    runs-on: ubuntu-latest
    steps:
      - name: Extract version
        run: |
          IMAGE_TAG=$(echo "${GITHUB_REF_NAME}" | sed 's#jmeter-test/##')
          echo "IMAGE_TAG=$IMAGE_TAG" >> $GITHUB_ENV

      - name: Checkout
        uses: actions/checkout@v4

      - name: Login to Docker Hub
        uses: docker/login-action@v3
        with:
          username: ${{ secrets.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKER_ACCESS_TOKEN }}

      - name: Build and push
        run: |
          docker build -t ${{ secrets.DOCKER_USERNAME }}/${{ env.APP_NAME }}:${{ env.IMAGE_TAG }} .
          docker push ${{ secrets.DOCKER_USERNAME }}/${{ env.APP_NAME }}:${{ env.IMAGE_TAG }}

      - name: Trigger Jenkins
        run: |
          curl -u "${{ secrets.JENKINS_USERNAME }}:${{ secrets.JENKINS_PASSWORD }}" \
            -X GET "${{ secrets.JENKINS_URL }}${{ secrets.JENKINS_JOB_PATH }}?token=${{ secrets.JENKINS_TOKEN }}&DOCKER_TAG=${{ env.IMAGE_TAG }}"
```

### 4.3 触发方式

```bash
git tag jmeter-test/v1.0.0
git push origin jmeter-test/v1.0.0
```

---

## 阶段五：端到端验证

| 步骤 | 验证点 |
|------|--------|
| 1 | Jenkins 能访问、登录正常 |
| 2 | Job 手动带参数 `DOCKER_TAG=latest` 构建成功 |
| 3 | 本地 `curl` 触发 Webhook 成功 |
| 4 | push tag 后 GitHub Actions 成功 |
| 5 | Jenkins 收到触发并执行部署 |

---

## 常见问题

### Jenkins 在内网，GitHub Actions 如何触发？

- 使用 **ngrok**、**frp** 等内网穿透，或
- 在公网部署一台 **GitHub Actions Self-Hosted Runner**，Runner 部署在内网或可访问 Jenkins 的网络中，由 Runner 调用 Jenkins

### 如何生成 Jenkins API Token？

**Manage Jenkins → Manage Users → 你的用户 → Configure → Add new Token**。

### Docker Hub 替代方案？

可使用 **GitHub Container Registry (ghcr.io)**，对应修改 login 和镜像地址即可。

---

## 路线图总览

```
Phase 1: 安装 Jenkins (Docker)
    ↓
Phase 2: 初始化、安装插件
    ↓
Phase 3: 创建参数化 Job (DOCKER_TAG, token)
    ↓
Phase 4: 配置 GitHub Secrets + Workflow
    ↓
Phase 5: push tag 触发端到端验证
```
