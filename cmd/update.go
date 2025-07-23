package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"alfred-tool/ui"
)

var updateCmd = &cobra.Command{
	Use:   "update <name>",
	Short: "修改SSH连接",
	Long:  `通过图形界面修改指定名称的SSH连接配置。`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		connectionName := args[0]
		if err := ui.ShowUpdateDialog(connectionName); err != nil {
			fmt.Printf("修改连接失败: %v\n", err)
			return
		}
	},
}