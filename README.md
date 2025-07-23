# Alfred Tool - SSH 连接管理工具

一个基于 Fyne GUI 的 SSH 连接管理控制台应用程序，使用 GORM 和 SQLite 进行数据存储。

## 功能特性

- **添加连接**: 通过优雅的 Fyne 表单界面添加 SSH 连接
- **搜索连接**: 根据连接名称或地址搜索已保存的连接
- **列表显示**: 显示所有已保存的 SSH 连接的简洁列表
- **多种认证**: 支持密码和私钥文件两种认证方式
- **Rsync 同步**: 配置和执行 rsync 文件同步任务
- **数据持久化**: 使用 SQLite 数据库存储连接信息和 rsync 配置

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
./alfred-tool add

# 搜索连接（根据名称或地址）
./alfred-tool search "myserver"
./alfred-tool search "192.168.1.100"

# 显示所有连接
./alfred-tool list

# 修改 SSH 连接（打开 GUI 表单）
./alfred-tool update "myserver"

# 删除 SSH 连接
./alfred-tool delete "myserver"
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

## 项目结构

```
alfred-tool/
├── main.go                    # 主程序入口
├── models/                   
│   ├── ssh_connection.go      # SSH 连接数据模型
│   └── rsync_config.go        # Rsync 配置数据模型
├── database/                 
│   └── database.go            # 数据库初始化和连接
├── services/                 
│   ├── ssh_service.go         # SSH 连接服务层
│   └── rsync_service.go       # Rsync 配置服务层
├── ui/                       
│   ├── view_dialog.go         # SSH 连接管理对话框
│   └── rsync_dialog.go        # Rsync 配置管理对话框
├── cmd/                      
│   ├── root.go                # 根命令
│   ├── add.go                 # SSH 连接添加命令
│   ├── list.go                # SSH 连接列表命令
│   ├── search.go              # SSH 连接搜索命令
│   ├── update.go              # SSH 连接更新命令
│   ├── delete.go              # SSH 连接删除命令
│   ├── rsync.go               # Rsync 主命令
│   ├── rsync_add.go           # Rsync 添加命令
│   ├── rsync_list.go          # Rsync 列表命令
│   ├── rsync_search.go        # Rsync 搜索命令
│   ├── rsync_update.go        # Rsync 更新命令
│   ├── rsync_delete.go        # Rsync 删除命令
│   └── rsync_run.go           # Rsync 执行命令
└── go.mod                     # Go 模块依赖
```

## 依赖项

- **Fyne v2**: GUI 框架
- **GORM**: ORM 框架
- **SQLite**: 数据库驱动
- **Cobra**: 命令行框架

## 数据存储

数据库文件自动创建在用户主目录下的 `.alfred-tool/connections.db` 路径，存储 SSH 连接和 Rsync 配置信息。

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

- 优雅的表单设计
- 根据密码类型智能切换输入框
- 文件选择器支持私钥文件选择
- 输入验证和错误提示
- 响应式布局设计