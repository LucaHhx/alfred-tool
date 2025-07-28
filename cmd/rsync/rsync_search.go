package rsync

import (
	"alfred-tool/services"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search [搜索词]",
	Short: "搜索rsync配置",
	Long:  `根据名称、SSH连接、路径或描述搜索rsync配置`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := args[0]

		configs, err := services.SearchRsyncConfigs(query)
		if err != nil {
			fmt.Printf("错误: %v\n", err)
			return
		}

		if len(configs) == 0 {
			fmt.Printf("没有找到匹配 '%s' 的rsync配置\n", query)
			return
		}

		fmt.Printf("找到 %d 个匹配 '%s' 的rsync配置:\n\n", len(configs), query)
		for _, config := range configs {
			fmt.Printf("名称: %s\n", config.Name)
			fmt.Printf("SSH连接: %s\n", config.SSHName)

			direction := "上传 (本地→服务器)"
			if config.Direction == "download" {
				direction = "下载 (服务器→本地)"
			}
			fmt.Printf("方向: %s\n", direction)

			fmt.Printf("本地路径: %s\n", config.LocalPath)
			fmt.Printf("远程路径: %s\n", config.RemotePath)

			if config.ExcludeRules != "" {
				fmt.Printf("排除规则: %s\n", strings.ReplaceAll(config.ExcludeRules, "\n", ", "))
			}

			if config.Options != "" {
				fmt.Printf("选项: %s\n", config.Options)
			}

			if config.Description != "" {
				fmt.Printf("描述: %s\n", config.Description)
			}

			fmt.Printf("使用次数: %d\n", config.UsageCount)
			fmt.Println("---")
		}
	},
}
