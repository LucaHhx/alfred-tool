package component

import (
	"bytes"
	"os/exec"
	"strings"
)

// ShowMacNativeFileDialog 显示原生macOS文件对话框以选择文件
func ShowMacNativeFileDialog() string {
	script := `POSIX path of (choose file)`
	cmd := exec.Command("osascript", "-e", script)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return ""
	}
	// mac 返回路径带冒号，需要转换为标准路径
	macPath := strings.TrimSpace(out.String())
	return macPath
}
