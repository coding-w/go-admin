# 使用 Alpine 版本的 Go 镜像作为构建阶段
FROM golang:alpine as builder

# 设置工作目录
WORKDIR /data/go-admin

# 将当前目录的所有文件复制到镜像中
COPY . .

# 配置 Go 环境变量并构建应用
# 启用 Go 模块支持
# 设置 Go 代理
# 禁用 Cgo
# 清理和下载依赖
# 编译应用，输出为 server
RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go env -w CGO_ENABLED=0 \
    && go env \
    && go mod tidy \
    && go build -o server .

# 使用更小的 Alpine 镜像作为最终镜像
FROM alpine:latest

# 设置维护者信息
LABEL MAINTAINER="wangxrz@163.com"

# 设置时区为上海
ENV TZ=Asia/Shanghai

# 更新包管理器并安装时区和 NTP 服务
# 设置时区
# 写入时区信息
RUN apk update && apk add --no-cache tzdata openntpd \
    && ln -sf /usr/share/zoneinfo/$TZ /etc/localtime \
    && echo $TZ > /etc/timezone

# 设置工作目录
WORKDIR /data/go-admin

# 从构建阶段复制编译后的二进制文件和配置文件
COPY --from=builder /data/go-admin/server ./
COPY --from=builder /data/go-admin/config.yaml ./

# 暴露应用监听的端口
EXPOSE 8888

# 设置容器启动时执行的命令
ENTRYPOINT ["./server", "-c", "config.yaml"]
