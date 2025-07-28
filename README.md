# Alfred Tool - SSH 连接管理工具

一个功能完整的 SSH 连接和服务器管理工具，支持 SSH 连接管理、Rsync 文件同步和服务部署管理。基于 Fyne GUI 和命令行双重接口，使用 GORM 和 SQLite 进行数据存储。

## 功能特性

### SSH 连接管理
- **添加连接**: 通过优雅的 Fyne 表单界面添加 SSH 连接
- **搜索连接**: 根据连接名称或地址搜索已保存的连接
- **列表显示**: 显示所有已保存的 SSH 连接的简洁列表
- **多种认证**: 支持密码和私钥文件两种认证方式
- **使用统计**: 自动记录连接使用次数

### Rsync 文件同步
- **配置管理**: 创建和管理 rsync 同步配置
- **双向同步**: 支持上传和下载文件同步
- **排除规则**: 支持文件排除模式配置
- **预览模式**: 支持 dry-run 预览同步操作
- **自定义选项**: 支持额外的 rsync 命令参数

### 服务管理 🆕
- **服务注册**: 记录服务器上部署的各种服务
- **服务详情**: 包含服务名称、类型、端口、路径等完整信息
- **关联管理**: 服务与 SSH 连接关联，便于管理
- **Markdown 输出**: 服务详情以美观的 Markdown 格式展示
- **搜索功能**: 快速搜索和定位服务

### 数据持久化
- 使用 SQLite 数据库存储所有配置信息
- 支持数据备份和恢复
- 自动数据库迁移

## 数据模型

### SSH 连接
每个 SSH 连接包含以下字段：
- `id`: 唯一标识符
- `name`: 连接名称
- `address`: 服务器地址
- `port`: 端口号（默认 22）
- `username`: 用户名
- `password_type`: 认证类型（password 或 keypath）
- `password`: 密码（当 password_type 为 password 时使用）
- `key_path`: 私钥文件路径（当 password_type 为 keypath 时使用）
- `description`: 连接描述
- `usage_count`: 使用次数

### Rsync 配置
每个 Rsync 配置包含以下字段：
- `id`: 唯一标识符
- `name`: 配置名称
- `ssh_name`: 关联的 SSH 连接名称
- `direction`: 传输方向（upload 或 download）
- `local_path`: 本地路径
- `remote_path`: 远程路径
- `exclude_rules`: 排除规则（换行分隔）
- `options`: 额外的 rsync 选项
- `description`: 配置描述
- `usage_count`: 使用次数

### 服务配置 🆕
每个服务配置包含以下字段：
- `id`: 唯一标识符
- `name`: 服务名称
- `server_name`: 服务器名称
- `description`: 服务简介
- `details`: 服务详细信息
- `status`: 服务状态（running/stopped/unknown）
- `port`: 服务端口号
- `service_type`: 服务类型（web/database/api等）
- `service_path`: 服务部署路径
- `config_path`: 配置文件路径
- `log_path`: 日志文件路径
- `ssh_connection_id`: 关联的 SSH 连接 ID

## 使用方法

### 编译项目

```bash
go mod tidy
go build -o alfred-tool
```

### 命令行使用

#### SSH 连接管理
```bash
# 添加新的 SSH 连接（打开 GUI 表单）
./alfred-tool ssh add

# 搜索连接（根据名称或地址）
./alfred-tool ssh search "myserver"
./alfred-tool ssh search "192.168.1.100"

# 显示所有连接
./alfred-tool ssh list

# 修改 SSH 连接（打开 GUI 表单）
./alfred-tool ssh update "myserver"

# 删除 SSH 连接
./alfred-tool ssh delete "myserver"

# 使用 SSH 连接（增加使用次数）
./alfred-tool ssh use "myserver"

# 同步配置到 ~/.ssh/config 文件
./alfred-tool ssh sync
```

#### Rsync 配置管理
```bash
# 添加新的 rsync 配置（打开 GUI 表单）
./alfred-tool rsync add

# 列出所有 rsync 配置
./alfred-tool rsync list

# 搜索 rsync 配置
./alfred-tool rsync search "backup"

# 修改 rsync 配置（打开 GUI 表单）
./alfred-tool rsync update "my-backup"

# 删除 rsync 配置
./alfred-tool rsync delete "my-backup"

# 执行 rsync 同步
./alfred-tool rsync run "my-backup"

# 预览 rsync 命令（不执行）
./alfred-tool rsync run "my-backup" --dry-run
```

#### 服务管理 🆕
```bash
# 添加新服务（打开 GUI 表单）
./alfred-tool service add

# 列出所有服务
./alfred-tool service list

# 搜索服务
./alfred-tool service search "nginx"

# 查看服务详情（Markdown 格式输出）
./alfred-tool service view 1

# 更新服务信息（打开 GUI 表单）
./alfred-tool service update 1

# 删除服务
./alfred-tool service delete 1
```

## 项目结构

```
alfred-tool/
├── main.go                    # 主程序入口
├── models/                   
│   ├── ssh_connection.go      # SSH 连接数据模型
│   ├── rsync_config.go        # Rsync 配置数据模型
│   └── service.go             # 服务数据模型
├── database/                 
│   └── database.go            # 数据库初始化和连接
├── services/                 
│   ├── ssh_service.go         # SSH 连接服务层
│   ├── rsync_service.go       # Rsync 配置服务层
│   └── service_service.go     # 服务管理服务层
├── ui/                       
│   ├── view_dialog.go         # SSH 连接管理对话框
│   ├── rsync_dialog.go        # Rsync 配置管理对话框
│   └── service_dialog.go      # 服务管理对话框
├── cmd/                      
│   ├── root.go                # 根命令
│   ├── ssh/                   # SSH 命令分组
│   │   ├── ssh.go             # SSH 主命令
│   │   ├── add.go             # SSH 连接添加命令
│   │   ├── list.go            # SSH 连接列表命令
│   │   ├── search.go          # SSH 连接搜索命令
│   │   ├── update.go          # SSH 连接更新命令
│   │   ├── delete.go          # SSH 连接删除命令
│   │   ├── use.go             # SSH 连接使用命令
│   │   └── sync.go            # SSH 配置同步命令
│   ├── rsync/                 # Rsync 命令分组
│   │   ├── rsync.go           # Rsync 主命令
│   │   ├── rsync_add.go       # Rsync 添加命令
│   │   ├── rsync_list.go      # Rsync 列表命令
│   │   ├── rsync_search.go    # Rsync 搜索命令
│   │   ├── rsync_update.go    # Rsync 更新命令
│   │   ├── rsync_delete.go    # Rsync 删除命令
│   │   └── rsync_run.go       # Rsync 执行命令
│   └── service/               # 服务管理命令分组
│       ├── service.go         # 服务管理主命令
│       ├── service_add.go     # 服务添加命令
│       ├── service_list.go    # 服务列表命令
│       ├── service_search.go  # 服务搜索命令
│       ├── service_view.go    # 服务详情命令
│       ├── service_update.go  # 服务更新命令
│       └── service_delete.go  # 服务删除命令
└── go.mod                     # Go 模块依赖
```

## 依赖项

- **Fyne v2**: GUI 框架
- **GORM**: ORM 框架
- **SQLite**: 数据库驱动
- **Cobra**: 命令行框架

## 数据存储

数据库文件自动创建在用户主目录下的 `.alfred-tool/connections.db` 路径，存储 SSH 连接、Rsync 配置和服务信息。

## Rsync 功能详情

### 支持的传输方向
- **上传**: 从本地同步文件到远程服务器
- **下载**: 从远程服务器同步文件到本地

### 支持的功能
- **排除规则**: 支持多个排除模式，每行一个
- **自定义选项**: 支持额外的 rsync 命令选项
- **预览模式**: 使用 `--dry-run` 预览同步命令
- **使用统计**: 自动记录配置使用次数

### 生成的 rsync 命令示例
```bash
# 上传示例
rsync -avz --progress --exclude "*.log" --exclude "*.tmp" -e "ssh -p 22 -i ~/.ssh/id_rsa" /local/path/ user@server:/remote/path/

# 下载示例
rsync -avz --progress --exclude "*.log" --exclude "*.tmp" -e "ssh -p 22 -i ~/.ssh/id_rsa" user@server:/remote/path/ /local/path/
```

## 界面特性

### GUI 界面
- **SSH 连接管理**: 优雅的表单设计，支持密码和私钥认证
- **Rsync 配置**: 完整的同步配置界面，实时命令预览
- **服务管理**: 直观的服务信息录入，支持SSH连接关联 🆕
- **智能表单**: 根据选择自动切换输入选项
- **文件选择器**: 支持私钥文件和目录选择
- **输入验证**: 实时验证和错误提示
- **响应式布局**: 适配不同屏幕尺寸

### 命令行界面
- 交互式输入表单
- 表格化数据展示
- Markdown 格式输出（服务详情）
- 彩色输出和进度提示
- 智能搜索和过滤
- 分组式命令结构，不同服务使用不同的子命令

## 服务管理功能详情 🆕

### 服务类型支持
- **Web 服务**: Nginx, Apache, Node.js 应用等
- **数据库服务**: MySQL, PostgreSQL, Redis, MongoDB 等
- **API 服务**: REST API, GraphQL, 微服务等
- **系统服务**: 系统守护进程、定时任务等
- **自定义服务**: 其他任意类型的服务

### 服务状态管理
- **运行中**: 服务正常运行
- **已停止**: 服务已停止
- **未知**: 服务状态未确定

### Markdown 输出示例
查看服务详情时将输出格式化的 Markdown 内容：

```markdown
# 服务详情

## 基本信息

| 字段 | 值 |
|------|---|
| ID | 1 |
| 服务名称 | **nginx** |
| 服务器名称 | **web-server-01** |
| 服务类型 | `web` |
| 端口号 | `80,443` |
| 服务状态 | `running` |

## 服务描述

**简介:** Nginx web 服务器

**详情:**

配置文件: /etc/nginx/nginx.conf
日志文件: /var/log/nginx/

## 关联SSH连接

| 字段 | 值 |
|------|-----|
| 连接名称 | **web-server** |
| 服务器地址 | `192.168.1.100:22` |
| 用户名 | `admin` |
```

## 最佳实践

### SSH 连接管理
1. 为每个服务器创建描述性的连接名称
2. 优先使用 SSH 密钥认证而非密码
3. 定期检查和更新连接信息
4. 利用使用统计了解连接频率

### Rsync 同步
1. 合理配置排除规则，避免同步不必要的文件
2. 使用 `--dry-run` 模式预览同步操作
3. 为重要数据创建双向备份配置
4. 定期测试同步配置的有效性

### 服务管理
1. 为每个服务提供详细的描述和文档
2. 记录完整的配置文件和日志路径
3. 定期更新服务状态信息
4. 利用 SSH 连接关联简化服务器管理