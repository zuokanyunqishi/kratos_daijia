#!/bin/bash
# 获取脚本所在的绝对目录
SCRIPT_DIR="$(realpath "$(dirname "$0")")"

# 定义基础目录
BASE_DIR="${SCRIPT_DIR}/deployments/k8s"

# 定义通用函数来启动服务
start_service() {
    local dir=$1
    local config_file=$2
    echo "Starting service in directory: $dir"
    cd "${BASE_DIR}/${dir}" || { echo "Failed to enter directory: $dir"; exit 1; }
    kubectl apply -f "$config_file" || { echo "Failed to apply $config_file"; exit 1; }
    # 假设其他文件是以 .yaml 结尾的文件，排除指定的配置文件
    for file in $(ls *.yaml | grep -v "$config_file"); do
        kubectl apply -f "$file" || { echo "Failed to apply $file"; exit 1; }
    done
    echo "Service in directory $dir started successfully."
}

# 步骤1: 启动 cluster 服务依赖
echo "Step 1: Starting cluster dependencies..."
cd "${BASE_DIR}/cluster" || { echo "Failed to enter cluster directory"; exit 1; }

# 检查文件是否存在
if [ ! -f "configs/amap.yaml" ]; then
    echo "Error: configs/amap.yaml does not exist"
    exit 1
fi

kubectl create configmap consul-config -n default --from-file=configs/amap.yaml || { echo "Failed to create configmap"; exit 1; }
kubectl apply -f k8s-deployment.yaml || { echo "Failed to apply k8s-deployment.yaml"; exit 1; }
echo "Cluster dependencies started successfully."

# 等待 consul 服务启动完毕，使用 kubectl wait 命令
echo "Waiting for consul service to be ready..."
kubectl rollout status StatefulSet/consul  -n default || { echo "Failed to wait for consul deployment"; exit 1; }
echo "Consul service is ready."



# 步骤2-6: 启动各个服务
start_service "customer" "kratos-customer-config.yaml"
start_service "driver" "kratos-driver-config.yaml"
start_service "map" "kratos-map-config.yaml"
start_service "valuation" "kratos-valuation-config.yaml"
start_service "verifyCode" "kratos-verifyCode-config.yaml"

echo "All services have been started successfully."

