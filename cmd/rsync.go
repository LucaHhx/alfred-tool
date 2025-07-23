package cmd

import (
	"github.com/spf13/cobra"
)

var rsyncCmd = &cobra.Command{
	Use:   "rsync",
	Short: "Rsync配置管理",
	Long: `管理rsync配置，支持本地和远程服务器之间的文件同步

rsync配置包含以下功能：
- 添加、修改、删除rsync配置
- 列出和搜索已保存的配置
- 执行rsync同步操作
- 支持上传和下载两种方向
- 支持排除规则和自定义选项`,
}

func init() {
	rootCmd.AddCommand(rsyncCmd)
}