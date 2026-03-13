# k3d 集群搭建（Windows）错误复盘

基于对话总结，整理在 Windows 上使用 k3d 搭建本地 Kubernetes 集群时暴露的问题及解决方案。

---

## 错误1：choco 命令未被识别

**背景**：按照方案在 Windows 上安装 k3d，使用 `choco install k3d` 命令。

**错误**：`'choco' is not recognized as an internal or external command, operable program or batch file.`

**原因**：Chocolatey 包管理器未安装在系统上，或安装后未将可执行路径加入环境变量（PATH），导致系统无法找到并执行 `choco` 命令。

**方案**：
1. **优先使用 winget**（Windows 10/11 内置）：
   ```powershell
   winget install k3d.k3d
   ```
2. 使用 Scoop 安装：
   ```powershell
   scoop install k3d
   ```
3. 手动下载：从 [k3d releases](https://github.com/k3d-io/k3d/releases) 下载 `k3d-windows-amd64.exe`，重命名为 `k3d.exe` 并加入 PATH。
4. 若需使用 Chocolatey：先以管理员身份安装 Chocolatey，再执行 `choco install k3d`。

---

## 错误2：kubectl 连接 API 超时

**背景**：k3d 集群已创建（`k3d cluster list` 显示 1/1 servers），context 已切换到 `k3d-mycluster`，执行 `kubectl get nodes` 时失败。

**错误**：
```
couldn't get current server API group list: Get "https://host.docker.internal:55656/api?timeout=32s": 
dial tcp 192.168.71.187:55656: connectex: A connection attempt failed because the connected party 
did not properly respond after a period of time, or established connection failed because connected 
host has failed to respond.
```

**原因**：
1. k3d 生成的 kubeconfig 默认使用 `host.docker.internal` 作为 API 地址；
2. 在 Windows 上，`host.docker.internal` 解析到宿主机 IP（如 192.168.71.187），kubectl 从宿主机发起连接时可能因网络栈、防火墙或 Docker Desktop 端口转发行为导致连接超时；
3. 端口 55656 虽已正确映射（`0.0.0.0:55656->6443/tcp`），但通过 `host.docker.internal` 访问时存在兼容性问题。

**方案**：
1. 将 kubeconfig 中的 server 改为 `127.0.0.1`：
   ```powershell
   kubectl config set-cluster k3d-mycluster --server=https://127.0.0.1:55656
   kubectl get nodes
   ```
2. 若仍失败，检查端口可达性：
   ```powershell
   Test-NetConnection -ComputerName 127.0.0.1 -Port 55656
   ```
3. 若防火墙拦截，添加放行规则（管理员 PowerShell）：
   ```powershell
   netsh advfirewall firewall add rule name="K3d K8s 55656" dir=in action=allow protocol=TCP localport=55656
   ```
4. **预防**：创建集群时使用固定端口，避免随机端口带来的额外复杂度：
   ```powershell
   k3d cluster delete mycluster
   k3d cluster create mycluster -p "6443:6443"
   ```
   此时 kubeconfig 通常使用 6443，连接更稳定。
