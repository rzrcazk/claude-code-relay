# 前端构建阶段
FROM node:18.18-alpine AS frontend-builder

WORKDIR /app/web

# 安装pnpm
RUN npm install -g pnpm

# 复制前端项目文件
COPY web/package.json web/pnpm-lock.yaml* ./

# 安装依赖（跳过prepare脚本避免husky问题）
RUN pnpm install --frozen-lockfile --prod --ignore-scripts

# 复制前端源码
COPY web/ ./

# 构建前端项目
RUN pnpm run build

# 后端构建阶段
FROM golang:1.21-alpine AS backend-builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w -s' -o claude-code-relay main.go

# 最终运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata wget

ENV TZ=Asia/Shanghai

WORKDIR /app

RUN mkdir -p /app/logs

# 复制后端可执行文件
COPY --from=backend-builder /app/claude-code-relay .

# 复制前端构建产物
COPY --from=frontend-builder /app/web/dist ./web/dist

COPY .env.example .env.example

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./claude-code-relay"]