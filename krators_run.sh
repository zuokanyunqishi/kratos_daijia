#!/bin/bash
# 一个窗口启动多个微服务, 并自动优化窗格布局。需要 tmux 支持。

# 配置需要排除的目录列表（支持空格）
exclude_dirs=(
    "docker-compose-deploy"
    "docs"
    "scripts"
    "configs"
    "."  # 显式排除当前目录
)

# 自动发现服务目录（排除当前目录和隐藏目录）
services=()
while IFS= read -r -d '' dir; do
    dir_name=$(basename "$dir")
    should_exclude=false

    # 检查是否在排除列表
    for exclude in "${exclude_dirs[@]}"; do
        if [[ "$dir_name" == "$exclude" ]]; then
            should_exclude=true
            break
        fi
    done

    # 额外排除隐藏目录（以.开头的目录）
    if [[ "$dir_name" == .* ]]; then
        should_exclude=true
    fi

    if ! $should_exclude; then
        services+=("$dir")
    fi
done < <(find . -mindepth 1 -maxdepth 1 -type d -print0)

# 检查是否找到服务
if [ ${#services[@]} -eq 0 ]; then
    echo "错误：未发现任何有效服务目录！"
    echo "已排除目录：${exclude_dirs[*]} 和所有隐藏目录"
    exit 1
fi

# 创建 tmux 会话
session_name="kratos_services"
tmux new-session -d -s "$session_name"

# 为每个服务创建窗格
for idx in "${!services[@]}"; do
    service_path="${services[$idx]}"

    if [ $idx -eq 0 ]; then
        tmux send-keys -t "$session_name" "cd \"$service_path\" && kratos run" C-m
    else
        tmux split-window -v -t "$session_name" "cd \"$service_path\" && kratos run"
        tmux select-layout -t "$session_name" tiled
    fi
done

# 优化布局并附加会话
tmux select-layout -t "$session_name" even-vertical
tmux attach -t "$session_name"
