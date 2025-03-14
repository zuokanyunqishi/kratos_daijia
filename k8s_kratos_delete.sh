#!/bin/bash

# 获取脚本所在的绝对目录
SCRIPT_DIR="$(realpath "$(dirname "$0")")"

# 定义基础目录
BASE_DIR="${SCRIPT_DIR}/deployments/k8s"

# 定义颜色代码
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 定义通用函数来删除服务
delete_service() {
    local dir=$1
    local config_file=$2
    echo "Deleting service in directory: $dir"
    cd "${BASE_DIR}/${dir}" || { echo -e "${RED}Failed to enter directory: $dir${NC}"; exit 1; }
    # 假设其他文件是以 .yaml 结尾的文件，排除指定的配置文件
    for file in $(ls *.yaml | grep -v "$config_file"); do
        delete_resource "$file"
    done
    delete_resource "$config_file"
    echo -e "${GREEN}Service in directory $dir deleted successfully.${NC}"
}

# 删除单个资源的函数
delete_resource() {
    local file=$1
    output=$(kubectl delete -f "$file" --ignore-not-found=true 2>&1)
    exit_code=$?
    if echo "$output" | grep -q "not found"; then
        echo -e "${GREEN}Resource in $file already does not exist.${NC}"
    elif [ $exit_code -ne 0 ]; then
        echo -e "${RED}Failed to delete $file: $output${NC}"
        exit 1
    else
        echo -e "${GREEN}Deleted $file successfully.${NC}"
    fi
}

# 删除单个资源的函数（适用于直接使用 kubectl delete 命令）
delete_resource_direct() {
    local resource_type=$1
    local resource_name=$2
    local namespace=$3
    output=$(kubectl delete "$resource_type" "$resource_name" -n "$namespace" --ignore-not-found=true 2>&1)
    exit_code=$?
    if echo "$output" | grep -q "not found"; then
        echo -e "${GREEN}Resource $resource_type/$resource_name already does not exist in namespace $namespace.${NC}"
    elif [ $exit_code -ne 0 ]; then
        echo -e "${RED}Failed to delete $resource_type/$resource_name in namespace $namespace: $output${NC}"
        exit 1
    else
        echo -e "${GREEN}Deleted $resource_type/$resource_name in namespace $namespace successfully.${NC}"
    fi
}

# 步骤1: 删除 cluster 服务依赖
echo "Step 1: Deleting cluster dependencies..."
cd "${BASE_DIR}/cluster" || { echo -e "${RED}Failed to enter cluster directory${NC}"; exit 1; }
delete_resource "k8s-deployment.yaml"
delete_resource_direct "configmap" "consul-config" "default"
echo -e "${GREEN}Cluster dependencies deleted successfully.${NC}"

# 步骤2-6: 删除各个服务
delete_service "customer" "kratos-customer-config.yaml"
delete_service "driver" "kratos-driver-config.yaml"
delete_service "map" "kratos-map-config.yaml"
delete_service "valuation" "kratos-valuation-config.yaml"
delete_service "verifyCode" "kratos-verifyCode-config.yaml"

echo -e "${GREEN}All services have been deleted successfully.${NC}"
