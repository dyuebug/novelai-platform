#!/bin/sh
set -e

# 设置默认值
export VITE_API_BASE_URL=${VITE_API_BASE_URL:-http://localhost:8080}

echo "Starting frontend container..."
echo "VITE_API_BASE_URL: $VITE_API_BASE_URL"

# 验证 VITE_API_BASE_URL 格式
if ! echo "$VITE_API_BASE_URL" | grep -qE '^https?://'; then
    echo "ERROR: VITE_API_BASE_URL must be a valid HTTP/HTTPS URL"
    echo "Current value: $VITE_API_BASE_URL"
    exit 1
fi

# 使用 envsubst 替换 Nginx 配置中的环境变量
echo "Generating Nginx configuration from template..."
envsubst '${VITE_API_BASE_URL}' < /etc/nginx/templates/default.conf.template > /etc/nginx/conf.d/default.conf

# 验证 Nginx 配置语法
echo "Validating Nginx configuration..."
nginx -t

if [ $? -ne 0 ]; then
    echo "ERROR: Nginx configuration validation failed"
    exit 1
fi

echo "Nginx configuration validated successfully"
echo "Starting Nginx..."

# 执行传入的命令（通常是 nginx -g "daemon off;"）
exec "$@"
