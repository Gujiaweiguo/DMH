#!/bin/bash
# DMH 独立容器部署脚本

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log() { echo -e "${GREEN}[$(date +'%H:%M:%S')]${NC} $1"; }
warn() { echo -e "${YELLOW}[WARNING]${NC} $1"; }
error() { echo -e "${RED}[ERROR]${NC} $1"; exit 1; }
info() { echo -e "${BLUE}[INFO]${NC} $1"; }

# 配置
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOYMENT_DIR="$(dirname "$SCRIPT_DIR")"
PROJECT_ROOT="$(dirname "$DEPLOYMENT_DIR")"
BACKEND_DIR="$PROJECT_ROOT/backend"
COMPOSE_FILE="$DEPLOYMENT_DIR/docker-compose.yml"

# 步骤1: 停止现有服务
log "步骤 1: 停止现有服务..."
cd "$DEPLOYMENT_DIR"

# 停止并删除旧容器
if docker compose -f "$COMPOSE_FILE" ps 2>/dev/null | grep -q "dmh-nginx\|dmh-api"; then
    info "停止旧容器..."
    docker compose -f "$COMPOSE_FILE" down
    log "旧容器已停止"
else
    info "没有运行中的 DMH 容器"
fi

# 停止独立进程
if pgrep -f "dmh-api" > /dev/null 2>&1; then
    info "停止独立进程..."
    cd "$BACKEND_DIR"
    if [ -f "deploy.sh" ]; then
        ./deploy.sh -t || warn "停止进程失败，继续..."
    fi
    log "独立进程已停止"
else
    info "没有运行中的独立进程"
fi

# 步骤2: 检查前端构建产物
log "步骤 2: 检查前端构建产物..."
if [ ! -d "$PROJECT_ROOT/frontend-admin/dist" ]; then
    warn "管理后台未构建，正在构建..."
    cd "$PROJECT_ROOT/frontend-admin"
    npm install
    npm run build
    log "管理后台构建完成"
else
    info "管理后台已构建"
fi

if [ ! -d "$PROJECT_ROOT/frontend-h5/dist" ]; then
    warn "H5前端未构建，正在构建..."
    cd "$PROJECT_ROOT/frontend-h5"
    npm install
    npm run build
    log "H5前端构建完成"
else
    info "H5前端已构建"
fi

# 步骤3: 构建后端镜像
log "步骤 3: 构建后端镜像..."
cd "$BACKEND_DIR"
docker build -t dmh-api:latest -f Dockerfile . || error "后端镜像构建失败"
log "后端镜像构建完成"

# 步骤4: 构建Nginx镜像
log "步骤 4: 构建 Nginx 镜像..."
cd "$DEPLOYMENT_DIR/nginx"
docker build -t dmh-nginx:latest -f Dockerfile . || error "Nginx镜像构建失败"
log "Nginx镜像构建完成"

# 步骤5: 启动 Docker Compose
log "步骤 5: 启动 Docker Compose..."
cd "$DEPLOYMENT_DIR"
docker compose -f "$COMPOSE_FILE" up -d || error "Docker Compose启动失败"

# 步骤6: 等待服务就绪
log "步骤 6: 等待服务就绪..."
info "等待 Nginx 服务..."
for i in {1..30}; do
    if curl -s http://localhost:3000 > /dev/null 2>&1; then
        log "✓ Nginx 服务已就绪"
        break
    fi
    if [ $i -eq 30 ]; then
        warn "Nginx 服务启动超时"
    fi
    sleep 2
    info "等待中... ($i/30)"
done

info "等待后端 API 服务..."
for i in {1..30}; do
    if curl -s http://localhost:8889 > /dev/null 2>&1; then
        log "✓ 后端 API 服务已就绪"
        break
    fi
    if [ $i -eq 30 ]; then
        warn "后端 API 服务启动超时"
    fi
    sleep 2
    info "等待中... ($i/30)"
done

# 步骤7: 验证服务
log "步骤 7: 验证服务..."
sleep 3

# 检查容器状态
info "检查容器状态..."
docker compose -f "$COMPOSE_FILE" ps

# 测试管理后台
info "测试管理后台 (http://localhost:3000)..."
if curl -s http://localhost:3000 | grep -q "html"; then
    log "✓ 管理后台正常"
else
    warn "管理后台可能有问题"
fi

# 测试H5前端
info "测试 H5 前端 (http://localhost:3100)..."
if curl -s http://localhost:3100 | grep -q "html"; then
    log "✓ H5 前端正常"
else
    warn "H5 前端可能有问题"
fi

# 测试API登录
info "测试 API 登录..."
LOGIN_RESULT=$(curl -s -X POST http://localhost:8889/api/v1/auth/login \
    -H "Content-Type: application/json" \
    -d '{"username":"admin","password":"123456"}')

if echo "$LOGIN_RESULT" | grep -q "token"; then
    log "✓ API 登录测试通过"
else
    warn "API 登录测试失败"
fi

# 步骤8: 显示管理命令
log "========================================="
log "独立容器部署完成！"
log "========================================="
info ""
info "服务访问地址："
info "  管理后台:      http://localhost:3000"
info "  H5 前端:       http://localhost:3100"
info "  后端 API:      http://localhost:8889"
info ""
info "容器管理命令："
info "  查看状态:      cd $DEPLOYMENT_DIR && docker compose ps"
info "  查看日志:      docker compose logs -f"
info "  重启服务:      docker compose restart"
info "  停止服务:      docker compose stop"
info "  启动服务:      docker compose start"
info "  完全清理:      docker compose down"
info ""
info "单容器日志："
info "  Nginx日志:     docker logs -f dmh-nginx"
info "  API日志:       docker logs -f dmh-api"
info ""
info "进入容器："
info "  Nginx容器:     docker exec -it dmh-nginx sh"
info "  API容器:       docker exec -it dmh-api sh"
info ""
warn "回滚到独立进程："
warn "  cd $DEPLOYMENT_DIR/scripts && ./rollback-containers.sh"
warn "========================================="
