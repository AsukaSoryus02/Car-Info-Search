#!/bin/bash

# 设置颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# 输出带颜色的信息函数
info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查Redis是否已安装
check_redis() {
    if ! command -v redis-server &> /dev/null; then
        error "Redis未安装，请先安装Redis"
        exit 1
    fi
    info "Redis已安装"
}

# 启动Redis服务
start_redis() {
    info "正在启动Redis服务..."
    if pgrep redis-server > /dev/null; then
        warning "Redis服务已在运行"
    else
        redis-server &
        if [ $? -eq 0 ]; then
            info "Redis服务启动成功"
        else
            error "Redis服务启动失败"
            exit 1
        fi
    fi
}

# 启动后端服务
start_backend() {
    info "正在启动后端服务..."
    cd "$(dirname "$0")/server"
    
    # 检查是否有编译好的可执行文件
    if [ -f "carrag-server" ] && [ -x "carrag-server" ]; then
        ./carrag-server &
    else
        # 如果没有编译好的可执行文件，则直接运行Go代码
        go run main.go &
    fi
    
    # 检查后端是否成功启动
    sleep 2
    if curl -s http://localhost:8080/api/cars > /dev/null 2>&1; then
        info "后端服务启动成功，运行在 http://localhost:8080"
    else
        warning "后端服务可能未成功启动，请检查日志"
    fi
    
    cd - > /dev/null
}

# 启动前端服务
start_frontend() {
    info "正在启动前端服务..."
    cd "$(dirname "$0")"
    npm run dev &
    
    # 等待前端服务启动
    sleep 5
    info "前端服务启动成功，运行在 http://localhost:3000"
    cd - > /dev/null
}

# 主函数
main() {
    echo "========== CarRag 一键启动脚本 =========="
    
    # 检查Redis
    check_redis
    
    # 启动Redis
    start_redis
    
    # 启动后端
    start_backend
    
    # 启动前端
    start_frontend
    
    echo "========================================="
    info "所有服务已启动"
    info "前端访问地址: http://localhost:3000"
    info "后端API地址: http://localhost:8080/api"
    echo "使用 Ctrl+C 可以终止所有服务"
    echo "========================================="
    
    # 等待用户按Ctrl+C
    trap "echo '正在关闭所有服务...' && pkill -f 'redis-server' && pkill -f 'carrag-server' && pkill -f 'node'" INT
    wait
}

# 执行主函数
main