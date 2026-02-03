#!/bin/bash
# generate-certs.sh
# 生成 mTLS 开发环境证书
# 仅用于开发和测试，生产环境请使用正式 CA 签发的证书

set -e

CERT_DIR="./certs"
DAYS=365

# 创建证书目录
mkdir -p "$CERT_DIR"

echo "=== 生成 CA 证书 ==="
openssl genrsa -out "$CERT_DIR/ca.key" 4096
openssl req -new -x509 -days $DAYS -key "$CERT_DIR/ca.key" -out "$CERT_DIR/ca.crt" \
    -subj "/C=CN/ST=Beijing/L=Beijing/O=NovelAI/OU=Dev/CN=NovelAI-CA"

echo "=== 生成服务端证书 (Python AI Service) ==="
openssl genrsa -out "$CERT_DIR/server.key" 2048
openssl req -new -key "$CERT_DIR/server.key" -out "$CERT_DIR/server.csr" \
    -subj "/C=CN/ST=Beijing/L=Beijing/O=NovelAI/OU=AIService/CN=localhost"

# 创建扩展配置
cat > "$CERT_DIR/server.ext" << EOF
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
subjectAltName = @alt_names

[alt_names]
DNS.1 = localhost
DNS.2 = ai-service
DNS.3 = python-ai-service
IP.1 = 127.0.0.1
IP.2 = ::1
EOF

openssl x509 -req -in "$CERT_DIR/server.csr" -CA "$CERT_DIR/ca.crt" -CAkey "$CERT_DIR/ca.key" \
    -CAcreateserial -out "$CERT_DIR/server.crt" -days $DAYS -extfile "$CERT_DIR/server.ext"

echo "=== 生成客户端证书 (Go Gateway) ==="
openssl genrsa -out "$CERT_DIR/client.key" 2048
openssl req -new -key "$CERT_DIR/client.key" -out "$CERT_DIR/client.csr" \
    -subj "/C=CN/ST=Beijing/L=Beijing/O=NovelAI/OU=Gateway/CN=go-gateway"

cat > "$CERT_DIR/client.ext" << EOF
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
extendedKeyUsage = clientAuth
EOF

openssl x509 -req -in "$CERT_DIR/client.csr" -CA "$CERT_DIR/ca.crt" -CAkey "$CERT_DIR/ca.key" \
    -CAcreateserial -out "$CERT_DIR/client.crt" -days $DAYS -extfile "$CERT_DIR/client.ext"

# 清理 CSR 和扩展文件
rm -f "$CERT_DIR"/*.csr "$CERT_DIR"/*.ext "$CERT_DIR"/*.srl

echo ""
echo "=== 证书生成完成 ==="
echo "CA 证书:     $CERT_DIR/ca.crt"
echo "服务端证书:  $CERT_DIR/server.crt"
echo "服务端私钥:  $CERT_DIR/server.key"
echo "客户端证书:  $CERT_DIR/client.crt"
echo "客户端私钥:  $CERT_DIR/client.key"
echo ""
echo "=== 配置说明 ==="
echo "Go Gateway (config.yaml):"
echo "  grpc:"
echo "    tls_enabled: true"
echo "    cert_file: ./certs/client.crt"
echo "    key_file: ./certs/client.key"
echo "    ca_file: ./certs/ca.crt"
echo ""
echo "Python AI Service (.env):"
echo "  GRPC_TLS_ENABLED=true"
echo "  GRPC_CERT_FILE=./certs/server.crt"
echo "  GRPC_KEY_FILE=./certs/server.key"
echo "  GRPC_CA_FILE=./certs/ca.crt"
