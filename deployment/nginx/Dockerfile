FROM nginx:1.25-alpine

# 安装 wget 用于健康检查
RUN apk add --no-cache wget

# 创建日志目录
RUN mkdir -p /var/log/nginx

# 复制配置文件
COPY conf.d/default.conf /etc/nginx/conf.d/default.conf

# 暴露端口
EXPOSE 3000 3100

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:3000 || exit 1

# 启动 nginx
CMD ["nginx", "-g", "daemon off;"]
