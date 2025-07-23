package cmd

import (
	"fmt"
	"alfred-tool/services"

	"github.com/spf13/cobra"
)

var rsyncDeleteCmd = &cobra.Command{
	Use:   "delete [配置名称]",
	Short: "删除rsync配置",
	Long:  `删除指定的rsync配置`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		configName := args[0]
		
		// 检查配置是否存在
		_, err := services.GetRsyncConfigByName(configName)
		if err != nil {
			fmt.Printf("错误: 配置 '%s' 不存在\n", configName)
			return
		}
		
		err = services.DeleteRsyncConfig(configName)
		if err != nil {
			fmt.Printf("删除失败: %v\n", err)
			return
		}
		
		fmt.Printf("rsync配置 '%s' 已删除\n", configName)
	},
}

func init() {
	rsyncCmd.AddCommand(rsyncDeleteCmd)
}