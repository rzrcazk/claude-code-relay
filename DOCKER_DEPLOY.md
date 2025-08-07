# Docker部署指南

## 快速开始

本项目提供两种Docker部署方式，根据你的环境选择合适的方案：

### 方式一：使用现有MySQL和Redis服务（推荐）

适用于已经有MySQL和Redis服务的环境，使用轻量级的单应用部署。

```bash
# 1. 复制环境变量模板
cp .env.example .env

# 2. 编辑.env文件，配置数据库连接信息
# 必需配置项：
# SESSION_SECRET=your-session-secret
# JWT_SECRET=your-jwt-secret  
# SALT=your-salt-value
# MYSQL_HOST=your-mysql-host
# MYSQL_PASSWORD=your-mysql-password
# REDIS_HOST=your-redis-host

# 3. 生成安全密钥(或随机生成字符串即可)
openssl rand -base64 32  # 用于SESSION_SECRET
openssl rand -base64 32  # 用于JWT_SECRET  
openssl rand -base64 16  # 用于SALT

# 4. 启动应用
docker-compose up -d

# 5. 查看状态
docker-compose ps
docker-compose logs -f app
```

### 方式二：一键部署全套服务

适用于新环境，自动创建MySQL、Redis和应用服务。

```bash
# 1. 复制环境变量模板（可选）
cp .env.example .env

# 2. 启动全套服务
docker-compose -f docker-compose-all.yml up -d

# 3. 查看所有服务状态
docker-compose -f docker-compose-all.yml ps

# 4. 查看日志
docker-compose -f docker-compose-all.yml logs -f
```

**默认配置：**
- 应用访问地址：http://localhost:10081
- MySQL数据库：claude_code_relay
- MySQL账号：claude / claude123456
- Redis：无密码，默认配置
- 管理员账号：admin / admin123

## 详细配置说明

### 环境变量配置

创建 `.env` 文件并配置以下参数：

```bash
# 基础配置
PORT=8080
GIN_MODE=release
HTTP_CLIENT_TIMEOUT=120

# 安全配置（生产环境必须）
SESSION_SECRET=your-session-secret-32-chars
JWT_SECRET=your-jwt-secret-32-chars
SALT=your-salt-16-chars

# MySQL配置（使用现有MySQL时）
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_USER=root
MYSQL_PASSWORD=your-mysql-password
MYSQL_DATABASE=claude_code_relay

# Redis配置（使用现有Redis时）
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=your-redis-password
REDIS_DB=0

# 日志配置
LOG_LEVEL=info
LOG_FILE=/app/logs/app.log
LOG_RECORD_API=true
LOG_RETENTION_MONTHS=3

# 管理员配置
DEFAULT_ADMIN_USERNAME=admin
DEFAULT_ADMIN_PASSWORD=your-admin-password
DEFAULT_ADMIN_EMAIL=admin@your-domain.com
```

### 数据持久化

#### 使用本地目录挂载

修改 `docker-compose-all.yml` 中的volumes配置：

```yaml
volumes:
  mysql_data:
    driver: local
    driver_opts:
      type: none
      o: bind
      device: ./data/mysql  # 本地MySQL数据目录
  redis_data:
    driver: local
    driver_opts:
      type: none
      o: bind
      device: ./data/redis  # 本地Redis数据目录
  app_logs:
    driver: local
    driver_opts:
      type: none
      o: bind
      device: ./logs        # 本地日志目录
```

启动前创建目录：

```bash
mkdir -p data/mysql data/redis logs
chmod 755 data/mysql data/redis logs
```

#### 使用Docker Volumes

```bash
# 查看所有volumes
docker volume ls

# 查看volume详情
docker volume inspect claude-code-relay_mysql_data

# 备份数据
docker run --rm -v claude-code-relay_mysql_data:/data -v $(pwd):/backup alpine tar czf /backup/mysql-backup.tar.gz -C /data .

# 恢复数据
docker run --rm -v claude-code-relay_mysql_data:/data -v $(pwd):/backup alpine tar xzf /backup/mysql-backup.tar.gz -C /data
```

## 运维管理

### 服务管理

```bash
# 启动所有服务
docker-compose -f docker-compose-all.yml up -d

# 重启特定服务
docker-compose -f docker-compose-all.yml restart app

# 停止所有服务
docker-compose -f docker-compose-all.yml stop

# 完全删除服务和数据
docker-compose -f docker-compose-all.yml down -v

# 查看服务状态
docker-compose -f docker-compose-all.yml ps

# 查看资源使用情况
docker stats claude-code-relay-app claude-code-relay-mysql claude-code-relay-redis
```

### 日志管理

```bash
# 查看应用日志
docker-compose logs -f app

# 查看MySQL日志
docker-compose logs -f mysql

# 查看Redis日志
docker-compose logs -f redis

# 查看日志文件
tail -f ./logs/app.log

# 清理日志
docker-compose exec app sh -c 'echo > /app/logs/app.log'
```

### 健康检查

```bash
# 检查应用健康状态
curl -f http://localhost:10081/health

# 检查数据库状态
curl -s http://localhost:10081/api/v1/status | jq .database_status

# 检查Redis状态
curl -s http://localhost:10081/api/v1/status | jq .redis_status

# 进入容器调试
docker-compose exec app sh
docker-compose exec mysql mysql -u root -p
docker-compose exec redis redis-cli
```

### 更新部署

```bash
# 拉取最新镜像
docker-compose pull

# 重新创建容器（保留数据）
docker-compose up -d --force-recreate

# 完整更新流程
docker-compose pull
docker-compose down
docker-compose up -d
```

## 性能调优

### MySQL优化

在 `docker-compose-all.yml` 中调整MySQL参数：

```yaml
command: >
  --default-authentication-plugin=mysql_native_password
  --character-set-server=utf8mb4
  --collation-server=utf8mb4_unicode_ci
  --default-time-zone='+08:00'
  --innodb-buffer-pool-size=512M        # 根据内存大小调整
  --max-connections=2000                 # 根据并发需求调整
  --innodb-log-file-size=128M
  --innodb-flush-log-at-trx-commit=2
  --slow-query-log=1
  --long-query-time=2
```

### Redis优化

```yaml
command: >
  redis-server
  --appendonly yes
  --appendfsync everysec
  --maxmemory 1gb                        # 根据可用内存调整
  --maxmemory-policy allkeys-lru
  --save 900 1 300 10 60 10000
  --tcp-keepalive 300
  --timeout 300
```

### 应用优化

```yaml
environment:
  GIN_MODE: release                      # 生产模式
  LOG_LEVEL: warn                        # 减少日志输出
  HTTP_CLIENT_TIMEOUT: 60                # 根据网络情况调整
deploy:
  resources:
    limits:
      cpus: '2.0'                        # 根据CPU核数调整
      memory: 2G                         # 根据内存大小调整
```

## 安全配置

### 网络安全

```yaml
# 自定义网络配置
networks:
  claude-network:
    driver: bridge
    ipam:
      config:
        - subnet: 172.20.0.0/16
          gateway: 172.20.0.1
    driver_opts:
      com.docker.network.bridge.name: claude-br0
      com.docker.network.bridge.enable_icc: "false"
```

### 容器安全

```yaml
security_opt:
  - no-new-privileges:true               # 禁止提升权限
  - seccomp:unconfined                   # 如需要可启用seccomp
user: "1000:1000"                        # 非root用户运行
read_only: true                          # 只读文件系统
tmpfs:
  - /tmp:noexec,nosuid,size=100m         # 临时目录
```

### 防火墙配置

```bash
# 仅允许必要端口访问
ufw allow 10081/tcp                      # 应用端口
ufw deny 3306/tcp                        # MySQL端口（仅内部访问）
ufw deny 6379/tcp                        # Redis端口（仅内部访问）
```

## 监控告警

### Prometheus监控

```yaml
# docker-compose-monitoring.yml
version: '3.8'
services:
  prometheus:
    image: prom/prometheus
    ports:
      - "9090:9090"
    volumes:
      - ./monitoring/prometheus.yml:/etc/prometheus/prometheus.yml
    networks:
      - claude-network

  grafana:
    image: grafana/grafana
    ports:
      - "3000:3000"
    environment:
      GF_SECURITY_ADMIN_PASSWORD: admin
    volumes:
      - grafana_data:/var/lib/grafana
    networks:
      - claude-network
```

### 基础告警

```bash
# CPU使用率监控
docker stats --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}" --no-stream

# 磁盘空间监控
df -h | grep -E "(mysql|redis|logs)"

# 服务可用性监控
curl -f http://localhost:10081/health || echo "Service Down" | mail -s "Alert" admin@example.com
```

## 故障排查

### 常见问题

1. **应用启动失败**
```bash
# 检查日志
docker-compose logs app

# 检查环境变量
docker-compose exec app env | grep -E "(MYSQL|REDIS|SECRET)"

# 检查网络连接
docker-compose exec app ping mysql
docker-compose exec app ping redis
```

2. **数据库连接失败**
```bash
# 检查MySQL状态
docker-compose exec mysql mysql -u root -pclaude123456 -e "SELECT 1"

# 测试连接
docker-compose exec app sh -c 'mysql -h mysql -u claude -pclaude123456 -e "SELECT 1"'
```

3. **Redis连接失败**
```bash
# 检查Redis状态
docker-compose exec redis redis-cli ping

# 测试连接
docker-compose exec app sh -c 'redis-cli -h redis ping'
```

### 紧急恢复

```bash
# 快速重启所有服务
docker-compose -f docker-compose-all.yml restart

# 重建应用容器
docker-compose -f docker-compose-all.yml up -d --force-recreate app

# 数据库恢复
docker-compose exec mysql mysql -u root -p < backup.sql

# 清理并重新部署
docker-compose -f docker-compose-all.yml down
docker system prune -f
docker-compose -f docker-compose-all.yml up -d
```

## 备份恢复

### 自动备份脚本

```bash
#!/bin/bash
# backup.sh

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="./backups/$DATE"
mkdir -p $BACKUP_DIR

# 备份MySQL数据
docker-compose exec mysql mysqldump -u root -pclaude123456 --all-databases > $BACKUP_DIR/mysql_backup.sql

# 备份Redis数据
docker-compose exec redis redis-cli --rdb /data/dump.rdb
docker cp claude-code-relay-redis:/data/dump.rdb $BACKUP_DIR/redis_backup.rdb

# 备份应用日志
cp -r logs $BACKUP_DIR/

# 备份配置文件
cp .env docker-compose*.yml $BACKUP_DIR/

echo "Backup completed: $BACKUP_DIR"
```

### 定时备份

```bash
# 添加到crontab
crontab -e

# 每天凌晨2点备份
0 2 * * * /path/to/backup.sh
```