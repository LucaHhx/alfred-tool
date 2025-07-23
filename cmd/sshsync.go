package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"alfred-tool/services"
)

var sshsyncCmd = &cobra.Command{
	Use:   "sshsync",
	Short: "同步数据库配置到SSH配置文件",
	Long:  `将数据库中的SSH连接配置同步到 ~/.ssh/config 文件中。`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := syncToSSHConfig(); err != nil {
			fmt.Printf("同步失败: %v\n", err)
			return
		}
		fmt.Println("SSH配置同步完成")
	},
}

func syncToSSHConfig() error {
	// 获取所有连接
	connections, err := services.ListAllConnections()
	if err != nil {
		return fmt.Errorf("获取连接列表失败: %w", err)
	}

	// 获取SSH配置文件路径
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("无法获取用户主目录: %w", err)
	}

	sshConfigPath := filepath.Join(homeDir, ".ssh", "config")
	
	// 确保.ssh目录存在
	sshDir := filepath.Dir(sshConfigPath)
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		return fmt.Errorf("无法创建SSH目录: %w", err)
	}

	// 生成新配置
	var configBuilder strings.Builder
	
	// 清空现有配置，只保留sshd管理的部分
	configBuilder.WriteString("# === SSHD MANAGED CONFIG START ===\n")
	
	// 为每个连接生成SSH配置
	for _, conn := range connections {
		configBuilder.WriteString(fmt.Sprintf("\nHost %s\n", conn.Name))
		configBuilder.WriteString(fmt.Sprintf("    HostName %s\n", conn.Address))
		configBuilder.WriteString(fmt.Sprintf("    Port %d\n", conn.Port))
		configBuilder.WriteString(fmt.Sprintf("    User %s\n", conn.Username))
		
		// 根据认证类型设置密钥路径
		if conn.PasswordType == "keypath" && conn.KeyPath != "" {
			configBuilder.WriteString(fmt.Sprintf("    IdentityFile %s\n", conn.KeyPath))
		}
		
		// 添加常用设置
		configBuilder.WriteString("    StrictHostKeyChecking no\n")
		configBuilder.WriteString("    UserKnownHostsFile /dev/null\n")
	}
	
	configBuilder.WriteString("\n# === SSHD MANAGED CONFIG END ===\n")

	// 写入配置文件
	if err := os.WriteFile(sshConfigPath, []byte(configBuilder.String()), 0600); err != nil {
		return fmt.Errorf("写入SSH配置文件失败: %w", err)
	}

	return nil
}

