#!/bin/bash

### 自签证通配符证书，用来模拟线上 tls 功能
# 设置变量
DOMAIN="devk8s.com"
WILDCARD="*.${DOMAIN}"
NAMESPACE="default"
SECRET_NAME="wildcard-tls-secret"

# 创建工作目录
mkdir -p certs && cd certs

# 1. 生成 CA 私钥和证书
echo "生成 CA 证书..."
openssl genrsa -out ca.key 2048
openssl req -x509 -new -nodes -key ca.key -sha256 -days 3650 -out ca.crt \
  -subj "/CN=MyLocalCA"

# 2. 生成通配符证书的私钥和 CSR
echo "生成通配符证书私钥和 CSR..."
openssl genrsa -out wildcard.${DOMAIN}.key 2048
openssl req -new -key wildcard.${DOMAIN}.key -out wildcard.${DOMAIN}.csr \
  -subj "/CN=${WILDCARD}"

# 3. 创建扩展文件并签署通配符证书
echo "签署通配符证书..."
cat > wildcard.ext <<EOF
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
subjectAltName = DNS:${WILDCARD}, DNS:${DOMAIN}
EOF

openssl x509 -req -in wildcard.${DOMAIN}.csr -CA ca.crt -CAkey ca.key \
  -CAcreateserial -out wildcard.${DOMAIN}.crt -days 365 -sha256 \
  -extfile wildcard.ext

# 4. 创建 Kubernetes TLS Secret
echo "创建 Kubernetes Secret..."
kubectl create secret tls ${SECRET_NAME} \
  --cert=wildcard.${DOMAIN}.crt \
  --key=wildcard.${DOMAIN}.key \
  -n ${NAMESPACE} \
  --dry-run=client -o yaml | kubectl apply -f -

# 5. 清理临时文件
rm -f wildcard.${DOMAIN}.csr wildcard.ext ca.srl

echo "完成！"
echo "证书文件位于 ./certs 目录下。"
echo "Secret '${SECRET_NAME}' 已创建在命名空间 '${NAMESPACE}' 中。"
echo "下一步："
echo "1. 配置 Ingress 使用 'secretName: ${SECRET_NAME}' 和 'hosts: ${WILDCARD}'。"
echo "2. 更新 /etc/hosts，例如：'127.0.0.1 grpc-driver.${DOMAIN}'。"
echo "3. 测试：curl -v --cacert ca.crt https://grpc-driver.${DOMAIN}"