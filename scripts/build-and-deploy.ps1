# 本地构建、生成镜像、导入 k3d 并重启应用
# 用法: .\scripts\build-and-deploy.ps1

$ErrorActionPreference = "Stop"
$ImageName = "jmeter-test:latest"
$ClusterName = "mycluster"
$DeploymentName = "jmeter-test-api"

# 重新加载 PATH，确保 k3d、kubectl 可在脚本中找到（IDE 终端有时 PATH 不完整）
$env:Path = [System.Environment]::GetEnvironmentVariable("Path", "Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path", "User")

if (-not (Get-Command k3d -ErrorAction SilentlyContinue)) {
    Write-Host "错误: 未找到 k3d，请在新开终端中重试，或手动执行后续步骤。" -ForegroundColor Red
    Write-Host "  k3d image import $ImageName -c $ClusterName" -ForegroundColor Yellow
    Write-Host "  kubectl rollout restart deployment $DeploymentName" -ForegroundColor Yellow
    exit 1
}

Write-Host "===== Step 1: Docker Build =====" -ForegroundColor Cyan
docker build -t $ImageName .

Write-Host "`n===== Step 2: Import to k3d =====" -ForegroundColor Cyan
& k3d image import $ImageName -c $ClusterName

Write-Host "`n===== Step 3: Restart Deployment =====" -ForegroundColor Cyan
kubectl rollout restart deployment $DeploymentName
kubectl rollout status deployment $DeploymentName

Write-Host "`n===== Done =====" -ForegroundColor Green
kubectl get pods -l app=jmeter-test
