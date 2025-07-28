package ssh

import (
	"fmt"

	"alfred-tool/ui"
	"github.com/spf13/cobra"
)

var UpdateCmd = &cobra.Command{
	Use:   "update <name>",
	Short: "更新SSH连接",
	Long:  `更新指定名称的SSH连接配置。`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		connectionName := args[0]
		if err := ui.ShowUpdateDialog(connectionName); err != nil {
			fmt.Printf("修改连接失败: %v\n", err)
			return
		}
	},
}
