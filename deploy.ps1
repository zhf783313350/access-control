# 生产环境部署脚本 (Windows PowerShell)

param(
    [Parameter(Position=0)]
    [string]$Command = "help"
)

$ProjectName = "access-control"
$ComposeFile = "docker-compose.prod.yml"

function Write-ColorOutput($ForegroundColor) {
    $fc = $host.UI.RawUI.ForegroundColor
    $host.UI.RawUI.ForegroundColor = $ForegroundColor
    if ($args) {
        Write-Output $args
    }
    $host.UI.RawUI.ForegroundColor = $fc
}

function Show-Banner {
    Write-Host "========================================" -ForegroundColor Green
    Write-Host "  Access Control 生产环境部署脚本" -ForegroundColor Green
    Write-Host "========================================" -ForegroundColor Green
}

function Show-Help {
    Write-Host "用法: .\deploy.ps1 [命令]"
    Write-Host ""
    Write-Host "命令:"
    Write-Host "  start       启动所有服务"
    Write-Host "  stop        停止所有服务"
    Write-Host "  restart     重启所有服务"
    Write-Host "  update      更新并重启服务"
    Write-Host "  logs        查看日志"
    Write-Host "  status      查看服务状态"
    Write-Host "  backup      备份数据库"
    Write-Host "  clean       清理未使用的资源"
    Write-Host "  help        显示此帮助"
}

function Start-Services {
    Write-Host "正在启动服务..." -ForegroundColor Green
    docker-compose -f $ComposeFile up -d
    Write-Host "服务启动完成" -ForegroundColor Green
    Show-Status
}

function Stop-Services {
    Write-Host "正在停止服务..." -ForegroundColor Yellow
    docker-compose -f $ComposeFile down
    Write-Host "服务已停止" -ForegroundColor Green
}

function Restart-Services {
    Write-Host "正在重启服务..." -ForegroundColor Yellow
    docker-compose -f $ComposeFile restart
    Write-Host "服务重启完成" -ForegroundColor Green
}

function Update-Services {
    Write-Host "正在更新服务..." -ForegroundColor Green
    
    # 拉取最新镜像
    docker-compose -f $ComposeFile pull
    
    # 重新构建并启动
    docker-compose -f $ComposeFile up -d --build
    
    # 清理旧镜像
    docker image prune -f
    
    Write-Host "服务更新完成" -ForegroundColor Green
}

function Show-Logs {
    docker-compose -f $ComposeFile logs -f --tail=100
}

function Show-Status {
    Write-Host "服务状态:" -ForegroundColor Green
    docker-compose -f $ComposeFile ps
    Write-Host ""
    Write-Host "资源使用:" -ForegroundColor Green
    docker stats --no-stream --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}"
}

function Backup-Database {
    $BackupDir = ".\backups"
    $Timestamp = Get-Date -Format "yyyyMMdd_HHmmss"
    $BackupFile = "$BackupDir\backup_$Timestamp.sql"
    
    if (!(Test-Path $BackupDir)) {
        New-Item -ItemType Directory -Path $BackupDir | Out-Null
    }
    
    Write-Host "正在备份数据库..." -ForegroundColor Green
    
    # 从 .env 读取数据库配置
    $env:DB_USER = "postgres"
    $env:DB_NAME = "access_control"
    
    docker-compose -f $ComposeFile exec -T postgres pg_dump -U $env:DB_USER $env:DB_NAME > $BackupFile
    
    Write-Host "数据库已备份到: $BackupFile" -ForegroundColor Green
}

function Clear-Resources {
    Write-Host "正在清理未使用的 Docker 资源..." -ForegroundColor Yellow
    docker system prune -f
    docker volume prune -f
    Write-Host "清理完成" -ForegroundColor Green
}

# 检查 .env 文件
if (!(Test-Path ".env")) {
    Write-Host "警告: .env 文件不存在，将从 .env.example 复制" -ForegroundColor Yellow
    Copy-Item ".env.example" ".env"
    Write-Host "请修改 .env 文件中的配置后重新运行此脚本" -ForegroundColor Red
    exit 1
}

Show-Banner

switch ($Command.ToLower()) {
    "start" { Start-Services }
    "stop" { Stop-Services }
    "restart" { Restart-Services }
    "update" { Update-Services }
    "logs" { Show-Logs }
    "status" { Show-Status }
    "backup" { Backup-Database }
    "clean" { Clear-Resources }
    default { Show-Help }
}
