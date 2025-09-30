# Temporal SDK Example

### SDK特性

- **模块化设计**: 支持业务模块的独立注册和管理
- **错误处理**: 完整的错误类型系统
- **重试策略**: 预定义的重试策略（标准、关键、最小、无重试）
- **配置管理**: 支持YAML配置文件
- **日志系统**: 基于zap的结构化日志

## 快速开始

### 1. 启动Temporal服务

```bash
# 使用Docker Compose启动Temporal服务
git clone https://github.com/temporalio/docker-compose.git
cd docker-compose
docker-compose up -d
```

## 配置说明

### config.yaml

```yaml
temporal:
  host_port: "localhost:7233"          # Temporal服务地址
  namespace: "order-processing"        # 命名空间
  identity: "order-worker"             # 工作者身份标识
  enabled: true                        # 是否启用

  workers: # Worker配置
    - task_queue: "order-processing"   # 任务队列名称
      max_concurrent_activities: 100   # 最大并发活动数
      max_concurrent_workflows: 50     # 最大并发工作流数
      enable_session_worker: false     # 是否启用会话Worker
      sticky_schedule_to_start: "10s"  # Sticky调度超时

    - task_queue: "notifications"      # 通知任务队列
      max_concurrent_activities: 200
      max_concurrent_workflows: 100

  retry: # 重试配置
    initial_interval: "1s"             # 初始重试间隔
    backoff_coefficient: 2.0           # 退避系数
    maximum_interval: "60s"            # 最大重试间隔
    maximum_attempts: 5                # 最大重试次数

  timeout: # 超时配置
    workflow_execution: "24h"          # 工作流执行超时
    workflow_task: "10s"               # 工作流任务超时
    activity_execution: "30m"          # 活动执行超时
    activity_start: "5s"               # 活动启动超时
    activity_heartbeat: "10s"          # 活动心跳超时
```
