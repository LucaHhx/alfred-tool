# Swift Dialog 编译和运行说明

## 概述

本目录包含一个 Swift 脚本 `dialog.swift`，用于创建动态对话框界面。该脚本使用 SwiftUI 和 AppKit 框架实现。

## 编译方式

### 方式一：直接运行（推荐）

Swift 脚本可以直接执行，无需预编译：

```bash
chmod +x dialog.swift
./dialog.swift config.json
```

或者：

```bash
swift dialog.swift config.json
```

### 方式二：编译为可执行文件

将 Swift 脚本编译为独立的可执行二进制文件：

```bash
swiftc dialog.swift -o dialog
```

编译后运行：

```bash
./dialog config.json
```

### 编译优化选项

- **优化编译**（生产环境推荐）：
  ```bash
  swiftc -O dialog.swift -o dialog
  ```

- **包含调试信息**：
  ```bash
  swiftc -g dialog.swift -o dialog
  ```

## 运行规则

### 命令行参数

脚本接受一个参数：配置文件路径或 JSON 字符串

```bash
./dialog <config_path_or_json>
```

- 如果未提供参数，默认使用 `config.json`
- 可以传递文件路径：`./dialog myconfig.json`
- 可以直接传递 JSON 字符串：`./dialog '{"windowTitle":"Test","fields":[]}'`

### 配置文件格式

配置文件为 JSON 格式，示例：

```json
{
  "windowTitle": "用户输入对话框",
  "windowWidth": 600,
  "windowHeight": 300,
  "okLabel": "确定",
  "cancelLabel": "取消",
  "alwaysOnTop": false,
  "fields": [
    {
      "type": "text",
      "label": "用户名",
      "bindingKey": "username",
      "defaultValue": "",
      "copy": true
    },
    {
      "type": "checkbox",
      "label": "记住我",
      "bindingKey": "remember",
      "defaultValue": "false"
    },
    {
      "type": "dropdown",
      "label": "选择角色",
      "bindingKey": "role",
      "options": ["管理员", "用户", "访客"],
      "defaultValue": "用户"
    },
    {
      "type": "filepicker",
      "label": "选择文件",
      "bindingKey": "filepath",
      "filePickerType": "file",
      "defaultValue": ""
    },
    {
      "type": "filepicker",
      "label": "选择文件夹",
      "bindingKey": "folderpath",
      "filePickerType": "folder",
      "defaultValue": ""
    }
  ]
}
```

### 支持的字段类型

- `text`: 单行文本输入框
- `checkbox`: 复选框
- `dropdown`: 下拉菜单
- `texteditor`: 可调整大小的多行文本编辑器
- `segmented`: 分段单选按钮组
- `filepicker`: 文件/文件夹选择器（支持浏览按钮和手动输入路径）
  - 通过 `filePickerType` 参数指定类型：`"file"` 选择文件，`"folder"` 选择文件夹

### 字段参数说明

- `copy` (可选): 设置为 `true` 时，在 `text` 或 `texteditor` 字段后显示"复制"按钮，可将内容复制到剪贴板
- `alwaysOnTop` (可选): 设置为 `true` 时，窗口将始终置顶显示在其他窗口之上

### 输出格式

当用户点击"确定"按钮时，程序会将所有字段值以 JSON 格式输出到标准输出：

```json
{"username":"john","remember":true,"role":"管理员"}
```

点击"取消"按钮时，程序直接退出，不输出任何内容。

## 系统要求

- macOS 10.15 (Catalina) 或更高版本
- Swift 5.0 或更高版本
- AppKit 和 SwiftUI 框架支持

## 许可证

MIT License - Copyright © 2025 andy4222

使用或改编此代码时必须注明原作者。