#!/bin/bash

# 获取脚本所在的绝对目录
SCRIPT_DIR="$(realpath "$(dirname "$0")")"

# 定义基础目录
BASE_DIR="${SCRIPT_DIR}/deployments/k8s"

# 定义颜色代码
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 定义通用函数来启动服务
start_service() {
    local dir=$1
    local config_file=$2
    echo "Starting service in directory: $dir"
    cd "${BASE_DIR}/${dir}" || { echo -e "${RED}Failed to enter directory: $dir${NC}"; exit 1; }
    kubectl apply -f "$config_file" || { echo -e "${RED}Failed to apply $config_file${NC}"; exit 1; }
    # 假设其他文件是以 .yaml 结尾的文件，排除指定的配置文件
    for file in $(ls *.yaml | grep -v "$config_file"); do
        kubectl apply -f "$file" || { echo -e "${RED}Failed to apply $file${NC}"; exit 1; }
    done
    echo -e "${GREEN}Service in directory $dir started successfully.${NC}"
}

# 步骤1: 启动 cluster 服务依赖
echo "Step 1: Starting cluster dependencies..."
cd "${BASE_DIR}/cluster" || { echo -e "${RED}Failed to enter cluster directory${NC}"; exit 1; }

# 检查文件是否存在
if [ ! -f "configs/amap.yaml" ]; then
    echo -e "${RED}Error: configs/amap.yaml does not exist${NC}"
    exit 1
fi

kubectl create configmap consul-config -n default --from-file=configs/amap.yaml || { echo -e "${RED}Failed to create configmap${NC}"; exit 1; }
kubectl apply -f k8s-deployment.yaml || { echo -e "${RED}Failed to apply k8s-deployment.yaml${NC}"; exit 1; }
echo -e "${GREEN}Cluster dependencies started successfully.${NC}"

# 等待 consul 服务启动完毕，使用 kubectl status 命令
echo "Waiting for consul service to be ready..."
kubectl rollout status StatefulSet/consul -n default || { echo -e "${RED}Failed to wait for consul deployment${NC}"; exit 1; }
echo -e "${GREEN}Consul service is ready.${NC}"


# 等待 mysql 服务启动完毕
echo "Waiting for mysql service to be ready..."
kubectl rollout status Deployment/mysql -n default || { echo -e "${RED}Failed to wait for mysql deployment${NC}"; exit 1; }
echo -e "${GREEN}MySQL service is ready.${NC}"

# 等待 redis 服务启动完毕
echo "Waiting for redis service to be ready..."
kubectl rollout status Deployment/redis -n default || { echo -e "${RED}Failed to wait for redis deployment${NC}"; exit 1; }
echo -e "${GREEN}Redis service is ready.${NC}"

# 等待 jaeger 服务启动完毕
echo "Waiting for jaeger service to be ready..."
kubectl rollout status Deployment/jaeger -n default || { echo -e "${RED}Failed to wait for jaeger deployment${NC}"; exit 1; }
echo -e "${GREEN}Jaeger service is ready.${NC}"

# 步骤2-6: 启动各个服务
start_service "customer" "kratos-customer-config.yaml"
start_service "driver" "kratos-driver-config.yaml"
start_service "map" "kratos-map-config.yaml"
start_service "valuation" "kratos-valuation-config.yaml"
start_service "verifyCode" "kratos-verifyCode-config.yaml"

echo -e "${GREEN}All services have been started successfully.${NC}"
