package rsync

import (
	"alfred-tool/services"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	dryRun bool
)

var runCmd = &cobra.Command{
	Use:   "run [配置名称]",
	Short: "执行rsync配置",
	Long:  `执行指定的rsync配置进行文件同步`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		configName := args[0]
		if dryRun {
			// 预览模式，只显示命令不执行
			cmdStr, err := services.DryRunRsyncConfig(configName)
			if err != nil {
				fmt.Printf("错误: %v\n", err)
				return
			}
			fmt.Printf("预览命令: %s\n", cmdStr)
			return
		}

		// 实际执行
		fmt.Printf("开始执行rsync配置: %s\n", configName)
		err := services.ExecuteRsyncConfig(configName)
		if err != nil {
			fmt.Printf("执行失败: %v\n", err)
			return
		}
		fmt.Printf("rsync配置 '%s' 执行完成\n", configName)
	},
}
