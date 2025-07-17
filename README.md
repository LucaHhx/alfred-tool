# SSH 连接管理工具

一个基于 Fyne GUI 的 SSH 连接管理控制台应用程序，使用 GORM 和 SQLite 进行数据存储。

## 功能特性

- **添加连接**: 通过优雅的 Fyne 表单界面添加 SSH 连接
- **搜索连接**: 根据连接名称或地址搜索已保存的连接
- **列表显示**: 显示所有已保存的 SSH 连接的简洁列表
- **多种认证**: 支持密码和私钥文件两种认证方式
- **数据持久化**: 使用 SQLite 数据库存储连接信息

## 数据模型

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

## 使用方法

### 编译项目

```bash
go mod tidy
go build -o sshd
```

### 命令行使用

```bash
# 添加新的 SSH 连接（打开 GUI 表单）
./sshd add

# 搜索连接（根据名称或地址）
./sshd search "myserver"
./sshd search "192.168.1.100"

# 显示所有连接
./sshd list
```

## 项目结构

```
sshd/
├── main.go              # 主程序入口
├── models/             
│   └── ssh_connection.go # SSH 连接数据模型
├── database/           
│   └── database.go      # 数据库初始化和连接
├── services/           
│   └── ssh_service.go   # SSH 连接服务层
├── ui/                 
│   └── add_dialog.go    # Fyne GUI 添加连接对话框
└── go.mod              # Go 模块依赖
```

## 依赖项

- **Fyne v2**: GUI 框架
- **GORM**: ORM 框架
- **SQLite**: 数据库驱动

## 数据存储

数据库文件自动创建在用户主目录下的 `.sshd/connections.db` 路径。

## 界面特性

- 优雅的表单设计
- 根据密码类型智能切换输入框
- 文件选择器支持私钥文件选择
- 输入验证和错误提示
- 响应式布局设计