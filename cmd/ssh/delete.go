package ssh

import (
	"fmt"

	"alfred-tool/services"

	"github.com/spf13/cobra"
)

var DeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "删除SSH连接",
	Long:  `删除指定名称的SSH连接配置。`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		connectionName := args[0]
		if err := services.DeleteConnection(connectionName); err != nil {
			fmt.Printf("删除连接失败: %v\n", err)
			return
		}

		fmt.Printf("连接 '%s' 已成功删除\n", connectionName)
	},
}
