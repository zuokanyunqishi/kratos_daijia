#!/bin/bash

# 获取脚本所在的绝对目录
SCRIPT_DIR="$(realpath "$(dirname "$0")")"

# 定义基础目录
BASE_DIR="${SCRIPT_DIR}/deployments/k8s"

# 定义通用函数来删除服务
delete_service() {
    local dir=$1
    local config_file=$2
    echo "Deleting service in directory: $dir"
    cd "${BASE_DIR}/${dir}" || { echo "Failed to enter directory: $dir"; exit 1; }
    # 假设其他文件是以 .yaml 结尾的文件，排除指定的配置文件
    for file in $(ls *.yaml | grep -v "$config_file"); do
        kubectl delete -f "$file" || { echo "Failed to delete $file"; exit 1; }
    done
    kubectl delete -f "$config_file" || { echo "Failed to delete $config_file"; exit 1; }
    echo "Service in directory $dir deleted successfully."
}

# 步骤1: 删除 cluster 服务依赖
echo "Step 1: Deleting cluster dependencies..."
cd "${BASE_DIR}/cluster" || { echo "Failed to enter cluster directory"; exit 1; }
kubectl delete -f k8s-deployment.yaml || { echo "Failed to delete k8s-deployment.yaml"; exit 1; }
kubectl delete configmap consul-config -n default || { echo "Failed to delete consul-config configmap"; exit 1; }
echo "Cluster dependencies deleted successfully."

# 步骤2-6: 删除各个服务
delete_service "customer" "kratos-customer-config.yaml"
delete_service "driver" "kratos-driver-config.yaml"
delete_service "map" "kratos-map-config.yaml"
delete_service "valuation" "kratos-valuation-config.yaml"
delete_service "verifyCode" "kratos-verifyCode-config.yaml"

echo "All services have been deleted successfully."

