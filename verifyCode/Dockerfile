# 构建阶段
FROM golang:1.24 AS builder
COPY . /src
WORKDIR /src

 #使用 .dockerignore 排除无关文件
RUN GOPROXY=https://goproxy.cn make build && \
    chmod +x ./bin/verifyCode && \
    ls -l ./bin  # 验证构建结果

RUN GOPROXY=https://goproxy.cn go install github.com/go-delve/delve/cmd/dlv@latest

# 最终镜像阶段
FROM debian:stable-slim
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/* \
    && apt-get autoremove -y

# 创建用户并设置工作目录
WORKDIR /app
COPY --from=builder /src/bin/verifyCode /app/

# 配置目录和权限
RUN mkdir -p /data/conf && chmod 755 /data/conf
# 安装 dlv 调试
COPY --from=builder /go/bin/dlv /usr/local/bin/dlv
# 暴露端口
EXPOSE 8600
EXPOSE 9600
 # dlv 调试端口
EXPOSE 2345
# 启动命令
CMD  ["./verifyCode", "-conf", "/data/conf"]

#CMD ["dlv", "exec", "./verifyCode", "--headless", "--listen=:2345", "--api-version=2", "--", "-conf", "/data/conf"]