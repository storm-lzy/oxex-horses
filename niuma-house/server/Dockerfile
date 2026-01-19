# 构建阶段
FROM golang:1.23-alpine AS builder

WORKDIR /build

# 配置 Alpine 国内镜像源
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 安装必要的构建工具
RUN apk add --no-cache git ca-certificates tzdata

# 配置 Go 国内代理
ENV GOPROXY=https://goproxy.cn,direct
ENV GO111MODULE=on

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建二进制文件
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o server ./cmd/main.go

# 运行阶段
FROM alpine:latest

WORKDIR /app

# 配置 Alpine 国内镜像源
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 安装必要的运行时依赖
RUN apk add --no-cache ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 从构建阶段复制二进制文件
COPY --from=builder /build/server .

# 创建配置目录
RUN mkdir -p /app/config

# 暴露端口
EXPOSE 8080

# 启动服务
CMD ["./server"]
