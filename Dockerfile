FROM golang:1.18-alpine AS builder

RUN apk add --no-cache git gcc musl-dev sqlite-dev

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download
COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o claude-code-relay main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata sqlite

ENV TZ=Asia/Shanghai

WORKDIR /app

RUN mkdir -p /app/data /app/logs

COPY --from=builder /app/claude-code-relay .

COPY .env.example .env.example

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./claude-code-relay"]