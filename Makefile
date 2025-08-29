BINARY_NAME=claude-code-relay
OUT_DIR=out
WEB_DIR=web

all: build

# 构建前端
build-web: clean
	cd ${WEB_DIR} && pnpm install && pnpm run build

# 构建后端（包含嵌入的前端资源）
build: prepare build-web
	GOOS=linux GOARCH=amd64 go build -o ${OUT_DIR}/${BINARY_NAME}-linux-amd64 main.go
	GOOS=windows GOARCH=amd64 go build -o ${OUT_DIR}/${BINARY_NAME}-windows-amd64.exe main.go
	GOOS=darwin GOARCH=amd64 go build -o ${OUT_DIR}/${BINARY_NAME}-darwin-amd64 main.go
	GOOS=linux GOARCH=arm64 go build -o ${OUT_DIR}/${BINARY_NAME}-linux-arm64 main.go
	GOOS=darwin GOARCH=arm64 go build -o ${OUT_DIR}/${BINARY_NAME}-darwin-arm64 main.go

prepare:
	mkdir -p ${OUT_DIR}

clean:
	rm -rf ${WEB_DIR}/dist

.PHONY: all build build-web prepare clean