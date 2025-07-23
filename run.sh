#!/bin/bash

# 设置UTF-8环境
export LANG=zh_CN.UTF-8
export LC_ALL=zh_CN.UTF-8

# 在macOS上设置中文字体
if [[ "$OSTYPE" == "darwin"* ]]; then
    # 尝试使用系统中文字体
    if [ -f "/System/Library/Fonts/PingFang.ttc" ]; then
        export FYNE_FONT="/System/Library/Fonts/PingFang.ttc"
    elif [ -f "/System/Library/Fonts/Helvetica.ttc" ]; then
        export FYNE_FONT="/System/Library/Fonts/Helvetica.ttc"
    fi
fi

# 编译并运行
go build -o alfred-tool .
./alfred-tool "$@"