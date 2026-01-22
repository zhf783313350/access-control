#!/bin/bash
# Kubernetes 部署脚本

set -e

NAMESPACE="access-control"
K8S_DIR="deploy/k8s"

echo "=========================================="
echo "  Access Control K8s 部署脚本"
echo "=========================================="

# 检查 kubectl 是否可用
if ! command -v kubectl &> /dev/null; then
    echo "错误: kubectl 未安装"
    exit 1
fi

# 函数: 显示帮助
show_help() {
    echo "用法: ./k8s-deploy.sh [命令]"
    echo ""
    echo "命令:"
    echo "  deploy      部署所有资源"
    echo "  delete      删除所有资源"
    echo "  status      查看部署状态"
    echo "  logs        查看 Pod 日志"
    echo "  restart     重启部署"
    echo "  scale N     扩缩容到 N 个副本"
    echo "  help        显示此帮助"
}

# 函数: 部署资源
deploy() {
    echo "正在部署到 Kubernetes..."
    
    # 1. 创建命名空间
    kubectl apply -f $K8S_DIR/namespace.yaml
    
    # 2. 创建 RBAC
    kubectl apply -f $K8S_DIR/rbac.yaml
    
    # 3. 创建 Secret (需要先手动修改)
    kubectl apply -f $K8S_DIR/secret.yaml
    
    # 4. 创建 ConfigMap
    kubectl apply -f $K8S_DIR/configmap.yaml
    
    # 5. 创建 Service
    kubectl apply -f $K8S_DIR/service.yaml
    
    # 6. 创建 Deployment
    kubectl apply -f $K8S_DIR/deployment.yaml
    
    # 7. 创建 HPA
    kubectl apply -f $K8S_DIR/hpa.yaml
    
    # 8. 创建 PDB
    kubectl apply -f $K8S_DIR/pdb.yaml
    
    # 9. 创建 Ingress
    kubectl apply -f $K8S_DIR/ingress.yaml
    
    echo "部署完成！"
    echo ""
    show_status
}

# 函数: 删除资源
delete() {
    echo "正在删除 Kubernetes 资源..."
    kubectl delete namespace $NAMESPACE --ignore-not-found
    echo "删除完成"
}

# 函数: 查看状态
show_status() {
    echo "Pod 状态:"
    kubectl get pods -n $NAMESPACE -o wide
    echo ""
    echo "Service 状态:"
    kubectl get svc -n $NAMESPACE
    echo ""
    echo "HPA 状态:"
    kubectl get hpa -n $NAMESPACE
    echo ""
    echo "Ingress 状态:"
    kubectl get ingress -n $NAMESPACE
}

# 函数: 查看日志
show_logs() {
    POD=$(kubectl get pods -n $NAMESPACE -l app=access-control-api -o jsonpath='{.items[0].metadata.name}')
    kubectl logs -f $POD -n $NAMESPACE
}

# 函数: 重启部署
restart() {
    echo "正在重启部署..."
    kubectl rollout restart deployment/access-control-api -n $NAMESPACE
    kubectl rollout status deployment/access-control-api -n $NAMESPACE
    echo "重启完成"
}

# 函数: 扩缩容
scale() {
    if [ -z "$1" ]; then
        echo "请指定副本数量"
        exit 1
    fi
    echo "正在扩缩容到 $1 个副本..."
    kubectl scale deployment/access-control-api --replicas=$1 -n $NAMESPACE
    echo "扩缩容完成"
}

# 主逻辑
case "$1" in
    deploy)
        deploy
        ;;
    delete)
        delete
        ;;
    status)
        show_status
        ;;
    logs)
        show_logs
        ;;
    restart)
        restart
        ;;
    scale)
        scale "$2"
        ;;
    help|*)
        show_help
        ;;
esac
