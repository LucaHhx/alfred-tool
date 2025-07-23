package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"alfred-tool/services"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "删除SSH连接",
	Long:  `删除指定名称的SSH连接配置。`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		connectionName := args[0]
		
		// 确认删除
		fmt.Printf("确定要删除连接 '%s' 吗？[y/N]: ", connectionName)
		var response string
		fmt.Scanln(&response)
		
		if response != "y" && response != "Y" && response != "yes" && response != "YES" {
			fmt.Println("操作已取消")
			return
		}
		
		if err := services.DeleteConnection(connectionName); err != nil {
			fmt.Printf("删除连接失败: %v\n", err)
			return
		}
		
		fmt.Printf("连接 '%s' 已成功删除\n", connectionName)
	},
}